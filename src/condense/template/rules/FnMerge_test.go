package rules

import (
	"testing"
	"reflect"
)

func TestFnMerge_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnMerge, "Fn::Merge", t)
}

func TestFnMerge_Passthru_NonArgsList(t *testing.T) {
	testRule_Passthru_NonArgsList(FnMerge, "Fn::Merge", t)
}

func TestFnMerge_Passthru_NonMap(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Merge": []interface{}{
			map[string]interface{}{"a": "firstValue"},
			"nonMap",
		},
	})

	newKey, newNode := FnMerge([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnMerge with an args-list containing a non-map data modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnMerge with an args-list containing a non-map data modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnMerge_Basic(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Merge": []interface{}{
			map[string]interface{}{"a": "firstValue", "b": "maskedValue"},
			map[string]interface{}{"b": "maskingValue", "c": "addedValue"},
			map[string]interface{}{"d": "otherAddedValue"},
		},
	})
	
	expected := interface{}(map[string]interface{}{
		"a": "firstValue",
		"b": "maskingValue",
		"c": "addedValue",
		"d": "otherAddedValue",
	})

	newKey, newNode := FnMerge([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnMerge modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnMerge did not result in the expected data (%#v instead of %#v)", newNode, expected)
	}
}
