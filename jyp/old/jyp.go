// JP - Json parser
// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package jp

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var floatBitsize = 32

type Keys_elems map[string]Elem
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
	ValObject      Keys_elems
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

func (elem Elem) jsonRender() string {
	return jsonRenderPrettyRecursive(elem, 1, "") // no prefix/pretty print param, simple rendering
}

func (elem Elem) jsonRenderPretty() string {
	return jsonRenderPrettyRecursive(elem, 1, "  ")
}

// the first level is 1 because the root object is 0.
// and everything in root is in level 1.

// not callable as object's method. Hidden/private fun
// private because I don't want to make mess in the callable functions
func jsonRenderPrettyRecursive(elem Elem, level int, pretty_print_prefix_block string) string {
	prefixInternalElems := ""
	if len(pretty_print_prefix_block) > 0 {
		prefixInternalElems = "\n" + strings.Repeat(pretty_print_prefix_block, level)
	}

	quote := "\""
	if elem.ValType == "bool" {
		if elem.ValBool {
			return "true"
		} else {
			return "false"
		}
	}
	if elem.ValType == "string" {
		return quote + elem.ValString + quote
	}
	if elem.ValType == "number_int" {
		return fmt.Sprintf("%d", elem.ValNumberInt)
	}
	if elem.ValType == "number_float" {
		return elem.ValString // when the float is created, string representation is given as param
	} // at the elem initialization

	// in a list, the order can be important so in the display you can't sort it.
	// TODO: pretty print of ARRAYS :-)
	if elem.ValType == "array" {
		accumulator := ""
		separator := ""
		for _, list_member := range elem.ValArray {
			accumulator = accumulator + separator + jsonRenderPrettyRecursive(list_member, level+1, pretty_print_prefix_block)
			separator = ","
		}
		return "[" + accumulator + "]"
	}

	if elem.ValType == "object" {
		accumulator := ""
		separator := ""

		// the output will be sorted in json output
		for _, key := range keysSortedFromObject(elem.ValObject) {
			accumulator = accumulator + separator + prefixInternalElems + quote + key + quote + ": " + jsonRenderPrettyRecursive(elem.ValObject[key], level+1, pretty_print_prefix_block)
			separator = ","
		}

		prefixParentClose := ""

		if len(pretty_print_prefix_block) > 0 {
			prefixParentClose = "\n"
			indentParent := strings.Repeat(pretty_print_prefix_block, max(level-1, 0))
			prefixParentClose += indentParent
		}

		return "{" + accumulator + prefixParentClose + "}"
	}
	return "json render error"
}

/////////////////////////////////////////////////////////////////////////////////////////////

func JsonParseSrc(src string) (Elem, error) {
	// fmt.Println("json_parse:" + src)

	elems_runes := elemRunesFromStr(src)
	// Elems_print_with_title(elems, "src")
	elems_structured := JsonParseElems(elems_runes)
	return elems_structured[0], nil // give back the first 'root' object
}

func JsonParseElems(elems Elem_list) Elem_list {
	elems = JsonCollectStringsInElems__removeSpaces(elems) // string detection is the first,
	// Elems_print_with_title(elems, "collect strings")

	elems = JsonCollectNumbersInElems(elems) // because strings can contain numbers
	// Elems_print_with_title(elems, "collect numbers")

	elems = JsonCollectScalarsInElems(elems) // or scalars, too
	// Elems_print_with_title(elems, "collect scalars")

	elems = JsonCollectArraysInElems(elems)
	// Elems_print_with_title(elems, "collect arrays")

	elems = JsonCollectObjectsInElems(elems)
	// Elems_print_with_title(elems, "collect objects")
	return elems
}

// ******************** array/object detection: ********************************
func JsonCollectArraysInElems(src Elem_list) Elem_list {
	return JsonStructureRangesAndHierarchiesInElems(src, '[', ']', "array")
}

func JsonCollectObjectsInElems(src Elem_list) Elem_list {
	//But: embedded lists can have embedded objects, too
	// at the beginnin here I have arrays only.
	for id, elemNow := range src {
		if elemNow.ValType == "array" {
			src[id].ValArray = JsonCollectObjectsInElems(elemNow.ValArray)
		}
	}

	// object detection in the top level.
	return JsonStructureRangesAndHierarchiesInElems(src, '{', '}', "object")
}

func comma_runes_removing(elems Elem_list) Elem_list {
	filtered := elemsNew()
	for _, elemNow := range elems {
		if !(elemNow.ValType == "rune" && elemNow.ValRune == ',') {
			filtered = append(filtered, elemNow)
		}
	}
	return filtered
}

func JsonStructureRangesAndHierarchiesInElems(src Elem_list, charOpen rune, charClose rune, valType string) Elem_list {
	src_pair_removed := src
	for {
		pos_last_opening_before_first_closing, pos_first_closing :=
			characterPositionFirstClosedPair(src_pair_removed, charOpen, charClose)
		if pos_last_opening_before_first_closing < 0 || pos_first_closing < 0 {
			return src_pair_removed
		} else {
			elem_pair := Elem{ValType: valType}
			elems_embedded := src_pair_removed[pos_last_opening_before_first_closing+1 : pos_first_closing]
			if valType == "array" {
				elem_pair.ValArray = comma_runes_removing(elems_embedded)
			} else {
				map_data := Keys_elems{}
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
			src_new := elemsCopy(src_pair_removed, 0, pos_last_opening_before_first_closing)
			src_new = append(src_new, elem_pair)
			src_new = append(src_new, elemsCopy(src_pair_removed, pos_first_closing+1, len(src_pair_removed))...)
			src_pair_removed = src_new
		}
	}
}

// ******************** array/object detection: ********************************

// ******************** scalar detection: true, false, null *************
// from more fixed runes it creates one Elem
// src can't contain strings! (strings can contain scalar words, too)
func JsonCollectScalarsInElems(src Elem_list) Elem_list {
	collector := elemsNew()
	runes := runesNew()

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
			elemTrue := ElemTrue()
			collector = append(collector, elemTrue)
		}
		if lastFourChar == "null" {
			idCut := idLast - 3
			collector = collector[:idCut]
			elemNull := ElemNull()
			collector = append(collector, elemNull)
		}
		if lastFiveChar == "false" {
			idCut := idLast - 4
			collector = collector[:idCut]
			elemFalse := ElemFalse()
			collector = append(collector, elemFalse)
		}
	}
	return collector
}

// ******************** end of scalar detection: ************************

// ********************* number detection *******************************
// from one or more rune it creates one Elem with collected digits
// src can't contain strings! (strings can contain numbers, too)
func JsonCollectNumbersInElems(src Elem_list) Elem_list {
	collector := elemsNew()
	runes := runesNew()

	for _, elemNow := range src {
		runeNow, isDigit := runeDigitInfo(elemNow)

		if elemUnprocessed(elemNow) && isDigit {
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
		collector = append(collector, elemNumberFromRunes(numberTxt))
	}
	return collector
}

func runeDigitInfo(elemNow Elem) (rune, bool) {
	digitSigns := "+-.0123456789"
	runeNow := elemNow.ValRune
	isDigit := strings.ContainsRune(digitSigns, runeNow)
	return runeNow, isDigit
}

// it can work if runes has elems, because it returns with an Elem
// and to determine the Elem minimum one rune is necessary
func elemNumberFromRunes(stringVal string) Elem {
	numType := numberTypeDetectFloatOrInt(stringVal)
	if numType == "number_int" {
		intVal, _ := strconv.Atoi(stringVal)
		return ElemInt(intVal)
	}
	floatVal := strToFloat(stringVal)
	return ElemFloat(stringVal, floatVal)
}

// ********************* end of JSON number detection *******************************
func strClosingQuote(inText bool, runeNow rune) bool {
	return inText && runeNow == '"'
}

// ********************* string detection *******************************************
// from one or more rune it creates one Elem with collected characters
func JsonCollectStringsInElems__removeSpaces(src Elem_list) Elem_list {
	var collector = elemsNew()
	var inText = false
	var runes = runesNew()

	for id, elemNow := range src {
		runeNow := elemNow.ValRune
		// fmt.Println(">>> runeNow", string(runeNow), "inText", inText)
		if strClosingQuote(inText, runeNow) && !elemIsEscapedInString(id, src) {
			inText = false
			collector = append(collector, ElemStr(string(runes)))
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

		// remove whitespaces if we are not in text:
		// if there is a real (non-whitespace) character
		if len(strings.TrimSpace(string(runeNow))) > 0 {
			collector = append(collector, elemNow)
		}
	}
	return collector
}

// ********************* end of string detection *************************************

//////////////////////////////////////////////////////////////////////////////////////

func ElemStr(value string) Elem {
	// example:
	// Elem{ValString: "age", Type: "string"},
	return Elem{ValString: value, ValType: "string"}
}

func ElemObject(values Keys_elems) Elem {
	return Elem{ValObject: values, ValType: "object"}
}

func ElemArray(values Elem_list) Elem {
	return Elem{ValArray: elemsCopyAll(values), ValType: "array"}
}

func ElemInt(value int) Elem {
	// return Elem{ValString: "5", Type: "number_int", ValNumberInt: 5},
	return Elem{ValString: strconv.Itoa(value), ValType: "number_int", ValNumberInt: value}
}

func ElemFloat(value_str_representation string, value_more_or_less_precise float64) Elem {
	// Elem{ValString: "7.6", Type: "number_float", ValNumberFloat: 7.599999904632568},
	return Elem{ValString: value_str_representation, ValType: "number_float", ValNumberFloat: value_more_or_less_precise}
}

func ElemTrue() Elem {
	return Elem{ValBool: true, ValType: "bool"}
}

func ElemFalse() Elem {
	return Elem{ValBool: false, ValType: "bool"}
}

func ElemNull() Elem {
	return Elem{ValType: "null"}
}

func elemRune(value rune) Elem {
	// example:
	// Elem{ValRune: ':', Type: "rune"},
	return Elem{ValRune: value, ValRuneString: string(value), ValType: "rune"}
}

//////////////////////////////////////////////////////////////////////////////////////

func elemsCopyAll(elems Elem_list) Elem_list {
	return elemsCopy(elems, 0, len(elems))
}

func elemsCopy(elems Elem_list, from_included int, to_excluded int) Elem_list {
	var collector = elemsNew()
	for i := from_included; i < to_excluded; i++ {
		collector = append(collector, elems[i])
	}
	return collector
}
func floatToString(value float64) string {
	return strconv.FormatFloat(value, 'E', -1, floatBitsize)
}
func strToFloat(value string) float64 {
	floatVal, _ := strconv.ParseFloat(value, floatBitsize)
	return floatVal
}

// in array, id is int. 0->value, 1->v2, 2->v3
// but in an object the id's are strings.
// it's easier to manage string id's only
func ElemPrint(id string, elem Elem, indent_level int) {
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
		data = floatToString(elem.ValNumberFloat)
	}
	if elem.ValType == "number_int" {
		data = strconv.Itoa(elem.ValNumberInt)
	}

	// print the current Elem's type and value
	fmt.Println(prefix, id, "--->", elem.ValType, data)

	if elem.ValType == "array" {
		ElemsPrint(elem.ValArray, indent_level+1)
	}

	if elem.ValType == "object" {
		for key, value_in_obj := range elem.ValObject {
			ElemPrint(key, value_in_obj, indent_level+1) // print the value for the key
		}
	}
}

func Elems_print_with_title(elems Elem_list, title string) {
	fmt.Println("===", title, "===")
	ElemsPrint(elems, 0)
}
func ElemsPrint(elems Elem_list, indent_level int) {
	for id, elem := range elems {
		ElemPrint(strconv.Itoa(id), elem, indent_level)
	}
}
func ElemPrintOne(elem Elem) {
	ElemPrint("0", elem, 0)
}
func indentation(level int) string {
	indentation := ""
	for i := 0; i < level; i++ {
		indentation = indentation + " "
	}
	return indentation
}

func runesNew() []rune {
	return make([]rune, 0)
}
func elemsNew() Elem_list {
	return make(Elem_list, 0)
}

func elemRunesFromStr(src string) Elem_list {
	var elems = make(Elem_list, len(src))
	for i, rune := range src {
		// fmt.Println(i, "->", string(rune))
		elems[i] = elemRune(rune)
	}
	return elems
}

func elemIsEscapedInString(positionOfDoubleQuote int, elems Elem_list) bool {
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

func numberTypeDetectFloatOrInt(number_txt string) string {
	for _, rune := range number_txt {
		if rune == '.' {
			return "number_float"
		}
	}
	return "number_int"
}

func elemUnprocessed(elem Elem) bool {
	return elem.ValType == "rune"
}

// Goal: find [ ]  { } pairs ....
// if 0 or positive num: the position of first ] Elem
// -1 means: src doesn't have the char
func characterPositionFirstClosedPair(src Elem_list, charOpen rune, charClose rune) (int, int) {
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

func ObjHasKey(dict Keys_elems, key string) bool {
	if _, ok := dict[key]; ok {
		return true
	}
	return false
}

func keysSortedFromObject(keyElemPairs Keys_elems) []string {
	keysSorted := []string{}
	for key := range keyElemPairs {
		keysSorted = append(keysSorted, key)
	}
	sort.Strings(keysSorted)
	return keysSorted
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}