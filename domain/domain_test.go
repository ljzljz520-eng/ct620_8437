package domain

import "testing"

func TestEntityValidation(t *testing.T) {
	a := NewAttachment("a", "c", "s", "x.pdf", "document", "k")
	if e := a.Validate(); e != nil {
		t.Fatal(e)
	}
	a.MarkProcessing()
	a.MarkComplete("ok")
	if !a.IsTerminal() {
		t.Fatal("not terminal")
	}
}
