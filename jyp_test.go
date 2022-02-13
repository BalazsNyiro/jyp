package jyp

import "testing"

func Test_string_detection_simple(t *testing.T) {
	// this is a source code representation, so " is in the string:
	//                                     `"............"`
	// in the detected value, there is the content WITHOUT " signs
	elems_with_runes := elems_from_str(`"name of king"`)
	elems_strings_detected := Json_collect_strings_in_elems__remove_spaces(elems_with_runes)
	wanted := elem{val_string: []rune("name of king"), val_type: "string"}
	compare_one_pair_received_wanted(elems_strings_detected[0], wanted, t)
}

func Test_string_detection_double(t *testing.T) {
	elems_with_runes := elems_from_str(`"name": "Bob", "age": 7`)
	elems_strings_detected := Json_collect_strings_in_elems__remove_spaces(elems_with_runes)
	wanted := []elem{
		elem{val_string: []rune("name"), val_type: "string"},
		elem{val_rune: ':', val_type: "rune"},
		elem{val_string: []rune("Bob"), val_type: "string"},
		elem{val_rune: ',', val_type: "rune"},
		elem{val_string: []rune("age"), val_type: "string"},
		elem{val_rune: ':', val_type: "rune"},
		elem{val_rune: '7', val_type: "rune"},
	}
	compare_receiveds_wanteds(elems_strings_detected, wanted, t)
}

func Test_string_detection_escaped_char(t *testing.T) {
	elems_with_runes := elems_from_str(`"he is \"Eduard\""`)
	elems_strings_detected := Json_collect_strings_in_elems__remove_spaces(elems_with_runes)
	wanted := elem{val_string: []rune("he is \\\"Eduard\\\""), val_type: "string"}
	compare_one_pair_received_wanted(elems_strings_detected[0], wanted, t)
}

func Test_number_int_detection(t *testing.T) {
	elems_with_runes := elems_from_str(`"price": 7.6, "age": 5`)
	elems_strings_detected := Json_collect_strings_in_elems__remove_spaces(elems_with_runes)
	elems_num_detected := Json_collect_numbers_in_elems(elems_strings_detected)
	elems_print(elems_num_detected)
	wanted := []elem{
		elem{val_string: []rune("price"), val_type: "string"},
		elem{val_rune: ':', val_type: "rune"},
		elem{val_string: []rune("7.6"), val_type: "number_float", val_number_float: 7.599999904632568},
		elem{val_rune: ',', val_type: "rune"},
		elem{val_string: []rune("age"), val_type: "string"},
		elem{val_rune: ':', val_type: "rune"},
		elem{val_string: []rune("5"), val_type: "number_int", val_number_int: 5},
	}
	compare_receiveds_wanteds(elems_num_detected, wanted, t)
}

////////////////////////////////////////////////////////////////////////

func compare_receiveds_wanteds(receiveds []elem, wanteds []elem, t *testing.T) {
	if len(receiveds) != len(wanteds) {
		t.Fatalf(`len(received_elems) != len(wanted_elems)`)
	}
	for i := 0; i < len(receiveds); i++ {
		compare_one_pair_received_wanted(receiveds[i], wanteds[i], t)
	}
}

func compare_one_pair_received_wanted(received elem, wanted elem, t *testing.T) {
	if !runes_are_similar(received.val_string, wanted.val_string) {
		t.Fatalf("\nreceived: %v\n  wanted: %v, error",
			received.val_string, wanted.val_string)
	}
	if received.val_rune != wanted.val_rune {
		t.Fatalf(`received rune = %v, wanted %v, error`,
			received.val_rune,
			wanted.val_rune)
	}
	if received.val_number_int != wanted.val_number_int {
		t.Fatalf(`received int = %v, wanted %v, error`, received.val_number_int, wanted.val_number_int)
	}
	if received.val_number_float != wanted.val_number_float {
		t.Fatalf(`received float= %v, wanted %v, error`, received.val_number_float, wanted.val_number_float)
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
