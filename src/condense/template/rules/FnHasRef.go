package rules

import (
	"condense/template"
	"deepalias"
	"fallbackmap"
)

func MakeFnHasRef(sources fallbackmap.Deep) template.Rule {
	return func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 {
			key = path[len(path)-1]
		}

		argInterface, ok := singleKey(node, "Fn::HasRef")
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

		_, ok = sources.Get(refpath)
		if ok {
			return key, interface{}(true)
		}

		return key, interface{}(false) // (ref not found)
	}
}
