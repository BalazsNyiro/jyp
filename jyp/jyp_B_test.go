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
	tokensTableB := tokensTableDetect_versionB(src)
	fmt.Println("token table creation time:", time.Since(timeStart))
	fmt.Println("tokensTableB")
	for _, tokenb := range tokensTableB {
		print_tokenB("tokenTable:", tokenb)
	}

	root, _, _ := JSON_B_structure_building(src, tokensTableB, 0)
	fmt.Println(root.repr(0))

	_ = root
	_ = tokensTableB
	_ = testName

	// base__print_tokenElems(tokensTableB)
}


//  go test -v -run  Test_structure_building_simple
func Test_structure_building(t *testing.T) {
	funName := ""
	testName := funName + "_basic_obj"

	src := `{"a": "A"}`
	tokensTableB := tokensTableDetect_versionB(src)
	root, _, _ := JSON_B_structure_building(src, tokensTableB, 0)
	compare_rune_rune(testName, '{', root.ValType, t)
	compare_int_int(testName, 1, len(root.ValObject), t) // has 1 elem
	compare_str_str(testName, "A", root.ValObject["a"].repr(), t)


	testName = funName + "_basic_arr"
	src = `["a", "A"]`
	tokensTableB = tokensTableDetect_versionB(src)
	root, _, _ = JSON_B_structure_building(src, tokensTableB, 0)
	compare_rune_rune(testName, '[', root.ValType, t)
	compare_int_int(testName, 2, len(root.ValArray), t) // has 1 elem
}
