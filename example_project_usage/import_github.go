package main

import "github.com/BalazsNyiro/jyp"
import "fmt"

func main() {
	elem_root, _ := jyp.Json_parse_src(`{"personal":{"city":"Paris", "cell": 123, "list": [1,2,"third"]}}`)
	jyp.Elem_print_one(elem_root)

	jyp.Elem_print("0", elem_root, 0)

  // native structure reading:
  fmt.Println(elem_root.ValObject["personal"].ValObject["list"].ValArray[2].ValString)

  // getter functions, same elem reading:
  fmt.Println(elem_root.Key("personal").Key("list").ArrayId(2).Str())

}

