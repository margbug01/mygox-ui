package bencho

import (
	"image"
	"image/color"
	"math"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
)

// ImageCompareOptions configures ImageCompare.
type ImageCompareOptions struct {
	NoLabels bool    // hides the Before and After tags
	Corner   float32 // 0 to 28; 20 by default
}

type imageCompareState struct {
	at       float32 // where the line rests, 0 (left) to 1 (right)
	held     bool
	flick    float32 // the drag's last speed, px/s
	lastX    float32 // the line's drawn x last frame, for its speed
	lastT    time.Time
	bend     float32 // how far the line's middle leads its ends, px
	initDone bool
}

const (
	imageCompareW     = 270
	imageCompareH     = 360
	imageCompareKnob  = 36
	imageCompareFlick = 700 // px/s past which a let-go glides to the edge
)

var (
	imageCompareOnce sync.Once
	imageCompareGray *ui.Bitmap
)

// imageCompareBefore is the "before" picture: genie's ridge with its
// colour taken out.
func imageCompareBefore() *ui.Bitmap {
	imageCompareOnce.Do(func() {
		genieArtwork()
		src := geniePixels[1]
		gray := image.NewRGBA(src.Bounds())
		for i := 0; i < len(src.Pix); i += 4 {
			l := uint8(0.3*float64(src.Pix[i]) + 0.59*float64(src.Pix[i+1]) + 0.11*float64(src.Pix[i+2]))
			l = uint8(float64(l)*0.62 + 22) // flatter and darker, like an old print
			gray.SetRGBA(i/4%src.Bounds().Dx(), i/4/src.Bounds().Dx(), color.RGBA{l, l, l, 255})
		}
		imageCompareGray = ui.NewBitmap(gray)
	})
	return imageCompareGray
}

// ImageCompare is a before and after slider over a picture: dragging the
// line uncovers one side or the other, a flick sends it gliding to the
// edge, and a click moves it there.
func ImageCompare(c *ui.Context, opts ...ImageCompareOptions) ui.Element {
	var o ImageCompareOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	corner := clamp32(o.Corner, 0, 28)
	after := genieArtwork()[1]
	before := imageCompareBefore()

	card := ui.Box(c).Size(imageCompareW, imageCompareH).Cursor(ui.CursorResizeEW)
	st := ui.Local(card, "image-compare", func() imageCompareState { return imageCompareState{at: 0.5} })
	px, _, over := card.PointerPosition()
	_, _, held := card.Dragged()
	if held {
		st.at = clamp32(px/imageCompareW, 0, 1)
		if vx, _ := card.DragVelocity(); vx != 0 {
			st.flick = vx
		}
	} else if st.held {
		if st.flick > imageCompareFlick {
			st.at = 1
		} else if st.flick < -imageCompareFlick {
			st.at = 0
		}
		st.flick = 0
	}
	st.held = held

	dur := 420 * time.Millisecond
	if held {
		dur = 90 * time.Millisecond
	}
	x := card.Spring("x", st.at*imageCompareW, ui.SpringOptions{Duration: dur, Bounce: 0.12})
	now := c.Now()
	if !st.initDone {
		st.lastX, st.lastT, st.initDone = x, now, true
	}
	if dt := float32(now.Sub(st.lastT).Seconds()); dt > 0 {
		v := (x - st.lastX) / dt
		st.bend += (clamp32(v*0.012, -16, 16) - st.bend) * min(1, dt*14)
		st.lastX, st.lastT = x, now
	}
	if math.Abs(float64(st.bend)) > 0.2 {
		c.AnimationFrame()
	}
	knobScale := 1.0 + 0.06*b2f(over && !held) - 0.24*b2f(held)
	knob := card.Spring("knob", knobScale, ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.4})

	card.Draw(func(p *ui.Painter, r ui.Rect) {
		// The pictures fill the card from a 144x192 window of the 192px art.
		src := ui.Rect{X: 24, Y: 0, W: 144, H: 192}
		lx, bend := r.X+x, st.bend
		p.Clip(r, corner, func() {
			p.ImageRegion(after, src, r, ui.ImageOptions{})
			var left, line ui.Path
			left.MoveTo(r.X, r.Y).LineTo(lx, r.Y).QuadTo(lx+2*bend, r.Y+r.H/2, lx, r.Y+r.H).LineTo(r.X, r.Y+r.H).Close()
			p.ClipPath(&left, func() { p.ImageRegion(before, src, r, ui.ImageOptions{}) })
			line.MoveTo(lx, r.Y).QuadTo(lx+2*bend, r.Y+r.H/2, lx, r.Y+r.H)
			p.StrokePath(&line, 2, ui.Hex("#ffffff"))
		})
		// The knob rides the line's middle.
		kx, ky := lx+bend, r.Y+r.H/2
		d := imageCompareKnob * knob
		k := ui.Rect{X: kx - d/2, Y: ky - d/2, W: d, H: d}
		p.Shadow(k, d/2, 0, 2, 8, 0, ui.RGBA(0, 0, 0, 0.18))
		p.Fill(k, ui.Hex("#ffffff"), d/2)
		s := knob
		ink := ui.Hex("#111111")
		p.StrokePolyline([][2]float32{{kx - 3*s, ky - 5*s}, {kx - 8*s, ky}, {kx - 3*s, ky + 5*s}}, 2, ink)
		p.StrokePolyline([][2]float32{{kx + 3*s, ky - 5*s}, {kx + 8*s, ky}, {kx + 3*s, ky + 5*s}}, 2, ink)
	})
	if !o.NoLabels {
		card.Children(func() {
			tag := func(label string, left, w float32) {
				ui.Box(c).Key(label).Absolute().Left(left).Top(12).Size(w, 28).Radius(14).
					Background(ui.RGBA(0, 0, 0, 0.32)).Center().Children(func() {
					ui.Text(c, label).FontSize(13).FontWeight(400).TextColor(ui.Hex("#ffffff")).SingleLine()
				})
			}
			tag("Before", 12, 62)
			tag("After", imageCompareW-12-52, 52)
		})
	}
	return card
}

func init() {
	registerDemo(
		"image-compare",
		"Drag the line to compare two pictures, or flick it and it glides to the edge.",
		func(c *ui.Context) { ImageCompare(c) },
	)
}
