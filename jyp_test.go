package jyp

import "testing"

func TestHelloEmpty(t *testing.T) {
	node_root, err := Json_parse("{\"key\": 1}")
	result_check(node_root, err, 10, nil, t)
}

func result_check(node int, err error, node_wanted int, err_wanted error, t *testing.T) {
	if node != node_wanted || err != err_wanted {
		t.Fatalf(`ret = %v, %v, want "", error`, node, err)
	}
}
