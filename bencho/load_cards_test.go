package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The cards develop in order, the first before the last, and settle on their content.
func TestLoadCardsDevelopInOrder(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(320, 290).Background(th.SurfaceSubtle).Children(func() {
			bencho.LoadCards(c, bencho.LoadCardsDefault())
		})
	}, 320, 290)
	step(tt, 600*time.Millisecond)
	early := tt.Image()
	if early.RGBAAt(160, 44) == early.RGBAAt(160, 236) {
		t.Fatal("the first card has arrived while the last has not")
	}
	step(tt, 3*time.Second)
	if !tt.HasText("Layout animation with FLIP") || !tt.HasText("motion") {
		t.Fatal("every card ends on its content")
	}
}
