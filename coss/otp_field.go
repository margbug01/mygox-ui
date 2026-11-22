package coss

import (
	"unicode"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// OTPField renders a segmented verification code input of length digits bound to code.
func OTPField(c *ui.Context, code *string, length int, groupSize ...int) ui.Element {
	if length <= 0 {
		length = 6
	}
	sepEvery := 0
	if len(groupSize) > 0 && groupSize[0] > 0 {
		sepEvery = groupSize[0]
	}
	t := theme.Of(c)
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}

	raw := ""
	if code != nil {
		raw = *code
	}
	runes := make([]rune, length)
	idx := 0
	for _, r := range raw {
		if idx >= length {
			break
		}
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			runes[idx] = r
			idx++
		}
	}

	return ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
		for i := 0; i < length; i++ {
			if sepEvery > 0 && i > 0 && i%sepEvery == 0 {
				ui.Box(c).Width(8).Height(2).Radius(1).Background(t.Border).Margin(0, 2)
			}
			slotVal := ""
			if runes[i] != 0 {
				slotVal = string(runes[i])
			}
			in := ui.TextInputBase(c, &slotVal).
				Size(36, 40).
				Radius(t.RadiusLG).
				Background(bg).
				Border(1, t.Input).
				TextColor(t.Foreground).
				TextAlign(ui.Center).
				FontSize(16).
				FontWeight(600).
				FocusRingColor(t.Ring)
			if in.Focused() {
				in.Border(1, t.Ring)
			}
			if in.Changed() && code != nil {
				rs := []rune(slotVal)
				if len(rs) > 0 {
					runes[i] = rs[len(rs)-1]
				} else {
					runes[i] = 0
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
