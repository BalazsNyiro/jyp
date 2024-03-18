/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

*/


package jyp

import "testing"

//  go test -v -run  Test_tokensTableDetect_versionB
func Test_tokensTableDetect_versionB(t *testing.T) {
	funName := "Test_tokensTableDetect_versionB"
	testName := funName + "_basic"

	src := `{"text":{"level2":[321,4.5,"string\"Escaped",true,false,null]}}`

	tokensTableB := tokensTableDetect_versionB(src)
	_ = tokensTableB
	_ = testName
}