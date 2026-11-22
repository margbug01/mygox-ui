package coss

import "github.com/egoist/mygo/ui"

// Collapsible renders a coss-styled disclosure built on ui.CollapsibleBase.
func Collapsible(c *ui.Context, open *bool, trigger func(p ui.CollapsibleParts), panel func()) ui.Element {
	root := ui.Column(c).WidthPercent(100)
	root.Children(func() {
		p := ui.CollapsibleBase(c, open)
		if trigger != nil {
			trigger(p)
		}
		if body := p.Panel(func() {
			if panel != nil {
				panel()
			}
		}); body != (ui.Element{}) {
			body.WidthPercent(100)
		}
	})
	return root
}
