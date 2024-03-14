/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

*/

package jyp

import "testing"

func Test_rune_runes_in_allowed_set(t *testing.T) {
	funName := "Test_rune_runes_in_allowed_set"
	testName := funName + "_base"

	allowedRunes := []rune("abc")
	runeInAllowedSet := base__validate_rune_are_in_allowed_set('a', allowedRunes)
	compare_bool_bool(testName, true, runeInAllowedSet, t)

	runeInAllowedSet = base__validate_rune_are_in_allowed_set('x', allowedRunes)
	compare_bool_bool(testName, false, runeInAllowedSet, t)

	runesInAllowedSet := base__validate_runes_are_in_allowed_set([]rune("cab"), allowedRunes)
	compare_bool_bool(testName, true, runesInAllowedSet , t)

	runesInAllowedSet = base__validate_runes_are_in_allowed_set([]rune("abba"), allowedRunes)
	compare_bool_bool(testName, true, runesInAllowedSet , t)

	runesInAllowedSet = base__validate_runes_are_in_allowed_set([]rune("notinset"), allowedRunes)
	compare_bool_bool(testName, false, runesInAllowedSet , t)
}

func Test_digitIntegerValue(t *testing.T) {
	funName := "Test_digitIntegerValue"
	testName := funName + "_base"

	val10Based, err := base__digit10BasedRune_integer_value('0')
	compare_bool_bool(testName, true, err==nil, t)
	compare_int_int(testName, 0, val10Based, t)

	val10Based, err = base__digit10BasedRune_integer_value('X')
	compare_bool_bool(testName,false, err==nil, t)
}


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
