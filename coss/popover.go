package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Popover attaches a floating popover to anchor while *open is true, as
// ui/popover.tsx PopoverPopup: rounded-lg border bg-popover shadow-lg/5,
// p-4, 14/20 text.
func Popover(c *ui.Context, anchor ui.Element, open *bool, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		panel.
			Padding(16).
			Radius(t.RadiusLG).
			Background(t.Popover).
			TextColor(t.PopoverForeground).
			Border(1, t.Border).
			Shadow(0, 10, 15, -3, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
			FontSize(14).
			FixedLineHeight(20).
			Children(fn)
	})
}

// PopoverTitle is PopoverTitle: font-semibold, leading-none, at size
// (text-lg 18 upstream; particles pick text-base 16 or text-sm 14).
func PopoverTitle(c *ui.Context, title string, size float32) ui.Element {
	return ui.Text(c, title).FontSize(size).FixedLineHeight(size).FontWeight(600)
}

// PopoverDescription is PopoverDescription: muted 14/20 text.
func PopoverDescription(c *ui.Context, text string) ui.Element {
	return ui.Text(c, text).Wrap().TextColor(theme.Of(c).MutedForeground)
}
