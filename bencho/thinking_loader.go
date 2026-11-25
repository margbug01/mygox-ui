package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// ThinkingLoaderOptions configures ThinkingLoader.
type ThinkingLoaderOptions struct {
	Step    time.Duration // how long each step shows, 1s to 4s; 2s by default
	NoShine bool          // no light sweeping across the label while it works
	Ask     string        // the question in the bubble above
	Steps   []string      // what it does, in order; the last is the finished label
}

type thinkingLoaderState struct {
	start time.Time
	open  bool
}

const thinkingLoaderH = 34

// ThinkingLoader is the small status pill an AI model shows while it works,
// under the question it was asked. Rays turn while the label rolls through
// each step, a light sweeping across it; at the end a check lands. Pressing
// the pill opens the steps.
func ThinkingLoader(c *ui.Context, opts ...ThinkingLoaderOptions) ui.Element {
	var o ThinkingLoaderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Step == 0 {
		o.Step = 2 * time.Second
	}
	o.Ask = orDefault(o.Ask, "Best CRM for a small team?")
	if len(o.Steps) == 0 {
		o.Steps = []string{"Understanding the question", "Reading 12 sources", "Comparing 5 CRMs", "Found 3 good fits"}
	}
	t := ThemeFrom(c)
	col := ui.Column(c).Width(340).MinHeight(232).Gap(16)
	st := ui.Local(col, "thinking_loader", func() thinkingLoaderState {
		// Open on the second step, as bencho.dev's still does.
		return thinkingLoaderState{start: c.Now().Add(-o.Step - o.Step/4)}
	})
	now := c.Now()
	working := len(o.Steps) - 1
	since := now.Sub(st.start)
	idx := min(int(since/o.Step), working)
	done := idx >= working
	if done && since > time.Duration(working)*o.Step+4*time.Second && !st.open {
		st.start = now // go round again
		since, idx, done = 0, 0, false
	}
	c.AnimationFrame()
	// The roll from the last label to this one.
	roll := clamp32(float32(since-time.Duration(idx)*o.Step)/float32(380*time.Millisecond), 0, 1)
	roll = 1 - (1-roll)*(1-roll)*(1-roll)
	k := clamp32(col.Spring("open", b2f(st.open), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.12}), 0, 1.05)
	phase := phaseSeconds(now)

	col.Children(func() {
		ui.Row(c).Key("ask").Children(func() {
			ui.Box(c).Grow(1)
			ui.Box(c).Height(34).Padding(0, 14).Radius(17).Background(lerpColor(t.SurfaceSubtle, t.SurfaceHover, 0.55)).Center().Children(func() {
				ui.Text(c, o.Ask).FontSize(14).TextColor(t.Text)
			})
		})
		pill := ui.Column(c).Key("pill").AlignSelf(ui.Start).Radius(17).Border(1, t.BorderStrong).Clip().Cursor(ui.CursorPointer).
			OnClick(func() { st.open = !st.open })
		pill.Children(func() {
			ui.Row(c).Height(thinkingLoaderH-2).Padding(0, 14, 0, 10).AlignItems(ui.Center).Children(func() {
				// Rays turning in steps, or the check that lands at the end.
				ui.Box(c).Size(18, 18).Margin(0, 10, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
					cx, cy := r.X+r.W/2, r.Y+r.H/2
					if done {
						land := clamp32(roll*1.3, 0, 1)
						var path ui.Path
						path.MoveTo(cx-5, cy).LineTo(cx-1.5, cy+3.5).LineTo(cx+5.5, cy-4)
						s := 0.6 + 0.4*land
						p.Transform(ui.Scaling(s, s), cx, cy, func() { p.StrokePath(&path, 1.8, t.Text.Alpha(land)) })
						return
					}
					turn := math.Floor(phase*10) / 10 * 2 * math.Pi
					for i := 0; i < 8; i++ {
						a := turn + float64(i)*math.Pi/4
						cs, sn := float32(math.Cos(a)), float32(math.Sin(a))
						fade := 0.3 + 0.7*float32(i)/7
						p.Line(cx+cs*3.5, cy+sn*3.5, cx+cs*7.5, cy+sn*7.5, 1.7, t.Text.Alpha(fade))
					}
				})
				// The label rolls up to the next step.
				ui.Box(c).Key("label").Height(20).Clip().Children(func() {
					cur := o.Steps[idx]
					ink := t.TextSubtle
					if done {
						ink = t.Text
					}
					if idx > 0 && roll < 1 {
						ui.Text(c, o.Steps[idx-1]).Key("prev").Absolute().Left(0).Top(-20 * roll).FontSize(15).TextColor(t.TextSubtle).Opacity(1 - roll).SingleLine()
					}
					label := ui.Box(c).Key(idx).Translate(0, 20*(1-roll)).Opacity(roll)
					label.Children(func() {
						ui.Text(c, cur).FontSize(15).TextColor(ink).SingleLine()
						if !done && !o.NoShine {
							// A light sweeping across the words.
							x := float32(math.Mod(phase*0.8, 1.4)-0.2) * 200
							ui.Box(c).Key("shine").Absolute().Left(x).Top(0).Size(36, 20).Clip().Children(func() {
								ui.Text(c, cur).Absolute().Left(-x).Top(0).FontSize(15).TextColor(t.Text.Alpha(0.85)).SingleLine()
							})
						}
					})
				})
				ui.Icon(c, icons.ChevronDown()).Size(14, 14).TextColor(t.TextSubtle).Margin(0, 0, 0, 10).Rotate(180 * clamp32(k, 0, 1))
			})
			if k < 0.02 {
				return
			}
			ui.Column(c).Key("steps").Height(float32(len(o.Steps))*26*k+6*k).Padding(0, 14, 0, 14).Opacity(clamp32(k, 0, 1)).Children(func() {
				for i, s := range o.Steps {
					ink := t.TextSubtle
					mark := icons.Check()
					switch {
					case i > idx:
						ink = t.TextSubtle.Alpha(0.5)
					case i == idx && !done:
						ink = t.Text
					}
					ui.Row(c).Key(i).Height(26).AlignItems(ui.Center).Children(func() {
						if i < idx || done {
							ui.Icon(c, mark).Size(13, 13).TextColor(t.TextSubtle).Margin(0, 10, 0, 2)
						} else {
							ui.Box(c).Size(5, 5).Radius(2.5).Background(ink).Margin(0, 14, 0, 6)
						}
						ui.Text(c, s).FontSize(13.5).TextColor(ink).SingleLine()
					})
				}
			})
		})
	})
	return col
}

func init() {
	registerDemo(
		"thinking-loader",
		"An AI status pill: rays turn while the label rolls through each step, a check lands at the end; press to see the steps.",
		func(c *ui.Context) { ThinkingLoader(c) },
	)
}
