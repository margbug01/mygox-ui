package particles

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed misc.go
var miscSource string

func init() {
	Register(Entry{ID: "p-kbd-1", Family: "kbd", Description: "Single keys and keyboard shortcut combinations", View: PKbd1, Source: extractFunc(miscSource, "PKbd1")})
	Register(Entry{ID: "p-separator-1", Family: "separator", Description: "Horizontal and vertical separators", View: PSeparator1, Source: extractFunc(miscSource, "PSeparator1")})
	Register(Entry{ID: "p-spinner-1", Family: "spinner", Description: "Basic loading spinner", View: PSpinner1, Source: extractFunc(miscSource, "PSpinner1")})
	Register(Entry{ID: "p-progress-1", Family: "progress", Description: "Animated progress bar", View: PProgress1, Source: extractFunc(miscSource, "PProgress1")})
	Register(Entry{ID: "p-progress-2", Family: "progress", Description: "Progress bar with label and percentage", View: PProgress2, Source: extractFunc(miscSource, "PProgress2")})
	Register(Entry{ID: "p-progress-3", Family: "progress", Description: "Progress bar with custom formatted value", View: PProgress3, Source: extractFunc(miscSource, "PProgress3")})
	Register(Entry{ID: "p-skeleton-2", Family: "skeleton", Description: "User card loading skeleton", View: PSkeleton2, Source: extractFunc(miscSource, "PSkeleton2")})
}

// PKbd1 displays single keys and keyboard shortcut combinations.
func PKbd1(c *ui.Context) {
	t := theme.Of(c)
	ui.Column(c).Gap(16).Children(func() {
		ui.Column(c).Gap(8).Children(func() {
			ui.Text(c, "Single keys:").FontSize(13).TextColor(t.MutedForeground)
			ui.Row(c).Gap(8).Children(func() {
				coss.Kbd(c, "K")
				coss.Kbd(c, "⌘")
				coss.Kbd(c, "⌃")
				coss.Kbd(c, "⇧")
			})
		})
		ui.Column(c).Gap(8).Children(func() {
			ui.Text(c, "Key combinations:").FontSize(13).TextColor(t.MutedForeground)
			ui.Row(c).Gap(8).Children(func() {
				coss.KbdGroup(c, func() { coss.Kbd(c, "⌘"); coss.Kbd(c, "K") })
				coss.KbdGroup(c, func() { coss.Kbd(c, "⌘"); coss.Kbd(c, "Shift"); coss.Kbd(c, "P") })
				coss.KbdGroup(c, func() { coss.Kbd(c, "Ctrl"); coss.Kbd(c, "Alt"); coss.Kbd(c, "Delete") })
			})
		})
	})
}

// PSeparator1 demonstrates horizontal and vertical separators.
func PSeparator1(c *ui.Context) {
	t := theme.Of(c)
	ui.Column(c).MaxWidth(288).Children(func() {
		ui.Column(c).Gap(4).Children(func() {
			ui.Text(c, "coss ui").FontSize(14).FontWeight(500).TextColor(t.Foreground)
			ui.Text(c, "Unstyled, accessible primitives for fast product UI and design systems.").
				FontSize(13).TextColor(t.MutedForeground)
		})
		coss.Separator(c).Margin(14, 0)
		ui.Row(c).AlignItems(ui.Center).Gap(14).FontSize(13).TextColor(t.Foreground).Children(func() {
			ui.Text(c, "Blog")
			coss.Separator(c, coss.Vertical)
			ui.Text(c, "Docs")
			coss.Separator(c, coss.Vertical)
			ui.Text(c, "Source")
			coss.Separator(c, coss.Vertical)
			ui.Text(c, "Releases")
		})
	})
}

// PSpinner1 renders a standalone coss spinner.
func PSpinner1(c *ui.Context) {
	coss.Spinner(c)
}

// PProgress1 renders an animated progress bar cycling smoothly across [20, 100].
func PProgress1(c *ui.Context) {
	box := ui.Box(c).WidthPercent(100)
	t := box.Loop("progress", 4*time.Second, ui.EaseInOut)
	box.Children(func() {
		coss.Progress(c, 20+80*t)
	})
}

// PProgress2 renders a labeled progress bar with percentage value.
func PProgress2(c *ui.Context) {
	t := theme.Of(c)
	ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
			ui.Text(c, "Export data").FontSize(13).FontWeight(500).TextColor(t.Foreground)
			ui.Text(c, "60%").FontSize(13).TextColor(t.MutedForeground)
		})
		coss.Progress(c, 60)
	})
}

// PProgress3 renders a progress bar with custom "502 / 512" formatted readout.
func PProgress3(c *ui.Context) {
	t := theme.Of(c)
	const val, maxVal = 502, 512
	ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
			ui.Text(c, "Upload").FontSize(13).FontWeight(500).TextColor(t.Foreground)
			ui.Text(c, fmt.Sprintf("%d / %d", val, maxVal)).FontSize(13).TextColor(t.MutedForeground)
		})
		coss.Progress(c, float32(val)/float32(maxVal)*100)
	})
}

// PSkeleton2 renders a static user-row loading skeleton.
func PSkeleton2(c *ui.Context) {
	ui.Row(c).WidthPercent(100).MaxWidth(368).AlignItems(ui.Center).Gap(14).Children(func() {
		coss.Skeleton(c).Size(40, 40).Radius(999).Shrink(0)
		ui.Column(c).Grow(1).Gap(6).Children(func() {
			coss.Skeleton(c).Height(14).WidthPercent(80)
			ui.Row(c).Gap(6).WidthPercent(70).Children(func() {
				coss.Skeleton(c).Height(12).Grow(1)
				coss.Skeleton(c).Height(12).Grow(1)
			})
		})
		coss.Skeleton(c).Height(24).Width(68).Shrink(0)
	})
}
