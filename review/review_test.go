package review

import (
	"courseattachments/domain"
	"courseattachments/store"
	"os"
	"strings"
	"testing"
)

func TestTeacherDashboard(t *testing.T) {
	path := "review.db"
	defer os.Remove(path)
	repo, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	a := domain.NewAttachment("a", "c", "s", "essay.pdf", "document", "systems")
	if err := repo.SaveAttachment(a); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveTask(domain.SummaryTask{ID: "task-a", AttachmentID: "a", State: domain.TaskQueued}); err != nil {
		t.Fatal(err)
	}
	dashboard, err := New(repo).Build(domain.SearchFilter{Keyword: "systems"})
	if err != nil || dashboard.NeedsReview != 1 || !dashboard.HasAttention() {
		t.Fatalf("dashboard=%+v err=%v", dashboard, err)
	}
	if len(dashboard.AttentionRows()) != 1 {
		t.Fatal("expected one attention row")
	}
	if len(dashboard.RowsForCourse("c")) != 1 || dashboard.CountByKind()["document"] != 1 {
		t.Fatal("course and kind summaries")
	}
	if csvText, err := dashboard.CSV(); err != nil || !strings.Contains(csvText, "attachment_id") {
		t.Fatalf("csv=%q err=%v", csvText, err)
	}
	if dashboard.SortByFilename().Rows[0].Attachment.Filename != "essay.pdf" {
		t.Fatal("filename sort")
	}
	if dashboard.DescribeKinds() != "document=1" || ReviewMessage(dashboard) == "" {
		t.Fatal("review description")
	}
	if len(dashboard.Filtered(domain.SearchFilter{CourseID: "missing"})) != 0 {
		t.Fatal("filter")
	}
}
