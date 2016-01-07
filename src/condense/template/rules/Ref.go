package rules

import (
	"fallbackmap"
	"deepalias"
	"condense/template"
)

func MakeRef(sources *fallbackmap.FallbackMap, rules *template.Rules) template.Rule {
	return func(path []interface{}, node interface{}) (interface{}, interface{}){
		key := interface{}(nil)
		if len(path) > 0 { key = path[len(path)-1] }

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
			newNode = template.Process(newNode, rules)
			return key, newNode
		}
		
		return key, node //passthru (ref not found)
	}
}