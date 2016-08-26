package rules

import (
	"condense/template"
	"encoding/json"
	"fmt"
	"golang.org/x/tools/godoc/vfs"
	"io"
	"path/filepath"
)

func MakeFnIncludeFile(opener vfs.Opener, rules *template.Rules) template.Rule {
	return func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 {
			key = path[len(path)-1]
		}

		argInterface, ok := singleKey(node, "Fn::IncludeFile")
		if !ok {
			return key, node //passthru
		}

		var argString string
		if argString, ok = argInterface.(string); !ok {
			return key, node //passthru
		}

		var absPath string
		var err error
		if absPath, err = filepath.Abs(argString); err != nil {
			panic(fmt.Errorf("Error opening imported file '%s': %s", argString, err))
		}

		var jsonStream io.Reader
		if jsonStream, err = opener.Open(absPath); err != nil {
			panic(fmt.Errorf("Error opening imported file '%s': %s", absPath, err))
		}

		dec := json.NewDecoder(jsonStream)
		includedTemplate := interface{}(nil)
		if err := dec.Decode(&includedTemplate); err != nil {
			panic(fmt.Errorf("Error loading imported file '%s': %s", argString, err))
		}

		key, generated := template.Walk(path, includedTemplate, rules)
		return key, interface{}(generated)
	}
}
