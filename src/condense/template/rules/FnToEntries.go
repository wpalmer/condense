package rules

func FnToEntries(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 { key = path[len(path)-1] }

	argsInterface, ok := singleKey(node, "Fn::ToEntries")
	if !ok {
		return key, node //passthru
	}

	var argMap map[string]interface{}
	if argMap, ok = argsInterface.(map[string]interface{}); !ok {
		return key, node //passthru
	}

	var entries []interface{}
	for deepKey, deepValue := range argMap {
		entries = append(entries, interface{}(map[string]interface{}{
			"key": deepKey,
			"value": deepValue,
		}))
	}

	return key, interface{}(entries)
}
