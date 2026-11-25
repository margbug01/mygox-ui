package bencho

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
)

// LiquidToggle is a squishy switch: tap it, or drag its thumb across, and
// the thumb stretches like a drop of liquid as it flies, while the track
// fills with ink and the thumb turns white. *on is its state. It is coss's
// Liquid Switch in bencho's look, 92 by 46.
func LiquidToggle(c *ui.Context, on *bool) ui.Element {
	t := ThemeFrom(c)
	off := t.BorderStrong
	if t.Dark {
		off = t.TextSubtle
	}
	return coss.Switch(c, on, coss.SwitchOptions{Liquid: true, Style: &coss.SwitchStyle{
		Width: 92, Height: 46, Thumb: 36, Inset: 5,
		Track: t.Surface, TrackOn: t.Text, ThumbOff: off, ThumbOn: t.Surface,
	}}).Label("Toggle")
}

func init() {
	registerDemo(
		"liq-toggle",
		"A squishy toggle: tap or drag it and the thumb stretches like a drop of liquid as it flies while the track fills with ink.",
		func(c *ui.Context) {
			holder := ui.Box(c)
			on := ui.Local(holder, "liquid_toggle_on", func() bool { return false })
			holder.Children(func() { LiquidToggle(c, on) })
		},
	)
}
