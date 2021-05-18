package internal

import (
	"sort"

	"github.com/sha1n/benchy/api"
)

// GetSortedScenarioIds returns a sorted array of scenario IDs for the specified api.Summary
func GetSortedScenarioIds(summary api.Summary) []api.ID {
	statsMap := summary.All()
	sortedIds := make([]api.ID, 0, len(statsMap))
	for k := range statsMap {
		sortedIds = append(sortedIds, k)
	}
	sort.Strings(sortedIds)

	return sortedIds
}
