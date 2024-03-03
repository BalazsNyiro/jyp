
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

	srcSep, tokensSep, errorsCollectedSep := json_separators_detect(src, tokensStartPositions, errorsCollected)
	//                                       `{"students":[{"name":"Bob", "age":12}{"name": "Eve", "age":34.56}]}`
	compare_string_string(testName, ` "students"   "name" "Bob"  "age" 12  "name"  "Eve"  "age" 34.56   `, srcSep, t)

	tokensDisplay(tokensSep)

	compare_int_int(testName, 15, len(tokensSep), t)

	/* because the separators are one char long elems, the start position and end position
	   are ALWAYS same, and the length of runes are 1, too.

	*/
	testOneElem := func (srcWanted string, positionInSrc int, tokensOneTest tokenTable_startPositionIndexed) {
		tokenNow := tokensOneTest[positionInSrc]
		compare_int_int(    testName, positionInSrc,         tokenNow.charPositionFirstInSourceCode,  t)
		compare_int_int(    testName, positionInSrc,         tokenNow.charPositionLastInSourceCode,   t)
		compare_int_int(    testName, 1,      len(tokenNow.runes), t)
		compare_runes_runes(testName, []rune(srcWanted),     tokenNow.runes,  t)

	}

	// testOneElem("{", 0, tokensSep)
	_ = testOneElem

	compare_int_int(testName, 0, tokensSep[0].charPositionFirstInSourceCode,  t)
	compare_int_int(testName, 0, tokensSep[0].charPositionLastInSourceCode,   t)
	compare_int_int(testName, 1, len(tokensSep[0].runes),   t)
	compare_runes_runes(testName, []rune("{"), tokensSep[0].runes, t)

	compare_int_int(testName, 11, tokensSep[11].charPositionFirstInSourceCode,  t)
	compare_int_int(testName, 11, tokensSep[11].charPositionLastInSourceCode,   t)
	compare_int_int(testName, 1, len(tokensSep[11].runes),   t)
	compare_runes_runes(testName, []rune(":"), tokensSep[11].runes, t)

	compare_int_int(testName, 12, tokensSep[12].charPositionFirstInSourceCode,  t)
	compare_int_int(testName, 12, tokensSep[12].charPositionLastInSourceCode,   t)
	compare_int_int(testName, 1, len(tokensSep[0].runes),   t)
	compare_runes_runes(testName, []rune("{"), tokensSep[0].runes, t)

	compare_int_int(testName, len(errorsCollectedSep), 0, t)
}

func Test_detect_strings(t *testing.T) {
	funName := "Test_detect_strings"


	////////////////////////////////////////////////////////////////////////////////////////////
	testName := funName + "_emptyString"
	src := `{"empty":""}`
	tokensStartPositions := tokenTable_startPositionIndexed{}
	errorsCollected := []error{}


	srcEmpty, tokensEmpty, errorsCollectedEmpty := json_string_detect(src, tokensStartPositions, errorsCollected)
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
	src2, tokens2, errorsCollected2 := json_string_detect(src, tokensStartPositions, errorsCollected)
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
	srcEsc, tokensEsc, errorsCollectedEsc := json_string_detect(srcEsc, tokensStartPositions, errorsCollected)
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
		fmt.Println(key, tokens[key])
	}
}
