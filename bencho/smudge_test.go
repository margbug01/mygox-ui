package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Every effect disturbs the picture under a stroke, and it settles back.
func TestSmudgeStrokesAndSettles(t *testing.T) {
	for _, e := range []bencho.SmudgeEffect{bencho.SmudgeWet, bencho.SmudgeRipple, bencho.SmudgeMosaic} {
		tt := ui.NewTester(func(c *ui.Context) {
			th := bencho.ApplyTheme(c, false)
			ui.Box(c).Size(300, 360).Background(th.SurfaceSubtle).Children(func() { bencho.Smudge(c, bencho.SmudgeOptions{Effect: e, Wobble: 1}) })
		}, 300, 360)
		tt.Move(-5, -5)
		step(tt, 50*time.Millisecond)
		still := tt.Image()
		for x := float32(60); x <= 240; x += 12 {
			tt.Move(x, 180)
			step(tt, 16*time.Millisecond)
		}
		if changed(still, tt.Image(), 360) == 0 {
			t.Fatalf("effect %d: a stroke disturbs the picture", e)
		}
		tt.Move(-5, -5)
		step(tt, 4*time.Second)
		if changed(still, tt.Image(), 360) > 3000 {
			t.Fatalf("effect %d: it settles back", e)
		}
	}
}
