package template

import (
	"testing"
	"reflect"
	"fmt"
)

func inputTypes() []interface{} {
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

func TestPassthruNoRules(t *testing.T) {
	testRules := Rules{}

	for _, input := range inputTypes() {
		newNode := Process(input, &testRules)

		if !reflect.DeepEqual(newNode, input) {
			t.Fatalf("Walking with no rules did not return a node equal to the input (%v instead of %v)", newNode, input)
		}
	}
}

func TestOverrideData(t *testing.T) {
	testRules := Rules{}
	replacement := "overridden"
	testRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 { key = path[len(path)-1] }

		return key, replacement
	})

	for _, input := range inputTypes() {
		newNode := Process(input, &testRules)

		if !reflect.DeepEqual(newNode, replacement) {
			t.Fatalf("Walking with an \"always replace\" rule did not return a node equal to the expected replacement (%v instead of %v)", newNode, replacement)
		}
	}
}

func TestOverrideMapKey(t *testing.T) {
	testRules := Rules{}
	replacement := "overridden"
	testRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 { key = path[len(path)-1] }

		if _, ok := key.(string); ok {
			return fmt.Sprintf("%s->%s", key.(string), replacement), node
		}

		return key, node
	})

	input := interface{}(map[string]interface{}{
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
	})

	expected := interface{}(map[string]interface{}{
		"string->overridden": interface{}("aString"),
		"bool->overridden": interface{}(true),
		"int->overridden": interface{}(1),
		"float->overridden": interface{}(1.0),
		"nil->overridden": interface{}(nil),
		"array->overridden": interface{}([]interface{}{
			interface{}("aString"),
			interface{}(true),
			interface{}(1),
			interface{}(1.0),
			interface{}(nil),
			[]interface{}{
				"a", "b", "c",
			},
		}),
		"map->overridden": interface{}(map[string]interface{}{
			"a->overridden": 1, "b->overridden": 2, "c->overridden": 3,
		}),
	})

	newNode := Process(input, &testRules)
	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("Walking with an \"always replace key\" rule did not return a node equal to the expected replacement (%v instead of %v)", newNode, expected)
	}
}

func TestSkipSome(t *testing.T) {
	testRules := Rules{}
	testRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 { key = path[len(path)-1] }

		if reflect.DeepEqual(path, []interface{}{"b"}) ||
		   reflect.DeepEqual(path, []interface{}{"d", "db"}) ||
		   reflect.DeepEqual(path, []interface{}{"arr", 1}) {
			return true, nil
		}

		return key, node
	})

	input := interface{}(map[string]interface{}{
		"a": 1, "b": 2, "c": 3, "d": map[string]interface{}{"da": 1, "db": 2, "dc": 3},
		"arr": []interface{}{"x", "y", "z"},
	})

	expected := interface{}(map[string]interface{}{
		"a": 1, "c": 3, "d": map[string]interface{}{"da": 1, "dc": 3},
		"arr": []interface{}{"x", "z"},
	})

	newNode := Process(input, &testRules)
	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("Walking with an \"skip\" rule did not return a node equal to the expected replacement (%v instead of %v)", newNode, expected)
	}
}

func TestSkipSomeEarly(t *testing.T) {
	testRules := Rules{}
	testRules.AttachEarly(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 { key = path[len(path)-1] }

		if reflect.DeepEqual(path, []interface{}{"c"}) ||
		   reflect.DeepEqual(path, []interface{}{"d"}) ||
		   reflect.DeepEqual(path, []interface{}{"arr"}) {
			return true, nil
		}

		if
		   reflect.DeepEqual(path, []interface{}{"d", "da"}) ||
		   reflect.DeepEqual(path, []interface{}{"d", "db"}) ||
		   reflect.DeepEqual(path, []interface{}{"d", "dc"}) ||
		   reflect.DeepEqual(path, []interface{}{"arr", 0}) ||
		   reflect.DeepEqual(path, []interface{}{"arr", 1}) ||
		   reflect.DeepEqual(path, []interface{}{"arr", 2}) {
			t.Fatalf("Walking with an \"early skip\" rule still processed deep node %v", path)
		}

		return key, node
	})

	input := interface{}(map[string]interface{}{
		"a": 1, "b": 2, "c": 3, "d": map[string]interface{}{"da": 1, "db": 2, "dc": 3},
		"arr": []interface{}{"x", "y", "z"},
	})

	expected := interface{}(map[string]interface{}{
		"a": 1, "b": 2,
	})

	newNode := Process(input, &testRules)
	if !reflect.DeepEqual(newNode, expected) {
		t.Fatalf("Walking with an \"early skip\" rule did not return a node equal to the expected replacement (%v instead of %v)", newNode, expected)
	}
}

func TestSkipEverything(t *testing.T) {
	testRules := Rules{}
	testRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		return true, nil
	})

	for _, input := range inputTypes() {
		newNode := Process(input, &testRules)

		if !reflect.DeepEqual(newNode, nil) {
			t.Fatalf("Walking with an \"always skip\" rule returned non-nil (%v instead of %v)", newNode, nil)
		}
	}
}

func TestUnsupportedType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// do nothing
		}
	}()

	input := interface{}(map[string]string{"a": "b"})

	testRules := Rules{}
	_ = Process(input, &testRules)
}

func TestMultiple(t *testing.T) {
	testRules := Rules{}
	testRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		return nil, "replacement-one"
	})
	testRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		if node != "replacement-one" {
			t.Fatalf("Walking with a multiple rules seems to have called rules out-of-order (entered with %v instead of %v)", node, "replacement-one")
		}

		return nil, "replacement-two"
	})

	newNode := Process(interface{}("a"), &testRules)
	if !reflect.DeepEqual(newNode, "replacement-two") {
		t.Fatalf("Walking with multiple rules did not return the expected result (%v instead of %v)", newNode, "replacement-two")
	}
}

func TestMultipleEarly(t *testing.T) {
	testRules := Rules{}
	testRules.AttachEarly(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		return nil, "replacement-one"
	})
	testRules.AttachEarly(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		if node != "replacement-one" {
			t.Fatalf("Walking with a multiple (early) rules seems to have called rules out-of-order (entered with %v instead of %v)", node, "replacement-one")
		}

		return nil, "replacement-two"
	})

	newNode := Process(interface{}("a"), &testRules)
	if !reflect.DeepEqual(newNode, "replacement-two") {
		t.Fatalf("Walking with multiple (early) rules did not return the expected result (%v instead of %v)", newNode, "replacement-two")
	}
}

func TestMultipleWrap(t *testing.T) {
	testRules := Rules{}
	testRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		return nil, "replacement-one"
	})
	testRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		if node != "replacement-one" {
			t.Fatalf("Walking with a multiple (wrapped) rules seems to have called rules out-of-order (entered with %v instead of %v)", node, "replacement-one")
		}

		return nil, "replacement-two"
	})

	testRulesWrapped := Rules{}
	testRulesWrapped.Attach( testRules.MakeEach() )

	newNode := Process(interface{}("a"), &testRulesWrapped)
	if !reflect.DeepEqual(newNode, "replacement-two") {
		t.Fatalf("Walking with multiple (wrapped) rules did not return the expected result (%v instead of %v)", newNode, "replacement-two")
	}
}

func TestMultipleEarlyWrap(t *testing.T) {
	testRules := Rules{}
	testRules.AttachEarly(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		return nil, "replacement-one"
	})
	testRules.AttachEarly(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		if node != "replacement-one" {
			t.Fatalf("Walking with a multiple (early, wrapped) rules seems to have called rules out-of-order (entered with %v instead of %v)", node, "replacement-one")
		}

		return nil, "replacement-two"
	})

	testRulesWrapped := Rules{}
	testRulesWrapped.AttachEarly( testRules.MakeEachEarly() )

	newNode := Process(interface{}("a"), &testRulesWrapped)
	if !reflect.DeepEqual(newNode, "replacement-two") {
		t.Fatalf("Walking with multiple (early, wrapped) rules did not return the expected result (%v instead of %v)", newNode, "replacement-two")
	}
}
