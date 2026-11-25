package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// SlingshotOptions configures Slingshot.
type SlingshotOptions struct {
	Power  float32 // how hard the band throws, 50 to 160 (%); 100 by default
	Bounce float32 // how lively the ball comes off the walls, 0 to 100; 30 by default
}

const (
	slingshotW     = 360
	slingshotH     = 260
	slingshotPouch = 38 // the band's round end, radius
	slingshotBall  = 29
	slingshotFree  = 70 // pulled this far, the band starts to fight back
)

type slingshotState struct {
	held       bool
	pullX      float32 // the ball off the anchor while held
	pullY      float32
	x, y       float32 // the ball off the anchor once let go
	vx, vy     float32
	flying     bool
	thrown     time.Time
	last       time.Time
	angle      float32 // the arrow's heading, degrees
	trembleFor float32
}

// Slingshot is a black ball in a grey band. Pulling it back stretches the
// band, the arrow on the ball turning to where it will go, and at the
// band's limit it trembles; let go and the ball flies, bounces off the
// walls and comes back into the band.
func Slingshot(c *ui.Context, opts ...SlingshotOptions) ui.Element {
	var o SlingshotOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Power == 0 {
		o.Power = 100
	}
	o.Power = clamp32(o.Power, 50, 160)
	if o.Bounce == 0 {
		o.Bounce = 30
	}
	o.Bounce = clamp32(o.Bounce, 0, 100)
	t := ThemeFrom(c)
	box := ui.Box(c).Size(slingshotW, slingshotH)
	st := ui.Local(box, "slingshot", func() slingshotState { return slingshotState{angle: -45} })
	now := c.Now()
	dt := float32(0)
	if !st.last.IsZero() {
		dt = min(float32(now.Sub(st.last).Seconds()), 0.05)
	}
	st.last = now
	const ax, ay = slingshotW / 2, slingshotH / 2

	px, py, _ := box.PointerPosition()
	_, _, pressed := box.Dragged()
	if pressed && !st.held && !st.flying && math.Hypot(float64(px-ax), float64(py-ay)) <= slingshotPouch+4 {
		st.held = true
	}
	if st.held && pressed {
		// The band gives freely, then ever less.
		dx, dy := px-ax, py-ay
		l := float32(math.Hypot(float64(dx), float64(dy)))
		k := float32(1)
		if l > slingshotFree {
			k = (slingshotFree + rubberBand(l-slingshotFree, 50)) / l
		}
		st.pullX, st.pullY = dx*k, dy*k
		st.trembleFor = clamp32((l-slingshotFree-60)/40, 0, 1)
		if l*k > 6 {
			st.angle = float32(math.Atan2(float64(-dy), float64(-dx)) * 180 / math.Pi)
		}
	}
	if st.held && !pressed {
		st.held = false
		if math.Hypot(float64(st.pullX), float64(st.pullY)) > 12 {
			s := 9.5 * o.Power / 100
			st.x, st.y = st.pullX, st.pullY
			st.vx, st.vy = -st.pullX*s, -st.pullY*s
			st.flying, st.thrown = true, now
		}
		st.pullX, st.pullY, st.trembleFor = 0, 0, 0
	}
	if st.flying && dt > 0 {
		slingshotFly(st, dt, now, 0.3+0.65*o.Bounce/100)
	}
	animating := st.flying || st.held

	// The ball, where it is now; the pull springs back when let go short.
	bx := box.Spring("pullx", st.pullX, ui.SpringOptions{Duration: 350 * time.Millisecond, Bounce: 0.4})
	by := box.Spring("pully", st.pullY, ui.SpringOptions{Duration: 350 * time.Millisecond, Bounce: 0.4})
	if st.held {
		bx, by = st.pullX, st.pullY
		if st.trembleFor > 0 {
			ph := float64(now.UnixNano()) / 1e9 * 95
			bx += float32(math.Sin(ph)) * 1.6 * st.trembleFor
			by += float32(math.Cos(ph*1.3)) * 1.6 * st.trembleFor
		}
	}
	if st.flying {
		bx, by = st.x, st.y
	}
	angle := st.angle
	if !st.held && !st.flying {
		angle = box.Spring("angle", -45, ui.SpringOptions{Duration: 500 * time.Millisecond})
	} else {
		box.Spring("angle", angle, ui.SpringOptions{Duration: 1})
	}
	if animating {
		c.AnimationFrame()
	}

	ball, band, ink := ui.Hex("#161819"), ui.Hex("#dbdbdb"), ui.Hex("#ffffff")
	if t.Dark {
		ball, band, ink = t.Text, ui.Hex("#3a3a3e"), t.Background
	}
	if math.Hypot(float64(px-ax-bx), float64(py-ay-by)) <= slingshotBall && !st.flying {
		box.Cursor(ui.CursorGrab)
	}
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy := r.X+ax, r.Y+ay
		// The band: the pouch at the anchor drawn out to the ball, narrowing
		// as it stretches. In flight the ball has left it behind.
		stretch := float32(0)
		if !st.flying {
			stretch = float32(math.Hypot(float64(bx), float64(by)))
		}
		if stretch < 1 {
			p.Fill(ui.Rect{X: cx - slingshotPouch, Y: cy - slingshotPouch, W: 2 * slingshotPouch, H: 2 * slingshotPouch}, band, slingshotPouch)
		} else {
			rb := slingshotPouch - clamp32(stretch*0.04, 0, 6)
			slingshotBand(p, cx, cy, slingshotPouch, cx+bx, cy+by, rb, band)
		}
		bcx, bcy := cx+bx, cy+by
		p.Fill(ui.Rect{X: bcx - slingshotBall, Y: bcy - slingshotBall, W: 2 * slingshotBall, H: 2 * slingshotBall}, ball, slingshotBall)
		p.Transform(ui.Rotation(angle), bcx, bcy, func() {
			p.Icon(icons.ArrowRight(), ui.Rect{X: bcx - 15.5, Y: bcy - 15.5, W: 31, H: 31}, ink)
		})
	})
	return box
}

// slingshotFly moves the thrown ball on by dt: it slows in the air and
// bounces off the walls, and after a moment a spring draws it home.
func slingshotFly(st *slingshotState, dt float32, now time.Time, restitution float32) {
	home := float32(now.Sub(st.thrown).Seconds()) > 0.9
	for s := float32(0); s < dt; s += 1.0 / 240 {
		h := min(float32(1.0/240), dt-s)
		if home {
			// Underdamped, so it settles into the band with a little wobble.
			st.vx += (-st.x*70 - st.vx*11) * h
			st.vy += (-st.y*70 - st.vy*11) * h
		} else {
			drag := float32(math.Exp(float64(-1.1 * h)))
			st.vx, st.vy = st.vx*drag, st.vy*drag
		}
		st.x += st.vx * h
		st.y += st.vy * h
		const mx, my = slingshotW/2 - slingshotBall, slingshotH/2 - slingshotBall
		if st.x < -mx || st.x > mx {
			st.x = clamp32(st.x, -mx, mx)
			st.vx = -st.vx * restitution
		}
		if st.y < -my || st.y > my {
			st.y = clamp32(st.y, -my, my)
			st.vy = -st.vy * restitution
		}
	}
	if math.Hypot(float64(st.vx), float64(st.vy)) > 40 {
		st.angle = float32(math.Atan2(float64(st.vy), float64(st.vx)) * 180 / math.Pi)
	}
	if home && math.Hypot(float64(st.x), float64(st.y)) < 0.4 && math.Hypot(float64(st.vx), float64(st.vy)) < 6 {
		st.flying, st.x, st.y, st.vx, st.vy = false, 0, 0, 0, 0
	}
}

// slingshotBand fills the hull of two circles: a round end of radius ra at
// (ax, ay) drawn out to one of radius rb at (bx, by).
func slingshotBand(p *ui.Painter, ax, ay, ra, bx, by, rb float32, col ui.Color) {
	th := math.Atan2(float64(by-ay), float64(bx-ax))
	var path ui.Path
	const n = 24
	for i := 0; i <= n; i++ { // the anchor's far half
		a := th + math.Pi/2 + math.Pi*float64(i)/n
		x, y := ax+ra*float32(math.Cos(a)), ay+ra*float32(math.Sin(a))
		if i == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
	for i := 0; i <= n; i++ { // the ball's far half
		a := th - math.Pi/2 + math.Pi*float64(i)/n
		path.LineTo(bx+rb*float32(math.Cos(a)), by+rb*float32(math.Sin(a)))
	}
	path.Close()
	p.FillPath(&path, col)
}

func init() {
	registerDemo(
		"slingshot",
		"Pull back the elastic band and fling the ball: it flies, bounces off the walls and comes back.",
		func(c *ui.Context) { Slingshot(c) },
	)
}
