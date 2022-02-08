package jyp

import "fmt"

type value struct {
	val_type         string // string, number_int, number_float, object, array, true, false, null
	val_string       string
	val_number_int   int
	val_number_float float64
	val_object       map[string]value
	val_array        []value
}

func obj_empty() map[string]value {
	obj_empty := make(map[string]value)
	return obj_empty
}

func Json_parse(src string) (map[string]value, error) {
	fmt.Println("json_parse:" + src)
	obj := obj_empty()
	val := value{val_type: "number_int", val_number_int: 1}
	obj["key"] = val
	return obj, nil
}
