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
	valRuneString  string // the Rune's strin representation, ONE char
	valString      string // if type==string, this value represents more characters
	valNumberInt   int
	valNumberFloat float64
	valObject      map[string]elem
	valArray       []elem
}

func Json_parse_src(src string) ([]elem, error) {
	fmt.Println("json_parse:" + src)

	elems := elems_from_str(src)
	// elems_print_with_title(elems, "src")
	elems = Json_parse_elems(elems)
	return elems, nil
}

func Json_parse_elems(elems []elem) []elem {
	elems = Json_collect_strings_in_elems__remove_spaces(elems) // string detection is the first,
	// elems_print_with_title(elems, "collect strings")

	elems = Json_collect_numbers_in_elems(elems) // because strings can contain numbers
	// elems_print_with_title(elems, "collect numbers")

	elems = Json_collect_scalars_in_elems(elems) // or scalars, too
	// elems_print_with_title(elems, "collect scalars")

	elems = Json_collect_arrays_in_elems(elems)
	// elems_print_with_title(elems, "collect arrays")

	elems = Json_collect_objects_in_elems(elems)
	// elems_print_with_title(elems, "collect objects")
	return elems
}

// ******************** array/object detection: ********************************
func Json_collect_arrays_in_elems(src []elem) []elem {
	return Json_structure_ranges_and_hierarchies_in_elems(src, '[', ']', "array")
}
func Json_collect_objects_recursive(src []elem) {

}
func Json_collect_objects_in_elems(src []elem) []elem {
	//But: embedded lists can have embedded objects, too
	// at the beginnin here I have arrays only.
	for id, elemNow := range src {
		if elemNow.valType == "array" {
			src[id].valArray = Json_collect_objects_in_elems(elemNow.valArray)
		}
	}

	// object detection in the top level.
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
				map_data := map[string]elem{}
				key := ""
				for _, elemNow := range elems_embedded {
					if key == "" && elemNow.valType == "string" {
						key = elemNow.valString
						continue
					}

					if key != "" && elemNow.valType != "rune" {
						map_data[key] = elemNow
						key = ""
					}
				}
				elem_pair.valObject = map_data
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
			elemTrue := elem_true()
			collector = append(collector, elemTrue)
		}
		if lastFourChar == "null" {
			idCut := idLast - 3
			collector = collector[:idCut]
			elemNull := elem_null()
			collector = append(collector, elemNull)
		}
		if lastFiveChar == "false" {
			idCut := idLast - 4
			collector = collector[:idCut]
			elemFalse := elem_false()
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
		collector = collector_append_possible_runes(collector, string(runes))
		collector = append(collector, elemNow) // save anything else
		runes = nil                            // clear the slice. if it had anything, the APPEND used it.
	}
	// save the info if digits are the last ones
	collector = collector_append_possible_runes(collector, string(runes))
	return collector
}

func collector_append_possible_runes(collector []elem, numberTxt string) []elem {
	if len(numberTxt) > 0 {
		collector = append(collector, _elem_number_from_runes(numberTxt))
	}
	return collector
}

func _rune_digit_info(elemNow elem) (rune, bool) {
	digitSigns := "+-.0123456789"
	runeNow := elemNow.valRune
	isDigit := strings.ContainsRune(digitSigns, runeNow)
	return runeNow, isDigit
}

// it can work if runes has elems, because it returns with an elem
// and to determine the elem minimum one rune is necessary
func _elem_number_from_runes(stringVal string) elem {
	numType := number_type_detect_float_or_int(stringVal)
	if numType == "number_int" {
		intVal, _ := strconv.Atoi(stringVal)
		return elem_number_int(intVal)
	}
	floatVal := str_to_float(stringVal)
	return elem_number_float(stringVal, floatVal)
}

// ********************* end of JSON number detection *******************************
func _str_closing_quote(inText bool, runeNow rune) bool {
	return inText && runeNow == '"'
}

// ********************* string detection *******************************************
// from one or more rune it creates one elem with collected characters
func Json_collect_strings_in_elems__remove_spaces(src []elem) []elem {
	var collector = elems_new()
	var inText = false
	var runes = runes_new()

	for id, elemNow := range src {
		runeNow := elemNow.valRune
		// fmt.Println(">>> runeNow", string(runeNow), "inText", inText)
		if _str_closing_quote(inText, runeNow) && !elem_is_escaped_in_string(id, src) {
			inText = false
			collector = append(collector, elem_string(string(runes)))
			runes = nil
			continue
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

func elem_string(value string) elem {
	// example:
	// elem{valString: "age", valType: "string"},
	return elem{valString: value, valType: "string"}
}

func elem_object(values map[string]elem) elem {
	return elem{valObject: values, valType: "object"}
}

func elem_array(values []elem) elem {
	return elem{valArray: elems_copy_all(values), valType: "array"}
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
	return elem{valRune: value, valRuneString: string(value), valType: "rune"}
}

//////////////////////////////////////////////////////////////////////////////////////

func elems_copy_all(elems []elem) []elem {
	return elems_copy(elems, 0, len(elems))
}

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

// in array, id is int. 0->value, 1->v2, 2->v3
// but in an object the id's are strings.
// it's easier to manage string id's only
func elem_print(id string, elem elem, indent_level int) {
	prefix := indentation(indent_level)
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

	// print the current elem's type and value
	fmt.Println(prefix, id, "--->", elem.valType, data)

	if elem.valType == "array" {
		elems_print(elem.valArray, indent_level+1)
	}

	if elem.valType == "object" {
		for key, value_in_obj := range elem.valObject {
			elem_print(key, value_in_obj, indent_level+1) // print the value for the key
		}
	}
}

func elems_print_with_title(elems []elem, title string) {
	fmt.Println("===", title, "===")
	elems_print(elems, 0)
}
func elems_print(elems []elem, indent_level int) {
	for id, elem := range elems {
		elem_print(strconv.Itoa(id), elem, indent_level)
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
		chars[i] = elem_rune(rune)
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

func number_type_detect_float_or_int(number_txt string) string {
	for _, rune := range number_txt {
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

func Obj_has_key(dict map[string]elem, key string) bool {
	if _, ok := dict[key]; ok {
		return true
	}
	return false
}

/*
	if elemNow.valType == "object" {
		// for keyNow, valNow := range elemNow.valObject {
		// // 	//elemNow.valObject[keyNow] = Json_collect_objects_in_elems(valNow)
		// }
	}

*/
