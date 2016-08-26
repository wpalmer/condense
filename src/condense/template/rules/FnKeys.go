package rules

func FnKeys(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 {
		key = path[len(path)-1]
	}

	argInterface, ok := singleKey(node, "Fn::Keys")
	if !ok {
		return key, node //passthru
	}

	var keys []interface{}
	var argMap map[string]interface{}
	if argMap, ok = argInterface.(map[string]interface{}); !ok {
		return key, node //passthru
	}

	for key := range argMap {
		keys = append(keys, key)
	}

	return key, interface{}(keys)
}
