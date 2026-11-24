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

// AlertDialog opens a confirmation dialog (ui/alert-dialog.tsx: max-w-lg,
// header and footer as Dialog's, no close button) while *spec.Open is true.
func AlertDialog(c *ui.Context, spec AlertDialogSpec) ui.Element {
	t := theme.Of(c)
	return ui.DialogBase(c, spec.Open, func(backdrop, panel ui.Element) {
		dialogBackdrop(backdrop)
		panel.
			WidthPercent(100).
			MaxWidth(512).
			Radius(t.Radius2XL).
			Border(1, t.Border).
			Clip()
		dialogPopup(c, t, panel, dialogContent{
			open:        spec.Open,
			title:       spec.Title,
			description: spec.Description,
			hideClose:   true,
			footer:      spec.Footer,
			bareFooter:  spec.BareFooter,
		})
	})
}
