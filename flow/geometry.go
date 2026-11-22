package flow

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// PortPosition computes the world-space center of portID on node n.
func PortPosition(n *Node, portID string) (Point, PortSide, bool) {
	if n == nil {
		return Point{}, PortRight, false
	}
	b := n.Bounds()
	var target *Port
	for i := range n.Ports {
		if n.Ports[i].ID == portID {
			target = &n.Ports[i]
			break
		}
	}
	if target == nil {
		return Point{X: b.X + b.W, Y: b.Y + b.H*0.5}, PortRight, false
	}

	// DatabaseSchemaNode aligns handles directly with their table rows.
	if n.Variant == VariantSchema && len(n.SchemaRows) > 0 {
		const (
			schemaHeaderH = float32(34)
			schemaRowH    = float32(26)
		)
		for idx, row := range n.SchemaRows {
			rowMidY := b.Y + schemaHeaderH + float32(idx)*schemaRowH + schemaRowH*0.5 + 2
			if row.InID == portID {
				return Point{X: b.X, Y: rowMidY}, PortLeft, true
			}
			if row.OutID == portID {
				return Point{X: b.X + b.W, Y: rowMidY}, PortRight, true
			}
		}
	}

	side := target.Side
	if side == "" {
		if target.Kind == PortInput {
			side = PortLeft
		} else {
			side = PortRight
		}
	}

	// Count sibling ports on the same side to space them evenly.
	total := 0
	idx := 0
	for i := range n.Ports {
		s := n.Ports[i].Side
		if s == "" {
			if n.Ports[i].Kind == PortInput {
				s = PortLeft
			} else {
				s = PortRight
			}
		}
		if s == side {
			if n.Ports[i].ID == portID {
				idx = total
			}
			total++
		}
	}
	if total == 0 {
		total = 1
	}

	// Offset below the 34px header when laying out left/right ports.
	headerH := float32(34)
	bodyH := b.H - headerH
	if bodyH < 28 {
		headerH = 0
		bodyH = b.H
	}
	frac := float32(idx+1) / float32(total+1)

	switch side {
	case PortLeft:
		return Point{X: b.X, Y: b.Y + headerH + bodyH*frac}, PortLeft, true
	case PortRight:
		return Point{X: b.X + b.W, Y: b.Y + headerH + bodyH*frac}, PortRight, true
	case PortTop:
		return Point{X: b.X + b.W*frac, Y: b.Y}, PortTop, true
	case PortBottom:
		return Point{X: b.X + b.W*frac, Y: b.Y + b.H}, PortBottom, true
	default:
		return Point{X: b.X + b.W, Y: b.Y + b.H*0.5}, PortRight, true
	}
}

func sideNormal(side PortSide) Point {
	switch side {
	case PortLeft:
		return Point{X: -1, Y: 0}
	case PortRight:
		return Point{X: 1, Y: 0}
	case PortTop:
		return Point{X: 0, Y: -1}
	case PortBottom:
		return Point{X: 0, Y: 1}
	default:
		return Point{X: 1, Y: 0}
	}
}

// EdgeGeometry holds the computed screen-space path, midpoint label position,
// and arrow geometry for an edge.
type EdgeGeometry struct {
	Path      *ui.Path
	LabelPos  Point
	Tip       Point
	Angle     float32
	Samples   []Point
}

// BuildEdgeGeometry constructs screen-space curve geometry for an edge.
func BuildEdgeGeometry(src, dst Point, srcSide, dstSide PortSide, kind EdgeKind) EdgeGeometry {
	if kind == "" {
		kind = EdgeBezier
	}
	switch kind {
	case EdgeStraight:
		return buildStraightGeometry(src, dst)
	case EdgeSmoothStep:
		return buildSmoothStepGeometry(src, dst, srcSide, dstSide)
	default:
		return buildBezierGeometry(src, dst, srcSide, dstSide)
	}
}

func buildStraightGeometry(src, dst Point) EdgeGeometry {
	p := new(ui.Path)
	p.MoveTo(src.X, src.Y)
	p.LineTo(dst.X, dst.Y)
	mid := Point{X: (src.X + dst.X) * 0.5, Y: (src.Y + dst.Y) * 0.5}
	angle := float32(math.Atan2(float64(dst.Y-src.Y), float64(dst.X-src.X)))
	return EdgeGeometry{
		Path:     p,
		LabelPos: mid,
		Tip:      dst,
		Angle:    angle,
		Samples:  []Point{src, mid, dst},
	}
}

func buildBezierGeometry(src, dst Point, srcSide, dstSide PortSide) EdgeGeometry {
	d := dist2D(src, dst)
	offset := clamp32(d*0.42, 36, 180)
	sn := sideNormal(srcSide)
	dn := sideNormal(dstSide)
	c1 := Point{X: src.X + sn.X*offset, Y: src.Y + sn.Y*offset}
	c2 := Point{X: dst.X + dn.X*offset, Y: dst.Y + dn.Y*offset}

	p := new(ui.Path)
	p.MoveTo(src.X, src.Y)
	p.CubeTo(c1.X, c1.Y, c2.X, c2.Y, dst.X, dst.Y)

	const steps = 20
	samples := make([]Point, steps+1)
	for i := 0; i <= steps; i++ {
		t := float32(i) / float32(steps)
		samples[i] = evalCubic(src, c1, c2, dst, t)
	}
	mid := evalCubic(src, c1, c2, dst, 0.5)
	prev := samples[steps-1]
	angle := float32(math.Atan2(float64(dst.Y-prev.Y), float64(dst.X-prev.X)))
	return EdgeGeometry{
		Path:     p,
		LabelPos: mid,
		Tip:      dst,
		Angle:    angle,
		Samples:  samples,
	}
}

func evalCubic(p0, p1, p2, p3 Point, t float32) Point {
	u := 1 - t
	tt := t * t
	uu := u * u
	uuu := uu * u
	ttt := tt * t
	return Point{
		X: uuu*p0.X + 3*uu*t*p1.X + 3*u*tt*p2.X + ttt*p3.X,
		Y: uuu*p0.Y + 3*uu*t*p1.Y + 3*u*tt*p2.Y + ttt*p3.Y,
	}
}

func buildSmoothStepGeometry(src, dst Point, srcSide, dstSide PortSide) EdgeGeometry {
	sn := sideNormal(srcSide)
	dn := sideNormal(dstSide)
	stub := float32(22)
	p1 := Point{X: src.X + sn.X*stub, Y: src.Y + sn.Y*stub}
	p2 := Point{X: dst.X + dn.X*stub, Y: dst.Y + dn.Y*stub}

	pts := []Point{src, p1}
	horiz := srcSide == PortLeft || srcSide == PortRight
	if horiz {
		if (srcSide == PortRight && p2.X >= p1.X) || (srcSide == PortLeft && p2.X <= p1.X) {
			midX := (p1.X + p2.X) * 0.5
			pts = append(pts, Point{X: midX, Y: p1.Y}, Point{X: midX, Y: p2.Y})
		} else {
			midY := (p1.Y + p2.Y) * 0.5
			if math.Abs(float64(p2.Y-p1.Y)) < 48 {
				midY = p1.Y - 56
			}
			pts = append(pts, Point{X: p1.X, Y: midY}, Point{X: p2.X, Y: midY})
		}
	} else {
		midY := (p1.Y + p2.Y) * 0.5
		pts = append(pts, Point{X: p1.X, Y: midY}, Point{X: p2.X, Y: midY})
	}
	pts = append(pts, p2, dst)

	// Deduplicate adjacent equal points.
	clean := make([]Point, 0, len(pts))
	for _, pt := range pts {
		if len(clean) == 0 || dist2D(clean[len(clean)-1], pt) > 0.5 {
			clean = append(clean, pt)
		}
	}

	p := new(ui.Path)
	p.MoveTo(clean[0].X, clean[0].Y)
	const cornerR = float32(10)
	for i := 1; i < len(clean)-1; i++ {
		prev, cur, next := clean[i-1], clean[i], clean[i+1]
		d1 := dist2D(prev, cur)
		d2 := dist2D(cur, next)
		r := cornerR
		if d1*0.5 < r {
			r = d1 * 0.5
		}
		if d2*0.5 < r {
			r = d2 * 0.5
		}
		if r < 1 || d1 < 1 || d2 < 1 {
			p.LineTo(cur.X, cur.Y)
			continue
		}
		inPt := Point{
			X: cur.X + (prev.X-cur.X)*(r/d1),
			Y: cur.Y + (prev.Y-cur.Y)*(r/d1),
		}
		outPt := Point{
			X: cur.X + (next.X-cur.X)*(r/d2),
			Y: cur.Y + (next.Y-cur.Y)*(r/d2),
		}
		p.LineTo(inPt.X, inPt.Y)
		p.QuadTo(cur.X, cur.Y, outPt.X, outPt.Y)
	}
	last := clean[len(clean)-1]
	p.LineTo(last.X, last.Y)

	mid := polylineMidpoint(clean)
	prev := clean[len(clean)-2]
	angle := float32(math.Atan2(float64(last.Y-prev.Y), float64(last.X-prev.X)))
	return EdgeGeometry{
		Path:     p,
		LabelPos: mid,
		Tip:      last,
		Angle:    angle,
		Samples:  clean,
	}
}

func polylineMidpoint(pts []Point) Point {
	if len(pts) == 0 {
		return Point{}
	}
	if len(pts) == 1 {
		return pts[0]
	}
	var total float32
	for i := 1; i < len(pts); i++ {
		total += dist2D(pts[i-1], pts[i])
	}
	half := total * 0.5
	var acc float32
	for i := 1; i < len(pts); i++ {
		seg := dist2D(pts[i-1], pts[i])
		if acc+seg >= half && seg > 0 {
			t := (half - acc) / seg
			return Point{
				X: pts[i-1].X + (pts[i].X-pts[i-1].X)*t,
				Y: pts[i-1].Y + (pts[i].Y-pts[i-1].Y)*t,
			}
		}
		acc += seg
	}
	return pts[len(pts)/2]
}

// DistanceToPolyline computes the minimum distance from point pt to a sampled curve.
func DistanceToPolyline(pt Point, samples []Point) float32 {
	if len(samples) < 2 {
		return float32(math.Inf(1))
	}
	best := float32(math.Inf(1))
	for i := 1; i < len(samples); i++ {
		d := distToSegment(pt, samples[i-1], samples[i])
		if d < best {
			best = d
		}
	}
	return best
}

func distToSegment(p, a, b Point) float32 {
	vx := b.X - a.X
	vy := b.Y - a.Y
	wx := p.X - a.X
	wy := p.Y - a.Y
	lenSq := vx*vx + vy*vy
	if lenSq <= 1e-6 {
		return dist2D(p, a)
	}
	t := clamp32((wx*vx+wy*vy)/lenSq, 0, 1)
	proj := Point{X: a.X + vx*t, Y: a.Y + vy*t}
	return dist2D(p, proj)
}

// SamplePolylineAt returns the point at arc-length fraction t in [0,1] along samples.
// Used by AnimatedSvgEdge to move a glowing shape smoothly along any edge geometry.
func SamplePolylineAt(samples []Point, t float32) Point {
	if len(samples) == 0 {
		return Point{}
	}
	if len(samples) == 1 || t <= 0 {
		return samples[0]
	}
	if t >= 1 {
		return samples[len(samples)-1]
	}
	var total float32
	for i := 1; i < len(samples); i++ {
		total += dist2D(samples[i-1], samples[i])
	}
	if total <= 1e-4 {
		return samples[0]
	}
	target := total * t
	var acc float32
	for i := 1; i < len(samples); i++ {
		seg := dist2D(samples[i-1], samples[i])
		if acc+seg >= target && seg > 0 {
			f := (target - acc) / seg
			return Point{
				X: samples[i-1].X + (samples[i].X-samples[i-1].X)*f,
				Y: samples[i-1].Y + (samples[i].Y-samples[i-1].Y)*f,
			}
		}
		acc += seg
	}
	return samples[len(samples)-1]
}

// ArrowPath builds a crisp filled triangular arrowhead pointing along angle at tip.
func ArrowPath(tip Point, angle, size float32) *ui.Path {
	cosA := float32(math.Cos(float64(angle)))
	sinA := float32(math.Sin(float64(angle)))
	backX := tip.X - cosA*size
	backY := tip.Y - sinA*size
	wing := size * 0.52
	left := Point{X: backX - sinA*wing, Y: backY + cosA*wing}
	right := Point{X: backX + sinA*wing, Y: backY - cosA*wing}

	p := new(ui.Path)
	p.MoveTo(tip.X, tip.Y)
	p.LineTo(left.X, left.Y)
	p.LineTo(right.X, right.Y)
	p.Close()
	return p
}
