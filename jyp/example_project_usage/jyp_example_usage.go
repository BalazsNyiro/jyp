package main

import (
	"fmt"
	"github.com/BalazsNyiro/jyp/jyp"
)

// run: go run jyp_example_usage.go

func main() {
	elem_root, errorsCollected := jyp.JsonParse(`{"personal":{"city":"Paris", "cell": 123, "money": 2.34, "list": [1,2,"third"]}}`)

	fmt.Println("errors collected:", errorsCollected)

	for key, val := range elem_root.ValObject["personal"].ValObject["list"].ValArray {
		fmt.Println("key:", key, "valueType:", val.ValType)
		// fmt.Println("key:", key)
	}

	// wanted usage: elem_root.get("personal", "list", "0").

}
