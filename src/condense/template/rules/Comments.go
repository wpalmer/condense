package rules

func ExcludeComments(path []interface{}, node interface{}) (interface{}, interface{}) {
	key := interface{}(nil)
	if len(path) > 0 { key = path[len(path)-1] }

	var ok bool
	var nodeMap map[string]interface{}

	if nodeMap, ok = node.(map[string]interface{}); !ok {
		return key, node
	}

	if _, ok = nodeMap["$comment"]; !ok {
		return key, node
	}

	if len(nodeMap) == 1 {
		return true, nil
	}

	delete(nodeMap, "$comment")
	return key, interface{}(nodeMap)
}
