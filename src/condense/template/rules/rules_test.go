package rules

import (
	"testing"
	"reflect"
	"condense/template"
)

func testdataInputTypes() []interface{} {
	return []interface{} {
		interface{}("aString"),
		interface{}(true),
		interface{}(1),
		interface{}(1.0),
		interface{}(nil),
		interface{}([]interface{}{
			interface{}("aString"),
			interface{}(true),
			interface{}(1),
			interface{}(1.0),
			interface{}(nil),
			[]interface{}{
				"a", "b", "c",
			},
		}),
		interface{}(map[string]interface{}{
			"string": interface{}("aString"),
			"bool": interface{}(true),
			"int": interface{}(1),
			"float": interface{}(1.0),
			"nil": interface{}(nil),
			"array": interface{}([]interface{}{
				interface{}("aString"),
				interface{}(true),
				interface{}(1),
				interface{}(1.0),
				interface{}(nil),
				[]interface{}{
					"a", "b", "c",
				},
			}),
			"map": interface{}(map[string]interface{}{
				"a": 1, "b": 2, "c": 3,
			}),
		}),
	}
}

func testRule_Passthru_NonMatching(aRule template.Rule, singleKey string, t *testing.T) {
	for _, input := range testdataInputTypes() {
		newKey, newNode := aRule([]interface{}{"x", "y"}, input)

		if newKey != "y" {
			t.Fatalf("%s on a non-matching node modified the path (%v instead of %v)", singleKey, newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("%s on a non-matching node modified the data (%v instead of %v)", singleKey, newNode, input)
		}
	}
}

func testRule_Passthru_NonArgsList(aRule template.Rule, singleKey string, t *testing.T) {
	input := interface{}(map[string]interface{}{
		singleKey: "non-argslist",
	})

	newKey, newNode := aRule([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("%s with a non-argslist modified the path (%v instead of %v)", singleKey, newKey, "y")
	}

	if !reflect.DeepEqual(newNode, input) {
		t.Fatalf("%s with an non-argslist modified the data (%v instead of %v)", singleKey, newNode, input)
	}
}
