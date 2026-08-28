package search

import (
	"courseattachments/domain"
	"courseattachments/store"
	"os"
	"testing"
)

func TestWorkflowSearch(t *testing.T) {
	p := "search.db"
	defer os.Remove(p)
	r, _ := store.Open(p)
	defer r.Close()
	r.SaveAttachment(domain.NewAttachment("a", "c", "s", "report.pdf", "document", "systems"))
	got, e := New(r).Query(domain.SearchFilter{Keyword: "system"})
	if e != nil || len(got) != 1 {
		t.Fatalf("%v %v", got, e)
	}
}
