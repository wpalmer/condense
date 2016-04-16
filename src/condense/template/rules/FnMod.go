package rules

func FnMod(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 { key = path[len(path)-1] }

	argsInterface, ok := singleKey(node, "Fn::Mod")
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

	for _, arg := range args {
		if _, ok = arg.(float64); !ok {
			return key, node //passthru
		}
	}

	return key, interface{}(float64(int(args[0].(float64)) % int(args[1].(float64))))
}
