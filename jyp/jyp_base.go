/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

*/

package jyp

import (
	"errors"
	"unicode"
)


func base__hexaRune_to_intVal(hexaChar rune) (int, error) { // TESTED
	hexaTable := map[rune]int{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'a': 10,
		'b': 11,
		'c': 12,
		'd': 13,
		'e': 14,
		'f': 15,
	}
	base10Val, keyInHexaTable := hexaTable[hexaChar]
	if keyInHexaTable {
		return base10Val, nil
	}
	return 0, errors.New("hexa char(" + string(hexaChar) + ") was not in hexa table")
}

func base__is_whitespace_rune(oneRune rune) bool { // TESTED
	/*  https://stackoverflow.com/questions/29038314/determining-whitespace-in-go
	'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).
	Other definitions of spacing characters are set by category Z and property Pattern_White_Space. */
	return unicode.IsSpace(oneRune)
}

