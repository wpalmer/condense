package rules

import (
	"testing"
	"reflect"
	"condense/template"
	"deepstack"
	"fallbackmap"
)

func testMakeGetAtt(stack deepstack.DeepStack, rules template.Rules) template.Rule {
	return MakeFnGetAtt(&stack, &rules)
}

func TestFnGetAtt_Passthru_NonMatching(t *testing.T) {
	fnGetAtt := testMakeGetAtt(deepstack.DeepStack{}, template.Rules{})
	testRule_Passthru_NonMatching(fnGetAtt, "Fn::GetAtt", t)
}

func TestFnGetAtt_Passthru_NonArgsList(t *testing.T) {
	fnGetAtt := testMakeGetAtt(deepstack.DeepStack{}, template.Rules{})
	testRule_Passthru_NonArgsList(fnGetAtt, "Fn::GetAtt", t)
}

func TestFnGetAtt_Passthru_Unbound(t *testing.T) {
	stack := deepstack.DeepStack{}
	templateRules := template.Rules{}

	fnGetAtt := MakeFnGetAtt(&stack, &templateRules)
	input := interface{}(map[string]interface{}{
		"Fn::GetAtt": []interface{}{"Unbound", "Value"},
	})

	newKey, newNode := fnGetAtt([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnGetAtt modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnGetAtt with no bound variables modified the data (%#v instead of %#v)", newNode, input)
	}
}

func TestFnGetAtt_Passthru_WrongNumberOfArguments(t *testing.T) {
	stack := deepstack.DeepStack{}
	stack.Push( fallbackmap.DeepMap(map[string]interface{}{
		"FakeResource": map[string]interface{}{
			"FakeProperty": "FakeValue",
		},
	}) )
	templateRules := template.Rules{}

	fnGetAtt := MakeFnGetAtt(&stack, &templateRules)

	inputs := []interface{}{
		[]interface{}{"FakeResource"},
		[]interface{}{"FakeResource", "FakeProperty", "tooMany"},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::GetAtt": input,
		})

		newKey, newNode := fnGetAtt([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnGetAtt modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnGetAtt with wrong number of arguments modified the data (%#v instead of %#v)", newNode, input)
		}
	}
}

func TestFnGetAtt_Passthru_NonStringArguments(t *testing.T) {
	stack := deepstack.DeepStack{}
	stack.Push( fallbackmap.DeepMap(map[string]interface{}{
		"FakeResource": map[string]interface{}{
			"FakeProperty": "FakeValue",
		},
	}) )
	templateRules := template.Rules{}
	
	inputs := []interface{}{
		[]interface{}{"FakeResource", 1},
		[]interface{}{1, "FakeProperty"},
	}

	for _, input := range inputs {
		fnGetAtt := MakeFnGetAtt(&stack, &templateRules)
		input := interface{}(map[string]interface{}{
			"Fn::GetAtt": input,
		})

		newKey, newNode := fnGetAtt([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnGetAtt modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnGetAtt with non-string arguments modified the data (%#v instead of %#v)", newNode, input)
		}
	}
}

func TestFnGetAtt_Basic(t *testing.T) {
	stack := deepstack.DeepStack{}
	stack.Push( fallbackmap.DeepMap(map[string]interface{}{
		"FakeResource": map[string]interface{}{
			"FakeProperty": map[string]interface{}{"FakeSub": "FakeValue"},
		},
	}) )
	templateRules := template.Rules{}
	
	inputs := []interface{}{
		[]interface{}{"FakeResource", "FakeProperty"},
		[]interface{}{"FakeResource", "FakeProperty.FakeSub"},
	}

	expected := []interface{}{
		map[string]interface{}{"FakeSub": "FakeValue"},
		"FakeValue",
	}

	for i, input := range inputs {
		fnGetAtt := MakeFnGetAtt(&stack, &templateRules)
		input := interface{}(map[string]interface{}{
			"Fn::GetAtt": input,
		})

		newKey, newNode := fnGetAtt([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnGetAtt modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, expected[i]) {
			t.Fatalf("FnGetAtt for %v did not return the expected result (%#v instead of %#v)", input, newNode, expected[i])
		}
	}
}

func TestFnGetAtt_ProcessBound(t *testing.T) {
	stack := deepstack.DeepStack{}
	stack.Push( fallbackmap.DeepMap(map[string]interface{}{
		"FakeResource": map[string]interface{}{
			"FakeProperty": map[string]interface{}{"FakeSub": "FakeValue"},
		},
	}) )
	templateRules := template.Rules{}
	templateRules.Attach( func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 { key = path[len(path)-1] }

		if stringVal, ok := node.(string); ok && stringVal == "FakeValue" {
			return key, interface{}("ProcessedFakeValue")
		}

		return key, node
	} )
	
	inputs := []interface{}{
		[]interface{}{"FakeResource", "FakeProperty.FakeSub"},
	}

	expected := []interface{}{
		"ProcessedFakeValue",
	}

	for i, input := range inputs {
		fnGetAtt := MakeFnGetAtt(&stack, &templateRules)
		input := interface{}(map[string]interface{}{
			"Fn::GetAtt": input,
		})

		newKey, newNode := fnGetAtt([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnGetAtt modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, expected[i]) {
			t.Fatalf("FnGetAtt for %v did not return the expected result (%#v instead of %#v)", input, newNode, expected[i])
		}
	}
}
