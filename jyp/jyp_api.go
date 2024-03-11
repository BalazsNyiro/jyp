/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

*/

package jyp
// TODO: newArray, newObject, newInt, newFloat, newBool....

func (v JSON_value) addKeyVal_key(key string, value JSON_value) {
	if v.ValType == "object" {
		objects := v.ValObject
		objects[key] = value
		v.ValObject = objects
	}
}

func (v JSON_value) addVal_key(value JSON_value) {
	if v.ValType == "array" {
		elems := v.ValArray
		elems = append(elems, value)
		v.ValArray = elems
	}
}

func (v JSON_value) newString(str string) JSON_value {
	return JSON_value{ValType: "string",
		CharPositionFirstInSourceCode: -1,
		CharPositionLastInSourceCode:  -1,
		Runes:                         []rune(str),
		AddedInGoCode:                 true,
	}
}
