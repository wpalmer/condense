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

type expectMoreCallback func(argsSoFar []interface{}) bool
type processCallback func(argsSoFar []interface{}, arg interface{}) (skip bool, newNode interface{})

func collectArgs(candidate interface{}, expectMore expectMoreCallback, process processCallback) ([]interface{}, bool) {
	var args []interface{}
	var ok bool

	if args, ok = candidate.([]interface{}); !ok {
		return nil, false // invalid args list
	}

	var collected []interface{}
	for _, arg := range args {
		skip, node := process(collected, arg)
		if skip {
			continue
		}

		if !expectMore(collected) {
			return nil, false // invalid args list: too many arguments
		}

		collected = append(collected, node)
	}

	if expectMore(collected) {
		return nil, false // invalid args list: not enough arguments
	}

	return collected, true
}
