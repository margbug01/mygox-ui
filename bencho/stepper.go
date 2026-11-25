package bencho

import (
	"math"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

const (
	dragStepperW, dragStepperH = 150, 46
	dragStepperMin             = 0
	dragStepperMax             = 999
	dragStepperPerStep         = 9                      // px of drag per step
	dragStepperRepeatAfter     = 400 * time.Millisecond // holding a button repeats after this
	dragStepperRoll            = 220 * time.Millisecond
)

type dragStepperState struct {
	from, dir int // the value before the last change and which way it went
	changedAt time.Time
	pressing  int // -1 or +1 while a button is held, 0 otherwise
	pressedAt time.Time
	repeats   int
	dragX     float32
	dragStart int
	dragging  bool
}

// DragStepper is a number stepper for small and large changes: tap − or +
// for a single step, hold one to repeat faster and faster, or drag across it
// to sweep through values.
func DragStepper(c *ui.Context, value *int) ui.Element {
	t := ThemeFrom(c)
	pill := ui.Box(c).Size(dragStepperW, dragStepperH).Radius(dragStepperH / 2).Background(t.Surface).Cursor(ui.CursorResizeEW)
	st := ui.Local(pill, "drag_stepper", func() dragStepperState { return dragStepperState{from: *value} })
	now := c.Now()
	set := func(v int) {
		v = min(max(v, dragStepperMin), dragStepperMax)
		if v == *value {
			return
		}
		st.from, st.dir, st.changedAt = *value, 1, now
		if v < *value {
			st.dir = -1
		}
		*value = v
	}

	// A drag across the pill sweeps; it starts once the pointer has moved a little.
	dx, _, _ := pill.Dragged()
	if pill.Pressed() {
		if !st.dragging && st.dragX == 0 {
			st.dragStart = *value
		}
		st.dragX += dx
		if math.Abs(float64(st.dragX)) > 4 {
			st.dragging = true
			st.pressing = 0
		}
		if st.dragging {
			set(st.dragStart + int(st.dragX/dragStepperPerStep))
		}
	} else {
		st.dragX, st.dragging = 0, false
	}

	// Holding a button repeats, faster the longer it is held.
	if st.pressing != 0 {
		held := now.Sub(st.pressedAt)
		if held >= dragStepperRepeatAfter {
			every := max(160*time.Millisecond-time.Duration(st.repeats)*12*time.Millisecond, 40*time.Millisecond)
			if held >= dragStepperRepeatAfter+time.Duration(st.repeats)*every {
				set(*value + st.pressing)
				st.repeats++
			}
		}
		c.AnimationFrame()
	}
	if now.Sub(st.changedAt) < dragStepperRoll {
		c.AnimationFrame()
	}
	lean := pill.Spring("ds_lean", clamp32(rubberBand(float32(math.Abs(float64(st.dragX))), 10), 0, 10)*float32(dragStepperSign(st.dragX)), ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.35})

	pill.Children(func() {
		for i, d := range []int{-1, 1} {
			btn := ui.Box(c).Key(i).Absolute().Left(float32(i)*(dragStepperW-46)).Top(0).Size(46, dragStepperH).Center().
				Cursor(ui.CursorPointer).Label([]string{"Decrease", "Increase"}[i])
			pressed := btn.Pressed() && !st.dragging
			switch {
			case pressed && st.pressing != d:
				st.pressing, st.pressedAt, st.repeats = d, now, 0
				set(*value + d) // the tap's single step lands on the press
			case !pressed && st.pressing == d:
				st.pressing = 0
			}
			hover := btn.Spring("hover", b2f(btn.Hovered()), ui.SpringOptions{Duration: 160 * time.Millisecond})
			push := btn.Spring("push", b2f(pressed), ui.SpringOptions{Duration: 140 * time.Millisecond, Bounce: 0.4})
			at := (*value <= dragStepperMin && d < 0) || (*value >= dragStepperMax && d > 0)
			icon := icons.Minus()
			if d > 0 {
				icon = icons.Plus()
			}
			s := 1 - 0.15*push
			btn.Draw(func(p *ui.Painter, r ui.Rect) {
				p.Fill(ui.Rect{X: r.X + 5, Y: r.Y + 5, W: 36, H: 36}, t.SurfaceSubtle.Alpha(hover*0.9), 18)
			})
			btn.Children(func() {
				col := t.Text
				if at {
					col = t.TextSubtle.Alpha(0.6)
				}
				ui.Icon(c, icon).Size(16, 16).TextColor(col).Scale(s, s).Translate(float32(d)*2*push, 0)
			})
		}
		// The number rolls in from the side it moved toward.
		k := float32(1)
		if !st.changedAt.IsZero() {
			k = clamp32(float32(now.Sub(st.changedAt))/float32(dragStepperRoll), 0, 1)
		}
		e := 1 - (1-k)*(1-k)
		num := ui.Box(c).Absolute().Left(46).Top(0).Size(dragStepperW-92, dragStepperH).Clip().Translate(lean*0.6, 0)
		num.Children(func() {
			if k < 1 {
				ui.Text(c, strconv.Itoa(st.from)).Key("old").FontSize(22).FontWeight(400).TextColor(t.Text).
					Absolute().Left(0).Top(dragStepperH/2-14).WidthPercent(100).TextAlign(ui.Center).
					Translate(0, -float32(st.dir)*18*e).Opacity(1 - e)
			}
			ui.Text(c, strconv.Itoa(*value)).Key("new").FontSize(22).FontWeight(400).TextColor(t.Text).
				Absolute().Left(0).Top(dragStepperH/2-14).WidthPercent(100).TextAlign(ui.Center).
				Translate(0, float32(st.dir)*18*(1-e)).Opacity(e)
		})
	})
	return pill
}

func dragStepperSign(v float32) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}

func init() {
	registerDemo(
		"stepper",
		"A number stepper for small and large changes: tap for a single step, or hold and drag to sweep through values quickly.",
		func(c *ui.Context) {
			holder := ui.Box(c)
			v := ui.Local(holder, "drag_stepper_value", func() int { return 24 })
			holder.Children(func() { DragStepper(c, v) })
		},
	)
}
