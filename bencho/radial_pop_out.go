package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// RadialPopOutOptions configures RadialPopOut.
type RadialPopOutOptions struct {
	Spread float32 // how far the ring opens, 60 to 105 (%); 100 by default
	Speed  float32 // 0 to 100; 50 by default
	Corner float32 // 0 to 30; 20 by default
}

const (
	radialPopOutW     = 470
	radialPopOutH     = 491
	radialPopOutCardW = 186
	radialPopOutCardH = 232
	radialPopOutRing  = 185 // the ring's radius at full spread
	radialPopOutSmall = 0.46
)

// radialPopOutRest is where each picture behind the front one sits in the
// stack: its offset from the middle and its turn.
var radialPopOutRest = [5]struct{ dx, dy, rot float32 }{
	{0, -9, -1.5}, {-25, 6, -8}, {24, 3, 8}, {-4, 7, -4}, {6, -3, 3},
}

// RadialPopOut is a picture with others stacked behind it. Pointing at it
// bursts the hidden pictures out around it in a ring; they fold back into
// the stack when the pointer leaves.
func RadialPopOut(c *ui.Context, opts ...RadialPopOutOptions) ui.Element {
	var o RadialPopOutOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Spread == 0 {
		o.Spread = 100
	}
	if o.Speed == 0 {
		o.Speed = 50
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Spread = clamp32(o.Spread, 60, 105)
	o.Speed = clamp32(o.Speed, 0, 100)
	o.Corner = clamp32(o.Corner, 0, 30)
	art := genieArtwork()
	pics := []int{3, 0, 1, 2, 4, 5} // behind, then the front one last

	stage := ui.Box(c).Size(radialPopOutW, radialPopOutH)
	st := ui.Local(stage, "radial_pop_out", func() bool { return false })
	cx, cy := float32(radialPopOutW/2), float32(radialPopOutH/2)
	px, py, over := stage.PointerPosition()
	ring := radialPopOutRing * o.Spread / 100
	// It opens over the front picture and stays open while the pointer is within the ring.
	inFront := float32(math.Abs(float64(px-cx))) < radialPopOutCardW/2 && float32(math.Abs(float64(py-cy))) < radialPopOutCardH/2
	near := float32(math.Hypot(float64(px-cx), float64(py-cy))) < ring+70
	*st = over && (inFront || (*st && near))

	dur := time.Duration(1000-7*o.Speed) * time.Millisecond
	stage.Children(func() {
		for k, pic := range pics {
			card := ui.Box(c).Key(pic).Absolute().Size(radialPopOutCardW, radialPopOutCardH)
			x, y, rot, scale := cx, cy, float32(0), float32(1)
			if k < len(radialPopOutRest) {
				// Each one leaves a beat after the one before it.
				out := card.Spring("out", b2f(*st), ui.SpringOptions{Duration: dur + time.Duration(k)*40*time.Millisecond, Bounce: 0.3})
				rest := radialPopOutRest[k]
				ang := (-90 + float64(k)*72) * math.Pi / 180
				sin, cos := math.Sincos(ang)
				x = cx + lerp32(rest.dx, ring*float32(cos), out)
				y = cy + lerp32(rest.dy, ring*float32(sin), out)
				rot = lerp32(rest.rot, -rest.rot*0.8, out)
				scale = lerp32(1, radialPopOutSmall, clamp32(out, 0, 1.2))
			} else {
				lift := card.Spring("lift", b2f(*st), ui.SpringOptions{Duration: dur, Bounce: 0.3})
				scale = 1 - 0.04*lift
			}
			radius := o.Corner / max(scale, 0.2)
			card.Left(x-radialPopOutCardW/2).Top(y-radialPopOutCardH/2).Rotate(rot).Scale(scale, scale).Radius(radius).
				Shadow(0, 4, 14, 0, ui.RGBA(0, 0, 0, 0.10)).
				Draw(func(p *ui.Painter, r ui.Rect) {
					w := float32(genieArtPx) * radialPopOutCardW / radialPopOutCardH
					src := ui.Rect{X: (genieArtPx - w) / 2, W: w, H: genieArtPx}
					p.Clip(r, radius, func() { p.ImageRegion(art[pic], src, r, ui.ImageOptions{}) })
				})
		}
	})
	return stage
}

func init() {
	registerDemo(
		"radial-pop-out",
		"A picture with others stacked behind it: pointing at it bursts the hidden ones out around it in a ring.",
		func(c *ui.Context) { RadialPopOut(c) },
	)
}
