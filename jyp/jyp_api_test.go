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
)
// Negative testcases/errors will be checked in a different file


//  go test -v -run Test_ObjPath_merged_expand__split_with_first_char
func Test_ObjPath_merged_expand__split_with_first_char(t *testing.T) {
	funName := "Test_ObjPath_merged_expand__split_with_first_char"
	testName := funName + "_basic"

	keys, _ := ObjPath_merged_expand__split_with_first_char("/a/b/c/d")
	compare_int_int(testName, 4, len(keys), t)
	compare_str_str(testName, "d", keys[3], t)
}

//  go test -v -run  TestJSON_value_AddKeyVal_path
func TestJSON_value_AddKeyVal_path(t *testing.T) {
	funName := "TestJSON_value_AddKeyVal_path"
	testName := funName + "_basic"

	src := `{"a": "A"}`
	root, _ := JsonParse(src)
	root.SetPath("/b", NewObj(), false)
	root.SetPath("/b/c", NewObj(), false)
	root.SetPath("/b/c/d", NewStr("Delta"), false)
	root.SetPath("/e/f/g/h", NewStr("Hugo"),true)
	root.SetPath("/array/a2", NewArr( NewNumInt(1), NewNumInt(2)), true)

	fmt.Println("root with new value, path-insert:")
	fmt.Println(root.Repr())

	compare_str_str(testName, "Delta", root.ValObject["b"].ValObject["c"].ValObject["d"].ValRunes, t)

	hVal, _ := root.GetPath("/e/f/g/h")
	compare_str_str(testName, "Hugo", hVal.ValRunes, t)

	arrVal, _ := root.GetPath("/array/a2")
	compare_int_int(testName, 2, arrVal.ValArray[1].ValNumberInt, t)

	hValWithKeys, _ := root.GetPathKeys([]string{"e", "f", "g", "h"})
	compare_str_str(testName, "Hugo", hValWithKeys.ValRunes, t)
}


//  go test -v -run Test_JsonParse_and_ObjPathKeys
func Test_JsonParse_and_ObjPathKeys(t *testing.T) {
	funName := "Test_JsonParse_and_ObjPathKeys"
	testName := funName + "_basic"
	_ = testName

	src := `{"a": {"b": {"c": "C"}}}`
	root, _ := JsonParse(src)
	elemC, _ := root.GetPath("/a/b/c")
	compare_str_str(testName, "C", elemC.ValRunes, t)
}

//  go test -v -run Test_JsonParse
func Test_JsonParse(t *testing.T) {
	funName := "Test_JsonParse"
	testName := funName + "_basic_obj"

	src := `{"a": "A"}`
	root, _ := JsonParse(src)
	compare_rune_rune(testName, '{', root.ValType, t)
	compare_int_int(testName, 1, len(root.ValObject), t) // has 1 elem
	compare_str_str(testName, "A", root.ValObject["a"].ValRunes, t)

	testName = funName + "_basic_arr"
	src = `["a", "A"]`
	root, _ = JsonParse(src)
	compare_rune_rune(testName, '[', root.ValType, t)
	compare_int_int(testName, 2, len(root.ValArray), t)          // has 1 elem
	compare_str_str(testName, "a", root.ValArray[0].ValRunes, t) // has 1 elem
	compare_str_str(testName, "A", root.ValArray[1].ValRunes, t) // has 1 elem
}
