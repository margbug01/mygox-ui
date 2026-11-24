package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// TableColumnSpec defines a column of a coss Table.
type TableColumnSpec struct {
	Title string
	// Header draws the header cell's content in place of Title.
	Header func()
	// Width is the column's least width. Columns share the table's width
	// in proportion to it, as an HTML table shares it by its columns'
	// content, or by their set widths with table-layout: fixed. A table
	// narrower than its columns' widths scrolls sideways.
	Width float32
	// Fixed keeps the column at Width.
	Fixed bool
	Align ui.Align
	// Checkbox drops the cell's padding toward the next column, as coss
	// does for the cells of row checkboxes.
	Checkbox bool
}

// TableOptions are a Table's variant and its optional parts.
type TableOptions struct {
	// Card draws the body as a bordered card, as coss's variant="card".
	Card bool
	// Selected marks rows as selected.
	Selected func(row int) bool
	// FooterLabel and FooterValue fill a footer row: the label across all
	// columns but the last, the value in the last.
	FooterLabel, FooterValue string
	// Caption is the muted line below the table.
	Caption string
}

// Table renders a coss data table: a header row of columns' titles, then
// rowCount rows whose cells cell draws.
func Table(c *ui.Context, columns []TableColumnSpec, rowCount int, cell func(row, col int), opts TableOptions) ui.Element {
	t := theme.Of(c)
	tint := ui.RGB(0, 0, 0)
	if t.Dark {
		tint = ui.RGB(255, 255, 255)
	}
	// Every row is a grid of these tracks, so the rows' columns line up.
	tracks := make([]ui.Track, len(columns))
	var minW float32
	for i, col := range columns {
		minW += col.Width
		tracks[i] = ui.Fr(col.Width).Min(col.Width)
		if col.Fixed {
			tracks[i] = ui.Fixed(col.Width)
		}
	}
	// cells lays out one row's cells; edge is the padding at the row's
	// ends, one less in the card's body to make room for its border.
	cells := func(edge, minH float32, fn func(col int)) {
		for i, col := range columns {
			start, end := float32(10), float32(10)
			if i == 0 {
				start = edge
			}
			if i == len(columns)-1 {
				end = edge
			}
			if col.Checkbox {
				end = 0
			}
			ui.Row(c).MinHeight(minH).MinWidth(0).AlignItems(ui.Center).Justify(col.Align).Padding(0, end, 0, start).
				Children(func() { fn(i) })
		}
	}
	text := func(s string, weight int, color ui.Color) {
		ui.Text(c, s).SingleLine().FontSize(14).FixedLineHeight(14).FontWeight(weight).TextColor(color)
	}
	bodyBg := t.Background
	if opts.Card {
		bodyBg = t.Card
	}
	return ui.Column(c).WidthPercent(100).MinWidth(minW).Children(func() {
		head := ui.Grid(c).ColumnTracks(tracks...)
		if !opts.Card {
			head.BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
		}
		head.Children(func() {
			cells(10, 40, func(i int) {
				if columns[i].Header != nil {
					columns[i].Header()
					return
				}
				text(columns[i].Title, 500, t.MutedForeground)
			})
		})

		body := ui.Column(c)
		edge := float32(10)
		if opts.Card {
			edge = 9
			body.Radius(t.RadiusXL).Border(1, t.Border).Background(t.Card).Clip().
				Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05))
		}
		body.Children(func() {
			for r := 0; r < rowCount; r++ {
				row := ui.Grid(c).ColumnTracks(tracks...).Padding(10, 0)
				if r < rowCount-1 {
					row.BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
				}
				if opts.Selected != nil && opts.Selected(r) {
					row.Background(theme.Mix(bodyBg, tint, 0.04))
				} else if row.Hovered() {
					row.Background(theme.Mix(bodyBg, tint, 0.02))
				}
				row.Children(func() {
					cells(edge, 0, func(i int) { cell(r, i) })
				})
			}
		})

		if opts.FooterLabel != "" {
			foot := ui.Row(c).Padding(14, 10).AlignItems(ui.Center).Justify(ui.SpaceBetween)
			if !opts.Card {
				foot.BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Background(theme.Mix(t.Card, tint, 0.02))
			}
			foot.Children(func() {
				text(opts.FooterLabel, 500, t.Foreground)
				text(opts.FooterValue, 500, t.Foreground)
			})
		}
		if opts.Caption != "" {
			ui.Text(c, opts.Caption).FontSize(14).FixedLineHeight(20).TextColor(t.MutedForeground).
				AlignSelf(ui.Center).Margin(16, 0, 0, 0)
		}
	})
}
