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
*/

package jyp

import (
	"errors"
	"fmt"
	"sort"
	"unicode"
)

var errorPrefix = "Error: "

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

func tokensTableDetect_structuralTokens_strings(srcStr string) tokenElems_B {
	tokenTable := tokenElems_B{}
	posUnknownBlockStart := -1 // used only if the token is longer than 1 char. numbers, false/true for example
	
	//////////// TOKEN ADD func ///////////////////////
	tokenAdd := func (typeOfToken rune, posFirst, posLast int) {
		// TODO: unknown token processing here, to avoid second loop? w
		tokenTable = append(tokenTable, tokenElem_B{tokenType: typeOfToken, posInSrcFirst: posFirst, posInSrcLast: posLast}  )
	} // func, tokenAdd //////////////////////////////

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
					tokenTable = append(tokenTable, tokenElem_B{tokenType: '?', posInSrcFirst: posUnknownBlockStart, posInSrcLast: pos-1}  )
					posUnknownBlockStart = -1
				}
			} else if runeNow == '{' || runeNow == '}' || runeNow == '[' || runeNow == ']' || runeNow == ',' || runeNow == ':' {
				if inUnknownBlock() {
					tokenTable = append(tokenTable, tokenElem_B{tokenType: '?', posInSrcFirst: posUnknownBlockStart, posInSrcLast: pos-1}  )
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

func base__is_whitespace_rune(oneRune rune) bool { // TESTED
	/*  https://stackoverflow.com/questions/29038314/determining-whitespace-in-go
		'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).
		Other definitions of spacing characters are set by category Z and property Pattern_White_Space. */
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


func getTextFromSrc(src string, token tokenElem_B, quoted bool) string {
	if quoted {
		return src[token.posInSrcFirst:token.posInSrcLast+1]
	}
	return src[token.posInSrcFirst+1:token.posInSrcLast]
}

func print_tokenB(prefix string, t tokenElem_B) {
	fmt.Println(prefix, string(t.tokenType), t.posInSrcFirst, t.posInSrcLast)
}


// return with pos only to avoid elem copy with reading/passing
// find the next token from allowed types
func build__find_next_token_pos(wantedThem bool, types []rune, posActual int, tokensTable tokenElems_B) (int, error) {
	var pos int
	if pos >= len(tokensTable) {
		return pos, errors.New("token position is bigger than last elem index")
	}
	if wantedThem { // want one from types:
		for pos = posActual; pos<len(tokensTable); pos++ {
			for _, wanted := range types {
				if tokensTable[pos].tokenType == wanted {
					// print_tokenB("wanted1:", tokensTable[pos])
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
				// print_tokenB("wanted2:", tokensTable[pos])
				return pos, nil
			}
		}
		return pos, errors.New("wanted token is not detected in table")
	}
}

// repeat the wanted unit prefix a few times
func prefixGen(oneUnitPrefix string, repeatNum int) string {
	if oneUnitPrefix == "" {
		return "" // if there is nothing to repeat
	}
	out := ""
	for i:=0; i<repeatNum; i++ {
		out += oneUnitPrefix
	}
	return out
}

func JSON_B_structure_building(src string, tokensTableB tokenElems_B, tokenPosStart int, errorsCollected []error) (JSON_value_B, []error, int) {
	problems := JSON_B_validation(tokensTableB)
	if tokenPosStart >= len(tokensTableB) {
		problems = append(problems, errors.New("wanted position index is higher than tokensTableB"))
	}
	if len(problems) > 0 {
		return JSON_value_B{}, problems, 0
	}
	elem := JSON_value_B{}
	var pos int

	for pos = tokenPosStart; pos<len(tokensTableB); pos++ {
		tokenNow := tokensTableB[pos]

		if tokenNow.tokenType == '{' {
			elem = NewObj_JSON_value_B()

			for ; pos <len(tokensTableB); pos++ { // detect children
				// todo: error handling
				pos, _ = build__find_next_token_pos(true, []rune{'"'}, pos, tokensTableB)

				// the next string key, the objKey is not quoted, but interpreted, too
				objKey := stringValueParsing_rawToInterpretedCharacters(getTextFromSrc(src, tokensTableB[pos], false), errorsCollected)

				// find the next : but don't do anything with that
				pos, _ = build__find_next_token_pos(true, []rune{':'}, pos+1, tokensTableB)

				// find the next ANY token, the new VALUE
				nextValueElem, _, posLastUsed := JSON_B_structure_building(src, tokensTableB, pos+1, errorsCollected)
				elem.ValObject[objKey] = nextValueElem
				pos = posLastUsed

				if pos+1 < len(tokensTableB) { // look forward:
					if tokensTableB[pos+1].tokenType == '}' {
						break
					}
				}
				pos, _ = build__find_next_token_pos(true, []rune{','}, pos+1, tokensTableB)
			} // for pos, internal children loop

		} else if tokenNow.tokenType == '?' {
			elem = NewString_JSON_value_quotedBothEnd("\"unknown_elem, maybe number or bool\"", errorsCollected)
			break

		} else if tokenNow.tokenType == '"' {
			elem = NewString_JSON_value_quotedBothEnd(getTextFromSrc(src, tokensTableB[pos], true), errorsCollected)
			break

		} else if tokenNow.tokenType == '[' {
			elem = NewArr_JSON_value_B()
			for ; pos < len(tokensTableB); pos++ { // detect children
				// find the next ANY token, the new VALUE
				nextValueElem, _, posLastUsed := JSON_B_structure_building(src, tokensTableB, pos+1, errorsCollected)
				elem.ValArray = append(elem.ValArray, nextValueElem)
				pos = posLastUsed

				if pos+1 < len(tokensTableB) { // look forward:
					if tokensTableB[pos+1].tokenType == ']' {
						break
					}
				}
				pos, _ = build__find_next_token_pos(true, []rune{','}, pos+1, tokensTableB)
			} // for pos, internal children loop
		} else if tokenNow.tokenType == '}' { break   // ascii:125,
		} else if tokenNow.tokenType == ']' { break } // elem prepared, exit
	} // for BIG loop

	return elem, problems, pos // ret with last used position
}

// TODO: newObject, newInt, newFloat, newBool....
func NewString_JSON_value_quotedBothEnd(text string, errorsCollected []error) JSON_value_B {
	// strictly have minimum one "opening....and...one..closing" quote!
	valString := stringValueParsing_rawToInterpretedCharacters( text[1:len(text)-1], errorsCollected)

	return JSON_value_B{
		ValType:      '"',
		ValStringRaw: text,
		ValString: valString,
	}
}

func NewObj_JSON_value_B() JSON_value_B {
	return JSON_value_B{
		ValType: '{',
		ValObject: map[string]JSON_value_B{},
	}
}

func NewArr_JSON_value_B() JSON_value_B {
	return JSON_value_B{
		ValType: '[',
		ValArray: []JSON_value_B{},
	}
}


type JSON_value_B struct {
	ValType rune

	// ...... these values represent a Json elem's value - and one of them is filled only.. ..........
	ValObject map[string]JSON_value_B
	ValArray  []JSON_value_B

	ValBool bool // true, false

	ValString   string     // the parsed string. \n means 1 char here, for example
	ValStringRaw   string  // exact copy of the original source code, that was parsed (non-interpreted. \t means 2 chars, not 1!
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


func base__separator_set_if_no_last_elem(position, length_numOfAllElems int, separator string) string {
	if position < length_numOfAllElems-1 {
		return separator
	}
	return ""
}

// set the string value from raw strings
// in orig soure code, \n means 2 chars: a backslash and 'n'.
// but if it is interpreted, that is one newline "\n" char.
func stringValueParsing_rawToInterpretedCharacters(src string, errorsCollected []error) string{ // TESTED

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
		runeNext1 := base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(src, pos+1)

		if runeActual != runeBackSlash { // a non-backSlash char
			valueFromRawSrcParsing = append(valueFromRawSrcParsing, runeActual)
			continue
		} else {
			// runeActual is \\ here, so ESCAPING started

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

	// errorsCollected is updated by default, it can be modified here in this func
	return string(valueFromRawSrcParsing)
}


// get the rune IF the index is really in the range of the src.
// return with ' ' space, IF the index is NOT in the range.
// reason: avoid never ending index checking, so do it only once
// the space can be answered because this func is used when a real char wanted to be detected,
// and if a space is returned, this has NO MEANING in that parse section
// this fun is NOT used in string detection - and other places whitespaces can be neglected, too
// getChar, with whitespace replace
func base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(src string, pos int) rune { // TESTED
	char := base__srcGetChar__safeOverindexing(src, pos)
	if base__is_whitespace_rune(char) {
		return ' ' // simplify everything. if the char is ANY whitespace char,
		// return with SPACE, this is not important in the source code parsing
	}
	return char
}

// getChar, no whitespace replace
func base__srcGetChar__safeOverindexing(src string, pos int) rune { // TESTED
	posPossibleMax := len(src) - 1  // if src is empty, max is -1,
	posPossibleMin := 0             // and the condition cannot be true here:
	if (pos >= posPossibleMin) && (pos <= posPossibleMax) {
		return []rune(src[pos:pos+1])[0]
	}
	return ' '
}

// runesSections were checked against illegal chars, so here digitRune is in 0123456789
// TODO: maybe can be removed if not used in the future in exponent number detection section
func base__digit10BasedRune_integer_value(digit10based rune) (int, error) {
	if digit10based == '0' {
		return 0, nil
	}
	if digit10based == '1' {
		return 1, nil
	}
	if digit10based == '2' {
		return 2, nil
	}
	if digit10based == '3' {
		return 3, nil
	}
	if digit10based == '4' {
		return 4, nil
	}
	if digit10based == '5' {
		return 5, nil
	}
	if digit10based == '6' {
		return 6, nil
	}
	if digit10based == '7' {
		return 7, nil
	}
	if digit10based == '8' {
		return 8, nil
	}
	if digit10based == '9' {
		return 9, nil
	}
	return 0, errors.New(errorPrefix + "rune (" + string(digit10based) + ")")
}

func base__hexaRune_to_intVal(hexaChar rune) (int, error) { // TESTED
	hexaTable := map[rune]int{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'a': 10,
		'b': 11,
		'c': 12,
		'd': 13,
		'e': 14,
		'f': 15,
	}
	base10Val, keyInHexaTable := hexaTable[hexaChar]
	if keyInHexaTable {
		return base10Val, nil
	}
	return 0, errors.New("hexa char(" + string(hexaChar) + ") was not in hexa table")
}

