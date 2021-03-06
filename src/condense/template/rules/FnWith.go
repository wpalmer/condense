package rules

import (
	"condense/template"
	"deepstack"
	"fallbackmap"
)

func MakeFnWith(sources *deepstack.DeepStack, outerRules *template.Rules) template.Rule {
	return func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 {
			key = path[len(path)-1]
		}

		raw, ok := singleKey(node, "Fn::With")
		if !ok {
			return key, node //passthru
		}

		args, ok := collectArgs(
			raw,
			func(argsSoFar []interface{}) bool {
				return len(argsSoFar) < 2
			},
			func(argsSoFar []interface{}, arg interface{}) (bool, interface{}) {
				// unconditionally process the argument, in case it needs to be skipped
				key, node := template.Walk(path, arg, outerRules)
				if skip, ok := key.(bool); ok && skip {
					return true, nil
				}

				if len(argsSoFar) == 1 {
					return false, arg // return unprocessed 2nd arg. It's a template.
				}

				return false, node
			},
		)

		if !ok {
			return key, node //passthru
		}

		var source map[string]interface{}
		if source, ok = args[0].(map[string]interface{}); !ok {
			return key, node //passthru
		}

		sources.Push(fallbackmap.DeepMap(source))
		innerTemplate := interface{}(args[1])
		key, generated := template.Walk(path, innerTemplate, outerRules)
		sources.PopDiscard()

		return key, interface{}(generated)
	}
}
