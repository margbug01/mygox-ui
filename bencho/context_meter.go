package bencho

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// ContextMeterStyle is one stacked bar, or a field of dots.
type ContextMeterStyle string

const (
	ContextMeterBar   ContextMeterStyle = "Bar"
	ContextMeterField ContextMeterStyle = "Field"
)

// ContextMeterColour is the meter in colour, or in steps of the card's ink.
type ContextMeterColour string

const (
	ContextMeterColourful ContextMeterColour = "Colour"
	ContextMeterMono      ContextMeterColour = "Mono"
)

// ContextMeterOptions configures ContextMeter.
type ContextMeterOptions struct {
	Style  ContextMeterStyle  // Bar by default
	Colour ContextMeterColour // Colour by default
}

type contextMeterKind struct {
	name      string
	tokens    float32 // thousands
	compacted float32 // thousands, after Compact
	col       ui.Color
}

var contextMeterKinds = [5]contextMeterKind{
	{"Chat", 108.6, 14.2, ui.Hex("#3570ff")},
	{"Files", 46.1, 46.1, ui.Hex("#d170f0")},
	{"Tools", 12, 4.8, ui.Hex("#ff6f19")},
	{"System", 5.5, 5.5, ui.Hex("#ffb619")},
	{"Memory", 3.4, 3.4, ui.Hex("#2ec04f")},
}

type contextMeterState struct {
	compact bool
	lit     int // the kind under the pointer, or -1
}

const (
	contextMeterW      = 309
	contextMeterH      = 249
	contextMeterWindow = 200 // thousands of tokens
	contextMeterBarW   = 277
	contextMeterGap    = 2.5
	contextMeterCols   = 22
	contextMeterRows   = 3
)

// ContextMeter shows a model's token window as one stacked bar, or as a
// field of dots, by kind of content. Pointing at a kind, in the meter or
// its chip, lights it up and shows its count; Compact squeezes the chat
// down and Restore brings it back.
func ContextMeter(c *ui.Context, opts ...ContextMeterOptions) ui.Element {
	var o ContextMeterOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	card := ui.Box(c).Size(contextMeterW, contextMeterH).Radius(20).Background(t.Surface).Clip()
	st := ui.Local(card, "context_meter", func() contextMeterState { return contextMeterState{lit: -1} })

	colOf := func(i int) ui.Color {
		if o.Colour == ContextMeterMono {
			return lerpColor(t.Text, t.Surface, float32(i)*0.17)
		}
		return contextMeterKinds[i].col
	}
	spring := ui.SpringOptions{Duration: 650 * time.Millisecond, Bounce: 0.2}
	var tokens [5]float32
	total := float32(0)
	for i, k := range contextMeterKinds {
		target := k.tokens
		if st.compact {
			target = k.compacted
		}
		tokens[i] = card.Spring(fmt.Sprint("tokens", i), target, spring)
		total += tokens[i]
	}
	var glow [5]float32 // 1 lit or nothing lit, fading to 0 when another is
	for i := range glow {
		on := st.lit < 0 || st.lit == i
		glow[i] = card.Spring(fmt.Sprint("glow", i), b2f(on), ui.SpringOptions{Duration: 250 * time.Millisecond})
	}
	freeGlow := card.Spring("glowfree", b2f(st.lit < 0), ui.SpringOptions{Duration: 250 * time.Millisecond})
	faded := func(col ui.Color, g float32) ui.Color { return lerpColor(lerpColor(t.Surface, col, 0.22), col, g) }

	lit := -1
	card.Children(func() {
		meter := ui.Box(c).Key("meter").Absolute().Left(15).Top(107).Size(contextMeterBarW, 38)
		if meter.Hovered() {
			if x, _, ok := meter.PointerPosition(); ok {
				if o.Style == ContextMeterField {
					col := int(clamp32(x/(contextMeterBarW/contextMeterCols), 0, contextMeterCols-1))
					lit = contextMeterDotKind(tokens, col*contextMeterRows)
				} else {
					avail := float32(contextMeterBarW - 5*contextMeterGap)
					acc := float32(0)
					for i := range tokens {
						acc += tokens[i] / contextMeterWindow * avail
						if x < acc+float32(i)*contextMeterGap+contextMeterGap/2 {
							lit = i
							break
						}
					}
				}
			}
		}

		ui.Text(c, "Context").FontSize(16).FontWeight(500).TextColor(t.Text).Absolute().Left(15).Top(19)
		label := "Compact"
		if st.compact {
			label = "Restore"
		}
		btn := ui.Box(c).Key("compact").Absolute().Right(17).Top(16).Height(28).Padding(0, 12, 0, 12).Radius(14).
			Background(t.SurfaceHover).Center().Cursor(ui.CursorPointer).OnClick(func() { st.compact = !st.compact })
		s := btn.Spring("press", 1-0.05*b2f(btn.Pressed()), ui.SpringOptions{Duration: 200 * time.Millisecond})
		btn.Scale(s, s).Children(func() {
			ui.Text(c, label).FontSize(14).FontWeight(500).TextColor(t.Text)
		})

		figure, sub := total, fmt.Sprintf("of %dK", contextMeterWindow)
		if st.lit >= 0 {
			figure, sub = tokens[st.lit], contextMeterKinds[st.lit].name
		}
		ui.Row(c).Absolute().Left(14).Top(55).Height(36).AlignItems(ui.End).Children(func() {
			ui.Text(c, fmt.Sprintf("%.1fK", figure)).FontSize(31).TextColor(t.Text).Height(36)
			ui.Text(c, sub).FontSize(12.5).TextColor(t.TextSubtle).Margin(0, 0, 5, 5)
		})

		meter.Draw(func(p *ui.Painter, r ui.Rect) {
			if o.Style == ContextMeterField {
				pitch := r.W / contextMeterCols
				dot := pitch - 3
				rowPitch := r.H / contextMeterRows
				for col := 0; col < contextMeterCols; col++ {
					for row := 0; row < contextMeterRows; row++ {
						k := contextMeterDotKind(tokens, col*contextMeterRows+row)
						fill := faded(t.SurfaceHover, freeGlow)
						if k >= 0 {
							fill = faded(colOf(k), glow[k])
						}
						x := r.X + float32(col)*pitch + 1.5
						y := r.Y + float32(row)*rowPitch + (rowPitch-dot)/2
						p.Fill(ui.Rect{X: x, Y: y, W: dot, H: dot}, fill, 3)
					}
				}
				return
			}
			avail := r.W - 5*contextMeterGap
			x := r.X
			for i := range tokens {
				w := tokens[i] / contextMeterWindow * avail
				if w > 0.3 {
					p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, faded(colOf(i), glow[i]), min(3.5, w/2))
				}
				x += w + contextMeterGap
			}
			if free := r.X + r.W - x; free > 0.3 {
				p.Fill(ui.Rect{X: x, Y: r.Y, W: free, H: r.H}, faded(lerpColor(t.SurfaceHover, t.Surface, 0.25), freeGlow), min(3.5, free/2))
			}
		})

		// The chips, two rows: Chat, Files and Tools, then System and Memory.
		rows := [2][]int{{0, 1, 2}, {3, 4}}
		for ri, row := range rows {
			ui.Row(c).Key(fmt.Sprint("chips", ri)).Absolute().Left(16).Top(163 + float32(ri)*38).Gap(7).Children(func() {
				for _, i := range row {
					k := contextMeterKinds[i]
					chip := ui.Row(c).Key(k.name).Height(30).Padding(0, 13, 0, 11).Radius(15).AlignItems(ui.Center).Gap(8)
					if chip.Hovered() {
						lit = i
					}
					bg := t.SurfaceSubtle
					if st.lit == i {
						bg = lerpColor(t.SurfaceSubtle, t.Text, 0.07)
					}
					chip.Background(bg).Children(func() {
						ui.Box(c).Size(8, 8).Radius(2).Background(colOf(i))
						ui.Text(c, k.name).FontSize(14).TextColor(t.Text)
					})
				}
			})
		}
	})
	if lit != st.lit {
		st.lit = lit
		c.AnimationFrame()
	}
	return card
}

// contextMeterDotKind is the kind dot n of the field belongs to, filling
// column by column, or -1 for free space.
func contextMeterDotKind(tokens [5]float32, n int) int {
	per := float32(contextMeterWindow) / (contextMeterCols * contextMeterRows)
	at := (float32(n) + 0.5) * per
	acc := float32(0)
	for i, v := range tokens {
		acc += v
		if at < acc {
			return i
		}
	}
	return -1
}

func init() {
	registerDemo(
		"context-meter",
		"A token window as one stacked bar or a field of dots; point at a kind to light it up, and Compact squeezes the chat down.",
		func(c *ui.Context) { ContextMeter(c) },
	)
}
