package particles

import (
	_ "embed"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
)

//go:embed alert.go
var alertSource string

func init() {
	Register(Entry{ID: "p-alert-1", Family: "alert", Description: "Basic alert", View: PAlert1, Source: extractFunc(alertSource, "PAlert1")})
	Register(Entry{ID: "p-alert-2", Family: "alert", Description: "Alert with icon", View: PAlert2, Source: extractFunc(alertSource, "PAlert2")})
	Register(Entry{ID: "p-alert-3", Family: "alert", Description: "Alert with icon and actions", View: PAlert3, Source: extractFunc(alertSource, "PAlert3")})
	Register(Entry{ID: "p-alert-4", Family: "alert", Description: "Info alert", View: PAlert4, Source: extractFunc(alertSource, "PAlert4")})
	Register(Entry{ID: "p-alert-5", Family: "alert", Description: "Success alert", View: PAlert5, Source: extractFunc(alertSource, "PAlert5")})
	Register(Entry{ID: "p-alert-6", Family: "alert", Description: "Warning alert", View: PAlert6, Source: extractFunc(alertSource, "PAlert6")})
	Register(Entry{ID: "p-alert-7", Family: "alert", Description: "Error alert", View: PAlert7, Source: extractFunc(alertSource, "PAlert7")})
	Register(Entry{ID: "p-alert-dialog-1", Family: "alert-dialog", Description: "Basic alert dialog", View: PAlertDialog1, Source: extractFunc(alertSource, "PAlertDialog1")})
	Register(Entry{ID: "p-alert-dialog-2", Family: "alert-dialog", Description: "Alert dialog with bare footer", View: PAlertDialog2, Source: extractFunc(alertSource, "PAlertDialog2")})
}

// PAlert1 renders a basic alert without an icon.
func PAlert1(c *ui.Context) {
	coss.Alert(c, "Heads up!").
		Description("Describe what can be done about it here.")
}

// PAlert2 renders a default alert with an info icon.
func PAlert2(c *ui.Context) {
	coss.Alert(c, "Heads up!").
		Icon(icons.Info()).
		Description("Describe what can be done about it here.")
}

// PAlert3 renders an alert with trailing Dismiss and Ok buttons.
func PAlert3(c *ui.Context) {
	coss.Alert(c, "Heads up!").
		Icon(icons.Info()).
		Description("Describe what can be done about it here.").
		Action(func() {
			coss.Button(c, "Dismiss").Size(coss.SizeXS).Variant(coss.ButtonGhost)
			coss.Button(c, "Ok").Size(coss.SizeXS)
		})
}

// PAlert4 renders an info semantic alert.
func PAlert4(c *ui.Context) {
	coss.Alert(c, "Heads up!").
		Variant(coss.AlertInfo).
		Icon(icons.Info()).
		Description("Describe what can be done about it here.")
}

// PAlert5 renders a success semantic alert.
func PAlert5(c *ui.Context) {
	coss.Alert(c, "Heads up!").
		Variant(coss.AlertSuccess).
		Icon(icons.CircleCheck()).
		Description("Describe what can be done about it here.")
}

// PAlert6 renders a warning semantic alert.
func PAlert6(c *ui.Context) {
	coss.Alert(c, "Heads up!").
		Variant(coss.AlertWarning).
		Icon(icons.TriangleAlert()).
		Description("Describe what can be done about it here.")
}

// PAlert7 renders an error semantic alert.
func PAlert7(c *ui.Context) {
	coss.Alert(c, "Heads up!").
		Variant(coss.AlertError).
		Icon(icons.CircleAlert()).
		Description("Describe what can be done about it here.")
}

// PAlertDialog1 renders a destructive confirmation dialog with a recessed footer.
func PAlertDialog1(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Button(c, "Delete Account").Variant(coss.ButtonDestructiveOutline).OnClick(func() {
			*open = true
		})
		coss.AlertDialog(c, coss.AlertDialogSpec{
			Open:        open,
			Title:       "Are you absolutely sure?",
			Description: "This action cannot be undone. This will permanently delete your account and remove your data from our servers.",
			Footer: func(closeFn func()) {
				coss.Button(c, "Cancel").Variant(coss.ButtonGhost).OnClick(closeFn)
				coss.Button(c, "Delete Account").Variant(coss.ButtonDestructive).OnClick(closeFn)
			},
		})
	})
}

// PAlertDialog2 renders a destructive confirmation dialog with a bare footer.
func PAlertDialog2(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Button(c, "Delete Account").Variant(coss.ButtonDestructiveOutline).OnClick(func() {
			*open = true
		})
		coss.AlertDialog(c, coss.AlertDialogSpec{
			Open:        open,
			Title:       "Are you absolutely sure?",
			Description: "This action cannot be undone. This will permanently delete your account and remove your data from our servers.",
			BareFooter:  true,
			Footer: func(closeFn func()) {
				coss.Button(c, "Cancel").Variant(coss.ButtonGhost).OnClick(closeFn)
				coss.Button(c, "Delete Account").Variant(coss.ButtonDestructive).OnClick(closeFn)
			},
		})
	})
}
