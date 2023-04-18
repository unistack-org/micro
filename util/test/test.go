package test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/codec"
	"golang.org/x/sync/errgroup"
)

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

var ErrUnknownContentType = errors.New("unknown content type")

type Extension struct {
	Ext []string
}

var (
	ExtToTypes = map[string][]string{
		"json":  {"application/json", "application/grpc+json"},
		"yaml":  {"application/yaml", "application/yml", "text/yaml", "text/yml"},
		"yml":   {"application/yaml", "application/yml", "text/yaml", "text/yml"},
		"proto": {"application/grpc", "application/grpc+proto", "application/proto"},
	}

	DefaultExts = []string{"csv", "json", "yaml", "yml", "proto"}
)

func clientCall(ctx context.Context, c client.Client, req client.Request, rsp interface{}) error {
	return nil
}

func NewResponseFromFile(rspfile string) (*codec.Frame, error) {
	rspbuf, err := os.ReadFile(rspfile)
	if err != nil {
		return nil, err
	}
	return &codec.Frame{Data: rspbuf}, nil
}

func NewRequestFromFile(c client.Client, reqfile string) (client.Request, error) {
	reqbuf, err := os.ReadFile(reqfile)
	if err != nil {
		return nil, err
	}

	endpoint := path.Base(path.Dir(reqfile))
	ext := getExt(reqfile)

	var ct string
	if cts, ok := ExtToTypes[ext]; ok {
		for _, t := range cts {
			if _, ok = c.Options().Codecs[t]; ok {
				ct = t
				break
			}
		}
	}

	if ct == "" {
		return nil, ErrUnknownContentType
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
	br := bufio.NewReader(r)

	for {
		s, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		} else if err == io.EOF && len(s) == 0 {
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
			if rows == nil {
				rows = m.NewRows(records[0])
			} else {
				for idx := 0; idx < len(records); idx++ {
					rows.FromCSVString(strings.Join(records[idx], ","))
				}
			}
			continue
		}

		if rows != nil {
			exp.WillReturnRows(rows)
			rows = nil
		}

		switch {
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

func RunWithClientExpectResults(ctx context.Context, c client.Client, m sqlmock.Sqlmock, dir string, exts []string) error {
	tcases, err := getFiles(dir, exts)
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

		for idx := 0; idx < len(tcase.reqfiles); idx++ {
			g.TryGo(func() error {
				req, err := NewRequestFromFile(c, tcase.reqfiles[idx])
				if err != nil {
					return err
				}
				rsp, err := NewResponseFromFile(tcase.rspfiles[idx])
				if err != nil {
					return err
				}
				data := &codec.Frame{}
				err = c.Call(gctx, req, data, client.WithContentType(req.ContentType()))
				if err != nil {
					return err
				}
				if !bytes.Equal(rsp.Data, data.Data) {
					return fmt.Errorf("rsp not equal test %s != %s", rsp.Data, data.Data)
				}
				return nil
			})
		}
	}
	return g.Wait()
}

func RunWithClientExpectErrors(ctx context.Context, c client.Client, dir string) error {
	g, gctx := errgroup.WithContext(ctx)
	if !strings.Contains(dir, "parallel") {
		g.SetLimit(1)
	}
	_ = gctx
	g.TryGo(func() error {
		// rsp := &codec.Frame{}
		// return c.Call(ctx, req, rsp, client.WithContentType(req.ContentType()))
		return nil
	})
	return g.Wait()
}

type Case struct {
	dbfiles  []string
	reqfiles []string
	rspfiles []string
}

func getFiles(dir string, exts []string) ([]Case, error) {
	var tcases []Case
	entries, err := os.ReadDir(dir)
	if len(entries) == 0 && err != nil {
		return tcases, err
	}

	if exts == nil {
		exts = DefaultExts
	}

	var dirs []string
	var dbfiles, reqfiles, rspfiles []string

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
					reqfiles = append(reqfiles, filepath.Join(dir, entry.Name()))
				case strings.HasSuffix(name, "_rsp"):
					rspfiles = append(rspfiles, filepath.Join(dir, entry.Name()))
				}
			}
		}
	}

	if len(reqfiles) > 0 && len(rspfiles) > 0 {
		tcases = append(tcases, Case{dbfiles: dbfiles, reqfiles: reqfiles, rspfiles: rspfiles})
	}

	for _, dir = range dirs {
		ntcases, err := getFiles(dir, exts)
		if len(ntcases) == 0 && err != nil {
			return tcases, err
		} else if len(ntcases) == 0 {
			continue
		}
		tcases = append(tcases, ntcases...)
	}

	return tcases, nil
}
