package particles

import (
	_ "embed"
	"fmt"
	"math"
	"math/rand/v2"
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
	Register(Entry{ID: "p-progress-bounce", Family: "progress", Description: "Progress bar that springs to each value and throws sparks going forward (mygox motion, not on coss.com)", View: PProgressBounce, Source: extractFunc(miscSource, "PProgressBounce")})
	Register(Entry{ID: "p-skeleton-2", Family: "skeleton", Description: "User card loading skeleton", View: PSkeleton2, Source: extractFunc(miscSource, "PSkeleton2")})
}

// PKbd1 displays single keys and keyboard shortcut combinations.
func PKbd1(c *ui.Context) {
	t := theme.Of(c)
	ui.Column(c).Gap(16).Children(func() {
		ui.Column(c).Gap(8).Children(func() {
			ui.Text(c, "Single keys:").FixedLineHeight(20).TextColor(t.MutedForeground)
			ui.Row(c).Gap(8).Children(func() {
				coss.Kbd(c, "K")
				coss.Kbd(c, "⌘")
				coss.Kbd(c, "⌃")
				coss.Kbd(c, "⇧")
			})
		})
		ui.Column(c).Gap(8).Children(func() {
			ui.Text(c, "Key combinations:").FixedLineHeight(20).TextColor(t.MutedForeground)
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
			ui.Text(c, "coss ui").FixedLineHeight(20).FontWeight(500).TextColor(t.Foreground)
			ui.Text(c, "Unstyled, accessible primitives for fast product UI and design systems.").
				FixedLineHeight(20).TextColor(t.MutedForeground)
		})
		coss.Separator(c).Margin(16, 0)
		ui.Row(c).AlignItems(ui.Center).Gap(16).FixedLineHeight(20).TextColor(t.Foreground).Children(func() {
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
	coss.Spinner(c).Size(24) // lucide renders 24px when no size class applies
}

// PProgress1 renders a progress bar that starts at 20 and grows by a
// random step every second until it is full.
// PProgressBounce is mygox's own: the bar springs to each new value,
// running a little past it, and throws sparks as it moves forward.
func PProgressBounce(c *ui.Context) {
	box := ui.Column(c).WidthPercent(100).MaxWidth(256).Gap(12).AlignItems(ui.Start)
	value := ui.Local(box, "value", func() float32 { return 20 })
	box.Children(func() {
		coss.Progress(c, *value, coss.ProgressOptions{Bounce: 40, Sparks: true})
		coss.Button(c, "Advance").Variant(coss.ButtonOutline).OnClick(func() {
			if *value += 20; *value > 100 {
				*value = 0
			}
		})
	})
}

func PProgress1(c *ui.Context) {
	box := ui.Box(c).WidthPercent(100).MaxWidth(256)
	value := ui.Local(box, "value", func() float32 { return 20 })
	next := ui.Local(box, "next", func() time.Time { return c.Now().Add(time.Second) })
	if *value < 100 {
		if !c.Now().Before(*next) {
			*value = min(100, float32(math.Round(float64(*value)+rand.Float64()*25)))
			*next = c.Now().Add(time.Second)
		}
		c.After(next.Sub(c.Now()))
	}
	box.Children(func() {
		coss.Progress(c, *value)
	})
}

// labeledProgress is Progress with ProgressLabel and ProgressValue above
// the track, as p-progress-2 and p-progress-3 compose them.
func labeledProgress(c *ui.Context, label, value string, percent float32) {
	ui.Column(c).WidthPercent(100).Gap(8).FixedLineHeight(20).Children(func() {
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(8).Children(func() {
			ui.Text(c, label).FontWeight(500)
			ui.Text(c, value).FontFeatures("tnum")
		})
		coss.Progress(c, percent)
	})
}

// PProgress2 renders a labeled progress bar with percentage value.
func PProgress2(c *ui.Context) {
	ui.Box(c).WidthPercent(100).MaxWidth(256).Children(func() {
		labeledProgress(c, "Export data", "60%", 60)
	})
}

// PProgress3 renders a progress bar with custom "502 / 512" formatted readout.
func PProgress3(c *ui.Context) {
	const val, maxVal = 502, 512
	ui.Box(c).WidthPercent(100).MaxWidth(256).Children(func() {
		labeledProgress(c, "Upload", fmt.Sprintf("%d / %d", val, maxVal), float32(val)/float32(maxVal)*100)
	})
}

// PSkeleton2 renders a static user-row loading skeleton.
func PSkeleton2(c *ui.Context) {
	ui.Row(c).WidthPercent(100).MaxWidth(368).AlignItems(ui.Center).Gap(16).Children(func() {
		coss.Skeleton(c).Size(40, 40).Radius(999).Shrink(0)
		ui.Column(c).Grow(1).Children(func() {
			coss.Skeleton(c).Height(16).WidthPercent(100).MaxWidth(216).Margin(2, 0)
			ui.Row(c).Gap(4).WidthPercent(100).MaxWidth(216).Children(func() {
				coss.Skeleton(c).Height(16).Grow(1).Basis(0).Margin(2, 0)
				coss.Skeleton(c).Height(16).Grow(1).Basis(0).Margin(2, 0)
			})
		})
		coss.Skeleton(c).Height(24).Width(68).Shrink(0)
	})
}
