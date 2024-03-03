// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

// this file is the implementation of the _standard_ json data format:
// https://www.json.org/json-en.html

// this song helped a lot to write this parser - respect:
// https://open.spotify.com/track/7znjTquY8gek1bKni5yzLG?si=3ae71af19f684d67

package jyp

import (
	"errors"
	"fmt"
	"strings"
)

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
	src, tokens, errorsCollected = json_detect_strings________(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_separators_____(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_true_false_null(src, tokens, errorsCollected)
	return elemRoot, nil
}


////////////////////// BASE FUNCTIONS ///////////////////////////////////////////////
func json_detect_strings________(src string, tokensStartPositions tokenTable_startPositionIndexed, errorsCollected []error) (string, tokenTable_startPositionIndexed, []error) {

	srcDetectedTokensRemoved := []rune{}
	// to find escaped \" \\\" sections in strings
	escapeBackSlashCounterBeforeCurrentChar := 0

	inStringDetection := false

	isEscaped := func() bool {
		return escapeBackSlashCounterBeforeCurrentChar % 2 != 0
	}

	var tokenNow Elem

	for posInSrc, runeActual := range src {

		if runeActual == '"' {
			if !inStringDetection {
					tokenNow = Elem{Type: "string"}
					inStringDetection = true
					tokenNow.charPositionFirstInSourceCode = posInSrc
					tokenNow.runes = append(tokenNow.runes, runeActual)
					srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, ' ')
					continue
			} else { // in string detection
				if ! isEscaped() {
					inStringDetection = false
					tokenNow.charPositionLastInSourceCode = posInSrc
					tokenNow.runes = append(tokenNow.runes, runeActual)
					tokensStartPositions[tokenNow.charPositionFirstInSourceCode] = tokenNow
					srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, ' ')
					continue
				}
			}
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

	} // for

	if inStringDetection {
		errorsCollected = append(errorsCollected, errors.New("non-closed string detected:"))
	}

	return string(srcDetectedTokensRemoved), tokensStartPositions, errorsCollected
}


func json_detect_separators_____(src string, tokensStartPositions tokenTable_startPositionIndexed, errorsCollected []error) (string, tokenTable_startPositionIndexed, []error) {
	srcDetectedTokensRemoved := []rune{}
	var tokenNow Elem

	for posInSrc, runeActual := range src {
		detectedType := ""

		if runeActual == '{' { detectedType = "objectOpen"  }
		if runeActual == '}' { detectedType = "objectClose" }
		if runeActual == '[' { detectedType = "arrayOpen"   }
		if runeActual == ']' { detectedType = "arrayClose"  }
		if runeActual == ',' { detectedType = "comma"       }
		if runeActual == ':' { detectedType = "colon"       }

		if detectedType == "" {
			// save the original rune, if it was not a detected char
			srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, runeActual)
		} else { // save Elem, if something important is detected
			tokenNow = Elem{Type: detectedType}
			tokenNow.charPositionFirstInSourceCode = posInSrc
			tokenNow.charPositionLastInSourceCode  = posInSrc
			tokenNow.runes = append(tokenNow.runes, runeActual)
			srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, ' ')
			tokensStartPositions[tokenNow.charPositionFirstInSourceCode] = tokenNow
		}
	} // for runeActual
	return string(srcDetectedTokensRemoved), tokensStartPositions, errorsCollected
}


/* this detection is AFTER string+separator detection.
   	in other words: only numbers and true/false/null values are left in the src.

	because the strings/separators are removed and replaced with space in the src, as placeholders,
    the true/false/null words are surrounded with spaces, as separators.
*/
func json_detect_true_false_null(src string, tokensStartPositions tokenTable_startPositionIndexed, errorsCollected []error) (string, tokenTable_startPositionIndexed, []error) {
	srcDetectedTokensRemoved := []rune{}
	var tokenNow Elem

	for posActual := 0; posActual < len(src); posActual++ {

		runeActual := src_get_char(src, posActual    )
		runeNext1  := src_get_char(src, posActual + 1)   // the real rune value IF the pos in the valid range of the src
		runeNext2  := src_get_char(src, posActual + 2)   // or space, if the index is bigger/lower than the valid range
		runeNext3  := src_get_char(src, posActual + 3)
		runeNext4  := src_get_char(src, posActual + 4)
		runeNext5  := src_get_char(src, posActual + 5)
		runeNext6  := src_get_char(src, posActual + 6)

		//                                               ' '            n             u           l        l         ' '         ' '  (closing space after the word)
		//                                               ' '            f             a           l        s          e          ' '  (closing space)
		word_ActualChar_plus_few_chars := string([]rune{runeActual, runeNext1, runeNext2, runeNext3, runeNext4, runeNext5, runeNext6})
		// A good question: why don't I use a simple string indexing? ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
		// because maybe I over index the src, so the wanted index is NOT in the valid range
		// because of this, runes are collected one by one, and if the index is NOT in the range, substituted with a meaningless SPACE

		detectedType := ""
		posFirst := 0
		posLast := 0

		// the word has to be detected WITH SPACE boundaries
		if strings.HasPrefix(word_ActualChar_plus_few_chars, " true ") {
			detectedType = "true"
			posFirst = posActual
			posLast = posActual + 5
		}

		if strings.HasPrefix(word_ActualChar_plus_few_chars, " false ") {
			detectedType = "false"
			posFirst = posActual
			posLast = posActual + 6
		}

		if strings.HasPrefix(word_ActualChar_plus_few_chars, " null ") {
			detectedType = "null"
			posFirst = posActual
			posLast = posActual + 5
		}
		fmt.Println("DETECTED TYPE", detectedType)
		if detectedType == "" {
			// save the original rune, if it was not a detected char
			srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, runeActual)
		} else { // save Elem, if something important is detected
			tokenNow = Elem{Type: detectedType}
			tokenNow.charPositionFirstInSourceCode = posFirst
			tokenNow.charPositionLastInSourceCode  = posLast

			for posDetected := posFirst; posDetected <=posLast; posDetected++ {
				// save all detected positions:
				tokenNow.runes = append(tokenNow.runes, ([]rune(src))[posDetected])

				// clear all detected positions from the src:
				srcDetectedTokensRemoved = append(srcDetectedTokensRemoved, ' ')

			}

			// set the actual position to the last detected pos,
			// because all chars were added to the Elem between posFirst->posLast,
			// so there is no reason to detect them again :-)
			posActual = posLast

			tokensStartPositions[tokenNow.charPositionFirstInSourceCode] = tokenNow
		}
	} // for runeActual
	return string(srcDetectedTokensRemoved), tokensStartPositions, errorsCollected
}

// get the rune IF the index is really in the range of the src.
// return with ' ' space, IF the index is NOT in the range.
// reason: avoid never ending index checking, so do it only once
// the space can be answered because this func is used when a real char wanted to be detected,
// and if a space is returned, this has NO MEANING in that parse section
func src_get_char(src string, pos int) rune {
	posPossibleMax := len(src)-1
	posPossibleMin := 0
	if len(src)	== 0 { // if the src is empty, posPossibleMax == -1, min cannot be bigger than max
		posPossibleMin = -1
	}
	if (pos >= posPossibleMin) && (pos <= posPossibleMax) {
		return ([]rune(src))[pos]
	}
	return ' '
}


/////////////////////// base functions /////////////////
