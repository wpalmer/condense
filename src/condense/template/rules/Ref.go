package rules

import (
	"condense/template"
	"deepalias"
	"fallbackmap"
)

func MakeRef(sources fallbackmap.Deep, rules *template.Rules) template.Rule {
	return func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 {
			key = path[len(path)-1]
		}

		argInterface, ok := singleKey(node, "Ref")
		if !ok {
			return key, node //passthru
		}

		var argString string
		if argString, ok = argInterface.(string); !ok {
			return key, node //passthru
		}

		var refpath []string
		for _, part := range deepalias.Split(argString) {
			refpath = append(refpath, part)
		}

		var newNode interface{}
		newNode, ok = sources.Get(refpath)
		if ok {
			var newKey interface{}
			newKey, newNode = template.Walk(path, newNode, rules)
			return newKey, newNode
		}

		return key, node //passthru (ref not found)
	}
}
