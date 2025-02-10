// Copyright (c) HashiCorp, Inc.

package utils

import "sort"

// Generic function to reorder elements based on a reference slice
func Reorder[T any, K comparable](
	resultItems []interface{},
	dataItems []T,
	extractKey func(T) K,
	extractResultKey func(map[string]interface{}) K,
) []interface{} {
	// Map to store the order of items based on data
	itemOrder := make(map[K]int)
	for i, item := range dataItems {
		itemOrder[extractKey(item)] = i
	}

	// Slices to store reordered and unmatched items
	ordered := make([]interface{}, 0, len(resultItems))
	unmatched := make([]interface{}, 0, len(resultItems))

	// Reorder resultItems based on dataItems
	for _, item := range resultItems {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			unmatched = append(unmatched, item) // Keep unmatched in original order
			continue
		}

		itemKey := extractResultKey(itemMap)
		if _, exists := itemOrder[itemKey]; exists {
			ordered = append(ordered, item)
		} else {
			unmatched = append(unmatched, item)
		}
	}

	// Sort ordered slice based on dataItems order
	sort.SliceStable(ordered, func(i, j int) bool {
		idI := extractResultKey(ordered[i].(map[string]interface{}))
		idJ := extractResultKey(ordered[j].(map[string]interface{}))
		return itemOrder[idI] < itemOrder[idJ]
	})

	// Append unmatched items at the end
	ordered = append(ordered, unmatched...)
	return ordered
}
