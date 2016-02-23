package rules

import (
	"condense/template"

	"fmt"
	"os"
	fspath "path"
)

type Stater interface{
	Stat(path string) (os.FileInfo, error)
}

func MakeFnFindFile(stater Stater) template.Rule {
	return func(path []interface{}, node interface{}) (interface{}, interface{}){
		key := interface{}(nil)
		if len(path) > 0 { key = path[len(path)-1] }

		argsInterface, ok := singleKey(node, "Fn::FindFile")
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

		var prefixes []interface{}
		if prefixes, ok = args[0].([]interface{}); !ok {
			return key, node //passthru
		}

		prefixStrings := []string{}
		for _, prefixInterface := range prefixes {
			if _, ok := prefixInterface.(string); !ok {
				return key, node //passthru
			}

			prefixStrings = append(prefixStrings, prefixInterface.(string))
		}

		var tail string
		if tail, ok = args[1].(string); !ok {
			return key, node //passthru
		}

		for _, prefix := range prefixStrings {
			fullpath := fspath.Join(prefix, tail)
			if _, err := stater.Stat(fullpath); err == nil {
				return key, interface{}(fullpath)
			}
		}

		panic(fmt.Errorf("Unable to locate file '%s' in %v", tail, prefixStrings))
	}
}
