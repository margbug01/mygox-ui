package bencho

import (
	"strconv"
	"time"
	"unicode"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// PasswordStrengthOptions configures PasswordStrength.
type PasswordStrengthOptions struct {
	Min    int     // the shortest password the first rule accepts, 6 to 16; 8 by default
	Corner float32 // 0 to 23; 20 by default
}

const (
	passwordStrengthW     = 300
	passwordStrengthField = 48
	passwordStrengthRow   = 25.75
)

// The meter's colour for one to four rules met.
var passwordStrengthTones = [...]string{"#fb2c36", "#fe9a00", "#f0b100", "#00bc7d"}

type passwordStrengthState struct {
	password string
	show     bool
}

// passwordStrengthRules reports which of the four rules pw meets: length,
// upper and lower case, a number, a symbol.
func passwordStrengthRules(pw string, minLen int) [4]bool {
	var upper, lower, digit, symbol bool
	for _, r := range pw {
		switch {
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsLower(r):
			lower = true
		case unicode.IsDigit(r):
			digit = true
		case !unicode.IsSpace(r):
			symbol = true
		}
	}
	return [4]bool{len([]rune(pw)) >= minLen, upper && lower, digit, symbol}
}

// PasswordStrength is a password field over a meter that fills as you type,
// one bar per rule met, and the rules themselves, each ticked off once the
// password meets it. The eye shows what was typed.
func PasswordStrength(c *ui.Context, opts ...PasswordStrengthOptions) ui.Element {
	var o PasswordStrengthOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Min == 0 {
		o.Min = 8
	}
	o.Min = min(max(o.Min, 6), 16)
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Corner = clamp32(o.Corner, 0, 23)
	t := ThemeFrom(c)

	col := ui.Column(c).Width(passwordStrengthW).GapY(8)
	st := ui.Local(col, "password_strength", func() passwordStrengthState { return passwordStrengthState{} })
	rules := passwordStrengthRules(st.password, o.Min)
	met := 0
	for _, ok := range rules {
		met += int(b2f(ok))
	}
	faint := lerpColor(t.TextSubtle, t.Surface, 0.45)
	tone := t.BorderStrong
	if met > 0 {
		tone = ui.Hex(passwordStrengthTones[met-1])
	}

	col.Children(func() {
		field := ui.Box(c).Height(passwordStrengthField).Radius(min(o.Corner, passwordStrengthField/2)).Background(t.Surface).Cursor(ui.CursorText)
		field.Children(func() {
			in := ui.TextInputBase(c, &st.password).Placeholder("Password").PlaceholderColor(faint).
				FontSize(16).TextColor(t.Text).Background(ui.Color{}).Padding(0).
				Absolute().Left(17).Top((passwordStrengthField-24)/2).Size(passwordStrengthW-17-48, 24)
			if !st.show {
				in.Password()
			}
			eye := ui.Box(c).Absolute().Left(passwordStrengthW-13-18-7).Top((passwordStrengthField-32)/2).Size(32, 32).Radius(8).
				Center().Cursor(ui.CursorPointer).KeepFocus().OnClick(func() { st.show = !st.show })
			icon := icons.Eye()
			if st.show {
				icon = icons.EyeOff()
			}
			eye.Children(func() { ui.Icon(c, icon).Size(18, 18).TextColor(lerpColor(faint, t.TextMuted, b2f(eye.Hovered()))) })
		})

		card := ui.Box(c).Height(162).Radius(o.Corner).Background(t.Surface)
		card.Children(func() {
			// The meter: one bar per rule met, filling left to right.
			for i := range 4 {
				fill := card.Spring("bar"+strconv.Itoa(i), b2f(i < met), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.2})
				ui.Box(c).Key("bar"+strconv.Itoa(i)).Absolute().Left(16+float32(i)*53).Top(23).Size(49, 4).Radius(2).
					Background(t.SurfaceHover).Draw(func(p *ui.Painter, r ui.Rect) {
					if w := r.W * clamp32(fill, 0, 1.05); w > 0.5 {
						p.Fill(ui.Rect{X: r.X, Y: r.Y, W: min(w, r.W), H: r.H}, tone, 2)
					}
				})
			}

			labels := [4]string{strconv.Itoa(o.Min) + "+ characters", "Upper and lower case", "A number", "A symbol"}
			for i, label := range labels {
				tick := card.Spring("rule"+strconv.Itoa(i), b2f(rules[i]), ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.4})
				on := clamp32(tick, 0, 1)
				y := 60 + float32(i)*passwordStrengthRow
				ui.Box(c).Key("ring"+strconv.Itoa(i)).Absolute().Left(16).Top(y-7).Size(14, 14).Draw(func(p *ui.Painter, r ui.Rect) {
					p.Stroke(passwordStrengthInset(r, 0.75), lerpColor(faint, t.Success, on), r.W/2, 1.5)
					if on > 0 {
						s := lerp32(0.3, 1, tick)
						d := ui.Rect{X: r.X + r.W/2*(1-s), Y: r.Y + r.H/2*(1-s), W: r.W * s, H: r.H * s}
						p.Fill(d, t.Success.Alpha(on), d.W/2)
						p.Icon(icons.Check(), passwordStrengthInset(d, d.W*0.18), ui.Hex("#ffffff").Alpha(on))
					}
				})
				ui.Row(c).Key("label" + strconv.Itoa(i)).Absolute().Left(38).Top(y - 11).Height(22).AlignItems(ui.Center).Children(func() {
					ui.Text(c, label).FontSize(13.75).TextColor(lerpColor(faint, t.Text, on))
				})
			}
		})
	})
	return col
}

func passwordStrengthInset(r ui.Rect, d float32) ui.Rect {
	return ui.Rect{X: r.X + d, Y: r.Y + d, W: r.W - 2*d, H: r.H - 2*d}
}

func init() {
	registerDemo(
		"password-strength",
		"A password field over a meter that fills as you type, with each rule ticked off once the password meets it.",
		func(c *ui.Context) { PasswordStrength(c) },
	)
}
