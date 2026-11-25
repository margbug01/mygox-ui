package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging the drum rolls the rows through the slot and it settles on the
// nearest one; scrolling turns it a row a notch; Currencies picks money.
func TestCountryWheelRollsIntoTheSlot(t *testing.T) {
	list := bencho.CountryWheelCountries
	chosen := "France"
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.CountryWheel(c, bencho.CountryWheelOptions{List: list, OnChange: func(n string) { chosen = n }})
	}, 290, 290)
	step(tt, 50*time.Millisecond)
	tt.Press(180, 200)
	step(tt, 50*time.Millisecond)
	tt.Move(180, 150)
	step(tt, 50*time.Millisecond)
	tt.Move(180, 88) // a little over two rows up
	tt.Release(180, 88)
	step(tt, 600*time.Millisecond)
	if chosen != "Italy" {
		t.Fatalf("dragging two rows up settles Italy in the slot, got %q", chosen)
	}
	tt.Scroll(150, 145, 0, -40)
	step(tt, 400*time.Millisecond)
	if chosen != "Germany" {
		t.Fatalf("scrolling up turns back a row, got %q", chosen)
	}
	list = bencho.CountryWheelCurrencies
	step(tt, 50*time.Millisecond)
	if !tt.HasText("Euro") || !tt.HasText("Pound sterling") {
		t.Fatalf("Currencies puts money on the wheel, got %v", tt.Texts())
	}
}
