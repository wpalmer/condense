package fallbackmap

type Deep interface {
	Get(path []string) (value interface{}, has_key bool)
}

type DeepMap map[string]interface{}
type DeepFunc func(path []string) (value interface{}, has_key bool)

func (deep DeepFunc) Get(path []string) (value interface{}, has_key bool) {
	return deep(path)
}

func (deep DeepMap) Get(path []string) (value interface{}, has_key bool) {
	if len(path) == 0 {
		return deep, true
	}

	next, ok := map[string]interface{}(deep)[path[0]]
	if !ok {
		return nil, false
	}

	if len(path[1:]) == 0 {
		return next, true
	}

	var next_deep Deep
	var next_map map[string]interface{}

	next_map, ok = next.(map[string]interface{})
	if ok {
		next_deep = DeepMap(next_map)
	} else {
		next_deep, ok = next.(Deep)
		if !ok {
			return nil, false
		}
	}

	return next_deep.Get(path[1:])
}

type FallbackMap struct {
	fallbacks []Deep
}

func (m *FallbackMap) Attach(fallback Deep) {
	m.fallbacks = append(m.fallbacks, fallback)
}

func (m *FallbackMap) Override(fallback Deep) {
	m.fallbacks = append([]Deep{fallback}, m.fallbacks...)
}

func (m *FallbackMap) Get(path []string) (value interface{}, has_key bool) {
	for _, fallback := range m.fallbacks {
		found_value, did_find := fallback.Get(path)
		if did_find {
			return found_value, did_find
		}
	}

	return nil, false
}
