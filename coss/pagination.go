package coss

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Pagination renders page controls bound to current (1-indexed) out of
// total pages: Previous, the pages next to the current one with an
// ellipsis for those left out, and Next, as ui/pagination.tsx.
func Pagination(c *ui.Context, current *int, total int) ui.Element {
	t := theme.Of(c)
	total = max(total, 1)
	cur := 1
	if current != nil {
		*current = min(max(*current, 1), total)
		cur = *current
	}
	ellipsis := func() {
		ui.Row(c).MinWidth(28).Center().Children(func() {
			ui.Icon(c, icons.Ellipsis()).Size(16, 16).TextColor(t.Foreground)
		})
	}
	return ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
		ButtonBox(c, func() {
			ui.Icon(c, icons.ChevronLeft()).Size(16, 16).Opacity(0.8).Margin(0, -2, 0, -4)
			ui.Text(c, "Previous").FontWeight(500)
		}).Variant(ButtonGhost).OnClick(func() {
			if current != nil && *current > 1 {
				*current--
			}
		})
		first, last := max(cur-1, 1), min(cur+1, total)
		if first > 1 {
			ellipsis()
		}
		for p := first; p <= last; p++ {
			page := p
			v := ButtonGhost
			if page == cur {
				v = ButtonOutline
			}
			ButtonBox(c, func() {
				ui.Text(c, fmt.Sprintf("%d", page)).FontWeight(500)
			}).Variant(v).Size(SizeIcon).OnClick(func() {
				if current != nil {
					*current = page
				}
			})
		}
		if last < total {
			ellipsis()
		}
		ButtonBox(c, func() {
			ui.Text(c, "Next").FontWeight(500)
			ui.Icon(c, icons.ChevronRight()).Size(16, 16).Opacity(0.8).Margin(0, -4, 0, -2)
		}).Variant(ButtonGhost).OnClick(func() {
			if current != nil && *current < total {
				*current++
			}
		})
	})
}
