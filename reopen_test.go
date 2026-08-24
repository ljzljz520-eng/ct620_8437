package courseattachments

import (
	"courseattachments/domain"
	"courseattachments/store"
	"os"
	"testing"
)

func TestWorkflowOne(t *testing.T) {
	p := "wf.db"
	defer os.Remove(p)
	r, _ := store.Open(p)
	r.SaveAttachment(domain.NewAttachment("x", "c", "s", "x.txt", "document", "k"))
	r.Close()
	r, _ = store.Open(p)
	defer r.Close()
	if _, e := r.GetAttachment("x"); e != nil {
		t.Fatal(e)
	}
}
