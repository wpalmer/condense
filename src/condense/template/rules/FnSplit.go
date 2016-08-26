package rules

import (
	"strings"
)

func FnSplit(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 {
		key = path[len(path)-1]
	}

	argsInterface, ok := singleKey(node, "Fn::Split")
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

	var glue string
	if glue, ok = args[0].(string); !ok {
		return key, node //passthru
	}

	var joined string
	if joined, ok = args[1].(string); !ok {
		return key, node //passthru
	}

	var pieces []interface{}
	for _, piece := range strings.Split(joined, glue) {
		pieces = append(pieces, interface{}(piece))
	}

	return key, interface{}(pieces)
}
