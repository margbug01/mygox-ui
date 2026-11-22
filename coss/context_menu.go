package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// ContextMenu attaches a coss-styled right-click context menu to target using ui.ContextMenuBase.
func ContextMenu(c *ui.Context, target ui.Element, open *bool, items []MenuEntry) {
	t := theme.Of(c)
	parts := ui.ContextMenuBase(c, target, open)
	renderMenuPopup(c, t, parts, items)
}
