package rules

func FnMerge(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 { key = path[len(path)-1] }

	argsInterface, ok := singleKey(node, "Fn::Merge")
	if !ok {
		return key, node //passthru
	}

	var args []interface{}
	if args, ok = argsInterface.([]interface{}); !ok {
		return key, node //passthru
	}

	merged := make(map[string]interface{})
	for _, argInterface := range args {
		var argMap map[string]interface{}
		if argMap, ok = argInterface.(map[string]interface{}); !ok {
			return key, node //can't merge non-maps, passthru
		}

		for deepKey, value := range argMap {
			merged[deepKey] = value
		}
	}

	return key, interface{}(merged)
}
