package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// DateReelOptions configures DateReel.
type DateReelOptions struct {
	H24    bool    // a 24 hour clock, no AM/PM drum
	Curve  float32 // a flat list at 0, a tight drum at 100; 40 by default
	Glide  float32 // how far a flick carries, 0 to 100; 50 by default
	Corner float32 // 0 to 32; 20 by default
	Today  time.Time
	Start  time.Time // the picked moment; tomorrow 9:00 by default
}

type dateReelDrum struct {
	pos, vel float32 // in rows
	held     bool
}

type dateReelState struct {
	drums [4]dateReelDrum // day, hour, minute, half
	last  time.Time
}

const (
	dateReelW   = 300
	dateReelH   = 204
	dateReelRow = 34
)

// DateReel is a date and time picker of drums you flick: each glides with
// momentum and snaps into place on a row, under a band marking the choice.
func DateReel(c *ui.Context, opts ...DateReelOptions) ui.Element {
	var o DateReelOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Curve == 0 {
		o.Curve = 40
	}
	if o.Glide == 0 {
		o.Glide = 50
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	if o.Today.IsZero() {
		o.Today = time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC)
	}
	if o.Start.IsZero() {
		o.Start = o.Today.AddDate(0, 0, 1).Add(9 * time.Hour)
	}
	t := ThemeFrom(c)
	card := ui.Box(c).Size(dateReelW, dateReelH).Radius(clamp32(o.Corner, 0, 32)).Background(t.Surface).Clip()
	st := ui.Local(card, "date_reel", func() dateReelState {
		s := dateReelState{last: c.Now()}
		h := o.Start.Hour()
		s.drums[0].pos = float32(o.Start.Sub(o.Today).Hours() / 24)
		if o.H24 {
			s.drums[1].pos = float32(h)
		} else {
			s.drums[1].pos = float32((h + 11) % 12) // 1 to 12, from 1
			s.drums[3].pos = float32(h / 12)
		}
		s.drums[2].pos = float32(o.Start.Minute() / 5)
		return s
	})
	now := c.Now()
	dt := clamp32(float32(now.Sub(st.last).Seconds()), 0, 0.05)
	st.last = now
	curve := clamp32(o.Curve, 0, 100) / 100
	glide := 0.5 + clamp32(o.Glide, 0, 100)/100*2.5

	type column struct {
		x, w   float32
		end    bool // right-aligned text
		min    int  // rows run from min to max (or round and round)
		max    int
		cyclic bool
		label  func(i int) string
	}
	cols := []column{
		{6, 130, true, -30, 365, false, func(i int) string {
			if i == 0 {
				return "Today"
			}
			return o.Today.AddDate(0, 0, i).Format("Mon, Jan 2")
		}},
		{150, 34, true, 0, 11, true, func(i int) string { return fmt.Sprint(i + 1) }},
		{200, 36, false, 0, 11, true, func(i int) string { return fmt.Sprintf("%02d", i*5) }},
		{241, 40, false, 0, 1, false, func(i int) string { return [2]string{"AM", "PM"}[i] }},
	}
	if o.H24 {
		cols[1] = column{176, 36, true, 0, 23, true, func(i int) string { return fmt.Sprintf("%02d", i) }}
		cols[2].x = 222
		cols = cols[:3]
	}

	moving := false
	card.Children(func() {
		mid := float32(dateReelH) / 2
		ui.Box(c).Key("band").Absolute().Left(12).Top(mid-18).Size(dateReelW-24, 36).Radius(10).Background(t.SurfaceSubtle)
		for ci, col := range cols {
			d := &st.drums[ci]
			drum := ui.Box(c).Key(ci).Absolute().Left(col.x).Top(0).Size(col.w, dateReelH).Cursor(ui.CursorGrab)
			_, dy, dragging := drum.Dragged()
			if dragging {
				move := -dy / dateReelRow
				d.pos += move
				if dt > 0 {
					d.vel = d.vel*0.4 + move/dt*0.6
				}
				d.held = true
			} else {
				if d.held {
					d.held = false
					d.vel *= glide / 1.5
				}
				if math.Abs(float64(d.vel)) > 0.6 {
					d.pos += d.vel * dt
					d.vel *= float32(math.Exp(-float64(dt) * 3.2 / float64(glide)))
				} else {
					d.vel = 0
					target := float32(math.Round(float64(d.pos)))
					d.pos += (target - d.pos) * min(1, dt*14)
					if math.Abs(float64(target-d.pos)) < 0.002 {
						d.pos = target
					}
				}
			}
			if !col.cyclic {
				lo, hi := float32(col.min), float32(col.max)
				if d.pos < lo || d.pos > hi {
					d.vel = 0
					d.pos += (clamp32(d.pos, lo, hi) - d.pos) * min(1, dt*16)
				}
			}
			if dragging || d.vel != 0 || d.pos != float32(math.Round(float64(d.pos))) {
				moving = true
			}
			pos := d.pos
			drum.Draw(func(p *ui.Painter, r ui.Rect) {
				first := int(math.Floor(float64(pos))) - 4
				for i := first; i <= first+8; i++ {
					if !col.cyclic && (i < col.min || i > col.max) {
						continue
					}
					label := i
					if col.cyclic {
						n := col.max - col.min + 1
						label = ((i-col.min)%n+n)%n + col.min
					}
					off := float32(i) - pos
					// On a drum, rows bunch up and fade toward its edges.
					ang := float64(off) * (0.08 + 0.25*float64(curve))
					if math.Abs(ang) > math.Pi/2 {
						continue
					}
					y := mid + float32(math.Sin(ang)/(0.08+0.25*float64(curve)))*dateReelRow
					y = lerp32(mid+off*dateReelRow, y, clamp32(curve*2.5, 0, 1))
					near := clamp32(1-float32(math.Abs(float64(off))), 0, 1)
					fade := float32(math.Cos(ang)) * clamp32(1.5-float32(math.Abs(float64(off)))*0.55, 0, 1)
					if fade <= 0.02 {
						continue
					}
					ink := lerpColor(t.TextSubtle, t.Text, near).Alpha(clamp32(fade, 0, 1))
					span := ui.Span{Text: col.label(label), Size: 16.5, Color: ink}
					w, h := p.MeasureText(0, span)
					x := r.X
					if col.end {
						x = r.X + r.W - w
					}
					sy := float32(math.Cos(ang))*0.25 + 0.75
					p.Transform(ui.Scaling(1, sy), x, r.Y+y, func() {
						p.RichText(x, r.Y+y-h/2, w+2, span)
					})
				}
			})
		}
	})
	if moving {
		c.AnimationFrame()
	}
	return card
}

func init() {
	registerDemo(
		"date-reel",
		"A date and time picker of drums you flick: each glides and snaps into place under the band.",
		func(c *ui.Context) { DateReel(c) },
	)
}
