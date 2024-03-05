// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

// this file is the implementation of the _standard_ json data format:
// https://www.json.org/json-en.html

// this song helped a lot to write this parser - respect:
// https://open.spotify.com/track/7znjTquY8gek1bKni5yzLG?si=3ae71af19f684d67

// in the code I intentionally avoid direct pointer usage - I think that is safer.

package jyp

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)


type JsonValueList []JsonValue
type JsonValueObject map[string]JsonValue

type JsonValue struct {
	Type string
	// possible types:
	// array, object,
	// bool, null, string, number_int, number_float,

	ValArray  JsonValueList
	ValObject JsonValueObject

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

type tokenTable_startPositionIndexed map[int]JsonValue

// if the src can be parsed, return with the JSON root object with nested elems, and err is nil.
func JsonParse(src string) (JsonValue, error) {
	elemRoot := JsonValue{}

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
		// fmt.Println("\n>>> one Token value Before detection:", token.ValString)
		token, errorsCollected = elem_string_value_validate_and_set(token, errorsCollected)
		// fmt.Println("<<< one Token value After detection:", token.ValString)

		token, errorsCollected = elem_number_value_validate_and_set(token, errorsCollected)
		// TODO: elem true|false|null value set?
		tokensUpdated[token.charPositionFirstInSourceCode] = token
	}
	return tokensUpdated, errorsCollected
}


// set the string value from raw strings
func elem_string_value_validate_and_set(token JsonValue, errorsCollected []error) (JsonValue, []error) { // TESTED

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

	fmt.Println("string JsonValue value detection:", src)
	runeBackSlash := '\\' // be careful: this is ONE \ char, only written with this expression

	for pos := 0; pos < len(src); pos++ {

		runeActual := src_get_char(src, pos)
		//fmt.Println("rune actual (string value set):", pos, string(runeActual), runeActual)
		runeNext1 := src_get_char(src, pos+1)

		if runeActual != runeBackSlash {  // a non-backSlash char
			valueFromRawSrcParsing = append(valueFromRawSrcParsing, runeActual)
			continue
		} else {
			// runeActual is \\ here, so ESCAPING started

			if runeNext1 == 'u' {
				// this is \u.... unicode code point - special situation,
				// because after the \u four other chars has to be handled

				runeNext2 := src_get_char(src, pos+2)
				runeNext3 := src_get_char(src, pos+3)
				runeNext4 := src_get_char(src, pos+4)
				runeNext5 := src_get_char(src, pos+5)


				base10_val_2, err2 := hexaRune_to_intVal(runeNext2)
				if err2 != nil {  errorsCollected = append(errorsCollected, err2)	}

				base10_val_3, err3 := hexaRune_to_intVal(runeNext3)
				if err3 != nil {  errorsCollected = append(errorsCollected, err3)	}

				base10_val_4, err4 := hexaRune_to_intVal(runeNext4)
				if err4 != nil {  errorsCollected = append(errorsCollected, err4)	}

				base10_val_5, err5 := hexaRune_to_intVal(runeNext5)
				if err5 != nil {  errorsCollected = append(errorsCollected, err5)	}


				unicodeVal_10Based := 0

				if err2 == nil && err3 == nil && err4 == nil && err5 == nil {
					unicodeVal_10Based = base10_val_2*16*16*16 +
						                 base10_val_3*16*16 +
						                 base10_val_4*16 +
						                 base10_val_5
				}
				runeFromHexaDigits := rune(unicodeVal_10Based)

				pos += 1+4 // one extra pos because of the u, and +4 because of the digits
				valueFromRawSrcParsing = append(valueFromRawSrcParsing, runeFromHexaDigits)


			} else { // the first detected char was a backslash, what is the second?
				// so this is a simple escaped char, for example: \" \t \b \n
				runeReal := '?'
				if runeNext1 == '"' {   // \" -> is a " char in a string
					runeReal = '"'      // in a string, this is an escaped " double quote char
				}
				if runeNext1 == runeBackSlash {  // in reality, these are the 2 chars: \\
					runeReal = '\\' // reverse solidus
				}
				if runeNext1 == '/'{ // a very special escaping: \/
					runeReal = '/'   // solidus
				}
				if runeNext1 == 'b'{ // This is the first good example for escaping:
					runeReal = '\b'  // in the src there were 2 chars: \ and b,
				} //  (backspace)    // and one char is inserted into the stringVal

				if runeNext1 == 'f'{ // formfeed
					runeReal = '\f'
				}

				if runeNext1 == 'n'{ // linefeed
					runeReal = '\n'
				}

				if runeNext1 == 'r'{ // carriage return
					runeReal = '\r'  //
				}

				if runeNext1 == 't'{ // horizontal tab
					runeReal = '\t'  //
				}

				pos += 1 // one extra pos increasing is necessary, because of
				// 2 chars were processed: the actual \ and the next one.

				valueFromRawSrcParsing = append(valueFromRawSrcParsing, runeReal)
			}
		} // else
	} // for

	fmt.Println("value from raw src parsing:", string(valueFromRawSrcParsing))
	token.ValString = string(valueFromRawSrcParsing)
	return token, errorsCollected
}



func elem_number_value_validate_and_set(token JsonValue, errorsCollected []error) (JsonValue, []error) {

	if token.Type != "number" { return token, errorsCollected } // don't modify non-number elems

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

	numberRunes := runes_copy(token.runes)

	// example number: -1234.567e-8
	isNegative := numberRunes[0] == '-'
	if isNegative { numberRunes = numberRunes[1:] } // remove - sign, if that is the first

	runesSectionInteger  := []rune{}   // and at the end, only integers remain...
	runesSectionFraction := []rune{}  // filled second
	runesSectionExponent := []rune{} // filled first
	////////////////// the main sections of the number



	///////////////// the main markers of the number
	isFractionDotUsed   := strings.ContainsRune(string(numberRunes), '.')
	isExponent_E_used   := strings.ContainsRune(string(numberRunes), 'E')
	isExponent_e_used   := strings.ContainsRune(string(numberRunes), 'e')


	/////////// go from back to forward: remove exponent part first
	if isExponent_e_used {
		numberRunes, runesSectionExponent = runes_split_at_pattern(numberRunes, 'e')
	}
	if isExponent_E_used {
		numberRunes, runesSectionExponent = runes_split_at_pattern(numberRunes, 'E')
	} /////////// if exponent is used, that is filled into the runeSectionExponentPart


	if isFractionDotUsed{
		numberRunes, runesSectionFraction = runes_split_at_pattern(numberRunes, '.')
		// after this, numberRunes lost the integer part.
	} // if FractionDot is used, split the runes :-)

	runesSectionInteger = numberRunes
	///////////////////

	lenErrorCollectorBeforeErrorDetection := len(errorsCollected)


	/////////// ERROR HANDLING ////////////
	// if the first digit is 0, there cannot be more digits.

	if len(runesSectionInteger) > 1 {
		if runesSectionInteger[0] == 0 { // if integer part starts with 0, there cannot be other digits after initial 0
			errorsCollected = append(errorsCollected, errors.New("digits after 0 in integer part: " + string(runesSectionInteger)))
		}
	}


	var digits09 = []rune("0123456789")
	if ! validate_runes_are_in_allowed_set(runesSectionInteger, digits09) {
		errorsCollected = append(errorsCollected, errors.New("illegal char in integer part: " + string(runesSectionInteger)))
	}

	if ! validate_runes_are_in_allowed_set(runesSectionFraction, digits09) {
		errorsCollected = append(errorsCollected, errors.New("illegal char in fraction part: " + string(runesSectionFraction)))
	}

	if len(runesSectionExponent) > 0 { // validate the first char in exponent section
		if ! validate_rune_are_in_allowed_set(runesSectionExponent[0], []rune{'+', '-'}) {
			errorsCollected = append(errorsCollected, errors.New("exponent part's first char is not +-: " + string(runesSectionExponent)))
		}
	}

	if len(runesSectionExponent) == 1 { // exponent section is too short
		errorsCollected = append(errorsCollected, errors.New("in exponent section +|- can be the FIRST char, then minimum one digit is necessary, and that is missing: " + string(runesSectionExponent)))
	}

	if len(runesSectionExponent) > 1 { // validate other chars in exponent section
		if ! validate_runes_are_in_allowed_set(runesSectionExponent[1:], digits09) {
			errorsCollected = append(errorsCollected, errors.New("illegal char after first char of exponent section: " + string(runesSectionExponent)))
		}
	}


	/////////////////////////// NUM CALCULATION, BASED ON DIGITS ////////////////////////////////
	thisIsValidNumber := lenErrorCollectorBeforeErrorDetection == len(errorsCollected)
	if thisIsValidNumber {

		// TODO: this section is too complicated, rewrite it.
		// cases: - only integer part,
		//        - int+fraction part,
		//        - int+exponent part
		//        - int+fraction+exponent part


		// ONLY INTEGER PART
		if len(runesSectionInteger) > 0 && len(runesSectionFraction) == 0 && len(runesSectionExponent) == 0 {
			numBase10, err := strconv.Atoi(string(runesSectionInteger));
			if err != nil {
				errorsCollected = append(errorsCollected, err)
			} else {
				token.ValNumberInt = numBase10
				token.Type = "number_integer"
			}
		}

		// ONLY INTEGER + FRACTION PART
		if len(runesSectionInteger) > 0 && len(runesSectionFraction) > 0 && len(runesSectionExponent) == 0 {
			numBase10, err := strconv.ParseFloat(string(runesSectionInteger), 64);
			if err != nil {
				errorsCollected = append(errorsCollected, err)
			} else {
				token.ValNumberFloat = numBase10
				token.Type = "number_float64"
			}
		}


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

// runesSections were checked against illegal chars, so here digitRune is in 0123456789
func digitIntegerValue(digit rune) int {
	unicode_code_point_zero_shift := '0' // '9' -> 9
	return int(digit - unicode_code_point_zero_shift)
}

// are the runes in the set?
func validate_runes_are_in_allowed_set(runes []rune, runesAllowed []rune) bool {
	for _, r := range runes {
		if ! validate_rune_are_in_allowed_set(r, runesAllowed) {
			return false
		}
	}
	return true
}


// is the rune in allowed set?
func validate_rune_are_in_allowed_set(runeValidated rune, runesAllowed []rune) bool {
	for _, r := range runesAllowed {
		if r == runeValidated {
			return true
		}
	}
	return false
}


// split once, at first occurance
func runes_split_at_pattern(runes []rune, splitterRune rune) ([]rune, []rune) {
	runesBefore := []rune{}
	runesAfter := []rune{}
	splitterDetected := false
	for _, r := range runes {
		if r == splitterRune {
			splitterDetected = true
			continue
		}
		if splitterDetected {
			runesAfter = append(runesAfter, r)
		} else {
			runesBefore = append(runesBefore, r)
		}
	}
	return runesBefore, runesAfter
}

// create a separated copy about original rune Slice
func runes_copy(runes []rune) []rune {
	runesNew := []rune{}
	for _, r := range runes {
		runesNew = append(runesNew, r)
	}
	return runesNew
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

	var tokenNow JsonValue

	for posInSrc, runeActual := range src {

		if runeActual == '"' {
			if !inStringDetection {
					tokenNow = JsonValue{Type: "string"}
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
	var tokenNow JsonValue

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
		} else { // save JsonValue, if something important is detected
			tokenNow = JsonValue{Type: detectedType}
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
			tokenNow := JsonValue{Type: detectedType}
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

		tokenNow := JsonValue{Type: "number"} // only numbers can be in the src now.
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


func TokensDisplay_startingCoords(tokens tokenTable_startPositionIndexed) {
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


func hexaRune_to_intVal(hexaChar rune) (int, error) {  // TESTED
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
/////////////////////// base functions /////////////////
