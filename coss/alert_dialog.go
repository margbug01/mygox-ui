package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// AlertDialogSpec configures a coss-styled confirmation dialog.
type AlertDialogSpec struct {
	Open        *bool
	Title       string
	Description string
	BareFooter  bool
	Footer      func(close func())
}

// AlertDialog opens a centered confirmation dialog with backdrop dimming when *spec.Open is true.
func AlertDialog(c *ui.Context, spec AlertDialogSpec) ui.Element {
	t := theme.Of(c)
	black := ui.RGB(0, 0, 0)
	return ui.DialogBase(c, spec.Open, func(backdrop, panel ui.Element) {
		backdrop.Background(theme.Alpha(black, 0.36))
		panel.
			Width(420).
			MaxWidthPercent(92).
			Radius(t.Radius2XL).
			Background(t.Popover).
			TextColor(t.PopoverForeground).
			Border(1, t.Border).
			Shadow(0, 12, 32, 0, theme.Alpha(black, 0.20)).
			Clip().
			Children(func() {
				ui.Column(c).WidthPercent(100).Padding(22, 22, 18, 22).Gap(6).Children(func() {
					if spec.Title != "" {
						ui.Text(c, spec.Title).FontSize(17).FontWeight(600).TextColor(t.PopoverForeground)
					}
					if spec.Description != "" {
						ui.Text(c, spec.Description).FontSize(14).TextColor(t.MutedForeground)
					}
				})
				if spec.Footer != nil {
					closeFn := func() {
						if spec.Open != nil {
							*spec.Open = false
						}
					}
					footer := ui.Row(c).
						WidthPercent(100).
						AlignItems(ui.Center).
						Justify(ui.End).
						Gap(8)
					if spec.BareFooter {
						footer.Padding(4, 22, 18, 22)
					} else {
						footer.Padding(14, 22).
							Background(theme.Alpha(t.Muted, 0.72)).
							BorderWidth(1, 0, 0, 0).
							BorderColor(t.Border)
					}
					footer.Children(func() { spec.Footer(closeFn) })
				}
			})
	})
}
