package fallbackmap

type FallbackMap struct {
	fallbacks []map[string]interface{}
}

func NewFallbackMap(raw map[string]interface{}) *FallbackMap {
	m := FallbackMap{
		fallbacks: []map[string]interface{}{raw},
	}
	return &m
}

func walk(p interface{}, path []string) (value interface{}, has_key bool) {
	if len(path) == 0 {
		return p, true
	}

	switch section := p.(type) {
	case map[string]interface{}:
		next, ok := section[path[0]]
		if ok {
			return walk(next, path[1:])
		}
	case *FallbackMap:
		return section.Get(path)
	}

	return nil, false
}

func (m *FallbackMap) Attach(fallback map[string]interface{}) {
	m.fallbacks = append(m.fallbacks, fallback)
}

func (m *FallbackMap) Override(fallback map[string]interface{}) {
	m.fallbacks = append([]map[string]interface{}{fallback}, m.fallbacks...)
}

func (m *FallbackMap) Get(path []string) (value interface{}, has_key bool) {
	for _, fallback := range m.fallbacks {
		found_value, did_find := walk(fallback, path)
		if did_find {
			return found_value, did_find
		}
	}

	return nil, false
}
