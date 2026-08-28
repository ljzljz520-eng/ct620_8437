package courseattachments

import (
	"context"
	"courseattachments/cli"
	"os"
	"testing"
)

func TestWorkflowThree(t *testing.T) {
	p := "int.db"
	defer os.Remove(p)
	a, _ := cli.New(p)
	defer a.Close()
	a.Seed()
	if e := a.Generate(context.Background(), "a1"); e != nil {
		t.Fatal(e)
	}
}
