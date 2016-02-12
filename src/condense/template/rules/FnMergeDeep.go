package rules

func deepMerge(target map[string]interface{}, source map[string]interface{}, depth int) {
	for key, value := range source {
		if depth > 0 {
			if _, ok := target[key]; ok {
				if _, ok := target[key].(map[string]interface{}); ok {
					if _, ok := value.(map[string]interface{}); ok {
						deepMerge(target[key].(map[string]interface{}), value.(map[string]interface{}), depth - 1)
						continue
					}
				}
			}
		}

		target[key] = value
	}
}

func FnMergeDeep(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 { key = path[len(path)-1] }

	argsInterface, ok := singleKey(node, "Fn::MergeDeep")
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

	var depthFloat float64
	if depthFloat, ok = args[0].(float64); !ok {
		return key, node //passthru
	}
	depth := int(depthFloat)

	var maps []interface{}
	if maps, ok = args[1].([]interface{}); !ok {
		return key, node //passthru
	}

	merged := make(map[string]interface{})
	for _, argInterface := range maps {
		var argMap map[string]interface{}
		if argMap, ok = argInterface.(map[string]interface{}); !ok {
			return key, node //can't merge non-maps, passthru
		}

		deepMerge(merged, argMap, depth)
	}

	return key, interface{}(merged)
}
