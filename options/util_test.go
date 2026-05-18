package options

import "testing"

// TestToInt8SliceE_nil covers the nil-input error branch.
func TestToInt8SliceE_nil(t *testing.T) {
	_, err := toInt8SliceE(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

// TestToInt8SliceE_direct covers the already-correct-type fast path.
func TestToInt8SliceE_direct(t *testing.T) {
	in := []int8{1, 2}
	out, err := toInt8SliceE(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("unexpected out %v", out)
	}
}

// TestToInt8SliceE_default covers the non-slice default error branch.
func TestToInt8SliceE_default(t *testing.T) {
	_, err := toInt8SliceE("notaslice")
	if err == nil {
		t.Fatal("expected error for non-slice input")
	}
}

// TestToInt16SliceE_nil covers the nil-input error branch.
func TestToInt16SliceE_nil(t *testing.T) {
	_, err := toInt16SliceE(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

// TestToInt16SliceE_default covers the non-slice default error branch.
func TestToInt16SliceE_default(t *testing.T) {
	_, err := toInt16SliceE("notaslice")
	if err == nil {
		t.Fatal("expected error for non-slice input")
	}
}

// TestToInt32SliceE_nil covers the nil-input error branch.
func TestToInt32SliceE_nil(t *testing.T) {
	_, err := toInt32SliceE(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

// TestToInt32SliceE_default covers the non-slice default error branch.
func TestToInt32SliceE_default(t *testing.T) {
	_, err := toInt32SliceE("notaslice")
	if err == nil {
		t.Fatal("expected error for non-slice input")
	}
}

// TestToInt64SliceE_nil covers the nil-input error branch.
func TestToInt64SliceE_nil(t *testing.T) {
	_, err := toInt64SliceE(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

// TestToInt64SliceE_default covers the non-slice default error branch.
func TestToInt64SliceE_default(t *testing.T) {
	_, err := toInt64SliceE("notaslice")
	if err == nil {
		t.Fatal("expected error for non-slice input")
	}
}

// TestToUintSliceE_nil covers the nil-input error branch.
func TestToUintSliceE_nil(t *testing.T) {
	_, err := toUintSliceE(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

// TestToUintSliceE_default covers the non-slice default error branch.
func TestToUintSliceE_default(t *testing.T) {
	_, err := toUintSliceE("notaslice")
	if err == nil {
		t.Fatal("expected error for non-slice input")
	}
}

// TestToUint8SliceE_nil covers the nil-input error branch.
func TestToUint8SliceE_nil(t *testing.T) {
	_, err := toUint8SliceE(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

// TestToUint8SliceE_default covers the non-slice default error branch.
func TestToUint8SliceE_default(t *testing.T) {
	_, err := toUint8SliceE("notaslice")
	if err == nil {
		t.Fatal("expected error for non-slice input")
	}
}

// TestToUint16SliceE_nil covers the nil-input error branch.
func TestToUint16SliceE_nil(t *testing.T) {
	_, err := toUint16SliceE(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

// TestToUint16SliceE_default covers the non-slice default error branch.
func TestToUint16SliceE_default(t *testing.T) {
	_, err := toUint16SliceE("notaslice")
	if err == nil {
		t.Fatal("expected error for non-slice input")
	}
}

// TestToUint32SliceE_nil covers the nil-input error branch.
func TestToUint32SliceE_nil(t *testing.T) {
	_, err := toUint32SliceE(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

// TestToUint32SliceE_default covers the non-slice default error branch.
func TestToUint32SliceE_default(t *testing.T) {
	_, err := toUint32SliceE("notaslice")
	if err == nil {
		t.Fatal("expected error for non-slice input")
	}
}

// TestToUint64SliceE_nil covers the nil-input error branch.
func TestToUint64SliceE_nil(t *testing.T) {
	_, err := toUint64SliceE(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

// TestToUint64SliceE_default covers the non-slice default error branch.
func TestToUint64SliceE_default(t *testing.T) {
	_, err := toUint64SliceE("notaslice")
	if err == nil {
		t.Fatal("expected error for non-slice input")
	}
}

// TestToFloat32SliceE_nil covers the nil-input error branch.
func TestToFloat32SliceE_nil(t *testing.T) {
	_, err := toFloat32SliceE(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

// TestToFloat32SliceE_default covers the non-slice default error branch.
func TestToFloat32SliceE_default(t *testing.T) {
	_, err := toFloat32SliceE("notaslice")
	if err == nil {
		t.Fatal("expected error for non-slice input")
	}
}

// TestToFloat64SliceE_nil covers the nil-input error branch.
func TestToFloat64SliceE_nil(t *testing.T) {
	_, err := toFloat64SliceE(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

// TestToFloat64SliceE_default covers the non-slice default error branch.
func TestToFloat64SliceE_default(t *testing.T) {
	_, err := toFloat64SliceE("notaslice")
	if err == nil {
		t.Fatal("expected error for non-slice input")
	}
}

// TestSetMap_nilSrc covers the nil-src error branch in setMap.
func TestSetMap_nilSrc(t *testing.T) {
	err := setMap(nil, "val")
	if err == nil {
		t.Fatal("expected error for nil src")
	}
}

// TestSetMap_nilDst covers the nil-dst error branch in setMap.
func TestSetMap_nilDst(t *testing.T) {
	err := setMap(map[string]string{}, nil)
	if err == nil {
		t.Fatal("expected error for nil dst")
	}
}

// TestSetMap_dstNonMap covers the []any{nonMap} -> return nil path in setMap.
func TestSetMap_dstNonMap(t *testing.T) {
	m := map[string]string{}
	err := setMap(m, []any{"notamap"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestSetMap_directNonMap covers the default-branch non-map dst -> return nil path.
func TestSetMap_directNonMap(t *testing.T) {
	m := map[string]string{}
	err := setMap(m, "notamap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
