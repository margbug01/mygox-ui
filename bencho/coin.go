package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// CoinShape is the outline of Coin.
type CoinShape string

const (
	CoinHex   CoinShape = "Hex" // the default
	CoinRound CoinShape = "Round"
)

// CoinMetal is what Coin is struck in.
type CoinMetal string

const (
	CoinSilver   CoinMetal = "Silver" // the default
	CoinGraphite CoinMetal = "Graphite"
)

// CoinOptions configures Coin.
type CoinOptions struct {
	Shape CoinShape // "Hex" by default
	Metal CoinMetal // "Silver" by default
	Float float32   // how much it bobs while it hangs there, 0 to 100; 50 by default (use Still for none)
	Still bool      // no bobbing, as Float 0 on bencho.dev
}

const (
	coinW     = 300
	coinH     = 300
	coinR     = 104 // the hexagon's corner radius from the centre
	coinRim   = 22  // the raised band round the face
	coinThick = 14
	coinLean  = 14 // degrees it leans toward the pointer at most
)

type coinState struct {
	turn    float64 // degrees spun about the upright axis
	spin    float64 // degrees a second
	last    time.Time
	settled bool
}

// coinMetals are a metal's face, rim, rim grooves, recess line and edge.
var coinMetals = map[CoinMetal][5]ui.Color{
	CoinSilver:   {ui.Hex("#bcbebf"), ui.Hex("#c2c3c6"), ui.Hex("#a3a4a8"), ui.Hex("#16181b"), ui.Hex("#45474b")},
	CoinGraphite: {ui.Hex("#4a4c51"), ui.Hex("#55575c"), ui.Hex("#3a3c40"), ui.Hex("#0b0c0e"), ui.Hex("#2a2b2f")},
}

// coinOutline is the coin's outline round its centre, r from it, as
// points: a hexagon standing on a corner with rounded corners, or a circle.
func coinOutline(shape CoinShape, r float32) [][2]float32 {
	var out [][2]float32
	if shape == CoinRound {
		for i := 0; i < 72; i++ {
			a := float64(i) * 2 * math.Pi / 72
			out = append(out, [2]float32{r * 0.93 * float32(math.Cos(a)), r * 0.93 * float32(math.Sin(a))})
		}
		return out
	}
	const soft = 0.2 // how much of each side the rounding takes
	var pts [6][2]float32
	for i := range pts {
		a := float64(i)*math.Pi/3 - math.Pi/2
		pts[i] = [2]float32{r * float32(math.Cos(a)), r * float32(math.Sin(a))}
	}
	for i := range pts {
		q, a, b := pts[i], pts[(i+5)%6], pts[(i+1)%6]
		in := [2]float32{q[0] + (a[0]-q[0])*soft, q[1] + (a[1]-q[1])*soft}
		out2 := [2]float32{q[0] + (b[0]-q[0])*soft, q[1] + (b[1]-q[1])*soft}
		for k := 0; k <= 8; k++ { // the rounded corner, as a quadratic curve
			u := float32(k) / 8
			out = append(out, [2]float32{
				(1-u)*(1-u)*in[0] + 2*(1-u)*u*q[0] + u*u*out2[0],
				(1-u)*(1-u)*in[1] + 2*(1-u)*u*q[1] + u*u*out2[1],
			})
		}
	}
	return out
}

// coinView turns, leans and sees in perspective a point of the coin, from
// its centre, onto the screen round (cx, cy).
type coinView struct {
	cy, sy, cx2, sx float32 // the turn's and the lean's cosine and sine
	ox, oy          float32
}

func newCoinView(yaw, lean, ox, oy float32) coinView {
	sy, cy := math.Sincos(float64(yaw) * math.Pi / 180)
	sx, cx := math.Sincos(float64(lean) * math.Pi / 180)
	return coinView{float32(cy), float32(sy), float32(cx), float32(sx), ox, oy}
}

func (v coinView) at(x, y, z float32) (float32, float32) {
	x, z = x*v.cy+z*v.sy, -x*v.sy+z*v.cy  // turned about the upright axis
	y, z = y*v.cx2-z*v.sx, y*v.sx+z*v.cx2 // leaned about the level one
	s := 700 / (700 - z)
	return v.ox + x*s, v.oy + y*s
}

func (v coinView) path(pts [][2]float32, z float32) *ui.Path {
	var p ui.Path
	for i, q := range pts {
		x, y := v.at(q[0], q[1], z)
		if i == 0 {
			p.MoveTo(x, y)
		} else {
			p.LineTo(x, y)
		}
	}
	return p.Close()
}

// Coin is a metal coin hanging in the air, bobbing a little and leaning
// toward the pointer. It can be grabbed and spun, or flung, when it keeps
// turning and slows to rest face on.
func Coin(c *ui.Context, opts ...CoinOptions) ui.Element {
	var o CoinOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Shape != CoinRound {
		o.Shape = CoinHex
	}
	if o.Metal != CoinGraphite {
		o.Metal = CoinSilver
	}
	if o.Float == 0 {
		o.Float = 50
	}
	if o.Still {
		o.Float = 0
	}
	o.Float = clamp32(o.Float, 0, 100)
	box := ui.Box(c).Size(coinW, coinH).Cursor(ui.CursorGrab)
	st := ui.Local(box, "coin", func() coinState { return coinState{settled: true} })
	now := c.Now()
	dt := 0.0
	if !st.last.IsZero() {
		dt = min(now.Sub(st.last).Seconds(), 0.05)
	}
	st.last = now

	px, py, over := box.PointerPosition()
	dx, _, pressed := box.Dragged()
	if pressed {
		box.Cursor(ui.CursorGrabbing)
		st.turn += float64(dx) * 0.9
		vx, _ := box.DragVelocity()
		st.spin, st.settled = float64(vx)*0.9, false
	} else if !st.settled {
		// Flung, it turns on, slowing, then comes to rest face on.
		st.turn += st.spin * dt
		st.spin *= math.Exp(-1.1 * dt)
		if math.Abs(st.spin) < 90 {
			st.settled = true
		}
	}
	turn := float32(st.turn)
	if st.settled && !pressed {
		face := float32(math.Round(st.turn/180) * 180)
		turn = box.Spring("settle", face, ui.SpringOptions{Duration: 700 * time.Millisecond, Bounce: 0.35})
	} else {
		box.Spring("settle", turn, ui.SpringOptions{Duration: 1})
	}

	// It leans toward the pointer, and away from nothing: at rest it is
	// turned a little to the right, showing its left edge.
	lx, ly := float32(8), float32(16)
	if over && !pressed {
		lx = -clamp32((py-coinH/2)/(coinH/2), -1, 1) * coinLean
		ly = clamp32((px-coinW/2)/(coinW/2), -1, 1) * coinLean
	}
	lean := ui.SpringOptions{Duration: 500 * time.Millisecond}
	leanX, leanY := box.Spring("leanx", lx, lean), box.Spring("leany", ly, lean)
	bob := float32(math.Sin(phaseSeconds(now)*1.7)) * 5 * o.Float / 100
	if o.Float > 0 || !st.settled || pressed {
		c.AnimationFrame()
	}

	m := coinMetals[o.Metal]
	outer, inner := coinOutline(o.Shape, coinR), coinOutline(o.Shape, coinR-coinRim)
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		yaw := float32(math.Mod(float64(turn+leanY), 360))
		v := newCoinView(yaw, leanX, r.X+coinW/2, r.Y+coinH/2+bob)
		// Which side faces the viewer.
		face := float32(coinThick / 2)
		if math.Cos(float64(yaw)*math.Pi/180)*math.Cos(float64(leanX)*math.Pi/180) < 0 {
			face = -face
		}
		// The milled edge, as a stack of slices, far to near, lighter
		// toward the face.
		const slices = 7
		for i := 0; i < slices; i++ {
			k := float32(i) / slices
			p.FillPath(v.path(outer, -face+2*face*k), lerpColor(m[4], m[2], k*0.7))
		}
		coinFace(p, v, outer, inner, face, m)
	})
	return box
}

// coinFace draws a face at depth z: the grooved rim and the flat recessed
// field.
func coinFace(p *ui.Painter, v coinView, outer, inner [][2]float32, z float32, m [5]ui.Color) {
	rim := v.path(outer, z)
	p.FillPathGradient(rim, ui.LinearGradient{From: m[1], To: lerpColor(m[1], m[2], 0.4), Angle: 160})
	// Fine grooves run across the rim, all the way round.
	var grooves ui.Path
	const n = 180
	for i := 0; i < n; i++ {
		a := float64(i) * 2 * math.Pi / n
		cs, sn := float32(math.Cos(a)), float32(math.Sin(a))
		x0, y0 := v.at(cs*(coinR-coinRim-2), sn*(coinR-coinRim-2), z)
		x1, y1 := v.at(cs*coinR, sn*coinR, z)
		grooves.MoveTo(x0, y0).LineTo(x1, y1)
	}
	p.ClipPath(rim, func() { p.StrokePath(&grooves, 0.6, m[2].Alpha(0.55)) })
	field := v.path(inner, z)
	p.FillPath(field, m[0])
	// The recess: a dark line where the field drops below the rim along
	// the top, a light one along the bottom.
	p.StrokePathGradient(field, 1.6, ui.LinearGradient{From: m[3], To: m[1], Angle: 180, Start: 0.05, End: 0.5})
}

func init() {
	registerDemo(
		"coin",
		"A metal coin hanging in the air that leans toward the pointer; grab it to spin it, or fling it and watch it slow to rest.",
		func(c *ui.Context) { Coin(c) },
	)
}
