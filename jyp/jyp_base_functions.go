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
	"fmt"
	"unicode"
)

// runesSections were checked against illegal chars, so here digitRune is in 0123456789
// TODO: maybe can be removed if not used in the future in exponent number detection section
func base__digit10BasedRune_integer_value(digit10based rune) (int, error) {
	if digit10based == '0' {
		return 0, nil
	}
	if digit10based == '1' {
		return 1, nil
	}
	if digit10based == '2' {
		return 2, nil
	}
	if digit10based == '3' {
		return 3, nil
	}
	if digit10based == '4' {
		return 4, nil
	}
	if digit10based == '5' {
		return 5, nil
	}
	if digit10based == '6' {
		return 6, nil
	}
	if digit10based == '7' {
		return 7, nil
	}
	if digit10based == '8' {
		return 8, nil
	}
	if digit10based == '9' {
		return 9, nil
	}
	return 0, errors.New(errorPrefix + "rune (" + string(digit10based) + ")")
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

// the string has whitespace chars only
func base__is_whitespace_string(src string) bool { // TESTED
	for _, runeFromStr := range src {
		if !base__is_whitespace_rune(runeFromStr) {
			return false
		}
	}
	return true
}

// TODO: TEST
func base__compare_runes_are_equal(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for pos, runeA := range a {
		runeB := b[pos]
		if runeA != runeB {
			return false
		}
	}
	return true
}

// create a separated copy about original rune Slice into a new variable (deepcopy)
func base__runes_copy(runes []rune) []rune { // TESTED
	runesNew := []rune{}
	for _, r := range runes {
		runesNew = append(runesNew, r)
	}
	return runesNew
}

// split once, at first occurance
func base__runes_split_at_pattern(runes []rune, splitterRune rune) ([]rune, []rune) {
	runesBefore := []rune{}
	runesAfter := []rune{}
	splitterDetected := false
	for _, r := range runes {
		if !splitterDetected && r == splitterRune { // split only at the first occurance
			splitterDetected = true
			continue
		}
		if splitterDetected {
			runesAfter = append(runesAfter, r)
		} else {
			runesBefore = append(runesBefore, r)
		}
	}
	return runesBefore, runesAfter
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

type wordInSrc struct {
	wordChars []rune
	posFirst  int
	posLast   int
}

// get the rune IF the index is really in the range of the src.
// return with ' ' space, IF the index is NOT in the range.
// reason: avoid never ending index checking, so do it only once
// the space can be answered because this func is used when a real char wanted to be detected,
// and if a space is returned, this has NO MEANING in that parse section
// this fun is NOT used in string detection - and other places whitespaces can be neglected, too
// getChar, with whitespace replace
func base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(src []rune, pos int) rune { // TESTED
	char := base__srcGetChar__safeOverindexing(src, pos)
	if base__is_whitespace_rune(char) {
		return ' ' // simplify everything. if the char is ANY whitespace char,
		// return with SPACE, this is not important in the source code parsing
	}
	return char
}

// getChar, no whitespace replace
func base__srcGetChar__safeOverindexing(src []rune, pos int) rune { // TESTED
	posPossibleMax := len(src) - 1  // if src is empty, max is -1,
	posPossibleMin := 0             // and the condition cannot be true here:
	if (pos >= posPossibleMin) && (pos <= posPossibleMax) {
		return src[pos]
	}
	return ' '
}

// give back words (plus posFirst/posLast info)
func base__src_get_whitespace_separated_words_posFirst_posLast(src []rune) []wordInSrc { // TESTED

	words := []wordInSrc{}

	wordChars := []rune{}
	posFirst := -1
	posLast := -1

	// posActual := -1, len(src) + 1: overindexing!
	// with this, I can be sure that minimum one space is detected first,
	// and minimum one space detected after the source code's normal chars!
	// with this solution, the last word detection can be closed with the last boundary space, in one
	// case, and I don't have to handle that later, in a second if/else condition

	// base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces() handles the overindexing
	for posActual := -1; posActual < len(src)+1; posActual++ {
		runeActual := base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(src, posActual)

		// the first and last chars, because of overindexing, are spaces, this is guaranteed!
		if runeActual == ' ' {  // means IF WHITESPACE: every whitespace is replaced to simple space in srcGetChar
			if len(wordChars) > 0 {
				word := wordInSrc{
					wordChars: wordChars,
					posFirst: posFirst,
					posLast:  posLast,
				}
				words = append(words, word)

				wordChars = []rune{}
				posFirst = -1
				posLast = -1
			}

		} else { // non-space rune:
			// save posFirst, posLast, and word-builder chars ///
			if len(wordChars) == 0 {
				posFirst = posActual
			}
			posLast = posActual
			wordChars = append(wordChars, runeActual)
		}

	}

	return words
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

// are the Runes in the set?
func base__validate_runes_are_in_allowed_set(runesToValidate, runesAllowed []rune) bool {
	for _, r := range runesToValidate {
		if !base__validate_rune_are_in_allowed_set(r, runesAllowed) {
			return false
		}
	}
	return true
}

// TODO: test
func base__runeInRunes(runeWanted rune, runesToValidate []rune) bool {
	for _, r := range runesToValidate {
		if r == runeWanted {
			return true
		}
	}
	return false
}

// test/Debug Helper - display Tokens table
func TokensDisplay_startingCoords(srcOrig []rune, tokens tokenTable_startPositionIndexed) {
	keys := local_tool__tokenTable_position_keys_sorted(tokens)

	fmt.Println("== Tokens Table display ==")
	for _, key := range keys{
		fmt.Println(string(srcOrig[tokens[key].CharPositionFirstInSourceCode:tokens[key].CharPositionLastInSourceCode+1]), key, tokens[key])
	}
}
