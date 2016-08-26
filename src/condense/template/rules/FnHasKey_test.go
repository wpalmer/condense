package rules

import (
	"reflect"
	"testing"
)

func TestFnHasKey_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnHasKey, "Fn::HasKey", t)
}

func TestFnHasKey_Passthru_NonArgsList(t *testing.T) {
	testRule_Passthru_NonArgsList(FnHasKey, "Fn::HasKey", t)
}

func TestFnHasKey_Passthru_WrongNumberOfArguments(t *testing.T) {
	inputs := []interface{}{
		[]interface{}{"a"},
		[]interface{}{"a", map[string]interface{}{}, "tooMany"},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::HasKey": input,
		})

		newKey, newNode := FnHasKey([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnHasKey modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnHasKey with wrong number of arguments modified the data (%#v instead of %#v)", newNode, input)
		}
	}
}

func TestFnHasKey_Passthru_NonMap(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::HasKey": []interface{}{"a", "NonMap"},
	})

	newKey, newNode := FnHasKey([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnHasKey with an args-list containing a non-map data modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnHasKey with an args-list containing a non-map data modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnHasKey_Passthru_NonStringKey(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::HasKey": []interface{}{1, map[string]interface{}{}},
	})

	newKey, newNode := FnJoin([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnHasKey with an args-list containing a non-string key modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnHasKey with an args-list containing a non-string key modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnHasKey_Basic(t *testing.T) {
	testData := [][]interface{}{
		[]interface{}{"a", map[string]interface{}{"a": 1}},
		[]interface{}{"a", map[string]interface{}{"b": 1}},
	}

	expected := []interface{}{
		true,
		false,
	}

	for i := range testData {
		input := interface{}(map[string]interface{}{
			"Fn::HasKey": testData[i],
		})

		newKey, newNode := FnHasKey([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnHasKey modified the path (%v instead of %v)", newKey, "y")
		}

		if newNode != expected[i] {
			t.Fatalf("FnHasKey of args-list %v did not return the expected result (%T(%v) instead of %T(%v))", input, newNode, newNode, expected[i], expected[i])
		}
	}
}
