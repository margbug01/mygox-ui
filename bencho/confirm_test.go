package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Delete turns into its own confirm; confirming deletes and offers an undo;
// an unanswered confirm backs off.
func TestInlineConfirmDeleteConfirmUndo(t *testing.T) {
	deleted, undone := 0, 0
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.InlineConfirm(c, bencho.InlineConfirmOptions{OnDelete: func() { deleted++ }, OnUndo: func() { undone++ }})
	}
	tt := ui.NewTester(view, 160, 44)
	tt.Frame()
	click := func(s string) {
		t.Helper()
		if err := tt.Click(s); err != nil {
			t.Fatalf("%v (texts %v)", err, tt.Texts())
		}
		step(tt, 500*time.Millisecond)
	}
	click("Delete")
	if !tt.HasText("Confirm delete") || deleted != 0 {
		t.Fatalf("Delete should ask to confirm first, texts %v deleted %d", tt.Texts(), deleted)
	}
	click("Confirm delete")
	if deleted != 1 || !tt.HasText("Undo") {
		t.Fatalf("confirming should delete and offer an undo, texts %v deleted %d", tt.Texts(), deleted)
	}
	click("Undo")
	if undone != 1 || !tt.HasText("Delete") || tt.HasText("Undo") {
		t.Fatalf("Undo should bring Delete back, texts %v", tt.Texts())
	}

	click("Delete")
	step(tt, 3500*time.Millisecond)
	if tt.HasText("Confirm delete") || deleted != 1 {
		t.Fatalf("an unanswered confirm should back off to Delete, texts %v", tt.Texts())
	}
}
