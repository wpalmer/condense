package rules

func FnFromEntries(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 { key = path[len(path)-1] }

	argsInterface, ok := singleKey(node, "Fn::FromEntries")
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
			return key, node //can't derive from non-map entry, passthru
		}

		var deepKeyInterface interface{}
		var deepKeyString string
		var deepValue interface{}

		if deepKeyInterface, ok = argMap["key"]; !ok {
			return key, node //can't derive from non {key,value} entry, passthru
		}

		if deepKeyString, ok = deepKeyInterface.(string); !ok {
			return key, node //can't derive from non-string key entry, passthru
		}

		if deepValue, ok = argMap["value"]; !ok {
			return key, node //can't derive from non {key,value} entry, passthru
		}

		merged[deepKeyString] = deepValue
	}

	return key, interface{}(merged)
}
