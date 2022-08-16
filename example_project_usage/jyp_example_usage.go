package main

import "github.com/BalazsNyiro/jyp"
import "fmt"

// run: go run jyp_example_usage.go

func main() {
	elem_root, _ := jyp.Json_parse_src(`{"personal":{"city":"Paris", "cell": 123, "money": 2.34, "list": [1,2,"third"]}}`)
	jyp.Elem_print_one(elem_root)

	// if the JSON structure is unknown for you, maybe you have to check ValType of elements.
	// if you read a known structure, the GETTERS are easier to read.

	// native structure reading:
	fmt.Println(elem_root.ValObject["personal"].ValObject["list"].ValArray[2].ValString)

	// getter functions, same elem reading (check GETTER FUNCS in jyp.go)
	fmt.Println(elem_root.Key("personal").Key("list").ArrayId(2).Str())
	fmt.Println(elem_root.Key("personal").Key("cell").Int())
	fmt.Println(elem_root.Key("personal").Key("money").Float())

	// add new elems into the structure - native solutions:
	elem_root.ValObject["new_string_in_root"] = jyp.ElemStr("New York")
	elem_root.ValObject["new_int_in_root"] = jyp.ElemInt(42)
	elem_root.ValObject["new_float_in_root"] = jyp.ElemFloat("56.78", 56.78)
	elem_root.ValObject["new_object_in_root"] = jyp.ElemObject(jyp.Keys_elems{
		"name": jyp.ElemStr("Eve"),
		"cell": jyp.ElemInt(123),
	})
	jyp.Elem_print_one(elem_root)

}
