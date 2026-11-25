package bencho

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// BlueprintLoaderOptions configures BlueprintLoader.
type BlueprintLoaderOptions struct {
	Speed float32 // how quickly the plan draws and the card arrives, 0 to 100; 50 by default
}

type blueprintLoaderState struct {
	start time.Time
}

// blueprintPart is one piece of the profile card, measured off the card itself.
type blueprintPart struct {
	r     ui.Rect
	round bool // a circle or a pill, not a box
	label string
}

var blueprintParts = []blueprintPart{
	{ui.Rect{X: 16, Y: 16, W: 56, H: 56}, true, "56"},
	{ui.Rect{X: 212, Y: 16, W: 72, H: 32}, true, "72×32"},
	{ui.Rect{X: 16, Y: 84, W: 140, H: 22}, false, ""},
	{ui.Rect{X: 16, Y: 106, W: 72, H: 18}, false, ""},
	{ui.Rect{X: 16, Y: 134, W: 238, H: 40}, false, "238"},
	{ui.Rect{X: 16, Y: 184, W: 206, H: 20}, false, ""},
	{ui.Rect{X: 16, Y: 216, W: 36, H: 20}, true, ""},
	{ui.Rect{X: 60, Y: 217, W: 142, H: 18}, false, ""},
}

// blueprintStroke draws the first k of a box's outline, clockwise from its top left.
func blueprintStroke(p *ui.Painter, r ui.Rect, k float32, col ui.Color) {
	if k <= 0 {
		return
	}
	pts := [][2]float32{{r.X, r.Y}, {r.X + r.W, r.Y}, {r.X + r.W, r.Y + r.H}, {r.X, r.Y + r.H}, {r.X, r.Y}}
	total := 2 * (r.W + r.H)
	left := k * total
	for i := 0; i < 4 && left > 0; i++ {
		a, b := pts[i], pts[i+1]
		seg := r.W
		if i%2 == 1 {
			seg = r.H
		}
		f := min(1, left/seg)
		p.Line(a[0], a[1], a[0]+(b[0]-a[0])*f, a[1]+(b[1]-a[1])*f, 1, col)
		left -= seg
	}
}

// BlueprintLoader draws the plan of a profile card on a dot grid: rules across
// at each part's top, then each part's outline with its measure, and finally
// blends the finished card in over it. Pressing draws it again.
func BlueprintLoader(c *ui.Context, opts ...BlueprintLoaderOptions) ui.Element {
	var o BlueprintLoaderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Speed == 0 {
		o.Speed = 50
	}
	t := ThemeFrom(c)
	card := ui.Box(c).Size(300, 252).Radius(20).Background(t.Surface).Clip().Cursor(ui.CursorPointer)
	st := ui.Local(card, "blueprint_loader", func() blueprintLoaderState { return blueprintLoaderState{start: c.Now()} })
	card.OnClick(func() { st.start = c.Now() })
	pace := lerp32(1.8, 0.6, clamp32(o.Speed, 0, 100)/100)
	since := float32(c.Now().Sub(st.start).Seconds()) / pace
	const plan, blend = 1.6, 0.7 // in paces
	if since < plan+blend {
		c.AnimationFrame()
	}
	in := clamp32((since-plan)/blend, 0, 1)
	in = in * in * (3 - 2*in)
	ink := ui.Hex("#3b82f6")
	if t.Dark {
		ink = ui.Hex("#60a5fa")
	}

	card.Children(func() {
		// The plan.
		ui.Box(c).Key("plan").Absolute().Left(0).Top(0).Size(300, 252).Opacity(1 - in).Draw(func(p *ui.Painter, r ui.Rect) {
			for y := r.Y + 6; y < r.Y+r.H; y += 12 {
				for x := r.X + 6; x < r.X+r.W; x += 12 {
					p.Fill(ui.Rect{X: x - 0.75, Y: y - 0.75, W: 1.5, H: 1.5}, t.BorderStrong, 0.75)
				}
			}
			for i, part := range blueprintParts {
				k := clamp32((since-float32(i)*0.12)/0.5, 0, 1)
				pr := ui.Rect{X: r.X + part.r.X, Y: r.Y + part.r.Y, W: part.r.W, H: part.r.H}
				// The rule across the card at the part's top.
				p.Line(r.X, pr.Y, r.X+r.W*min(1, k*1.6), pr.Y, 1, ink.Alpha(0.18))
				if part.round {
					if k > 0 {
						p.Stroke(pr, ink.Alpha(0.85*k), min(pr.W, pr.H)/2, 1)
					}
				} else {
					p.Fill(pr, ink.Alpha(0.06*k), 3)
					blueprintStroke(p, pr, k, ink.Alpha(0.85))
				}
				if part.label != "" && k > 0.6 {
					span := ui.Span{Text: part.label, Size: 9, Color: ink.Alpha((k - 0.6) / 0.4)}
					w, _ := p.MeasureText(0, span)
					p.RichText(pr.X+pr.W/2-w/2, pr.Y+pr.H+2, w+2, span)
				}
			}
		})
		// The finished card, blending in.
		ui.Column(c).Key("card").Absolute().Left(0).Top(0).Size(300, 252).Padding(16).Opacity(in).Children(func() {
			ui.Row(c).AlignItems(ui.Start).Children(func() {
				ui.Box(c).Size(56, 56).Draw(func(p *ui.Painter, r ui.Rect) {
					p.FillGradient(r, ui.LinearGradient{From: ui.Hex("#36b0ff"), To: ui.Hex("#5b45dd"), Angle: 150}, 28)
				})
				ui.Box(c).Grow(1)
				ui.Box(c).Size(72, 32).Radius(16).Background(t.Text).Center().Children(func() {
					ui.Text(c, "Follow").FontSize(14).FontWeight(500).TextColor(t.Surface)
				})
			})
			ui.Text(c, "Mara Lindqvist").FontSize(17).FontWeight(500).TextColor(t.Text).Margin(10, 0, 0, 0)
			ui.Text(c, "@maralind").FontSize(14).TextColor(t.TextSubtle).Margin(1, 0, 0, 0)
			ui.Text(c, "Designing motion at a small studio in Lisbon. Springs over curves, always.").
				FontSize(14.5).LineHeight(1.38).TextColor(t.TextMuted).Margin(10, 0, 0, 0)
			ui.Row(c).Gap(16).Margin(10, 0, 0, 0).Children(func() {
				for _, s := range [][2]string{{"312", "Following"}, {"1,284", "Followers"}} {
					ui.Row(c).Key(s[1]).Gap(4).Children(func() {
						ui.Text(c, s[0]).FontSize(14).TextColor(t.Text)
						ui.Text(c, s[1]).FontSize(14).TextColor(t.TextSubtle)
					})
				}
			})
			ui.Row(c).AlignItems(ui.Center).Margin(10, 0, 0, 0).Children(func() {
				ui.Box(c).Size(36, 20).Draw(func(p *ui.Painter, r ui.Rect) {
					p.FillGradient(ui.Rect{X: r.X, Y: r.Y, W: 20, H: 20}, ui.LinearGradient{From: ui.Hex("#c77dff"), To: ui.Hex("#ff5d8f"), Angle: 150}, 10)
					p.Fill(ui.Rect{X: r.X + 13, Y: r.Y - 1.5, W: 23, H: 23}, t.Surface, 11.5)
					p.FillGradient(ui.Rect{X: r.X + 14.5, Y: r.Y, W: 20, H: 20}, ui.LinearGradient{From: ui.Hex("#ff6b6b"), To: ui.Hex("#ff9f68"), Angle: 150}, 10)
				})
				ui.Text(c, fmt.Sprintf("Followed by %s and %s", "Kai", "Inês")).FontSize(13.5).TextColor(t.TextSubtle).Margin(0, 0, 0, 8)
			})
		})
	})
	return card
}

func init() {
	registerDemo(
		"blueprint-loader",
		"Draws the plan of a profile card on a dot grid, rules and measured boxes, then blends the finished card in.",
		func(c *ui.Context) { BlueprintLoader(c) },
	)
}
