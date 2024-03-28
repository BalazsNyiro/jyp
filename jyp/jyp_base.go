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


// repeat the wanted unit prefix a few times
func base__prefixGenerator_for_repr(oneUnitPrefix string, repeatNum int) string { // TESTED
	if oneUnitPrefix == "" {
		return "" // if there is nothing to repeat
	}
	out := ""
	for i:=0; i<repeatNum; i++ {
		out += oneUnitPrefix
	}
	return out
}

// first/last char removal can be important with "strings" if quotes are not important
func base__read_sourceCode_section_basedOnTokenPositions(src []rune, token tokenElem, removeFirstLastChar bool) []rune{ // TESTED
	if !removeFirstLastChar {
		return src[token.posInSrcFirst:token.posInSrcLast+1]
	}
	return src[token.posInSrcFirst+1:token.posInSrcLast]
}

// for list printing, set comma as a separator if NOT the last elem is printed.
// after the last elem, separator has to be empty
func base__separator_set_if_no_last_elem(position, length_numOfAllElems int, separator string) string { // TESTED
	if position < length_numOfAllElems-1 {
		return separator
	}
	return ""
}


// get the rune IF the index is really in the range of the src.
// return with ' ' space, IF the index is NOT in the range.
// reason: avoid never ending index checking, so do it only once
// the space can be answered because this fun is used when a real char wanted to be detected,
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


/* https://stackoverflow.com/questions/30263607/how-to-get-a-single-unicode-character-from-string
If the string is encoded in UTF-8, there is no direct way to access the nth rune of the string,
because the size of the runes (in bytes) is not constant.

performance: the string->rune conversion used too often. is it possible to do the conversion once?
so src has to be a string

*/
// getChar, no whitespace replace
func base__srcGetChar__safeOverindexing(src []rune, pos int) rune { //DEEP-TESTED
	posPossibleMax := len(src) - 1  // if src is empty, max is -1,
	posPossibleMin := 0             // and the condition cannot be true here:
	if (pos >= posPossibleMin) && (pos <= posPossibleMax) {
		return src[pos]
	}
	return ' '
}


