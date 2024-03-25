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
	tokensTableB := stepA__tokensTableDetect_structuralTokens_strings_L1(src)
	fmt.Println("token table creation time:", time.Since(timeStart))

	fmt.Println("tokensTableB")
	for _, tokenb := range tokensTableB {
		tokenb.print("tokenTable:")
	}

	errorsCollected := stepB__JSON_B_validation_L1(tokensTableB)
	root, _ := stepC__JSON_B_structure_building__L1(src, tokensTableB, 0, errorsCollected)
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
	tokensTableB := stepA__tokensTableDetect_structuralTokens_strings_L1(src)
	errorsCollected = stepB__JSON_B_validation_L1(tokensTableB)
	root, _ := stepC__JSON_B_structure_building__L1(src, tokensTableB, 0, errorsCollected)
	compare_rune_rune(testName, '{', root.ValType, t)
	compare_int_int(testName, 1, len(root.ValObject), t) // has 1 elem
	compare_str_str(testName, "A", root.ValObject["a"].ValString, t)

	testName = funName + "_basic_arr"
	src = `["a", "A"]`
	tokensTableB = stepA__tokensTableDetect_structuralTokens_strings_L1(src)
	errorsCollected = stepB__JSON_B_validation_L1(tokensTableB)
	root, _ = stepC__JSON_B_structure_building__L1(src, tokensTableB, 0, errorsCollected)
	compare_rune_rune(testName, '[', root.ValType, t)
	compare_int_int(testName, 2, len(root.ValArray), t) // has 1 elem
	compare_str_str(testName, "a", root.ValArray[0].ValString, t) // has 1 elem
	compare_str_str(testName, "A", root.ValArray[1].ValString, t) // has 1 elem
}


//  go test -v -run  Test_structure_building_complex
func Test_structure_building_complex(t *testing.T) {
	funName := "Test_structure_building_complex"
	testName := funName + "_base"
	errorsCollected := []error{}

	src := `{"a": "A", "arr": ["0", "1", "2"], "obj": {"key": ["val"]} }`
	tokensTableB := stepA__tokensTableDetect_structuralTokens_strings_L1(src)
	tokensTableB.print()
	errorsCollected = stepB__JSON_B_validation_L1(tokensTableB)
	root, _ := stepC__JSON_B_structure_building__L1(src, tokensTableB, 0, errorsCollected)
	fmt.Println(root.Repr())
	compare_rune_rune(testName, '{', root.ValType, t)
	compare_int_int(testName, 3, len(root.ValObject), t) // has 1 elem
	compare_str_str(testName, "A", root.ValObject["a"].ValString, t)

	array := root.ValObject["arr"]
	compare_rune_rune(testName, '[', array.ValType, t)
	compare_int_int(testName, 3, len(array.ValArray), t) // has 1 elem
	compare_str_str(testName, "0", array.ValArray[0].ValString, t) // has 1 elem
	compare_str_str(testName, "1", array.ValArray[1].ValString, t) // has 1 elem
	compare_str_str(testName, "2", array.ValArray[2].ValString, t) // has 1 elem

	obj := root.ValObject["obj"]
	compare_rune_rune(testName, '{', obj.ValType, t)

	val := obj.ValObject["key"].ValArray[0]
	compare_rune_rune(testName, '"', val.ValType, t)
	compare_str_str(testName, "val", val.ValString, t) // has 1 elem
}


//  go test -v -run  Test_numbers_int
func Test_numbers_int(t *testing.T) {
	funName := "Test_numbers_int"
	testName := funName + "_base"
	errorsCollected := []error{}

	src := `{"age": -123, "favouriteNums": [4, 5, 6] } }`
	tokensTableB := stepA__tokensTableDetect_structuralTokens_strings_L1(src)
	tokensTableB.print()

	errorsCollected = stepB__JSON_B_validation_L1(tokensTableB)
	root, _ := stepC__JSON_B_structure_building__L1(src, tokensTableB, 0, errorsCollected)
	fmt.Println(root.Repr())

	compare_rune_rune(testName, '{', root.ValType, t)

	age := root.ValObject["age"]
	compare_int_int(testName, -123, age.ValNumberInt, t)

	nums := root.ValObject["favouriteNums"]
	compare_int_int(testName, 4, nums.ValArray[0].ValNumberInt, t)
	compare_int_int(testName, 5, nums.ValArray[1].ValNumberInt, t)
	compare_int_int(testName, 6, nums.ValArray[2].ValNumberInt, t)
}


//  go test -v -run  Test_numbers_float
func Test_numbers_float(t *testing.T) {
	funName := "Test_numbers_float"
	testName := funName + "_base"
	errorsCollected := []error{}

	src := `{"celsiusDegrees": [0.12, -3.45, 6.789] }`
	tokensTableB := stepA__tokensTableDetect_structuralTokens_strings_L1(src)
	tokensTableB.print()

	errorsCollected = stepB__JSON_B_validation_L1(tokensTableB)
	root, _ := stepC__JSON_B_structure_building__L1(src, tokensTableB, 0, errorsCollected)
	fmt.Println(root.Repr())

	compare_rune_rune(testName, '{', root.ValType, t)

	celsiusDegrees := root.ValObject["celsiusDegrees"]
	compare_flt_flt(testName, 0.12, celsiusDegrees.ValArray[0].ValNumberFloat, t)
	compare_flt_flt(testName, -3.45,celsiusDegrees.ValArray[1].ValNumberFloat, t)
	compare_flt_flt(testName, 6.789, celsiusDegrees.ValArray[2].ValNumberFloat, t)
}
