package review

import (
	"courseattachments/domain"
	"courseattachments/store"
	"fmt"
	"sort"
)

type TeacherQueue struct {
	TeacherID string
	Rows      []Row
}

func (s *Service) QueueForTeacher(teacherID string, filter domain.SearchFilter) (TeacherQueue, error) {
	if teacherID == "" {
		return TeacherQueue{}, fmt.Errorf("teacher id required")
	}
	dashboard, err := s.Build(filter)
	if err != nil {
		return TeacherQueue{}, err
	}
	rows := dashboard.AttentionRows()
	sort.Slice(rows, func(i, j int) bool { return rows[i].Attachment.CreatedAt < rows[j].Attachment.CreatedAt })
	return TeacherQueue{TeacherID: teacherID, Rows: rows}, nil
}

func (q TeacherQueue) Next() (Row, bool) {
	if len(q.Rows) == 0 {
		return Row{}, false
	}
	return q.Rows[0], true
}

func (q TeacherQueue) IDs() []string {
	ids := make([]string, 0, len(q.Rows))
	for _, row := range q.Rows {
		ids = append(ids, row.Attachment.ID)
	}
	return ids
}

func (q TeacherQueue) Empty() bool { return len(q.Rows) == 0 }

func TeacherSummary(repo *store.BoltRepository, teacherID string, filter domain.SearchFilter) (string, error) {
	queue, err := New(repo).QueueForTeacher(teacherID, filter)
	if err != nil {
		return "", err
	}
	if queue.Empty() {
		return "no review items", nil
	}
	return fmt.Sprintf("teacher=%s pending=%d next=%s", queue.TeacherID, len(queue.Rows), queue.Rows[0].Attachment.DisplayLabel()), nil
}
