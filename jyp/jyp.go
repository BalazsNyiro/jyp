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
	"strconv"
)

var errorPrefix = "Error: "


// the type... constants are simple flags for tokens/JSON_values, to know what is stored in them
const typeIsUnknown byte = 0

const typeObjectOpen byte = 1
const typeObjectClose byte = 2
const typeObject byte = 3

const typeArrayOpen byte = 4
const typeArrayClose byte = 5
const typeArray byte = 6

const typeComma byte = 7
const typeColon byte = 8

const typeString byte = 9
const typeNull byte = 10
const typeBool byte = 11

const typeNumberInt byte = 12
const typeNumberFloat64 byte = 13
const typeNumber_exactTypeIsNotSet byte = 14


/*
const typeIsUnknown = "typeIsUnknown "

const typeObjectOpen = "typeObjectOpen"
const typeObjectClose = "typeObjectClose"
const typeObject = "typeObject"

const typeArrayOpen = "typeArrayOpen"
const typeArrayClose = "typeArrayClose"
const typeArray = "typeArray"

const typeComma = "typeComma"
const typeColon = "typeColon"

const typeString = "typeString "
const typeNull = "typeNull"
const typeBool = "typeBool"

const typeNumberInt = "typeNumberInt"
const typeNumberFloat64 = "typeNumberFloat64"
const typeNumber_exactTypeIsNotSet = "typeNumber_exactTypeIsNotSet"

 */

type JSON_oneObject map[string]JSON_value
type JSON_oneArray []JSON_value

// JSON value id -> Value pairs.
var global_valObjects =  map[int]JSON_oneObject{}
var global_valArrays = map[int]JSON_oneArray{}
var global_valBools = map[int] bool{}
var global_valStrings = map[int] string{}
var global_valNumInts = map[int] int{}
var global_valNumFloats = map[int] float64{}

// json elem: an elem, based on json.org definition
var global_JSON_ELEM_CONTAINER =  map[int]JSON_value{}

type JSON_value struct {
	Id     int   // unique id of the value
	ValType byte
	CharPositionFirstInSourceCode int // 0: the first char in source code, 1: 2nd...
	CharPositionLastInSourceCode  int // 0: the first char in source code, 1: 2nd...

	/////// CONNECTIONS ///
	IdParent int // the id of the parent elem - internal attrib to build structures

	LevelInObjectStructure int // later it helps to print the Json Value :-)
}
func (v JSON_value) ValString() string {
	return global_valStrings[v.Id]
}

func (v JSON_value) ValNumberInt() int {
	return global_valNumInts[v.Id]
}

func (v JSON_value) ValNumberFloat() float64{
	return global_valNumFloats[v.Id]
}

func (v JSON_value) ValArray() JSON_oneArray {
	return global_valArrays[v.Id]
}

func (v JSON_value) ValObject() JSON_oneObject {
	return global_valObjects[v.Id]
}

type tokenTable_startPositionIndexed_containerId map[int]int

func objectHierarchyBuilding(tokens tokenTable_startPositionIndexed_containerId, errorsCollected []error) JSON_value {
	var elemRoot JSON_value

	positionKeys_of_tokens := local_tool__tokenTable_position_keys_sorted(tokens)

	if len(positionKeys_of_tokens) < 1 {
		errorsCollected = append(errorsCollected, errors.New("emtpy source code, no tokens"))
		return elemRoot
	}

	keyFirst := positionKeys_of_tokens[0]
	if global_JSON_ELEM_CONTAINER[tokens[keyFirst]].ValType != typeObjectOpen {
		errorsCollected = append(errorsCollected, errors.New("the first token has to be 'objectOpen' in JSON source code"))
		return elemRoot
	}
	///////////////////////////////////////////////////////////////

	idParent := -1 // id can be 0 or bigger, so -1 is a non-existing parent id (root elem doesn't have parent
	keyStrings_of_collectors__filledIfObjectKey_emptyIfParentIsArray := map[int]string{} // if the parent is an object, elems can be inserted with keys.
	lastDetectedStringKey__inObject := ""

	// tokenKeys are charPosition based numbers, they are not continuous.
	for _, tokenPositionKey := range positionKeys_of_tokens {
		fmt.Println("token position key:", tokenPositionKey)
		tokenActualId := tokens[tokenPositionKey]

		if global_JSON_ELEM_CONTAINER[tokenActualId].ValType == typeComma { continue } // placeholders
		if global_JSON_ELEM_CONTAINER[tokenActualId].ValType == typeColon { continue }

		///////////////// SIMPLE JSON VALUES ARE SAVED DIRECTLY INTO THEIR PARENT /////////////////////////
		// in json objects, the first string is always the key. then the next elem is the value
		// detect this situation: string key in an object:
		if idParent >= 0 { // the first root elem doesn't have parents, so IdParent == -1
			if lastDetectedStringKey__inObject == "" {
				if global_JSON_ELEM_CONTAINER[idParent].ValType == typeObject {
					if global_JSON_ELEM_CONTAINER[tokenActualId].ValType == typeString {
						lastDetectedStringKey__inObject = global_valStrings[tokenActualId]
						continue
					} // == string
				} // == object
			} // == ""
		} // >= 0

		//////////////////////////////////////////////////////////////////////
		if global_JSON_ELEM_CONTAINER[tokenActualId].ValType == typeObjectOpen || global_JSON_ELEM_CONTAINER[tokenActualId].ValType == typeArrayOpen {
			// the id is important ONLY for the children - when they are inserted
			// into the parent containers. So when the container elem is parsed,
			// the id can be re-used later.
			id := len(global_JSON_ELEM_CONTAINER) // get the next free id in the database
			// container: array|object,
			levelInObjectStructure := 0
			if idParent >= 0 {
				levelInObjectStructure = global_JSON_ELEM_CONTAINER[idParent].LevelInObjectStructure + 1
			}

			containerNew := JSON_value{
				IdParent:                      idParent,
				Id:                            id,
				CharPositionFirstInSourceCode: global_JSON_ELEM_CONTAINER[tokenActualId].CharPositionFirstInSourceCode,
				LevelInObjectStructure:        levelInObjectStructure,
			}
			if global_JSON_ELEM_CONTAINER[tokenActualId].ValType == typeObjectOpen {
				containerNew.ValType = typeObject
			} else { // arrayOpen
				containerNew.ValType = typeArray
			}

			// the container has to be saved, because every new elem will be inserted
			// later, and at the close point, the whole container is handled as one
			global_JSON_ELEM_CONTAINER[id] = containerNew // single value.

			/* 	at this point, key has to be saved, because when the container is CLOSED,
			   	at that moment the key will be used, to insert the obj into te parent.
				if the container is an 'array' then keywords are not used, so they are empty
			*/
			keyStrings_of_collectors__filledIfObjectKey_emptyIfParentIsArray[id] = lastDetectedStringKey__inObject // save the key (it can be empty, or filled!)
			lastDetectedStringKey__inObject = ""
			idParent = id // this new array is the new parent for the next elems
			continue
		} // openers

		////////////////////////////////////////////////////////////////////////////////////////////////
		/////////////////// SIMPLE VALUES or container-closers ////////////////////////////////////////
		// 	"bool", "null", "string", "number_integer", "number_float64", "objectClose", "arrayClose"

		isCloserToken := global_JSON_ELEM_CONTAINER[tokenActualId].ValType == typeObjectClose || global_JSON_ELEM_CONTAINER[tokenActualId].ValType == typeArrayClose

		var value JSON_value

		if isCloserToken {
			if idParent == 0 { // closerToken IN root elem, which id==0
				elemRoot = global_JSON_ELEM_CONTAINER[0] //  read elem 0
				break                    // the exit point of the processing
			}
			// handle the container as a single value - and restore it's keyString.
			value = global_JSON_ELEM_CONTAINER[idParent] // the actual parent is closed, handle it as ONE value
			delete(global_JSON_ELEM_CONTAINER, idParent) // and remove the actual elemContainer from containers

			idParent = value.IdParent // the new parent after the obj close is the UPPER level parent
			lastDetectedStringKey__inObject = keyStrings_of_collectors__filledIfObjectKey_emptyIfParentIsArray[value.Id]
			delete(keyStrings_of_collectors__filledIfObjectKey_emptyIfParentIsArray, value.Id)
		} else {
			value = global_JSON_ELEM_CONTAINER[tokenActualId]
			value.IdParent = idParent
			value.LevelInObjectStructure = global_JSON_ELEM_CONTAINER[idParent].LevelInObjectStructure + 1
		} // notCloser

		///////////////// update the parent container with the new elem ////////////////////////////
		parent := global_JSON_ELEM_CONTAINER[idParent]
		parent.CharPositionLastInSourceCode = value.CharPositionLastInSourceCode
		// ^^^ the tokenCloser's last position is saved with this!
		if parent.ValType == typeObject { // objects are the typical containers, so this is the first in the checklist
			global_valObjects[parent.IdParent][lastDetectedStringKey__inObject] = value
			lastDetectedStringKey__inObject = "" // clear the keyName, we used that for the current object
		} else
		if parent.ValType == typeArray {
			global_valArrays[parent.IdParent] = append(global_valArrays[parent.IdParent], value)
		}
		global_JSON_ELEM_CONTAINER[idParent] = parent // save back the updated parent

	} // for, tokenNum, tokenPositionKey
	return elemRoot
}

// //////////////////// VALUE setter FUNCTIONS ///////////////////////////////////////////////
func valueValidationsSettings_inTokens(srcOrig []rune, tokens tokenTable_startPositionIndexed_containerId, errorsCollected []error) []error {
	// tokens is a MAP, and the keys are the token positions in src.
	// so there are gaps between keys!
	for _, tokenId := range tokens {

		if global_JSON_ELEM_CONTAINER[tokenId].ValType == typeString {
			valueValidateAndSetElemString(srcOrig, tokenId, errorsCollected)
		}
		/*
		if tokenOrig.ValType == typeNumber_exactTypeIsNotSet {
			tokenUpdated := valueValidateAndSetElemNumber(srcOrig, tokenOrig, errorsCollected)
			tokens[tokenStartPosInSrc] = tokenUpdated
		}

		 */
		// TODO: elem true|false|null value set?

	}
	return errorsCollected
}

// set the string value from raw strings
func valueValidateAndSetElemString(srcOrig []rune, tokenId int, errorsCollected []error) { // TESTED

	if global_JSON_ELEM_CONTAINER[tokenId].ValType != typeString {
		return
	} // don't modify non-string tokens

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
	for pos := global_JSON_ELEM_CONTAINER[tokenId].CharPositionFirstInSourceCode+1; pos <= global_JSON_ELEM_CONTAINER[tokenId].CharPositionLastInSourceCode-1; pos++ {

		runeActual := base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(srcOrig, pos)
		//fmt.Println("rune actual (string value set):", pos, string(runeActual), runeActual)
		runeNext1 := base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(srcOrig, pos+1)

		if runeActual != runeBackSlash { // a non-backSlash char
			valueFromRawSrcParsing = append(valueFromRawSrcParsing, runeActual)
			continue
		} else {
			// runeActual is \\ here, so ESCAPING started

			if runeNext1 == 'u' {
				// this is \u.... unicode code point - special situation,
				// because after the \u four other chars has to be handled

				runeNext2 := base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(srcOrig, pos+2)
				runeNext3 := base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(srcOrig, pos+3)
				runeNext4 := base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(srcOrig, pos+4)
				runeNext5 := base__srcGetChar__safeOverindexing__spaceGivenBackForAllWhitespaces(srcOrig, pos+5)

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

	// fmt.Println("value from raw src parsing:", string(valueFromRawSrcParsing))
	global_valStrings[tokenId] = string(valueFromRawSrcParsing)
}

func valueValidateAndSetElemNumber(srcOrig []rune, tokenNow JSON_value, errorsCollected []error) JSON_value{

	if tokenNow.ValType != typeNumber_exactTypeIsNotSet {
		return tokenNow
	} // don't modify non-number elems

	/*
			digits      mean: 0123456789
			digits19    mean:  123456789
			eE          mean: e|E
		    plusMinus   mean: +|-
			minus       mean: -
			fractionDot mean: .

			A number's sections:
			  v maybeMinusSign
		       vvvv part integer-digits
			       v fraction point
			        vvv  part fraction digits
		               v exponentEeLetter ----------\
			            v exponentPlusMinus ---------| exponent_section
		                 v exponentDigits ----------/
			  -1234.567e-8

			- maybeMinus: optional
			- fractionPoint: optional
			- exponentSection: optional
	*/

	// dividerBecauseOfFractionPoint := 0 // 10^0 = 1.
	// in case of 12.3: divider = 10^-1
	// in case of 1.23: divider = 10^-2

	numberRunes := base__runes_copy(srcOrig[tokenNow.CharPositionFirstInSourceCode:tokenNow.CharPositionLastInSourceCode+1])

	// example number: -1234.567e-8
	isNegative := numberRunes[0] == '-'
	if isNegative {
		numberRunes = numberRunes[1:]
	} // remove - sign, if that is the first

	runesSectionInteger := []rune{}  // and at the end, only integers remain...
	runesSectionFraction := []rune{} // filled second
	runesSectionExponent := []rune{} // filled first
	////////////////// the main sections of the number

	///////////////// the main markers of the number
	isFractionDotUsed := base__runeInRunes('.', numberRunes)
	isExponent_E_used := base__runeInRunes('E', numberRunes)
	isExponent_e_used := base__runeInRunes('e', numberRunes)

	/////////// go from back to forward: remove exponent part first
	if isExponent_e_used {
		numberRunes, runesSectionExponent = base__runes_split_at_pattern(numberRunes, 'e')
	}
	if isExponent_E_used {
		numberRunes, runesSectionExponent = base__runes_split_at_pattern(numberRunes, 'E')
	} /////////// if exponent is used, that is filled into the runeSectionExponentPart

	if isFractionDotUsed {
		numberRunes, runesSectionFraction = base__runes_split_at_pattern(numberRunes, '.')
		// after this, numberRunes lost the integer part.
	} // if FractionDot is used, split the Runes :-)

	runesSectionInteger = numberRunes
	///////////////////

	lenErrorCollectorBeforeErrorDetection := len(errorsCollected)

	/////////// ERROR HANDLING ////////////
	// if the first digit is 0, there cannot be more digits.

	if len(runesSectionInteger) > 1 {
		if runesSectionInteger[0] == 0 { // if integer part starts with 0, there cannot be other digits after initial 0
			errorsCollected = append(errorsCollected, errors.New("digits after 0 in integer part: "+string(runesSectionInteger)))
		}
	}

	var digits09 = []rune("0123456789")
	if !base__validate_runes_are_in_allowed_set(runesSectionInteger, digits09) {
		errorsCollected = append(errorsCollected, errors.New("illegal char in integer part: "+string(runesSectionInteger)))
	}

	if !base__validate_runes_are_in_allowed_set(runesSectionFraction, digits09) {
		errorsCollected = append(errorsCollected, errors.New("illegal char in fraction part: "+string(runesSectionFraction)))
	}

	if len(runesSectionExponent) > 0 { // validate the first char in exponent section
		if !base__validate_rune_are_in_allowed_set(runesSectionExponent[0], []rune{'+', '-'}) {
			errorsCollected = append(errorsCollected, errors.New("exponent part's first char is not +-: "+string(runesSectionExponent)))
		}
	}

	if len(runesSectionExponent) == 1 { // exponent section is too short
		errorsCollected = append(errorsCollected, errors.New("in exponent section +|- can be the FIRST char, then minimum one digit is necessary, and that is missing: "+string(runesSectionExponent)))
	}

	if len(runesSectionExponent) > 1 { // validate other chars in exponent section
		if !base__validate_runes_are_in_allowed_set(runesSectionExponent[1:], digits09) {
			errorsCollected = append(errorsCollected, errors.New("illegal char after first char of exponent section: "+string(runesSectionExponent)))
		}
	}

	/////////////////////////// NUM CALCULATION, BASED ON DIGITS ////////////////////////////////
	thisIsValidNumber := lenErrorCollectorBeforeErrorDetection == len(errorsCollected)
	if thisIsValidNumber {

		// cases: - only integer part,
		//        - int+fraction part,
		//        - int+exponent part
		//        - int+fraction+exponent part

		// ONLY INTEGER PART
		if len(runesSectionInteger) > 0 && len(runesSectionFraction) == 0 && len(runesSectionExponent) == 0 {
			numBase10, err := strconv.Atoi(string(runesSectionInteger))
			if err != nil {
				errorsCollected = append(errorsCollected, err)
			} else {
				if isNegative {
					numBase10 = -numBase10
				}
				global_valNumInts[tokenNow.Id] = numBase10
				tokenNow.ValType = typeNumberInt
			}
		}

		// ONLY INTEGER + FRACTION PART
		if len(runesSectionInteger) > 0 && len(runesSectionFraction) > 0 && len(runesSectionExponent) == 0 {
			numBase10, err := strconv.ParseFloat(string(runesSectionInteger)+"."+string(runesSectionFraction), 64)
			if err != nil {
				errorsCollected = append(errorsCollected, err)
			} else {
				if isNegative {
					numBase10 = -numBase10
				}
				global_valNumFloats[tokenNow.Id] = numBase10
				tokenNow.ValType = typeNumberFloat64
			}
		}

		// TODO: this section is too complicated, rewrite it.

		/*

			// if isNegative { multiplier = -1}

			integerValue := 0

			// calculate the exact value, based on elements of the number
			for _, r := range string(runesSectionInteger) + string(runesSectionFraction) {
				integerValue = integerValue * 10 // shift the value with one decimal place left
				integerValue += digitIntegerValue(r)
			}

			if runesSectionExponent[0] == '+' {
				for _, eDigit := range runesSectionExponent[1:] {
					eDigitVal := digitIntegerValue(eDigit)
					integerValue = integerValue * 10
					if eDigitVal > 0 {
						integerValue = integerValue * eDigitVal
					}
				}
			}

			////////////////////////////////////////////////////////////////////
			// then negative exponent and fraction points has to be handled, too
			divider := 0
			if runesSectionExponent[0] == '-' {
				for _, eDigit := range runesSectionExponent[1:] {
					eDigitVal := digitIntegerValue(eDigit)
					divider = divider * 10
					if eDigitVal > 0 {
						integerValue = integerValue * eDigitVal
					}
				}
			}

			if len(runesSectionFraction) > 0 {
				divider = divider - len(runesSectionFraction)  // divide the num with 10, 100, 1000...
			}


		*/

		// numberValue := multiplier * ()
	}
	return tokenNow
}

// ////////////////////  DETECTIONS  ///////////////////////////////////////////////
// Documented in program plan
func jsonDetect_strings______(src []rune, tokensStartPositions tokenTable_startPositionIndexed_containerId, errorsCollected []error) { // TESTED

	// to find escaped \" \\\" sections in strings
	isEscaped := false
	inStringDetection := false
	var tokenNow JSON_value

	for posInSrc, runeActual := range src {

		if runeActual == '"' {
			if !inStringDetection { // if at " char handling, we are NOT in string
					inStringDetection = true

					tokenNow = JSON_value{
						ValType: typeString,
						CharPositionFirstInSourceCode: posInSrc,
						CharPositionLastInSourceCode: posInSrc,
					}
					src[posInSrc] = ' '
					continue

			} else { // inStringDetection
				if !isEscaped { // a non-escaped " char in a string detection
					inStringDetection = false          // is the end of the string

					tokenNow.Id = len(global_JSON_ELEM_CONTAINER)
					tokenNow.CharPositionLastInSourceCode = posInSrc
					global_JSON_ELEM_CONTAINER[tokenNow.Id] = tokenNow
					tokensStartPositions[tokenNow.CharPositionFirstInSourceCode] = tokenNow.Id // save token
					global_valStrings[tokenNow.Id] = string(src[tokenNow.CharPositionFirstInSourceCode:tokenNow.CharPositionLastInSourceCode+1])
					src[posInSrc] = ' '
					continue
				}
			}
			// BE CAREFUL, there is a 3rd option!
			// if inStringDetection && isEscaped() -- which is handled as part of a string: inStringDetection
			// and this is different from the previous two, so don't change the structure!

		} // if " is detected, everything is handled in the conditions

		if inStringDetection {
			if runeActual == '\\' {
				isEscaped = ! isEscaped
			} else { // the escape series ended :-)
				isEscaped = false
			}
			src[posInSrc] = ' ' // remove detected char from src
		}

	} // for, runeActual
	if inStringDetection {
		errorsCollected = append(errorsCollected, errors.New("non-closed string detected:"))
	}
} // detect strings

func jsonDetect_separators___(src []rune, tokensStartPositions tokenTable_startPositionIndexed_containerId, errorsCollected []error) { // TESTED
	var tokenNow JSON_value

	detectedType := typeIsUnknown
	for posInSrc, runeActual := range src {

		if runeActual == ' ' {           // space is a very often filler after string detection.
			detectedType = typeIsUnknown // so the program doesn't have to check everything,
		} else                           // if a string is detected (other cases are not checked
		if runeActual == '{' {
			detectedType = typeObjectOpen
		} else
		if runeActual == '}' {
			detectedType = typeObjectClose
		} else
		if runeActual == '[' {
			detectedType = typeArrayOpen
		} else
		if runeActual == ']' {
			detectedType = typeArrayClose
		} else
		if runeActual == ',' {
			detectedType = typeComma
		} else
		if runeActual == ':' {
			detectedType = typeColon
		}

		if detectedType == typeIsUnknown {
			// keep the original rune, if it was not a detected char
		} else { // save token, if something important is detected
			tokenNow = JSON_value{Id: len(global_JSON_ELEM_CONTAINER), ValType: detectedType}
			tokenNow.CharPositionFirstInSourceCode = posInSrc
			tokenNow.CharPositionLastInSourceCode = posInSrc
			src[posInSrc] = ' '
			tokensStartPositions[tokenNow.CharPositionFirstInSourceCode] = tokenNow.Id
			global_JSON_ELEM_CONTAINER[tokenNow.Id] = tokenNow

			detectedType = typeIsUnknown
			// set back the type to unknown, if token is handled
		}
	} // for runeActual
} // separators....


/*
	 this detection is AFTER string+separator detection.
	   	in other words: only numbers and true/false/null values are left in the src.

		because the strings/separators are removed and replaced with space in the src, as placeholders,
	    the true/false/null words are surrounded with spaces, as separators.
*/
func jsonDetect_trueFalseNull(src []rune, tokensStartPositions tokenTable_startPositionIndexed_containerId, errorsCollected []error) { // TESTED

	// copy the original structure, not use the same variable
	// the detected word runesInSrc will be deleted from here.

	detectedType := typeIsUnknown // 3 types of word can be detected in this fun

	charsTrue := []rune("true")
	charsFalse := []rune("false")
	charsNull := []rune("null")
	for _, wordOne := range base__src_get_whitespace_separated_words_posFirst_posLast(src) {

		if base__compare_runes_are_equal(wordOne.wordChars, charsTrue) {
			detectedType = typeBool
		} else
		if base__compare_runes_are_equal(wordOne.wordChars, charsFalse) {
			detectedType = typeBool
		} else
		if base__compare_runes_are_equal(wordOne.wordChars, charsNull) {
			detectedType = typeNull
		}

		if detectedType != typeIsUnknown {
			tokenNow := JSON_value{Id: len(global_JSON_ELEM_CONTAINER), ValType: detectedType}
			if detectedType == typeBool {
				global_valBools[tokenNow.Id] = base__compare_runes_are_equal(wordOne.wordChars, charsTrue)
			}
			tokenNow.CharPositionFirstInSourceCode = wordOne.posFirst
			tokenNow.CharPositionLastInSourceCode = wordOne.posLast

			for posDetected := wordOne.posFirst; posDetected <= wordOne.posLast; posDetected++ {
				// clear detected positions from the src:
				src[posDetected] = ' '
				// only detected word runesInSrc are removed from the storage, where ALL original src is inserted in the first step
			}
			tokensStartPositions[tokenNow.CharPositionFirstInSourceCode] = tokenNow.Id
			global_JSON_ELEM_CONTAINER[tokenNow.Id] = tokenNow

			detectedType = typeIsUnknown // set the default value again ONLY if it was not unknown, in this case
		}
		/*	 it can be a number too, in case of else - so it is not an error, if typeIsUnknown */
	}
}

// words are detected here, and I can hope only that they are numbers - later they will be validated
func jsonDetect_numbers______(src []rune, tokensStartPositions tokenTable_startPositionIndexed_containerId, errorsCollected []error) { // TESTED

	for _, wordOne := range base__src_get_whitespace_separated_words_posFirst_posLast(src) {

		tokenNow := JSON_value{Id: len(global_JSON_ELEM_CONTAINER), ValType: typeNumber_exactTypeIsNotSet} // only numbers can be in the src now.
		tokenNow.CharPositionFirstInSourceCode = wordOne.posFirst
		tokenNow.CharPositionLastInSourceCode = wordOne.posLast
		tokensStartPositions[tokenNow.CharPositionFirstInSourceCode] = tokenNow.Id
		global_JSON_ELEM_CONTAINER[tokenNow.Id] = tokenNow
	}
}

/////////////////////////// local tools: supporter funcs, not for main logic /////////////////////////////

// tokenTable keys are character positions in JSON source code (positive integers)
func local_tool__tokenTable_position_keys_sorted(tokens tokenTable_startPositionIndexed_containerId) []int {
	keys := make([]int, 0, len(tokens))
	for k := range tokens {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func (v JSON_value) local_tool__updateLevelForChildren() {
	// when new elems are added dynamically, the level of the objects needs to be synchronised,
	// because the newly inserted object level can be determined ONLY when they are inserted

	if v.ValType == typeArray {
		arrayUpdated := []JSON_value{}
		for _, elem := range global_valArrays[v.Id] {
			elem.LevelInObjectStructure = v.LevelInObjectStructure + 1
			arrayUpdated = append(arrayUpdated, elem)
		}
		global_valArrays[v.Id] = arrayUpdated
	}
	if v.ValType == typeObject {
		objectsUpdated := map[string]JSON_value{}

		for key, elem := range global_valObjects[v.Id] {
			elem.LevelInObjectStructure = v.LevelInObjectStructure + 1
			objectsUpdated[key] = elem
		}
		global_valObjects[v.Id] = objectsUpdated
	}

}
