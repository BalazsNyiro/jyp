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

