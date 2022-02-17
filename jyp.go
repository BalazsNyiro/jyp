// JYP - Json/Yaml Parser
// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package jyp

import (
	"fmt"
	"strconv"
	"strings"
)

var floatBitsize = 32

type elem struct {
	valType string
	// rune, string, number_int, number_float,
	// object, array, bool, null

	valBool        bool // true, false
	valRune        rune
	valString      string
	valNumberInt   int
	valNumberFloat float64
	valObject      map[string]elem
	valArray       []elem
}

func Json_parse(src string) (elem, error) {
	fmt.Println("json_parse:" + src)
	elems := elems_from_str(src)

	elems = Json_collect_strings_in_elems__remove_spaces(elems) // string detection is the first,
	elems = Json_collect_numbers_in_elems(elems)                // because strings can contain numbers
	elems = Json_collect_scalars_in_elems(elems)                // or scalars, too
	elems = Json_collect_arrays_in_elems(elems)
	elems = Json_collect_objects_in_elems(elems)

	elems_print(elems, 0)
	return elems[0], nil
}

// ******************** array/object detection: ********************************
func Json_collect_arrays_in_elems(src []elem) []elem {
	return Json_structure_ranges_and_hierarchies_in_elems(src, '[', ']', "array")
}

func Json_collect_objects_in_elems(src []elem) []elem {
	return Json_structure_ranges_and_hierarchies_in_elems(src, '{', '}', "object")
}

func Json_structure_ranges_and_hierarchies_in_elems(src []elem, charOpen rune, charClose rune, valType string) []elem {
	src_pair_removed := src
	for {
		pos_last_opening_before_first_closing, pos_first_closing :=
			character_position_first_closed_pair(src_pair_removed, charOpen, charClose)
		if pos_last_opening_before_first_closing < 0 || pos_first_closing < 0 {
			return src_pair_removed
		} else {
			elem_pair := elem{valType: valType}
			elems_embedded := src_pair_removed[pos_last_opening_before_first_closing+1 : pos_first_closing]
			if valType == "array" {
				elem_pair.valArray = elems_embedded
			} else {
				map_representation := map[string]elem{}
				key := ""
				for _, elemNow := range elems_embedded {
					if key == "" && elemNow.valType == "string" {
						key = elemNow.valString
						continue
					}
					// TODO: detect value
					//if key != "" && {}

				}

				elem_pair.valObject = map_representation
			}
			src_new := elems_copy(src_pair_removed, 0, pos_last_opening_before_first_closing)
			src_new = append(src_new, elem_pair)
			src_new = append(src_new, elems_copy(src_pair_removed, pos_first_closing+1, len(src_pair_removed))...)
			src_pair_removed = src_new
		}
	}
}

// ******************** array/object detection: ********************************

// ******************** scalar detection: true, false, null *************
// from more fixed runes it creates one elem
// src can't contain strings! (strings can contain scalar words, too)
func Json_collect_scalars_in_elems(src []elem) []elem {
	collector := elems_new()
	runes := runes_new()

	for id, elemNow := range src {
		runes = append(runes, elemNow.valRune) // collect all runes
		collector = append(collector, elemNow) // a shortest JSON code that can contain a scalar is this: {"a":null}
		if id < 5 {
			continue
		} // false needs 5 chars,

		lastFourChar := string(runes[id-3 : id+1])
		lastFiveChar := string(runes[id-4 : id+1])

		idLast := (len(collector) - 1)
		if lastFourChar == "true" {
			/* in slice operators, the TO id is excluded.
			   collector[:idLast] means: remove the last elems from collectors.
			   collector[:idLast-3] means: remove the last AND the prev 3, so the last 4.
			*/
			idCut := idLast - 3 // remove the last elems from collector
			collector = collector[:idCut]
			elemTrue := elem{valType: "bool", valBool: true}
			collector = append(collector, elemTrue)
		}
		if lastFourChar == "null" {
			idCut := idLast - 3
			collector = collector[:idCut]
			elemNull := elem{valType: "null"}
			collector = append(collector, elemNull)
		}
		if lastFiveChar == "false" {
			idCut := idLast - 4
			collector = collector[:idCut]
			elemFalse := elem{valType: "bool", valBool: false}
			collector = append(collector, elemFalse)
		}
	}
	return collector
}

// ******************** end of scalar detection: ************************

// ********************* number detection *******************************
// from one or more rune it creates one elem with collected digits
// src can't contain strings! (strings can contain numbers, too)
func Json_collect_numbers_in_elems(src []elem) []elem {
	collector := elems_new()
	runes := runes_new()

	for _, elemNow := range src {
		runeNow, isDigit := _rune_digit_info(elemNow)

		if elem_unprocessed(elemNow) && isDigit {
			runes = append(runes, runeNow)
			continue
		}
		collector, runes = collector_append_possible_runes(collector, runes)
		collector = append(collector, elemNow) // save anything else
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

func _rune_digit_info(elemNow elem) (rune, bool) {
	digitSigns := "+-.0123456789"
	runeNow := elemNow.valRune
	isDigit := strings.ContainsRune(digitSigns, runeNow)
	return runeNow, isDigit
}

// it can work if runes has elems, because it returns with an elem
// and to determine the elem minimum one rune is necessary
func _elem_number_from_runes(runes []rune) elem {
	numType := number_type(runes)
	stringVal := string(runes)
	if numType == "number_int" {
		intVal, _ := strconv.Atoi(stringVal)
		return elem{valString: string(runes), valType: numType, valNumberInt: intVal}
	}
	floatVal := str_to_float(stringVal)
	return elem{valString: string(runes), valType: numType, valNumberFloat: floatVal}
}

// ********************* end of JSON number detection *******************************

// ********************* string detection *******************************************
// from one or more rune it creates one elem with collected characters
func Json_collect_strings_in_elems__remove_spaces(src []elem) []elem {
	var collector = elems_new()
	var inText = false
	var runes = runes_new()

	for id, elemNow := range src {
		runeNow := elemNow.valRune

		if inText && runeNow == '"' {
			escaped := elem_is_escaped_in_string(id, src)

			if !escaped {
				inText = false
				collector = append(collector,
					elem{valString: string(runes), valType: "string"})
				runes = nil
				continue
			}
		}
		if inText {
			runes = append(runes, runeNow)
			continue
		}
		if runeNow == '"' {
			inText = true
			continue
		}
		if runeNow != ' ' {
			collector = append(collector, elemNow)
		}
	}
	return collector
}

// ********************* end of string detection *************************************

//////////////////////////////////////////////////////////////////////////////////////
func elems_copy(elems []elem, from_included int, to_excluded int) []elem {
	var collector = elems_new()
	for i := from_included; i < to_excluded; i++ {
		collector = append(collector, elems[i])
	}
	return collector
}
func float_to_string(value float64) string {
	return strconv.FormatFloat(value, 'E', -1, floatBitsize)
}
func str_to_float(value string) float64 {
	floatVal, _ := strconv.ParseFloat(value, floatBitsize)
	return floatVal
}
func elems_print(elems []elem, indent int) {
	prefix := indentation(indent)
	for i, elem := range elems {
		data := ""
		if elem.valType == "array" {
			data = "[...]"
		}
		if elem.valType == "null" {
			data = "null"
		}
		if elem.valType == "bool" {
			if elem.valBool {
				data = "true"
			} else {
				data = "false"
			}
		}
		if elem.valType == "string" {
			data = elem.valString
		}
		if elem.valType == "rune" {
			data = string(elem.valRune)
		}
		if elem.valType == "number_float" {
			data = float_to_string(elem.valNumberFloat)
		}
		if elem.valType == "number_int" {
			data = strconv.Itoa(elem.valNumberInt)
		}
		fmt.Println(prefix, i, "--->", elem.valType, data)
		if elem.valType == "array" {
			elems_print(elem.valArray, indent+1)
		}
	}
}
func indentation(level int) string {
	indentation := ""
	for i := 0; i < level; i++ {
		indentation = indentation + " "
	}
	return indentation
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
		chars[i] = elem{valRune: rune, valType: "rune"}
	}
	return chars
}

func elem_is_escaped_in_string(positionOfDoubleQuote int, elems []elem) bool {
	posChecked := positionOfDoubleQuote
	escaped := false
	for {
		posChecked-- // move to the previous elem
		if posChecked < 0 {
			return escaped
		}
		if elems[posChecked].valRune != '\\' {
			return escaped
		}
		// val_rune == \  so flip escaped...
		escaped = !escaped
	}
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
	return elem.valType == "rune"
}

// Goal: find [ ]  { } pairs ....
// if 0 or positive num: the position of first ] elem
// -1 means: src doesn't have the char
func character_position_first_closed_pair(src []elem, charOpen rune, charClose rune) (int, int) {
	posOpen := -1
	for id, elemNow := range src {
		if elemNow.valRune == charOpen {
			posOpen = id
		}
		if elemNow.valRune == charClose {
			return posOpen, id
		}
	}
	return posOpen, -1
}
