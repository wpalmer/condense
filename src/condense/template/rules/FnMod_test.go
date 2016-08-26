package rules

import (
	"reflect"
	"testing"
)

func TestFnMod_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnMod, "Fn::Mod", t)
}

func TestFnMod_Passthru_NonArgsList(t *testing.T) {
	testRule_Passthru_NonArgsList(FnMod, "Fn::Mod", t)
}

func TestFnMod_Passthru_WrongNumberOfArguments(t *testing.T) {
	inputs := [][]interface{}{
		[]interface{}{float64(15), float64(5), float64(2)},
		[]interface{}{float64(15)},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::Mod": input,
		})

		newKey, newNode := FnMod([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnMod modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnMod of wrong-sized args-list %v modified the data (%v instead of %v)", newNode, input)
		}
	}
}

func TestFnMod_Passthru_NonNumber(t *testing.T) {
	inputs := [][]interface{}{
		[]interface{}{float64(1.0), "non-number"},
		[]interface{}{"non-number", float64(1.0)},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::Mod": input,
		})

		newKey, newNode := FnMod([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnMod with an args-list containing a non-number modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnMod with an args-list containing a non-number modified the data (%v instead of %v)", newNode, input)
		}
	}
}

func TestFnMod_Basic(t *testing.T) {
	testData := [][]interface{}{
		[]interface{}{float64(5.0), float64(2.0)},
	}

	expected := []interface{}{
		float64(1.0),
	}

	for i := range testData {
		input := interface{}(map[string]interface{}{
			"Fn::Mod": testData[i],
		})

		newKey, newNode := FnMod([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnMod modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, expected[i]) {
			t.Fatalf("FnMod of args-list %v did not return the expected result (%T(%v) instead of %T(%v))", input, newNode, newNode, expected[i], expected[i])
		}
	}
}
