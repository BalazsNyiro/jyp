/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

*/

package jyp

import "testing"


func Test_runes_copy(t *testing.T) {
	funName := "Test_runes_copy"
	testName := funName + "_base"

	sample := "broken_mirror"
	runesCopy := []rune(sample)
	compare_runes_runes(testName, []rune(sample), runesCopy, t)
}



func Test_is_whitespace_string_rune(t *testing.T) {
	funName := "Test_is_whitespace_string_rune"

	src := "abc\t\n12 -1.2"

	// whitespace string tests
	testName := funName + "_simpleStringWithWhitespaceEnding"
	isWhitespace := base__is_whitespace_string(src[0:5])
	compare_bool_bool(testName, false, isWhitespace, t)

	testName = funName + "_simpleStringOnlyWhitespace"
	isWhitespace = base__is_whitespace_string(src[3:5])
	compare_bool_bool(testName, true, isWhitespace, t)

	// whitespace rune tests
	testName = funName + "_simpleRuneWhitespace"
	runeSelected := rune(src[4])
	isWhitespace = base__is_whitespace_rune(runeSelected)
	compare_bool_bool(testName, true, isWhitespace, t)

	testName = funName + "_simpleRuneNonWhitespace"
	runeSelected = rune(src[6])
	isWhitespace = base__is_whitespace_rune(runeSelected)
	compare_bool_bool(testName, false, isWhitespace, t)
	compare_rune_rune(testName, '2', runeSelected, t)
}

// go test -v -run Test_hexaRune_to_intVal
func Test_hexaRune_to_intVal(t *testing.T) {
	funName := "Test_hexaRune_to_intVal"
	testName := funName + "hexa2int_conversation"

	intValDetected, err := base__hexaRune_to_intVal('b')
	compare_bool_bool(testName, true, err == nil, t)
	compare_int_int(testName, 11, intValDetected, t)

	intValDetected, err = base__hexaRune_to_intVal('m')
	compare_bool_bool(testName, true, err != nil, t)
}

func Test_separator_set_if_no_last_elem(t *testing.T) {
	funName := "Test_separator_set_if_no_last_elem"
	testName := funName + "_base"

	sep := ","
	elems := []int{0, 1, 2}

	separator := base__separator_set_if_no_last_elem(3, len(elems), sep)
	compare_str_str(testName, "", separator, t)

	separator = base__separator_set_if_no_last_elem(2, len(elems), sep)
	compare_str_str(testName, "", separator, t)

	separator = base__separator_set_if_no_last_elem(1, len(elems), sep)
	compare_str_str(testName, ",", separator, t)
}
