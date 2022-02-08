package jyp

import "testing"

func TestObjKey(t *testing.T) {
	object_root, err := Json_parse(`{"key": 1}`)

	val_wanted := value{val_type: "number_int", val_number_int: 1}
	result_check(object_root["key"], err, val_wanted, nil, t)
}

func result_check(value_received value, err error, value_wanted value, err_wanted error, t *testing.T) {
	if value_received.val_number_int != value_wanted.val_number_int || err != err_wanted {
		t.Fatalf(`ret = %v, %v, want "", error`, value_received, err)
	}
}
