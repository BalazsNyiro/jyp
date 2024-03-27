/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

 */

package jyp

import (
	"errors"
	"strconv"
	"strings"
)

func JsonParse(srcStr string) (JSON_value, []error) {

	tokensTableB := stepA__tokensTableDetect_structuralTokens_strings_L1(srcStr)
	errorsCollected := stepB__JSON_validation_L1(tokensTableB)
	elemRoot, _ := stepC__JSON_structure_building__L1(srcStr, tokensTableB, 0, errorsCollected)

	return elemRoot, errorsCollected
}

// example usage: Repr(2) means: use "  " 2 spaces as indentation
// if ind
// otherwise a simple formatted output
func (v JSON_value) Repr(indentationLength ...int) string {
	if len(indentationLength) == 0 { // so no param is passed
		return v.Repr_tuned("", 0)
	}
	if indentationLength[0] < 1 { // a 0 or negative param is passed
		return v.Repr_tuned("", 0)
	}
	indentation := base__prefixGenerator_for_repr(" ", indentationLength[0])
	return v.Repr_tuned(indentation, 0)
}

// tunable repr: with this, tabulator can be used for example instead of spaces as indent,
// level 0 means left align - if higher level is used, the output will be moved to right on the screen
func (v JSON_value) Repr_tuned(indent string, level int) string {
	prefix := "" // indentOneLevelPrefix
	prefix2 := "" // indentTwoLevelPrefix
	newLine := ""
	colon := ":"

	if len(indent) > 0 {
		prefix = base__prefixGenerator_for_repr(indent, level)
		prefix2 = base__prefixGenerator_for_repr(indent, level+1)
		newLine = "\n"
		colon = ": "
	}

	if v.ValType == '"' {
		return "\""+v.ValString + "\""
	} else

	if v.ValType == 'I' {
		return strconv.Itoa(v.ValNumberInt)
	} else

	if v.ValType == 'F' {
		return strconv.FormatFloat(v.ValNumberFloat, 'f', -1, 64)
	} else

	if v.ValType == 'b' {
		if v.ValBool {
			return "true"
		}
		return "false"
	} else

	if v.ValType == 'n' {
		return "null"
	} else

	if v.ValType == '{' {
		out := prefix + "{" + newLine
		for counter, childKey := range v.ValObject_keys_sorted() {
			comma := base__separator_set_if_no_last_elem(counter, len(v.ValObject), ",")
			childVal := v.ValObject[childKey]
			out += prefix2 + "\""+childKey+"\"" + colon + childVal.Repr_tuned(indent, level+1) + comma + newLine
		}
		out += prefix + "}"
		return out
	} else

	if v.ValType == '[' {
		out := prefix + "[" + newLine
		for counter, child := range v.ValArray {
			comma := base__separator_set_if_no_last_elem(counter, len(v.ValArray), ",")
			out += prefix2 + indent + child.Repr_tuned(indent, level+1) + comma + newLine
		}
		out += prefix + "]"
		return out
	}
	return ""
}


func NewNull() JSON_value {
	return JSON_value{
		ValType: 'n',
	}
}

func NewBool(val bool) JSON_value {
	return JSON_value{
		ValType: 'b', // bool
		ValBool: val,
	}
}

func NewNumInt(num int) JSON_value { // TESTED
	return JSON_value{
		ValType:      'I', // Integer
		ValNumberInt: num,
	}
}

func NewNumFloat(num float64) JSON_value {
	return JSON_value{
		ValType:      'F',    // Float
		ValNumberFloat: num,
	}
}


var exampleNewString string = "backslashTab\t or newline \n represents 1 char, not two"
// new string object from simple Go string (interpreted)
func NewStr(text string) JSON_value { // TESTED
	return JSON_value{
		ValType:      '"',
		ValString: text,
	}
}

// raw: the string needs to be interpreted. "a\tb": \t represents 2 chars, it needs to be interpreted.
func NewString__rawToInterpreted__QuotedBothEnd(text string, errorsCollected []error) JSON_value {
	// strictly have minimum one "opening....and...one..closing" quote!
	return JSON_value{
		ValType:      '"',
		ValString: stringValueParsing_rawToInterpretedCharacters_L2( text[1:len(text)-1], errorsCollected),
	}
}

func NewObj() JSON_value { // TESTED
	return JSON_value{
		ValType: '{',
		ValObject: map[string]JSON_value{},
	}
}

// newArray creation, with possible children
func NewArr(elems ...JSON_value) JSON_value { // TESTED
	children := []JSON_value{}
	for _, elem := range elems {
		children = append(children, elem)
	}
	return JSON_value{
		ValType: '[',
		ValArray: children,
	}
}

//////////////////////////////////////////////////////////////////////////////////////

// if autoCreateChildren == true, a complex path can be added, and if the Children doesn't exist,
// the func creates it recursively
func (v JSON_value) SetPath(keysMerged string, value JSON_value, autoCreateChildren bool) error { // TESTED
	if v.ValType ==  '{' {
		keys, err:= ObjPath_merged_expand__split_with_first_char(keysMerged)
		if err != nil {
			return err
		}

		if len(keys) == 1 {
			return v.AddKeyVal(keys[0], value)
		}

		if len(keys) > 1 {

			children, isChildInObj := v.ValObject[keys[0]]
			if autoCreateChildren {
				if !isChildInObj {
					children = NewObj()
				}
			} else {  // not autoCreateChildren
				if ! isChildInObj {
					return errors.New("unknown key: " + keys[0])
				}
			}

			separator := string(keysMerged[0])
			pathAfterFirstKey := separator+strings.Join(keys[1:], separator)
			err2 := children.SetPath(pathAfterFirstKey, value, autoCreateChildren)
			if err2 != nil {
				return err2
			}
			v.ValObject[keys[0]] = children
		}

		return nil
	}
	return errors.New(errorPrefix + "add value into non-object")
}


func (v JSON_value) AddKeyVal(key string, value JSON_value) error {
	if v.ValType ==  '{' {
		objects := v.ValObject
		objects[key] = value
		v.ValObject = objects
		return nil
	}
	return errors.New(errorPrefix + "add value into non-object")
}

// add value into an ARRAY
func (v JSON_value) AddVal_into_array(value JSON_value) error {
	if v.ValType ==  '[' {
		elems := v.ValArray
		elems = append(elems, value)
		v.ValArray = elems
		return nil
	}
	return errors.New(errorPrefix + "add value into non-array")
}



// new value has to be inserted into the array
func (v JSON_value) Arr(index int) (JSON_value, error) { // TESTED
	// ask ONE indexed elem from an array

	var valueEmpty JSON_value
	indexMax := len(v.ValArray) - 1
	if index > indexMax {
		return valueEmpty, errors.New(errorPrefix + "index ("+strconv.Itoa(index)+") is not in array")
	}

	valueCollected := v.ValArray[index]
	return valueCollected, nil
}



//////////////////////////////////////////////////////////////////////////////////////

func (v JSON_value) GetPath(keysMerged string) (JSON_value, error) { // TESTED
	// object reader with merged string keys (first character is the key elem separator
	// elem_root.GetPath("/personal/list")     separator: /
	// elem_root.GetPath("|personal|list")     separator: |
	// elem_root.GetPath(">personal>list")     separator: |
	// the separator can be any character.
	var valueEmpty JSON_value

	if len(keysMerged) < 2 {
		return valueEmpty, errors.New(errorPrefix + "missing separator and key(s) in merged GetPath")
	}
	// possible errors are handled with len(...)<2
	keys, _ := ObjPath_merged_expand__split_with_first_char(keysMerged)
	// fmt.Println("KEYS:", keys, len(keys))
	return v.GetPathKeys(keys)
}



func ObjPath_merged_expand__split_with_first_char(path string) ([]string, error){ // TESTED
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

func (v JSON_value) GetPathKeys(keysEmbedded []string) (JSON_value, error) { // TESTED
	// object reader with separated string keys:  elem_root.GetPathKeys([]string{"personal", "list"})
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
	if valueCollected.ValType !=  '{' {
		return valueEmpty, errors.New(errorPrefix + keysEmbedded[0] + "-> child is not object, key cannot be used")
	}
	return valueCollected.GetPathKeys(keysEmbedded[1:])
}
