package bencho

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// ResizablePanelsLayout picks how many panes ResizablePanels shows.
type ResizablePanelsLayout string

const (
	ResizablePanelsTwo   ResizablePanelsLayout = "2" // sidebar and content
	ResizablePanelsThree ResizablePanelsLayout = "3" // details beside them too
)

// ResizablePanelsOptions configures ResizablePanels.
type ResizablePanelsOptions struct {
	Panels ResizablePanelsLayout // ResizablePanelsTwo by default
	Corner float32               // 0 to 28; 20 by default
}

type resizablePanelsState struct {
	size     [2]float32 // the sidebar's and the details' settled widths; 0 is folded
	raw      [2]float32 // where a divider's drag has taken its pane, unresisted
	dragging [2]bool
}

const (
	resizablePanelsW       = 300
	resizablePanelsH       = 180
	resizablePanelsSideMin = 80  // a side pane resists below this…
	resizablePanelsFold    = 36  // …and folds away this far past it
	resizablePanelsMainMin = 100 // the content pane never gets narrower
	resizablePanelsStep    = 8   // an arrow key's move
)

var resizablePanelsBars = map[string][]float32{
	"Sidebar": {0.72, 0.56, 0.64},
	"Content": {0.92, 0.8, 0.86, 0.54},
	"Details": {0.8, 0.6, 0.7},
}

// ResizablePanels is a card split into a sidebar and content (and details,
// with three panels) by draggable dividers. Dragging a divider trades width
// between its panes; a side pane resists at its minimum and folds away when
// pulled past it, and drags back out from the edge. A focused divider moves
// with the arrow keys.
func ResizablePanels(c *ui.Context, opts ...ResizablePanelsOptions) ui.Element {
	var o ResizablePanelsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	three := o.Panels == ResizablePanelsThree
	n := 1 // dividers
	if three {
		n = 2
	}
	t := ThemeFrom(c)
	card := ui.Box(c).Size(resizablePanelsW, resizablePanelsH).Radius(clamp32(o.Corner, 0, 28)).Background(t.Surface).Clip()
	st := ui.Local(card, "resizable_panels", func() resizablePanelsState {
		if three {
			return resizablePanelsState{size: [2]float32{96, 88}}
		}
		return resizablePanelsState{size: [2]float32{128, 0}}
	})
	room := float32(resizablePanelsW - n) // the panes' widths add up to this

	// A side pane's width for an unresisted drag to raw, with the other side
	// pane at other: rubber at both ends, folded well past its minimum.
	shape := func(raw, other float32) float32 {
		hi := room - other - resizablePanelsMainMin
		switch {
		case raw > hi:
			return hi + rubberBand(raw-hi, 24)
		case raw >= resizablePanelsSideMin:
			return raw
		case raw < resizablePanelsSideMin-resizablePanelsFold:
			return 0
		}
		return resizablePanelsSideMin - rubberBand(resizablePanelsSideMin-raw, 16)
	}
	// Let go, a resisting pane springs back to its minimum and a folding one
	// folds.
	settle := func(raw, other float32) float32 {
		switch {
		case raw < resizablePanelsSideMin-resizablePanelsFold:
			return 0
		case raw < resizablePanelsSideMin:
			return resizablePanelsSideMin
		}
		return min(raw, room-other-resizablePanelsMainMin)
	}

	// Each frame's targets, then springs: quick under the finger, softer after.
	var w [2]float32
	for i := 0; i < n; i++ {
		target := st.size[i]
		if st.dragging[i] {
			target = shape(st.raw[i], st.size[1-i])
		}
		opt := ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.12}
		if st.dragging[i] && target > 0 {
			opt = ui.SpringOptions{Duration: 60 * time.Millisecond}
		}
		w[i] = max(card.Spring(fmt.Sprint("w", i), target, opt), 0)
	}
	main := room - w[0] - w[1]

	type pane struct {
		name string
		x, w float32
	}
	panes := []pane{{"Sidebar", 0, w[0]}, {"Content", w[0] + 1, main}}
	if three {
		panes = append(panes, pane{"Details", w[0] + 1 + main + 1, w[1]})
	}

	card.Children(func() {
		for _, pn := range panes {
			ui.Box(c).Key(pn.name).Absolute().Left(pn.x).Top(0).Size(max(pn.w, 0), resizablePanelsH).Clip().Children(func() {
				ui.Text(c, pn.name).Absolute().Left(12).Top(15).FontSize(14).FontWeight(500).TextColor(t.Text).SingleLine()
				ui.Text(c, fmt.Sprintf("%.0f px", pn.w)).Absolute().Left(12).Top(35).FontSize(13.5).TextColor(lerpColor(t.TextSubtle, t.Surface, 0.15)).SingleLine()
				for k, f := range resizablePanelsBars[pn.name] {
					ui.Box(c).Key(k).Absolute().Left(12).Top(68+float32(k)*14).Size(max(f*(pn.w-24), 0), 6).Radius(3).Background(t.SurfaceSubtle)
				}
			})
		}
		for i := 0; i < n; i++ {
			x := w[0] // the sidebar's divider
			if i == 1 {
				x = w[0] + 1 + main
			}
			// Dragging right widens the sidebar but narrows the details.
			dir := float32(1 - 2*i)
			ui.Box(c).Key(fmt.Sprint("line", i)).Absolute().Left(x).Top(0).Size(1, resizablePanelsH).Background(t.Border)
			handle := ui.Box(c).Key(fmt.Sprint("divider", i)).Absolute().Left(x-6).Top(0).Size(13, resizablePanelsH).
				Cursor(ui.CursorResizeEW).Focusable()
			dx, _, dragging := handle.Dragged()
			if dragging && !st.dragging[i] {
				st.raw[i] = st.size[i] // a folded pane comes back once pulled past the fold
			}
			if dragging {
				st.raw[i] += dx * dir
			} else if st.dragging[i] {
				st.size[i] = settle(st.raw[i], st.size[1-i])
			}
			st.dragging[i] = dragging
			key := func(by float32) {
				s := st.size[i] + by
				switch {
				case st.size[i] == 0 && by > 0:
					s = resizablePanelsSideMin
				case s < resizablePanelsSideMin:
					s = 0
				}
				st.size[i] = clamp32(s, 0, room-st.size[1-i]-resizablePanelsMainMin)
			}
			if handle.Shortcut(0, ui.KeyRight) {
				key(resizablePanelsStep * dir)
			}
			if handle.Shortcut(0, ui.KeyLeft) {
				key(-resizablePanelsStep * dir)
			}
			on := handle.Spring("on", b2f(handle.Hovered() || dragging || handle.Focused()), ui.SpringOptions{Duration: 160 * time.Millisecond})
			handle.Children(func() {
				ui.Box(c).Absolute().Left(4.5).Top(resizablePanelsH/2-14).Size(4, 28).Radius(2).
					Background(t.BorderStrong).Opacity(on).Scale(1, 0.6+0.4*on)
			})
		}
	})
	return card
}

func init() {
	registerDemo(
		"resizable-panels",
		"Resizable panels with draggable dividers that also work from the keyboard: the panes trade width, resist at their minimum and fold away past it.",
		func(c *ui.Context) { ResizablePanels(c) },
	)
}
