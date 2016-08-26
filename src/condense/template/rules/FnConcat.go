package rules

func FnConcat(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 {
		key = path[len(path)-1]
	}

	argsInterface, ok := singleKey(node, "Fn::Concat")
	if !ok {
		return key, node //passthru
	}

	var args []interface{}
	if args, ok = argsInterface.([]interface{}); !ok {
		return key, node //passthru
	}

	var concatenated []interface{}
	for _, argInterface := range args {
		var argArray []interface{}
		if argArray, ok = argInterface.([]interface{}); !ok {
			return key, node //can't concatenate non-arrays, passthru
		}

		for _, item := range argArray {
			concatenated = append(concatenated, item)
		}
	}

	return key, interface{}(concatenated)
}
