package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// SweepButtonOptions configures SweepButton.
type SweepButtonOptions struct {
	Speed  float32 // seconds for one lap, 1.5 to 10; 4.5 by default
	Glow   float32 // how strong the colour is, 10 to 100; 30 by default
	Line   bool    // a thin rim line in the sweep's colours
	Corner float32 // 0 to 24; 24 by default
	Label  string  // "Get started" by default
	// OnClick runs when the button is pressed.
	OnClick func()
}

const (
	sweepButtonW = 148
	sweepButtonH = 48
)

// sweepButtonHues are the soft colours travelling round the edge, each a
// third of a lap apart.
var sweepButtonHues = [3]struct {
	col ui.Color
	at  float64
}{{ui.Hex("#ff8fae"), 0}, {ui.Hex("#5fe3a1"), 0.3}, {ui.Hex("#c9a8ff"), 0.62}}

// SweepButton is a call to action whose edge carries a soft sweep of colour
// travelling round it, one lap every Speed seconds.
func SweepButton(c *ui.Context, opts ...SweepButtonOptions) ui.Element {
	var o SweepButtonOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Speed == 0 {
		o.Speed = 4.5
	}
	if o.Glow == 0 {
		o.Glow = 30
	}
	if o.Corner == 0 {
		o.Corner = 24
	}
	o.Speed = clamp32(o.Speed, 1.5, 10)
	o.Glow = clamp32(o.Glow, 10, 100)
	o.Corner = clamp32(o.Corner, 0, 24)
	o.Label = orDefault(o.Label, "Get started")
	t := ThemeFrom(c)

	btn := ui.Box(c).Size(sweepButtonW, sweepButtonH).Radius(o.Corner).Background(t.Surface).Center().
		Cursor(ui.CursorPointer).Label(o.Label)
	if o.OnClick != nil {
		btn.OnClick(o.OnClick)
	}
	press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 160 * time.Millisecond})
	btn.Scale(1-0.03*press, 1-0.03*press)
	lap := mod1(phaseSeconds(c.Now()) / float64(o.Speed))
	c.AnimationFrame()
	alpha := min(o.Glow/100*1.5, 1)

	btn.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, t.Surface, o.Corner)
		p.Clip(r, o.Corner, func() {
			for _, h := range sweepButtonHues {
				x, y := sweepButtonEdge(r, o.Corner, mod1(lap+h.at))
				const rad = 46
				p.FillRadialGradient(ui.Rect{X: x - rad, Y: y - rad, W: 2 * rad, H: 2 * rad}, ui.RadialGradient{Radius: rad, Stops: []ui.GradientStop{
					{At: 0, Color: h.col.Alpha(alpha)}, {At: 0.5, Color: h.col.Alpha(alpha * 0.45)}, {At: 1, Color: h.col.Alpha(0)},
				}}, 0)
			}
		})
		if o.Line {
			stops := []ui.GradientStop{}
			for _, h := range sweepButtonHues {
				stops = append(stops, ui.GradientStop{At: float32(h.at), Color: h.col.Alpha(min(alpha*1.6, 1))})
			}
			stops = append(stops, ui.GradientStop{At: 1, Color: sweepButtonHues[0].col.Alpha(min(alpha*1.6, 1))})
			p.StrokePathConicGradient(sweepButtonRim(r, o.Corner), 1.2, ui.ConicGradient{Angle: float32(lap*360) + 90, Stops: stops})
		}
	})
	btn.Children(func() {
		ui.Text(c, o.Label).FontSize(15.5).FontWeight(500).TextColor(t.Text).SingleLine()
	})
	return btn
}

// sweepButtonEdge is the point a fraction s of the way round the button's
// rounded outline, clockwise from the middle of its top edge.
func sweepButtonEdge(r ui.Rect, rad float32, s float64) (float32, float32) {
	rad = min(rad, r.W/2, r.H/2)
	sw, sh := float64(r.W-2*rad), float64(r.H-2*rad)
	arc := math.Pi / 2 * float64(rad)
	total := 2*sw + 2*sh + 4*arc
	d := s * total
	cx0, cy0 := float64(r.X+rad), float64(r.Y+rad)
	cx1, cy1 := float64(r.X+r.W-rad), float64(r.Y+r.H-rad)
	corner := func(cx, cy, from float64, d float64) (float32, float32) {
		a := from + d/float64(max(rad, 0.001))
		return float32(cx + float64(rad)*math.Cos(a)), float32(cy + float64(rad)*math.Sin(a))
	}
	segs := []struct {
		n  float64
		at func(d float64) (float32, float32)
	}{
		{sw / 2, func(d float64) (float32, float32) { return float32(float64(r.X+r.W/2) + d), r.Y }},
		{arc, func(d float64) (float32, float32) { return corner(cx1, cy0, -math.Pi/2, d) }},
		{sh, func(d float64) (float32, float32) { return r.X + r.W, float32(cy0 + d) }},
		{arc, func(d float64) (float32, float32) { return corner(cx1, cy1, 0, d) }},
		{sw, func(d float64) (float32, float32) { return float32(cx1 - d), r.Y + r.H }},
		{arc, func(d float64) (float32, float32) { return corner(cx0, cy1, math.Pi/2, d) }},
		{sh, func(d float64) (float32, float32) { return r.X, float32(cy1 - d) }},
		{arc, func(d float64) (float32, float32) { return corner(cx0, cy0, math.Pi, d) }},
		{sw / 2, func(d float64) (float32, float32) { return float32(cx0 + d), r.Y }},
	}
	for _, sg := range segs {
		if d <= sg.n {
			return sg.at(d)
		}
		d -= sg.n
	}
	return r.X + r.W/2, r.Y
}

// sweepButtonRim is the button's rounded outline as a path.
func sweepButtonRim(r ui.Rect, rad float32) *ui.Path {
	var p ui.Path
	const steps = 96
	for i := 0; i <= steps; i++ {
		x, y := sweepButtonEdge(r, rad, float64(i)/steps)
		if i == 0 {
			p.MoveTo(x, y)
		} else {
			p.LineTo(x, y)
		}
	}
	p.Close()
	return &p
}

func init() {
	registerDemo(
		"sweep-button",
		"A Get started button whose edge carries a soft sweep of colour travelling round it, one lap every few seconds.",
		func(c *ui.Context) { SweepButton(c) },
	)
}
