package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// CreateMenuItem is one thing the create menu can make.
type CreateMenuItem struct {
	Label string
	Icon  *ui.SVG
}

// CreateMenuOptions configures CreateMenu.
type CreateMenuOptions struct {
	Corner   float32 // the open menu's corner, 0 to 40; 20 by default
	OnCreate func(it CreateMenuItem)
}

type createMenuState struct {
	open    bool
	toggled time.Time
}

const (
	createMenuPillW = 102
	createMenuPillH = 38
	createMenuW     = 212
	createMenuRow   = 36
	createMenuPad   = 11
)

// Lucide's file-text, sheet and grid-2x2 (ISC licence), which coss/icons doesn't carry.
func createMenuIcon(paths string) *ui.SVG {
	return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">` + paths + `</svg>`))
}

// CreateMenuDefaultItems is the menu bencho.dev shows.
func CreateMenuDefaultItems() []CreateMenuItem {
	return []CreateMenuItem{
		{Label: "Document", Icon: createMenuIcon(`<path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"/><path d="M14 2v4a2 2 0 0 0 2 2h4"/><path d="M10 9H8"/><path d="M16 13H8"/><path d="M16 17H8"/>`)},
		{Label: "Spreadsheet", Icon: createMenuIcon(`<rect width="18" height="18" x="3" y="3" rx="2"/><path d="M3 9h18"/><path d="M9 3v18"/>`)},
		{Label: "Board", Icon: createMenuIcon(`<rect width="18" height="18" x="3" y="3" rx="2"/><path d="M3 12h18"/><path d="M12 3v18"/>`)},
		{Label: "Folder", Icon: icons.Folder()},
	}
}

// CreateMenu starts as a single "+ Create" pill. Pressing it spreads the pill
// into the full menu around the same centre, as if the menu was inside it all
// along; picking an item, or pressing the pill's spot again, pulls it back in.
func CreateMenu(c *ui.Context, items []CreateMenuItem, opts ...CreateMenuOptions) ui.Element {
	var o CreateMenuOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	t := ThemeFrom(c)
	menuH := float32(len(items))*createMenuRow + 2*createMenuPad
	holder := ui.Box(c).Size(createMenuW, menuH)
	st := ui.Local(holder, "create_menu", func() createMenuState { return createMenuState{} })

	// Width leads and height follows a beat behind, which gives the spread its
	// gooey stretch; both overshoot a little on the way open.
	open := b2f(st.open)
	kw := holder.Spring("w", open, ui.SpringOptions{Duration: 460 * time.Millisecond, Bounce: 0.28 * float64(open)})
	kh := holder.Spring("h", open, ui.SpringOptions{Duration: 560 * time.Millisecond, Bounce: 0.22 * float64(open)})
	kc := clamp32(holder.Spring("content", open, ui.SpringOptions{Duration: 300 * time.Millisecond}), 0, 1)
	w := lerp32(createMenuPillW, createMenuW, kw)
	h := lerp32(createMenuPillH, menuH, kh)
	corner := lerp32(createMenuPillH/2, min(o.Corner, menuH/2), clamp32(kh, 0, 1))
	since := c.Now().Sub(st.toggled)
	if since < 700*time.Millisecond {
		c.AnimationFrame()
	}

	holder.Children(func() {
		blob := ui.Box(c).Absolute().Left((createMenuW-w)/2).Top((menuH-h)/2).Size(w, h).Radius(corner).
			Background(t.Surface).Clip()
		press := blob.Spring("press", b2f(blob.Pressed() && !st.open), ui.SpringOptions{Duration: 140 * time.Millisecond})
		blob.Scale(1-0.04*press, 1-0.04*press)
		if !st.open {
			blob.Cursor(ui.CursorPointer).OnClick(func() {
				st.open = true
				st.toggled = c.Now()
			})
		}
		blob.Children(func() {
			// The pill's label shrinks into the middle as the menu spreads.
			if kc < 0.6 {
				s := 1 - 0.3*kc/0.6
				ui.Row(c).Key("pill").Absolute().Left((w-createMenuPillW)/2).Top((h-createMenuPillH)/2).
					Size(createMenuPillW, createMenuPillH).Center().Opacity(1-kc/0.6).Scale(s, s).Children(func() {
					ui.Icon(c, icons.Plus()).Size(14, 14).TextColor(t.Text).Margin(0, 7, 0, 0)
					ui.Text(c, "Create").FontSize(14).TextColor(t.Text).SingleLine()
				})
			}
			if kc <= 0.05 {
				return
			}
			// The menu sits still about the centre; the blob's edge uncovers it.
			for i, it := range items {
				k := clamp32((kc-0.25-float32(i)*0.08)/0.5, 0, 1)
				row := ui.Row(c).Key(i).Absolute().
					Left((w-createMenuW)/2+8).Top((h-menuH)/2+createMenuPad+float32(i)*createMenuRow+(1-k)*6).
					Size(createMenuW-16, createMenuRow).Padding(0, 0, 0, 20).AlignItems(ui.Center).
					Radius(min(o.Corner, 12)).Opacity(k).Cursor(ui.CursorPointer)
				hover := row.Spring("hover", b2f(row.Hovered()), ui.SpringOptions{Duration: 140 * time.Millisecond})
				row.Background(t.SurfaceSubtle.Alpha(hover)).OnClick(func() {
					st.open = false
					st.toggled = c.Now()
					if o.OnCreate != nil {
						o.OnCreate(it)
					}
				}).Children(func() {
					ui.Icon(c, it.Icon).Size(16, 16).TextColor(t.TextMuted).Margin(0, 12, 0, 0)
					ui.Text(c, it.Label).FontSize(14).TextColor(t.TextMuted).SingleLine()
				})
			}
		})
	})
	return holder
}

func init() {
	registerDemo(
		"liq-create",
		"A + Create pill that spreads into the full create menu, as if the menu was inside it all along.",
		func(c *ui.Context) { CreateMenu(c, CreateMenuDefaultItems()) },
	)
}
