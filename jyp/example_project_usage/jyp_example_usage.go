package main

import (
	"fmt"
	"github.com/BalazsNyiro/jyp/jyp"
)

// run: go run jyp_example_usage.go

func main() {
	elem_root, errorsCollected := jyp.JsonParse(`{"personal":{"city":"Paris", "cell": 123, "money": 2.34, "list": [1,2,"third"]}}`)

	fmt.Println("errors collected:", errorsCollected)

	for key, val := range elem_root.ValObject["personal"].ValObject {
		// fmt.Println("key:", key, "val:", val)
		fmt.Println("key:", key)
		_ = val
	}

}
