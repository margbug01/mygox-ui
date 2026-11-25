package bencho_test

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Refresh moves the price: the new digits rise in one after another, and a
// fall shows as a red badge.
func TestPriceFlowRefreshRisesNewDigitsAndColoursChange(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.PriceFlow(c)
	}
	tt := ui.NewTester(view, 300, 230)
	tt.Frame()
	ink := func() int { // dark pixels in the price
		img := tt.Image()
		n := 0
		for y := 58; y < 102; y++ {
			for x := 0; x < 175; x++ {
				if img.RGBAAt(x, y).G < 90 {
					n++
				}
			}
		}
		return n
	}
	settled := ink()
	if err := tt.Click("Refresh"); err != nil {
		t.Fatal(err)
	}
	step(tt, 64*time.Millisecond)
	if early := ink(); early > settled*9/10 {
		t.Fatalf("just after Refresh the new digits should still be rising in: %d dark pixels, %d at rest", early, settled)
	}
	step(tt, time.Second)
	if got := strings.Join(tt.Texts(), ""); !strings.Contains(got, "$118.47") || !strings.Contains(got, "4.64%") {
		t.Fatalf("after Refresh want $118.47 down 4.64%%, got %q", got)
	}
	r, ok := tt.Find("4.64%")
	if !ok {
		t.Fatal("no change badge")
	}
	if px := tt.Image().RGBAAt(int(r.X)-4, int(r.Y+r.H/2)); px.R < 180 || px.G > 120 {
		t.Fatalf("a fall should turn the badge red, got %v", px)
	}
}
