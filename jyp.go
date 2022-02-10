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

type key_value_pair struct {
	key   string
	value value
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

	var chars = make([]rune, len(src))
	for i, rune := range src {
		fmt.Println(i, "->", string(rune))
		chars[i] = rune
	}

	return obj, nil
}

/*
func Json_recursive_object_finder(src string) (value, error) {
	fmt.Println("recursive object finder:" + src)
	for i, c := range src {
		fmt.Println(i, " => ", string(c))
	}
}

*/
