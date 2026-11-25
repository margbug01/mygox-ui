package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// SeekOptions configures Seek.
type SeekOptions struct {
	Corner   float32 // 0 to 24px; 24 by default
	Give     float32 // how far it leans toward the cursor, 0 to 100; 50 by default
	Spring   float32 // how the width settles, 0 to 100; 20 by default
	Width    float32 // the open field, 200 to 360px; 240 by default
	OnSearch func(query string)
}

const (
	seekH          = 48
	seekBlockW     = 266
	seekBlockH     = 88
	seekCloseGrace = 300 * time.Millisecond
)

type seekState struct {
	open     bool
	openedAt time.Time
	query    string
}

// Seek is a search bar that starts as a single icon. Press the icon and it
// opens into a search field as one object, the width settling on a spring;
// Escape, or leaving it empty, folds it back.
func Seek(c *ui.Context, opts ...SeekOptions) ui.Element {
	var o SeekOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 24
	}
	if o.Give == 0 {
		o.Give = 50
	}
	if o.Spring == 0 {
		o.Spring = 20
	}
	if o.Width == 0 {
		o.Width = 240
	}
	o.Corner, o.Give, o.Spring, o.Width = clamp32(o.Corner, 0, 24), clamp32(o.Give, 0, 100), clamp32(o.Spring, 0, 100), clamp32(o.Width, 200, 360)
	t := ThemeFrom(c)

	holder := ui.Box(c).Size(max(seekBlockW, o.Width+26), seekBlockH)
	st := ui.Local(holder, "seek", func() seekState { return seekState{} })
	now := c.Now()
	w := holder.Spring("seek_w", seekH+(o.Width-seekH)*b2f(st.open), ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: float64(o.Spring) / 100 * 0.7})
	hw := max(seekBlockW, o.Width+26)

	holder.Children(func() {
		bar := ui.Box(c).Absolute().Left(hw/2-w/2).Top(seekBlockH/2-seekH/2).Size(w, seekH).
			Radius(min(o.Corner, seekH/2)).Background(t.Surface).Clip()
		px, py, over := bar.PointerPosition()
		lean := func(v, size float32) float32 { // toward the cursor, a few px at most
			if !over || st.open {
				return 0
			}
			return clamp32((v-size/2)/(size/2), -1, 1) * 4 * o.Give / 50
		}
		lx := bar.Spring("lean_x", lean(px, w), ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.3})
		ly := bar.Spring("lean_y", lean(py, seekH), ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.3})
		hover := bar.Spring("hover", b2f(over && !st.open), ui.SpringOptions{Duration: 200 * time.Millisecond})
		press := bar.Spring("press", b2f(bar.Pressed() && !st.open), ui.SpringOptions{Duration: 160 * time.Millisecond, Bounce: 0.3})
		s := 1 + 0.04*hover - 0.1*press
		bar.Translate(lx, ly).Scale(s, s)
		if !st.open {
			bar.Cursor(ui.CursorPointer).Label("Search").OnClick(func() {
				st.open, st.openedAt = true, now
			})
		}
		bar.Children(func() {
			iconCol := lerpColor(t.TextMuted, t.Text, hover)
			ui.Icon(c, icons.Search()).Size(18, 18).TextColor(iconCol).Absolute().Left(15).Top(seekH/2 - 9)
			if !st.open && w < seekH+8 {
				return
			}
			in := ui.TextInput(c, &st.query).Placeholder("Search").PlaceholderColor(t.TextSubtle).
				FontSize(15).TextColor(t.Text).Background(ui.Color{}).Border(0, ui.Color{}).Padding(0).
				Absolute().Left(45).Top(seekH/2 - 11).Width(max(o.Width-45-18, 0)).Height(22).
				Opacity(clamp32((w-seekH)/(o.Width-seekH)*1.6-0.3, 0, 1))
			if !st.open {
				return
			}
			if now.Sub(st.openedAt) < 50*time.Millisecond {
				in.Focus()
			}
			switch {
			case in.Shortcut(0, ui.KeyEscape):
				st.open, st.query = false, ""
			case in.Submitted() && o.OnSearch != nil:
				o.OnSearch(st.query)
			case !in.Focused() && st.query == "" && now.Sub(st.openedAt) > seekCloseGrace:
				st.open = false
			}
		})
	})
	if st.open && now.Sub(st.openedAt) <= seekCloseGrace {
		c.AnimationFrame()
	}
	return holder
}

func init() {
	registerDemo(
		"seek",
		"A search bar that starts as a single icon: press it and it opens into a search field as one object, its width settling on a spring.",
		func(c *ui.Context) { Seek(c) },
	)
}
