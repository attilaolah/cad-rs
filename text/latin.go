// Package text provides string replacers for Cyrillic / Gaj's Latin.
package text

import "strings"

var (
	// Cyrillic to Latin map (with digraphs).
	cyr2lat = map[rune]rune{
		'А': 'A',
		'Б': 'B',
		'В': 'V',
		'Г': 'G',
		'Д': 'D',
		'Ђ': 'Đ',
		'Е': 'E',
		'Ж': 'Ž',
		'З': 'Z',
		'И': 'I',
		'Ј': 'J',
		'К': 'K',
		'Л': 'L',
		'Љ': 'ǈ',
		'М': 'M',
		'Н': 'N',
		'Њ': 'ǋ',
		'О': 'O',
		'П': 'P',
		'Р': 'R',
		'С': 'S',
		'Т': 'T',
		'Ћ': 'Ć',
		'У': 'U',
		'Ф': 'F',
		'Х': 'H',
		'Ц': 'C',
		'Ч': 'Č',
		'Џ': 'ǅ',
		'Ш': 'Š',
	}

	// Map for removing digraphs.
	digraphs = map[rune]string{
		'ǅ': "Dž",
		'ǈ': "Lj",
		'ǋ': "Nj",
	}

	// ASCII transliteration map, mostly based on the "Serbian" column here:
	// https://en.wikipedia.org/wiki/Scientific_transliteration_of_Cyrillic#Table
	// Digraphs already covered by the "digraphs" map above are excluded.
	ascii = map[rune]string{
		'Đ': "Dj",  // Serbian, official
		'Ž': "Zh",  // Chech, unofficial
		'Ć': "Tj",  // Serbian, unofficial
		'Č': "Tsh", // Chech, unofficial
		'Š': "Sh",  // Chech, unofficial
	}

	// Azbuka, the Serbian Cyrillic alphabet.
	Azbuka = func(m map[rune]rune) []rune {
		ret := []rune{}
		for k := range m {
			ret = append(ret, k)
		}
		return ret
	}(cyr2lat)

	// ToLatin is a Cyrillic to Latin string replacer.
	// Uses Unicode digraphs so the result can be transliterated back.
	// Upper-case is translated to mixed-case digraphs. Use strings.ToUpper() to convert to all-caps.
	ToLatin = strings.NewReplacer(func(m map[rune]rune) []string {
		ret := []string{}
		for k, v := range m {
			ret = append(ret, strings.ToUpper(string(k)), strings.ToTitle(string(v)))
			ret = append(ret, strings.ToLower(string(k)), strings.ToLower(string(v)))
		}
		return ret
	}(cyr2lat)...)

	// ToASCII transliterates text to ASCII.
	// This should only be used to generate e.g. web-safe filenames.
	ToASCII = strings.NewReplacer(func(m map[rune]string) []string {
		ret := []string{}
		for k, v := range m {
			ret = append(ret, strings.ToUpper(string(k)), strings.ToTitle(v))
			ret = append(ret, strings.ToLower(string(k)), strings.ToLower(v))
		}
		return ret
	}(ascii)...)

	// RemoveDigraphs removes digraphs from the text.
	// E.g. the digraph 'ǆ' becomes two characters: 'd' followed by 'ž'.
	RemoveDigraphs = strings.NewReplacer(func(m map[rune]string) []string {
		ret := []string{}
		for k, v := range m {
			ret = append(ret, strings.ToUpper(string(k)), strings.ToUpper(v))
			ret = append(ret, strings.ToTitle(string(k)), strings.ToTitle(v))
			ret = append(ret, strings.ToLower(string(k)), strings.ToLower(v))
		}
		return ret
	}(digraphs)...)
)
