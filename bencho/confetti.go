package bencho

import (
	"math"
	"math/rand"
	"time"

	"github.com/egoist/mygo/ui"
)

// ConfettiBurstOptions configures ConfettiBurst.
type ConfettiBurstOptions struct {
	Amount int     // pieces in one burst, 10 to 240; 90 by default
	Spread float32 // how wide the burst fans out, in degrees, 10 to 180; 90 by default
	Mono   bool    // black and greys instead of colours
	Label  string  // "Celebrate" by default
}

type confettiPiece struct {
	vx, vy     float32 // launch velocity, px/s
	w, h       float32
	spin, flip float32 // turns per second
	phase      float32
	sway       float32
	col        ui.Color
}

type confettiBurst struct {
	at     time.Time
	pieces []confettiPiece
}

type confettiState struct {
	bursts []confettiBurst
	rng    *rand.Rand
}

const (
	confettiW    = 300
	confettiH    = 280
	confettiLife = 3.2 // seconds a burst lasts
	confettiDrag = 2.4
	confettiFall = 560
)

var confettiColours = []ui.Color{
	ui.Hex("#ff4d6d"), ui.Hex("#ffb703"), ui.Hex("#3a86ff"), ui.Hex("#06d6a0"),
	ui.Hex("#8338ec"), ui.Hex("#fb5607"), ui.Hex("#ff70a6"),
}

// Lucide's party-popper (ISC licence), which coss/icons doesn't carry.
var confettiPopper = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5.8 11.3 2 22l10.7-3.79"/><path d="M4 3h.01"/><path d="M22 8h.01"/><path d="M15 2h.01"/><path d="M22 20h.01"/><path d="m22 2-2.24.75a2.9 2.9 0 0 0-1.96 3.12c.1.86-.57 1.63-1.45 1.63h-.38c-.86 0-1.6.6-1.76 1.44L14 10"/><path d="m22 13-.82-.33c-.86-.34-1.82.2-1.98 1.11c-.11.7-.72 1.22-1.43 1.22H17"/><path d="m11 2 .33.82c.34.86-.2 1.82-1.11 1.98C9.52 4.9 9 5.52 9 6.23V7"/><path d="M11 13c1.93 1.93 2.83 4.17 2 5-.83.83-3.07-.07-5-2-1.93-1.93-2.83-4.17-2-5 .83-.83 3.07.07 5 2Z"/></svg>`))

// ConfettiBurst is a Celebrate button. Pressing it fires confetti that
// climbs, spreads and drifts down across the block, each piece flipping and
// swaying as it falls; presses in a row pile bursts on top of each other.
func ConfettiBurst(c *ui.Context, opts ...ConfettiBurstOptions) ui.Element {
	var o ConfettiBurstOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Amount == 0 {
		o.Amount = 90
	}
	o.Amount = min(max(o.Amount, 10), 240)
	if o.Spread == 0 {
		o.Spread = 90
	}
	o.Label = orDefault(o.Label, "Celebrate")
	t := ThemeFrom(c)
	stage := ui.Box(c).Size(confettiW, confettiH)
	st := ui.Local(stage, "confetti", func() confettiState { return confettiState{rng: rand.New(rand.NewSource(1))} })
	now := c.Now()
	live := st.bursts[:0]
	for _, b := range st.bursts {
		if now.Sub(b.at).Seconds() < confettiLife {
			live = append(live, b)
		}
	}
	st.bursts = live
	if len(st.bursts) > 0 {
		c.AnimationFrame()
	}
	fire := func() {
		b := confettiBurst{at: c.Now()}
		half := float64(clamp32(o.Spread, 10, 180)) / 2 * math.Pi / 180
		for i := 0; i < o.Amount; i++ {
			r := st.rng
			a := -math.Pi/2 + (r.Float64()*2-1)*half
			v := 300 + 260*r.Float64()
			col := confettiColours[r.Intn(len(confettiColours))]
			if o.Mono {
				col = lerpColor(ui.Hex("#111111"), ui.Hex("#b5b5b5"), r.Float32())
			}
			b.pieces = append(b.pieces, confettiPiece{
				vx: float32(math.Cos(a) * v), vy: float32(math.Sin(a) * v),
				w: 5 + 4*r.Float32(), h: 3 + 2*r.Float32(),
				spin: -1.5 + 3*r.Float32(), flip: 1 + 2.5*r.Float32(), phase: r.Float32() * 6.28,
				sway: 6 + 14*r.Float32(), col: col,
			})
		}
		st.bursts = append(st.bursts, b)
	}

	stage.Children(func() {
		bg, fg := ui.Hex("#17181a"), ui.Hex("#ffffff")
		if t.Dark {
			bg, fg = ui.Hex("#f4f4f5"), ui.Hex("#17181a")
		}
		btn := ui.Row(c).Key("button").Absolute().Left((confettiW-139)/2).Top((confettiH-44)/2).Size(139, 44).Radius(22).
			Background(bg).Center().Cursor(ui.CursorPointer).OnClick(fire)
		press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond, Bounce: 0.3})
		btn.Scale(1-0.06*press, 1-0.06*press).Children(func() {
			ui.Icon(c, confettiPopper).Size(17, 17).TextColor(fg).Margin(0, 10, 0, 0)
			ui.Text(c, o.Label).FontSize(16).FontWeight(500).TextColor(fg)
		})
	})
	// Pieces fly over everything, the button included.
	stage.DrawOver(func(p *ui.Painter, r ui.Rect) {
		ox, oy := r.X+r.W/2, r.Y+r.H/2-14
		for _, b := range st.bursts {
			tt := float32(now.Sub(b.at).Seconds())
			decay := float32(math.Exp(-confettiDrag * float64(tt)))
			fade := clamp32((confettiLife-tt)/0.6, 0, 1)
			for _, pc := range b.pieces {
				// Linear drag with gravity, solved in closed form.
				x := pc.vx * (1 - decay) / confettiDrag
				term := float32(confettiFall / confettiDrag)
				y := term*tt + (pc.vy-term)*(1-decay)/confettiDrag
				x += pc.sway * float32(math.Sin(float64(tt*3+pc.phase))) * (1 - decay)
				flip := float32(math.Cos(float64(tt*pc.flip*2*math.Pi + pc.phase)))
				cx, cy := ox+x, oy+y
				p.Transform(ui.Rotation(360*pc.spin*tt+pc.phase*57), cx, cy, func() {
					p.Transform(ui.Scaling(1, max(0.08, float32(math.Abs(float64(flip))))), cx, cy, func() {
						col := pc.col
						if flip < 0 {
							col = lerpColor(col, ui.Hex("#000000"), 0.18)
						}
						p.Fill(ui.Rect{X: cx - pc.w/2, Y: cy - pc.h/2, W: pc.w, H: pc.h}, col.Alpha(fade), 1)
					})
				})
			}
		}
	})
	return stage
}

func init() {
	registerDemo(
		"confetti",
		"A Celebrate button whose confetti climbs, spreads and drifts down, each piece flipping as it falls.",
		func(c *ui.Context) { ConfettiBurst(c) },
	)
}
