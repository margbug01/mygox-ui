package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing the hand deals the top card round to the back, changing what is on top.
func TestCardFanDealsTheTopCardBack(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.CardFan(c)
	}, 510, 353)
	step(tt, 700*time.Millisecond)
	before := tt.Image()
	tt.ClickAt(380, 200)
	tt.Move(5, 5) // point away, so no card is lifted
	step(tt, 1200*time.Millisecond)
	if before.RGBAAt(390, 220) == tt.Image().RGBAAt(390, 220) {
		t.Fatal("a new card sits on top of the hand")
	}
}
