package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// TiltCardOptions configures TiltCard.
type TiltCardOptions struct {
	Tilt   float32 // how far it turns under you, 0 to 20 degrees; 10 by default
	Shade  float32 // how deep the dent reads, 0 to 100; 60 by default
	Corner float32 // the card's radius, 0 to 40; 20 by default
}

const (
	tiltCardW = 260
	tiltCardH = 320
)

// TiltCard is a photo card that sinks under the pointer instead of rising
// toward it: the side you point at turns away from you and a soft dent of
// shade follows the pointer.
func TiltCard(c *ui.Context, opts ...TiltCardOptions) ui.Element {
	var o TiltCardOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Tilt == 0 {
		o.Tilt = 10
	}
	if o.Shade == 0 {
		o.Shade = 60
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Tilt = clamp32(o.Tilt, 0, 20)
	o.Shade = clamp32(o.Shade, 0, 100)
	o.Corner = clamp32(o.Corner, 0, 40)
	art := genieArtwork()[0]

	// The pointer is read on a still frame round the card, so the turn cannot feed back into it.
	frame := ui.Box(c).Size(tiltCardW, tiltCardH)
	px, py, over := frame.PointerPosition()
	nx, ny := float32(0), float32(0)
	if over {
		nx, ny = clamp32(px/tiltCardW*2-1, -1, 1), clamp32(py/tiltCardH*2-1, -1, 1)
	}
	spring := ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.15}
	sx := frame.Spring("x", nx, spring)
	sy := frame.Spring("y", ny, spring)
	in := frame.Spring("in", b2f(over), spring)
	dentX := frame.Spring("dx", px, ui.SpringOptions{Duration: 160 * time.Millisecond})
	dentY := frame.Spring("dy", py, ui.SpringOptions{Duration: 160 * time.Millisecond})

	frame.Children(func() {
		s := 1 - 0.025*in
		// The side under the pointer goes away from you.
		ui.Box(c).Size(tiltCardW, tiltCardH).Radius(o.Corner).
			RotateY(o.Tilt*sx).RotateX(-o.Tilt*sy).Perspective(700).Scale(s, s).
			Shadow(0, 10-4*in, 26-8*in, 0, ui.RGBA(0, 0, 0, 0.13)).
			Draw(func(p *ui.Painter, r ui.Rect) {
				p.Clip(r, o.Corner, func() {
					w := float32(genieArtPx) * tiltCardW / tiltCardH
					p.ImageRegion(art, ui.Rect{X: (genieArtPx - w) / 2, W: w, H: genieArtPx}, r, ui.ImageOptions{})
					if in > 0.01 {
						const rad = 150
						k := o.Shade / 100 * 0.42 * clamp32(in, 0, 1)
						p.FillRadialGradient(ui.Rect{X: r.X + dentX - rad, Y: r.Y + dentY - rad, W: 2 * rad, H: 2 * rad}, ui.RadialGradient{Radius: rad, Stops: []ui.GradientStop{
							{At: 0, Color: ui.RGBA(40, 40, 50, k)}, {At: 0.45, Color: ui.RGBA(40, 40, 50, k*0.5)}, {At: 1, Color: ui.RGBA(40, 40, 50, 0)},
						}}, 0)
					}
				})
			})
	})
	return frame
}

func init() {
	registerDemo(
		"tilt",
		"A photo card that sinks under the pointer instead of rising toward it: the side you point at turns away and a soft dent follows.",
		func(c *ui.Context) { TiltCard(c) },
	)
}
