package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Slider renders a coss-styled horizontal slider built on ui.SliderBase.
func Slider(c *ui.Context, value *float64, lo, hi float64) ui.Element {
	t := theme.Of(c)
	sl := ui.SliderBase(c, value, lo, hi).
		WidthPercent(100).
		Height(20).
		AlignItems(ui.Center).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring)

	frac := float32(0)
	if value != nil && hi > lo {
		f := (*value - lo) / (hi - lo)
		if f < 0 {
			f = 0
		}
		if f > 1 {
			f = 1
		}
		frac = float32(f)
	}

	sl.Draw(func(p *ui.Painter, r ui.Rect) {
		cy := r.Y + r.H/2
		trackH := float32(4)
		trackRect := ui.Rect{X: r.X, Y: cy - trackH/2, W: r.W, H: trackH}
		p.Fill(trackRect, t.Muted, trackH/2)
		fillW := r.W * frac
		if fillW > 0 {
			p.Fill(ui.Rect{X: r.X, Y: cy - trackH/2, W: fillW, H: trackH}, t.Primary, trackH/2)
		}
		thumbR := float32(8)
		cx := r.X + fillW
		if cx < r.X+thumbR {
			cx = r.X + thumbR
		}
		if cx > r.X+r.W-thumbR {
			cx = r.X + r.W - thumbR
		}
		thumb := ui.Rect{X: cx - thumbR, Y: cy - thumbR, W: thumbR * 2, H: thumbR * 2}
		p.Fill(thumb, t.Background, thumbR)
		p.Stroke(thumb, t.Primary, thumbR, 1.5)
	})
	return sl
}
