package cli

import (
	"os"
	"testing"
)

func TestCLICommands(t *testing.T) {
	p := "cli.db"
	defer os.Remove(p)
	a, e := New(p)
	if e != nil {
		t.Fatal(e)
	}
	defer a.Close()
	if e = a.Seed(); e != nil {
		t.Fatal(e)
	}
	if _, e = a.SearchText("systems"); e != nil {
		t.Fatal(e)
	}
}
