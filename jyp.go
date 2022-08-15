// JYP - Json/Yaml Parser
// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package jyp

import (
	"fmt"
	"strconv"
	"strings"
)

var floatBitsize = 32

type keys_elems map[string]Elem
type Elem_list []Elem

type Elem struct {
	ValType string
	// rune, string, number_int, number_float,
	// object, array, bool, null

	ValBool        bool // true, false
	ValRune        rune
	ValRuneString  string // the Rune's string representation, ONE char
	ValString      string // if type==string, this value represents more characters
	ValNumberInt   int
	ValNumberFloat float64
	ValObject      keys_elems
	ValArray       Elem_list
}

// basically this is the native reading solution of an embedded elem struct/list, 3rd element:
// {"personal":{"city":"Paris", "cell": 123, "list": [1,2,"third"]}}
// fmt.Println(elem_root.ValObject["personal"].ValObject["list"].ValArray[2].ValString)

// because in json an object can have string keys only and a list can have integer keys only,
// there is a wrapper solution to simplify reading:

// GETTER FUNCS
func (elem Elem) Key(key string) Elem {
	return elem.ValObject[key]
}

func (elem Elem) Array() Elem_list {
	return elem.ValArray
}

func (elem Elem) ArrayId(index int) Elem {
	return elem.ValArray[index]
}

func (elem Elem) Str() string {
	return elem.ValString
}

func (elem Elem) Int() int {
	return elem.ValNumberInt
}

func (elem Elem) Float() float64 {
	return elem.ValNumberFloat
}

func (elem Elem) Bool() bool {
	return elem.ValBool
}

/////////////////////////////////////////////////////////////////////////////////////////////

func Json_parse_src(src string) (Elem, error) {
	fmt.Println("json_parse:" + src)

	elems_runes := elem_runes_from_str(src)
	// Elems_print_with_title(elems, "src")
	elems_structured := Json_parse_elems(elems_runes)
	return elems_structured[0], nil // give back the first 'root' object
}

func Json_parse_elems(elems Elem_list) Elem_list {
	elems = Json_collect_strings_in_elems__remove_spaces(elems) // string detection is the first,
	// Elems_print_with_title(elems, "collect strings")

	elems = Json_collect_numbers_in_elems(elems) // because strings can contain numbers
	// Elems_print_with_title(elems, "collect numbers")

	elems = Json_collect_scalars_in_elems(elems) // or scalars, too
	// Elems_print_with_title(elems, "collect scalars")

	elems = Json_collect_arrays_in_elems(elems)
	// Elems_print_with_title(elems, "collect arrays")

	elems = Json_collect_objects_in_elems(elems)
	// Elems_print_with_title(elems, "collect objects")
	return elems
}

// ******************** array/object detection: ********************************
func Json_collect_arrays_in_elems(src Elem_list) Elem_list {
	return Json_structure_ranges_and_hierarchies_in_elems(src, '[', ']', "array")
}

func Json_collect_objects_in_elems(src Elem_list) Elem_list {
	//But: embedded lists can have embedded objects, too
	// at the beginnin here I have arrays only.
	for id, elemNow := range src {
		if elemNow.ValType == "array" {
			src[id].ValArray = Json_collect_objects_in_elems(elemNow.ValArray)
		}
	}

	// object detection in the top level.
	return Json_structure_ranges_and_hierarchies_in_elems(src, '{', '}', "object")
}

func comma_runes_removing(elems Elem_list) Elem_list {
	filtered := elems_new()
	for _, elemNow := range elems {
		if !(elemNow.ValType == "rune" && elemNow.ValRune == ',') {
			filtered = append(filtered, elemNow)
		}
	}
	return filtered
}

func Json_structure_ranges_and_hierarchies_in_elems(src Elem_list, charOpen rune, charClose rune, valType string) Elem_list {
	src_pair_removed := src
	for {
		pos_last_opening_before_first_closing, pos_first_closing :=
			character_position_first_closed_pair(src_pair_removed, charOpen, charClose)
		if pos_last_opening_before_first_closing < 0 || pos_first_closing < 0 {
			return src_pair_removed
		} else {
			elem_pair := Elem{ValType: valType}
			elems_embedded := src_pair_removed[pos_last_opening_before_first_closing+1 : pos_first_closing]
			if valType == "array" {
				elem_pair.ValArray = comma_runes_removing(elems_embedded)
			} else {
				map_data := keys_elems{}
				key := ""
				for _, elemNow := range elems_embedded {
					if key == "" && elemNow.ValType == "string" {
						key = elemNow.ValString
						continue
					}

					if key != "" && elemNow.ValType != "rune" {
						map_data[key] = elemNow
						key = ""
					}
				}
				elem_pair.ValObject = map_data
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
// from more fixed runes it creates one Elem
// src can't contain strings! (strings can contain scalar words, too)
func Json_collect_scalars_in_elems(src Elem_list) Elem_list {
	collector := elems_new()
	runes := runes_new()

	for id, elemNow := range src {
		runes = append(runes, elemNow.ValRune) // collect all runes
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
// from one or more rune it creates one Elem with collected digits
// src can't contain strings! (strings can contain numbers, too)
func Json_collect_numbers_in_elems(src Elem_list) Elem_list {
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

func collector_append_possible_runes(collector Elem_list, numberTxt string) Elem_list {
	if len(numberTxt) > 0 {
		collector = append(collector, _elem_number_from_runes(numberTxt))
	}
	return collector
}

func _rune_digit_info(elemNow Elem) (rune, bool) {
	digitSigns := "+-.0123456789"
	runeNow := elemNow.ValRune
	isDigit := strings.ContainsRune(digitSigns, runeNow)
	return runeNow, isDigit
}

// it can work if runes has elems, because it returns with an Elem
// and to determine the Elem minimum one rune is necessary
func _elem_number_from_runes(stringVal string) Elem {
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
// from one or more rune it creates one Elem with collected characters
func Json_collect_strings_in_elems__remove_spaces(src Elem_list) Elem_list {
	var collector = elems_new()
	var inText = false
	var runes = runes_new()

	for id, elemNow := range src {
		runeNow := elemNow.ValRune
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

func elem_string(value string) Elem {
	// example:
	// Elem{ValString: "age", ValType: "string"},
	return Elem{ValString: value, ValType: "string"}
}

func elem_object(values keys_elems) Elem {
	return Elem{ValObject: values, ValType: "object"}
}

func elem_array(values Elem_list) Elem {
	return Elem{ValArray: elems_copy_all(values), ValType: "array"}
}

func elem_number_int(value int) Elem {
	// return Elem{ValString: "5", ValType: "number_int", ValNumberInt: 5},
	return Elem{ValString: strconv.Itoa(value), ValType: "number_int", ValNumberInt: value}
}

func elem_number_float(value_str_representation string, value_more_or_less_precise float64) Elem {
	// Elem{ValString: "7.6", ValType: "number_float", ValNumberFloat: 7.599999904632568},
	return Elem{ValString: value_str_representation, ValType: "number_float", ValNumberFloat: value_more_or_less_precise}
}

func elem_true() Elem {
	return Elem{ValBool: true, ValType: "bool"}
}

func elem_false() Elem {
	return Elem{ValBool: false, ValType: "bool"}
}

func elem_null() Elem {
	return Elem{ValType: "null"}
}

func elem_rune(value rune) Elem {
	// example:
	// Elem{ValRune: ':', ValType: "rune"},
	return Elem{ValRune: value, ValRuneString: string(value), ValType: "rune"}
}

//////////////////////////////////////////////////////////////////////////////////////

func elems_copy_all(elems Elem_list) Elem_list {
	return elems_copy(elems, 0, len(elems))
}

func elems_copy(elems Elem_list, from_included int, to_excluded int) Elem_list {
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
func Elem_print(id string, elem Elem, indent_level int) {
	prefix := indentation(indent_level)
	data := ""
	if elem.ValType == "array" {
		data = "[...]"
	}
	if elem.ValType == "null" {
		data = "null"
	}
	if elem.ValType == "bool" {
		if elem.ValBool {
			data = "true"
		} else {
			data = "false"
		}
	}
	if elem.ValType == "string" {
		data = elem.ValString
	}
	if elem.ValType == "rune" {
		data = string(elem.ValRune)
	}
	if elem.ValType == "number_float" {
		data = float_to_string(elem.ValNumberFloat)
	}
	if elem.ValType == "number_int" {
		data = strconv.Itoa(elem.ValNumberInt)
	}

	// print the current Elem's type and value
	fmt.Println(prefix, id, "--->", elem.ValType, data)

	if elem.ValType == "array" {
		Elems_print(elem.ValArray, indent_level+1)
	}

	if elem.ValType == "object" {
		for key, value_in_obj := range elem.ValObject {
			Elem_print(key, value_in_obj, indent_level+1) // print the value for the key
		}
	}
}

func Elems_print_with_title(elems Elem_list, title string) {
	fmt.Println("===", title, "===")
	Elems_print(elems, 0)
}
func Elems_print(elems Elem_list, indent_level int) {
	for id, elem := range elems {
		Elem_print(strconv.Itoa(id), elem, indent_level)
	}
}
func Elem_print_one(elem Elem) {
	Elem_print("0", elem, 0)
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
func elems_new() Elem_list {
	return make(Elem_list, 0)
}

func elem_runes_from_str(src string) Elem_list {
	var elems = make(Elem_list, len(src))
	for i, rune := range src {
		// fmt.Println(i, "->", string(rune))
		elems[i] = elem_rune(rune)
	}
	return elems
}

func elem_is_escaped_in_string(positionOfDoubleQuote int, elems Elem_list) bool {
	posChecked := positionOfDoubleQuote
	escaped := false
	for {
		posChecked-- // move to the previous Elem
		if posChecked < 0 {
			return escaped
		}
		if elems[posChecked].ValRune != '\\' {
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

func elem_unprocessed(elem Elem) bool {
	return elem.ValType == "rune"
}

// Goal: find [ ]  { } pairs ....
// if 0 or positive num: the position of first ] Elem
// -1 means: src doesn't have the char
func character_position_first_closed_pair(src Elem_list, charOpen rune, charClose rune) (int, int) {
	posOpen := -1
	for id, elemNow := range src {
		if elemNow.ValRune == charOpen {
			posOpen = id
		}
		if elemNow.ValRune == charClose {
			return posOpen, id
		}
	}
	return posOpen, -1
}

func Obj_has_key(dict keys_elems, key string) bool {
	if _, ok := dict[key]; ok {
		return true
	}
	return false
}
