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

func Json_parse(src string) (elem, error) {
	fmt.Println("json_parse:" + src)
	chars := elems_from_str(src)
	collector := Json_collect_strings_in_elems__remove_spaces(chars)
	collector = Json_collect_numbers_in_elems(collector)
	elems_print(collector)
	return collector[0], nil
}

func Json_collect_numbers_in_elems(src []elem) []elem {
	var collector = elems_new()
	for id, elem_now := range src {
		if elem_now.val_type == "rune" {
			fmt.Println("id", id)
		}
		collector = append(collector, elem_now)
	}
	return collector
}

func Json_collect_strings_in_elems__remove_spaces(src []elem) []elem {
	var collector = elems_new()
	var in_text = false
	var runes = runes_new()

	for id, elem_now := range src {
		rune_now := elem_now.val_rune

		if in_text && rune_now == '"' {
			escaped := elem_is_escaped_in_string(id, src)

			if !escaped {
				in_text = false
				collector = append(collector,
					elem{val_string: runes, val_type: "string"})
				runes = runes_new()
				continue
			}
		}

		if in_text {
			runes = append(runes, rune_now)
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
	return collector
}

////////////////////////////////////////////////////////////////////////////////////
func elems_print(elems []elem) {
	for i, elem := range elems {
		if elem.val_type == "string" {
			fmt.Println(i, "--->", elem.val_type, string(elem.val_string))
		} else {
			fmt.Println(i, "--->", elem.val_type, string(elem.val_rune))
		}
	}
}

func runes_new() []rune {
	return make([]rune, 0)
}
func elems_new() []elem {
	return make([]elem, 0)
}

func elems_from_str(src string) []elem {
	var chars = make([]elem, len(src))
	for i, rune := range src {
		// fmt.Println(i, "->", string(rune))
		chars[i] = elem{val_rune: rune, val_type: "rune"}
	}
	return chars
}

func elem_is_escaped_in_string(position_of_double_quote int, elems []elem) bool {
	pos_checked := position_of_double_quote
	escaped := false
	for {
		pos_checked-- // move to the previous elem
		if pos_checked < 0 {
			return escaped
		}
		if elems[pos_checked].val_rune != '\\' {
			return escaped
		}
		// val_rune == \  so flip escaped...
		escaped = !escaped
	}
	return escaped
}
