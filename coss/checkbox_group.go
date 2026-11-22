package coss

import "github.com/egoist/mygo/ui"

// CheckboxGroup arranges related checkboxes vertically with standard spacing.
func CheckboxGroup(c *ui.Context, fn func()) ui.Element {
	return ui.Column(c).WidthPercent(100).Gap(10).Children(fn)
}
