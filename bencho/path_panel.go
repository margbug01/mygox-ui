package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// PathPanelTone picks the colour row's swatches.
type PathPanelTone string

const (
	PathPanelColour PathPanelTone = "Colour" // bright swatches, orange first
	PathPanelMono   PathPanelTone = "Mono"   // shades of ink
)

// PathPanelOptions configures PathPanel.
type PathPanelOptions struct {
	Tone PathPanelTone // PathPanelColour by default
}

type pathPanelState struct {
	links    float32 // 1 to pathPanelMaxLinks
	swatch   int
	inverted bool
}

const (
	pathPanelW        = 300
	pathPanelPad      = 6
	pathPanelRowH     = 48
	pathPanelGap      = 4
	pathPanelMaxLinks = 12
)

var (
	pathPanelColours = []ui.Color{ui.Hex("#ff6a3c"), ui.Hex("#3d7bff"), ui.Hex("#17b26a"), ui.Hex("#a35cff"), ui.Hex("#ffc21a")}
	pathPanelInks    = []float32{0, 0.3, 0.55, 0.78} // how far each ink swatch is from Text towards the row
)

// PathPanel is a compact settings panel where each row is its own control:
// drag along Links to set it like a slider, tap Color to cycle through its
// swatches, and tap Inverted to flip its switch.
func PathPanel(c *ui.Context, opts ...PathPanelOptions) ui.Element {
	var o PathPanelOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	card := ui.Column(c).Width(pathPanelW).Padding(pathPanelPad, pathPanelPad, pathPanelPad, pathPanelPad).Gap(pathPanelGap).
		Radius(18).Background(t.Surface)
	st := ui.Local(card, "path_panel", func() pathPanelState { return pathPanelState{links: 6} })
	swatches := len(pathPanelColours)
	if o.Tone == PathPanelMono {
		swatches = len(pathPanelInks)
	}
	swatch := func(i int) ui.Color {
		if o.Tone == PathPanelMono {
			return lerpColor(t.Text, t.SurfaceHover, pathPanelInks[i%len(pathPanelInks)])
		}
		return pathPanelColours[i%len(pathPanelColours)]
	}
	rowW := float32(pathPanelW - 2*pathPanelPad)
	label := func(p *ui.Painter, r ui.Rect, s string, col ui.Color) {
		span := ui.Span{Text: s, Size: 14.5, Weight: 400, Color: col}
		_, h := p.MeasureText(0, span)
		p.RichText(r.X+17, r.Y+(r.H-h)/2, 0, span)
	}
	row := func(key string) ui.Element {
		el := ui.Box(c).Key(key).Size(rowW, pathPanelRowH).Cursor(ui.CursorPointer)
		hover := el.Spring("hover", b2f(el.Hovered() || el.Pressed()), ui.SpringOptions{Duration: 160 * time.Millisecond})
		el.Background(lerpColor(t.SurfaceSubtle, t.SurfaceHover, hover*0.35)).Radius(12)
		return el
	}

	card.Children(func() {
		links := row("links").Cursor(ui.CursorResizeEW)
		_, _, dragging := links.Dragged()
		if dragging || links.Pressed() {
			x, _, _ := links.PointerPosition()
			st.links = float32(math.Round(float64(clamp32(x/rowW*pathPanelMaxLinks, 1, pathPanelMaxLinks))))
		}
		fill := links.Spring("fill", st.links/pathPanelMaxLinks*rowW, ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.2})
		grip := links.Spring("grip", b2f(dragging || links.Pressed()), ui.SpringOptions{Duration: 180 * time.Millisecond})
		links.Draw(func(p *ui.Painter, r ui.Rect) {
			f := ui.Rect{X: r.X, Y: r.Y, W: clamp32(fill, 24, rowW), H: r.H}
			p.Fill(f, lerpColor(t.SurfaceHover, t.BorderStrong, grip*0.5), 12)
			value := fmt.Sprint(int(st.links))
			vs := ui.Span{Text: value, Size: 15, Weight: 400, Color: t.Text}
			vw, vh := p.MeasureText(0, vs)
			vx := r.X + r.W - 16 - vw
			for x := r.X + 32; x < vx-14; x += 33 {
				if x > f.X+f.W+10 {
					p.Fill(ui.Rect{X: x - 1.5, Y: r.Y + r.H/2 - 1.5, W: 3, H: 3}, lerpColor(t.BorderStrong, t.TextSubtle, 0.4), 1.5)
				}
			}
			label(p, r, "Links", t.TextMuted)
			p.RichText(vx, r.Y+(r.H-vh)/2, 0, vs)
		})

		color := row("color").OnClick(func() { st.swatch = (st.swatch + 1) % swatches })
		pop := color.Spring("pop", float32(st.swatch), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.5})
		color.Draw(func(p *ui.Painter, r ui.Rect) {
			label(p, r, "Color", t.TextMuted)
			// The swatch swells as it changes and settles back.
			d := float32(math.Abs(float64(pop - float32(st.swatch))))
			rad := 12 * (1 + 0.25*min(d, 1))
			cx, cy := r.X+r.W-16-12, r.Y+r.H/2
			p.Fill(ui.Rect{X: cx - rad, Y: cy - rad, W: 2 * rad, H: 2 * rad}, swatch(st.swatch), rad)
		})

		inv := row("inverted").OnClick(func() { st.inverted = !st.inverted })
		on := inv.Spring("on", b2f(st.inverted), ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.3})
		inv.Draw(func(p *ui.Painter, r ui.Rect) {
			label(p, r, "Inverted", t.TextMuted)
			track := ui.Rect{X: r.X + r.W - 16 - 40, Y: r.Y + r.H/2 - 12, W: 40, H: 24}
			p.Fill(track, lerpColor(t.SurfaceHover, t.Text, clamp32(on, 0, 1)), 12)
			k := track.X + 2 + on*16
			knob := lerpColor(ui.Hex("#ffffff"), t.Surface, clamp32(on, 0, 1)*b2f(t.Dark))
			p.Fill(ui.Rect{X: k, Y: track.Y + 2, W: 20, H: 20}, knob, 10)
		})
	})
	return card
}

func init() {
	registerDemo(
		"path-panel",
		"A compact settings panel where each row is its own control: drag a row like a slider, tap the colour row to cycle swatches, flip a switch.",
		func(c *ui.Context) { PathPanel(c) },
	)
}
