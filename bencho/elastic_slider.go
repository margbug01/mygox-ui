package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// ElasticSliderOptions configures ElasticSlider.
type ElasticSliderOptions struct {
	Min        float32
	Max        float32
	Step       float32
	Width      float32
	Height     float32
	MaxStretch float32
	Label      string
	Unit       string
	Accent     ui.Color
	OnChange   func(float32)
}

// ElasticSlider renders a draggable horizontal slider that stretches with
// rubber-band resistance when pulled past either end and snaps back with a spring.
func ElasticSlider(c *ui.Context, value *float32, opts ...ElasticSliderOptions) ui.Element {
	var o ElasticSliderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Max <= o.Min {
		o.Min = 0
		o.Max = 100
	}
	if o.Width <= 0 {
		o.Width = 240
	}
	if o.Height <= 0 {
		o.Height = 36
	}
	if o.MaxStretch <= 0 {
		o.MaxStretch = 22
	}
	accent := colorOr(o.Accent, t.Accent)

	cur := o.Min
	if value != nil {
		cur = clamp32(*value, o.Min, o.Max)
	}

	wrap := ui.Column(c).
		Width(o.Width).
		Gap(8)

	wrap.Children(func() {
		headerRow := ui.Row(c).
			WidthPercent(100).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center)

		track := ui.Box(c).
			Key("track").
			Width(o.Width).
			Height(o.Height).
			Cursor(ui.CursorPointer)

		bounds := track.LayoutBounds()
		px, _, _ := track.PointerPosition()
		_, _, dragging := track.Dragged()
		active := dragging || track.Pressed()

		var targetStretch float32
		if (active || track.Clicked()) && bounds.W > 0 {
			raw := px / bounds.W
			if raw < 0 {
				targetStretch = -rubberBand(-px, o.MaxStretch)
				raw = 0
			} else if raw > 1 {
				targetStretch = rubberBand(px-bounds.W, o.MaxStretch)
				raw = 1
			}
			next := o.Min + raw*(o.Max-o.Min)
			if o.Step > 0 {
				steps := float32(math.Round(float64((next - o.Min) / o.Step)))
				next = clamp32(o.Min+steps*o.Step, o.Min, o.Max)
			}
			if next != cur {
				cur = next
				if value != nil {
					*value = cur
				}
				if o.OnChange != nil {
					o.OnChange(cur)
				}
			}
		}

		stretch := track.Spring("elastic_stretch", targetStretch, ui.SpringOptions{
			Duration: 300 * time.Millisecond,
			Bounce:   0.42,
		})
		pressTarget := float32(0)
		if active {
			pressTarget = 1
		}
		pressAmt := track.Spring("elastic_press", pressTarget, ui.SpringOptions{
			Duration: 200 * time.Millisecond,
			Bounce:   0.22,
		})

		frac := clamp32((cur-o.Min)/(o.Max-o.Min), 0, 1)

		headerRow.Children(func() {
			label := o.Label
			if label == "" {
				label = "Intensity"
			}
			ui.Text(c, label).FontSize(12).FontWeight(500).TextColor(t.TextMuted)
			ui.Text(c, fmt.Sprintf("%.0f%s", cur, o.Unit)).
				FontSize(12).
				FontWeight(600).
				TextColor(t.Text)
		})

		track.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			absStretch := stretch
			if absStretch < 0 {
				absStretch = -absStretch
			}
			squash := clamp32(absStretch/o.MaxStretch, 0, 1) * 4.5
			pillH := r.H - squash + pressAmt*2
			pillY := r.Y + (r.H-pillH)*0.5
			pillX := r.X
			pillW := r.W + absStretch
			if stretch < 0 {
				pillX = r.X + stretch
			}
			radius := pillH * 0.5

			pillRect := ui.Rect{X: pillX, Y: pillY, W: pillW, H: pillH}
			p.Fill(pillRect, t.SurfaceSubtle, radius)
			p.Stroke(pillRect, t.BorderStrong, radius, 1)

			fillW := pillW * frac
			if fillW > 0 {
				if fillW < pillH && frac > 0 {
					fillW = lerp32(pillH*0.65, pillH, frac)
					if fillW > pillW {
						fillW = pillW
					}
				}
				fillRect := ui.Rect{X: pillX, Y: pillY, W: fillW, H: pillH}
				p.Fill(fillRect, accent, radius)
			}

			knobX := pillX + clamp32(fillW-12, 8, pillW-12)
			knobH := pillH * 0.42
			knobRect := ui.Rect{
				X: knobX,
				Y: pillY + (pillH-knobH)*0.5,
				W: 3.5,
				H: knobH,
			}
			knobCol := ui.Hex("#ffffff").Alpha(0.85)
			if frac < 0.12 && !t.Dark {
				knobCol = t.TextMuted.Alpha(0.65)
			}
			p.Fill(knobRect, knobCol, 2)
		})
	})

	return wrap
}

type elasticDemoState struct {
	val float32
}

func init() {
	registerDemo(
		"elastic-slider",
		"Rubber-band slider that stretches past 0% and 100% and snaps back on release.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_elastic", func() elasticDemoState {
				return elasticDemoState{val: 62}
			})
			holder.Children(func() {
				ElasticSlider(c, &st.val, ElasticSliderOptions{
					Min:   0,
					Max:   100,
					Step:  1,
					Width: 230,
					Label: "Volume",
					Unit:  "%",
				})
			})
		},
	)
}
