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
			if base__is_whitespace_rune(runeNow) {
				if posUnknownBlockStart == -1 {
					// skip the whitespaces from tokens, don't do anything
				} else { // whitespace AFTER an unknown token
					tokenTable = append(tokenTable, tokenElem_B{tokenType: '?', posInSrcFirst: posUnknownBlockStart, posInSrcLast: pos-1}  )
					posUnknownBlockStart = -1
				}
			} else if runeNow == '{' || runeNow == '}' || runeNow == '[' || runeNow == ']' || runeNow == ',' || runeNow == ':' {
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


func getTextFromSrc(src string, token tokenElem_B) string {
	return src[token.posInSrcFirst:token.posInSrcLast+1]
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

func JSON_B_structure_building(src string, tokensTableB tokenElems_B, tokenPosStart int) (JSON_value_B, []error, int) {
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
				objKey := getTextFromSrc(src, tokensTableB[pos]) // next string key

				// find the next : but don't do anything with that
				pos, _ = build__find_next_token_pos(true, []rune{':'}, pos+1, tokensTableB)

				// find the next ANY token, the new VALUE
				nextValueElem, _, posLastUsed := JSON_B_structure_building(src, tokensTableB, pos+1)
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
			elem = NewString_JSON_value_quotedBothEnd("\"unknown_elem, maybe number or bool\"")
			break

		} else if tokenNow.tokenType == '"' {
			elem = NewString_JSON_value_quotedBothEnd(getTextFromSrc(src, tokensTableB[pos]))
			break

		} else if tokenNow.tokenType == '[' {
			elem = NewArr_JSON_value_B()
			for ; pos < len(tokensTableB); pos++ { // detect children
				// find the next ANY token, the new VALUE
				nextValueElem, _, posLastUsed := JSON_B_structure_building(src, tokensTableB, pos+1)
				elem.ValArray = append(elem.ValArray, nextValueElem)
				pos = posLastUsed

				if pos+1 < len(tokensTableB) { // look forward:
					if tokensTableB[pos+1].tokenType == ']' {
						break
					}
				}
				pos, _ = build__find_next_token_pos(true, []rune{','}, pos+1, tokensTableB)
			} // for pos, internal children loop
		} else if tokenNow.tokenType == '}' { break
		} else if tokenNow.tokenType == ']' { break } // elem prepared, exit
	} // for BIG loop

	return elem, problems, pos // ret with last used position
}

// TODO: newObject, newInt, newFloat, newBool....
func NewString_JSON_value_quotedBothEnd(text string) JSON_value_B {
	// strictly have minimum one "opening....and...one..closing" quote!
	return JSON_value_B{
		ValType:      '"',
		ValStringRaw: text,
		ValString: text[1:len(text)-1],
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


// example usage: repr(2) means: use "  " 2 spaces as indentation
// if ind
// otherwise a simple formatted output
func (v JSON_value_B) repr(indentationLength ...int) string {
	if len(indentationLength) == 0 { // so no param is passed
		return v.repr_fine("", 0)
	}
	if indentationLength[0] < 1 { // a 0 or negative param is passed
		return v.repr_fine("", 0)
	}
	indentation := prefixGen(" ", indentationLength[0])
	return v.repr_fine(indentation, 0)
}

// tunable repr
func (v JSON_value_B) repr_fine(indent string, level int) string {
	prefix := "" // indentOneLevelPrefix
	prefix2 := "" // indentTwoLevelPrefix
	newLine := ""
	if len(indent) > 0 {
		prefix = prefixGen(indent, level)
		prefix2 = prefixGen(indent, level+1)
		newLine = "\n"
	}

	if v.ValType == '"' {
		return "\""+v.ValString + "\""
	}

	if v.ValType == '{' {
		out := prefix + "{" + newLine
		for counter, childKey := range v.ValObject_keys_sorted() {
			comma := base__separator_set_if_no_last_elem(counter, len(v.ValObject), ",")
			childVal := v.ValObject[childKey]
			out += prefix2 + childKey + ":" + " " + childVal.repr_fine(indent, level+1) + comma + newLine
		}
		out += prefix + "}" + newLine
		return out
	}
	// TODO: comma after values
	if v.ValType == '[' {
		out := prefix + "[" + newLine
		for counter, child := range v.ValArray {
			comma := base__separator_set_if_no_last_elem(counter, len(v.ValArray), ",")
			out += prefix2 + indent + child.repr_fine(indent, level+1) + comma + newLine
		}
		out += prefix + "]" + newLine
		return out
	}
	return ""
}

func base__separator_set_if_no_last_elem(position, length_numOfAllElems int, separator string) string {
	if position < length_numOfAllElems-1 {
		return separator
	}
	return ""
}
