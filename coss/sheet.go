package coss

import "github.com/egoist/mygo/ui"

// SheetSpec configures a side flyout sheet panel (defaults to right side).
type SheetSpec struct {
	Open        *bool
	Side        DrawerSide
	Title       string
	Description string
	Body        func(close func())
	Footer      func(close func())
}

// Sheet renders a coss-styled side sheet panel.
func Sheet(c *ui.Context, spec SheetSpec) ui.Element {
	side := spec.Side
	if side == DrawerBottom {
		side = DrawerRight
	}
	return Drawer(c, DrawerSpec{
		Open:        spec.Open,
		Side:        side,
		Title:       spec.Title,
		Description: spec.Description,
		Body:        spec.Body,
		Footer:      spec.Footer,
	})
}
