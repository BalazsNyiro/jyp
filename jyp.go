package jyp

import "fmt"

type value struct {
	val_type   string // string, number, object, array, true, false, null
	val_string string
	val_number int
	val_object map[string]value
	val_array  []value
}

func obj() map[string]value {
	obj_empty := make(map[string]value)
	return obj_empty
}

func Json_parse(src string) (int, error) {
	fmt.Println("json_parse:" + src)
	// obj := obj()
	return len(src), nil
}
