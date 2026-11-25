package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// ResizeFrameSnap is whether a resize snaps to common ratios.
type ResizeFrameSnap string

const (
	ResizeFrameSnapOn  ResizeFrameSnap = "On"
	ResizeFrameSnapOff ResizeFrameSnap = "Off"
)

// ResizeFrameOptions configures ResizeFrame.
type ResizeFrameOptions struct {
	Snap   ResizeFrameSnap // ratio snap; On by default
	Corner float32         // the frame's corner radius, 0 to 24; 0 by default
}

// resizeFrameSide is a set of the frame's sides a drag moves.
type resizeFrameSide uint8

const (
	resizeFrameLeft resizeFrameSide = 1 << iota
	resizeFrameRight
	resizeFrameTop
	resizeFrameBottom
)

type resizeFrameState struct {
	l, t, r, b float32         // the frame's edges in the stage
	grab       resizeFrameSide // the sides the current drag moves
	offX, offY float32         // pointer minus the grabbed edges at press
	ratio      string          // the ratio it snapped to, "" when free
	dragging   bool
}

const (
	resizeFrameW      = 440
	resizeFrameH      = 330
	resizeFrameMargin = 12
	resizeFrameMinW   = 60
	resizeFrameMinH   = 40
	resizeFrameHandle = 9
	resizeFrameReach  = 8 // how near an edge a press grabs it
)

var resizeFrameRatios = []struct {
	name string
	w, h float32
}{{"1:1", 1, 1}, {"4:3", 4, 3}, {"3:2", 3, 2}, {"16:9", 16, 9}, {"3:4", 3, 4}, {"2:3", 2, 3}, {"9:16", 9, 16}}

// resizeFrameHit is the sides a press at x, y grabs: a corner, an edge or
// none.
func resizeFrameHit(st *resizeFrameState, x, y float32) resizeFrameSide {
	in := func(v, lo, hi float32) bool { return v >= lo-resizeFrameReach && v <= hi+resizeFrameReach }
	if !in(x, st.l, st.r) || !in(y, st.t, st.b) {
		return 0
	}
	near := func(v, edge float32) bool { return float32(math.Abs(float64(v-edge))) <= resizeFrameReach }
	var s resizeFrameSide
	if near(x, st.l) {
		s |= resizeFrameLeft
	} else if near(x, st.r) {
		s |= resizeFrameRight
	}
	if near(y, st.t) {
		s |= resizeFrameTop
	} else if near(y, st.b) {
		s |= resizeFrameBottom
	}
	return s
}

func resizeFrameCursor(s resizeFrameSide) ui.Cursor {
	switch s {
	case resizeFrameLeft | resizeFrameTop, resizeFrameRight | resizeFrameBottom:
		return ui.CursorResizeNWSE
	case resizeFrameRight | resizeFrameTop, resizeFrameLeft | resizeFrameBottom:
		return ui.CursorResizeNESW
	case resizeFrameLeft, resizeFrameRight:
		return ui.CursorResizeEW
	case resizeFrameTop, resizeFrameBottom:
		return ui.CursorResizeNS
	}
	return ui.CursorDefault
}

// resizeFrameDrag moves the grabbed sides to the pointer, keeping the frame
// inside the stage and above its least size, then snaps the ratio.
func resizeFrameDrag(st *resizeFrameState, x, y float32, snap bool) {
	g := st.grab
	switch {
	case g&resizeFrameLeft != 0:
		st.l = clamp32(x-st.offX, resizeFrameMargin, st.r-resizeFrameMinW)
	case g&resizeFrameRight != 0:
		st.r = clamp32(x-st.offX, st.l+resizeFrameMinW, resizeFrameW-resizeFrameMargin)
	}
	switch {
	case g&resizeFrameTop != 0:
		st.t = clamp32(y-st.offY, resizeFrameMargin, st.b-resizeFrameMinH)
	case g&resizeFrameBottom != 0:
		st.b = clamp32(y-st.offY, st.t+resizeFrameMinH, resizeFrameH-resizeFrameMargin)
	}
	st.ratio = ""
	if !snap {
		return
	}
	w, h := st.r-st.l, st.b-st.t
	for _, rt := range resizeFrameRatios {
		want := rt.w / rt.h
		if float32(math.Abs(float64(w/h/want-1))) > 0.035 {
			continue
		}
		// An edge drag fits the side it moves; a corner fits the height.
		horizontalOnly := g&(resizeFrameTop|resizeFrameBottom) == 0
		switch {
		case horizontalOnly && g&resizeFrameLeft != 0:
			st.l = st.r - h*want
		case horizontalOnly:
			st.r = st.l + h*want
		case g&resizeFrameTop != 0:
			st.t = st.b - w/want
		default:
			st.b = st.t + w/want
		}
		st.ratio = rt.name
		return
	}
}

// ResizeFrame is a frame with eight handles: pulling an edge or corner
// resizes it, the live width and height show beside the moving sides, and
// near a common ratio it snaps to it.
func ResizeFrame(c *ui.Context, opts ...ResizeFrameOptions) ui.Element {
	var o ResizeFrameOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Snap == "" {
		o.Snap = ResizeFrameSnapOn
	}
	corner := clamp32(o.Corner, 0, 24)
	t := ThemeFrom(c)

	stage := ui.Box(c).Size(resizeFrameW, resizeFrameH)
	st := ui.Local(stage, "resize-frame", func() resizeFrameState {
		return resizeFrameState{l: 70, t: 65, r: 370, b: 265}
	})
	px, py, over := stage.PointerPosition()
	_, _, held := stage.Dragged()
	switch {
	case held && !st.dragging:
		st.dragging = true
		st.grab = resizeFrameHit(st, px, py)
		st.offX, st.offY = 0, 0
		if st.grab&resizeFrameLeft != 0 {
			st.offX = px - st.l
		} else if st.grab&resizeFrameRight != 0 {
			st.offX = px - st.r
		}
		if st.grab&resizeFrameTop != 0 {
			st.offY = py - st.t
		} else if st.grab&resizeFrameBottom != 0 {
			st.offY = py - st.b
		}
	case held && st.grab != 0:
		resizeFrameDrag(st, px, py, o.Snap == ResizeFrameSnapOn)
	case !held:
		st.dragging, st.grab, st.ratio = false, 0, ""
	}
	hover := st.grab
	if hover == 0 && over {
		hover = resizeFrameHit(st, px, py)
	}
	stage.Cursor(resizeFrameCursor(hover))

	sp := ui.SpringOptions{Duration: 110 * time.Millisecond}
	l := stage.Spring("l", st.l, sp)
	tp := stage.Spring("t", st.t, sp)
	r := stage.Spring("r", st.r, sp)
	b := stage.Spring("b", st.b, sp)
	show := stage.Spring("show", b2f(st.grab != 0), ui.SpringOptions{Duration: 180 * time.Millisecond})
	handles := []resizeFrameSide{
		resizeFrameLeft | resizeFrameTop, resizeFrameTop, resizeFrameRight | resizeFrameTop, resizeFrameRight,
		resizeFrameRight | resizeFrameBottom, resizeFrameBottom, resizeFrameLeft | resizeFrameBottom, resizeFrameLeft,
	}
	grow := make([]float32, len(handles))
	for i, s := range handles {
		grow[i] = stage.Spring(fmt.Sprint("h", s), b2f(hover == s), ui.SpringOptions{Duration: 160 * time.Millisecond, Bounce: 0.3})
	}

	stage.Draw(func(p *ui.Painter, box ui.Rect) {
		l, tp, r, b := box.X+l, box.Y+tp, box.X+r, box.Y+b
		frame := ui.Rect{X: l, Y: tp, W: r - l, H: b - tp}
		p.Fill(frame, t.Text.Alpha(0.045), corner)
		p.Stroke(frame, t.Text, corner, 1)
		for i, s := range handles {
			x, y := (l+r)/2, (tp+b)/2
			if s&resizeFrameLeft != 0 {
				x = l
			} else if s&resizeFrameRight != 0 {
				x = r
			}
			if s&resizeFrameTop != 0 {
				y = tp
			} else if s&resizeFrameBottom != 0 {
				y = b
			}
			if s&(resizeFrameLeft|resizeFrameRight) != 0 && s&(resizeFrameTop|resizeFrameBottom) != 0 {
				// Corner handles sit on the curve when the frame is rounded.
				in := corner * (1 - 1/math.Sqrt2)
				if s&resizeFrameLeft != 0 {
					x += in
				} else {
					x -= in
				}
				if s&resizeFrameTop != 0 {
					y += in
				} else {
					y -= in
				}
			}
			n := resizeFrameHandle + 3*grow[i]
			h := ui.Rect{X: x - n/2, Y: y - n/2, W: n, H: n}
			fill := t.Surface
			if st.grab == s {
				fill = t.Text
			}
			p.Fill(h, fill, 2.5)
			p.Stroke(h, t.Text, 2.5, 1)
		}
	})
	stage.Children(func() {
		if show < 0.01 {
			return
		}
		pill := func(key, label string, x, y float32) {
			ui.Box(c).Key(key).Absolute().Left(x).Top(y).Size(56, 20).Radius(6).Background(t.Text).
				Opacity(show).Center().Children(func() {
				ui.Text(c, label).FontSize(11).FontWeight(500).TextColor(t.Surface).SingleLine()
			})
		}
		w, h := st.r-st.l, st.b-st.t
		midY, midX := (tp+b)/2-10, (l+r)/2-28
		if st.grab&resizeFrameLeft != 0 {
			pill("w", fmt.Sprintf("W %.0f", w), l-56-10, midY)
		} else if st.grab&resizeFrameRight != 0 {
			pill("w", fmt.Sprintf("W %.0f", w), r+10, midY)
		}
		if st.grab&resizeFrameTop != 0 {
			pill("h", fmt.Sprintf("H %.0f", h), midX, tp-20-10)
		} else if st.grab&resizeFrameBottom != 0 {
			pill("h", fmt.Sprintf("H %.0f", h), midX, b+10)
		}
		if st.ratio != "" {
			ui.Box(c).Key("ratio").Absolute().Left(l).Top(tp).Size(r-l, b-tp).Center().Opacity(show).Children(func() {
				ui.Text(c, st.ratio).FontSize(13).FontWeight(500).TextColor(t.TextMuted).SingleLine()
			})
		}
	})
	return stage
}

func init() {
	registerDemo(
		"resize-frame",
		"A resizable frame whose handles show the live width and height on the moving sides, snapping to common ratios.",
		func(c *ui.Context) { ResizeFrame(c) },
	)
}
