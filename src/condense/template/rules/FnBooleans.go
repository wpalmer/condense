package rules

import (
	"reflect"
)

func FnIf(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 { key = path[len(path)-1] }

	argsInterface, ok := singleKey(node, "Fn::If")
	if !ok {
		return key, node //passthru
	}

	var args []interface{}
	if args, ok = argsInterface.([]interface{}); !ok {
		return key, node //passthru
	}

	if len(args) != 3 {
		return key, node //passthru
	}

	var condition bool
	if condition, ok = args[0].(bool); !ok {
		return key, node //passthru
	}

	if condition {
		return key, args[1]
	}

	return key, args[2]
}

func FnEquals(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 { key = path[len(path)-1] }

	argsInterface, ok := singleKey(node, "Fn::Equals")
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

	return key, interface{}(reflect.DeepEqual(args[0], args[1]))
}

func FnAnd(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 { key = path[len(path)-1] }

	argsInterface, ok := singleKey(node, "Fn::And")
	if !ok {
		return key, node //passthru
	}

	var args []interface{}
	if args, ok = argsInterface.([]interface{}); !ok {
		return key, node //passthru
	}

	if len(args) < 2 {
		return key, node //passthru
	}

	for _, arg := range args {
		var argBool bool
		if argBool, ok = arg.(bool); !ok {
			return key, node //passthru
		}

		if !argBool {
			return key, interface{}(false)
		}
	}

	return key, interface{}(true)
}

func FnOr(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 { key = path[len(path)-1] }

	argsInterface, ok := singleKey(node, "Fn::Or")
	if !ok {
		return key, node //passthru
	}

	var args []interface{}
	if args, ok = argsInterface.([]interface{}); !ok {
		return key, node //passthru
	}

	if len(args) < 2 {
		return key, node //passthru
	}

	hasTrue := false
	for _, arg := range args {
		var argBool bool
		if argBool, ok = arg.(bool); !ok {
			return key, node //passthru
		}

		hasTrue = hasTrue || argBool
	}

	return key, interface{}(hasTrue)
}

func FnNot(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 { key = path[len(path)-1] }

	argsInterface, ok := singleKey(node, "Fn::Not")
	if !ok {
		return key, node //passthru
	}

	var args []interface{}
	if args, ok = argsInterface.([]interface{}); !ok {
		return key, node //passthru
	}

	if len(args) != 1 {
		return key, node //passthru
	}

	arg := args[0]
	var argBool bool
	if argBool, ok = arg.(bool); !ok {
		return key, node //passthru
	}

	return key, interface{}(!argBool)
}
