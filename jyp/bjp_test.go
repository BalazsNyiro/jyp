
package bjp

import (
	"fmt"
	"testing"
)

func Test_detect_next_opener(t *testing.T) {
	funName := "Test_detect_next_opener-closer_"

	src := "\n\t\r { }\r\n\t "

	testName := funName + "object"
	nextOpenerDetectedType, posOpener := detectNextOpenerTypeFromBeginning(src)
	nextCloserDetectedType, posCloser := detectNextCloserTypeFromEnd(src)

	compare_string_string(testName, nextOpenerDetectedType, "object", t)
	compare_string_string(testName, nextCloserDetectedType, "object", t)
	compare_int_int(testName, posOpener, 4, t)
	compare_int_int(testName, posCloser, 6, t)


	src = " 123 \r\n\t "
	testName = funName + "number"
	nextOpenerDetectedType, posOpener = detectNextOpenerTypeFromBeginning(src)
	nextCloserDetectedType, posCloser = detectNextCloserTypeFromEnd(src)

	compare_string_string(testName, nextOpenerDetectedType, "number", t)
	compare_string_string(testName, nextCloserDetectedType, "number", t)
	compare_int_int(testName, posOpener, 1, t)
	compare_int_int(testName, posCloser, 3, t)


	src = " \"text\" \r\n\t "
	testName = funName + "text"
	nextOpenerDetectedType, posOpener = detectNextOpenerTypeFromBeginning(src)
	nextCloserDetectedType, posCloser = detectNextCloserTypeFromEnd(src)

	compare_string_string(testName, nextOpenerDetectedType, "string", t)
	compare_string_string(testName, nextCloserDetectedType, "string", t)
	compare_int_int(testName, posOpener, 1, t)
	compare_int_int(testName, posCloser, 6, t)


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
