package coss

import "github.com/egoist/mygo/ui"

// Form renders a vertical form container with standard spacing matching coss form.tsx.
func Form(c *ui.Context, fn func()) ui.Element {
	return ui.Column(c).WidthPercent(100).Gap(16).Children(fn)
}
