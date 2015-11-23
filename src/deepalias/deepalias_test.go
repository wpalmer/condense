package deepalias

import "testing"
import "fallbackmap"

func testGetInt(i fallbackmap.Deep, path []string, value int, t *testing.T) {
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

func testGetString(i fallbackmap.Deep, path []string, value string, t *testing.T) {
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

// With no aliases involved, acts as fallbackmap.Deep
func TestGetBasic(t *testing.T) {
	i := DeepAlias{fallbackmap.DeepMap(map[string]interface{}{
		"foo": 1,
		"bar": 2,
	})}

	testGetInt(i, []string{"foo"}, 1, t)
	testGetInt(i, []string{"bar"}, 2, t)
}

func TestGetAliased(t *testing.T) {
	i := DeepAlias{fallbackmap.DeepMap(map[string]interface{}{
		"foo": 1,
		"bar": 2,
		"anAlias": "foo",
		"anotherAlias": "bar",
	})}

	testGetInt(i, []string{"[anAlias]"}, 1, t)
	testGetInt(i, []string{"[anotherAlias]"}, 2, t)
}

// With no aliases involved, acts as fallbackmap.Deep
func TestGetDeepBasic(t *testing.T) {
	i := DeepAlias{fallbackmap.DeepMap(map[string]interface{}{
		"foo": 1,
		"a": map[string]interface{}{
			"aa": 2,
			"ab": map[string]interface{}{
				"aba": 3,
				"abb": 4,
			},
		},
	})}

	testGetInt(i, []string{"a", "aa"}, 2, t)
	testGetInt(i, []string{"a", "ab", "aba"}, 3, t)
}

func TestGetDeepAlias(t *testing.T) {
	i := DeepAlias{fallbackmap.DeepMap(map[string]interface{}{
		"foo": 1,
		"a": map[string]interface{}{
			"aa": 2,
			"ab": map[string]interface{}{
				"aba": 3,
				"anAlias": "bar",
			},
			"bar": 4,
		},
	})}

	testGetInt(i, []string{"a", "[a.ab.anAlias]"}, 4, t)
}
