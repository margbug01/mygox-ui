package bencho

import (
	"math"
	"math/rand"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// CrumpleOptions configures Crumple.
type CrumpleOptions struct {
	Tilt    float32 // how far the receipt tips under the pointer, 0 to 15 degrees; 5 by default
	Power   float32 // how hard a flick throws the ball, 40 to 200 (%); 100 by default
	NoLabel bool    // no Move to Bin label while the ball is held over the bin
}

const (
	crumpleW       = 300
	crumpleH       = 380
	crumpleX       = 62 // the receipt's left
	crumpleY       = 28 // its top
	crumpleRW      = 176
	crumpleRH      = 324
	crumpleBall    = 34 // the ball's radius
	crumpleCols    = 6
	crumpleRows    = 10
	crumpleSqueeze = 550 * time.Millisecond
	crumplePrint   = 1100 * time.Millisecond
)

// crumplePhase is where the receipt is in its life.
type crumplePhase int

const (
	crumpleFlat     crumplePhase = iota // on the table, to be pressed
	crumpleBalling                      // scrunching up
	crumpleBallFree                     // a ball, to be picked up and thrown
	crumpleBinned                       // in the bin; a fresh one prints
)

type crumpleState struct {
	phase      crumplePhase
	since      time.Time
	last       time.Time
	x, y       float32 // the ball's centre
	vx, vy     float32
	held       bool
	grabX      float32 // the pointer off the ball's centre when picked up
	grabY      float32
	ball       [][2]float32 // each mesh vertex in the ball, off its centre
	shade      []float32    // each facet's shade
	seed       int64
	wasPressed bool
}

var crumpleBin = ui.Rect{X: 232, Y: 306, W: 52, H: 60}

// crumpleLine is a line of the receipt: text on the left, right or in
// the middle, at y (its middle) from the receipt's top.
type crumpleLine struct {
	y           float32
	left, right string
	mid         string
	size        float32
	weight      int
}

var crumpleLines = []crumpleLine{
	{y: 16.5, mid: "BENCHO", size: 15, weight: 700},
	{y: 31, mid: "BLOCK STORE", size: 6, weight: 700},
	{y: 39, mid: "12 COMPONENT LANE", size: 4.9},
	{y: 46, mid: "STORE 0042   TILL 03", size: 4.9},
	{y: 57, left: "17/10/2026", right: "21:40", size: 4.9},
	{y: 64, left: "ORDER #48213", right: "CASHIER: GOAT", size: 4.9},
	{y: 72, left: "__________________________________________", size: 4.9},
	{y: 82, left: "QTY ITEM", right: "EUR", size: 4.4, weight: 700},
	{y: 90.5, left: "1   COUNTRY WHEEL", right: "4.00", size: 5.5},
	{y: 99.5, left: "2   FOLDER", right: "6.50", size: 5.5},
	{y: 108.5, left: "1   COIN", right: "1.00", size: 5.5},
	{y: 117, left: "1   FOGGY WORD", right: "3.20", size: 5.5},
	{y: 126, left: "3   SLOSH", right: "2.40", size: 5.5},
	{y: 134.5, left: "1   CRUMPLE", right: "0.00", size: 5.5},
	{y: 144.5, left: "__________________________________________", size: 4.9},
	{y: 154, left: "SUBTOTAL", right: "17.10", size: 5.5},
	{y: 162, left: "VAT 21%", right: "3.59", size: 4.9},
	{y: 172.5, left: "TOTAL", right: "EUR 20.69", size: 9.6, weight: 700},
	{y: 185.5, left: "__________________________________________", size: 4.9},
	{y: 196, left: "VISA  **** **** **** 4417", right: "20.69", size: 4.9},
	{y: 203, left: "CONTACTLESS", right: "APPROVED", size: 4.9},
	{y: 210, left: "AUTH 0A73F2", right: "REF 99120", size: 4.9},
	{y: 219.5, left: "ITEMS 9", right: "CHANGE 0.00", size: 4.9},
	{y: 227.5, left: "__________________________________________", size: 4.9},
	{y: 267, mid: "4 006381 333931 48213", size: 4.9},
	{y: 278.5, mid: "THANK YOU FOR SHOPPING", size: 6.5, weight: 700},
	{y: 287, mid: "KEEP THIS RECEIPT", size: 4.9},
	{y: 293.5, mid: "RETURNS WITHIN 30 DAYS", size: 4.9},
	{y: 300, mid: "BENCHO.DEV", size: 4.9, weight: 700},
}

// Crumple is a long white till receipt. Pressing it scrunches it into a
// ball of paper; the ball can be picked up and thrown, and in the bin in
// the corner it goes, and a fresh receipt prints. The receipt tips a
// little under the pointer.
func Crumple(c *ui.Context, opts ...CrumpleOptions) ui.Element {
	var o CrumpleOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Tilt == 0 {
		o.Tilt = 5
	}
	o.Tilt = clamp32(o.Tilt, 0, 15)
	if o.Power == 0 {
		o.Power = 100
	}
	o.Power = clamp32(o.Power, 40, 200)
	t := ThemeFrom(c)
	box := ui.Box(c).Size(crumpleW, crumpleH).Cursor(ui.CursorPointer)
	st := ui.Local(box, "crumple", func() crumpleState { return crumpleState{seed: 1} })
	now := c.Now()
	dt := float32(0)
	if !st.last.IsZero() {
		dt = min(float32(now.Sub(st.last).Seconds()), 0.05)
	}
	st.last = now
	since := now.Sub(st.since)

	px, py, over := box.PointerPosition()
	_, _, pressed := box.Dragged()
	vx, vy := box.DragVelocity()
	down := pressed && !st.wasPressed
	st.wasPressed = pressed
	switch st.phase {
	case crumpleFlat:
		if down && px >= crumpleX && px <= crumpleX+crumpleRW && py >= crumpleY && py <= crumpleY+crumpleRH {
			crumpleNewBall(st)
			st.phase, st.since = crumpleBalling, now
			st.x, st.y = crumpleX+crumpleRW/2, crumpleY+crumpleRH/2
			st.held, st.grabX, st.grabY = true, px-st.x, py-st.y
		}
	case crumpleBalling:
		if since >= crumpleSqueeze {
			st.phase, st.since = crumpleBallFree, now
		}
	case crumpleBinned:
		if since >= 500*time.Millisecond+crumplePrint {
			st.phase, st.since = crumpleFlat, now
		}
	}
	if st.phase == crumpleBallFree || st.phase == crumpleBalling {
		if down && math.Hypot(float64(px-st.x), float64(py-st.y)) <= crumpleBall+8 {
			st.held, st.grabX, st.grabY = true, px-st.x, py-st.y
		}
		if st.held && !pressed {
			st.held = false
			k := o.Power / 100
			st.vx, st.vy = vx*k, vy*k
		}
		if st.held {
			// Pulled toward the pointer as it scrunches, then carried.
			k := clamp32(float32(since)/float32(crumpleSqueeze), 0, 1)
			if st.phase == crumpleBallFree {
				k = 1
			}
			st.x = lerp32(crumpleX+crumpleRW/2, px-st.grabX*(1-k), k)
			st.y = lerp32(crumpleY+crumpleRH/2, py-st.grabY*(1-k), k)
			st.vx, st.vy = 0, 0
		} else if st.phase == crumpleBallFree && dt > 0 {
			if crumpleFall(st, dt) {
				st.phase, st.since = crumpleBinned, now
			}
		}
	}
	overBin := st.held && st.phase == crumpleBallFree && crumpleNear(st.x, st.y, 18)
	c.AnimationFrame()

	// The receipt tips toward the pointer while it lies flat.
	tx, ty := float32(0), float32(0)
	if over && st.phase == crumpleFlat && !pressed {
		tx = -clamp32((py-crumpleH/2)/(crumpleRH/2), -1, 1) * o.Tilt
		ty = clamp32((px-crumpleW/2)/(crumpleRW/2), -1, 1) * o.Tilt
	}
	spr := ui.SpringOptions{Duration: 450 * time.Millisecond}
	tiltX, tiltY := box.Spring("tiltx", tx, spr), box.Spring("tilty", ty, spr)
	binShown := box.Spring("bin", b2f(st.phase != crumpleFlat), ui.SpringOptions{Duration: 350 * time.Millisecond})
	binGrow := box.Spring("bingrow", b2f(overBin), ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.4})
	label := box.Spring("label", b2f(overBin && !o.NoLabel), ui.SpringOptions{Duration: 200 * time.Millisecond})

	box.Draw(func(p *ui.Painter, r ui.Rect) {
		paper := ui.Rect{X: r.X + crumpleX, Y: r.Y + crumpleY, W: crumpleRW, H: crumpleRH}
		// The bin, behind the ball.
		if binShown > 0.01 {
			b := crumpleBin
			b.X, b.Y = r.X+b.X, r.Y+b.Y
			s := 1 + 0.08*binGrow
			p.Transform(ui.Scaling(s, s), b.X+b.W/2, b.Y+b.H, func() {
				p.Fill(ui.Rect{X: b.X + 4, Y: b.Y + 8, W: b.W - 8, H: b.H - 8}, t.Surface.Alpha(binShown), 8)
				p.Stroke(ui.Rect{X: b.X + 4, Y: b.Y + 8, W: b.W - 8, H: b.H - 8}, t.BorderStrong.Alpha(binShown), 8, 1.5)
				p.Fill(ui.Rect{X: b.X, Y: b.Y, W: b.W, H: 8}, t.BorderStrong.Alpha(binShown), 4)
				p.Icon(icons.Trash(), ui.Rect{X: b.X + b.W/2 - 10, Y: b.Y + 22, W: 20, H: 20}, t.TextSubtle.Alpha(binShown))
			})
			if label > 0.01 {
				span := ui.Span{Text: "Move to Bin", Size: 11, Weight: 500, Color: t.AccentText.Alpha(label)}
				w, h := p.MeasureText(0, span)
				lx, ly := min(b.X+b.W/2-w/2-8, r.X+r.W-w-20), b.Y-h-20
				p.Fill(ui.Rect{X: lx, Y: ly, W: w + 16, H: h + 8}, t.Accent.Alpha(label), (h+8)/2)
				p.RichText(lx+8, ly+4, w+2, span)
			}
		}
		switch st.phase {
		case crumpleFlat:
			p.Transform(ui.RotationX(tiltX).Then(ui.RotationY(tiltY)).Then(ui.Perspective(900)), paper.X+paper.W/2, paper.Y+paper.H/2, func() {
				crumpleReceipt(p, paper, 1)
			})
		case crumpleBinned:
			// The ball drops into the bin, then the next receipt rolls out.
			if since < 300*time.Millisecond {
				k := float32(since) / float32(300*time.Millisecond)
				b := crumpleBin
				crumpleDrawBall(p, st, r.X+lerp32(st.x, b.X+b.W/2, k), r.Y+lerp32(st.y, b.Y+b.H/2, k), 1-0.6*k, 1)
			}
			if since > 500*time.Millisecond {
				k := clamp32(float32(since-500*time.Millisecond)/float32(crumplePrint), 0, 1)
				k = 1 - (1-k)*(1-k)
				p.Clip(ui.Rect{X: r.X, Y: paper.Y - 2, W: r.W, H: r.H}, 0, func() {
					crumpleReceipt(p, ui.Rect{X: paper.X, Y: paper.Y - paper.H*(1-k), W: paper.W, H: paper.H}, 1)
				})
			}
		default:
			k := float32(1)
			if st.phase == crumpleBalling {
				k = clamp32(float32(since)/float32(crumpleSqueeze), 0, 1)
			}
			crumpleDrawSqueeze(p, st, paper, r.X+st.x, r.Y+st.y, k)
		}
	})
	return box
}

// crumpleNear tells whether (x, y) is within pad of the bin's mouth.
func crumpleNear(x, y, pad float32) bool {
	b := crumpleBin
	return x >= b.X-pad && x <= b.X+b.W+pad && y >= b.Y-pad && y <= b.Y+b.H+pad
}

// crumpleFall moves a thrown ball on by dt under gravity, off the walls
// and floor; it reports whether the ball went into the bin.
func crumpleFall(st *crumpleState, dt float32) bool {
	for s := float32(0); s < dt; s += 1.0 / 240 {
		h := min(float32(1.0/240), dt-s)
		st.vy += 2000 * h
		st.x += st.vx * h
		st.y += st.vy * h
		b := crumpleBin
		if st.x > b.X+6 && st.x < b.X+b.W-6 && st.y > b.Y && st.y < b.Y+b.H {
			return true
		}
		if st.x < crumpleBall {
			st.x, st.vx = crumpleBall, -st.vx*0.45
		}
		if st.x > crumpleW-crumpleBall {
			st.x, st.vx = crumpleW-crumpleBall, -st.vx*0.45
		}
		if st.y < crumpleBall {
			st.y, st.vy = crumpleBall, -st.vy*0.45
		}
		if st.y > crumpleH-crumpleBall {
			st.y, st.vy = crumpleH-crumpleBall, -st.vy*0.4
			st.vx *= float32(math.Exp(float64(-4 * h))) // rolls to a stop
		}
	}
	return false
}

// crumpleNewBall picks where each mesh vertex ends up in the ball, and how
// each facet is lit, a new way each time.
func crumpleNewBall(st *crumpleState) {
	rng := rand.New(rand.NewSource(st.seed))
	st.seed++
	st.ball, st.shade = st.ball[:0], st.shade[:0]
	for j := 0; j <= crumpleRows; j++ {
		for i := 0; i <= crumpleCols; i++ {
			u, v := float64(i)/crumpleCols*2-1, float64(j)/crumpleRows*2-1
			// The square onto the disc, then knocked about.
			x, y := u*math.Sqrt(1-v*v/2), v*math.Sqrt(1-u*u/2)
			x += (rng.Float64() - 0.5) * 0.35
			y += (rng.Float64() - 0.5) * 0.35
			st.ball = append(st.ball, [2]float32{float32(x) * crumpleBall, float32(y) * crumpleBall})
		}
	}
	for k := 0; k < crumpleCols*crumpleRows*2; k++ {
		st.shade = append(st.shade, rng.Float32())
	}
}

// crumpleDrawSqueeze draws the receipt k of the way to a ball centred at
// (cx, cy): its mesh drawn in, each facet shading as it creases, the
// print fading as it goes.
func crumpleDrawSqueeze(p *ui.Painter, st *crumpleState, paper ui.Rect, cx, cy, k float32) {
	e := 1 - (1-k)*(1-k)*(1-k)
	ox, oy := paper.X+paper.W/2, paper.Y+paper.H/2
	vert := func(i, j int) (float32, float32) {
		b := st.ball[j*(crumpleCols+1)+i]
		rx := paper.X + paper.W*float32(i)/crumpleCols
		ry := paper.Y + paper.H*float32(j)/crumpleRows
		return lerp32(rx, cx+b[0], e), lerp32(ry, cy+b[1], e)
	}
	crumpleMesh(p, st, vert, k)
	if k < 0.5 {
		s := 1 - 0.8*e
		p.Transform(ui.Scaling(s, s).Then(ui.Translation(lerp32(0, cx-ox, e), lerp32(0, cy-oy, e))), ox, oy, func() {
			crumpleText(p, paper, 1-2*k)
		})
	}
}

// crumpleDrawBall draws the finished ball at (cx, cy), scaled by s.
func crumpleDrawBall(p *ui.Painter, st *crumpleState, cx, cy, s, k float32) {
	crumpleMesh(p, st, func(i, j int) (float32, float32) {
		b := st.ball[j*(crumpleCols+1)+i]
		return cx + b[0]*s, cy + b[1]*s
	}, k)
}

// crumpleMesh fills the facets between the vertices vert gives, shaded k
// of the way to fully creased; a handful of paths, one per shade.
func crumpleMesh(p *ui.Painter, st *crumpleState, vert func(i, j int) (float32, float32), k float32) {
	const levels = 5
	var fills [levels]ui.Path
	var creases ui.Path
	f := 0
	tri := func(ax, ay, bx, by, cx, cy float32) {
		l := min(int(st.shade[f]*levels), levels-1)
		fills[l].MoveTo(ax, ay).LineTo(bx, by).LineTo(cx, cy).Close()
		creases.MoveTo(ax, ay).LineTo(bx, by).LineTo(cx, cy)
		f++
	}
	for j := 0; j < crumpleRows; j++ {
		for i := 0; i < crumpleCols; i++ {
			ax, ay := vert(i, j)
			bx, by := vert(i+1, j)
			cx, cy := vert(i, j+1)
			dx, dy := vert(i+1, j+1)
			if (i+j)%2 == 0 {
				tri(ax, ay, bx, by, dx, dy)
				tri(ax, ay, dx, dy, cx, cy)
			} else {
				tri(ax, ay, bx, by, cx, cy)
				tri(bx, by, dx, dy, cx, cy)
			}
		}
	}
	for l := range fills {
		shade := float32(l) / (levels - 1) * 0.3 * k
		p.FillPath(&fills[l], lerpColor(ui.Hex("#ffffff"), ui.Hex("#8a8a8a"), shade))
	}
	if k > 0.05 {
		p.StrokePath(&creases, 0.6, ui.Hex("#9a9a9a").Alpha(0.2*k))
	}
}

// crumpleReceipt draws the receipt in r, its print at alpha a.
func crumpleReceipt(p *ui.Painter, r ui.Rect, a float32) {
	// Torn top and bottom edges.
	var edge ui.Path
	const tooth, depth = 4, 3
	edge.MoveTo(r.X, r.Y+depth)
	for x := float32(0); x < r.W; x += tooth {
		edge.LineTo(r.X+x+tooth/2, r.Y).LineTo(r.X+min(x+tooth, r.W), r.Y+depth)
	}
	edge.LineTo(r.X+r.W, r.Y+r.H-depth)
	for x := r.W; x > 0; x -= tooth {
		edge.LineTo(r.X+x-tooth/2, r.Y+r.H).LineTo(r.X+max(x-tooth, 0), r.Y+r.H-depth)
	}
	edge.Close()
	p.FillPath(&edge, ui.Hex("#ffffff"))
	crumpleText(p, r, a)
}

// crumpleText prints the receipt's lines and barcode into r.
func crumpleText(p *ui.Painter, r ui.Rect, a float32) {
	if a <= 0 {
		return
	}
	ink, bold := ui.Hex("#4a4a4a").Alpha(a), ui.Hex("#242424").Alpha(a)
	const pad = 15
	for _, l := range crumpleLines {
		span := ui.Span{Font: "monospace", Size: l.size, Weight: max(l.weight, 400), Color: ink}
		if l.weight >= 700 {
			span.Color = bold
		}
		put := func(s string, x func(w float32) float32) {
			span.Text = s
			w, h := p.MeasureText(0, span)
			p.RichText(x(w), r.Y+l.y-h/2, w+2, span)
		}
		if l.mid != "" {
			put(l.mid, func(w float32) float32 { return r.X + r.W/2 - w/2 })
		}
		if l.left != "" {
			put(l.left, func(float32) float32 { return r.X + pad })
		}
		if l.right != "" {
			put(l.right, func(w float32) float32 { return r.X + r.W - pad - w })
		}
	}
	// The barcode: bars of one to three widths, the same every time.
	var bars ui.Path
	rng := rand.New(rand.NewSource(48213))
	for x := float32(22); x < r.W-22; {
		w := float32(1+rng.Intn(3)) * 0.45
		bars.MoveTo(r.X+x, r.Y+242).LineTo(r.X+x+w, r.Y+242).LineTo(r.X+x+w, r.Y+265).LineTo(r.X+x, r.Y+265).Close()
		x += w + float32(1+rng.Intn(4))*0.9
	}
	p.FillPath(&bars, bold)
}

func init() {
	registerDemo(
		"crumple",
		"A till receipt that scrunches into a paper ball when pressed; throw the ball in the bin and a fresh one prints.",
		func(c *ui.Context) { Crumple(c) },
	)
}
