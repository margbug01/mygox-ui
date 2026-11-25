package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// IconBarAxis is the way IconBar runs.
type IconBarAxis string

const (
	IconBarRow    IconBarAxis = "row"
	IconBarColumn IconBarAxis = "column"
)

// IconBarItem is one tab of IconBar.
type IconBarItem struct {
	Label string
	Icon  *ui.SVG
}

// IconBarOptions configures IconBar.
type IconBarOptions struct {
	Corner float32     // 0 to 26; 26 by default
	Dilate float32     // how far the pill stretches in transit, 0 to 100; 100 by default
	Bounce float32     // how hard it lands, 0 to 100; 50 by default
	Speed  float32     // 0 to 100; 50 by default
	Axis   IconBarAxis // row by default
}

// IconBarDefaultItems are the tabs bencho.dev shows.
func IconBarDefaultItems() []IconBarItem {
	return []IconBarItem{
		{"Home", icons.House()}, {"Search", icons.Search()}, {"Files", icons.Folder()},
		{"Saved", icons.Bookmark()}, {"Profile", icons.UserRound()},
	}
}

const (
	iconBarSlot  = 40
	iconBarPad   = 5
	iconBarThick = 50
	iconBarPill  = 38
)

type iconBarState struct {
	forward bool // the last move went right (or down)
}

// IconBar is a tab bar of icons. Picking a tab stretches the selection pill
// across to it, its leading edge first, and it lands with a small bounce.
func IconBar(c *ui.Context, items []IconBarItem, selected *int, opts ...IconBarOptions) ui.Element {
	var o IconBarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 26
	}
	o.Corner = clamp32(o.Corner, 0, 26)
	if o.Dilate == 0 {
		o.Dilate = 100
	}
	if o.Bounce == 0 {
		o.Bounce = 50
	}
	if o.Speed == 0 {
		o.Speed = 50
	}
	if o.Axis == "" {
		o.Axis = IconBarRow
	}
	t := ThemeFrom(c)
	*selected = min(max(*selected, 0), len(items)-1)
	column := o.Axis == IconBarColumn

	long := float32(len(items))*iconBarSlot + 2*iconBarPad
	bar := ui.Box(c).Radius(min(o.Corner, iconBarThick/2)).Background(t.Surface)
	if column {
		bar.Size(iconBarThick, long)
	} else {
		bar.Size(long, iconBarThick)
	}
	st := ui.Local(bar, "icon_bar", func() iconBarState { return iconBarState{} })

	// The pill's two ends run on their own springs: the leading one quick,
	// the trailing one slower by as much as Dilate says, so it stretches.
	dur := lerp32(700, 260, clamp32(o.Speed, 0, 100)/100)
	lag := 1 + clamp32(o.Dilate, 0, 100)/100*0.9
	bounce := float64(clamp32(o.Bounce, 0, 100)) / 100 * 0.5
	lead := ui.SpringOptions{Duration: time.Duration(dur) * time.Millisecond, Bounce: bounce}
	trail := ui.SpringOptions{Duration: time.Duration(dur*lag) * time.Millisecond, Bounce: bounce * 0.6}
	start := iconBarPad + float32(*selected)*iconBarSlot + (iconBarSlot-iconBarPill)/2
	headOpts, tailOpts := lead, trail // head is the far end along the axis
	if !st.forward {
		headOpts, tailOpts = trail, lead
	}
	near := bar.Spring("near", start, tailOpts)
	far := bar.Spring("far", start+iconBarPill, headOpts)
	if far-near < iconBarPill*0.6 { // keep a pill shape through any overshoot
		mid := (near + far) / 2
		near, far = mid-iconBarPill*0.3, mid+iconBarPill*0.3
	}

	bar.Children(func() {
		pill := ui.Box(c).Absolute().Radius(min(o.Corner, iconBarPill/2)).Background(t.SurfaceSubtle)
		across := float32(iconBarThick-iconBarPill) / 2
		if column {
			pill.Left(across).Top(near).Size(iconBarPill, far-near)
		} else {
			pill.Left(near).Top(across).Size(far-near, iconBarPill)
		}
		for i, it := range items {
			on := i == *selected
			at := iconBarPad + float32(i)*iconBarSlot
			tab := ui.Box(c).Key(i).Absolute().Center().Cursor(ui.CursorPointer).
				Role(ui.RoleTab).Label(it.Label).AccessSelected(on).
				OnClick(func() {
					st.forward = i > *selected
					*selected = i
				})
			if column {
				tab.Left(0).Top(at).Size(iconBarThick, iconBarSlot)
			} else {
				tab.Left(at).Top(0).Size(iconBarSlot, iconBarThick)
			}
			lit := tab.Spring("lit", b2f(on), ui.SpringOptions{Duration: 220 * time.Millisecond})
			hover := tab.Spring("hover", b2f(tab.Hovered()), ui.SpringOptions{Duration: 150 * time.Millisecond})
			press := tab.Spring("press", b2f(tab.Pressed()), ui.SpringOptions{Duration: 120 * time.Millisecond})
			ink := lerpColor(lerpColor(lerpColor(t.TextSubtle, t.SurfaceHover, 0.35), t.TextMuted, hover), t.Text, lit)
			tab.Children(func() {
				ui.Icon(c, it.Icon).Size(20, 20).TextColor(ink).Scale(1-0.1*press, 1-0.1*press)
			})
		}
	})
	return bar
}

func init() {
	registerDemo(
		"icon-bar",
		"A tab bar of icons whose selection pill stretches across to the tab picked and lands with a small bounce.",
		func(c *ui.Context) {
			holder := ui.Box(c)
			sel := ui.Local(holder, "demo", func() int { return 0 })
			holder.Children(func() { IconBar(c, IconBarDefaultItems(), sel) })
		},
	)
}
