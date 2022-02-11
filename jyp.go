package jyp

import "fmt"

type elem struct {
	val_type string
	// rune, string, number_int, number_float,
	// object, array, true, false, null

	val_rune         rune
	val_string       string
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
	Json_recursive_object_finder(chars)

	return obj, nil
}

func Json_recursive_object_finder(src []elem) (elem, error) {
	for i, elem := range src {
		if elem.val_type == "rune" {
			fmt.Println(i, " => ", elem.val_type, string(elem.val_rune))
		}
	}
	return src[0], nil
}
