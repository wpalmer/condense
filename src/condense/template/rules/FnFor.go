package rules

import (
	"condense/template"
	"deepstack"
	"fallbackmap"
)

func MakeFnFor(sources *deepstack.DeepStack, templateRules *template.Rules) template.Rule {
	return func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 {
			key = path[len(path)-1]
		}

		raw, ok := singleKey(node, "Fn::For")
		if !ok {
			return key, node //passthru
		}

		args, ok := collectArgs(
			raw,
			func(argsSoFar []interface{}) bool { return len(argsSoFar) < 3 },
			func(argsSoFar []interface{}, arg interface{}) (skip bool, newNode interface{}) {
				// unconditionally process the argument, in case it needs to be skipped
				key, node := template.Walk(path, arg, templateRules)
				if skip, ok := key.(bool); ok && skip {
					return true, nil
				}

				if len(argsSoFar) == 2 {
					return false, arg // return unprocessed 3rd arg. It's a template.
				}

				return false, node
			},
		)
		if !ok {
			return key, node //passthru
		}

		var refNames []interface{}
		var refName interface{}

		if refNames, ok = args[0].([]interface{}); ok {
			if len(refNames) == 1 {
				refNames = []interface{}{nil, refNames[0]}
			} else if len(refNames) != 2 {
				return key, node //passthru
			}
		} else {
			refNames = []interface{}{nil, args[0]}
		}

		for _, refName = range refNames {
			if _, ok = refName.(string); !ok && refName != nil {
				return key, node //passthru
			}
		}

		valuesInterface := args[1]

		var values []interface{}
		if values, ok = valuesInterface.([]interface{}); !ok {
			return key, node //passthru
		}

		loopTemplate := interface{}(args[2])

		generated := []interface{}{}
		for deepIndex, value := range values {
			refMap := make(map[string]interface{})
			if refNames[0] != nil {
				refMap[refNames[0].(string)] = float64(deepIndex)
			}

			if refNames[1] != nil {
				refMap[refNames[1].(string)] = value
			}

			deepPath := make([]interface{}, len(path)+1)
			copy(deepPath, path)
			deepPath[cap(deepPath)-1] = interface{}(deepIndex)

			sources.Push(fallbackmap.DeepMap(refMap))

			newIndex, processed := template.Walk(deepPath, loopTemplate, templateRules)
			sources.PopDiscard()

			if newIndex != nil {
				generated = append(generated, processed)
			}
		}

		return key, interface{}(generated)
	}
}
