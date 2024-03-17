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


	globalTokensStartPositions = tokenTable_startPositionIndexed{}

	/*
	files := []string{"large-file_03percent.json", "large-file_06percent.json", "large-file_12percent.json", "large-file_25percent.json", "large-file_50percent.json", "large-file.json"}

	for _, file := range files {
		srcStr := file_read_for_tests(file)
		globalSrc = []rune(srcStr)
		globalErrorsCollected = []error{}

		start1 := time.Now()
		jsonDetect_strings()
		time_str := time.Since(start1)
		fmt.Println("time_str", time_str, file)

	}

	 */

	// 	srcStr := strings.Repeat(srcEverything, 100)
	// https://raw.githubusercontent.com/json-iterator/test-data/master/large-file.json
	srcStr := file_read_for_tests("large-file.json")
	// python3 json.loads() speed: 0.24469351768493652 sec
	// my speed: 3.82s (2024 Marc 16)



	globalSrc = []rune(srcStr)
	globalTokensStartPositions = tokenTable_startPositionIndexed{}
	globalErrorsCollected = []error{}

	start1 := time.Now()
	jsonDetect_strings()
	time_str := time.Since(start1)
	fmt.Println("time_str", time_str)

	startSep := time.Now()
	jsonDetect_separators()
	time_separators := time.Since(startSep)

	startBool := time.Now()
	jsonDetect_trueFalseNull()
	time_bool := time.Since(startBool)

	startNum := time.Now()
	jsonDetect_numbers()
	time_num := time.Since(startNum)

	startValid := time.Now()
	valueValidationsSettings_inTokens()
	time_valid := time.Since(startValid)

	startHierarchy := time.Now()
	elemRoot := objectHierarchyBuilding()
	time_hierarchy := time.Since(startHierarchy)



	fmt.Println("time_sep", time_separators)
	fmt.Println("time_bol", time_bool)
	fmt.Println("time_num", time_num)
	fmt.Println("time_val", time_valid)
	fmt.Println("time_hie", time_hierarchy)

	/* basic big result:
	time_str 185.342µs
	time_sep 116.447µs
	time_bol 18.02011ms
	time_num 18.412806ms
	time_val 443.194µs
	time_hie 1.073941ms

	# a little bigger json time result:
	time_str 1.152284ms
	time_sep 1.425252ms
	time_bol 1.433805656s
	time_num 1.121081367s
	time_val 3.130319ms
	time_hie 20.419823ms

	*/

	_ = elemRoot
}





//  go test -v -run Test_object_hierarchy_building
func Test_object_hierarchy_building(t *testing.T) {
	funName := "Test_object_hierarchy_building"
	testName := funName + "_basic"

	globalTokensStartPositions = tokenTable_startPositionIndexed{}
 	src := `{"embedded":{"level2": [3,4.5,"stringAtEnd"]}}`

	elemRoot, errorsCollected := JsonParse(src)
	jsonVal1, _ := elemRoot.ObjPathKeys([]string{"embedded", "level2"})
	jsonVal2, _ := elemRoot.ObjPath("/embedded/level2")
	compare_int_int(testName, 0, len(errorsCollected), t)
	compare_str_str(testName, jsonVal1.ValArray[2].repr(), jsonVal2.ValArray[2].repr(), t)

	valThird, _ := jsonVal1.Arr(2)
	compare_str_str(testName, `"stringAtEnd"`, valThird.repr(), t)
	fmt.Println(elemRoot.repr())

	elemRoot.AddKeyVal_into_object("newKey", NewString_JSON_value("newVal"))
	jsonVal3, _ := elemRoot.ObjPath("/newKey")
	compare_str_str(testName, `"newVal"`, jsonVal3.repr(), t)
	fmt.Println(elemRoot.repr(4))

	elemRoot.AddKeyVal_path_into_object("/embedded/level2", NewString_JSON_value("overwritten"))
	jsonVal4, _ := elemRoot.ObjPath("/embedded/level2")
	compare_str_str(testName, `"overwritten"`, jsonVal4.repr(), t)
	fmt.Println(elemRoot.repr(4))
}


////////////////////////////////////////////////////////////////////////////////////////////

func Test_parse_number_integer(t *testing.T) {
	funName := "Test_parse_number_integer"
	testName := funName + "_basic"
	globalTokensStartPositions = tokenTable_startPositionIndexed{}

	srcStr := `{"int":123, "float": 456.78, "intNegative": -9, "floatNegative": -0.12}`
	globalSrc= []rune(srcStr)
	globalErrorsCollected = []error{}

	jsonDetect_strings()
	jsonDetect_separators()
	jsonDetect_trueFalseNull()
	jsonDetect_numbers()

	valueValidationsSettings_inTokens()
	TokensDisplay_startingCoords(globalTokensStartPositions)

	compare_int_int(testName, 17,     len(globalTokensStartPositions),               t)
	compare_str_str(testName, "int",   globalTokensStartPositions[ 1].valString,      t)
	compare_int_int(testName, 123,    globalTokensStartPositions[ 7].valNumberInt,   t)
	compare_flt_flt(testName, 456.78, globalTokensStartPositions[21].valNumberFloat, t)
	compare_int_int(testName, -9,     globalTokensStartPositions[44].valNumberInt,   t)
	compare_flt_flt(testName, -0.12,  globalTokensStartPositions[65].valNumberFloat, t)
}

//  go test -v -run   Test_token_validate_and_value_set_for_strings
func Test_token_validate_and_value_set_for_strings(t *testing.T) {
	funName := "Test_token_validate_and_value_set_for_strings"

	testName := funName + "_escaped__quotes__reverseSolidus"
	globalTokensStartPositions = tokenTable_startPositionIndexed{}

	srcStr := `{"quote":"\"Assume a virtue, if you have it not.\"\nShakespeare", "source": "http:\/\/www.quotationspage.com\/quotes\/William_Shakespeare\/"}`
	globalSrc = []rune(srcStr)
	globalErrorsCollected = []error{}

	jsonDetect_strings()
	jsonDetect_separators()
	jsonDetect_trueFalseNull()
	jsonDetect_numbers()

	valueValidationsSettings_inTokens()
	// at this point, string tokens' real value is parsed - but there are no embedded structures yet
	// TokensDisplay_startingCoords(tokens)

	compare_int_int(testName, 9, len(globalTokensStartPositions), t)
	compare_str_str(testName, "quote",  globalTokensStartPositions[1].valString, t)
	compare_str_str(testName, `"Assume a virtue, if you have it not."`+"\nShakespeare",    globalTokensStartPositions[9].valString,  t)
	compare_str_str(testName, "http://www.quotationspage.com/quotes/William_Shakespeare/", globalTokensStartPositions[76].valString, t)



	////////////////////////////////////////////////////////////////////////////////////////////
	testName = funName + "_all"
	src2:=`{"quotation":      "\" text\"", 
            "reverseSolidus": "\\ reverseSolidus", 
            "solidus":        "\/ solidus", 
            "backspace":      "\b backspace", 
            "formFeed":       "\f formFeed", 
            "lineFeed":       "\n lineFeed", 
            "carriageReturn": "\r carriageReturn", 
            "horizontalTab":  "\t horizontalTab",
			"unicodeChar":    "\u0022"
			"unicodeChar2":   "\u00e4"
}`

	globalSrc = []rune(src2)
	globalTokensStartPositions = tokenTable_startPositionIndexed{}
	globalErrorsCollected = []error{}

	jsonDetect_strings()
	jsonDetect_separators()
	jsonDetect_trueFalseNull()
	jsonDetect_numbers()

	valueValidationsSettings_inTokens()
	// TokensDisplay_startingCoords(tokens)
	compare_str_str(testName, `" text"`,            globalTokensStartPositions[ 19].valString,  t)
	compare_str_str(testName, "\\ reverseSolidus",  globalTokensStartPositions[ 63].valString,  t)
	compare_str_str(testName, "/ solidus",          globalTokensStartPositions[115].valString, t)
	compare_str_str(testName, "\b backspace",       globalTokensStartPositions[160].valString, t)
	compare_str_str(testName, "\f formFeed",        globalTokensStartPositions[207].valString, t)
	compare_str_str(testName, "\n lineFeed",        globalTokensStartPositions[253].valString, t)

	compare_str_str(testName, "\r carriageReturn",  globalTokensStartPositions[299].valString, t)
	compare_str_str(testName, "\t horizontalTab",   globalTokensStartPositions[351].valString, t)
	compare_str_str(testName, "\"",                 globalTokensStartPositions[392].valString, t)
	compare_str_str(testName, "ä",                  globalTokensStartPositions[422].valString, t)
}








////////////////////////////////////////////////////////////////////////////////////////////
func Test_detect_numbers(t *testing.T) {
	funName := "Test_detect_numbers"
	testName := funName + "_basic"
	globalTokensStartPositions = tokenTable_startPositionIndexed{}


	srcStr := `{"age":123, "balance": -456.78, "problems": 0, "loan": -1.2E+3, "otherNum": 0.1e-4}`
	globalSrc = []rune(srcStr)
	globalErrorsCollected = []error{}

	jsonDetect_strings()
	jsonDetect_separators()
	jsonDetect_trueFalseNull()
	jsonDetect_numbers()
	// TokensDisplay_startingCoords(tokens)
	compare_int_int(testName, 21, len(globalTokensStartPositions), t)

	compare_str_str(testName, "123",     string(globalTokensStartPositions[7].runesInSrc), t)
	compare_str_str(testName, "-456.78", string(globalTokensStartPositions[23].runesInSrc), t)
	compare_str_str(testName, "0",       string(globalTokensStartPositions[44].runesInSrc), t)
	compare_str_str(testName, "-1.2E+3", string(globalTokensStartPositions[55].runesInSrc), t)
	compare_str_str(testName, "0.1e-4",  string(globalTokensStartPositions[76].runesInSrc), t)
}






func Test_true_false_null(t *testing.T) {
	funName := "Test_true_false_null"
	testName := funName + "_basic"

	srcStr := `{"name":"Bob","money":123,"boy":true,"girl":false,"age":null}`
	globalSrc = []rune(srcStr)
	srcLenOrig := len(globalSrc)

	globalErrorsCollected = []error{}
	globalTokensStartPositions = tokenTable_startPositionIndexed{}

	jsonDetect_strings()
	jsonDetect_separators()
	jsonDetect_trueFalseNull()

	// the orig src len has to be equal with the cleaned/received one's length:
	compare_int_int(testName, srcLenOrig, len(globalSrc), t)
	// TokensDisplay_startingCoords(tokens)
	compare_runes_runes(testName, []rune(`                      123                                    `), globalSrc, t)
	// compare_int_int(testName, 20 , len(tokens), t)
}

func Test_separators_detect(t *testing.T) {
	funName := "Test_separators_detect"

	testName := funName + "_basic"
	globalTokensStartPositions = tokenTable_startPositionIndexed{}

	srcStr := `{"students":[{"name":"Bob", "age":12}{"name": "Eve", "age":34.56}]}`
	globalSrc = []rune(srcStr)
	globalErrorsCollected = []error{}

	jsonDetect_separators()
	//                              `{"students":[{"name":"Bob", "age":12}{"name": "Eve", "age":34.56}]}`
	compare_runes_runes(testName, []rune(` "students"   "name" "Bob"  "age" 12  "name"  "Eve"  "age" 34.56   `), globalSrc, t)
	compare_int_int(testName, 15, len(globalTokensStartPositions), t)

	/* because the separators are one char long elems, the start position and end position
	   are ALWAYS same, and the length of Runes are 1, too. */
	testOneElem := func (srcWanted string, positionInSrc int) {
		tokenNow := globalTokensStartPositions[positionInSrc]
		compare_int_int(    testName, positionInSrc,         tokenNow.charPositionFirstInSourceCode,  t)
		compare_int_int(    testName, positionInSrc,         tokenNow.charPositionLastInSourceCode,   t)
		compare_int_int(    testName, 1,      len(tokenNow.runesInSrc), t)
		compare_runes_runes(testName, []rune(srcWanted),     tokenNow.runesInSrc,  t)
	}

	testOneElem("{", 0  )
	testOneElem(":", 11 )
	testOneElem("[", 12 )
	testOneElem("{", 13 )
	testOneElem(":", 20 )
	testOneElem(",", 26 )
	testOneElem(":", 33 )
	testOneElem("}", 36 )
	testOneElem("{", 37 )
	testOneElem(":", 44 )
	testOneElem(",", 51 )
	testOneElem(":", 58 )
	testOneElem("}", 64 )
	testOneElem("]", 65 )
	testOneElem("}", 66 )

	compare_int_int(testName, len(globalErrorsCollected), 0, t)
}

func Test_detect_strings(t *testing.T) {
	funName := "Test_detect_strings"


	////////////////////////////////////////////////////////////////////////////////////////////
	testName := funName + "_emptyString"
	srcStr := `{"empty":""}`
	globalSrc= []rune(srcStr)
	srcLenOrig := len(globalSrc)

	globalTokensStartPositions = tokenTable_startPositionIndexed{}
	globalErrorsCollected = []error{}


	jsonDetect_strings()
	// after token detection, the parsed section is removed;
	//                                       `{"empty":""}`, t)
	compare_runes_runes(testName, []rune(`{       :  }`), globalSrc, t)
	compare_int_int(testName, srcLenOrig, len(globalSrc), t)

	compare_int_int(testName, len(globalTokensStartPositions), 2, t) // 3 strings were detected
	compare_int_int(testName, 1, globalTokensStartPositions[1].charPositionFirstInSourceCode,  t)
	compare_int_int(testName, 7, globalTokensStartPositions[1].charPositionLastInSourceCode,  t)
	compare_runes_runes(testName, []rune(`"empty"`), globalTokensStartPositions[1].runesInSrc, t)
	compare_runes_runes(testName, []rune(`""`), globalTokensStartPositions[9].runesInSrc, t)

	compare_int_int(testName, len(globalErrorsCollected), 0, t)




	////////////////////////////////////////////////////////////////////////////////////////////
	testName = funName + "_simpleStringDetect"
	globalTokensStartPositions = tokenTable_startPositionIndexed{}

	srcStr = `{"name":"Bob", "age": 42}`
	globalSrc = []rune(srcStr)
	srcLenOrig = len(globalSrc)
	globalTokensStartPositions = tokenTable_startPositionIndexed{}
	globalErrorsCollected = []error{}

	// tokens are indexed by the first char where they were detected
	jsonDetect_strings()
	//                              `{"name":"Bob", "age": 42}`
	// after token detection, the parsed section is removed;
	compare_runes_runes(testName, []rune(`{      :     ,      : 42}`), globalSrc, t)
	compare_int_int(testName, srcLenOrig, len(globalSrc), t)

	compare_int_int(testName, 3, len(globalTokensStartPositions), t) // 3 strings were detected
	compare_int_int(testName, 1, globalTokensStartPositions[1].charPositionFirstInSourceCode,  t)
	compare_int_int(testName, 6, globalTokensStartPositions[1].charPositionLastInSourceCode, t)
	compare_runes_runes(testName, []rune(`"name"`), globalTokensStartPositions[1].runesInSrc, t)
	compare_int_int(testName, len(globalErrorsCollected), 0, t)


	////////////////////////////////////////////////////////////////////////////////////////////
	testName = funName + "_escape"
	globalTokensStartPositions = tokenTable_startPositionIndexed{}

	srcEscStr := `{"name \"of\" the \t\\\"rose\n\"":"red"}`
	globalSrc = []rune(srcEscStr)
	srcLenOrig = len(globalSrc)
	globalErrorsCollected = []error{}

	// tokens are indexed by the first char where they were detected
	jsonDetect_strings()

	//                              `{"name \"of\" the \t\\\"rose\n\"":"red"}`
	compare_runes_runes(testName, []rune(`{                                :     }`), globalSrc, t)
	compare_int_int(testName, srcLenOrig, len(globalSrc), t)
	compare_int_int(testName, 1, globalTokensStartPositions[1].charPositionFirstInSourceCode, t)
	compare_int_int(testName, 32, globalTokensStartPositions[1].charPositionLastInSourceCode, t)
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

// test/Debug Helper - display Tokens table
func TokensDisplay_startingCoords(tokens tokenTable_startPositionIndexed) {
	keys := local_tool__tokenTable_position_keys_sorted()

	fmt.Println("== Tokens Table display ==")
	for _, key := range keys{
		fmt.Println(string(tokens[key].runesInSrc), key, tokens[key])
	}
}

func file_read_for_tests(fn string) string {
	dat, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	return string(dat)
}