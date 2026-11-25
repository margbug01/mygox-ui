package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

func TestGlideTipFollowsThePointerAndDrainsAway(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.GlideTip(c)
	}
	tt := ui.NewTester(view, 184, 80)
	tt.Frame()
	if len(tt.Texts()) != 0 {
		t.Fatalf("at rest want no tooltip, got %v", tt.Texts())
	}
	tt.Move(68, 60) // the share button
	step(tt, 500*time.Millisecond)
	if !tt.HasText("Share") {
		t.Fatalf("pointing at share want its tooltip, got %v", tt.Texts())
	}
	tt.Move(164, 60) // along the row to delete
	step(tt, 500*time.Millisecond)
	if !tt.HasText("Delete") || tt.HasText("Share") {
		t.Fatalf("moving to delete want the tooltip to follow, got %v", tt.Texts())
	}
	tt.Move(300, 300)
	step(tt, 600*time.Millisecond)
	if len(tt.Texts()) != 0 {
		t.Fatalf("after the pointer leaves want the tooltip gone, got %v", tt.Texts())
	}
}
