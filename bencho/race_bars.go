package bencho

import (
	"fmt"
	"sort"
	"time"

	"github.com/egoist/mygo/ui"
)

// RaceBarItem is one contender: a yearly series from FirstYear on.
type RaceBarItem struct {
	Label  string // its badge shows the first letter
	Color  ui.Color
	Values []float32 // one a year
}

// RaceBarsData is a race to show.
type RaceBarsData struct {
	Title     string
	FirstYear int
	Items     []RaceBarItem
	Format    func(v float32) string
	Total     func(v float32) string // the combined figure; Format when nil
}

// RaceBarsOptions configures RaceBars.
type RaceBarsOptions struct {
	Mono bool    // grey bars with the leader in ink
	Pace float32 // how fast the months go by, 0.4 to 2.5; 1 by default
}

type raceBarsState struct {
	t       float32 // 0 to 1 across the years
	playing bool
	last    time.Time
	ranks   map[string]float32
}

const (
	raceBarsW    = 360
	raceBarsRow  = 32
	raceBarsBarH = 26
	raceBarsMaxW = 222
	raceBarsYear = 1.0 // seconds per year at pace 1
)

// RaceBarsMarketCap is the biggest companies by market cap, in trillions of dollars.
func RaceBarsMarketCap() RaceBarsData {
	return RaceBarsData{
		Title: "Market cap", FirstYear: 2016,
		Format: func(v float32) string { return fmt.Sprintf("$%.2fT", v) },
		Total:  func(v float32) string { return fmt.Sprintf("$%.1fT", v) },
		Items: []RaceBarItem{
			{"Apple", ui.Hex("#2d6bff"), []float32{0.61, 0.86, 0.75, 1.29, 2.25, 2.91, 2.07, 2.99, 3.79, 3.85, 4.98}},
			{"Nvidia", ui.Hex("#16c25c"), []float32{0.06, 0.12, 0.08, 0.14, 0.32, 0.73, 0.36, 1.22, 3.29, 4.4, 4.5}},
			{"Microsoft", ui.Hex("#00b4d8"), []float32{0.48, 0.66, 0.78, 1.2, 1.68, 2.52, 1.79, 2.79, 3.13, 3.6, 3.94}},
			{"Google", ui.Hex("#ffb703"), []float32{0.53, 0.73, 0.72, 0.92, 1.18, 1.92, 1.15, 1.75, 2.32, 2.9, 3.47}},
			{"SpaceX", ui.Hex("#ff3f96"), []float32{0.01, 0.02, 0.03, 0.04, 0.07, 0.1, 0.13, 0.18, 0.35, 0.8, 1.82}},
			{"Meta", ui.Hex("#9050ff"), []float32{0.33, 0.51, 0.38, 0.59, 0.78, 0.92, 0.32, 0.91, 1.48, 1.7, 1.56}},
			{"Tesla", ui.Hex("#ff3b30"), []float32{0.03, 0.05, 0.05, 0.08, 0.67, 1.06, 0.39, 0.79, 1.29, 1.1, 1.42}},
			{"Broadcom", ui.Hex("#ff8a1f"), []float32{0.07, 0.1, 0.1, 0.13, 0.18, 0.27, 0.23, 0.52, 1.08, 1.3, 1.4}},
		},
	}
}

// RaceBarsPopulation is the most populous countries, in people.
func RaceBarsPopulation() RaceBarsData {
	lerp := func(a, b float32) []float32 {
		v := make([]float32, 11)
		for i := range v {
			v[i] = a + (b-a)*float32(i)/10
		}
		return v
	}
	china := lerp(1.379, 1.412)
	for i := 7; i < 11; i++ {
		china[i] = 1.412 - float32(i-6)*0.003
	}
	return RaceBarsData{
		Title: "Population", FirstYear: 2016,
		Format: func(v float32) string {
			if v >= 1 {
				return fmt.Sprintf("%.2fB", v)
			}
			return fmt.Sprintf("%.0fM", v*1000)
		},
		Items: []RaceBarItem{
			{"India", ui.Hex("#ff9933"), lerp(1.324, 1.46)},
			{"China", ui.Hex("#de2910"), china},
			{"United States", ui.Hex("#3c3b6e"), lerp(0.323, 0.345)},
			{"Indonesia", ui.Hex("#ce1126"), lerp(0.261, 0.285)},
			{"Pakistan", ui.Hex("#01411c"), lerp(0.203, 0.255)},
			{"Nigeria", ui.Hex("#008751"), lerp(0.19, 0.24)},
			{"Brazil", ui.Hex("#009c3b"), lerp(0.206, 0.213)},
			{"Bangladesh", ui.Hex("#006a4e"), lerp(0.161, 0.176)},
		},
	}
}

// at reads a yearly series at t, from 0 (the first year) to 1 (the last).
func (it RaceBarItem) at(t float32) float32 {
	f := t * float32(len(it.Values)-1)
	i := int(f)
	if i >= len(it.Values)-1 {
		return it.Values[len(it.Values)-1]
	}
	k := f - float32(i)
	return it.Values[i] + (it.Values[i+1]-it.Values[i])*k
}

// RaceBars is an animated bar chart race. Bars grow, the year counts up and
// rows swap places as one overtakes another; play runs it from the start,
// and dragging the timeline scrubs through the years.
func RaceBars(c *ui.Context, data RaceBarsData, opts ...RaceBarsOptions) ui.Element {
	var o RaceBarsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Pace == 0 {
		o.Pace = 1
	}
	t := ThemeFrom(c)
	card := ui.Column(c).Width(raceBarsW).Padding(20, 20, 16, 20).Radius(24).Background(t.Surface)
	st := ui.Local(card, "race_bars", func() raceBarsState { return raceBarsState{t: 1, ranks: map[string]float32{}} })
	now := c.Now()
	years := float32(len(data.Items[0].Values) - 1)
	if st.playing {
		dt := float32(now.Sub(st.last).Seconds())
		st.t += dt * clamp32(o.Pace, 0.4, 2.5) / (years * raceBarsYear)
		if st.t >= 1 {
			st.t, st.playing = 1, false
		}
		c.AnimationFrame()
	}
	st.last = now

	type row struct {
		item RaceBarItem
		v    float32
	}
	rows := make([]row, len(data.Items))
	sum, top := float32(0), float32(0)
	for i, it := range data.Items {
		rows[i] = row{it, it.at(st.t)}
		sum += rows[i].v
		top = max(top, rows[i].v)
	}
	order := append([]row(nil), rows...)
	sort.SliceStable(order, func(a, b int) bool { return order[a].v > order[b].v })
	rank := map[string]int{}
	for i, r := range order {
		rank[r.item.Label] = i
	}
	year := data.FirstYear + int(st.t*years+0.0001)

	card.Children(func() {
		ui.Row(c).AlignItems(ui.Start).Children(func() {
			ui.Column(c).Grow(1).Gap(2).Children(func() {
				ui.Text(c, data.Title).FontSize(17).FontWeight(500).TextColor(t.Text)
				ui.Row(c).Gap(4).Children(func() {
					ui.Text(c, "Combined").FontSize(14).TextColor(t.TextSubtle)
					total := data.Format
					if data.Total != nil {
						total = data.Total
					}
					ui.Text(c, total(sum)).FontSize(14).TextColor(t.Text)
				})
			})
			ui.Text(c, fmt.Sprint(year)).FontSize(30).FontWeight(500).TextColor(t.Text)
		})
		ui.Box(c).Key("rows").Height(float32(len(rows))*raceBarsRow).Margin(16, 0, 0, 0).Children(func() {
			for _, r := range rows {
				el := ui.Row(c).Key(r.item.Label).Absolute().Left(0).Height(raceBarsBarH).AlignItems(ui.Center)
				y := el.Spring("rank", float32(rank[r.item.Label]), ui.SpringOptions{Duration: 420 * time.Millisecond})
				el.Top(y * raceBarsRow)
				col := r.item.Color
				if o.Mono {
					col = lerpColor(t.SurfaceHover, t.BorderStrong, 0.6)
					if rank[r.item.Label] == 0 {
						col = t.Text
					}
				}
				w := max(8, raceBarsMaxW*r.v/max(top, 0.0001))
				el.Children(func() {
					ui.Box(c).Size(raceBarsBarH, raceBarsBarH).Radius(raceBarsBarH/2).Background(col).Border(1, t.Border).Center().Children(func() {
						ui.Text(c, r.item.Label[:1]).FontSize(12).FontWeight(600).TextColor(ui.Hex("#ffffff"))
					})
					ui.Box(c).Size(w, raceBarsBarH).Radius(8).Background(col).Margin(0, 0, 0, 6)
					ui.Text(c, data.Format(r.v)).FontSize(14).FontWeight(500).TextColor(t.Text).Margin(0, 0, 0, 8)
				})
			}
		})
		ui.Row(c).Height(36).AlignItems(ui.Center).Margin(14, 0, 0, 0).Children(func() {
			track := ui.Box(c).Key("track").Grow(1).Height(20).Cursor(ui.CursorPointer)
			if _, _, dragging := track.Dragged(); dragging || track.Pressed() {
				x, _, _ := track.PointerPosition()
				if w := track.LayoutBounds().W; w > 0 {
					st.t, st.playing = clamp32(x/w, 0, 1), false
				}
			}
			track.Draw(func(p *ui.Painter, r ui.Rect) {
				y := r.Y + r.H/2 - 1.5
				p.Fill(ui.Rect{X: r.X, Y: y, W: r.W, H: 3}, t.SurfaceHover, 1.5)
				p.Fill(ui.Rect{X: r.X, Y: y, W: r.W * st.t, H: 3}, t.Text, 1.5)
			})
			btn := ui.Box(c).Key("play").Size(28, 28).Radius(14).Margin(0, 0, 0, 14).Background(t.SurfaceSubtle).Center().
				Cursor(ui.CursorPointer).OnClick(func() {
				if !st.playing && st.t >= 1 {
					st.t = 0
				}
				st.playing = !st.playing
				st.last = c.Now()
			})
			btn.Children(func() {
				if st.playing {
					ui.Row(c).Gap(3).Children(func() {
						ui.Box(c).Size(3, 11).Radius(1).Background(t.Text)
						ui.Box(c).Size(3, 11).Radius(1).Background(t.Text)
					})
				} else {
					ui.Box(c).Size(12, 13).Margin(0, 0, 0, 2).Draw(func(p *ui.Painter, r ui.Rect) {
						var tri ui.Path
						tri.MoveTo(r.X+1, r.Y+1).LineTo(r.X+r.W, r.Y+r.H/2).LineTo(r.X+1, r.Y+r.H-1).Close()
						p.FillPath(&tri, t.Text)
					})
				}
			})
		})
	})
	return card
}

func init() {
	registerDemo(
		"race-bars",
		"A bar chart race: bars grow, the year counts up and rows swap places; play it or scrub the timeline.",
		func(c *ui.Context) { RaceBars(c, RaceBarsMarketCap()) },
	)
}
