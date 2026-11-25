package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
)

// ProgressTrackerOptions configures ProgressTracker; the zero value is
// bencho.dev's default.
type ProgressTrackerOptions struct {
	Title    string  // "Your first week" by default
	Bounce   float32 // how far the bar runs past its mark before it settles, 0 to 100; 40 by default
	NoSparks bool    // no sparks thrown off the bar's head as it moves forward
	OnChange func(done int)
}

var progressTrackerTasks = []string{"Pick a profile photo", "Say hi to the team", "Ship your first project", "Break something (safely)", "Celebrate with a snack"}

const (
	progressTrackerW, progressTrackerH = 320, 278
	progressTrackerBarY                = 58
	progressTrackerRow0                = 94.5 // the first row's middle
	progressTrackerPitch               = 38.4
	progressTrackerBurstLife           = 750 * time.Millisecond
)

var progressTrackerFrom, progressTrackerTo = ui.Hex("#756af4"), ui.Hex("#a452d8")

type progressTrackerState struct {
	done      []bool
	allDoneAt time.Time
}

// ProgressTracker is a getting started checklist under a progress bar.
// Ticking a task crosses it out, the bar springs forward throwing sparks and
// the count rolls up; finishing them all sets off a small burst on every
// row. Tick a task again to bring it back.
func ProgressTracker(c *ui.Context, tasks []string, opts ...ProgressTrackerOptions) ui.Element {
	var o ProgressTrackerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Title == "" {
		o.Title = "Your first week"
	}
	if o.Bounce == 0 {
		o.Bounce = 40
	}
	t := ThemeFrom(c)
	now := c.Now()

	card := ui.Box(c).Size(progressTrackerW, progressTrackerH).Radius(24).Background(t.Surface)
	st := ui.Local(card, "progress_tracker", func() progressTrackerState {
		d := make([]bool, len(tasks))
		if len(d) > 0 {
			d[0] = true
		}
		return progressTrackerState{done: d}
	})
	count := 0
	for _, d := range st.done {
		if d {
			count++
		}
	}
	toggle := func(i int) {
		st.done[i] = !st.done[i]
		n := 0
		for _, d := range st.done {
			if d {
				n++
			}
		}
		if st.done[i] && n == len(tasks) {
			st.allDoneAt = now
		}
		if o.OnChange != nil {
			o.OnChange(n)
		}
	}

	roll := card.Spring("count", float32(count), ui.SpringOptions{Duration: 420 * time.Millisecond})
	burstAge := now.Sub(st.allDoneAt)
	if count == len(tasks) && burstAge < progressTrackerBurstLife {
		c.AnimationFrame()
	}

	card.Draw(func(p *ui.Painter, r ui.Rect) {
		// The count, its done number rolling up or down.
		_, th := p.MeasureText(0, ui.Span{Text: "0", Size: 15})
		tw, _ := p.MeasureText(0, ui.Span{Text: fmt.Sprintf("/%d", len(tasks)), Size: 15})
		dw, _ := p.MeasureText(0, ui.Span{Text: "0", Size: 15})
		right := r.X + r.W - 20
		p.RichText(right-tw, r.Y+21, 0, ui.Span{Text: fmt.Sprintf("/%d", len(tasks)), Size: 15, Color: t.TextSubtle})
		digits := ui.Rect{X: right - tw - dw, Y: r.Y + 21, W: dw, H: th}
		lo := float32(math.Floor(float64(roll)))
		f := roll - lo
		digit := func(n int, y float32) { // right-aligned against the slash
			sp := ui.Span{Text: fmt.Sprint(n), Size: 15, Color: t.TextSubtle}
			w, _ := p.MeasureText(0, sp)
			p.RichText(digits.X+digits.W-w, y, 0, sp)
		}
		p.Clip(ui.Rect{X: digits.X - 4, Y: digits.Y, W: digits.W + 4, H: digits.H}, 0, func() {
			digit(int(lo), digits.Y-f*th)
			if f > 0.01 {
				digit(int(lo)+1, digits.Y+(1-f)*th)
			}
		})
	})

	card.Children(func() {
		// The bar, springing forward and throwing sparks.
		coss.Progress(c, float32(count*100)/float32(max(1, len(tasks))), coss.ProgressOptions{
			Bounce: o.Bounce, Sparks: !o.NoSparks,
			Style: &coss.ProgressStyle{Height: 8, Track: t.SurfaceHover, From: progressTrackerFrom, To: progressTrackerTo},
		}).Absolute().Left(20).Top(progressTrackerBarY).Width(progressTrackerW - 40)
		ui.Text(c, o.Title).Absolute().Left(20).Top(19).FontSize(18).FontWeight(500).TextColor(t.Text)
		ui.Text(c, fmt.Sprintf("%d/%d", count, len(tasks))).Absolute().Left(0).Top(0).Opacity(0) // the count as text, for screen readers
		for i, task := range tasks {
			progressTrackerRow(c, t, i, task, st.done[i], func() { toggle(i) }, count == len(tasks), burstAge)
		}
	})
	return card
}

func progressTrackerRow(c *ui.Context, t Theme, i int, task string, done bool, toggle func(), all bool, burstAge time.Duration) {
	cy := progressTrackerRow0 + float32(i)*progressTrackerPitch
	row := ui.Row(c).Key(task).Absolute().Left(10).Top(cy-17).Size(progressTrackerW-20, 34).Padding(0, 10).Gap(12).AlignItems(ui.Center).Cursor(ui.CursorPointer).OnClick(toggle)
	k := row.Spring("done", b2f(done), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.35})
	strike := row.Spring("strike", b2f(done), ui.SpringOptions{Duration: 420 * time.Millisecond})
	press := row.Spring("press", b2f(row.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
	row.Scale(1-0.015*press, 1-0.015*press)
	row.Children(func() {
		ui.Box(c).Size(20, 20).Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy := r.X+10, r.Y+10
			p.Stroke(r, t.TextSubtle.Alpha(0.55*(1-clamp32(k, 0, 1))), 10, 1.8)
			if s := 20 * clamp32(k, 0, 1.15); s > 0.5 {
				p.Fill(ui.Rect{X: cx - s/2, Y: cy - s/2, W: s, H: s}, t.Text, s/2)
				var check ui.Path
				check.Cap(ui.RoundCap).Join(ui.RoundJoin).MoveTo(cx-4, cy+0.2).LineTo(cx-1.2, cy+3).LineTo(cx+4.2, cy-2.6)
				check.Trim(0, clamp32(k, 0, 1))
				p.StrokePath(&check, 2, t.Surface)
			}
			// Finishing them all bursts every row.
			if all && burstAge < progressTrackerBurstLife {
				b := float32(burstAge.Seconds()/progressTrackerBurstLife.Seconds()) - float32(i)*0.06
				if b > 0 && b < 1 {
					for j := range 10 {
						a := float64(j)/10*2*math.Pi + float64(i)
						d := 12 + 14*b
						x, y := cx+float32(math.Cos(a))*d, cy+float32(math.Sin(a))*d
						s := 3 * (1 - b)
						p.Fill(ui.Rect{X: x - s/2, Y: y - s/2, W: s, H: s}, lerpColor(progressTrackerFrom, progressTrackerTo, float32(j%4)/3).Alpha(1-b), s/2)
					}
				}
			}
		})
		ink := lerpColor(t.Text, t.TextSubtle.Alpha(0.85), clamp32(strike, 0, 1))
		ui.Text(c, task).FontSize(15).TextColor(ink).SingleLine().DrawOver(func(p *ui.Painter, r ui.Rect) {
			if strike > 0.01 {
				y := r.Y + r.H*0.55
				p.Line(r.X-1, y, r.X-1+(r.W+2)*clamp32(strike, 0, 1), y, 1.5, ink)
			}
		})
	})
}

func init() {
	registerDemo(
		"progress-tracker",
		"A getting started checklist: ticking a task crosses it out, the bar springs forward and the count rolls up; finishing them all bursts every row.",
		func(c *ui.Context) {
			ProgressTracker(c, progressTrackerTasks)
		},
	)
}
