package domain

import "strings"

type SearchFilter struct{ CourseID, StudentID, Keyword string }

func NewSearchFilter(course, student, keyword string) SearchFilter {
	return SearchFilter{CourseID: strings.TrimSpace(course), StudentID: strings.TrimSpace(student), Keyword: NormalizeKeyword(keyword)}
}

func (f SearchFilter) Matches(a Attachment) bool {
	if f.CourseID != "" && a.CourseID != f.CourseID {
		return false
	}
	if f.StudentID != "" && a.StudentID != f.StudentID {
		return false
	}
	if f.Keyword != "" {
		hay := strings.ToLower(a.Filename + " " + a.Keywords + " " + a.Summary)
		return strings.Contains(hay, strings.ToLower(f.Keyword))
	}
	return true
}

func (f SearchFilter) Empty() bool {
	return f.CourseID == "" && f.StudentID == "" && NormalizeKeyword(f.Keyword) == ""
}

func (f SearchFilter) Describe() string {
	parts := []string{}
	if f.CourseID != "" {
		parts = append(parts, "course="+f.CourseID)
	}
	if f.StudentID != "" {
		parts = append(parts, "student="+f.StudentID)
	}
	if f.Keyword != "" {
		parts = append(parts, "keyword="+NormalizeKeyword(f.Keyword))
	}
	if len(parts) == 0 {
		return "all attachments"
	}
	return strings.Join(parts, ", ")
}
func NormalizeKeyword(v string) string { return strings.ToLower(strings.TrimSpace(v)) }
func TokenizeKeywords(v string) []string {
	raw := strings.Fields(strings.ToLower(v))
	out := make([]string, 0, len(raw))
	for _, x := range raw {
		if len(x) > 1 {
			out = append(out, x)
		}
	}
	return out
}

func UniqueKeywords(values ...string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		for _, token := range TokenizeKeywords(value) {
			if !seen[token] {
				seen[token] = true
				out = append(out, token)
			}
		}
	}
	return out
}

func KeywordHistogram(attachments []Attachment) map[string]int {
	histogram := map[string]int{}
	for _, attachment := range attachments {
		for _, keyword := range UniqueKeywords(attachment.Keywords) {
			histogram[keyword]++
		}
	}
	return histogram
}
func (a Attachment) KeywordScore(k string) int {
	score := 0
	for _, x := range TokenizeKeywords(a.Keywords + " " + a.Filename + " " + a.Summary) {
		if x == NormalizeKeyword(k) {
			score += 2
		} else if strings.Contains(x, NormalizeKeyword(k)) {
			score++
		}
	}
	return score
}

func (a Attachment) KeywordMatches(k string) []string {
	needle := NormalizeKeyword(k)
	if needle == "" {
		return nil
	}
	seen := map[string]bool{}
	matched := []string{}
	for _, token := range TokenizeKeywords(a.Keywords + " " + a.Filename + " " + a.Summary) {
		if strings.Contains(token, needle) && !seen[token] {
			seen[token] = true
			matched = append(matched, token)
		}
	}
	return matched
}
