package ingest

import (
	"courseattachments/domain"
	"courseattachments/store"
	"fmt"
	"strings"
)

type Service struct{ Repo *store.BoltRepository }

func New(r *store.BoltRepository) *Service { return &Service{Repo: r} }
func (s *Service) RegisterCourse(c domain.Course) error {
	if e := c.Validate(); e != nil {
		return e
	}
	return s.Repo.SaveCourse(c)
}
func (s *Service) RegisterStudent(st domain.Student) error {
	if e := st.Validate(); e != nil {
		return e
	}
	return s.Repo.SaveStudent(st)
}
func (s *Service) Submit(a domain.Attachment) error {
	a.Filename = NormalizeName(a.Filename)
	if a.Kind == "" {
		a.Kind = CanonicalKind(a.Filename)
	}
	if e := a.Validate(); e != nil {
		return e
	}
	if a.Status == "" {
		a.Status = domain.StatusUnprocessed
	}
	return s.Repo.SaveAttachment(a)
}
func (s *Service) SubmitWithTask(a domain.Attachment) (domain.SummaryTask, error) {
	if e := s.Submit(a); e != nil {
		return domain.SummaryTask{}, e
	}
	t := domain.SummaryTask{ID: "task-" + a.ID, AttachmentID: a.ID, State: domain.TaskQueued}
	if err := s.Repo.SaveTask(t); err != nil {
		return t, err
	}
	_, err := s.Repo.AppendEvent(a.ID, "submitted", "attachment accepted")
	return t, err
}
func ValidateKind(name string) string {
	return CanonicalKind(strings.ToLower(strings.TrimSpace(name)))
}
func EnsureExists(r *store.BoltRepository, id string) error {
	if _, e := r.GetAttachment(id); e != nil {
		return fmt.Errorf("attachment %s: %w", id, e)
	}
	return nil
}

func (s *Service) RegisterBatch(course domain.Course, student domain.Student, attachments []domain.Attachment) (int, error) {
	if err := s.RegisterCourse(course); err != nil {
		return 0, err
	}
	if err := s.RegisterStudent(student); err != nil {
		return 0, err
	}
	accepted := 0
	for _, attachment := range attachments {
		if attachment.CourseID == "" {
			attachment.CourseID = course.ID
		}
		if attachment.StudentID == "" {
			attachment.StudentID = student.ID
		}
		if _, err := s.SubmitWithTask(attachment); err != nil {
			return accepted, err
		}
		accepted++
	}
	return accepted, nil
}

func (s *Service) ReplaceKeywords(id, keywords string) error {
	if strings.TrimSpace(keywords) == "" {
		return fmt.Errorf("keywords required")
	}
	return s.Repo.UpdateAttachment(id, func(a *domain.Attachment) { a.Keywords = strings.Join(domain.UniqueKeywords(keywords), " ") })
}

func (s *Service) ValidateSubmission(a domain.Attachment) []string {
	errors := []string{}
	if a.ID == "" {
		errors = append(errors, "id")
	}
	if a.CourseID == "" {
		errors = append(errors, "course")
	}
	if a.StudentID == "" {
		errors = append(errors, "student")
	}
	if !SupportedExtension(a.Filename) {
		errors = append(errors, "extension")
	}
	if a.Filename != NormalizeName(a.Filename) {
		errors = append(errors, "filename")
	}
	return errors
}
