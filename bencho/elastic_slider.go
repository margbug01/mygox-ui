package bencho

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
)

// ElasticSliderOptions configures ElasticSlider.
type ElasticSliderOptions struct {
	Label   string  // "Opacity" by default
	Width   float32 // 260 by default
	Stretch float32 // how far the pill can be pulled past an end, 0 to 100; 50 by default
	Bounce  float32 // how lively the snap back is, 0 to 100; 50 by default
}

// ElasticSlider is a pill you drag to set *value, 0 to 100: coss's Slider
// with Stretch, drawn as bencho's pill, filled up to the value, with the
// label and the percent inside.
func ElasticSlider(c *ui.Context, value *float64, opts ...ElasticSliderOptions) ui.Element {
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
	t := ThemeFrom(c)
	draw := func(p *ui.Painter, box ui.Rect, frac float32) {
		percent := fmt.Sprintf("%.0f%%", *value)
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
	}
	return coss.Slider(c, value, 0, 100, coss.SliderOptions{Stretch: o.Stretch, Bounce: o.Bounce, Draw: draw}).
		Step(1).Size(o.Width, 48).MarginY(0).PaddingX(0).Cursor(ui.CursorResizeEW)
}

func init() {
	registerDemo(
		"elastic-slider",
		"An opacity pill you drag; pull past either end and it stretches like rubber, then springs back with a little overshoot.",
		func(c *ui.Context) {
			holder := ui.Box(c)
			v := ui.Local(holder, "elastic_value", func() float64 { return 64 })
			holder.Children(func() { ElasticSlider(c, v) })
		},
	)
}
