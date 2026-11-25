package bencho_test

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

func TestPianoPlaysPressedGlidedAndTypedKeys(t *testing.T) {
	var notes []string
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Piano(c, bencho.PianoOptions{OnNote: func(n string) { notes = append(notes, n) }})
	}
	tt := ui.NewTester(view, 304, 142)
	tt.Frame()
	rest := tt.Image().RGBAAt(98, 120)
	tt.Move(98, 120) // E
	tt.Press(98, 120)
	step(tt, 200*time.Millisecond)
	if sunk := tt.Image().RGBAAt(98, 120); sunk == rest {
		t.Fatalf("a pressed key should sink and darken, still %v", sunk)
	}
	tt.Move(170, 120) // glide on to G
	step(tt, 50*time.Millisecond)
	tt.Release(170, 120)
	step(tt, 300*time.Millisecond)
	if back := tt.Image().RGBAAt(98, 120); back != rest {
		t.Fatalf("a let-go key should come back up, got %v want %v", back, rest)
	}
	// While pointing at the piano, W plays C#.
	tt.KeyDown(ui.KeyW)
	tt.Frame()
	tt.KeyUp(ui.KeyW)
	tt.Frame()
	if want := []string{"E4", "G4", "C#4"}; !slices.Equal(notes, want) {
		t.Fatalf("want notes %v, got %v", want, notes)
	}
}
