package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// LoadCard is one result in LoadCards.
type LoadCard struct {
	Title, Source string
	From, To      ui.Color // the source's dot
}

// LoadCardsOptions configures LoadCards.
type LoadCardsOptions struct {
	Pace time.Duration // time between one card developing and the next, 120 to 800ms; 380ms by default
	Blur float32       // how soft a layer is as it comes and goes, 0 to 14px; 6 by default
}

type loadCardsState struct {
	start time.Time
}

const (
	loadCardsW    = 320
	loadCardsCard = 88
)

// LoadCardsDefault is the search results bencho.dev shows.
func LoadCardsDefault() []LoadCard {
	return []LoadCard{
		{Title: "Springs that keep their velocity", Source: "motion", From: ui.Hex("#9b6bff"), To: ui.Hex("#6b45f0")},
		{Title: "Reduced motion, done properly", Source: "web.dev", From: ui.Hex("#6cc9ff"), To: ui.Hex("#3aa6f2")},
		{Title: "Layout animation with FLIP", Source: "css-tricks", From: ui.Hex("#ff63c8"), To: ui.Hex("#f0508a")},
	}
}

// loadCardsSoft draws a layer that is blur px out of focus: faded, with ghost
// copies spread around it.
func loadCardsSoft(c *ui.Context, key any, blur float32, alpha float32, build func(ghost bool, a float32)) {
	ui.Box(c).Key(key).Absolute().Left(0).Top(0).WidthPercent(100).HeightPercent(100).Opacity(alpha).Children(func() {
		build(false, 1-0.5*clamp32(blur/8, 0, 1))
		if blur < 0.3 {
			return
		}
		for g := 0; g < 6; g++ {
			a := float64(g) * math.Pi / 3
			ui.Box(c).Key(g).Absolute().Left(float32(math.Cos(a)) * blur * 0.6).Top(float32(math.Sin(a)) * blur * 0.6).
				WidthPercent(100).HeightPercent(100).Children(func() { build(true, 0.18) })
		}
	})
}

// LoadCards brings in three result cards one after another. Each arrives as a
// skeleton easing out of a blur, a light running over its bars; then the bars
// blur away as the real content comes into focus. Pressing them loads again.
func LoadCards(c *ui.Context, cards []LoadCard, opts ...LoadCardsOptions) ui.Element {
	var o LoadCardsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Pace == 0 {
		o.Pace = 380 * time.Millisecond
	}
	if o.Blur == 0 {
		o.Blur = 6
	}
	t := ThemeFrom(c)
	list := ui.Column(c).Width(loadCardsW).Gap(8).Cursor(ui.CursorPointer)
	st := ui.Local(list, "load_cards", func() loadCardsState { return loadCardsState{start: c.Now()} })
	list.OnClick(func() { st.start = c.Now() })
	now := c.Now()
	since := float32(now.Sub(st.start).Seconds())
	pace := float32(o.Pace.Seconds())
	const arrive, hold, develop = 0.45, 0.9, 0.5 // seconds
	if since < pace*float32(len(cards))+arrive+hold+develop {
		c.AnimationFrame()
	}
	bone := lerpColor(t.SurfaceSubtle, t.SurfaceHover, 0.45)
	shine := float32(math.Mod(float64(since)*1.3, 1.6)) - 0.3

	list.Children(func() {
		for i, cd := range cards {
			at := since - float32(i)*pace
			in := clamp32(at/arrive, 0, 1)          // the skeleton coming out of the blur
			dev := clamp32((at-hold)/develop, 0, 1) // the content coming into focus
			if in <= 0 {
				ui.Box(c).Key(i).Height(loadCardsCard)
				continue
			}
			card := ui.Box(c).Key(i).Height(loadCardsCard).Radius(20).Background(t.Surface).Clip().
				Opacity(clamp32(in*1.5, 0, 1)).Translate(0, (1-in)*10)
			card.Children(func() {
				if dev < 1 {
					loadCardsSoft(c, "bones", o.Blur*(1-in+dev), 1-dev, func(ghost bool, a float32) {
						ui.Box(c).Absolute().Left(16).Top(22).Size(220, 12).Radius(6).Background(bone.Alpha(a))
						ui.Box(c).Absolute().Left(16).Top(53).Size(20, 20).Radius(10).Background(bone.Alpha(a))
						ui.Box(c).Absolute().Left(44).Top(57).Size(70, 12).Radius(6).Background(bone.Alpha(a))
						if !ghost {
							// A light running across the bars.
							ui.Box(c).Absolute().Left(0).Top(0).WidthPercent(100).HeightPercent(100).Draw(func(p *ui.Painter, r ui.Rect) {
								x := r.X + r.W*shine
								p.FillGradient(ui.Rect{X: x - 40, Y: r.Y, W: 80, H: r.H}, ui.LinearGradient{
									From: t.Surface.Alpha(0), To: t.Surface.Alpha(0.55), Angle: 90, Start: 0, End: 0.5,
								}, 0)
								p.FillGradient(ui.Rect{X: x, Y: r.Y, W: 40, H: r.H}, ui.LinearGradient{
									From: t.Surface.Alpha(0.55), To: t.Surface.Alpha(0), Angle: 90,
								}, 0)
							})
						}
					})
				}
				if dev > 0 {
					loadCardsSoft(c, "content", o.Blur*(1-dev), dev, func(ghost bool, a float32) {
						ui.Text(c, cd.Title).Absolute().Left(16).Top(17).FontSize(15.5).FontWeight(500).TextColor(t.Text.Alpha(a)).SingleLine()
						ui.Box(c).Absolute().Left(16).Top(53).Size(20, 20).Draw(func(p *ui.Painter, r ui.Rect) {
							p.FillGradient(r, ui.LinearGradient{From: cd.From.Alpha(a), To: cd.To.Alpha(a), Angle: 150}, 10)
						})
						ui.Text(c, cd.Source).Absolute().Left(44).Top(53).FontSize(14).TextColor(t.TextSubtle.Alpha(a)).SingleLine()
					})
				}
			})
		}
	})
	return list
}

func init() {
	registerDemo(
		"load-cards",
		"Three result cards that arrive as skeletons out of a blur and develop into their content one after another.",
		func(c *ui.Context) { LoadCards(c, LoadCardsDefault()) },
	)
}
