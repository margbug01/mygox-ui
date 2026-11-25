package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pointing at a target grows the sales box inside it and the header follows;
// Year switches to yearly figures.
func TestPillarsGrowsTheHoveredBox(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Pillars(c)
	}, 360, 364)
	step(tt, time.Second)
	blue := func(x, y int) bool {
		p := tt.Image().RGBAAt(x, y)
		return p.B > 200 && p.R < 120
	}
	if !tt.HasText("$2,536.99") || !blue(322, 250) || blue(147, 250) {
		t.Fatal("rests on June, its sales box filled")
	}
	tt.Move(147, 220) // March's box
	step(tt, time.Second)
	if !tt.HasText("$2,950.75") || !tt.HasText("$3,400") || !blue(147, 250) || blue(322, 250) {
		t.Fatal("pointing at March grows its box and empties June's")
	}
	tt.Click("Year")
	tt.Move(5, 5)
	step(tt, time.Second)
	if !tt.HasText("2026") || !tt.HasText("$30,441.90") {
		t.Fatal("Year shows yearly sales")
	}
}
