package parser

import (
	"fmt"
)

func MapKeys(m map[string]any) []string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// search the value of a nested map by a list of keys.  Send each key is as an individual string
func NestedMapValue(m map[string]any, keys ...string) (any, error) {
	var current any = m

	for _, key := range keys {
		// If current is a map, assert it to map[string]interface{}
		if mapValue, ok := current.(map[string]any); ok {
			current = mapValue[key]
		} else {
			return nil, fmt.Errorf("key '%s' not found or not a map at level", key)
		}
	}

	return current, nil
}
