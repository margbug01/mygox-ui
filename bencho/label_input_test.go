package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The label sits inside the field until it is pressed, then lifts onto its
// top edge and the field takes typing.
func TestLabelInputLiftsLabelOnFocus(t *testing.T) {
	email := ""
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(280, 65).Background(th.SurfaceSubtle).Children(func() { bencho.LabelInput(c, &email) })
	}, 280, 65)
	tt.Frame()
	mid := func() float32 {
		r, ok := tt.Find("Email")
		if !ok {
			t.Fatal("the label is gone")
		}
		return r.Y + r.H/2
	}
	if y := mid(); y < 32 || y > 44 {
		t.Fatalf("at rest the label sits in the field's middle (38), at %v", y)
	}
	tt.ClickAt(140, 38)
	step(tt, 600*time.Millisecond)
	if y := mid(); y < 9 || y > 15 {
		t.Fatalf("pressed, the label lifts onto the top edge (12), at %v", y)
	}
	tt.Type("sam@studio.com")
	step(tt, 100*time.Millisecond)
	if email != "sam@studio.com" {
		t.Fatalf("the field takes typing, got %q", email)
	}
}
