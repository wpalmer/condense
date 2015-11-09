package inputs

// A key->value store, to which namespaced values can be added.
//
// When looking up a value (via Get([]string)), the primary namespace is
// searched, followed by all other "Attached" namespaces. In this way,
// additional sources of Input data can be Attached indiscriminently after
// initial data has been loaded, while allowing the initial data to take
// priority.
type Inputs map[string]interface{}

func merge(a map[string]interface{}, b map[string]interface{}) {
	for k, spec := range b {
		a[k] = spec
	}
}

func NewInputs(raw map[string]interface{}) *Inputs {
	inputs := Inputs(raw)
	return &inputs
}

func (inputs *Inputs) Map() map[string]interface{} {
	return map[string]interface{}(*inputs)
}

func (inputs *Inputs) Merge(other *Inputs) {
	merge(map[string]interface{}(*inputs), map[string]interface{}(*other))
}

func (inputs *Inputs) Attach(namespace string, other *Inputs) {
	m := inputs.Map()

	_, ok := m["::"]
	if !ok {
		m["::"] = make(map[string]interface{})
	}

	m["::"].(map[string]interface{})[namespace] = other
}

func (inputs *Inputs) Namespace(namespace string) (*Inputs, bool) {
	m := inputs.Map()
	namespaces_generic, ok := m["::"]
	if !ok {
		return nil, false
	}

	namespaces, ok := namespaces_generic.(map[string]interface{})
	if !ok {
		return nil, false
	}

	namespaced_generic, ok := namespaces[namespace]
	if !ok {
		return nil, false
	}

	namespaced, ok := namespaced_generic.(*Inputs)
	if !ok {
		return nil, false
	}

	return namespaced, true
}

func get(p interface{}, path []string) (interface{}, bool) {
	if len(path) == 0 {
		return p, true
	}

	switch section := p.(type) {
	case map[string]interface{}:
		next, ok := section[path[0]]
		if ok {
			final, ok := get(next, path[1:])
			if ok {
				return final, true
			}
		}
	case *Inputs:
		cast := map[string]interface{}(*section)
		final, ok := get(cast, path)
		if ok {
			return final, true
		}
	}

	// Path not found within "this" namespace, so search for Namespaces

	if path[0] == "::" {
		// If we were already in a Namespace search, it has failed.
		return nil, false
	}
	
	namespaces_generic, ok := get(p, []string{"::"})
	if !ok {
		// no Namespaces map defined, so path has not been found.
		return nil, false
	}

	namespaces, ok := namespaces_generic.(map[string]interface{})
	if !ok {
		// Namespaces map was not the expected type, so path has not been found.
		return nil, false
	}

	for _, namespace := range namespaces {
		final, ok := get(namespace, path)
		if ok {
			return final, true
		}
	}

	// Path found in neither "this" namespace, nor defined sub-namespaces
	return nil, false
}

// Note: intentionally does not handle lookup-by-array-index
func (inputs *Inputs) Get(path []string) (interface{}, bool) {
	return get(inputs, path)
}
