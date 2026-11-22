package particles

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed autocomplete.go
var autocompleteSource string

func init() {
	Register(Entry{ID: "p-autocomplete-1", Family: "autocomplete", Description: "Basic autocomplete", View: PAutocomplete1, Source: extractFunc(autocompleteSource, "PAutocomplete1")})
	Register(Entry{ID: "p-autocomplete-2", Family: "autocomplete", Description: "Disabled autocomplete", View: PAutocomplete2, Source: extractFunc(autocompleteSource, "PAutocomplete2")})
	Register(Entry{ID: "p-autocomplete-3", Family: "autocomplete", Description: "Small autocomplete", View: PAutocomplete3, Source: extractFunc(autocompleteSource, "PAutocomplete3")})
	Register(Entry{ID: "p-autocomplete-4", Family: "autocomplete", Description: "Large autocomplete", View: PAutocomplete4, Source: extractFunc(autocompleteSource, "PAutocomplete4")})
	Register(Entry{ID: "p-autocomplete-5", Family: "autocomplete", Description: "Autocomplete with label", View: PAutocomplete5, Source: extractFunc(autocompleteSource, "PAutocomplete5")})
	Register(Entry{ID: "p-autocomplete-6", Family: "autocomplete", Description: "Autocomplete with inline autofill", View: PAutocomplete6, Source: extractFunc(autocompleteSource, "PAutocomplete6")})
	Register(Entry{ID: "p-autocomplete-7", Family: "autocomplete", Description: "Autocomplete auto-highlighting first match", View: PAutocomplete7, Source: extractFunc(autocompleteSource, "PAutocomplete7")})
	Register(Entry{ID: "p-autocomplete-8", Family: "autocomplete", Description: "Autocomplete with clear button", View: PAutocomplete8, Source: extractFunc(autocompleteSource, "PAutocomplete8")})
	Register(Entry{ID: "p-autocomplete-9", Family: "autocomplete", Description: "Autocomplete with trigger and clear buttons", View: PAutocomplete9, Source: extractFunc(autocompleteSource, "PAutocomplete9")})
	Register(Entry{ID: "p-autocomplete-10", Family: "autocomplete", Description: "Autocomplete with grouped items", View: PAutocomplete10, Source: extractFunc(autocompleteSource, "PAutocomplete10")})
	Register(Entry{ID: "p-autocomplete-11", Family: "autocomplete", Description: "Autocomplete with limited result count", View: PAutocomplete11, Source: extractFunc(autocompleteSource, "PAutocomplete11")})
	Register(Entry{ID: "p-autocomplete-12", Family: "autocomplete", Description: "Autocomplete with async loading state", View: PAutocomplete12, Source: extractFunc(autocompleteSource, "PAutocomplete12")})
	Register(Entry{ID: "p-autocomplete-13", Family: "autocomplete", Description: "Autocomplete inside a form", View: PAutocomplete13, Source: extractFunc(autocompleteSource, "PAutocomplete13")})
	Register(Entry{ID: "p-autocomplete-14", Family: "autocomplete", Description: "Search autocomplete form", View: PAutocomplete14, Source: extractFunc(autocompleteSource, "PAutocomplete14")})
	Register(Entry{ID: "p-autocomplete-15", Family: "autocomplete", Description: "Pill-shaped autocomplete", View: PAutocomplete15, Source: extractFunc(autocompleteSource, "PAutocomplete15")})
	Register(Entry{ID: "p-autocomplete-16", Family: "autocomplete", Description: "Address autocomplete with location pin", View: PAutocomplete16, Source: extractFunc(autocompleteSource, "PAutocomplete16")})
}

var acFruits = []string{"Apple", "Banana", "Blueberry", "Cherry", "Grape", "Mango", "Orange", "Peach", "Pear", "Strawberry"}

func PAutocomplete1(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	q := ui.Local(host, "q", func() string { return "" })
	host.Children(func() { coss.Autocomplete(c, q, acFruits, "Search fruits…") })
}

func PAutocomplete2(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	q := ui.Local(host, "q", func() string { return "Apple" })
	host.Children(func() { coss.Autocomplete(c, q, acFruits).Disabled(true).Opacity(0.64) })
}

func PAutocomplete3(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	q := ui.Local(host, "q", func() string { return "" })
	host.Children(func() { coss.Autocomplete(c, q, acFruits, "Small…").Height(28).FontSize(13) })
}

func PAutocomplete4(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	q := ui.Local(host, "q", func() string { return "" })
	host.Children(func() { coss.Autocomplete(c, q, acFruits, "Large…").Height(38) })
}

func PAutocomplete5(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	q := ui.Local(host, "q", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Favorite fruit")
			coss.Autocomplete(c, q, acFruits, "Type to search…")
		})
	})
}

func PAutocomplete6(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	q := ui.Local(host, "q", func() string { return "Ap" })
	host.Children(func() { coss.Autocomplete(c, q, acFruits, "Autofill…") })
}

func PAutocomplete7(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	q := ui.Local(host, "q", func() string { return "" })
	host.Children(func() { coss.Autocomplete(c, q, acFruits, "Auto-highlight first…") })
}

func PAutocomplete8(c *ui.Context) {
	host := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(6)
	q := ui.Local(host, "q", func() string { return "Mango" })
	host.Children(func() {
		coss.Autocomplete(c, q, acFruits, "Search…").Grow(1)
		if *q != "" {
			coss.IconButton(c, icons.X(), "Clear").Size(coss.SizeIconSM).Variant(coss.ButtonGhost).OnClick(func() { *q = "" })
		}
	})
}

func PAutocomplete9(c *ui.Context) {
	host := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(4)
	q := ui.Local(host, "q", func() string { return "Cherry" })
	host.Children(func() {
		coss.Autocomplete(c, q, acFruits, "Search or pick…").Grow(1)
		if *q != "" {
			coss.IconButton(c, icons.X(), "Clear").Size(coss.SizeIconSM).Variant(coss.ButtonGhost).OnClick(func() { *q = "" })
		}
		coss.IconButton(c, icons.ChevronDown(), "Toggle list").Size(coss.SizeIconSM).Variant(coss.ButtonOutline)
	})
}

func PAutocomplete10(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	q := ui.Local(host, "q", func() string { return "" })
	grouped := []string{"Citrus — Orange", "Citrus — Lemon", "Berries — Blueberry", "Berries — Strawberry", "Stone — Peach"}
	host.Children(func() { coss.Autocomplete(c, q, grouped, "Search grouped fruits…") })
}

func PAutocomplete11(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	q := ui.Local(host, "q", func() string { return "" })
	host.Children(func() { coss.Autocomplete(c, q, acFruits[:5], "Top 5 matches…") })
}

func PAutocomplete12(c *ui.Context) {
	host := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8)
	q := ui.Local(host, "q", func() string { return "" })
	host.Children(func() {
		coss.Autocomplete(c, q, acFruits, "Async search…").Grow(1)
		if strings.TrimSpace(*q) != "" {
			coss.Spinner(c).Size(14)
		}
	})
}

func PAutocomplete13(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	q := ui.Local(host, "q", func() string { return "" })
	saved := ui.Local(host, "saved", func() string { return "" })
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() {
				coss.FieldLabel(c, "Framework")
				coss.Autocomplete(c, q, []string{"MyGo", "Next.js", "Remix", "Astro", "SvelteKit"}, "Pick a framework…")
			})
			coss.Button(c, "Submit").Size(coss.SizeSM).OnClick(func() { *saved = *q })
			if *saved != "" {
				ui.Text(c, fmt.Sprintf("Selected: %s", *saved)).FontSize(12).TextColor(t.SuccessForeground)
			}
		})
	})
}

func PAutocomplete14(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	q := ui.Local(host, "q", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() { ui.Icon(c, icons.Search()).Size(14, 14) })
			coss.InputGroupInput(c, q).Placeholder("Search documentation…")
			coss.Button(c, "Go").Size(coss.SizeXS)
		})
	})
}

func PAutocomplete15(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	q := ui.Local(host, "q", func() string { return "" })
	host.Children(func() { coss.Autocomplete(c, q, acFruits, "Pill search…").Radius(999) })
}

func PAutocomplete16(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	q := ui.Local(host, "q", func() string { return "" })
	addrs := []string{
		"1600 Amphitheatre Pkwy, Mountain View, CA",
		"1 Apple Park Way, Cupertino, CA",
		"350 Fifth Ave, New York, NY",
		"221B Baker St, London",
	}
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Street address")
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(6).Children(func() {
				ui.Icon(c, icons.MapPin()).Size(15, 15).TextColor(theme.Of(c).MutedForeground)
				coss.Autocomplete(c, q, addrs, "Search address…").Grow(1)
			})
		})
	})
}
