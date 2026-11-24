package coss

import (
	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// DrawerSide sets which edge a Drawer slides in from.
type DrawerSide uint8

const (
	DrawerBottom DrawerSide = iota
	DrawerRight
	DrawerLeft
	DrawerTop
)

// DrawerVariant is DrawerPopup's variant.
type DrawerVariant uint8

const (
	// DrawerDefault rounds the edge facing the window's middle.
	DrawerDefault DrawerVariant = iota
	// DrawerInset floats the drawer 16px off the window's edges, rounded
	// and bordered all round.
	DrawerInset
	// DrawerStraight has no rounded corners.
	DrawerStraight
)

// DrawerSpec configures a coss-styled edge Drawer built on ui.DialogBase.
type DrawerSpec struct {
	Open    *bool
	Side    DrawerSide
	Variant DrawerVariant
	// ShowBar draws the drag bar (DrawerPopup showBar); ShowClose the
	// close button at end-2 top-2 (showCloseButton).
	ShowBar   bool
	ShowClose bool
	// Center is DrawerHeader className="text-center" with a centered footer.
	Center bool
	// Height fixes a bottom or top drawer's height (a snap point); Tall
	// makes it the window's height less 48 (calc(100dvh-3rem)).
	Height float32
	Tall   bool
	// SnapPoints are the heights a bottom drawer rests at, the first as it
	// opens: DIPs, or fractions of the room the window gives it up to 1,
	// as DrawerRoot's snapPoints. A drag goes one point at a time
	// (snapToSequentialPoints); past the lowest it closes the drawer.
	SnapPoints  []float32
	Title       string
	Description string
	Body        func(close func())
	Footer      func(close func())
	BareFooter  bool
}

// Drawer renders an edge drawer as ui/drawer.tsx: full width at the
// bottom (8px below the window's top at most) or top (48px clear),
// w-[calc(100%-48px)] max-w-md and full height at the sides, a border on
// its inner edge, with Dialog's header, panel and footer.
func Drawer(c *ui.Context, spec DrawerSpec) ui.Element {
	t := theme.Of(c)
	return ui.DialogBase(c, spec.Open, func(backdrop, panel ui.Element) {
		dialogBackdrop(backdrop)
		backdrop.Padding(0)
		inset := spec.Variant == DrawerInset
		r := t.Radius2XL
		if spec.Variant == DrawerStraight {
			r = 0
		}
		panel.Clip().BorderColor(t.Border)
		// The backdrop is a column: Justify places top and bottom drawers,
		// AlignSelf side ones.
		switch spec.Side {
		case DrawerBottom:
			backdrop.Justify(ui.End).Padding(8, 0, 0, 0)
			panel.WidthPercent(100).BorderWidth(1, 0, 0, 0).Radius(r, r, 0, 0)
			if spec.Tall {
				backdrop.Padding(48, 0, 0, 0)
				panel.HeightPercent(100)
			}
			if inset {
				backdrop.Padding(8, 16, 16, 16)
			}
		case DrawerTop:
			backdrop.Justify(ui.Start).Padding(0, 0, 48, 0)
			panel.WidthPercent(100).BorderWidth(0, 0, 1, 0).Radius(0, 0, r, r)
			if inset {
				backdrop.Padding(16, 16, 48, 16)
			}
		case DrawerLeft:
			panel.AlignSelf(ui.Start).WidthPercent(100).MaxWidth(448).HeightPercent(100).
				Margin(0, 48, 0, 0).BorderWidth(0, 1, 0, 0).Radius(0, r, r, 0)
		case DrawerRight:
			panel.AlignSelf(ui.End).WidthPercent(100).MaxWidth(448).HeightPercent(100).
				Margin(0, 0, 0, 48).BorderWidth(0, 0, 0, 1).Radius(r, 0, 0, r)
		}
		if inset {
			if spec.Side == DrawerLeft || spec.Side == DrawerRight {
				backdrop.Padding(16)
			}
			panel.Border(1, t.Border).Radius(t.Radius2XL)
		}
		if spec.Height > 0 {
			panel.Height(spec.Height)
		}
		drawerSwipe(c, spec, backdrop, panel)
		var overlay func()
		if spec.ShowBar {
			// has-data-[slot=drawer-bar]: 8px on the bar's side.
			switch spec.Side {
			case DrawerBottom:
				panel.Padding(8, 0, 0, 0)
			case DrawerTop:
				panel.Padding(0, 0, 8, 0)
			case DrawerLeft:
				panel.Padding(0, 8, 0, 0)
			case DrawerRight:
				panel.Padding(0, 0, 0, 8)
			}
			overlay = func() { drawerBar(c, t, spec.Side) }
		}
		dialogPopup(c, t, panel, dialogContent{
			open:        spec.Open,
			title:       spec.Title,
			description: spec.Description,
			hideClose:   !spec.ShowClose,
			body:        spec.Body,
			footer:      spec.Footer,
			bareFooter:  spec.BareFooter,
			growBody:    spec.Tall || spec.Height > 0 || len(spec.SnapPoints) > 0,
			center:      spec.Center,
			overlay:     overlay,
		})
	})
}

// drawerDrag is a drawer's offset toward its edge and the rest it goes to.
type drawerDrag struct {
	off  ui.Point
	rest int
	init bool
}

// drawerSwipe lets the drawer be dragged toward its edge, as base-ui's
// Drawer swipes: let go far or fast enough, it slides out and closes;
// else it goes back, or to the nearest of SnapPoints. The backdrop fades
// as it goes.
func drawerSwipe(c *ui.Context, spec DrawerSpec, backdrop, panel ui.Element) {
	_, winH := c.Size()
	// sign is +1 where the edge is down or right, the axis 0 for y.
	sign, vertical := float32(1), true
	switch spec.Side {
	case DrawerTop:
		sign = -1
	case DrawerLeft:
		sign, vertical = -1, false
	case DrawerRight:
		vertical = false
	}
	b := panel.LayoutBounds()
	size := b.H
	if !vertical {
		size = b.W
	}
	// rests are offsets toward the edge, smallest first; the last closes.
	var rests []float32
	if len(spec.SnapPoints) > 0 && spec.Side == DrawerBottom {
		full := winH - 8
		panel.Height(full)
		size = full
		for i := len(spec.SnapPoints) - 1; i >= 0; i-- {
			v := spec.SnapPoints[i]
			if v <= 1 {
				v *= full
			}
			rests = append(rests, full-min(v, full))
		}
	} else {
		rests = []float32{0}
	}
	closed := len(rests)
	rests = append(rests, size+16)
	first := closed - 1
	if len(spec.SnapPoints) == 0 {
		first = 0
	}

	st := ui.Local(panel, "drawer-drag", func() drawerDrag { return drawerDrag{} })
	if !st.init && size > 0 {
		st.init, st.rest = true, first
		st.off = axisPoint(vertical, sign*rests[first])
	}
	sequential := len(spec.SnapPoints) > 0
	o := ui.DraggableOptions{
		Snap: func(p ui.Point) ui.Point {
			d := sign * p.Y
			if !vertical {
				d = sign * p.X
			}
			lo, hi := 0, len(rests)-1
			if sequential {
				lo, hi = max(st.rest-1, 0), min(st.rest+1, len(rests)-1)
			}
			best := lo
			for i := lo; i <= hi; i++ {
				if abs32(rests[i]-d) < abs32(rests[best]-d) {
					best = i
				}
			}
			st.rest = best
			return axisPoint(vertical, sign*rests[best])
		},
	}
	if vertical {
		o.Axis = ui.Vertical
		lo := min(0, sign*rests[closed])
		o.Bounds = &ui.Rect{Y: lo, H: rests[closed]}
	} else {
		o.Axis = ui.Horizontal
		lo := min(0, sign*rests[closed])
		o.Bounds = &ui.Rect{X: lo, W: rests[closed]}
	}
	panel.Draggable(&st.off, o)

	shown := panel.DragOffset()
	d := sign * shown.Y
	if !vertical {
		d = sign * shown.X
	}
	p := float32(0)
	if span := rests[closed] - rests[first]; span > 0 {
		p = min(max((d-rests[first])/span, 0), 1)
	}
	if p > 0 {
		backdrop.Material(dimBlur{glass.Blur{Radius: 8 * (1 - p)}, theme.Alpha(ui.RGB(0, 0, 0), 0.32*(1-p))})
	}
	if st.rest == closed && !panel.Dragging() && !panel.Transitioning() && spec.Open != nil {
		*spec.Open = false
		*st = drawerDrag{}
	}
}

func axisPoint(vertical bool, v float32) ui.Point {
	if vertical {
		return ui.Point{Y: v}
	}
	return ui.Point{X: v}
}

func abs32(v float32) float32 { return max(v, -v) }

// drawerBar is DrawerBar: a 48x4 bg-input pill in a p-3 strip along the
// drawer's open edge (4x48 for side drawers).
func drawerBar(c *ui.Context, t theme.Tokens, side DrawerSide) {
	bar := ui.Row(c).Absolute().Padding(12).Center()
	pill := func(w, h float32) {
		bar.Children(func() { ui.Box(c).Size(w, h).Radius(999).Background(t.Input) })
	}
	switch side {
	case DrawerBottom:
		bar.Top(0).Left(0).Right(0)
		pill(48, 4)
	case DrawerTop:
		bar.Bottom(0).Left(0).Right(0)
		pill(48, 4)
	case DrawerLeft:
		bar.Right(0).Top(0).Bottom(0)
		pill(4, 48)
	case DrawerRight:
		bar.Left(0).Top(0).Bottom(0)
		pill(4, 48)
	}
}
