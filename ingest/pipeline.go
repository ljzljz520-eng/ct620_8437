package ingest

import (
	"courseattachments/domain"
	"courseattachments/store"
	"fmt"
	"strings"
)

type IntakeReceipt struct {
	AttachmentID string
	TaskID       string
	Kind         string
	Display      string
}

type Pipeline struct {
	Service *Service
	Repo    *store.BoltRepository
}

func NewPipeline(repo *store.BoltRepository) *Pipeline {
	return &Pipeline{Service: New(repo), Repo: repo}
}

func (p *Pipeline) Accept(course domain.Course, student domain.Student, filename, keywords string) (IntakeReceipt, error) {
	if !SupportedExtension(filename) {
		return IntakeReceipt{}, fmt.Errorf("unsupported attachment extension")
	}
	clean := NormalizeName(filename)
	attachment := domain.NewAttachment("", course.ID, student.ID, clean, CanonicalKind(clean), keywords)
	attachment.ID = p.nextID(course.ID, student.ID, clean)
	if err := p.Service.RegisterCourse(course); err != nil {
		return IntakeReceipt{}, err
	}
	if err := p.Service.RegisterStudent(student); err != nil {
		return IntakeReceipt{}, err
	}
	task, err := p.Service.SubmitWithTask(attachment)
	if err != nil {
		return IntakeReceipt{}, err
	}
	return IntakeReceipt{AttachmentID: attachment.ID, TaskID: task.ID, Kind: attachment.Kind, Display: attachment.DisplayLabel()}, nil
}

func (p *Pipeline) nextID(course, student, filename string) string {
	base := strings.ReplaceAll(strings.ToLower(course+"-"+student+"-"+filename), " ", "-")
	for index := 1; ; index++ {
		candidate := fmt.Sprintf("%s-%d", base, index)
		if _, err := p.Repo.GetAttachment(candidate); err != nil {
			return candidate
		}
	}
}

func (p *Pipeline) AttachmentsForStudent(studentID string) ([]domain.Attachment, error) {
	return p.Repo.FindByStudent(studentID)
}

func (p *Pipeline) ValidateCourseStudent(course domain.Course, student domain.Student) error {
	if err := course.Validate(); err != nil {
		return err
	}
	if err := student.Validate(); err != nil {
		return err
	}
	return nil
}

func (p *Pipeline) ReceiptLabel(receipt IntakeReceipt) string {
	return fmt.Sprintf("%s [%s] task=%s", receipt.Display, receipt.Kind, receipt.TaskID)
}
