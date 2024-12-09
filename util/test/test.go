package test

import (
	"bufio"
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/codec"
	"go.unistack.org/micro/v3/errors"
	"go.unistack.org/micro/v3/metadata"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var ErrUnknownContentType = fmt.Errorf("unknown content type")

type Extension struct {
	Ext []string
}

var (
	// ExtToTypes map file extension to content type
	ExtToTypes = map[string][]string{
		"json":  {"application/json", "application/grpc+json"},
		"yaml":  {"application/yaml", "application/yml", "text/yaml", "text/yml"},
		"yml":   {"application/yaml", "application/yml", "text/yaml", "text/yml"},
		"proto": {"application/grpc", "application/grpc+proto", "application/proto"},
	}
	// DefaultExts specifies default file extensions to load data
	DefaultExts = []string{"csv", "json", "yaml", "yml", "proto"}
	// Codecs map to detect codec for test file or request content type
	Codecs map[string]codec.Codec

	// ResponseCompareFunc used to compare actual response with test case data
	ResponseCompareFunc = func(expectRsp []byte, testRsp interface{}, expectCodec codec.Codec, testCodec codec.Codec) error {
		var err error

		expectMap := make(map[string]interface{})
		if err = expectCodec.Unmarshal(expectRsp, &expectMap); err != nil {
			return fmt.Errorf("failed to unmarshal err: %w", err)
		}

		testMap := make(map[string]interface{})
		switch v := testRsp.(type) {
		case *codec.Frame:
			if err = testCodec.Unmarshal(v.Data, &testMap); err != nil {
				return fmt.Errorf("failed to unmarshal err: %w", err)
			}
		case *errors.Error:
			if err = expectCodec.Unmarshal([]byte(v.Error()), &testMap); err != nil {
				return fmt.Errorf("failed to unmarshal err: %w", err)
			}
		case error:
			st, ok := status.FromError(v)
			if !ok {
				return v
			}
			me := errors.Parse(st.Message())
			if me.Code != 0 {
				if err = expectCodec.Unmarshal([]byte(me.Error()), &testMap); err != nil {
					return fmt.Errorf("failed to unmarshal err: %w", err)
				}
				break
			}
			for _, se := range st.Details() {
				switch ne := se.(type) {
				case proto.Message:
					var buf []byte
					if buf, err = testCodec.Marshal(ne); err != nil {
						return fmt.Errorf("failed to marshal err: %w", err)
					}
					if err = testCodec.Unmarshal(buf, &testMap); err != nil {
						return fmt.Errorf("failed to unmarshal err: %w", err)
					}
				default:
					return st.Err()
				}
			}
		case interface{ GRPCStatus() *status.Status }:
			st := v.GRPCStatus()
			me := errors.Parse(st.Message())
			if me.Code != 0 {
				if err = expectCodec.Unmarshal([]byte(me.Error()), &testMap); err != nil {
					return fmt.Errorf("failed to unmarshal err: %w", err)
				}
				break
			}
		case *status.Status:
			me := errors.Parse(v.Message())
			if me.Code != 0 {
				if err = expectCodec.Unmarshal([]byte(me.Error()), &testMap); err != nil {
					return fmt.Errorf("failed to unmarshal err: %w", err)
				}
				break
			}
			for _, se := range v.Details() {
				switch ne := se.(type) {
				case proto.Message:
					buf, err := testCodec.Marshal(ne)
					if err != nil {
						return fmt.Errorf("failed to marshal err: %w", err)
					}
					if err = testCodec.Unmarshal(buf, &testMap); err != nil {
						return fmt.Errorf("failed to unmarshal err: %w", err)
					}
				default:
					return v.Err()
				}
			}
		}

		if !reflect.DeepEqual(expectMap, testMap) {
			return fmt.Errorf("test: %s != rsp: %s", expectMap, testMap)
		}

		return nil
	}
)

func FromCSVString(columns []*sqlmock.Column, rows *sqlmock.Rows, s string) *sqlmock.Rows {
	res := strings.NewReader(strings.TrimSpace(s))
	csvReader := csv.NewReader(res)

	for {
		res, err := csvReader.Read()
		if err != nil || res == nil {
			break
		}

		var row []driver.Value
		for i, v := range res {
			item := CSVColumnParser(strings.TrimSpace(v))
			if null, nullOk := columns[i].IsNullable(); null && nullOk && item == nil {
				row = append(row, nil)
			} else {
				row = append(row, item)
			}

		}
		rows = rows.AddRow(row...)
	}

	return rows
}

func CSVColumnParser(s string) []byte {
	switch {
	case strings.ToLower(s) == "null":
		return nil
	case s == "":
		return nil
	}
	return []byte(s)
}

func NewResponseFromFile(rspfile string) (*codec.Frame, error) {
	rspbuf, err := os.ReadFile(rspfile)
	if err != nil {
		return nil, err
	}
	return &codec.Frame{Data: rspbuf}, nil
}

func getCodec(codecs map[string]codec.Codec, ext string) (codec.Codec, error) {
	var c codec.Codec
	if cts, ok := ExtToTypes[ext]; ok {
		for _, t := range cts {
			if c, ok = codecs[t]; ok {
				return c, nil
			}
		}
	}
	return nil, ErrUnknownContentType
}

func getContentType(codecs map[string]codec.Codec, ext string) (string, error) {
	if cts, ok := ExtToTypes[ext]; ok {
		for _, t := range cts {
			if _, ok = codecs[t]; ok {
				return t, nil
			}
		}
	}
	return "", ErrUnknownContentType
}

func getExt(name string) string {
	ext := filepath.Ext(name)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}
	return ext
}

func getNameWithoutExt(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}

func NewRequestFromFile(c client.Client, reqfile string) (client.Request, error) {
	reqbuf, err := os.ReadFile(reqfile)
	if err != nil {
		return nil, err
	}

	endpoint := path.Base(path.Dir(reqfile))
	if idx := strings.Index(endpoint, "_"); idx > 0 {
		endpoint = endpoint[idx+1:]
	}
	ext := getExt(reqfile)

	ct, err := getContentType(c.Options().Codecs, ext)
	if err != nil {
		return nil, err
	}

	req := c.NewRequest("test", endpoint, &codec.Frame{Data: reqbuf}, client.RequestContentType(ct))

	return req, nil
}

func SQLFromFile(m sqlmock.Sqlmock, name string) error {
	fp, err := os.Open(name)
	if err != nil {
		return err
	}
	defer fp.Close()
	return SQLFromReader(m, fp)
}

func SQLFromBytes(m sqlmock.Sqlmock, buf []byte) error {
	return SQLFromReader(m, bytes.NewReader(buf))
}

func SQLFromString(m sqlmock.Sqlmock, buf string) error {
	return SQLFromReader(m, strings.NewReader(buf))
}

func SQLFromReader(m sqlmock.Sqlmock, r io.Reader) error {
	var rows *sqlmock.Rows
	var exp *sqlmock.ExpectedQuery
	var columns []*sqlmock.Column

	br := bufio.NewReader(r)

	for {
		s, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		} else if err == io.EOF && len(s) == 0 {
			if rows != nil && exp != nil {
				exp.WillReturnRows(rows)
			}
			return nil
		}

		if s[0] != '#' {
			r := csv.NewReader(strings.NewReader(s))
			r.Comma = ','
			var records [][]string
			records, err = r.ReadAll()
			if err != nil {
				return err
			}
			if rows == nil && len(columns) > 0 {
				rows = m.NewRowsWithColumnDefinition(columns...)
			} else {
				for idx := 0; idx < len(records); idx++ {
					if len(columns) == 0 {
						return fmt.Errorf("csv file not valid, does not have %q line", "# columns ")
					}
					rows = FromCSVString(columns, rows, strings.Join(records[idx], ","))
				}
			}
			continue
		}

		if rows != nil {
			exp.WillReturnRows(rows)
			rows = nil
		}

		switch {
		case strings.HasPrefix(strings.ToLower(s[2:]), "columns"):
			for _, field := range strings.Split(s[2+len("columns")+1:], ",") {
				args := strings.Split(field, "|")

				column := sqlmock.NewColumn(args[0]).Nullable(false)

				if len(args) > 1 {
					for _, arg := range args {
						switch arg {
						case "BOOLEAN", "BOOL":
							column = column.OfType("BOOL", false)
						case "NUMBER", "DECIMAL":
							column = column.OfType("DECIMAL", float64(0.0)).WithPrecisionAndScale(10, 4)
						case "VARCHAR":
							column = column.OfType("VARCHAR", nil)
						case "NULL":
							column = column.Nullable(true)
						}
					}
				}

				columns = append(columns, column)
			}
		case strings.HasPrefix(strings.ToLower(s[2:]), "begin"):
			m.ExpectBegin()
		case strings.HasPrefix(strings.ToLower(s[2:]), "commit"):
			m.ExpectCommit()
		case strings.HasPrefix(strings.ToLower(s[2:]), "rollback"):
			m.ExpectRollback()
		case strings.HasPrefix(strings.ToLower(s[2:]), "exec "):
			m.ExpectExec(s[2+len("exec "):])
		case strings.HasPrefix(strings.ToLower(s[2:]), "query "):
			exp = m.ExpectQuery(s[2+len("query "):])
		}
	}
}

func Run(ctx context.Context, c client.Client, m sqlmock.Sqlmock, dir string, exts []string) error {
	tcases, err := GetCases(dir, exts)
	if err != nil {
		return err
	}

	g, gctx := errgroup.WithContext(ctx)
	if !strings.Contains(dir, "parallel") {
		g.SetLimit(1)
	}

	for _, tcase := range tcases {

		for _, dbfile := range tcase.dbfiles {
			if err = SQLFromFile(m, dbfile); err != nil {
				return err
			}
		}

		tc := tcase
		g.Go(func() error {
			var xrid string
			var gerr error

			treq, err := NewRequestFromFile(c, tc.reqfile)
			if err != nil {
				gerr = fmt.Errorf("failed to read request from file %s err: %w", tc.reqfile, err)
				return gerr
			}

			xrid = fmt.Sprintf("%s-%d", treq.Endpoint(), time.Now().Unix())

			defer func() {
				if gerr == nil {
					fmt.Printf("test %s xrid: %s status: success\n", filepath.Dir(tc.reqfile), xrid)
				} else {
					fmt.Printf("test %s xrid: %s status: failure error: %v\n", filepath.Dir(tc.reqfile), xrid, err)
				}
			}()

			data := &codec.Frame{}
			md := metadata.New(1)
			md.Set("X-Request-Id", xrid)
			cerr := c.Call(metadata.NewOutgoingContext(gctx, md), treq, data, client.WithContentType(treq.ContentType()))

			var rspfile string

			if tc.errfile != "" {
				rspfile = tc.errfile
			} else if tc.rspfile != "" {
				rspfile = tc.rspfile
			} else {
				gerr = fmt.Errorf("errfile and rspfile is empty")
				return gerr
			}

			expectRsp, err := NewResponseFromFile(rspfile)
			if err != nil {
				gerr = fmt.Errorf("failed to read response from file %s err: %w", rspfile, err)
				return gerr
			}

			testCodec, err := getCodec(Codecs, getExt(tc.reqfile))
			if err != nil {
				gerr = fmt.Errorf("failed to get response file codec err: %w", err)
				return gerr
			}

			expectCodec, err := getCodec(Codecs, getExt(rspfile))
			if err != nil {
				gerr = fmt.Errorf("failed to get response file codec err: %w", err)
				return gerr
			}

			if cerr == nil && tc.errfile != "" {
				gerr = fmt.Errorf("expected err %s not happened", expectRsp.Data)
				return gerr
			} else if cerr != nil && tc.errfile != "" {
				if err = ResponseCompareFunc(expectRsp.Data, cerr, expectCodec, testCodec); err != nil {
					gerr = err
					return gerr
				}
			} else if cerr != nil && tc.errfile == "" {
				gerr = cerr
				return gerr
			} else if cerr == nil && tc.errfile == "" {
				if err = ResponseCompareFunc(expectRsp.Data, data, expectCodec, testCodec); err != nil {
					gerr = err
					return gerr
				}
			}

			/*
				cf, err := getCodec(c.Options().Codecs, getExt(tc.rspfile))
				if err != nil {
					return err
				}
			*/

			return nil
		})

	}

	return g.Wait()
}

type Case struct {
	reqfile string
	rspfile string
	errfile string
	dbfiles []string
}

func GetCases(dir string, exts []string) ([]Case, error) {
	var tcases []Case
	entries, err := os.ReadDir(dir)
	if len(entries) == 0 && err != nil {
		return tcases, err
	}

	if exts == nil {
		exts = DefaultExts
	}

	var dirs []string
	var dbfiles []string
	var reqfile, rspfile, errfile string

	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, filepath.Join(dir, entry.Name()))
			continue
		}
		if info, err := entry.Info(); err != nil {
			return tcases, err
		} else if !info.Mode().IsRegular() {
			continue
		}

		for _, ext := range exts {
			if getExt(entry.Name()) == ext {
				name := getNameWithoutExt(entry.Name())
				switch {
				case strings.HasSuffix(name, "_db"):
					dbfiles = append(dbfiles, filepath.Join(dir, entry.Name()))
				case strings.HasSuffix(name, "_req"):
					reqfile = filepath.Join(dir, entry.Name())
				case strings.HasSuffix(name, "_rsp"):
					rspfile = filepath.Join(dir, entry.Name())
				case strings.HasSuffix(name, "_err"):
					errfile = filepath.Join(dir, entry.Name())
				}
			}
		}
	}

	if reqfile != "" && (rspfile != "" || errfile != "") {
		tcases = append(tcases, Case{dbfiles: dbfiles, reqfile: reqfile, rspfile: rspfile, errfile: errfile})
	}

	for _, dir = range dirs {
		ntcases, err := GetCases(dir, exts)
		if len(ntcases) == 0 && err != nil {
			return tcases, err
		} else if len(ntcases) == 0 {
			continue
		}
		tcases = append(tcases, ntcases...)
	}

	return tcases, nil
}
