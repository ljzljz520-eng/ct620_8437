package fixture

import "courseattachments/domain"

func Course() domain.Course {
	return domain.Course{ID: "course-1", Code: "CS101", Title: "Systems", Active: true}
}
func Student() domain.Student {
	return domain.Student{ID: "student-1", Name: "Ada", Email: "ada@example.test"}
}
func Attachment(id, name, keywords string) domain.Attachment {
	return domain.NewAttachment(id, "course-1", "student-1", name, "document", keywords)
}
func Attachments() []domain.Attachment {
	return []domain.Attachment{Attachment("a1", "report.pdf", "systems concurrency"), Attachment("a2", "diagram.png", "systems architecture"), Attachment("a3", "bundle.zip", "submission archive")}
}
func IDs() []string { return []string{"a1", "a2", "a3"} }

func AttachmentsForCourse(courseID, studentID string) []domain.Attachment {
	result := make([]domain.Attachment, 0, len(Attachments()))
	for _, attachment := range Attachments() {
		attachment.CourseID = courseID
		attachment.StudentID = studentID
		result = append(result, attachment)
	}
	return result
}

func CompleteAttachment(id string) domain.Attachment {
	attachment := Attachment(id, "report.pdf", "systems concurrency")
	attachment.MarkProcessing()
	attachment.MarkComplete("deterministic fixture summary")
	return attachment
}

func PendingAttachment(id string) domain.Attachment { return Attachment(id, "notes.txt", "review") }

func CourseWithCode(id, code string) domain.Course {
	course := Course()
	course.ID = id
	course.Code = code
	return course
}

func StudentWithName(id, name string) domain.Student {
	student := Student()
	student.ID = id
	student.Name = name
	student.Email = id + "@example.test"
	return student
}
