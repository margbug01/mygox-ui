package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Label renders a form control label matching coss label.tsx.
func Label(c *ui.Context, text string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, text).
		SingleLine().
		FontSize(14).
		FixedLineHeight(16).
		FontWeight(500).
		TextColor(t.Foreground)
}
