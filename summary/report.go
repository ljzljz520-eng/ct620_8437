package summary

import (
	"context"
	"courseattachments/domain"
	"courseattachments/store"
	"fmt"
	"strings"
)

type Report struct {
	Attachment domain.Attachment
	Task       domain.SummaryTask
	Events     []domain.AuditEvent
}

func NewReport(attachment domain.Attachment, task domain.SummaryTask, events []domain.AuditEvent) Report {
	return Report{Attachment: attachment, Task: task, Events: append([]domain.AuditEvent(nil), events...)}
}

func BuildPreview(a domain.Attachment) string {
	if a.Summary != "" {
		return a.Summary
	}
	return "Summary pending for " + a.Filename
}

func BuildDetailedPreview(a domain.Attachment) string {
	parts := []string{a.DisplayLabel(), a.StatusLabel()}
	if a.Summary != "" {
		parts = append(parts, a.Summary)
	}
	return strings.Join(parts, " | ")
}
func CanCancel(a domain.Attachment) bool { return a.Status == domain.StatusProcessing }
func EnsureStopped(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
func ProgressLabel(p int) string {
	if p >= 100 {
		return "complete"
	}
	if p <= 0 {
		return "pending"
	}
	return "processing"
}

func (r Report) Ready() bool {
	return r.Attachment.Status == domain.StatusComplete && r.Task.State == domain.TaskFinished
}

func (r Report) Cancelled() bool {
	return r.Task.State == domain.TaskCancelled || r.Attachment.Status == domain.StatusCancelled
}

func (r Report) Timeline() string {
	parts := []string{}
	for _, event := range r.Events {
		parts = append(parts, fmt.Sprintf("%d %s", event.Sequence, event.Action))
	}
	return strings.Join(parts, " -> ")
}

func (r Report) Text() string {
	return fmt.Sprintf("%s\nstate=%s progress=%d\n%s", BuildDetailedPreview(r.Attachment), r.Task.State, r.Task.Progress, r.Timeline())
}

func LoadReport(repo *store.BoltRepository, id string) (Report, error) {
	a, err := repo.GetAttachment(id)
	if err != nil {
		return Report{}, err
	}
	task, err := repo.GetTask("task-" + id)
	if err != nil {
		return Report{}, err
	}
	events, err := repo.ListEvents(id)
	if err != nil {
		return Report{}, err
	}
	return NewReport(a, task, events), nil
}
