package coss

import (
	"unicode"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// OTPFieldSpec configures an OTPField. The zero value takes 6 digits.
type OTPFieldSpec struct {
	// Length is the number of slots; 6 when 0.
	Length int
	// Group, when set, puts an OTPFieldSeparator after every Group slots.
	Group int
	// Large is size="lg".
	Large bool
	// Placeholder shows in each empty slot that is not focused.
	Placeholder string
	// Mask obscures the characters, as mask.
	Mask bool
	// Accept reports whether a slot takes r; digits only when nil.
	Accept  func(r rune) bool
	Invalid bool
}

// OTPField renders a one-time code input of separate character slots bound
// to code, as ui/otp-field.tsx's OTPField; typing moves to the next slot.
func OTPField(c *ui.Context, code *string, spec OTPFieldSpec) ui.Element {
	t := theme.Of(c)
	length := spec.Length
	if length <= 0 {
		length = 6
	}
	accept := spec.Accept
	if accept == nil {
		accept = unicode.IsDigit
	}
	size, fontSize := float32(32), float32(14)
	if spec.Large {
		size, fontSize = 36, 16
	}
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}

	runes := make([]rune, length)
	i := 0
	for _, r := range *code {
		if i < length && accept(r) {
			runes[i] = r
			i++
		}
	}

	row := ui.Row(c).AlignItems(ui.Center).Gap(8)
	focusNext := ui.Local(row, "focusNext", func() int { return -1 })
	return row.Children(func() {
		for i := range length {
			if spec.Group > 0 && i > 0 && i%spec.Group == 0 {
				OTPFieldSeparator(c)
			}
			slot := ""
			if runes[i] != 0 {
				slot = string(runes[i])
			}
			in := ui.TextInputBase(c, &slot).
				Size(size, size).
				Shrink(0).
				Radius(t.RadiusLG).
				Background(bg).
				Border(1, t.Input).
				TextColor(t.Foreground).
				CaretColor(t.Foreground).
				TextAlign(ui.Center).
				FontSize(fontSize)
			if spec.Mask {
				in.Password()
			}
			switch {
			case in.Focused():
				in.Border(1, t.Ring)
			case spec.Invalid:
				in.Border(1, theme.Alpha(t.Destructive, 0.36))
			case !t.Dark:
				in.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04))
			}
			if spec.Placeholder != "" && !in.Focused() {
				in.Placeholder(spec.Placeholder).PlaceholderColor(t.MutedForeground)
			}
			if *focusNext == i {
				in.Focus()
				*focusNext = -1
			}
			if in.Changed() {
				runes[i] = 0
				for _, r := range slot {
					if accept(r) {
						runes[i] = r
					}
				}
				if runes[i] != 0 && i+1 < length {
					*focusNext = i + 1
				}
				out := make([]rune, 0, length)
				for _, r := range runes {
					if r != 0 {
						out = append(out, r)
					}
				}
				*code = string(out)
			}
		}
	})
}

// OTPFieldSeparator is the short bar between groups of OTPField slots.
func OTPFieldSeparator(c *ui.Context) ui.Element {
	return ui.Box(c).Size(12, 2).Radius(1).Background(theme.Of(c).Input).Shrink(0)
}
