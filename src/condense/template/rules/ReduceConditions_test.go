package rules

import (
	"reflect"
	"testing"
)

func TestReduceConditions_Passthru_NonCondition(t *testing.T) {
	paths := [][]interface{}{
		[]interface{}{"x"},
		[]interface{}{"x", "y"},
		[]interface{}{"Conditions"},
		[]interface{}{"Conditions", float64(1.0)},
		[]interface{}{float64(1.0), "Conditions"},
	}

	input := interface{}(true)

	for _, path := range paths {
		newKey, newNode := ReduceConditions(path, input)

		if newKey != path[len(path)-1] {
			t.Fatalf("ReduceConditions on a non-Condition node modified the path (%v instead of %v)", newKey, path[len(path)-1])
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("ReduceConditions on a non-Condition node modified the data (%v instead of %v)", newNode, input)
		}
	}
}

func TestReduceConditions_Passthru_NonBool(t *testing.T) {
	input := interface{}("nonBool")
	newKey, newNode := ReduceConditions([]interface{}{"Conditions", "Foo"}, input)

	if newKey != "Foo" {
		t.Fatalf("ReduceConditions on a non-bool node modified the path (%v instead of %v)", newKey, "Foo")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("ReduceConditions on a non-bool node modified the data (%v instead of %v)", newNode, input)
	}
}

func TestReduceConditions_Basic(t *testing.T) {
	inputs := []interface{}{
		true,
		false,
	}

	expected := []interface{}{
		map[string]interface{}{"Fn::Equals": interface{}([]interface{}{"1", "1"})},
		map[string]interface{}{"Fn::Equals": interface{}([]interface{}{"0", "1"})},
	}

	for i, input := range inputs {
		newKey, newNode := ReduceConditions([]interface{}{"Conditions", "Foo"}, input)

		if newKey != "Foo" {
			t.Fatalf("ReduceConditions modified the path (%v instead of %v)", newKey, "Foo")
		}

		if !reflect.DeepEqual(newNode, expected[i]) {
			t.Fatalf("ReduceConditions of '%#v' did not produce the expected results (%#v instead of %#v)", input, newNode, expected[i])
		}
	}
}
