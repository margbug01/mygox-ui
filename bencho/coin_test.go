package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A fling sets the coin spinning: it narrows as it turns edge on, keeps
// turning after letting go, and slows to rest face on.
func TestCoinFlingSpinsAndSettles(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 300).Background(th.SurfaceSubtle).Children(func() { bencho.Coin(c, bencho.CoinOptions{Still: true}) })
	}, 300, 300)
	step(tt, 100*time.Millisecond)
	width := func() int { // the coin across its middle
		img := tt.Image()
		bg, n := img.RGBAAt(2, 150), 0
		for x := 0; x < 300; x++ {
			if img.RGBAAt(x, 150) != bg {
				n++
			}
		}
		return n
	}
	rest := width()
	if rest < 170 || rest > 200 {
		t.Fatalf("the coin is about 183px across at rest, have %d", rest)
	}
	tt.Press(150, 150)
	for x := float32(150); x <= 270; x += 30 {
		tt.Move(x, 150)
		step(tt, 16*time.Millisecond)
	}
	tt.Release(270, 150)
	narrowest := rest
	for i := 0; i < 40; i++ {
		step(tt, 16*time.Millisecond)
		narrowest = min(narrowest, width())
	}
	if narrowest > rest/2 {
		t.Fatalf("a fling spins the coin edge on, narrowest %d of %d", narrowest, rest)
	}
	step(tt, 5*time.Second)
	if w := width(); w < rest-4 || w > rest+4 {
		t.Fatalf("the coin comes to rest face on, %d across, want %d", w, rest)
	}
}
