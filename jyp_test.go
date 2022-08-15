// JYP - Json/Yaml Parser
// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package jyp

import (
	"fmt"
	"testing"
)

func Test_string_detection_simple(t *testing.T) {
	// this is a source code representation, so " is in the string:
	//                                     `"............"`
	// in the detected value, there is the content WITHOUT " signs
	elems_with_runes := elem_runes_from_str(`"name of king"`)
	elems_strings_detected := Json_collect_strings_in_elems__remove_spaces(elems_with_runes)
	wanted := Elem{ValString: "name of king", ValType: "string"}
	compare_one_pair_received_wanted(elems_strings_detected[0], wanted, t)
}

func Test_string_detection_double(t *testing.T) {
	elems_with_runes := elem_runes_from_str(`"name": "Bob", "age": 7`)
	elems_strings_detected := Json_collect_strings_in_elems__remove_spaces(elems_with_runes)
	wanted := elem_list{
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
	elems_with_runes := elem_runes_from_str(`"he is \"Eduard\""`)
	elems_strings_detected := Json_collect_strings_in_elems__remove_spaces(elems_with_runes)
	wanted := elem_string("he is \\\"Eduard\\\"")
	compare_one_pair_received_wanted(elems_strings_detected[0], wanted, t)
}

func Test_number_int_detection(t *testing.T) {
	elems_with_runes := elem_runes_from_str(`"price": 7.6, "age": 5`)
	elems_strings_detected := Json_collect_strings_in_elems__remove_spaces(elems_with_runes)
	elems_num_detected := Json_collect_numbers_in_elems(elems_strings_detected)
	Elems_print(elems_num_detected, 0)
	wanted := elem_list{
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
	elems := elem_runes_from_str(`"True": true, "False": false, "age": null `)
	elems = Json_collect_strings_in_elems__remove_spaces(elems)
	elems = Json_collect_scalars_in_elems(elems)
	// Elems_print(elems, 0)
	wanted := elem_list{
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
	elems := elem_runes_from_str(`"name": "Bob", "scores": [4, 6], "friends": [["Eve", 16], ["Joe", 42]], "key": "val"`)
	elems = Json_collect_strings_in_elems__remove_spaces(elems)
	elems = Json_collect_numbers_in_elems(elems)
	array := Json_collect_arrays_in_elems(elems)
	fmt.Println("(1) arrays detected:", len(array))
	Elems_print(array, 0)

	wanted := elem_list{
		elem_string("name"),
		elem_rune(':'),
		elem_string("Bob"),
		elem_rune(','),
		elem_string("scores"),
		elem_rune(':'),
		elem_array(
			elem_list{
				elem_number_int(4),
				elem_number_int(6),
			},
		),
		elem_rune(','),
		elem_string("friends"),
		elem_rune(':'),
		elem_array(elem_list{
			elem_array(elem_list{
				elem_string("Eve"),
				elem_number_int(16),
			}),
			elem_array(elem_list{
				elem_string("Joe"),
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
	elem_root, _ := Json_parse_src(`{"personal":{"city":"Paris", "cell": 123}}`)
	fmt.Println("Test_object_detection")
	Elem_print_one(elem_root)
	wanted := elem_object(keys_elems{
		"personal": elem_object(keys_elems{
			"city": elem_string("Paris"),
			"cell": elem_number_int(123),
		}),
	})
	compare_one_pair_received_wanted(elem_root, wanted, t)
}

func Test_json_1(t *testing.T) {
	elem_root, _ := Json_parse_src(`{"name": "Bob", "friends": [ {"name":"Eve", "cell": 123, "age": 21} ]}`)
	fmt.Println(" Test_json_1")
	Elem_print_one(elem_root)
	wanted := elem_object(keys_elems{
		"name": elem_string("Bob"),
		"friends": elem_array(elem_list{
			elem_object(keys_elems{
				"name": elem_string("Eve"),
				"cell": elem_number_int(123),
				"age":  elem_number_int(21),
			}),
		}),
	})
	compare_one_pair_received_wanted(elem_root, wanted, t)
}

func Test_complex_big(t *testing.T) {
	elem_root, _ := Json_parse_src(`{"name": "Bob", "friends": [ {"name":"Eve", "scores":[1,2]}, {"name":"Joe", "scores":[3,4]} ]}`)
	fmt.Println("Test_complex_big")
	Elem_print_one(elem_root)
	wanted := elem_object(keys_elems{
		"name": elem_string("Bob"),
		"friends": elem_array(elem_list{
			elem_object(keys_elems{
				"name": elem_string("Eve"),
				"scores": elem_array(elem_list{
					elem_number_int(1),
					elem_number_int(2),
				}),
			}),
			elem_object(keys_elems{
				"name": elem_string("Joe"),
				"scores": elem_array(elem_list{
					elem_number_int(3),
					elem_number_int(4),
				}),
			}),
		}),
	})
	compare_one_pair_received_wanted(elem_root, wanted, t)
}

/////////////////////////////////////////////////////////////////////////

func compare_receivedElems_wantedElems(receiveds elem_list, wanteds elem_list, t *testing.T) {
	// compare only the lenth of the top level!
	if len(receiveds) != len(wanteds) {
		fmt.Println(">>> === compare, len !=   ===========")
		Elems_print(receiveds, 0)
		fmt.Println("    .......................    ")
		Elems_print(wanteds, 0)
		fmt.Println("<<< === compare, len !=   ===========")
		t.Fatalf(`len(received_elems %v) != len(wanted_elems %v) `,
			len(receiveds), len(wanteds))

	}
	for i := 0; i < len(receiveds); i++ {
		compare_one_pair_received_wanted(receiveds[i], wanteds[i], t)
	}
}

func _compare_two_objects(objA Elem, objB Elem, t *testing.T) {
	for keyReceived, _ := range objA.ValObject {
		if Obj_has_key(objB.ValObject, keyReceived) == false {
			t.Fatalf(`wanted object doesn't have key' %v error`, keyReceived)
		}
		compare_one_pair_received_wanted(objA.ValObject[keyReceived], objB.ValObject[keyReceived], t)
	}
}

func compare_one_pair_received_wanted(received Elem, wanted Elem, t *testing.T) {
	if received.ValString != wanted.ValString {
		t.Fatalf("\nreceived: %v\n  wanted: %v, error",
			received.ValString, wanted.ValString)
	}
	if received.ValRune != wanted.ValRune {
		t.Fatalf(`received rune = %v, wanted %v, error`,
			received.ValRune,
			wanted.ValRune)
	}
	if received.ValNumberInt != wanted.ValNumberInt {
		t.Fatalf(`received int = %v, wanted %v, error`, received.ValNumberInt, wanted.ValNumberInt)
	}
	if received.ValNumberFloat != wanted.ValNumberFloat {
		t.Fatalf(`received float= %v, wanted %v, error`, received.ValNumberFloat, wanted.ValNumberFloat)
	}
	if received.ValBool != wanted.ValBool {
		t.Fatalf(`received bool = %v, wanted %v, error`, received.ValBool, wanted.ValBool)
	}

	if received.ValType == "array" {
		compare_receivedElems_wantedElems(received.ValArray, wanted.ValArray, t)
	}

	if received.ValType == "object" {
		_compare_two_objects(received, wanted, t) // check based on received object
		_compare_two_objects(wanted, received, t) // check based on wanted object (from other direction)
	}
}
