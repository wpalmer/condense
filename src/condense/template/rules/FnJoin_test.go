package rules

import (
	"testing"
	"reflect"
)

func TestFnJoin_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnJoin, "Fn::Join", t)
}

func TestFnJoin_Passthru_NonArgsList(t *testing.T) {
	testRule_Passthru_NonArgsList(FnJoin, "Fn::Join", t)
}

func TestFnJoin_Passthru_WrongNumberOfArguments(t *testing.T) {
	inputs := []interface{}{
		[]interface{}{","},
		[]interface{}{",", []interface{}{}, "tooMany"},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::Join": input,
		})

		newKey, newNode := FnJoin([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnJoin modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnJoin with wrong number of arguments modified the data (%#v instead of %#v)", newNode, input)
		}
	}
}

func TestFnJoin_Passthru_NonDataList(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Join": []interface{}{",", "NonList"},
	})

	newKey, newNode := FnJoin([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnJoin with an args-list containing a non-list data modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnJoin with an args-list containing a non-list data modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnJoin_Passthru_NonStringGlue(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Join": []interface{}{1, []interface{}{}},
	})

	newKey, newNode := FnJoin([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnJoin with an args-list containing a non-string glue modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnJoin with an args-list containing a non-string glue modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnJoin_Passthru_NonBasicData(t *testing.T) {
	inputs := []interface{}{
		[]interface{}{",", []interface{}{[]interface{}{1}}},
		[]interface{}{",", []interface{}{map[string]interface{}{"a":1}}},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::Join": input,
		})

		newKey, newNode := FnJoin([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnJoin modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnJoin with wrong number of arguments modified the data (%#v instead of %#v)", newNode, input)
		}
	}
}

func TestFnJoin_Basic(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Join": []interface{}{",", []interface{}{"a", 1, float64(1.5)}},
	})
	
	expected := "a,1,1"

	newKey, newNode := FnJoin([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnJoin modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnJoin did not produce the expected results (%#v instead of %#v)", newNode, expected)
	}
}
