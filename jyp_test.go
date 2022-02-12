package jyp

import "testing"

func TestObjKey(t *testing.T) {
	object_root, err := Json_parse(`{"age": 7, "friends": ["Bob", "Eve"]}`)

	val_wanted := elem{val_type: "number_int", val_number_int: 0}
	result_check(object_root, err, val_wanted, nil, t)
}

func result_check(value_received elem, err error, value_wanted elem, err_wanted error, t *testing.T) {
	if value_received.val_number_int != value_wanted.val_number_int || err != err_wanted {
		t.Fatalf(`ret = %v, %v, want "", error`, value_received, err)
	}
}
