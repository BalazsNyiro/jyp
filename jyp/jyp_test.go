
package jyp

import (
	"fmt"
	"testing"
)

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

func compare_bool_bool(testName string, wanted bool, received bool, t *testing.T) {
	if wanted != received {
		t.Fatalf("\nError, different bool comparison %s wanted: %t, received: %t", testName, wanted, received)
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

func compare_rune_rune(callerInfo string, runeWanted, runeReceived rune, t *testing.T) {
	if runeWanted != runeReceived {
		errMsg := fmt.Sprintf("\nErr (%s) rune <>rune:\n  wanted -->>%s<<-- ??\nreceived -->>%s<<--\n\n", callerInfo, string(runeWanted), string(runeReceived))
		t.Fatalf(errMsg)
	}
}
