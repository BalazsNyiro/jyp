/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

*/

package jyp

import (
	"errors"
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

type tokenElems_B []tokenElem_B

func tokensTableDetect_versionB(srcStr string) tokenElems_B {
	tokenTable := tokenElems_B{}
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
				// so it can be a number, true/false, or
				// whitespaces:
				if posUnknownBlockStart == -1 {
					posUnknownBlockStart = pos // save block start
					// standard Json has to be closed with known
				}
			} // unknown token

		} else { // inString:
			///////////////////// CLOSING administration ////////////////
			if stringCloseAtEnd {
				tokenAdd('"', posStringStart, pos)
				posStringStart = -1
			} else { // not string closing
				if runeNow == '\\' {
					isEscaped = ! isEscaped
				} else { // the escape series ended :-)
					isEscaped = false
				}
			}
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

func base__print_tokenElems(tokenElems tokenElems_B) {
	for _, tokenElem := range tokenElems {
		fmt.Println(
			string(tokenElem.tokenType),
			fmt.Sprintf("%2d", tokenElem.posInSrcFirst),
			fmt.Sprintf("%2d", tokenElem.posInSrcLast),
			)
	}
}

func JSON_B_validation(tokenTableB tokenElems_B) []error {
	// TODO: loop over the table, find incorrect {..} [..], ".." pairs,
	// incorrect numbers, everything that can be a problem
	// so after this, we can be sure that every elem are in pairs.
	return []error{}
}


func getTextFromSrc(src string, token tokenElem_B) string {
	return src[token.posInSrcFirst:token.posInSrcLast+1]
}

func print_tokenB(prefix string, t tokenElem_B) {
	fmt.Println(prefix, string(t.tokenType), t.posInSrcFirst, t.posInSrcLast)
}


// return with pos only to avoid elem copy with reading/passing
func build__find_next_token_pos(typeWanted rune, posActual int, tokensTable tokenElems_B) (int, error) {
	var pos int
	for pos = posActual+1; pos<len(tokensTable); pos++ {
		if tokensTable[pos].tokenType == typeWanted {
			return pos, nil
		}
	}
	return pos, errors.New("wanted token is not detected in table")
}

func JSON_B_structure_building(src string, tokensTableB tokenElems_B, tokenPosStart int) (JSON_value_B, []error, int) {
	errors := JSON_B_validation(tokensTableB)
	if len(errors) > 0 {
		return JSON_value_B{}, errors, 0
	}
	elem := JSON_value_B{}
	var pos int
	for pos = tokenPosStart; pos<len(tokensTableB); pos++ {
		tokenNow := tokensTableB[pos]
		print_tokenB("for head", tokenNow)

		if tokenNow.tokenType == '{' {
			elem = NewObj_JSON_value_B()
			var objKey string
			var posInObj int

			// find the next string
			/*
			for posInObj = pos+1; posInObj<len(tokensTableB); posInObj++ {
				if tokensTableB[posInObj].tokenType == '"' {
					objKey = getTextFromSrc(src, tokensTableB[posInObj])
					fmt.Println("find first key:", objKey)
					break // read the next string, that will be the key
			}} // for, posInObj

			 */
			// todo: error handling
			posInObj, _ = build__find_next_token_pos('"', pos, tokensTableB)
			objKey = getTextFromSrc(src, tokensTableB[posInObj])

			fmt.Println("find first key:", objKey)

			// find the next :
			posInObj, _ = build__find_next_token_pos(':', pos, tokensTableB)
			print_tokenB("detected COLON:", tokensTableB[posInObj])
			/*
			for posInObj = posInObj+1; posInObj<len(tokensTableB); posInObj++ {
				if tokensTableB[posInObj].tokenType == ':' {
					fmt.Println(": detected")
					break
				}
			}

			 */


			// find the next ANY token
			for posInObj = posInObj+1; posInObj<len(tokensTableB); posInObj++ {
				if tokensTableB[posInObj].tokenType == '"' ||
					tokensTableB[posInObj].tokenType == '{' {
					fmt.Println("next possible thing that we can handle, detected")
					break
				}
			}
			tokenNext := tokensTableB[posInObj]
			print_tokenB("tokenNext", tokenNext)

			////////////// VALUE HANDLING ////////////
			// handle embedded objects
			if tokenNext.tokenType == '{' {
				objEmbedded, errorsEmbedded, posEmbeddedLastUsed := JSON_B_structure_building(src, tokensTableB, posInObj+1)
				// todo error handling
				_ = errorsEmbedded
				elem.ValObject[objKey] = objEmbedded
				posInObj = posEmbeddedLastUsed
			}

			if tokenNext.tokenType == '"' {
				elem.ValObject[objKey] = NewString_JSON_value_B(
					getTextFromSrc(src, tokensTableB[posInObj]))
			}

			pos = posInObj

			if tokensTableB[posInObj].tokenType == '}' {
				break
			}
		} // handle objects
	} // for BIG loop

	return elem, errors, pos
}

// TODO: newArray, newObject, newInt, newFloat, newBool....
func NewString_JSON_value_B(text string) JSON_value_B {
	return JSON_value_B{
		ValType: '"',
		ValString: text, // now it is a raw string only
		// TODO: this string has to be interpreted,
	}
}

func NewObj_JSON_value_B() JSON_value_B {
	return JSON_value_B{
		ValType: '{',
		ValObject: map[string]JSON_value_B{},
	}
}

type JSON_value_B struct {
	ValType rune

	// ...... these values represent a Json elem's value - and one of them is filled only.. ..........
	ValObject map[string]JSON_value_B
	ValArray  []JSON_value_B

	ValBool bool // true, false

	ValString      string  // a string JSON value is stored here (filled ONLY if ValType is string)
	ValNumberInt   int     // an integer JSON value is stored here
	ValNumberFloat float64 // a float JSON value is saved here
}
