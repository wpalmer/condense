package rules

import (
	"testing"
	"reflect"
)

func TestFnIf_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnIf, "Fn::If", t)
}

func TestFnIf_Passthru_NonArgsList(t *testing.T) {
	testRule_Passthru_NonArgsList(FnIf, "Fn::If", t)
}

func TestFnIf_Passthru_WrongNumberOfArguments(t *testing.T) {
	inputs := [][]interface{}{
		[]interface{}{true, "aString"},
		[]interface{}{true, "aString", "aNotherString", "aThirdString"},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::If": input,
		})

		newKey, newNode := FnIf([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnIf modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnIf of wrong-sized args-list %v modified the data (%v instead of %v)", newNode, input)
		}
	}
}

func TestFnIf_Passthru_NonBooleanFirst(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::If": []interface{}{"aNonBool", "aString", "aNotherString"},
	})

	newKey, newNode := FnIf([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnIf with non-boolean first argument modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnIf with non-boolean first argument modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnIf_Basic(t *testing.T) {
	inputs := [][]interface{}{
		[]interface{}{true, "TrueValue", "FalseValue"},
		[]interface{}{false, "TrueValue", "FalseValue"},
	}

	expected := []interface{}{
		"TrueValue",
		"FalseValue",
	}

	for i, input := range inputs {
		newKey, newNode := FnIf([]interface{}{"x", "y"}, map[string]interface{}{"Fn::If": input})
		if newKey != "y" {
			t.Fatalf("FnIf modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, expected[i]) {
			t.Fatalf("FnIf with args-list %v did not return the the expected value (%v instead of %v)", input, newNode, expected[i])
		}
	}
}

func TestFnEquals_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnEquals, "Fn::Equals", t)
}

func TestFnEquals_Passthru_NonArgsList(t *testing.T) {
	testRule_Passthru_NonArgsList(FnEquals, "Fn::Equals", t)
}

func TestFnEquals_Passthru_WrongNumberOfArguments(t *testing.T) {
	inputs := [][]interface{}{
		[]interface{}{1},
		[]interface{}{1, 2, 3},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::Equals": input,
		})

		newKey, newNode := FnEquals([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnEquals modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnEquals of wrong-sized args-list %v modified the data (%v instead of %v)", newNode, input)
		}
	}
}

func TestFnEquals_Basic(t *testing.T) {
	inputs := [][]interface{}{
		[]interface{}{1, 1},
		[]interface{}{1, 2},
		[]interface{}{
			map[string]interface{}{"a": 1, "b": map[string]interface{}{"c": 2, "d": 3}},
			map[string]interface{}{"a": 1, "b": map[string]interface{}{"c": 2, "d": 3}},
		},
		[]interface{}{
			map[string]interface{}{"a": 1, "b": map[string]interface{}{"c": 2, "d": 3}},
			map[string]interface{}{"a": 1, "b": map[string]interface{}{"c": 2, "d": 4}},
		},
	}

	expected := []interface{}{
		true,
		false,
		true,
		false,
	}

	for i, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::Equals": input,
		})

		newKey, newNode := FnEquals([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnEquals modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, expected[i]) {
			t.Fatalf("FnEquals of %v did not return the expected result (%v instead of %v)", input, newNode, expected[i])
		}
	}
}

func TestFnAnd_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnAnd, "Fn::And", t)
}

func TestFnAnd_Passthru_NonArgsList(t *testing.T) {
	testRule_Passthru_NonArgsList(FnAnd, "Fn::And", t)
}

func TestFnAnd_Passthru_WrongNumberOfArguments(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::And": []interface{}{true},
	})

	newKey, newNode := FnAnd([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnAnd modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnAnd of wrong-sized args-list %v modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnAnd_Passthru_NonBool(t *testing.T) {
	inputs := [][]interface{}{
		[]interface{}{1, 1},
		[]interface{}{true, 1},
		[]interface{}{true, true, 1},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::And": input,
		})

		newKey, newNode := FnAnd([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnAnd modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnAnd of non-bool args modified the data (%v instead of %v)", newNode, input)
		}
	}
}

func TestFnAnd_Basic(t *testing.T) {
	inputs := [][]interface{}{
		[]interface{}{false, false},
		[]interface{}{false, true},
		[]interface{}{true, false},
		[]interface{}{true, true},
		[]interface{}{false, false, false},
		[]interface{}{false, false, true},
		[]interface{}{false, true, false},
		[]interface{}{false, true, true},
		[]interface{}{true, false, false},
		[]interface{}{true, false, true},
		[]interface{}{true, true, false},
		[]interface{}{true, true, true},
	}

	expected := []interface{}{
		false,
		false,
		false,
		true,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		true,
	}

	for i, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::And": input,
		})

		newKey, newNode := FnAnd([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnAnd modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, expected[i]) {
			t.Fatalf("FnAnd of %v did not return the expected result (%v instead of %v)", input, newNode, expected[i])
		}
	}
}

func TestFnOr_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnOr, "Fn::Or", t)
}

func TestFnOr_Passthru_NonArgsList(t *testing.T) {
	testRule_Passthru_NonArgsList(FnOr, "Fn::Or", t)
}

func TestFnOr_Passthru_WrongNumberOfArguments(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Or": []interface{}{true},
	})

	newKey, newNode := FnOr([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnOr modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnOr of wrong-sized args-list %v modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnOr_Passthru_NonBool(t *testing.T) {
	inputs := [][]interface{}{
		[]interface{}{1, 1},
		[]interface{}{true, 1},
		[]interface{}{true, true, 1},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::Or": input,
		})

		newKey, newNode := FnOr([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnOr modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnOr of non-bool args modified the data (%v instead of %v)", newNode, input)
		}
	}
}

func TestFnOr_Basic(t *testing.T) {
	inputs := [][]interface{}{
		[]interface{}{false, false},
		[]interface{}{false, true},
		[]interface{}{true, false},
		[]interface{}{true, true},
		[]interface{}{false, false, false},
		[]interface{}{false, false, true},
		[]interface{}{false, true, false},
		[]interface{}{false, true, true},
		[]interface{}{true, false, false},
		[]interface{}{true, false, true},
		[]interface{}{true, true, false},
		[]interface{}{true, true, true},
	}

	expected := []interface{}{
		false,
		true,
		true,
		true,
		false,
		true,
		true,
		true,
		true,
		true,
		true,
		true,
	}

	for i, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::Or": input,
		})

		newKey, newNode := FnOr([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnOr modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, expected[i]) {
			t.Fatalf("FnOr of %v did not return the expected result (%v instead of %v)", input, newNode, expected[i])
		}
	}
}

func TestFnNot_Passthru_NonMatching(t *testing.T) {
	testRule_Passthru_NonMatching(FnNot, "Fn::Not", t)
}

func TestFnNot_Passthru_NonBool(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"Fn::Not": 0,
	})

	newKey, newNode := FnNot([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnNot modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnNot of non-bool arg modified the data (%v instead of %v)", newNode, input)
	}
}

func TestFnNot_Basic(t *testing.T) {
	inputs := []interface{}{
		false,
		true,
	}

	expected := []interface{}{
		true,
		false,
	}

	for i, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::Not": input,
		})

		newKey, newNode := FnNot([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnNot modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, expected[i]) {
			t.Fatalf("FnNot of %v did not return the expected result (%v instead of %v)", input, newNode, expected[i])
		}
	}
}
