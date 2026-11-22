package particles

import (
	_ "embed"
	"fmt"
	"sort"
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

type invoiceRow struct {
	id, status, method, amount string
}

var sampleInvoices = []invoiceRow{
	{"INV001", "Paid", "Credit Card", "$250.00"},
	{"INV002", "Pending", "PayPal", "$150.00"},
	{"INV003", "Unpaid", "Bank Transfer", "$350.00"},
	{"INV004", "Paid", "Credit Card", "$450.00"},
}

func renderInvoiceTable(c *ui.Context) {
	t := theme.Of(c)
	cols := []coss.TableColumnSpec{
		{Title: "Invoice", Width: 80},
		{Title: "Status", Width: 84},
		{Title: "Method"},
		{Title: "Amount", Width: 76, Align: ui.End},
	}
	coss.Table(c, cols, len(sampleInvoices), func(row, col int) {
		r := sampleInvoices[row]
		switch col {
		case 0:
			ui.Text(c, r.id).FontSize(13).FontWeight(500).TextColor(t.Foreground)
		case 1:
			v := coss.BadgeSuccess
			if r.status == "Pending" {
				v = coss.BadgeWarning
			} else if r.status == "Unpaid" {
				v = coss.BadgeError
			}
			coss.Badge(c, r.status).Variant(v).Size(coss.BadgeSizeSM)
		case 2:
			ui.Text(c, r.method).FontSize(13).TextColor(t.MutedForeground)
		case 3:
			ui.Text(c, r.amount).FontSize(13).FontWeight(500).TextColor(t.Foreground)
		}
	})
}

func PTable1(c *ui.Context) { renderInvoiceTable(c) }

func PTable2(c *ui.Context) {
	coss.Frame(c, func() {
		coss.FrameHeader(c, func() {
			coss.FrameTitle(c, "Recent Invoices")
			coss.FrameDescription(c, "A list of your recent billing transactions.")
		})
		renderInvoiceTable(c)
	})
}

func PTable3(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).Gap(8)
	sel := ui.Local(host, "sel", func() [4]bool { return [4]bool{true, false, false, false} })
	host.Children(func() {
		cols := []coss.TableColumnSpec{
			{Title: "", Width: 28},
			{Title: "Invoice", Width: 80},
			{Title: "Status"},
			{Title: "Amount", Width: 76, Align: ui.End},
		}
		coss.Table(c, cols, len(sampleInvoices), func(row, col int) {
			r := sampleInvoices[row]
			switch col {
			case 0:
				coss.Checkbox(c, &sel[row])
			case 1:
				ui.Text(c, r.id).FontSize(13).FontWeight(500).TextColor(t.Foreground)
			case 2:
				ui.Text(c, r.status).FontSize(13).TextColor(t.MutedForeground)
			case 3:
				ui.Text(c, r.amount).FontSize(13).TextColor(t.Foreground)
			}
		})
		n := 0
		for _, b := range sel {
			if b {
				n++
			}
		}
		ui.Text(c, fmt.Sprintf("%d of %d rows selected", n, len(sampleInvoices))).FontSize(12).TextColor(t.MutedForeground)
	})
}

func PTable4(c *ui.Context) {
	type flight struct{ code, route, status string }
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).Gap(10)
	asc := ui.Local(host, "asc", func() bool { return true })
	page := ui.Local(host, "page", func() int { return 1 })

	rows := []flight{
		{"UA 402", "SFO → JFK", "On Time"},
		{"DL 118", "LAX → ORD", "Boarding"},
		{"BA 284", "SFO → LHR", "Delayed"},
		{"NH 007", "SFO → NRT", "On Time"},
	}
	sort.Slice(rows, func(i, j int) bool {
		if *asc {
			return rows[i].code < rows[j].code
		}
		return rows[i].code > rows[j].code
	})

	host.Children(func() {
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
				ui.Icon(c, icons.PlaneTakeoff()).Size(15, 15).TextColor(t.MutedForeground)
				ui.Text(c, "Departures").FontSize(14).FontWeight(600).TextColor(t.Foreground)
			})
			ic := icons.ChevronUp()
			if !*asc {
				ic = icons.ChevronDown()
			}
			coss.Button(c, "Sort Flight").Size(coss.SizeXS).Variant(coss.ButtonOutline).TrailingIcon(ic).OnClick(func() {
				*asc = !*asc
			})
		})
		cols := []coss.TableColumnSpec{
			{Title: "Flight", Width: 84},
			{Title: "Route"},
			{Title: "Status", Width: 84, Align: ui.End},
		}
		start := (*page - 1) * 2
		end := start + 2
		if end > len(rows) {
			end = len(rows)
		}
		sub := rows[start:end]
		coss.Table(c, cols, len(sub), func(row, col int) {
			f := sub[row]
			switch col {
			case 0:
				ui.Text(c, f.code).FontSize(13).FontWeight(500).TextColor(t.Foreground)
			case 1:
				ui.Text(c, f.route).FontSize(13).TextColor(t.MutedForeground)
			case 2:
				coss.Badge(c, f.status).Variant(coss.BadgeOutline).Size(coss.BadgeSizeSM)
			}
		})
		coss.Pagination(c, page, 2)
	})
}

func PTable5(c *ui.Context) { coss.Card(c, func() { renderInvoiceTable(c) }) }
func PTable6(c *ui.Context) { coss.Frame(c, func() { PTable3(c) }) }
func PTable7(c *ui.Context) { coss.Card(c, func() { coss.CardPanel(c, func() { renderInvoiceTable(c) }) }) }
func PTable8(c *ui.Context) { coss.Frame(c, func() { coss.FramePanel(c, func() { PTable4(c) }) }) }

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
