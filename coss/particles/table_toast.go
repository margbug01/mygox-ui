package particles

import (
	_ "embed"
	"fmt"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed table_toast.go
var tableToastSource string

func init() {
	Register(Entry{ID: "p-table-1", Family: "table", Description: "Basic invoice table", View: PTable1, Source: extractFunc(tableToastSource, "PTable1")})
	Register(Entry{ID: "p-table-2", Family: "table", Description: "Frame wrapping card-style table", View: PTable2, Source: extractFunc(tableToastSource, "PTable2")})
	Register(Entry{ID: "p-table-3", Family: "table", Description: "Selectable table with row checkboxes", View: PTable3, Source: extractFunc(tableToastSource, "PTable3")})
	Register(Entry{ID: "p-table-4", Family: "table", Description: "Sortable & paginated flights table", View: PTable4, Source: extractFunc(tableToastSource, "PTable4")})
	Register(Entry{ID: "p-table-5", Family: "table", Description: "Card-style elevated table", View: PTable5, Source: extractFunc(tableToastSource, "PTable5")})
	Register(Entry{ID: "p-table-6", Family: "table", Description: "Card-framed selectable table", View: PTable6, Source: extractFunc(tableToastSource, "PTable6")})
	Register(Entry{ID: "p-table-7", Family: "table", Description: "Card-framed invoice summary table", View: PTable7, Source: extractFunc(tableToastSource, "PTable7")})
	Register(Entry{ID: "p-table-8", Family: "table", Description: "Card-framed sortable flights table", View: PTable8, Source: extractFunc(tableToastSource, "PTable8")})

	Register(Entry{ID: "p-toast-1", Family: "toast", Description: "Default stacked toast", View: PToast1, Source: extractFunc(tableToastSource, "PToast1")})
	Register(Entry{ID: "p-toast-2", Family: "toast", Description: "Semantic toasts (success, error, info, warning)", View: PToast2, Source: extractFunc(tableToastSource, "PToast2")})
	Register(Entry{ID: "p-toast-3", Family: "toast", Description: "Loading-state toast", View: PToast3, Source: extractFunc(tableToastSource, "PToast3")})
	Register(Entry{ID: "p-toast-4", Family: "toast", Description: "Toast with Undo action button", View: PToast4, Source: extractFunc(tableToastSource, "PToast4")})
	Register(Entry{ID: "p-toast-5", Family: "toast", Description: "Async promise toast flow", View: PToast5, Source: extractFunc(tableToastSource, "PToast5")})
	Register(Entry{ID: "p-toast-6", Family: "toast", Description: "Variable height stacked toasts", View: PToast6, Source: extractFunc(tableToastSource, "PToast6")})
	Register(Entry{ID: "p-toast-7", Family: "toast", Description: "Anchored toast after copy", View: PToast7, Source: extractFunc(tableToastSource, "PToast7")})
	Register(Entry{ID: "p-toast-8", Family: "toast", Description: "Anchored error toast on failure", View: PToast8, Source: extractFunc(tableToastSource, "PToast8")})
	Register(Entry{ID: "p-toast-9", Family: "toast", Description: "Long-running download toast with cancel", View: PToast9, Source: extractFunc(tableToastSource, "PToast9")})
	Register(Entry{ID: "p-toast-10", Family: "toast", Description: "Deduplicated success toast", View: PToast10, Source: extractFunc(tableToastSource, "PToast10")})
	Register(Entry{ID: "p-toast-11", Family: "toast", Description: "Deduplicated error toast", View: PToast11, Source: extractFunc(tableToastSource, "PToast11")})
	Register(Entry{ID: "p-toast-12", Family: "toast", Description: "Anchored save success toast", View: PToast12, Source: extractFunc(tableToastSource, "PToast12")})
	Register(Entry{ID: "p-toast-13", Family: "toast", Description: "Anchored save error toast", View: PToast13, Source: extractFunc(tableToastSource, "PToast13")})
}

type project struct {
	name, status, team, budget string
}

var sampleProjects = []project{
	{"Website Redesign", "Paid", "Frontend Team", "$12,500"},
	{"Mobile App", "Unpaid", "Mobile Team", "$8,750"},
	{"API Integration", "Pending", "Backend Team", "$5,200"},
	{"Database Migration", "Paid", "DevOps Team", "$3,800"},
	{"User Dashboard", "Paid", "UX Team", "$7,200"},
	{"Security Audit", "Failed", "Security Team", "$2,100"},
}

// statusBadge is the outline badge with a colored dot that the project
// and flight tables show a status in.
func statusBadge(c *ui.Context, status string) {
	t := theme.Of(c)
	dot := theme.Alpha(t.MutedForeground, 0.64)
	switch status {
	case "Paid", "On Time":
		dot = ui.Hex("#10b981")
	case "Pending", "Delayed":
		dot = ui.Hex("#f59e0b")
	case "Failed", "Cancelled":
		dot = ui.Hex("#ef4444")
	case "Boarding":
		dot = ui.Hex("#3b82f6")
	}
	coss.Badge(c, status).Variant(coss.BadgeOutline).Dot(dot)
}

// projectCell draws the project tables' cells; col counts from the
// Project column.
func projectCell(c *ui.Context, p project, col int) {
	t := theme.Of(c)
	switch col {
	case 0:
		ui.Text(c, p.name).SingleLine().FixedLineHeight(14).FontWeight(500).TextColor(t.Foreground)
	case 1:
		statusBadge(c, p.status)
	case 2:
		ui.Text(c, p.team).SingleLine().FixedLineHeight(14).TextColor(t.Foreground)
	case 3:
		ui.Text(c, p.budget).SingleLine().FixedLineHeight(14).TextColor(t.Foreground)
	}
}

// projectColumns are the project tables' columns, widths in proportion to
// their content as an HTML table's.
var projectColumns = []coss.TableColumnSpec{
	{Title: "Project", Width: 153},
	{Title: "Status", Width: 86},
	{Title: "Team", Width: 118},
	{Title: "Budget", Width: 74, Align: ui.End},
}

func projectTable(c *ui.Context, opts coss.TableOptions) {
	opts.FooterLabel, opts.FooterValue = "Total Budget", "$39,550"
	coss.Table(c, projectColumns, len(sampleProjects), func(row, col int) {
		projectCell(c, sampleProjects[row], col)
	}, opts)
}

// tableCardFrame is coss's CardFrame: a bordered card whose recessed
// surface a table fills edge to edge.
func tableCardFrame(c *ui.Context, fn func()) {
	t := theme.Of(c)
	ui.Column(c).WidthPercent(100).Radius(t.Radius2XL).Border(1, t.Border).Background(t.Card).Clip().
		Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).Children(func() {
		ui.Column(c).WidthPercent(100).Background(theme.Alpha(t.Muted, 0.72)).Children(func() {
			// The table's body overlaps the frame's border, as coss's -m-px.
			ui.Column(c).AlignSelf(ui.Stretch).Margin(-1).Children(fn)
		})
	})
}

// tablePreview is the width coss.com's card gives the table particles:
// full, up to 896.
func tablePreview(c *ui.Context, fn func()) {
	ui.Column(c).WidthPercent(100).MaxWidth(896).Children(fn)
}

func PTable1(c *ui.Context) {
	tablePreview(c, func() { projectTable(c, coss.TableOptions{Caption: "A list of current projects."}) })
}

func PTable2(c *ui.Context) {
	tablePreview(c, func() {
		coss.Frame(c, func() { projectTable(c, coss.TableOptions{Card: true}) })
	})
}

// selectableProjectTable is p-table-3's table, with a checkbox per row.
func selectableProjectTable(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	sel := ui.Local(host, "sel", func() []bool { return make([]bool, len(sampleProjects)) })
	all := !slices.Contains(*sel, false)
	cols := append([]coss.TableColumnSpec{{Width: 26, Fixed: true, Checkbox: true, Header: func() {
		before := all
		coss.Checkbox(c, &all)
		if all != before {
			for i := range *sel {
				(*sel)[i] = all
			}
		}
	}}}, projectColumns...)
	host.Children(func() {
		coss.Table(c, cols, len(sampleProjects), func(row, col int) {
			if col == 0 {
				coss.Checkbox(c, &(*sel)[row])
				return
			}
			projectCell(c, sampleProjects[row], col-1)
		}, coss.TableOptions{
			Card:        true,
			Selected:    func(row int) bool { return (*sel)[row] },
			FooterLabel: "Total Budget",
			FooterValue: "$39,550",
		})
	})
}

func PTable3(c *ui.Context) {
	tablePreview(c, func() { coss.Frame(c, func() { selectableProjectTable(c) }) })
}

type flight struct {
	code, departure, duration, arrival, destination, status, terminal, gate string
}

var sampleFlights = []flight{
	{"AA1234", "08:30", "5h 15m", "11:45", "Los Angeles", "On Time", "1", "A12"},
	{"DL5678", "14:20", "4h 50m", "17:10", "San Francisco", "Delayed", "2", "B24"},
	{"UA9012", "10:15", "3h 15m", "13:30", "Miami", "On Time", "1", "C8"},
	{"SW3456", "16:45", "2h 35m", "18:20", "Seattle", "On Time", "3", "D15"},
	{"JB7890", "09:00", "5h 30m", "12:30", "Salt Lake City", "Cancelled", "2", "E3"},
	{"AS2345", "11:30", "2h 45m", "14:15", "Phoenix", "On Time", "1", "F7"},
	{"HA6789", "13:00", "5h 30m", "20:30", "Las Vegas", "Delayed", "2", "G12"},
	{"FX0123", "07:15", "1h 45m", "09:00", "Dallas", "Boarding", "1", "H5"},
	{"WN4567", "06:00", "2h 30m", "08:30", "Denver", "Boarding", "2", "I9"},
	{"B61234", "12:45", "2h 35m", "15:20", "Portland", "On Time", "3", "J14"},
	{"NK8901", "15:30", "3h 15m", "18:45", "Atlanta", "On Time", "1", "K6"},
	{"F92345", "09:45", "2h 15m", "12:00", "Chicago", "Delayed", "2", "L11"},
	{"SY6789", "11:00", "3h 15m", "14:15", "Boston", "On Time", "1", "M3"},
	{"G40123", "13:30", "3h 15m", "16:45", "New York", "On Time", "3", "N8"},
	{"YX5678", "08:00", "3h 20m", "11:20", "Washington", "Delayed", "2", "O12"},
	{"4U9012", "10:30", "3h 20m", "13:50", "Orlando", "Delayed", "1", "P5"},
	{"QF3456", "14:00", "2h 30m", "16:30", "Houston", "On Time", "3", "Q9"},
	{"LH7890", "07:30", "2h 30m", "10:00", "Minneapolis", "Cancelled", "2", "R7"},
	{"KL2345", "16:15", "3h 15m", "19:30", "Detroit", "Cancelled", "1", "S4"},
	{"AF6789", "12:00", "3h 10m", "15:10", "Philadelphia", "On Time", "3", "T16"},
	{"BA0123", "09:15", "3h 10m", "12:25", "Charlotte", "On Time", "2", "U10"},
	{"IB4567", "15:45", "2h 15m", "18:00", "Nashville", "Delayed", "1", "V8"},
	{"EK8901", "11:45", "2h 15m", "14:00", "Austin", "Cancelled", "3", "W13"},
	{"QR2345", "13:15", "3h 25m", "16:40", "Tampa", "On Time", "2", "X6"},
	{"TK6789", "08:45", "2h 45m", "11:30", "Raleigh", "On Time", "1", "Y11"},
	{"VS3456", "10:00", "2h 45m", "12:45", "Indianapolis", "On Time", "2", "Z4"},
	{"LX7890", "17:30", "2h 30m", "20:00", "Kansas City", "Delayed", "3", "A8"},
	{"OS1234", "12:30", "2h 50m", "15:20", "Columbus", "On Time", "1", "B19"},
	{"SN5678", "18:00", "2h 15m", "20:15", "Milwaukee", "On Time", "2", "C22"},
	{"TP9012", "19:15", "2h 15m", "21:30", "Memphis", "On Time", "3", "D6"},
}

// flightTimes is a flight's "06:00 ○--- 2h 30m ---○ 08:30" cell.
func flightTimes(c *ui.Context, f flight) {
	t := theme.Of(c)
	cancelled := f.status == "Cancelled"
	// leg is the departure or arrival mark: a ring and a dashed line.
	leg := func(ringFirst bool) {
		ui.Row(c).AlignItems(ui.Center).Gap(2).Opacity(0.5).Children(func() {
			ring := func() { ui.Box(c).Size(6, 6).Radius(999).Border(1, t.MutedForeground) }
			if ringFirst {
				ring()
			}
			ui.Box(c).Width(12).Height(1).BorderWidth(1, 0, 0, 0).BorderColor(t.MutedForeground).BorderStyle(ui.BorderDashed)
			if !ringFirst {
				ring()
			}
		})
	}
	time := func(s string, color ui.Color) {
		tx := ui.Text(c, s).SingleLine().FixedLineHeight(14).FontFeatures("tnum").TextColor(color)
		if cancelled {
			tx.Strikethrough()
		}
	}
	row := ui.Row(c).AlignItems(ui.Center).Gap(6)
	if cancelled {
		row.Opacity(0.5)
	}
	row.Children(func() {
		dep, fg := t.Foreground, t.Foreground
		if cancelled {
			dep, fg = t.MutedForeground, t.MutedForeground
		} else if f.status == "Delayed" {
			dep = t.WarningForeground
		}
		time(f.departure, dep)
		leg(true)
		time(f.duration, t.MutedForeground)
		leg(false)
		time(f.arrival, fg)
	})
}

// flightsTable is p-table-4's sortable, paginated table of departures,
// sorted by departure time to start with, inside frame, with the paging
// controls in footer.
func flightsTable(c *ui.Context, frame func(fn func()), footer func(fn func()) ui.Element) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	type sorting struct {
		col  int
		desc bool
	}
	sortBy := ui.Local(host, "sort", func() sorting { return sorting{col: 2} })
	page := ui.Local(host, "page", func() string { return "1" })
	sel := ui.Local(host, "sel", func() map[string]bool { return map[string]bool{} })

	key := func(f flight) string {
		return [...]string{"", f.code, f.departure, f.destination, f.status, f.terminal, f.gate}[sortBy.col]
	}
	rows := slices.Clone(sampleFlights)
	slices.SortStableFunc(rows, func(a, b flight) int {
		if sortBy.desc {
			return strings.Compare(key(b), key(a))
		}
		return strings.Compare(key(a), key(b))
	})
	const pageSize = 10
	pageIdx, _ := strconv.Atoi(*page)
	pageRows := rows[(pageIdx-1)*pageSize : min(pageIdx*pageSize, len(rows))]

	titles := []string{"", "Flight", "Time", "Destination", "Status", "Terminal", "Gate"}
	widths := []float32{28, 80, 220, 180, 120, 90, 80}
	cols := make([]coss.TableColumnSpec, len(titles))
	for i := range cols {
		cols[i] = coss.TableColumnSpec{Width: widths[i], Header: func() {
			if i == 0 {
				all := !slices.ContainsFunc(pageRows, func(f flight) bool { return !(*sel)[f.code] })
				before := all
				coss.Checkbox(c, &all)
				if all != before {
					for _, f := range pageRows {
						(*sel)[f.code] = all
					}
				}
				return
			}
			h := ui.ButtonBase(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(8).Cursor(ui.CursorPointer)
			h.Children(func() {
				ui.Text(c, titles[i]).SingleLine().FixedLineHeight(14).FontWeight(500).TextColor(t.MutedForeground)
				if sortBy.col == i {
					ic := icons.ChevronUp()
					if sortBy.desc {
						ic = icons.ChevronDown()
					}
					ui.Icon(c, ic).Size(16, 16).Opacity(0.8).TextColor(t.MutedForeground)
				}
			})
			if h.Clicked() {
				if sortBy.col == i {
					sortBy.desc = !sortBy.desc
				} else {
					*sortBy = sorting{col: i}
				}
			}
		}}
	}
	cols[0].Checkbox = true

	host.Children(func() {
		frame(func() {
			coss.Table(c, cols, len(pageRows), func(row, col int) {
				f := pageRows[row]
				switch col {
				case 0:
					v := (*sel)[f.code]
					coss.Checkbox(c, &v)
					(*sel)[f.code] = v
				case 1:
					ui.Text(c, f.code).SingleLine().FixedLineHeight(14).Font("monospace").FontWeight(500).TextColor(t.MutedForeground)
				case 2:
					flightTimes(c, f)
				case 3:
					ui.Text(c, f.destination).SingleLine().FixedLineHeight(14).FontWeight(500).TextColor(t.Foreground)
				case 4:
					statusBadge(c, f.status)
				case 5:
					coss.Badge(c, f.terminal).Variant(coss.BadgeOutline).Size(coss.BadgeSizeLG).Icon(icons.PlaneTakeoff())
				case 6:
					ui.Text(c, f.gate).SingleLine().FixedLineHeight(14).TextColor(t.Foreground)
				}
			}, coss.TableOptions{Card: true, Selected: func(row int) bool { return (*sel)[pageRows[row].code] }})
			footer(func() {
				ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
					ui.Text(c, "Viewing").SingleLine().TextColor(t.MutedForeground)
					var ranges []selectOption
					for i := 0; i*pageSize < len(rows); i++ {
						ranges = append(ranges, selectOption{Value: strconv.Itoa(i + 1), Label: fmt.Sprintf("%d-%d", i*pageSize+1, min((i+1)*pageSize, len(rows)))})
					}
					cossSelect(c, page, ranges, 74)
					ui.Row(c).Gap(4).Children(func() {
						ui.Text(c, "of").SingleLine().TextColor(t.MutedForeground)
						ui.Text(c, strconv.Itoa(len(rows))).SingleLine().FontWeight(500).TextColor(t.Foreground)
						ui.Text(c, "results").SingleLine().TextColor(t.MutedForeground)
					})
				})
				ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
					coss.Button(c, "Previous").Variant(coss.ButtonOutline).Size(coss.SizeSM).Disabled(pageIdx == 1).OnClick(func() {
						*page = strconv.Itoa(pageIdx - 1)
					})
					coss.Button(c, "Next").Variant(coss.ButtonOutline).Size(coss.SizeSM).Disabled(pageIdx*pageSize >= len(rows)).OnClick(func() {
						*page = strconv.Itoa(pageIdx + 1)
					})
				})
			}).Padding(8)
		})
	})
}

func PTable4(c *ui.Context) {
	tablePreview(c, func() {
		flightsTable(c, func(fn func()) { coss.Frame(c, fn) }, func(fn func()) ui.Element { return coss.FrameFooter(c, fn) })
	})
}

func PTable5(c *ui.Context) {
	tablePreview(c, func() { projectTable(c, coss.TableOptions{Card: true}) })
}

func PTable6(c *ui.Context) {
	tablePreview(c, func() { tableCardFrame(c, func() { selectableProjectTable(c) }) })
}

func PTable7(c *ui.Context) {
	tablePreview(c, func() {
		tableCardFrame(c, func() { projectTable(c, coss.TableOptions{Card: true}) })
	})
}

func PTable8(c *ui.Context) {
	tablePreview(c, func() {
		flightsTable(c, func(fn func()) { tableCardFrame(c, fn) }, func(fn func()) ui.Element {
			return ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(fn)
		})
	})
}

// anchoredToastShowing reports whether an anchored toast shown until
// *until still shows, and asks for the frame that hides it.
func anchoredToastShowing(c *ui.Context, until *time.Time) bool {
	left := until.Sub(c.Now())
	if left <= 0 {
		return false
	}
	c.After(left)
	return true
}

func PToast1(c *ui.Context) {
	coss.ToastViewport(c)
	coss.Button(c, "Default Toast").Variant(coss.ButtonOutline).OnClick(func() {
		c.AddToast(ui.Toast{Title: "Event has been created", Description: "Monday, January 3rd at 6:00pm"})
	})
}

func PToast2(c *ui.Context) {
	coss.ToastViewport(c)
	toasts := []ui.Toast{
		{Type: "success", Title: "Success!", Description: "Your changes have been saved."},
		{Type: "error", Title: "Uh oh! Something went wrong.", Description: "There was a problem with your request."},
		{Type: "info", Title: "Heads up!", Description: "You can add components to your app using the cli."},
		{Type: "warning", Title: "Warning!", Description: "Your session is about to expire."},
	}
	ui.Row(c).Wrap().Gap(8).Children(func() {
		for _, ts := range toasts {
			label := strings.ToUpper(ts.Type[:1]) + ts.Type[1:] + " Toast"
			coss.Button(c, label).Variant(coss.ButtonOutline).OnClick(func() { c.AddToast(ts) })
		}
	})
}

func PToast3(c *ui.Context) {
	coss.ToastViewport(c)
	coss.Button(c, "Loading Toast").Variant(coss.ButtonOutline).OnClick(func() {
		c.AddToast(ui.Toast{Type: "loading", Title: "Loading…", Description: "Please wait while we process your request."})
	})
}

func PToast4(c *ui.Context) {
	coss.ToastViewport(c)
	coss.Button(c, "Perform Action").Variant(coss.ButtonOutline).OnClick(func() {
		c.AddToast(ui.Toast{
			Type:        "success",
			Title:       "Action performed",
			Description: "You can undo this action.",
			Timeout:     1000 * time.Second,
			Action:      "Undo",
			OnAction: func() {
				c.AddToast(ui.Toast{Type: "info", Title: "Action undone", Description: "The action has been reverted."})
			},
		})
	})
}

func PToast5(c *ui.Context) {
	coss.ToastViewport(c)
	host := ui.Box(c)
	type promise struct {
		id      string
		settles time.Time
	}
	p := ui.Local(host, "p", func() promise { return promise{} })
	if p.id != "" {
		if c.Now().Before(p.settles) {
			c.After(p.settles.Sub(c.Now()))
		} else if rand.Float64() > 0.3 {
			c.AddToast(ui.Toast{ID: p.id, Type: "success", Title: "This is a success toast!", Description: "Success: Data loaded successfully"})
			p.id = ""
		} else {
			c.AddToast(ui.Toast{ID: p.id, Type: "error", Title: "Something went wrong", Description: "Please try again."})
			p.id = ""
		}
	}
	host.Children(func() {
		coss.Button(c, "Run Promise").Variant(coss.ButtonOutline).OnClick(func() {
			p.id = c.AddToast(ui.Toast{Type: "loading", Title: "Loading…", Description: "The promise is loading.", Timeout: -1})
			p.settles = c.Now().Add(2 * time.Second)
		})
	})
}

func PToast6(c *ui.Context) {
	coss.ToastViewport(c)
	texts := []string{
		"Short message.",
		"A bit longer message that spans two lines.",
		"This is a longer description that intentionally takes more vertical space to demonstrate stacking with varying heights.",
		"An even longer description that should span multiple lines so we can verify the clamped collapsed height and smooth expansion animation when hovering or focusing the viewport.",
	}
	host := ui.Box(c)
	count := ui.Local(host, "count", func() int { return 0 })
	host.Children(func() {
		coss.Button(c, "With Varying Heights").Variant(coss.ButtonOutline).OnClick(func() {
			*count++
			c.AddToast(ui.Toast{Title: fmt.Sprintf("Toast %d created", *count), Description: texts[rand.IntN(len(texts))]})
		})
	})
}

func PToast7(c *ui.Context) {
	host := ui.Box(c)
	until := ui.Local(host, "until", func() time.Time { return time.Time{} })
	copied := anchoredToastShowing(c, until)
	host.Children(func() {
		ic := icons.Copy()
		if copied {
			ic = icons.Check()
		}
		btn := coss.IconButton(c, ic, "Copy link").Variant(coss.ButtonOutline).Disabled(copied).OnClick(func() {
			c.WriteClipboard("https://coss.com")
			*until = c.Now().Add(2 * time.Second)
		})
		cossTooltip(c, btn.Element, "Copy to clipboard")
		if copied {
			coss.AnchoredToast(c, btn.Element, ui.Toast{Title: "Copied!"}, 4, true)
		}
	})
}

func PToast8(c *ui.Context) {
	host := ui.Box(c)
	type form struct{ submitted, failed time.Time }
	f := ui.Local(host, "f", func() form { return form{} })
	now := c.Now()
	submitting := now.Before(f.submitted.Add(2 * time.Second))
	if submitting {
		c.After(f.submitted.Add(2 * time.Second).Sub(now))
	} else if !f.submitted.IsZero() && f.failed.IsZero() {
		f.failed = now
	}
	errorUntil := f.failed.Add(5 * time.Second)
	showing := !f.failed.IsZero() && anchoredToastShowing(c, &errorUntil)
	host.Children(func() {
		label := "Submit"
		if submitting {
			label = "Submitting…"
		}
		btn := coss.Button(c, label).Variant(coss.ButtonOutline).Loading(submitting).Disabled(submitting).OnClick(func() {
			*f = form{submitted: c.Now()}
		})
		if showing {
			coss.AnchoredToast(c, btn.Element, ui.Toast{
				Type:        "error",
				Title:       "Error submitting form",
				Description: "The server is not responding. Please try again later.",
			}, 4, false)
		}
	})
}

func PToast9(c *ui.Context) {
	coss.ToastViewport(c)
	host := ui.Box(c)
	type job struct {
		id       string
		ready    time.Time
		progress int
		tick     time.Time
	}
	j := ui.Local(host, "job", func() job { return job{} })
	if j.id != "" {
		now := c.Now()
		for !now.Before(j.tick) {
			j.progress = min(99, j.progress+rand.IntN(9)+2)
			j.tick = j.tick.Add(300 * time.Millisecond)
		}
		if now.Before(j.ready) {
			c.After(min(j.ready.Sub(now), j.tick.Sub(now)))
		} else {
			if rand.Float64() > 0.2 {
				c.AddToast(ui.Toast{ID: j.id, Type: "success", Title: "Download started", Description: "Your file is now downloading."})
			} else {
				c.AddToast(ui.Toast{ID: j.id, Type: "error", Title: "Failed to generate report", Description: "Please try again later."})
			}
			*j = job{}
		}
	}
	host.Children(func() {
		generating := j.id != ""
		var b *coss.ButtonElement
		if generating {
			b = coss.Button(c, fmt.Sprintf("Loading… %2d%%", j.progress)).Variant(coss.ButtonOutline)
		} else {
			b = coss.Button(c, "Download").Variant(coss.ButtonOutline).LeadingIcon(icons.Download())
		}
		b.Disabled(generating).OnClick(func() {
			id := c.AddToast(ui.Toast{
				Type:        "loading",
				Title:       "Generating report…",
				Description: "Your download will begin once ready.",
				Timeout:     -1,
				Action:      "Cancel",
				OnAction: func() {
					c.AddToast(ui.Toast{Type: "info", Title: "Cancelled", Description: "Report generation was cancelled."})
					*j = job{}
				},
			})
			now := c.Now()
			*j = job{id: id, ready: now.Add(4 * time.Second), tick: now.Add(300 * time.Millisecond)}
		})
	})
}

func PToast10(c *ui.Context) {
	coss.ToastViewport(c)
	coss.Button(c, "One Success Toast").Variant(coss.ButtonOutline).OnClick(func() {
		c.AddToast(ui.Toast{
			ID:          "coss-demo-dedup-toast",
			Type:        "success",
			Title:       "Saved",
			Description: "Repeated clicks update this toast instead of stacking another.",
		})
	})
}

func PToast11(c *ui.Context) {
	coss.ToastViewport(c)
	coss.Button(c, "One Error Toast").Variant(coss.ButtonOutline).OnClick(func() {
		c.AddToast(ui.Toast{
			ID:          "coss-demo-error-toast",
			Type:        "error",
			Title:       "Something went wrong",
			Description: "Repeated clicks update this toast; errors use a shake animation.",
		})
	})
}

// saveDraftToast is p-toast-12 and p-toast-13: a Save button that shows
// toast beside it for two seconds.
func saveDraftToast(c *ui.Context, toast ui.Toast) {
	host := ui.Box(c)
	until := ui.Local(host, "until", func() time.Time { return time.Time{} })
	showing := anchoredToastShowing(c, until)
	host.Children(func() {
		btn := coss.IconButton(c, icons.Save(), "Save").Variant(coss.ButtonOutline).OnClick(func() {
			*until = c.Now().Add(2 * time.Second)
		})
		cossTooltip(c, btn.Element, "Save")
		if showing {
			coss.AnchoredToast(c, btn.Element, toast, 6, true)
		}
	})
}

func PToast12(c *ui.Context) { saveDraftToast(c, ui.Toast{Title: "Draft saved"}) }

func PToast13(c *ui.Context) {
	saveDraftToast(c, ui.Toast{Type: "error", Title: "Couldn't save draft"})
}
