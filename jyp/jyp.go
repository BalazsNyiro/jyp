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
	"sort"
	"strconv"
	"strings"
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

// token: a structural elem of the source code, maybe without real meaning (comma separator, for example)
type token struct {
	valType byte

	valBool        bool
	valString      string
	valNumberInt   int
	valNumberFloat float64

	charPositionFirstInSourceCode int // 0: the first char in source code, 1: 2nd...
	charPositionLastInSourceCode  int // 0: the first char in source code, 1: 2nd...
	runes                         []rune
}

type JSON_value struct {
	ValType byte // the type of the value

	// ...............................................................................................
	// ...... these values represent a Json elem's value - and one of them is filled only.. ..........
	ValObject map[string]JSON_value
	ValArray  []JSON_value

	ValBool bool // true, false

	ValString      string  // a string JSON value is stored here (filled ONLY if ValType is string)
	ValNumberInt   int     // an integer JSON value is stored here
	ValNumberFloat float64 // a float JSON value is saved here
	// ...............................................................................................

	//////// PARSING SECTION: detection from the JSON source code /////
	CharPositionFirstInSourceCode int // 0: the first char in source code, 1: 2nd...
	CharPositionLastInSourceCode  int // 0: the first char in source code, 1: 2nd...
	Runes                         []rune
	AddedInGoCode                 bool

	/////// CONNECTIONS ///
	idParent int // the id of the parent elem - internal attrib to build structures
	idSelf   int // the id of this actual elem

	LevelInObjectStructure int // later it helps to print the Json Value :-)
}

type tokenTable_startPositionIndexed map[int]token

func objectHierarchyBuilding(tokens tokenTable_startPositionIndexed, errorsCollected []error) (JSON_value, []error) {
	var elemRoot JSON_value

	positionKeys_of_tokens := local_tool__tokenTable_position_keys_sorted(tokens)

	if len(positionKeys_of_tokens) < 1 {
		errorsCollected = append(errorsCollected, errors.New("emtpy source code, no tokens"))
		return elemRoot, errorsCollected
	}

	keyFirst := positionKeys_of_tokens[0]
	if tokens[keyFirst].valType != typeObjectOpen {
		errorsCollected = append(errorsCollected, errors.New("the first token has to be 'objectOpen' in JSON source code"))
		return elemRoot, errorsCollected
	}
	///////////////////////////////////////////////////////////////

	idParent := -1 // id can be 0 or bigger, so -1 is a non-existing parent id (root elem doesn't have parent
	containers := map[int]JSON_value{}
	keyStrings_of_collectors__filledIfObjectKey_emptyIfParentIsArray := map[int]string{} // if the parent is an object, elems can be inserted with keys.
	lastDetectedStringKey__inObject := ""

	// tokenKeys are charPosition based numbers, they are not continuous.
	for _, tokenPositionKey := range positionKeys_of_tokens {

		tokenActual := tokens[tokenPositionKey]

		if tokenActual.valType == typeComma { continue } // placeholders
		if tokenActual.valType == typeColon { continue }

		///////////////// SIMPLE JSON VALUES ARE SAVED DIRECTLY INTO THEIR PARENT /////////////////////////
		// in json objects, the first string is always the key. then the next elem is the value
		// detect this situation: string key in an object:
		if idParent >= 0 { // the first root elem doesn't have parents, so idParent == -1
			if lastDetectedStringKey__inObject == "" {
				if containers[idParent].ValType == typeObject {
					if tokenActual.valType == typeString {
						lastDetectedStringKey__inObject = tokenActual.valString
						continue
					} // == string
				} // == object
			} // == ""
		} // >= 0

		//////////////////////////////////////////////////////////////////////
		if tokenActual.valType == typeObjectOpen || tokenActual.valType == typeArrayOpen {
			// the id is important ONLY for the children - when they are inserted
			// into the parent containers. So when the container elem is parsed,
			// the id can be re-used later.
			id := len(containers) // get the next free id in the database
			// container: array|object,
			levelInObjectStructure := 0
			if idParent >= 0 {
				levelInObjectStructure = containers[idParent].LevelInObjectStructure + 1
			}

			containerNew := JSON_value{
				idParent:                      idParent,
				idSelf:                        id,
				CharPositionFirstInSourceCode: tokenActual.charPositionFirstInSourceCode,
				LevelInObjectStructure:        levelInObjectStructure,
			}
			if tokenActual.valType == typeObjectOpen {
				containerNew.ValType = typeObject
				containerNew.ValObject = map[string]JSON_value{}
			} else { // arrayOpen
				containerNew.ValType = typeArray
				containerNew.ValArray = []JSON_value{}
			}

			// the container has to be saved, because every new elem will be inserted
			// later, and at the close point, the whole container is handled as one
			containers[id] = containerNew // single value.

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

		isCloserToken := tokenActual.valType == typeObjectClose || tokenActual.valType == typeArrayClose

		var value JSON_value

		if isCloserToken {
			if idParent == 0 { // closerToken IN root elem, which id==0
				elemRoot = containers[0] //  read elem 0
				break                    // the exit point of the processing
			}
			// handle the container as a single value - and restore it's keyString.
			value = containers[idParent] // the actual parent is closed, handle it as ONE value
			delete(containers, idParent) // and remove the actual elemContainer from containers

			idParent = value.idParent // the new parent after the obj close is the UPPER level parent
			lastDetectedStringKey__inObject = keyStrings_of_collectors__filledIfObjectKey_emptyIfParentIsArray[value.idSelf]
			delete(keyStrings_of_collectors__filledIfObjectKey_emptyIfParentIsArray, value.idSelf)
		} else {
			value = JSON_value{ValType: tokenActual.valType,
				CharPositionFirstInSourceCode: tokenActual.charPositionFirstInSourceCode,
				CharPositionLastInSourceCode:  tokenActual.charPositionLastInSourceCode,
				Runes:                         tokenActual.runes,
				idParent:                      idParent,
				LevelInObjectStructure:        containers[idParent].LevelInObjectStructure + 1,
				// idSelf is not filled, the whole id conception is used ONLY to find parents
				// during parsing.
			}

			//if tokenActual.valType == typeNull {
			//	_ = "null type has no value, don't store it"
			// }

			if tokenActual.valType == typeBool {
				value.ValBool = tokenActual.valBool
			}
			if tokenActual.valType == typeString {
				value.ValString = tokenActual.valString
			}
			if tokenActual.valType == typeNumberInt {
				value.ValNumberInt = tokenActual.valNumberInt
			}
			if tokenActual.valType == typeNumberFloat64 {
				value.ValNumberFloat = tokenActual.valNumberFloat
			}

		} // notCloser

		///////////////// update the parent container with the new elem ////////////////////////////
		parent := containers[idParent]
		parent.CharPositionLastInSourceCode = value.CharPositionLastInSourceCode
		// ^^^ the tokenCloser's last position is saved with this!
		if parent.ValType == typeArray {
			elems := parent.ValArray
			elems = append(elems, value)
			parent.ValArray = elems
		}
		if parent.ValType == typeObject {
			parent_valObjects := parent.ValObject
			parent_valObjects[lastDetectedStringKey__inObject] = value
			lastDetectedStringKey__inObject = "" // clear the keyName, we used that for the current object
			parent.ValObject = parent_valObjects
		}
		containers[idParent] = parent // save back the updated parent

	} // for, tokenNum, tokenPositionKey
	return elemRoot, errorsCollected
}

// //////////////////// VALUE setter FUNCTIONS ///////////////////////////////////////////////
func valueValidationsSettings_inTokens(tokens tokenTable_startPositionIndexed, errorsCollected []error) (tokenTable_startPositionIndexed, []error) {
	tokensUpdated := tokenTable_startPositionIndexed{}
	for _, tokenOne := range tokens {
		tokenOne, errorsCollected = value_validate_and_set__elem_string(tokenOne, errorsCollected)
		tokenOne, errorsCollected = value_validate_and_set__elem_number(tokenOne, errorsCollected)
		// TODO: elem true|false|null value set?
		tokensUpdated[tokenOne.charPositionFirstInSourceCode] = tokenOne
	}
	return tokensUpdated, errorsCollected
}

// set the string value from raw strings
func value_validate_and_set__elem_string(token token, errorsCollected []error) (token, []error) { // TESTED

	if token.valType != typeString {
		return token, errorsCollected
	} // don't modify non-string tokens

	/* Tasks:
	- is it a valid string?
	- convert special char representations to real chars

	the func works typically with 2 chars, for example: \t
	but sometime with 6: \u0123, so I need to look forward for the next 5 chars
	*/

	src := token.runes
	src = src[1 : len(src)-1] // "remove opening/closing quotes from the string value"

	valueFromRawSrcParsing := []rune{}

	// fmt.Println("string token value detection:", src)
	runeBackSlash := '\\' // be careful: this is ONE \ char, only written with this expression

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
				}
				if runeNext1 == runeBackSlash { // in reality, these are the 2 chars: \\
					runeReal = '\\' // reverse solidus
				}
				if runeNext1 == '/' { // a very special escaping: \/
					runeReal = '/' // solidus
				}
				if runeNext1 == 'b' { // This is the first good example for escaping:
					runeReal = '\b' // in the src there were 2 chars: \ and b,
				} //  (backspace)    // and one char is inserted into the stringVal

				if runeNext1 == 'f' { // formfeed
					runeReal = '\f'
				}

				if runeNext1 == 'n' { // linefeed
					runeReal = '\n'
				}

				if runeNext1 == 'r' { // carriage return
					runeReal = '\r' //
				}

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
	token.valString = string(valueFromRawSrcParsing)
	return token, errorsCollected
}

func value_validate_and_set__elem_number(token token, errorsCollected []error) (token, []error) {

	if token.valType != typeNumber_exactTypeIsNotSet {
		return token, errorsCollected
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

	numberRunes := base__runes_copy(token.runes)

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
	isFractionDotUsed := strings.ContainsRune(string(numberRunes), '.')
	isExponent_E_used := strings.ContainsRune(string(numberRunes), 'E')
	isExponent_e_used := strings.ContainsRune(string(numberRunes), 'e')

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
				token.valNumberInt = numBase10
				token.valType = typeNumberInt
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
				token.valNumberFloat = numBase10
				token.valType = typeNumberFloat64
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

	return token, errorsCollected
}

// ////////////////////  DETECTIONS  ///////////////////////////////////////////////
// Documented in program plan
func jsonDetect_strings______(src []rune, tokensStartPositions tokenTable_startPositionIndexed, errorsCollected []error) ([]rune, tokenTable_startPositionIndexed, []error) { // TESTED

	srcDetectedTokensRemoved := []rune{}
	// to find escaped \" \\\" sections in strings
	escapeBackSlashCounterBeforeCurrentChar := 0

	inStringDetection := false

	isEscaped := func() bool {
		return escapeBackSlashCounterBeforeCurrentChar%2 != 0
	}

	var tokenNow token

	for posInSrc, runeActual := range src {

		if runeActual == '"' {
			if !inStringDetection { // if at " char handling, we are NOT in string
				inStringDetection = true

				tokenNow = token{valType: typeString}
				tokenNow.charPositionFirstInSourceCode = posInSrc
				tokenNow.runes = append(tokenNow.runes, runeActual)

				srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, ' ')
				continue
			}
			if inStringDetection && !isEscaped() { // a non-escaped " char in a string detection
				inStringDetection = false          // is the end of the string

				tokenNow.charPositionLastInSourceCode = posInSrc
				tokenNow.runes = append(tokenNow.runes, runeActual)
				tokensStartPositions[tokenNow.charPositionFirstInSourceCode] = tokenNow // save token

				srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, ' ')
				continue
			}

			// BE CAREFUL, there is a 3rd option!
			// if inStringDetection && isEscaped() -- which is handled as part of a string: inStringDetection
			// and this is different from the previous two, so don't change the structure!

		} // if " is detected, everything is handled in the conditions

		if inStringDetection {
			tokenNow.runes = append(tokenNow.runes, runeActual)

			if runeActual == '\\' {
				escapeBackSlashCounterBeforeCurrentChar++
			} else { // the escape series ended :-)
				escapeBackSlashCounterBeforeCurrentChar = 0
			}

			// add empty placeholder where the token was detected
			srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, ' ')
		} else {
			// save the original rune, if it was not in a string
			srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, runeActual)
		}

	} // for, runeActual

	if inStringDetection {
		errorsCollected = append(errorsCollected, errors.New("non-closed string detected:"))
	}

	return srcDetectedTokensRemoved, tokensStartPositions, errorsCollected
} // detect strings

func jsonDetect_separators___(src []rune, tokensStartPositions tokenTable_startPositionIndexed, errorsCollected []error) ([]rune, tokenTable_startPositionIndexed, []error) { // TESTED
	srcDetectedTokensRemoved := []rune{}
	var tokenNow token

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
			// save the original rune, if it was not a detected char
			srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, runeActual)
		} else { // save token, if something important is detected
			tokenNow = token{valType: detectedType}
			tokenNow.charPositionFirstInSourceCode = posInSrc
			tokenNow.charPositionLastInSourceCode = posInSrc
			tokenNow.runes = append(tokenNow.runes, runeActual)
			srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, ' ')
			tokensStartPositions[tokenNow.charPositionFirstInSourceCode] = tokenNow

			detectedType = typeIsUnknown
			// set back the type to unknown, if token is handled
		}
	} // for runeActual
	return srcDetectedTokensRemoved, tokensStartPositions, errorsCollected
} // separators....


/*
	 this detection is AFTER string+separator detection.
	   	in other words: only numbers and true/false/null values are left in the src.

		because the strings/separators are removed and replaced with space in the src, as placeholders,
	    the true/false/null words are surrounded with spaces, as separators.
*/
func jsonDetect_trueFalseNull(src []rune, tokensStartPositions tokenTable_startPositionIndexed, errorsCollected []error) ([]rune, tokenTable_startPositionIndexed, []error) { // TESTED

	srcDetectedTokensRemoved := base__runes_copy(src)
	// copy the original structure, not use the same variable
	// the detected word runes will be deleted from here.

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
			tokenNow := token{valType: detectedType}
			if detectedType == typeBool {
				tokenNow.valBool = base__compare_runes_are_equal(wordOne.wordChars, charsTrue)
			}
			tokenNow.charPositionFirstInSourceCode = wordOne.posFirst
			tokenNow.charPositionLastInSourceCode = wordOne.posLast

			for posDetected := wordOne.posFirst; posDetected <= wordOne.posLast; posDetected++ {
				// save all detected positions:
				tokenNow.runes = append(tokenNow.runes, (src)[posDetected])
				// clear detected positions from the src:
				srcDetectedTokensRemoved[posDetected] = ' '
				// only detected word runes are removed from the storage, where ALL original src is inserted in the first step

			}
			tokensStartPositions[tokenNow.charPositionFirstInSourceCode] = tokenNow

			detectedType = typeIsUnknown // set the default value again ONLY if it was not unknown, in this case
		}
		/*	 it can be a number too, in case of else - so it is not an error, if typeIsUnknown */
	}
	return srcDetectedTokensRemoved, tokensStartPositions, errorsCollected
}

// words are detected here, and I can hope only that they are numbers - later they will be validated
func jsonDetect_numbers______(src []rune, tokensStartPositions tokenTable_startPositionIndexed, errorsCollected []error) ([]rune, tokenTable_startPositionIndexed, []error) { // TESTED

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// there is no reason to give back the cleaned src, because this is the last function and src is NOT USED IN THE PROCESS ANYMORE
	// srcDetectedTokensRemoved := base__runes_copy(src) // copy the original structure, not use the same variable
	srcDetectedTokensRemoved := []rune{} // I would like to keep the standard return structure, so return with an empty struct
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	for _, wordOne := range base__src_get_whitespace_separated_words_posFirst_posLast(src) {

		tokenNow := token{valType: typeNumber_exactTypeIsNotSet} // only numbers can be in the src now.
		tokenNow.charPositionFirstInSourceCode = wordOne.posFirst
		tokenNow.charPositionLastInSourceCode = wordOne.posLast

		for posDetected := wordOne.posFirst; posDetected <= wordOne.posLast; posDetected++ {
			// save all detected positions:
			tokenNow.runes = append(tokenNow.runes, (src)[posDetected])
			// clear detected positions from the src:

			/////////////////////////////////////////
			//// don't use this, nobody reads src later
			// srcDetectedTokensRemoved[posDetected] = ' '
			/////////////////////////////////////////


		}
		tokensStartPositions[tokenNow.charPositionFirstInSourceCode] = tokenNow
	}
	return srcDetectedTokensRemoved, tokensStartPositions, errorsCollected
}

/////////////////////////// local tools: supporter funcs, not for main logic /////////////////////////////

// tokenTable keys are character positions in JSON source code (positive integers)
func local_tool__tokenTable_position_keys_sorted(tokens tokenTable_startPositionIndexed) []int {
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
		for _, elem := range v.ValArray {
			elem.LevelInObjectStructure = v.LevelInObjectStructure + 1
			arrayUpdated = append(arrayUpdated, elem)
		}
		v.ValArray = arrayUpdated
	}
	if v.ValType == typeObject {
		objectsUpdated := map[string]JSON_value{}

		for key, elem := range v.ValObject {
			elem.LevelInObjectStructure = v.LevelInObjectStructure + 1
			objectsUpdated[key] = elem
		}
		v.ValObject = objectsUpdated
	}

}
