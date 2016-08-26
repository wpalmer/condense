package rules

import (
	"condense/template"
	"deepstack"
	"fallbackmap"
	"reflect"
	"testing"
)

func testMakeRef(stack deepstack.DeepStack, rules template.Rules) template.Rule {
	return MakeRef(&stack, &rules)
}

func TestRef_Passthru_NonMatching(t *testing.T) {
	ref := testMakeRef(deepstack.DeepStack{}, template.Rules{})
	testRule_Passthru_NonMatching(ref, "Ref", t)
}

func TestRef_Passthru_Unbound(t *testing.T) {
	stack := deepstack.DeepStack{}
	templateRules := template.Rules{}

	ref := MakeRef(&stack, &templateRules)
	input := interface{}(map[string]interface{}{
		"Ref": "UnboundValue",
	})

	newKey, newNode := ref([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("Ref modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("Ref with no bound variables modified the data (%#v instead of %#v)", newNode, input)
	}
}

func TestRef_Passthru_NonStringArguments(t *testing.T) {
	stack := deepstack.DeepStack{}
	templateRules := template.Rules{}

	ref := MakeRef(&stack, &templateRules)
	input := interface{}(map[string]interface{}{
		"Ref": float64(1),
	})

	newKey, newNode := ref([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("Ref modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("Ref with non-string argument modified the data (%#v instead of %#v)", newNode, input)
	}
}

func TestRef_Basic(t *testing.T) {
	stack := deepstack.DeepStack{}
	stack.Push(fallbackmap.DeepMap(map[string]interface{}{
		"BoundVar": "BoundValue",
	}))
	templateRules := template.Rules{}

	expected := "BoundValue"

	ref := MakeRef(&stack, &templateRules)
	input := interface{}(map[string]interface{}{
		"Ref": "BoundVar",
	})

	newKey, newNode := ref([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("Ref modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("Ref for %v did not return the expected result (%#v instead of %#v)", input, newNode, expected)
	}
}
