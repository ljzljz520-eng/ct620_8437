package fixture

import "testing"

func TestFixtureData(t *testing.T) {
	if len(Attachments()) != 3 || len(IDs()) != 3 {
		t.Fatal("fixtures")
	}
}
