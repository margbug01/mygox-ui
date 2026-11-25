package particles

import (
	_ "embed"
	"fmt"
	"slices"
	"strconv"

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
	Register(Entry{ID: "p-slider-elastic", Family: "slider", Description: "Slider that stretches past its ends (mygox motion, not on coss.com)", View: PSliderElastic, Source: extractFunc(sliderSource, "PSliderElastic")})
	Register(Entry{ID: "p-slider-23", Family: "slider", Description: "Slider in form", View: PSlider23, Source: extractFunc(sliderSource, "PSlider23")})
}

// sliderHost is the coss.com card's preview width for sliders: full width
// up to 256.
func sliderHost(c *ui.Context) ui.Element { return ui.Column(c).WidthPercent(100).MaxWidth(256) }

// sliderHeader is the label and value row above a slider, mb-2.
func sliderHeader(c *ui.Context, label func(), value string) {
	t := theme.Of(c)
	ui.Row(c).WidthPercent(100).Margin(0, 0, 8, 0).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(4).Children(func() {
		label()
		ui.Text(c, value).TextColor(t.Foreground).FontFeatures("tnum")
	})
}

// sliderCaption is the muted text-xs medium text of slider scales.
func sliderCaption(c *ui.Context, s string) ui.Element {
	return ui.Text(c, s).FontSize(12).FixedLineHeight(16).FontWeight(500).TextColor(theme.Of(c).MutedForeground).SingleLine()
}

// sliderNumber is the small number input beside a slider, kept in sync
// with *v both ways.
func sliderNumber(c *ui.Context, v *float64, lo, hi float64, width float32) {
	box := ui.Box(c).Width(width)
	type state struct {
		shown float64
		text  string
	}
	s := ui.Local(box, "s", func() state { return state{*v, strconv.Itoa(int(*v))} })
	if s.shown != *v {
		*s = state{*v, strconv.Itoa(int(*v))}
	}
	box.Children(func() {
		coss.Input(c, &s.text).Size(coss.InputSM).WidthPercent(100).FontFeatures("tnum").TextAlign(ui.Center)
	})
	if n, err := strconv.ParseFloat(s.text, 64); err == nil && n != *v {
		*v = min(max(n, lo), hi)
		s.shown = *v
	}
}

// PSliderElastic is mygox's own: pulled past an end, the slider stretches
// like rubber, and springs back when let go.
func PSliderElastic(c *ui.Context) {
	host := sliderHost(c)
	val := ui.Local(host, "val", func() float64 { return 50 })
	host.Children(func() { coss.Slider(c, val, 0, 100, coss.SliderOptions{Stretch: 50}).Step(1) })
}

func PSlider1(c *ui.Context) {
	host := sliderHost(c)
	val := ui.Local(host, "val", func() float64 { return 50 })
	host.Children(func() { coss.Slider(c, val, 0, 100) })
}

func PSlider2(c *ui.Context) {
	host := sliderHost(c)
	val := ui.Local(host, "val", func() float64 { return 50 })
	host.Children(func() {
		sliderHeader(c, func() { coss.FieldLabel(c, "Opacity") }, fmt.Sprint(int(*val)))
		coss.Slider(c, val, 0, 100)
	})
}

func PSlider3(c *ui.Context) {
	host := sliderHost(c)
	val := ui.Local(host, "val", func() float64 { return 50 })
	host.Children(func() { coss.Slider(c, val, 0, 100).Disabled(true) })
}

func PSlider4(c *ui.Context) {
	host := sliderHost(c)
	val := ui.Local(host, "val", func() float64 { return 15 })
	host.Children(func() {
		coss.Slider(c, val, 5, 35).Step(1)
		ui.Row(c).WidthPercent(100).Margin(16, 0, 0, 0).Justify(ui.SpaceBetween).Children(func() {
			for _, s := range []string{"5 GB", "20 GB", "35 GB"} {
				sliderCaption(c, s)
			}
		})
	})
}

func PSlider5(c *ui.Context) {
	t := theme.Of(c)
	host := sliderHost(c)
	val := ui.Local(host, "val", func() float64 { return 5 })
	host.Children(func() {
		coss.Slider(c, val, 0, 12).Step(1)
		ui.Row(c).WidthPercent(100).Margin(12, 0, 0, 0).PaddingX(10).Justify(ui.SpaceBetween).Children(func() {
			for i := 0; i <= 12; i++ {
				ui.Column(c).Width(24).MarginX(-12).AlignItems(ui.Center).Gap(8).Children(func() {
					h := float32(4)
					if i%2 != 0 {
						h = 2
					}
					ui.Box(c).Size(1, h).Background(theme.Alpha(t.MutedForeground, 0.72))
					l := sliderCaption(c, strconv.Itoa(i))
					if i%2 != 0 {
						l.Opacity(0)
					}
				})
			}
		})
	})
}

func PSlider6(c *ui.Context) {
	host := sliderHost(c)
	val := ui.Local(host, "val", func() float64 { return 50 })
	host.Children(func() {
		ui.Row(c).WidthPercent(100).Margin(0, 0, 12, 0).Justify(ui.SpaceBetween).Children(func() {
			sliderCaption(c, "Low")
			sliderCaption(c, "High")
		})
		coss.Slider(c, val, 0, 100).Step(10)
	})
}

func sliderRange(c *ui.Context, init []float64, collision coss.SliderCollision) {
	host := sliderHost(c)
	vals := ui.Local(host, "vals", func() []float64 { return init })
	host.Children(func() { coss.RangeSlider(c, *vals, 0, 100, false, collision) })
}

func PSlider7(c *ui.Context)  { sliderRange(c, []float64{25, 75}, coss.SliderPush) }
func PSlider8(c *ui.Context)  { sliderRange(c, []float64{20, 50, 80}, coss.SliderPush) }
func PSlider9(c *ui.Context)  { sliderRange(c, []float64{25, 75}, coss.SliderNone) }
func PSlider10(c *ui.Context) { sliderRange(c, []float64{25, 75}, coss.SliderSwap) }

func PSlider11(c *ui.Context) {
	t := theme.Of(c)
	host := sliderHost(c)
	val := ui.Local(host, "val", func() float64 { return 25 })
	host.Children(func() {
		sliderHeader(c, func() { coss.FieldLabel(c, "Volume") }, fmt.Sprint(int(*val)))
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8).Children(func() {
			ui.Icon(c, icons.VolumeX()).Size(16, 16).TextColor(t.Foreground).Opacity(0.8).Shrink(0)
			ui.Box(c).Grow(1).Children(func() { coss.Slider(c, val, 0, 100) })
			ui.Icon(c, icons.Volume2()).Size(16, 16).TextColor(t.Foreground).Opacity(0.8).Shrink(0)
		})
	})
}

func PSlider12(c *ui.Context) {
	host := sliderHost(c)
	val := ui.Local(host, "val", func() float64 { return 25 })
	host.Children(func() {
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(16).Children(func() {
			ui.Box(c).Grow(1).Children(func() { coss.Slider(c, val, 0, 150).Step(1) })
			sliderNumber(c, val, 0, 150, 48)
		})
	})
}

func PSlider13(c *ui.Context) {
	host := sliderHost(c)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{0, 20} })
	host.Children(func() {
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8).Children(func() {
			v := *vals
			sliderNumber(c, &v[0], 0, v[1], 40)
			ui.Box(c).Grow(1).Children(func() { coss.RangeSlider(c, v, 0, 50, false, coss.SliderPush) })
			sliderNumber(c, &v[1], v[0], 50, 40)
		})
	})
}

func PSlider14(c *ui.Context) {
	host := sliderHost(c)
	val := ui.Local(host, "val", func() float64 { return 100 })
	host.Children(func() {
		ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
			coss.FieldLabel(c, fmt.Sprintf("%d credits/mo", int(*val))).FontFeatures("tnum")
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8).Children(func() {
				coss.IconButton(c, icons.Minus(), "Decrease value").Size(coss.SizeIcon).Variant(coss.ButtonOutline).
					Disabled(*val <= 0).OnClick(func() { *val = max(*val-5, 0) })
				ui.Box(c).Grow(1).Children(func() { coss.Slider(c, val, 0, 200).Step(5) })
				coss.IconButton(c, icons.Plus(), "Increase value").Size(coss.SizeIcon).Variant(coss.ButtonOutline).
					Disabled(*val >= 200).OnClick(func() { *val = min(*val+5, 200) })
			})
		})
	})
}

func formatDollars(n int) string {
	s := strconv.Itoa(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return "$" + s
}

func PSlider15(c *ui.Context) {
	const lo, hi = 5.0, 1240.0
	host := sliderHost(c)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{lo, hi} })
	host.Children(func() {
		price := func(v float64) string {
			if v == hi {
				return formatDollars(int(v)) + "+"
			}
			return formatDollars(int(v))
		}
		coss.Fieldset(c, func() {
			coss.FieldsetLegend(c, fmt.Sprintf("From %s to %s", price((*vals)[0]), price((*vals)[1]))).FontFeatures("tnum")
			coss.RangeSlider(c, *vals, lo, hi, false, coss.SliderPush)
		}).Gap(12)
	})
}

func PSlider16(c *ui.Context) {
	labels := []string{"Awful", "Poor", "Okay", "Good", "Amazing"}
	host := sliderHost(c)
	val := ui.Local(host, "val", func() float64 { return 3 })
	host.Children(func() {
		sliderHeader(c, func() { coss.FieldLabel(c, "Rate your experience") }, labels[int(*val)-1])
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8).Children(func() {
			ui.Text(c, "😡").FontSize(24).FixedLineHeight(32)
			ui.Box(c).Grow(1).Children(func() { coss.Slider(c, val, 1, 5).Step(1) })
			ui.Text(c, "😍").FontSize(24).FixedLineHeight(32)
		})
	})
}

func PSlider17(c *ui.Context) {
	host := ui.Box(c)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{50} })
	host.Children(func() { coss.RangeSlider(c, *vals, 0, 100, true, coss.SliderPush) })
}

func PSlider18(c *ui.Context) {
	host := ui.Row(c).Height(160).AlignItems(ui.Center).Justify(ui.Center)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{25, 75} })
	host.Children(func() { coss.RangeSlider(c, *vals, 0, 100, true, coss.SliderPush) })
}

func PSlider19(c *ui.Context) {
	host := ui.Column(c).AlignItems(ui.Center).Gap(16)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{25} })
	host.Children(func() {
		coss.RangeSlider(c, *vals, 0, 100, true, coss.SliderPush)
		sliderNumber(c, &(*vals)[0], 0, 100, 64)
	})
}

func PSlider20(c *ui.Context) {
	host := ui.Row(c).Height(192).Justify(ui.Center).Gap(32)
	bands := ui.Local(host, "bands", func() [][]float64 { return [][]float64{{2}, {1}, {-1}, {-3}, {2}} })
	host.Children(func() {
		for _, b := range *bands {
			coss.RangeSlider(c, b, -5, 5, true, coss.SliderPush)
		}
	})
}

func PSlider21(c *ui.Context) {
	t := theme.Of(c)
	host := sliderHost(c)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{-2, 4, 2} })
	host.Children(func() {
		coss.Fieldset(c, func() {
			coss.FieldsetLegend(c, "Object position")
			ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
				for i, axis := range []string{"X", "Y", "Z"} {
					ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8).Children(func() {
						ui.Text(c, axis).Width(12).FontWeight(500).TextColor(t.MutedForeground)
						ui.Box(c).Grow(1).Children(func() { coss.Slider(c, &(*vals)[i], -10, 10).Step(1) })
						sliderNumber(c, &(*vals)[i], -10, 10, 64)
					})
				}
			})
			coss.Button(c, "Reset").Variant(coss.ButtonOutline).LeadingIcon(icons.RotateCcw()).
				OnClick(func() { *vals = []float64{0, 0, 0} }).WidthPercent(100)
		}).Gap(16)
	})
}

// sliderPrices are p-slider-22's item prices.
var sliderPrices = []float64{
	80, 95, 110, 125, 130, 140, 145, 150, 155, 165, 175, 185, 195, 205, 215, 225, 235, 245, 255, 260,
	265, 270, 275, 280, 285, 290, 290, 295, 295, 295, 298, 299, 300, 305, 310, 315, 320, 325, 330, 335,
	340, 345, 350, 355, 360, 365, 365, 375, 380, 385, 390, 395, 400, 405, 410, 415, 420, 425, 430, 435,
	440, 445, 450, 455, 460, 465, 470, 475, 480, 485, 490, 495, 495, 498, 499, 500, 500, 500, 515, 530,
	545, 560, 575, 590, 605, 620, 635, 650, 655, 660, 665, 670, 675, 680, 685, 690, 695, 700, 700, 700,
	700, 700, 700, 725, 750, 775, 800, 815, 830, 845, 845, 845, 870, 875, 880, 885, 890, 895, 898, 900,
}

func PSlider22(c *ui.Context) {
	t := theme.Of(c)
	const lo, hi, ticks = 80.0, 900.0, 40
	step := (hi - lo) / ticks
	counts := make([]int, ticks)
	for _, p := range sliderPrices {
		if i := int((p - lo) / step); i < ticks {
			counts[i]++
		}
	}
	maxCount := slices.Max(counts)

	host := sliderHost(c).Gap(16)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{200, 780} })
	host.Children(func() {
		v := *vals
		inRange := 0
		for _, p := range sliderPrices {
			if p >= v[0] && p <= v[1] {
				inRange++
			}
		}
		ui.Column(c).WidthPercent(100).Children(func() {
			ui.Row(c).WidthPercent(100).Height(48).AlignItems(ui.End).PaddingX(12).Children(func() {
				for i, n := range counts {
					rMin, rMax := lo+float64(i)*step, lo+float64(i+1)*step
					col := theme.Alpha(t.Primary, 0.2)
					if inRange > 0 && rMin <= v[1] && rMax >= v[0] {
						col = theme.Alpha(t.Primary, 0.5)
					}
					ui.Box(c).Grow(1).Basis(0).Height(48 * float32(n) / float32(maxCount)).PaddingX(1).Children(func() {
						ui.Box(c).WidthPercent(100).HeightPercent(100).Background(col)
					})
				}
			})
			coss.RangeSlider(c, v, lo, hi, false, coss.SliderPush)
		})
		ui.Row(c).WidthPercent(100).Gap(16).Children(func() {
			for i := range 2 {
				ui.Box(c).Grow(1).Basis(0).Children(func() { sliderPrice(c, &v[i], lo, hi) })
			}
		})
		coss.Button(c, fmt.Sprintf("Show %d items", inRange)).Variant(coss.ButtonOutline).WidthPercent(100)
	})
}

// sliderPrice is p-slider-22's "$" input group, kept in sync with *v.
func sliderPrice(c *ui.Context, v *float64, lo, hi float64) {
	t := theme.Of(c)
	box := ui.Box(c).WidthPercent(100)
	type state struct {
		shown float64
		text  string
	}
	s := ui.Local(box, "s", func() state { return state{*v, strconv.Itoa(int(*v))} })
	if s.shown != *v {
		*s = state{*v, strconv.Itoa(int(*v))}
	}
	box.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() { ui.Text(c, "$").FontSize(14).TextColor(t.MutedForeground) })
			coss.InputGroupInput(c, &s.text).FontFeatures("tnum")
		})
	})
	if n, err := strconv.ParseFloat(s.text, 64); err == nil && n != *v {
		*v = min(max(n, lo), hi)
		s.shown = *v
	}
}

func PSlider23(c *ui.Context) {
	host := sliderHost(c)
	vals := ui.Local(host, "vals", func() []float64 { return []float64{25, 75} })
	host.Children(func() {
		coss.Form(c, func() {
			coss.Fieldset(c, func() {
				ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
					ui.Column(c).WidthPercent(100).Children(func() {
						sliderHeader(c, func() { coss.FieldsetLegend(c, "Volume") }, fmt.Sprintf("%d – %d", int((*vals)[0]), int((*vals)[1])))
						coss.RangeSlider(c, *vals, 0, 100, false, coss.SliderPush)
					})
					coss.FieldDescription(c, "Choose a value between 0 and 100")
				})
			}).Gap(12)
			coss.Button(c, "Submit").WidthPercent(100)
		})
	})
}
