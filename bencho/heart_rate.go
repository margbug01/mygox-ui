package bencho

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
)

// HeartRateOptions configures HeartRate.
type HeartRateOptions struct {
	Corner float32 // 0 to 50px; 20 by default
}

// heartRateHours is each hour's lowest and highest beats per minute, from
// midnight; the last is now.
var heartRateHours = [24][2]float32{
	{57.6, 65.4}, {54.9, 62.3}, {52.5, 59.2}, {53.7, 57.6}, {54.1, 60.7}, {57.2, 63.8},
	{55.6, 69.3}, {62, 75.6}, {65.8, 81.4}, {63.8, 77.9}, {61.2, 71.7}, {65.8, 82.9},
	{69.7, 89.2}, {77.5, 95.8}, {72.4, 87.2}, {65.4, 79}, {59.4, 73.1}, {64.6, 81.4},
	{70.1, 85.7}, {66.2, 77.5}, {60.4, 71.7}, {63.4, 73.1}, {65.8, 75.6}, {67.8, 75.6},
}

const (
	heartRateW     = 256
	heartRateH     = 252
	heartRatePeak  = 97 // at 1PM
	heartRateMin   = 51 // at 3AM
	heartRateLeft  = 25 // midnight's column
	heartRatePitch = 9.04
	heartRateTop   = 129
	heartRateFloor = 210.5
)

type heartRateState struct {
	beat time.Time // when the heart last beat
}

// heartRateHeart is Lucide's heart.
var heartRateHeart = sync.OnceValue(func() *ui.SVG {
	return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 9.5a5.5 5.5 0 0 1 9.591-3.676.56.56 0 0 0 .818 0A5.49 5.49 0 0 1 22 9.5c0 2.29-1.5 4-3 5.5l-5.492 5.313a2 2 0 0 1-3 .019L5 15c-1.5-1.5-3-3.2-3-5.5"/></svg>`))
})

func heartRateY(bpm float32) float32 { return 264.9 - 1.2826*bpm }

// HeartRate is a dark heart rate widget over a chart of the whole day, one
// bar of lowest to highest per hour. Scrubbing through the hours shows the
// rate for that moment, and the heart beats at it.
func HeartRate(c *ui.Context, opts ...HeartRateOptions) ui.Element {
	var o HeartRateOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	corner := float32(20)
	if o.Corner != 0 {
		corner = clamp32(o.Corner, 0, 50)
	}
	if o.Corner < 0 {
		corner = 0
	}
	t := ThemeFrom(c)
	fill := ui.Hex("#0d0f10")
	white, grey, candle, grid := ui.Hex("#ffffff"), ui.Hex("#8f9092"), ui.Hex("#5a5c5e"), ui.Hex("#1c1e1f")
	card := ui.Box(c).Size(heartRateW, heartRateH).Radius(corner).Background(fill).Clip()
	if t.Dark {
		card.Border(1, t.Border)
	}
	st := ui.Local(card, "heart_rate", func() heartRateState { return heartRateState{} })

	// The hour under the pointer, or now.
	hour := len(heartRateHours) - 1
	if card.Hovered() {
		if x, y, ok := card.PointerPosition(); ok && y > heartRateTop-20 && y < heartRateFloor+30 {
			hour = int(clamp32((x-heartRateLeft)/heartRatePitch+0.5, 0, 23))
		}
	}
	h := heartRateHours[hour]
	bpm := float32(math.Round(float64((h[0] + h[1]) / 2)))

	// The heart beats once every 60/bpm seconds: a quick swell and settle.
	now := c.Now()
	period := time.Duration(float64(time.Minute) / float64(bpm))
	if st.beat.IsZero() || now.Sub(st.beat) > 2*period {
		st.beat = now
	}
	for now.Sub(st.beat) >= period {
		st.beat = st.beat.Add(period)
	}
	since := float32(now.Sub(st.beat).Seconds())
	pulse := float32(0)
	if since < 0.3 {
		pulse = float32(math.Sin(float64(since/0.3) * math.Pi))
	}
	c.AnimationFrame()

	card.Draw(func(p *ui.Painter, r ui.Rect) {
		for i := range heartRateHours {
			x := r.X + heartRateLeft + float32(i)*heartRatePitch
			p.Fill(ui.Rect{X: x - 0.5, Y: r.Y + heartRateTop, W: 1, H: heartRateFloor - heartRateTop}, grid, 0)
		}
		for i, hr := range heartRateHours {
			x := r.X + heartRateLeft + float32(i)*heartRatePitch
			top, bottom := r.Y+heartRateY(hr[1]), r.Y+heartRateY(hr[0])
			col := candle
			if i == hour {
				col = white
			}
			p.Fill(ui.Rect{X: x - 2, Y: top, W: 4, H: bottom - top}, col, 2)
		}
		px, py := r.X+heartRateLeft+13*heartRatePitch, r.Y+heartRateY(heartRatePeak)
		p.Fill(ui.Rect{X: px - 3, Y: py - 3, W: 6, H: 6}, white, 3)
		mx, my := r.X+heartRateLeft+3*heartRatePitch, r.Y+heartRateY(heartRateMin)
		p.Fill(ui.Rect{X: mx - 3.5, Y: my - 3.5, W: 7, H: 7}, fill, 3.5)
		p.Stroke(ui.Rect{X: mx - 3, Y: my - 3, W: 6, H: 6}, white, 3, 1.2)
	})

	card.Children(func() {
		ui.Row(c).Absolute().Left(22).Top(13).Height(50).AlignItems(ui.End).Children(func() {
			ui.Text(c, fmt.Sprint(bpm)).FontSize(41).TextColor(white).Height(50)
			ui.Text(c, "BPM").FontSize(16).TextColor(grey).Margin(0, 0, 6, 5)
		})
		s := 1 + 0.16*pulse
		ui.Icon(c, heartRateHeart()).Absolute().Left(214).Top(21).Size(23, 23).TextColor(lerpColor(grey, white, pulse)).Scale(s, s)

		for i, row := range [2]struct {
			label string
			value int
			solid bool
		}{{"PEAK:", heartRatePeak, true}, {"MIN:", heartRateMin, false}} {
			ui.Row(c).Key(row.label).Absolute().Left(20).Top(74 + float32(i)*27.5).Height(18).AlignItems(ui.Center).Children(func() {
				dot := ui.Box(c).Size(6, 6).Radius(3)
				if row.solid {
					dot.Background(white)
				} else {
					dot.Border(1.2, white)
				}
				ui.Text(c, row.label).FontSize(12).LetterSpacing(0.8).TextColor(grey).Margin(0, 0, 0, 9)
				ui.Text(c, fmt.Sprintf("%d BPM", row.value)).FontSize(12).LetterSpacing(0.6).TextColor(white).Margin(0, 0, 0, 8)
			})
		}
		for i, label := range [4]string{"12AM", "6AM", "12PM", "6PM"} {
			x := 34 + float32(i)*53.5
			ui.Text(c, label).Key(label).FontSize(10.5).TextColor(grey).Absolute().Left(x - 20).Top(218).Width(40).TextAlign(ui.Center)
		}
	})
	return card
}

func init() {
	registerDemo(
		"heart-rate",
		"A heart rate widget over a chart of the day; scrub the hours and the heart beats at that moment's rate.",
		func(c *ui.Context) { HeartRate(c) },
	)
}
