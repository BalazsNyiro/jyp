package jyp

import (
	"fmt"
	"strconv"
	"strings"
)

var float_bitsize = 32

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
	elems := elems_from_str(src)

	elems = Json_collect_strings_in_elems__remove_spaces(elems) // string detection is the first,
	elems = Json_collect_numbers_in_elems(elems)                // because strings can contain numbers
	elems = Json_collect_scalars_in_elems(elems)                // or scalars, too

	elems_print(elems)
	return elems[0], nil
}

// ******************** scalar detection: true, false, null *************
func Json_collect_scalars_in_elems(src []elem) []elem {
	collector := elems_new()
	for _, elem_now := range src {
		collector = append(collector, elem_now)
	}
	return collector
}

// ******************** end of scalar detection: ************************

// ********************* number detection *******************************
// from one or more rune it creates one elem with collected digits
// requirement: not unprocessed strings in the code (strings can contain numbers, too)
func Json_collect_numbers_in_elems(src []elem) []elem {
	collector := elems_new()
	runes := runes_new()

	for _, elem_now := range src {
		rune_now, is_digit := _rune_digit_info(elem_now)

		if elem_unprocessed(elem_now) && is_digit {
			runes = append(runes, rune_now)
			continue
		}
		collector, runes = collector_append_possible_runes(collector, runes)
		collector = append(collector, elem_now) // save anything else
	}
	// save the info if digits are the last ones
	collector, _ = collector_append_possible_runes(collector, runes)
	return collector
}

func collector_append_possible_runes(collector []elem, runes []rune) ([]elem, []rune) {
	if len(runes) > 0 {
		collector = append(collector, _elem_number_from_runes(runes))
		runes = nil // clear the slice?
	}
	return collector, runes
}

func _rune_digit_info(elem_now elem) (rune, bool) {
	digit_signs := "+-.0123456789"
	rune_now := elem_now.val_rune
	is_digit := strings.ContainsRune(digit_signs, rune_now)
	return rune_now, is_digit
}

// it can work if runes has elems, because it returns with an elem
// and to determine the elem minimum one rune is necessary
func _elem_number_from_runes(runes []rune) elem {
	num_type := number_type(runes)
	string_val := string(runes)
	if num_type == "number_int" {
		int_val, _ := strconv.Atoi(string_val)
		return elem{val_string: runes, val_type: num_type, val_number_int: int_val}
	}
	float_val, _ := strconv.ParseFloat(string_val, float_bitsize)
	return elem{val_string: runes, val_type: num_type, val_number_float: float_val}
}

// ********************* end of JSON number detection *******************************

// from one or more rune it creates one elem with collected characters
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
				runes = nil
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
		data := ""
		if elem.val_type == "string" {
			data = string(elem.val_string)
		}
		if elem.val_type == "rune" {
			data = string(elem.val_rune)
		}
		if elem.val_type == "number_float" {
			data = strconv.FormatFloat(elem.val_number_float, 'E', -1, float_bitsize)
		}
		if elem.val_type == "number_int" {
			data = strconv.Itoa(elem.val_number_int)
		}
		fmt.Println(i, "--->", elem.val_type, data)
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

func elem_unprocessed(elem elem) bool {
	return elem.val_type == "rune"
}
