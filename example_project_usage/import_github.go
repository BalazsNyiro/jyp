package main

import "github.com/BalazsNyiro/jyp"
import "fmt"

// run: go run import_github.go

func main() {
	elem_root, _ := jyp.Json_parse_src(`{"personal":{"city":"Paris", "cell": 123, "money": 2.34, "list": [1,2,"third"]}}`)
	jyp.Elem_print_one(elem_root)

	jyp.Elem_print("0", elem_root, 0)

	// if the JSON structure is unknown for you, maybe you have to check ValType of elements.
	// if you read a known structure, the GETTERS are easier to read.

	// native structure reading:
	fmt.Println(elem_root.ValObject["personal"].ValObject["list"].ValArray[2].ValString)

	// getter functions, same elem reading (check GETTER FUNCS in jyp.go)
	fmt.Println(elem_root.Key("personal").Key("list").ArrayId(2).Str())
	fmt.Println(elem_root.Key("personal").Key("cell").Int())
	fmt.Println(elem_root.Key("personal").Key("money").Float())

}
