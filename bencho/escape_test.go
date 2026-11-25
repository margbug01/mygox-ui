package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Reaching for the button makes it slide away; after four dodges it gives in
// and a press on it lands.
func TestEscapeButtonDodgesThenGivesIn(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		ui.Box(c).Size(340, 280).Children(func() { bencho.EscapeButton(c) })
	}, 340, 280)
	tt.Frame()
	for i := 0; i < 4; i++ {
		r, ok := tt.Find("Touch me")
		if !ok {
			t.Fatalf("dodge %d: no button, got %v", i, tt.Texts())
		}
		x, y := r.X+r.W/2, r.Y+r.H/2
		tt.Move(x, y)
		step(tt, 700*time.Millisecond)
		moved, _ := tt.Find("Touch me")
		if mx, my := moved.X+moved.W/2, moved.Y+moved.H/2; (mx-x)*(mx-x)+(my-y)*(my-y) < 60*60 {
			t.Fatalf("dodge %d: the button slides away from the pointer, it only went from (%v,%v) to (%v,%v)", i, x, y, mx, my)
		}
		tt.Move(335, 275) // off the stage
		step(tt, 100*time.Millisecond)
	}
	r, _ := tt.Find("Touch me")
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	step(tt, 300*time.Millisecond)
	if !tt.HasText("You got me") {
		t.Fatalf("after four dodges it gives in, got %v", tt.Texts())
	}
}
