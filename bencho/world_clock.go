package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// WorldClockFormat is a 24 or 12 hour clock.
type WorldClockFormat string

const (
	WorldClock24h WorldClockFormat = "24h"
	WorldClock12h WorldClockFormat = "12h"
)

// WorldClockOptions configures WorldClock.
type WorldClockOptions struct {
	Clock  WorldClockFormat // 24h by default
	Corner float32          // 0 to 32px; 20 by default, negative for square
	// OnHome is called with the city tapped to make it home.
	OnHome func(city string)
}

const (
	worldClockSize  = 300
	worldClockRow   = 58.5
	worldClockRuler = 4 * worldClockRow // where the ruler starts
	worldClockHour  = 29.3              // ruler px an hour
	worldClockStart = 18*60 + 40        // home's time when it opens, in minutes
)

// The card is always the dark fill bencho.dev shows it in.
var (
	worldClockBg    = ui.Hex("#0d0f10")
	worldClockHome  = ui.Hex("#1b1d1e")
	worldClockInk   = ui.Hex("#f5f7f8")
	worldClockMuted = ui.Hex("#86888a")
)

type worldClockCity struct {
	name string
	utc  int // hours from UTC
}

var worldClockCities = []worldClockCity{{"Tbilisi", 4}, {"San Francisco", -7}, {"New York", -4}, {"London", 0}}

type worldClockState struct {
	home  int     // the city on top
	shift float32 // minutes the ruler has been scrubbed by
}

func worldClockTime(min int, f WorldClockFormat) string {
	min = ((min % 1440) + 1440) % 1440
	h, m := min/60, min%60
	if f == WorldClock12h {
		half := "AM"
		if h >= 12 {
			half = "PM"
		}
		return fmt.Sprintf("%d:%02d %s", (h+11)%12+1, m, half)
	}
	return fmt.Sprintf("%02d:%02d", h, m)
}

// WorldClock shows home's time alongside three cities. Dragging the ruler
// at the bottom shifts them all together; tapping a city makes it home.
func WorldClock(c *ui.Context, opts ...WorldClockOptions) ui.Element {
	var o WorldClockOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	corner := float32(20)
	if o.Corner != 0 {
		corner = clamp32(o.Corner, 0, 32)
	}

	card := ui.Box(c).Size(worldClockSize, worldClockSize).Radius(corner).Background(worldClockBg).Clip()
	st := ui.Local(card, "world_clock", func() worldClockState { return worldClockState{} })
	homeUTC := worldClockCities[st.home].utc
	// Home's time, in minutes, snapped to five as the ruler turns.
	now := worldClockStart + int(math.Round(float64(st.shift)/5))*5

	card.Children(func() {
		order := []int{st.home}
		for i := range worldClockCities {
			if i != st.home {
				order = append(order, i)
			}
		}
		for row, i := range order {
			city := worldClockCities[i]
			r := ui.Row(c).Key(city.name).Absolute().Left(0).Top(float32(row)*worldClockRow).Size(worldClockSize, worldClockRow).
				AlignItems(ui.Center).Padding(0, 16.5, 0, 16.5)
			if row == 0 {
				r.Background(worldClockHome)
			} else {
				r.Cursor(ui.CursorPointer).OnClick(func() {
					st.home = i
					if o.OnHome != nil {
						o.OnHome(city.name)
					}
				})
				r.DrawOver(func(p *ui.Painter, b ui.Rect) {
					p.Fill(ui.Rect{X: b.X, Y: b.Y + b.H - 1, W: b.W, H: 1}, worldClockHome, 0)
				})
			}
			diff := city.utc - homeUTC
			r.Children(func() {
				ui.Text(c, city.name).FontSize(20.5).TextColor(worldClockInk)
				if row > 0 {
					sign := "+"
					if diff < 0 {
						sign = "−"
					}
					ui.Text(c, fmt.Sprintf("%s%dh", sign, max(diff, -diff))).FontSize(11.5).TextColor(worldClockMuted).Margin(3, 0, 0, 9)
				}
				ui.Box(c).Grow(1)
				ui.Text(c, worldClockTime(now+diff*60, o.Clock)).FontSize(20.5).TextColor(worldClockInk)
			})
		}

		ruler := ui.Box(c).Key("ruler").Absolute().Left(0).Top(worldClockRuler).Size(worldClockSize, worldClockSize-worldClockRuler).
			Cursor(ui.CursorResizeEW)
		if dx, _, dragging := ruler.Dragged(); dragging {
			st.shift -= dx / worldClockHour * 60
		}
		at := ruler.Spring("at", float32(now), ui.SpringOptions{Duration: 200 * time.Millisecond})
		ruler.Draw(func(p *ui.Painter, r ui.Rect) {
			mid := r.X + worldClockSize/2
			fade := func(x float32) float32 {
				return clamp32(1-float32(math.Abs(float64(x-mid)))/(worldClockSize/2+10), 0, 1)
			}
			first := int(math.Floor(float64(at)/15)) - 24
			for q := first; q <= first+48; q++ {
				x := mid + (float32(q*15)-at)/60*worldClockHour
				if x < r.X-4 || x > r.X+r.W+4 {
					continue
				}
				f := fade(x)
				// The ruler bends away a little at its ends, like a drum.
				lift := 0.0002 * (x - mid) * (x - mid)
				h, top := float32(8.5), r.Y+41-lift
				if q%4 == 0 {
					h, top = 18.5, r.Y+34-lift
					hour := ((q/4)%24 + 24) % 24
					label := fmt.Sprint(hour)
					if o.Clock == WorldClock12h {
						label = fmt.Sprint((hour+11)%12 + 1)
					}
					near := clamp32(1-float32(math.Abs(float64(x-mid)))/worldClockHour, 0, 1)
					span := ui.Span{Text: label, Size: 12, Weight: 400, Color: worldClockInk.Alpha(0.15 + 0.5*f*f + 0.2*near)}
					w, _ := p.MeasureText(0, span)
					p.RichText(x-w/2, r.Y+16-lift, 0, span)
				}
				p.Fill(ui.Rect{X: x - 1, Y: top, W: 2, H: h}, worldClockInk.Alpha(0.1+0.42*f*f), 1)
			}
			p.Fill(ui.Rect{X: mid - 1, Y: r.Y + 32, W: 2, H: 23.5}, worldClockInk, 1)
		})
	})
	return card
}

func init() {
	registerDemo(
		"world-clock",
		"A world clock with a time scrubber: your time alongside three cities, dragging the ruler shifts them all together, and tapping a city makes it home.",
		func(c *ui.Context) { WorldClock(c) },
	)
}
