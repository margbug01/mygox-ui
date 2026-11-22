package coss

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// NumberFieldSpec configures a coss-styled numeric stepper input.
type NumberFieldSpec struct {
	Value    *float64
	Min      float64
	Max      float64
	Step     float64
	Decimals int
	Disabled bool
}

// NumberField renders a numeric input with decrement and increment buttons.
func NumberField(c *ui.Context, spec NumberFieldSpec) ui.Element {
	t := theme.Of(c)
	step := spec.Step
	if step == 0 {
		step = 1
	}
	hasBounds := spec.Max > spec.Min

	clamp := func(v float64) float64 {
		if hasBounds {
			if v < spec.Min {
				v = spec.Min
			}
			if v > spec.Max {
				v = spec.Max
			}
		}
		return v
	}

	val := float64(0)
	if spec.Value != nil {
		val = *spec.Value
	}
	fmtStr := "%." + strconv.Itoa(spec.Decimals) + "f"
	text := fmt.Sprintf(fmtStr, val)

	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	row := ui.Row(c).
		Height(32).
		AlignItems(ui.Center).
		Radius(t.RadiusLG).
		Background(bg).
		Border(1, t.Input).
		Clip()
	if spec.Disabled {
		row.Opacity(0.64)
	}
	row.Children(func() {
		IconButton(c, icons.Minus(), "Decrement").
			Variant(ButtonGhost).
			Size(SizeIcon).
			Radius(0).
			Disabled(spec.Disabled || (hasBounds && val <= spec.Min)).
			OnClick(func() {
				if spec.Value != nil {
					*spec.Value = clamp(*spec.Value - step)
				}
			})
		ui.Box(c).Width(1).Height(20).Background(t.Border)
		in := ui.TextInputBase(c, &text).
			Width(72).
			Height(30).
			Padding(0, 8).
			Background(ui.Transparent).
			TextColor(t.Foreground).
			TextAlign(ui.Center).
			FontSize(14).
			Disabled(spec.Disabled)
		if in.Changed() && spec.Value != nil {
			if parsed, err := strconv.ParseFloat(text, 64); err == nil {
				*spec.Value = clamp(parsed)
			}
		}
		ui.Box(c).Width(1).Height(20).Background(t.Border)
		IconButton(c, icons.Plus(), "Increment").
			Variant(ButtonGhost).
			Size(SizeIcon).
			Radius(0).
			Disabled(spec.Disabled || (hasBounds && val >= spec.Max)).
			OnClick(func() {
				if spec.Value != nil {
					*spec.Value = clamp(*spec.Value + step)
				}
			})
	})
	return row
}
