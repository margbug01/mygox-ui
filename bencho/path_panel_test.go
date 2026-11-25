package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging along Links fills it like a slider, tapping Color cycles its
// swatch and tapping Inverted flips the switch.
func TestPathPanelRowsAreControls(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.PathPanel(c)
	}, 300, 164)
	step(tt, 100*time.Millisecond)
	px := func(x, y int) [3]uint8 {
		c := tt.Image().RGBAAt(x, y)
		return [3]uint8{c.R, c.G, c.B}
	}
	row, fill := px(220, 15), px(40, 15)
	if row == fill {
		t.Fatal("Links starts half filled")
	}
	tt.Press(100, 30)
	tt.Move(180, 30)
	tt.Move(270, 30)
	step(tt, 100*time.Millisecond)
	tt.Release(270, 30)
	step(tt, 500*time.Millisecond)
	if px(220, 15) != fill {
		t.Fatal("dragging along Links fills it to the pointer")
	}

	swatch := px(266, 82)
	tt.ClickAt(150, 82)
	step(tt, 600*time.Millisecond)
	if px(266, 82) == swatch {
		t.Fatal("tapping Color moves to the next swatch")
	}

	knob := px(250, 134)
	tt.ClickAt(150, 134)
	step(tt, 600*time.Millisecond)
	if px(250, 134) == knob || px(266, 134) != knob {
		t.Fatal("tapping Inverted slides the switch's knob across")
	}
}
