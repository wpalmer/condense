package deepalias

import "testing"
import "fallbackmap"
import "reflect"

func testGetNil(i fallbackmap.Deep, path []string, t *testing.T) {
	v, ok := i.Get(path)

	if v != nil {
		t.Fatalf("getting of %v returned nil", path)
	}
	
	if ok {
		t.Fatalf("getting of %v returned a result", path)
	}
}

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

func TestSplit(t *testing.T) {
	s := "a.b.c.[d.e.f.g].h.i.j.[k.l.m].n.o"
	split := Split(s)

	expected := []string{"a","b","c","[d.e.f.g]","h","i","j","[k.l.m]","n","o"}
	if !reflect.DeepEqual(split, expected) {
		t.Fatalf("splitting of %v did not return expected result (%v instead of %v",
			s,
			split,
			expected,
		)
	}
}

func TestSplitNested(t *testing.T) {
	s := "a.b.c.[d.[e.f].g].h.i.j.[k.l.m].n.o"
	split := Split(s)

	expected := []string{"a","b","c","[d.[e.f].g]","h","i","j","[k.l.m]","n","o"}
	if !reflect.DeepEqual(split, expected) {
		t.Fatalf("splitting of %v did not return expected result (%v instead of %v",
			s,
			split,
			expected,
		)
	}
}

func TestSplitAbormal(t *testing.T) {
	s := "a.b.c.[d.e.f.g].h.i.j.[k.l.m"
	split := Split(s)

	expected := []string{"a","b","c","[d.e.f.g]","h","i","j","[k", "l", "m"}
	if !reflect.DeepEqual(split, expected) {
		t.Fatalf("splitting of %v did not return expected result (%v instead of %v",
			s,
			split,
			expected,
		)
	}
}

func TestSplitAbormalNested(t *testing.T) {
	s := "a.b.c.[d.e.f.g].h.i.j.[k.[l.m]"
	split := Split(s)

	expected := []string{"a","b","c","[d.e.f.g]","h","i","j","[k", "[l", "m]"}
	if !reflect.DeepEqual(split, expected) {
		t.Fatalf("splitting of %v did not return expected result (%v instead of %v",
			s,
			split,
			expected,
		)
	}
}

// With no aliases involved, returns nothing (to prevent infinite recursion)
func TestGetBasic(t *testing.T) {
	i := DeepAlias{fallbackmap.DeepMap(map[string]interface{}{
		"foo": 1,
		"bar": 2,
	})}

	testGetNil(i, []string{"foo"}, t)
	testGetNil(i, []string{"bar"}, t)
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

// With no aliases involved, returns nothing (to prevent infinite recursion)
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

	testGetNil(i, []string{"a", "aa"}, t)
	testGetNil(i, []string{"a", "ab", "aba"}, t)
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
