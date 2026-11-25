package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// EscapeButtonOptions configures EscapeButton.
type EscapeButtonOptions struct {
	Corner       float32 // 0 to 28; 28 by default
	Skittishness float32 // how hard it pushes away, 0 to 100; 55 by default
	Patience     int     // dodges before it gives in, 1 to 8; 4 by default
	Radius       float32 // how close you need to get, 60 to 200; 120 by default
	Label        string  // "Touch me" by default
}

type escapeButtonState struct {
	x, y   float32 // where the button's middle is heading
	near   bool    // the pointer is inside the button's reach
	dodges int
	got    time.Time // when it was finally pressed
}

const (
	escapeButtonW    = 300
	escapeButtonH    = 240
	escapeButtonBtnW = 122
	escapeButtonBtnH = 44
)

// EscapeButton is a button that slides away from the cursor when you try to
// reach it, and after a few dodges gives up and lets itself be pressed.
func EscapeButton(c *ui.Context, opts ...EscapeButtonOptions) ui.Element {
	var o EscapeButtonOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 28
	}
	if o.Skittishness == 0 {
		o.Skittishness = 55
	}
	if o.Patience == 0 {
		o.Patience = 4
	}
	if o.Radius == 0 {
		o.Radius = 120
	}
	o.Corner = clamp32(o.Corner, 0, 28)
	o.Skittishness = clamp32(o.Skittishness, 0, 100)
	o.Patience = min(max(o.Patience, 1), 8)
	o.Radius = clamp32(o.Radius, 60, 200)
	o.Label = orDefault(o.Label, "Touch me")
	t := ThemeFrom(c)

	stage := ui.Box(c).Size(escapeButtonW, escapeButtonH)
	st := ui.Local(stage, "escape_button", func() escapeButtonState {
		return escapeButtonState{x: escapeButtonW / 2, y: escapeButtonH / 2}
	})
	gotFor := c.Now().Sub(st.got)
	if !st.got.IsZero() && gotFor > 1800*time.Millisecond {
		st.got, st.dodges = time.Time{}, 0 // ready to run again
	}
	tired := st.dodges >= o.Patience

	// Reaching for it: once the pointer comes close it slides to the nearest
	// free spot at least Radius away, until its patience runs out.
	px, py, over := stage.PointerPosition()
	reach := 18 + o.Radius*0.12
	close := over && px > st.x-escapeButtonBtnW/2-reach && px < st.x+escapeButtonBtnW/2+reach &&
		py > st.y-escapeButtonBtnH/2-reach && py < st.y+escapeButtonBtnH/2+reach
	if close && !st.near && !tired {
		st.x, st.y = escapeButtonDodge(st.x, st.y, px, py, o.Radius, o.Skittishness)
		st.dodges++
	}
	st.near = close

	spring := ui.SpringOptions{Duration: time.Duration(460-2.6*o.Skittishness) * time.Millisecond, Bounce: 0.25}
	x := stage.Spring("x", st.x, spring)
	y := stage.Spring("y", st.y, spring)
	label := o.Label
	if !st.got.IsZero() {
		label = "You got me"
	}

	stage.Children(func() {
		btn := ui.Box(c).Absolute().Left(x-escapeButtonBtnW/2).Top(y-escapeButtonBtnH/2).Size(escapeButtonBtnW, escapeButtonBtnH).
			Radius(min(o.Corner, escapeButtonBtnH/2)).Background(t.Surface).Center().Cursor(ui.CursorPointer).
			OnClick(func() {
				if st.dodges >= o.Patience {
					st.got = c.Now()
				}
			})
		press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
		btn.Scale(1-0.04*press, 1-0.04*press)
		btn.Children(func() {
			ui.Text(c, label).FontSize(15.5).FontWeight(500).TextColor(t.TextMuted).SingleLine()
		})
	})
	if !st.got.IsZero() {
		c.AnimationFrame()
	}
	return stage
}

// escapeButtonDodge picks where the button flees to from (x, y) with the
// pointer at (px, py): the nearest spot inside the stage at least radius from
// the pointer, or the farthest one when none is that far. Skittishness pushes
// the pick further out.
func escapeButtonDodge(x, y, px, py, radius, skittish float32) (float32, float32) {
	const m = 6
	minX, maxX := float32(escapeButtonBtnW/2+m), float32(escapeButtonW-escapeButtonBtnW/2-m)
	minY, maxY := float32(escapeButtonBtnH/2+m), float32(escapeButtonH-escapeButtonBtnH/2-m)
	want := radius * (0.55 + skittish/100*0.6)
	bestX, bestY, best := x, y, float32(math.MaxFloat32)
	farX, farY, far := x, y, float32(-1)
	for j := 0; j <= 8; j++ {
		for i := 0; i <= 8; i++ {
			cx, cy := lerp32(minX, maxX, float32(i)/8), lerp32(minY, maxY, float32(j)/8)
			d := float32(math.Hypot(float64(cx-px), float64(cy-py)))
			move := float32(math.Hypot(float64(cx-x), float64(cy-y)))
			if d > far {
				farX, farY, far = cx, cy, d
			}
			if d >= want && move < best {
				bestX, bestY, best = cx, cy, move
			}
		}
	}
	if best == float32(math.MaxFloat32) {
		return farX, farY
	}
	return bestX, bestY
}

func init() {
	registerDemo(
		"escape",
		"A Touch me button that slides away from the cursor when you try to reach it, and after a few dodges gives up.",
		func(c *ui.Context) { EscapeButton(c) },
	)
}
