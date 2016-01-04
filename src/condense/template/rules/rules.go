package rules

import (
	"strings"
	"fallbackmap"
	"deepalias"
	"condense/template"
)

func isEqualString(candidate interface{}, test string) bool {
	var ok bool
	var candidateString string

	if candidateString, ok = candidate.(string); !ok {
		return false
	}

	return candidateString == test
}

func singleKey(candidate interface{}, test string) (interface{}, bool) {
	var ok bool
	var candidateMap map[string]interface{}

	if candidateMap, ok = candidate.(map[string]interface{}); !ok {
		return nil, false
	}
	
	if len(candidateMap) != 1 {
		return nil, false
	}

	var value interface{}
	value, ok = candidateMap[test]
	return value, ok
}

func ExcludeComments(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 { key = path[len(path)-1] }

	var ok bool
	var nodeMap map[string]interface{}

	if nodeMap, ok = node.(map[string]interface{}); !ok {
		return key, node
	}

	if _, ok = nodeMap["$comment"]; !ok {
		return key, node
	}

	if len(nodeMap) == 1 {
		return nil, nil
	}

	delete(nodeMap, "$comment")
	return key, interface{}(nodeMap)
}

func FnJoin(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 { key = path[len(path)-1] }

	argsInterface, ok := singleKey(node, "Fn::Join")
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

	var glue string
	if glue, ok = args[0].(string); !ok {
		return key, node //passthru
	}

	var pieces []interface{}
	if pieces, ok = args[1].([]interface{}); !ok {
		return key, node //passthru
	}

	var piecesStrings []string
	for _, piece := range pieces {
		var pieceString string

		if pieceString, ok = piece.(string); !ok {
			return key, node //passthru
		}

		piecesStrings = append(piecesStrings, pieceString)
	}

	joined := strings.Join(piecesStrings, glue)
	return key, interface{}(joined)
}

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
			newNode = template.Process(newNode, rules)
			return key, newNode
		}
		
		return key, node //passthru (ref not found)
	}
}

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
