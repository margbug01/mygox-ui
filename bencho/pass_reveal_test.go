package bencho_test

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing the eye turns the dots into letters from the first one on;
// pressing again turns them back into dots from the last one.
func TestPasswordRevealUnmasksInOrderAndBack(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.PasswordReveal(c, "sunsetRain7")
	}
	tt := ui.NewTester(view, 300, 52)
	tt.Frame()
	// The first and the last character cells of the field.
	cell := func(img *image.RGBA, x0 int) []uint8 {
		var px []uint8
		for y := 14; y < 38; y++ {
			for x := x0; x < x0+10; x++ {
				px = append(px, img.RGBAAt(x, y).G)
			}
		}
		return px
	}
	same := func(a, b []uint8) bool { return string(a) == string(b) }
	dots := tt.Image()
	first := cell(dots, 23)

	if err := tt.Click("Show password"); err != nil {
		t.Fatal(err)
	}
	step(tt, 350*time.Millisecond)
	mid := tt.Image()
	if same(cell(mid, 23), first) {
		t.Fatal("350ms into the reveal the first dot should be a letter")
	}
	step(tt, 900*time.Millisecond)
	shown := tt.Image()
	if same(mid.Pix, shown.Pix) {
		t.Fatal("350ms into the reveal the last dots should not have turned yet")
	}
	if !tt.HasText("Hide password") {
		t.Fatalf("revealed, the eye should offer to hide, got %v", tt.Texts())
	}

	if err := tt.Click("Hide password"); err != nil {
		t.Fatal(err)
	}
	step(tt, 350*time.Millisecond)
	if back := tt.Image(); same(cell(back, 23), first) || same(back.Pix, shown.Pix) {
		t.Fatal("350ms into hiding, the last letters should be dots again and the first still a letter")
	}
	step(tt, 900*time.Millisecond)
	if !same(cell(tt.Image(), 23), first) {
		t.Fatal("hidden again, the first character should be a dot")
	}
}
