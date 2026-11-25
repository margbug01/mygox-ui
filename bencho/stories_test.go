package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

func TestStoriesPlaysSkipsAndPausesOnHold(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Stories(c)
	}
	tt := ui.NewTester(view, 280, 400)
	tt.Frame()
	if !tt.HasText("Day 1") {
		t.Fatalf("want the first page, got %v", tt.Texts())
	}
	// The page moves on by itself once its 3s are up.
	step(tt, 3600*time.Millisecond)
	if !tt.HasText("Day 2") || tt.HasText("Day 1") {
		t.Fatalf("after 3s want Day 2, got %v", tt.Texts())
	}
	// A tap on the right skips ahead.
	tt.ClickAt(220, 200)
	step(tt, 600*time.Millisecond)
	if !tt.HasText("Day 3") || tt.HasText("Day 2") {
		t.Fatalf("after a tap on the right want Day 3, got %v", tt.Texts())
	}
	// Holding pauses: 4s held keeps Day 3, and letting go does not skip.
	tt.Move(220, 200)
	tt.Press(220, 200)
	step(tt, 4*time.Second)
	tt.Release(220, 200)
	step(tt, 600*time.Millisecond)
	if !tt.HasText("Day 3") {
		t.Fatalf("holding should pause on Day 3, got %v", tt.Texts())
	}
	// A tap on the left goes back.
	tt.ClickAt(40, 200)
	step(tt, 600*time.Millisecond)
	if !tt.HasText("Day 2") || tt.HasText("Day 3") {
		t.Fatalf("after a tap on the left want Day 2, got %v", tt.Texts())
	}
}
