package rules

import (
	"testing"
	"reflect"
)

func TestFnLength_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnLength, "Fn::Length", t)
}

func TestFnLength_Passthru_NonList(t *testing.T) {
	testRule_Passthru_NonArgsList(FnLength, "Fn::Length", t)
}

func TestFnLength_Basic(t *testing.T) {
	testData := [][]interface{}{
		[]interface{}{},
		[]interface{}{"a","b","c"},
	}

	expected := []interface{}{
		interface{}(float64(0)),
		interface{}(float64(3)),
	}

	for i := range testData {
		input := interface{}(map[string]interface{}{
			"Fn::Length": testData[i],
		})

		newKey, newNode := FnLength([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnLength modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, expected[i]) {
			t.Fatalf("FnLength of list %v did not return the expected result (%T(%v) instead of %T(%v))", input, newNode, newNode, expected[i], expected[i])
		}
	}
}
