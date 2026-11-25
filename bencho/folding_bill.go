package bencho

import (
	"math"
	"sort"
	"time"

	"github.com/egoist/mygo/ui"
)

// FoldingBillPrint is what is printed on FoldingBill.
type FoldingBillPrint string

const (
	FoldingBillPlain  FoldingBillPrint = "Plain" // the default
	FoldingBillDollar FoldingBillPrint = "Print" // the dollar print
)

// FoldingBillOptions configures FoldingBill.
type FoldingBillOptions struct {
	Print     FoldingBillPrint // "Plain" by default
	Stiffness float32          // silk at 0, card at 100; 40 by default
	Settle    float32          // how quickly it smooths out and comes home, 0 to 100
}

const (
	foldingBillW    = 340
	foldingBillH    = 240
	foldingBillCols = 30
	foldingBillRows = 10
	foldingBillBW   = 350 // the bill, flat
	foldingBillBH   = 141
	foldingBillStep = 1.0 / 120
)

type foldingBillVec [3]float32

type foldingBillState struct {
	pos, prev, rest []foldingBillVec
	links           [][3]float32 // a, b, rest length
	stiff           []float32    // each link's strength; below zero, a bending link
	grabbed         int
	last            time.Time
	released        time.Time
}

// foldingBillView turns the bill's plane to where it is seen from: turned
// away on the right, rolled a little, in perspective.
func foldingBillView(v foldingBillVec) foldingBillVec {
	const yaw, roll = 25 * math.Pi / 180, 4.3 * math.Pi / 180
	sy, cy := float32(math.Sin(yaw)), float32(math.Cos(yaw))
	x, z := v[0]*cy+v[2]*sy, -v[0]*sy+v[2]*cy
	sr, cr := float32(math.Sin(roll)), float32(math.Cos(roll))
	x, y := x*cr-v[1]*sr, x*sr+v[1]*cr
	return foldingBillVec{x, y, z}
}

// foldingBillScreen puts a viewed point on the block.
func foldingBillScreen(v foldingBillVec) (float32, float32) {
	s := 700 / (700 - v[2])
	return 178 + v[0]*s, 128 + v[1]*s
}

// foldingBillUnview turns a move seen on screen at depth z back into the
// bill's plane.
func foldingBillUnview(dx, dy, z float32) foldingBillVec {
	s := (700 - z) / 700
	dx, dy = dx*s, dy*s
	const yaw, roll = 25 * math.Pi / 180, 4.3 * math.Pi / 180
	sr, cr := float32(math.Sin(roll)), float32(math.Cos(roll))
	x, y := dx*cr+dy*sr, -dx*sr+dy*cr
	sy, cy := float32(math.Sin(yaw)), float32(math.Cos(yaw))
	return foldingBillVec{x * cy, y, x * sy}
}

func foldingBillNew() foldingBillState {
	s := foldingBillState{grabbed: -1}
	idx := func(i, j int) int { return j*(foldingBillCols+1) + i }
	for j := 0; j <= foldingBillRows; j++ {
		for i := 0; i <= foldingBillCols; i++ {
			u, v := float32(i)/foldingBillCols, float32(j)/foldingBillRows
			// Not quite flat: a soft wave along it catches the light.
			z := 14 * float32(math.Sin(math.Pow(float64(u), 1.5)*math.Pi*1.5))
			p := foldingBillVec{(u - 0.5) * foldingBillBW, (v - 0.5) * foldingBillBH, z}
			s.pos, s.prev, s.rest = append(s.pos, p), append(s.prev, p), append(s.rest, p)
		}
	}
	link := func(a, b int, class float32) {
		pa, pb := s.rest[a], s.rest[b]
		d := float32(math.Sqrt(float64((pa[0]-pb[0])*(pa[0]-pb[0]) + (pa[1]-pb[1])*(pa[1]-pb[1]) + (pa[2]-pb[2])*(pa[2]-pb[2]))))
		s.links = append(s.links, [3]float32{float32(a), float32(b), d})
		s.stiff = append(s.stiff, class)
	}
	for j := 0; j <= foldingBillRows; j++ {
		for i := 0; i <= foldingBillCols; i++ {
			if i < foldingBillCols {
				link(idx(i, j), idx(i+1, j), 1)
			}
			if j < foldingBillRows {
				link(idx(i, j), idx(i, j+1), 1)
			}
			if i < foldingBillCols && j < foldingBillRows {
				link(idx(i, j), idx(i+1, j+1), 0.6)
				link(idx(i+1, j), idx(i, j+1), 0.6)
			}
			// Bending: links over two cells, as strong as the stiffness.
			if i+2 <= foldingBillCols {
				link(idx(i, j), idx(i+2, j), -1)
			}
			if j+2 <= foldingBillRows {
				link(idx(i, j), idx(i, j+2), -1)
			}
		}
	}
	return s
}

// step runs the cloth on one tick: home pulls each point back by home.
func (s *foldingBillState) step(bend, home float32) {
	for i := range s.pos {
		p, q := s.pos[i], s.prev[i]
		n := foldingBillVec{}
		for k := 0; k < 3; k++ {
			n[k] = p[k] + (p[k]-q[k])*0.96 + (s.rest[i][k]-p[k])*home
		}
		s.prev[i], s.pos[i] = p, n
	}
	for it := 0; it < 5; it++ {
		for k, l := range s.links {
			w := s.stiff[k]
			if w < 0 {
				w = bend
			}
			a, b := int(l[0]), int(l[1])
			pa, pb := &s.pos[a], &s.pos[b]
			dx, dy, dz := pb[0]-pa[0], pb[1]-pa[1], pb[2]-pa[2]
			d := float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
			if d < 1e-4 {
				continue
			}
			k2 := (d - l[2]) / d * 0.5 * w
			if a == s.grabbed {
				pb[0], pb[1], pb[2] = pb[0]-dx*k2*2, pb[1]-dy*k2*2, pb[2]-dz*k2*2
				continue
			}
			if b == s.grabbed {
				pa[0], pa[1], pa[2] = pa[0]+dx*k2*2, pa[1]+dy*k2*2, pa[2]+dz*k2*2
				continue
			}
			pa[0], pa[1], pa[2] = pa[0]+dx*k2, pa[1]+dy*k2, pa[2]+dz*k2
			pb[0], pb[1], pb[2] = pb[0]-dx*k2, pb[1]-dy*k2, pb[2]-dz*k2
		}
	}
}

// FoldingBill is a silvery banknote hanging in the air, seen turned away a
// little. It is cloth: grab it anywhere and drag, and it folds and bunches
// after the hand; let go and it smooths out and comes home.
func FoldingBill(c *ui.Context, opts ...FoldingBillOptions) ui.Element {
	var o FoldingBillOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Print != FoldingBillDollar {
		o.Print = FoldingBillPlain
	}
	if o.Stiffness == 0 {
		o.Stiffness = 40
	}
	o.Stiffness = clamp32(o.Stiffness, 0, 100)
	o.Settle = clamp32(o.Settle, 0, 100)
	t := ThemeFrom(c)
	box := ui.Box(c).Size(foldingBillW, foldingBillH).Cursor(ui.CursorGrab)
	st := ui.Local(box, "folding_bill", foldingBillNew)
	now := c.Now()
	dt := float32(0)
	if !st.last.IsZero() {
		dt = min(float32(now.Sub(st.last).Seconds()), 0.05)
	}
	st.last = now

	px, py, _ := box.PointerPosition()
	dx, dy, pressed := box.Dragged()
	if pressed && st.grabbed < 0 {
		// Take hold of the point of the bill nearest the pointer.
		best := float32(1e9)
		for i, p := range st.pos {
			x, y := foldingBillScreen(foldingBillView(p))
			if d := (x-px)*(x-px) + (y-py)*(y-py); d < best {
				best, st.grabbed = d, i
			}
		}
		if best > 40*40 {
			st.grabbed = -2 // pressed off the bill
		}
	}
	if pressed && st.grabbed >= 0 {
		box.Cursor(ui.CursorGrabbing)
		g := &st.pos[st.grabbed]
		m := foldingBillUnview(dx, dy, foldingBillView(*g)[2])
		for k := 0; k < 3; k++ {
			g[k] += m[k]
		}
		// Carried toward the viewer a little, so it folds over the rest.
		g[2] += float32(math.Hypot(float64(dx), float64(dy))) * 0.15
		st.prev[st.grabbed] = *g
	}
	if !pressed && st.grabbed != -1 {
		if st.grabbed >= 0 {
			st.released = now
		}
		st.grabbed = -1
	}
	// Held, it barely pulls home; let go, ever more firmly as it smooths.
	home := float32(0.0025)
	if st.grabbed < 0 && !st.released.IsZero() {
		since := float32(now.Sub(st.released).Seconds())
		home = min(0.004+0.012*since*(1+o.Settle/25), 0.08)
	}
	bend := 0.02 + 0.5*o.Stiffness/100
	for s := float32(0); s < dt; s += foldingBillStep {
		st.step(bend, home)
	}
	moving := pressed
	for i := range st.pos {
		for k := 0; k < 3; k++ {
			if math.Abs(float64(st.pos[i][k]-st.rest[i][k])) > 0.05 {
				moving = true
			}
		}
	}
	if moving || dt == 0 {
		c.AnimationFrame()
	} else {
		copy(st.pos, st.rest)
		copy(st.prev, st.rest)
	}

	dark, light := ui.Hex("#98a0a5"), ui.Hex("#eef0f2")
	if o.Print == FoldingBillDollar {
		dark, light = ui.Hex("#5f7a66"), ui.Hex("#dfeadf")
	}
	if t.Dark {
		dark, light = lerpColor(dark, t.Background, 0.35), lerpColor(light, ui.Hex("#9aa0a6"), 0.25)
	}
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		type quad struct {
			pts   [4][2]float32
			depth float32
			shade float32
		}
		view := make([]foldingBillVec, len(st.pos))
		for i, q := range st.pos {
			view[i] = foldingBillView(q)
		}
		idx := func(i, j int) int { return j*(foldingBillCols+1) + i }
		quads := make([]quad, 0, foldingBillCols*foldingBillRows)
		light3 := foldingBillVec{-0.5, -0.2, 0.84} // from the front left
		hx, hy, hz := light3[0], light3[1], light3[2]+1
		hl := float32(math.Sqrt(float64(hx*hx + hy*hy + hz*hz)))
		// Each point's shade, from the facets round it, so the cloth is
		// lit smoothly rather than facet by facet.
		shade := make([]float32, len(view))
		at := func(i, j int) foldingBillVec {
			return view[idx(min(max(i, 0), foldingBillCols), min(max(j, 0), foldingBillRows))]
		}
		for j := 0; j <= foldingBillRows; j++ {
			for i := 0; i <= foldingBillCols; i++ {
				u, v := at(i+1, j), at(i-1, j)
				w, z := at(i, j+1), at(i, j-1)
				e1 := foldingBillVec{u[0] - v[0], u[1] - v[1], u[2] - v[2]}
				e2 := foldingBillVec{w[0] - z[0], w[1] - z[1], w[2] - z[2]}
				n := foldingBillVec{e1[1]*e2[2] - e1[2]*e2[1], e1[2]*e2[0] - e1[0]*e2[2], e1[0]*e2[1] - e1[1]*e2[0]}
				l := float32(math.Sqrt(float64(n[0]*n[0] + n[1]*n[1] + n[2]*n[2])))
				if l > 0 {
					n = foldingBillVec{n[0] / l, n[1] / l, n[2] / l}
				}
				diffuse := float32(math.Abs(float64(n[0]*light3[0] + n[1]*light3[1] + n[2]*light3[2])))
				// A sheen where the cloth turns the light toward the viewer.
				spec := float32(math.Pow(math.Abs(float64((n[0]*hx+n[1]*hy+n[2]*hz)/hl)), 24))
				shade[idx(i, j)] = clamp32((diffuse-0.33)/0.36+spec*0.3, 0, 1)
			}
		}
		for j := 0; j < foldingBillRows; j++ {
			for i := 0; i < foldingBillCols; i++ {
				var q quad
				for k, n := range [4]int{idx(i, j), idx(i+1, j), idx(i+1, j+1), idx(i, j+1)} {
					x, y := foldingBillScreen(view[n])
					q.pts[k] = [2]float32{r.X + x, r.Y + y}
					q.depth += view[n][2] / 4
					q.shade += shade[n] / 4
				}
				quads = append(quads, q)
			}
		}
		sort.SliceStable(quads, func(a, b int) bool { return quads[a].depth < quads[b].depth })
		// Far to near; runs of facets of one shade go in one path.
		const levels = 28
		var path ui.Path
		run := -1
		flush := func() {
			if run >= 0 {
				col := lerpColor(dark, light, float32(run)/(levels-1))
				p.FillPath(&path, col)
				p.StrokePath(&path, 0.7, col) // closes the hairline seams between facets
			}
			path = ui.Path{}
		}
		for _, q := range quads {
			lv := min(int(q.shade*levels), levels-1)
			if lv != run {
				flush()
				run = lv
			}
			path.MoveTo(q.pts[0][0], q.pts[0][1])
			for k := 1; k < 4; k++ {
				path.LineTo(q.pts[k][0], q.pts[k][1])
			}
			path.Close()
		}
		flush()
	})
	return box
}

func init() {
	registerDemo(
		"folding-bill",
		"A cloth banknote: grab it, drag it and fold it, then let it smooth out and come home.",
		func(c *ui.Context) { FoldingBill(c) },
	)
}
