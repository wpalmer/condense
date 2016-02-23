package rules

import (
	"testing"
	"reflect"
)

func TestFnToEntries_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnToEntries, "Fn::ToEntries", t)
}

func TestFnToEntries_Passthru_NonMap(t *testing.T) {
	input := interface{}(map[string]interface{}{"Fn::ToEntries": "nonMap"})
	newKey, newNode := FnToEntries([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("ToEntries with a non-map modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("ToEntries with a non-map modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnToEntries_Basic(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::ToEntries": map[string]interface{}{"a": 1, "b": "foo", "c": 3.0},
	})

	expected := map[string]interface{}{
		"a": map[string]interface{}{"key": "a", "value": 1},
		"b": map[string]interface{}{"key": "b", "value": "foo"},
		"c": map[string]interface{}{"key": "c", "value": 3.0},
	}

	newKey, newNode := FnToEntries([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("ToEntries modified the path (%v instead of %v)", newKey, "y")
	}

	if _, ok := newNode.([]interface{}); !ok {
		t.Fatalf("ToEntries did not result in an array (%#v instead)", newNode)
	}

	for _, expected_data := range expected {
		found := false
		for _, newNodeEntry := range newNode.([]interface{}) {
			if reflect.DeepEqual(newNodeEntry, expected_data) {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("ToEntries as %#v did not result in the expected data (%#v not present in %#v)", input, expected_data, newNode)
		}
	}
}
