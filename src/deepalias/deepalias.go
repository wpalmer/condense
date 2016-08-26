package deepalias

import (
	"fallbackmap"
	"strings"
)

type DeepAlias struct {
	fallbackmap.Deep
}

func Split(pathString string) []string {
	var final []string
	var next []string
	var components []string
	var nested int

	components = strings.Split(pathString, ".")
	for _, component := range components {
		next = append(next, component)
		if strings.HasPrefix(component, "[") {
			nested++
		}

		if strings.HasSuffix(component, "]") {
			nested--
		}

		if nested == 0 {
			final = append(final, strings.Join(next, "."))
			next = []string{}
		}
	}

	if len(next) > 0 {
		final = append(final, next...)
	}

	return final
}

func DeAlias(path []string, deep fallbackmap.Deep) ([]string, bool) {
	did_translate := false
	translated := []string{}
	for _, component := range path {
		if strings.HasPrefix(component, "[") && strings.HasSuffix(component, "]") {
			alias_path := strings.Split(component[1:len(component)-1], ".")
			dereferenced_component, found := deep.Get(alias_path)
			if found {
				dereferenced_component_string, ok := dereferenced_component.(string)
				if ok {
					component = dereferenced_component_string
					did_translate = true
				}
			}
		}

		translated = append(translated, component)
	}

	return translated, did_translate
}

func (alias DeepAlias) Get(path []string) (value interface{}, ok bool) {
	if len(path) == 0 {
		return alias, true
	}

	translated, did_translate := DeAlias(path, alias.Deep)
	if !did_translate {
		return nil, false
	}

	return alias.Deep.Get(translated)
}
