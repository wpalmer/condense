package rules

import (
	"condense/template"
	"testing"
	"path"
	"reflect"

	"golang.org/x/tools/godoc/vfs/mapfs"
)

func testMakeFnFindFile(files map[string]string) template.Rule {
	fs := mapfs.New(files)
	return MakeFnFindFile(fs)
}

func TestFnFindFile_Passthru_NonMatching(t *testing.T) {
	fnFindFile := testMakeFnFindFile(map[string]string{})
	testRule_Passthru_NonMatching(fnFindFile, "Fn::FindFile", t)
}

func TestFnFindFile_Passthru_NonArgsList(t *testing.T) {
	fnFindFile := testMakeFnFindFile(map[string]string{})
	testRule_Passthru_NonArgsList(fnFindFile, "Fn::FindFile", t)
}

func TestFnFindFile_Passthru_NonStringListFirstArgument(t *testing.T) {
	fnFindFile := testMakeFnFindFile(map[string]string{"1": "{\"content\": 1}"})

	inputs := []interface{}{
		"nonList",
		[]interface{}{"valid", float64(1)},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::FindFile": []interface{}{input, "tail.json"},
		})

		newKey, newNode := fnFindFile([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnFindFile modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnFindFile with non-string-list first argument modified the data (%#v instead of %#v)", newNode, input)
		}
	}
}

func TestFnFindFile_Passthru_NonStringTail(t *testing.T) {
	fnFindFile := testMakeFnFindFile(map[string]string{"a/theFile": "content"})

	input := interface{}(map[string]interface{}{
		"Fn::FindFile": []interface{}{[]interface{}{"a"}, float64(1)},
	})

	newKey, newNode := fnFindFile([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnFindFile modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("FnFindFile with non-string tail argument modified the data (%#v instead of %#v)", newNode, input)
	}
}

func TestFnFindFile_Passthru_WrongNumberOfArguments(t *testing.T) {
	fnFindFile := testMakeFnFindFile(map[string]string{"a/theFile": "content"})

	inputs := []interface{}{
		[]interface{}{[]interface{}{"a"}},
		[]interface{}{[]interface{}{"a"}, "theFile", "tooMany"},
	}

	for _, input := range inputs {
		input := interface{}(map[string]interface{}{
			"Fn::FindFile": input,
		})

		newKey, newNode := fnFindFile([]interface{}{"x", "y"}, input)
		if newKey != "y" {
			t.Fatalf("FnFindFile modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("FnFindFile with wrong number of arguments modified the data (%#v instead of %#v)", newNode, input)
		}
	}
}

func TestFnFindFile_Panic_BadFilename(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// do nothing
		}
	}()

	fnFindFile := testMakeFnFindFile(map[string]string{"b/theFile": "contents"})

	input := interface{}(map[string]interface{}{
		"Fn::FindFile": []interface{}{[]interface{}{"/a", "/b"}, "nonexistantFile"},
	})

	_, _ = fnFindFile([]interface{}{"x", "y"}, input)
	t.Fatalf("Searching for a non-existant file did not panic")
}

func TestFnFindFile_Basic(t *testing.T) {
	fnFindFile := testMakeFnFindFile(map[string]string{
		"b/theFile": "contents",
		"c/theFile": "contents",
	})
	input := interface{}(map[string]interface{}{
		"Fn::FindFile": []interface{}{[]interface{}{"/a", "/b", "/c"}, "theFile"},
	})

	_, _ = fnFindFile([]interface{}{"x", "y"}, input)

	expected := interface{}(path.Join("/b", "theFile"))
	newKey, newNode := fnFindFile([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("FnFindFile modified the path (%v instead of %v)", newKey, "y")
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("FnFindFile did not return the expected result (%#v instead of %#v)", newNode, expected)
	}
}
