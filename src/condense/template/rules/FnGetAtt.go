package rules

import (
	"fallbackmap"
	"deepalias"
	"condense/template"
)

func MakeFnGetAtt(sources *fallbackmap.FallbackMap, rules *template.Rules) template.Rule {
	return func(path []interface{}, node interface{}) (interface{}, interface{}){
		key := interface{}(nil)
		if len(path) > 0 { key = path[len(path)-1] }

		argsInterface, ok := singleKey(node, "Fn::GetAtt")
		if !ok {
			return key, node //passthru
		}

		var args []interface{}
		if args, ok = argsInterface.([]interface{}); !ok {
			return key, node //passthru
		}

		if len(args) != 2 {
			return key, node //passthru
		}

		var refpath []string
		for _, arg := range args {
			var argString string
			if argString, ok = arg.(string); !ok {
				return key, node //passthru
			}

			for _, part := range deepalias.Split(argString) {
				refpath = append(refpath, part)
			}
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
