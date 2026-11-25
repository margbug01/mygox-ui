package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// glideTipLink is Lucide's link-2, the level chain bencho.dev shows.
var glideTipLink = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 17H7A5 5 0 0 1 7 7h2"/><path d="M15 7h2a5 5 0 1 1 0 10h-2"/><line x1="8" x2="16" y1="12" y2="12"/></svg>`))

var glideTipButtons = []struct {
	label string
	icon  func() *ui.SVG
}{
	{"Copy link", func() *ui.SVG { return glideTipLink }},
	{"Share", icons.Share},
	{"Save", icons.Bookmark},
	{"Delete", icons.Trash},
}

const (
	glideTipD    = 40 // button diameter
	glideTipStep = 48 // between button left edges
	glideTipW    = glideTipStep*3 + glideTipD
	glideTipH    = 80
	glideTipPill = 28 // tooltip height
)

type glideTipState struct {
	at      int     // the button the tooltip last stood over
	entered int     // times the pointer came onto the row; a new entry starts the tooltip on its button
	pillW   float32 // the tooltip's width, measured when it was last painted
}

// GlideTip is a row of four icon buttons sharing one liquid tooltip. It grows
// out of the button you point at, flows along the row to the next one as you
// move, and drains back into its button when the pointer leaves.
func GlideTip(c *ui.Context) ui.Element {
	t := ThemeFrom(c)
	row := ui.Box(c).Size(glideTipW, glideTipH)
	st := ui.Local(row, "glide_tip", func() glideTipState { return glideTipState{} })

	hover := -1
	if px, py, over := row.PointerPosition(); over && py >= glideTipH-glideTipD-4 {
		i := int(px / glideTipStep)
		if i >= 0 && i < len(glideTipButtons) && px-float32(i*glideTipStep) <= glideTipD+4 {
			hover = i
		}
	}
	grow := row.Spring("grow", b2f(hover >= 0), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.25})
	if hover >= 0 {
		if grow < 0.05 && hover != st.at {
			st.entered++
		}
		st.at = hover
	}
	center := func(i int) float32 { return float32(i*glideTipStep) + glideTipD/2 }
	x := row.Spring(st.entered, center(st.at), ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.2})
	label := glideTipButtons[st.at].label
	// The pill stays inside the row; at the ends the neck leans to its button.
	px := clamp32(x, st.pillW/2, glideTipW-st.pillW/2)

	row.Children(func() {
		for i, b := range glideTipButtons {
			btn := ui.Box(c).Key(i).Absolute().Left(float32(i*glideTipStep)).Top(glideTipH-glideTipD).Size(glideTipD, glideTipD).
				Radius(glideTipD / 2).Background(t.Surface).Center().Cursor(ui.CursorPointer)
			press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 160 * time.Millisecond})
			btn.Scale(1-0.06*press, 1-0.06*press)
			btn.Children(func() { ui.Icon(c, b.icon()).Size(20, 20).TextColor(t.Text) })
		}
		// The liquid: a pill above the row joined to its button by a neck
		// that pinches in the middle and leans while the pill catches up.
		base := center(st.at)
		ui.Box(c).Absolute().Left(0).Top(0).Size(glideTipW, glideTipH).Draw(func(p *ui.Painter, r ui.Rect) {
			if grow < 0.01 {
				return
			}
			g := clamp32(grow, 0, 1.2)
			w, _ := p.MeasureText(0, ui.Span{Text: label, Size: 13, Weight: 500})
			st.pillW = w + 24
			pw := st.pillW * clamp32(g, 0, 1)
			ph := glideTipPill * clamp32(g, 0, 1)
			top := lerp32(glideTipH-glideTipD/2, 4, g)
			pill := ui.Rect{X: r.X + px - pw/2, Y: r.Y + top, W: pw, H: ph}
			p.Fill(pill, t.Surface, ph/2)
			neckTop, neckBot := r.Y+top+ph-2, r.Y+glideTipH-glideTipD+2
			if neckBot <= neckTop {
				return
			}
			half := 9 * clamp32(g, 0, 1)
			pinch := half * 0.45
			midY := (neckTop + neckBot) / 2
			tx, bx := r.X+px, r.X+base
			mx := (tx + bx) / 2
			var neck ui.Path
			neck.MoveTo(tx-half, neckTop).QuadTo(mx-pinch, midY, bx-half, neckBot).
				LineTo(bx+half, neckBot).QuadTo(mx+pinch, midY, tx+half, neckTop).Close()
			p.FillPath(&neck, t.Surface)
		})
		if grow < 0.5 {
			return
		}
		ui.Text(c, label).FontSize(13).FontWeight(500).TextColor(t.Text).Opacity(clamp32(grow*2-1, 0, 1)).
			Absolute().Left(px - 60).Top(4 + (glideTipPill-18)/2).Width(120).TextAlign(ui.Center)
	})
	return row
}

func init() {
	registerDemo(
		"glide-tip",
		"A liquid tooltip that grows out of the icon button you point at, flows along the row to the next one, and drains away when the pointer leaves.",
		func(c *ui.Context) { GlideTip(c) },
	)
}
