/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

*/

package jyp

import "testing"

// go test -v -run Test_hexaRune_to_intVal
func Test_hexaRune_to_intVal(t *testing.T) {
	funName := "Test_hexaRune_to_intVal"
	testName := funName + "hexa2int_conversation"

	type hexaRuneTest struct {
		hexaRune rune
		isValidHexaRune bool
		decimalWantedVal int
	}
	testElems := []hexaRuneTest{
		{'0', true, 0},
		{'1', true, 1},
		{'2', true, 2},
		{'3', true, 3},
		{'4', true, 4},
		{'5', true, 5},
		{'6', true, 6},
		{'7', true, 7},
		{'8', true, 8},
		{'9', true, 9},
		{'a', true, 10},
		{'b', true, 11},
		{'c', true, 12},
		{'d', true, 13},
		{'e', true, 14},
		{'f', true, 15},
	}

	for _, hexaRuneTestCase := range testElems { // test all possible elems
		intValDetected, err := base__hexaRune_to_intVal(hexaRuneTestCase.hexaRune)
		compare_bool_bool(testName, hexaRuneTestCase.isValidHexaRune, err == nil, t)
		compare_int_int(testName, hexaRuneTestCase.decimalWantedVal, intValDetected, t)
	}
	// one manual test
	intValDetected, err := base__hexaRune_to_intVal('b')
	compare_bool_bool(testName, true, err == nil, t)
	compare_int_int(testName, 11, intValDetected, t)

	intValDetected, err = base__hexaRune_to_intVal('m')
	compare_bool_bool(testName, true, err != nil, t)
}
