package ingest

import (
	"courseattachments/domain"
	"courseattachments/store"
	"os"
	"testing"
)

func TestWorkflowIngest(t *testing.T) {
	p := "ingest.db"
	defer os.Remove(p)
	r, _ := store.Open(p)
	defer r.Close()
	s := New(r)
	if e := s.RegisterCourse(domain.Course{ID: "c", Code: "C", Title: "Course"}); e != nil {
		t.Fatal(e)
	}
	if e := s.RegisterStudent(domain.Student{ID: "s", Name: "Stu", Email: "s@x"}); e != nil {
		t.Fatal(e)
	}
	if _, e := s.SubmitWithTask(domain.NewAttachment("a", "c", "s", "x.pdf", "document", "hello")); e != nil {
		t.Fatal(e)
	}
}
