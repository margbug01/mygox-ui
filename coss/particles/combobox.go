package particles

import (
	_ "embed"
	"fmt"
	"slices"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed combobox.go
var comboboxSource string

func init() {
	Register(Entry{ID: "p-combobox-1", Family: "combobox", Description: "Basic combobox", View: PCombobox1, Source: extractFunc(comboboxSource, "PCombobox1")})
	Register(Entry{ID: "p-combobox-2", Family: "combobox", Description: "Disabled combobox", View: PCombobox2, Source: extractFunc(comboboxSource, "PCombobox2")})
	Register(Entry{ID: "p-combobox-3", Family: "combobox", Description: "Small combobox", View: PCombobox3, Source: extractFunc(comboboxSource, "PCombobox3")})
	Register(Entry{ID: "p-combobox-4", Family: "combobox", Description: "Large combobox", View: PCombobox4, Source: extractFunc(comboboxSource, "PCombobox4")})
	Register(Entry{ID: "p-combobox-5", Family: "combobox", Description: "Combobox with label", View: PCombobox5, Source: extractFunc(comboboxSource, "PCombobox5")})
	Register(Entry{ID: "p-combobox-6", Family: "combobox", Description: "Combobox auto-highlighting first match", View: PCombobox6, Source: extractFunc(comboboxSource, "PCombobox6")})
	Register(Entry{ID: "p-combobox-7", Family: "combobox", Description: "Combobox with clear button", View: PCombobox7, Source: extractFunc(comboboxSource, "PCombobox7")})
	Register(Entry{ID: "p-combobox-8", Family: "combobox", Description: "Combobox with grouped options", View: PCombobox8, Source: extractFunc(comboboxSource, "PCombobox8")})
	Register(Entry{ID: "p-combobox-9", Family: "combobox", Description: "Multi-select combobox with chips", View: PCombobox9, Source: extractFunc(comboboxSource, "PCombobox9")})
	Register(Entry{ID: "p-combobox-10", Family: "combobox", Description: "Combobox with search input inside popup", View: PCombobox10, Source: extractFunc(comboboxSource, "PCombobox10")})
	Register(Entry{ID: "p-combobox-11", Family: "combobox", Description: "Combobox in a form", View: PCombobox11, Source: extractFunc(comboboxSource, "PCombobox11")})
	Register(Entry{ID: "p-combobox-12", Family: "combobox", Description: "Multi-select combobox in a form", View: PCombobox12, Source: extractFunc(comboboxSource, "PCombobox12")})
	Register(Entry{ID: "p-combobox-13", Family: "combobox", Description: "Combobox with leading search icon", View: PCombobox13, Source: extractFunc(comboboxSource, "PCombobox13")})
	Register(Entry{ID: "p-combobox-14", Family: "combobox", Description: "Multi-select combobox with search addon", View: PCombobox14, Source: extractFunc(comboboxSource, "PCombobox14")})
	Register(Entry{ID: "p-combobox-15", Family: "combobox", Description: "Pill-shaped combobox", View: PCombobox15, Source: extractFunc(comboboxSource, "PCombobox15")})
	Register(Entry{ID: "p-combobox-16", Family: "combobox", Description: "Timezone combobox", View: PCombobox16, Source: extractFunc(comboboxSource, "PCombobox16")})
	Register(Entry{ID: "p-combobox-17", Family: "combobox", Description: "Timezone combobox with search icon", View: PCombobox17, Source: extractFunc(comboboxSource, "PCombobox17")})
	Register(Entry{ID: "p-combobox-18", Family: "combobox", Description: "Combobox with select-style trigger", View: PCombobox18, Source: extractFunc(comboboxSource, "PCombobox18")})
	Register(Entry{ID: "p-combobox-19", Family: "combobox", Description: "Multi-select combobox with stacked chips", View: PCombobox19, Source: extractFunc(comboboxSource, "PCombobox19")})
	Register(Entry{ID: "p-combobox-20", Family: "combobox", Description: "Multi-select combobox with divided chips", View: PCombobox20, Source: extractFunc(comboboxSource, "PCombobox20")})
}

var cbFruits = []string{"Apple", "Banana", "Orange", "Grape", "Strawberry", "Mango", "Pineapple", "Kiwi", "Peach", "Pear"}

// cbCountries are p-combobox-10's countries, first the "Select country"
// entry it starts on.
var cbCountries = []string{
	"Select country",
	"Afghanistan",
	"Albania",
	"Algeria",
	"Andorra",
	"Angola",
	"Argentina",
	"Armenia",
	"Australia",
	"Austria",
	"Azerbaijan",
	"Bahamas",
	"Bahrain",
	"Bangladesh",
	"Barbados",
	"Belarus",
	"Belgium",
	"Belize",
	"Benin",
	"Bhutan",
	"Bolivia",
	"Bosnia and Herzegovina",
	"Botswana",
	"Brazil",
	"Brunei",
	"Bulgaria",
	"Burkina Faso",
	"Burundi",
	"Cambodia",
	"Cameroon",
	"Canada",
	"Cape Verde",
	"Central African Republic",
	"Chad",
	"Chile",
	"China",
	"Colombia",
	"Comoros",
	"Congo",
	"Costa Rica",
	"Croatia",
	"Cuba",
	"Cyprus",
	"Czech Republic",
	"Denmark",
	"Djibouti",
	"Dominica",
	"Dominican Republic",
	"Ecuador",
	"Egypt",
	"El Salvador",
	"Equatorial Guinea",
	"Eritrea",
	"Estonia",
	"Ethiopia",
	"Fiji",
	"Finland",
	"France",
	"Gabon",
	"Gambia",
	"Georgia",
	"Germany",
	"Ghana",
	"Greece",
	"Grenada",
	"Guatemala",
	"Guinea",
	"Guinea-Bissau",
	"Guyana",
	"Haiti",
	"Honduras",
	"Hungary",
	"Iceland",
	"India",
	"Indonesia",
	"Iran",
	"Iraq",
	"Ireland",
	"Israel",
	"Italy",
	"Jamaica",
	"Japan",
	"Jordan",
	"Kazakhstan",
	"Kenya",
	"Kuwait",
	"Kyrgyzstan",
	"Laos",
	"Latvia",
	"Lebanon",
	"Lesotho",
	"Liberia",
	"Libya",
	"Liechtenstein",
	"Lithuania",
	"Luxembourg",
	"Madagascar",
	"Malawi",
	"Malaysia",
	"Maldives",
	"Mali",
	"Malta",
	"Marshall Islands",
	"Mauritania",
	"Mauritius",
	"Mexico",
	"Micronesia",
	"Moldova",
	"Monaco",
	"Mongolia",
	"Montenegro",
	"Morocco",
	"Mozambique",
	"Myanmar",
	"Namibia",
	"Nauru",
	"Nepal",
	"Netherlands",
	"New Zealand",
	"Nicaragua",
	"Niger",
	"Nigeria",
	"North Korea",
	"North Macedonia",
	"Norway",
	"Oman",
	"Pakistan",
	"Palau",
	"Palestine",
	"Panama",
	"Papua New Guinea",
	"Paraguay",
	"Peru",
	"Philippines",
	"Poland",
	"Portugal",
	"Qatar",
	"Romania",
	"Russia",
	"Rwanda",
	"Samoa",
	"San Marino",
	"Saudi Arabia",
	"Senegal",
	"Serbia",
	"Seychelles",
	"Sierra Leone",
	"Singapore",
	"Slovakia",
	"Slovenia",
	"Solomon Islands",
	"Somalia",
	"South Africa",
	"South Korea",
	"South Sudan",
	"Spain",
	"Sri Lanka",
	"Sudan",
	"Suriname",
	"Sweden",
	"Switzerland",
	"Syria",
	"Taiwan",
	"Tajikistan",
	"Tanzania",
	"Thailand",
	"Timor-Leste",
	"Togo",
	"Tonga",
	"Trinidad and Tobago",
	"Tunisia",
	"Turkey",
	"Turkmenistan",
	"Tuvalu",
	"Uganda",
	"Ukraine",
	"United Arab Emirates",
	"United Kingdom",
	"United States",
	"Uruguay",
	"Uzbekistan",
	"Vanuatu",
	"Vatican City",
	"Venezuela",
	"Vietnam",
	"Yemen",
	"Zambia",
	"Zimbabwe",
}

// fruitCombobox is the combobox most particles show: an input over the
// fruits matching its text.
func fruitCombobox(c *ui.Context, v *string, p coss.ComboboxInputProps, empty string) coss.ComboboxInputElement {
	cb := coss.ComboboxInput(c, v, p)
	matches := coss.ComboboxFilter(cb.Parts, *v, cbFruits)
	coss.ComboboxPopup(c, cb.Parts, func() {
		if len(matches) == 0 {
			coss.ComboboxEmpty(c, empty)
		}
		for _, it := range matches {
			coss.ComboboxItem(c, cb.Parts, it, it == *v, func() { ui.Text(c, it).SingleLine() })
		}
	})
	if it, ok := cb.Parts.Chosen(); ok {
		*v = it
	}
	return cb
}

// fruitChips is the multiple combobox of p-combobox-9, 12 and 14: chips
// for the chosen fruits, then an input over the fruits matching it.
func fruitChips(c *ui.Context, chosen *[]string, text *string, placeholder string, startIcon bool) {
	t := theme.Of(c)
	coss.ComboboxChips(c, func() {
		if startIcon {
			ui.Icon(c, icons.Search()).Size(16, 16).Margin(0, 2, 0, 6).TextColor(theme.Alpha(t.Foreground, 0.8))
		}
		for _, it := range *chosen {
			coss.ComboboxChip(c, it, func() {
				*chosen = slices.DeleteFunc(*chosen, func(s string) bool { return s == it })
			})
		}
		ph := placeholder
		if len(*chosen) > 0 {
			ph = ""
		}
		parts := coss.ComboboxChipsInput(c, text, ph, len(*chosen) > 0)
		matches := coss.ComboboxFilter(parts, *text, cbFruits)
		coss.ComboboxPopup(c, parts, func() {
			if len(matches) == 0 {
				coss.ComboboxEmpty(c, "No items found.")
			}
			for _, it := range matches {
				coss.ComboboxItem(c, parts, it, slices.Contains(*chosen, it), func() { ui.Text(c, it).SingleLine() })
			}
		})
		if it, ok := parts.Chosen(); ok {
			if i := slices.Index(*chosen, it); i >= 0 {
				*chosen = slices.Delete(*chosen, i, i+1)
			} else {
				*chosen = append(*chosen, it)
			}
			*text = ""
		}
	})
}

func PCombobox1(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		fruitCombobox(c, v, coss.ComboboxInputProps{Placeholder: "Select a item…"}, "No items found.")
	})
}

func PCombobox2(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "Orange" })
	host.Children(func() {
		fruitCombobox(c, v, coss.ComboboxInputProps{Placeholder: "Select an item…", Disabled: true}, "No items found.")
	})
}

func PCombobox3(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		fruitCombobox(c, v, coss.ComboboxInputProps{Placeholder: "Select an item...", Size: coss.InputSM}, "No results found.")
	})
}

func PCombobox4(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		fruitCombobox(c, v, coss.ComboboxInputProps{Placeholder: "Select an item...", Size: coss.InputLG}, "No results found.")
	})
}

func PCombobox5(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).MaxWidth(256).Gap(8)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Label(c, "Fruits")
		fruitCombobox(c, v, coss.ComboboxInputProps{Placeholder: "Select an item..."}, "No results found.")
	})
}

func PCombobox6(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		fruitCombobox(c, v, coss.ComboboxInputProps{Placeholder: "Select an item..."}, "No results found.")
	})
}

func PCombobox7(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		fruitCombobox(c, v, coss.ComboboxInputProps{Placeholder: "Select a item…", ShowClear: true}, "No items found.")
	})
}

func PCombobox8(c *ui.Context) {
	groups := []struct {
		name  string
		items []string
	}{
		{"Status", []string{"Open", "In progress", "Blocked", "Resolved", "Closed"}},
		{"Priority", []string{"Low", "Medium", "High", "Urgent"}},
		{"Team", []string{"Design", "Frontend", "Backend", "DevOps", "QA", "Mobile", "Data", "Security", "Platform", "Infrastructure", "Product", "Marketing", "Sales", "Support", "Research", "Content", "Analytics", "Operations", "Finance", "HR", "Legal", "Growth", "Partner", "Community", "Docs", "Localization", "Accessibility", "SRE", "Release", "Architecture", "UX", "UI", "Management"}},
	}
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		cb := coss.ComboboxInput(c, v, coss.ComboboxInputProps{Placeholder: "e.g. feature"})
		coss.ComboboxPopup(c, cb.Parts, func() {
			shown := 0
			for _, g := range groups {
				matches := coss.ComboboxFilter(cb.Parts, *v, g.items)
				if len(matches) == 0 {
					continue
				}
				if shown > 0 {
					coss.SelectSeparator(c)
				}
				shown++
				coss.SelectGroupLabel(c, g.name)
				for _, it := range matches {
					coss.ComboboxItem(c, cb.Parts, it, it == *v, func() { ui.Text(c, it).SingleLine() })
				}
			}
			if shown == 0 {
				coss.ComboboxEmpty(c, "No tags found.")
			}
		})
		if it, ok := cb.Parts.Chosen(); ok {
			*v = it
		}
	})
}

func PCombobox9(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	chosen := ui.Local(host, "chosen", func() []string { return []string{"Apple", "Strawberry"} })
	text := ui.Local(host, "text", func() string { return "" })
	host.Children(func() { fruitChips(c, chosen, text, "Select a item...", false) })
}

// comboboxTriggerButton is the trigger of p-combobox-10: an outline button
// showing the value at its start and the chevrons at its end.
func comboboxTriggerButton(c *ui.Context, value string) *coss.ButtonElement {
	t := theme.Of(c)
	btn := coss.ButtonBox(c, func() {
		ui.Text(c, value).SingleLine().FontWeight(400).Grow(1)
		ui.Icon(c, icons.ChevronsUpDown()).Size(16, 16).MarginEnd(-4).TextColor(theme.Alpha(t.Foreground, 0.8))
	}).Variant(coss.ButtonOutline)
	btn.WidthPercent(100).Justify(ui.SpaceBetween)
	return btn
}

func PCombobox10(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return cbCountries[0] })
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		btn := comboboxTriggerButton(c, *v).OnClick(func() { *open = !*open })
		coss.ComboboxSearchPopup(c, btn.Element, open, v, cbCountries, "e.g. United Kingdom", "No countries found.")
	})
}

func PCombobox11(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).MaxWidth(256).Gap(16)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Favorite item")
			fruitCombobox(c, v, coss.ComboboxInputProps{Placeholder: "Select an item..."}, "No results found.")
		})
		coss.Button(c, "Submit").WidthPercent(100)
	})
}

func PCombobox12(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).MaxWidth(256).Gap(16)
	chosen := ui.Local(host, "chosen", func() []string { return nil })
	text := ui.Local(host, "text", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Favorite items")
			fruitChips(c, chosen, text, "Select items…", false)
		})
		coss.Button(c, "Submit").WidthPercent(100)
	})
}

func PCombobox13(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		fruitCombobox(c, v, coss.ComboboxInputProps{Placeholder: "Search items…", StartIcon: icons.Search()}, "No items found.")
	})
}

func PCombobox14(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	chosen := ui.Local(host, "chosen", func() []string { return []string{"Apple", "Grape"} })
	text := ui.Local(host, "text", func() string { return "" })
	host.Children(func() { fruitChips(c, chosen, text, "Select a item...", true) })
}

func PCombobox15(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		fruitCombobox(c, v, coss.ComboboxInputProps{Placeholder: "Select a item..."}, "No items found.").Radius(999)
	})
}

func cbTimezoneLabels() []string {
	out := make([]string, len(selectTimezones))
	for i, tz := range selectTimezones {
		out[i] = tz.Label
	}
	return out
}

func PCombobox16(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "(GMT+0) Europe/London" })
	host.Children(func() {
		zones := cbTimezoneLabels()
		cb := coss.ComboboxInput(c, v, coss.ComboboxInputProps{Placeholder: "Select timezone..."})
		matches := coss.ComboboxFilter(cb.Parts, *v, zones)
		coss.ComboboxPopup(c, cb.Parts, func() {
			if len(matches) == 0 {
				coss.ComboboxEmpty(c, "No timezones found.")
			}
			for _, it := range matches {
				coss.ComboboxItem(c, cb.Parts, it, it == *v, func() { ui.Text(c, it).SingleLine() })
			}
		})
		if it, ok := cb.Parts.Chosen(); ok {
			*v = it
		}
	})
}

// comboboxSelectButton is the trigger of p-combobox-17 and 18, a
// SelectButton: the select trigger's look on a plain button.
func comboboxSelectButton(c *ui.Context, open *bool, value, placeholder string) ui.Element {
	btn := ui.ButtonBase(c)
	coss.SelectTrigger(c, btn, coss.SelectDefault, func() {
		if value == "" {
			ui.Text(c, placeholder).SingleLine().TextColor(theme.Alpha(theme.Of(c).MutedForeground, 0.72))
			return
		}
		coss.SelectValue(c, value, "")
	})
	if btn.Clicked() {
		*open = !*open
	}
	return btn
}

func PCombobox17(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		btn := comboboxSelectButton(c, open, *v, "Select timezone")
		coss.ComboboxSearchPopup(c, btn, open, v, cbTimezoneLabels(), "e.g. Europe/London", "No timezones found.")
	})
}

func PCombobox18(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		btn := comboboxSelectButton(c, open, *v, "Select a fruit")
		coss.ComboboxSearchPopup(c, btn, open, v, cbFruits, "Search fruits...", "No items found.")
	})
}

type cbMember struct {
	avatar                    ui.ImageSource
	initials, label, priority string
	weight                    int
}

var cbTeam = []cbMember{
	{avatarFK, "JH", "Jenny Hamilton", "Highest", 200},
	{avatarNE, "PS", "Paul Smith", "Medium", 100},
	{avatarEW, "LW", "Luna Wyen", "High", 150},
	{nil, "AC", "Alex Chen", "Low", 100},
	{nil, "SJ", "Sarah Johnson", "Medium", 50},
	{nil, "ED", "Emma Davis", "Lowest", 100},
}

// teamCombobox is p-combobox-19 and 20: a combobox adding team members
// above the list of those added, which divided lays out as one bordered
// list instead of bordered rows.
func teamCombobox(c *ui.Context, divided bool) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).MaxWidth(320).Gap(8)
	chosen := ui.Local(host, "chosen", func() []string { return []string{"Jenny Hamilton", "Paul Smith"} })
	text := ui.Local(host, "text", func() string { return "" })
	host.Children(func() {
		labels := make([]string, len(cbTeam))
		for i, m := range cbTeam {
			labels[i] = m.label
		}
		cb := coss.ComboboxInput(c, text, coss.ComboboxInputProps{Placeholder: "Add team members…", StartIcon: icons.Search()})
		matches := coss.ComboboxFilter(cb.Parts, *text, labels)
		coss.ComboboxPopup(c, cb.Parts, func() {
			if len(matches) == 0 {
				coss.ComboboxEmpty(c, "No team members found.")
			}
			for _, it := range matches {
				coss.ComboboxItem(c, cb.Parts, it, slices.Contains(*chosen, it), func() { ui.Text(c, it).SingleLine() })
			}
		})
		if it, ok := cb.Parts.Chosen(); ok {
			if i := slices.Index(*chosen, it); i >= 0 {
				*chosen = slices.Delete(*chosen, i, i+1)
			} else {
				*chosen = append(*chosen, it)
			}
			*text = ""
		}
		if len(*chosen) == 0 {
			return
		}
		list := ui.Column(c).WidthPercent(100).Gap(8)
		if divided {
			list.Gap(0).Radius(t.RadiusLG).Border(1, t.Border)
		}
		list.Children(func() {
			for i, name := range *chosen {
				m := cbTeam[slices.Index(labels, name)]
				row := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8).Padding(4, 4, 4, 8).FontSize(14)
				if divided {
					if i > 0 {
						row.BorderWidth(1, 0, 0, 0).BorderColor(t.Border)
					}
				} else {
					row.Radius(t.RadiusLG).Border(1, t.Input)
				}
				row.Children(func() {
					av := coss.Avatar(c, m.initials).Size(20)
					if m.avatar != nil {
						av.Image(m.avatar)
					}
					ui.Text(c, m.label).SingleLine().Ellipsis("…").FontWeight(500).Grow(1)
					badge := coss.Badge(c, m.priority).Variant(coss.BadgeOutline)
					weight := ui.Text(c, fmt.Sprintf("%d%%", m.weight)).FontFeatures("tnum").TextColor(t.MutedForeground)
					if divided {
						badge.Size(coss.BadgeSizeSM)
						weight.FontSize(12)
					}
					coss.IconButton(c, icons.X(), "Remove "+m.label).Variant(coss.ButtonGhost).Size(coss.SizeIconXS).
						OnClick(func() { *chosen = slices.DeleteFunc(*chosen, func(s string) bool { return s == name }) })
				})
			}
		})
	})
}

func PCombobox19(c *ui.Context) { teamCombobox(c, false) }

func PCombobox20(c *ui.Context) { teamCombobox(c, true) }
