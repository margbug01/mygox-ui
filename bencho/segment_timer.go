package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// SegmentTimerOptions configures SegmentTimer.
type SegmentTimerOptions struct {
	Start  int     // the minutes it is set to when it opens, 1 to 60; 25 by default
	Bounce float32 // how much give the ticks, the time and the buttons have, 0 to 100; 50 by default, negative for none
	Corner float32 // 0 to 40px; 20 by default, negative for square
	// OnDone is called when a countdown reaches 0:00.
	OnDone func()
}

const (
	segmentTimerW  = 260
	segmentTimerH  = 317
	segmentTimerCX = 130 // the dial's centre
	segmentTimerCY = 127.75
	segmentTimerR  = 97 // the ticks' outer end
)

type segmentTimerState struct {
	left    float32   // seconds left when last set, paused or started
	running time.Time // when the countdown started, or zero
	bumped  time.Time // when the time last jumped, for its pop
}

// SegmentTimer is a focus timer on a dial of sixty ticks, one a minute.
// Drag round the ring to set the time; play counts it down tick by tick,
// +5 adds five minutes and the reset button goes back to the start.
func SegmentTimer(c *ui.Context, opts ...SegmentTimerOptions) ui.Element {
	var o SegmentTimerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	startMin := 25
	if o.Start != 0 {
		startMin = int(clamp32(float32(o.Start), 1, 60))
	}
	prop := func(v, def, hi float32) float32 {
		if v == 0 {
			return def
		}
		return clamp32(v, 0, hi)
	}
	bounce := float64(prop(o.Bounce, 50, 100)) / 100
	corner := prop(o.Corner, 20, 40)
	t := ThemeFrom(c)
	now := c.Now()

	card := ui.Box(c).Size(segmentTimerW, segmentTimerH).Radius(corner).Background(t.Surface)
	st := ui.Local(card, "segment_timer", func() segmentTimerState { return segmentTimerState{left: float32(startMin) * 60} })

	left := st.left
	if !st.running.IsZero() {
		left = st.left - float32(now.Sub(st.running).Seconds())
		if left <= 0 {
			st.left, st.running, left = 0, time.Time{}, 0
			if o.OnDone != nil {
				o.OnDone()
			}
		} else {
			c.AnimationFrame()
		}
	}
	running := !st.running.IsZero()
	set := func(sec float32) {
		st.left = clamp32(sec, 0, 3600)
		if running {
			st.running = now
		}
		st.bumped = now
	}
	lit := int(math.Ceil(float64(left) / 60))

	card.Children(func() {
		dial := ui.Box(c).Key("dial").Absolute().Left(segmentTimerCX-112).Top(segmentTimerCY-112).Size(224, 224).Cursor(ui.CursorPointer)
		_, _, dragging := dial.Dragged()
		if x, y, ok := dial.PointerPosition(); ok && (dragging || dial.Pressed()) {
			// Clockwise from twelve, a minute a tick.
			a := math.Atan2(float64(x-112), float64(112-y)) * 180 / math.Pi
			if a < 0 {
				a += 360
			}
			if m := max(1, min(60, int(math.Ceil(a/6)))); float32(m*60) != st.left || running {
				set(float32(m * 60))
			}
		}
		var ticks [60]float32
		for i := range ticks {
			ticks[i] = dial.Spring(fmt.Sprint("tick", i), b2f(i < lit), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: bounce * 0.6})
		}
		dial.Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy := r.X+112, r.Y+112
			var off ui.Path
			off.Cap(ui.RoundCap)
			for i, k := range ticks {
				s, co := math.Sincos(float64(6*i+3) * math.Pi / 180)
				in := float32(segmentTimerR - 12.5 + 1.25)
				if i%5 == 4 {
					in -= 4
				}
				in -= 3 * max(k-1, 0) // overshoot reaches in a little
				out := float32(segmentTimerR - 1.25)
				x0, y0 := cx+in*float32(s), cy-in*float32(co)
				x1, y1 := cx+out*float32(s), cy-out*float32(co)
				if k <= 0.01 {
					off.MoveTo(x0, y0).LineTo(x1, y1)
					continue
				}
				var tick ui.Path
				tick.Cap(ui.RoundCap)
				tick.MoveTo(x0, y0).LineTo(x1, y1)
				p.StrokePath(&tick, 2.5, lerpColor(t.SurfaceHover, t.Text, clamp32(k, 0, 1)))
			}
			p.StrokePath(&off, 2.5, t.SurfaceHover)
		})

		pop := 1 + 0.06*float32(bounce)*float32(math.Max(0, 1-now.Sub(st.bumped).Seconds()/0.25))
		if now.Sub(st.bumped) < 250*time.Millisecond {
			c.AnimationFrame()
		}
		sec := int(math.Ceil(float64(left)))
		ui.Text(c, fmt.Sprintf("%d:%02d", sec/60, sec%60)).Key("time").Absolute().Left(segmentTimerCX-70).Top(segmentTimerCY-25.5).Size(140, 56).
			TextAlign(ui.Center).FontSize(41).FontWeight(500).TextColor(t.Text).Scale(pop, pop)

		button := func(key string, x, d float32, bg ui.Color, on bool, fn func(), paint func(b ui.Element)) {
			b := ui.Box(c).Key(key).Absolute().Left(x-d/2).Top(270-d/2).Size(d, d).Radius(d / 2).Background(bg).Center()
			if on {
				b.Cursor(ui.CursorPointer).OnClick(fn)
			}
			press := b.Spring("press", 1-0.1*b2f(on && b.Pressed()), ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: bounce * 0.7})
			b.Scale(press, press)
			paint(b)
		}
		button("add", 70.5, 35, t.SurfaceSubtle, true, func() { set(min(left+300, 3600)) }, func(b ui.Element) {
			b.Children(func() { ui.Text(c, "+5").FontSize(14.5).FontWeight(500).TextColor(t.Text) })
		})
		button("play", 130, 48, t.Text, left > 0, func() {
			if running {
				st.left, st.running = left, time.Time{}
			} else {
				st.running = now
			}
		}, func(b ui.Element) {
			b.Draw(func(p *ui.Painter, r ui.Rect) {
				cx, cy := r.X+r.W/2, r.Y+r.H/2
				if running {
					p.Fill(ui.Rect{X: cx - 6, Y: cy - 7.5, W: 4.5, H: 15}, t.Surface, 1.5)
					p.Fill(ui.Rect{X: cx + 1.5, Y: cy - 7.5, W: 4.5, H: 15}, t.Surface, 1.5)
					return
				}
				var tri ui.Path
				tri.Join(ui.RoundJoin)
				tri.MoveTo(cx-4.5, cy-7.5).LineTo(cx+7.5, cy).LineTo(cx-4.5, cy+7.5).Close()
				p.FillPath(&tri, t.Surface)
				p.StrokePath(&tri, 2.5, t.Surface)
			})
		})
		changed := running || left != float32(startMin)*60
		button("reset", 188, 36, t.SurfaceSubtle, changed, func() {
			st.left, st.running, st.bumped = float32(startMin)*60, time.Time{}, now
		}, func(b ui.Element) {
			col := t.TextSubtle.Alpha(0.6)
			if changed {
				col = t.Text
			}
			b.Children(func() { ui.Icon(c, icons.RotateCcw()).Size(17, 17).TextColor(col) })
		})
	})
	return card
}

func init() {
	registerDemo(
		"segment-timer",
		"A focus timer with a draggable dial of sixty ticks, one for each minute; drag around the ring to set the time, then play counts it down.",
		func(c *ui.Context) { SegmentTimer(c) },
	)
}
