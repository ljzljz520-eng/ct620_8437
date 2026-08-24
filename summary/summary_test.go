package summary

import (
	"context"
	"courseattachments/domain"
	"courseattachments/store"
	"os"
	"testing"
)

func TestCancelledAttachmentSummaryStops(t *testing.T) {
	p := "summary.db"
	defer os.Remove(p)
	r, _ := store.Open(p)
	defer r.Close()
	r.SaveAttachment(domain.NewAttachment("a", "c", "s", "x.pdf", "document", "k"))
	e := New(r)
	if x := e.Start(context.Background(), "a"); x != nil {
		t.Fatal(x)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_ = e.Generate(ctx, "a")
	a, _ := r.GetAttachment("a")
	if a.Status == domain.StatusComplete {
		t.Fatal("cancelled task completed")
	}
}
