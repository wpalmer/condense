package rules

func isEqualString(candidate interface{}, test string) bool {
	var ok bool
	var candidateString string

	if candidateString, ok = candidate.(string); !ok {
		return false
	}

	return candidateString == test
}

func singleKey(candidate interface{}, test string) (interface{}, bool) {
	var ok bool
	var candidateMap map[string]interface{}

	if candidateMap, ok = candidate.(map[string]interface{}); !ok {
		return nil, false
	}
	
	if len(candidateMap) != 1 {
		return nil, false
	}

	var value interface{}
	value, ok = candidateMap[test]
	return value, ok
}
