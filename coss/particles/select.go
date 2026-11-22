package particles

import (
	_ "embed"

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
	{Value: "mygo", Label: "MyGo"},
}

func PSelect1(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "next" })
	host.Children(func() { coss.Select(c, v, selFrameworks) })
}

func PSelect2(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "next" })
	host.Children(func() { coss.Select(c, v, selFrameworks).Height(28).FontSize(13) })
}

func PSelect3(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "next" })
	host.Children(func() { coss.Select(c, v, selFrameworks).Height(38) })
}

func PSelect4(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "next" })
	host.Children(func() { coss.Select(c, v, selFrameworks).Disabled(true).Opacity(0.64) })
}

func PSelect5(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() { coss.Select(c, v, selFrameworks, "Choose framework…") })
}

func PSelect6(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "go" })
	opts := []coss.SelectOption[string]{
		{Value: "go", Label: "Backend · Go"},
		{Value: "rust", Label: "Backend · Rust"},
		{Value: "ts", Label: "Frontend · TypeScript"},
		{Value: "swift", Label: "Mobile · Swift"},
	}
	host.Children(func() { coss.Select(c, v, opts) })
}

func PSelect7(c *ui.Context) {
	host := ui.Row(c).Gap(6).Wrap()
	sel := ui.Local(host, "sel", func() [4]bool { return [4]bool{true, false, true, false} })
	names := []string{"Next.js", "Vite", "Astro", "MyGo"}
	host.Children(func() {
		for i, name := range names {
			idx := i
			v := coss.BadgeSecondary
			if sel[idx] {
				v = coss.BadgeDefault
			}
			coss.BadgeButton(c, name).Variant(v).OnClick(func() { sel[idx] = !sel[idx] })
		}
	})
}

func PSelect8(c *ui.Context) {
	host := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8)
	v := ui.Local(host, "v", func() string { return "usb" })
	host.Children(func() {
		ui.Icon(c, icons.Cable()).Size(15, 15).TextColor(theme.Of(c).MutedForeground)
		coss.Select(c, v, []coss.SelectOption[string]{
			{Value: "usb", Label: "USB-C"},
			{Value: "thunderbolt", Label: "Thunderbolt 4"},
			{Value: "hdmi", Label: "HDMI 2.1"},
		}).Grow(1)
	})
}

func PSelect9(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8)
	v := ui.Local(host, "v", func() string { return "code" })
	iconMap := map[string]*ui.SVG{"code": icons.Code2(), "globe": icons.Globe(), "layers": icons.Layers(), "zap": icons.Zap()}
	host.Children(func() {
		if ic, ok := iconMap[*v]; ok {
			ui.Icon(c, ic).Size(15, 15).TextColor(t.MutedForeground)
		}
		coss.Select(c, v, []coss.SelectOption[string]{
			{Value: "code", Label: "Engineering"},
			{Value: "globe", Label: "Web Platform"},
			{Value: "layers", Label: "Design System"},
			{Value: "zap", Label: "Performance"},
		}).Grow(1)
	})
}

func PSelect10(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() int { return 101 })
	host.Children(func() {
		coss.Select(c, v, []coss.SelectOption[int]{
			{Value: 101, Label: "Object #101 · Starter"},
			{Value: 202, Label: "Object #202 · Pro"},
			{Value: 303, Label: "Object #303 · Enterprise"},
		})
	})
}

func PSelect11(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "mygo" })
	saved := ui.Local(host, "saved", func() string { return "" })
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() {
				coss.FieldLabel(c, "Primary framework")
				coss.Select(c, v, selFrameworks)
			})
			coss.Button(c, "Submit").Size(coss.SizeSM).OnClick(func() { *saved = *v })
			if *saved != "" {
				ui.Text(c, "Submitted: "+*saved).FontSize(12).TextColor(t.SuccessForeground)
			}
		})
	})
}

func PSelect12(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "pro" })
	host.Children(func() {
		coss.Select(c, v, []coss.SelectOption[string]{
			{Value: "free", Label: "Free"},
			{Value: "pro", Label: "Pro"},
			{Value: "legacy", Label: "Legacy (unavailable)", Disabled: true},
		})
	})
}

func PSelect13(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "shanghai" })
	host.Children(func() {
		coss.Select(c, v, []coss.SelectOption[string]{
			{Value: "shanghai", Label: "Asia/Shanghai (UTC+8)"},
			{Value: "tokyo", Label: "Asia/Tokyo (UTC+9)"},
			{Value: "london", Label: "Europe/London (UTC+0)"},
			{Value: "nyc", Label: "America/New_York (UTC-5)"},
		})
	})
}

func PSelect14(c *ui.Context) {
	host := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8)
	v := ui.Local(host, "v", func() string { return "active" })
	dotCol := ui.Hex("#10b981")
	switch *v {
	case "pending":
		dotCol = ui.Hex("#f59e0b")
	case "failed":
		dotCol = ui.Hex("#ef4444")
	}
	host.Children(func() {
		ui.Box(c).Size(8, 8).Radius(999).Background(dotCol).Shrink(0)
		coss.Select(c, v, []coss.SelectOption[string]{
			{Value: "active", Label: "Active"},
			{Value: "pending", Label: "Pending"},
			{Value: "failed", Label: "Failed"},
		}).Grow(1)
	})
}

func PSelect15(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "mygo" })
	host.Children(func() { coss.Select(c, v, selFrameworks).Radius(999) })
}

func PSelect16(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8)
	v := ui.Local(host, "v", func() string { return "mygo" })
	host.Children(func() {
		ui.Text(c, "Runtime:").FontSize(13).TextColor(t.MutedForeground).Shrink(0)
		coss.Select(c, v, selFrameworks).Grow(1)
	})
}

func PSelect17(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "cn" })
	host.Children(func() {
		coss.Select(c, v, []coss.SelectOption[string]{
			{Value: "cn", Label: "🇨🇳 China"},
			{Value: "us", Label: "🇺🇸 United States"},
			{Value: "jp", Label: "🇯🇵 Japan"},
			{Value: "gb", Label: "🇬🇧 United Kingdom"},
		})
	})
}

func PSelect18(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "fast" })
	host.Children(func() {
		coss.Select(c, v, []coss.SelectOption[string]{
			{Value: "fast", Label: "Fast — Low latency responses"},
			{Value: "balanced", Label: "Balanced — Everyday coding"},
			{Value: "deep", Label: "Deep — Complex reasoning"},
		})
	})
}

func PSelect19(c *ui.Context) {
	host := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8)
	v := ui.Local(host, "v", func() string { return "SJ" })
	host.Children(func() {
		coss.Avatar(c, *v).Size(24)
		coss.Select(c, v, []coss.SelectOption[string]{
			{Value: "SJ", Label: "Sarah Johnson"},
			{Value: "AR", Label: "Alex Rivera"},
			{Value: "MA", Label: "Mark Andersson"},
		}).Grow(1)
	})
}

func PSelect20(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8)
	v := ui.Local(host, "v", func() string { return "SJ" })
	handles := map[string]string{"SJ": "@sarahj", "AR": "@arivera", "MA": "@marka"}
	host.Children(func() {
		coss.Avatar(c, *v).Size(28)
		ui.Column(c).Grow(1).Gap(2).Children(func() {
			coss.Select(c, v, []coss.SelectOption[string]{
				{Value: "SJ", Label: "Sarah Johnson"},
				{Value: "AR", Label: "Alex Rivera"},
				{Value: "MA", Label: "Mark Andersson"},
			})
			ui.Text(c, handles[*v]).FontSize(11).TextColor(t.MutedForeground)
		})
	})
}

func PSelect21(c *ui.Context) {
	host := ui.Box(c).Width(150)
	v := ui.Local(host, "v", func() string { return "mygo" })
	host.Children(func() { coss.Select(c, v, selFrameworks) })
}

func PSelect22(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "mygo" })
	host.Children(func() {
		coss.Select(c, v, selFrameworks).Border(1, theme.Alpha(t.Info, 0.45)).Background(theme.Alpha(t.Info, 0.06))
	})
}

func PSelect23(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "mygo" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "UI Framework")
			coss.Select(c, v, selFrameworks)
		})
	})
}

func PSelect24(c *ui.Context) {
	host := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8)
	v := ui.Local(host, "v", func() string { return "ocean" })
	palettes := map[string][3]ui.Color{
		"ocean":  {ui.Hex("#0284c7"), ui.Hex("#38bdf8"), ui.Hex("#bae6fd")},
		"forest": {ui.Hex("#059669"), ui.Hex("#34d399"), ui.Hex("#a7f3d0")},
		"sunset": {ui.Hex("#ea580c"), ui.Hex("#fb923c"), ui.Hex("#fed7aa")},
	}
	host.Children(func() {
		sw := palettes[*v]
		ui.Row(c).Gap(3).Shrink(0).Children(func() {
			for _, col := range sw {
				ui.Box(c).Size(10, 10).Radius(999).Background(col)
			}
		})
		coss.Select(c, v, []coss.SelectOption[string]{
			{Value: "ocean", Label: "Ocean"},
			{Value: "forest", Label: "Forest"},
			{Value: "sunset", Label: "Sunset"},
		}).Grow(1)
	})
}
