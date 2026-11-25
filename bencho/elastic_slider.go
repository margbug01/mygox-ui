package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// ElasticSliderOptions configures ElasticSlider.
type ElasticSliderOptions struct {
	Label    string  // "Opacity" by default
	Width    float32 // 260 by default
	Stretch  float32 // how far the pill can be pulled past an end, 0 to 100; 50 by default
	Bounce   float32 // how lively the snap back is, 0 to 100; 50 by default
	OnChange func(float32)
}

const elasticH = 48

// ElasticSlider is a pill you drag to set *value, 0 to 100. Pull past either
// end and it stretches like rubber, then springs back with a little overshoot.
func ElasticSlider(c *ui.Context, value *float32, opts ...ElasticSliderOptions) ui.Element {
	var o ElasticSliderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Label == "" {
		o.Label = "Opacity"
	}
	if o.Width <= 0 {
		o.Width = 260
	}
	if o.Stretch == 0 {
		o.Stretch = 50
	}
	if o.Bounce == 0 {
		o.Bounce = 50
	}
	t := ThemeFrom(c)

	pill := ui.Box(c).Size(o.Width, elasticH).Cursor(ui.CursorResizeEW)
	px, _, _ := pill.PointerPosition()
	_, _, dragging := pill.Dragged()
	var pull float32 // signed px past an end, rubber-banded
	if dragging || pill.Pressed() {
		maxPull := o.Stretch / 100 * 40
		switch {
		case px < 0:
			pull = -rubberBand(-px, maxPull)
		case px > o.Width:
			pull = rubberBand(px-o.Width, maxPull)
		}
		if next := float32(math.Round(float64(clamp32(px/o.Width, 0, 1) * 100))); next != *value {
			*value = next
			if o.OnChange != nil {
				o.OnChange(next)
			}
		}
	}
	pull = pill.Spring("elastic_pull", pull, ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: float64(o.Bounce) / 100 * 0.6})
	frac := pill.Spring("elastic_value", *value/100, ui.SpringOptions{Duration: 160 * time.Millisecond})
	percent := fmt.Sprintf("%.0f%%", *value)

	pill.Draw(func(p *ui.Painter, r ui.Rect) {
		stretch := float32(math.Abs(float64(pull)))
		h := elasticH - stretch*0.18
		box := ui.Rect{X: r.X + min(pull, 0), Y: r.Y + (elasticH-h)/2, W: o.Width + stretch, H: h}
		fill := box
		fill.W = max(frac*box.W, 0)
		label := func(c ui.Color) {
			_, th := p.MeasureText(0, ui.Span{Text: o.Label, Size: 14, Weight: 400})
			y := box.Y + (box.H-th)/2
			p.RichText(box.X+18, y, 0, ui.Span{Text: o.Label, Size: 14, Weight: 400, Color: c})
			w, _ := p.MeasureText(0, ui.Span{Text: percent, Size: 14, Weight: 500})
			p.RichText(box.X+box.W-18-w, y, 0, ui.Span{Text: percent, Size: 14, Weight: 500, Color: c})
		}
		p.Fill(box, t.Surface, box.H/2)
		label(t.Text)
		if fill.W > 0 {
			p.Clip(fill, min(box.H/2, fill.W/2), func() {
				p.Fill(fill, t.Text, box.H/2)
				label(t.Surface)
			})
		}
	})
	return pill
}

func init() {
	registerDemo(
		"elastic-slider",
		"An opacity pill you drag; pull past either end and it stretches like rubber, then springs back with a little overshoot.",
		func(c *ui.Context) {
			holder := ui.Box(c)
			v := ui.Local(holder, "elastic_value", func() float32 { return 64 })
			holder.Children(func() { ElasticSlider(c, v) })
		},
	)
}
