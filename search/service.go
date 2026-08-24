package search

import (
	"courseattachments/domain"
	"courseattachments/store"
	"sort"
)

type Service struct{ Repo *store.BoltRepository }
type Result struct {
	Attachment domain.Attachment
	Score      int
	Matches    []string
}

func New(r *store.BoltRepository) *Service { return &Service{Repo: r} }
func (s *Service) Query(f domain.SearchFilter) ([]Result, error) {
	as, e := s.Repo.ListAttachments()
	if e != nil {
		return nil, e
	}
	out := []Result{}
	for _, a := range as {
		if f.Matches(a) {
			out = append(out, Result{Attachment: a, Score: a.KeywordScore(f.Keyword), Matches: a.KeywordMatches(f.Keyword)})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].Attachment.ID < out[j].Attachment.ID
		}
		return out[i].Score > out[j].Score
	})
	return out, nil
}
func (s *Service) Summaries(f domain.SearchFilter) ([]string, error) {
	rs, e := s.Query(f)
	if e != nil {
		return nil, e
	}
	o := []string{}
	for _, r := range rs {
		o = append(o, r.Attachment.Filename+": "+r.Attachment.Summary)
	}
	return o, nil
}
func ParseFilter(course, student, keyword string) domain.SearchFilter {
	return domain.NewSearchFilter(course, student, keyword)
}

func (s *Service) QueryPage(filter domain.SearchFilter, offset, limit int) ([]Result, int, error) {
	results, err := s.Query(filter)
	if err != nil {
		return nil, 0, err
	}
	total := len(results)
	if offset < 0 {
		offset = 0
	}
	if offset >= total {
		return []Result{}, total, nil
	}
	end := total
	if limit > 0 && offset+limit < end {
		end = offset + limit
	}
	return results[offset:end], total, nil
}

func (s *Service) Suggestions(prefix string, limit int) ([]string, error) {
	attachments, err := s.Repo.ListAttachments()
	if err != nil {
		return nil, err
	}
	prefix = domain.NormalizeKeyword(prefix)
	seen := map[string]bool{}
	result := []string{}
	for _, attachment := range attachments {
		for _, keyword := range domain.UniqueKeywords(attachment.Keywords) {
			if (prefix == "" || len(keyword) >= len(prefix) && keyword[:len(prefix)] == prefix) && !seen[keyword] {
				seen[keyword] = true
				result = append(result, keyword)
			}
		}
	}
	sort.Strings(result)
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

func (s *Service) CourseSummary(courseID string) (map[string]int, error) {
	results, err := s.Query(domain.SearchFilter{CourseID: courseID})
	if err != nil {
		return nil, err
	}
	return GroupByStatus(results), nil
}
