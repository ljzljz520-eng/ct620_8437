package search

import (
	"courseattachments/domain"
	"sort"
)

type Facet struct {
	Value string
	Count int
}

func CourseFacets(results []Result) []Facet {
	counts := map[string]int{}
	for _, result := range results {
		counts[result.Attachment.CourseID]++
	}
	return sortedFacets(counts)
}

func StudentFacets(results []Result) []Facet {
	counts := map[string]int{}
	for _, result := range results {
		counts[result.Attachment.StudentID]++
	}
	return sortedFacets(counts)
}

func StatusFacets(results []Result) []Facet {
	counts := map[string]int{}
	for _, result := range results {
		counts[result.Attachment.Status]++
	}
	return sortedFacets(counts)
}

func sortedFacets(counts map[string]int) []Facet {
	result := make([]Facet, 0, len(counts))
	for value, count := range counts {
		result = append(result, Facet{Value: value, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Value < result[j].Value
		}
		return result[i].Count > result[j].Count
	})
	return result
}

func FilterByStatus(results []Result, status string) []Result {
	filtered := []Result{}
	for _, result := range results {
		if result.Attachment.Status == status {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

func FilterByKind(results []Result, kind string) []Result {
	filtered := []Result{}
	for _, result := range results {
		if result.Attachment.Kind == kind {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

func PreviewRows(results []Result) []string {
	rows := make([]string, 0, len(results))
	for _, result := range results {
		label := result.Attachment.DisplayLabel()
		if result.Attachment.Summary == "" {
			label += " | pending"
		} else {
			label += " | " + result.Attachment.Summary
		}
		rows = append(rows, label)
	}
	return rows
}

func CompleteOnly(results []Result) []Result { return FilterByStatus(results, domain.StatusComplete) }
