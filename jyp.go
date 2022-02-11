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

// first
func Json_object_finder(src []elem) (elem, error) {

	// ********** find basic string elems *****************
	var src_with_string_elems = make([]elem, len(src))
	var in_text = false
	var collector = elem{val_string: make([]rune, 4), val_type: "string"}
	for i, elem := range src {
		if in_text {
			collector.val_string = append(collector.val_string, elem.val_rune)
		} else {
		}
		if elem.val_type == "rune" {
			fmt.Println(i, " => ", elem.val_type, string(elem.val_rune))
		}
	}

	for i, elem := range src_with_string_elems {
		fmt.Println(i, "--->", elem.val_type)
	}
	fmt.Println(string(collector.val_string))
	return src[0], nil
}
