package rules

import (
	"condense/template"
	"deepstack"
	"fallbackmap"
	"testing"
	"reflect"
)

func testMakeFnFor(stack deepstack.DeepStack, rules template.Rules) template.Rule {
	return MakeFnFor(&stack, &rules)
}

func TestFnFor_Passthru_NonMatching(t *testing.T) {
	fnFor := testMakeFnFor(deepstack.DeepStack{}, template.Rules{})
	testRule_Passthru_NonMatching(fnFor, "Fn::For", t)
}

func TestFnFor_Passthru_NonArgsList(t *testing.T) {
	fnFor := testMakeFnFor(deepstack.DeepStack{}, template.Rules{})
	testRule_Passthru_NonArgsList(fnFor, "Fn::For", t)
}

func TestFnFor_Passthru_WrongNumberOfArguments(t *testing.T) {
	fnFor := testMakeFnFor(deepstack.DeepStack{}, template.Rules{})

	inputs := [][]interface{}{
		[]interface{}{"aRefName"},
		[]interface{}{"aRefName", []interface{}{1,2}},
		[]interface{}{"aRefName", []interface{}{1,2}, "aTemplate", "tooMany"},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::For": input,
		})

		newKey, newNode := fnFor([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnFor modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnFor of wrong-sized args-list %v modified the data (%v instead of %v)", newNode, input)
		}
	}
}

func TestFnFor_Passthru_BadRefNames(t *testing.T) {
	fnFor := testMakeFnFor(deepstack.DeepStack{}, template.Rules{})

	inputs := [][]interface{}{
		[]interface{}{ []interface{}{}, []interface{}{1}, "aTemplate" },
		[]interface{}{ []interface{}{"key", "value", "tooMany"}, []interface{}{1}, "aTemplate" },
		[]interface{}{ []interface{}{"key", 1}, []interface{}{1}, "aTemplate" },
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::For": input,
		})

		newKey, newNode := fnFor([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnFor modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnFor with arguments containing a bad refNames list %v modified the data (%v instead of %v)", input, newNode, input)
		}
	}
}

func TestFnFor_Passthru_BadValues(t *testing.T) {
	fnFor := testMakeFnFor(deepstack.DeepStack{}, template.Rules{})

	input := interface{}(map[string]interface{}{
		"Fn::For": []interface{}{ []interface{}{"key", "value"}, "badValue", "aTemplate" },
	})

	newKey, newNode := fnFor([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnFor modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnFor with arguments containing a bad values list %v modified the data (%v instead of %v)", input, newNode, input)
	}
}

func TestFnFor_EmptyValues(t *testing.T) {
	fnFor := testMakeFnFor(deepstack.DeepStack{}, template.Rules{})

	input := interface{}(map[string]interface{}{
		"Fn::For": []interface{}{ []interface{}{"key", "value"}, []interface{}{}, "aTemplate" },
	})

	expected := []interface{}{}
	newKey, newNode := fnFor([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnFor modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnFor with empty values did not return an empty array (%#v instead of %#v)", newNode, expected)
	}
}

func TestFnFor_NoTemplateRules(t *testing.T) {
	fnFor := testMakeFnFor(deepstack.DeepStack{}, template.Rules{})

	input := interface{}(map[string]interface{}{
		"Fn::For": []interface{}{ []interface{}{"key", "value"}, []interface{}{"a", "b"}, "aTemplate" },
	})

	expected := []interface{}{"aTemplate", "aTemplate"}
	newKey, newNode := fnFor([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnFor modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnFor with no template rules did not return an array of unmodified templates (%#v instead of %#v)", newNode, expected)
	}
}

func TestFnFor_BasicRule(t *testing.T) {
	stack := deepstack.DeepStack{}
	templateRules := template.Rules{}
	templateRules.Attach( func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 { key = path[len(path)-1] }

		if key == 1 {
			return key, interface{}("replaced")
		}

		return key, node
	} )
	
	fnFor := MakeFnFor(&stack, &templateRules)
	input := interface{}(map[string]interface{}{
		"Fn::For": []interface{}{ []interface{}{"key", "value"}, []interface{}{"a", "b"}, "aTemplate" },
	})

	expected := []interface{}{"aTemplate", "replaced"}
	newKey, newNode := fnFor([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnFor modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnFor with did not return the expected result (%#v instead of %#v)", newNode, expected)
	}
}

func TestFnFor_PreProcessArgs(t *testing.T) {
	stack := deepstack.DeepStack{}
	templateRules := template.Rules{}
	templateRules.Attach( func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 { key = path[len(path)-1] }

		aString, ok := node.(string)
		if !ok {
			return key, node
		}

		if aString == "skip" {
			return true, nil
		}

		if aString == "refNames" {
			return key, []interface{}{"key", "value"}
		}

		if aString == "values" {
			return key, []interface{}{"a"}
		}

		if aString == "template" {
			_, hasKey := stack.Get([]string{"value"})
			if hasKey {
				return key, "processed-template"
			}

			return key, "preprocessed-template"
		}

		return key, node
	} )
	
	fnFor := MakeFnFor(&stack, &templateRules)
	input := interface{}(map[string]interface{}{
		"Fn::For": []interface{}{ "skip", "refNames", "skip", "values", "skip", "template", "skip" },
	})

	expected := []interface{}{"processed-template"}
	newKey, newNode := fnFor([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnFor modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnFor with did not return the expected result (%#v instead of %#v)", newNode, expected)
	}
}

func TestFnFor_StackWithArray(t *testing.T) {
	stack := deepstack.DeepStack{}
	stack.Push( fallbackmap.DeepMap(map[string]interface{}{"outer": "outerValue", "masked": "outerMasked"}) )

	testRefNames := []interface{}{
		[]interface{}{"key", "masked"},
		[]interface{}{nil, "masked"},
		[]interface{}{"key", nil},
		[]interface{}{nil, nil},
		[]interface{}{"masked"},
		"masked",
	}

	expected := []interface{}{
		[]interface{}{
			map[string]interface{}{
				"outer": []interface{}{"outerValue", true},
				"masked": []interface{}{"innerMasking", true},
				"key": []interface{}{float64(0), true},
			},
		},
		[]interface{}{
			map[string]interface{}{
				"outer": []interface{}{"outerValue", true},
				"masked": []interface{}{"innerMasking", true},
			},
		},
		[]interface{}{
			map[string]interface{}{
				"outer": []interface{}{"outerValue", true},
				"masked": []interface{}{"outerMasked", true},
				"key": []interface{}{float64(0), true},
			},
		},
		[]interface{}{
			map[string]interface{}{
				"outer": []interface{}{"outerValue", true},
				"masked": []interface{}{"outerMasked", true},
			},
		},
		[]interface{}{
			map[string]interface{}{
				"outer": []interface{}{"outerValue", true},
				"masked": []interface{}{"innerMasking", true},
			},
		},
		[]interface{}{
			map[string]interface{}{
				"outer": []interface{}{"outerValue", true},
				"masked": []interface{}{"innerMasking", true},
			},
		},
	}

	for i, refNames := range testRefNames {
		input := interface{}(map[string]interface{}{
			"Fn::For": []interface{}{
				refNames,
				[]interface{}{"innerMasking"},
				"aTemplate",
			},
		})

		templateRules := template.Rules{}
		templateRules.Attach( func(path []interface{}, node interface{}) (interface{}, interface{}) {
			key := interface{}(nil)
			if len(path) > 0 { key = path[len(path)-1] }
			if stringVal, ok := node.(string); !ok || stringVal != "aTemplate" {
				return key, node
			}

			generated := map[string]interface{}{}
			for binding, _ := range expected[i].([]interface{})[0].(map[string]interface{}) {
				value, has_key := stack.Get([]string{binding})
				generated[ binding ] = []interface{}{ value, has_key }
			}

			return key, generated
		} )

		fnFor := MakeFnFor(&stack, &templateRules)
		newKey, newNode := fnFor([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnFor modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, expected[i]) {
			t.Fatalf("FnFor did not have the correct stack values with refNames %v during templateRule (%#v instead of %#v)",
				refNames,
				newNode,
				expected[i],
			)
		}
	}
}
