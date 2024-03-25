/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.


this file is the implementation of the _standard_ json data format:
https://www.json.org/json-en.html

these songs helped a lot to write this parser - respect:
  - Drummatix /туманами/
  - Mari Samuelsen /Sequence (four)/

in the code I intentionally avoid direct pointer usage - I think that is safer:
  - for goroutines
  - if json blocks are read and inserted into other json block, pointers are not useful,
    because they can have side effects. Every value has to be COPIED.


This module: is the main logic of json parsing



TODOS:
 - stepB, validation: {} [] "" pairings, missing commas, colons?
 - error handling, use errorsCollected everywhere
 - special number handling (e, hexa)
 - dedicated tests for all files against errors, for all functions
*/

package jyp

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

var errorPrefix = "Error: "

type tokenElem struct {
	tokenType rune /* one rune is stored here to represent a unit in the source code
                      { objOpen            123  (charCode)
                      } objClose           125
                      [ arrayOpen           91
                      ] arrayClose          93
                      , comma               44
                      : colon               58
                      " string              34

                      0 number              48    number is a general type, at this point we don't know is it an integer or a float or other
	                  I integer
                      F float64

                      b bool                98 // general collector type for true/false types
                      t true               116
                      f false              102

                      n null               110
                      ? not identified,     63
	                    only saved: later the type can be defined
	*/
	posInSrcFirst int
	posInSrcLast  int
}

func (t tokenElem) print(prefix string) {
	fmt.Println(prefix, string(t.tokenType), t.posInSrcFirst, t.posInSrcLast)
}

type tokenElems []tokenElem
func (tokens tokenElems) print() {
	for _, tokenElem := range tokens {
		fmt.Println(
			string(tokenElem.tokenType),
			fmt.Sprintf("%2d", tokenElem.posInSrcFirst),
			fmt.Sprintf("%2d", tokenElem.posInSrcLast),
		)
	}
}

// What if this is re-organised? JSON_value_B has only an Id,
// and every different type has a special storage?
// with that idea, 30ms can be saved, but the code becomes more complex
type JSON_value_B struct {
	ValType rune

	// ...... these values represent a Json elem's value - and one of them is filled only.. ..........
	ValObject map[string]JSON_value_B
	ValArray  []JSON_value_B

	ValBool bool // true, false

	ValString   string     // the parsed string. \n means 1 char here, for example
	ValNumberInt   int     // an integer JSON value is stored here
	ValNumberFloat float64 // a float JSON value is saved here
}

func (v JSON_value_B) ValObject_keys_sorted() []string{
	keys := make([]string, 0, len(v.ValObject))
	for k, _ := range v.ValObject {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}



func stepA__tokensTableDetect_structuralTokens_strings_L1(srcStr string) tokenElems {
	tokenTable := tokenElems{}
	posUnknownBlockStart := -1 // used only if the token is longer than 1 char. numbers, false/true for example
	
	//////////// TOKEN ADD //////////////////////////
	tokenAdd := func (typeOfToken rune, posFirst, posLast int) {
		if typeOfToken == '?' {
			textInSrc := srcStr[posFirst:posLast+1]
			if textInSrc == "true" { typeOfToken = 't' } else
			if textInSrc == "false" { typeOfToken = 'f' } else
			if textInSrc == "null" { typeOfToken = 'n' } else {
				typeOfToken = '0' // if eveything else is processed, only number can be the last choice
			}
		}

		tokenTable = append(tokenTable, tokenElem{tokenType: typeOfToken, posInSrcFirst: posFirst, posInSrcLast: posLast}  )
	} ////////// TOKEN ADD //////////////////////////

	inUnknownBlock := func () bool { return posUnknownBlockStart != -1	}
	
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
			if base__is_whitespace_rune(runeNow) {
				if ! inUnknownBlock() {
					// skip the whitespaces from tokens if the pos is NOT in unknown block,
					// so don't start an unknown block with a whitespace
				} else { // whitespace AFTER an unknown token
					// save the previously detected unknownBlock,
					// and skip the whitespace

					tokenAdd('?', posUnknownBlockStart, pos-1)
					posUnknownBlockStart = -1
				}
			} else if runeNow == '{' || runeNow == '}' || runeNow == '[' || runeNow == ']' || runeNow == ',' || runeNow == ':' {
				if inUnknownBlock() {
					tokenAdd('?', posUnknownBlockStart, pos-1)
					posUnknownBlockStart = -1
				}
				tokenAdd(runeNow, pos, pos)
			} else {
				// not in string, and not json structural char and not whitespace
				// so it can be a number, true/false/null or
				// whitespaces:
				if ! inUnknownBlock() {
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

func stepB__JSON_B_validation_L1(tokenTableB tokenElems) []error {
	// TODO: loop over the table, find incorrect {..} [..], ".." pairs,
	// incorrect numbers, everything that can be a problem
	// so after this, we can be sure that every elem are in pairs.
	return []error{}
}


// L1: Level 1. A higher level is a more general fun, a lower level is a tool, lib func, or something small
func stepC__JSON_B_structure_building__L1(src string, tokensTableB tokenElems, tokenPosStart int, errorsCollected []error) (JSON_value_B, int) {
	if tokenPosStart >= len(tokensTableB) {
		errorsCollected= append(errorsCollected, errors.New("wanted position index is higher than tokensTableB"))
	}

	if len(errorsCollected) > 0 {
		return JSON_value_B{}, 0
	}
	elem := JSON_value_B{}
	var pos int

	for pos = tokenPosStart; pos<len(tokensTableB); pos++ {
		tokenNow := tokensTableB[pos]

		if tokenNow.tokenType == '"' {
			elem = NewString_JSON_value_quotedBothEnd(base__read_sourceCode_section_basedOnTokenPositions(src, tokensTableB[pos], false), errorsCollected)
			break

		} else if tokenNow.tokenType == '0' { // general number detection
			textInSrc := base__read_sourceCode_section_basedOnTokenPositions(src, tokensTableB[pos], false)
			// detect simple integers first
			i, err := strconv.Atoi(textInSrc)  // it can't interpret Scientific nums!
			if err == nil { // it was really an integer...
				elem = NewNumInt(i)
				break
			}

			f, err := strconv.ParseFloat(textInSrc, 64)
			if err == nil { // it was really an integer...
				elem = NewNumFloat(f)
				break
			}

			// worst case, I don't know what is this, so insert it as a string
			elem = NewString_JSON_value_quotedBothEnd(base__read_sourceCode_section_basedOnTokenPositions(src, tokensTableB[pos], false), errorsCollected)
			break

		} else if tokenNow.tokenType == '{' {
				elem = NewObj_JSON_value_B()

				for ; pos <len(tokensTableB); { // detect children
					pos, _ = token_find_next__L2(true, []rune{'"'}, pos+1, tokensTableB)

					// the next string key, the objKey is not quoted, but interpreted, too
					objKey := stringValueParsing_rawToInterpretedCharacters_L2(base__read_sourceCode_section_basedOnTokenPositions(src, tokensTableB[pos], true), errorsCollected)

					// find the next : but don't do anything with that
					pos, _ = token_find_next__L2(true, []rune{':'}, pos+1, tokensTableB)

					// find the next ANY token, the new VALUE
					nextValueElem, posLastUsed := stepC__JSON_B_structure_building__L1(src, tokensTableB, pos+1, errorsCollected)
					elem.ValObject[objKey] = nextValueElem
					pos = posLastUsed

					if pos+1 < len(tokensTableB) { // look forward:
						if tokensTableB[pos+1].tokenType == '}' {
							break
						}
					}
					pos, _ = token_find_next__L2(true, []rune{','}, pos+1, tokensTableB)
				} // for pos, internal children loop

		} else if tokenNow.tokenType == '[' {
			elem = NewArr_JSON_value_B()
			for ; pos < len(tokensTableB);  { // detect children
				// find the next ANY token, the new VALUE
				nextValueElem, posLastUsed := stepC__JSON_B_structure_building__L1(src, tokensTableB, pos+1, errorsCollected)
				elem.ValArray = append(elem.ValArray, nextValueElem)
				pos = posLastUsed

				if pos+1 < len(tokensTableB) { // look forward:
					if tokensTableB[pos+1].tokenType == ']' {  // if the next elem is ], this is the last child,
						break // and stop the children detection, leave the detection for loop
					}
				}
				pos, _ = token_find_next__L2(true, []rune{','}, pos+1, tokensTableB)
			} // for pos, internal children loop

		} else if tokenNow.tokenType == '}' { break   // ascii:125,
		} else if tokenNow.tokenType == ']' { break   // elem prepared, exit

		} else if tokenNow.tokenType == 't' {
			elem = NewBool(true)
			break

		} else if tokenNow.tokenType == 'f' {
			elem = NewBool(false)
			break

		} else if tokenNow.tokenType == 'n' {
			elem = NewNull()
			break
		}

		/*	This is not possible anymore, every possible token type is detected now
			} else if tokenNow.tokenType == '?' {
				elem = NewString_JSON_value_quotedBothEnd("\"unknown_elem, maybe number or bool\"", errorsCollected)
				break
		*/
	} // for BIG loop

	return elem, pos // ret with last used position
}



// return with pos only to avoid elem copy with reading/passing
// find the next token from allowed types
// one token, or more than one token can be searched
func token_find_next__L2(wantSomethingFromTypes bool, types []rune, posActual int, tokensTable tokenElems) (int, error) {
	var pos int
	if pos >= len(tokensTable) {
		return pos, errors.New("token position is bigger than last elem index")
	}
	if wantSomethingFromTypes { // want one from types:
		for pos = posActual; pos<len(tokensTable); pos++ {
			for _, wanted := range types {
				if tokensTable[pos].tokenType == wanted {
					// base__print_tokenB("wanted1:", tokensTable[pos])
					return pos, nil
				}
			}
		}
		return pos, errors.New("wanted token is not detected in table")

	} else { // want something that is NOT in typeList
		for pos = posActual; pos<len(tokensTable); pos++ {
			actualTypeIsNonWanted := false
			for _, nonWantedType := range types {
				if tokensTable[pos].tokenType == nonWantedType {
					actualTypeIsNonWanted = true// all nonWanted has to be checked
					break
				}
			}
			if ! actualTypeIsNonWanted {
				// base__print_tokenB("wanted2:", tokensTable[pos])
				return pos, nil
			}
		}
		return pos, errors.New("wanted token is not detected in table")
	}
}


// set the string value from raw strings
// in orig soure code, \n means 2 chars: a backslash and 'n'.
// but if it is interpreted, that is one newline "\n" char.
func stringValueParsing_rawToInterpretedCharacters_L2(src string, errorsCollected []error) string{ // TESTED

	/* Tasks:
	- is it a valid string?
	- convert special char representations to real chars

	the func works typically with 2 chars, for example: \t
	but sometime with 6: \u0123, so I need to look forward for the next 5 chars
	*/

	valueFromRawSrcParsing := []rune{}

	// fmt.Println("string token value detection:", src)
	runeBackSlash := '\\' // be careful: this is ONE \ char, only written with this expression

	// pos start + 1: strings has initial " in runes
	// post end -1  closing " after string content
	for pos := 0; pos < len(src); pos++ {

		runeActual := base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(src, pos)
		//fmt.Println("rune actual (string value set):", pos, string(runeActual), runeActual)

		if runeActual != runeBackSlash { // a non-backSlash char
			valueFromRawSrcParsing = append(valueFromRawSrcParsing, runeActual)
			continue
		} else {
			// runeActual is \\ here, so ESCAPING started
			runeNext1 := base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(src, pos+1)

			if runeNext1 == 'u' {
				// this is \u.... unicode code point - special situation,
				// because after the \u four other chars has to be handled

				runeNext2 := base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(src, pos+2)
				runeNext3 := base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(src, pos+3)
				runeNext4 := base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(src, pos+4)
				runeNext5 := base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(src, pos+5)

				base10_val_2, err2 := base__hexaRune_to_intVal(runeNext2)
				if err2 != nil {
					errorsCollected = append(errorsCollected, err2)
				}

				base10_val_3, err3 := base__hexaRune_to_intVal(runeNext3)
				if err3 != nil {
					errorsCollected = append(errorsCollected, err3)
				}

				base10_val_4, err4 := base__hexaRune_to_intVal(runeNext4)
				if err4 != nil {
					errorsCollected = append(errorsCollected, err4)
				}

				base10_val_5, err5 := base__hexaRune_to_intVal(runeNext5)
				if err5 != nil {
					errorsCollected = append(errorsCollected, err5)
				}

				unicodeVal_10Based := 0

				if err2 == nil && err3 == nil && err4 == nil && err5 == nil {
					unicodeVal_10Based = base10_val_2*16*16*16 +
						base10_val_3*16*16 +
						base10_val_4*16 +
						base10_val_5
				}
				runeFromHexaDigits := rune(unicodeVal_10Based)

				pos += 1 + 4 // one extra pos because of the u, and +4 because of the digits
				valueFromRawSrcParsing = append(valueFromRawSrcParsing, runeFromHexaDigits)

			} else { // the first detected char was a backslash, what is the second?
				// so this is a simple escaped char, for example: \" \t \b \n
				runeReal := '?'
				if runeNext1 == '"' { // \" -> is a " char in a string
					runeReal = '"' // in a string, this is an escaped " double quote char
				} else
				if runeNext1 == runeBackSlash { // in reality, these are the 2 chars: \\
					runeReal = '\\' // reverse solidus
				} else
				if runeNext1 == '/' { // a very special escaping: \/
					runeReal = '/' // solidus
				} else
				if runeNext1 == 'b' { // This is the first good example for escaping:
					runeReal = '\b' // in the src there were 2 chars: \ and b,
				} else //  (backspace)    // and one char is inserted into the stringVal
				if runeNext1 == 'f' { // formfeed
					runeReal = '\f'
				} else
				if runeNext1 == 'n' { // linefeed
					runeReal = '\n'
				} else
				if runeNext1 == 'r' { // carriage return
					runeReal = '\r' //
				} else
				if runeNext1 == 't' { // horizontal tab
					runeReal = '\t' //
				}

				pos += 1 // one extra pos increasing is necessary, because of
				// 2 chars were processed: the actual \ and the next one.

				valueFromRawSrcParsing = append(valueFromRawSrcParsing, runeReal)
			}
		} // else
	} // for

	return string(valueFromRawSrcParsing)
}

