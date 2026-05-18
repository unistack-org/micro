package text

import "testing"

func TestDetectEncoding(t *testing.T) {
	t.Run("empty string returns all zeros", func(t *testing.T) {
		result := DetectEncoding("")
		if len(result) == 0 {
			t.Fatal("expected non-empty map")
		}
		for k, v := range result {
			if v != 0 {
				t.Fatalf("expected 0 score for empty input, key %q got %d", k, v)
			}
		}
	})

	t.Run("returns all six charsets", func(t *testing.T) {
		result := DetectEncoding("test")
		expected := []string{"UTF-8", "CP1251", "KOI8-R", "IBM866", "ISO-8859-5", "MAC"}
		for _, k := range expected {
			if _, ok := result[k]; !ok {
				t.Fatalf("missing charset key %q in result", k)
			}
		}
	})

	t.Run("ASCII-only string", func(t *testing.T) {
		result := DetectEncoding("Hello, world!")
		if result == nil {
			t.Fatal("expected non-nil map")
		}
		if _, ok := result["UTF-8"]; !ok {
			t.Fatal("expected UTF-8 key in result")
		}
	})

	t.Run("non-empty string with high bytes", func(t *testing.T) {
		// Use a string with bytes in 128-256 range to exercise charset detection
		b := make([]byte, 10)
		for i := range b {
			b[i] = byte(200 + i)
		}
		result := DetectEncoding(string(b))
		if result == nil {
			t.Fatal("expected non-nil result")
		}
	})

	t.Run("UTF-8 uppercase detection", func(t *testing.T) {
		// UTF-8 uppercase: 0xD0 (208) followed by 0x90-0xAF (144-175) or 0x91 (129)
		b := []byte{0xD0, 0x90} // Cyrillic A in UTF-8
		result := DetectEncoding(string(b))
		if result["UTF-8"] == 0 {
			t.Fatal("expected non-zero UTF-8 score for uppercase")
		}
	})

	t.Run("UTF-8 lowercase detection", func(t *testing.T) {
		// UTF-8 lowercase: 0xD0 (208) followed by 0xB0-0xBF (176-191) or 0x91 (129)
		// OR 0xD1 (209) followed by 0x80-0x8F (128-143)
		b := []byte{0xD0, 0xB0} // Cyrillic a in UTF-8
		result := DetectEncoding(string(b))
		if result["UTF-8"] == 0 {
			t.Fatal("expected non-zero UTF-8 score for lowercase")
		}
	})

	t.Run("CP1251 charset detection", func(t *testing.T) {
		// CP1251 lowercase: 224-255 or 184
		b := []byte{0xE0} // 224
		result := DetectEncoding(string(b))
		if result["CP1251"] == 0 {
			t.Fatal("expected non-zero CP1251 score")
		}
	})

	t.Run("KOI8-R charset detection", func(t *testing.T) {
		// KOI8-R lowercase: 192-223 or 163
		b := []byte{0xC0} // 192
		result := DetectEncoding(string(b))
		if result["KOI8-R"] == 0 {
			t.Fatal("expected non-zero KOI8-R score")
		}
	})

	t.Run("IBM866 charset detection", func(t *testing.T) {
		// IBM866 lowercase: 160-175 or 224-240
		b := []byte{0xA0} // 160
		result := DetectEncoding(string(b))
		if result["IBM866"] == 0 {
			t.Fatal("expected non-zero IBM866 score")
		}
	})

	t.Run("ISO-8859-5 charset detection", func(t *testing.T) {
		// ISO-8859-5 lowercase: 208-239 or 161
		b := []byte{0xD0} // 208
		result := DetectEncoding(string(b))
		if result["ISO-8859-5"] == 0 {
			t.Fatal("expected non-zero ISO-8859-5 score")
		}
	})

	t.Run("MAC charset detection", func(t *testing.T) {
		// MAC lowercase: 222-254
		b := []byte{0xDE} // 222
		result := DetectEncoding(string(b))
		if result["MAC"] == 0 {
			t.Fatal("expected non-zero MAC score")
		}
	})

	t.Run("bytes outside 128-256 range ignored", func(t *testing.T) {
		result1 := DetectEncoding("Hello")
		result2 := DetectEncoding("Hello\x00\x01\x7F\x100")
		// Both should return all zeros since no valid cyrillic bytes
		for k := range result1 {
			if result1[k] != result2[k] {
				t.Fatalf("expected same scores, got different for %q", k)
			}
		}
	})

	t.Run("CP1251 uppercase detection", func(t *testing.T) {
		// CP1251 uppercase: 192-223 or 168
		b := []byte{0xC0} // 192
		result := DetectEncoding(string(b))
		if result["CP1251"] == 0 {
			t.Fatal("expected non-zero CP1251 score for uppercase")
		}
	})

	t.Run("multiple charset scoring", func(t *testing.T) {
		// This should hit multiple charset detections
		b := []byte{0xD0, 0xB0, 0xC0, 0xA0}
		result := DetectEncoding(string(b))
		totalScore := 0
		for _, v := range result {
			totalScore += v
		}
		if totalScore == 0 {
			t.Fatal("expected non-zero total score")
		}
	})
}
