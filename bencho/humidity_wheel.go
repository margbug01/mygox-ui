package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// HumidityWheelOptions configures HumidityWheel.
type HumidityWheelOptions struct {
	Wave  float32 // how much the lit band ripples, 0 (a still gauge) to 7; 2.6 by default, negative for still
	Ticks int     // density, 24 to 96; 60 by default
	Sweep float32 // degrees the ring spans, 240 to 340; 300 by default
	// Value is where it starts, 0 to 100; 15 by default.
	Value    int
	OnChange func(value int)
}

const (
	humidityWheelSize  = 250
	humidityWheelMid   = 124.75
	humidityWheelIn    = 64   // the ticks' inner end
	humidityWheelOut   = 92.5 // and outer, unlit
	humidityWheelReach = 12.5 // how much further the lit band's head reaches
)

var humidityWheelAmber = ui.Hex("#e5a224")

type humidityWheelState struct {
	value int
	born  time.Time
}

// HumidityWheel is a circular dial of ticks with a gap at the bottom. Drag
// round the ring to set the value; the lit band follows the finger,
// reaching further out toward its head, and ripples with a slow wave.
func HumidityWheel(c *ui.Context, opts ...HumidityWheelOptions) ui.Element {
	var o HumidityWheelOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	wave := float32(2.6)
	if o.Wave != 0 {
		wave = clamp32(o.Wave, 0, 7)
	}
	ticks := 60
	if o.Ticks != 0 {
		ticks = int(clamp32(float32(o.Ticks), 24, 96))
	}
	sweep := float32(300)
	if o.Sweep != 0 {
		sweep = clamp32(o.Sweep, 240, 340)
	}
	start := 15
	if o.Value != 0 {
		start = int(clamp32(float32(o.Value), 0, 100))
	}
	t := ThemeFrom(c)
	now := c.Now()

	dial := ui.Box(c).Size(humidityWheelSize, humidityWheelSize).Cursor(ui.CursorPointer)
	st := ui.Local(dial, "humidity_wheel", func() humidityWheelState { return humidityWheelState{value: start, born: now} })
	first := 180 + (360-sweep)/2 // the first tick's angle, clockwise from twelve: the gap is at six

	_, _, dragging := dial.Dragged()
	if x, y, ok := dial.PointerPosition(); ok && (dragging || dial.Pressed()) {
		a := float32(math.Atan2(float64(x-humidityWheelMid), float64(humidityWheelMid-y)) * 180 / math.Pi)
		a = float32(math.Mod(float64(a-first)+720, 360)) // from the first tick
		if a > sweep {                                   // in the gap: the nearer end
			a = sweep * b2f(a < sweep+(360-sweep)/2)
		}
		if v := int(math.Round(float64(a / sweep * 100))); v != st.value {
			st.value = v
			if o.OnChange != nil {
				o.OnChange(v)
			}
		}
	}
	// How many ticks are lit.
	head := dial.Spring("head", float32(st.value)/100*float32(ticks), ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.2})
	phase := float32(now.Sub(st.born).Seconds())
	if wave > 0 {
		c.AnimationFrame()
	}

	dial.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy := r.X+humidityWheelMid, r.Y+humidityWheelMid
		var off ui.Path
		off.Cap(ui.RoundCap)
		step := sweep / float32(ticks-1)
		for k := range ticks {
			fk := float32(k)
			s, co := math.Sincos(float64(first+fk*step) * math.Pi / 180)
			lit := clamp32(head-fk, 0, 1)
			out := float32(humidityWheelOut)
			if lit > 0 {
				// Longer toward the head, rippling as the wave runs along.
				ramp := 0.28 + 0.72*clamp32(fk/max(head, 1), 0, 1)
				ripple := wave / 2.6 * 2 * float32(math.Sin(float64(phase*3-fk*0.7)))
				out += lit * (humidityWheelReach*ramp + ripple)
			}
			in := float32(humidityWheelIn + 1.25)
			out -= 1.25
			x0, y0 := cx+in*float32(s), cy-in*float32(co)
			x1, y1 := cx+out*float32(s), cy-out*float32(co)
			if lit == 0 {
				off.MoveTo(x0, y0).LineTo(x1, y1)
				continue
			}
			var tick ui.Path
			tick.Cap(ui.RoundCap)
			tick.MoveTo(x0, y0).LineTo(x1, y1)
			p.StrokePath(&tick, 2.5, lerpColor(t.TextSubtle.Alpha(0.38), humidityWheelAmber, lit))
		}
		p.StrokePath(&off, 2.5, t.TextSubtle.Alpha(0.38))
	})
	dial.Children(func() {
		ui.Row(c).Absolute().Left(0).Top(humidityWheelMid-22).Size(humidityWheelSize, 44).AlignItems(ui.Center).Justify(ui.Center).
			Children(func() {
				ui.Text(c, fmt.Sprint(st.value)).FontSize(30).FontWeight(500).TextColor(t.Text)
				ui.Text(c, "%").FontSize(28).TextColor(t.Text).Margin(0, 0, 0, 3)
			})
	})
	return dial
}

func init() {
	registerDemo(
		"humidity",
		"A circular dial slider with a ring of ticks; drag around the ring to set a value and the lit band follows your finger, rippling with an optional wave.",
		func(c *ui.Context) { HumidityWheel(c) },
	)
}
