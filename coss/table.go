package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// TableColumnSpec defines a header column in a coss Table.
type TableColumnSpec struct {
	Title string
	Width float32 // 0 means flexible Grow(1)
	Align ui.Align
}

// Table renders a clean coss-styled data table with header and custom cell renderer.
func Table(c *ui.Context, columns []TableColumnSpec, rowCount int, cell func(row, col int)) ui.Element {
	t := theme.Of(c)
	return ui.Column(c).
		WidthPercent(100).
		Radius(t.RadiusXL).
		Background(t.Card).
		Border(1, t.Border).
		Clip().
		Children(func() {
			// Header
			ui.Row(c).
				WidthPercent(100).
				AlignItems(ui.Center).
				Padding(10, 14).
				Gap(12).
				Background(theme.Alpha(t.Muted, 0.6)).
				BorderWidth(0, 0, 1, 0).
				BorderColor(t.Border).
				Children(func() {
					for _, col := range columns {
						hdr := ui.Row(c).AlignItems(ui.Center).Justify(col.Align)
						if col.Width > 0 {
							hdr.Width(col.Width).Shrink(0)
						} else {
							hdr.Grow(1).MinWidth(0)
						}
						hdr.Children(func() {
							ui.Text(c, col.Title).FontSize(12).FontWeight(600).TextColor(t.MutedForeground).SingleLine()
						})
					}
				})

			// Rows
			for r := 0; r < rowCount; r++ {
				rowIdx := r
				rowEl := ui.Row(c).
					WidthPercent(100).
					AlignItems(ui.Center).
					Padding(10, 14).
					Gap(12)
				if rowIdx < rowCount-1 {
					rowEl.BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
				}
				if rowEl.Hovered() {
					rowEl.Background(theme.Alpha(t.Accent, 0.6))
				}
				rowEl.Children(func() {
					for colIdx, col := range columns {
						cIdx := colIdx
						cellEl := ui.Row(c).AlignItems(ui.Center).Justify(col.Align)
						if col.Width > 0 {
							cellEl.Width(col.Width).Shrink(0)
						} else {
							cellEl.Grow(1).MinWidth(0)
						}
						cellEl.Children(func() {
							if cell != nil {
								cell(rowIdx, cIdx)
							}
						})
					}
				})
			}
		})
}
