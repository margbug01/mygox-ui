package bencho

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
)

// ProgressOptions configures Progress. bencho.dev's block has no props;
// these only let a caller start it elsewhere and hear it change.
type ProgressOptions struct {
	Value    float32 // where it starts, 0 to 100; 66 by default
	OnChange func(value int)
}

const (
	progressW     = 300
	progressH     = 176
	progressTicks = 34
	progressFirst = 14 // the first tick's centre
	progressPitch = 8.23
	progressFoot  = 163.5
	progressLift  = 20 // how much taller a tick stands once it is reached
)

// progressRest is each tick's height before it is reached: a slow swell,
// a dip and a fall.
var progressRest = [progressTicks]float32{
	50.5, 53.5, 57.5, 62.5, 62.5, 63.5, 62.5, 61.5, 57.5, 53.5, 48, 44.5, 42, 37.5, 37.5, 37.5, 37.5,
	41.5, 45.5, 49.5, 52.5, 57.5, 55.5, 59.5, 59.5, 59.5, 58, 55.5, 51.5, 47.5, 43.5, 39.5, 35.5, 35.5,
}

// Progress is a tick slider made of a row of small marks. Scrubbing along
// the ticks sets the value: the ticks it reaches turn ink and stand up
// taller as the pointer passes, and the figure follows. It is coss's
// Slider, in steps of 1 over the ticks, drawn as the ticks.
func Progress(c *ui.Context, opts ...ProgressOptions) ui.Element {
	var o ProgressOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	start := float32(66)
	if o.Value != 0 {
		start = clamp32(o.Value, 0, 100)
	}
	t := ThemeFrom(c)
	holder := ui.Box(c)
	value := ui.Local(holder, "progress", func() float64 { return float64(int(start)) })
	heard := ui.Local(holder, "progress_heard", func() float64 { return *value })
	var lift [progressTicks]float32
	draw := func(p *ui.Painter, r ui.Rect, _ float32) {
		for i, h := range progressRest {
			k := clamp32(lift[i], 0, 1)
			hh := h + progressLift*lift[i]
			x := r.X + progressFirst + float32(i)*progressPitch
			p.Fill(ui.Rect{X: x - 2, Y: r.Y + progressFoot - hh, W: 4, H: hh}, lerpColor(t.SurfaceHover, t.Text, k), 2)
		}
	}
	holder.Children(func() {
		// The value runs from half a pitch before the first tick to half a
		// pitch past the last: 0 before the first tick, 100 past the last.
		left := float32(progressFirst - progressPitch/2)
		card := coss.Slider(c, value, 0, 100, coss.SliderOptions{Draw: draw}).Step(1).
			Size(progressW, progressH).MarginY(0).Padding(0, progressW-left-progressPitch*progressTicks, 0, left).
			Radius(20).Background(t.Surface).Cursor(ui.CursorResizeEW)
		if *value != *heard {
			*heard = *value
			if o.OnChange != nil {
				o.OnChange(int(*value))
			}
		}
		reached := int(float32(*value) / 100 * progressTicks)
		for i := range lift {
			lift[i] = card.Spring(fmt.Sprint("tick", i), b2f(i < reached), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.35})
		}
		card.Children(func() {
			ui.Row(c).Absolute().Left(12).Top(8).Height(52).AlignItems(ui.End).Children(func() {
				ui.Text(c, fmt.Sprint(int(*value))).FontSize(40).FontWeight(300).TextColor(t.Text).Height(52)
				ui.Text(c, "%").FontSize(16).TextColor(lerpColor(t.TextSubtle, t.Surface, 0.35)).Margin(0, 0, 8, 5)
			})
		})
	})
	return holder
}

func init() {
	registerDemo(
		"progress",
		"A tick slider made of small marks; scrub along the ticks to set a value and watch it change as you go.",
		func(c *ui.Context) { Progress(c) },
	)
}
