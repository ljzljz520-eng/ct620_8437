package courseattachments

import (
	"context"
	"courseattachments/domain"
	"courseattachments/store"
	"courseattachments/summary"
	"os"
	"testing"
)

func TestCancelledAttachmentSummaryStops(t *testing.T) {
	path := "cancel-regression.db"
	defer os.Remove(path)
	repo, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	if err := repo.SaveAttachment(domain.NewAttachment("cancel-a", "course-1", "student-1", "essay.pdf", "document", "cancellation")); err != nil {
		t.Fatal(err)
	}
	engine := summary.New(repo)
	if err := engine.Start(context.Background(), "cancel-a"); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_ = engine.Generate(ctx, "cancel-a")
	attachment, err := repo.GetAttachment("cancel-a")
	if err != nil {
		t.Fatal(err)
	}
	if attachment.Status != domain.StatusUnprocessed {
		t.Fatalf("expected unprocessed after cancellation, got %s", attachment.Status)
	}
}
