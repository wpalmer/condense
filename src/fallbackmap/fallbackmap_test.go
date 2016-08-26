package fallbackmap

import "testing"

func testGetInt(i Deep, path []string, value int, t *testing.T) {
	v, ok := i.Get(path)
	if !ok {
		t.Fatalf("getting of %v did not return a result", path)
	}

	vi, ok := v.(int)
	if !ok {
		t.Fatalf("getting of %v did not return an int (got %v instead)", path, v)
	}

	if vi != value {
		t.Fatalf("getting of %v did not return expected result (%v)", path, value)
	}
}

func testGetString(i Deep, path []string, value string, t *testing.T) {
	v, ok := i.Get(path)
	if !ok {
		t.Fatalf("getting of %v did not return a result", path)
	}

	vs, ok := v.(string)
	if !ok {
		t.Fatalf("getting of %v did not return an string (got %v instead)", path, v)
	}

	if vs != value {
		t.Fatalf("getting of %v did not return expected result (%v)", path, value)
	}
}

func TestGetBasic(t *testing.T) {
	i := &FallbackMap{[]Deep{DeepMap(map[string]interface{}{
		"foo": 1,
		"bar": 2,
	})}}

	testGetInt(i, []string{"foo"}, 1, t)
	testGetInt(i, []string{"bar"}, 2, t)
}

func TestGetDeep(t *testing.T) {
	i := &FallbackMap{[]Deep{DeepMap(map[string]interface{}{
		"foo": 1,
		"a": map[string]interface{}{
			"aa": 2,
			"ab": map[string]interface{}{
				"aba": 3,
				"abb": 4,
			},
		},
	})}}

	testGetInt(i, []string{"a", "aa"}, 2, t)
	testGetInt(i, []string{"a", "ab", "aba"}, 3, t)
}

func TestGetNegative(t *testing.T) {
	i := &FallbackMap{[]Deep{DeepMap(map[string]interface{}{
		"foo": 1,
		"bar": 2,
	})}}

	_, ok := i.Get([]string{"baz"})
	if ok {
		t.Fatalf("Get(...) of non-existant value reported success")
	}
}

func TestGetNegative_NamespaceCollectionEmpty(t *testing.T) {
	i := &FallbackMap{[]Deep{DeepMap(map[string]interface{}{
		"foo": 1,
		"bar": 2,
	})}}

	_, ok := i.Get([]string{"baz"})
	if ok {
		t.Fatalf(
			"Get(...) of non-existant value reported success %s",
			"(empty namespace collection defined)",
		)
	}
}

func TestAttach(t *testing.T) {
	i := &FallbackMap{[]Deep{DeepMap(map[string]interface{}{
		"foo": 1,
		"bar": 2,
	})}}

	o := map[string]interface{}{
		"bar": 3,
		"baz": 4,
	}

	i.Attach(DeepMap(o))

	testGetInt(i, []string{"foo"}, 1, t)
	testGetInt(i, []string{"bar"}, 2, t)
	testGetInt(i, []string{"baz"}, 4, t)
}

func TestDeepAttach(t *testing.T) {
	i := &FallbackMap{[]Deep{DeepMap(map[string]interface{}{
		"a": 1,
		"aStack": map[string]interface{}{
			"Outputs": map[string]interface{}{
				"overriddenOutputValue": "overridden",
			},
		},
	})}}

	o := map[string]interface{}{
		"aStack": map[string]interface{}{
			"Outputs": map[string]interface{}{
				"overriddenOutputValue": "was-not-overridden",
				"nonOverriddenOutputValue": "never-tried-to-override",
			},
		},
	}

	i.Attach(DeepMap(o))

	testGetString(i,
		[]string{"aStack", "Outputs", "overriddenOutputValue"},
		"overridden",
		t,
	)

	testGetString(i,
		[]string{"aStack", "Outputs", "nonOverriddenOutputValue"},
		"never-tried-to-override",
		t,
	)
}

func TestDeepSingle(t *testing.T) {
	deepSingle := NewDeepSingle([]string{"xx", "yy"}, map[string]interface{}{
		"a": "aValue",
		"b": "bValue",
	})

	testGetString(&deepSingle,
		[]string{"xx", "yy", "a"},
		"aValue",
		t,
	)

	testGetString(&deepSingle,
		[]string{"xx", "yy", "a"},
		"aValue",
		t,
	)

	testGetString(&deepSingle,
		[]string{"xx", "yy", "b"},
		"bValue",
		t,
	)
}
