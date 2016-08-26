package rules

import (
	"condense/template"
	"deepstack"
	"fallbackmap"
	"reflect"
	"testing"
)

func testMakeFnWith(stack deepstack.DeepStack, rules template.Rules) template.Rule {
	return MakeFnWith(&stack, &rules)
}

func TestFnWith_Passthru_NonMatching(t *testing.T) {
	fnWith := testMakeFnWith(deepstack.DeepStack{}, template.Rules{})
	testRule_Passthru_NonMatching(fnWith, "Fn::With", t)
}

func TestFnWith_Passthru_NonArgsList(t *testing.T) {
	fnWith := testMakeFnWith(deepstack.DeepStack{}, template.Rules{})
	testRule_Passthru_NonArgsList(fnWith, "Fn::With", t)
}

func TestFnWith_Passthru_WrongNumberOfArguments(t *testing.T) {
	fnWith := testMakeFnWith(deepstack.DeepStack{}, template.Rules{})

	inputs := [][]interface{}{
		[]interface{}{map[string]interface{}{"a": "one"}},
		[]interface{}{map[string]interface{}{"a": "one"}, "aTemplate", "tooMany"},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::With": input,
		})

		newKey, newNode := fnWith([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnWith modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnWith of wrong-sized args-list %v modified the data (%v instead of %v)", input, newNode, input)
		}
	}
}

func TestFnWith_Passthru_NonMapBindings(t *testing.T) {
	fnWith := testMakeFnWith(deepstack.DeepStack{}, template.Rules{})

	input := interface{}(map[string]interface{}{
		"Fn::With": []interface{}{"nonMap", "aTemplate"},
	})

	newKey, newNode := fnWith([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnWith modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnWith with arguments containing a non-map binding %v modified the data (%v instead of %v)", input, newNode, input)
	}
}

func TestFnWith_NoTemplateRules(t *testing.T) {
	fnWith := testMakeFnWith(deepstack.DeepStack{}, template.Rules{})

	input := interface{}(map[string]interface{}{
		"Fn::With": []interface{}{map[string]interface{}{"a": "one"}, "aTemplate"},
	})

	expected := interface{}("aTemplate")
	newKey, newNode := fnWith([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnWith modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnWith having no rules did not return an unmodified template (%#v instead of %#v)", newNode, expected)
	}
}

func TestFnWith_BasicRule(t *testing.T) {
	stack := deepstack.DeepStack{}
	templateRules := template.Rules{}
	templateRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 {
			key = path[len(path)-1]
		}

		return key, interface{}(map[string]interface{}{"a": "one"})
	})

	fnWith := MakeFnWith(&stack, &templateRules)
	input := interface{}(map[string]interface{}{
		"Fn::With": []interface{}{map[string]interface{}{"a": "one"}, "aTemplate"},
	})

	expected := interface{}(map[string]interface{}{"a": "one"})
	newKey, newNode := fnWith([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnWith modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnWith with did not return the expected result (%#v instead of %#v)", newNode, expected)
	}
}

func TestFnWith_PreProcessArgs(t *testing.T) {
	stack := deepstack.DeepStack{}
	templateRules := template.Rules{}
	templateRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 {
			key = path[len(path)-1]
		}

		aString, ok := node.(string)
		if !ok {
			return key, node
		}

		if aString == "skip" {
			return true, nil
		}

		if aString == "values" {
			return key, map[string]interface{}{"a": "one"}
		}

		if aString == "template" {
			_, hasKey := stack.Get([]string{"a"})
			if hasKey {
				return key, "processed-template"
			}

			return key, "preprocessed-template"
		}

		return key, node
	})

	fnWith := MakeFnWith(&stack, &templateRules)
	input := interface{}(map[string]interface{}{
		"Fn::With": []interface{}{"skip", "values", "skip", "template", "skip"},
	})

	expected := "processed-template"
	newKey, newNode := fnWith([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnWith modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnWith with did not return the expected result (%#v instead of %#v)", newNode, expected)
	}
}

func TestFnWith_StackWithArray(t *testing.T) {
	stack := deepstack.DeepStack{}
	stack.Push(fallbackmap.DeepMap(map[string]interface{}{"outer": "outer-value", "masked": "masked-value"}))

	input := interface{}(map[string]interface{}{
		"Fn::With": []interface{}{
			map[string]interface{}{
				"masked": "masking-value",
				"inner":  "inner-value",
			},
			map[string]interface{}{
				"outer":     "replace-with-outer",
				"masked":    "replace-with-masked",
				"inner":     "replace-with-inner",
				"untouched": "stay-the-same",
			},
		},
	})

	expected := interface{}(map[string]interface{}{
		"outer":     "outer-value",
		"masked":    "masking-value",
		"inner":     "inner-value",
		"untouched": "stay-the-same",
	})

	templateRules := template.Rules{}
	templateRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 {
			key = path[len(path)-1]
		}

		newNode := make(map[string]interface{})
		if nodeMap, ok := node.(map[string]interface{}); ok {
			if _, ok := node.(map[string]interface{})["untouched"]; ok {
				for key, value := range nodeMap {
					newValue, hasKey := stack.Get([]string{key})
					if hasKey {
						newNode[key] = newValue
					} else {
						newNode[key] = value
					}
				}
				return key, newNode
			}
		}

		return key, node
	})

	fnWith := MakeFnWith(&stack, &templateRules)
	newKey, newNode := fnWith([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnWith modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnWith did not have the correct stack values during templateRule (%#v instead of %#v)",
			newNode,
			expected,
		)
	}
}
