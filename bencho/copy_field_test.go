package bencho_test

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The copy button puts the code on the clipboard and shows a check for a
// moment before the copy icon comes back.
func TestCopyFieldCopiesCodeAndShowsCheck(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.CopyField(c)
	}
	tt := ui.NewTester(view, 296, 296)
	tt.Frame()
	tt.Move(5, 290) // away from the button, so hover does not tint it
	step(tt, 300*time.Millisecond)
	icon := func() *[256]uint8 { // the button's 16x16 glyph, as grey levels
		var g [256]uint8
		img := tt.Image()
		for y := range 16 {
			for x := range 16 {
				g[y*16+x] = img.RGBAAt(264+x, 16+y).G
			}
		}
		return &g
	}
	before := icon()

	tt.ClickAt(272, 24)
	step(tt, 500*time.Millisecond)
	if got := tt.Clipboard(); !strings.Contains(got, "<Confetti") || !strings.HasPrefix(got, "import") {
		t.Fatalf("clipboard holds %q, want the code", got)
	}
	tt.Move(5, 290)
	step(tt, 300*time.Millisecond)
	if *icon() == *before {
		t.Fatal("after copying the button should show a check, not the copy icon")
	}
	step(tt, 1500*time.Millisecond)
	if *icon() != *before {
		t.Fatal("the copy icon should come back after a moment")
	}
}
