package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// GenerateOptions configures Generate.
type GenerateOptions struct {
	Pace   time.Duration // from the press to the finished picture, 3s to 8s; 5s by default
	Ripple float32       // how liquid the surface is while the picture forms, 0 to 100; 50 by default
	Corner float32       // corner radius, 0 to 40; 20 by default
}

const generateW, generateH = 270, 330

type generateState struct {
	started time.Time // zero until Generate is pressed
	runs    int       // pictures made so far; each new one turns a little
}

// Generate is an AI image loader: press Generate and a picture forms under
// frosted glass, broad shapes first and fine detail last, a badge counting
// the percent. Press the finished picture to make another.
func Generate(c *ui.Context, opts ...GenerateOptions) ui.Element {
	var o GenerateOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Pace == 0 {
		o.Pace = 5 * time.Second
	}
	if o.Ripple == 0 {
		o.Ripple = 50
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Pace = min(max(o.Pace, 3*time.Second), 8*time.Second)
	o.Ripple, o.Corner = clamp32(o.Ripple, 0, 100), clamp32(o.Corner, 0, 40)
	t := ThemeFrom(c)
	now := c.Now()

	card := ui.Box(c).Size(generateW, generateH).Radius(o.Corner).Background(t.Surface).Clip()
	st := ui.Local(card, "generate", func() generateState { return generateState{} })
	started := !st.started.IsZero()
	p := float32(0)
	if started {
		p = clamp32(float32(now.Sub(st.started))/float32(o.Pace), 0, 1)
		if p < 1 {
			c.AnimationFrame()
		}
	}
	shown := card.Spring("shown", b2f(started), ui.SpringOptions{Duration: 420 * time.Millisecond})
	done := p >= 1
	if done {
		card.Cursor(ui.CursorPointer).OnClick(func() { st.started, st.runs = c.Now(), st.runs+1 })
	}

	card.Children(func() {
		if shown > 0.01 {
			paper := ui.Hex("#e9e7e8")
			ink := lerpColor(ui.Hex("#6a6a6a"), ui.Hex("#1c1c1e"), clamp32(p*1.3, 0, 1))
			if t.Dark {
				paper, ink = ui.Hex("#2a292b"), lerpColor(ui.Hex("#8a8a8a"), ui.Hex("#f2f2f2"), clamp32(p*1.3, 0, 1))
			}
			ui.Box(c).Absolute().Left(0).Top(0).Size(generateW, generateH).Background(paper).Opacity(clamp32(shown, 0, 1))
			// Broad shapes come first, under a frost that clears as it forms.
			frost := 26 * float32(math.Pow(float64(1-p), 1.4))
			wobble := o.Ripple / 100 * (1 - p)
			s := phaseSeconds(now)
			runs := st.runs
			ui.Box(c).Absolute().Left(0).Top(0).Size(generateW, generateH).Opacity(clamp32(shown, 0, 1) * clamp32(0.25+p*4, 0, 1)).Blur(frost).
				Draw(func(pt *ui.Painter, r ui.Rect) {
					generateBloom(pt, r, ink, wobble, s, runs, 0)
				})
			// Fine detail comes last: veins and seeds on the petals.
			if d := clamp32((p-0.55)/0.45, 0, 1); d > 0 {
				ui.Box(c).Absolute().Left(0).Top(0).Size(generateW, generateH).Opacity(d).Blur(4 * (1 - d)).
					Draw(func(pt *ui.Painter, r ui.Rect) {
						generateBloom(pt, r, lerpColor(ink, paper, 0.55), wobble, s, runs, 1)
					})
			}
			if !done || now.Sub(st.started) < o.Pace+900*time.Millisecond {
				if done {
					c.After(o.Pace + 900*time.Millisecond - now.Sub(st.started))
				}
				ui.Box(c).Absolute().Left(12).Bottom(12).Height(22).PaddingX(8).Radius(11).Background(t.Surface).Center().
					Opacity(clamp32(shown, 0, 1)).Children(func() {
					ui.Text(c, fmt.Sprintf("%d%%", int(p*100))).FontSize(11).FontWeight(500).TextColor(t.Text)
				})
			}
		}
		if shown > 0.99 {
			return
		}
		btn := ui.Box(c).Absolute().Left(generateW/2-50).Top(generateH/2-21).Size(100, 42).Radius(21).Background(t.SurfaceSubtle).
			Center().Cursor(ui.CursorPointer).Opacity(1 - clamp32(shown, 0, 1))
		if !started {
			btn.OnClick(func() { st.started = c.Now() })
		}
		press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
		btn.Scale(1-0.05*press, 1-0.05*press)
		btn.Children(func() { ui.Text(c, "Generate").FontSize(15).FontWeight(500).TextColor(t.Text) })
	})
	return card
}

// generateBloom paints the picture, an ink flower of six petals swept to
// one side: layer 0 its body, layer 1 its veins and seeds. wobble bends the
// petals like liquid while the picture forms.
func generateBloom(p *ui.Painter, r ui.Rect, ink ui.Color, wobble float32, s float64, run, layer int) {
	cx, cy := r.X+r.W*0.52, r.Y+r.H*0.47
	turn := float64(run) * 0.7
	petals := []struct{ a, l, w float64 }{
		{-1.95, 150, 30}, {-1.2, 140, 34}, {-0.35, 120, 30}, {0.55, 120, 30}, {1.5, 160, 26}, {2.35, 150, 28}, {3.1, 115, 22},
	}
	for i, pt := range petals {
		a := pt.a + turn + float64(wobble)*0.22*math.Sin(s*2.2+float64(i)*1.3)
		l := float32(pt.l) * (1 + wobble*0.08*float32(math.Sin(s*3+float64(i))))
		w := float32(pt.w)
		dx, dy := float32(math.Cos(a)), float32(math.Sin(a))
		nx, ny := -dy, dx
		tipX, tipY := cx+dx*l, cy+dy*l
		// Each petal swells near the middle and curls at the tip.
		curl := float32(18) * float32(math.Sin(float64(i)*2.1))
		if layer == 0 {
			var path ui.Path
			path.MoveTo(cx+nx*w*0.5, cy+ny*w*0.5).
				CubeTo(cx+dx*l*0.35+nx*w, cy+dy*l*0.35+ny*w, cx+dx*l*0.8+nx*w*0.6+nx*curl, cy+dy*l*0.8+ny*w*0.6+ny*curl, tipX+nx*curl, tipY+ny*curl).
				CubeTo(cx+dx*l*0.8-nx*w*0.4+nx*curl, cy+dy*l*0.8-ny*w*0.4+ny*curl, cx+dx*l*0.35-nx*w, cy+dy*l*0.35-ny*w, cx-nx*w*0.5, cy-ny*w*0.5).Close()
			p.FillPath(&path, ink)
			continue
		}
		var vein ui.Path
		vein.MoveTo(cx+dx*12, cy+dy*12).QuadTo(cx+dx*l*0.6+nx*curl*0.4, cy+dy*l*0.6+ny*curl*0.4, tipX+nx*curl*0.9-dx*10, tipY+ny*curl*0.9-dy*10)
		p.StrokePath(&vein, 1.4, ink)
	}
	if layer == 0 {
		p.Fill(ui.Rect{X: cx - 34, Y: cy - 34, W: 68, H: 68}, ink, 34)
		return
	}
	for k := range 9 {
		a := float64(k) * 2 * math.Pi / 9
		p.Fill(ui.Rect{X: cx + float32(math.Cos(a))*16 - 2, Y: cy + float32(math.Sin(a))*16 - 2, W: 4, H: 4}, ink, 2)
	}
}

func init() {
	registerDemo(
		"generate",
		"An AI image generation loader: press Generate and a picture forms under frosted glass, broad shapes first and fine detail last.",
		func(c *ui.Context) { Generate(c) },
	)
}
