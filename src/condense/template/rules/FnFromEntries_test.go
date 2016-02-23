package rules

import (
	"testing"
	"reflect"
)

func TestFnFromEntries_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnFromEntries, "Fn::FromEntries", t)
}

func TestFnFromEntries_Passthru_NonArgsList(t *testing.T) {
	testRule_Passthru_NonArgsList(FnFromEntries, "Fn::FromEntries", t)
}

func TestFnFromEntries_Passthru_NonKVMap(t *testing.T) {
	inputs := []interface{}{
		map[string]interface{}{"nonKey": "a", "value": "b"},
		map[string]interface{}{"key": "a", "nonValue": "b"},
		map[string]interface{}{"key": 1, "value": "b"},
		"nonMap",
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::FromEntries": []interface{}{
				map[string]interface{}{"key": "validKey", "value": "validVal"},
				input,
			},
		})

		newKey, newNode := FnFromEntries([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FromEntries with an args-list containing a non-map{key:...,value:...} modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FromEntries with an args-list containing a non-map{key:...,value:...} modified the data (%v instead of %v)", newNode, input)
		}
	}
}

func TestFnFromEntries_Basic(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::FromEntries": []interface{}{
			map[string]interface{}{"key": "firstKey", "value": "firstValue"},
			map[string]interface{}{"key": "secondKey", "value": "secondValue"},
		},
	})

	expected := map[string]interface{}{
		"firstKey": "firstValue",
		"secondKey": "secondValue",
	}

	newKey, newNode := FnFromEntries([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FromEntries with modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FromEntries as %#v did not result in the expected data (%#v instead of %#v)", input, newNode, expected)
	}
}
