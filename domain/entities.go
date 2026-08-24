package domain

import (
	"fmt"
	"strings"
)

type Attachment struct {
	ID, CourseID, StudentID, Filename, Kind, Keywords, Status, Summary string
	CreatedAt                                                          int64
}
type SummaryTask struct {
	ID, AttachmentID, State, Error string
	Progress                       int
}
type Course struct {
	ID, Code, Title string
	Active          bool
}
type Student struct{ ID, Name, Email string }

type AuditEvent struct {
	ID, AttachmentID, Action, Detail string
	Sequence                         int
}

const (
	StatusUnprocessed = "unprocessed"
	StatusProcessing  = "processing"
	StatusComplete    = "complete"
	StatusCancelled   = "cancelled"
)

const (
	TaskQueued    = "queued"
	TaskRunning   = "running"
	TaskCancelled = "cancelled"
	TaskFinished  = "finished"
)

func (a Attachment) Validate() error {
	if a.ID == "" || a.CourseID == "" || a.StudentID == "" {
		return fmt.Errorf("missing attachment identity")
	}
	if a.Filename == "" {
		return fmt.Errorf("missing filename")
	}
	if a.Kind != "document" && a.Kind != "image" && a.Kind != "archive" {
		return fmt.Errorf("unsupported kind")
	}
	if strings.TrimSpace(a.Filename) != a.Filename {
		return fmt.Errorf("filename must be normalized")
	}
	if a.CreatedAt < 0 {
		return fmt.Errorf("invalid creation marker")
	}
	return nil
}
func (c Course) Validate() error {
	if c.ID == "" || c.Code == "" || c.Title == "" {
		return fmt.Errorf("invalid course")
	}
	return nil
}
func (s Student) Validate() error {
	if s.ID == "" || s.Name == "" || s.Email == "" {
		return fmt.Errorf("invalid student")
	}
	return nil
}
func (t SummaryTask) Validate() error {
	if t.ID == "" || t.AttachmentID == "" {
		return fmt.Errorf("invalid task")
	}
	return nil
}

func (e AuditEvent) Validate() error {
	if e.ID == "" || e.AttachmentID == "" || e.Action == "" {
		return fmt.Errorf("invalid audit event")
	}
	if e.Sequence < 1 {
		return fmt.Errorf("invalid audit sequence")
	}
	return nil
}
func NewAttachment(id, course, student, name, kind, keywords string) Attachment {
	return Attachment{ID: id, CourseID: course, StudentID: student, Filename: name, Kind: kind, Keywords: keywords, Status: StatusUnprocessed}
}
func (a Attachment) IsTerminal() bool {
	return a.Status == StatusComplete || a.Status == StatusCancelled
}
func (a *Attachment) MarkProcessing() {
	if a.Status == StatusUnprocessed {
		a.Status = StatusProcessing
	}
}
func (a *Attachment) MarkComplete(summary string) {
	if a.Status == StatusProcessing {
		a.Status = StatusComplete
		a.Summary = summary
	}
}
func (a *Attachment) MarkCancelled() {
	if a.Status == StatusProcessing {
		a.Status = StatusUnprocessed
	}
}

func (a Attachment) HasKeyword(keyword string) bool {
	needle := NormalizeKeyword(keyword)
	if needle == "" {
		return true
	}
	for _, token := range TokenizeKeywords(a.Keywords + " " + a.Filename) {
		if token == needle {
			return true
		}
	}
	return false
}

func (a Attachment) StatusLabel() string {
	switch a.Status {
	case StatusUnprocessed:
		return "Waiting for summary"
	case StatusProcessing:
		return "Summary in progress"
	case StatusComplete:
		return "Summary ready"
	case StatusCancelled:
		return "Summary cancelled"
	default:
		return "Unknown status"
	}
}
