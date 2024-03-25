/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

*/


package jyp

import "testing"

//  go test -v -run Test_JsonParse_and_ObjPathKeys
func Test_JsonParse_and_ObjPathKeys(t *testing.T) {
	funName := "Test_JsonParse_and_ObjPathKeys"
	testName := funName + "_basic"
	_ = testName

	src := `{"a": {"b": {"c": "C"}}}`
	root, _ := JsonParse(src)
	elemC, _ := root.ObjPath("/a/b/c")
	compare_str_str(testName, "C", elemC.ValString, t)
}

//  go test -v -run Test_JsonParse
func Test_JsonParse(t *testing.T) {
	funName := "Test_JsonParse"
	testName := funName + "_basic_obj"

	src := `{"a": "A"}`
	root, _ := JsonParse(src)
	compare_rune_rune(testName, '{', root.ValType, t)
	compare_int_int(testName, 1, len(root.ValObject), t) // has 1 elem
	compare_str_str(testName, "A", root.ValObject["a"].ValString, t)

	testName = funName + "_basic_arr"
	src = `["a", "A"]`
	root, _ = JsonParse(src)
	compare_rune_rune(testName, '[', root.ValType, t)
	compare_int_int(testName, 2, len(root.ValArray), t) // has 1 elem
	compare_str_str(testName, "a", root.ValArray[0].ValString, t) // has 1 elem
	compare_str_str(testName, "A", root.ValArray[1].ValString, t) // has 1 elem
}
