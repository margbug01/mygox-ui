package coss

import (
	"math"
	"slices"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// sliderThumb is the thumb's size; the control lays out as its 4px track
// and the thumb overhangs it, as in coss.
const sliderThumb = 16

// Slider renders coss's horizontal slider for one value, built on
// ui.SliderBase.
func Slider(c *ui.Context, value *float64, lo, hi float64) ui.Element {
	t := theme.Of(c)
	sl := ui.SliderBase(c, value, lo, hi).
		WidthPercent(100).
		Height(sliderThumb).
		MarginY(-(sliderThumb - 4) / 2).
		PaddingX(sliderThumb / 2).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring)
	v := *value
	sl.Draw(func(p *ui.Painter, r ui.Rect) { drawSlider(p, t, r, []float64{v}, lo, hi, false) })
	return sl
}

// SliderCollision is what a dragged thumb does on meeting another, after
// Base UI's thumbCollisionBehavior.
type SliderCollision uint8

const (
	SliderPush SliderCollision = iota
	SliderNone
	SliderSwap
)

// RangeSlider renders coss's slider with a thumb for each of values, in
// steps of 1, horizontal or vertical. A vertical one is as tall as its
// parent and at least 176.
func RangeSlider(c *ui.Context, values []float64, lo, hi float64, vertical bool, collision SliderCollision) ui.Element {
	t := theme.Of(c)
	sl := ui.Box(c).Cursor(ui.CursorPointer)
	if vertical {
		sl.Width(sliderThumb).MarginX(-(sliderThumb - 4) / 2).HeightPercent(100).MinHeight(176)
	} else {
		sl.WidthPercent(100).Height(sliderThumb).MarginY(-(sliderThumb - 4) / 2)
	}
	active := ui.Local(sl, "active", func() int { return -1 })
	if !sl.Pressed() {
		*active = -1
	} else if px, py, _ := sl.PointerPosition(); len(values) > 0 {
		b := sl.Bounds()
		f := float64(px-sliderThumb/2) / float64(b.W-sliderThumb)
		if vertical {
			f = 1 - float64(py-sliderThumb/2)/float64(b.H-sliderThumb)
		}
		val := math.Round(lo + min(max(f, 0), 1)*(hi-lo))
		if *active < 0 {
			*active = 0
			for i, v := range values {
				if math.Abs(v-val) < math.Abs(values[*active]-val) {
					*active = i
				}
			}
		}
		i := *active
		values[i] = val
		switch collision {
		case SliderPush:
			for j := i + 1; j < len(values); j++ {
				values[j] = max(values[j], values[j-1])
			}
			for j := i - 1; j >= 0; j-- {
				values[j] = min(values[j], values[j+1])
			}
		case SliderSwap:
			slices.Sort(values)
			*active = slices.Index(values, val)
		}
	}
	snap := slices.Clone(values)
	sl.Draw(func(p *ui.Painter, r ui.Rect) { drawSlider(p, t, r, snap, lo, hi, vertical) })
	return sl
}

// drawSlider paints coss's track, indicator and thumbs in r, the thumbs
// kept inside it (thumbAlignment="edge").
func drawSlider(p *ui.Painter, t theme.Tokens, r ui.Rect, values []float64, lo, hi float64, vertical bool) {
	const half = sliderThumb / 2
	length := r.W
	if vertical {
		length = r.H
	}
	// pos maps a value to its thumb's center along the track, from the start
	// (left, or bottom when vertical).
	pos := func(v float64) float32 {
		f := float32(0)
		if hi > lo {
			f = float32(min(max((v-lo)/(hi-lo), 0), 1))
		}
		return half + f*(length-sliderThumb)
	}
	rect := func(from, to, thick float32) ui.Rect {
		if vertical {
			return ui.Rect{X: r.X + (r.W-thick)/2, Y: r.Y + r.H - to, W: thick, H: to - from}
		}
		return ui.Rect{X: r.X + from, Y: r.Y + (r.H-thick)/2, W: to - from, H: thick}
	}
	p.Fill(rect(2, length-2, 4), t.Input, 2)
	from, to := float32(2), pos(values[0])
	if len(values) > 1 {
		from, to = pos(slices.Min(values)), pos(slices.Max(values))
	}
	if to > from {
		p.Fill(rect(from, to, 4), t.Primary, 2)
	}
	for _, v := range values {
		c := pos(v)
		th := rect(c-half, c+half, sliderThumb)
		p.Shadow(th, half, 0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05))
		p.Fill(th, ui.RGB(255, 255, 255), half)
		p.Stroke(th, t.Input, half, 1)
	}
}
