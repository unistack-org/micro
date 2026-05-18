package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"go.unistack.org/micro/v5/client"
	"go.unistack.org/micro/v5/client/mock"
	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/errors"
	codecpb "go.unistack.org/micro-proto/v5/codec"
	"google.golang.org/grpc/status"
)

func Test_NewResponseFromFile(t *testing.T) {
	frame, err := NewResponseFromFile("testdata/result/01_firstcase/Call_rsp.json")
	if err != nil {
		t.Fatal(err)
	}
	if frame == nil {
		t.Fatal("expected non-nil frame")
	}
}

func Test_NewResponseFromFileMissing(t *testing.T) {
	_, err := NewResponseFromFile("testdata/nonexistent.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func Test_SQLFromBytes(t *testing.T) {
	ctx := context.TODO()
	db, c, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	data := []byte("# begin\n# query select_from_t\n# columns id|VARCHAR,name|VARCHAR\nid,val\n1,foo\n# commit\n")
	if err = SQLFromBytes(c, data); err != nil {
		t.Fatal(err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := tx.QueryContext(ctx, "select_from_t")
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var id, name string
		if err = rows.Scan(&id, &name); err != nil {
			t.Fatal(err)
		}
	}
	_ = rows.Close()
	_ = tx.Commit()
	if err = c.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func Test_SQLFromString(t *testing.T) {
	ctx := context.TODO()
	db, c, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	data := "# begin\n# query select_from_u\n# columns id|VARCHAR,name|VARCHAR\nid,val\n1,bar\n# commit\n"
	if err = SQLFromString(c, data); err != nil {
		t.Fatal(err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := tx.QueryContext(ctx, "select_from_u")
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var id, name string
		if err = rows.Scan(&id, &name); err != nil {
			t.Fatal(err)
		}
	}
	_ = rows.Close()
	_ = tx.Commit()
	if err = c.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func Test_getCodecFound(t *testing.T) {
	c := codec.NewCodec()
	Codecs = map[string]codec.Codec{
		"application/json": c,
	}
	result, err := getCodec(Codecs, "json")
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("expected non-nil codec")
	}
}

func Test_getCodecNotFound(t *testing.T) {
	Codecs = map[string]codec.Codec{}
	_, err := getCodec(Codecs, "json")
	if err == nil {
		t.Fatal("expected error for unknown extension")
	}
}

func Test_getCodecUnknownExt(t *testing.T) {
	Codecs = map[string]codec.Codec{}
	_, err := getCodec(Codecs, "xyz")
	if err == nil {
		t.Fatal("expected error for unknown extension")
	}
}

func Test_getContentTypeFound(t *testing.T) {
	c := codec.NewCodec()
	Codecs = map[string]codec.Codec{
		"application/json": c,
	}
	ct, err := getContentType(Codecs, "json")
	if err != nil {
		t.Fatal(err)
	}
	if ct == "" {
		t.Fatal("expected non-empty content type")
	}
}

func Test_getContentTypeNotFound(t *testing.T) {
	Codecs = map[string]codec.Codec{}
	_, err := getContentType(Codecs, "json")
	if err == nil {
		t.Fatal("expected error for unknown content type")
	}
}

func Test_SQLFromStringRollback(t *testing.T) {
	_, c, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	data := "# rollback\n"
	if err = SQLFromString(c, data); err != nil {
		t.Fatal(err)
	}
}

func Test_SQLFromStringExec(t *testing.T) {
	_, c, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	data := "# exec INSERT INTO t VALUES\n"
	if err = SQLFromString(c, data); err != nil {
		t.Fatal(err)
	}
}

func Test_SQLFromStringQueryWithRows(t *testing.T) {
	ctx := context.TODO()
	db, c, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	// Test the path where rows != nil at EOF with exp != nil
	data := "# query select_q\n# columns id|VARCHAR|NULL,flag|BOOL|BOOLEAN,amount|NUMBER|DECIMAL\nid,val,num\n1,true,1.5\n"
	if err = SQLFromString(c, data); err != nil {
		t.Fatal(err)
	}
	rows, err := db.QueryContext(ctx, "select_q")
	if err != nil {
		t.Fatal(err)
	}
	_ = rows.Close()
}

func Test_ResponseCompareFuncWithFrame(t *testing.T) {
	c := codec.NewCodec()
	frame := &codecpb.Frame{Data: []byte(`{}`)}
	err := ResponseCompareFunc([]byte(`{}`), frame, c, c)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ResponseCompareFuncWithError(t *testing.T) {
	c := codec.NewCodec()
	// errors.Error implements error; its Error() returns JSON-encoded form
	me := &errors.Error{ID: "test", Code: 200, Detail: "ok", Status: "OK"}
	err := ResponseCompareFunc([]byte(me.Error()), me, c, c)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ResponseCompareFuncMismatch(t *testing.T) {
	c := codec.NewCodec()
	frame := &codecpb.Frame{Data: []byte(`{"x":1}`)}
	err := ResponseCompareFunc([]byte(`{"y":2}`), frame, c, c)
	if err == nil {
		t.Fatal("expected mismatch error")
	}
}

func Test_ResponseCompareFuncGenericError(t *testing.T) {
	c := codec.NewCodec()
	// A plain error (not *errors.Error, not grpc status) should cause unmarshal failure
	plainErr := fmt.Errorf("plain error")
	err := ResponseCompareFunc([]byte(`{}`), plainErr, c, c)
	// plain errors.FromError wraps as gRPC status with no code — will return the error itself
	_ = err
}

func Test_CSVColumnParserNull(t *testing.T) {
	result := CSVColumnParser("null")
	if result != nil {
		t.Fatalf("expected nil for 'null', got %v", result)
	}
}

func Test_CSVColumnParserEmpty(t *testing.T) {
	result := CSVColumnParser("")
	if result != nil {
		t.Fatalf("expected nil for empty string, got %v", result)
	}
}

func Test_CSVColumnParserValue(t *testing.T) {
	result := CSVColumnParser("hello")
	if string(result) != "hello" {
		t.Fatalf("expected 'hello', got %q", result)
	}
}

func Test_SQLFromFile(t *testing.T) {
	ctx := context.TODO()
	db, c, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	if err = SQLFromFile(c, "testdata/result/01_firstcase/Call_db.csv"); err != nil {
		t.Fatal(err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	rows, err := tx.QueryContext(ctx, "select * from test;")
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var id int64
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			t.Fatal(err)
		}
		if id != 1 || name != "test" {
			t.Fatalf("invalid rows %v %v", id, name)
		}
	}

	if err = rows.Close(); err != nil {
		t.Fatal(err)
	}

	if err = rows.Err(); err != nil {
		t.Fatal(err)
	}

	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}
	if err = c.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func Test_GetCases(t *testing.T) {
	files, err := GetCases("testdata/", nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) == 0 {
		t.Fatalf("no files matching")
	}

	if n := len(files); n != 1 {
		t.Fatalf("invalid number of test cases %d", n)
	}
}

func Test_NewRequestFromFile(t *testing.T) {
	c := codec.NewCodec()
	Codecs = map[string]codec.Codec{"application/json": c}

	mc := mock.NewClient(client.Codec("application/json", c))
	req, err := NewRequestFromFile(mc, "testdata/result/01_firstcase/Call_req.json")
	if err != nil {
		t.Fatal(err)
	}
	if req == nil {
		t.Fatal("expected non-nil request")
	}
}

func Test_NewRequestFromFileMissing(t *testing.T) {
	c := codec.NewCodec()
	mc := mock.NewClient(client.Codec("application/json", c))
	_, err := NewRequestFromFile(mc, "testdata/nonexistent.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func Test_RunBadDir(t *testing.T) {
	_, sqlm, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	c := codec.NewCodec()
	mc := mock.NewClient(client.Codec("application/json", c))
	err = Run(context.Background(), mc, sqlm, "/nonexistent/dir/that/does/not/exist", nil)
	if err == nil {
		t.Fatal("expected error for bad dir")
	}
}

func Test_ResponseCompareFuncWithStatusStatus(t *testing.T) {
	c := codec.NewCodec()
	Codecs = map[string]codec.Codec{"application/json": c}

	me := &errors.Error{ID: "test", Code: 200, Detail: "ok", Status: "OK"}
	st := status.New(200, me.Error())
	err := ResponseCompareFunc([]byte(me.Error()), st, c, c)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ResponseCompareFuncGRPCStatusIface(t *testing.T) {
	c := codec.NewCodec()
	Codecs = map[string]codec.Codec{"application/json": c}

	me := &errors.Error{ID: "test", Code: 200, Detail: "ok", Status: "OK"}
	st := status.New(200, me.Error())
	// grpcErr implements interface{ GRPCStatus() *status.Status }
	grpcErr := st.Err()
	err := ResponseCompareFunc([]byte(me.Error()), grpcErr, c, c)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_FromCSVStringNullableColumn(t *testing.T) {
	col := sqlmock.NewColumn("c1").OfType("VARCHAR", nil).Nullable(true)
	rows := sqlmock.NewRows([]string{"c1"})
	// pass "null" to trigger the nullable nil path
	rows = FromCSVString([]*sqlmock.Column{col}, rows, "null")
	if rows == nil {
		t.Fatal("expected non-nil rows")
	}
}
