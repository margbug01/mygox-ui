package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// NotifyOptions configures Notify.
type NotifyOptions struct {
	Corner   float32 // 0 to 22px; 22 by default
	Label    string  // "Notify me" by default
	OnChange func(signedUp bool)
}

const (
	notifyH     = 44
	notifySwing = 1400 * time.Millisecond
)

var notifyBellFilled = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.268 21a2 2 0 0 0 3.464 0"/><path fill="currentColor" d="M3.262 15.326A1 1 0 0 0 4 17h16a1 1 0 0 0 .74-1.673C19.41 13.956 18 12.499 18 8A6 6 0 0 0 6 8c0 4.499-1.411 5.956-2.738 7.326"/></svg>`))

type notifyState struct {
	on bool
	at time.Time
}

// Notify is a notify-me button with a small bell. Press it and the bell
// swings while the button stays pressed in, to show you are signed up; press
// again to sign out.
func Notify(c *ui.Context, opts ...NotifyOptions) ui.Element {
	var o NotifyOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 22
	}
	o.Corner = clamp32(o.Corner, 0, 22)
	o.Label = orDefault(o.Label, "Notify me")
	t := ThemeFrom(c)

	btn := ui.Box(c).Height(notifyH).Padding(0, 20).Row().Center().Gap(9).Radius(o.Corner).Cursor(ui.CursorPointer).Label(o.Label)
	st := ui.Local(btn, "notify", func() notifyState { return notifyState{} })
	now := c.Now()
	btn.OnClick(func() {
		st.on, st.at = !st.on, now
		if o.OnChange != nil {
			o.OnChange(st.on)
		}
	})
	// Held down while pressed, and kept down once signed up.
	down := btn.Spring("notify_down", max(b2f(st.on), b2f(btn.Pressed())), ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.3})
	hover := btn.Spring("notify_hover", b2f(btn.Hovered()), ui.SpringOptions{Duration: 160 * time.Millisecond})
	s := 1 + 0.02*hover*(1-down) - 0.04*down
	btn.Background(lerpColor(t.Surface, t.SurfaceHover, down*0.45)).Border(1, t.BorderStrong.Alpha(down)).Scale(s, s)

	// The bell swings on a damped pendulum after a sign-up.
	angle := float32(0)
	if st.on && now.Sub(st.at) < notifySwing {
		k := now.Sub(st.at).Seconds()
		angle = float32(22 * math.Sin(k*2*math.Pi*2.6) * math.Exp(-k*3.2))
		c.AnimationFrame()
	}
	btn.Children(func() {
		ui.Box(c).Size(18, 18).Draw(func(p *ui.Painter, r ui.Rect) {
			bell := icons.Bell()
			if st.on {
				bell = notifyBellFilled
			}
			p.Transform(ui.Rotation(angle), r.X+r.W/2, r.Y+1, func() { p.Icon(bell, r, t.Text) })
		})
		ui.Text(c, o.Label).FontSize(14).FontWeight(400).TextColor(t.Text)
	})
	return btn
}

func init() {
	registerDemo(
		"toasts",
		"A notify-me button with a small bell: press it and the bell swings while the button stays pressed to show you are signed up.",
		func(c *ui.Context) { Notify(c) },
	)
}
