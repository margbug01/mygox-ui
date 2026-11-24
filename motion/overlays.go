package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// ============================================================================
// 23. Accordion
// ============================================================================

// AccordionItem defines a single collapsible entry in Accordion.
type AccordionItem struct {
	Value   string
	Title   string
	Content func()
}

// AccordionOptions configures Accordion.
type AccordionOptions struct {
	AllowCollapse bool
	MaxItemHeight float32
	Spring        ui.SpringOptions
}

// Accordion renders a vertical stack of spring-animated expandable items
// separated by crisp zinc hairline dividers.
func Accordion(c *ui.Context, items []AccordionItem, expanded *string, opts ...AccordionOptions) ui.Element {
	var o AccordionOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.MaxItemHeight <= 0 {
		o.MaxItemHeight = 140
	}
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 280, Damping: 26, Mass: 0.7})

	current := ""
	if expanded != nil {
		current = *expanded
	}

	col := ui.Column(c).WidthPercent(100)
	col.Children(func() {
		for i, item := range items {
			it := item
			isOpen := current == it.Value
			rowBox := ui.Column(c).
				Key(it.Value).
				WidthPercent(100)
			if i < len(items)-1 {
				rowBox.BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
			}

			p := clamp32(rowBox.Spring("acc-p", boolFloat(isOpen), spring), 0, 1)

			rowBox.Children(func() {
				hdr := ui.Row(c).
					WidthPercent(100).
					Padding(10, 2).
					AlignItems(ui.Center).
					Justify(ui.SpaceBetween).
					Cursor(ui.CursorPointer)
				hdr.Children(func() {
					ui.Text(c, it.Title).FontSize(13).FontWeight(500).TextColor(t.Text).SingleLine()
					ui.Box(c).
						Width(18).
						Height(18).
						Center().
						Rotate(p * 180).
						Children(func() {
							ui.Icon(c, icons.ChevronDown()).Size(16, 16).TextColor(t.TextMuted)
						})
				})
				if hdr.Clicked() && expanded != nil {
					if isOpen {
						if o.AllowCollapse {
							*expanded = ""
						}
					} else {
						*expanded = it.Value
					}
				}

				if p > 0.01 && it.Content != nil {
					ui.Box(c).
						WidthPercent(100).
						MaxHeight(p * o.MaxItemHeight).
						Clip().
						Opacity(p).
						Translate(0, (1-p)*-4).
						Children(func() {
							ui.Box(c).WidthPercent(100).Padding(0, 2, 10, 2).Children(it.Content)
						})
				}
			})
		}
	})
	return col
}

// ============================================================================
// 24. Disclosure
// ============================================================================

// DisclosureOptions configures Disclosure.
type DisclosureOptions struct {
	MaxHeight float32
	Radius    float32
	Spring    ui.SpringOptions
}

// Disclosure renders a single bordered card that smoothly expands and collapses
// its content panel with spring physics.
func Disclosure(c *ui.Context, triggerLabel string, open *bool, content func(), opts ...DisclosureOptions) ui.Element {
	var o DisclosureOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.MaxHeight <= 0 {
		o.MaxHeight = 160
	}
	if o.Radius <= 0 {
		o.Radius = 10
	}
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 280, Damping: 26, Mass: 0.7})

	isOpen := open != nil && *open
	card := ui.Column(c).
		WidthPercent(100).
		Radius(o.Radius).
		Background(t.Surface).
		Border(1, t.Border).
		Padding(0, 14).
		Clip()

	p := clamp32(card.Spring("disc-p", boolFloat(isOpen), spring), 0, 1)

	card.Children(func() {
		trig := ui.Row(c).
			WidthPercent(100).
			Padding(10, 0).
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween).
			Cursor(ui.CursorPointer)
		trig.Children(func() {
			ui.Text(c, triggerLabel).FontSize(13).FontWeight(500).TextColor(t.Text).SingleLine()
			ui.Box(c).Width(18).Height(18).Center().Rotate(p * 180).Children(func() {
				ui.Icon(c, icons.ChevronDown()).Size(16, 16).TextColor(t.TextMuted)
			})
		})
		if trig.Clicked() && open != nil {
			*open = !*open
		}

		if p > 0.01 && content != nil {
			ui.Box(c).
				WidthPercent(100).
				MaxHeight(p * o.MaxHeight).
				Clip().
				Opacity(p).
				Translate(0, (1-p)*-4).
				Children(func() {
					ui.Box(c).
						WidthPercent(100).
						Padding(8, 0, 12, 0).
						BorderWidth(1, 0, 0, 0).
						BorderColor(t.Border).
						Children(content)
				})
		}
	})
	return card
}

// ============================================================================
// 25. Dialog
// ============================================================================

// DialogOptions configures Dialog.
type DialogOptions struct {
	Title       string
	Description string
	MaxWidth    float32
	Spring      ui.SpringOptions
}

// Dialog renders a spring scale+fade modal overlay via ui.DialogBase with a
// restrained monochrome zinc panel and backdrop.
func Dialog(c *ui.Context, open *bool, body func(), opts ...DialogOptions) ui.Element {
	var o DialogOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.MaxWidth <= 0 {
		o.MaxWidth = 400
	}
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 340, Damping: 26, Mass: 0.7})

	return ui.DialogBase(c, open, func(backdrop, panel ui.Element) {
		p := clamp32(panel.Spring("dlg-p", 1, spring), 0, 1)
		if p == 0 {
			// Ensure initial frame starts from subtle scale-in
			panel.Jump("dlg-p", 0.01)
			p = panel.Spring("dlg-p", 1, spring)
		}

		backdrop.Background(ui.RGBA(9, 9, 11, 0.52*p))

		sc := lerp32(0.94, 1.0, p)
		panel.
			Width(o.MaxWidth).
			MaxWidthPercent(92).
			Radius(14).
			Background(t.Surface).
			Border(1, t.Border).
			Padding(20).
			Shadow(0, 20, 48, -12, ui.RGBA(0, 0, 0, 0.35)).
			Opacity(p).
			Scale(sc, sc).
			Translate(0, (1-p)*10).
			Children(func() {
				ui.Column(c).WidthPercent(100).Gap(12).Children(func() {
					ui.Row(c).WidthPercent(100).AlignItems(ui.Start).Justify(ui.SpaceBetween).Gap(12).Children(func() {
						ui.Column(c).Grow(1).MinWidth(0).Gap(4).Children(func() {
							if o.Title != "" {
								ui.Text(c, o.Title).FontSize(15).FontWeight(600).TextColor(t.Text)
							}
							if o.Description != "" {
								ui.Text(c, o.Description).FontSize(12.5).TextColor(t.TextMuted).LineHeight(1.45)
							}
						})
						closeBtn := ui.Box(c).
							Width(26).
							Height(26).
							Radius(6).
							Center().
							Cursor(ui.CursorPointer)
						if closeBtn.Hovered() {
							closeBtn.Background(t.SurfaceHover)
						}
						closeBtn.Children(func() {
							ui.Text(c, "×").FontSize(16).FontWeight(500).TextColor(t.TextMuted).SingleLine()
						})
						if closeBtn.Clicked() && open != nil {
							*open = false
						}
					})
					if body != nil {
						body()
					}
				})
			})
	})
}

// ============================================================================
// 26. MorphingDialog
// ============================================================================

// MorphingDialogOptions configures MorphingDialog.
type MorphingDialogOptions struct {
	ModalWidth float32
	Radius     float32
	Spring     ui.SpringOptions
}

// MorphingDialog renders a clickable preview card that smoothly morphs into an
// expanded centered detail modal on click.
func MorphingDialog(c *ui.Context, open *bool, triggerCard func(), expandedContent func(), opts ...MorphingDialogOptions) ui.Element {
	var o MorphingDialogOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.ModalWidth <= 0 {
		o.ModalWidth = 440
	}
	if o.Radius <= 0 {
		o.Radius = 12
	}
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 280, Damping: 26, Mass: 0.75})

	isOpen := open != nil && *open
	trigger := ui.Box(c).
		Radius(o.Radius).
		Background(t.Surface).
		Border(1, t.Border).
		Clip().
		Cursor(ui.CursorPointer)
	if trigger.Hovered() {
		trigger.Border(1, t.BorderStrong)
	}
	if isOpen {
		trigger.Opacity(0.35)
	}
	trigger.Children(func() {
		if triggerCard != nil {
			triggerCard()
		}
	})
	if trigger.Clicked() && open != nil {
		*open = true
	}

	ui.DialogBase(c, open, func(backdrop, panel ui.Element) {
		st := ui.Local(panel, "md-init", func() bool { return false })
		if !*st {
			*st = true
			panel.Jump("md-p", 0)
		}
		p := clamp32(panel.Spring("md-p", 1, spring), 0, 1)

		backdrop.Background(ui.RGBA(9, 9, 11, 0.56*p))
		sc := lerp32(0.84, 1.0, p)
		panel.
			Width(o.ModalWidth).
			MaxWidthPercent(92).
			Radius(o.Radius + 4).
			Background(t.Surface).
			Border(1, t.Border).
			Clip().
			Shadow(0, 24, 60, -12, ui.RGBA(0, 0, 0, 0.4)).
			Opacity(p).
			Scale(sc, sc).
			Translate(0, (1-p)*22).
			Children(func() {
				ui.Box(c).WidthPercent(100).Children(func() {
					if expandedContent != nil {
						expandedContent()
					}
					closeBtn := ui.Box(c).
						Absolute().
						Top(12).
						Right(12).
						Width(28).
						Height(28).
						Radius(999).
						Background(t.SurfaceSubtle).
						Border(1, t.Border).
						Center().
						Cursor(ui.CursorPointer)
					if closeBtn.Hovered() {
						closeBtn.Background(t.SurfaceHover)
					}
					closeBtn.Children(func() {
						ui.Text(c, "×").FontSize(15).FontWeight(600).TextColor(t.Text).SingleLine()
					})
					if closeBtn.Clicked() && open != nil {
						*open = false
					}
				})
			})
	})
	return trigger
}

// ============================================================================
// 27. MorphingPopover
// ============================================================================

// MorphingPopoverOptions configures MorphingPopover.
type MorphingPopoverOptions struct {
	Width  float32
	Radius float32
	Spring ui.SpringOptions
}

// MorphingPopover renders a compact trigger button that smoothly morphs/expands
// into a floating popover form anchored to the trigger.
func MorphingPopover(c *ui.Context, open *bool, triggerLabel string, content func(), opts ...MorphingPopoverOptions) ui.Element {
	var o MorphingPopoverOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 280
	}
	if o.Radius <= 0 {
		o.Radius = 10
	}
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 340, Damping: 28, Mass: 0.65})

	trigger := ui.Box(c).
		Padding(7, 14).
		Radius(8).
		Background(t.Surface).
		Border(1, t.Border).
		Cursor(ui.CursorPointer)
	if trigger.Hovered() {
		trigger.Background(t.SurfaceSubtle).Border(1, t.BorderStrong)
	}
	trigger.Children(func() {
		ui.Text(c, triggerLabel).FontSize(12.5).FontWeight(500).TextColor(t.Text).SingleLine()
	})
	if trigger.Clicked() && open != nil {
		*open = !*open
	}

	ui.PopoverBase(c, trigger, open, func(panel ui.Element) {
		st := ui.Local(panel, "mp-init", func() bool { return false })
		if !*st {
			*st = true
			panel.Jump("mp-p", 0)
		}
		p := clamp32(panel.Spring("mp-p", 1, spring), 0, 1)
		sc := lerp32(0.90, 1.0, p)

		panel.
			Width(o.Width).
			Radius(o.Radius).
			Background(t.Surface).
			Border(1, t.Border).
			Padding(14).
			Shadow(0, 12, 32, -8, ui.RGBA(0, 0, 0, 0.24)).
			Opacity(p).
			Scale(sc, sc).
			Translate(0, (1-p)*-6).
			Children(func() {
				if content != nil {
					content()
				}
			})
	})
	return trigger
}

// ============================================================================
// 28. TransitionPanel
// ============================================================================

// TransitionPanelOptions configures TransitionPanel.
type TransitionPanelOptions struct {
	Duration time.Duration
	Blur     bool
}

type transPanelState struct {
	active  int
	changed time.Time
	hasInit bool
}

// TransitionPanel smoothly crossfades and slides (with optional blur) between
// indexed content panels whenever activeIndex changes.
func TransitionPanel(c *ui.Context, activeIndex int, panels []func(), opts ...TransitionPanelOptions) ui.Element {
	var o TransitionPanelOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Duration <= 0 {
		o.Duration = 260 * time.Millisecond
	}

	host := ui.Box(c).WidthPercent(100)
	if len(panels) == 0 {
		return host
	}
	if activeIndex < 0 {
		activeIndex = 0
	}
	if activeIndex >= len(panels) {
		activeIndex = len(panels) - 1
	}

	now := c.Now()
	st := ui.Local(host, "tp-state", func() transPanelState {
		return transPanelState{active: activeIndex, changed: now.Add(-time.Second), hasInit: true}
	})
	if st.active != activeIndex {
		st.active = activeIndex
		st.changed = now
	}

	elapsed := now.Sub(st.changed)
	p := float32(1)
	if elapsed < o.Duration {
		c.AnimationFrame()
		p = float32(easeOutCubic(float64(elapsed) / float64(o.Duration)))
	}

	host.Children(func() {
		pane := ui.Box(c).
			Key(activeIndex).
			WidthPercent(100).
			Opacity(p).
			Translate(0, (1-p)*8)
		if o.Blur && p < 0.98 {
			pane.Blur((1 - p) * 3.5)
		}
		pane.Children(func() {
			if panels[activeIndex] != nil {
				panels[activeIndex]()
			}
		})
	})
	return host
}
