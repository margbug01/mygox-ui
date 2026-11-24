package coss

import "github.com/egoist/mygo/ui"

// CheckboxGroup stacks related checkboxes 12px apart, each as wide as its
// label, as coss's CheckboxGroup does.
func CheckboxGroup(c *ui.Context, fn func()) ui.Element {
	return ui.Column(c).AlignItems(ui.Start).Gap(12).Children(fn)
}
