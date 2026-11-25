package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// SpecInspectorOptions configures SpecInspector. bencho.dev's block has no
// props; Corner is the card's, 20 by default.
type SpecInspectorOptions struct {
	Corner float32
}

type specInspectorState struct {
	on       [3]bool
	pressing bool
	x0, y0   float32 // where the frame was started, stage coordinates
	x1, y1   float32 // its other corner
}

type specInspectorRow struct {
	title, note string
	from, to    string // the avatar's gradient
}

const (
	specInspectorW     = 380
	specInspectorH     = 300
	specInspectorCardX = 40
	specInspectorCardY = 50
	specInspectorCardW = 300
	specInspectorCardH = 200
	specInspectorPad   = 16
	specInspectorRow0  = 68 // the first row's middle, card coordinates
	specInspectorPitch = 48
	specInspectorDot   = 32
	specInspectorTextX = 60
	specInspectorTogX  = 244
	specInspectorTogW  = 40
	specInspectorTogH  = 24
)

var (
	specInspectorRows = []specInspectorRow{
		{"Mentions", "When someone tags you", "#8c74f7", "#6a4ff5"},
		{"Replies", "Replies to your threads", "#72eaf4", "#4fcdf2"},
		{"Digest", "A summary every Monday", "#f77fd2", "#f1499f"},
	}
	specInspectorInk   = ui.Hex("#1f55e0") // blueprint paper
	specInspectorFrame = ui.Hex("#3b7bf6")
)

// specInspectorText is one line of the card: where it sits and how it is set,
// so the blueprint can measure what the card really draws.
type specInspectorText struct {
	s      string
	x, y   float32
	size   float32
	weight int
}

func specInspectorTexts() []specInspectorText {
	ts := []specInspectorText{{"Notifications", specInspectorPad, 17, 16, 500}}
	for i, r := range specInspectorRows {
		cy := float32(specInspectorRow0 + i*specInspectorPitch)
		ts = append(ts,
			specInspectorText{r.title, specInspectorTextX, cy - 19, 14, 500},
			specInspectorText{r.note, specInspectorTextX, cy - 1, 12.5, 400})
	}
	return ts
}

// SpecInspector is a settings card that shows its blueprint under a frame.
// Dragging out a frame anywhere on it turns the part inside into blue paper
// with the card's boxes outlined and measured off its real layout; letting go
// melts the blueprint back into the card. The switches still flip on a click.
func SpecInspector(c *ui.Context, opts ...SpecInspectorOptions) ui.Element {
	var o SpecInspectorOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	corner := clamp32(o.Corner, 0, 28)
	t := ThemeFrom(c)
	stage := ui.Box(c).Size(specInspectorW, specInspectorH).Cursor(ui.CursorCrosshair)
	st := ui.Local(stage, "spec_inspector", func() specInspectorState {
		return specInspectorState{on: [3]bool{true, true, false}}
	})

	x, y, _ := stage.PointerPosition()
	x, y = clamp32(x, 0, specInspectorW), clamp32(y, 0, specInspectorH)
	if stage.Pressed() {
		if !st.pressing {
			st.x0, st.y0 = x, y
		}
		st.x1, st.y1 = x, y
	}
	st.pressing = stage.Pressed()
	frame := ui.Rect{X: min(st.x0, st.x1), Y: min(st.y0, st.y1), W: float32(math.Abs(float64(st.x1 - st.x0))), H: float32(math.Abs(float64(st.y1 - st.y0)))}
	big := frame.W > 6 || frame.H > 6
	show := clamp32(stage.Spring("show", b2f(st.pressing && big), ui.SpringOptions{Duration: 260 * time.Millisecond}), 0, 1)
	if !big {
		show = 0
	}
	texts := specInspectorTexts()

	stage.Children(func() {
		ui.Box(c).Absolute().Left(specInspectorCardX).Top(specInspectorCardY).Size(specInspectorCardW, specInspectorCardH).
			Radius(corner).Background(t.Surface).Children(func() {
			for i, tx := range texts {
				col := t.Text
				if i > 0 && i%2 == 0 {
					col = t.TextSubtle
				}
				ui.Text(c, tx.s).Key(i).Absolute().Left(tx.x).Top(tx.y).FontSize(tx.size).FontWeight(tx.weight).TextColor(col).SingleLine()
			}
			for i, r := range specInspectorRows {
				cy := float32(specInspectorRow0 + i*specInspectorPitch)
				ui.Box(c).Key("dot"+r.title).Absolute().Left(specInspectorPad).Top(cy-specInspectorDot/2).Size(specInspectorDot, specInspectorDot).
					Draw(func(p *ui.Painter, b ui.Rect) {
						p.FillRadialGradient(b, ui.RadialGradient{
							Center: ui.Point{X: 0.35, Y: 0.3}, Radius: specInspectorDot * 0.8,
							Stops: []ui.GradientStop{{Color: ui.Hex(r.from), At: 0}, {Color: ui.Hex(r.to), At: 1}},
						}, specInspectorDot/2)
					})
				tog := ui.Box(c).Key("tog"+r.title).Absolute().Left(specInspectorTogX).Top(cy-specInspectorTogH/2).
					Size(specInspectorTogW, specInspectorTogH).Radius(specInspectorTogH / 2).Cursor(ui.CursorPointer).
					OnClick(func() { st.on[i] = !st.on[i] })
				k := clamp32(tog.Spring("on", b2f(st.on[i]), ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: 0.2}), -0.1, 1.1)
				tog.Background(lerpColor(lerpColor(t.SurfaceSubtle, t.Border, 0.6), t.Text, clamp32(k, 0, 1)))
				tog.Children(func() {
					ui.Box(c).Absolute().Left(2+16*k).Top(2).Size(20, 20).Radius(10).Background(t.Surface)
				})
			}
		})
	})

	stage.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if show <= 0.01 {
			return
		}
		fr := ui.Rect{X: r.X + frame.X, Y: r.Y + frame.Y, W: frame.W, H: frame.H}
		card := ui.Rect{X: r.X + specInspectorCardX, Y: r.Y + specInspectorCardY, W: specInspectorCardW, H: specInspectorCardH}
		p.Clip(fr, 0, func() {
			p.Clip(card, corner, func() { specInspectorBlueprint(p, card, texts, st.on, show) })
		})
		// The frame: a blue outline, square handles and its size under it.
		blue := specInspectorFrame.Alpha(show)
		p.Fill(fr, blue.Alpha(0.06*show), 0)
		p.Stroke(fr, blue, 0, 1)
		for _, h := range [][2]float32{{0, 0}, {1, 0}, {0, 1}, {1, 1}} {
			hx, hy := fr.X+fr.W*h[0], fr.Y+fr.H*h[1]
			hr := ui.Rect{X: hx - 3.5, Y: hy - 3.5, W: 7, H: 7}
			p.Fill(hr, ui.Hex("#ffffff").Alpha(show), 1)
			p.Stroke(hr, blue, 1, 1)
		}
		label := ui.Span{Text: fmt.Sprintf("%.0f × %.0f", frame.W, frame.H), Size: 11, Weight: 600, Color: ui.Hex("#ffffff").Alpha(show)}
		lw, lh := p.MeasureText(200, label)
		ly := fr.Y + fr.H + 6
		if ly+lh+4 > r.Y+r.H {
			ly = fr.Y - lh - 10
		}
		pill := ui.Rect{X: fr.X + fr.W/2 - lw/2 - 6, Y: ly, W: lw + 12, H: lh + 4}
		p.Fill(pill, blue, (lh+4)/2)
		p.RichText(pill.X+6, pill.Y+2, lw+4, label)
	})
	if show > 0.01 && show < 0.99 {
		c.AnimationFrame()
	}
	return stage
}

// specInspectorBlueprint draws the card as a blueprint over card: blue paper
// with a grid, every box outlined in white and the spacing measured.
func specInspectorBlueprint(p *ui.Painter, card ui.Rect, texts []specInspectorText, on [3]bool, a float32) {
	p.Fill(card, specInspectorInk.Alpha(a), 0)
	white := ui.Hex("#ffffff")
	line := white.Alpha(0.85 * a)
	faint := white.Alpha(0.12 * a)
	var grid ui.Path
	for gx := card.X + 10; gx < card.X+card.W; gx += 10 {
		grid.MoveTo(gx, card.Y).LineTo(gx, card.Y+card.H)
	}
	for gy := card.Y + 10; gy < card.Y+card.H; gy += 10 {
		grid.MoveTo(card.X, gy).LineTo(card.X+card.W, gy)
	}
	p.StrokePath(&grid, 1, faint)

	// The padding box, dashed.
	var dash ui.Path
	in := ui.Rect{X: card.X + specInspectorPad, Y: card.Y + specInspectorPad, W: card.W - 2*specInspectorPad, H: card.H - 2*specInspectorPad}
	dash.MoveTo(in.X, in.Y).LineTo(in.X+in.W, in.Y).LineTo(in.X+in.W, in.Y+in.H).LineTo(in.X, in.Y+in.H).Close().Dash(0, 3, 3)
	p.StrokePath(&dash, 1, white.Alpha(0.45*a))

	tag := func(cx, cy float32, s string) {
		sp := ui.Span{Text: s, Size: 9.5, Weight: 600, Color: specInspectorInk.Alpha(a)}
		w, h := p.MeasureText(200, sp)
		bg := ui.Rect{X: cx - w/2 - 3, Y: cy - h/2 - 1, W: w + 6, H: h + 2}
		p.Fill(bg, white.Alpha(0.95*a), 3)
		p.RichText(bg.X+3, bg.Y+1, w+4, sp)
	}

	// Text boxes, measured from what the card draws.
	for _, tx := range texts {
		sp := ui.Span{Text: tx.s, Size: tx.size, Weight: tx.weight, Color: white.Alpha(0.9 * a)}
		w, h := p.MeasureText(1000, sp)
		p.RichText(card.X+tx.x, card.Y+tx.y, w+4, sp)
		p.Stroke(ui.Rect{X: card.X + tx.x - 1, Y: card.Y + tx.y, W: w + 2, H: h}, white.Alpha(0.35*a), 1, 1)
	}
	tag(card.X+card.W-36, card.Y+26, fmt.Sprintf("%d × %d", int(card.W), int(card.H)))

	var marks ui.Path
	for i := range specInspectorRows {
		cy := card.Y + float32(specInspectorRow0+i*specInspectorPitch)
		dot := ui.Rect{X: card.X + specInspectorPad, Y: cy - specInspectorDot/2, W: specInspectorDot, H: specInspectorDot}
		p.Stroke(dot, line, specInspectorDot/2, 1.2)
		tog := ui.Rect{X: card.X + specInspectorTogX, Y: cy - specInspectorTogH/2, W: specInspectorTogW, H: specInspectorTogH}
		p.Stroke(tog, line, specInspectorTogH/2, 1.2)
		kx := tog.X + 2 + 16*b2f(on[i])
		p.Stroke(ui.Rect{X: kx, Y: tog.Y + 2, W: 20, H: 20}, line, 10, 1)
		// Gaps: card edge to avatar, avatar to text.
		marks.MoveTo(card.X, cy).LineTo(dot.X, cy)
		marks.MoveTo(dot.X+dot.W, cy).LineTo(card.X+specInspectorTextX, cy)
		if i == 0 {
			tag(card.X+specInspectorPad/2, cy-8, "16")
			tag(dot.X+dot.W/2, dot.Y+dot.H+7, "32")
			tag(tog.X+tog.W/2, tog.Y-8, "40 × 24")
		}
	}
	// Row pitch, measured down the right of the text column.
	px := card.X + 226
	r0 := card.Y + specInspectorRow0
	marks.MoveTo(px, r0).LineTo(px, r0+specInspectorPitch)
	marks.MoveTo(px-3, r0).LineTo(px+3, r0)
	marks.MoveTo(px-3, r0+specInspectorPitch).LineTo(px+3, r0+specInspectorPitch)
	p.StrokePath(&marks, 1, line)
	tag(px, r0+specInspectorPitch/2, "48")
}

func init() {
	registerDemo(
		"spec-inspector",
		"A settings card that reveals its blueprint: drag a frame over it and the part inside turns into a measured blueprint, melting back when you let go.",
		func(c *ui.Context) { SpecInspector(c) },
	)
}
