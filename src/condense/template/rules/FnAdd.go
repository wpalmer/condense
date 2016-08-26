package rules

func FnAdd(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 {
		key = path[len(path)-1]
	}

	argsInterface, ok := singleKey(node, "Fn::Add")
	if !ok {
		return key, node //passthru
	}

	var args []interface{}
	if args, ok = argsInterface.([]interface{}); !ok {
		return key, node //passthru
	}

	total := float64(0)
	for _, arg := range args {
		var argFloat float64

		if argFloat, ok = arg.(float64); !ok {
			return key, node //passthru
		}

		total += argFloat
	}

	return key, interface{}(total)
}
