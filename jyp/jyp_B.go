/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

*/

package jyp

import "fmt"

type tokenTable_verB struct {
	tokenType rune /* one rune is stored here to represent a unit in the source code
                      [ arrayOpen
                      ] arrayClose
                      { objOpen
                      } objClose
                      , comma
                      : colon
                      " string
                      0 digit

	*/
}

func tokensTableDetect_versionB(srcStr string) string {
	inString := false
	isEscaped := false

	for pos, runeNow := range srcStr {

		stringCloser := false
		if runeNow == '"' {
			if ! inString {
				inString = true
			} else { // in string processing:
				if ! isEscaped {
					stringCloser = true
				}
			}
		}

		fmt.Println("pos:", pos, string(runeNow), inString)

		if inString {
			if runeNow == '\\' {
				isEscaped = ! isEscaped
			} else { // the escape series ended :-)
				isEscaped = false
			}
		}
		if stringCloser { inString = false }
	}
	return "XXX"
}