package store

import (
	"courseattachments/domain"
	"os"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	p := "reopen.db"
	defer os.Remove(p)
	r, e := Open(p)
	if e != nil {
		t.Fatal(e)
	}
	r.SaveCourse(domain.Course{ID: "c", Code: "C", Title: "Course"})
	r.SaveStudent(domain.Student{ID: "s", Name: "Stu", Email: "s@x"})
	r.SaveAttachment(domain.NewAttachment("a", "c", "s", "x.pdf", "document", "k"))
	r.SaveTask(domain.SummaryTask{ID: "task-a", AttachmentID: "a", State: domain.TaskQueued, Progress: 0})
	r.Close()
	r, e = Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer r.Close()
	if _, e = r.GetAttachment("a"); e != nil {
		t.Fatal(e)
	}
	if _, e = r.GetTask("task-a"); e != nil {
		t.Fatal(e)
	}
}
