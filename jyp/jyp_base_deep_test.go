/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

*/

package jyp

import "testing"

// go test -v -run  Test_base__srcGetChar__safeOverindexing__deep
func Test_base__srcGetChar__safeOverindexing__deep(t *testing.T) {
	funName := "Test_base__srcGetChar__safeOverindexing"
	testName := funName + "_base"

	txt := []rune("") // this is the smallest string that can be passed
	// posPossibleMax == -1
	// posPossibleMin == 0

	// the posWanted >= 0 && posWanted <= -1 ?
	charRead := base__srcGetChar__safeOverindexing(txt, 0)
	compare_rune_rune(testName,' ', charRead, t)

	// what if we get negative index?
	charRead = base__srcGetChar__safeOverindexing(txt, -1)
	// in that case, posWanted>=0 is the guard.
	compare_rune_rune(testName,' ', charRead, t)

	txt = []rune("ww")
	charRead = base__srcGetChar__safeOverindexing(txt, 999)
	compare_rune_rune(testName,' ', charRead, t)

	txt = []rune("abc") // and a normal test, not corner-case
	charRead = base__srcGetChar__safeOverindexing(txt, 2)
	compare_rune_rune(testName,'c', charRead, t)
}
