package bencho_test

import (
	"image/color"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging the turned-down corner peels the page back and uncovers the
// note underneath; let go, the page falls flat again.
func TestPagePeelUncoversTheNoteAndFallsBack(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(260, 320).Background(th.SurfaceSubtle).Children(func() { bencho.PagePeel(c) })
	}, 260, 320)
	step(tt, 100*time.Millisecond)
	note := func(c color.RGBA) bool { return c.R < 80 && c.G > c.R+10 } // dark teal
	if c := tt.Image().RGBAAt(256, 4); c.R == 255 {
		t.Fatal("at rest the corner is turned down off the stage")
	}
	if note(tt.Image().RGBAAt(225, 45)) {
		t.Fatal("at rest the note is hidden")
	}
	tt.Press(245, 25)
	for i := float32(0); i <= 10; i++ {
		tt.Move(245-i*12, 25+i*12)
		step(tt, 16*time.Millisecond)
	}
	step(tt, 100*time.Millisecond)
	if !note(tt.Image().RGBAAt(225, 45)) {
		t.Fatalf("peeling the corner back shows the note, have %v", tt.Image().RGBAAt(225, 45))
	}
	tt.Release(125, 145)
	step(tt, 1200*time.Millisecond)
	if note(tt.Image().RGBAAt(225, 45)) {
		t.Fatal("let go, the page falls flat over the note")
	}
}
