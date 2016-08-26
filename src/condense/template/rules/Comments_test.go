package rules

import (
	"reflect"
	"testing"
)

func inputTypes() []interface{} {
	return []interface{}{
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
			"bool":   interface{}(true),
			"int":    interface{}(1),
			"float":  interface{}(1.0),
			"nil":    interface{}(nil),
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

func TestExcludeComments_Passthru_NonComment(t *testing.T) {
	for _, input := range testdataInputTypes() {
		newKey, newNode := ExcludeComments([]interface{}{"x", "y"}, input)

		if newKey != "y" {
			t.Fatalf("ExcludeComments on a non-comment modified the path (%v instead of %v)", newKey, "y")
		}

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("ExcludeComments on a non-comment modified the data (%v instead of %v)", newNode, input)
		}
	}
}

func TestExcludeComments_FromMap(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"a":        "noExclude-a",
		"$comment": "exclude",
		"b":        "noExclude-b",
	})

	expected := interface{}(map[string]interface{}{
		"a": "noExclude-a",
		"b": "noExclude-b",
	})

	newKey, newNode := ExcludeComments([]interface{}{"x", "y"}, input)
	if newKey != "y" {
		t.Fatalf("ExcludeComments on a map containing a $comment key modified the path (%v instead of %v)", newKey, nil)
	}

	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("ExcludeComments on a map containing a $comment key did not return a map excluding the key (%v instead of %v)", newNode, expected)
	}
}

func TestExcludeComments_Full(t *testing.T) {
	input := interface{}(map[string]interface{}{
		"$comment": "exclude",
	})

	newKey, _ := ExcludeComments([]interface{}{}, input)
	if skip, ok := newKey.(bool); !ok || !skip {
		t.Fatalf("ExcludeComments on a map containing only a $comment key did not output a skip (%v,%v instead of %v,%v)", skip, ok, true, true)
	}
}
