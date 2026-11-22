package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// DialogSpec configures a coss-styled modal dialog.
type DialogSpec struct {
	Open        *bool
	Title       string
	Description string
	Width       float32
	HideClose   bool
	Body        func(close func())
	Footer      func(close func())
}

// Dialog renders a coss-styled modal dialog built on ui.DialogBase.
func Dialog(c *ui.Context, spec DialogSpec) ui.Element {
	t := theme.Of(c)
	black := ui.RGB(0, 0, 0)
	w := spec.Width
	if w <= 0 {
		w = 460
	}
	closeFn := func() {
		if spec.Open != nil {
			*spec.Open = false
		}
	}
	return ui.DialogBase(c, spec.Open, func(backdrop, panel ui.Element) {
		backdrop.Background(theme.Alpha(black, 0.36))
		panel.
			Width(w).
			MaxWidthPercent(92).
			Radius(t.Radius2XL).
			Background(t.Popover).
			TextColor(t.PopoverForeground).
			Border(1, t.Border).
			Shadow(0, 14, 36, 0, theme.Alpha(black, 0.22)).
			Clip().
			Children(func() {
				if spec.Title != "" || spec.Description != "" || !spec.HideClose {
					ui.Row(c).WidthPercent(100).AlignItems(ui.Start).Justify(ui.SpaceBetween).Padding(20, 20, 12, 20).Gap(12).Children(func() {
						ui.Column(c).Grow(1).MinWidth(0).Gap(4).Children(func() {
							if spec.Title != "" {
								ui.Text(c, spec.Title).FontSize(17).FontWeight(600).TextColor(t.PopoverForeground)
							}
							if spec.Description != "" {
								ui.Text(c, spec.Description).FontSize(13).TextColor(t.MutedForeground)
							}
						})
						if !spec.HideClose {
							IconButton(c, icons.X(), "Close").Size(SizeIconSM).Variant(ButtonGhost).OnClick(closeFn)
						}
					})
				}
				if spec.Body != nil {
					ui.Column(c).WidthPercent(100).Padding(8, 20, 20, 20).Gap(12).Children(func() {
						spec.Body(closeFn)
					})
				}
				if spec.Footer != nil {
					ui.Row(c).
						WidthPercent(100).
						AlignItems(ui.Center).
						Justify(ui.End).
						Gap(8).
						Padding(14, 20).
						Background(theme.Alpha(t.Muted, 0.72)).
						BorderWidth(1, 0, 0, 0).
						BorderColor(t.Border).
						Children(func() { spec.Footer(closeFn) })
				}
			})
	})
}
