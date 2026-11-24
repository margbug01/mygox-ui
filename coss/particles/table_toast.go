package particles

import (
	_ "embed"
	"fmt"
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

func PToast1(c *ui.Context) {
	host := ui.Column(c).Center().Gap(10)
	shown := ui.Local(host, "shown", func() bool { return true })
	host.Children(func() {
		coss.Button(c, "Show Toast").Variant(coss.ButtonOutline).OnClick(func() { *shown = true })
		if *shown {
			coss.ToastCard(c, coss.ToastSpec{
				Title:       "Event has been created",
				Description: "Monday, November 23, 2026 at 9:00 AM",
				OnClose:     func() { *shown = false },
			})
		}
	})
}

func PToast2(c *ui.Context) {
	host := ui.Column(c).Center().Gap(10)
	kind := ui.Local(host, "kind", func() coss.ToastVariant { return coss.ToastSuccess })
	host.Children(func() {
		ui.Row(c).Gap(6).Wrap().Children(func() {
			coss.Button(c, "Success").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *kind = coss.ToastSuccess })
			coss.Button(c, "Error").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *kind = coss.ToastError })
			coss.Button(c, "Info").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *kind = coss.ToastInfo })
			coss.Button(c, "Warning").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *kind = coss.ToastWarning })
		})
		titles := map[coss.ToastVariant]string{
			coss.ToastSuccess: "Changes saved",
			coss.ToastError:   "Something went wrong",
			coss.ToastInfo:    "New update available",
			coss.ToastWarning: "Approaching rate limit",
		}
		coss.ToastCard(c, coss.ToastSpec{Variant: *kind, Title: titles[*kind], Description: "Semantic toast notification preview."})
	})
}

func PToast3(c *ui.Context) {
	coss.ToastCard(c, coss.ToastSpec{
		Variant:     coss.ToastLoading,
		Title:       "Uploading files…",
		Description: "Please wait while your assets finish syncing.",
	})
}

func PToast4(c *ui.Context) {
	host := ui.Column(c).Center().Gap(8)
	undone := ui.Local(host, "undone", func() bool { return false })
	host.Children(func() {
		title := "Message archived"
		if *undone {
			title = "Archive undone"
		}
		coss.ToastCard(c, coss.ToastSpec{
			Title:       title,
			Description: "Moved to your Archive folder.",
			ActionLabel: "Undo",
			OnAction:    func() { *undone = true },
		})
	})
}

func PToast5(c *ui.Context) {
	type st struct {
		running bool
		done    bool
		until   time.Time
	}
	host := ui.Column(c).Center().Gap(10)
	s := ui.Local(host, "s", func() st { return st{} })
	if s.running {
		if c.Now().After(s.until) {
			s.running = false
			s.done = true
		} else {
			c.After(150 * time.Millisecond)
		}
	}
	host.Children(func() {
		coss.Button(c, "Run Promise").Variant(coss.ButtonOutline).OnClick(func() {
			s.running = true
			s.done = false
			s.until = c.Now().Add(1200 * time.Millisecond)
		})
		if s.running {
			coss.ToastCard(c, coss.ToastSpec{Variant: coss.ToastLoading, Title: "Deploying…", Description: "Building release artifact."})
		} else if s.done {
			coss.ToastCard(c, coss.ToastSpec{Variant: coss.ToastSuccess, Title: "Deployed!", Description: "Production release is live."})
		}
	})
}

func PToast6(c *ui.Context) {
	ui.Column(c).Gap(8).Children(func() {
		coss.ToastCard(c, coss.ToastSpec{Title: "Quick note"})
		coss.ToastCard(c, coss.ToastSpec{Variant: coss.ToastInfo, Title: "Workspace synced", Description: "All 56 primitives and 530 particles are indexed locally for fast browsing."})
	})
}

func PToast7(c *ui.Context) {
	host := ui.Column(c).Center().Gap(8)
	copied := ui.Local(host, "copied", func() bool { return false })
	host.Children(func() {
		ic := icons.Copy()
		if *copied {
			ic = icons.Check()
		}
		coss.Button(c, "Copy snippet").Variant(coss.ButtonOutline).LeadingIcon(ic).OnClick(func() {
			c.WriteClipboard("go get github.com/margbug01/cossgo")
			*copied = true
		})
		if *copied {
			coss.Badge(c, "Copied!").Variant(coss.BadgeSuccess)
		}
	})
}

func PToast8(c *ui.Context) {
	host := ui.Column(c).Center().Gap(8)
	failed := ui.Local(host, "failed", func() bool { return true })
	host.Children(func() {
		coss.Button(c, "Trigger Error").Variant(coss.ButtonDestructiveOutline).OnClick(func() { *failed = true })
		if *failed {
			coss.ToastCard(c, coss.ToastSpec{Variant: coss.ToastError, Title: "Request failed", Description: "Could not reach upstream server.", OnClose: func() { *failed = false }})
		}
	})
}

func PToast9(c *ui.Context) {
	host := ui.Column(c).Center().Gap(8)
	active := ui.Local(host, "active", func() bool { return true })
	host.Children(func() {
		coss.Button(c, "Download report").Variant(coss.ButtonOutline).LeadingIcon(icons.Download()).OnClick(func() { *active = true })
		if *active {
			coss.ToastCard(c, coss.ToastSpec{
				Variant:     coss.ToastLoading,
				Title:       "Generating report…",
				Description: "Your download will start shortly.",
				ActionLabel: "Cancel",
				OnAction:    func() { *active = false },
			})
		}
	})
}

func PToast10(c *ui.Context) {
	host := ui.Column(c).Center().Gap(8)
	count := ui.Local(host, "n", func() int { return 1 })
	host.Children(func() {
		coss.Button(c, "Trigger Success").Variant(coss.ButtonOutline).OnClick(func() { *count++ })
		coss.ToastCard(c, coss.ToastSpec{Variant: coss.ToastSuccess, Title: fmt.Sprintf("Saved (×%d)", *count), Description: "Deduplicated toast updates in place."})
	})
}

func PToast11(c *ui.Context) {
	host := ui.Column(c).Center().Gap(8)
	count := ui.Local(host, "n", func() int { return 1 })
	host.Children(func() {
		coss.Button(c, "Trigger Error").Variant(coss.ButtonOutline).OnClick(func() { *count++ })
		coss.ToastCard(c, coss.ToastSpec{Variant: coss.ToastError, Title: fmt.Sprintf("Sync failed (×%d)", *count), Description: "Deduplicated error toast."})
	})
}

func PToast12(c *ui.Context) {
	host := ui.Column(c).Center().Gap(8)
	saved := ui.Local(host, "saved", func() bool { return true })
	host.Children(func() {
		coss.Button(c, "Save").Variant(coss.ButtonOutline).LeadingIcon(icons.Save()).OnClick(func() { *saved = true })
		if *saved {
			coss.ToastCard(c, coss.ToastSpec{Variant: coss.ToastSuccess, Title: "Changes saved"})
		}
	})
}

func PToast13(c *ui.Context) {
	host := ui.Column(c).Center().Gap(8)
	errShown := ui.Local(host, "err", func() bool { return true })
	host.Children(func() {
		coss.Button(c, "Save").Variant(coss.ButtonOutline).LeadingIcon(icons.Save()).OnClick(func() { *errShown = true })
		if *errShown {
			coss.ToastCard(c, coss.ToastSpec{Variant: coss.ToastError, Title: "Failed to save changes"})
		}
	})
}
