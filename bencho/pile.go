package bencho

import "math"

// pile is a small 2D physics world for Gravity tags and Falling text: each
// body is a capsule (a pill, or a disc when its spine has no length) whose two
// spine ends are Verlet particles held apart by a stick, so bodies tip, roll
// and stack without a full rigid-body solver.
type pile struct {
	bodies  []pileBody
	w, h    float32 // the box they live in; open at the top
	gravity float32 // px/s²
}

type pileBody struct {
	ax, ay, bx, by     float32 // spine ends now
	pax, pay, pbx, pby float32 // spine ends a step ago
	r                  float32 // half the pill's height
	spine              float32 // the stick's length: the pill's width less its height
	asleep             bool    // not yet dropped in
	still              bool    // settled: frozen until something knocks it
	rest               int     // steps in a row it has barely moved
}

const pileStep = 1.0 / 120

// pileBodyAt makes a w by h pill centred at (x, y), turned by angle radians.
func pileBodyAt(x, y, w, h, angle float32) pileBody {
	r := h / 2
	spine := max(0, w-h)
	cs, sn := float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))
	ax, ay := x-cs*spine/2, y-sn*spine/2
	bx, by := x+cs*spine/2, y+sn*spine/2
	return pileBody{ax: ax, ay: ay, bx: bx, by: by, pax: ax, pay: ay, pbx: bx, pby: by, r: r, spine: spine}
}

// centre and angle of a body, for drawing.
func (b *pileBody) centre() (x, y, angle float32) {
	return (b.ax + b.bx) / 2, (b.ay + b.by) / 2, float32(math.Atan2(float64(b.by-b.ay), float64(b.bx-b.ax)))
}

// push moves the body's spine ends by (dx, dy) split as at t (0 = all on a, 1 = all on b).
func (b *pileBody) push(dx, dy, t float32) {
	b.ax += dx * (1 - t)
	b.ay += dy * (1 - t)
	b.bx += dx * t
	b.by += dy * t
}

// closest finds t along segment a→b nearest to the point (px, py).
func pileClosest(ax, ay, bx, by, px, py float32) float32 {
	dx, dy := bx-ax, by-ay
	l := dx*dx + dy*dy
	if l < 1e-6 {
		return 0
	}
	return clamp32(((px-ax)*dx+(py-ay)*dy)/l, 0, 1)
}

// advance runs the world on by dt seconds in fixed steps. pointer, when its
// radius is above zero, shoves bodies out of its way, and its velocity
// (vx, vy) carries into what it touches.
func (p *pile) advance(dt float32, px, py, pr, vx, vy float32) {
	steps := int(math.Ceil(float64(dt / pileStep)))
	steps = min(steps, 12)
	for s := 0; s < steps; s++ {
		p.step(px, py, pr, vx, vy)
	}
}

func (p *pile) step(px, py, pr, vx, vy float32) {
	const h = pileStep
	g := p.gravity * h * h
	start := make([][4]float32, len(p.bodies))
	for i := range p.bodies {
		b := &p.bodies[i]
		start[i] = [4]float32{b.ax, b.ay, b.bx, b.by}
		if b.asleep || b.still {
			continue
		}
		// Bodies barely moving come to rest, so a settled pile stays still.
		va := (b.ax-b.pax)*(b.ax-b.pax) + (b.ay-b.pay)*(b.ay-b.pay)
		vb := (b.bx-b.pbx)*(b.bx-b.pbx) + (b.by-b.pby)*(b.by-b.pby)
		if va < 0.0009 && vb < 0.0009 {
			b.pax, b.pay, b.pbx, b.pby = b.ax, b.ay, b.bx, b.by
		}
		// Verlet, with a little air drag.
		nax, nay := b.ax+(b.ax-b.pax)*0.995, b.ay+(b.ay-b.pay)*0.995+g
		nbx, nby := b.bx+(b.bx-b.pbx)*0.995, b.by+(b.by-b.pby)*0.995+g
		b.pax, b.pay, b.pbx, b.pby = b.ax, b.ay, b.bx, b.by
		b.ax, b.ay, b.bx, b.by = nax, nay, nbx, nby
	}
	for it := 0; it < 6; it++ {
		for i := range p.bodies {
			b := &p.bodies[i]
			if b.asleep {
				continue
			}
			// The pointer wakes what it touches.
			if pr > 0 && b.still {
				t := pileClosest(b.ax, b.ay, b.bx, b.by, px, py)
				cx, cy := b.ax+(b.bx-b.ax)*t, b.ay+(b.by-b.ay)*t
				if float32(math.Hypot(float64(cx-px), float64(cy-py))) < pr+b.r {
					b.still, b.rest = false, 0
				}
			}
			if b.still {
				continue
			}
			// The stick.
			dx, dy := b.bx-b.ax, b.by-b.ay
			d := float32(math.Hypot(float64(dx), float64(dy)))
			if d > 1e-4 {
				k := (d - b.spine) / d / 2
				b.ax += dx * k
				b.ay += dy * k
				b.bx -= dx * k
				b.by -= dy * k
			} else if b.spine > 0 {
				b.bx += b.spine
			}
			// The floor and walls, with friction on the floor.
			for e := 0; e < 2; e++ {
				x, y, ox := &b.ax, &b.ay, &b.pax
				if e == 1 {
					x, y, ox = &b.bx, &b.by, &b.pbx
				}
				if *y > p.h-b.r {
					*y = p.h - b.r
					*ox += (*x - *ox) * 0.3
				}
				if *x < b.r {
					*x = b.r
				}
				if *x > p.w-b.r {
					*x = p.w - b.r
				}
			}
			// The pointer.
			if pr > 0 {
				t := pileClosest(b.ax, b.ay, b.bx, b.by, px, py)
				cx, cy := b.ax+(b.bx-b.ax)*t, b.ay+(b.by-b.ay)*t
				nx, ny := cx-px, cy-py
				d := float32(math.Hypot(float64(nx), float64(ny)))
				if d < pr+b.r && d > 1e-4 {
					over := (pr + b.r - d) / d
					b.push(nx*over, ny*over, t)
					// A fast sweep throws them.
					b.push(vx*h*0.6, vy*h*0.6, 0.5)
				}
			}
		}
		// Bodies against bodies.
		for i := range p.bodies {
			a := &p.bodies[i]
			if a.asleep {
				continue
			}
			for j := i + 1; j < len(p.bodies); j++ {
				b := &p.bodies[j]
				if b.asleep {
					continue
				}
				pileCollide(a, b)
			}
		}
	}
	// Bodies that have barely moved for a while settle.
	for i := range p.bodies {
		b := &p.bodies[i]
		if b.asleep || b.still {
			continue
		}
		m := float32(math.Abs(float64(b.ax-start[i][0]))) + float32(math.Abs(float64(b.ay-start[i][1]))) +
			float32(math.Abs(float64(b.bx-start[i][2]))) + float32(math.Abs(float64(b.by-start[i][3])))
		if m < 0.3 {
			b.rest++
		} else {
			b.rest = 0
		}
		if b.rest > 50 {
			b.still = true
			b.pax, b.pay, b.pbx, b.pby = b.ax, b.ay, b.bx, b.by
		}
	}
}

// pileCollide pushes two capsules apart where their spines come closer than
// their radii.
func pileCollide(a, b *pileBody) {
	if a.still && b.still {
		return
	}
	// Closest points between the two spines, by refining from each end.
	ta := pileClosest(a.ax, a.ay, a.bx, a.by, (b.ax+b.bx)/2, (b.ay+b.by)/2)
	tb := float32(0)
	for k := 0; k < 3; k++ {
		px, py := a.ax+(a.bx-a.ax)*ta, a.ay+(a.by-a.ay)*ta
		tb = pileClosest(b.ax, b.ay, b.bx, b.by, px, py)
		qx, qy := b.ax+(b.bx-b.ax)*tb, b.ay+(b.by-b.ay)*tb
		ta = pileClosest(a.ax, a.ay, a.bx, a.by, qx, qy)
	}
	px, py := a.ax+(a.bx-a.ax)*ta, a.ay+(a.by-a.ay)*ta
	qx, qy := b.ax+(b.bx-b.ax)*tb, b.ay+(b.by-b.ay)*tb
	nx, ny := qx-px, qy-py
	d := float32(math.Hypot(float64(nx), float64(ny)))
	reach := a.r + b.r
	if d >= reach {
		return
	}
	if d < 1e-4 {
		nx, ny, d = 0, 1, 1
	}
	over := (reach - d) / d / 2
	switch {
	case a.still && reach-d < 1.5: // a settled body holds; the other gives
		b.push(nx*over*2, ny*over*2, tb)
	case b.still && reach-d < 1.5:
		a.push(-nx*over*2, -ny*over*2, ta)
	default: // a hard knock wakes a settled body
		a.still, b.still = false, false
		a.push(-nx*over, -ny*over, ta)
		b.push(nx*over, ny*over, tb)
	}
	// Friction: touching bodies lose some of their motion, so they come to
	// rest leaning on each other rather than sliding flat.
	for _, q := range []*pileBody{a, b} {
		q.pax += (q.ax - q.pax) * 0.12
		q.pay += (q.ay - q.pay) * 0.12
		q.pbx += (q.bx - q.pbx) * 0.12
		q.pby += (q.by - q.pby) * 0.12
	}
}
