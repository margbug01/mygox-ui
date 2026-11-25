package bencho

import (
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
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
// are gone. It is coss's ToastStack with bencho's toasts in its look.
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
	holder := ui.Box(c)
	st := ui.Local(holder, "toast_stack", func() toastStackState { return toastStackState{} })
	if st.gone == [3]bool{true, true, true} {
		if wait := toastStackRefill - c.Now().Sub(st.lastGone); wait > 0 {
			c.After(wait)
		} else {
			st.gone = [3]bool{}
		}
	}
	// The toasts still up, newest last.
	var up []int
	var toasts []ui.Toast
	for i := range toastStackToasts {
		if !st.gone[i] {
			up = append(up, i)
			toasts = append(toasts, ui.Toast{ID: strconv.Itoa(i)})
		}
	}
	content := func(k int, reveal float32) {
		i := up[k]
		ts := toastStackToasts[i]
		col := t.TextSubtle
		if ts.hue != "" {
			col = ui.Hex(ts.hue)
		}
		ui.Row(c).Absolute().Left(0).Top(0).WidthPercent(100).Height(toastStackToast).PaddingX(16).Gap(12).AlignItems(ui.Center).
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
			ui.Box(c).Size(24, 24).Radius(12).Center().Cursor(ui.CursorPointer).
				OnClick(func() { st.gone[i], st.lastGone = true, c.Now() }).Children(func() {
				ui.Icon(c, icons.X()).Size(16, 16).TextColor(t.TextSubtle)
			})
		})
	}
	holder.Children(func() {
		coss.ToastStack(c, toasts, coss.ToastStackOptions{
			Width: toastStackW, Height: toastStackToast, Slots: len(toastStackToasts), Peek: o.Peek, Gap: o.Gap, Rest: 60,
			Corner: o.Corner, Background: t.Surface, Border: t.Border, Content: content,
		})
	})
	return holder
}

func init() {
	registerDemo(
		"toast-stack",
		"Three stacked toasts, the two behind peeking out above the front one; pointing at them springs them apart into a list.",
		func(c *ui.Context) { ToastStack(c) },
	)
}
