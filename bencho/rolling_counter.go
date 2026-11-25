package bencho

import (
	"math"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"
)

// RollingCounterOptions configures RollingCounter.
type RollingCounterOptions struct {
	Bounce   float32       // how lively a digit lands, 0 to 100; 45 by default
	Stagger  time.Duration // wait between one column and the next, 0 to 120ms; 45ms by default
	Blur     float32       // how much a spinning digit smears, 0 to 100; 50 by default
	Corner   float32       // 0 to 28; 16 by default
	Value    int           // 0 to 999; 95 by default
	OnChange func(v int)
}

// rollingCounterMark is the value from a moment on; the columns read it
// one stagger after another, so a carry ripples leftwards.
type rollingCounterMark struct {
	at    time.Time
	value int
}

type rollingCounterState struct {
	value   int
	marks   []rollingCounterMark // oldest first
	drag    float32              // drag px not yet turned into steps
	moved   float32              // drag px since time last moved on
	vel     float32              // steps per second, for the fling
	coast   float32              // fractional steps of the fling
	last    time.Time
	held    int // -1 or +1 while a button is held
	nextAt  time.Time
	holdAt  time.Time
	wasDrag bool
}

const (
	rollingCounterW     = 156
	rollingCounterH     = 56
	rollingCounterInset = 5
	rollingCounterBtnW  = 40
	rollingCounterMax   = 999
	rollingCounterCols  = 3
	rollingCounterPitch = 46 // px between one digit and the next on a reel
	rollingCounterStep  = 14 // px of drag per step
	rollingCounterFont  = 31
)

// RollingCounter is an odometer style number input between a minus and a
// plus. Tap a button to step, hold it to speed up, or drag the number up
// or right to count up and fling it; the digits roll and carry column by
// column.
func RollingCounter(c *ui.Context, opts ...RollingCounterOptions) ui.Element {
	var o RollingCounterOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Bounce == 0 {
		o.Bounce = 45
	}
	if o.Stagger == 0 {
		o.Stagger = 45 * time.Millisecond
	}
	if o.Blur == 0 {
		o.Blur = 50
	}
	if o.Corner == 0 {
		o.Corner = 16
	}
	if o.Value == 0 {
		o.Value = 95
	}
	stagger := min(max(o.Stagger, 0), 120*time.Millisecond)
	corner := clamp32(o.Corner, 0, 28)
	t := ThemeFrom(c)
	now := c.Now()

	card := ui.Box(c).Size(rollingCounterW, rollingCounterH).Background(t.Surface).Radius(corner)
	st := ui.Local(card, "rolling_counter", func() rollingCounterState {
		v := min(max(o.Value, 0), rollingCounterMax)
		return rollingCounterState{value: v, marks: []rollingCounterMark{{value: v}}, last: now}
	})
	dt := clamp32(float32(now.Sub(st.last).Seconds()), 0, 0.05)
	st.last = now
	set := func(v int) bool {
		v = min(max(v, 0), rollingCounterMax)
		if v == st.value {
			return false
		}
		st.value = v
		st.marks = append(st.marks, rollingCounterMark{now, v})
		if o.OnChange != nil {
			o.OnChange(v)
		}
		return true
	}
	// Forget marks every column has already read.
	for len(st.marks) > 1 && !st.marks[1].at.After(now.Add(-rollingCounterCols*stagger)) {
		st.marks = st.marks[1:]
	}
	lagged := func(k int) int {
		at := now.Add(-time.Duration(k) * stagger)
		for i := len(st.marks) - 1; i > 0; i-- {
			if !st.marks[i].at.After(at) {
				return st.marks[i].value
			}
		}
		return st.marks[0].value
	}

	button := func(dir int, left float32) {
		b := ui.Box(c).Key(dir).Absolute().Left(left).Top(rollingCounterInset).
			Size(rollingCounterBtnW, rollingCounterH-2*rollingCounterInset).Cursor(ui.CursorPointer)
		if b.Pressed() {
			if st.held != dir {
				// Down: one step now, then repeats that come faster the longer it is held.
				st.held, st.holdAt, st.nextAt = dir, now, now.Add(400*time.Millisecond)
				set(st.value + dir)
			} else if !now.Before(st.nextAt) {
				set(st.value + dir)
				gap := 140 * math.Exp(-now.Sub(st.holdAt).Seconds()*1.1)
				st.nextAt = now.Add(time.Duration(max(gap, 28) * float64(time.Millisecond)))
			}
			c.AnimationFrame()
		} else if st.held == dir {
			st.held = 0
		}
		hover := b.Spring("hover", b2f(b.Hovered()), ui.SpringOptions{Duration: 160 * time.Millisecond})
		press := b.Spring("press", b2f(b.Pressed()), ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.3})
		b.Background(lerpColor(t.SurfaceSubtle, t.SurfaceHover, hover)).Radius(max(corner-rollingCounterInset, 0)).
			Scale(1-0.08*press, 1-0.08*press)
		b.Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			p.Fill(ui.Rect{X: cx - 5.5, Y: cy - 0.6, W: 11, H: 1.2}, t.Text, 0)
			if dir > 0 {
				p.Fill(ui.Rect{X: cx - 0.6, Y: cy - 5.5, W: 1.2, H: 11}, t.Text, 0)
			}
		})
	}

	card.Children(func() {
		button(-1, rollingCounterInset)
		button(1, rollingCounterW-rollingCounterInset-rollingCounterBtnW)

		areaW := float32(rollingCounterW - 2*rollingCounterInset - 2*rollingCounterBtnW)
		area := ui.Box(c).Absolute().Left(rollingCounterInset+rollingCounterBtnW).Top(rollingCounterInset).
			Size(areaW, rollingCounterH-2*rollingCounterInset).Cursor(ui.CursorGrab)
		dx, dy, dragging := area.Dragged()
		if dragging {
			// Up or right counts up.
			move := dx - dy
			st.drag += move
			n := int(st.drag / rollingCounterStep)
			st.drag -= float32(n) * rollingCounterStep
			set(st.value + n)
			st.moved += move
			if dt > 0 {
				st.vel = st.vel*0.5 + st.moved/rollingCounterStep/dt*0.5
				st.moved = 0
			}
		} else {
			if st.wasDrag {
				st.drag, st.coast, st.moved = 0, 0, 0
			}
			if st.vel != 0 {
				st.coast += st.vel * dt
				n := int(st.coast)
				st.coast -= float32(n)
				if n != 0 && !set(st.value+n) {
					st.vel = 0
				}
				st.vel *= float32(math.Exp(-float64(dt) * 3.2))
				if math.Abs(float64(st.vel)) < 2 {
					st.vel = 0
				}
				c.AnimationFrame()
			}
		}
		st.wasDrag = dragging
		if now.Sub(st.marks[len(st.marks)-1].at) < rollingCounterCols*stagger {
			c.AnimationFrame()
		}

		roll := ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: float64(clamp32(o.Bounce, 0, 100)) / 100 * 0.55}
		colW, _ := c.MeasureText(0, ui.Span{Text: "0", Size: rollingCounterFont})
		var pos, width [rollingCounterCols]float32
		total := float32(0)
		for k := range rollingCounterCols {
			v := lagged(k)
			p10 := int(math.Pow10(k))
			// Each reel counts in whole turns of the one before, so a carry rolls forward.
			pos[k] = area.Spring(k, float32(v/p10), roll)
			show := k == 0 || v >= p10
			width[k] = area.Spring(10+k, b2f(show)*colW, ui.SpringOptions{Duration: 320 * time.Millisecond})
			total += width[k]
		}
		smearK := clamp32(o.Blur, 0, 100) / 100
		area.Draw(func(p *ui.Painter, r ui.Rect) {
			x := r.X + (r.W-total)/2
			cy := r.Y + r.H/2
			for k := rollingCounterCols - 1; k >= 0; k-- {
				w := width[k]
				if w < 0.5 {
					continue
				}
				col := ui.Rect{X: x, Y: r.Y, W: w, H: r.H}
				x += w
				target := float32(lagged(k) / int(math.Pow10(k)))
				smear := clamp32(float32(math.Abs(float64(target-pos[k])))/2, 0, 1) * smearK
				fade := clamp32(w/colW, 0, 1)
				base := float32(math.Floor(float64(pos[k])))
				p.Clip(col, 0, func() {
					for i := base - 1; i <= base+2; i++ {
						d := int(mod1(float64(i)/10)*10+0.5) % 10
						y := cy + (i-pos[k])*rollingCounterPitch
						draw := func(oy, a float32) {
							span := ui.Span{Text: strconv.Itoa(d), Size: rollingCounterFont, Color: t.Text.Alpha(a * fade)}
							tw, th := p.MeasureText(0, span)
							p.RichText(col.X+(colW-tw)/2-(colW-w)/2, y-th/2+oy, colW, span)
						}
						draw(0, 1-0.45*smear)
						if smear > 0.02 {
							draw(-6*smear, 0.28*smear)
							draw(6*smear, 0.28*smear)
						}
					}
				})
			}
			// The reel fades out above and below the window.
			fadeH := float32(12)
			p.FillGradient(ui.Rect{X: r.X, Y: r.Y, W: r.W, H: fadeH}, ui.LinearGradient{From: t.Surface, To: t.Surface.Alpha(0), Angle: 180}, 0)
			p.FillGradient(ui.Rect{X: r.X, Y: r.Y + r.H - fadeH, W: r.W, H: fadeH}, ui.LinearGradient{From: t.Surface.Alpha(0), To: t.Surface, Angle: 180}, 0)
		})
	})
	return card
}

func init() {
	registerDemo(
		"rolling-counter",
		"An odometer style number input between a minus and a plus: tap, hold to speed up, or drag and fling the number, and the digits roll and carry column by column.",
		func(c *ui.Context) { RollingCounter(c) },
	)
}
