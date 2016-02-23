package rules

import (
	"testing"
	"reflect"
)

func TestFnAdd_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnAdd, "Fn::Add", t)
}

func TestFnAdd_Passthru_NonArgsList(t *testing.T) {
	testRule_Passthru_NonArgsList(FnAdd, "Fn::Add", t)
}

func TestFnAdd_Passthru_NonNumber(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Add": []interface{}{1.0, "non-number"},
	})

	newKey, newNode := FnAdd([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnAdd with an args-list containing a non-number modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnAdd with an args-list containing a non-number modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnAdd_Basic(t *testing.T) {
	testData := [][]interface{}{
		[]interface{}{1.0},
		[]interface{}{1.0, 2.0},
		[]interface{}{1.0, 2.1, 3.3},
		[]interface{}{1.0, 2.1, 3.3, 4.3},
	}

	expected := []interface{}{
		float64(1.0),
		float64(3.0),
		float64(6.4),
		float64(10.7),
	}

	for i := range testData {
		input := interface{}(map[string]interface{}{
			"Fn::Add": testData[i],
		})

		newKey, newNode := FnAdd([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnAdd modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, expected[i]) {
			t.Fatalf("FnAdd of args-list %v did not return the expected result (%T(%v) instead of %T(%v))", input, newNode, newNode, expected[i], expected[i])
		}
	}
}
