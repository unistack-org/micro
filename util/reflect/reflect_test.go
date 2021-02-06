package reflect

import (
	"net/url"
	"testing"
)

func TestURLSliceVars(t *testing.T) {
	u, err := url.Parse("http://localhost/v1/test/call/my_name?key=arg1&key=arg2&key=arg3")
	if err != nil {
		t.Fatal(err)
	}

	mp, err := URLMap(u.RawQuery)
	if err != nil {
		t.Fatal(err)
	}

	v, ok := mp["key"]
	if !ok {
		t.Fatalf("key not exists: %#+v", mp)
	}

	vm, ok := v.([]interface{})
	if !ok {
		t.Fatalf("invalid key value")
	}

	if len(vm) != 3 {
		t.Fatalf("missing key value: %#+v", mp)
	}
}

func TestURLVars(t *testing.T) {
	u, err := url.Parse("http://localhost/v1/test/call/my_name?req=key&arg1=arg1&arg2=12345&nested.string_args=str1&nested.string_args=str2&arg2=54321")
	if err != nil {
		t.Fatal(err)
	}

	mp, err := URLMap(u.RawQuery)
	if err != nil {
		t.Fatal(err)
	}
	_ = mp
}
