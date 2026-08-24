package review

import (
	"courseattachments/domain"
	"courseattachments/search"
	"courseattachments/store"
	"courseattachments/summary"
	"fmt"
	"sort"
	"strings"
)

type Row struct {
	Attachment domain.Attachment
	Task       domain.SummaryTask
	Score      int
	Attention  bool
}

type Dashboard struct {
	Filter      domain.SearchFilter
	Rows        []Row
	Total       int
	Ready       int
	NeedsReview int
	Courses     []search.Facet
	Students    []search.Facet
	Statuses    []search.Facet
}

type Service struct {
	Repo   *store.BoltRepository
	Search *search.Service
}

func New(repo *store.BoltRepository) *Service {
	return &Service{Repo: repo, Search: search.New(repo)}
}

func (s *Service) Build(filter domain.SearchFilter) (Dashboard, error) {
	results, err := s.Search.Query(filter)
	if err != nil {
		return Dashboard{}, err
	}
	dashboard := Dashboard{Filter: filter, Total: len(results)}
	dashboard.Courses = search.CourseFacets(results)
	dashboard.Students = search.StudentFacets(results)
	dashboard.Statuses = search.StatusFacets(results)
	for _, result := range results {
		task, taskErr := s.Repo.GetTask("task-" + result.Attachment.ID)
		if taskErr != nil {
			task = domain.SummaryTask{ID: "task-" + result.Attachment.ID, AttachmentID: result.Attachment.ID, State: domain.TaskQueued}
		}
		row := Row{Attachment: result.Attachment, Task: task, Score: result.Score}
		row.Attention = summary.SummaryNeedsReview(row.Attachment, row.Task)
		if summary.SummaryReady(row.Attachment, row.Task) {
			dashboard.Ready++
		}
		if row.Attention {
			dashboard.NeedsReview++
		}
		dashboard.Rows = append(dashboard.Rows, row)
	}
	sort.SliceStable(dashboard.Rows, func(i, j int) bool {
		if dashboard.Rows[i].Attention != dashboard.Rows[j].Attention {
			return dashboard.Rows[i].Attention
		}
		if dashboard.Rows[i].Score != dashboard.Rows[j].Score {
			return dashboard.Rows[i].Score > dashboard.Rows[j].Score
		}
		return dashboard.Rows[i].Attachment.ID < dashboard.Rows[j].Attachment.ID
	})
	return dashboard, nil
}

func (d Dashboard) AttentionRows() []Row {
	rows := make([]Row, 0)
	for _, row := range d.Rows {
		if row.Attention {
			rows = append(rows, row)
		}
	}
	return rows
}

func (d Dashboard) ReadyRows() []Row {
	rows := make([]Row, 0)
	for _, row := range d.Rows {
		if summary.SummaryReady(row.Attachment, row.Task) {
			rows = append(rows, row)
		}
	}
	return rows
}

func (d Dashboard) Render() string {
	lines := []string{fmt.Sprintf("filter=%s total=%d ready=%d needs-review=%d", d.Filter.Describe(), d.Total, d.Ready, d.NeedsReview)}
	for _, row := range d.Rows {
		marker := "ready"
		if row.Attention {
			marker = "review"
		}
		lines = append(lines, fmt.Sprintf("%s | %s | %s | score=%d", marker, row.Attachment.DisplayLabel(), row.Attachment.StatusLabel(), row.Score))
	}
	return strings.Join(lines, "\n")
}

func (d Dashboard) StatusCounts() map[string]int {
	counts := map[string]int{}
	for _, facet := range d.Statuses {
		counts[facet.Value] = facet.Count
	}
	return counts
}

func (d Dashboard) HasAttention() bool { return d.NeedsReview > 0 }
