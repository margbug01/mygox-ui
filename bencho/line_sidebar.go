package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// LineSidebarOptions configures LineSidebar.
type LineSidebarOptions struct {
	Reach  float32 // how many rows the pull spreads across, 0 to 100; 50 by default
	Length float32 // how far a line reaches out, 0 to 100; 40 by default
}

type lineSidebarState struct {
	active int
}

const (
	lineSidebarW     = 300
	lineSidebarRow   = 34
	lineSidebarTop   = 22.5 // the first row's middle
	lineSidebarX     = 16
	lineSidebarShort = 22 // a line at rest
)

// LineSidebarDefaultSections is the table of contents bencho.dev shows.
func LineSidebarDefaultSections() []string {
	return []string{"Overview", "Components", "Animations", "Backgrounds", "Showcase", "Playground", "Templates", "Changelog"}
}

// LineSidebar is a table of contents drawn like a ruler. Running the pointer
// down it makes the nearest lines reach out while their names lean in, the
// pull falling off across the neighbours; pressing a row makes it current.
func LineSidebar(c *ui.Context, sections []string, opts ...LineSidebarOptions) ui.Element {
	var o LineSidebarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Reach == 0 {
		o.Reach = 50
	}
	if o.Length == 0 {
		o.Length = 40
	}
	o.Reach = clamp32(o.Reach, 0, 100)
	o.Length = clamp32(o.Length, 0, 100)
	t := ThemeFrom(c)
	h := lineSidebarTop + 35.5 + float32(max(len(sections)-1, 0))*lineSidebarRow
	nav := ui.Box(c).Size(lineSidebarW, h)
	st := ui.Local(nav, "line_sidebar", func() lineSidebarState { return lineSidebarState{} })

	// The pull is a bell around the pointer, Reach rows wide.
	_, py, over := nav.PointerPosition()
	sigma := 0.35 + o.Reach/100*1.6 // in rows
	reach := o.Length / 40 * 36     // how much longer a fully pulled line grows
	grey := lerpColor(t.TextSubtle, t.SurfaceSubtle, 0.15)
	if t.Dark {
		grey = lerpColor(t.TextSubtle, t.TextMuted, 0.4)
	}

	nav.Children(func() {
		for i := 0; i+1 < len(sections); i++ {
			ui.Box(c).Key(fmt.Sprint("tick", i)).Absolute().Left(lineSidebarX).Top(lineSidebarTop+(float32(i)+0.5)*lineSidebarRow).
				Size(8, 1).Background(t.BorderStrong)
		}
		for i, name := range sections {
			y := lineSidebarTop + float32(i)*lineSidebarRow
			pull := float32(0)
			if over {
				d := float64((py - y) / lineSidebarRow / sigma)
				pull = float32(math.Exp(-d * d / 2))
			}
			row := ui.Box(c).Key(i).Absolute().Left(0).Top(y-lineSidebarRow/2).Size(lineSidebarW, lineSidebarRow).
				Cursor(ui.CursorPointer).OnClick(func() { st.active = i })
			on := i == st.active
			p := row.Spring("pull", max(pull, b2f(on)), ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: 0.15})
			ink := row.Spring("ink", b2f(on), ui.SpringOptions{Duration: 220 * time.Millisecond})
			lineW := lineSidebarShort + reach*p
			lean := 6 * p
			lineCol := lerpColor(grey, t.Text, ink)
			numCol := lerpColor(lerpColor(t.TextSubtle, t.SurfaceSubtle, 0.45), t.TextMuted, max(ink, p*0.6))
			nameCol := lerpColor(grey, t.Text, max(ink, p*0.35))
			row.Children(func() {
				ui.Box(c).Absolute().Left(lineSidebarX).Top(lineSidebarRow/2-0.75).Size(lineW, 1.5).Background(lineCol)
				ui.Row(c).Absolute().Left(lineSidebarX + lineW + 10 + lean).Top(0).Height(lineSidebarRow).AlignItems(ui.Center).Children(func() {
					ui.Text(c, fmt.Sprintf("%02d", i+1)).FontSize(12.5).FontFeatures("tnum").TextColor(numCol).SingleLine().Width(21)
					ui.Text(c, name).FontSize(16.2).FontWeight(400+100*int(ink+0.5)).TextColor(nameCol).SingleLine().Margin(0, 0, 0, 5.5)
				})
			})
		}
	})
	return nav
}

func init() {
	registerDemo(
		"line-sidebar",
		"A table of contents drawn like a ruler: the lines nearest the pointer reach out and their names lean in, falling off across the neighbours.",
		func(c *ui.Context) { LineSidebar(c, LineSidebarDefaultSections()) },
	)
}
