package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing Notify signs you up: the bell swings and then rests, and the
// button stays pressed in until you press it again.
func TestNotifyBellSwingsAndButtonStaysPressed(t *testing.T) {
	signed := false
	view := func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(200, 60).Center().Background(th.SurfaceSubtle).Children(func() {
			bencho.Notify(c, bencho.NotifyOptions{OnChange: func(on bool) { signed = on }})
		})
	}
	tt := ui.NewTester(view, 200, 60)
	tt.Frame()
	r, ok := tt.Find("Notify me") // the button, which the label names
	if !ok {
		t.Fatal("no Notify me")
	}
	face := func() [3]uint8 { // the button's fill, right of its label
		px := tt.Image().RGBAAt(int(r.X+r.W)-10, int(r.Y+r.H/2))
		return [3]uint8{px.R, px.G, px.B}
	}
	bell := func() int { // where the dark pixels at the button's left are: the bell, wherever it swings
		img := tt.Image()
		n := 0
		for y := int(r.Y); y < int(r.Y+r.H); y++ {
			for x := int(r.X) + 8; x < int(r.X)+44; x++ {
				if img.RGBAAt(x, y).G < 100 {
					n += x * y
				}
			}
		}
		return n
	}
	rest := face()
	tt.Move(2, 2)

	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	step(tt, 96*time.Millisecond)
	a := bell()
	step(tt, 96*time.Millisecond)
	if !signed || bell() == a {
		t.Fatalf("pressing should sign up and swing the bell (signed=%v)", signed)
	}
	step(tt, 2*time.Second)
	a = bell()
	step(tt, 200*time.Millisecond)
	if bell() != a {
		t.Fatal("the bell should come to rest")
	}
	if face() == rest {
		t.Fatal("signed up, the button should stay pressed in")
	}
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	tt.Move(2, 2)
	step(tt, time.Second)
	if signed || face() != rest {
		t.Fatalf("pressing again should sign out and let the button up (signed=%v)", signed)
	}
}
