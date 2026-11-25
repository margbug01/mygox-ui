package bencho

import (
	"image"
	"math"
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
	foldingBillW     = 340
	foldingBillH     = 240
	foldingBillCols  = 30
	foldingBillRows  = 10
	foldingBillBW    = 350 // the bill, flat
	foldingBillBH    = 141
	foldingBillStep  = 1.0 / 120
	foldingBillScale = 2   // pixels per point in the bill's picture
	foldingBillRound = 7   // the bill's corner radius
	foldingBillRim   = 2.2 // the edge, darker, as if it had some thickness
)

type foldingBillVec [3]float32

type foldingBillState struct {
	pos, prev, rest []foldingBillVec
	links           [][3]float32 // a, b, rest length
	stiff           []float32    // each link's strength; below zero, a bending link
	grabbed         int
	last            time.Time
	released        time.Time

	img    *image.RGBA // the bill drawn at foldingBillScale, redrawn only while it moves
	bitmap *ui.Bitmap
	zbuf   []float32
	dark   bool
	print  FoldingBillPrint
	atRest bool // the picture shows the bill at rest
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
			z := 12 * float32(math.Sin(math.Pow(float64(u), 1.35)*math.Pi*1.85))
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

	dark, light := ui.Hex("#9aa2a8"), ui.Hex("#eef0f2")
	if o.Print == FoldingBillDollar {
		dark, light = ui.Hex("#5f7a66"), ui.Hex("#dfeadf")
	}
	if t.Dark {
		dark, light = lerpColor(dark, t.Background, 0.35), lerpColor(light, ui.Hex("#9aa0a6"), 0.25)
	}
	if st.img == nil || moving || !st.atRest || st.dark != t.Dark || st.print != o.Print {
		foldingBillRender(st, dark, light)
		st.dark, st.print, st.atRest = t.Dark, o.Print, !moving
	}
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		p.ImageRegion(st.bitmap, ui.Rect{}, r, ui.ImageOptions{})
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

// foldingBillRender draws the cloth into st.img: each point is lit from the
// facets round it, and the light is blended smoothly across each facet, so
// the bill shades like a sheet rather than a set of tiles; a depth buffer
// keeps folds that come toward the viewer in front, and the bill's outline
// is a rounded rectangle with a slightly darker rim.
func foldingBillRender(st *foldingBillState, dark, light ui.Color) {
	const sc = foldingBillScale
	w, h := foldingBillW*sc, foldingBillH*sc
	if st.img == nil {
		st.img = image.NewRGBA(image.Rect(0, 0, w, h))
		st.zbuf = make([]float32, w*h)
	}
	pix := st.img.Pix
	for i := range pix {
		pix[i] = 0
	}
	for i := range st.zbuf {
		st.zbuf[i] = -1e9
	}
	idx := func(i, j int) int { return j*(foldingBillCols+1) + i }
	view := make([]foldingBillVec, len(st.pos))
	sx := make([]float32, len(st.pos))
	sy := make([]float32, len(st.pos))
	for i, q := range st.pos {
		view[i] = foldingBillView(q)
		x, y := foldingBillScreen(view[i])
		sx[i], sy[i] = x*sc, y*sc
	}
	light3 := foldingBillVec{-0.5, -0.2, 0.84} // from the front left
	hx, hy, hz := light3[0], light3[1], light3[2]+1
	hl := float32(math.Sqrt(float64(hx*hx + hy*hy + hz*hz)))
	at := func(i, j int) foldingBillVec {
		return view[idx(min(max(i, 0), foldingBillCols), min(max(j, 0), foldingBillRows))]
	}
	shade := make([]float32, len(view))
	for j := 0; j <= foldingBillRows; j++ {
		for i := 0; i <= foldingBillCols; i++ {
			u, v := at(i+1, j), at(i-1, j)
			a, b := at(i, j+1), at(i, j-1)
			e1 := foldingBillVec{u[0] - v[0], u[1] - v[1], u[2] - v[2]}
			e2 := foldingBillVec{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
			n := foldingBillVec{e1[1]*e2[2] - e1[2]*e2[1], e1[2]*e2[0] - e1[0]*e2[2], e1[0]*e2[1] - e1[1]*e2[0]}
			if l := float32(math.Sqrt(float64(n[0]*n[0] + n[1]*n[1] + n[2]*n[2]))); l > 0 {
				n = foldingBillVec{n[0] / l, n[1] / l, n[2] / l}
			}
			diffuse := float32(math.Abs(float64(n[0]*light3[0] + n[1]*light3[1] + n[2]*light3[2])))
			spec := float32(math.Pow(math.Abs(float64((n[0]*hx+n[1]*hy+n[2]*hz)/hl)), 24))
			shade[idx(i, j)] = clamp32((diffuse-0.33)/0.36+spec*0.3, 0, 1)
		}
	}
	dr, dg, db := float32(dark.R), float32(dark.G), float32(dark.B)
	lr, lg, lb := float32(light.R), float32(light.G), float32(light.B)
	// The rounded outline, in the bill's own measure: below zero inside.
	edge := func(u, v float32) float32 {
		qx := float32(math.Abs(float64(u-0.5)))*foldingBillBW - (foldingBillBW/2 - foldingBillRound)
		qy := float32(math.Abs(float64(v-0.5)))*foldingBillBH - (foldingBillBH/2 - foldingBillRound)
		ox, oy := max(qx, 0), max(qy, 0)
		return float32(math.Sqrt(float64(ox*ox+oy*oy))) + min(max(qx, qy), 0) - foldingBillRound
	}
	tri := func(a, b, c int, ua, va, ub, vb, uc, vc float32) {
		ax, ay, bx, by, cx, cy := sx[a], sy[a], sx[b], sy[b], sx[c], sy[c]
		area := (bx-ax)*(cy-ay) - (by-ay)*(cx-ax)
		if math.Abs(float64(area)) < 1e-3 {
			return
		}
		x0, x1 := int(max(min(ax, bx, cx), 0)), int(min(max(ax, bx, cx)+1, float32(w-1)))
		y0, y1 := int(max(min(ay, by, cy), 0)), int(min(max(ay, by, cy)+1, float32(h-1)))
		for y := y0; y <= y1; y++ {
			py := float32(y) + 0.5
			for x := x0; x <= x1; x++ {
				px := float32(x) + 0.5
				wa := ((bx-px)*(cy-py) - (by-py)*(cx-px)) / area
				wb := ((cx-px)*(ay-py) - (cy-py)*(ax-px)) / area
				wc := 1 - wa - wb
				if wa < -0.002 || wb < -0.002 || wc < -0.002 {
					continue
				}
				z := wa*view[a][2] + wb*view[b][2] + wc*view[c][2]
				k := y*w + x
				if z <= st.zbuf[k] {
					continue
				}
				u := wa*ua + wb*ub + wc*uc
				v := wa*va + wb*vb + wc*vc
				d := edge(u, v) // bill points; about sc pixels each
				cover := clamp32(-d*sc, 0, 1)
				if cover <= 0 {
					continue
				}
				st.zbuf[k] = z
				sh := wa*shade[a] + wb*shade[b] + wc*shade[c]
				if rim := -d; rim < foldingBillRim {
					sh *= 0.72 + 0.28*rim/foldingBillRim
				}
				o := k * 4
				pix[o] = uint8((dr + (lr-dr)*sh) * cover)
				pix[o+1] = uint8((dg + (lg-dg)*sh) * cover)
				pix[o+2] = uint8((db + (lb-db)*sh) * cover)
				pix[o+3] = uint8(255 * cover)
			}
		}
	}
	for j := 0; j < foldingBillRows; j++ {
		for i := 0; i < foldingBillCols; i++ {
			u0, u1 := float32(i)/foldingBillCols, float32(i+1)/foldingBillCols
			v0, v1 := float32(j)/foldingBillRows, float32(j+1)/foldingBillRows
			a, b, c, d := idx(i, j), idx(i+1, j), idx(i+1, j+1), idx(i, j+1)
			tri(a, b, c, u0, v0, u1, v0, u1, v1)
			tri(a, c, d, u0, v0, u1, v1, u0, v1)
		}
	}
	if st.bitmap == nil {
		st.bitmap = ui.NewBitmap(st.img)
	} else {
		st.bitmap.Update(st.img)
	}
}
