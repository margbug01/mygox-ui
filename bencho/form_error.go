package bencho

import (
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// FormErrorOptions configures FormError.
type FormErrorOptions struct {
	Placeholder string             // "you@studio.com" by default
	OnSubmit    func(email string) // called with an address that passes
}

const (
	formErrorW      = 300
	formErrorH      = 52
	formErrorTop    = 44
	formErrorCorner = 19
	formErrorBtn    = 40
)

type formErrorState struct {
	email   string
	problem string    // what is wrong with the last try, "" when none
	shakeAt time.Time // when the last wrong try came
	done    bool      // the last try passed
}

// formErrorCheck returns what is wrong with an email address, or "".
func formErrorCheck(email string) string {
	email = strings.TrimSpace(email)
	at := strings.IndexByte(email, '@')
	switch {
	case email == "":
		return "Enter your email first"
	case at <= 0 || strings.Count(email, "@") > 1:
		return "An email needs an @"
	case strings.ContainsAny(email, " ,;"):
		return "No spaces in an email"
	}
	domain := email[at+1:]
	if dot := strings.LastIndexByte(domain, '.'); dot <= 0 || dot >= len(domain)-2 {
		return "That domain looks unfinished"
	}
	return ""
}

// FormError is an email field with an arrow that checks it: a wrong address
// shakes the field, rings it red and says what is wrong under it, until the
// next key; a right one turns the arrow into a green check.
func FormError(c *ui.Context, opts ...FormErrorOptions) ui.Element {
	var o FormErrorOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	o.Placeholder = orDefault(o.Placeholder, "you@studio.com")
	t := ThemeFrom(c)
	now := c.Now()

	root := ui.Box(c).Size(formErrorW, 140)
	st := ui.Local(root, "form_error", func() formErrorState { return formErrorState{} })
	submit := func() {
		st.problem = formErrorCheck(st.email)
		st.done = st.problem == ""
		if !st.done {
			st.shakeAt = now
			return
		}
		if o.OnSubmit != nil {
			o.OnSubmit(strings.TrimSpace(st.email))
		}
	}

	// A wrong try shakes the field side to side, dying out in half a second.
	shake := float32(0)
	if since := now.Sub(st.shakeAt).Seconds(); !st.shakeAt.IsZero() && since < 0.6 {
		shake = float32(9 * math.Exp(-since*7) * math.Sin(since*2*math.Pi*6.5))
		c.AnimationFrame()
	}
	bad := root.Spring("bad", b2f(st.problem != ""), ui.SpringOptions{Duration: 220 * time.Millisecond})
	done := root.Spring("done", b2f(st.done), ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.35})

	root.Children(func() {
		field := ui.Box(c).Absolute().Left(0).Top(formErrorTop).Size(formErrorW, formErrorH).Translate(shake, 0).
			Radius(formErrorCorner).Background(t.Surface).Cursor(ui.CursorText)
		if bad > 0.01 {
			field.Border(1.5, t.Danger.Alpha(clamp32(bad, 0, 1)))
		}
		field.Children(func() {
			ui.TextInputBase(c, &st.email).Placeholder(o.Placeholder).PlaceholderColor(lerpColor(t.TextSubtle, t.Surface, 0.4)).
				FontSize(16).TextColor(t.Text).Background(ui.Color{}).Padding(0).
				Absolute().Left(16).Top((formErrorH-24)/2).Size(formErrorW-16-formErrorBtn-20, 24).
				OnTextKey(func(k ui.TextKey) bool {
					if k.Key == ui.KeyEnter {
						submit()
						return true
					}
					// Any other key puts the error away.
					st.problem, st.done = "", false
					return false
				}).
				OnChange(func() { st.problem, st.done = "", false })

			btn := ui.Box(c).Absolute().Left(formErrorW-6-formErrorBtn).Top(6).Size(formErrorBtn, formErrorBtn).
				Radius(formErrorCorner - 6).Cursor(ui.CursorPointer).KeepFocus().OnClick(submit)
			press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 160 * time.Millisecond})
			bg := lerpColor(t.Accent, t.Danger, clamp32(bad, 0, 1))
			bg = lerpColor(bg, t.Success, clamp32(done, 0, 1))
			btn.Background(bg).Scale(1-0.06*press, 1-0.06*press).Children(func() {
				ui.Icon(c, icons.ArrowRight()).Absolute().Left(11).Top(11).Size(18, 18).TextColor(t.AccentText).
					Opacity(clamp32(1-done, 0, 1)).Translate(done*10, 0)
				ui.Icon(c, icons.Check()).Absolute().Left(11).Top(11).Size(18, 18).TextColor(ui.Hex("#ffffff")).
					Opacity(clamp32(done, 0, 1)).Scale(lerp32(0.4, 1, done), lerp32(0.4, 1, done))
			})
		})

		// What is wrong drops in under the field; a pass says so in green.
		msg, col, show := st.problem, t.Danger, bad
		if st.done {
			msg, col, show = "You're on the list", t.Success, done
		}
		if msg != "" {
			ui.Row(c).Absolute().Left(16).Top(formErrorTop+formErrorH+10).Height(18).AlignItems(ui.Center).
				Opacity(clamp32(show, 0, 1)).Translate(0, (1-clamp32(show, 0, 1))*-6).Children(func() {
				icon := icons.CircleAlert()
				if st.done {
					icon = icons.CircleCheck()
				}
				ui.Icon(c, icon).Size(14, 14).TextColor(col).Margin(0, 6, 0, 0)
				ui.Text(c, msg).FontSize(13).FontWeight(500).TextColor(col)
			})
		}
	})
	return root
}

func init() {
	registerDemo(
		"form-error",
		"An email field whose arrow checks it: a wrong address shakes the field, rings it red and says what is wrong; a right one turns the arrow into a green check.",
		func(c *ui.Context) { FormError(c) },
	)
}
