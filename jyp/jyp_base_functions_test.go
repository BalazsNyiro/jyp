/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

*/

package jyp

import "testing"


//  go test -v -run Test_hexaRune_to_intVal
func Test_hexaRune_to_intVal(t *testing.T) {
	funName := "Test_hexaRune_to_intVal"
	testName := funName + "hexa2int_conversation"

	intValDetected, err := hexaRune_to_intVal('b')
	compare_bool_bool(testName, true, err == nil, t)
	compare_int_int(testName, 11, intValDetected, t)

	intValDetected, err = hexaRune_to_intVal('m')
	compare_bool_bool(testName, true, err != nil, t)
}


func Test_separator_set_if_no_last_elem(t *testing.T) {
	funName := "Test_separator_set_if_no_last_elem"
	testName := funName + "_base"

	sep := ","
	elems := []int{0,1,2}

	separator := separator_set_if_no_last_elem(3, len(elems), sep)
	compare_str_str(testName, "", separator, t)

	separator = separator_set_if_no_last_elem(2, len(elems), sep)
	compare_str_str(testName, "", separator, t)

	separator = separator_set_if_no_last_elem(1, len(elems), sep)
	compare_str_str(testName, ",", separator, t)
}