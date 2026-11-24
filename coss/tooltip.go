package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Side is the side of an anchor a popup opens on.
type Side uint8

const (
	SideTop Side = iota
	SideRight
	SideBottom
	SideLeft
)

// Tooltip attaches a hover/focus tooltip to anchor, as ui/tooltip.tsx
// TooltipPopup: rounded-md border bg-popover 12/16 text, px-2 py-1,
// shadow-md/5, 4px from anchor on side (top when not given). Where side
// has no room, it opens on the other side.
func Tooltip(c *ui.Context, anchor ui.Element, text string, side ...Side) ui.Element {
	t := theme.Of(c)
	s := SideTop
	if len(side) > 0 {
		s = side[0]
	}
	return ui.TooltipBase(c, anchor, func(tip ui.Element) {
		switch s {
		case SideRight:
			tip.AttachTo(anchor, ui.AnchorRight, ui.AnchorLeft).Margin(0, 0, 0, 4)
		case SideBottom:
			tip.AttachTo(anchor, ui.AnchorBottom, ui.AnchorTop).Margin(4, 0, 0, 0)
		case SideLeft:
			tip.AttachTo(anchor, ui.AnchorLeft, ui.AnchorRight).Margin(0, 4, 0, 0)
		default:
			tip.AttachTo(anchor, ui.AnchorTop, ui.AnchorBottom).Margin(0, 0, 4, 0)
		}
		tip.
			Padding(4, 8).
			Radius(t.RadiusMD).
			Background(t.Popover).
			TextColor(t.PopoverForeground).
			Border(1, t.Border).
			FontSize(12).
			FixedLineHeight(16).
			Shadow(0, 4, 6, -1, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
			Children(func() {
				ui.Text(c, text).Wrap().Balance()
			})
	})
}
