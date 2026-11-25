package coss

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// ProgressOptions are a Progress bar's optional motion and look.
type ProgressOptions struct {
	// Bounce springs the bar to each new value, running past it by as much,
	// 0 to 100, and up to 8% past it, before it settles; 0 moves it at once.
	Bounce float32
	// Sparks throws sparks off the bar's head as it moves forward.
	Sparks bool
	// Style draws the bar in another design; nil is coss's look.
	Style *ProgressStyle
}

// ProgressStyle is the look of a Progress bar: a track Height tall, filled
// with a gradient from From to To.
type ProgressStyle struct {
	Height   float32
	Track    ui.Color
	From, To ui.Color
}

type progressState struct {
	last    float32   // the value, as last built
	sparkAt time.Time // when it last moved forward
}

const progressSparkLife = 550 * time.Millisecond

// Progress renders a coss-styled linear progress bar where value is in [0, 100].
func Progress(c *ui.Context, value float32, opts ...ProgressOptions) ui.Element {
	t := theme.Of(c)
	if value < 0 {
		value = 0
	}
	if value > 100 {
		value = 100
	}
	var o ProgressOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Bounce == 0 && !o.Sparks && o.Style == nil {
		return ui.Box(c).
			WidthPercent(100).
			Height(6).
			Radius(999).
			Background(t.Input).
			Clip().
			Children(func() {
				ui.Box(c).
					HeightPercent(100).
					WidthPercent(value).
					Background(t.Primary)
			})
	}
	s := o.Style
	if s == nil {
		s = &ProgressStyle{Height: 6, Track: t.Input, From: t.Primary, To: t.Primary}
	}
	bar := ui.Box(c).WidthPercent(100).Height(s.Height)
	now := c.Now()
	st := ui.Local(bar, "progress", func() progressState { return progressState{last: value} })
	if value > st.last {
		st.sparkAt = now
	}
	st.last = value
	frac := value / 100
	if o.Bounce > 0 {
		frac = bar.Spring("bar", frac, ui.SpringOptions{Duration: 650 * time.Millisecond, Bounce: float64(min(max(o.Bounce, 0), 100)) / 100 * 0.7})
	}
	sparkAge := now.Sub(st.sparkAt)
	sparks := o.Sparks && !st.sparkAt.IsZero() && sparkAge < progressSparkLife
	if sparks {
		c.AnimationFrame()
	}
	bar.Draw(func(p *ui.Painter, track ui.Rect) {
		rad := track.H / 2
		p.Fill(track, s.Track, rad)
		if w := track.W * min(max(frac, 0), 1.08); w > 0.5 {
			fill := track
			fill.W = max(w, track.H)
			p.FillGradient(fill, ui.LinearGradient{Angle: 90, Stops: []ui.GradientStop{{Color: s.From, At: 0}, {Color: s.To, At: 1}}}, rad)
		}
		if !sparks {
			return
		}
		k := float32(sparkAge.Seconds() / progressSparkLife.Seconds())
		hx, hy := track.X+track.W*min(max(frac, 0), 1), track.Y+rad
		for j := range 8 {
			a := float64(j)/8*2*math.Pi + 0.4
			d := 4 + 22*k
			x := hx + float32(math.Cos(a))*d*1.3
			y := hy + float32(math.Sin(a))*d*0.8
			sz := 3.2 * (1 - k)
			p.Fill(ui.Rect{X: x - sz/2, Y: y - sz/2, W: sz, H: sz}, lerpColor(s.From, s.To, float32(j%3)/2).Alpha(1-k), sz/2)
		}
	})
	return bar
}
