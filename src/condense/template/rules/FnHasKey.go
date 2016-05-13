package rules

func FnHasKey(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 { key = path[len(path)-1] }

	argsInterface, ok := singleKey(node, "Fn::HasKey")
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

	var argKey string
	if argKey, ok = args[0].(string); !ok {
		return key, node //passthru
	}

	var argMap map[string]interface{}
	if argMap, ok = args[1].(map[string]interface{}); !ok {
		return key, node //passthru
	}

	_, ok = argMap[argKey]
	return key, interface{}(ok)
}
