package bencho

import (
	"math"
	"strings"
	"time"
	"unicode"

	"github.com/egoist/mygo/ui"
)

// OneTimeCodeAnswer is what a full code gets back.
type OneTimeCodeAnswer string

const (
	OneTimeCodeAccept OneTimeCodeAnswer = "Accept"
	OneTimeCodeReject OneTimeCodeAnswer = "Reject"
)

// OneTimeCodeLength is how many digits the code has.
type OneTimeCodeLength string

const (
	OneTimeCode4 OneTimeCodeLength = "4"
	OneTimeCode6 OneTimeCodeLength = "6"
)

// OneTimeCodeOptions configures OneTimeCode.
type OneTimeCodeOptions struct {
	Answer OneTimeCodeAnswer // Accept by default
	Length OneTimeCodeLength // 4 by default
	Corner float32           // 0 to 22; 12 by default
	// OnComplete is called with each full code.
	OnComplete func(code string)
}

const (
	oneTimeCodeBoxW = 35
	oneTimeCodeBoxH = 44
	oneTimeCodeGap  = 8
)

type oneTimeCodeState struct {
	code     string
	want     bool              // a press asked for the focus
	answer   OneTimeCodeAnswer // what the full code got back, "" while it is not full
	answered time.Time         // when it got it
	clear    bool              // empty the field on the next build
}

// OneTimeCode is a code input that looks like separate boxes but behaves as
// one field: typing fills them in turn, pasting fills them all and Backspace
// steps back. A full code is accepted, the digits turning green in a wave,
// or rejected, the boxes shaking red before they empty.
func OneTimeCode(c *ui.Context, opts ...OneTimeCodeOptions) ui.Element {
	var o OneTimeCodeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Answer == "" {
		o.Answer = OneTimeCodeAccept
	}
	n := 4
	if o.Length == OneTimeCode6 {
		n = 6
	}
	if o.Corner == 0 {
		o.Corner = 12
	}
	o.Corner = clamp32(o.Corner, 0, 22)
	t := ThemeFrom(c)
	now := c.Now()

	w := float32(n)*oneTimeCodeBoxW + float32(n-1)*oneTimeCodeGap
	root := ui.Box(c).Size(w, oneTimeCodeBoxH).Cursor(ui.CursorText)
	st := ui.Local(root, "one_time_code", func() oneTimeCodeState { return oneTimeCodeState{} })
	root.OnClick(func() { st.want = true })

	// Keep digits only, as many as there are boxes; a pasted code fills them all.
	if d := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, st.code); d != st.code || len(d) > n {
		st.code = d[:min(len(d), n)]
	}
	switch {
	case st.answer == "" && len(st.code) == n:
		st.answer, st.answered = o.Answer, now
		if o.OnComplete != nil {
			o.OnComplete(st.code)
		}
	case st.answer == OneTimeCodeAccept && len(st.code) < n: // editing an accepted code starts over
		st.answer = ""
	case st.answer == OneTimeCodeReject && now.Sub(st.answered) > 900*time.Millisecond: // empties for another try
		st.answer, st.clear = "", true
	}
	since := now.Sub(st.answered)
	rejected := st.answer == OneTimeCodeReject
	shake := float32(0)
	if rejected && since < 600*time.Millisecond {
		s := since.Seconds()
		shake = float32(7 * math.Exp(-s*6) * math.Sin(s*2*math.Pi*7))
		c.AnimationFrame()
	}
	ok := root.Spring("ok", b2f(st.answer == OneTimeCodeAccept), ui.SpringOptions{Duration: 400 * time.Millisecond})
	bad := root.Spring("bad", b2f(rejected), ui.SpringOptions{Duration: 180 * time.Millisecond})

	root.Children(func() {
		field := ui.TextInputBase(c, &st.code).Absolute().Left(0).Top(0).Size(1, 1).Opacity(0).MaxLength(n)
		if st.clear {
			st.code, st.clear = "", false
		}
		if st.want {
			field.Focus()
			st.want = false
		}
		focused := field.Focused()
		if focused {
			c.AnimationFrame() // the caret blinks
		}

		for i := range n {
			filled := i < len(st.code)
			active := focused && i == min(len(st.code), n-1) && st.answer == ""
			box := ui.Box(c).Key(i).Absolute().Left(float32(i)*(oneTimeCodeBoxW+oneTimeCodeGap)+shake).Top(0).
				Size(oneTimeCodeBoxW, oneTimeCodeBoxH).Radius(o.Corner).Background(t.Surface)
			pop := box.Spring("pop", b2f(filled), ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.45})
			// Accepted, the digits turn green one after another.
			wave := clamp32((float32(since.Seconds())-float32(i)*0.06)/0.2, 0, 1) * ok
			ink := lerpColor(lerpColor(t.Text, t.Success, wave), t.Danger, bad)
			if bad > 0.01 {
				box.Border(1.5, t.Danger.Alpha(bad))
			}
			box.Children(func() {
				if filled {
					ui.Row(c).Size(oneTimeCodeBoxW, oneTimeCodeBoxH).Center().
						Scale(lerp32(0.5, 1, pop), lerp32(0.5, 1, pop)).Opacity(clamp32(pop*1.4, 0, 1)).
						Translate(0, -3*float32(math.Sin(float64(wave)*math.Pi))).Children(func() {
						ui.Text(c, st.code[i:i+1]).FontSize(20).FontWeight(600).TextColor(ink)
					})
				}
				if active && (!filled || i == n-1 && len(st.code) < n) {
					blink := float32(math.Cos(float64(now.UnixMilli()%1000)/1000*2*math.Pi)*0.5 + 0.5)
					ui.Box(c).Absolute().Left((oneTimeCodeBoxW-1.5)/2).Top(13).Size(1.5, 18).Background(t.Text).Opacity(0.25 + 0.75*blink)
				}
			})
		}
	})
	if st.answer != "" && since < time.Second {
		c.AnimationFrame()
	}
	return root
}

func init() {
	registerDemo(
		"one-time-code",
		"A code input that looks like four boxes but behaves as one field: typing, pasting and Backspace all work, and a full code is accepted in green or shaken off in red.",
		func(c *ui.Context) { OneTimeCode(c) },
	)
}
