package rules

import (
	"testing"
	"reflect"
)

func TestFnMergeDeep_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnMergeDeep, "Fn::MergeDeep", t)
}

func TestFnMergeDeep_Passthru_NonArgsList(t *testing.T) {
	testRule_Passthru_NonArgsList(FnMergeDeep, "Fn::MergeDeep", t)
}

func TestFnMergeDeep_Passthru_WrongNumberOfArguments(t *testing.T) {
	inputs := []interface{}{
		[]interface{}{float64(1)},
		[]interface{}{float64(1), map[string]interface{}{"a": "firstValue"}, "tooMany"},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::MergeDeep": input,
		})

		newKey, newNode := FnMergeDeep([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnMergeDeep modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnMergeDeep with wrong number of arguments modified the data (%#v instead of %#v)", newNode, input)
		}
	}
}

func TestFnMergeDeep_Passthru_NonFloatDepth(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::MergeDeep": []interface{}{
			"nonFloat",
			[]interface{}{
				map[string]interface{}{"a": "firstValue"},
				map[string]interface{}{"b": "secondValue"},
			},
		},
	})

	newKey, newNode := FnMergeDeep([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnMergeDeep with an args-list containing a non-array-of-maps data modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnMergeDeep with an args-list containing a non-array-of-maps data modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnMergeDeep_Passthru_NonMapArray(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::MergeDeep": []interface{}{
			float64(0),
			"nonArray",
		},
	})

	newKey, newNode := FnMergeDeep([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnMergeDeep with an args-list containing a non-array-of-maps data modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnMergeDeep with an args-list containing a non-array-of-maps data modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnMergeDeep_Passthru_NonMap(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::MergeDeep": []interface{}{
			float64(0),
			[]interface{}{
				map[string]interface{}{"a": "firstValue"},
				"nonMap",
			},
		},
	})

	newKey, newNode := FnMergeDeep([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnMergeDeep with an args-list containing a non-map data modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnMergeDeep with an args-list containing a non-map data modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnMergeDeep_Basic(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::MergeDeep": []interface{}{
			float64(0),
			[]interface{}{
				map[string]interface{}{"a": "firstValue", "b": "maskedValue"},
				map[string]interface{}{"b": "maskingValue", "c": "addedValue"},
				map[string]interface{}{"d": "otherAddedValue"},
			},
		},
	})
	
	expected := interface{}(map[string]interface{}{
		"a": "firstValue",
		"b": "maskingValue",
		"c": "addedValue",
		"d": "otherAddedValue",
	})

	newKey, newNode := FnMergeDeep([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnMergeDeep modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnMergeDeep did not result in the expected data (%#v instead of %#v)", newNode, expected)
	}
}

func TestFnMergeDeep_Depth(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"a": "firstValue",
			"b": "maskedValue",
			"deepOne": map[string]interface{}{
				"a": "firstValue",
				"b": "maskedValue",
				"deepTwo": map[string]interface{}{
					"a": "firstValue",
					"b": "maskedValue",
				},
			},
		},
		map[string]interface{}{
			"b": "maskingValue",
			"c": "addedValue",
			"deepOne": map[string]interface{}{
				"b": "maskingValue",
				"c": "addedValue",
				"deepTwo": map[string]interface{}{
					"b": "maskingValue",
					"c": "addedValue",
				},
			},
		},
		map[string]interface{}{
			"d": "otherAddedValue",
			"deepOne": map[string]interface{}{
				"d": "otherAddedValue",
				"deepTwo": map[string]interface{}{
					"d": "otherAddedValue",
				},
			},
		},
	}
	
	expected := []interface{}{
		map[string]interface{}{
			"a": "firstValue",
			"b": "maskingValue",
			"c": "addedValue",
			"d": "otherAddedValue",
			"deepOne": map[string]interface{}{
				"d": "otherAddedValue",
				"deepTwo": map[string]interface{}{
					"d": "otherAddedValue",
				},
			},
		},
		map[string]interface{}{
			"a": "firstValue",
			"b": "maskingValue",
			"c": "addedValue",
			"d": "otherAddedValue",
			"deepOne": map[string]interface{}{
				"a": "firstValue",
				"b": "maskingValue",
				"c": "addedValue",
				"d": "otherAddedValue",
				"deepTwo": map[string]interface{}{
					"d": "otherAddedValue",
				},
			},
		},
		map[string]interface{}{
			"a": "firstValue",
			"b": "maskingValue",
			"c": "addedValue",
			"d": "otherAddedValue",
			"deepOne": map[string]interface{}{
				"a": "firstValue",
				"b": "maskingValue",
				"c": "addedValue",
				"d": "otherAddedValue",
				"deepTwo": map[string]interface{}{
					"a": "firstValue",
					"b": "maskingValue",
					"c": "addedValue",
					"d": "otherAddedValue",
				},
			},
		},
	}

	for i, expected_value := range expected {
		newKey, newNode := FnMergeDeep(
			[]interface{}{"x", "y"},
			interface{}(map[string]interface{}{
				"Fn::MergeDeep": []interface{}{ interface{}(float64(i)), input },
			}),
		)

		if newKey != "y" {
			t.Fatalf("FnMergeDeep [%d] modified the path (%v instead of %v)", i, newKey, "y")
		}

		if !reflect.DeepEqual(newNode, expected_value) {
			t.Fatalf("FnMergeDeep [%d] did not result in the expected data (\n%v\n instead of \n%v\n)", i, newNode, expected_value)
		}
	}
}
