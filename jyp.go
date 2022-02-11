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

	for _, elem_now := range src {
		rune_now := elem_now.val_rune

		if in_text && rune_now == '"' {
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

		if rune_now == '"' {
			in_text = true
			continue
		}
		if rune_now != ' ' {
			collector = append(collector, elem_now)
		}
	}

	for i, elem := range collector {

		if elem.val_type == "string" {
			fmt.Println(i, "--->", elem.val_type, string(elem.val_string))
		} else {
			fmt.Println(i, "--->", elem.val_type, string(elem.val_rune))
		}
	}
	fmt.Println(string(runes))

	return src[0], nil
}
