package coss

import (
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// NumberFieldSpec configures a coss-styled numeric stepper input.
type NumberFieldSpec struct {
	Value *float64
	// Min and Max bound the value when Max > Min.
	Min, Max float64
	// Step is how much the buttons and scrubbing change the value; 1 when 0.
	Step float64
	// Format shows the value; the shortest decimal when nil.
	Format   func(float64) string
	Size     InputSize
	Disabled bool
	// Label, when set, shows above the group as NumberFieldScrubArea:
	// dragging it sideways scrubs the value.
	Label string
}

// NumberFieldElement is a NumberField: the field, and its Group, the
// bordered row of the buttons and the input.
type NumberFieldElement struct {
	ui.Element
	Group ui.Element
}

// NumberField renders a full-width numeric input between decrement and
// increment buttons, as ui/number-field.tsx's NumberField.
func NumberField(c *ui.Context, spec NumberFieldSpec) NumberFieldElement {
	t := theme.Of(c)
	step := spec.Step
	if step == 0 {
		step = 1
	}
	set := func(v float64) {
		if spec.Max > spec.Min {
			v = math.Min(math.Max(v, spec.Min), spec.Max)
		}
		*spec.Value = math.Round(v*1e9) / 1e9
	}
	format := spec.Format
	if format == nil {
		format = func(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }
	}
	h, px := float32(32), float32(11)
	switch spec.Size {
	case InputSM:
		h, px = 28, 9
	case InputLG:
		h = 36
	}
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}

	field := ui.Column(c).WidthPercent(100).Gap(8)
	var group ui.Element
	field.Children(func() {
		if spec.Label != "" {
			scrub := ui.Row(c).Cursor(ui.CursorResizeEW)
			acc := ui.Local(scrub, "scrub", func() float32 { return 0 })
			scrub.Children(func() { Label(c, spec.Label).LineHeight(16.0 / 14) })
			if dx, _, ok := scrub.Dragged(); ok && !spec.Disabled {
				*acc += dx
				for ; *acc >= 2; *acc -= 2 {
					set(*spec.Value + step)
				}
				for ; *acc <= -2; *acc += 2 {
					set(*spec.Value - step)
				}
			}
		}
		group = ui.Row(c).
			WidthPercent(100).
			Height(h).
			AlignItems(ui.Center).
			Radius(t.RadiusLG).
			Background(bg).
			Border(1, t.Input).
			FontSize(14).
			TextColor(t.Foreground).
			Clip()
		if group.FocusWithin() {
			group.Border(1, t.Ring)
		}
		if spec.Disabled {
			group.Opacity(0.64)
		} else if !t.Dark {
			group.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04))
		}
		group.Children(func() {
			stepper := func(ic *ui.SVG, delta float64) {
				btn := ui.ButtonBase(c).HeightPercent(100).Padding(0, px).Center().Shrink(0).
					Disabled(spec.Disabled).DimDisabled(false)
				if !spec.Disabled {
					btn.Cursor(ui.CursorPointer)
					if btn.Hovered() {
						btn.Background(t.Accent)
					}
				}
				btn.Children(func() { ui.Icon(c, ic).Size(16, 16) })
				if btn.Clicked() {
					set(*spec.Value + delta)
				}
			}
			stepper(icons.Minus(), -step)
			text := format(*spec.Value)
			in := ui.TextInputBase(c, &text).
				Grow(1).
				MinWidth(0).
				HeightPercent(100).
				Padding(0, px).
				Background(ui.Transparent).
				TextColor(t.Foreground).
				CaretColor(t.Foreground).
				TextAlign(ui.Center).
				FontSize(14).
				FontFeatures("tnum").
				Disabled(spec.Disabled).
				DimDisabled(false)
			if in.Changed() {
				if v, err := strconv.ParseFloat(text, 64); err == nil {
					set(v)
				}
			}
			stepper(icons.Plus(), step)
		})
	})
	return NumberFieldElement{Element: field, Group: group}
}
