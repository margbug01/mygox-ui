package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Resting on someone keeps their path lit and fades the others; Near folds
// the outer ring away and All grows it back.
func TestBlastRadiusLightsPathAndFoldsRing(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.BlastRadius(c)
	}, 300, 330)
	step(tt, 50*time.Millisecond)
	pink := func() bool { // the friend straight below you
		p := tt.Image().RGBAAt(150, 228)
		return p.R > 200 && p.G < 150
	}
	if !pink() || !tt.HasText("21") {
		t.Fatal("rests with everyone in colour, 21 people")
	}
	tt.Move(197, 201) // the blue friend at four o'clock
	step(tt, 400*time.Millisecond)
	if blue := tt.Image().RGBAAt(197, 201); pink() || blue.B < 200 || blue.R > 140 {
		t.Fatalf("the blue friend stays lit (%v) while the others fade", blue)
	}
	tt.Move(5, 325)
	tt.Click("Near")
	step(tt, time.Second)
	outer := tt.Image().RGBAAt(150, 74) // a contact on the outer ring
	if !tt.HasText("7") || outer.R < 240 || outer.B < 240 {
		t.Fatalf("Near folds the outer ring away (%v), 7 people", outer)
	}
	tt.Click("All")
	step(tt, time.Second)
	if p := tt.Image().RGBAAt(150, 74); p.R > 200 || !tt.HasText("21") {
		t.Fatalf("All grows it back (%v)", p)
	}
}
