package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// StepPlayerOptions configures StepPlayer.
type StepPlayerOptions struct {
	Wake   float32        // how far back the ripple runs when a step is done, 0 to 100; 100 by default
	Pace   time.Duration  // how long a step takes to fill, 1s to 6s; 3s by default
	Steps  int            // 3 to 8; 5 by default
	OnStep func(step int) // called with the step it moves to, from 0
}

const (
	stepPlayerH      = 36
	stepPlayerBar    = 34 // the current step's bar
	stepPlayerDot    = 6.5
	stepPlayerGap    = 10
	stepPlayerPad    = 15
	stepPlayerRipple = 320 * time.Millisecond
)

// stepPlayerPlay and stepPlayerPause are solid, as bencho.dev's are.
var (
	stepPlayerPlay  = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M7 4.5v15a1 1 0 0 0 1.5.86l12.5-7.5a1 1 0 0 0 0-1.72L8.5 3.64A1 1 0 0 0 7 4.5Z"/></svg>`))
	stepPlayerPause = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><rect x="5" y="4" width="5" height="16" rx="1.5"/><rect x="14" y="4" width="5" height="16" rx="1.5"/></svg>`))
)

type stepPlayerState struct {
	step     int
	elapsed  time.Duration // into the current step
	playing  bool
	last     time.Time
	movedAt  time.Time // when a step was last done, for the ripple
	finished bool
}

// StepPlayer is a step indicator that plays through a sequence: press play
// and the current step fills as a bar, then eases on to the next dot, a
// ripple running back through the steps done. Press a dot to jump to it.
func StepPlayer(c *ui.Context, opts ...StepPlayerOptions) ui.Element {
	var o StepPlayerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Wake == 0 {
		o.Wake = 100
	}
	if o.Pace == 0 {
		o.Pace = 3 * time.Second
	}
	if o.Steps == 0 {
		o.Steps = 5
	}
	o.Wake = clamp32(o.Wake, 0, 100)
	o.Pace = min(max(o.Pace, time.Second), 6*time.Second)
	o.Steps = min(max(o.Steps, 3), 8)
	t := ThemeFrom(c)
	now := c.Now()

	pillW := 2*stepPlayerPad + stepPlayerBar + float32(o.Steps-1)*(stepPlayerDot+stepPlayerGap)
	root := ui.Row(c).Height(stepPlayerH).Gap(8)
	st := ui.Local(root, "step_player", func() stepPlayerState { return stepPlayerState{} })
	moveTo := func(step int) {
		st.step, st.elapsed, st.movedAt = step, 0, c.Now()
		if o.OnStep != nil {
			o.OnStep(step)
		}
	}
	if st.playing {
		if !st.last.IsZero() {
			st.elapsed += min(now.Sub(st.last), 100*time.Millisecond)
		}
		if st.elapsed >= o.Pace {
			if st.step == o.Steps-1 {
				st.elapsed, st.playing, st.finished = o.Pace, false, true
			} else {
				moveTo(st.step + 1)
			}
		}
		c.AnimationFrame()
	}
	st.last = now
	fill := float32(st.elapsed) / float32(o.Pace)
	since := now.Sub(st.movedAt)
	if since < stepPlayerRipple+time.Duration(o.Steps)*80*time.Millisecond {
		c.AnimationFrame()
	}

	root.Children(func() {
		ui.Box(c).Size(pillW, stepPlayerH).Radius(stepPlayerH / 2).Background(t.Surface).Children(func() {
			x := float32(stepPlayerPad)
			reach := o.Wake / 100 * float32(o.Steps)
			for i := range o.Steps {
				item := ui.Box(c).Key(i).Absolute().Top(0).Height(stepPlayerH).Cursor(ui.CursorPointer).
					OnClick(func() { st.finished = false; moveTo(i) })
				w := item.Spring("w", b2f(i == st.step)*(stepPlayerBar-stepPlayerDot)+stepPlayerDot, ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.2})
				item.Left(x - stepPlayerGap/2).Width(w + stepPlayerGap)
				x += w + stepPlayerGap
				// The ripple runs back from the step just done, fading as it goes.
				bump := float32(0)
				if back := float32(st.step - 1 - i); i < st.step && back < reach {
					tt := since - time.Duration(back)*80*time.Millisecond
					if tt > 0 && tt < stepPlayerRipple {
						bump = float32(math.Sin(math.Pi*float64(tt)/float64(stepPlayerRipple))) * (1 - back/max(reach, 1))
					}
				}
				done, current := i < st.step, i == st.step
				f := fill
				item.Draw(func(p *ui.Painter, r ui.Rect) {
					h := stepPlayerDot * (1 + 0.45*bump)
					body := ui.Rect{X: r.X + stepPlayerGap/2, Y: r.Y + r.H/2 - h/2, W: r.W - stepPlayerGap, H: h}
					if body.W < h {
						body.X -= (h - body.W) / 2
						body.W = h
					}
					switch {
					case done:
						p.Fill(body, t.Text, h/2)
					case current:
						p.Fill(body, t.BorderStrong, h/2)
						if w := body.W * f; w > 0.5 {
							p.Clip(body, h/2, func() { p.Fill(ui.Rect{X: body.X, Y: body.Y, W: w, H: h}, t.Text, 0) })
						}
					default:
						p.Fill(body, t.BorderStrong, h/2)
					}
				})
			}
		})
		btn := ui.Box(c).Size(stepPlayerH, stepPlayerH).Radius(stepPlayerH / 2).Background(t.Surface).Center().Cursor(ui.CursorPointer).
			OnClick(func() {
				if st.finished {
					st.finished = false
					moveTo(0)
				}
				st.playing = !st.playing
			})
		press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
		btn.Scale(1-0.08*press, 1-0.08*press)
		icon, col := stepPlayerPlay, t.Text
		if st.playing {
			icon = stepPlayerPause
		}
		if btn.Pressed() {
			col = t.TextSubtle
		}
		btn.Children(func() {
			ui.Icon(c, icon).Size(14, 14).TextColor(col)
		})
	})
	return root
}

func init() {
	registerDemo(
		"step-player",
		"A step progress indicator that plays through a sequence: press play and each step fills as a bar, then eases on to the next dot.",
		func(c *ui.Context) { StepPlayer(c) },
	)
}
