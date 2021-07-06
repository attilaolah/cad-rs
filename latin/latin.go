// Package latin provides string replacers for Cyrillic / Gaj's Latin.
package latin

import "strings"

// ToLatin is a Cyrillic to Latin string replacer.
// Uses Unicode digraphs so the result can be transliterated back.
// Upper-case is translated to mixed-case digraphs. Use strings.ToUpper() to convert to all-caps.
var ToLatin = strings.NewReplacer(func(m map[rune]rune) []string {
	ret := []string{}
	for k, v := range m {
		ret = append(ret, string(k), string(v))
	}
	return ret
}(map[rune]rune{
	// Upper-case:
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
	// Lower-case:
	'а': 'a',
	'б': 'b',
	'в': 'v',
	'г': 'g',
	'д': 'd',
	'ђ': 'đ',
	'е': 'e',
	'ж': 'ž',
	'з': 'z',
	'и': 'i',
	'ј': 'j',
	'к': 'k',
	'л': 'l',
	'љ': 'ǉ',
	'м': 'm',
	'н': 'n',
	'њ': 'ǌ',
	'о': 'o',
	'п': 'p',
	'р': 'r',
	'с': 's',
	'т': 't',
	'ћ': 'ć',
	'у': 'u',
	'ф': 'f',
	'х': 'h',
	'ц': 'c',
	'ч': 'č',
	'џ': 'ǆ',
	'ш': 'š',
})...)

// RemoveDigraphs removes digraphs from the text.
// E.g. the digraph 'ǆ' becomes two characters: 'd' followed by 'ž'.
var RemoveDigraphs = strings.NewReplacer(func(m map[rune]string) []string {
	ret := []string{}
	for k, v := range m {
		ret = append(ret, string(k), v)
	}
	return ret
}(map[rune]string{
	// All-caps:
	'Ǆ': "DŽ",
	'Ǉ': "LJ",
	'Ǌ': "NJ",
	// Upper-case:
	'ǅ': "Dž",
	'ǈ': "Lj",
	'ǋ': "Nj",
	// Lower-case:
	'ǆ': "dž",
	'ǉ': "lj",
	'ǌ': "nj",
})...)
