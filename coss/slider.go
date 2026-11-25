package coss

import (
	"math"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// sliderThumb is the thumb's size; the control lays out as its 4px track
// and the thumb overhangs it, as in coss.
const sliderThumb = 16

// SliderOptions are a Slider's optional motion and look.
type SliderOptions struct {
	// Stretch lets the slider be pulled past either end, stretching like
	// rubber, up to 40px at 100, and squashing as it does; let go, it
	// springs back, and the value glides to each new place. 0 is off.
	Stretch float32
	// Bounce is how lively the snap back is, 0 to 100; 50 when 0.
	Bounce float32
	// Draw paints the slider in another design: r is its box, stretched
	// and squashed by the pull, and frac how far along the value is.
	Draw func(p *ui.Painter, r ui.Rect, frac float32)
}

// Slider renders coss's horizontal slider for one value, built on
// ui.SliderBase.
func Slider(c *ui.Context, value *float64, lo, hi float64, opts ...SliderOptions) ui.Element {
	var o SliderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := theme.Of(c)
	sl := ui.SliderBase(c, value, lo, hi).
		WidthPercent(100).
		Height(sliderThumb).
		MarginY(-(sliderThumb - 4) / 2).
		PaddingX(sliderThumb / 2).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring)
	v := *value
	if o.Stretch == 0 && o.Draw == nil {
		sl.Draw(func(p *ui.Painter, r ui.Rect) { drawSlider(p, t, r, []float64{v}, lo, hi, false) })
		return sl
	}
	frac := float32(0)
	if hi > lo {
		frac = float32((v - lo) / (hi - lo))
	}
	var pull float32 // px past an end, signed, rubber-banded
	if o.Stretch > 0 {
		if o.Bounce == 0 {
			o.Bounce = 50
		}
		px, _, _ := sl.PointerPosition()
		if _, _, dragging := sl.Dragged(); dragging || sl.Pressed() {
			w, most := sl.Bounds().W, o.Stretch/100*40
			switch {
			case px < 0:
				pull = -rubberBand(-px, most)
			case px > w:
				pull = rubberBand(px-w, most)
			}
		}
		pull = sl.Spring("elastic_pull", pull, ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: float64(o.Bounce) / 100 * 0.6})
		frac = sl.Spring("elastic_value", frac, ui.SpringOptions{Duration: 160 * time.Millisecond})
	}
	draw := o.Draw
	if draw == nil {
		draw = func(p *ui.Painter, r ui.Rect, frac float32) {
			drawSlider(p, t, r, []float64{lo + float64(frac)*(hi-lo)}, lo, hi, false)
		}
	}
	sl.Draw(func(p *ui.Painter, r ui.Rect) {
		stretch := float32(math.Abs(float64(pull)))
		h := r.H - stretch*0.18
		draw(p, ui.Rect{X: r.X + min(pull, 0), Y: r.Y + (r.H-h)/2, W: r.W + stretch, H: h}, frac)
	})
	return sl
}

// rubberBand eases over, px past an end, so that it nears limit but never
// reaches it.
func rubberBand(over, limit float32) float32 {
	if limit <= 0 {
		return 0
	}
	return limit * (1 - 1/(1+over/limit))
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
