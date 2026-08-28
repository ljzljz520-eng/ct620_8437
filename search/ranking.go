package search

import (
	"courseattachments/domain"
	"sort"
	"strings"
)

func Rank(results []Result) []Result {
	sort.SliceStable(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Attachment.ID < results[j].Attachment.ID
		}
		return results[i].Score > results[j].Score
	})
	return results
}
func GroupByStatus(results []Result) map[string]int {
	o := map[string]int{}
	for _, r := range results {
		o[r.Attachment.Status]++
	}
	return o
}
func FilterComplete(results []Result) []Result {
	o := []Result{}
	for _, r := range results {
		if r.Attachment.Status == domain.StatusComplete {
			o = append(o, r)
		}
	}
	return o
}
func Names(results []Result) []string {
	o := []string{}
	for _, r := range results {
		o = append(o, r.Attachment.Filename)
	}
	return o
}

func Top(results []Result, limit int) []Result {
	Rank(results)
	if limit <= 0 || limit >= len(results) {
		return results
	}
	return results[:limit]
}

func ScoreTotal(results []Result) int {
	total := 0
	for _, result := range results {
		total += result.Score
	}
	return total
}

func WithMinimumScore(results []Result, minimum int) []Result {
	filtered := []Result{}
	for _, result := range results {
		if result.Score >= minimum {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

func Explain(result Result) string {
	if len(result.Matches) == 0 {
		return result.Attachment.Filename + " (no keyword boost)"
	}
	return result.Attachment.Filename + " matched " + strings.Join(result.Matches, ",")
}
