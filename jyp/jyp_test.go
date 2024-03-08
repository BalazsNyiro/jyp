
package jyp

import (
	"fmt"
	"testing"
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



//  go test -v -run Test_object_hierarchy_building
func Test_object_hierarchy_building(t *testing.T) {
	funName := "Test_object_hierarchy_building"
	testName := funName + "_basic"

	src := `{"int":123, "float": 456.78, "intNegative": -9, "floatNegative": -0.12}`

	elemRoot, errorsCollected := JsonParse(src)
	fmt.Println("REPR:")
	print(elemRoot.repr(2))

	fmt.Println("=====================")
	fmt.Println("rootType", elemRoot.ValType)

	// TODO: the main root object and hierarchy building is in progress

	_ = elemRoot
	_ = errorsCollected
	_ = testName
}


////////////////////////////////////////////////////////////////////////////////////////////

func Test_parse_number_integer(t *testing.T) {
	funName := "Test_parse_number_integer"
	testName := funName + "_basic"

	src := `{"int":123, "float": 456.78, "intNegative": -9, "floatNegative": -0.12}`
	tokens := tokenTable_startPositionIndexed{}
	errorsCollected := []error{}

	src, tokens, errorsCollected = json_detect_strings________(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_separators_____(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_true_false_null(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_numbers________(src, tokens, errorsCollected)

	tokens, errorsCollected = tokens_validations_value_settings(tokens, errorsCollected)
	TokensDisplay_startingCoords(tokens)

	compare_int_int(testName, 17,     len(tokens),               t)
	compare_str_str(testName, "int",  tokens[ 1].valString,      t)
	compare_int_int(testName, 123,    tokens[ 7].valNumberInt,   t)
	compare_flt_flt(testName, 456.78, tokens[21].valNumberFloat, t)
	compare_int_int(testName, -9,     tokens[44].valNumberInt,   t)
	compare_flt_flt(testName, -0.12,  tokens[65].valNumberFloat, t)
}

//  go test -v -run   Test_token_validate_and_value_set_for_strings
func Test_token_validate_and_value_set_for_strings(t *testing.T) {
	funName := "Test_detect_numbers"

	testName := funName + "_escaped__quotes__reverseSolidus"

	src := `{"quote":"\"Assume a virtue, if you have it not.\"\nShakespeare", "source": "http:\/\/www.quotationspage.com\/quotes\/William_Shakespeare\/"}`
	tokens := tokenTable_startPositionIndexed{}
	errorsCollected := []error{}

	src, tokens, errorsCollected = json_detect_strings________(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_separators_____(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_true_false_null(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_numbers________(src, tokens, errorsCollected)

	tokens, errorsCollected = tokens_validations_value_settings(tokens, errorsCollected)
	// at this point, string tokens' real value is parsed - but there are no embedded structures yet
	// TokensDisplay_startingCoords(tokens)

	compare_int_int(testName, 9, len(tokens), t)
	compare_str_str(testName, "quote",  tokens[1].valString, t)
	compare_str_str(testName, `"Assume a virtue, if you have it not."`+"\nShakespeare",    tokens[9].valString,  t)
	compare_str_str(testName, "http://www.quotationspage.com/quotes/William_Shakespeare/", tokens[76].valString, t)



	////////////////////////////////////////////////////////////////////////////////////////////
	testName = funName + "_all"
	src = `{"quotation":      "\" text\"", 
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

	tokens = tokenTable_startPositionIndexed{}
	errorsCollected = []error{}

	src, tokens, errorsCollected = json_detect_strings________(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_separators_____(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_true_false_null(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_numbers________(src, tokens, errorsCollected)

	tokens, errorsCollected = tokens_validations_value_settings(tokens, errorsCollected)
	// TokensDisplay_startingCoords(tokens)
	compare_str_str(testName, `" text"`,            tokens[19].valString,  t)
	compare_str_str(testName, "\\ reverseSolidus",  tokens[63].valString,  t)
	compare_str_str(testName, "/ solidus",          tokens[115].valString, t)
	compare_str_str(testName, "\b backspace",       tokens[160].valString, t)
	compare_str_str(testName, "\f formFeed",        tokens[207].valString, t)
	compare_str_str(testName, "\n lineFeed",        tokens[253].valString, t)

	compare_str_str(testName, "\r carriageReturn",  tokens[299].valString, t)
	compare_str_str(testName, "\t horizontalTab",   tokens[351].valString, t)
	compare_str_str(testName, "\"",                 tokens[392].valString, t)
	compare_str_str(testName, "ä",                  tokens[422].valString, t)
}








////////////////////////////////////////////////////////////////////////////////////////////
func Test_detect_numbers(t *testing.T) {
	funName := "Test_detect_numbers"
	testName := funName + "_basic"


	src := `{"age":123, "balance": -456.78, "problems": 0, "loan": -1.2E+3, "otherNum": 0.1e-4}`
	tokens := tokenTable_startPositionIndexed{}
	errorsCollected := []error{}

	src, tokens, errorsCollected = json_detect_strings________(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_separators_____(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_true_false_null(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_numbers________(src, tokens, errorsCollected)
	// TokensDisplay_startingCoords(tokens)
	compare_int_int(testName, 21, len(tokens), t)

	compare_str_str(testName, "123",     string(tokens[7].runes), t)
	compare_str_str(testName, "-456.78", string(tokens[23].runes), t)
	compare_str_str(testName, "0",       string(tokens[44].runes), t)
	compare_str_str(testName, "-1.2E+3", string(tokens[55].runes), t)
	compare_str_str(testName, "0.1e-4",  string(tokens[76].runes), t)
}



func Test_src_get_words(t *testing.T) {
	funName := "Test_src_get_words"
	testName := funName + "_basic"

	whitepaceSeparatedString := "abc\t\n12 -1.2"
	words := src_get_whitespace_separated_words_posFirst_posLast(whitepaceSeparatedString)

	// how many words are detected?
	compare_int_int(testName, 3, len(words), t)

	compare_int_int(testName,    5, words[1].posFirst, t)
	compare_int_int(testName,    6, words[1].posLast,  t)
	compare_str_str(testName, "12", words[1].word,     t)

	compare_int_int(testName,      8, words[2].posFirst, t)
	compare_int_int(testName,     11, words[2].posLast,  t)
	compare_str_str(testName, "-1.2", words[2].word,     t)
}

func Test_src_get_char(t *testing.T) {
	funName := "Test_src_get_char"

	src:= "abc\t\n12 -1.2"

	testName := funName + "_overindexNegative"
	charSelected := src_get_char(src, -2)
	compare_rune_rune(testName, ' ', charSelected, t)

	testName = funName + "_overindexPositive"
	charSelected = src_get_char(src, 9999999)
	compare_rune_rune(testName, ' ', charSelected, t)

	testName = funName + "_whitespaceConversion"
	charSelected = src_get_char(src, 3)
	compare_rune_rune(testName, ' ', charSelected, t)

	testName = funName + "_normalSelection"
	charSelected = src_get_char(src, 2)
	compare_rune_rune(testName, 'c', charSelected, t)
}


func Test_is_whitespace_string_rune(t *testing.T) {
	funName := "Test_is_whitespace_string_rune"

	src := "abc\t\n12 -1.2"

	// whitespace string tests
	testName := funName + "_simpleStringWithWhitespaceEnding"
	isWhitespace := is_whitespace_string(src[0:5])
	compare_bool_bool(testName, false, isWhitespace, t)

	testName = funName + "_simpleStringOnlyWhitespace"
	isWhitespace = is_whitespace_string(src[3:5])
	compare_bool_bool(testName, true, isWhitespace, t)


	// whitespace rune tests
	testName = funName + "_simpleRuneWhitespace"
	runeSelected := rune(src[4])
	isWhitespace = is_whitespace_rune(runeSelected)
	compare_bool_bool(testName, true, isWhitespace, t)

	testName = funName + "_simpleRuneNonWhitespace"
	runeSelected = rune(src[6])
	isWhitespace = is_whitespace_rune(runeSelected)
	compare_bool_bool(testName, false, isWhitespace, t)
	compare_rune_rune(testName, '2', runeSelected, t)
}


func Test_true_false_null(t *testing.T) {
	funName := "Test_true_false_null"
	testName := funName + "_basic"

	src := `{"name":"Bob","money":123,"boy":true,"girl":false,"age":null}`
	srcLenOrig := len(src)

	tokens := tokenTable_startPositionIndexed{}
	errorsCollected := []error{}

	src, tokens, errorsCollected = json_detect_strings________(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_separators_____(src, tokens, errorsCollected)
	src, tokens, errorsCollected = json_detect_true_false_null(src, tokens, errorsCollected)

	// the orig src len has to be equal with the cleaned/received one's length:
	compare_int_int(testName, srcLenOrig, len(src), t)
	// TokensDisplay_startingCoords(tokens)
	compare_str_str(testName, `                      123                                    `, src, t)
	// compare_int_int(testName, 20 , len(tokens), t)

	_ = funName
	_ = testName
	_ = srcLenOrig
}

func Test_separators_detect(t *testing.T) {
	funName := "Test_separators_detect"

	testName := funName + "_basic"
	src := `{"students":[{"name":"Bob", "age":12}{"name": "Eve", "age":34.56}]}`
	tokensStartPositions := tokenTable_startPositionIndexed{}
	errorsCollected := []error{}

	srcSep, tokensSep, errorsCollectedSep := json_detect_separators_____(src, tokensStartPositions, errorsCollected)
	//                              `{"students":[{"name":"Bob", "age":12}{"name": "Eve", "age":34.56}]}`
	compare_str_str(testName, ` "students"   "name" "Bob"  "age" 12  "name"  "Eve"  "age" 34.56   `, srcSep, t)
	compare_int_int(testName, 15, len(tokensSep), t)

	/* because the separators are one char long elems, the start position and end position
	   are ALWAYS same, and the length of Runes are 1, too. */
	testOneElem := func (srcWanted string, positionInSrc int) {
		tokenNow := tokensSep[positionInSrc]
		compare_int_int(    testName, positionInSrc,         tokenNow.charPositionFirstInSourceCode,  t)
		compare_int_int(    testName, positionInSrc,         tokenNow.charPositionLastInSourceCode,   t)
		compare_int_int(    testName, 1,      len(tokenNow.runes), t)
		compare_runes_runes(testName, []rune(srcWanted),     tokenNow.runes,  t)
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

	compare_int_int(testName, len(errorsCollectedSep), 0, t)
}

func Test_detect_strings(t *testing.T) {
	funName := "Test_detect_strings"


	////////////////////////////////////////////////////////////////////////////////////////////
	testName := funName + "_emptyString"
	src := `{"empty":""}`
	srcLenOrig := len(src)

	tokensStartPositions := tokenTable_startPositionIndexed{}
	errorsCollected := []error{}


	srcEmpty, tokensEmpty, errorsCollectedEmpty := json_detect_strings________(src, tokensStartPositions, errorsCollected)
	// after token detection, the parsed section is removed;
	//                                       `{"empty":""}`, t)
	compare_str_str(testName, `{       :  }`, srcEmpty, t)
	compare_int_int(testName, srcLenOrig, len(srcEmpty), t)

	compare_int_int(testName, len(tokensEmpty), 2, t) // 3 strings were detected
	compare_int_int(testName, 1, tokensEmpty[1].charPositionFirstInSourceCode,  t)
	compare_int_int(testName, 7, tokensEmpty[1].charPositionLastInSourceCode,  t)
	compare_runes_runes(testName, []rune(`"empty"`), tokensEmpty[1].runes, t)
	compare_runes_runes(testName, []rune(`""`), tokensEmpty[9].runes, t)

	compare_int_int(testName, len(errorsCollectedEmpty), 0, t)




	////////////////////////////////////////////////////////////////////////////////////////////
	testName = funName + "_simpleStringDetect"
	src = `{"name":"Bob", "age": 42}`
	srcLenOrig = len(src)
	tokensStartPositions = tokenTable_startPositionIndexed{}
	errorsCollected = []error{}

	// tokens are indexed by the first char where they were detected
	src2, tokens2, errorsCollected2 := json_detect_strings________(src, tokensStartPositions, errorsCollected)
	//                              `{"name":"Bob", "age": 42}`
	// after token detection, the parsed section is removed;
	compare_str_str(testName, `{      :     ,      : 42}`, src2, t)
	compare_int_int(testName, srcLenOrig, len(src2), t)

	compare_int_int(testName, 3, len(tokens2), t) // 3 strings were detected
	compare_int_int(testName, 1, tokens2[1].charPositionFirstInSourceCode,  t)
	compare_int_int(testName, 6, tokens2[1].charPositionLastInSourceCode, t)
	compare_runes_runes(testName, []rune(`"name"`), tokens2[1].runes, t)
	compare_int_int(testName, len(errorsCollected2), 0, t)


	////////////////////////////////////////////////////////////////////////////////////////////
	testName = funName + "_escape"
	srcEsc := `{"name \"of\" the \t\\\"rose\n\"":"red"}`
	srcLenOrig = len(srcEsc)
	tokensStartPositions = tokenTable_startPositionIndexed{}
	errorsCollected = []error{}

	// tokens are indexed by the first char where they were detected
	srcEsc, tokensEsc, errorsCollectedEsc := json_detect_strings________(srcEsc, tokensStartPositions, errorsCollected)
	_ = tokensEsc
	_ = errorsCollectedEsc

	//                              `{"name \"of\" the \t\\\"rose\n\"":"red"}`
	compare_str_str(testName, `{                                :     }`, srcEsc, t)
	compare_int_int(testName, srcLenOrig, len(srcEsc), t)
	compare_int_int(testName, 1, tokensEsc[1].charPositionFirstInSourceCode, t)
	compare_int_int(testName, 32, tokensEsc[1].charPositionLastInSourceCode, t)
}


//  go test -v -run Test_hexaRune_to_intVal
func Test_hexaRune_to_intVal(t *testing.T) {
	funName := "Test_hexaRune_to_intVal"
	testName := funName + "hexa2int_conversation"

	intValDetected, err := hexaRune_to_intVal('b')
	compare_bool_bool(testName, true, err == nil, t)
	compare_int_int(testName, 11, intValDetected, t)

	intValDetected, err = hexaRune_to_intVal('m')
	compare_bool_bool(testName, true, err != nil, t)
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

