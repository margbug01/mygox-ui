package bencho

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"
)

// SleepDialSnap is the grid, in minutes, the handles land on.
type SleepDialSnap string

const (
	SleepDialSnap5  SleepDialSnap = "5"
	SleepDialSnap15 SleepDialSnap = "15"
	SleepDialSnap30 SleepDialSnap = "30"
)

// SleepDialOptions configures SleepDial.
type SleepDialOptions struct {
	Snap    SleepDialSnap // 15 by default
	Density int           // ticks round the day, 24 to 96; 48 by default
	Reach   float32       // how far the ticks in the range reach, 44 to 68px; 56 by default
	// OnChange is called with the range, in minutes from midnight.
	OnChange func(from, to int)
}

const (
	sleepDialSize   = 288
	sleepDialMid    = 144
	sleepDialIn     = 79.5 // where the ticks in the range start
	sleepDialRing   = 96   // the handles' path
	sleepDialHandle = 10   // the handles' radius
)

type sleepDialState struct {
	from, to int // minutes from midnight
	grab     int // the handle held: 0 none, 1 from, 2 to
}

func sleepDialAngle(min float32) float64 { return float64(min) / 1440 * 2 * math.Pi }

// SleepDial is a circular range slider over a day: two handles on a dial
// of ticks. Drag either handle round the ring to set the range between
// them; the ticks inside it stand dark and reach further, and the middle
// tells how long it is.
func SleepDial(c *ui.Context, opts ...SleepDialOptions) ui.Element {
	var o SleepDialOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	snap := 15
	switch o.Snap {
	case SleepDialSnap5:
		snap = 5
	case SleepDialSnap30:
		snap = 30
	}
	ticks := 48
	if o.Density != 0 {
		ticks = int(clamp32(float32(o.Density), 24, 96))
	}
	reach := float32(56)
	if o.Reach != 0 {
		reach = clamp32(o.Reach, 44, 68)
	}
	t := ThemeFrom(c)

	dial := ui.Box(c).Size(sleepDialSize, sleepDialSize).Cursor(ui.CursorPointer)
	st := ui.Local(dial, "sleep_dial", func() sleepDialState { return sleepDialState{from: 23 * 60, to: 6*60 + 30} })

	_, _, dragging := dial.Dragged()
	held := dragging || dial.Pressed()
	if x, y, ok := dial.PointerPosition(); ok && held {
		at := func(min int) (float32, float32) {
			s, co := math.Sincos(sleepDialAngle(float32(min)))
			return sleepDialMid + sleepDialRing*float32(s), sleepDialMid - sleepDialRing*float32(co)
		}
		if st.grab == 0 { // the nearer handle, if the press is on or by one
			fx, fy := at(st.from)
			tx, ty := at(st.to)
			df := math.Hypot(float64(x-fx), float64(y-fy))
			dt := math.Hypot(float64(x-tx), float64(y-ty))
			switch {
			case df <= dt && df < 28:
				st.grab = 1
			case dt < 28:
				st.grab = 2
			default:
				st.grab = -1 // not on a handle: nothing to drag
			}
		}
		a := math.Atan2(float64(x-sleepDialMid), float64(sleepDialMid-y))
		if a < 0 {
			a += 2 * math.Pi
		}
		min := int(math.Round(a/(2*math.Pi)*1440/float64(snap))) * snap % 1440
		from, to := st.from, st.to
		switch st.grab {
		case 1:
			from = min
		case 2:
			to = min
		}
		if from != to && (from != st.from || to != st.to) {
			st.from, st.to = from, to
			if o.OnChange != nil {
				o.OnChange(from, to)
			}
		}
	}
	if !held {
		st.grab = 0
	}

	span := (st.to - st.from + 1440) % 1440
	from, to := float32(st.from), float32(st.from+span) // to, counting on past midnight

	dial.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy := r.X+sleepDialMid, r.Y+sleepDialMid
		var off ui.Path
		off.Cap(ui.RoundCap)
		var on ui.Path
		on.Cap(ui.RoundCap)
		for k := range ticks {
			min := float32(k) * 1440 / float32(ticks)
			cardinal := k%(ticks/4) == 0
			s, co := math.Sincos(sleepDialAngle(min))
			path, in, out := &off, float32(84), float32(103)
			if cardinal {
				in, out = 80, 106
			}
			if inside := float32(math.Mod(float64(min-from)+2880, 1440)); inside > 0 && inside < to-from {
				path, in, out = &on, sleepDialIn, sleepDialIn+reach/2-4
				if cardinal {
					out += 4
				}
			}
			path.MoveTo(cx+in*float32(s), cy-in*float32(co)).LineTo(cx+out*float32(s), cy-out*float32(co))
		}
		p.StrokePath(&off, 2.5, t.TextSubtle.Alpha(0.42))
		p.StrokePath(&on, 3, t.Text)
		for i, m := range [2]float32{from, to} {
			s, co := math.Sincos(sleepDialAngle(m))
			hx, hy := cx+sleepDialRing*float32(s), cy-sleepDialRing*float32(co)
			// A soft shadow without a blur filter: rings fading out.
			for k := float32(6); k >= 1; k-- {
				rr := sleepDialHandle + k*1.6
				p.Fill(ui.Rect{X: hx - rr, Y: hy - rr + 2, W: 2 * rr, H: 2 * rr}, ui.RGBA(0, 0, 0, 0.022), rr)
			}
			p.Fill(ui.Rect{X: hx - sleepDialHandle, Y: hy - sleepDialHandle, W: 2 * sleepDialHandle, H: 2 * sleepDialHandle}, ui.Hex("#ffffff"), sleepDialHandle)
			if i == 0 {
				p.Fill(ui.Rect{X: hx - 2, Y: hy - 2, W: 4, H: 4}, ui.Hex("#808080"), 2)
			}
		}
	})
	dial.Children(func() {
		ui.Row(c).Absolute().Left(0).Top(sleepDialMid-18).Size(sleepDialSize, 44).AlignItems(ui.Baseline).Justify(ui.Center).Children(func() {
			unit := func(n int, u string, gap float32) {
				ui.Text(c, fmt.Sprint(n)).Key(u+"n").FontSize(30).FontWeight(500).TextColor(t.Text).Margin(0, 0, 0, gap)
				ui.Text(c, u).Key(u).FontSize(19).TextColor(t.TextSubtle.Alpha(0.8)).Margin(0, 0, 0, 4)
			}
			unit(span/60, "h", 0)
			unit(span%60, "m", 9)
		})
	})
	return dial
}

func init() {
	registerDemo(
		"sleep",
		"A circular range slider with two handles on a dial; drag either handle around the ring to set the range between them, the handles snapping to steps of five, fifteen or thirty minutes.",
		func(c *ui.Context) { SleepDial(c) },
	)
}
