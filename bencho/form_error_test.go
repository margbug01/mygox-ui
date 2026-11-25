package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A wrong address shakes the field and says what is wrong until the next
// key; a right one is sent and the arrow turns into a check.
func TestFormErrorShakesOnWrongEmail(t *testing.T) {
	var sent string
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 140).Background(th.SurfaceSubtle).Children(func() {
			bencho.FormError(c, bencho.FormErrorOptions{OnSubmit: func(e string) { sent = e }})
		})
	}, 300, 140)
	tt.Frame()
	left := func() int { // the field's left edge along its middle
		img := tt.Image()
		for x := range 300 {
			if img.RGBAAt(x, 70) != img.RGBAAt(150, 10) {
				return x
			}
		}
		return -1
	}
	tt.ClickAt(274, 70) // the arrow, with nothing typed
	moved := 0
	for range 20 {
		step(tt, 16*time.Millisecond)
		moved = max(moved, left())
	}
	if moved < 3 || !tt.HasText("Enter your email first") {
		t.Fatalf("an empty try shakes the field (moved %dpx) and says so, shown %q", moved, tt.Texts())
	}
	step(tt, time.Second)
	if l := left(); l != 0 {
		t.Fatalf("the shake dies out back in place, left edge at %d", l)
	}

	tt.ClickAt(100, 70)
	tt.Type("sam@studio")
	step(tt, 300*time.Millisecond)
	if tt.HasText("Enter your email first") {
		t.Fatal("typing puts the error away")
	}
	tt.Key(0, ui.KeyEnter)
	step(tt, 300*time.Millisecond)
	if !tt.HasText("That domain looks unfinished") || sent != "" {
		t.Fatalf("an unfinished domain is caught, shown %q", tt.Texts())
	}
	tt.Type(".com")
	tt.ClickAt(274, 70)
	step(tt, 500*time.Millisecond)
	if sent != "sam@studio.com" || !tt.HasText("You're on the list") {
		t.Fatalf("a right address is sent, sent %q, shown %q", sent, tt.Texts())
	}
}
