package particles

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed date_picker.go
var datePickerSource string

func init() {
	Register(Entry{ID: "p-date-picker-1", Family: "date-picker", Description: "Basic date picker", View: PDatePicker1, Source: extractFunc(datePickerSource, "PDatePicker1")})
	Register(Entry{ID: "p-date-picker-2", Family: "date-picker", Description: "Date range picker", View: PDatePicker2, Source: extractFunc(datePickerSource, "PDatePicker2")})
	Register(Entry{ID: "p-date-picker-3", Family: "date-picker", Description: "Date picker with field label", View: PDatePicker3, Source: extractFunc(datePickerSource, "PDatePicker3")})
	Register(Entry{ID: "p-date-picker-4", Family: "date-picker", Description: "Date picker with quick presets", View: PDatePicker4, Source: extractFunc(datePickerSource, "PDatePicker4")})
	Register(Entry{ID: "p-date-picker-5", Family: "date-picker", Description: "Date picker with manual text input", View: PDatePicker5, Source: extractFunc(datePickerSource, "PDatePicker5")})
	Register(Entry{ID: "p-date-picker-6", Family: "date-picker", Description: "Date picker closing on select", View: PDatePicker6, Source: extractFunc(datePickerSource, "PDatePicker6")})
	Register(Entry{ID: "p-date-picker-7", Family: "date-picker", Description: "Multiple dates picker", View: PDatePicker7, Source: extractFunc(datePickerSource, "PDatePicker7")})
	Register(Entry{ID: "p-date-picker-8", Family: "date-picker", Description: "Date picker with select-style trigger", View: PDatePicker8, Source: extractFunc(datePickerSource, "PDatePicker8")})
	Register(Entry{ID: "p-date-picker-9", Family: "date-picker", Description: "Two-month range date picker", View: PDatePicker9, Source: extractFunc(datePickerSource, "PDatePicker9")})
}

func PDatePicker1(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() {
		btn := coss.Button(c, d.Format("Jan 02, 2006")).
			Variant(coss.ButtonOutline).
			LeadingIcon(icons.Calendar()).
			OnClick(func() { *open = !*open })
		coss.Popover(c, btn.Element, open, func() { coss.Calendar(c, d) })
	})
}

func PDatePicker2(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() {
		label := fmt.Sprintf("%s – %s", d.Format("Nov 02"), d.AddDate(0, 0, 5).Format("Nov 02, 2006"))
		btn := coss.Button(c, label).
			Variant(coss.ButtonOutline).
			LeadingIcon(icons.Calendar()).
			OnClick(func() { *open = !*open })
		coss.Popover(c, btn.Element, open, func() { coss.Calendar(c, d) })
	})
}

func PDatePicker3(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	open := ui.Local(host, "open", func() bool { return false })
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Date of birth")
			btn := coss.Button(c, d.Format("Jan 02, 2006")).
				Variant(coss.ButtonOutline).
				LeadingIcon(icons.Calendar()).
				OnClick(func() { *open = !*open })
			coss.Popover(c, btn.Element, open, func() { coss.Calendar(c, d) })
		})
	})
}

func PDatePicker4(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() {
		btn := coss.Button(c, d.Format("Jan 02, 2006")).
			Variant(coss.ButtonOutline).
			LeadingIcon(icons.Calendar()).
			OnClick(func() { *open = !*open })
		coss.Popover(c, btn.Element, open, func() {
			ui.Column(c).Gap(8).Children(func() {
				ui.Row(c).Gap(6).Children(func() {
					coss.Button(c, "Today").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *d = calRefDate })
					coss.Button(c, "Tomorrow").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *d = calRefDate.AddDate(0, 0, 1) })
					coss.Button(c, "Next week").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *d = calRefDate.AddDate(0, 0, 7) })
				})
				coss.Calendar(c, d)
			})
		})
	})
}

func PDatePicker5(c *ui.Context) {
	host := ui.Box(c).Width(220)
	open := ui.Local(host, "open", func() bool { return false })
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	txt := ui.Local(host, "txt", func() string { return calRefDate.Format("2006-01-02") })
	host.Children(func() {
		grp := coss.InputGroup(c, func() {
			coss.InputGroupInput(c, txt).Placeholder("YYYY-MM-DD")
			coss.IconButton(c, icons.Calendar(), "Open calendar").Size(coss.SizeIconXS).Variant(coss.ButtonGhost).OnClick(func() { *open = !*open })
		})
		coss.Popover(c, grp, open, func() { coss.Calendar(c, d) })
	})
}

func PDatePicker6(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	prev := *d
	host.Children(func() {
		btn := coss.Button(c, d.Format("Jan 02, 2006")).
			Variant(coss.ButtonOutline).
			LeadingIcon(icons.Calendar()).
			OnClick(func() { *open = !*open })
		coss.Popover(c, btn.Element, open, func() {
			coss.Calendar(c, d)
			if !d.Equal(prev) {
				*open = false
			}
		})
	})
}

func PDatePicker7(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() {
		btn := coss.Button(c, "2 dates selected").
			Variant(coss.ButtonOutline).
			LeadingIcon(icons.Calendar()).
			OnClick(func() { *open = !*open })
		coss.Popover(c, btn.Element, open, func() { coss.Calendar(c, d) })
	})
}

func PDatePicker8(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).Width(200)
	open := ui.Local(host, "open", func() bool { return false })
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() {
		btn := coss.ButtonBox(c, func() {
			ui.Text(c, d.Format("Jan 02, 2006")).FontSize(14).Grow(1)
			ui.Icon(c, icons.ChevronsUpDown()).Size(14, 14).TextColor(t.MutedForeground)
		}).Variant(coss.ButtonOutline).OnClick(func() { *open = !*open })
		btn.WidthPercent(100)
		coss.Popover(c, btn.Element, open, func() { coss.Calendar(c, d) })
	})
}

func PDatePicker9(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	d1 := ui.Local(host, "d1", func() time.Time { return calRefDate })
	d2 := ui.Local(host, "d2", func() time.Time { return calRefDate.AddDate(0, 1, 0) })
	host.Children(func() {
		btn := coss.Button(c, fmt.Sprintf("%s – %s", d1.Format("Nov 02"), d2.Format("Dec 02"))).
			Variant(coss.ButtonOutline).
			LeadingIcon(icons.Calendar()).
			OnClick(func() { *open = !*open })
		coss.Popover(c, btn.Element, open, func() {
			ui.Row(c).Gap(10).Children(func() {
				coss.Calendar(c, d1)
				coss.Calendar(c, d2)
			})
		})
	})
}
