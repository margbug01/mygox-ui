package coss

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// SelectOption represents one selectable choice in a Select dropdown.
type SelectOption[T comparable] struct {
	Value    T
	Label    string
	Icon     *ui.SVG // drawn before the label by a folding select
	Disabled bool
}

// SelectOptions are a Select's placeholder and its optional motion.
type SelectOptions struct {
	Placeholder string
	// Fold unfolds the options out of the trigger: the choice stays where
	// it is while the others spring out below it one after another, and a
	// pick slides up from its row into the trigger as the rest fold away.
	Fold bool
	// Bounce is how far the folding card springs past its size as it
	// opens, 0 to 100; 30 when 0.
	Bounce float32
	// Style draws a folding select in another design; nil is coss's look.
	Style *SelectStyle
}

// SelectStyle is the look of a folding select: its trigger and every row
// are Height tall, and the card it opens into is the trigger grown.
type SelectStyle struct {
	Width              float32 // 0: as wide as its parent
	Height, PaddingX   float32
	FontSize           float32
	Corner, OpenCorner float32 // shut, and open
	Background, Border ui.Color
	Value, Option      ui.Color // the choice's label, and the others'
	Icon, Hover        ui.Color
}

// SelectSize selects the trigger height and padding, as ui/select.tsx's size.
type SelectSize uint8

const (
	SelectDefault SelectSize = iota
	SelectSM
	SelectLG
)

// SelectElement is the trigger of a Select.
type SelectElement struct {
	ui.Element
}

// Size sets the trigger's height and padding.
func (s SelectElement) Size(size SelectSize) SelectElement {
	sizeSelectTrigger(s.Element, size)
	return s
}

// Select renders a coss-styled dropdown selector built on ui.SelectBase,
// showing each option's Label. Compose SelectTrigger, SelectPopup and
// SelectItem for options with richer content.
func Select[T comparable](c *ui.Context, selected *T, options []SelectOption[T], opts ...SelectOptions) SelectElement {
	var o SelectOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Fold {
		return SelectElement{foldSelect(c, selected, options, o)}
	}
	parts := ui.SelectBase(c, selected)
	label := ""
	for _, opt := range options {
		if opt.Value == *selected {
			label = opt.Label
			break
		}
	}
	trigger := SelectTrigger(c, parts.Trigger, SelectDefault, func() { SelectValue(c, label, o.Placeholder) })
	SelectPopup(c, parts, func() {
		for _, opt := range options {
			SelectItem(c, parts, opt.Value, opt.Value == *selected, opt.Disabled, func() {
				ui.Text(c, opt.Label).SingleLine()
			})
		}
	})
	return SelectElement{trigger}
}

type foldSelectState[T comparable] struct {
	want    bool // open, as the user last asked
	shown   bool // the popup's flag, which the select sets as it opens and closes
	set     bool // shown as handed to the select, to tell its changes
	showing T    // the choice the open card was built around
	slideAt time.Time
	slideY  float32 // how far below the trigger the pick started
}

const foldSelectPad = 8 // below the open card's last row

// foldSelect is a Select with Fold: shut, its trigger draws the card; open,
// and until it has folded away, the card is the popup, laid over the
// trigger.
func foldSelect[T comparable](c *ui.Context, selected *T, options []SelectOption[T], o SelectOptions) ui.Element {
	t := theme.Of(c)
	s := SelectStyle{
		Height: 32, PaddingX: 11, FontSize: 14, Corner: t.RadiusLG, OpenCorner: t.RadiusLG,
		Background: t.Popover, Border: t.Input, Value: t.Foreground, Option: t.MutedForeground,
		Icon: t.MutedForeground, Hover: t.Accent,
	}
	if o.Style != nil {
		s = *o.Style
	}
	if o.Bounce == 0 {
		o.Bounce = 30
	}
	holder := ui.Box(c).Height(s.Height)
	if s.Width > 0 {
		holder.Width(s.Width)
	} else {
		holder.WidthPercent(100)
	}
	w := s.Width
	if w == 0 {
		w = holder.Bounds().W
	}
	st := ui.Local(holder, "fold", func() foldSelectState[T] { return foldSelectState[T]{showing: *selected} })

	// The select opens the popup from its trigger and closes it on a pick,
	// Escape or a press outside; the card follows, and stays up while it
	// folds away.
	if st.shown != st.set {
		st.want = st.shown
	}
	if *selected != st.showing {
		st.slideY = s.Height * float32(foldSelectRank(options, st.showing, *selected))
		st.slideAt = c.Now()
		st.showing = *selected
	}
	open := holder.Spring("open", b2f(st.want), ui.SpringOptions{
		Duration: 440 * time.Millisecond,
		Bounce:   float64(min(max(o.Bounce, 0), 100)) / 100 * 0.5 * float64(b2f(st.want)),
	})
	slide := float32(0)
	if !st.slideAt.IsZero() {
		k := min(max(float32(c.Now().Sub(st.slideAt))/float32(380*time.Millisecond), 0), 1)
		e := 1 - (1-k)*(1-k)*(1-k)
		slide = st.slideY * (1 - e)
		if k < 1 {
			c.AnimationFrame()
		} else {
			st.slideAt = time.Time{}
		}
	}
	st.shown = st.want || open != 0 || slide != 0
	st.set = st.shown

	cur := -1
	for i, opt := range options {
		if opt.Value == *selected {
			cur = i
			break
		}
	}
	others := float32(len(options))
	if cur >= 0 {
		others--
	}
	row := func(r ui.Element, label string, icon *ui.SVG, y float32, col ui.Color, weight int) ui.Element {
		r.Absolute().Left(0).Top(y).Size(w, s.Height).Padding(0, s.PaddingX).AlignItems(ui.Center)
		r.Children(func() {
			if icon != nil {
				ui.Icon(c, icon).Size(16, 16).TextColor(s.Icon).Margin(0, 12, 0, 0)
			}
			ui.Text(c, label).FontSize(s.FontSize).FontWeight(weight).TextColor(col).SingleLine()
		})
		return r
	}
	// card styles box as the card open by open, with the choice slide below
	// its place; item makes the rows of the options, keyed in kc.
	card := func(box ui.Element, open, slide float32, item func(kc *ui.Context, v T) ui.Element) {
		box.Size(w, s.Height+max(open, 0)*(others*s.Height+foldSelectPad)).
			Radius(s.Corner + (s.OpenCorner-s.Corner)*min(max(open, 0), 1)).
			Background(s.Background).Clip()
		if s.Border.A != 0 {
			box.Border(1, s.Border)
		}
		if open > 0.02 {
			box.Shadow(0, 8, 24, 0, ui.RGBA(0, 0, 0, 0.08*min(max(open, 0), 1)))
		}
		box.Children(func() {
			if cur >= 0 {
				row(item(c.Key("value"), *selected), options[cur].Label, options[cur].Icon, slide, s.Value, 500).Cursor(ui.CursorPointer)
			} else {
				row(ui.Row(c), o.Placeholder, nil, 0, s.Option, 500)
			}
			ui.Icon(c, icons.ChevronDown()).Absolute().Left(w-s.PaddingX-16).Top((s.Height-16)/2).Size(16, 16).
				TextColor(s.Icon).Rotate(180 * min(max(open, 0), 1))
			j := 0
			for i, opt := range options {
				if i == cur {
					continue
				}
				pos := float32(j + 1)
				j++
				k := min(max(open*(others+1)-(pos-1)*0.6, 0), 1)
				if k <= 0 {
					continue
				}
				r := row(item(c.Key(i), opt.Value), opt.Label, opt.Icon, s.Height*pos-(1-k)*s.Height*0.5, s.Option, 400).
					Opacity(k).Cursor(ui.CursorPointer)
				if opt.Disabled {
					r.Disabled(true).DimDisabled(false).Opacity(0.64 * k)
					continue
				}
				// The keys' highlight shows while open; folding away, it
				// stays on its index, a row the pick has moved under it.
				hover := r.Spring("hover", b2f(r.Hovered() || st.want && r.Highlighted()), ui.SpringOptions{Duration: 140 * time.Millisecond})
				r.Background(s.Hover.Alpha(hover))
			}
		})
	}

	var trigger ui.Element
	holder.Children(func() {
		parts := ui.SelectBaseOpen(c, selected, &st.shown)
		trigger = parts.Trigger.WidthPercent(100).HeightPercent(100).Cursor(ui.CursorPointer).DimDisabled(false)
		if !st.shown {
			trigger.Children(func() {
				card(ui.Box(c).Absolute().Left(0).Top(0), 0, 0, func(kc *ui.Context, _ T) ui.Element { return ui.Row(kc) })
			})
		}
		parts.Popup(func(panel ui.Element) {
			card(panel.AttachTo(trigger, ui.AnchorTopStart, ui.AnchorTopStart), open, slide, func(_ *ui.Context, v T) ui.Element { return parts.Item(v) })
		})
	})
	return trigger
}

// foldSelectRank is the row, counting the trigger's as 0, that picked had
// in the card open around shown.
func foldSelectRank[T comparable](options []SelectOption[T], shown, picked T) int {
	rank := 0
	for _, opt := range options {
		if opt.Value == shown {
			continue
		}
		rank++
		if opt.Value == picked {
			return rank
		}
	}
	return 0
}

func b2f(b bool) float32 {
	if b {
		return 1
	}
	return 0
}

// SelectTrigger styles a ui.SelectBase trigger and fills it with value,
// then the chevrons-up-down icon. A disabled trigger shows at 0.64 opacity
// once the caller sets Opacity(0.64): mygo's own dimming is off.
func SelectTrigger(c *ui.Context, trigger ui.Element, size SelectSize, value func()) ui.Element {
	t := theme.Of(c)
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	trigger.
		WidthPercent(100).
		Radius(t.RadiusLG).
		Background(bg).
		Border(1, t.Input).
		Justify(ui.SpaceBetween).
		AlignItems(ui.Center).
		TextColor(t.Foreground).
		FontSize(14).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring).
		DimDisabled(false)
	sizeSelectTrigger(trigger, size)
	if !t.Dark {
		trigger.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04))
	}
	trigger.Children(func() {
		ui.Row(c).Grow(1).MinWidth(0).AlignItems(ui.Center).Gap(8).Children(value)
		ui.Icon(c, icons.ChevronsUpDown()).Size(16, 16).MarginEnd(-4).
			TextColor(theme.Alpha(t.Foreground, 0.8)).Shrink(0)
	})
	return trigger
}

func sizeSelectTrigger(trigger ui.Element, size SelectSize) {
	switch size {
	case SelectSM:
		trigger.MinHeight(28).Padding(0, 9).Gap(6)
	case SelectLG:
		trigger.MinHeight(36).Padding(0, 11).Gap(8)
	default:
		trigger.MinHeight(32).Padding(0, 11).Gap(8)
	}
}

// SelectValue shows the chosen label, or placeholder in the muted color
// when label is empty.
func SelectValue(c *ui.Context, label, placeholder string) ui.Element {
	t := theme.Of(c)
	if label == "" {
		return ui.Text(c, placeholder).SingleLine().Ellipsis("…").TextColor(t.MutedForeground)
	}
	return ui.Text(c, label).SingleLine().Ellipsis("…")
}

// SelectPopup builds the select's popup panel and the options fn adds with
// SelectItem, SelectGroupLabel and SelectSeparator.
func SelectPopup[T comparable](c *ui.Context, parts ui.SelectParts[T], fn func()) {
	parts.Popup(func(panel ui.Element) { stylePopupPanel(c, panel).Padding(4).Children(fn) })
}

// stylePopupPanel gives a select or combobox popup its frame: rounded-lg,
// border, bg-popover, shadow-lg/5.
func stylePopupPanel(c *ui.Context, panel ui.Element) ui.Element {
	t := theme.Of(c)
	return panel.
		Radius(t.RadiusLG).
		Background(t.Popover).
		TextColor(t.PopoverForeground).
		Border(1, t.Border).
		Shadow(0, 10, 15, -3, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
		Shadow(0, 4, 6, -4, theme.Alpha(ui.RGB(0, 0, 0), 0.05))
}

// SelectItem builds the option for value: a check column marking the
// chosen one, then fn's content. A disabled option shows at 0.64 opacity.
func SelectItem[T comparable](c *ui.Context, parts ui.SelectParts[T], value T, checked, disabled bool, fn func()) ui.Element {
	return styleOptionItem(c, parts.Item(value), checked, disabled, fn)
}

// styleOptionItem lays out a select or combobox option: a check column,
// then fn's content, highlighted with the accent.
func styleOptionItem(c *ui.Context, item ui.Element, checked, disabled bool, fn func()) ui.Element {
	t := theme.Of(c)
	item.
		WidthPercent(100).
		MinHeight(28).
		Padding(4, 16, 4, 8).
		Gap(8).
		Radius(t.RadiusSM).
		AlignItems(ui.Center).
		FontSize(14).
		Disabled(disabled).
		DimDisabled(false)
	if disabled {
		item.Opacity(0.64)
	} else if item.Highlighted() || item.Hovered() {
		item.Background(t.Accent).TextColor(t.AccentForeground)
	}
	item.Children(func() {
		ui.Box(c).Size(16, 16).Shrink(0).Children(func() {
			if checked {
				ui.Icon(c, icons.Check()).Size(16, 16)
			}
		})
		ui.Row(c).Grow(1).MinWidth(0).AlignItems(ui.Center).Children(fn)
	})
	return item
}

// SelectGroupLabel titles a group of options in a SelectPopup.
func SelectGroupLabel(c *ui.Context, text string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, text).Padding(6, 8).FontSize(12).FontWeight(500).TextColor(t.MutedForeground)
}

// SelectSeparator divides groups of options in a SelectPopup.
func SelectSeparator(c *ui.Context) ui.Element {
	t := theme.Of(c)
	return ui.Box(c).Height(1).Margin(4, 8).Background(t.Border)
}

// SelectLabel labels a select from above, as ui/select.tsx's SelectLabel
// outside a Field.
func SelectLabel(c *ui.Context, text string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, text).SingleLine().FontSize(14).FontWeight(500).TextColor(t.Foreground).Margin(0, 0, 8, 0)
}
