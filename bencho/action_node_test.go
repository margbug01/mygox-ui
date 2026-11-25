package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pointing at the card brings its buttons out of the corner one after
// another, and leaving folds them away; a button reports its action.
func TestActionNodeButtonsComeOutOnHover(t *testing.T) {
	var got string
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(340, 300).Background(th.SurfaceSubtle).Children(func() {
			bencho.ActionNode(c, bencho.ActionNodeOptions{OnAction: func(n string) { got = n }})
		})
	}, 340, 300)
	// out tells whether a white button sits on the spot centred at (x, y),
	// looking below its icon, rather than the grey stage.
	out := func(x, y int) bool { return tt.Image().RGBAAt(x, y+12).R == 255 }
	step(tt, 50*time.Millisecond)
	if out(220, 22) || out(302, 106) {
		t.Fatal("the buttons rest folded away")
	}
	tt.Move(130, 140)
	step(tt, 260*time.Millisecond)
	if !out(220, 22) || out(302, 106) {
		t.Fatal("they come out one after another: the first is out before the last")
	}
	step(tt, 600*time.Millisecond)
	if !out(268, 22) || !out(302, 106) {
		t.Fatal("then all four are out")
	}
	tt.Move(220, 22)
	tt.ClickAt(220, 22)
	if got != "Link" {
		t.Fatalf("pressing the first button links, got %q", got)
	}
	tt.Move(100, 290)
	step(tt, time.Second)
	if out(220, 22) || out(302, 106) {
		t.Fatal("leaving folds them back into the corner")
	}
}
