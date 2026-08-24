package cli

import (
	"context"
	"courseattachments/domain"
	"courseattachments/summary"
)

func (a *App) IngestOne(v domain.Attachment) error { _, e := a.Ingest.SubmitWithTask(v); return e }
func (a *App) Query(course, student, keyword string) ([]string, error) {
	return a.Search.Summaries(domain.SearchFilter{CourseID: course, StudentID: student, Keyword: keyword})
}
func (a *App) Cancel(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return a.Summary.Cancel(id)
	}
}
func (a *App) Generate(ctx context.Context, id string) error { return a.Summary.Generate(ctx, id) }

func (a *App) Start(ctx context.Context, id string) error { return a.Summary.Start(ctx, id) }

func (a *App) Status(id string) (string, error) { return a.Summary.Status(id) }

func (a *App) Preview(id string) (string, error) {
	a2, err := a.Repo.GetAttachment(id)
	if err != nil {
		return "", err
	}
	return summary.BuildDetailedPreview(a2), nil
}

func (a *App) Suggestions(prefix string, limit int) ([]string, error) {
	return a.Search.Suggestions(prefix, limit)
}

func (a *App) Delete(id string) error { return a.Repo.DeleteAttachment(id) }
