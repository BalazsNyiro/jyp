package jyp

import "testing"

func Test_string_detection_simple(t *testing.T) {
	// this is a source code representation, so " is in the string:
	//                                     `"............"`
	// in the detected value, there is the content WITHOUT " signs
	elems_with_runes := elems_from_str(`"name of king"`)
	elems_strings_detected := Json_collect_strings_in_elems__remove_spaces(elems_with_runes)
	wanted := elem{valString: []rune("name of king"), valType: "string"}
	compare_one_pair_received_wanted(elems_strings_detected[0], wanted, t)
}

func Test_string_detection_double(t *testing.T) {
	elems_with_runes := elems_from_str(`"name": "Bob", "age": 7`)
	elems_strings_detected := Json_collect_strings_in_elems__remove_spaces(elems_with_runes)
	wanted := []elem{
		elem{valString: []rune("name"), valType: "string"},
		elem{valRune: ':', valType: "rune"},
		elem{valString: []rune("Bob"), valType: "string"},
		elem{valRune: ',', valType: "rune"},
		elem{valString: []rune("age"), valType: "string"},
		elem{valRune: ':', valType: "rune"},
		elem{valRune: '7', valType: "rune"},
	}
	compare_receiveds_wanteds(elems_strings_detected, wanted, t)
}

func Test_string_detection_escaped_char(t *testing.T) {
	elems_with_runes := elems_from_str(`"he is \"Eduard\""`)
	elems_strings_detected := Json_collect_strings_in_elems__remove_spaces(elems_with_runes)
	wanted := elem{valString: []rune("he is \\\"Eduard\\\""), valType: "string"}
	compare_one_pair_received_wanted(elems_strings_detected[0], wanted, t)
}

func Test_number_int_detection(t *testing.T) {
	elems_with_runes := elems_from_str(`"price": 7.6, "age": 5`)
	elems_strings_detected := Json_collect_strings_in_elems__remove_spaces(elems_with_runes)
	elems_num_detected := Json_collect_numbers_in_elems(elems_strings_detected)
	elems_print(elems_num_detected)
	wanted := []elem{
		elem{valString: []rune("price"), valType: "string"},
		elem{valRune: ':', valType: "rune"},
		elem{valString: []rune("7.6"), valType: "number_float", valNumberFloat: 7.599999904632568},
		elem{valRune: ',', valType: "rune"},
		elem{valString: []rune("age"), valType: "string"},
		elem{valRune: ':', valType: "rune"},
		elem{valString: []rune("5"), valType: "number_int", valNumberInt: 5},
	}
	compare_receiveds_wanteds(elems_num_detected, wanted, t)
}

func Test_scalar_detection(t *testing.T) {
	elems := elems_from_str(`"age": null, "True": true, "False": false`)
	elems = Json_collect_strings_in_elems__remove_spaces(elems)
	elems = Json_collect_scalars_in_elems(elems)
	elems_print(elems)
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
	if !runes_are_similar(received.valString, wanted.valString) {
		t.Fatalf("\nreceived: %v\n  wanted: %v, error",
			received.valString, wanted.valString)
	}
	if received.valRune != wanted.valRune {
		t.Fatalf(`received rune = %v, wanted %v, error`,
			received.valRune,
			wanted.valRune)
	}
	if received.valNumberInt != wanted.valNumberInt {
		t.Fatalf(`received int = %v, wanted %v, error`, received.valNumberInt, wanted.valNumberInt)
	}
	if received.valNumberFloat != wanted.valNumberFloat {
		t.Fatalf(`received float= %v, wanted %v, error`, received.valNumberFloat, wanted.valNumberFloat)
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
