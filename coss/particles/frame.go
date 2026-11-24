package particles

import (
	_ "embed"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed frame.go
var frameSource string

func init() {
	Register(Entry{ID: "p-frame-1", Family: "frame", Description: "Basic frame with header, panel, and footer", View: PFrame1, Source: extractFunc(frameSource, "PFrame1")})
	Register(Entry{ID: "p-frame-2", Family: "frame", Description: "Collapsible frame with delete action", View: PFrame2, Source: extractFunc(frameSource, "PFrame2")})
	Register(Entry{ID: "p-frame-3", Family: "frame", Description: "Frame with separated panels", View: PFrame3, Source: extractFunc(frameSource, "PFrame3")})
	Register(Entry{ID: "p-frame-4", Family: "frame", Description: "Frame with stacked panels and separator", View: PFrame4, Source: extractFunc(frameSource, "PFrame4")})
	Register(Entry{ID: "p-empty-1", Family: "empty", Description: "Empty state with icon media and actions", View: PEmpty1, Source: extractFunc(frameSource, "PEmpty1")})
	Register(Entry{ID: "p-skeleton-1", Family: "skeleton", Description: "Animated loading skeleton transitioning to user cards", View: PSkeleton1, Source: extractFunc(frameSource, "PSkeleton1")})
}

// frameSection is the h2 + p pair the frame particles put in a panel.
func frameSection(c *ui.Context, title string) {
	t := theme.Of(c)
	ui.Text(c, title).FontWeight(600).TextColor(t.Foreground)
	ui.Text(c, "Section description").TextColor(t.MutedForeground)
}

func PFrame1(c *ui.Context) {
	t := theme.Of(c)
	coss.Frame(c, func() {
		coss.FrameHeader(c, func() {
			coss.FrameTitle(c, "Section header")
			coss.FrameDescription(c, "Brief description about the section")
		})
		coss.FramePanel(c, func() { frameSection(c, "Section title") })
		coss.FrameFooter(c, func() {
			ui.Text(c, "Footer").FixedLineHeight(20).TextColor(t.MutedForeground)
		})
	})
}

func PFrame2(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Frame(c, func() {
			var p ui.CollapsibleParts
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Padding(8).Children(func() {
				p = ui.CollapsibleBase(c, open)
				tr := p.Trigger.AlignItems(ui.Center).Gap(8).Height(32).Padding(0, 11).Radius(t.RadiusLG).FontWeight(500).Cursor(ui.CursorPointer)
				if tr.Hovered() {
					tr.Background(t.Accent)
				}
				tr.Children(func() {
					ui.Icon(c, icons.ChevronDown()).Size(16, 16).Opacity(0.8).Margin(0, -2).Rotate(180 * p.Progress())
					ui.Text(c, "Section header")
				})
				coss.IconButton(c, icons.Trash(), "Delete").Variant(coss.ButtonGhost)
			})
			p.Panel(func() {
				coss.FramePanel(c, func() { frameSection(c, "Section title") })
			})
		})
	})
}

func PFrame3(c *ui.Context) {
	coss.Frame(c, func() {
		coss.FrameHeader(c, func() {
			coss.FrameTitle(c, "Section header")
			coss.FrameDescription(c, "Brief description about the section")
		})
		coss.FramePanel(c, func() { frameSection(c, "Separated panel") })
		coss.FramePanel(c, func() { frameSection(c, "Separated panel") })
	})
}

func PFrame4(c *ui.Context) {
	coss.Frame(c, func() {
		coss.FrameHeader(c, func() {
			coss.FrameTitle(c, "Section header")
			coss.FrameDescription(c, "Brief description about the section")
		})
		coss.FramePanel(c, func() {
			ui.Column(c).WidthPercent(100).Padding(20).Children(func() { frameSection(c, "Stacked panel") })
			coss.Separator(c)
			ui.Column(c).WidthPercent(100).Padding(20).Children(func() { frameSection(c, "Stacked panel") })
		}).Padding(0)
	})
}

func PEmpty1(c *ui.Context) {
	coss.Empty(c, func() {
		coss.EmptyHeader(c, func() {
			coss.EmptyMedia(c, icons.Route())
			coss.EmptyTitle(c, "No upcoming meetings")
			coss.EmptyDescription(c, "Create a meeting to get started.")
		})
		coss.EmptyContent(c, func() {
			ui.Row(c).Gap(8).Children(func() {
				coss.Button(c, "Create meeting").Size(coss.SizeSM)
				coss.Button(c, "View docs").Size(coss.SizeSM).Variant(coss.ButtonOutline).LeadingIcon(icons.Book())
			})
		})
	})
}

func PSkeleton1(c *ui.Context) {
	t := theme.Of(c)
	type user struct {
		delay     time.Duration
		fallback  string
		name      string
		role      string
		followers string
	}
	users := []user{
		{1500 * time.Millisecond, "SJ", "Sarah Johnson", "Design Engineer", "15k followers"},
		{2200 * time.Millisecond, "MA", "Mark Bennett Andersson", "Product Designer", "8k followers"},
		{1800 * time.Millisecond, "AR", "Alex Rivera", "UI/UX Designer", "12k followers"},
	}

	root := ui.Column(c).WidthPercent(100).MaxWidth(368).Gap(24)
	started := ui.Local(root, "started", func() time.Time { return c.Now() })
	elapsed := c.Now().Sub(*started)
	if elapsed < 2300*time.Millisecond {
		c.After(200 * time.Millisecond)
	}
	root.Children(func() {
		for _, u := range users {
			if elapsed < u.delay {
				PSkeleton2(c)
				continue
			}
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(16).Children(func() {
				coss.Avatar(c, u.fallback).Size(40)
				ui.Column(c).Grow(1).MinWidth(0).Gap(4).Children(func() {
					ui.Text(c, u.name).SingleLine().Ellipsis("…").FixedLineHeight(20).FontWeight(500).TextColor(t.Foreground)
					ui.Row(c).AlignItems(ui.Center).Gap(12).TextColor(t.MutedForeground).FontSize(12).FixedLineHeight(16).Children(func() {
						ui.Text(c, u.role).SingleLine()
						ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
							ui.Icon(c, icons.UsersRound()).Size(12, 12)
							ui.Text(c, u.followers).SingleLine()
						})
					})
				})
				coss.Button(c, "Follow").Size(coss.SizeXS).LeadingIcon(icons.UserRoundPlus())
			})
		}
	})
}
