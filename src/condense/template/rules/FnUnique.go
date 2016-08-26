package rules

import (
	"reflect"
)

func FnUnique(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 {
		key = path[len(path)-1]
	}

	argsInterface, ok := singleKey(node, "Fn::Unique")
	if !ok {
		return key, node //passthru
	}

	var args []interface{}
	if args, ok = argsInterface.([]interface{}); !ok {
		return key, node //passthru
	}

	var filtered []interface{}
	for i := range args {
		found := false
		for j := range filtered {
			if reflect.DeepEqual(&filtered[j], &args[i]) {
				found = true
				break
			}
		}

		if !found {
			filtered = append(filtered, args[i])
		}
	}

	return key, interface{}(filtered)
}
