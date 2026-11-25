package bencho

import (
	"math"
	"math/rand"
	"time"

	"github.com/egoist/mygo/ui"
)

// TicketOptions configures Ticket.
type TicketOptions struct {
	Jelly  float32 // how soft the neck and the wobble are, 0 to 100; 50 by default
	Joints int     // how many tiny joints hold the stub on, 3 to 14; 7 by default
	Corner float32 // the ticket's corner radius, 0 to 28 px; 20 by default (use Square for none)
	Square bool    // square corners, as Corner 0 on bencho.dev
}

const (
	ticketW    = 300
	ticketH    = 186
	ticketX    = 10 // the ticket's left
	ticketY    = 15 // its top
	ticketTW   = 280
	ticketTH   = 156
	ticketTear = 210 // the perforation
)

type ticketPhase int

const (
	ticketOn     ticketPhase = iota // held on by at least one joint
	ticketFree                      // torn off, in the hand
	ticketFall                      // let go, dropping away
	ticketReturn                    // a fresh stub settling in
)

type ticketState struct {
	phase        ticketPhase
	broken       int // joints torn, from the top
	held         bool
	grabX, grabY float32
	pullX, pullY float32
	angle        float32 // the stub's turn about its bottom corner while held on
	since        time.Time
	fallX, fallY float32
	fallVX       float32
	fallVY       float32
	fallTurn     float32
	last         time.Time
}

// ticketBars are the stub's barcode: each bar's left and width, px.
var ticketBars = func() [][2]float32 {
	rng := rand.New(rand.NewSource(204))
	var bars [][2]float32
	for x := float32(0); x < 52; {
		w := float32(1+rng.Intn(3)) * 1.1
		bars = append(bars, [2]float32{x, min(w, 52-x)})
		x += w + float32(1+rng.Intn(2))*1.6
	}
	return bars
}()

// Ticket is a boarding pass, BCN to TYO, with a seat stub on the right
// held on by a perforation of tiny joints. Dragging the stub tears it
// away corner to corner, the joint at the tear stretching like jelly
// before it snaps, one at a time; torn off, the stub follows the hand and,
// let go, drops away while a fresh one settles in. Let go before it is
// torn through and the stub wobbles back into line.
func Ticket(c *ui.Context, opts ...TicketOptions) ui.Element {
	var o TicketOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Jelly == 0 {
		o.Jelly = 50
	}
	o.Jelly = clamp32(o.Jelly, 0, 100)
	if o.Joints == 0 {
		o.Joints = 7
	}
	o.Joints = min(max(o.Joints, 3), 14)
	if o.Corner == 0 {
		o.Corner = 20
	}
	if o.Square {
		o.Corner = 0
	}
	o.Corner = clamp32(o.Corner, 0, 28)
	t := ThemeFrom(c)
	box := ui.Box(c).Size(ticketW, ticketH)
	st := ui.Local(box, "ticket", func() ticketState { return ticketState{} })
	now := c.Now()
	dt := float32(0)
	if !st.last.IsZero() {
		dt = min(float32(now.Sub(st.last).Seconds()), 0.05)
	}
	st.last = now
	jointY := func(i int) float32 { // joint i's height, from the top
		lo, hi := float32(ticketY)+max(o.Corner, 10)*0.6, float32(ticketY+ticketTH)-max(o.Corner, 10)*0.6
		return lo + (hi-lo)*(float32(i)+0.5)/float32(o.Joints)
	}

	px, py, over := box.PointerPosition()
	_, _, pressed := box.Dragged()
	onStub := over && px > ticketTear && px < ticketX+ticketTW && py > ticketY && py < ticketY+ticketTH
	if onStub && (st.phase == ticketOn || st.phase == ticketReturn) {
		box.Cursor(ui.CursorGrab)
	}
	if pressed && !st.held && onStub && (st.phase == ticketOn || st.phase == ticketReturn) {
		st.held, st.grabX, st.grabY = true, px, py
		st.phase = ticketOn
	}
	soft := 0.5 + o.Jelly/100 // how far a joint stretches before it snaps
	pivotX, pivotY := float32(ticketTear), float32(ticketY+ticketTH)
	angleOf := func(x, y float32) float64 { return math.Atan2(float64(y-pivotY), float64(x-pivotX)) }
	if st.held && pressed {
		box.Cursor(ui.CursorGrabbing)
		st.pullX, st.pullY = px-st.grabX, py-st.grabY
		if st.phase == ticketOn {
			// The stub turns about its bottom corner after the hand; each
			// joint still holding stretches by how far it is carried, and
			// snaps past what the jelly allows, the top one first.
			turn := float32((angleOf(px, py) - angleOf(st.grabX, st.grabY)) * 180 / math.Pi)
			st.angle = clamp32(turn, -3, 75)
			for st.broken < o.Joints {
				h := pivotY - jointY(st.broken)
				if h*float32(math.Sin(float64(st.angle)*math.Pi/180)) < 13*soft {
					break
				}
				st.broken++
			}
			if st.broken == o.Joints {
				st.phase = ticketFree
				st.grabX, st.grabY = px, py
				st.pullX, st.pullY = 0, 0
			}
		}
	}
	if st.held && !pressed {
		st.held = false
		if st.phase == ticketFree {
			vx, vy := box.DragVelocity()
			st.phase, st.since = ticketFall, now
			st.fallX, st.fallY = st.pullX, st.pullY
			st.fallVX, st.fallVY = vx*0.4, min(vy*0.4, -120)
			st.fallTurn = 0
		}
		st.pullX, st.pullY = 0, 0
	}
	switch st.phase {
	case ticketFall:
		st.fallVY += 1800 * dt
		st.fallX += st.fallVX * dt
		st.fallY += st.fallVY * dt
		st.fallTurn += st.fallVX * 0.15 * dt
		if now.Sub(st.since) > 900*time.Millisecond {
			st.phase, st.since, st.broken = ticketReturn, now, 0
		}
	case ticketReturn:
		if now.Sub(st.since) > 500*time.Millisecond {
			st.phase = ticketOn
		}
	}
	// The stub turns about the joint at the tear, and its pull stretches
	// the jelly there; let go, it wobbles back.
	wob := ui.SpringOptions{Duration: 600 * time.Millisecond, Bounce: 0.25 + 0.5*float64(o.Jelly)/100}
	ta := float32(0)
	if st.held && st.phase == ticketOn {
		ta = st.angle
	}
	angle := box.Spring("angle", ta, wob)
	if st.held {
		angle = ta
	}
	pullX, pullY := st.pullX, st.pullY
	if st.phase != ticketOn || st.held {
		c.AnimationFrame()
	}

	muted, ink, dots := t.TextSubtle, t.Text, ui.Hex("#cdcdcd")
	if t.Dark {
		dots = t.BorderStrong
	}
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		ox, oy := r.X, r.Y
		// The pass.
		main := ui.Rect{X: ox + ticketX, Y: oy + ticketY, W: ticketTear - ticketX, H: ticketTH}
		p.FillRadii(main, [4]float32{o.Corner, 5, 5, o.Corner}, t.Surface)
		text := func(x, y float32, s string, size float32, w int, col ui.Color, spacing float32) {
			span := ui.Span{Text: s, Size: size, Weight: w, Color: col, LetterSpacing: spacing}
			_, h := p.MeasureText(0, span)
			p.RichText(ox+x, oy+y-h/2, 200, span)
		}
		text(30, 40, "BOARDING PASS", 10.5, 400, muted, 0.4)
		text(30, 76, "BCN", 32, 400, ink, 0)
		text(102, 75, "→", 17, 400, muted, 0)
		text(130, 76, "TYO", 32, 400, ink, 0)
		for i, col := range [][2]string{{"Flight", "BN 204"}, {"Gate", "12"}, {"Boards", "09:40"}} {
			x := []float32{30, 93, 133}[i]
			text(x, 123, col[0], 10.5, 400, muted, 0)
			text(x, 146, col[1], 12.5, 500, ink, 0)
		}
		// The joints still holding, as the perforation's dots.
		for y := float32(ticketY + 7); y < ticketY+ticketTH-4; y += 7 {
			p.Fill(ui.Rect{X: ox + ticketTear - 1.6, Y: oy + y - 1.6, W: 3.2, H: 3.2}, dots, 1.6)
		}

		// The stub, wherever it is.
		stub := func(alpha float32) {
			sr := ui.Rect{X: ox + ticketTear, Y: oy + ticketY, W: ticketX + ticketTW - ticketTear, H: ticketTH}
			p.FillRadii(sr, [4]float32{5, o.Corner, o.Corner, 5}, t.Surface.Alpha(alpha))
			cx := sr.X + sr.W/2
			for _, l := range []struct {
				s    string
				y    float32
				size float32
				col  ui.Color
			}{{"SEAT", 40, 10.5, muted}, {"3A", 74, 31, ink}} {
				span := ui.Span{Text: l.s, Size: l.size, Color: l.col.Alpha(alpha), LetterSpacing: 0.4 * b2f(l.size < 12)}
				w, h := p.MeasureText(0, span)
				p.RichText(cx-w/2, oy+l.y-h/2, w+2, span)
			}
			var bars ui.Path
			for _, b := range ticketBars {
				x := cx - 26 + b[0]
				bars.MoveTo(x, oy+121).LineTo(x+b[1], oy+121).LineTo(x+b[1], oy+150).LineTo(x, oy+150).Close()
			}
			p.FillPath(&bars, ink.Alpha(alpha))
		}
		switch st.phase {
		case ticketOn:
			// Turned about its bottom corner; the joints still holding
			// stretch from the pass to the stub like jelly.
			var necks ui.Path
			for i := st.broken; i < o.Joints; i++ {
				y := oy + jointY(i)
				ex, ey := ticketRotate(ox+ticketTear, y, ox+pivotX, oy+pivotY, angle)
				necks.MoveTo(ox+ticketTear-1, y).LineTo(ex+1, ey)
			}
			if angle > 0.5 {
				p.StrokePath(&necks, 3.2, t.Surface)
			}
			// Where it has torn, a hairline of the stage shows.
			if st.broken > 0 && angle < 0.5 {
				p.Fill(ui.Rect{X: ox + ticketTear - 0.5, Y: oy + ticketY, W: 1, H: jointY(st.broken-1) - ticketY + 4}, t.SurfaceSubtle, 0)
			}
			p.Transform(ui.Rotation(angle), ox+pivotX, oy+pivotY, func() { stub(1) })
		case ticketFree:
			p.Transform(ui.Rotation(clamp32(pullX*0.25, -30, 30)).Then(ui.Translation(pullX, pullY)), ox+ticketTear, oy+ticketY+ticketTH/2, func() { stub(1) })
		case ticketFall:
			a := 1 - clamp32(float32(now.Sub(st.since).Seconds())/0.9, 0, 1)
			p.Transform(ui.Rotation(st.fallTurn).Then(ui.Translation(st.fallX, st.fallY)), ox+ticketTear, oy+ticketY+ticketTH/2, func() { stub(a) })
		case ticketReturn:
			k := clamp32(float32(now.Sub(st.since).Seconds())/0.5, 0, 1)
			k = 1 - (1-k)*(1-k)
			p.Transform(ui.Translation(0, (1-k)*10), 0, 0, func() { stub(k) })
		}
	})
	return box
}

// ticketRotate turns (x, y) about (cx, cy) by degrees.
func ticketRotate(x, y, cx, cy, degrees float32) (float32, float32) {
	s, c := math.Sincos(float64(degrees) * math.Pi / 180)
	dx, dy := x-cx, y-cy
	return cx + dx*float32(c) - dy*float32(s), cy + dx*float32(s) + dy*float32(c)
}

func init() {
	registerDemo(
		"ticket",
		"A boarding pass whose seat stub tears off as you drag it, one jelly joint at a time.",
		func(c *ui.Context) { Ticket(c) },
	)
}
