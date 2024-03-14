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

func (v JSON_value) AddKeyVal_path_into_object(keysMerged string, value JSON_value) error {
	if v.ValType == "object" {
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

		v.updateLevelForChildren()
		return nil
	}
	return errors.New(errorPrefix + "add value into non-object")
}


func (v JSON_value) AddKeyVal_into_object(key string, value JSON_value) error {
	if v.ValType == "object" {
		objects := v.ValObject
		objects[key] = value
		v.ValObject = objects

		v.updateLevelForChildren()
		return nil
	}
	return errors.New(errorPrefix + "add value into non-object")
}

// add value into an ARRAY
func (v JSON_value) AddVal_into_array(value JSON_value) error {
	if v.ValType == "array" {
		elems := v.ValArray
		elems = append(elems, value)
		v.ValArray = elems

		v.updateLevelForChildren()
		return nil
	}
	return errors.New(errorPrefix + "add value into non-array")
}

// TODO: newArray, newObject, newInt, newFloat, newBool....
func NewString_JSON_value(str string) JSON_value {
	return JSON_value{ValType: "string",
		CharPositionFirstInSourceCode: -1,
		CharPositionLastInSourceCode:  -1,
		Runes:                         []rune(`"`+str+`"`),  // strings have "..." boundaries in runes,
		AddedInGoCode:                 true,                 // because in the Json source code the container is "..."
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
	if valueCollected.ValType != "object" {
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

