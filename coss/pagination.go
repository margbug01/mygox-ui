package coss

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Pagination renders page controls bound to current (1-indexed) out of total pages.
func Pagination(c *ui.Context, current *int, total int) ui.Element {
	t := theme.Of(c)
	if total < 1 {
		total = 1
	}
	cur := 1
	if current != nil {
		if *current < 1 {
			*current = 1
		}
		if *current > total {
			*current = total
		}
		cur = *current
	}
	return ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
		Button(c, "Previous").
			Variant(ButtonGhost).
			Size(SizeSM).
			LeadingIcon(icons.ChevronLeft()).
			Disabled(cur <= 1).
			OnClick(func() {
				if current != nil && *current > 1 {
					*current--
				}
			})

		for p := 1; p <= total; p++ {
			page := p
			if total > 7 && page != 1 && page != total && (page < cur-1 || page > cur+1) {
				if page == 2 || page == total-1 {
					ui.Text(c, "…").FontSize(13).TextColor(t.MutedForeground).Padding(0, 6)
				}
				continue
			}
			v := ButtonGhost
			if page == cur {
				v = ButtonOutline
			}
			Button(c, fmt.Sprintf("%d", page)).
				Variant(v).
				Size(SizeSM).
				MinWidth(28).
				OnClick(func() {
					if current != nil {
						*current = page
					}
				})
		}

		Button(c, "Next").
			Variant(ButtonGhost).
			Size(SizeSM).
			TrailingIcon(icons.ChevronRight()).
			Disabled(cur >= total).
			OnClick(func() {
				if current != nil && *current < total {
					*current++
				}
			})
	})
}
