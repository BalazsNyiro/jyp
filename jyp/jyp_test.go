
package jyp

import (
	"fmt"
	"sort"
	"testing"
)


func Test_separators_detect(t *testing.T) {
	funName := "Test_separators_detect"

	testName := funName + "_basic"
	src := `{"students":[{"name":"Bob", "age":12}{"name": "Eve", "age":34.56}]}`
	tokensStartPositions := tokenTable_startPositionIndexed{}
	errorsCollected := []error{}

	srcSep, tokensSep, errorsCollectedSep := json_detect_separators(src, tokensStartPositions, errorsCollected)
	//                                       `{"students":[{"name":"Bob", "age":12}{"name": "Eve", "age":34.56}]}`
	compare_string_string(testName, ` "students"   "name" "Bob"  "age" 12  "name"  "Eve"  "age" 34.56   `, srcSep, t)
	compare_int_int(testName, 15, len(tokensSep), t)

	/* because the separators are one char long elems, the start position and end position
	   are ALWAYS same, and the length of runes are 1, too. */
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
	tokensStartPositions := tokenTable_startPositionIndexed{}
	errorsCollected := []error{}


	srcEmpty, tokensEmpty, errorsCollectedEmpty := json_detect_strings(src, tokensStartPositions, errorsCollected)
	// after token detection, the parsed section is removed;
	//                                       `{"empty":""}`, t)
	compare_string_string(testName, `{       :  }`, srcEmpty, t)

	compare_int_int(testName, len(tokensEmpty), 2, t) // 3 strings were detected
	compare_int_int(testName, 1, tokensEmpty[1].charPositionFirstInSourceCode,  t)
	compare_int_int(testName, 7, tokensEmpty[1].charPositionLastInSourceCode,  t)
	compare_runes_runes(testName, []rune(`"empty"`), tokensEmpty[1].runes, t)
	compare_runes_runes(testName, []rune(`""`), tokensEmpty[9].runes, t)

	compare_int_int(testName, len(errorsCollectedEmpty), 0, t)




	////////////////////////////////////////////////////////////////////////////////////////////
	testName = funName + "_simpleStringDetect"
	src = `{"name":"Bob", "age": 42}`
	tokensStartPositions = tokenTable_startPositionIndexed{}
	errorsCollected = []error{}

	// tokens are indexed by the first char where they were detected
	src2, tokens2, errorsCollected2 := json_detect_strings(src, tokensStartPositions, errorsCollected)
	//                                       `{"name":"Bob", "age": 42}`
	// after token detection, the parsed section is removed;
	compare_string_string(testName, `{      :     ,      : 42}`, src2, t)

	compare_int_int(testName, 3, len(tokens2), t)  // 3 strings were detected
	compare_int_int(testName, 1, tokens2[1].charPositionFirstInSourceCode,  t)
	compare_int_int(testName, 6, tokens2[1].charPositionLastInSourceCode, t)
	compare_runes_runes(testName, []rune(`"name"`), tokens2[1].runes, t)
	compare_int_int(testName, len(errorsCollected2), 0, t)


	////////////////////////////////////////////////////////////////////////////////////////////
	testName = funName + "_escape"
	srcEsc := `{"name \"of\" the \t\\\"rose\n\"":"red"}`
	print("escaped src:", srcEsc, "\n")
	tokensStartPositions = tokenTable_startPositionIndexed{}
	errorsCollected = []error{}

	// tokens are indexed by the first char where they were detected
	srcEsc, tokensEsc, errorsCollectedEsc := json_detect_strings(srcEsc, tokensStartPositions, errorsCollected)
	_ = tokensEsc
	_ = errorsCollectedEsc

	//                                       `{"name \"of\" the \t\\\"rose\n\"":"red"}`
	compare_string_string(testName, `{                                :     }`, srcEsc, t)
	compare_int_int(testName, 1, tokensEsc[1].charPositionFirstInSourceCode,  t)
	compare_int_int(testName, 32, tokensEsc[1].charPositionLastInSourceCode, t)
}

//////////////////////////// TEST BASE FUNCS ///////////////////
func compare_int_int(testName string, wantedNum int, received int, t *testing.T) {
	if wantedNum != received {
		t.Fatalf("\nError in %s wanted: %d, received: %d", testName, wantedNum, received)
	}
}

func compare_string_string(callerInfo, strWanted, strReceived string, t *testing.T) {
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

func tokensDisplay(tokens tokenTable_startPositionIndexed) {
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
