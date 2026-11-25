package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// HubDiagramOptions configures HubDiagram.
type HubDiagramOptions struct {
	Nodes  int     // 4 to 6; 6 by default
	Spread float32 // how far the satellites sit from the hub, 130 to 160px; 140 by default
	Corner float32 // 0 to 20px; 20 by default
}

type hubDiagramNode struct {
	name string
	icon func() *ui.SVG
}

func hubDiagramIcon(body string) func() *ui.SVG {
	var svg *ui.SVG
	return func() *ui.SVG {
		if svg == nil {
			svg = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">` + body + `</svg>`))
		}
		return svg
	}
}

// Lucide's cpu, shield-check and chart-column.
var (
	hubDiagramCPU    = hubDiagramIcon(`<rect width="16" height="16" x="4" y="4" rx="2"/><rect width="6" height="6" x="9" y="9" rx="1"/><path d="M15 2v2"/><path d="M15 20v2"/><path d="M2 15h2"/><path d="M2 9h2"/><path d="M20 15h2"/><path d="M20 9h2"/><path d="M9 2v2"/><path d="M9 20v2"/>`)
	hubDiagramShield = hubDiagramIcon(`<path d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z"/><path d="m9 12 2 2 4-4"/>`)
	hubDiagramChart  = hubDiagramIcon(`<path d="M3 3v16a2 2 0 0 0 2 2h16"/><path d="M18 17V9"/><path d="M13 17V5"/><path d="M8 17v-3"/>`)
)

// The satellites clockwise from the top; fewer nodes drop from the end of
// the second list first.
var hubDiagramNodes = [6]hubDiagramNode{
	{"Search", icons.Search}, {"Cache", icons.Zap}, {"Storage", icons.Database},
	{"Metrics", hubDiagramChart}, {"Auth", hubDiagramShield}, {"CDN", icons.Globe},
}

const (
	hubDiagramW    = 460
	hubDiagramH    = 290
	hubDiagramHub  = 60
	hubDiagramPill = 40
	hubDiagramLap  = 2.2 // seconds a packet takes along a spoke
)

// HubDiagram is a service in the middle and the things it connects to
// around it, packets running out along the spokes. Pointing at a node
// lights up its line.
func HubDiagram(c *ui.Context, opts ...HubDiagramOptions) ui.Element {
	var o HubDiagramOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	n := 6
	if o.Nodes != 0 {
		n = int(clamp32(float32(o.Nodes), 4, 6))
	}
	spread := float32(140)
	if o.Spread != 0 {
		spread = clamp32(o.Spread, 130, 160)
	}
	corner := float32(20)
	if o.Corner != 0 {
		corner = clamp32(o.Corner, 0, 20)
	}
	if o.Corner < 0 {
		corner = 0
	}
	t := ThemeFrom(c)
	cx, cy := float32(hubDiagramW/2), float32(hubDiagramH/2)
	pos := func(i int) (float32, float32) {
		a := (-90 + float64(i)*360/float64(n)) * math.Pi / 180
		return cx + spread*float32(math.Cos(a)), cy + spread*0.72*float32(math.Sin(a))
	}

	root := ui.Box(c).Size(hubDiagramW, hubDiagramH)
	lit := ui.Local(root, "hub_diagram", func() int { return -1 })
	var glow [6]float32
	for i := 0; i < n; i++ {
		glow[i] = root.Spring(fmt.Sprint("glow", i), b2f(*lit == i), ui.SpringOptions{Duration: 250 * time.Millisecond})
	}
	phase := phaseSeconds(c.Now())
	c.AnimationFrame()

	root.Draw(func(p *ui.Painter, r ui.Rect) {
		hx, hy := r.X+cx, r.Y+cy
		for i := 0; i < n; i++ {
			x, y := pos(i)
			x, y = r.X+x, r.Y+y
			col := lerpColor(t.BorderStrong, t.Text, glow[i])
			p.StrokePolyline([][2]float32{{hx, hy}, {x, y}}, 1.5, col)
			// A packet: a short streak, brightest at its head.
			dx, dy := x-hx, y-hy
			l := float32(math.Hypot(float64(dx), float64(dy)))
			ux, uy := dx/l, dy/l
			from, to := float32(hubDiagramHub/2), l-hubDiagramPill/2
			speed := 1 + glow[i]
			u := float32(mod1(phase*float64(speed)/hubDiagramLap + float64(i)*0.37))
			head := from + (to-from)*u
			fade := clamp32(min(u, 1-u)*8, 0, 1)
			ink := lerpColor(lerpColor(t.TextSubtle, t.Surface, 0.2), t.Text, glow[i])
			for s := 0; s < 6; s++ {
				a0, a1 := head-float32(s+1)*2.6, head-float32(s)*2.6
				if a0 < from {
					break
				}
				p.StrokePolyline([][2]float32{{hx + ux*a0, hy + uy*a0}, {hx + ux*a1, hy + uy*a1}}, 2, ink.Alpha(fade*(1-float32(s)/6)))
			}
		}
	})

	hovered := -1
	root.Children(func() {
		ui.Box(c).Key("hub").Absolute().Left(cx-hubDiagramHub/2).Top(cy-hubDiagramHub/2).Size(hubDiagramHub, hubDiagramHub).
			Radius(corner).Background(t.Surface).Center().Children(func() {
			ui.Icon(c, hubDiagramCPU()).Size(24, 24).TextColor(t.Text)
		})
		for i := 0; i < n; i++ {
			nd := hubDiagramNodes[i]
			x, y := pos(i)
			pill := ui.Row(c).Key(nd.name).Height(hubDiagramPill).Padding(0, 16, 0, 16).Radius(hubDiagramPill / 2).
				Background(t.Surface).AlignItems(ui.Center).Gap(9)
			if pill.Hovered() {
				hovered = i
			}
			w := pill.LayoutBounds().W
			if w == 0 {
				w = 100
			}
			s := pill.Spring("lift", 1+0.04*glow[i], ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.3})
			pill.Absolute().Left(x-w/2).Top(y-hubDiagramPill/2).Scale(s, s).Children(func() {
				ui.Icon(c, nd.icon()).Size(16, 16).TextColor(lerpColor(lerpColor(t.TextSubtle, t.Surface, 0.15), t.Text, glow[i]))
				ui.Text(c, nd.name).FontSize(14.5).TextColor(t.Text)
			})
		}
	})
	*lit = hovered
	return root
}

func init() {
	registerDemo(
		"hub-diagram",
		"A service hub and what it connects to, packets running along the spokes; pointing at a node lights its line.",
		func(c *ui.Context) { HubDiagram(c) },
	)
}
