package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// OptionWheelOptions configures OptionWheel.
type OptionWheelOptions struct {
	Tilt float32 // how hard the words lean back on the drum, 0 to 100; 60 by default, negative for none
	Blur float32 // how far out of focus the far words go, 0 to 10px; 4 by default, negative for none
	// Items are the words on the wheel; ten cities by default, New York chosen.
	Items    []string
	OnChange func(index int)
}

const (
	optionWheelW      = 300
	optionWheelH      = 320
	optionWheelLeft   = 43  // the chosen word's left edge
	optionWheelRadius = 121 // the drum's radius
	optionWheelStep   = 25.8 * math.Pi / 180
	optionWheelFont   = 38
	optionWheelLine   = 46 // the height of a line of it
)

var optionWheelCities = []string{"London", "Rome", "Tokyo", "Paris", "Berlin", "New York", "Lisbon", "Seoul", "Sydney", "Cairo"}

type optionWheelState struct {
	target int     // the chosen word, counting on past either end
	pos    float32 // where the wheel is, in words
	scroll float32 // scroll not yet turned into a word
	last   time.Time
}

// OptionWheel is a curved drum of words. The chosen one sits in the middle
// while the rest lean away on the drum, softer and out of focus. Drag,
// scroll or use the arrow keys to spin it.
func OptionWheel(c *ui.Context, opts ...OptionWheelOptions) ui.Element {
	var o OptionWheelOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	// Zero means the default; pass a negative value for none.
	prop := func(v, def, hi float32) float32 {
		if v == 0 {
			return def
		}
		return clamp32(v, 0, hi)
	}
	tilt, blur := prop(o.Tilt, 60, 100), prop(o.Blur, 4, 10)
	items := o.Items
	start := 0
	if len(items) == 0 {
		items, start = optionWheelCities, 5
	}
	n := len(items)
	t := ThemeFrom(c)
	now := c.Now()

	wheel := ui.Box(c).Size(optionWheelW, optionWheelH).Cursor(ui.CursorResizeNS).Focusable().FocusRing(false)
	st := ui.Local(wheel, "option_wheel", func() optionWheelState {
		return optionWheelState{target: start, pos: float32(start), last: now}
	})
	choose := func(i int) {
		if i == st.target {
			return
		}
		st.target = i
		if o.OnChange != nil {
			o.OnChange(((i % n) + n) % n)
		}
	}

	wheel.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputScroll {
			return false
		}
		st.scroll += ev.DY
		for st.scroll >= 36 {
			st.scroll -= 36
			choose(st.target + 1)
		}
		for st.scroll <= -36 {
			st.scroll += 36
			choose(st.target - 1)
		}
		return true
	})
	if wheel.Shortcut(0, ui.KeyDown) {
		choose(st.target + 1)
	}
	if wheel.Shortcut(0, ui.KeyUp) {
		choose(st.target - 1)
	}

	dt := float32(now.Sub(st.last).Seconds())
	st.last = now
	_, dy, dragging := wheel.Dragged()
	switch {
	case dragging:
		// The words follow the pointer, a word per word's pitch at the middle.
		st.pos -= dy / 52
		choose(int(math.Round(float64(st.pos))))
	case wheel.Pressed():
	default:
		st.pos += (float32(st.target) - st.pos) * (1 - float32(math.Exp(float64(-dt*14))))
		if math.Abs(float64(float32(st.target)-st.pos)) > 0.002 {
			c.AnimationFrame()
		} else {
			st.pos = float32(st.target)
		}
	}

	pos := st.pos
	wheel.Children(func() {
		mid := int(math.Round(float64(pos)))
		for i := mid - 4; i <= mid+4; i++ {
			th := (float32(i) - pos) * optionWheelStep
			if math.Abs(float64(th)) >= math.Pi/2*0.98 {
				continue
			}
			s, co := math.Sincos(float64(th))
			sin, cos := float32(s), float32(co)
			far := float32(math.Abs(float64(th))) / (math.Pi / 2) // 0 in the middle, 1 at the drum's edge
			lean := tilt / 60
			x := optionWheelLeft - 38*(1-cos)*lean
			cy := optionWheelH/2 + optionWheelRadius*sin
			soft := blur * far * 2
			word := items[((i%n)+n)%n]
			box := ui.Box(c).Key(i).Absolute().Left(x).Top(cy-optionWheelLine/2).Size(optionWheelW, optionWheelLine).
				TransformOrigin(0.2, 0.5).Scale(0.6+0.4*cos, max(cos, 0.35)).
				Rotate(th * 180 / math.Pi * 0.3 * lean).Opacity(clamp32(1-far*0.93, 0, 1))
			box.Children(func() {
				text := func(k int, dx, dy, a float32) {
					ui.Text(c, word).Key(k).Absolute().Left(dx).Top(dy).Height(optionWheelLine).FontSize(optionWheelFont).FontWeight(500).
						TextColor(t.Text.Alpha(a))
				}
				if soft < 0.4 {
					text(0, 0, 0, 1)
					return
				}
				// A soft look without a blur filter: a ring of faint copies.
				text(0, 0, 0, 0.3)
				for k := 1; k <= 8; k++ {
					a := float64(k) * math.Pi / 4
					r := soft * (0.5 + 0.5*float32(k%2))
					text(k, r*float32(math.Cos(a)), r*float32(math.Sin(a)), 0.13)
				}
			})
		}
	})
	return wheel
}

func init() {
	registerDemo(
		"option-wheel",
		"A curved drum picker for words: the chosen one sits in the middle while the rest tilt away, softer and out of focus; spin it by dragging, scrolling or the arrow keys.",
		func(c *ui.Context) { OptionWheel(c) },
	)
}
