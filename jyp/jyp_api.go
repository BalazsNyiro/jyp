/*
Copyright (c) 2024, Balazs Nyiro, balazs.nyiro.ca@gmail.com
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

*/

package jyp

func (v JSON_value) addKeyVal(key string, value JSON_value) {
	if v.ValType == "object" {
		objects := v.ValObject
		objects[key] = value
		v.ValObject = objects
	}
	v.updateLevelForChildren()
}

func (v JSON_value) addVal_key(value JSON_value) {
	if v.ValType == "array" {
		elems := v.ValArray
		elems = append(elems, value)
		v.ValArray = elems
	}
	v.updateLevelForChildren()
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
