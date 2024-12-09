package text

func DetectEncoding(text string) map[string]int {
	charsets := map[string]int{
		"UTF-8":      0,
		"CP1251":     0,
		"KOI8-R":     0,
		"IBM866":     0,
		"ISO-8859-5": 0,
		"MAC":        0,
	}

	if len(text) == 0 {
		return charsets
	}

	utflower := 7
	utfupper := 5
	lowercase := 3
	uppercase := 1
	lastSimb := 0

	for a := 0; a < len(text); a++ {
		char := int(text[a])

		// non-russian characters
		if char < 128 || char > 256 {
			continue
		}

		// UTF-8
		if (lastSimb == 208) && ((char > 143 && char < 176) || char == 129) {
			charsets["UTF-8"] += (utfupper * 2)
		}
		if ((lastSimb == 208) && ((char > 175 && char < 192) || char == 145)) || (lastSimb == 209 && char > 127 && char < 144) {
			charsets["UTF-8"] += (utflower * 2)
		}

		// CP1251
		if (char > 223 && char < 256) || char == 184 {
			charsets["CP1251"] += lowercase
		}
		if (char > 191 && char < 224) || char == 168 {
			charsets["CP1251"] += uppercase
		}

		// KOI8-R
		if (char > 191 && char < 224) || char == 163 {
			charsets["KOI8-R"] += lowercase
		}
		if (char > 222 && char < 256) || char == 179 {
			charsets["KOI8-R"] += uppercase
		}

		// IBM866
		if (char > 159 && char < 176) || (char > 223 && char < 241) {
			charsets["IBM866"] += lowercase
		}
		if (char > 127 && char < 160) || char == 241 {
			charsets["IBM866"] += uppercase
		}

		// ISO-8859-5
		if (char > 207 && char < 240) || char == 161 {
			charsets["ISO-8859-5"] += lowercase
		}
		if (char > 175 && char < 208) || char == 241 {
			charsets["ISO-8859-5"] += uppercase
		}

		// MAC
		if char > 221 && char < 255 {
			charsets["MAC"] += lowercase
		}
		if char > 127 && char < 160 {
			charsets["MAC"] += uppercase
		}

		lastSimb = char
	}

	return charsets
}
