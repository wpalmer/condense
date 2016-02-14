package template
import (
	"fmt"
)

type Rule func(path []interface{}, node interface{}) (newKey interface{}, newNode interface{})
type Rules struct {
	Early []Rule
	Depth []Rule
}

func (r *Rules) AttachEarly(rule Rule) {
	r.Early = append(r.Early, rule)
}

func (r *Rules) Attach(rule Rule) {
	r.Depth = append(r.Depth, rule)
}

func eachRule(path []interface{}, node interface{}, rules []Rule) (newKey interface{}, newNode interface{}) {
	newPath := make([]interface{}, len(path))
	copy(newPath, path)

	newNode = node
	newKey = interface{}(nil)
	if len(newPath) > 0 {
		newKey = newPath[len(newPath)-1]
	}
	for _, rule := range rules {
		newKey, newNode = rule(newPath, newNode)
		if skip, ok := newKey.(bool); ok && skip {
			return true, nil
		}

		if len(newPath) > 0 {
			newPath[len(newPath)-1] = newKey
		}
	}

	return newKey, newNode
}

func (r *Rules) MakeEachEarly() Rule {
	return func (path []interface{}, node interface{}) (newKey interface{}, newNode interface{}) {
		return eachRule(path, node, r.Early)
	}
}

func (r *Rules) MakeEach() Rule {
	return func (path []interface{}, node interface{}) (newKey interface{}, newNode interface{}) {
		return eachRule(path, node, r.Depth)
	}
}

func Walk(path []interface{}, node interface{}, rules *Rules) (newKey interface{}, newNode interface{}) {
	newPath := make([]interface{}, len(path))
	copy(newPath, path)

	newNode = node
	newKey = interface{}(nil)
	if len(newPath) > 0 {
		newKey = newPath[len(newPath)-1]
	}

	newKey, newNode = eachRule(newPath, newNode, rules.Early)
	if skip, ok := newKey.(bool); ok && skip {
		return true, nil
	}

	if len(newPath) > 0 {
		newPath[len(newPath)-1] = newKey
	}

	switch typed := newNode.(type) {
	default:
		panic(fmt.Sprintf("unknown type: %T\n", typed))
	case []interface{}:
		filtered := []interface{}{}
		for deepIndex, deepNode := range typed {
			newDeepPath := make([]interface{}, len(newPath)+1)
			copy(newDeepPath, newPath)
			newDeepPath[cap(newDeepPath)-1] = deepIndex

			newDeepIndex := interface{}(deepIndex)
			newDeepNode := deepNode
			newDeepIndex, newDeepNode = Walk(newDeepPath, newDeepNode, rules)

			if skip, ok := newDeepIndex.(bool); !ok || !skip {
				filtered = append(filtered, newDeepNode)
			}
		}

		newNode = interface{}(filtered)
	case map[string]interface{}:
		filtered := make(map[string]interface{})
		for deepKey, deepNode := range typed {
			newDeepPath := make([]interface{}, len(newPath)+1)
			copy(newDeepPath, newPath)
			newDeepPath[cap(newDeepPath)-1] = deepKey

			newDeepKey := interface{}(deepKey)
			newDeepNode := deepNode
			newDeepKey, newDeepNode = Walk(newDeepPath, newDeepNode, rules)

			if skip, ok := newDeepKey.(bool); !ok || !skip {
				filtered[newDeepKey.(string)] = newDeepNode
			}
		}

		newNode = interface{}(filtered)
	case string:
	case bool:
	case int:
	case float64:
	case nil:
	}

	newKey, newNode = eachRule(newPath, newNode, rules.Depth)
	return newKey, newNode
}

func Process(node interface{}, rules *Rules) interface{} {
	emptyPath := []interface{}{}
	newKey, processed := Walk(emptyPath, node, rules)

	if skip, ok := newKey.(bool); ok && skip {
		return interface{}(nil)
	}

	return processed
}
