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

func (v JSON_value) addKeyVal_path(keysMerged string, value JSON_value) error {
	if v.ValType == "object" {
		keys, err:= ObjPath_merged_expand__split_with_first_char(keysMerged)
		if err != nil {
			return err
		}

		if len(keys) == 1 {
			return v.addKeyVal(keys[0], value)
		}

		if len(keys) > 1 {
			object := v.ValObject[keys[0]]

			separator := string(keysMerged[0])
			pathAfterFirstKey := separator+strings.Join(keys[1:], separator)
			err2 := object.addKeyVal_path(pathAfterFirstKey, value)
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


func (v JSON_value) addKeyVal(key string, value JSON_value) error {
	if v.ValType == "object" {
		objects := v.ValObject
		objects[key] = value
		v.ValObject = objects

		v.updateLevelForChildren()
		return nil
	}
	return errors.New(errorPrefix + "add value into non-object")
}

func (v JSON_value) addVal_key(value JSON_value) error {
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
func newString(str string) JSON_value {
	return JSON_value{ValType: "string",
		CharPositionFirstInSourceCode: -1,
		CharPositionLastInSourceCode:  -1,
		Runes:                         []rune(`"`+str+`"`),  // strings have "..." boundaries in runes,
		AddedInGoCode:                 true,                 // because in the Json source code the container is "..."
	}
}
