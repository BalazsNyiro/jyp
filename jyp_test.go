// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package jyp

import (
	"fmt"
	"strconv"
	"testing"
)

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
		elem_string("name"),
		elem_rune(':'),
		elem_string("Bob"),
		elem_rune(','),
		elem_string("age"),
		elem_rune(':'),
		elem_rune('7'),
	}
	compare_receiveds_wanteds(elems_strings_detected, wanted, t)
}

func Test_string_detection_escaped_char(t *testing.T) {
	elems_with_runes := elems_from_str(`"he is \"Eduard\""`)
	elems_strings_detected := Json_collect_strings_in_elems__remove_spaces(elems_with_runes)
	wanted := elem_string("he is \\\"Eduard\\\"")
	compare_one_pair_received_wanted(elems_strings_detected[0], wanted, t)
}

func Test_number_int_detection(t *testing.T) {
	elems_with_runes := elems_from_str(`"price": 7.6, "age": 5`)
	elems_strings_detected := Json_collect_strings_in_elems__remove_spaces(elems_with_runes)
	elems_num_detected := Json_collect_numbers_in_elems(elems_strings_detected)
	elems_print(elems_num_detected, 0)
	wanted := []elem{
		elem_string("price"),
		elem_rune(':'),
		elem_number_float("7.6", 7.599999904632568),
		elem_rune(','),
		elem_string("age"),
		elem_rune(':'),
		elem_number_int(5),
	}
	compare_receiveds_wanteds(elems_num_detected, wanted, t)
}

func Test_scalar_detection(t *testing.T) {
	elems := elems_from_str(`"True": true, "False": false, "age": null `)
	elems = Json_collect_strings_in_elems__remove_spaces(elems)
	elems = Json_collect_scalars_in_elems(elems)
	// elems_print(elems, 0)
	wanted := []elem{
		elem_string("True"),
		elem_rune(':'),
		elem_true(),
		elem_rune(','),
		elem_string("False"),
		elem_rune(':'),
		elem_false(),
		elem_rune(','),
		elem_string("age"),
		elem_rune(':'),
		elem_null(),
	}
	compare_receiveds_wanteds(elems, wanted, t)
}

func Test_array_detection(t *testing.T) {
	elems := elems_from_str(`"name": "Bob", "scores": [4, 6], "friends": ["Eve", "Joe", 42], "key": "val"`)
	elems = Json_collect_strings_in_elems__remove_spaces(elems)
	elems = Json_collect_numbers_in_elems(elems)
	array := Json_collect_arrays_in_elems(elems)
	fmt.Println("arrays detected:")
	elems_print(array, 0)
	/*
		wanted := []elem{
			elem{valString: []rune("True"), valType: "string"},
			elem{valRune: ':', valType: "rune"},
			elem{valBool: true, valType: "bool"},
			elem{valRune: ',', valType: "rune"},
			elem{valString: []rune("False"), valType: "string"},
			elem{valRune: ':', valType: "rune"},
			elem{valBool: false, valType: "bool"},
			elem{valRune: ',', valType: "rune"},
			elem{valString: []rune("age"), valType: "string"},
			elem{valRune: ':', valType: "rune"},
			elem{valType: "null"},
		}

	*/
	// compare_receiveds_wanteds(elems, wanted, t)
}

////////////////////////////////////////////////////////////////////////

func elem_number_int(value int) elem {
	// return elem{valString: []rune("5"), valType: "number_int", valNumberInt: 5},
	return elem{valString: []rune(strconv.Itoa(value)), valType: "number_int", valNumberInt: value}
}

func elem_number_float(value_str_representation string, value_more_or_less_precise float64) elem {
	// elem{valString: []rune("7.6"), valType: "number_float", valNumberFloat: 7.599999904632568},
	return elem{valString: []rune(value_str_representation), valType: "number_float", valNumberFloat: value_more_or_less_precise}

}

func elem_string(value string) elem {
	// example:
	// elem{valString: []rune("age"), valType: "string"},
	return elem{valString: []rune(value), valType: "string"}
}
func elem_true() elem {
	return elem{valBool: true, valType: "bool"}
}

func elem_false() elem {
	return elem{valBool: false, valType: "bool"}
}

func elem_null() elem {
	return elem{valType: "null"}
}

func elem_rune(value rune) elem {
	// example:
	// elem{valRune: ':', valType: "rune"},
	return elem{valRune: value, valType: "rune"}
}

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
	if received.valBool != wanted.valBool {
		t.Fatalf(`received bool = %v, wanted %v, error`, received.valBool, wanted.valBool)
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
