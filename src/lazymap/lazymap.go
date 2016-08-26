package lazymap

import (
	"condense/template"
	"fallbackmap"
	"strings"
)

type LazyMap struct {
	deep      fallbackmap.Deep
	processed map[string]fallbackmap.Deep
	rules     *template.Rules
}

func NewLazyMap(deep fallbackmap.Deep, rules *template.Rules) LazyMap {
	return LazyMap{
		deep:      deep,
		processed: make(map[string]fallbackmap.Deep),
		rules:     rules,
	}
}

func (lazy LazyMap) Get(path []string) (value interface{}, has_key bool) {
	if len(path) == 0 {
		return lazy, true
	}

	processed, ok := lazy.processed[strings.Join(path, ".")]
	if ok {
		value, has_key = processed.Get(path)
		if has_key {
			lazy.processed[strings.Join(path, ".")] = fallbackmap.NewDeepSingle(path, value)
			return value, true
		}

		lazy.processed[strings.Join(path, ".")] = fallbackmap.DeepNil
		return nil, false
	}

	value, has_key = lazy.deep.Get(path)
	if has_key {
		var newKey interface{}
		newKey, value = template.Walk([]interface{}{path[len(path)-1]}, value, lazy.rules)
		if newKey == nil {
			return nil, false
		}
		if newKey == path[len(path)-1] {
			lazy.processed[strings.Join(path, ".")] = fallbackmap.NewDeepSingle(path, value)
			return value, true
		}

		return nil, false
	}

	if len(path) == 1 {
		return nil, false
	}

	head := path[:len(path)-1]
	_, has_key = lazy.Get(head)
	if has_key {
		value, has_key = lazy.processed[strings.Join(head, ".")].Get(path)
		if has_key {
			lazy.processed[strings.Join(path, ".")] = fallbackmap.NewDeepSingle(path, value)
			return value, true
		}
	}

	return nil, false
}
