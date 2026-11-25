package bencho_test

import (
	"image/color"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing the palette morphs a new set in, the first colour changing before
// the last.
func TestPaletteMorphsInNewSetOneColourAfterAnother(t *testing.T) {
	var got []ui.Color
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Palette(c, bencho.PaletteOptions{OnChange: func(cs []ui.Color) { got = cs }})
	}
	tt := ui.NewTester(view, 300, 78)
	tt.Frame()
	swatch := func(i int) color.RGBA { return tt.Image().RGBAAt(36+57*i, 39) }
	first, last := swatch(0), swatch(3)

	tt.ClickAt(120, 39)
	step(tt, 240*time.Millisecond)
	if len(got) != 4 {
		t.Fatalf("pressing should pick a new set of 4, got %v", got)
	}
	if swatch(0) == first {
		t.Fatal("the first colour should already be changing")
	}
	if swatch(3) != last {
		t.Fatal("the last colour should not have started changing yet")
	}
	step(tt, 2*time.Second)
	for i, want := range got {
		if px := swatch(i); px != (color.RGBA{want.R, want.G, want.B, 255}) {
			t.Fatalf("swatch %d is %v, want the new set's %v", i, px, want)
		}
	}
}
