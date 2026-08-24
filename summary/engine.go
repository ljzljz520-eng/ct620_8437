package summary

import (
	"context"
	"courseattachments/domain"
	"courseattachments/store"
	"fmt"
)

type Engine struct{ Repo *store.BoltRepository }

func New(r *store.BoltRepository) *Engine { return &Engine{Repo: r} }
func (e *Engine) Start(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	attachment, err := e.Repo.GetAttachment(id)
	if err != nil {
		return err
	}
	if attachment.IsTerminal() && attachment.Status != domain.StatusCancelled {
		return fmt.Errorf("attachment %s already complete", id)
	}
	if err := e.Repo.UpdateAttachment(id, func(a *domain.Attachment) { a.MarkProcessing() }); err != nil {
		return err
	}
	if err := e.Repo.SaveTask(domain.SummaryTask{ID: "task-" + id, AttachmentID: id, State: domain.TaskRunning, Progress: 0}); err != nil {
		return err
	}
	return e.Repo.RecordTransition(id, "started", "summary generation started")
}
func (e *Engine) Generate(ctx context.Context, id string) error {
	a, err := e.Repo.GetAttachment(id)
	if err != nil {
		return err
	}
	if a.Status == domain.StatusUnprocessed {
		if err := e.Start(ctx, id); err != nil {
			return err
		}
		a, err = e.Repo.GetAttachment(id)
		if err != nil {
			return err
		}
	}
	if a.Status != domain.StatusProcessing {
		return fmt.Errorf("attachment %s is not processing", id)
	}
	for step := 1; step <= 4; step++ {
		select {
		case <-ctx.Done():
			if err := e.Repo.UpdateAttachment(id, func(a *domain.Attachment) { a.MarkCancelled() }); err != nil {
				return err
			}
			if err := e.Repo.UpdateTask("task-"+id, func(task *domain.SummaryTask) {
				task.State = domain.TaskCancelled
				task.Error = ctx.Err().Error()
				task.Progress = (step - 1) * 25
			}); err != nil {
				return err
			}
			if err := e.Repo.RecordTransition(id, "cancel-observed", "generation context cancelled"); err != nil {
				return err
			}
			// Stop generation: keep the attachment unprocessed instead of committing success.
			return ctx.Err()
		default:
		}
		if err := e.Repo.UpdateTask("task-"+id, func(task *domain.SummaryTask) { task.State = domain.TaskRunning; task.Progress = step * 25 }); err != nil {
			return err
		}
	}
	result := fmt.Sprintf("Simulated summary for %s (%s): %s", a.Filename, a.Kind, summarizeKeywords(a.Keywords))
	if err := e.Repo.UpdateAttachment(id, func(value *domain.Attachment) { value.MarkComplete(result) }); err != nil {
		return err
	}
	if err := e.Repo.UpdateTask("task-"+id, func(task *domain.SummaryTask) { task.State = domain.TaskFinished; task.Progress = 100 }); err != nil {
		return err
	}
	return e.Repo.RecordTransition(id, "completed", "summary persisted")
}
func (e *Engine) cancel(id string) error {
	if _, err := e.Repo.GetAttachment(id); err != nil {
		return err
	}
	if err := e.Repo.UpdateAttachment(id, func(a *domain.Attachment) { a.MarkCancelled() }); err != nil {
		return err
	}
	if err := e.Repo.UpdateTask("task-"+id, func(task *domain.SummaryTask) {
		task.State = domain.TaskCancelled
		task.Error = context.Canceled.Error()
		task.Progress = 0
	}); err != nil {
		return err
	}
	_ = e.Repo.RecordTransition(id, "cancelled", "teacher cancelled summary")
	return context.Canceled
}
func (e *Engine) Cancel(id string) error { return e.cancel(id) }
func (e *Engine) Status(id string) (string, error) {
	a, e2 := e.Repo.GetAttachment(id)
	return a.Status, e2
}

func (e *Engine) Task(id string) (domain.SummaryTask, error) { return e.Repo.GetTask("task-" + id) }

func (e *Engine) Pending(id string) bool {
	a, err := e.Repo.GetAttachment(id)
	return err == nil && a.Status == domain.StatusUnprocessed
}

func (e *Engine) CanRetry(id string) bool {
	a, err := e.Repo.GetAttachment(id)
	return err == nil && (a.Status == domain.StatusUnprocessed || a.Status == domain.StatusCancelled)
}

func (e *Engine) Reset(id string) error {
	if err := e.Repo.UpdateAttachment(id, func(a *domain.Attachment) { a.Status = domain.StatusUnprocessed; a.Summary = "" }); err != nil {
		return err
	}
	return e.Repo.UpdateTask("task-"+id, func(task *domain.SummaryTask) { task.State = domain.TaskQueued; task.Error = ""; task.Progress = 0 })
}

func (e *Engine) Report(id string) (Report, error) {
	a, err := e.Repo.GetAttachment(id)
	if err != nil {
		return Report{}, err
	}
	task, err := e.Task(id)
	if err != nil {
		return Report{}, err
	}
	events, err := e.Repo.ListEvents(id)
	if err != nil {
		return Report{}, err
	}
	return NewReport(a, task, events), nil
}

func summarizeKeywords(keywords string) string {
	tokens := domain.UniqueKeywords(keywords)
	if len(tokens) == 0 {
		return "no keywords"
	}
	if len(tokens) > 3 {
		tokens = tokens[:3]
	}
	result := tokens[0]
	for _, token := range tokens[1:] {
		result += ", " + token
	}
	return result
}
