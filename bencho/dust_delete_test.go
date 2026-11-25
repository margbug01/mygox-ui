package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Binning the photo blows it away as dust, and the next photo takes its place.
func TestDustDeleteBinsToTheNextPhoto(t *testing.T) {
	for _, effect := range []bencho.DustDeleteEffect{bencho.DustDeletePour, bencho.DustDeleteAsh, bencho.DustDeleteWind, bencho.DustDeletePieces} {
		tt := ui.NewTester(func(c *ui.Context) {
			bencho.ApplyTheme(c, false)
			bencho.DustDelete(c, bencho.DustDeleteOptions{Effect: effect, Speed: 100, Grain: 6})
		}, 260, 345)
		step(tt, 50*time.Millisecond)
		first := tt.Image()
		tt.ClickAt(228, 313)
		step(tt, 300*time.Millisecond)
		mid := tt.Image()
		step(tt, 1500*time.Millisecond)
		last := tt.Image()
		if changed(first, mid, 290) == 0 || changed(first, last, 290) == 0 || changed(mid, last, 290) == 0 {
			t.Fatalf("effect %d: the photo turns to dust and the next one grows in", effect)
		}
	}
}
