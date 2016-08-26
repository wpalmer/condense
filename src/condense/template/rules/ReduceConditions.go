package rules

func ReduceConditions(path []interface{}, node interface{}) (interface{}, interface{}) {
	var ok bool
	key := interface{}(nil)
	if len(path) > 0 {
		key = path[len(path)-1]
	}

	if len(path) < 2 || !isEqualString(path[0], "Conditions") {
		return key, node //passthru
	}

	if _, ok = path[1].(string); !ok {
		// not a named condition
		return key, node //passthru
	}

	var conditionState bool
	if conditionState, ok = node.(bool); !ok {
		// not a bool, already okay
		return key, node //passthru
	}

	if conditionState {
		return key, map[string]interface{}{
			"Fn::Equals": interface{}([]interface{}{
				interface{}("1"),
				interface{}("1"),
			}),
		}
	}

	return key, map[string]interface{}{
		"Fn::Equals": interface{}([]interface{}{
			interface{}("0"),
			interface{}("1"),
		}),
	}
}
