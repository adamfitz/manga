package compare

import (
	"encoding/json"
	"fmt"
)

func CompareJSONArrays(jsonA, jsonB string) ([]string, error) {
	/* 
		Performs a comparison of two JSON arrays and returns the elements that are present in the first array but not in the second array.
	*/
	var listA, listB []string

	// Unmarshal the first JSON string into a slice
	err := json.Unmarshal([]byte(jsonA), &listA)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON A: %w", err)
	}

	// Unmarshal the second JSON string into a slice
	err = json.Unmarshal([]byte(jsonB), &listB)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON B: %w", err)
	}

	// Find differences: elements in A but not in B
	diff := []string{}
	// using a set for effeciancy
	setB := make(map[string]struct{}, len(listB))
	for _, item := range listB {
		setB[item] = struct{}{}
	}
	for _, item := range listA {
		if _, found := setB[item]; !found {
			diff = append(diff, item)
		}
	}

	return diff, nil
}