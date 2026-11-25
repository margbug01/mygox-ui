package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// RadialMenuOptions configures RadialMenu.
type RadialMenuOptions struct {
	Count   int           // how many tools are on the fan, 2 to 6; 5 by default
	Stagger time.Duration // how far apart they deal out, 0 to 70ms; 28ms by default
	Radius  float32       // 44 to 104; 68 by default
	OnPick  func(i int)
}

type radialMenuState struct {
	open     bool
	toggled  time.Time
	selected int
}

const (
	radialMenuButton = 56
	radialMenuTool   = 44
)

func radialMenuIcon(paths string) *ui.SVG {
	return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">` + paths + `</svg>`))
}

// Lucide's eraser, type and square (ISC licence), which coss/icons doesn't carry.
var radialMenuTools = []*ui.SVG{
	icons.Pencil(),
	radialMenuIcon(`<path d="m7 21-4.3-4.3c-1-1-1-2.5 0-3.4l9.6-9.6c1-1 2.5-1 3.4 0l5.6 5.6c1 1 1 2.5 0 3.4L13 21"/><path d="M22 21H7"/><path d="m5 11 9 9"/>`),
	radialMenuIcon(`<path d="M4 7V4h16v3"/><path d="M9 20h6"/><path d="M12 4v16"/>`),
	radialMenuIcon(`<rect width="18" height="18" x="3" y="3" rx="2"/>`),
	icons.Image(),
	icons.Circle(),
}

// RadialMenu keeps tools tucked behind one round button. Pressing it turns
// the plus into a cross and deals the tools out around it in a half circle,
// one after another; pressing a tool picks it and folds the fan away.
func RadialMenu(c *ui.Context, opts ...RadialMenuOptions) ui.Element {
	var o RadialMenuOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Count == 0 {
		o.Count = 5
	}
	o.Count = min(max(o.Count, 2), len(radialMenuTools))
	if o.Stagger == 0 {
		o.Stagger = 28 * time.Millisecond
	}
	if o.Radius == 0 {
		o.Radius = 68
	}
	t := ThemeFrom(c)
	size := 2*o.Radius + radialMenuTool
	holder := ui.Box(c).Size(size, size)
	st := ui.Local(holder, "radial_menu", func() radialMenuState { return radialMenuState{selected: -1} })
	now := c.Now()
	since := now.Sub(st.toggled)
	cx, cy := size/2, size/2

	holder.Children(func() {
		for i := 0; i < o.Count; i++ {
			// Deal out left to right on the way open, and back right to left.
			order := i
			if !st.open {
				order = o.Count - 1 - i
			}
			wait := time.Duration(order) * o.Stagger
			target := b2f(st.open && since >= wait) + b2f(!st.open && since < wait)
			if st.toggled.IsZero() {
				target = 0
			}
			if since < wait+o.Stagger {
				c.AnimationFrame()
			}
			tool := ui.Box(c).Key(i)
			k := tool.Spring("out", target, ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.3})
			if k < 0.01 && !st.open {
				continue
			}
			a := math.Pi * (1 - float64(i)/float64(o.Count-1)) // from the left, over the top, to the right
			x := cx + float32(math.Cos(a))*o.Radius*k
			y := cy - float32(math.Sin(a))*o.Radius*k
			on := i == st.selected
			hover := tool.Spring("hover", b2f(tool.Hovered()), ui.SpringOptions{Duration: 160 * time.Millisecond, Bounce: 0.3})
			bg, fg := t.Surface, t.Text
			if on {
				bg, fg = t.Text, t.Surface
			}
			s := clamp32(k, 0, 1.2) * (1 + 0.08*hover)
			tool.Absolute().Left(x-radialMenuTool/2).Top(y-radialMenuTool/2).Size(radialMenuTool, radialMenuTool).
				Radius(radialMenuTool/2).Background(bg).Center().Scale(s, s).Opacity(clamp32(k*1.4, 0, 1)).
				Shadow(0, 4, 12, 0, ui.RGBA(0, 0, 0, 0.06)).Cursor(ui.CursorPointer).
				OnClick(func() {
					st.selected = i
					st.open = false
					st.toggled = c.Now()
					if o.OnPick != nil {
						o.OnPick(i)
					}
				}).
				Children(func() {
					ui.Icon(c, radialMenuTools[i]).Size(18, 18).TextColor(fg)
				})
		}

		btn := ui.Box(c).Key("button").Absolute().Left(cx-radialMenuButton/2).Top(cy-radialMenuButton/2).
			Size(radialMenuButton, radialMenuButton).Radius(radialMenuButton / 2).Background(t.Surface).Center().
			Cursor(ui.CursorPointer).OnClick(func() {
			st.open = !st.open
			st.toggled = c.Now()
		})
		turn := btn.Spring("turn", b2f(st.open), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.25})
		press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
		btn.Scale(1-0.06*press, 1-0.06*press).Children(func() {
			ui.Icon(c, icons.Plus()).Size(20, 20).TextColor(t.Text).Rotate(45 * turn)
		})
	})
	return holder
}

func init() {
	registerDemo(
		"radial",
		"A round button that deals its tools out around it in a half circle.",
		func(c *ui.Context) { RadialMenu(c) },
	)
}
