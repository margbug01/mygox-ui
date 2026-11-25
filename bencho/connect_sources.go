package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// ConnectSource is one tool on the Integrations card.
type ConnectSource struct {
	Name      string
	Color     ui.Color
	Synced    float32   // thousands of items, shown as "4.1k synced"
	Bars      []float32 // the activity shape, 0 to 1, one per bar
	Connected bool      // its starting state
}

type connectSourcesState struct {
	on      []bool
	changed []time.Time
}

const (
	connectSourcesW    = 300
	connectSourcesTile = 130
	connectSourcesTall = 100
	connectSourcesGap  = 8
)

// ConnectSourcesDefault is the card bencho.dev shows.
func ConnectSourcesDefault() []ConnectSource {
	return []ConnectSource{
		{Name: "Mail", Color: ui.Hex("#9b7bea"), Synced: 4.1, Connected: true, Bars: []float32{.15, .15, .25, .35, .55, .8, .95, .95, .75, .5, .3, .2, .15, .15, .15, .15}},
		{Name: "Calendar", Color: ui.Hex("#4d86e6"), Synced: 3.9, Connected: true, Bars: []float32{.25, .25, .3, .4, .6, .75, .85, .95, 1, .95, .9, .8, .7, .55, .4, .3}},
		{Name: "Files", Color: ui.Hex("#97b2d6"), Synced: 6.2, Connected: true, Bars: []float32{.1, .15, .15, .15, .15, .15, .15, .2, .3, .55, .8, 1, .95, .7, .4, .3}},
		{Name: "Chat", Color: ui.Hex("#efc069"), Synced: 2.4, Connected: true, Bars: []float32{.15, .15, .2, .25, .4, .65, .95, .9, .65, .4, .25, .15, .15, .15, .15, .15}},
		{Name: "Docs", Color: ui.Hex("#6fd8a2"), Synced: 1.8, Bars: []float32{.4, .7, .95, .95, .65, .4, .2, .15, .15, .15, .15, .15, .15, .15, .15, .15}},
		{Name: "Tasks", Color: ui.Hex("#e27ab9"), Synced: 3.3, Bars: []float32{.15, .15, .2, .15, .2, .3, .55, .85, 1, .85, .55, .3, .15, .15, .15, .15}},
	}
}

// ConnectSources is an Integrations card: a tile per tool with its live
// activity bars and sync count, a count of connected tools and a setup bar.
// The plus on a quiet tile connects it, its bars rising one after another
// as the count climbs; the dot on a connected tile disconnects it.
func ConnectSources(c *ui.Context, sources []ConnectSource) ui.Element {
	t := ThemeFrom(c)
	card := ui.Column(c).Width(connectSourcesW).Padding(16).Radius(24).Background(t.Surface)
	st := ui.Local(card, "connect_sources", func() connectSourcesState {
		s := connectSourcesState{on: make([]bool, len(sources)), changed: make([]time.Time, len(sources))}
		for i, src := range sources {
			s.on[i] = src.Connected
		}
		return s
	})
	now := c.Now()
	n := 0
	for _, on := range st.on {
		n += int(b2f(on))
	}
	quiet := t.BorderStrong
	if t.Dark {
		quiet = ui.RGBA(255, 255, 255, 0.14)
	}

	card.Children(func() {
		ui.Row(c).AlignItems(ui.Start).Children(func() {
			ui.Column(c).Grow(1).Gap(3).Children(func() {
				ui.Text(c, "Integrations").FontSize(15).FontWeight(500).TextColor(t.Text)
				ui.Text(c, "Connect your tools").FontSize(12.5).TextColor(t.TextSubtle)
			})
			ui.Box(c).Height(28).Padding(0, 10).Radius(10).Background(t.SurfaceSubtle).Center().Children(func() {
				ui.Text(c, fmt.Sprintf("%d of %d", n, len(sources))).FontSize(13).TextColor(t.Text)
			})
		})

		ui.Row(c).Wrap().Gap(connectSourcesGap).Margin(14, 0, 0, 0).Children(func() {
			for i, src := range sources {
				on := st.on[i]
				tile := ui.Column(c).Key(i).Size(connectSourcesTile, connectSourcesTall).Padding(12).Radius(14).
					Border(1, t.Border)
				k := clamp32(tile.Spring("on", b2f(on), ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: 0.2}), 0, 1.15)
				since := float32(now.Sub(st.changed[i])) / float32(time.Second)
				if !st.changed[i].IsZero() && since < 1.4 {
					c.AnimationFrame()
				}
				ink := lerpColor(t.TextSubtle, t.Text, clamp32(k, 0, 1))
				tile.Children(func() {
					ui.Row(c).Height(20).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Size(7, 7).Radius(1.5).Background(src.Color).Margin(0, 9, 0, 0)
						ui.Text(c, src.Name).FontSize(13.5).TextColor(ink).Grow(1).SingleLine()
						btn := ui.Box(c).Key("status").Size(20, 20).Radius(10).Center().Cursor(ui.CursorPointer).OnClick(func() {
							st.on[i] = !on
							st.changed[i] = c.Now()
						})
						hover := btn.Spring("hover", b2f(btn.Hovered()), ui.SpringOptions{Duration: 140 * time.Millisecond})
						btn.Background(lerpColor(t.SurfaceHover, src.Color.Alpha(0.18), clamp32(k, 0, 1))).
							Scale(1+0.08*hover, 1+0.08*hover)
						btn.Children(func() {
							if k > 0.5 {
								d := 8 * clamp32(k, 0, 1.15)
								ui.Box(c).Key("dot").Size(d, d).Radius(d / 2).Background(src.Color)
							} else {
								ui.Icon(c, icons.Plus()).Key("plus").Size(14, 14).TextColor(t.TextMuted).Opacity(1 - k*2)
							}
						})
					})
					// The bars: quiet grey until connected, then each rises in turn in the tool's colour.
					ui.Box(c).Key("bars").Size(connectSourcesTile-24, 26).Margin(10, 0, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
						step := r.W / float32(len(src.Bars))
						for j, h := range src.Bars {
							bk := k
							if !st.changed[i].IsZero() && on {
								bk = clamp32((since-float32(j)*0.03)/0.4, 0, 1)
							}
							live := float32(1)
							if on {
								live = 1 + 0.08*float32(math.Sin(phaseSeconds(now)*2.2+float64(j)*0.6))
							}
							bh := max(4, r.H*h*live)
							col := lerpColor(quiet, src.Color, bk)
							p.Fill(ui.Rect{X: r.X + float32(j)*step + (step-3)/2, Y: r.Y + r.H - bh, W: 3, H: bh}, col, 1.5)
						}
						if on {
							p.AnimationFrame()
						}
					})
					label := "Not connected"
					if on || k > 0.5 {
						count := src.Synced
						if !st.changed[i].IsZero() && on {
							count *= clamp32(since/1.1, 0, 1)
						}
						label = fmt.Sprintf("%.1fk synced", count)
					}
					ui.Text(c, label).FontSize(12.5).TextColor(t.TextSubtle).Margin(8, 0, 0, 0).SingleLine()
				})
			}
		})

		ui.Row(c).Height(28).AlignItems(ui.Center).Margin(8, 0, 0, 0).Children(func() {
			ui.Text(c, "Setup").FontSize(12.5).TextColor(t.TextSubtle).Margin(0, 12, 0, 0)
			for i, src := range sources {
				seg := ui.Box(c).Key(i).Size(26, 4).Radius(2).Margin(0, 4, 0, 0)
				k := clamp32(seg.Spring("on", b2f(st.on[i]), ui.SpringOptions{Duration: 400 * time.Millisecond}), 0, 1)
				seg.Background(lerpColor(t.SurfaceHover, src.Color, k))
			}
			ui.Box(c).Grow(1)
			pct := card.Spring("pct", float32(52+8*n), ui.SpringOptions{Duration: 700 * time.Millisecond})
			ui.Text(c, fmt.Sprintf("%d%%", int(pct+0.5))).FontSize(13.5).FontWeight(500).TextColor(t.Text)
		})
	})
	return card
}

func init() {
	registerDemo(
		"connect-sources",
		"An Integrations card of tool tiles with live activity bars; the plus connects a tool as its bars rise and the setup bar fills.",
		func(c *ui.Context) { ConnectSources(c, ConnectSourcesDefault()) },
	)
}
