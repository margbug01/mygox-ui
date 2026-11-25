package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// ToastStackOptions configures ToastStack.
type ToastStackOptions struct {
	Peek   float32 // px of each tucked toast that shows, 6 to 24; 14 by default
	Gap    float32 // px between toasts once fanned out, 4 to 20; 10 by default
	Corner float32 // corner radius, 0 to 32; 20 by default
}

// toastStackCircleX is Lucide's circle-x.
var toastStackCircleX = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="m15 9-6 6"/><path d="m9 9 6 6"/></svg>`))

// toastStackToasts run from the oldest, at the top of the fanned list, to the
// newest, in front.
var toastStackToasts = []struct {
	title, body string
	icon        func() *ui.SVG
	hue         string // "" for the muted ink
	retry       bool
}{
	{"brand-kit.zip uploaded", "Saved to Launch week.", icons.CircleCheck, "#30a46c", false},
	{"Kai joined your bench", "They can comment and copy.", icons.Info, "", false},
	{"Payment declined", "Try another card.", func() *ui.SVG { return toastStackCircleX }, "#e5484d", true},
}

const (
	toastStackW      = 300
	toastStackToast  = 64 // a toast's height
	toastStackRefill = 1200 * time.Millisecond
)

type toastStackState struct {
	gone     [3]bool
	lastGone time.Time
}

// ToastStack is three toasts at the edge: the newest in front, the two behind
// peeking out above it. Point at them and they spring apart into a list;
// close one with its cross and it slips away, and the stack refills once all
// are gone.
func ToastStack(c *ui.Context, opts ...ToastStackOptions) ui.Element {
	var o ToastStackOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Peek == 0 {
		o.Peek = 14
	}
	if o.Gap == 0 {
		o.Gap = 10
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Peek, o.Gap, o.Corner = clamp32(o.Peek, 6, 24), clamp32(o.Gap, 4, 20), clamp32(o.Corner, 0, 32)
	t := ThemeFrom(c)
	n := len(toastStackToasts)
	h := float32(n)*toastStackToast + float32(n-1)*o.Gap

	root := ui.Box(c).Size(toastStackW, h)
	st := ui.Local(root, "toast_stack", func() toastStackState { return toastStackState{} })
	if st.gone == [3]bool{true, true, true} {
		if wait := toastStackRefill - c.Now().Sub(st.lastGone); wait > 0 {
			c.After(wait)
		} else {
			st.gone = [3]bool{}
		}
	}
	fan := root.Spring("fan", b2f(root.Hovered()), ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.25})
	restFront := h - toastStackToast - 60 // the front toast's top at rest: 88 of 212 at the defaults

	// The toasts still up, newest last; depth counts back from the front one.
	var up []int
	for i := range n {
		if !st.gone[i] {
			up = append(up, i)
		}
	}
	root.Children(func() {
		for k, i := range up {
			depth := float32(len(up) - 1 - k)
			pos := float32(n - len(up) + k) // fanned slot, the list kept at the bottom
			rest := restFront - depth*o.Peek
			y := lerp32(rest, pos*(toastStackToast+o.Gap), fan)
			shrink := (1 - fan) * depth * 7
			toastStackOne(c, t, o, i, y, shrink, depth, fan, func() { st.gone[i], st.lastGone = true, c.Now() })
		}
	})
	return root
}

// toastStackOne draws toast i at y, narrowed by shrink px a side; a toast
// behind the front one (depth > 0) shows only its outline until fanned.
func toastStackOne(c *ui.Context, t Theme, o ToastStackOptions, i int, y, shrink, depth, fan float32, dismiss func()) {
	ts := toastStackToasts[i]
	reveal := float32(1)
	if depth > 0 {
		reveal = clamp32(fan, 0, 1)
	}
	box := ui.Box(c).Key(i).Absolute().Left(shrink).Top(y).Size(toastStackW-2*shrink, toastStackToast)
	// A toast slides in from above when it first shows.
	born := ui.Local(box, "born", func() bool { return false })
	enter := box.Spring("enter", b2f(*born), ui.SpringOptions{Duration: 360 * time.Millisecond, Bounce: 0.2})
	if !*born {
		*born = true
		c.AnimationFrame()
	}
	box.Translate(0, (1-enter)*-16).Opacity(clamp32(enter, 0, 1))
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, t.Surface.Alpha(reveal), o.Corner)
		if reveal < 1 {
			p.Stroke(r, t.Border.Alpha(1-reveal), o.Corner, 1)
		}
	})
	if reveal < 0.05 {
		return
	}
	box.Children(func() {
		col := t.TextSubtle
		if ts.hue != "" {
			col = ui.Hex(ts.hue)
		}
		ui.Row(c).Absolute().Left(0).Top(0).Size(toastStackW-2*shrink, toastStackToast).PaddingX(16).Gap(12).AlignItems(ui.Center).
			Opacity(reveal).Children(func() {
			ui.Icon(c, ts.icon()).Size(18, 18).TextColor(col)
			ui.Column(c).Grow(1).Gap(2).Children(func() {
				ui.Text(c, ts.title).FontSize(14).FontWeight(500).TextColor(t.Text)
				ui.Text(c, ts.body).FontSize(13).TextColor(t.TextSubtle)
			})
			if ts.retry {
				btn := ui.Box(c).Height(28).PaddingX(12).Radius(14).Background(t.SurfaceSubtle).Center().Cursor(ui.CursorPointer)
				press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
				btn.Scale(1-0.06*press, 1-0.06*press)
				btn.Children(func() { ui.Text(c, "Retry").FontSize(13).FontWeight(500).TextColor(t.Text) })
			}
			ui.Box(c).Size(24, 24).Radius(12).Center().Cursor(ui.CursorPointer).OnClick(dismiss).Children(func() {
				ui.Icon(c, icons.X()).Size(16, 16).TextColor(t.TextSubtle)
			})
		})
	})
}

func init() {
	registerDemo(
		"toast-stack",
		"Three stacked toasts, the two behind peeking out above the front one; pointing at them springs them apart into a list.",
		func(c *ui.Context) { ToastStack(c) },
	)
}
