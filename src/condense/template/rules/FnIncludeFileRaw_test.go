package rules

import (
	"condense/template"
	"reflect"
	"testing"

	"golang.org/x/tools/godoc/vfs/mapfs"
)

func testMakeFnIncludeFileRaw(files map[string]string) template.Rule {
	fs := mapfs.New(files)
	return MakeFnIncludeFileRaw(fs)
}

func TestMakeFnIncludeFileRaw_Passthru_NonMatching(t *testing.T) {
	fnIncludeFileRaw := testMakeFnIncludeFileRaw(map[string]string{})
	testRule_Passthru_NonMatching(fnIncludeFileRaw, "Fn::IncludeFileRaw", t)
}

func TestMakeFnIncludeFileRaw_Passthru_NonStringArgument(t *testing.T) {
	fnIncludeFileRaw := testMakeFnIncludeFileRaw(map[string]string{"1": "{\"content\": 1}"})

	input := interface{}(map[string]interface{}{
		"Fn::IncludeFileRaw": float64(1),
	})

	newKey, newNode := fnIncludeFileRaw([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnIncludeFileRaw modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnIncludeFileRaw with non-string argument modified the data (%#v instead of %#v)", newNode, input)
	}
}

func TestMakeFnIncludeFileRaw_Panic_BadFilename(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// do nothing
		}
	}()

	fnIncludeFileRaw := testMakeFnIncludeFileRaw(map[string]string{"a": "{\"content\": 1}"})

	input := interface{}(map[string]interface{}{
		"Fn::IncludeFileRaw": "/b",
	})

	_, _ = fnIncludeFileRaw([]interface{}{"x", "y"}, input)
	t.Fatalf("Including a non-existant file did not panic")
}

func TestMakeFnIncludeFileRaw_Basic(t *testing.T) {
	fnIncludeFileRaw := testMakeFnIncludeFileRaw(map[string]string{"a": "File contents"})

	input := interface{}(map[string]interface{}{
		"Fn::IncludeFileRaw": "/a",
	})

	expected := "File contents"
	newKey, newNode := fnIncludeFileRaw([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnIncludeFileRaw modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnIncludeFileRaw with no template rules did not return an unmodified file contents (%#v instead of %#v)", newNode, expected)
	}
}
