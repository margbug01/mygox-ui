package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Improve swaps the selection for the rewrite after the think time; Explain
// opens the pill into a note and its cross closes it.
func TestSelectionActionsImproveAndExplain(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.SelectionActions(c, "Before.", "old words", "new words", "Why it reads better.",
			bencho.SelectionActionsOptions{Think: 500 * time.Millisecond})
	}, 320, 240)
	step(tt, 50*time.Millisecond)
	tt.Click("Improve")
	step(tt, 300*time.Millisecond)
	if !tt.HasText("Improving") || !tt.HasText("old") {
		t.Fatal("Improve thinks over the old words first")
	}
	step(tt, 1500*time.Millisecond)
	if tt.HasText("old") || !tt.HasText("new") {
		t.Fatal("then the rewrite takes their place")
	}
	tt.Click("Explain")
	step(tt, 600*time.Millisecond)
	if !tt.HasText("Why it reads better.") || tt.HasText("Explain") {
		t.Fatal("Explain opens the pill into the note")
	}
}
