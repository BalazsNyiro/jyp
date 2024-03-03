// jsonB - Balazs' Json parser
// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

// this file is the implementation of the _standard_ json data format:
// https://www.json.org/json-en.html

package jyp

import "errors"

const ABC_lower string = "abcdefghijklmnopqrstuvwxyz"
const ABC_upper string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

const digitZero string = "0"
const digits_19 string = "123456789"
const digits_09 string = "0123456789"

const whitespaces string = "\r\n\t "

const objOpen  string = "{"
const arrClose string = "}"
const arrOpen  string = "["
const objClose string = "]"

const separatorComma string = ","
const separatorColon string = ":"
const separatorDot   string = "."
const separatorExponent string = "eE"

const separatorMinus     string = "-"
const separatorPlusMinus string = "+-"



type Elems []Elem
type ElemMap map[string]Elem

type Elem struct {
	Type string
	// possible types:
	// array, object,
	// bool, null, string, number_int, number_float,

	ValArray       Elems
	ValObject      ElemMap

	ValBool        bool // true, false
	isNull		   bool // if true, then the value is null

	ValString      string
	ValNumberInt   int
	ValNumberFloat float64

	//////// PARSING SECTION: detection from the JSON source code /////
	charPositionFirstInSourceCode int   // 0: the first char in source code, 1: 2nd...
	charPositionLastInSourceCode  int   // 0: the first char in source code, 1: 2nd...
	runes []rune
}

type tokenTable_startPositionIndexed map[int]Elem

// if the src can be parsed, return with the JSON root object with nested elems, and err is nil.
func JsonParse(src string) (Elem, error) {
	elemRoot := Elem{}

	errorsCollected := []error{}
	tokens := tokenTable_startPositionIndexed{}

	// a simple rule - inputs:  src, tokens, errors are inputs,
	//                 outputs: src, tokens, errors
	// the src is always less and less, as tokens are detected
	// the tokens table has more and more elems, as the src sections are parsed
	// at the end, src is total empty (if everything goes well) - and we don't have errors, too
	src, tokens, errorsCollected = json_string_detect(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_separators_detect(src, tokens, errorsCollected)
	return elemRoot, nil
}


////////////////////// BASE FUNCTIONS ///////////////////////////////////////////////
func json_string_detect(src string, tokensStartPositions tokenTable_startPositionIndexed, errorsCollected []error) (string, tokenTable_startPositionIndexed, []error) {

	srcDetectedTokensRemoved := []rune{}
	// to find escaped \" \\\" sections in strings
	escapeBackSlashCounterBeforeCurrentChar := 0

	inStringDetection := false

	isEscaped := func() bool {
		return escapeBackSlashCounterBeforeCurrentChar % 2 != 0
	}

	var tokenNow Elem

	for posInSrc, r := range src {

		if r == '"' {
			if !inStringDetection {
					tokenNow = Elem{Type: "string"}
					inStringDetection = true
					tokenNow.charPositionFirstInSourceCode = posInSrc
					tokenNow.runes = append(tokenNow.runes, r)
					srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, ' ')
					continue
			} else { // in string detection
				if ! isEscaped() {
					inStringDetection = false
					tokenNow.charPositionLastInSourceCode = posInSrc
					tokenNow.runes = append(tokenNow.runes, r)
					tokensStartPositions[tokenNow.charPositionFirstInSourceCode] = tokenNow
					srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, ' ')
					continue
				}
			}
		} // if " is detected, everything is handled in the conditions


		if inStringDetection {
			tokenNow.runes = append(tokenNow.runes, r)

			if r == '\\' {
				escapeBackSlashCounterBeforeCurrentChar++
			} else { // the escape series ended :-)
				escapeBackSlashCounterBeforeCurrentChar = 0
			}

			// add empty placeholder where the token was detected
			srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, ' ')
		} else {
			// save the original rune, if it was not in a string
			srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, r)
		}

	} // for

	if inStringDetection {
		errorsCollected = append(errorsCollected, errors.New("non-closed string detected:"))
	}

	return string(srcDetectedTokensRemoved), tokensStartPositions, errorsCollected
}


func json_separators_detect(src string, tokensStartPositions tokenTable_startPositionIndexed, errorsCollected []error) (string, tokenTable_startPositionIndexed, []error) {
	srcDetectedTokensRemoved := []rune{}
	var tokenNow Elem

	for posInSrc, r := range src {
		detectedType := ""

		if r == '{' { detectedType = "objectOpen"  }
		if r == '}' { detectedType = "objectClose" }
		if r == '[' { detectedType = "arrayOpen"   }
		if r == ']' { detectedType = "arrayClose"  }
		if r == ',' { detectedType = "comma"       }
		if r == ':' { detectedType = "colon"       }

		if detectedType == "" {
			// save the original rune, if it was not a detected char
			srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, r)
		} else { // save Elem, if something important is detected
			tokenNow = Elem{Type: detectedType}
			tokenNow.charPositionFirstInSourceCode = posInSrc
			tokenNow.runes = append(tokenNow.runes, r)
			srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, ' ')
			continue
		}
	} // for r
	return string(srcDetectedTokensRemoved), tokensStartPositions, errorsCollected
}



/////////////////////// base functions /////////////////
