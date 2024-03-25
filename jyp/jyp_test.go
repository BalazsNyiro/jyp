/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

TODO, tests:
non-closed string error detection
incorrect numbers test (start with 0)
incorrect atoms test (something else than true/false/null)
*/

package jyp

import (
	"fmt"
	"os"
	"testing"
	"time"
	"unicode/utf8"
)

var srcEverything string = `{
    "whatIsThis": "a global, general json structure that has everything - the final test",

    "trueKey"  : true,
    "falseKey" : false,
    "nullKey"  : null,
    "stringKey": "str",

    "numbers": {
        "zeroPos":  0,
        "zeroNeg": -0,
        "intPos":   8,
        "intNeg":  -8,

        "what is this 0":  "zero started nums are different cases in Json spec",
        "floatPosZero":  0.1023,
        "floatNegZero": -0.4056,

        "what is this 1":  "non-zero started nums are different cases in Json spec",
        "floatPosOne":   1.809 ,
        "floatNegOne":  -1.5067,


        "what is this 2": "exponent without fraction",
        "exp_minus_e_plus":   -4e+3,
        "exp_minus_e_minus":  -4e-4,
        "exp_plus_e_plus":     4e+5,
        "exp_plus_e_minus":    4e-6,

        "exp_minus_E_plus":   -4E+3,
        "exp_minus_E_minus":  -4E-4,
        "exp_plus_E_plus":     4E+5,
        "exp_plus_E_minus":    4E-6,

        "what is this 3": "exponent with fraction",
        "exp_minus_e_plus__fract":  -1.1e+3,
        "exp_minus_e_minus_fract":  -2.2e-4,
        "exp_plus_e_plus___fract":   3.3e+5,
        "exp_plus_e_minus__fract":   4.4e-6,

        "exp_minus_E_plus__fract":  -5.5E+3,
        "exp_minus_E_minus_fract":  -6.6E-4,
        "exp_plus_E_plus___fract":   7.7E+5,
        "exp_plus_E_minus__fract":   8.8E-6,

        "what is this 4": "zero exponents",
        "exp_minus_e_plus__zero":  -1.1e+0,
        "exp_plus__e_plus__zero":   2.1e+0,
        "exp_minus_e_minus_zero":  -1.1e-0,
        "exp_plus__e_minus_zero":   2.1e-0,
    },

    "array_with_everything": ["str", -2, 0, 3, -4.5, -0, 6.7,
                              {"obj_in_array": ["embeddedArr":   null,
                                                "embeddedTrue":  true,
                                                "embeddedFalse": false,
                                               ]}
    ],

    "stringsAllPossibleOption": [
        "usedSource":             "https://www.json.org/json-en.html",
        "simple":                 "text",

        "quotation_mark": `     + "quote: \"wisdom\"" + `,
        "reverse_solidus": `    + "reversed: '\\' "   + `,

		"solidusExplanationUrl":  "https://groups.google.com/g/opensocial-and-gadgets-spec/c/FkLsC-2blbo?pli=1",
		"solidusExplanation":     "http:\/\/example.org" is the right way to encode a URL in a JSON string",
        "solidus":                "solidus:  http:\/\/example.org",

        "backspace": `          + "a\bb" + `,
        "formfeed": `           + "a\fb" + `,
        "linefeed": `           + "a\nb" + `,
        "carriage_return": `    + "a\rb" + `,
        "horizontal_tab": `     + "a\tb" + `,
        "4hex digits":            "quotation mark digit: \u0022",
    ]
}`

// TODO: test json src with errors!



//  go test -v -run Test_speed
func Test_speed(t *testing.T) {
	funName := "Test_speed"
	testName := funName + "_basic"
	_ = testName

	/*
		files := []string{"large-file_03percent.json", "large-file_06percent.json", "large-file_12percent.json", "large-file_25percent.json", "large-file_50percent.json", "large-file.json"}

		for _, file := range files {
			srcStr := file_read_to_string(file)
			src := []rune(srcStr)
			tokens := tokenTable_startPositionIndexed{}
			errorsCollected := []error{}

			start1 := time.Now()
			jsonDetect_strings______(src, tokens, errorsCollected)
			time_str := time.Since(start1)
			fmt.Println("time_str", time_str, file)

		}

	*/

	// 	srcStr := strings.Repeat(srcEverything, 100)
	// https://raw.githubusercontent.com/json-iterator/test-data/master/large-file.json

	timeReadFileStart := time.Now()
	srcStr := file_read_to_string("large-file.json")
	fmt.Println("time read file to string:", time.Since(timeReadFileStart))

	timeSimpleStringPassing := time.Now()
	tokensTableB := stepA__tokensTableDetect_structuralTokens_strings_L1(srcStr)
	fmt.Println("time tokensTableDetect structuralTokens:", time.Since(timeSimpleStringPassing))

	timeStructure := time.Now()
	errorsCollected := stepB__JSON_B_validation_L1(tokensTableB)
	root, _ := stepC__JSON_B_structure_building__L1(srcStr, tokensTableB, 0, errorsCollected)
	fmt.Println("time structure:", time.Since(timeStructure))
	_ = root

	// python3 json.loads() speed: 0.24469351768493652 sec
	// my speed: 3.82s (2024 Marc 16)
	//           3.47s (2024 Marc 17)
	//           1.56s (2024 Marc 25)

}




//////////////////////////// TEST BASE FUNCS ///////////////////
func compare_int_int(testName string, wantedNum, received int, t *testing.T) {
	if wantedNum != received {
		t.Fatalf("\nError in %s wanted: %d, received: %d", testName, wantedNum, received)
	}
}

func compare_flt_flt(testName string, wantedNum, received float64, t *testing.T) {
	if wantedNum != received {
		t.Fatalf("\nError in %s wanted: %f, received: %f", testName, wantedNum, received)
	}
}


func compare_bool_bool(testName string, wanted, received bool, t *testing.T) {
	if wanted != received {
		t.Fatalf("\nError, different bool comparison %s wanted: %t, received: %t", testName, wanted, received)
	}
}

func compare_str_str(callerInfo, strWanted, strReceived string, t *testing.T) {
	if strWanted != strReceived {
		t.Fatalf("\nErr String difference (%s):\n  wanted -->>%s<<-- ??\nreceived -->>%s<<--\n\n", callerInfo, strWanted, strReceived)
	}
}

func compare_runes_runes(callerInfo string, runesWanted, runesReceived []rune, t *testing.T) {
	errMsg := fmt.Sprintf("\nErr (%s) []rune <>[]rune:\n  wanted -->>%s<<-- ??\nreceived -->>%s<<--\n\n", callerInfo, string(runesWanted), string(runesReceived))
	if len(runesWanted) != len(runesReceived) {
		t.Fatalf(errMsg)
		return
	}

	for pos, runeWanted:= range runesWanted {
		if runeWanted != runesReceived[pos] {
			t.Fatalf(errMsg)
			return
		}
	}
}

func compare_rune_rune(callerInfo string, runeWanted, runeReceived rune, t *testing.T) {
	if runeWanted != runeReceived {
		errMsg := fmt.Sprintf("\nErr (%s) rune <>rune:\n  wanted -->>%s<<-- ??\nreceived -->>%s<<--\n\n", callerInfo, string(runeWanted), string(runeReceived))
		t.Fatalf(errMsg)
	}
}


func file_read_to_string(fn string) string {
	dat, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	return string(dat)
}

func file_read_to_runes(fn string) []rune {
	bytes, err := os.ReadFile(fn)
	if err != nil { panic(err) }
	runes := []rune{}
	// https://pkg.go.dev/unicode/utf8#DecodeRune
	for len(bytes) > 0 {
		r, size := utf8.DecodeRune(bytes)
		bytes = bytes[size:]
		runes = append(runes, r)
	}
	return runes
}
