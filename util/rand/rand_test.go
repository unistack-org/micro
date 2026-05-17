package rand

import (
	"testing"
)

func TestInt63(t *testing.T) {
	r := &Rand{}
	const maxInt63 int64 = 1<<62 + (1<<62 - 1) // 2^63 - 1
	for i := 0; i < 100; i++ {
		val := r.Int63()
		if val < 0 {
			t.Fatalf("Int63 returned negative value: %d", val)
		}
		if val > maxInt63 {
			t.Fatalf("Int63 returned value > max: %d", val)
		}
	}
}

func TestInt31(t *testing.T) {
	r := &Rand{}
	const maxInt31 int32 = 1<<30 + (1<<30 - 1) // 2^31 - 1
	for i := 0; i < 100; i++ {
		val := r.Int31()
		if val < 0 {
			t.Fatalf("Int31 returned negative value: %d", val)
		}
		if val > maxInt31 {
			t.Fatalf("Int31 returned value > max: %d", val)
		}
	}
}

func TestInt(t *testing.T) {
	r := &Rand{}
	for i := 0; i < 100; i++ {
		val := r.Int()
		if val < 0 {
			t.Fatalf("Int returned negative value: %d", val)
		}
	}
}

func TestUint32(t *testing.T) {
	r := &Rand{}
	for i := 0; i < 100; i++ {
		val := r.Uint32()
		if val < 0 {
			t.Fatalf("Uint32 returned negative value: %d", val)
		}
	}
}

func TestUint64(t *testing.T) {
	r := &Rand{}
	for i := 0; i < 100; i++ {
		val := r.Uint64()
		if val < 0 {
			t.Fatalf("Uint64 returned negative value: %d", val)
		}
	}
}

func TestFloat64(t *testing.T) {
	r := &Rand{}
	for i := 0; i < 100; i++ {
		val := r.Float64()
		if val < 0 || val >= 1 {
			t.Fatalf("Float64 returned value outside [0, 1): %v", val)
		}
	}
}

func TestFloat32(t *testing.T) {
	r := &Rand{}
	for i := 0; i < 100; i++ {
		val := r.Float32()
		if val < 0 || val >= 1 {
			t.Fatalf("Float32 returned value outside [0, 1): %v", val)
		}
	}
}

func TestIntnSmall(t *testing.T) {
	r := &Rand{}
	n := 10
	for i := 0; i < 100; i++ {
		val := r.Intn(n)
		if val < 0 || val >= n {
			t.Fatalf("Intn(%d) returned value outside [0, %d): %d", n, n, val)
		}
	}
}

func TestIntnLarge(t *testing.T) {
	r := &Rand{}
	n := 1<<31 + 1000
	for i := 0; i < 50; i++ {
		val := r.Intn(n)
		if val < 0 || val >= n {
			t.Fatalf("Intn(%d) returned value outside [0, %d): %d", n, n, val)
		}
	}
}

func TestInt31nRegular(t *testing.T) {
	r := &Rand{}
	n := int32(100)
	for i := 0; i < 100; i++ {
		val := r.Int31n(n)
		if val < 0 || val >= n {
			t.Fatalf("Int31n(%d) returned value outside [0, %d): %d", n, n, val)
		}
	}
}

func TestInt31nPowerOfTwo(t *testing.T) {
	r := &Rand{}
	n := int32(1024) // power of two
	for i := 0; i < 100; i++ {
		val := r.Int31n(n)
		if val < 0 || val >= n {
			t.Fatalf("Int31n(%d) returned value outside [0, %d): %d", n, n, val)
		}
	}
}

func TestInt63n(t *testing.T) {
	r := &Rand{}
	n := int64(1000000)
	for i := 0; i < 100; i++ {
		val := r.Int63n(n)
		if val < 0 || val >= n {
			t.Fatalf("Int63n(%d) returned value outside [0, %d): %d", n, n, val)
		}
	}
}

func TestShuffleNormal(t *testing.T) {
	r := &Rand{}
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	originalLen := len(slice)

	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

	if len(slice) != originalLen {
		t.Fatalf("Shuffle changed slice length from %d to %d", originalLen, len(slice))
	}

	// Verify all elements are still present
	seen := make(map[int]bool)
	for _, v := range slice {
		if v < 0 || v >= 10 {
			t.Fatalf("Shuffle produced invalid value: %d", v)
		}
		seen[v] = true
	}

	if len(seen) != 10 {
		t.Fatalf("Shuffle lost or duplicated elements, seen %d unique values", len(seen))
	}
}

func TestShuffleZero(t *testing.T) {
	r := &Rand{}
	slice := []int{1, 2, 3}

	// Should not panic with n=0
	defer func() {
		if recover() != nil {
			t.Fatalf("Shuffle panicked with n=0")
		}
	}()

	r.Shuffle(0, func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

	// Slice should be unchanged
	if len(slice) != 3 {
		t.Fatalf("Shuffle with n=0 changed slice length")
	}
}

func BenchmarkInt63(b *testing.B) {
	r := &Rand{}
	for i := 0; i < b.N; i++ {
		_ = r.Int63()
	}
}

func BenchmarkInt31(b *testing.B) {
	r := &Rand{}
	for i := 0; i < b.N; i++ {
		_ = r.Int31()
	}
}

func BenchmarkFloat64(b *testing.B) {
	r := &Rand{}
	for i := 0; i < b.N; i++ {
		_ = r.Float64()
	}
}

func BenchmarkIntn(b *testing.B) {
	r := &Rand{}
	for i := 0; i < b.N; i++ {
		_ = r.Intn(1000)
	}
}

func BenchmarkShuffle(b *testing.B) {
	r := &Rand{}
	slice := make([]int, 100)
	for i := range slice {
		slice[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Shuffle(len(slice), func(i, j int) {
			slice[i], slice[j] = slice[j], slice[i]
		})
	}
}
