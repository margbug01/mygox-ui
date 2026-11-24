package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Kbd renders a keyboard shortcut keycap matching coss kbd.tsx.
func Kbd(c *ui.Context, key string) ui.Element {
	return kbdBox(c).Children(func() {
		ui.Text(c, key).SingleLine()
	})
}

// KbdIcon renders a keycap holding an icon at kbd.tsx's svg size-3.
func KbdIcon(c *ui.Context, icon *ui.SVG) ui.Element {
	return kbdBox(c).Children(func() {
		ui.Icon(c, icon).Size(12, 12)
	})
}

func kbdBox(c *ui.Context) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		Center().
		Shrink(0).
		Height(20).
		MinWidth(20).
		Padding(0, 4).
		Radius(4).
		Gap(4).
		Background(t.Muted).
		TextColor(t.MutedForeground).
		FontSize(12).
		FontWeight(500)
}

// KbdGroup arranges multiple Kbd badges horizontally with tight spacing.
func KbdGroup(c *ui.Context, fn func()) ui.Element {
	return ui.Row(c).AlignItems(ui.Center).Gap(4).Shrink(0).Children(fn)
}
