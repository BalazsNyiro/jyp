// JYP - Json/Yaml Parser
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
	wanted := elem{valString: "name of king", valType: "string"}
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
	compare_receivedElems_wantedElems(elems_strings_detected, wanted, t)
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
	compare_receivedElems_wantedElems(elems_num_detected, wanted, t)
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
	compare_receivedElems_wantedElems(elems, wanted, t)
}

/*
 VERY IMPORTANT: in array detection I can't remove `,` runes, because
 `{` and `}` runes can be there, too, and they have their own `,` runes, too
 after object detection, I can return to the lists to remove the raw, unused runes

So in these test, they are not completely correct jsons, they are in an evolutionary level,
as the processing takes step by step, the result will be more better and better
*/
// here we don't use object detection, only scalars and array
func Test_array_detection(t *testing.T) {
	elems := elems_from_str(`"name": "Bob", "scores": [4, 6], "friends": [["Eve", 16], ["Joe", 42]], "key": "val"`)
	elems = Json_collect_strings_in_elems__remove_spaces(elems)
	elems = Json_collect_numbers_in_elems(elems)
	array := Json_collect_arrays_in_elems(elems)
	fmt.Println("(1) arrays detected:", len(array))
	elems_print(array, 0)

	wanted := []elem{
		elem_string("name"),
		elem_rune(':'),
		elem_string("Bob"),
		elem_rune(','),
		elem_string("scores"),
		elem_rune(':'),
		elem_array(
			[]elem{
				elem_number_int(4),
				elem_rune(','),
				elem_number_int(6),
			},
		),
		elem_rune(','),
		elem_string("friends"),
		elem_rune(':'),
		elem_array([]elem{
			elem_array([]elem{
				elem_string("Eve"),
				elem_rune(','),
				elem_number_int(16),
			}),
			elem_rune(','),
			elem_array([]elem{
				elem_string("Joe"),
				elem_rune(','),
				elem_number_int(42),
			}),
		},
		),
		elem_rune(','),
		elem_string("key"),
		elem_rune(':'),
		elem_string("val"),
	}
	compare_receivedElems_wantedElems(array, wanted, t)
}

//////////////// COMPLETE JSON TESTS ////////////////////////////////////

func Test_object_detection(t *testing.T) {
	elems, _ := Json_parse(`{"personal":{"city":"Paris", "cell": 123}}`)
	fmt.Println("Test_object_detection")
	elems_print(elems, 0)
	wanted := elem_object(map[string]elem{
		"personal": elem_object(map[string]elem{
			"city": elem_string("Paris"),
			"cell": elem_number_int(123),
		}),
	})
	compare_one_pair_received_wanted(elems[0], wanted, t)
}

// elem_root, _ := Json_parse(`{"name": "Bob", "friends": [{"name":Eve", "scores":[1,2]}, {"name":Bob", "scores":[3,4]}]})`)
func Test_json_1(t *testing.T) {
	elems, _ := Json_parse(`{"name": "Bob", "friends": [ {"name":"Eve"} ]}`)
	fmt.Println(" Test_json_1")
	elems_print(elems, 0)
	// wanted := elem_object(map[string]elem{
	// 	"personal": elem_object(map[string]elem{
	// 		"city": elem_string("Paris"),
	// 		"cell": elem_number_int(123),
	// 	}),
	// })
	// compare_one_pair_received_wanted(elem_root, wanted, t)
}
func Test_complex_big(t *testing.T) {
	return
	elems, _ := Json_parse(`{"name": "Bob", "friends": [ {"name":Eve", "scores":[1,2]}, {"name":Bob", "scores":[3,4]} ]})`)
	fmt.Println("Test_complex_big")
	elems_print(elems, 0)
	// wanted := elem_object(map[string]elem{
	// 	"personal": elem_object(map[string]elem{
	// 		"city": elem_string("Paris"),
	// 		"cell": elem_number_int(123),
	// 	}),
	// })
	// compare_one_pair_received_wanted(elem_root, wanted, t)
}

////////////////////////////////////////////////////////////////////////
func elem_object(values map[string]elem) elem {
	return elem{valObject: values, valType: "object"}
}

func elem_array(values []elem) elem {
	return elem{valArray: elems_copy(values, 0, len(values)), valType: "array"}
}

func elem_number_int(value int) elem {
	// return elem{valString: "5", valType: "number_int", valNumberInt: 5},
	return elem{valString: strconv.Itoa(value), valType: "number_int", valNumberInt: value}
}

func elem_number_float(value_str_representation string, value_more_or_less_precise float64) elem {
	// elem{valString: "7.6", valType: "number_float", valNumberFloat: 7.599999904632568},
	return elem{valString: value_str_representation, valType: "number_float", valNumberFloat: value_more_or_less_precise}
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

/////////////////////////////////////////////////////////////////////////

func compare_receivedElems_wantedElems(receiveds []elem, wanteds []elem, t *testing.T) {
	// compare only the lenth of the top level!
	if len(receiveds) != len(wanteds) {
		fmt.Println(">>> === compare, len !=   ===========")
		elems_print(receiveds, 0)
		fmt.Println("    .......................    ")
		elems_print(wanteds, 0)
		fmt.Println("<<< === compare, len !=   ===========")
		t.Fatalf(`len(received_elems %v) != len(wanted_elems %v) `,
			len(receiveds), len(wanteds))

	}
	for i := 0; i < len(receiveds); i++ {
		compare_one_pair_received_wanted(receiveds[i], wanteds[i], t)
	}
}

func _compare_two_objects(objA elem, objB elem, t *testing.T) {
	for keyReceived, _ := range objA.valObject {
		if Obj_has_key(objB.valObject, keyReceived) == false {
			t.Fatalf(`wanted object doesn't have key' %v error`, keyReceived)
		}
		compare_one_pair_received_wanted(objA.valObject[keyReceived], objB.valObject[keyReceived], t)
	}
}

func compare_one_pair_received_wanted(received elem, wanted elem, t *testing.T) {
	if received.valString != wanted.valString {
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

	if received.valType == "array" {
		compare_receivedElems_wantedElems(received.valArray, wanted.valArray, t)
	}

	if received.valType == "object" {
		_compare_two_objects(received, wanted, t) // check based on received object
		_compare_two_objects(wanted, received, t) // check based on wanted object (from other direction)
	}
}
