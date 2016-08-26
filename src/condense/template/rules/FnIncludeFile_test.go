package rules

import (
	"condense/template"
	"reflect"
	"testing"

	"golang.org/x/tools/godoc/vfs/mapfs"
)

func testMakeFnIncludeFile(files map[string]string, rules template.Rules) template.Rule {
	fs := mapfs.New(files)
	return MakeFnIncludeFile(fs, &rules)
}

func TestMakeFnIncludeFile_Passthru_NonMatching(t *testing.T) {
	fnIncludeFile := testMakeFnIncludeFile(map[string]string{}, template.Rules{})
	testRule_Passthru_NonMatching(fnIncludeFile, "Fn::IncludeFile", t)
}

func TestMakeFnIncludeFile_Passthru_NonStringArgument(t *testing.T) {
	fnIncludeFile := testMakeFnIncludeFile(map[string]string{"1": "{\"content\": 1}"}, template.Rules{})

	input := interface{}(map[string]interface{}{
		"Fn::IncludeFile": float64(1),
	})

	newKey, newNode := fnIncludeFile([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnIncludeFile modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnIncludeFile with non-string argument modified the data (%#v instead of %#v)", newNode, input)
	}
}

func TestMakeFnIncludeFile_Panic_BadFilename(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// do nothing
		}
	}()

	fnIncludeFile := testMakeFnIncludeFile(map[string]string{"a": "{\"content\": 1}"}, template.Rules{})

	input := interface{}(map[string]interface{}{
		"Fn::IncludeFile": "/b",
	})

	_, _ = fnIncludeFile([]interface{}{"x", "y"}, input)
	t.Fatalf("Including a non-existant file did not panic")
}

func TestMakeFnIncludeFile_Panic_NonJSON(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// do nothing
		}
	}()

	fnIncludeFile := testMakeFnIncludeFile(map[string]string{"a": "nonJSON"}, template.Rules{})

	input := interface{}(map[string]interface{}{
		"Fn::IncludeFile": "/a",
	})

	_, _ = fnIncludeFile([]interface{}{"x", "y"}, input)
	t.Fatalf("Including a non-JSON file did not panic")
}

func TestMakeFnIncludeFile_NoTemplateRules(t *testing.T) {
	fnIncludeFile := testMakeFnIncludeFile(map[string]string{"a": "{\"content\": 1}"}, template.Rules{})

	input := interface{}(map[string]interface{}{
		"Fn::IncludeFile": "/a",
	})

	expected := interface{}(map[string]interface{}{"content": float64(1)})
	newKey, newNode := fnIncludeFile([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnIncludeFile modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnIncludeFile with no template rules did not return an unmodified file contents (%#v instead of %#v)", newNode, expected)
	}
}

func TestMakeFnIncludeFile_BasicRule(t *testing.T) {
	templateRules := template.Rules{}
	templateRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 {
			key = path[len(path)-1]
		}

		if key == "content" {
			return key, interface{}("replaced")
		}

		return key, node
	})

	fs := mapfs.New(map[string]string{"a": "{\"content\": 1}"})
	fnIncludeFile := MakeFnIncludeFile(fs, &templateRules)
	input := interface{}(map[string]interface{}{
		"Fn::IncludeFile": "/a",
	})

	expected := interface{}(map[string]interface{}{"content": "replaced"})
	newKey, newNode := fnIncludeFile([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnIncludeFile modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnIncludeFile did not return the expected result (%#v instead of %#v)", newNode, expected)
	}
}
