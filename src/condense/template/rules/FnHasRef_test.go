package rules

import (
	"condense/template"
	"fallbackmap"
	"reflect"
	"testing"
)

func testMakeFnHasRef(deep fallbackmap.Deep) template.Rule {
	return MakeFnHasRef(deep)
}

func TestFnHasRef_Passthru_NonMatching(t *testing.T) {
	hasRef := testMakeFnHasRef(fallbackmap.DeepNil)
	testRule_Passthru_NonMatching(hasRef, "Fn::HasRef", t)
}

func TestFnHasRef_Passthru_NonStringArguments(t *testing.T) {
	deep := fallbackmap.DeepNil

	hasRef := MakeFnHasRef(&deep)
	input := interface{}(map[string]interface{}{
		"Ref": float64(1),
	})

	newKey, newNode := hasRef([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("Ref modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("Ref with non-string argument modified the data (%#v instead of %#v)", newNode, input)
	}
}

func TestFnHasRef_Basic(t *testing.T) {
	deep := fallbackmap.NewDeepSingle([]string{"BoundValue"}, "aValue")

	hasRef := MakeFnHasRef(&deep)
	inputs := map[string]bool{
		"UnboundValue": false,
		"BoundValue":   true,
	}

	for key, expected := range inputs {
		input := interface{}(map[string]interface{}{"Fn::HasRef": key})
		newKey, newNode := hasRef([]interface{}{"x", "y"}, input)

		if newKey != "y" {
			t.Fatalf("HasRef modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, expected) {
			t.Fatalf("HasRef of %v did not return %#v (returned %#v instead)", key, expected, newNode)
		}
	}
}
