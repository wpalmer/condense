package rules

import (
	"testing"
	"reflect"
)

func TestFnConcat_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnConcat, "Fn::Concat", t)
}

func TestFnConcat_Passthru_NonArgsList(t *testing.T) {
	testRule_Passthru_NonArgsList(FnConcat, "Fn::Concat", t)
}

func TestFnConcat_Passthru_NonArray(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Concat": []interface{}{
			[]interface{}{1,2},
			3,
		},
	})

	newKey, newNode := FnConcat([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnConcat with an args-list containing a non-array modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnConcat with an args-list containing a non-array modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnConcat_Basic(t *testing.T) {
	testData := [][]interface{}{
		[]interface{}{ []interface{}{1,2} },
		[]interface{}{ []interface{}{1,2}, []interface{}{3,4} },
		[]interface{}{ []interface{}{1,2}, []interface{}{3}, []interface{}{4,5} },
	}

	expected := []interface{}{
		[]interface{}{1,2},
		[]interface{}{1,2,3,4},
		[]interface{}{1,2,3,4,5},
	}

	for i := range testData {
		input := interface{}(map[string]interface{}{
			"Fn::Concat": testData[i],
		})

		newKey, newNode := FnConcat([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnConcat modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, expected[i]) {
			t.Fatalf("FnConcat of args-list %v did not return the expected result (%T(%v) instead of %T(%v))", input, newNode, newNode, expected[i], expected[i])
		}
	}
}
