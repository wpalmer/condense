package rules

import (
	"fallbackmap"
	"condense/template"
	"deepalias"
)

func MakeFnWith(sources *fallbackmap.FallbackMap, outerRules *template.Rules) template.Rule {
	return func (path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 { key = path[len(path)-1] }

		argsInterface, ok := singleKey(node, "Fn::With")
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

		var source map[string]interface{}
		if source, ok = args[0].(map[string]interface{}); !ok {
			return key, node //passthru
		}

		innerTemplate := interface{}(args[1])
		templateSources := fallbackmap.FallbackMap{}
		templateRules := template.Rules{}

		templateSources.Attach(fallbackmap.DeepMap(source))
		templateSources.Attach(deepalias.DeepAlias{&templateSources})
		templateSources.Attach(sources)

		templateRules.Attach(MakeRef(&templateSources, &templateRules))
		templateRules.Attach(MakeFnGetAtt(&templateSources, &templateRules))
		templateRules.Attach(outerRules.MakeEach())
		templateRules.AttachEarly(outerRules.MakeEachEarly())

		key, generated := template.Walk(path, innerTemplate, &templateRules)
		return key, interface{}(generated)
	}
}
