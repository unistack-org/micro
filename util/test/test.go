package test

import (
	"encoding/csv"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/codec"
)

var ErrUnknownContentType = errors.New("unknown content type")

type Extension struct {
	Ext []string
}

var ExtToTypes = map[string][]string{
	"json":  {"application/json", "application/grpc+json"},
	"yaml":  {"application/yaml", "application/yml", "text/yaml", "text/yml"},
	"yml":   {"application/yaml", "application/yml", "text/yaml", "text/yml"},
	"proto": {"application/grpc", "application/grpc+proto", "application/proto"},
}

func NewRequestFromFile(c client.Client, reqfile string) (client.Request, error) {
	reqbuf, err := os.ReadFile(reqfile)
	if err != nil {
		return nil, err
	}

	endpoint := path.Base(path.Dir(reqfile))
	ext := path.Ext(reqfile)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}

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

func NewSQLRowsFromFile(c sqlmock.Sqlmock, file string) (*sqlmock.Rows, error) {
	fp, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	r := csv.NewReader(fp)
	r.Comma = '\t'
	r.Comment = '#'

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	rows := c.NewRows(records[0])
	for idx := 1; idx < len(records); idx++ {
		rows.FromCSVString(strings.Join(records[idx], ";"))
	}

	return rows, nil
}
