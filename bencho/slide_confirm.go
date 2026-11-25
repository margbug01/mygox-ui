package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// SlideConfirmOptions configures SlideConfirm.
type SlideConfirmOptions struct {
	Width     float32 // the track's length, 220 to 380; 280 by default
	Corner    float32 // 0 to 28; 28 by default
	Speed     float32 // how fast it takes over once let go, 0 to 100; 50 by default
	Label     string  // "Slide to confirm" by default
	Done      string  // "Confirmed" by default
	OnConfirm func()
}

type slideConfirmPhase int

const (
	slideConfirmIdle      slideConfirmPhase = iota // at rest, or springing back
	slideConfirmHeld                               // following the pointer
	slideConfirmFinishing                          // let go past the mark: it slides on by itself
	slideConfirmDone                               // at the end, confirmed
)

type slideConfirmState struct {
	phase  slideConfirmPhase
	x, vel float32 // the handle's offset along the track and its speed
	from   float32 // where the take-over started
	at     time.Time
	last   time.Time
}

const (
	slideConfirmH     = 56
	slideConfirmInset = 4
	slideConfirmKnob  = 48
	slideConfirmMark  = 0.6 // share of the travel past which a release confirms
)

// SlideConfirm is a slide to confirm button that asks for a deliberate
// gesture: push the handle across the track, and once it passes the mark
// it takes over and finishes the slide by itself; let go short of it and
// it springs back.
func SlideConfirm(c *ui.Context, opts ...SlideConfirmOptions) ui.Element {
	var o SlideConfirmOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Width == 0 {
		o.Width = 280
	}
	if o.Corner == 0 {
		o.Corner = 28
	}
	if o.Speed == 0 {
		o.Speed = 50
	}
	if o.Label == "" {
		o.Label = "Slide to confirm"
	}
	if o.Done == "" {
		o.Done = "Confirmed"
	}
	w := clamp32(o.Width, 220, 380)
	corner := clamp32(o.Corner, 0, 28)
	travel := w - slideConfirmKnob - 2*slideConfirmInset
	t := ThemeFrom(c)
	now := c.Now()

	track := ui.Box(c).Size(w, slideConfirmH)
	st := ui.Local(track, "slide_confirm", func() slideConfirmState { return slideConfirmState{last: now} })
	dt := clamp32(float32(now.Sub(st.last).Seconds()), 0, 0.05)
	st.last = now

	track.Children(func() {
		knob := ui.Box(c).Absolute().Top(slideConfirmInset).Size(slideConfirmKnob, slideConfirmKnob)
		dx, _, dragging := knob.Dragged()
		switch {
		case dragging && (st.phase == slideConfirmIdle || st.phase == slideConfirmHeld):
			st.phase = slideConfirmHeld
			st.x, st.vel = clamp32(st.x+dx, 0, travel), 0
		case st.phase == slideConfirmHeld:
			if st.x >= travel*slideConfirmMark {
				st.phase, st.from, st.at = slideConfirmFinishing, st.x, now
			} else {
				st.phase = slideConfirmIdle
			}
		}
		switch st.phase {
		case slideConfirmIdle:
			// Back to the start on a spring.
			for range 4 {
				h := dt / 4
				st.vel += (-st.x*420 - st.vel*26) * h
				st.x += st.vel * h
			}
			if st.x > 0.2 || math.Abs(float64(st.vel)) > 1 {
				c.AnimationFrame()
			} else {
				st.x, st.vel = 0, 0
			}
		case slideConfirmFinishing:
			dur := lerp32(0.9, 0.14, clamp32(o.Speed, 0, 100)/100)
			k := clamp32(float32(now.Sub(st.at).Seconds())/dur, 0, 1)
			e := 1 - (1-k)*(1-k)*(1-k)
			st.x = st.from + (travel-st.from)*e
			if k >= 1 {
				st.phase, st.at = slideConfirmDone, now
				if o.OnConfirm != nil {
					o.OnConfirm()
				}
			}
			c.AnimationFrame()
		case slideConfirmDone:
			st.x = travel
			if now.Sub(st.at) > 2200*time.Millisecond {
				st.phase = slideConfirmIdle
			}
			c.AnimationFrame()
		}

		done := track.Spring("done", b2f(st.phase == slideConfirmDone), ui.SpringOptions{Duration: 360 * time.Millisecond, Bounce: 0.25})
		press := knob.Spring("press", b2f(knob.Pressed() || st.phase == slideConfirmHeld), ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.3})
		hover := knob.Spring("hover", b2f(knob.Hovered()), ui.SpringOptions{Duration: 180 * time.Millisecond})
		progress := st.x / travel
		knobR := max(corner-slideConfirmInset, 0)

		track.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, t.Surface, corner)
			label := func(s string, col ui.Color, a, shift float32) {
				if a <= 0.01 {
					return
				}
				span := ui.Span{Text: s, Size: 14, Weight: 400, Color: col.Alpha(a)}
				tw, th := p.MeasureText(0, span)
				p.RichText(r.X+(r.W-tw)/2+shift, r.Y+(r.H-th)/2, tw+2, span)
			}
			label(o.Label, t.TextSubtle, (1-progress*1.7)*(1-done), progress*14)
			label(o.Done, t.Text, done, -(slideConfirmKnob/2)*done)
		})
		knob.Left(slideConfirmInset+st.x).Cursor(ui.CursorGrab).
			Scale(1+0.04*hover-0.07*press, 1+0.04*hover-0.07*press)
		knob.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, t.Accent, knobR)
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			if a := 1 - done; a > 0.01 {
				s := 1 - done*0.4
				p.StrokePolyline([][2]float32{{cx - 7.5*s, cy}, {cx + 7.5*s, cy}}, 1.7, t.AccentText.Alpha(a))
				p.StrokePolyline([][2]float32{{cx + 0.2*s, cy - 7.3*s}, {cx + 7.5*s, cy}, {cx + 0.2*s, cy + 7.3*s}}, 1.7, t.AccentText.Alpha(a))
			}
			if done > 0.01 {
				s := 0.6 + 0.4*done
				p.StrokePolyline([][2]float32{{cx - 7*s, cy + 0.5*s}, {cx - 2.2*s, cy + 5*s}, {cx + 7.5*s, cy - 5*s}}, 1.9, t.AccentText.Alpha(clamp32(done, 0, 1)))
			}
		})
	})
	return track
}

func init() {
	registerDemo(
		"slide-confirm",
		"A slide to confirm button: push the handle across the track, and once it passes the mark it takes over and finishes the slide by itself.",
		func(c *ui.Context) { SlideConfirm(c) },
	)
}
