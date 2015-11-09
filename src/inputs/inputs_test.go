package inputs

import "testing"

func testGetInt(i *Inputs, path []string, value int, t *testing.T) {
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

func testGetString(i *Inputs, path []string, value string, t *testing.T) {
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
	i := NewInputs(map[string]interface{}{
		"foo": 1,
		"bar": 2,
	})

	testGetInt(i, []string{"foo"}, 1, t)
	testGetInt(i, []string{"bar"}, 2, t)
}

func TestGetDeep(t *testing.T) {
	i := NewInputs(map[string]interface{}{
		"foo": 1,
		"a": map[string]interface{}{
			"aa": 2,
			"ab": map[string]interface{}{
				"aba": 3,
				"abb": 4,
			},
		},
	})

	testGetInt(i, []string{"a", "aa"}, 2, t)
	testGetInt(i, []string{"a", "ab", "aba"}, 3, t)
}

func TestGetNegative(t *testing.T) {
	i := NewInputs(map[string]interface{}{
		"foo": 1,
		"bar": 2,
	})

	_, ok := i.Get([]string{"baz"})
	if ok {
		t.Fatalf("Get(...) of non-existant value reported success")
	}
}

func TestGetNegative_NamespaceCollectionInvalid(t *testing.T) {
	i := NewInputs(map[string]interface{}{
		"foo": 1,
		"bar": 2,
		"::": "InvalidCollection",
	})

	_, ok := i.Get([]string{"baz"})
	if ok {
		t.Fatalf(
			"Get(...) of non-existant value reported success %s",
			"(invalid namespace collection defined)",
		)
	}
}

func TestGetNegative_NamespaceCollectionEmpty(t *testing.T) {
	i := NewInputs(map[string]interface{}{
		"foo": 1,
		"bar": 2,
		"::": map[string]interface{}{},
	})

	_, ok := i.Get([]string{"baz"})
	if ok {
		t.Fatalf(
			"Get(...) of non-existant value reported success %s",
			"(empty namespace collection defined)",
		)
	}
}

func TestMerge(t *testing.T) {
	i := NewInputs(map[string]interface{}{
		"foo": 1,
		"bar": 2,
	})
	
	o := NewInputs(map[string]interface{}{
		"bar": 3,
		"baz": 4,
	})
	
	i.Merge(o)

	testGetInt(i, []string{"foo"}, 1, t)
	testGetInt(i, []string{"bar"}, 3, t)
	testGetInt(i, []string{"baz"}, 4, t)
}

func TestAttach(t *testing.T) {
	i := NewInputs(map[string]interface{}{
		"foo": 1,
		"bar": 2,
	})
	
	o := NewInputs(map[string]interface{}{
		"bar": 3,
		"baz": 4,
	})
	
	i.Attach("A", o)
	
	if ns, _ := i.Namespace("A"); ns != o {
		t.Fatalf("Retrieved namespace was not the same one which was attached")
	}

	testGetInt(i, []string{"foo"}, 1, t)
	testGetInt(i, []string{"bar"}, 2, t)
	testGetInt(i, []string{"baz"}, 4, t)
	testGetInt(i, []string{"::", "A", "bar"}, 3, t)
}

func TestNamespaceNegative_NoNamespaces(t *testing.T){
	i := NewInputs(map[string]interface{}{
		"foo": 1,
		"bar": 2,
	})

	if _, ok := i.Namespace("A"); ok {
		t.Fatalf("Namespace(...) reported success with no namespaces defined")
	}
}

func TestNamespaceNegative_MissingNamespace(t *testing.T){
	i := NewInputs(map[string]interface{}{
		"foo": 1,
		"bar": 2,
	})
	
	o := NewInputs(map[string]interface{}{
		"bar": 3,
		"baz": 4,
	})

	i.Attach("A", o)

	if _, ok := i.Namespace("B"); ok {
		t.Fatalf("Namespace(...) reported success when namespace not defined")
	}
}

func TestNamespaceNegative_NotANamespaceCollection(t *testing.T){
	i := NewInputs(map[string]interface{}{
		"foo": 1,
		"bar": 2,
		"::": "NotANamespaceCollection",
	})

	if _, ok := i.Namespace("A"); ok {
		t.Fatalf("Namespace(...) reported success when namespaces not a map")
	}
}

func TestNamespaceNegative_NotANamespace(t *testing.T){
	i := NewInputs(map[string]interface{}{
		"foo": 1,
		"bar": 2,
		"::": map[string]interface{}{
			"A": "NotANamespace",
		},
	})

	if _, ok := i.Namespace("A"); ok {
		t.Fatalf("Namespace(...) reported success when namespace not an Inputs")
	}
}

func TestDeepAttach(t *testing.T) {
	i := NewInputs(map[string]interface{}{
		"a": 1,
		"aStack": map[string]interface{}{
			"Outputs": map[string]interface{}{
				"overriddenOutputValue": "overridden",
			},
		},
	})

	o := NewInputs(map[string]interface{}{
		"aStack": map[string]interface{}{
			"Outputs": map[string]interface{}{
				"overriddenOutputValue": "was-not-overridden",
				"nonOverriddenOutputValue": "never-tried-to-override",
			},
		},
	})

	i.Attach("aStack", o)

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
