package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Fieldset groups related fields with an optional legend heading. Like
// coss's, it adds no spacing of its own: callers set Gap.
func Fieldset(c *ui.Context, fn func()) ui.Element {
	return ui.Column(c).WidthPercent(100).Children(fn)
}

// FieldsetLegend renders the heading of a Fieldset.
func FieldsetLegend(c *ui.Context, text string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, text).FontSize(16).FixedLineHeight(24).FontWeight(600).TextColor(t.Foreground)
}
