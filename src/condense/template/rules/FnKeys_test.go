package rules

import (
	"testing"
	"reflect"
)

func TestFnKeys_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnKeys, "Fn::Keys", t)
}

func TestFnKeys_Passthru_NonMap(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Keys": 0,
	})

	newKey, newNode := FnKeys([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnKeys modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnKeys of non-map arg modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnKeys_Basic(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Keys": map[string]interface{}{"a": 1, "b": 2, "c": 3},
	})

	newKey, newNode := FnKeys([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnKeys modified the path (%v instead of %v)", newKey, "y")
	}

	newNodeArray, ok := newNode.([]interface{})
	if !ok {
		t.Fatalf("FnKeys of %v did not return an array (%v instead)", input, newNode)
	}

	if len(newNodeArray) != 3 {
		t.Fatalf("FnKeys of %v did not return all input keys (%v instead)", input, newNode)
	}

	resultMap := map[string]bool{}
	for _, v := range newNodeArray {
		vString, ok := v.(string)
		if !ok {
			t.Fatalf("FnKeys of %v returned a non-string key (%v)", input, v)
		}

		resultMap[ vString ] = true
	}

	expected := map[string]bool{
		"a": true,
		"b": true,
		"c": true,
	}

	if !reflect.DeepEqual(resultMap, expected) {
		t.Fatalf("FnKeys of %v did not return the expected keys (%v instead of %v)", input, resultMap, expected)
	}
}
