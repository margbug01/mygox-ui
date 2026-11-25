package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// LineSphereOptions configures LineSphere.
type LineSphereOptions struct {
	Lines int     // strokes in each shape, 300 to 1200; 720 by default
	Reach float32 // how wide a field the pointer pushes, as a share of the frame, 0.3 to 1.4; 0.8 by default
	Sway  float32 // how far the sphere's tips breathe in and out, 0 to 30; 3 by default
}

type lineSphereState struct {
	shape, from int
	at          time.Time
	born        time.Time
}

const (
	lineSphereW   = 340
	lineSphereH   = 320
	lineSphereR   = 112
	lineSphereSeg = 7 // points along each stroke
)

// lineSphereStroke gives stroke i of n in shape s at time now, as points on
// or about the unit sphere (x right, y down, z toward the viewer).
func lineSphereStroke(s, i, n int, now float64, sway float32, out *[lineSphereSeg][3]float64) {
	// A golden spiral spreads the strokes evenly.
	g := math.Pi * (3 - math.Sqrt(5))
	y := 1 - 2*(float64(i)+0.5)/float64(n)
	rr := math.Sqrt(1 - y*y)
	th := g * float64(i)
	px, pz := math.Cos(th)*rr, math.Sin(th)*rr
	h := math.Mod(float64(i)*0.61803, 1)
	switch s {
	case 0: // globe: arcs flowing from a pole over the surface
		// The pole sits up and to the left, toward the viewer.
		pole := [3]float64{-0.7, -0.45, 0.55}
		// Arc from the pole toward (px, y, pz), along the great circle.
		d := pole[0]*px + pole[1]*y + pole[2]*pz
		tx, ty, tz := px-d*pole[0], y-d*pole[1], pz-d*pole[2]
		tl := math.Sqrt(tx*tx+ty*ty+tz*tz) + 1e-9
		tx, ty, tz = tx/tl, ty/tl, tz/tl
		length := (0.25 + 0.6*h) * math.Pi
		for k := 0; k < lineSphereSeg; k++ {
			a := length * float64(k) / float64(lineSphereSeg-1)
			c, sn := math.Cos(a), math.Sin(a)
			out[k] = [3]float64{c*pole[0] + sn*tx, c*pole[1] + sn*ty, c*pole[2] + sn*tz}
		}
	case 1: // sphere: spikes standing out from the surface, breathing
		breathe := 1 + float64(sway)/100*math.Sin(now*1.6+h*6)
		for k := 0; k < lineSphereSeg; k++ {
			r := 0.82 + (0.18+0.12*h)*breathe*float64(k)/float64(lineSphereSeg-1)
			out[k] = [3]float64{px * r, y * r, pz * r}
		}
	default: // ring: short strokes round a flat, tilted ring
		a := th
		rad := 0.62 + 0.38*h
		for k := 0; k < lineSphereSeg; k++ {
			b := a + 0.22*float64(k)/float64(lineSphereSeg-1)
			x, z := math.Cos(b)*rad, math.Sin(b)*rad
			// Tilt the ring toward the viewer.
			out[k] = [3]float64{x, z * 0.42, z * 0.9}
		}
	}
}

// LineSphere is a globe of fine lines that turns slowly. A press flows it
// into the next shape (globe, spiky sphere, flat ring), and the pointer parts
// the lines it passes through.
func LineSphere(c *ui.Context, opts ...LineSphereOptions) ui.Element {
	var o LineSphereOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Lines == 0 {
		o.Lines = 720
	}
	o.Lines = min(max(o.Lines, 300), 1200)
	if o.Reach == 0 {
		o.Reach = 0.8
	}
	if o.Sway == 0 {
		o.Sway = 3
	}
	t := ThemeFrom(c)
	area := ui.Box(c).Size(lineSphereW, lineSphereH).Cursor(ui.CursorPointer)
	st := ui.Local(area, "line_sphere", func() lineSphereState { return lineSphereState{born: c.Now()} })
	area.OnClick(func() {
		st.from, st.shape, st.at = st.shape, (st.shape+1)%3, c.Now()
	})
	c.AnimationFrame()
	now := phaseSeconds(c.Now())
	morph := float64(1)
	if !st.at.IsZero() {
		morph = math.Min(1, c.Now().Sub(st.at).Seconds()/1.2)
		morph = morph * morph * (3 - 2*morph)
	}
	px, py, over := area.PointerPosition()
	reach := float64(clamp32(o.Reach, 0.3, 1.4)) * lineSphereW / 2.6
	push := area.Spring("push", b2f(over), ui.SpringOptions{Duration: 400 * time.Millisecond})

	area.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy := float64(r.X+r.W/2), float64(r.Y+r.H/2)
		spin := c.Now().Sub(st.born).Seconds() * 0.18
		cs, sn := math.Cos(spin), math.Sin(spin)
		var a, b [lineSphereSeg][3]float64
		pts := make([][2]float32, lineSphereSeg)
		// Strokes go into eight paths by depth, each drawn once in its own
		// shade mixed toward the paper: a few large paths, not hundreds.
		paper := t.SurfaceSubtle
		var buckets [8]ui.Path
		var heads [8]ui.Path // the last stretch of each stroke, darker toward its tip
		type tip struct{ x, y, fade float32 }
		var tips []tip
		for i := 0; i < o.Lines; i++ {
			lineSphereStroke(st.from, i, o.Lines, now, o.Sway, &a)
			lineSphereStroke(st.shape, i, o.Lines, now, o.Sway, &b)
			var tipX, tipY, depth float64
			for k := 0; k < lineSphereSeg; k++ {
				x := a[k][0] + (b[k][0]-a[k][0])*morph
				y := a[k][1] + (b[k][1]-a[k][1])*morph
				z := a[k][2] + (b[k][2]-a[k][2])*morph
				// Turn about the vertical axis.
				x, z = x*cs+z*sn, -x*sn+z*cs
				sx, sy := cx+x*lineSphereR, cy+y*lineSphereR
				// The pointer parts the lines.
				if push > 0.01 {
					dx, dy := sx-float64(r.X+px), sy-float64(r.Y+py)
					d := math.Hypot(dx, dy)
					if d < reach && d > 0.01 {
						f := math.Pow(1-d/reach, 2) * 34 * float64(push)
						sx += dx / d * f
						sy += dy / d * f
					}
				}
				pts[k] = [2]float32{float32(sx), float32(sy)}
				tipX, tipY, depth = sx, sy, z
			}
			// Nearer strokes are darker; the far side fades into the paper.
			fade := float32(0.12 + 0.55*(depth+1)/2)
			bi := min(7, max(0, int((depth+1)/2*8)))
			bk := &buckets[bi]
			bk.MoveTo(pts[0][0], pts[0][1])
			for k := 1; k < lineSphereSeg-2; k++ {
				bk.LineTo(pts[k][0], pts[k][1])
			}
			hd := &heads[bi]
			hd.MoveTo(pts[lineSphereSeg-3][0], pts[lineSphereSeg-3][1])
			hd.LineTo(pts[lineSphereSeg-2][0], pts[lineSphereSeg-2][1])
			hd.LineTo(pts[lineSphereSeg-1][0], pts[lineSphereSeg-1][1])
			if depth > -0.35 {
				tips = append(tips, tip{float32(tipX), float32(tipY), fade})
			}
		}
		for b := range buckets {
			fade := float32(0.12 + 0.55*(float32(b)+0.5)/8)
			p.StrokePath(&buckets[b], 0.55, lerpColor(paper, t.Text, fade*0.32))
			p.StrokePath(&heads[b], 0.7, lerpColor(paper, t.Text, fade*0.9))
		}
		for _, tp := range tips {
			d := 1.3 + 2*tp.fade
			p.Fill(ui.Rect{X: tp.x - d/2, Y: tp.y - d/2, W: d, H: d}, lerpColor(paper, t.Text, 0.35+tp.fade), d/2)
		}
	})
	return area
}

func init() {
	registerDemo(
		"line-sphere",
		"A turning globe of fine lines that flows into a spiky sphere and a flat ring on a press; the pointer parts the lines.",
		func(c *ui.Context) { LineSphere(c) },
	)
}
