// BJP - Balazs' Json parser
// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

// this file is the implementation of the _standard_ json data format:
// https://www.json.org/json-en.html

package bjp

import (
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
	ValType string
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
}

// if the src can be parsed, return with the JSON root object with nested elems, and err is nil.
func parseSrc(src string) (Elem, error) {
	elemRoot := Elem{}




	return elemRoot, nil
}


////////////////////// BASE FUNCTIONS ///////////////////////////////////////////////
const message_no_more_char_in_src = "no_more_character_in_json_src"

// what is the next Json elem? object/array/string/ number, true, false, null / no_more_character_in_json_src
// and give back the position of the detected char
func detectNextOpenerTypeFromBeginning(src string) (string, int) {
	pos := -1
	for _, runeNow := range src {
		pos += 1

		if runeNow == '{' {return "object", pos}
		if runeNow == '[' {return "array" , pos}
		if runeNow == '"' {return "string", pos}

		// the numbers, or the true/false/null values can be incorrect,
		// the token's correctness is tested in the parse step

		// a number can be started with - or a digit (0, or non-zero-and-every-digits)
		if strings.Contains(digits_09, string(runeNow)) { return "number", pos}
		if runeNow == '-' {return "number", pos}

		if runeNow == 't' {return "true" , pos}
		if runeNow == 'f' {return "false", pos}
		if runeNow == 'n' {return "null" , pos}

		// we are NOT in a string, that is handled separatedly! a string can have whitespaces, too
		if strings.Contains(whitespaces, string(runeNow)) {
			continue
		}
	}
	return message_no_more_char_in_src, pos
}


// and give back the position of the detected char
func detectNextCloserTypeFromEnd(src string) (string, int) {
	lenSrc := len(src)
	if lenSrc == 0 {
		return message_no_more_char_in_src, -1 // no source code
	}

	lastPos := lenSrc - 1
	pos := lastPos
	for ; pos >= 0; pos-- {
		runeNow := src[pos]

		if runeNow == '}' {return "object", pos}
		if runeNow == ']' {return "array" , pos}
		if runeNow == '"' {return "string", pos}

		// a number always ends with a digit
		if strings.Contains(digits_09, string(runeNow)) { return "number", pos}

		// we are NOT in a string, that is handled separatedly! a string can have whitespaces, too
		if strings.Contains(whitespaces, string(runeNow)) {
			continue
		}
		if runeNow == 'e' { // maybe truE, maybe falsE
			posPrev := pos-1
			if posPrev >= 0 {
				runePrev := src[posPrev]
				if runePrev == 's' { //falSe, the 2nd char from back in falSe
					return "false" , pos
				}

				if runePrev == 'u' { //trUe, the 2nd char from back in trUe
					return "true" , pos
				}

			}
		}
		if runeNow == 'f' {return "false", pos}
		if runeNow == 'l' {return "null" , pos}
	}
	return message_no_more_char_in_src, pos
}



/////////////////////// base functions /////////////////
