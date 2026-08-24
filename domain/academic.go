package domain

import (
	"fmt"
	"sort"
	"strings"
)

type AttachmentInput struct {
	Filename string
	Kind     string
	Bytes    int
	Keywords string
}

type SubmissionPolicy struct {
	MaxBytes       int
	AllowedKinds   []string
	RequireKeyword bool
}

func DefaultSubmissionPolicy() SubmissionPolicy {
	return SubmissionPolicy{MaxBytes: 25 * 1024 * 1024, AllowedKinds: []string{"document", "image", "archive"}, RequireKeyword: false}
}

func (p SubmissionPolicy) AllowsKind(kind string) bool {
	for _, allowed := range p.AllowedKinds {
		if strings.EqualFold(strings.TrimSpace(allowed), strings.TrimSpace(kind)) {
			return true
		}
	}
	return false
}

func (p SubmissionPolicy) Validate(input AttachmentInput) error {
	if strings.TrimSpace(input.Filename) == "" {
		return fmt.Errorf("filename required")
	}
	if input.Bytes < 0 {
		return fmt.Errorf("attachment size cannot be negative")
	}
	if p.MaxBytes > 0 && input.Bytes > p.MaxBytes {
		return fmt.Errorf("attachment exceeds %d bytes", p.MaxBytes)
	}
	if !p.AllowsKind(input.Kind) {
		return fmt.Errorf("attachment kind %q is not allowed", input.Kind)
	}
	if p.RequireKeyword && len(TokenizeKeywords(input.Keywords)) == 0 {
		return fmt.Errorf("at least one keyword is required")
	}
	return nil
}

type CourseCatalog struct {
	Courses  []Course
	Students []Student
}

func NewCourseCatalog(courses []Course, students []Student) CourseCatalog {
	c := CourseCatalog{Courses: append([]Course(nil), courses...), Students: append([]Student(nil), students...)}
	sort.Slice(c.Courses, func(i, j int) bool { return c.Courses[i].ID < c.Courses[j].ID })
	sort.Slice(c.Students, func(i, j int) bool { return c.Students[i].ID < c.Students[j].ID })
	return c
}

func (c CourseCatalog) FindCourse(id string) (Course, bool) {
	for _, course := range c.Courses {
		if course.ID == id {
			return course, true
		}
	}
	return Course{}, false
}

func (c CourseCatalog) FindStudent(id string) (Student, bool) {
	for _, student := range c.Students {
		if student.ID == id {
			return student, true
		}
	}
	return Student{}, false
}

func (c CourseCatalog) ValidateAttachment(a Attachment) error {
	if _, ok := c.FindCourse(a.CourseID); !ok {
		return fmt.Errorf("course %s is not registered", a.CourseID)
	}
	if _, ok := c.FindStudent(a.StudentID); !ok {
		return fmt.Errorf("student %s is not registered", a.StudentID)
	}
	return a.Validate()
}

func (c CourseCatalog) Labels() []string {
	labels := make([]string, 0, len(c.Courses)+len(c.Students))
	for _, course := range c.Courses {
		labels = append(labels, course.Label())
	}
	for _, student := range c.Students {
		labels = append(labels, student.Label())
	}
	sort.Strings(labels)
	return labels
}

func NormalizeKeywordsForDisplay(keywords string) string {
	return strings.Join(UniqueKeywords(keywords), ", ")
}
