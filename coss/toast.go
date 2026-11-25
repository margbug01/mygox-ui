package coss

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// ToastViewport gives the window's toasts, added with Context.AddToast,
// coss's look: popover cards in the bottom right corner, newest last,
// with an icon for the Types "success", "error", "info", "warning" and
// "loading". Call it in a view that adds toasts; the first call in a frame
// builds them.
func ToastViewport(c *ui.Context) {
	t := theme.Of(c)
	ui.ToastViewportBase(c, func(viewport ui.Element, toasts []ui.Toast) {
		viewport.Padding(32).AlignItems(ui.End).Gap(12)
		for _, ts := range toasts {
			parts := ui.ToastBase(c, ts)
			parts.Root.Width(360).MaxWidthPercent(100).Radius(t.RadiusLG).Border(1, t.Border).
				Background(t.Popover).TextColor(t.PopoverForeground).
				Shadow(0, 10, 15, -3, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
				Children(func() {
					var action func()
					if ts.Action != "" {
						action = func() {
							parts.ActionButton().Height(24).Padding(0, 7).Radius(t.RadiusMD).
								Background(t.Primary).TextColor(t.PrimaryForeground).Cursor(ui.CursorPointer).
								Children(func() {
									ui.Text(c, ts.Action).SingleLine().FontSize(12).FontWeight(500)
								})
						}
					}
					toastBody(c, ts, action)
				})
		}
	})
}

// ToastStackOptions shape a ToastStack; the zero value is coss's look.
type ToastStackOptions struct {
	Width, Height float32 // of a toast
	Slots         int     // the most toasts it holds, which its height is for; 3 when 0
	Peek          float32 // px of each tucked toast that shows; 14 when 0
	Gap           float32 // px between the toasts fanned out; 10 when 0
	Rest          float32 // px from the front toast, while they are tucked, to the stack's bottom
	Corner        float32
	Background    ui.Color
	Border        ui.Color // a tucked toast's outline
	Ring          ui.Color // drawn round the toasts shown, as coss's border
	// Content builds toast k in its card in place of coss's icon, title and
	// description; reveal is how far it shows.
	Content func(k int, reveal float32)
}

// ToastStack lays toasts, oldest first, out as a stack in place: the newest
// in front, and those before it narrower and tucked behind, Peek px of each
// showing above it, as an outline. With the pointer over the stack they
// spring apart into a list at its bottom, Gap apart, and a toast slides in
// from above as it first shows. A toast's ID names it, its index when
// empty.
func ToastStack(c *ui.Context, toasts []ui.Toast, o ToastStackOptions) ui.Element {
	t := theme.Of(c)
	if o.Width == 0 {
		o.Width, o.Height, o.Corner, o.Background, o.Border, o.Ring = 360, 66, t.RadiusLG, t.Popover, t.Border, t.Border
	}
	if o.Slots == 0 {
		o.Slots = 3
	}
	if o.Peek == 0 {
		o.Peek = 14
	}
	if o.Gap == 0 {
		o.Gap = 10
	}
	content := o.Content
	if content == nil {
		content = func(k int, reveal float32) {
			ui.Column(c).WidthPercent(100).HeightPercent(100).Justify(ui.Center).TextColor(t.PopoverForeground).Opacity(reveal).
				Children(func() { toastBody(c, toasts[k], nil) })
		}
	}
	n, slots := float32(len(toasts)), float32(o.Slots)
	h := slots*o.Height + (slots-1)*o.Gap
	root := ui.Box(c).Size(o.Width, h)
	fan := root.Spring("fan", b2f(root.Hovered()), ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.25})
	restFront := h - o.Height - o.Rest
	root.Children(func() {
		for k, ts := range toasts {
			depth := n - 1 - float32(k)
			pos := slots - n + float32(k) // its place fanned out, the list kept at the bottom
			rest := restFront - depth*o.Peek
			y := rest + (pos*(o.Height+o.Gap)-rest)*fan
			shrink := (1 - fan) * depth * 7
			reveal := float32(1)
			if depth > 0 {
				reveal = min(max(fan, 0), 1)
			}
			var key any = k
			if ts.ID != "" {
				key = ts.ID
			}
			box := ui.Box(c).Key(key).Absolute().Left(shrink).Top(y).Size(o.Width-2*shrink, o.Height)
			born := ui.Local(box, "born", func() bool { return false })
			enter := box.Spring("enter", b2f(*born), ui.SpringOptions{Duration: 360 * time.Millisecond, Bounce: 0.2})
			if !*born {
				*born = true
				c.AnimationFrame()
			}
			box.Translate(0, (1-enter)*-16).Opacity(min(max(enter, 0), 1))
			box.Draw(func(p *ui.Painter, r ui.Rect) {
				p.Fill(r, o.Background.Alpha(reveal), o.Corner)
				if o.Ring.A != 0 {
					p.Stroke(r, o.Ring.Alpha(reveal), o.Corner, 1)
				}
				if reveal < 1 {
					p.Stroke(r, o.Border.Alpha(1-reveal), o.Corner, 1)
				}
			})
			if reveal < 0.05 {
				continue
			}
			box.Children(func() { content(k, reveal) })
		}
	})
	return root
}

// AnchoredToast shows ts above anchor, sideOffset apart: as a tooltip of
// its title with tooltip, else as a toast card at most 256 wide. Call it
// while the toast shows.
func AnchoredToast(c *ui.Context, anchor ui.Element, ts ui.Toast, sideOffset float32, tooltip bool) {
	t := theme.Of(c)
	ui.Overlay(c, func() {
		p := ui.Column(c).AttachTo(anchor, ui.AnchorTop, ui.AnchorBottom).Margin(sideOffset, 0, 0, 0).
			Border(1, t.Border).Background(t.Popover).TextColor(t.PopoverForeground)
		if tooltip {
			p.Radius(t.RadiusMD).Padding(4, 8).Shadow(0, 4, 6, -1, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).Children(func() {
				ui.Text(c, ts.Title).SingleLine().FontSize(12).FixedLineHeight(16)
			})
			return
		}
		p.MaxWidth(256).Radius(t.RadiusLG).Shadow(0, 10, 15, -3, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).Children(func() {
			toastBody(c, ts, nil)
		})
	})
}

// toastBody is a toast card's content: its Type's icon, title and
// description, and the action button action builds.
func toastBody(c *ui.Context, ts ui.Toast, action func()) {
	t := theme.Of(c)
	ui.Row(c).Padding(12, 14).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(6).Children(func() {
		ui.Row(c).Gap(8).MinWidth(0).Children(func() {
			// The icon sits in a box a line high, centred on the title.
			icon := func(svg *ui.SVG, color ui.Color) {
				ui.Row(c).Height(20).AlignItems(ui.Center).Shrink(0).Children(func() {
					ui.Icon(c, svg).Size(16, 16).TextColor(color)
				})
			}
			switch ts.Type {
			case "success":
				icon(icons.CircleCheck(), t.Success)
			case "error":
				icon(icons.CircleAlert(), t.Destructive)
			case "info":
				icon(icons.Info(), t.Info)
			case "warning":
				icon(icons.TriangleAlert(), t.Warning)
			case "loading":
				ui.Row(c).Height(20).AlignItems(ui.Center).Shrink(0).Opacity(0.8).Children(func() {
					Spinner(c).Size(16)
				})
			}
			ui.Column(c).MinWidth(0).Gap(2).Children(func() {
				ui.Text(c, ts.Title).FontSize(14).FixedLineHeight(20).FontWeight(500)
				if ts.Description != "" {
					ui.Text(c, ts.Description).FontSize(14).FixedLineHeight(20).TextColor(t.MutedForeground)
				}
			})
		})
		if action != nil {
			action()
		}
	})
}
