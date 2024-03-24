/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

*/


package jyp

import (
	"fmt"
	"testing"
	"time"
)

//  go test -v -run  Test_tokensTableDetect_versionB
func Test_tokensTableDetect_versionB(t *testing.T) {
	funName := "Test_tokensTableDetect_versionB"
	testName := funName + "_basic"

	//src := `{"text":{"level2":[321,4.5,"string\"Escaped",true,false,null]}}`
	src := ""
	// src := file_read_to_string("large-file.json")


	timeStart := time.Now()
	// src = `{"a": "b"}`
	src = `{"a": "A", "b1": {"b2":"B2"}, "c":"C", "list":["k", "bh"]}`
	tokensTableB := tokensTableDetect_structuralTokens_strings(src)
	fmt.Println("token table creation time:", time.Since(timeStart))
	fmt.Println("tokensTableB")
	for _, tokenb := range tokensTableB {
		print_tokenB("tokenTable:", tokenb)
	}

	errorsCollected := []error{}
	root, _ := JSON_B_structure_building(src, tokensTableB, 0, errorsCollected)
	fmt.Println(root.Repr(2))

	_ = root
	_ = tokensTableB
	_ = testName

	// base__print_tokenElems(tokensTableB)
}


//  go test -v -run  Test_structure_building
func Test_structure_building(t *testing.T) {
	funName := "Test_structure_building"
	testName := funName + "_basic_obj"
	errorsCollected := []error{}

	src := `{"a": "A"}`
	tokensTableB := tokensTableDetect_structuralTokens_strings(src)
	errorsCollected = JSON_B_validation(tokensTableB)
	root, _ := JSON_B_structure_building(src, tokensTableB, 0, errorsCollected)
	compare_rune_rune(testName, '{', root.ValType, t)
	compare_int_int(testName, 1, len(root.ValObject), t) // has 1 elem
	compare_str_str(testName, "A", root.ValObject["a"].ValString, t)

	testName = funName + "_basic_arr"
	src = `["a", "A"]`
	tokensTableB = tokensTableDetect_structuralTokens_strings(src)
	errorsCollected = JSON_B_validation(tokensTableB)
	root, _ = JSON_B_structure_building(src, tokensTableB, 0, errorsCollected)
	compare_rune_rune(testName, '[', root.ValType, t)
	compare_int_int(testName, 2, len(root.ValArray), t) // has 1 elem
	compare_str_str(testName, "a", root.ValArray[0].ValString, t) // has 1 elem
	compare_str_str(testName, "A", root.ValArray[1].ValString, t) // has 1 elem

}


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
	for _, hexaRuneTestCase := range testElems {

		intValDetected, err := base__hexaRune_to_intVal(hexaRuneTestCase.hexaRune)
		compare_bool_bool(testName, hexaRuneTestCase.isValidHexaRune, err == nil, t)
		compare_int_int(testName, hexaRuneTestCase.decimalWantedVal, intValDetected, t)

	}
	intValDetected, err := base__hexaRune_to_intVal('b')
	compare_bool_bool(testName, true, err == nil, t)
	compare_int_int(testName, 11, intValDetected, t)

	intValDetected, err = base__hexaRune_to_intVal('m')
	compare_bool_bool(testName, true, err != nil, t)
}
