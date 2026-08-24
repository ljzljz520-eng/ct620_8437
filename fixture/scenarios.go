package fixture

import (
	"courseattachments/domain"
	"fmt"
)

func ScenarioAttachments() map[string]domain.Attachment {
	o := map[string]domain.Attachment{}
	for _, a := range Attachments() {
		o[a.ID] = a
	}
	return o
}
func ScenarioNames() []string {
	a := Attachments()
	o := []string{}
	for _, v := range a {
		o = append(o, v.Filename)
	}
	return o
}
func Describe(a domain.Attachment) string { return fmt.Sprintf("%s:%s:%s", a.ID, a.Kind, a.Status) }
func ValidKinds() []string                { return []string{"document", "image", "archive"} }

func ScenarioCourses() []domain.Course {
	return []domain.Course{Course(), CourseWithCode("course-2", "MATH201")}
}

func ScenarioStudents() []domain.Student {
	return []domain.Student{Student(), StudentWithName("student-2", "Grace")}
}

func ScenarioSummary(id string) domain.SummaryTask {
	return domain.SummaryTask{ID: "task-" + id, AttachmentID: id, State: domain.TaskQueued, Progress: 0}
}

func ScenarioLabels() []string {
	labels := []string{}
	for _, attachment := range Attachments() {
		labels = append(labels, attachment.DisplayLabel())
	}
	return labels
}
