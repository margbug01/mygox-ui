package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// ThinkingOptions configures ThinkingBlock.
type ThinkingOptions struct {
	Pace   time.Duration // how long each step of the reasoning takes, 500ms to 2.4s; 1.1s by default
	Corner float32       // 0 to 26; 20 by default; the trace sits on no card, so it rounds only the step highlight
}

type thinkingState struct {
	start  time.Time
	open   bool // the folded trace opened again
	opened bool // open was set by a press, not by thinking
}

// ThinkingDefaultSteps is the reasoning bencho.dev shows.
func ThinkingDefaultSteps() []string {
	return []string{
		"Reading the request: a stack of cards that opens into a list.",
		"Three should peek out underneath, a little smaller each.",
		"Opening fans them into a column and keeps their order.",
		"A spring with a little bounce will make it feel physical.",
	}
}

// ThinkingBlock is an AI thinking state with its reasoning trace. Steps arrive
// one by one under a shimmering Thinking, the newest dark and the earlier ones
// greyed along a thread; when the last is done the trace folds into a single
// "Thought for" line that opens it up again.
func ThinkingBlock(c *ui.Context, steps []string, opts ...ThinkingOptions) ui.Element {
	var o ThinkingOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Pace == 0 {
		o.Pace = 1100 * time.Millisecond
	}
	t := ThemeFrom(c)
	col := ui.Column(c).Width(300).MinHeight(250) // the trace grows down from a still heading
	st := ui.Local(col, "thinking", func() thinkingState { return thinkingState{start: c.Now()} })
	now := c.Now()
	since := now.Sub(st.start)
	shown := min(len(steps), int(since/o.Pace)+1)
	done := since >= time.Duration(len(steps))*o.Pace
	if !done {
		c.AnimationFrame()
	}
	open := !done || st.open
	k := clamp32(col.Spring("open", b2f(open), ui.SpringOptions{Duration: 420 * time.Millisecond}), 0, 1)
	secs := int(math.Round(float64(time.Duration(len(steps))*o.Pace) / float64(time.Second)))

	col.Children(func() {
		head := ui.Row(c).Key("head").Height(28).AlignItems(ui.Center).Cursor(ui.CursorPointer).OnClick(func() {
			if done {
				st.open = !st.open
			}
		})
		head.Children(func() {
			label := "Thinking"
			ink := t.TextSubtle
			if done {
				label = fmt.Sprintf("Thought for %ds", secs)
			} else {
				// A light moving through the word while it works.
				ink = lerpColor(t.TextSubtle, t.Text, float32(0.5+0.5*math.Sin(phaseSeconds(now)*4)))
			}
			ui.Text(c, label).FontSize(15.5).FontWeight(500).TextColor(ink)
			ui.Icon(c, icons.ChevronDown()).Size(15, 15).TextColor(t.TextSubtle).Margin(0, 0, 0, 5).Rotate(-90 * (1 - k))
		})
		if k < 0.01 {
			return
		}
		trace := ui.Column(c).Key("trace").Margin(8, 0, 0, 0).Gap(12).Opacity(k).Clip()
		if done {
			trace.Height(float32(len(steps)) * 50 * k)
		}
		trace.Children(func() {
			for i := 0; i < shown; i++ {
				current := i == shown-1 && !done
				arrive := clamp32(float32(since-time.Duration(i)*o.Pace)/float32(400*time.Millisecond), 0, 1)
				ink := t.TextSubtle
				dot := lerpColor(t.SurfaceHover, t.BorderStrong, 0.3)
				if current {
					ink, dot = t.Text, t.Text
				}
				last := i == shown-1
				row := ui.Row(c).Key(i).AlignItems(ui.Start).Opacity(arrive).Translate(0, (1-arrive)*6)
				row.Children(func() {
					ui.Box(c).Width(22).AlignSelf(ui.Stretch).Draw(func(p *ui.Painter, r ui.Rect) {
						p.Fill(ui.Rect{X: r.X + 2, Y: r.Y + 7, W: 7, H: 7}, dot, 3.5)
						if !last {
							// The thread down to the next step.
							p.Fill(ui.Rect{X: r.X + 5, Y: r.Y + 17, W: 1, H: r.H + 4}, lerpColor(t.SurfaceHover, t.BorderStrong, 0.2), 0)
						}
					})
					ui.Text(c, steps[i]).FontSize(14.5).LineHeight(1.35).TextColor(ink).Grow(1).MinWidth(0)
				})
			}
		})
	})
	return col
}

func init() {
	registerDemo(
		"thinking",
		"An AI reasoning trace whose steps arrive one by one, then fold into a single line that opens them again.",
		func(c *ui.Context) { ThinkingBlock(c, ThinkingDefaultSteps()) },
	)
}
