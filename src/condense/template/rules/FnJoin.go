package rules

import (
	"fmt"
	"strings"
)

func FnJoin(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 {
		key = path[len(path)-1]
	}

	argsInterface, ok := singleKey(node, "Fn::Join")
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

	var pieces []interface{}
	if pieces, ok = args[1].([]interface{}); !ok {
		return key, node //passthru
	}

	var piecesStrings []string
	for _, piece := range pieces {
		var pieceString string

		switch typed := piece.(type) {
		default:
			return key, node //passthru
		case string:
			pieceString = typed
		case int:
			pieceString = fmt.Sprintf("%d", typed)
		case float64:
			pieceString = fmt.Sprintf("%d", int(typed))
		}

		piecesStrings = append(piecesStrings, pieceString)
	}

	joined := strings.Join(piecesStrings, glue)
	return key, interface{}(joined)
}
