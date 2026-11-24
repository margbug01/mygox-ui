package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Meter renders a coss Meter for value in [lo, hi]: a full-width column
// with, when label or valueText is set, a label row 8px over an 8px
// bg-input track filled in the primary color.
func Meter(c *ui.Context, value, lo, hi float64, label, valueText string) ui.Element {
	t := theme.Of(c)
	pct := float32(0)
	if hi > lo {
		pct = float32(min(max((value-lo)/(hi-lo), 0), 1) * 100)
	}
	return ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
		if label != "" || valueText != "" {
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(8).Children(func() {
				ui.Text(c, label).FontSize(14).FixedLineHeight(20).FontWeight(500).TextColor(t.Foreground)
				ui.Text(c, valueText).FontSize(14).FixedLineHeight(20).FontFeatures("tnum").TextColor(t.Foreground)
			})
		}
		ui.Box(c).WidthPercent(100).Height(8).Background(t.Input).Clip().Children(func() {
			ui.Box(c).HeightPercent(100).WidthPercent(pct).Background(t.Primary)
		})
	})
}
