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
	"unicode"
)

type tokenElem_B struct {
	tokenType rune /* one rune is stored here to represent a unit in the source code
                      { objOpen
                      } objClose
                      [ arrayOpen
                      ] arrayClose
                      , comma
                      : colon
                      " string
                      0 digit
                      t true
                      f false
                      n null
                      ? not identified, only saved: later the type can be defined
	*/
	posInSrcFirst int
	posInSrcLast  int
}

type tokenElems []tokenElem_B

func tokensTableDetect_versionB(srcStr string) tokenElems{
	tokenTable := tokenElems{}
	posUnknownBlockStart := -1 // used only if the token is longer than 1 char. numbers, false/true for example
	
	//////////// TOKEN ADD func ///////////////////////
	tokenAdd := func (typeOfToken rune, posFirst, posLast int) {
		if posUnknownBlockStart != -1 {  // JSON has to be in containers {...} or [...] so it is closed with a known elem
			tokenTable = append(tokenTable, tokenElem_B{tokenType: '?', posInSrcFirst: posUnknownBlockStart, posInSrcLast: posFirst-1}  )
			posUnknownBlockStart = -1
		}
		tokenTable = append(tokenTable, tokenElem_B{tokenType: typeOfToken, posInSrcFirst: posFirst, posInSrcLast: posLast}  )
	} // func, tokenAdd //////////////////////////////
	
	
	posStringStart := -1  //////////////////////////////////////////
	inString := func () bool { // if string start position detected,
		return posStringStart != -1    // we are in String detection
	} //////////////////////////////////////////////////////////////
	/*
	inStringDebug := func(runeNow rune) string {
		info := " "
		if inString() || runeNow == '"' { info = "S" } // display S for debugging, when inString==true
		return info
	}
	 */
	isEscaped := false

	for pos, runeNow := range srcStr {

		stringCloseAtEnd := false
		if runeNow == '"' {
			if ! inString() {
				posStringStart = pos // posStringStart is modified only if interval is started
			} else { // in string processing:
				if ! isEscaped { // and not escaped:
					stringCloseAtEnd = true
				} // string can be closed only at the end of the codeBlock, not here.
			}     // so the closing " is part of the string section
		} /////////////////////////////////////

		// detect tokens:
		if ! inString() { // json structural chars:
			if runeNow == '{' || runeNow == '}' || runeNow == '[' || runeNow == ']' || runeNow == ',' || runeNow == ':' {
				tokenAdd(runeNow, pos, pos)
			} else {
				// not in string, and not json structural char
				// skip whitespaces:
				if ! base__is_whitespace_rune(runeNow) && posUnknownBlockStart == -1 {
					posUnknownBlockStart = pos // save block start
					// standard Json has to be closed with known
				}
			} // unknown token
		} // not inString


		///////////////////// CLOSING administration ////////////////
		if inString() {
			if runeNow == '\\' {
				isEscaped = ! isEscaped
			} else { // the escape series ended :-)
				isEscaped = false
			}
		}

		if stringCloseAtEnd {
			tokenAdd('"', posStringStart, pos)
			posStringStart = -1
		}
		// fmt.Println(fmt.Sprintf("pos: %2d", pos), string(runeNow), inStringDebug(runeNow), " token:", tokenAddedInForLoop)

	} // for, tokenTable
	return tokenTable
}

func base__is_whitespace_rune(oneRune rune) bool { // TESTED
	/*
		https://stackoverflow.com/questions/29038314/determining-whitespace-in-go
		func IsSpace

		func IsSpace(r rune) bool

		IsSpace reports whether the rune is a space character as defined by Unicode's White Space property; in the Latin-1 space this is

		'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).

		Other definitions of spacing characters are set by category Z and property Pattern_White_Space.
	*/
	return unicode.IsSpace(oneRune)
}

func base__print_tokenElems(tokenElems tokenElems) {
	for _, tokenElem := range tokenElems {
		fmt.Println(
			string(tokenElem.tokenType),
			fmt.Sprintf("%2d", tokenElem.posInSrcFirst),
			fmt.Sprintf("%2d", tokenElem.posInSrcLast),
			)
	}
}