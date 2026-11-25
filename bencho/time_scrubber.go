package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// TimeScrubberOptions configures TimeScrubber.
type TimeScrubberOptions struct {
	Step     int     // what it settles on, in minutes: 5, 15 or 30; 15 by default
	Momentum float32 // how far a fling carries, 0 to 100; 50 by default
	H24      bool    // a 24 hour clock
	Corner   float32 // 0 to 32; 20 by default
	Start    int     // minutes after midnight; 9:30 by default
	OnChange func(minutes int)
}

type timeScrubberState struct {
	value    float64 // minutes after midnight, unsettled
	vel      float64 // minutes per second
	last     time.Time
	wasDrag  bool
	reported int
	dragging bool    // the ruler was being dragged last frame
	pending  float32 // how far it was dragged then, in px
}

const (
	timeScrubberW     = 300
	timeScrubberPitch = 10.0 // px per five minutes
)

// timeScrubberLabel formats minutes after midnight on a 12 or 24 hour clock.
func timeScrubberLabel(m int, h24 bool) (clock, half string) {
	m = ((m % 1440) + 1440) % 1440
	h, mm := m/60, m%60
	if h24 {
		return fmt.Sprintf("%02d:%02d", h, mm), ""
	}
	half = "AM"
	if h >= 12 {
		half = "PM"
	}
	if h = h % 12; h == 0 {
		h = 12
	}
	return fmt.Sprintf("%d:%02d", h, mm), half
}

// TimeScrubber is a time picker on a ruler under a red needle. Dragging the
// ruler sets the time; a fling coasts on and settles on the nearest step.
func TimeScrubber(c *ui.Context, opts ...TimeScrubberOptions) ui.Element {
	var o TimeScrubberOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Step == 0 {
		o.Step = 15
	}
	if o.Momentum == 0 {
		o.Momentum = 50
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	if o.Start == 0 {
		o.Start = 9*60 + 30
	}
	t := ThemeFrom(c)
	card := ui.Column(c).Width(timeScrubberW).Height(141).Padding(18, 0, 0, 0).Radius(clamp32(o.Corner, 0, 32)).
		Background(t.Surface).AlignItems(ui.Center).Clip()
	st := ui.Local(card, "time_scrubber", func() timeScrubberState {
		return timeScrubberState{value: float64(o.Start), last: c.Now(), reported: o.Start}
	})
	now := c.Now()
	dt := math.Min(now.Sub(st.last).Seconds(), 0.05)
	st.last = now
	step := float64(o.Step)
	glide := 1.2 + float64(clamp32(o.Momentum, 0, 100))/100*4 // how slowly a fling dies away

	// The ruler reports its drag as it is built, below; it lands here next frame.
	if dx, dragging := st.pending, st.dragging; dragging {
		st.pending = 0
		d := -float64(dx) / timeScrubberPitch * 5
		st.value += d
		if dt > 0 {
			st.vel = st.vel*0.5 + d/dt*0.5
		}
		st.wasDrag = true
	} else {
		if st.wasDrag {
			st.wasDrag = false
		}
		if math.Abs(st.vel) > 8 {
			// Coasting.
			st.value += st.vel * dt
			st.vel *= math.Exp(-dt * 6 / glide)
		} else {
			// Settling on the nearest step.
			st.vel = 0
			target := math.Round(st.value/step) * step
			st.value += (target - st.value) * math.Min(1, dt*12)
			if math.Abs(target-st.value) < 0.02 {
				st.value = target
			}
		}
	}
	if st.dragging || st.wasDrag || st.vel != 0 || st.value != math.Round(st.value/step)*step {
		c.AnimationFrame()
	}
	shown := int(math.Round(st.value))
	if settled := int(math.Round(st.value/step) * step); !st.wasDrag && st.vel == 0 && settled != st.reported {
		st.reported = settled
		if o.OnChange != nil {
			o.OnChange(((settled % 1440) + 1440) % 1440)
		}
	}
	clock, half := timeScrubberLabel(shown, o.H24)
	needle := ui.Hex("#ef4444")

	card.Children(func() {
		ui.Row(c).AlignItems(ui.End).Gap(6).Children(func() {
			ui.Text(c, clock).FontSize(31).FontWeight(500).TextColor(t.Text)
			if half != "" {
				ui.Text(c, half).FontSize(15).FontWeight(500).TextColor(t.TextSubtle).Margin(0, 0, 6, 0)
			}
		})
		ruler := ui.Box(c).Key("ruler").Width(timeScrubberW).Height(70).Cursor(ui.CursorGrab)
		dx, _, dragging := ruler.Dragged()
		st.pending += dx
		st.dragging = dragging
		ruler.Margin(8, 0, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
			cx := r.X + r.W/2
			base := r.Y + 26 // the ticks hang either side of this line
			first := int(math.Floor(st.value/5)) - 16
			for i := first; i <= first+33; i++ {
				m := i * 5
				x := cx + float32((float64(m)-st.value)/5*timeScrubberPitch)
				fade := clamp32(1-float32(math.Abs(float64(x-cx)))/(r.W/2-10), 0, 1)
				if fade <= 0 {
					continue
				}
				h, ink := float32(11), t.TextMuted
				switch {
				case ((m%60)+60)%60 == 0:
					h, ink = 22, t.Text
				case ((m%15)+15)%15 == 0:
					h, ink = 15, t.Text
				}
				p.Fill(ui.Rect{X: x - 1, Y: base - 11, W: 2, H: h}, ink.Alpha(fade), 1)
				if ((m%60)+60)%60 == 0 {
					label := fmt.Sprintf("%02d:00", ((m/60)%24+24)%24)
					if !o.H24 {
						hh := ((m/60)%24 + 24) % 24
						suffix := " AM"
						if hh >= 12 {
							suffix = " PM"
						}
						if hh%12 == 0 {
							hh = 12
						} else {
							hh %= 12
						}
						label = fmt.Sprint(hh) + suffix
					}
					span := ui.Span{Text: label, Size: 12, Color: t.TextSubtle.Alpha(fade)}
					w, _ := p.MeasureText(0, span)
					p.RichText(x-w/2, base+20, w+2, span)
				}
			}
			p.Fill(ui.Rect{X: cx - 1.5, Y: base - 18, W: 3, H: 36}, needle, 1.5)
		})
	})
	return card
}

func init() {
	registerDemo(
		"time-scrubber",
		"A time picker on a ruler under a red needle: drag to set the time, or fling it and it coasts to the nearest step.",
		func(c *ui.Context) { TimeScrubber(c) },
	)
}
