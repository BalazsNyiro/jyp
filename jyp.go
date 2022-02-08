package jyp

import "fmt"

func main() {
	fmt.Println("hello world")
}

func Json_parse(src string) (int, error) {
	fmt.Println("json_parse:" + src)
	return len(src), nil
}
