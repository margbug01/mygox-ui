package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// EdgeLightColour is the colour the rim catches.
type EdgeLightColour string

const (
	EdgeLightCandle EdgeLightColour = "Candle"
	EdgeLightRose   EdgeLightColour = "Rosé"
	EdgeLightScreen EdgeLightColour = "Screen"
)

// EdgeLightOptions configures EdgeLight.
type EdgeLightOptions struct {
	Colour   EdgeLightColour // Rosé by default
	Edge     float32         // how far out from the middle the light starts, 0 to 80 (%); 40 by default
	Glow     float32         // how strong the light is, 10 to 100, 50 being the reference; 80 by default
	Tilt     float32         // how far the print leans toward the pointer, 0 to 100; 20 by default
	Title    string          // "The Last Scroll" by default
	Subtitle string          // "Oil on canvas · 1841" by default
}

const (
	edgeLightW     = 340
	edgeLightH     = 360
	edgeLightCardW = 238
	edgeLightCardH = 294
	edgeLightR     = 22
)

// EdgeLight is a framed painting whose border glows toward the cursor: as
// the pointer moves out to an edge a short arc of light catches the rim and
// travels round with it, and the print leans toward the pointer.
func EdgeLight(c *ui.Context, opts ...EdgeLightOptions) ui.Element {
	var o EdgeLightOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	o.Colour = EdgeLightColour(orDefault(string(o.Colour), string(EdgeLightRose)))
	if o.Edge == 0 {
		o.Edge = 40
	}
	if o.Glow == 0 {
		o.Glow = 80
	}
	if o.Tilt == 0 {
		o.Tilt = 20
	}
	o.Edge = clamp32(o.Edge, 0, 80)
	o.Glow = clamp32(o.Glow, 10, 100)
	o.Tilt = clamp32(o.Tilt, 0, 100)
	o.Title = orDefault(o.Title, "The Last Scroll")
	o.Subtitle = orDefault(o.Subtitle, "Oil on canvas · 1841")
	light := map[EdgeLightColour]ui.Color{
		EdgeLightCandle: ui.Hex("#ffae3d"),
		EdgeLightRose:   ui.Hex("#ff6f9f"),
		EdgeLightScreen: ui.Hex("#7cc0ff"),
	}[o.Colour]
	if light == (ui.Color{}) {
		light = ui.Hex("#ff6f9f")
	}
	t := ThemeFrom(c)
	art := genieArtwork()

	stage := ui.Box(c).Size(edgeLightW, edgeLightH).Center()
	// Where the pointer is, as a fraction of the card's half size from its middle.
	px, py, over := stage.PointerPosition()
	nx := (px - edgeLightW/2) / (edgeLightCardW / 2)
	ny := (py - edgeLightH/2) / (edgeLightCardH / 2)
	out := float32(math.Max(math.Abs(float64(nx)), math.Abs(float64(ny))))
	near := float32(0)
	if over {
		near = clamp32((out-o.Edge/100)/max(1-o.Edge/100, 0.05), 0, 1)
	}
	angle := float32(math.Atan2(float64(nx*edgeLightCardW), float64(-ny*edgeLightCardH))) * 180 / math.Pi
	spring := ui.SpringOptions{Duration: 260 * time.Millisecond}
	lit := stage.Spring("lit", near, spring)
	ang := edgeLightUnwrap(stage, angle)
	tiltX := stage.Spring("tx", b2f(over)*clamp32(nx, -1.2, 1.2), spring)
	tiltY := stage.Spring("ty", b2f(over)*clamp32(ny, -1.2, 1.2), spring)
	lean := o.Tilt / 100 * 12

	stage.Children(func() {
		card := ui.Column(c).Size(edgeLightCardW, edgeLightCardH).Radius(edgeLightR).Background(t.Surface).
			Border(1, lerpColor(t.Border, t.BorderStrong, 0.5)).Padding(9, 9, 0, 9).
			RotateY(lean * tiltX).RotateX(-lean * tiltY).Perspective(900)
		card.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if lit < 0.01 {
				return
			}
			k := lit * o.Glow / 50
			path := edgeLightRim(r, edgeLightR)
			arc := float32(70)
			g := func(a float32) ui.ConicGradient {
				col := light.Alpha(clamp32(a*k, 0, 1))
				return ui.ConicGradient{Angle: ang - arc, Stops: []ui.GradientStop{
					{At: 0, Color: col.Alpha(0)}, {At: arc / 360, Color: col}, {At: 2 * arc / 360, Color: col.Alpha(0)}, {At: 1, Color: col.Alpha(0)},
				}}
			}
			p.StrokePathConicGradient(path, 14, g(0.10))
			p.StrokePathConicGradient(path, 7, g(0.22))
			p.StrokePathConicGradient(path, 3, g(0.55))
			p.StrokePathConicGradient(path, 1.4, g(1))
		})
		card.Children(func() {
			ui.Box(c).Size(edgeLightCardW-20, edgeLightCardW-20).Draw(func(p *ui.Painter, r ui.Rect) {
				p.Clip(r, 14, func() { p.ImageRegion(art[0], ui.Rect{}, r, ui.ImageOptions{}) })
			})
			ui.Text(c, o.Title).FontSize(16).FontWeight(500).TextColor(t.Text).SingleLine().Margin(13, 0, 0, 6)
			ui.Text(c, o.Subtitle).FontSize(14).TextColor(lerpColor(t.TextSubtle, t.Surface, 0.15)).SingleLine().Margin(4, 0, 0, 6)
		})
	})
	return stage
}

// edgeLightUnwrap springs toward angle a (degrees) the shorter way round.
func edgeLightUnwrap(e ui.Element, a float32) float32 {
	st := ui.Local(e, "edge_light_angle", func() float32 { return a })
	d := float32(math.Mod(float64(a-*st)+540, 360) - 180)
	*st += d
	return e.Spring("arc", *st, ui.SpringOptions{Duration: 220 * time.Millisecond})
}

// edgeLightRim is the card's rounded outline.
func edgeLightRim(r ui.Rect, rad float32) *ui.Path {
	const k = 0.5523
	var p ui.Path
	x0, y0, x1, y1 := r.X, r.Y, r.X+r.W, r.Y+r.H
	p.MoveTo(x0+rad, y0).LineTo(x1-rad, y0)
	p.CubeTo(x1-rad+rad*k, y0, x1, y0+rad-rad*k, x1, y0+rad)
	p.LineTo(x1, y1-rad)
	p.CubeTo(x1, y1-rad+rad*k, x1-rad+rad*k, y1, x1-rad, y1)
	p.LineTo(x0+rad, y1)
	p.CubeTo(x0+rad-rad*k, y1, x0, y1-rad+rad*k, x0, y1-rad)
	p.LineTo(x0, y0+rad)
	p.CubeTo(x0, y0+rad-rad*k, x0+rad-rad*k, y0, x0+rad, y0)
	p.Close()
	return &p
}

func init() {
	registerDemo(
		"edge-light",
		"A framed painting whose border glows toward the cursor: a short arc of light catches the rim near the pointer and travels round with it.",
		func(c *ui.Context) { EdgeLight(c) },
	)
}
