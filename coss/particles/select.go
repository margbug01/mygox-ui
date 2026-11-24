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

//go:embed select.go
var selectSource string

func init() {
	Register(Entry{ID: "p-select-1", Family: "select", Description: "Basic select", View: PSelect1, Source: extractFunc(selectSource, "PSelect1")})
	Register(Entry{ID: "p-select-2", Family: "select", Description: "Small select", View: PSelect2, Source: extractFunc(selectSource, "PSelect2")})
	Register(Entry{ID: "p-select-3", Family: "select", Description: "Large select", View: PSelect3, Source: extractFunc(selectSource, "PSelect3")})
	Register(Entry{ID: "p-select-4", Family: "select", Description: "Disabled select", View: PSelect4, Source: extractFunc(selectSource, "PSelect4")})
	Register(Entry{ID: "p-select-5", Family: "select", Description: "Select without item alignment", View: PSelect5, Source: extractFunc(selectSource, "PSelect5")})
	Register(Entry{ID: "p-select-6", Family: "select", Description: "Select with grouped options", View: PSelect6, Source: extractFunc(selectSource, "PSelect6")})
	Register(Entry{ID: "p-select-7", Family: "select", Description: "Multiple select", View: PSelect7, Source: extractFunc(selectSource, "PSelect7")})
	Register(Entry{ID: "p-select-8", Family: "select", Description: "Select with leading icon", View: PSelect8, Source: extractFunc(selectSource, "PSelect8")})
	Register(Entry{ID: "p-select-9", Family: "select", Description: "Select options with icons", View: PSelect9, Source: extractFunc(selectSource, "PSelect9")})
	Register(Entry{ID: "p-select-10", Family: "select", Description: "Select with object values", View: PSelect10, Source: extractFunc(selectSource, "PSelect10")})
	Register(Entry{ID: "p-select-11", Family: "select", Description: "Select in a form", View: PSelect11, Source: extractFunc(selectSource, "PSelect11")})
	Register(Entry{ID: "p-select-12", Family: "select", Description: "Select with disabled items", View: PSelect12, Source: extractFunc(selectSource, "PSelect12")})
	Register(Entry{ID: "p-select-13", Family: "select", Description: "Timezone select", View: PSelect13, Source: extractFunc(selectSource, "PSelect13")})
	Register(Entry{ID: "p-select-14", Family: "select", Description: "Status select with colored dot", View: PSelect14, Source: extractFunc(selectSource, "PSelect14")})
	Register(Entry{ID: "p-select-15", Family: "select", Description: "Pill-shaped select trigger", View: PSelect15, Source: extractFunc(selectSource, "PSelect15")})
	Register(Entry{ID: "p-select-16", Family: "select", Description: "Select with left text label", View: PSelect16, Source: extractFunc(selectSource, "PSelect16")})
	Register(Entry{ID: "p-select-17", Family: "select", Description: "Select with country flags", View: PSelect17, Source: extractFunc(selectSource, "PSelect17")})
	Register(Entry{ID: "p-select-18", Family: "select", Description: "Select with option descriptions", View: PSelect18, Source: extractFunc(selectSource, "PSelect18")})
	Register(Entry{ID: "p-select-19", Family: "select", Description: "Select with user avatars", View: PSelect19, Source: extractFunc(selectSource, "PSelect19")})
	Register(Entry{ID: "p-select-20", Family: "select", Description: "Rich select with avatars and usernames", View: PSelect20, Source: extractFunc(selectSource, "PSelect20")})
	Register(Entry{ID: "p-select-21", Family: "select", Description: "Auto-width select", View: PSelect21, Source: extractFunc(selectSource, "PSelect21")})
	Register(Entry{ID: "p-select-22", Family: "select", Description: "Select with custom tint border", View: PSelect22, Source: extractFunc(selectSource, "PSelect22")})
	Register(Entry{ID: "p-select-23", Family: "select", Description: "Select with field label", View: PSelect23, Source: extractFunc(selectSource, "PSelect23")})
	Register(Entry{ID: "p-select-24", Family: "select", Description: "Color palette select with swatches", View: PSelect24, Source: extractFunc(selectSource, "PSelect24")})
}

var selFrameworks = []coss.SelectOption[string]{
	{Value: "next", Label: "Next.js"},
	{Value: "vite", Label: "Vite"},
	{Value: "astro", Label: "Astro"},
}

func PSelect1(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "next" })
	host.Children(func() { coss.Select(c, v, selFrameworks) })
}

func PSelect2(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() { coss.Select(c, v, selFrameworks, "Select framework").Size(coss.SelectSM) })
}

func PSelect3(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() { coss.Select(c, v, selFrameworks, "Select framework").Size(coss.SelectLG) })
}

func PSelect4(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() { coss.Select(c, v, selFrameworks, "Select framework").Disabled(true).Opacity(0.64) })
}

func PSelect5(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() { coss.Select(c, v, selFrameworks, "Select framework") })
}

func PSelect6(c *ui.Context) {
	frontend := selFrameworks
	backend := []coss.SelectOption[string]{
		{Value: "express", Label: "Express"},
		{Value: "nestjs", Label: "NestJS"},
		{Value: "fastify", Label: "Fastify"},
		{Value: "django", Label: "Django"},
		{Value: "flask", Label: "Flask"},
		{Value: "rails", Label: "Rails"},
	}
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		label := ""
		for _, o := range append(append([]coss.SelectOption[string]{}, frontend...), backend...) {
			if o.Value == *v {
				label = o.Label
			}
		}
		parts := ui.SelectBase(c, v)
		coss.SelectTrigger(c, parts.Trigger, coss.SelectDefault, func() { coss.SelectValue(c, label, "Select framework") })
		coss.SelectPopup(c, parts, func() {
			for i, group := range []struct {
				name  string
				items []coss.SelectOption[string]
			}{{"Frontend", frontend}, {"Backend", backend}} {
				if i > 0 {
					coss.SelectSeparator(c)
				}
				coss.SelectGroupLabel(c, group.name)
				for _, o := range group.items {
					coss.SelectItem(c, parts, o.Value, o.Value == *v, false, func() { ui.Text(c, o.Label).SingleLine() })
				}
			}
		})
	})
}

func PSelect7(c *ui.Context) {
	languages := []coss.SelectOption[string]{
		{Value: "cpp", Label: "C++"},
		{Value: "csharp", Label: "C#"},
		{Value: "go", Label: "Go"},
		{Value: "java", Label: "Java"},
		{Value: "javascript", Label: "JavaScript"},
		{Value: "php", Label: "PHP"},
		{Value: "python", Label: "Python"},
		{Value: "rust", Label: "Rust"},
		{Value: "swift", Label: "Swift"},
		{Value: "typescript", Label: "TypeScript"},
	}
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	chosen := ui.Local(host, "chosen", func() []string { return []string{"javascript", "typescript"} })
	pick := ui.Local(host, "pick", func() string { return "" })
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		labelOf := map[string]string{}
		for _, o := range languages {
			labelOf[o.Value] = o.Label
		}
		text := "Select languages…"
		if len(*chosen) > 0 {
			text = labelOf[(*chosen)[0]]
			if len(*chosen) > 1 {
				text += fmt.Sprintf(" (+%d more)", len(*chosen)-1)
			}
		}
		parts := ui.SelectBaseOpen(c, pick, open)
		coss.SelectTrigger(c, parts.Trigger, coss.SelectDefault, func() { coss.SelectValue(c, text, "") })
		coss.SelectPopup(c, parts, func() {
			for _, o := range languages {
				coss.SelectItem(c, parts, o.Value, slices.Contains(*chosen, o.Value), false, func() { ui.Text(c, o.Label).SingleLine() })
			}
		})
		if *pick != "" {
			if i := slices.Index(*chosen, *pick); i >= 0 {
				*chosen = slices.Delete(*chosen, i, i+1)
			} else {
				*chosen = append(*chosen, *pick)
			}
			*pick = ""
			*open = true
		}
	})
}

func PSelect8(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "next" })
	host.Children(func() {
		label := ""
		for _, o := range selFrameworks {
			if o.Value == *v {
				label = o.Label
			}
		}
		parts := ui.SelectBase(c, v)
		coss.SelectTrigger(c, parts.Trigger, coss.SelectDefault, func() {
			ui.Icon(c, icons.Cable()).Size(16, 16).TextColor(theme.Alpha(t.Foreground, 0.8))
			coss.SelectValue(c, label, "")
		})
		coss.SelectPopup(c, parts, func() {
			for _, o := range selFrameworks {
				coss.SelectItem(c, parts, o.Value, o.Value == *v, false, func() { ui.Text(c, o.Label).SingleLine() })
			}
		})
	})
}

func PSelect9(c *ui.Context) {
	t := theme.Of(c)
	items := []struct {
		icon         *ui.SVG
		label, value string
	}{
		{icons.Layers(), "Components", "components"},
		{icons.Zap(), "Performance", "performance"},
		{icons.Globe(), "Network", "network"},
		{icons.Code2(), "Development", "development"},
	}
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "components" })
	host.Children(func() {
		parts := ui.SelectBase(c, v)
		coss.SelectTrigger(c, parts.Trigger, coss.SelectDefault, func() {
			for _, it := range items {
				if it.value == *v {
					ui.Icon(c, it.icon).Size(16, 16).TextColor(theme.Alpha(t.Foreground, 0.8))
					coss.SelectValue(c, it.label, "")
				}
			}
		})
		coss.SelectPopup(c, parts, func() {
			for _, it := range items {
				coss.SelectItem(c, parts, it.value, it.value == *v, false, func() {
					ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
						ui.Icon(c, it.icon).Size(16, 16)
						ui.Text(c, it.label).SingleLine()
					})
				})
			}
		})
	})
}

// selectDescribed shows a label above a muted description, as the two-line
// values of p-select-10 and p-select-18.
func selectDescribed(c *ui.Context, label, description string) {
	ui.Column(c).MinWidth(0).Children(func() {
		ui.Text(c, label).SingleLine().Ellipsis("…").LineHeight(20.0 / 14)
		ui.Text(c, description).SingleLine().Ellipsis("…").FontSize(12).LineHeight(16.0 / 12).TextColor(theme.Of(c).MutedForeground)
	})
}

func PSelect10(c *ui.Context) {
	items := []struct{ description, label, value string }{
		{"npx create-next-app", "Next.js", "next"},
		{"npm create vite@latest", "Vite", "vite"},
		{"npm create astro@latest", "Astro", "astro"},
		{"npx create-remix", "Remix", "remix"},
	}
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "next" })
	host.Children(func() {
		parts := ui.SelectBase(c, v)
		coss.SelectTrigger(c, parts.Trigger, coss.SelectDefault, func() {
			for _, it := range items {
				if it.value == *v {
					selectDescribed(c, it.label, it.description)
				}
			}
		}).PaddingY(4)
		coss.SelectPopup(c, parts, func() {
			for _, it := range items {
				coss.SelectItem(c, parts, it.value, it.value == *v, false, func() { selectDescribed(c, it.label, it.description) })
			}
		})
	})
}

func PSelect11(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).MaxWidth(256).Gap(16)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Framework")
			coss.Select(c, v, selFrameworks, "Select a framework")
			coss.FieldDescription(c, "Pick your favorite.")
		})
		coss.Button(c, "Submit").WidthPercent(100)
	})
}

func PSelect12(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "next" })
	host.Children(func() {
		coss.Select(c, v, []coss.SelectOption[string]{
			{Value: "next", Label: "Next.js"},
			{Value: "vite", Label: "Vite"},
			{Value: "astro", Label: "Astro (coming soon)", Disabled: true},
			{Value: "remix", Label: "Remix (coming soon)", Disabled: true},
			{Value: "nuxt", Label: "Nuxt"},
		})
	})
}

// selectTimezones lists zones as upstream labels them, "(GMT±h) Area/City",
// sorted by offset; upstream lists every zone Intl knows.
var selectTimezones = []coss.SelectOption[string]{
	{Value: "Pacific/Honolulu", Label: "(GMT-10) Pacific/Honolulu"},
	{Value: "America/Anchorage", Label: "(GMT-9) America/Anchorage"},
	{Value: "America/Los_Angeles", Label: "(GMT-8) America/Los Angeles"},
	{Value: "America/Denver", Label: "(GMT-7) America/Denver"},
	{Value: "America/Chicago", Label: "(GMT-6) America/Chicago"},
	{Value: "America/New_York", Label: "(GMT-5) America/New York"},
	{Value: "America/Sao_Paulo", Label: "(GMT-3) America/Sao Paulo"},
	{Value: "Europe/London", Label: "(GMT+0) Europe/London"},
	{Value: "Europe/Berlin", Label: "(GMT+1) Europe/Berlin"},
	{Value: "Europe/Paris", Label: "(GMT+1) Europe/Paris"},
	{Value: "Africa/Cairo", Label: "(GMT+2) Africa/Cairo"},
	{Value: "Europe/Moscow", Label: "(GMT+3) Europe/Moscow"},
	{Value: "Asia/Dubai", Label: "(GMT+4) Asia/Dubai"},
	{Value: "Asia/Kolkata", Label: "(GMT+5:30) Asia/Kolkata"},
	{Value: "Asia/Shanghai", Label: "(GMT+8) Asia/Shanghai"},
	{Value: "Asia/Tokyo", Label: "(GMT+9) Asia/Tokyo"},
	{Value: "Australia/Sydney", Label: "(GMT+11) Australia/Sydney"},
	{Value: "Pacific/Auckland", Label: "(GMT+13) Pacific/Auckland"},
}

func PSelect13(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "Europe/London" })
	host.Children(func() { coss.Select(c, v, selectTimezones) })
}

func PSelect14(c *ui.Context) {
	items := []struct {
		color        ui.Color
		label, value string
	}{
		{ui.Hex("#10b981"), "Completed", "completed"},
		{ui.Hex("#3b82f6"), "In Progress", "in-progress"},
		{ui.Hex("#f59e0b"), "Pending", "pending"},
		{ui.Hex("#6b7280"), "Cancelled", "cancelled"},
		{ui.Hex("#ef4444"), "Failed", "failed"},
	}
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "completed" })
	host.Children(func() {
		parts := ui.SelectBase(c, v)
		coss.SelectTrigger(c, parts.Trigger, coss.SelectDefault, func() {
			for _, it := range items {
				if it.value == *v {
					ui.Box(c).Size(8, 8).Radius(999).Background(it.color).Shrink(0)
					coss.SelectValue(c, it.label, "")
				}
			}
		})
		coss.SelectPopup(c, parts, func() {
			for _, it := range items {
				coss.SelectItem(c, parts, it.value, it.value == *v, false, func() {
					ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
						ui.Box(c).Size(8, 8).Radius(999).Background(it.color).Shrink(0)
						ui.Text(c, it.label).SingleLine()
					})
				})
			}
		})
	})
}

func PSelect15(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "active" })
	host.Children(func() {
		coss.Select(c, v, []coss.SelectOption[string]{
			{Value: "active", Label: "Active"},
			{Value: "inactive", Label: "Inactive"},
			{Value: "archived", Label: "Archived"},
		}).Radius(999)
	})
}

func PSelect16(c *ui.Context) {
	t := theme.Of(c)
	items := []coss.SelectOption[string]{
		{Value: "javascript", Label: "JavaScript"},
		{Value: "typescript", Label: "TypeScript"},
		{Value: "python", Label: "Python"},
		{Value: "go", Label: "Go"},
	}
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "javascript" })
	host.Children(func() {
		parts := ui.SelectBase(c, v)
		coss.SelectTrigger(c, parts.Trigger, coss.SelectDefault, func() {
			for _, o := range items {
				if o.Value == *v {
					ui.RichText(c).SingleLine().Children(func() {
						ui.Text(c, "Language: ").TextColor(t.MutedForeground)
						ui.Text(c, o.Label)
					})
				}
			}
		})
		coss.SelectPopup(c, parts, func() {
			for _, o := range items {
				coss.SelectItem(c, parts, o.Value, o.Value == *v, false, func() { ui.Text(c, o.Label).SingleLine() })
			}
		})
	})
}

func PSelect17(c *ui.Context) {
	type country struct{ flag, label, value string }
	continents := []struct {
		name  string
		items []country
	}{
		{"America", []country{{"🇺🇸", "United States", "us"}, {"🇨🇦", "Canada", "ca"}, {"🇲🇽", "Mexico", "mx"}}},
		{"Europe", []country{{"🇬🇧", "United Kingdom", "gb"}, {"🇫🇷", "France", "fr"}, {"🇩🇪", "Germany", "de"}}},
		{"Asia", []country{{"🇨🇳", "China", "cn"}, {"🇯🇵", "Japan", "jp"}, {"🇮🇳", "India", "in"}}},
	}
	flagged := func(it country) {
		ui.Row(c).AlignItems(ui.Center).Gap(8).MinWidth(0).Children(func() {
			ui.Text(c, it.flag).FontSize(16).LineHeight(1)
			ui.Text(c, it.label).SingleLine().Ellipsis("…")
		})
	}
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "ca" })
	host.Children(func() {
		parts := ui.SelectBase(c, v)
		coss.SelectTrigger(c, parts.Trigger, coss.SelectDefault, func() {
			for _, g := range continents {
				for _, it := range g.items {
					if it.value == *v {
						flagged(it)
					}
				}
			}
		})
		coss.SelectPopup(c, parts, func() {
			for i, g := range continents {
				if i > 0 {
					coss.SelectSeparator(c)
				}
				coss.SelectGroupLabel(c, g.name)
				for _, it := range g.items {
					coss.SelectItem(c, parts, it.value, it.value == *v, false, func() { flagged(it) })
				}
			}
		})
	})
}

func PSelect18(c *ui.Context) {
	items := []struct{ description, label, value string }{
		{"Ideal for individuals", "Standard Plan", "standard"},
		{"For professional users", "Pro Plan", "pro"},
		{"Built for large teams", "Enterprise Plan", "enterprise"},
	}
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "pro" })
	host.Children(func() {
		parts := ui.SelectBase(c, v)
		coss.SelectTrigger(c, parts.Trigger, coss.SelectDefault, func() {
			for _, it := range items {
				if it.value == *v {
					coss.SelectValue(c, it.label, "")
				}
			}
		})
		coss.SelectPopup(c, parts, func() {
			for _, it := range items {
				coss.SelectItem(c, parts, it.value, it.value == *v, false, func() { selectDescribed(c, it.label, it.description) })
			}
		})
	})
}

var selectUsers = []struct {
	avatar                           *ui.SVG
	initials, label, username, value string
}{
	{avatarPortrait1, "JH", "Jenny Hamilton", "@jennycodes", "jenny"},
	{avatarPortrait2, "PS", "Paul Smith", "@paulsmith", "paul"},
	{avatarPortrait3, "LW", "Luna Wyen", "@wyen.luna", "luna"},
}

func PSelect19(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "jenny" })
	host.Children(func() {
		user := func(i int) {
			u := selectUsers[i]
			ui.Row(c).AlignItems(ui.Center).Gap(8).MinWidth(0).Children(func() {
				coss.Avatar(c, u.initials).Size(20).Image(u.avatar)
				ui.Text(c, u.label).SingleLine().Ellipsis("…")
			})
		}
		parts := ui.SelectBase(c, v)
		coss.SelectTrigger(c, parts.Trigger, coss.SelectDefault, func() {
			for i, u := range selectUsers {
				if u.value == *v {
					user(i)
				}
			}
		})
		coss.SelectPopup(c, parts, func() {
			coss.SelectGroupLabel(c, "Impersonate user")
			for i, u := range selectUsers {
				coss.SelectItem(c, parts, u.value, u.value == *v, false, func() { user(i) })
			}
		})
	})
}

func PSelect20(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "jenny" })
	host.Children(func() {
		user := func(i int) {
			u := selectUsers[i]
			ui.Row(c).AlignItems(ui.Center).Gap(8).MinWidth(0).Children(func() {
				coss.Avatar(c, u.initials).Image(u.avatar)
				ui.Column(c).MinWidth(0).Children(func() {
					ui.Text(c, u.label).SingleLine().Ellipsis("…").FontWeight(500).LineHeight(20.0 / 14)
					ui.Text(c, u.username).SingleLine().Ellipsis("…").FontSize(12).LineHeight(16.0 / 12).TextColor(t.MutedForeground)
				})
			})
		}
		parts := ui.SelectBase(c, v)
		coss.SelectTrigger(c, parts.Trigger, coss.SelectDefault, func() {
			for i, u := range selectUsers {
				if u.value == *v {
					user(i)
				}
			}
		}).PaddingY(6)
		coss.SelectPopup(c, parts, func() {
			for i, u := range selectUsers {
				coss.SelectItem(c, parts, u.value, u.value == *v, false, func() { user(i) }).PaddingY(6)
			}
		})
	})
}

func PSelect21(c *ui.Context) {
	host := ui.Box(c)
	v := ui.Local(host, "v", func() string { return "next" })
	host.Children(func() { coss.Select(c, v, selFrameworks).MinWidth(144) })
}

func PSelect22(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "next" })
	host.Children(func() {
		coss.Select(c, v, selFrameworks).Border(1, ui.Transparent).Background(t.Muted).NoShadow()
	})
}

func PSelect23(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "apple" })
	host.Children(func() {
		coss.SelectLabel(c, "Fruits")
		coss.Select(c, v, []coss.SelectOption[string]{
			{Value: "apple", Label: "Apple"},
			{Value: "banana", Label: "Banana"},
			{Value: "orange", Label: "Orange"},
		})
	})
}

func PSelect24(c *ui.Context) {
	t := theme.Of(c)
	palettes := []struct {
		label, value string
		colors       [3]string
	}{
		{"Mono", "mono", [3]string{"#18181B", "#A1A1AA", "#FAFAFA"}},
		{"Porcelain", "porcelain", [3]string{"#AA8F73", "#C8B79A", "#889FA3"}},
		{"Ember", "ember", [3]string{"#B94724", "#DB8649", "#EDD0A0"}},
		{"Terra", "terra", [3]string{"#A64F3C", "#C78F70", "#896577"}},
		{"Rosé", "rose", [3]string{"#9B3F60", "#C98693", "#E9B9A5"}},
		{"Dusk", "dusk", [3]string{"#BE835B", "#8E7C9F", "#6E97AE"}},
		{"Orchard", "orchard", [3]string{"#63794B", "#ADA66B", "#83A18A"}},
		{"Petrol", "petrol", [3]string{"#326A76", "#C97868", "#C9AD81"}},
		{"Lagoon", "lagoon", [3]string{"#247E92", "#4CAFA3", "#B6D8C7"}},
		{"Borealis", "borealis", [3]string{"#267F69", "#527FB8", "#9A83BF"}},
		{"Cobalt", "cobalt", [3]string{"#244E9A", "#607DA9", "#B5C9DD"}},
		{"Iris", "iris", [3]string{"#7361A3", "#AC87A5", "#DBC3CA"}},
	}
	host := ui.Column(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "dusk" })
	host.Children(func() {
		swatched := func(i int) {
			p := palettes[i]
			ui.Row(c).AlignItems(ui.Center).Gap(8).MinWidth(0).Children(func() {
				ui.Row(c).Shrink(0).Children(func() {
					for k, col := range p.colors {
						dot := ui.Box(c).Size(18, 18).Radius(999).Background(ui.Hex(col)).Border(1, t.Background)
						if k > 0 {
							dot.MarginStart(-6)
						}
					}
				})
				ui.Text(c, p.label).SingleLine().Ellipsis("…")
			})
		}
		coss.SelectLabel(c, "Color palette")
		parts := ui.SelectBase(c, v)
		coss.SelectTrigger(c, parts.Trigger, coss.SelectDefault, func() {
			for i, p := range palettes {
				if p.value == *v {
					swatched(i)
				}
			}
		})
		coss.SelectPopup(c, parts, func() {
			for i, p := range palettes {
				coss.SelectItem(c, parts, p.value, p.value == *v, false, func() { swatched(i) })
			}
		})
	})
}
