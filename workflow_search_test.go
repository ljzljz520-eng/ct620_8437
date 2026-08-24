package courseattachments

import (
	"courseattachments/cli"
	"os"
	"testing"
)

func TestWorkflowTwo(t *testing.T) {
	p := "ws.db"
	defer os.Remove(p)
	a, _ := cli.New(p)
	defer a.Close()
	a.Seed()
	if _, e := a.Query("course-1", "student-1", "architecture"); e != nil {
		t.Fatal(e)
	}
}
