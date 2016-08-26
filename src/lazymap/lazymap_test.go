package lazymap

import (
	"testing"
	//	"reflect"
	"condense/template"
	"fallbackmap"
	"strings"
)

func Test_Lazy_callsDeepGet(t *testing.T) {
	templateRules := template.Rules{}
	lazy := NewLazyMap(fallbackmap.DeepFunc(func(path []string) (value interface{}, has_key bool) {
		return "aValue", true
	}), &templateRules)

	for path := range [][]string{[]string{"a"}, []string{"a", "b"}, []string{"a", "b", "c"}} {
		result, has_key := lazy.Get([]string{"anything"})
		if !has_key {
			t.Fatalf("Get() of a map which always returns did not claim to have a key for path %#v", path)
		}

		resultString, ok := result.(string)
		if !ok {
			t.Fatalf("Get() of a map which always returns did not return the correct type of value for path %#v", path)
		}

		if resultString != "aValue" {
			t.Fatalf("Get() of a map which always returns did not return the correct value for path %#v", path)
		}
	}
}

func Test_Lazy_appliesTemplate(t *testing.T) {
	testRules := template.Rules{}
	testRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 {
			key = path[len(path)-1]
		}

		if key == "replaceThis" {
			return key, "replacement"
		}

		return key, node
	})

	lazy := NewLazyMap(fallbackmap.DeepMap(map[string]interface{}{
		"a": "anA",
		"b": interface{}(map[string]interface{}{
			"innerA":      "anInnerA",
			"replaceThis": "this should be replaced",
		}),
		"replaceThis": "this should also be replaced",
	}), &testRules)

	for path, expected := range map[string]string{
		"a":             "anA",
		"b.innerA":      "anInnerA",
		"b.replaceThis": "replacement",
		"replaceThis":   "replacement",
	} {
		result, has_key := lazy.Get(strings.Split(path, "."))
		if !has_key {
			t.Fatalf("Get() of a map which always returns did not claim to have a key for path %#v", path)
		}

		resultString, ok := result.(string)
		if !ok {
			t.Fatalf("Get() of a map which always returns did not return the correct type of value for path %#v", path)
		}

		if resultString != expected {
			t.Fatalf("Get() of a map which always returns did not return the expected value '%s' (got '%s' instead) for path %#v", expected, resultString, path)
		}
	}
}

func Test_Lazy_traversesTemplateRuleResults(t *testing.T) {
	testRules := template.Rules{}
	testRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 {
			key = path[len(path)-1]
		}

		if key == "replaceThis" {
			return key, map[string]interface{}{
				"replacedDeep": "replacedDeepValue",
			}
		}

		return key, node
	})

	lazy := NewLazyMap(fallbackmap.DeepMap(map[string]interface{}{
		"a": "anA",
		"b": interface{}(map[string]interface{}{
			"innerA":      "anInnerA",
			"replaceThis": "this should be replaced",
		}),
		"replaceThis": "this should also be replaced",
	}), &testRules)

	for path, expected := range map[string]string{
		"b.replaceThis.replacedDeep": "replacedDeepValue",
		"replaceThis.replacedDeep":   "replacedDeepValue",
	} {
		result, has_key := lazy.Get(strings.Split(path, "."))
		if !has_key {
			t.Fatalf("Get() of a rule-modified map did not claim to have a key for path %#v", path)
		}

		resultString, ok := result.(string)
		if !ok {
			t.Fatalf("Get() of a rule-modified map did not return the correct type of value for path %#v", path)
		}

		if resultString != expected {
			t.Fatalf("Get() of a rule-modified map did not return the expected value '%s' (got '%s' instead) for path %#v", expected, resultString, path)
		}
	}
}
