package jyp

import "testing"

func Test_string_detection_simple(t *testing.T) {
	elems_with_runes := elems_from_str(`"name of king"`)
	elems_strings_detected, _ := Json_string_find_in_elems__remove_spaces(elems_with_runes)
	wanted := elem{val_string: []rune("name of king"), val_type: "string"}
	check_elem__string_rune(elems_strings_detected[0], wanted, t)
}

func Test_string_detection_double(t *testing.T) {
	elems_with_runes := elems_from_str(`"\"name\"": "Bob", "age": 7`)
	elems_strings_detected, _ := Json_string_find_in_elems__remove_spaces(elems_with_runes)
	wanted := []elem{
		elem{val_string: []rune("\\\"name\\\""), val_type: "string"},
		elem{val_rune: ':', val_type: "rune"},
		elem{val_string: []rune("Bob"), val_type: "string"},
		elem{val_rune: ',', val_type: "rune"},
		elem{val_string: []rune("age"), val_type: "string"},
		elem{val_rune: ':', val_type: "rune"},
		elem{val_rune: '7', val_type: "rune"},
	}
	check_elems__string_rune(elems_strings_detected, wanted, t)
}
func check_elems__string_rune(receiveds []elem, wanteds []elem, t *testing.T) {
	if len(receiveds) != len(wanteds) {
		t.Fatalf(`len(received_elems) != len(wanted_elems)`)
	}
	for i := 0; i < len(receiveds); i++ {
		check_elem__string_rune(receiveds[i], wanteds[i], t)
	}
}

func check_elem__string_rune(received elem, wanted elem, t *testing.T) {
	if !runes_are_similar(received.val_string, wanted.val_string) {
		t.Fatalf("\nreceived: %v\n  wanted: %v, error",
			received.val_string, wanted.val_string)
	}
	if received.val_rune != wanted.val_rune {
		t.Fatalf(`received rune = %v, wanted %v, error`,
			received.val_rune,
			wanted.val_rune)
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
