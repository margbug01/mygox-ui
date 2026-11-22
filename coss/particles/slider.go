package particles

import (
	_ "embed"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed slider.go
var sliderSource string

func init() {
	Register(Entry{ID: "p-slider-1", Family: "slider", Description: "Basic slider", View: PSlider1, Source: extractFunc(sliderSource, "PSlider1")})
	Register(Entry{ID: "p-slider-2", Family: "slider", Description: "Slider with label and value", View: PSlider2, Source: extractFunc(sliderSource, "PSlider2")})
	Register(Entry{ID: "p-slider-3", Family: "slider", Description: "Disabled slider", View: PSlider3, Source: extractFunc(sliderSource, "PSlider3")})
	Register(Entry{ID: "p-slider-4", Family: "slider", Description: "Slider with reference labels", View: PSlider4, Source: extractFunc(sliderSource, "PSlider4")})
	Register(Entry{ID: "p-slider-5", Family: "slider", Description: "Slider with ticks", View: PSlider5, Source: extractFunc(sliderSource, "PSlider5")})
	Register(Entry{ID: "p-slider-6", Family: "slider", Description: "Slider with labels above", View: PSlider6, Source: extractFunc(sliderSource, "PSlider6")})
	Register(Entry{ID: "p-slider-7", Family: "slider", Description: "Range slider", View: PSlider7, Source: extractFunc(sliderSource, "PSlider7")})
	Register(Entry{ID: "p-slider-8", Family: "slider", Description: "Slider with 3 thumbs", View: PSlider8, Source: extractFunc(sliderSource, "PSlider8")})
	Register(Entry{ID: "p-slider-9", Family: "slider", Description: "Range slider with collision behavior none", View: PSlider9, Source: extractFunc(sliderSource, "PSlider9")})
	Register(Entry{ID: "p-slider-10", Family: "slider", Description: "Range slider with collision behavior swap", View: PSlider10, Source: extractFunc(sliderSource, "PSlider10")})
	Register(Entry{ID: "p-slider-11", Family: "slider", Description: "Slider with icons", View: PSlider11, Source: extractFunc(sliderSource, "PSlider11")})
	Register(Entry{ID: "p-slider-12", Family: "slider", Description: "Slider with input", View: PSlider12, Source: extractFunc(sliderSource, "PSlider12")})
	Register(Entry{ID: "p-slider-13", Family: "slider", Description: "Range slider with inputs", View: PSlider13, Source: extractFunc(sliderSource, "PSlider13")})
	Register(Entry{ID: "p-slider-14", Family: "slider", Description: "Slider with increment and decrement buttons", View: PSlider14, Source: extractFunc(sliderSource, "PSlider14")})
	Register(Entry{ID: "p-slider-15", Family: "slider", Description: "Price range slider", View: PSlider15, Source: extractFunc(sliderSource, "PSlider15")})
	Register(Entry{ID: "p-slider-16", Family: "slider", Description: "Emoji rating slider", View: PSlider16, Source: extractFunc(sliderSource, "PSlider16")})
	Register(Entry{ID: "p-slider-17", Family: "slider", Description: "Vertical slider", View: PSlider17, Source: extractFunc(sliderSource, "PSlider17")})
	Register(Entry{ID: "p-slider-18", Family: "slider", Description: "Vertical range slider", View: PSlider18, Source: extractFunc(sliderSource, "PSlider18")})
	Register(Entry{ID: "p-slider-19", Family: "slider", Description: "Vertical slider with input", View: PSlider19, Source: extractFunc(sliderSource, "PSlider19")})
	Register(Entry{ID: "p-slider-20", Family: "slider", Description: "Equalizer with vertical sliders", View: PSlider20, Source: extractFunc(sliderSource, "PSlider20")})
	Register(Entry{ID: "p-slider-21", Family: "slider", Description: "Object position sliders with reset", View: PSlider21, Source: extractFunc(sliderSource, "PSlider21")})
	Register(Entry{ID: "p-slider-22", Family: "slider", Description: "Price slider with histogram", View: PSlider22, Source: extractFunc(sliderSource, "PSlider22")})
	Register(Entry{ID: "p-slider-23", Family: "slider", Description: "Slider in form", View: PSlider23, Source: extractFunc(sliderSource, "PSlider23")})
}

type thumbCollision uint8

const (
	collisionClamp thumbCollision = iota
	collisionNone
	collisionSwap
)

func clampFloat(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func multiThumbSlider(c *ui.Context, values []float64, lo, hi float64, mode thumbCollision) ui.Element {
	t := theme.Of(c)
	sl := ui.Box(c).WidthPercent(100).Height(20).Cursor(ui.CursorPointer)
	activeIdx := ui.Local(sl, "activeThumb", func() int { return -1 })

	if sl.Pressed() {
		if px, _, ok := sl.PointerPosition(); ok {
			b := sl.Bounds()
			if b.W > 0 && len(values) > 0 {
				frac := clampFloat(float64(px-b.X)/float64(b.W), 0, 1)
				val := math.Round(lo + frac*(hi-lo))
				if *activeIdx < 0 || *activeIdx >= len(values) {
					best := 0
					bestDist := math.Abs(values[0] - val)
					for i := 1; i < len(values); i++ {
						if d := math.Abs(values[i] - val); d < bestDist {
							best = i
							bestDist = d
						}
					}
					*activeIdx = best
				}
				i := *activeIdx
				switch mode {
				case collisionNone:
					values[i] = val
				case collisionSwap:
					values[i] = val
					sort.Float64s(values)
				default:
					minV, maxV := lo, hi
					if i > 0 {
						minV = values[i-1]
					}
					if i+1 < len(values) {
						maxV = values[i+1]
					}
					values[i] = clampFloat(val, minV, maxV)
				}
			}
		}
	} else {
		*activeIdx = -1
	}

	snap := make([]float64, len(values))
	copy(snap, values)

	sl.Draw(func(p *ui.Painter, r ui.Rect) {
		cy := r.Y + r.H/2
		trackH := float32(4)
		p.Fill(ui.Rect{X: r.X, Y: cy - trackH/2, W: r.W, H: trackH}, t.Muted, trackH/2)
		if len(snap) > 0 && hi > lo {
			minF := float32(1)
			maxF := float32(0)
			for _, v := range snap {
				f := float32(clampFloat((v-lo)/(hi-lo), 0, 1))
				if f < minF {
					minF = f
				}
				if f > maxF {
					maxF = f
				}
			}
			if len(snap) == 1 {
				minF = 0
			}
			x0 := r.X + r.W*minF
			x1 := r.X + r.W*maxF
			if x1 > x0 {
				p.Fill(ui.Rect{X: x0, Y: cy - trackH/2, W: x1 - x0, H: trackH}, t.Primary, trackH/2)
			}
			thumbR := float32(8)
			for _, v := range snap {
				f := float32(clampFloat((v-lo)/(hi-lo), 0, 1))
				cx := r.X + r.W*f
				if cx < r.X+thumbR {
					cx = r.X + thumbR
				}
				if cx > r.X+r.W-thumbR {
					cx = r.X + r.W - thumbR
				}
				thumb := ui.Rect{X: cx - thumbR, Y: cy - thumbR, W: thumbR * 2, H: thumbR * 2}
				p.Fill(thumb, t.Background, thumbR)
				p.Stroke(thumb, t.Primary, thumbR, 1.5)
			}
		}
	})
	return sl
}

func verticalMultiSlider(c *ui.Context, values []float64, lo, hi float64, height float32) ui.Element {
	t := theme.Of(c)
	sl := ui.Box(c).Size(20, height).Cursor(ui.CursorPointer)
	activeIdx := ui.Local(sl, "vActiveThumb", func() int { return -1 })

	if sl.Pressed() {
		if _, py, ok := sl.PointerPosition(); ok {
			b := sl.Bounds()
			if b.H > 0 && len(values) > 0 {
				frac := clampFloat(1-float64(py-b.Y)/float64(b.H), 0, 1)
				val := math.Round(lo + frac*(hi-lo))
				if *activeIdx < 0 || *activeIdx >= len(values) {
					best := 0
					bestDist := math.Abs(values[0] - val)
					for i := 1; i < len(values); i++ {
						if d := math.Abs(values[i] - val); d < bestDist {
							best = i
							bestDist = d
						}
					}
					*activeIdx = best
				}
				i := *activeIdx
				minV, maxV := lo, hi
				if i > 0 {
					minV = values[i-1]
				}
				if i+1 < len(values) {
					maxV = values[i+1]
				}
				values[i] = clampFloat(val, minV, maxV)
			}
		}
	} else {
		*activeIdx = -1
	}

	snap := make([]float64, len(values))
	copy(snap, values)

	sl.Draw(func(p *ui.Painter, r ui.Rect) {
		cx := r.X + r.W/2
		trackW := float32(4)
		p.Fill(ui.Rect{X: cx - trackW/2, Y: r.Y, W: trackW, H: r.H}, t.Muted, trackW/2)
		if len(snap) > 0 && hi > lo {
			minF := float32(1)
			maxF := float32(0)
			for _, v := range snap {
				f := float32(clampFloat((v-lo)/(hi-lo), 0, 1))
				if f < minF {
					minF = f
				}
				if f > maxF {
					maxF = f
				}
			}
			if len(snap) == 1 {
				minF = 0
			}
			yBottom := r.Y + r.H*(1-minF)
			yTop := r.Y + r.H*(1-maxF)
			if yBottom > yTop {
				p.Fill(ui.Rect{X: cx - trackW/2, Y: yTop, W: trackW, H: yBottom - yTop}, t.Primary, trackW/2)
			}
			thumbR := float32(8)
			for _, v := range snap {
				f := float32(clampFloat((v-lo)/(hi-lo), 0, 1))
				cy := r.Y + r.H*(1-f)
				if cy < r.Y+thumbR {
					cy = r.Y + thumbR
				}
				if cy > r.Y+r.H-thumbR {
					cy = r.Y + r.H - thumbR
				}
				thumb := ui.Rect{X: cx - thumbR, Y: cy - thumbR, W: thumbR * 2, H: thumbR * 2}
				p.Fill(thumb, t.Background, thumbR)
				p.Stroke(thumb, t.Primary, thumbR, 1.5)
			}
		}
	})
	return sl
}

func formatDollars(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return "$" + s
	}
	var b strings.Builder
	b.WriteByte('$')
	rem := len(s) % 3
	if rem > 0 {
		b.WriteString(s[:rem])
		if len(s) > rem {
			b.WriteByte(',')
		}
	}
	for i := rem; i < len(s); i += 3 {
		if i > rem {
			b.WriteByte(',')
		}
		b.WriteString(s[i : i+3])
	}
	return b.String()
}

func PSlider1(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() float64 { return 50 })
	host.Children(func() {
		coss.Slider(c, val, 0, 100)
	})
}

func PSlider2(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() float64 { return 50 })
	host.Children(func() {
		coss.Field(c, func() {
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
				coss.FieldLabel(c, "Opacity")
				ui.Text(c, fmt.Sprintf("%d", int(math.Round(*val)))).FontSize(13).TextColor(t.MutedForeground).FontFeatures("tnum")
			})
			coss.Slider(c, val, 0, 100)
		})
	})
}

func PSlider3(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() float64 { return 50 })
	host.Children(func() {
		coss.Slider(c, val, 0, 100).Disabled(true).Opacity(0.5)
	})
}

func PSlider4(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).MaxWidth(280).Gap(10)
	val := ui.Local(host, "val", func() float64 { return 15 })
	host.Children(func() {
		coss.Slider(c, val, 5, 35)
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
			ui.Text(c, "5 GB").FontSize(12).FontWeight(500).TextColor(t.MutedForeground)
			ui.Text(c, "20 GB").FontSize(12).FontWeight(500).TextColor(t.MutedForeground)
			ui.Text(c, "35 GB").FontSize(12).FontWeight(500).TextColor(t.MutedForeground)
		})
	})
}

func PSlider5(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).MaxWidth(300).Gap(8)
	val := ui.Local(host, "val", func() float64 { return 5 })
	host.Children(func() {
		coss.Slider(c, val, 0, 12).Step(1)
		ui.Row(c).WidthPercent(100).AlignItems(ui.Start).Justify(ui.SpaceBetween).Padding(0, 6).Children(func() {
			for i := 0; i <= 12; i++ {
				ui.Column(c).Center().Gap(4).Width(14).Children(func() {
					tickH := float32(6)
					if i%2 != 0 {
						tickH = 3
					}
					ui.Box(c).Size(1, tickH).Background(theme.Alpha(t.MutedForeground, 0.72))
					lbl := ""
					if i%2 == 0 {
						lbl = strconv.Itoa(i)
					}
					ui.Text(c, lbl).SingleLine().FontSize(11).FontWeight(500).TextColor(t.MutedForeground)
				})
			}
		})
	})
}

func PSlider6(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).MaxWidth(280).Gap(8)
	val := ui.Local(host, "val", func() float64 { return 50 })
	host.Children(func() {
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
			ui.Text(c, "Low").FontSize(12).FontWeight(500).TextColor(t.MutedForeground)
			ui.Text(c, "High").FontSize(12).FontWeight(500).TextColor(t.MutedForeground)
		})
		coss.Slider(c, val, 0, 100).Step(10)
	})
}

func PSlider7(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{25, 75} })
	host.Children(func() {
		multiThumbSlider(c, *vals, 0, 100, collisionClamp)
	})
}

func PSlider8(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{20, 50, 80} })
	host.Children(func() {
		multiThumbSlider(c, *vals, 0, 100, collisionClamp)
	})
}

func PSlider9(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{25, 75} })
	host.Children(func() {
		multiThumbSlider(c, *vals, 0, 100, collisionNone)
	})
}

func PSlider10(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{25, 75} })
	host.Children(func() {
		multiThumbSlider(c, *vals, 0, 100, collisionSwap)
	})
}

func PSlider11(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).MaxWidth(280).Gap(8)
	val := ui.Local(host, "val", func() float64 { return 25 })
	host.Children(func() {
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
			coss.FieldLabel(c, "Volume")
			ui.Text(c, fmt.Sprintf("%d", int(math.Round(*val)))).FontSize(13).TextColor(t.MutedForeground).FontFeatures("tnum")
		})
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(10).Children(func() {
			ui.Icon(c, icons.VolumeX()).Size(16, 16).TextColor(t.Foreground).Opacity(0.8).Shrink(0)
			coss.Slider(c, val, 0, 100).Grow(1)
			ui.Icon(c, icons.Volume2()).Size(16, 16).TextColor(t.Foreground).Opacity(0.8).Shrink(0)
		})
	})
}

func PSlider12(c *ui.Context) {
	host := ui.Row(c).WidthPercent(100).MaxWidth(300).AlignItems(ui.Center).Gap(12)
	val := ui.Local(host, "val", func() float64 { return 25 })
	text := ui.Local(host, "text", func() string { return "25" })
	host.Children(func() {
		prev := *val
		coss.Slider(c, val, 0, 150).Grow(1)
		if *val != prev {
			*text = strconv.Itoa(int(math.Round(*val)))
		}
		prevTxt := *text
		coss.Input(c, text).Size(coss.InputSM).Width(56).FontFeatures("tnum").TextAlign(ui.Center)
		if *text != prevTxt {
			if n, err := strconv.ParseFloat(strings.TrimSpace(*text), 64); err == nil {
				*val = clampFloat(n, 0, 150)
			}
		}
	})
}

func PSlider13(c *ui.Context) {
	host := ui.Row(c).WidthPercent(100).MaxWidth(320).AlignItems(ui.Center).Gap(8)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{0, 20} })
	loTxt := ui.Local(host, "loTxt", func() string { return "0" })
	hiTxt := ui.Local(host, "hiTxt", func() string { return "20" })
	host.Children(func() {
		prevLo, prevHi := (*vals)[0], (*vals)[1]
		prevLoTxt, prevHiTxt := *loTxt, *hiTxt

		coss.Input(c, loTxt).Size(coss.InputSM).Width(48).FontFeatures("tnum").TextAlign(ui.Center)
		multiThumbSlider(c, *vals, 0, 50, collisionClamp).Grow(1)
		coss.Input(c, hiTxt).Size(coss.InputSM).Width(48).FontFeatures("tnum").TextAlign(ui.Center)

		if (*vals)[0] != prevLo {
			*loTxt = strconv.Itoa(int(math.Round((*vals)[0])))
		} else if *loTxt != prevLoTxt {
			if n, err := strconv.ParseFloat(strings.TrimSpace(*loTxt), 64); err == nil {
				(*vals)[0] = clampFloat(n, 0, (*vals)[1])
			}
		}
		if (*vals)[1] != prevHi {
			*hiTxt = strconv.Itoa(int(math.Round((*vals)[1])))
		} else if *hiTxt != prevHiTxt {
			if n, err := strconv.ParseFloat(strings.TrimSpace(*hiTxt), 64); err == nil {
				(*vals)[1] = clampFloat(n, (*vals)[0], 50)
			}
		}
	})
}

func PSlider14(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	val := ui.Local(host, "val", func() float64 { return 100 })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, fmt.Sprintf("%d credits/mo", int(math.Round(*val)))).FontFeatures("tnum")
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8).Children(func() {
				coss.IconButton(c, icons.Minus(), "Decrease value").
					Variant(coss.ButtonOutline).
					Disabled(*val <= 0).
					OnClick(func() { *val = clampFloat(*val-5, 0, 200) })
				coss.Slider(c, val, 0, 200).Step(5).Grow(1)
				coss.IconButton(c, icons.Plus(), "Increase value").
					Variant(coss.ButtonOutline).
					Disabled(*val >= 200).
					OnClick(func() { *val = clampFloat(*val+5, 0, 200) })
			})
		})
	})
}

func PSlider15(c *ui.Context) {
	const lo, hi = 5.0, 1240.0
	host := ui.Box(c).WidthPercent(100).MaxWidth(300)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{lo, hi} })
	host.Children(func() {
		v0 := int(math.Round((*vals)[0]))
		v1 := int(math.Round((*vals)[1]))
		endStr := formatDollars(v1)
		if v1 >= int(hi) {
			endStr += "+"
		}
		coss.Fieldset(c, func() {
			coss.FieldsetLegend(c, fmt.Sprintf("From %s to %s", formatDollars(v0), endStr)).FontFeatures("tnum")
			multiThumbSlider(c, *vals, lo, hi, collisionClamp)
		})
	})
}

func PSlider16(c *ui.Context) {
	t := theme.Of(c)
	labels := []string{"Awful", "Poor", "Okay", "Good", "Amazing"}
	host := ui.Column(c).WidthPercent(100).MaxWidth(280).Gap(8)
	val := ui.Local(host, "val", func() float64 { return 3 })
	host.Children(func() {
		idx := int(math.Round(*val)) - 1
		if idx < 0 {
			idx = 0
		}
		if idx >= len(labels) {
			idx = len(labels) - 1
		}
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
			coss.FieldLabel(c, "Rate your experience")
			ui.Text(c, labels[idx]).FontSize(13).TextColor(t.Foreground)
		})
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(10).Children(func() {
			ui.Text(c, "😡").FontSize(20)
			coss.Slider(c, val, 1, 5).Step(1).Grow(1)
			ui.Text(c, "😍").FontSize(20)
		})
	})
}

func PSlider17(c *ui.Context) {
	host := ui.Box(c)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{50} })
	host.Children(func() {
		verticalMultiSlider(c, *vals, 0, 100, 140)
	})
}

func PSlider18(c *ui.Context) {
	host := ui.Box(c)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{25, 75} })
	host.Children(func() {
		verticalMultiSlider(c, *vals, 0, 100, 150)
	})
}

func PSlider19(c *ui.Context) {
	host := ui.Column(c).Center().Gap(12)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{25} })
	txt := ui.Local(host, "txt", func() string { return "25" })
	host.Children(func() {
		prev := (*vals)[0]
		verticalMultiSlider(c, *vals, 0, 100, 140)
		if (*vals)[0] != prev {
			*txt = strconv.Itoa(int(math.Round((*vals)[0])))
		}
		prevTxt := *txt
		coss.Input(c, txt).Size(coss.InputSM).Width(64).FontFeatures("tnum").TextAlign(ui.Center)
		if *txt != prevTxt {
			if n, err := strconv.ParseFloat(strings.TrimSpace(*txt), 64); err == nil {
				(*vals)[0] = clampFloat(n, 0, 100)
			}
		}
	})
}

func PSlider20(c *ui.Context) {
	t := theme.Of(c)
	labels := []string{"60 Hz", "250 Hz", "1k", "4k", "16k"}
	host := ui.Row(c).Center().Gap(24)
	bands := ui.Local(host, "bands", func() [][]float64 {
		return [][]float64{{2}, {1}, {-1}, {-3}, {2}}
	})
	host.Children(func() {
		for i, lbl := range labels {
			ui.Column(c).Center().Gap(8).Children(func() {
				verticalMultiSlider(c, (*bands)[i], -5, 5, 140)
				ui.Text(c, lbl).SingleLine().FontSize(11).TextColor(t.MutedForeground)
			})
		}
	})
}

func PSlider21(c *ui.Context) {
	t := theme.Of(c)
	axes := []string{"X", "Y", "Z"}
	host := ui.Box(c).WidthPercent(100).MaxWidth(300)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{-2, 4, 2} })
	txts := ui.Local(host, "txts", func() []string { return []string{"-2", "4", "2"} })
	host.Children(func() {
		coss.Fieldset(c, func() {
			coss.FieldsetLegend(c, "Object position")
			ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
				for i, axis := range axes {
					idx := i
					ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8).Children(func() {
						ui.Text(c, axis).Width(14).FontSize(12).FontWeight(500).TextColor(t.MutedForeground)
						prev := (*vals)[idx]
						coss.Slider(c, &(*vals)[idx], -10, 10).Step(1).Grow(1)
						if (*vals)[idx] != prev {
							(*txts)[idx] = strconv.Itoa(int(math.Round((*vals)[idx])))
						}
						prevT := (*txts)[idx]
						coss.Input(c, &(*txts)[idx]).Size(coss.InputSM).Width(56).FontFeatures("tnum").TextAlign(ui.Center)
						if (*txts)[idx] != prevT {
							if n, err := strconv.ParseFloat(strings.TrimSpace((*txts)[idx]), 64); err == nil {
								(*vals)[idx] = clampFloat(n, -10, 10)
							}
						}
					})
				}
			})
			coss.Button(c, "Reset").Variant(coss.ButtonOutline).LeadingIcon(icons.RotateCcw()).
				OnClick(func() {
					(*vals)[0], (*vals)[1], (*vals)[2] = 0, 0, 0
					(*txts)[0], (*txts)[1], (*txts)[2] = "0", "0", "0"
				}).WidthPercent(100)
		})
	})
}

var sliderHistogramCounts = []int{
	1, 1, 2, 2, 3, 2, 2, 2, 3, 4,
	6, 4, 3, 3, 4, 3, 2, 3, 3, 5,
	6, 3, 2, 2, 2, 2, 2, 3, 4, 5,
	7, 3, 2, 2, 2, 3, 4, 3, 2, 3,
}

func PSlider22(c *ui.Context) {
	t := theme.Of(c)
	const lo, hi = 80.0, 900.0
	const bins = 40
	step := (hi - lo) / bins
	maxBin := 7

	host := ui.Column(c).WidthPercent(100).MaxWidth(320).Gap(14)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{200, 780} })
	loTxt := ui.Local(host, "loTxt", func() string { return "200" })
	hiTxt := ui.Local(host, "hiTxt", func() string { return "780" })

	host.Children(func() {
		v0, v1 := (*vals)[0], (*vals)[1]
		inRangeCount := 0
		for i, cnt := range sliderHistogramCounts {
			rMin := lo + float64(i)*step
			rMax := rMin + step
			if rMin <= v1 && rMax >= v0 {
				inRangeCount += cnt
			}
		}

		ui.Column(c).WidthPercent(100).Gap(4).Children(func() {
			ui.Row(c).WidthPercent(100).Height(48).AlignItems(ui.End).Padding(0, 8).Gap(1).Children(func() {
				for i, cnt := range sliderHistogramCounts {
					rMin := lo + float64(i)*step
					rMax := rMin + step
					selected := inRangeCount > 0 && rMin <= v1 && rMax >= v0
					col := theme.Alpha(t.Primary, 0.20)
					if selected {
						col = theme.Alpha(t.Primary, 0.55)
					}
					h := float32(cnt) / float32(maxBin) * 44
					if h < 4 {
						h = 4
					}
					ui.Box(c).Grow(1).Height(h).Radius(1).Background(col)
				}
			})
			prevLo, prevHi := (*vals)[0], (*vals)[1]
			multiThumbSlider(c, *vals, lo, hi, collisionClamp)
			if (*vals)[0] != prevLo {
				*loTxt = strconv.Itoa(int(math.Round((*vals)[0])))
			}
			if (*vals)[1] != prevHi {
				*hiTxt = strconv.Itoa(int(math.Round((*vals)[1])))
			}
		})

		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(12).Children(func() {
			prevLoT := *loTxt
			coss.InputGroup(c, func() {
				coss.InputGroupAddon(c, func() { ui.Text(c, "$").SingleLine() })
				coss.InputGroupInput(c, loTxt).FontFeatures("tnum")
			}).Grow(1)
			if *loTxt != prevLoT {
				if n, err := strconv.ParseFloat(strings.TrimSpace(*loTxt), 64); err == nil {
					(*vals)[0] = clampFloat(n, lo, (*vals)[1])
				}
			}
			prevHiT := *hiTxt
			coss.InputGroup(c, func() {
				coss.InputGroupAddon(c, func() { ui.Text(c, "$").SingleLine() })
				coss.InputGroupInput(c, hiTxt).FontFeatures("tnum")
			}).Grow(1)
			if *hiTxt != prevHiT {
				if n, err := strconv.ParseFloat(strings.TrimSpace(*hiTxt), 64); err == nil {
					(*vals)[1] = clampFloat(n, (*vals)[0], hi)
				}
			}
		})

		coss.Button(c, fmt.Sprintf("Show %d items", inRangeCount)).Variant(coss.ButtonOutline).WidthPercent(100)
	})
}

func PSlider23(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(300)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{25, 75} })
	submitted := ui.Local(host, "submitted", func() string { return "" })
	host.Children(func() {
		coss.Form(c, func() {
			coss.Fieldset(c, func() {
				coss.Field(c, func() {
					ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
						coss.FieldsetLegend(c, "Volume")
						ui.Text(c, fmt.Sprintf("%d – %d", int(math.Round((*vals)[0])), int(math.Round((*vals)[1])))).
							FontSize(13).TextColor(t.MutedForeground).FontFeatures("tnum")
					})
					multiThumbSlider(c, *vals, 0, 100, collisionClamp)
					coss.FieldDescription(c, "Choose a value between 0 and 100")
				})
			})
			coss.Button(c, "Submit").OnClick(func() {
				*submitted = fmt.Sprintf("Volume: %d, %d", int(math.Round((*vals)[0])), int(math.Round((*vals)[1])))
			})
			if *submitted != "" {
				ui.Text(c, *submitted).FontSize(12).TextColor(t.MutedForeground)
			}
		})
	})
}
