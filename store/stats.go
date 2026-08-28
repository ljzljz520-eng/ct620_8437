package store

import (
	"courseattachments/domain"
	"sort"
)

type AttachmentStats struct {
	Total, Pending, Processing, Complete, Cancelled int
	ByKind                                          map[string]int
}

func (r *BoltRepository) Stats() (AttachmentStats, error) {
	attachments, err := r.ListAttachments()
	if err != nil {
		return AttachmentStats{}, err
	}
	stats := AttachmentStats{ByKind: map[string]int{}}
	for _, attachment := range attachments {
		stats.Total++
		stats.ByKind[attachment.Kind]++
		switch attachment.Status {
		case domain.StatusUnprocessed:
			stats.Pending++
		case domain.StatusProcessing:
			stats.Processing++
		case domain.StatusComplete:
			stats.Complete++
		case domain.StatusCancelled:
			stats.Cancelled++
		}
	}
	return stats, nil
}

func (s AttachmentStats) ReadyRatio() float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.Complete) / float64(s.Total)
}

func (s AttachmentStats) Kinds() []string {
	result := make([]string, 0, len(s.ByKind))
	for kind := range s.ByKind {
		result = append(result, kind)
	}
	sort.Strings(result)
	return result
}

func (r *BoltRepository) RepairMissingTasks() (int, error) {
	attachments, err := r.ListAttachments()
	if err != nil {
		return 0, err
	}
	repaired := 0
	for _, attachment := range attachments {
		if _, taskErr := r.GetTask("task-" + attachment.ID); taskErr == nil {
			continue
		}
		task := domain.SummaryTask{ID: "task-" + attachment.ID, AttachmentID: attachment.ID, State: domain.TaskQueued}
		if err := r.SaveTask(task); err != nil {
			return repaired, err
		}
		repaired++
	}
	return repaired, nil
}
