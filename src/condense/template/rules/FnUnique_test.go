package rules

import (
	"testing"
	"reflect"
)

func TestFnUnique_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnUnique, "Fn::Unique", t)
}

func TestFnUnique_Passthru_NonArray(t *testing.T) {
	input := interface{}(map[string]interface{}{"Fn::Unique": "nonArray"})
	newKey, newNode := FnUnique([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnUnique with a non-array modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnUnique with a non-array modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnUnique_Basic(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Unique": []interface{}{
			"a",
			"b",
			1,
			map[string]interface{}{"a": 1},
			"c",
			"b",
			map[string]interface{}{"a": 2},
			"d",
			1,
			"e",
			map[string]interface{}{"a": 1},
		},
	})

	expected := []interface{}{
		"a",
		"b",
		1,
		map[string]interface{}{"a": 1},
		"c",
		map[string]interface{}{"a": 2},
		"d",
		"e",
	}

	newKey, newNode := FnUnique([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnUnique modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnUnique as %#v did not result in the expected data (%#v instead of %#v)", input, newNode, expected)
	}
}
