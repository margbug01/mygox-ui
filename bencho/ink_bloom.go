package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// InkBloomOptions configures InkBloom.
type InkBloomOptions struct {
	Edge  float32 // how ragged the spreading edge is, 0 to 100; 50 by default
	Speed float32 // how fast a bloom fills the frame, 0 to 100; 50 by default
}

type inkBloomState struct {
	from, to int
	at       time.Time
	x, y     float32 // where the bloom started, in the photo
	seed     int
}

const (
	inkBloomW   = 272
	inkBloomH   = 352
	inkBloomPad = 8
)

// inkBloomCrop is the part of a square picture that fills the photo's portrait shape.
func inkBloomCrop() ui.Rect {
	pw, ph := float32(inkBloomW-2*inkBloomPad), float32(inkBloomH-2*inkBloomPad)
	w := float32(genieArtPx) * pw / ph
	return ui.Rect{X: (genieArtPx - w) / 2, W: w, H: genieArtPx}
}

// inkBloomBlob is a ragged round edge of radius r about (cx, cy).
func inkBloomBlob(cx, cy, r, rag float32, seed int, phase float64) *ui.Path {
	var path ui.Path
	const n = 72
	for i := 0; i <= n; i++ {
		a := float64(i) / n * 2 * math.Pi
		wob := genieNoise(math.Cos(a)*2.2+5+phase, math.Sin(a)*2.2+5, seed)*2 - 1
		rr := r * (1 + rag*0.2*float32(wob))
		x, y := cx+rr*float32(math.Cos(a)), cy+rr*float32(math.Sin(a))
		if i == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
	path.Close()
	return &path
}

// InkBloom is an image slider in a white frame. Pressing the photo blooms the
// next picture out from that point like ink spreading through water, with a
// ragged edge and a thin lens bending the picture behind; the round arrows
// bloom from their side, back or forward.
func InkBloom(c *ui.Context, opts ...InkBloomOptions) ui.Element {
	var o InkBloomOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Edge == 0 {
		o.Edge = 50
	}
	if o.Speed == 0 {
		o.Speed = 50
	}
	art := genieArtwork()
	t := ThemeFrom(c)
	frame := ui.Box(c).Size(inkBloomW, inkBloomH).Radius(20).Background(t.Surface)
	st := ui.Local(frame, "ink_bloom", func() inkBloomState { return inkBloomState{} })
	dur := time.Duration(lerp32(1500, 500, clamp32(o.Speed, 0, 100)/100)) * time.Millisecond
	k := float32(1)
	if !st.at.IsZero() {
		k = clamp32(float32(c.Now().Sub(st.at))/float32(dur), 0, 1)
		if k < 1 {
			c.AnimationFrame()
		} else {
			st.from, st.at = st.to, time.Time{}
		}
	}
	bloom := func(x, y float32, step int) {
		if !st.at.IsZero() {
			st.from = st.to
		}
		st.to = (st.from + step + len(art)) % len(art)
		st.x, st.y, st.at = x, y, c.Now()
		st.seed++
	}
	pw, ph := float32(inkBloomW-2*inkBloomPad), float32(inkBloomH-2*inkBloomPad)
	rag := clamp32(o.Edge, 0, 100) / 100

	frame.Children(func() {
		photo := ui.Box(c).Key("photo").Absolute().Left(inkBloomPad).Top(inkBloomPad).Size(pw, ph).Cursor(ui.CursorPointer)
		photo.OnClick(func() {
			x, y, _ := photo.PointerPosition()
			bloom(x, y, 1)
		})
		photo.Draw(func(p *ui.Painter, r ui.Rect) {
			crop := inkBloomCrop()
			p.Clip(r, 14, func() {
				p.ImageRegion(art[st.from], crop, r, ui.ImageOptions{})
				if st.at.IsZero() {
					return
				}
				cx, cy := r.X+st.x, r.Y+st.y
				far := float32(0)
				for _, q := range [][2]float32{{r.X, r.Y}, {r.X + r.W, r.Y}, {r.X, r.Y + r.H}, {r.X + r.W, r.Y + r.H}} {
					far = max(far, float32(math.Hypot(float64(q[0]-cx), float64(q[1]-cy))))
				}
				e := 1 - (1-k)*(1-k)*(1-k)
				rad := far * 1.3 * e
				phase := float64(k) * 1.5
				// The lens: just outside the edge the old picture swells.
				p.ClipPath(inkBloomBlob(cx, cy, rad+7, rag, st.seed, phase), func() {
					p.Transform(ui.Scaling(1.07, 1.07), cx, cy, func() {
						p.ImageRegion(art[st.from], crop, r, ui.ImageOptions{})
					})
				})
				edge := inkBloomBlob(cx, cy, rad, rag, st.seed, phase)
				p.ClipPath(edge, func() {
					p.ImageRegion(art[st.to], crop, r, ui.ImageOptions{})
				})
				p.StrokePath(edge, 1.2, ui.RGBA(255, 255, 255, 0.35*(1-k)))
			})
		})
		arrow := func(key string, icon *ui.SVG, x float32, step int) {
			b := ui.Box(c).Key(key).Absolute().Left(x-17).Top(inkBloomH/2-17).Size(34, 34).Radius(17).
				Background(ui.Hex("#ffffff")).Center().Cursor(ui.CursorPointer).
				OnClick(func() { bloom(x-inkBloomPad, inkBloomH/2-inkBloomPad, step) })
			press := b.Spring("press", b2f(b.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
			b.Scale(1-0.08*press, 1-0.08*press).Children(func() {
				ui.Icon(c, icon).Size(16, 16).TextColor(ui.Hex("#111111"))
			})
		}
		arrow("prev", icons.ChevronLeft(), 37, -1)
		arrow("next", icons.ChevronRight(), inkBloomW-37, 1)
	})
	return frame
}

func init() {
	registerDemo(
		"ink-bloom",
		"An image slider where the next picture blooms out from the press like ink in water.",
		func(c *ui.Context) { InkBloom(c) },
	)
}
