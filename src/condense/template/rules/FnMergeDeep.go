package rules

func deepMerge(depth int, maps ...map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for _, oneMap := range maps {
		for key, _ := range oneMap {
			if depth > 0 {
				if _, ok := merged[key]; ok {
					if _, ok := merged[key].(map[string]interface{}); ok {
						if _, ok := oneMap[key].(map[string]interface{}); ok {
							merged[key] = deepMerge(depth - 1,
								merged[key].(map[string]interface{}),
								oneMap[key].(map[string]interface{}),
							)
							continue
						}
					}
				}
			}

			merged[key] = oneMap[key]
		}
	}

	return merged
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

	var mapsInterface []interface{}
	if mapsInterface, ok = args[1].([]interface{}); !ok {
		return key, node //passthru
	}

	var merged map[string]interface{}
	maps := []map[string]interface{}{ }

	for _, argInterface := range mapsInterface {
		var argMap map[string]interface{}
		if argMap, ok = argInterface.(map[string]interface{}); !ok {
			return key, node //can't merge non-maps, passthru
		}

		maps = append(maps, argMap)
	}

	merged = deepMerge(depth, maps...)
	return key, interface{}(merged)
}
