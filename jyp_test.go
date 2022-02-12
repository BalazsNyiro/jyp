package jyp

import "testing"

func Test_string_detection_simple(t *testing.T) {
	elems_with_runes := elems_from_str(`"name of king"`)
	elems_strings_detected, _ := Json_string_finder_in_elems(elems_with_runes)
	wanted := elem{val_rune: 'n', val_string: []rune("name of king"), val_type: "string"}
	result_check_val_string(elems_strings_detected[0], wanted, t)
}

func result_check_val_string(value_received elem, value_wanted elem, t *testing.T) {
	if !runes_are_similar(value_received.val_string, value_wanted.val_string) {
		t.Fatalf(`received = %v, want %v, error`, value_received, value_wanted)
	}
}

func runes_are_similar(runes1 []rune, runes2 []rune) bool {
	if len(runes1) != len(runes2) {
		return false
	}
	if len(runes1) == 0 {
		return true
	}
	for i := 0; i < len(runes1); i++ {
		if runes1[i] != runes2[i] {
			return false
		}
	}
	return true
}
