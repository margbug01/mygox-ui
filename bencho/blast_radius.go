package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// BlastRadiusColour is the people in colour, or in greys.
type BlastRadiusColour string

const (
	BlastRadiusColourful BlastRadiusColour = "Colour"
	BlastRadiusMono      BlastRadiusColour = "Mono"
)

// BlastRadiusHops is how far out the graph opens.
type BlastRadiusHops string

const (
	BlastRadiusAll  BlastRadiusHops = "All"
	BlastRadiusNear BlastRadiusHops = "Near"
)

// BlastRadiusShape is the nodes as circles or rounded squares.
type BlastRadiusShape string

const (
	BlastRadiusCircle BlastRadiusShape = "Circle"
	BlastRadiusSquare BlastRadiusShape = "Square"
)

// BlastRadiusOptions configures BlastRadius.
type BlastRadiusOptions struct {
	Colour   BlastRadiusColour // Colour by default
	Hops     BlastRadiusHops   // where the graph starts, All by default
	Shape    BlastRadiusShape  // Circle by default
	NoCounts bool              // hides the close, new and people counts under the graph
	Corner   float32           // 0 to 28px; 20 by default
}

// blastRadiusFriend is someone one hop out, and the people they know.
type blastRadiusFriend struct {
	angle float32 // degrees clockwise from east
	col   ui.Color
	close bool
	outer []blastRadiusContact
}

type blastRadiusContact struct {
	angle float32
	new   bool
}

var blastRadiusNetwork = []blastRadiusFriend{
	{-90, ui.Hex("#9a7cf6"), true, []blastRadiusContact{{-107, true}, {-90, false}, {-73, false}}},
	{-30, ui.Hex("#9bb8da"), true, []blastRadiusContact{{-45, true}, {-29, false}, {-14.5, true}}},
	{30, ui.Hex("#5c8bf5"), true, []blastRadiusContact{{20, true}, {42, false}}},
	{90, ui.Hex("#ec78b9"), false, []blastRadiusContact{{104, false}, {76, false}}},
	{150, ui.Hex("#6fdfaf"), false, []blastRadiusContact{{160, false}, {138, false}}},
	{210, ui.Hex("#f4c567"), true, []blastRadiusContact{{220, false}, {196, true}}},
}

type blastRadiusState struct {
	near bool
}

const (
	blastRadiusW     = 300
	blastRadiusH     = 330
	blastRadiusCX    = 150
	blastRadiusCY    = 174
	blastRadiusInner = 54
	blastRadiusOuter = 100
)

// blastRadiusDashedCircle strokes a circle in short dashes.
func blastRadiusDashedCircle(p *ui.Painter, cx, cy, r, dash, gap float32, col ui.Color) {
	step := (dash + gap) / r
	for a := float32(0); a < 2*math.Pi; a += step {
		b := a + dash/r
		p.StrokePolyline([][2]float32{{cx + r*float32(math.Cos(float64(a))), cy + r*float32(math.Sin(float64(a)))},
			{cx + r*float32(math.Cos(float64(b))), cy + r*float32(math.Sin(float64(b)))}}, 1, col)
	}
}

// blastRadiusDashedLine strokes a line in short dashes.
func blastRadiusDashedLine(p *ui.Painter, x0, y0, x1, y1 float32, col ui.Color) {
	dx, dy := x1-x0, y1-y0
	l := float32(math.Hypot(float64(dx), float64(dy)))
	for s := float32(0); s < l; s += 6 {
		e := min(s+3.5, l)
		p.StrokePolyline([][2]float32{{x0 + dx*s/l, y0 + dy*s/l}, {x0 + dx*e/l, y0 + dy*e/l}}, 1, col)
	}
}

func blastRadiusAt(cx, cy, r, deg float32) (float32, float32) {
	a := float64(deg) * math.Pi / 180
	return cx + r*float32(math.Cos(a)), cy + r*float32(math.Sin(a))
}

// BlastRadius is a radial network graph: you in the middle, the people you
// know ring by ring around you. Resting on someone lights the path back to
// you and fades the rest; Near folds the graph to the first ring and All
// grows the outer one back.
func BlastRadius(c *ui.Context, opts ...BlastRadiusOptions) ui.Element {
	var o BlastRadiusOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	corner := float32(20)
	if o.Corner != 0 {
		corner = clamp32(o.Corner, 0, 28)
	}
	t := ThemeFrom(c)
	card := ui.Box(c).Size(blastRadiusW, blastRadiusH).Radius(corner).Background(t.Surface).Clip()
	st := ui.Local(card, "blast_radius", func() blastRadiusState { return blastRadiusState{near: o.Hops == BlastRadiusNear} })
	grow := card.Spring("outer", b2f(!st.near), ui.SpringOptions{Duration: 650 * time.Millisecond, Bounce: 0.2})
	outerR := lerp32(blastRadiusInner, blastRadiusOuter, grow)

	// The person under the pointer: friend f, and contact k of theirs or -1.
	hf, hk := -1, -1
	if card.Hovered() {
		if x, y, ok := card.PointerPosition(); ok {
			best := float32(14 * 14)
			for f, fr := range blastRadiusNetwork {
				fx, fy := blastRadiusAt(blastRadiusCX, blastRadiusCY, blastRadiusInner, fr.angle)
				if d := (x-fx)*(x-fx) + (y-fy)*(y-fy); d < best {
					best, hf, hk = d, f, -1
				}
				if grow < 0.5 {
					continue
				}
				for k, ct := range fr.outer {
					ox, oy := blastRadiusAt(blastRadiusCX, blastRadiusCY, outerR, ct.angle)
					if d := (x-ox)*(x-ox) + (y-oy)*(y-oy); d < best {
						best, hf, hk = d, f, k
					}
				}
			}
		}
	}
	dim := card.Spring("dim", b2f(hf >= 0), ui.SpringOptions{Duration: 250 * time.Millisecond})
	litK := hk

	colOf := func(f int) ui.Color {
		if o.Colour == BlastRadiusMono {
			return lerpColor(t.Text, t.Surface, 0.3+0.08*float32(f))
		}
		return blastRadiusNetwork[f].col
	}
	onPath := func(f, k int) bool { return hf == f && (k < 0 || litK == k) }
	fadeTo := func(col ui.Color, lit bool) ui.Color {
		if lit {
			return col
		}
		return lerpColor(col, t.Surface, 0.75*dim)
	}
	node := func(p *ui.Painter, x, y, r float32, col ui.Color) {
		rad := r
		if o.Shape == BlastRadiusSquare {
			rad = r * 0.4
		}
		p.Fill(ui.Rect{X: x - r, Y: y - r, W: 2 * r, H: 2 * r}, col, rad)
	}
	ring := func(p *ui.Painter, x, y, r float32, col ui.Color) {
		rad := r
		if o.Shape == BlastRadiusSquare {
			rad = r * 0.45
		}
		p.Stroke(ui.Rect{X: x - r, Y: y - r, W: 2 * r, H: 2 * r}, col, rad, 1.5)
	}

	card.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy := r.X+blastRadiusCX, r.Y+blastRadiusCY
		guide := lerpColor(t.BorderStrong, t.Surface, 0.15)
		blastRadiusDashedCircle(p, cx, cy, blastRadiusInner, 3, 3, guide)
		if grow > 0.01 {
			blastRadiusDashedCircle(p, cx, cy, outerR, 3, 3, guide.Alpha(clamp32(grow, 0, 1)))
		}
		// Lines first, then the people over them.
		for f, fr := range blastRadiusNetwork {
			fx, fy := blastRadiusAt(cx, cy, blastRadiusInner, fr.angle)
			if grow > 0.01 {
				for k, ct := range fr.outer {
					ox, oy := blastRadiusAt(cx, cy, outerR, ct.angle)
					col := fadeTo(t.Text, onPath(f, k))
					if ct.new {
						blastRadiusDashedLine(p, fx, fy, ox, oy, fadeTo(lerpColor(t.TextSubtle, t.Surface, 0.2), onPath(f, k)).Alpha(grow))
					} else {
						p.StrokePolyline([][2]float32{{fx, fy}, {ox, oy}}, 1+0.5*b2f(onPath(f, k)), col.Alpha(clamp32(grow, 0, 1)))
					}
				}
			}
			p.StrokePolyline([][2]float32{{cx, cy}, {fx, fy}}, 1.5+0.5*b2f(onPath(f, -1)), fadeTo(t.Text, onPath(f, -1)))
		}
		for f, fr := range blastRadiusNetwork {
			fx, fy := blastRadiusAt(cx, cy, blastRadiusInner, fr.angle)
			if grow > 0.01 {
				for k, ct := range fr.outer {
					ox, oy := blastRadiusAt(cx, cy, outerR, ct.angle)
					a := clamp32(grow, 0, 1)
					if ct.new {
						p.Fill(ui.Rect{X: ox - 6.5, Y: oy - 6.5, W: 13, H: 13}, t.Surface.Alpha(a), 6.5)
						blastRadiusDashedCircle(p, ox, oy, 6.5, 1.6, 1.6, fadeTo(t.TextSubtle, onPath(f, k)).Alpha(a))
					}
					rr := float32(3.5) + 1.5*b2f(onPath(f, k))
					node(p, ox, oy, rr*a, fadeTo(colOf(f), onPath(f, k)).Alpha(a))
				}
			}
			lit := onPath(f, -1) || hf == f
			if fr.close {
				p.Fill(ui.Rect{X: fx - 10, Y: fy - 10, W: 20, H: 20}, t.Surface, 10)
				ring(p, fx, fy, 9.5, fadeTo(t.Text, lit))
			}
			node(p, fx, fy, 6+b2f(hf == f && litK < 0), fadeTo(colOf(f), lit))
		}
		// You, ringed, in the middle.
		p.Fill(ui.Rect{X: cx - 21, Y: cy - 21, W: 42, H: 42}, t.Surface, 21)
		ring(p, cx, cy, 20.5, t.Text)
		youR := float32(16.5)
		youRad := youR
		if o.Shape == BlastRadiusSquare {
			youRad = youR * 0.4
		}
		you := ui.RadialGradient{Center: ui.Point{X: -0.2, Y: -0.25}, Radius: 2 * youR, Stops: []ui.GradientStop{
			{At: 0, Color: ui.Hex("#c879fb")}, {At: 0.6, Color: ui.Hex("#bb45ee")}, {At: 1, Color: ui.Hex("#a93ae0")},
		}}
		if o.Colour == BlastRadiusMono {
			you.Stops = []ui.GradientStop{{At: 0, Color: lerpColor(t.Text, t.Surface, 0.35)}, {At: 1, Color: t.Text}}
		}
		p.FillRadialGradient(ui.Rect{X: cx - youR, Y: cy - youR, W: 2 * youR, H: 2 * youR}, you, youRad)
	})

	closeN, newN, people := 0, 0, 1+len(blastRadiusNetwork)
	for _, fr := range blastRadiusNetwork {
		if fr.close {
			closeN++
		}
		if !st.near {
			people += len(fr.outer)
			for _, ct := range fr.outer {
				if ct.new {
					newN++
				}
			}
		}
	}

	card.Children(func() {
		ui.Text(c, "Network").FontSize(15.5).FontWeight(500).TextColor(t.Text).Absolute().Left(16).Top(15)
		ui.Text(c, "People around you").FontSize(13.5).TextColor(t.TextSubtle).Absolute().Left(16).Top(33)

		seg := ui.Box(c).Key("hops").Absolute().Left(198).Top(17).Size(86, 32).Radius(16).Background(t.SurfaceSubtle)
		px := seg.Spring("pill", b2f(!st.near), ui.SpringOptions{Duration: 350 * time.Millisecond, Bounce: 0.15})
		seg.Draw(func(p *ui.Painter, r ui.Rect) {
			x := lerp32(3, 49, px)
			w := lerp32(45, 34, px)
			p.Fill(ui.Rect{X: r.X + x, Y: r.Y + 3, W: w, H: r.H - 6}, t.Surface, 13)
		})
		seg.Children(func() {
			for i, label := range [2]string{"Near", "All"} {
				near := i == 0
				col := t.TextSubtle
				if near == st.near {
					col = t.Text
				}
				left, w := float32(3), float32(45)
				if !near {
					left, w = 49, 34
				}
				ui.Box(c).Key(label).Absolute().Left(left).Top(0).Size(w, 32).Center().Cursor(ui.CursorPointer).
					OnClick(func() { st.near = near }).Children(func() {
					ui.Text(c, label).FontSize(13.5).TextColor(col)
				})
			}
		})

		if o.NoCounts {
			return
		}
		ui.Row(c).Absolute().Left(17).Top(295).Height(22).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(10, 10).Radius(5).Border(1.5, t.Text).Center().Children(func() {
				ui.Box(c).Size(4, 4).Radius(2).Background(lerpColor(t.TextSubtle, t.Surface, 0.1))
			})
			ui.Text(c, "Close").FontSize(13.5).TextColor(t.TextSubtle).Margin(0, 0, 0, 6)
			ui.Text(c, fmt.Sprint(closeN)).FontSize(13.5).FontWeight(500).TextColor(t.Text).Margin(0, 0, 0, 5)
		})
		ui.Row(c).Absolute().Left(121).Top(295).Height(22).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(11, 11).Draw(func(p *ui.Painter, r ui.Rect) {
				blastRadiusDashedCircle(p, r.X+5.5, r.Y+5.5, 5, 1.5, 1.5, t.TextSubtle)
				p.Fill(ui.Rect{X: r.X + 4, Y: r.Y + 4, W: 3, H: 3}, t.TextSubtle, 1.5)
			})
			ui.Text(c, "New").FontSize(13.5).TextColor(t.TextSubtle).Margin(0, 0, 0, 6)
			ui.Text(c, fmt.Sprint(newN)).FontSize(13.5).FontWeight(500).TextColor(t.Text).Margin(0, 0, 0, 5)
		})
		ui.Row(c).Absolute().Right(17).Top(295).Height(22).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "People").FontSize(13.5).TextColor(t.TextSubtle)
			ui.Text(c, fmt.Sprint(people)).FontSize(13.5).FontWeight(500).TextColor(t.Text).Margin(0, 0, 0, 5)
		})
	})
	return card
}

func init() {
	registerDemo(
		"blast-radius",
		"A radial network graph: you in the middle, ring by ring; resting on a node lights its path back to you, and Near or All folds or grows the outer ring.",
		func(c *ui.Context) { BlastRadius(c) },
	)
}
