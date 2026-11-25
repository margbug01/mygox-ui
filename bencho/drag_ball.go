package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// DragBallOptions configures DragBall.
type DragBallOptions struct {
	Stretch float32 // how much it deforms as it moves, 0 to 70; 36 by default (use Rigid for none)
	Give    float32 // how quickly it rounds back out, 10 to 100; 50 by default
	Grip    float32 // how hard it squeezes when held, 0 to 100; 50 by default (use Loose for none)
	Size    float32 // the ball, 32 to 84 px; 56 by default
	Rigid   bool    // no stretch, as Stretch 0
	Loose   bool    // no squeeze, as Grip 0
}

const (
	dragBallW = 300
	dragBallH = 200
)

type dragBallState struct {
	held         bool
	offX, offY   float32
	lastX, lastY float32
	vx, vy       float32
	pull         float32 // the stretch a frame ago
	last         time.Time
}

// DragBall is a soft white ball. Pointing at it swells it; holding it
// squeezes it; dragging it about stretches it along the way it goes, and
// let go it springs home and rounds back out.
func DragBall(c *ui.Context, opts ...DragBallOptions) ui.Element {
	var o DragBallOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	def := func(v *float32, d, lo, hi float32, zero bool) {
		if *v == 0 {
			*v = d
		}
		if zero {
			*v = 0
		}
		*v = clamp32(*v, lo, hi)
	}
	def(&o.Stretch, 36, 0, 70, o.Rigid)
	def(&o.Give, 50, 10, 100, false)
	def(&o.Grip, 50, 0, 100, o.Loose)
	def(&o.Size, 56, 32, 84, false)
	t := ThemeFrom(c)
	box := ui.Box(c).Size(dragBallW, dragBallH)
	st := ui.Local(box, "drag_ball", func() dragBallState { return dragBallState{lastX: dragBallW / 2, lastY: dragBallH / 2} })
	now := c.Now()
	r := o.Size / 2

	px, py, over := box.PointerPosition()
	_, _, pressed := box.Dragged()
	onBall := over && math.Hypot(float64(px-st.lastX), float64(py-st.lastY)) <= float64(r+4)
	if pressed && !st.held && onBall {
		st.held, st.offX, st.offY = true, px-st.lastX, py-st.lastY
	}
	if !pressed {
		st.held = false
	}
	if onBall || st.held {
		box.Cursor(ui.CursorGrab)
	}
	tx, ty := float32(dragBallW/2), float32(dragBallH/2)
	follow := ui.SpringOptions{Duration: 550 * time.Millisecond, Bounce: 0.45}
	if st.held {
		box.Cursor(ui.CursorGrabbing)
		tx, ty = clamp32(px-st.offX, r, dragBallW-r), clamp32(py-st.offY, r, dragBallH-r)
		follow = ui.SpringOptions{Duration: 140 * time.Millisecond, Bounce: 0.1}
	}
	x, y := box.Spring("x", tx, follow), box.Spring("y", ty, follow)
	// How fast it goes, from where it was a frame ago.
	if dt := float32(now.Sub(st.last).Seconds()); !st.last.IsZero() && dt > 0 {
		st.vx, st.vy = (x-st.lastX)/dt, (y-st.lastY)/dt
	}
	st.last, st.lastX, st.lastY = now, x, y

	// The shape: swelling under the pointer, squeezed in the hand,
	// stretched along the way it moves; it rounds back as fast as its give.
	shape := ui.SpringOptions{Duration: time.Duration(float32(time.Second) * (0.95 - 0.75*o.Give/100)), Bounce: 0.45}
	swell := box.Spring("swell", b2f(onBall && !st.held), shape)
	grip := box.Spring("grip", b2f(st.held)*o.Grip/100, shape)
	speed := float32(math.Hypot(float64(st.vx), float64(st.vy)))
	// It takes up a stretch at once and rounds back as its give lets it.
	want := min(speed/700, 1) * o.Stretch / 100
	grow := shape
	if want > st.pull {
		grow = ui.SpringOptions{Duration: 160 * time.Millisecond}
	}
	pull := box.Spring("pull", want, grow)
	st.pull = pull
	angle := float32(math.Atan2(float64(st.vy), float64(st.vx)) * 180 / math.Pi)
	if speed < 5 {
		angle = 0
	}
	sx := (1 + 0.07*swell) * (1 - 0.12*grip) * (1 + 0.9*pull)
	sy := (1 + 0.07*swell) * (1 + 0.08*grip) / (1 + 0.45*pull)
	if speed > 1 || st.held {
		c.AnimationFrame()
	}

	ball := t.Surface
	if t.Dark {
		ball = ui.Hex("#f2f2f2")
	}
	box.Draw(func(p *ui.Painter, rr ui.Rect) {
		cx, cy := rr.X+x, rr.Y+y
		tr := ui.Rotation(-angle).Then(ui.Scaling(sx, sy)).Then(ui.Rotation(angle))
		// The squeeze is about the upright axis, so it stays upright
		// whichever way the ball last moved.
		if pull < 0.01 {
			tr = ui.Scaling(sx, sy)
		}
		p.Transform(tr, cx, cy, func() {
			p.Fill(ui.Rect{X: cx - r, Y: cy - r, W: 2 * r, H: 2 * r}, ball, r)
		})
	})
	return box
}

func init() {
	registerDemo(
		"drag-ball",
		"A soft ball to drag about: it swells under the pointer, squeezes when held and stretches as it moves.",
		func(c *ui.Context) { DragBall(c) },
	)
}
