package bencho

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// TapeOptions configures Tape.
type TapeOptions struct {
	NoSnap bool // lets go wherever it is, not on the nearest quarter step
	NoTape bool // hides the ticks that run at the chosen speed
	// Speed is where it starts, 0.75 to 2; 1 by default.
	Speed    float32
	OnChange func(speed float32)
}

const (
	tapeW      = 300
	tapeH      = 35
	tapeTrackW = 242
	tapeR      = 17   // the pill's end radius: the thumb sits that far in
	tapeMin    = 0.75 // speed at the left end
	tapeMax    = 2
	tapeTick   = 8  // from tick to tick
	tapeRun    = 30 // px a second the tape runs at 1x
)

func tapeX(speed float32) float32 {
	return tapeR + (speed-tapeMin)/(tapeMax-tapeMin)*(tapeTrackW-2*tapeR)
}

// tapeLabel writes 1 as 1.0×, 1.25 as 1.25×.
func tapeLabel(speed float32) string {
	s := strings.TrimRight(fmt.Sprintf("%.2f", speed), "0")
	if strings.HasSuffix(s, ".") {
		s += "0"
	}
	return s + "×"
}

type tapeState struct {
	speed float32
	run   float32 // how far the tape has run, px
	last  time.Time
}

// Tape is a playback speed slider whose filled track is a tape of ticks
// running past at the chosen speed. Drag along it to set the speed; let go
// and it settles on the nearest quarter step.
func Tape(c *ui.Context, opts ...TapeOptions) ui.Element {
	var o TapeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	start := float32(1)
	if o.Speed != 0 {
		start = clamp32(o.Speed, tapeMin, tapeMax)
	}
	t := ThemeFrom(c)
	now := c.Now()

	root := ui.Row(c).Size(tapeW, tapeH).AlignItems(ui.Center)
	st := ui.Local(root, "tape", func() tapeState { return tapeState{speed: start, last: now} })
	set := func(v float32) {
		if v != st.speed {
			st.speed = v
			if o.OnChange != nil {
				o.OnChange(v)
			}
		}
	}

	root.Children(func() {
		track := ui.Box(c).Size(tapeTrackW, tapeH).Cursor(ui.CursorResizeEW)
		_, _, dragging := track.Dragged()
		held := dragging || track.Pressed()
		if x, _, ok := track.PointerPosition(); ok && held {
			u := clamp32((x-tapeR)/(tapeTrackW-2*tapeR), 0, 1)
			set(float32(math.Round(float64(tapeMin+u*(tapeMax-tapeMin))*100)) / 100)
		} else if !o.NoSnap {
			set(float32(math.Round(float64(st.speed*4))) / 4)
		}
		thumb := track.Spring("thumb", tapeX(st.speed), ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: 0.2})

		if !o.NoTape {
			st.run += float32(now.Sub(st.last).Seconds()) * st.speed * tapeRun
			st.run = float32(math.Mod(float64(st.run), tapeTick))
			c.AnimationFrame()
		}
		st.last = now
		run := st.run

		track.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, t.SurfaceHover, tapeH/2)
			// The quarter steps past the fill.
			for v := float32(tapeMin); v < tapeMax; v += 0.25 {
				if x := tapeX(v); x > thumb+tapeR {
					p.Fill(ui.Rect{X: r.X + x - 1.5, Y: r.Y + tapeH/2 - 1.5, W: 3, H: 3}, t.TextSubtle.Alpha(0.55), 1.5)
				}
			}
			fill := ui.Rect{X: r.X, Y: r.Y, W: thumb + tapeR, H: tapeH}
			p.Fill(fill, t.Text, tapeH/2)
			if !o.NoTape {
				// Ticks run left through the fill, fading in at the thumb and out
				// at the left end.
				for x := thumb - run; x > 0; x -= tapeTick {
					a := 0.33 * clamp32(x/24, 0, 1) * clamp32((thumb-x)/10, 0, 1)
					p.Fill(ui.Rect{X: r.X + x - 1, Y: r.Y + 11, W: 2, H: 12}, t.Surface.Alpha(a), 1)
				}
			}
			p.Fill(ui.Rect{X: r.X + thumb - 1.5, Y: r.Y + 8.5, W: 3, H: 16.5}, t.Surface, 1.5)
		})
		ui.Text(c, tapeLabel(st.speed)).Width(tapeW - tapeTrackW).TextAlign(ui.End).FontSize(20).TextColor(t.TextSubtle)
	})
	return root
}

func init() {
	registerDemo(
		"tape",
		"A playback speed slider with a moving tape: the filled track's ticks run past at the speed you choose, and it settles on quarter steps when let go.",
		func(c *ui.Context) { Tape(c) },
	)
}
