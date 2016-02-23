package rules

import (
	"testing"
	"reflect"
)

func TestFnSplit_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnSplit, "Fn::Split", t)
}

func TestFnSplit_Passthru_NonArgsList(t *testing.T) {
	testRule_Passthru_NonArgsList(FnSplit, "Fn::Split", t)
}

func TestFnSplit_Passthru_WrongNumberOfArguments(t *testing.T) {
	inputs := []interface{}{
		[]interface{}{","},
		[]interface{}{",", "a,string,to,split", "tooMany"},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::Split": input,
		})

		newKey, newNode := FnSplit([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnSplit modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnSplit with wrong number of arguments modified the data (%#v instead of %#v)", newNode, input)
		}
	}
}

func TestFnSplit_Passthru_NonStringData(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Split": []interface{}{",", 1},
	})

	newKey, newNode := FnSplit([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnSplit with an args-list containing a non-string data modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnSplit with an args-list containing a non-string data modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnSplit_Passthru_NonStringGlue(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Split": []interface{}{1, []interface{}{}},
	})

	newKey, newNode := FnSplit([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnSplit with an args-list containing a non-string glue modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnSplit with an args-list containing a non-string glue modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnSplit_Basic(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Split": []interface{}{",", "a,b,c,d"},
	})
	
	expected := []interface{}{"a", "b", "c", "d"}

	newKey, newNode := FnSplit([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnSplit modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnSplit did not produce the expected results (%#v instead of %#v)", newNode, expected)
	}
}
