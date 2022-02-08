package jyp

import "testing"

func TestHelloEmpty(t *testing.T) {
	object_root, err := Json_parse(`{"key": 1}`)
	result_check(object_root, err, 10, nil, t)
}

func result_check(object int, err error, object_wanted int, err_wanted error, t *testing.T) {
	if object != object_wanted || err != err_wanted {
		t.Fatalf(`ret = %v, %v, want "", error`, object, err)
	}
}
