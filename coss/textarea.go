package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Textarea creates a coss-styled multi-line text input built on ui.TextAreaBase.
func Textarea(c *ui.Context, value *string) ui.Element {
	t := theme.Of(c)
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	el := ui.TextAreaBase(c, value).
		WidthPercent(100).
		MinHeight(72).
		Padding(8, 11).
		Radius(t.RadiusLG).
		Background(bg).
		TextColor(t.Foreground).
		PlaceholderColor(theme.Alpha(t.MutedForeground, 0.72)).
		CaretColor(t.Foreground).
		FontSize(14).
		Border(1, t.Input).
		FocusRingColor(t.Ring)
	if el.Focused() {
		el.Border(1, t.Ring)
	}
	return el
}
