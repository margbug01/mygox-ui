package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// WakeSliderOptions configures WakeSlider.
type WakeSliderOptions struct {
	Bars int     // how many bars make the track, 12 to 64; 36 by default
	Wake float32 // how high the bars rise behind a fast drag, 0 to 100; 60 by default, negative for none
	// Value is where it starts, 0 to 100; 50 by default.
	Value    float32
	OnChange func(value int)
}

const (
	wakeSliderW      = 324
	wakeSliderH      = 84
	wakeSliderTrackX = 16  // the first bar's left
	wakeSliderTrackW = 248 // from the first bar's left to the last one's right
	wakeSliderBarW   = 3
	wakeSliderBarH   = 20
)

type wakeSliderState struct {
	value  int
	lift   []float32 // each bar's rise, 0 to 1, and how fast it moves
	speed  []float32
	lastX  float32
	lastAt time.Time
	held   bool
}

// WakeSlider is a range slider with no handle, made of thin bars: press
// or drag anywhere along it to set the value. Drag fast and the bars you
// sweep past spring up in a wake behind your hand.
func WakeSlider(c *ui.Context, opts ...WakeSliderOptions) ui.Element {
	var o WakeSliderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	bars := 36
	if o.Bars != 0 {
		bars = int(clamp32(float32(o.Bars), 12, 64))
	}
	wake := float32(60)
	if o.Wake != 0 {
		wake = clamp32(o.Wake, 0, 100)
	}
	start := float32(50)
	if o.Value != 0 {
		start = clamp32(o.Value, 0, 100)
	}
	t := ThemeFrom(c)
	now := c.Now()

	track := ui.Box(c).Size(wakeSliderW, wakeSliderH).Cursor(ui.CursorResizeEW)
	st := ui.Local(track, "wake_slider", func() wakeSliderState { return wakeSliderState{value: int(start)} })
	if len(st.lift) != bars {
		st.lift, st.speed = make([]float32, bars), make([]float32, bars)
	}
	pitch := (wakeSliderTrackW - wakeSliderBarW) / float32(bars-1)
	barX := func(i int) float32 { return wakeSliderTrackX + float32(i)*pitch + wakeSliderBarW/2 }

	dt := float32(now.Sub(st.lastAt).Seconds())
	_, _, dragging := track.Dragged()
	held := dragging || track.Pressed()
	if x, _, ok := track.PointerPosition(); ok && held {
		u := clamp32((x-wakeSliderTrackX)/wakeSliderTrackW, 0, 1)
		if v := int(math.Round(float64(u * 100))); v != st.value {
			st.value = v
			if o.OnChange != nil {
				o.OnChange(v)
			}
		}
		// The bars swept since the last frame rise with the hand's speed.
		if st.held && dt > 0 {
			speed := abs(x-st.lastX) / dt
			kick := clamp32(speed/1400, 0, 1)
			lo, hi := min(st.lastX, x), max(st.lastX, x)
			for i := range st.lift {
				if bx := barX(i); bx >= lo && bx <= hi {
					st.speed[i] = max(st.speed[i], kick*22)
				}
			}
		}
		if dt > 0 || !st.held {
			st.lastX = x
		}
	}
	st.held = held
	if dt > 0 {
		st.lastAt = now
	}

	// Each bar's rise is a little spring pulled back to rest.
	moving := false
	step := min(dt, 1.0/30)
	for i := range st.lift {
		a := -260*st.lift[i] - 15*st.speed[i]
		st.speed[i] += a * step
		st.lift[i] += st.speed[i] * step
		if abs(st.lift[i]) > 0.002 || abs(st.speed[i]) > 0.01 {
			moving = true
		} else {
			st.lift[i], st.speed[i] = 0, 0
		}
	}
	if moving || held {
		c.AnimationFrame()
	}

	filled := int(math.Round(float64(st.value) / 100 * float64(bars)))
	rise := wake / 100 * 26
	lift := st.lift
	track.Draw(func(p *ui.Painter, r ui.Rect) {
		cy := r.Y + wakeSliderH/2
		for i := range bars {
			h := wakeSliderBarH + max(lift[i], -0.4)*rise
			col := t.TextSubtle.Alpha(0.5)
			if i < filled {
				col = t.Text
			}
			p.Fill(ui.Rect{X: r.X + barX(i) - wakeSliderBarW/2, Y: cy - h/2, W: wakeSliderBarW, H: h}, col, wakeSliderBarW/2)
		}
	})
	track.Children(func() {
		ui.Text(c, fmt.Sprint(st.value)).Absolute().Left(272).Top(wakeSliderH/2-12).Size(30, 24).TextAlign(ui.End).
			FontSize(16).TextColor(t.TextMuted)
	})
	return track
}

func abs(v float32) float32 { return float32(math.Abs(float64(v))) }

func init() {
	registerDemo(
		"wake-slider",
		"A thumbless range slider made of thin bars; drag fast and the bars you sweep past spring up in a wake behind your hand.",
		func(c *ui.Context) { WakeSlider(c) },
	)
}
