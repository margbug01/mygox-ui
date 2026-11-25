package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// DotSwarmSource is one row: an alert source and its activity through the day.
type DotSwarmSource struct {
	Name  string
	Color ui.Color
	Peaks [][3]float32 // gaussian bumps: hour, spread in hours, height 0 to 1
}

// DotSwarmOptions configures DotSwarm.
type DotSwarmOptions struct {
	Mono    bool // all rows in ink
	Density int  // how tightly the dots are packed, 1 to 5; 4 by default
	Rows    int  // how many sources are shown, 2 to 5; 4 by default
}

type dotSwarmState struct {
	px   float32 // the pointer over the grid last frame
	over bool
}

const (
	dotSwarmW     = 300
	dotSwarmLines = 7 // dot lines per row: a middle line and three either side
)

// DotSwarmDefault is bencho.dev's week of traffic.
func DotSwarmDefault() []DotSwarmSource {
	return []DotSwarmSource{
		{"Web", ui.Hex("#9b6cf0"), [][3]float32{{12, 2.6, 1}, {16.5, 1.5, 0.35}}},
		{"Mobile", ui.Hex("#4c82e8"), [][3]float32{{12.5, 3.2, 0.9}, {9, 1.2, 0.4}, {19.5, 1.4, 0.35}}},
		{"Email", ui.Hex("#f0bd62"), [][3]float32{{11.5, 1.6, 1}, {9.5, 1, 0.4}}},
		{"API", ui.Hex("#6fd7ad"), [][3]float32{{2.5, 1.1, 1}, {23.6, 0.8, 0.6}}},
		{"Push", ui.Hex("#ff7a59"), [][3]float32{{8, 1.5, 0.6}, {20, 1.6, 0.8}}},
	}
}

func (s DotSwarmSource) at(hour float32) float32 {
	v := float32(0)
	for _, p := range s.Peaks {
		d := (hour - p[0]) / p[1]
		v += p[2] * float32(math.Exp(float64(-d*d/2)))
	}
	return clamp32(v, 0, 1)
}

// DotSwarm is activity by hour drawn as swarms of dots round a line, one row
// per source. Running the pointer along it swells the column under it and
// reads out that half hour; at rest it marks the peak.
func DotSwarm(c *ui.Context, sources []DotSwarmSource, opts ...DotSwarmOptions) ui.Element {
	var o DotSwarmOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Density == 0 {
		o.Density = 4
	}
	if o.Rows == 0 {
		o.Rows = 4
	}
	rows := sources[:min(max(o.Rows, 2), len(sources))]
	t := ThemeFrom(c)
	cols := 25 + 6*min(max(o.Density, 1), 5) // 49 at density 4: every half hour
	pitch := float32(dotSwarmW-32) / float32(cols-1)
	dot := pitch * 0.62
	rowH := float32(dotSwarmLines-1)*6 + 14

	// The peak half hour, and the visits in a column.
	visits := func(col int) int {
		h := float32(col) / float32(cols-1) * 24
		n := float32(0)
		for _, s := range rows {
			n += s.at(h)
		}
		return int(n * 180)
	}
	peak := 0
	for i := range cols {
		if visits(i) > visits(peak) {
			peak = i
		}
	}

	card := ui.Column(c).Width(dotSwarmW).Padding(16).Radius(24).Background(t.Surface)
	st := ui.Local(card, "dot_swarm", func() dotSwarmState { return dotSwarmState{} })
	target := float32(peak)
	if st.over {
		target = clamp32(float32(math.Round(float64(st.px/pitch))), 0, float32(cols-1))
	}
	col := card.Spring("col", target, ui.SpringOptions{Duration: 260 * time.Millisecond})
	at := int(math.Round(float64(col)))
	hour := float32(at) / float32(cols-1) * 24

	card.Children(func() {
		ui.Row(c).AlignItems(ui.Start).Children(func() {
			ui.Column(c).Grow(1).Gap(2).Children(func() {
				ui.Text(c, "Activity by hour").FontSize(15.5).FontWeight(500).TextColor(t.Text)
				ui.Text(c, "Last 7 days").FontSize(13).TextColor(t.TextSubtle)
			})
			ui.Box(c).Height(28).Padding(0, 11).Radius(14).Background(t.SurfaceSubtle).Center().Children(func() {
				ui.Text(c, "UTC").FontSize(12.5).FontWeight(500).TextColor(t.TextMuted)
			})
		})
		// The readout over the scrubbed column.
		ui.Box(c).Key("readout").Height(36).Margin(12, 0, 0, 0).Children(func() {
			clock := fmt.Sprintf("%02d:%02d", int(hour), int(math.Mod(float64(hour), 1)*60+0.5))
			tip := ui.Row(c).Absolute().Top(0).Height(24).Padding(0, 9).Gap(4).Radius(12).Border(1, t.BorderStrong).
				Background(t.Surface).AlignItems(ui.Center)
			w := tip.LayoutBounds().W
			tip.Left(clamp32(col*pitch-w/2, 0, dotSwarmW-32-w)).Children(func() {
				if at == peak {
					ui.Text(c, "Peak").FontSize(12).TextColor(t.Text)
				}
				ui.Text(c, clock).FontSize(12).TextColor(t.Text)
				ui.Text(c, "·").FontSize(12).TextColor(t.TextSubtle)
				ui.Text(c, fmt.Sprintf("%d visits", visits(at))).FontSize(12).TextColor(t.Text)
			})
		})
		grid := ui.Box(c).Key("grid")
		st.px, _, st.over = grid.PointerPosition()
		grid.Height(float32(len(rows))*rowH + float32(len(rows)-1)*8).Draw(func(p *ui.Painter, r ui.Rect) {
			gx := r.X + col*pitch
			p.Fill(ui.Rect{X: gx - 2.5, Y: r.Y - 12, W: 5, H: 5}, t.Text, 2.5)
			for y := r.Y - 6; y < r.Y+r.H; y += 4 {
				p.Fill(ui.Rect{X: gx - 0.5, Y: y, W: 1, H: 2}, t.TextSubtle.Alpha(0.6), 0)
			}
			idle := lerpColor(t.SurfaceSubtle, t.SurfaceHover, 0.55)
			for ri, s := range rows {
				top := r.Y + float32(ri)*(rowH+8)
				ink := s.Color
				if o.Mono {
					ink = t.Text
				}
				for i := 0; i < cols; i++ {
					v := s.at(float32(i) / float32(cols-1) * 24)
					reach := int(math.Round(float64(v * 3.4))) // dots either side of the middle line
					swell := clamp32(1-float32(math.Abs(float64(float32(i)-col)))/1.6, 0, 1)
					for l := 0; l < dotSwarmLines; l++ {
						off := l - dotSwarmLines/2
						on := off == 0 || v > 0.04 && dotSwarmAbs(off) <= reach
						c := idle
						if on {
							c = ink
						}
						d := dot * (1 + 0.5*swell)
						cx, cy := r.X+float32(i)*pitch, top+7+float32(l)*6
						p.Fill(ui.Rect{X: cx - d/2, Y: cy - d/2, W: d, H: d}, c, d/2)
					}
				}
			}
		})
		ui.Row(c).Margin(10, 0, 0, 0).Children(func() {
			for i, h := range []string{"00", "06", "12", "18", "24"} {
				if i > 0 {
					ui.Box(c).Grow(1)
				}
				ui.Text(c, h).Key(h).FontSize(12).TextColor(t.TextSubtle)
			}
		})
		ui.Row(c).Wrap().Margin(16, 0, 0, 0).Children(func() {
			total := float32(0)
			shares := make([]float32, len(rows))
			for i, s := range rows {
				for h := float32(0); h < 24; h += 0.5 {
					shares[i] += s.at(h)
				}
				total += shares[i]
			}
			for i, s := range rows {
				ink := s.Color
				if o.Mono {
					ink = t.Text
				}
				ui.Row(c).Key(s.Name).Width((dotSwarmW - 32) / 2).Height(22).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(8, 8).Radius(2).Background(ink).Margin(0, 7, 0, 0)
					ui.Text(c, s.Name).FontSize(13.5).TextColor(t.TextSubtle).Margin(0, 5, 0, 0)
					ui.Text(c, fmt.Sprintf("%d%%", int(shares[i]/total*100*0.78+0.5))).FontSize(13.5).FontWeight(500).TextColor(t.Text)
				})
			}
		})
	})
	return card
}

func dotSwarmAbs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func init() {
	registerDemo(
		"dot-swarm",
		"Activity by hour as swarms of dots round a line, one row per source; run along it to read each half hour.",
		func(c *ui.Context) { DotSwarm(c, DotSwarmDefault()) },
	)
}
