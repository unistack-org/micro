package flow

import (
	"reflect"
	"testing"
)

func FuzzMarshall(f *testing.F) {
	f.Fuzz(func(t *testing.T, ref []byte) {
		rm := RawMessage(ref)

		b, err := rm.MarshalJSON()
		if err != nil {
			t.Errorf("Error MarshalJSON: %s", err)
		}

		if !reflect.DeepEqual(ref, b) {
			t.Errorf("Error. Expected '%s', was '%s'", ref, b)
		}
	})
}

func FuzzUnmarshall(f *testing.F) {
	f.Fuzz(func(t *testing.T, ref string) {
		b := []byte(ref)
		rm := RawMessage(b)

		if err := rm.UnmarshalJSON(b); err != nil {
			t.Errorf("Error UnmarshalJSON: %s", err)
		}

		if ref != string(rm) {
			t.Errorf("Error. Expected '%s', was '%s'", ref, rm)
		}
	})
}
