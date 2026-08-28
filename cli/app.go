package cli

import (
	"context"
	"courseattachments/audit"
	"courseattachments/domain"
	"courseattachments/ingest"
	"courseattachments/review"
	"courseattachments/search"
	"courseattachments/store"
	"courseattachments/summary"
	"fmt"
)

type App struct {
	Repo    *store.BoltRepository
	Ingest  *ingest.Service
	Search  *search.Service
	Summary *summary.Engine
	Audit   *audit.Ledger
	Review  *review.Service
}

func New(path string) (*App, error) {
	r, e := store.Open(path)
	if e != nil {
		return nil, e
	}
	return &App{Repo: r, Ingest: ingest.New(r), Search: search.New(r), Summary: summary.New(r), Audit: audit.New(r), Review: review.New(r)}, nil
}
func (a *App) Close() error { return a.Repo.Close() }
func (a *App) Seed() error {
	if e := a.Ingest.RegisterCourse(domain.Course{ID: "course-1", Code: "CS101", Title: "Systems", Active: true}); e != nil {
		return e
	}
	if e := a.Ingest.RegisterStudent(domain.Student{ID: "student-1", Name: "Ada", Email: "ada@example.test"}); e != nil {
		return e
	}
	for _, v := range []domain.Attachment{domain.NewAttachment("a1", "course-1", "student-1", "report.pdf", "document", "systems concurrency"), domain.NewAttachment("a2", "course-1", "student-1", "diagram.png", "image", "architecture"), domain.NewAttachment("a3", "course-1", "student-1", "bundle.zip", "archive", "submission archive")} {
		if _, e := a.Ingest.SubmitWithTask(v); e != nil {
			return e
		}
	}
	return nil
}
func (a *App) RunSummary(id string) error { return a.Summary.Generate(context.Background(), id) }
func (a *App) SearchText(k string) (string, error) {
	r, e := a.Search.Summaries(domain.SearchFilter{Keyword: k})
	if e != nil {
		return "", e
	}
	return fmt.Sprint(r), nil
}

func (a *App) SearchResults(course, student, keyword string) ([]search.Result, error) {
	return a.Search.Query(search.ParseFilter(course, student, keyword))
}

func (a *App) Stats() (store.AttachmentStats, error) { return a.Repo.Stats() }

func (a *App) Report(id string) (summary.Report, error) { return a.Summary.Report(id) }

func (a *App) AuditTimeline(id string) (string, error) {
	events, err := a.Audit.Events(id)
	if err != nil {
		return "", err
	}
	return audit.Timeline(events), nil
}

func (a *App) TeacherDashboard(teacherID, course, student, keyword string) (string, error) {
	dashboard, err := a.Review.Build(domain.NewSearchFilter(course, student, keyword))
	if err != nil {
		return "", err
	}
	if teacherID == "" {
		return "", fmt.Errorf("teacher id required")
	}
	return dashboard.Render(), nil
}

func (a *App) Intake(course domain.Course, student domain.Student, filename, keywords string) (string, error) {
	pipeline := ingest.NewPipeline(a.Repo)
	receipt, err := pipeline.Accept(course, student, filename, keywords)
	if err != nil {
		return "", err
	}
	return pipeline.ReceiptLabel(receipt), nil
}
