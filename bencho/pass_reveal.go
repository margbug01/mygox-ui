package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// PasswordRevealOptions configures PasswordReveal. bencho.dev's block has no
// props; OnToggle reports each press of the eye.
type PasswordRevealOptions struct {
	OnToggle func(revealed bool)
}

const (
	passwordRevealStagger = 45 * time.Millisecond // between one character and the next
	passwordRevealDot     = 7                     // bullet diameter
	passwordRevealSlot    = 12                    // bullet advance
	passwordRevealFont    = 16
)

type passwordRevealState struct {
	revealed  bool
	toggledAt time.Time
}

// PasswordReveal is a password field beside a round eye button. Press the eye
// and each dot turns into its letter in order; press again and they go back
// to dots from the last one.
func PasswordReveal(c *ui.Context, password string, opts ...PasswordRevealOptions) ui.Element {
	var o PasswordRevealOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	chars := []rune(password)
	n := len(chars)

	row := ui.Row(c).Gap(8).AlignItems(ui.Center)
	st := ui.Local(row, "password_reveal", func() passwordRevealState { return passwordRevealState{} })

	// Character i shows its letter once the wave reaches it: from the first
	// when revealing, from the last when hiding.
	shown := make([]bool, n)
	since := c.Now().Sub(st.toggledAt)
	var next time.Duration = -1
	for i := range chars {
		order := i
		if !st.revealed {
			order = n - 1 - i
		}
		at := time.Duration(order) * passwordRevealStagger
		reached := st.toggledAt.IsZero() || since >= at
		shown[i] = st.revealed == reached
		if !reached && (next < 0 || at-since < next) {
			next = at - since
		}
	}
	if next >= 0 {
		c.After(next)
	}

	row.Children(func() {
		field := ui.Box(c).Size(240, 52).Radius(26).Background(t.Surface).Label("Password")
		ks := make([]float32, n)
		for i := range chars {
			ks[i] = field.Spring(i, b2f(shown[i]), ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.25})
		}
		field.Draw(func(p *ui.Painter, r ui.Rect) {
			cy := r.Y + r.H/2
			x := r.X + 24
			p.Clip(r, r.H/2, func() {
				for i, ch := range chars {
					k := ks[i]
					sp := ui.Span{Text: string(ch), Size: passwordRevealFont, Weight: 500, Color: t.Text}
					gw, gh := p.MeasureText(0, sp)
					adv := lerp32(passwordRevealSlot, gw+1, clamp32(k, 0, 1))
					if d := passwordRevealDot * (1 - k); d > 0.2 {
						cx := x + passwordRevealDot/2
						p.Fill(ui.Rect{X: cx - d/2, Y: cy - d/2, W: d, H: d}, t.Text, d/2)
					}
					if a := clamp32(k, 0, 1); a > 0.01 {
						sp.Color = t.Text.Alpha(a)
						p.RichText(x, cy-gh/2+(1-k)*5, 0, sp)
					}
					x += adv
				}
			})
		})

		eye := icons.EyeOff()
		label := "Show password"
		if st.revealed {
			eye, label = icons.Eye(), "Hide password"
		}
		btn := ui.Box(c).Size(52, 52).Radius(26).Background(t.Surface).Center().Cursor(ui.CursorPointer).Label(label).OnClick(func() {
			st.revealed = !st.revealed
			st.toggledAt = c.Now()
			if o.OnToggle != nil {
				o.OnToggle(st.revealed)
			}
		})
		press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 180 * time.Millisecond})
		btn.Scale(1-0.06*press, 1-0.06*press)
		btn.Children(func() {
			ui.Icon(c, eye).Size(17, 17).TextColor(t.Text)
		})
	})
	return row
}

func init() {
	registerDemo(
		"pass-reveal",
		"A password field whose dots turn into letters one at a time when you press the eye, and back from the last one when you press again.",
		func(c *ui.Context) {
			PasswordReveal(c, "sunsetRain7")
		},
	)
}
