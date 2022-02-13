package jyp

import (
	"fmt"
	"strings"
)

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
func _rune_digit_info(elem_now elem) (rune, bool) {
	digit_signs := "+-.0123456789"
	rune_now := elem_now.val_rune
	is_digit := strings.ContainsRune(digit_signs, rune_now)
	return rune_now, is_digit
}
func elem_unprocessed(elem elem) bool {
	return elem.val_type == "rune"
}

func Json_collect_numbers_in_elems(src []elem) []elem {
	collector := elems_new()
	runes := runes_new()

	for _, elem_now := range src {

		if elem_unprocessed(elem_now) {
			rune_now, is_digit := _rune_digit_info(elem_now)

			if is_digit {
				runes = append(runes, rune_now)
				continue
			}

			if !is_digit && len(runes) > 0 {
				collector = append(collector, elem{val_string: runes, val_type: number_type(runes)})
				runes = runes_new()
			}
		} // unprocessed
		// { } [ ] : , true false null can be here, for example
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

func number_type(runes []rune) string {
	for _, rune := range runes {
		if rune == '.' {
			return "number_float"
		}
	}
	return "number_int"
}
