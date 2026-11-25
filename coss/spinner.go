package coss

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

var spinnerIconSVG = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>`))

// SpinnerElement is a coss-styled spinning loader icon (LoaderCircle).
type SpinnerElement struct {
	ui.Element
	ink *ui.Color // a drawn Kind's ink, which Color sets
}

// SpinnerKind is what a Spinner draws.
type SpinnerKind uint8

const (
	SpinnerCircle SpinnerKind = iota // coss's LoaderCircle, turning
	// SpinnerTrace is a dash running round its own faint circle, its
	// length breathing as it goes.
	SpinnerTrace
	// SpinnerTraceSquare runs the dash round a rounded square.
	SpinnerTraceSquare
	// SpinnerMorph is a square that rounds into a circle as it turns, and
	// back again.
	SpinnerMorph
)

// SpinnerOptions are what a Spinner draws, and how fast.
type SpinnerOptions struct {
	Kind  SpinnerKind
	Speed float32 // of the drawn Kinds; 1 when 0
	// Ink and Track color the drawn Kinds, the dash or the shape and the
	// trace's track; the theme's foreground and a faint one when unset.
	Ink, Track ui.Color
}

// Spinner renders a smoothly rotating 16x16 loading spinner that inherits
// text color, or with a Kind, draws that.
func Spinner(c *ui.Context, opts ...SpinnerOptions) *SpinnerElement {
	var o SpinnerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Kind == SpinnerCircle {
		el := ui.Icon(c, spinnerIconSVG).Size(16, 16).Shrink(0)
		angle := el.Loop("spin", 800*time.Millisecond, ui.Linear) * 360
		el.Rotate(angle)
		return &SpinnerElement{Element: el}
	}
	t := theme.Of(c)
	if o.Speed == 0 {
		o.Speed = 1
	}
	if o.Ink.A == 0 {
		o.Ink = t.Foreground
	}
	if o.Track.A == 0 {
		o.Track = theme.Alpha(t.Foreground, 0.12)
	}
	// The drawn kinds run on the wall clock, so that those side by side
	// keep in step.
	c.AnimationFrame()
	now := float64(c.Now().UnixNano()) / 1e9 * float64(o.Speed)
	s := &SpinnerElement{ink: &o.Ink}
	s.Element = ui.Box(c).Size(16, 16).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
		if o.Kind == SpinnerMorph {
			drawMorphSpinner(p, r, now, *s.ink)
		} else {
			drawTraceSpinner(p, r, now, o.Kind == SpinnerTrace, *s.ink, o.Track)
		}
	})
	return s
}

// drawTraceSpinner draws in r, at time now, a dash running round its
// track, a circle or a rounded square; its sizes are for r 28 wide.
func drawTraceSpinner(p *ui.Painter, r ui.Rect, now float64, circle bool, ink, track ui.Color) {
	f := r.W / 28
	w := 2.5 * f
	s := float64(r.W - w)
	rad := 8.0 * float64(f)
	if circle {
		rad = s / 2
	}
	p.Stroke(ui.Rect{X: r.X + w/2, Y: r.Y + w/2, W: float32(s), H: float32(s)}, track, float32(rad), w)
	head := now * 0.7
	length := 0.12 + 0.1*(0.5+0.5*math.Sin(now*2*math.Pi*0.6))
	var path ui.Path
	for i := 0; i <= 24; i++ {
		x, y := traceSpinnerPoint(head-length+length*float64(i)/24, s, rad)
		px, py := r.X+w/2+float32(x), r.Y+w/2+float32(y)
		if i == 0 {
			path.MoveTo(px, py)
		} else {
			path.LineTo(px, py)
		}
	}
	path.Cap(ui.RoundCap)
	p.StrokePath(&path, w, ink)
}

// traceSpinnerPoint is the point at fraction f round a rounded square of
// side s and corner rad, clockwise from the top middle.
func traceSpinnerPoint(f, s, rad float64) (float64, float64) {
	straight := s - 2*rad
	arc := math.Pi / 2 * rad
	side := straight + arc
	d := math.Mod(f, 1) * 4 * side
	// Start at the top middle, going right.
	d += straight / 2
	for corner := 0; ; corner++ {
		if d < straight {
			switch corner % 4 {
			case 0:
				return rad + d, 0
			case 1:
				return s, rad + d
			case 2:
				return s - rad - d, s
			default:
				return 0, s - rad - d
			}
		}
		d -= straight
		if d < arc {
			a := d / rad
			switch corner % 4 {
			case 0:
				return s - rad + rad*math.Sin(a), rad - rad*math.Cos(a)
			case 1:
				return s - rad + rad*math.Cos(a), s - rad + rad*math.Sin(a)
			case 2:
				return rad - rad*math.Sin(a), s - rad + rad*math.Cos(a)
			default:
				return rad - rad*math.Cos(a), rad - rad*math.Sin(a)
			}
		}
		d -= arc
	}
}

// drawMorphSpinner draws in r, at time now, a square rounding into a circle
// as it turns; its sizes are for r 28 wide.
func drawMorphSpinner(p *ui.Painter, r ui.Rect, now float64, ink ui.Color) {
	f := r.W / 28
	k := 0.5 - 0.5*math.Cos(now*2*math.Pi*0.5)
	s := 20 * f
	least := 5 * f
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	p.Transform(ui.Rotation(float32(k*90)), cx, cy, func() {
		p.Fill(ui.Rect{X: cx - s/2, Y: cy - s/2, W: s, H: s}, ink, least+(s/2-least)*float32(k))
	})
}

// Size sets both width and height of the spinner in DIPs.
func (s *SpinnerElement) Size(px float32) *SpinnerElement {
	s.Element.Size(px, px)
	return s
}

// Color overrides the spinner stroke color.
func (s *SpinnerElement) Color(col ui.Color) *SpinnerElement {
	s.Element.TextColor(col)
	if s.ink != nil {
		*s.ink = col
	}
	return s
}
