package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// CommandBarOptions configures CommandBar.
type CommandBarOptions struct {
	Corner float32 // 0 to 28; 28 by default
	Bridge float32 // how far the button stretches before it breaks free, 1 to 12; 3 by default
	Edge   float32 // how sharp the liquid boundary is, 12 to 30; 22 by default
	OnRun  func(command string)
}

const (
	commandBarW      = 300
	commandBarH      = 56
	commandBarFieldW = 274
	commandBarBtn    = 48
)

// The waveform's bars, as bencho.dev draws them.
var commandBarBars = [...]float32{5, 9, 12, 7, 10, 5}

type commandBarState struct {
	text      string
	listening bool
	since     time.Time
}

// CommandBar is a command field with a voice button fused into its end.
// Pressing the button pulls it out of the field on a liquid bridge that
// stretches and breaks, leaving it to listen as a round button of its own;
// pressing it again (or Escape) slides it back and it fuses in. Enter runs
// what was typed.
func CommandBar(c *ui.Context, opts ...CommandBarOptions) ui.Element {
	var o CommandBarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 28
	}
	o.Corner = clamp32(o.Corner, 0, 28)
	if o.Bridge == 0 {
		o.Bridge = 3
	}
	o.Bridge = clamp32(o.Bridge, 1, 12)
	if o.Edge == 0 {
		o.Edge = 22
	}
	o.Edge = clamp32(o.Edge, 12, 30)
	t := ThemeFrom(c)
	now := c.Now()

	root := ui.Box(c).Size(commandBarW, commandBarH)
	st := ui.Local(root, "command_bar", func() commandBarState { return commandBarState{} })
	toggle := func() { st.listening, st.since = !st.listening, now }

	p := root.Spring("out", b2f(st.listening), ui.SpringOptions{Duration: 620 * time.Millisecond, Bounce: 0.3})
	// Out, the button sits just past where the bridge breaks.
	reach := o.Bridge * 3
	fieldW := lerp32(commandBarFieldW, commandBarW-commandBarBtn-max(10, reach+4), clamp32(p, 0, 1))
	r := lerp32(20, commandBarBtn/2, clamp32(p, 0, 1))
	cx := lerp32(commandBarFieldW-28, commandBarW-commandBarBtn/2, p)
	if st.listening {
		c.AnimationFrame() // the bars move while it listens
	}
	listen := clamp32(root.Spring("listen", b2f(st.listening), ui.SpringOptions{Duration: 300 * time.Millisecond}), 0, 1)

	root.Children(func() {
		// The field and the button as one liquid shape: while they are
		// close a neck joins them, thinning as they part, until it breaks.
		ui.Box(c).Absolute().Left(0).Top(0).Size(commandBarW, commandBarH).PassThrough().Draw(func(pt *ui.Painter, rr ui.Rect) {
			fr := ui.Rect{X: rr.X, Y: rr.Y, W: fieldW, H: commandBarH}
			pt.Fill(fr, t.Surface, min(o.Corner, commandBarH/2))
			bx, cy := rr.X+cx, rr.Y+commandBarH/2
			pt.Fill(ui.Rect{X: bx - r, Y: cy - r, W: 2 * r, H: 2 * r}, t.Surface, r)
			gap := bx - r - fr.X - fr.W
			if gap > -r && gap < reach {
				k := clamp32(1-gap/reach, 0, 1) // 1 touching, 0 broken
				waist := r * k * k * lerp32(0.75, 0.35, (o.Edge-12)/18)
				x0, x1 := fr.X+fr.W-min(o.Corner, commandBarH/2)*0.5, bx-r*0.4
				mid := (x0 + x1) / 2
				hold := r * 0.85
				var neck ui.Path
				neck.MoveTo(x0, cy-hold).QuadTo(mid, cy-waist, x1, cy-hold).
					LineTo(x1, cy+hold).QuadTo(mid, cy+waist, x0, cy+hold).Close()
				pt.FillPath(&neck, t.Surface)
			}
			// The waveform: still at rest, moving while it listens.
			ink := lerpColor(lerpColor(t.TextSubtle, t.Text, 0.15), t.Danger, listen)
			phase := float64(now.Sub(st.since).Seconds())
			for i, h := range commandBarBars {
				if listen > 0 {
					h = lerp32(h, 4+10*float32(0.5+0.5*math.Sin(phase*9+float64(i)*1.3)), listen)
				}
				x := bx - 7.5 + float32(i)*3
				pt.Fill(ui.Rect{X: x - 0.75, Y: cy - h/2, W: 1.5, H: h}, ink, 0.75)
			}
		})

		ui.TextInputBase(c, &st.text).Placeholder("Ask anything...").PlaceholderColor(lerpColor(t.TextSubtle, t.Surface, 0.08)).
			FontSize(16).TextColor(t.Text).Background(ui.Color{}).Padding(0).
			Absolute().Left(20).Top((commandBarH-24)/2).Size(fieldW-20-44, 24).
			OnTextKey(func(k ui.TextKey) bool {
				switch k.Key {
				case ui.KeyEnter:
					if k.Text() == "" {
						return false
					}
					if o.OnRun != nil {
						o.OnRun(k.Text())
					}
					k.SetText("")
					return true
				case ui.KeyEscape:
					if st.listening {
						toggle()
						return true
					}
				}
				return false
			})

		btn := ui.Box(c).Absolute().Left(cx-r).Top(commandBarH/2-r).Size(2*r, 2*r).Radius(r).
			Cursor(ui.CursorPointer).KeepFocus().Label("Voice").OnClick(toggle)
		press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 120 * time.Millisecond})
		btn.Background(t.SurfaceSubtle.Alpha(max(b2f(btn.Hovered())*0.6*clamp32(p, 0, 1), press*0.5)))
	})
	return root
}

func init() {
	registerDemo(
		"command",
		"A command field with a voice button fused into its end: pressing it pulls the button out on a liquid bridge that stretches and breaks, and pressing again fuses it back.",
		func(c *ui.Context) { CommandBar(c) },
	)
}
