package rules

import (
	"fallbackmap"
	"condense/template"
)

func MakeFnFor(sources *fallbackmap.FallbackMap, templateRules *template.Rules) template.Rule {
	return func (path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 { key = path[len(path)-1] }

		argsInterface, ok := singleKey(node, "Fn::For")
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

		var refName string
		if refName, ok = args[0].(string); !ok {
			return key, node //passthru
		}

		var values []interface{}
		if values, ok = args[1].([]interface{}); !ok {
			return key, node //passthru
		}

		loopTemplate := interface{}(args[2])

		var generated []interface{}
		for deepIndex, value := range values {
			loopTemplateSources := fallbackmap.FallbackMap{}
			loopTemplateRules := template.Rules{}
			
			loopTemplateSources.Attach(fallbackmap.DeepFunc(func(path []string) (interface{}, bool) {
				if len(path) != 1 || path[0] != refName {
					return nil, false
				}

				return value, true
			}))
			loopTemplateSources.Attach(sources)

			loopTemplateRules.Attach(MakeRef(&loopTemplateSources, &loopTemplateRules))
			loopTemplateRules.Attach(templateRules.MakeEach())
			loopTemplateRules.AttachEarly(templateRules.MakeEachEarly())

			deepPath := make([]interface{}, len(path)+1)
			copy(deepPath, path)
			deepPath[cap(deepPath)-1] = interface{}(deepIndex)

			newIndex, processed := template.Walk(deepPath, loopTemplate, &loopTemplateRules)
			if newIndex != nil {
				generated = append(generated, processed)
			}
		}

		return key, interface{}(generated)
	}
}
