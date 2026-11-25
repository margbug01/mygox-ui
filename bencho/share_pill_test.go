package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pointing at a face names it; the link button copies the file's link; and
// Invite opens an email field in place that sends and closes again.
func TestSharePillNamesCopiesAndInvites(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.SharePill(c, bencho.SharePillDefaultMembers(), bencho.SharePillOptions{Link: "https://example.com/f/1"})
	}, 300, 122)
	tt.Frame()
	tt.Move(100, 82)
	step(tt, 300*time.Millisecond)
	if !tt.HasText("Ava Chen") {
		t.Fatalf("pointing at the third face names it, got %v", tt.Texts())
	}
	tt.ClickAt(182, 82)
	step(tt, 200*time.Millisecond)
	if tt.Clipboard() != "https://example.com/f/1" || !tt.HasText("Link copied") {
		t.Fatalf("the link button copies the link, got %q and %v", tt.Clipboard(), tt.Texts())
	}
	if err := tt.Click("Invite"); err != nil {
		t.Fatal(err)
	}
	step(tt, 400*time.Millisecond)
	tt.Type("ines@example.com")
	if err := tt.Click("Send"); err != nil {
		t.Fatal(err)
	}
	step(tt, 400*time.Millisecond)
	if !tt.HasText("Invite sent") || !tt.HasText("Invite") || tt.HasText("Send") {
		t.Fatalf("sending closes the field and says so, got %v", tt.Texts())
	}
}
