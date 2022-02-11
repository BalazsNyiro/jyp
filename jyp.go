package jyp

import "fmt"

type elem struct {
	val_type string
	// rune, string, number_int, number_float,
	// object, array, true, false, null

	val_rune         rune
	val_string       []rune
	val_number_int   int
	val_number_float float64
	val_object       map[string]elem
	val_array        []elem
}

func obj_empty() map[string]elem {
	obj_empty := make(map[string]elem)
	return obj_empty
}

func Json_parse(src string) (map[string]elem, error) {
	fmt.Println("json_parse:" + src)
	obj := obj_empty()
	val := elem{val_type: "number_int", val_number_int: 1}
	obj["key"] = val

	var chars = make([]elem, len(src))
	for i, rune := range src {
		fmt.Println(i, "->", string(rune))
		chars[i] = elem{val_rune: rune, val_type: "rune"}
	}
	Json_object_finder(chars)

	return obj, nil
}

func runes_new() []rune {
	return make([]rune, 0)
}
func elems_new(size int) []elem {
	return make([]elem, size)
}

func Json_object_finder(src []elem) (elem, error) {

	// ********** find basic string elems *****************
	var collector = elems_new(len(src))
	var in_text = false
	var runes = runes_new()

	for i, elem_now := range src {
		char := string(elem_now.val_rune)

		if in_text && char == "\"" {
			in_text = false
			collector = append(collector,
				elem{val_string: runes, val_type: "string"})
			runes = runes_new()
			continue
		}

		if in_text {
			runes = append(runes, elem_now.val_rune)
			continue
		}

		if char == "\"" {
			in_text = true
			continue
		} else {
			collector = append(collector, elem_now)
		}
		runes = runes_new()

		if elem_now.val_type == "rune" {
			fmt.Println(i, " => ", elem_now.val_type, string(elem_now.val_rune))
		}
	}

	for i, elem := range collector {
		fmt.Println(i, "--->", elem.val_type, string(elem.val_rune), string(elem.val_string))
	}
	fmt.Println(string(runes))
	// ********** find basic string elems *****************

	return src[0], nil
}
