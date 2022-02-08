package jyp

import "fmt"

type value struct {
	val_type   string // string, number, object, array, true, false, null
	val_string string
	val_number int
	val_object map[string]value
	val_array  []value
}

func Json_parse(src string) (int, error) {
	fmt.Println("json_parse:" + src)
	return len(src), nil
}
