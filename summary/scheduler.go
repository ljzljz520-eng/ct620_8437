package summary

import (
	"context"
	"courseattachments/domain"
	"courseattachments/store"
)

type Scheduler struct {
	Engine *Engine
	Repo   *store.BoltRepository
}

func NewScheduler(repo *store.BoltRepository) *Scheduler {
	return &Scheduler{Engine: New(repo), Repo: repo}
}

func (s *Scheduler) RunOne(ctx context.Context, id string) error {
	if err := s.Engine.Start(ctx, id); err != nil {
		return err
	}
	return s.Engine.Generate(ctx, id)
}

func (s *Scheduler) RunPending(ctx context.Context, limit int) ([]string, error) {
	pending, err := s.Repo.FindPending()
	if err != nil {
		return nil, err
	}
	completed := []string{}
	for _, attachment := range pending {
		if limit > 0 && len(completed) >= limit {
			break
		}
		if err := s.RunOne(ctx, attachment.ID); err != nil {
			return completed, err
		}
		completed = append(completed, attachment.ID)
	}
	return completed, nil
}

func (s *Scheduler) QueueSize() (int, error) {
	pending, err := s.Repo.FindPending()
	if err != nil {
		return 0, err
	}
	return len(pending), nil
}

func (s *Scheduler) CancelAndReport(id string) (Report, error) {
	_ = s.Engine.Cancel(id)
	return s.Engine.Report(id)
}

func TaskStateLabel(task domain.SummaryTask) string {
	switch task.State {
	case domain.TaskQueued:
		return "queued"
	case domain.TaskRunning:
		return "running"
	case domain.TaskCancelled:
		return "cancelled"
	case domain.TaskFinished:
		return "finished"
	default:
		return "unknown"
	}
}
