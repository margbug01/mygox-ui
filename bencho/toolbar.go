package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// CanvasToolbarOptions configures CanvasToolbar.
type CanvasToolbarOptions struct {
	Corner float32 // 0 to 25; 20 by default
	// OnTool is called with the tool picked: "Select", "Frame", "Text",
	// "Pen", "Embed", or the shape tool's shape ("Rectangle", "Ellipse",
	// "Triangle", "Arrow").
	OnTool func(tool string)
}

type canvasToolbarTool struct {
	name string
	icon *ui.SVG
}

func canvasToolbarSVG(body string) *ui.SVG {
	return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">` + body + `</svg>`))
}

// Outline glyphs in Lucide's stroke style.
var (
	canvasToolbarTools = []canvasToolbarTool{
		{"Select", canvasToolbarSVG(`<path d="M4.5 4.5 19.5 10.6l-6.4 2-2 6.4z"/>`)},
		{"Frame", canvasToolbarSVG(`<path d="M9 3.5v17M15 3.5v17M3.5 9h17M3.5 15h17"/>`)},
		{"Shape", nil}, // drawn with the last shape picked
		{"Text", canvasToolbarSVG(`<path d="M12 4.5v15M5 7.5v-3h14v3M9 19.5h6"/>`)},
		{"Pen", canvasToolbarSVG(`<circle cx="18.5" cy="5.5" r="2"/><circle cx="5.5" cy="18.5" r="2"/><path d="M5.5 16.5C5.5 10 10 5.5 16.5 5.5"/>`)},
		{"Embed", canvasToolbarSVG(`<path d="m15 7 5 5-5 5M9 7l-5 5 5 5"/>`)},
	}
	canvasToolbarShapes = []canvasToolbarTool{
		{"Rectangle", canvasToolbarSVG(`<rect x="4.5" y="4.5" width="15" height="15" rx="3"/>`)},
		{"Ellipse", canvasToolbarSVG(`<circle cx="12" cy="12" r="8"/>`)},
		{"Triangle", canvasToolbarSVG(`<path d="M12 4.5 20 19H4z"/>`)},
		{"Arrow", canvasToolbarSVG(`<path d="M5 19 19 5M10 5h9v9"/>`)},
	}
)

const (
	canvasToolbarW    = 264
	canvasToolbarH    = 50
	canvasToolbarTop  = 50
	canvasToolbarBtn  = 38
	canvasToolbarStep = 40
	canvasToolbarPad  = 7
	canvasToolbarBlue = "#0e99ff"
)

type canvasToolbarState struct {
	tool  int // index into canvasToolbarTools
	shape int // the shape tool's last shape
	menu  bool
}

// canvasToolbarX is where tool i's button starts: the divider before Pen
// takes a little room.
func canvasToolbarX(i float32) float32 {
	x := canvasToolbarPad + i*canvasToolbarStep
	if i > 3 {
		x += clamp32(i-3, 0, 1) * 13
	}
	return x
}

// CanvasToolbar is a floating toolbar for a drawing canvas. A blue tile
// slides to the tool pressed; the shape tool remembers the last shape
// picked, and pressing it again (or its corner) opens the shapes to pick
// another.
func CanvasToolbar(c *ui.Context, opts ...CanvasToolbarOptions) ui.Element {
	var o CanvasToolbarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Corner = clamp32(o.Corner, 0, 25)
	t := ThemeFrom(c)
	blue := ui.Hex(canvasToolbarBlue)

	root := ui.Box(c).Size(canvasToolbarW, 150)
	st := ui.Local(root, "canvas_toolbar", func() canvasToolbarState { return canvasToolbarState{} })
	pick := func(name string) {
		if o.OnTool != nil {
			o.OnTool(name)
		}
	}
	slide := root.Spring("tile", float32(st.tool), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.25})
	menu := root.Spring("menu", b2f(st.menu), ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: 0.25})
	btnR := min(o.Corner*0.75, canvasToolbarBtn/2)

	root.Children(func() {
		ui.Box(c).Absolute().Left(0).Top(canvasToolbarTop).Size(canvasToolbarW, canvasToolbarH).
			Radius(min(o.Corner*1.25, canvasToolbarH/2)).Background(t.Surface).Children(func() {
			ui.Box(c).Absolute().Left(canvasToolbarX(slide)).Top((canvasToolbarH-canvasToolbarBtn)/2).
				Size(canvasToolbarBtn, canvasToolbarBtn).Radius(btnR).Background(blue)
			ui.Box(c).Absolute().Left(canvasToolbarX(4)-7).Top((canvasToolbarH-20)/2).Size(1, 20).Background(t.Border)

			for i, tool := range canvasToolbarTools {
				on := i == st.tool
				icon := tool.icon
				if icon == nil {
					icon = canvasToolbarShapes[st.shape].icon
				}
				btn := ui.Box(c).Key(i).Absolute().Left(canvasToolbarX(float32(i))).Top((canvasToolbarH-canvasToolbarBtn)/2).
					Size(canvasToolbarBtn, canvasToolbarBtn).Radius(btnR).Cursor(ui.CursorPointer).
					Role(ui.RoleRadio).Label(tool.name).AccessSelected(on).
					OnClick(func() {
						if tool.icon == nil && st.tool == i {
							st.menu = !st.menu
							return
						}
						st.tool, st.menu = i, false
						if tool.icon == nil {
							pick(canvasToolbarShapes[st.shape].name)
							return
						}
						pick(tool.name)
					})
				hover := btn.Spring("hover", b2f(btn.Hovered() && !on), ui.SpringOptions{Duration: 150 * time.Millisecond})
				press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 120 * time.Millisecond})
				lit := clamp32(1-max(slide-float32(i), float32(i)-slide), 0, 1)
				btn.Background(t.SurfaceSubtle.Alpha(hover)).Scale(1-0.08*press, 1-0.08*press).Children(func() {
					ui.Icon(c, icon).Absolute().Left(8).Top(8).Size(22, 22).TextColor(lerpColor(t.Text, ui.Hex("#ffffff"), lit))
					if tool.icon == nil { // the corner mark: this tool has more
						ui.Box(c).Absolute().Left(canvasToolbarBtn-11).Top(canvasToolbarBtn-11).Size(5, 5).Draw(func(p *ui.Painter, r ui.Rect) {
							var tri ui.Path
							tri.MoveTo(r.X+r.W, r.Y).LineTo(r.X+r.W, r.Y+r.H).LineTo(r.X, r.Y+r.H).Close()
							p.FillPath(&tri, lerpColor(t.TextSubtle, ui.Hex("#ffffff"), lit*0.8))
						})
					}
				})
			}
		})

		// The shapes open under the shape tool.
		if menu > 0.01 {
			mw := float32(len(canvasToolbarShapes))*36 + 8
			mx := clamp32(canvasToolbarX(2)+canvasToolbarBtn/2-mw/2, 0, canvasToolbarW-mw)
			ui.Box(c).Absolute().Left(mx).Top(canvasToolbarTop+canvasToolbarH+6).Size(mw, 44).Radius(min(o.Corner, 22)).
				Background(t.Surface).Shadow(0, 6, 18, 0, ui.RGBA(0, 0, 0, 0.08)).
				Opacity(clamp32(menu, 0, 1)).Scale(lerp32(0.9, 1, menu), lerp32(0.9, 1, menu)).Translate(0, (1-menu)*-8).
				Children(func() {
					for i, s := range canvasToolbarShapes {
						on := i == st.shape
						b := ui.Box(c).Key(s.name).Absolute().Left(4+float32(i)*36).Top(4).Size(36, 36).Radius(min(o.Corner*0.6, 18)).
							Cursor(ui.CursorPointer).Label(s.name).OnClick(func() {
							st.shape, st.menu = i, false
							pick(s.name)
						})
						b.Background(lerpColor(t.Surface, t.SurfaceSubtle, max(b2f(b.Hovered()), b2f(on)))).Children(func() {
							ui.Icon(c, s.icon).Absolute().Left(9).Top(9).Size(18, 18).TextColor(t.Text)
						})
					}
				})
		}
	})
	return root
}

func init() {
	registerDemo(
		"toolbar",
		"A floating canvas toolbar: a blue tile slides to the tool pressed, and the shape tool remembers the last shape picked.",
		func(c *ui.Context) { CanvasToolbar(c) },
	)
}
