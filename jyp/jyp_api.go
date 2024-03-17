/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.


api functions: the often called user supporter functions from the importer program
*/

package jyp

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// if the src can be parsed, return with the JSON root object with nested elems, and err is nil.
func JsonParse(srcStr string) (JSON_value, []error) {

	var errorsCollected []error
	tokens := tokenTable_startPositionIndexed{}
	src := []rune(srcStr)
	srcOrig := []rune(srcStr)

	// the src is always less and less, as tokens are detected
	// the tokens table has more and more elems, as the src sections are parsed
	// at the end, src is total empty (if everything goes well) - and we don't have errors, too

	// only strings can have errors at this parsing step, but the src|tokens|errors are
	// lead through every fun, as a standard solution - so the possibility is open to throw an error everywhere.

	// here maybe the tokens|errorsCollected ret val handling could be removed,
	// but with this, it is clearer what is happening in the fun - so I use this form.
	// in other words: represent if the structure is changed in the function.
	jsonDetect_strings______(src, tokens, errorsCollected)
	jsonDetect_separators___(src, tokens, errorsCollected)
	jsonDetect_trueFalseNull(src, tokens, errorsCollected)
	jsonDetect_numbers______(src, tokens, errorsCollected)

	// at this point, Numbers are not validated - the ruins are collected only,
	// and the lists/objects doesn't have embedded structures - it has to be built, too.
	// src has to be empty, or contain only whitespaces.


	// set correct string values, based on raw rune src.
	// example: "\u0022quote\u0022"'s real form: `"quote"`,
	// so the raw source has to be interpreted (escaped chars, unicode chars)
	valueValidationsSettings_inTokens(srcOrig, tokens, errorsCollected)
	fmt.Println("TokenTable after detections")
	TokensDisplay_startingCoords(srcOrig, tokens)
	elemRoot := objectHierarchyBuilding(tokens, errorsCollected)

	return elemRoot, errorsCollected
}

func (v JSON_value) AddKeyVal_path_into_object(keysMerged string, value JSON_value) error {
	if v.ValType == typeObject {
		keys, err:= ObjPath_merged_expand__split_with_first_char(keysMerged)
		if err != nil {
			return err
		}

		if len(keys) == 1 {
			return v.AddKeyVal_into_object(keys[0], value)
		}

		if len(keys) > 1 {
			object := v.ValObject[keys[0]]

			separator := string(keysMerged[0])
			pathAfterFirstKey := separator+strings.Join(keys[1:], separator)
			err2 := object.AddKeyVal_path_into_object(pathAfterFirstKey, value)
			if err2 != nil {
				return err2
			}
			v.ValObject[keys[0]] = object
		}

		v.local_tool__updateLevelForChildren()
		return nil
	}
	return errors.New(errorPrefix + "add value into non-object")
}


func (v JSON_value) AddKeyVal_into_object(key string, value JSON_value) error {
	if v.ValType == typeObject {
		objects := v.ValObject
		objects[key] = value
		v.ValObject = objects

		v.local_tool__updateLevelForChildren()
		return nil
	}
	return errors.New(errorPrefix + "add value into non-object")
}

// add value into an ARRAY
func (v JSON_value) AddVal_into_array(value JSON_value) error {
	if v.ValType == typeArray {
		elems := v.ValArray
		elems = append(elems, value)
		v.ValArray = elems

		v.local_tool__updateLevelForChildren()
		return nil
	}
	return errors.New(errorPrefix + "add value into non-array")
}

// TODO: newArray, newObject, newInt, newFloat, newBool....
func NewString_JSON_value(str string) JSON_value {
	return JSON_value{ValType: typeString,
		CharPositionFirstInSourceCode: -1,
		CharPositionLastInSourceCode:  -1,
		AddedInGoCode:                 true,                 // because in the Json source code the container is "..."
		ValString: str,
	}
}

func ObjPath_merged_expand__split_with_first_char(path string) ([]string, error){
	if len(path) < 1 {
		return []string{}, errors.New("separator is NOT defined")
	}
	if len(path) < 2 { // minimum one path elem is necessary, that we want to read or write
		// if there is nothing after the separator, the path is empty
		return []string{}, errors.New("separator and minimum one path elem are NOT defined")
	}
	separatorChar := path[0]
	return strings.Split(path, string(separatorChar))[1:], nil
	// so the first empty elem has to be removed (empty string), this is the reason of [1:]
	/*
		if you try to use this:  '/embedded/level2' then before the first separator, an empty string will be in elems
		separator: /
		>>> ''           // EMPTY STRING
		>>> 'embedded'
		>>> 'level2'
		for _, key := range keys {
			print(fmt.Sprintf(">>> '%s' \n", key))
		}
	*/
}


func (v JSON_value) ObjPath(keysMerged string) (JSON_value, error) {
	// object reader with merged string keys (first character is the key elem separator
	// elem_root.ObjPath("/personal/list")     separator: /
	// elem_root.ObjPath("|personal|list")     separator: |
	// elem_root.ObjPath(">personal>list")     separator: |
	// the separator can be any character.
	var valueEmpty JSON_value

	if len(keysMerged) < 2 {
		return valueEmpty, errors.New(errorPrefix + "missing separator and key(s) in merged ObjPath")
	}
	// possible errors are handled with len(...)<2
	keys, _ := ObjPath_merged_expand__split_with_first_char(keysMerged)
	// fmt.Println("KEYS:", keys, len(keys))
	return v.ObjPathKeys(keys)
}

func (v JSON_value) ObjPathKeys(keysEmbedded []string) (JSON_value, error) {
	// object reader with separated string keys:  elem_root.ObjPathKeys([]string{"personal", "list"})
	var valueEmpty JSON_value

	if len(keysEmbedded) < 1 {
		return valueEmpty, errors.New(errorPrefix + "missing object keys (no keys are passed)")
	}

	// minimum 1 key is received
	valueCollected, keyFirstIsKnownInObject := v.ValObject[keysEmbedded[0]]
	if ! keyFirstIsKnownInObject {
		return valueEmpty, errors.New(errorPrefix + "unknown object key (key:"+keysEmbedded[0]+")")
	}

	if len(keysEmbedded) == 1 {
		if keyFirstIsKnownInObject {
			return valueCollected, nil
		}
	}

	// len(keys) > 1
	if valueCollected.ValType != typeObject {
		return valueEmpty, errors.New(errorPrefix + keysEmbedded[0] + "-> child is not object, key cannot be used")
	}
	return valueCollected.ObjPathKeys(keysEmbedded[1:])
}

func (v JSON_value) Arr(index int) (JSON_value, error) {
	// ask ONE indexed elem from an array

	var valueEmpty JSON_value
	indexMax := len(v.ValArray) - 1
	if index > indexMax {
		return valueEmpty, errors.New(errorPrefix + "index ("+strconv.Itoa(index)+") is not in array")
	}

	valueCollected := v.ValArray[index]
	return valueCollected, nil
}

// an ALWAYS string representation of the value
// if indentation > 0: pretty print, with passed indentation per level
// if indentation <= 0, inline print
// zero or one param is accepted. repr() means repr(0), when there is NO indentation (for simple values that is fine)
func (v JSON_value) repr(indentationByUser ...int) string {
	prefix := ""      // dense/inline mode is default, so no prefix
	prefixChildOfObj := ""      // dense/inline mode is default, so no prefix
	lineEnd := ""     // no line ending
	objectKeyValSeparator := ":"  // and tight separator

	indentation := 0
	if len(indentationByUser) > 0 {
		indentation = indentationByUser[0]
	}


	if indentation >= 1 {
		lineEnd = "\n"  // inline print if no indentaion

		prefixFiller := " "
		prefix = strings.Repeat(prefixFiller, v.LevelInObjectStructure*indentation)
		prefixChildOfObj = strings.Repeat(prefixFiller, (v.LevelInObjectStructure+1)*indentation)
		objectKeyValSeparator = ": " // separator with space
	}

	if v.ValType == typeObject || v.ValType == typeArray {
		var charOpen  string
		var charClose string
		var reprValue string

		if v.ValType == typeObject {
			charOpen = "{"
			charClose = "}"

			counter := 0
			for _, childKey := range v.ValObject_keys_sorted() {
				childVal := v.ValObject[childKey]
				comma := base__separator_set_if_no_last_elem(counter, len(v.ValObject), ",")
				reprValue += prefixChildOfObj + "\"" + childKey + "\"" + objectKeyValSeparator + childVal.repr(indentation) + comma + lineEnd
				counter ++
			}
		} else {
			charOpen = "["
			charClose = "]"
			for counter, childVal := range v.ValArray {
				comma := base__separator_set_if_no_last_elem(counter, len(v.ValArray), ",")
				reprValue += prefixChildOfObj + childVal.repr(indentation) + comma + lineEnd
			}
		}

		extraNewlineAfterRootElemPrint := ""
		if v.idSelf == 0 {
			extraNewlineAfterRootElemPrint = "\n"
		}
		return prefix + charOpen + lineEnd + reprValue + prefix + charClose + extraNewlineAfterRootElemPrint

	} else {
		// simple value, not a container
		if v.ValType == typeString { return "\"" + v.ValString + "\"" }
		if v.ValType == typeNull { return "null" }
		if v.ValType == typeBool {
			if v.ValBool { return "true"}
			return "false"
		}
		if v.ValType == typeNumberInt{ return strconv.Itoa(v.ValNumberInt) }
		if v.ValType == typeNumberFloat64{ return strconv.FormatFloat(v.ValNumberFloat, 'f', 0, 64) }
	}
	return "?"
}

func (v JSON_value) ValObject_keys_sorted() []string{
	keys := make([]string, 0, len(v.ValObject))
	for k, _ := range v.ValObject {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
