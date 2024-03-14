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
// are the Runes in the set?
func base__validate_runes_are_in_allowed_set(runesToValidate, runesAllowed []rune) bool {
	for _, r := range runesToValidate {
		if ! base__validate_rune_are_in_allowed_set(r, runesAllowed) {
			return false
		}
	}
	return true
}


// is the rune in allowed set?
func base__validate_rune_are_in_allowed_set(runeValidated rune, runesAllowed []rune) bool {
	for _, r := range runesAllowed {
		if r == runeValidated {
			return true
		}
	}
	return false
}

// runesSections were checked against illegal chars, so here digitRune is in 0123456789
// TODO: maybe can be removed if not used in the future in exponent number detection section
func base__digit10BasedRune_integer_value(digit10based rune) (int, error) {
	if digit10based == '0' { return 0, nil }
	if digit10based == '1' { return 1, nil }
	if digit10based == '2' { return 2, nil }
	if digit10based == '3' { return 3, nil }
	if digit10based == '4' { return 4, nil }
	if digit10based == '5' { return 5, nil }
	if digit10based == '6' { return 6, nil }
	if digit10based == '7' { return 7, nil }
	if digit10based == '8' { return 8, nil }
	if digit10based == '9' { return 9, nil }
	return 0, errors.New(errorPrefix + "rune ("+string(digit10based)+")")
}


// create a separated copy about original rune Slice into a new variable (deepcopy)
func base__runes_copy(runes []rune) []rune {  // TESTED
	runesNew := []rune{}
	for _, r := range runes {
		runesNew = append(runesNew, r)
	}
	return runesNew
}


// the string has whitespace chars only
func base__is_whitespace_string(src string) bool { // TESTED
	for _, runeFromStr := range src {
		if ! base__is_whitespace_rune(runeFromStr) {
			return false
		}
	}
	return true
}

// the rune is a whitespace char
func base__is_whitespace_rune(oneRune rune) bool { // TESTED
	/*
		https://stackoverflow.com/questions/29038314/determining-whitespace-in-go
		func IsSpace

		func IsSpace(r rune) bool

		IsSpace reports whether the rune is a space character as defined by Unicode's White Space property; in the Latin-1 space this is

		'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).

		Other definitions of spacing characters are set by category Z and property Pattern_White_Space.
	*/
	return unicode.IsSpace(oneRune)
}

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

// return with a separator if position is not last elem, (position is before the last)
// or with empty string if last elem is reached
// position is 0 based
func base__separator_set_if_no_last_elem(position, length_numOfAllElems int, separator string) string {
	if position < length_numOfAllElems-1 {
		return separator
	}
	return ""
}
