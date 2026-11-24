package particles

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

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
	Register(Entry{ID: "p-autocomplete-6", Family: "autocomplete", Description: "Autocomplete autofilling the input with the highlighted item", View: PAutocomplete6, Source: extractFunc(autocompleteSource, "PAutocomplete6")})
	Register(Entry{ID: "p-autocomplete-7", Family: "autocomplete", Description: "Autocomplete auto highlighting the first option", View: PAutocomplete7, Source: extractFunc(autocompleteSource, "PAutocomplete7")})
	Register(Entry{ID: "p-autocomplete-8", Family: "autocomplete", Description: "Autocomplete with clear button", View: PAutocomplete8, Source: extractFunc(autocompleteSource, "PAutocomplete8")})
	Register(Entry{ID: "p-autocomplete-9", Family: "autocomplete", Description: "Autocomplete with trigger and clear buttons", View: PAutocomplete9, Source: extractFunc(autocompleteSource, "PAutocomplete9")})
	Register(Entry{ID: "p-autocomplete-10", Family: "autocomplete", Description: "Autocomplete with grouped items", View: PAutocomplete10, Source: extractFunc(autocompleteSource, "PAutocomplete10")})
	Register(Entry{ID: "p-autocomplete-11", Family: "autocomplete", Description: "Autocomplete with limited number of results", View: PAutocomplete11, Source: extractFunc(autocompleteSource, "PAutocomplete11")})
	Register(Entry{ID: "p-autocomplete-12", Family: "autocomplete", Description: "Autocomplete with async items loading", View: PAutocomplete12, Source: extractFunc(autocompleteSource, "PAutocomplete12")})
	Register(Entry{ID: "p-autocomplete-13", Family: "autocomplete", Description: "Autocomplete form", View: PAutocomplete13, Source: extractFunc(autocompleteSource, "PAutocomplete13")})
	Register(Entry{ID: "p-autocomplete-14", Family: "autocomplete", Description: "Autocomplete with search icon", View: PAutocomplete14, Source: extractFunc(autocompleteSource, "PAutocomplete14")})
	Register(Entry{ID: "p-autocomplete-15", Family: "autocomplete", Description: "Pill-shaped autocomplete", View: PAutocomplete15, Source: extractFunc(autocompleteSource, "PAutocomplete15")})
	Register(Entry{ID: "p-autocomplete-16", Family: "autocomplete", Description: "Address autocomplete with sample addresses", View: PAutocomplete16, Source: extractFunc(autocompleteSource, "PAutocomplete16")})
}

// acBasic is the fruit autocomplete most particles show, max-w-64 wide.
func acBasic(c *ui.Context, initial string, p coss.AutocompleteProps) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	q := ui.Local(host, "q", func() string { return initial })
	host.Children(func() { coss.AutocompleteWith(c, q, cbFruits, p, "No items found.") })
}

func PAutocomplete1(c *ui.Context) {
	acBasic(c, "", coss.AutocompleteProps{Placeholder: "Search items…"})
}

func PAutocomplete2(c *ui.Context) {
	acBasic(c, "", coss.AutocompleteProps{Placeholder: "Search items…", Disabled: true})
}

func PAutocomplete3(c *ui.Context) {
	acBasic(c, "", coss.AutocompleteProps{Placeholder: "Search items…", Size: coss.InputSM})
}

func PAutocomplete4(c *ui.Context) {
	acBasic(c, "", coss.AutocompleteProps{Placeholder: "Search items…", Size: coss.InputLG})
}

func PAutocomplete5(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).MaxWidth(256).Gap(8)
	q := ui.Local(host, "q", func() string { return "" })
	host.Children(func() {
		coss.Label(c, "Fruits")
		coss.AutocompleteWith(c, q, cbFruits, coss.AutocompleteProps{Placeholder: "Search items…"}, "No items found.")
	})
}

// PAutocomplete6 is mode="both": the popup filters and the text autofills
// with the highlighted item.
func PAutocomplete6(c *ui.Context) {
	acBasic(c, "", coss.AutocompleteProps{Placeholder: "Search items…"})
}

// PAutocomplete7 is autoHighlight: the first match is highlighted.
func PAutocomplete7(c *ui.Context) {
	acBasic(c, "", coss.AutocompleteProps{Placeholder: "Search items…"})
}

func PAutocomplete8(c *ui.Context) {
	acBasic(c, "", coss.AutocompleteProps{Placeholder: "Search items…", ShowClear: true})
}

func PAutocomplete9(c *ui.Context) {
	acBasic(c, "", coss.AutocompleteProps{Placeholder: "Search items…", ShowClear: true, ShowTrigger: true})
}

var acTagGroups = []struct {
	name string
	tags []string
}{
	{"Status", []string{"Open", "In progress", "Blocked", "Resolved", "Closed"}},
	{"Priority", []string{"Low", "Medium", "High", "Urgent"}},
	{"Team", []string{"Design", "Frontend", "Backend", "DevOps", "QA", "Mobile", "Data", "Security", "Platform", "Infrastructure", "Product", "Marketing", "Sales", "Support", "Research", "Content", "Analytics", "Operations", "Finance", "HR", "Legal", "Growth", "Partner", "Community", "Docs", "Localization", "Accessibility", "SRE", "Release", "Architecture", "UX", "UI", "Management"}},
}

func PAutocomplete10(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	q := ui.Local(host, "q", func() string { return "" })
	host.Children(func() {
		ac := coss.AutocompleteInput(c, q, coss.AutocompleteProps{Placeholder: "e.g. feature"})
		coss.ComboboxPopup(c, ac.Parts, func() {
			n := 0
			for i, g := range acTagGroups {
				tags := coss.ComboboxFilter(ac.Parts, *q, g.tags)
				if len(tags) == 0 {
					continue
				}
				n += len(tags)
				coss.SelectGroupLabel(c, g.name)
				for _, tag := range tags {
					coss.AutocompleteItem(c, ac.Parts, tag, func() { ui.Text(c, tag).SingleLine() })
				}
				if i < len(acTagGroups)-1 {
					coss.SelectSeparator(c)
				}
			}
			if n == 0 {
				coss.ComboboxEmpty(c, "No tags found.")
			}
		})
		if v, ok := ac.Parts.Chosen(); ok {
			*q = v
		}
	})
}

var acLanguages = []string{"JavaScript", "TypeScript", "Python", "Java", "C#", "C++", "C", "Go", "Rust", "Ruby", "PHP", "Swift", "Kotlin", "Scala", "Elixir", "Haskell", "Dart", "Objective-C", "Julia", "R", "Perl", "Lua", "OCaml", "F#"}

func PAutocomplete11(c *ui.Context) {
	const limit = 7
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	q := ui.Local(host, "q", func() string { return "" })
	host.Children(func() {
		ac := coss.AutocompleteInput(c, q, coss.AutocompleteProps{Placeholder: "e.g. feature"})
		matches := coss.ComboboxFilter(ac.Parts, *q, acLanguages)
		coss.ComboboxPopup(c, ac.Parts, func() {
			if len(matches) == 0 {
				coss.ComboboxEmpty(c, "No tags found.")
			}
			for i, lang := range matches {
				if i == limit {
					break
				}
				coss.AutocompleteItem(c, ac.Parts, lang, func() { ui.Text(c, lang).SingleLine() })
			}
			if more := len(matches) - limit; more > 0 {
				coss.AutocompleteStatus(c, func() {
					ui.Text(c, fmt.Sprintf("+%d more (keep typing to narrow down)", more))
				})
			}
		})
		if v, ok := ac.Parts.Chosen(); ok {
			*q = v
		}
	})
}

// acTwoLine is a suggestion with a medium title above an xs muted line.
func acTwoLine(c *ui.Context, title, sub string) {
	ui.Column(c).Grow(1).MinWidth(0).Children(func() {
		ui.Text(c, title).FontWeight(500).SingleLine().Ellipsis("…").LineHeight(20.0 / 14)
		ui.Text(c, sub).FontSize(12).TextColor(theme.Of(c).MutedForeground).SingleLine().Ellipsis("…").LineHeight(16.0 / 12)
	})
}

// acSearch shows the suggestions in items matching the text, with a status
// line naming how many matched; the popup only opens on text.
func acSearch(c *ui.Context, p coss.AutocompleteProps, items [][2]string, noun, notFound string) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	q := ui.Local(host, "q", func() string { return "" })
	host.Children(func() {
		ac := coss.AutocompleteInput(c, q, p)
		query := strings.ToLower(strings.TrimSpace(*q))
		var matches [][2]string
		for _, it := range items {
			if query != "" && (strings.Contains(strings.ToLower(it[0]), query) || strings.Contains(strings.ToLower(it[1]), query)) {
				matches = append(matches, it)
			}
		}
		if query == "" {
			return
		}
		coss.ComboboxPopup(c, ac.Parts, func() {
			coss.AutocompleteStatus(c, func() {
				if len(matches) == 0 {
					ui.Text(c, fmt.Sprintf(notFound, *q)).FontSize(14).FontWeight(400)
					return
				}
				s := "s"
				if len(matches) == 1 {
					s = ""
				}
				ui.Text(c, fmt.Sprintf("%d %s%s found", len(matches), noun, s))
			})
			for _, it := range matches {
				coss.AutocompleteItem(c, ac.Parts, it[0], func() { acTwoLine(c, it[0], it[1]) })
			}
		})
		if v, ok := ac.Parts.Chosen(); ok {
			*q = v
		}
	})
}

var acMovies = [][2]string{
	{"The Shawshank Redemption", "1994"},
	{"The Godfather", "1972"},
	{"The Dark Knight", "2008"},
	{"The Godfather Part II", "1974"},
	{"12 Angry Men", "1957"},
	{"Pulp Fiction", "1994"},
	{"Forrest Gump", "1994"},
	{"Inception", "2010"},
}

func PAutocomplete12(c *ui.Context) {
	acSearch(c, coss.AutocompleteProps{Placeholder: "e.g. Pulp Fiction or 1994"}, acMovies,
		"result", `Movie or year "%s" does not exist in the Top 100 IMDb movies`)
}

func PAutocomplete13(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	q := ui.Local(host, "q", func() string { return "" })
	invalid := ui.Local(host, "invalid", func() bool { return false })
	until := ui.Local(host, "loadingUntil", func() time.Time { return time.Time{} })
	loading := c.Now().Before(*until)
	if loading {
		c.After(150 * time.Millisecond)
	}
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() {
				coss.FieldLabel(c, "Favorite item")
				ac := coss.AutocompleteWith(c, q, cbFruits, coss.AutocompleteProps{Placeholder: "Search items…"}, "No items found.")
				if *invalid && *q == "" {
					ac.Border(1, theme.Alpha(t.Destructive, 0.36))
					coss.FieldError(c, "Please select a item.")
				}
			})
			coss.Button(c, "Submit").Loading(loading).OnClick(func() {
				*invalid = *q == ""
				if !*invalid {
					*until = c.Now().Add(800 * time.Millisecond)
				}
			})
		})
	})
}

func PAutocomplete14(c *ui.Context) {
	acBasic(c, "", coss.AutocompleteProps{Placeholder: "Search items…", StartIcon: icons.Search()})
}

func PAutocomplete15(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	q := ui.Local(host, "q", func() string { return "" })
	host.Children(func() {
		coss.AutocompleteWith(c, q, cbFruits, coss.AutocompleteProps{Placeholder: "Search items..."}, "No items found.").Radius(9999)
	})
}

var acAddresses = [][2]string{
	{"1600 Amphitheatre Parkway", "Mountain View, CA 94043, USA"},
	{"1600 Pennsylvania Avenue NW", "Washington, DC 20500, USA"},
	{"350 Fifth Avenue", "New York, NY 10118, USA"},
	{"221B Baker Street", "London NW1 6XE, United Kingdom"},
	{"Champ de Mars, 5 Avenue Anatole France", "75007 Paris, France"},
	{"Piazza del Colosseo, 1", "00184 Roma RM, Italy"},
	{"Platz der Republik 1", "11011 Berlin, Germany"},
	{"1 Macquarie Street", "Sydney NSW 2000, Australia"},
}

func PAutocomplete16(c *ui.Context) {
	acSearch(c, coss.AutocompleteProps{Placeholder: "Enter an address", StartIcon: icons.MapPin()}, acAddresses,
		"suggestion", `No addresses found for "%s"`)
}
