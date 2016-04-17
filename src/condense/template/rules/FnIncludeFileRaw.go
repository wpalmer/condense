package rules

import (
	"condense/template"
	"fmt"
	"golang.org/x/tools/godoc/vfs"
	"io"
	"io/ioutil"
	"path/filepath"
)

func MakeFnIncludeFileRaw(opener vfs.Opener) template.Rule {
	return func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 {
			key = path[len(path)-1]
		}

		argInterface, ok := singleKey(node, "Fn::IncludeFileRaw")
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

		var dataStream io.Reader
		if dataStream, err = opener.Open(absPath); err != nil {
			panic(fmt.Errorf("Error opening imported file '%s': %s", absPath, err))
		}

		var data []byte
		if data, err = ioutil.ReadAll(dataStream); err != nil {
			panic(fmt.Errorf("Error loading imported file '%s': %s", argString, err))
		}

		return key, interface{}(string(data))
	}
}
