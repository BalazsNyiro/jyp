// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

// this file is the implementation of the _standard_ json data format:
// https://www.json.org/json-en.html

// this song helped a lot to write this parser - respect:
// https://open.spotify.com/track/7znjTquY8gek1bKni5yzLG?si=3ae71af19f684d67

package jyp

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)


type Elems []tokenElem
type ElemMap map[string]tokenElem

type tokenElem struct {
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

type tokenTable_startPositionIndexed map[int]tokenElem

// if the src can be parsed, return with the JSON root object with nested elems, and err is nil.
func JsonParse(src string) (tokenElem, error) {
	elemRoot := tokenElem{}

	errorsCollected := []error{}
	tokens := tokenTable_startPositionIndexed{}

	// a simple rule - inputs:  src, tokens, errors are inputs,
	//                 outputs: src, tokens, errors
	// the src is always less and less, as tokens are detected
	// the tokens table has more and more elems, as the src sections are parsed
	// at the end, src is total empty (if everything goes well) - and we don't have errors, too

	// only strings can have errors at this parsing step, but the src|tokens|errors are
	// lead through every fun, as a standard solution - so the possibility is open to throw an error everywhere.

	// here maybe the tokens|errorsCollected ret val handling could be removed,
	// but with this, it is clearer what is happening in the fun - so I use this form.
	// in other words: represent if the structure is changed in the function.
	src, tokens, errorsCollected = json_detect_strings________(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_separators_____(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_true_false_null(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_numbers________(src, tokens, errorsCollected)

	// at this point, Numbers are not validated - the ruins are collected only,
	// and the lists/objects doesn't have embedded structures - it has to be built, too.
	// src has to be empty, or contain only whitespaces.

	TokensDisplay(tokens)

	// set correct string values, based on raw rune src.
	// example: "\u0022quote\u0022"'s real form: `"quote"`,
	// so the raw source has to be interpreted (escaped chars, unicode chars)
	tokens, errorsCollected = tokens_validations_value_settings(tokens, errorsCollected)


	return elemRoot, nil
}


////////////////////// VALUE setter FUNCTIONS ///////////////////////////////////////////////
func tokens_validations_value_settings(tokens tokenTable_startPositionIndexed, errorsCollected []error) (tokenTable_startPositionIndexed, []error) {
	tokensUpdated := tokenTable_startPositionIndexed{}
	for _, token := range tokens {
		fmt.Println("\n>>> one Token value Before detection:", token.ValString)
		token, errorsCollected = token_string_value_validate_and_set(token, errorsCollected)
		fmt.Println("<<< one Token value After detection:", token.ValString)
		tokensUpdated[token.charPositionFirstInSourceCode] = token
	}
	return tokensUpdated, errorsCollected
}


// set the string value from raw strings
func token_string_value_validate_and_set(token tokenElem, errorsCollected []error) (tokenElem, []error) {

	if token.Type != "string" {
		return token, errorsCollected
	} // don't modify non-string tokens

	/* Tasks:
	 - is it a valid string?
	 - convert special char representations to real chars

	 the func works typically with 2 chars, for example: \t
	 but sometime with 6: \u0123, so I need to look forward for the next 5 chars
	*/

	src := string(token.runes)
	src = src[1:len(src)-1]  // "remove opening/closing quotes from the string value"

	valueFromRawSrcParsing := []rune{}

	fmt.Println("string tokenElem value detection:", src)

	for pos := 0; pos < len(src); pos++ {

		runeActual := src_get_char(src, pos)
		// runeNext1 := src_get_char(src, pos+1)
		// runeNext2 := src_get_char(src, pos+2)
		// runeNext3 := src_get_char(src, pos+3)
		// runeNext4 := src_get_char(src, pos+4)
		// runeNext5 := src_get_char(src, pos+5)

		if runeActual != '\\' {  // a non-backSlash char
			valueFromRawSrcParsing = append(valueFromRawSrcParsing, runeActual)
			continue
		}

	}
	fmt.Println("value from raw src parsing:", string(valueFromRawSrcParsing))
	token.ValString = string(valueFromRawSrcParsing)
	return token, errorsCollected
}




////////////////////// BASE FUNCTIONS ///////////////////////////////////////////////
func json_detect_strings________(src string, tokensStartPositions tokenTable_startPositionIndexed, errorsCollected []error) (string, tokenTable_startPositionIndexed, []error) { // TESTED

	srcDetectedTokensRemoved := []rune{}
	// to find escaped \" \\\" sections in strings
	escapeBackSlashCounterBeforeCurrentChar := 0

	inStringDetection := false

	isEscaped := func() bool {
		return escapeBackSlashCounterBeforeCurrentChar % 2 != 0
	}

	var tokenNow tokenElem

	for posInSrc, runeActual := range src {

		if runeActual == '"' {
			if !inStringDetection {
					tokenNow = tokenElem{Type: "string"}
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


func json_detect_separators_____(src string, tokensStartPositions tokenTable_startPositionIndexed, errorsCollected []error) (string, tokenTable_startPositionIndexed, []error) { // TESTED
	srcDetectedTokensRemoved := []rune{}
	var tokenNow tokenElem

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
		} else { // save tokenElem, if something important is detected
			tokenNow = tokenElem{Type: detectedType}
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
func json_detect_true_false_null(src string, tokensStartPositions tokenTable_startPositionIndexed, errorsCollected []error) (string, tokenTable_startPositionIndexed, []error) { // TESTED
	srcDetectedTokensRemoved := []rune(src)

	for _, wordOne := range src_get_whitespace_separated_words_posFirst_posLast(src) {

		detectedType := "" // 3 types of word can be detected in this fun
		if wordOne.word == "true"  { detectedType = "true"  }
		if wordOne.word == "false" { detectedType = "false" }
		if wordOne.word == "null"  { detectedType = "false" }

		if detectedType != "" {
			tokenNow := tokenElem{Type: detectedType}
			tokenNow.charPositionFirstInSourceCode = wordOne.posFirst
			tokenNow.charPositionLastInSourceCode  = wordOne.posLast

			for posDetected := wordOne.posFirst; posDetected <= wordOne.posLast; posDetected++ {
				// save all detected positions:
				tokenNow.runes = append(tokenNow.runes, ([]rune(src))[posDetected])
				// clear detected positions from the src:
				srcDetectedTokensRemoved[posDetected] = ' '
			}
			tokensStartPositions[tokenNow.charPositionFirstInSourceCode] = tokenNow
		}
	}
	return string(srcDetectedTokensRemoved), tokensStartPositions, errorsCollected
}


// words are detected here, and I can hope only that they are numbers - later they will be validated
func json_detect_numbers________(src string, tokensStartPositions tokenTable_startPositionIndexed, errorsCollected []error) (string, tokenTable_startPositionIndexed, []error) { // TESTED
	srcDetectedTokensRemoved := []rune(src)

	for _, wordOne := range src_get_whitespace_separated_words_posFirst_posLast(src) {

		tokenNow := tokenElem{Type: "number"} // only numbers can be in the src now.
		tokenNow.charPositionFirstInSourceCode = wordOne.posFirst
		tokenNow.charPositionLastInSourceCode  = wordOne.posLast

		for posDetected := wordOne.posFirst; posDetected <= wordOne.posLast; posDetected++ {
			// save all detected positions:
			tokenNow.runes = append(tokenNow.runes, ([]rune(src))[posDetected])
			// clear detected positions from the src:
			srcDetectedTokensRemoved[posDetected] = ' '
		}
		tokensStartPositions[tokenNow.charPositionFirstInSourceCode] = tokenNow
	}
	return string(srcDetectedTokensRemoved), tokensStartPositions, errorsCollected
}




////////////////////////////////////
type word struct {
	word string
	posFirst int
	posLast int
}

func src_get_whitespace_separated_words_posFirst_posLast(src string) []word { // TESTED
	words := []word{}

	wordChars := []rune{}
	posFirst  := -1
	posLast   := -1

	// posActual := -1, len(src) + 1: overindexing!
	// with this, I can be sure that minimum one space is detected first,
	// and minimum one space detected after the source code's normal chars!
	// with this solution, the last word detection can be closed with the last boundary space, in one
	// case, and I don't have to handle that later, in a second if/else condition

	// src_get_char() handles the overindexing
	for posActual := -1; posActual < len(src)+1; posActual++ {
		runeActual := src_get_char(src, posActual)

		// the first and last chars, because of overindexing, are spaces, this is guaranteed!
		if is_whitespace_rune(runeActual) {
			if len(wordChars) > 0 {
				word := word{
					word    : string(wordChars),
					posFirst: posFirst,
					posLast : posLast,
				}
				words = append(words, word)
			}
			wordChars = []rune{}
			posFirst  = -1
			posLast   = -1

		} else {
			// save posFirst, posLast, and word-builder chars ///
			if len(wordChars) == 0 {
				posFirst = posActual
			}
			posLast = posActual
			wordChars = append(wordChars, runeActual)
		}

	}


	return words
}
////////////////////////////////////

// get the rune IF the index is really in the range of the src.
// return with ' ' space, IF the index is NOT in the range.
// reason: avoid never ending index checking, so do it only once
// the space can be answered because this func is used when a real char wanted to be detected,
// and if a space is returned, this has NO MEANING in that parse section
// this fun is NOT used in string detection - and other places whitespaces can be neglected, too
func src_get_char(src string, pos int) rune {  // TESTED
	posPossibleMax := len(src)-1
	posPossibleMin := 0
	if len(src)	== 0 { // if the src is empty, posPossibleMax == -1, min cannot be bigger than max
		posPossibleMin = -1
	}
	if (pos >= posPossibleMin) && (pos <= posPossibleMax) {
		charSelected := ([]rune(src))[pos]
		if is_whitespace_rune(charSelected) {
			charSelected = ' '
			// simplify everything. if the char is a whitespace, return with SPACE
		}
		return charSelected
	}
	return ' '
}

// the string has whitespace chars only
func is_whitespace_string(src string) bool { // TESTED
	return strings.TrimSpace(src) == ""
}

// the rune is a whitespace char
func is_whitespace_rune(oneRune rune) bool { // TESTED
	return is_whitespace_string(string([]rune{oneRune}))
}


func TokensDisplay(tokens tokenTable_startPositionIndexed) {
	keys := make([]int, 0, len(tokens))
	for k := range tokens {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	fmt.Println("== Tokens Table display ==")
	for _, key := range keys{
		fmt.Println(string(tokens[key].runes), key, tokens[key])
	}
}
/////////////////////// base functions /////////////////
