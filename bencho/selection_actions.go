package bencho

import (
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// SelectionActionsOptions configures SelectionActions.
type SelectionActionsOptions struct {
	Think time.Duration // how long the rewrite takes to arrive, 300ms to 2s; 900ms by default
}

type selectionActionsState struct {
	asked    time.Time // when Improve was pressed
	improved bool
	note     bool
}

const (
	selectionActionsW    = 320
	selectionActionsLine = 22
	selectionActionsGap  = 3.85 // a space at 14px
)

// Lucide's message-circle and lightbulb (ISC licence), which coss/icons doesn't carry.
var (
	selectionActionsExplain = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><path d="M7.9 20A9 9 0 1 0 4 16.1L2 22Z"/></svg>`))
	selectionActionsImprove = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><path d="M15 14c.2-1 .7-1.7 1.5-2.5 1-.9 1.5-2.2 1.5-3.5A6 6 0 0 0 6 8c0 1 .2 2.2 1.5 3.5.7.7 1.3 1.5 1.5 2.5"/><path d="M9 18h6"/><path d="M10 22h4"/></svg>`))
)

// SelectionActions shows a paragraph with a sentence selected and a pill of AI
// actions under it. Improve thinks for a moment, blurs the selection away and
// brings the rewrite in out of a blur, word by word; Explain, or the chevron,
// opens the pill into a note about the selection.
func SelectionActions(c *ui.Context, before, selected, rewrite, explanation string, opts ...SelectionActionsOptions) ui.Element {
	var o SelectionActionsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Think == 0 {
		o.Think = 900 * time.Millisecond
	}
	t := ThemeFrom(c)
	col := ui.Column(c).Width(selectionActionsW).AlignItems(ui.Center)
	st := ui.Local(col, "selection_actions", func() selectionActionsState { return selectionActionsState{} })
	now := c.Now()
	highlight := ui.Hex("#d8d9ef")
	if t.Dark {
		highlight = ui.RGBA(140, 140, 255, 0.28)
	}

	// Where Improve is: thinking, then the rewrite arriving word by word.
	thinking := !st.asked.IsZero() && now.Sub(st.asked) < o.Think
	arriving := float32(0)
	if !st.asked.IsZero() && !thinking {
		st.improved = true
		arriving = float32(now.Sub(st.asked)-o.Think) / float32(time.Second)
	}
	sel := strings.Fields(selected)
	if st.improved {
		sel = strings.Fields(rewrite)
	}
	if thinking || (st.improved && arriving < 1.6) {
		c.AnimationFrame()
	}

	word := func(key any, w string, ink, bg ui.Color, blur, alpha, gap float32) {
		b := ui.Box(c).Key(key).Height(selectionActionsLine).Padding(0, gap, 0, 0).Background(bg).Opacity(alpha)
		b.Children(func() {
			ui.Text(c, w).FontSize(14).LineHeight(1.571).TextColor(ink.Alpha(1 - 0.5*blur))
			// A blur, faked with ghost copies spread about the word.
			for g := 0; blur > 0.02 && g < 4; g++ {
				a := float64(g) * math.Pi / 2
				ui.Text(c, w).Key(g).Absolute().Left(float32(math.Cos(a)) * 3 * blur).Top(float32(math.Sin(a)) * 3 * blur).
					FontSize(14).LineHeight(1.571).TextColor(ink.Alpha(0.22 * blur))
			}
		})
	}

	col.Children(func() {
		ui.Row(c).Key("text").Width(selectionActionsW).Wrap().Children(func() {
			for i, w := range strings.Fields(before) {
				word(i, w, t.TextSubtle, ui.Color{}, 0, 1, selectionActionsGap)
			}
			for i, w := range sel {
				blur, alpha := float32(0), float32(1)
				switch {
				case thinking:
					// The selection breathes while the rewrite is thought up.
					p := float32(0.5 + 0.5*math.Sin(phaseSeconds(now)*6-float64(i)*0.5))
					blur, alpha = 0.25+0.35*p, 0.75
				case st.improved:
					k := clamp32((arriving-float32(i)*0.06)/0.45, 0, 1)
					blur, alpha = 1-k, k
				}
				gap := float32(selectionActionsGap)
				if i == len(sel)-1 {
					gap = 0 // the highlight stops at the full stop
				}
				word([2]int{1, i}, w, t.Text, highlight, blur, alpha, gap)
			}
		})

		// The pill, or the note it opens into.
		k := clamp32(col.Spring("note", b2f(st.note), ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.2}), 0, 1.1)
		pill := ui.Column(c).Key("pill").Margin(18, 0, 0, 0).Width(lerp32(232, 300, k)).Radius(20).Background(t.Surface).Clip()
		pill.Children(func() {
			if k < 0.5 {
				ui.Row(c).Key("actions").Height(40).Padding(0, 10, 0, 10).AlignItems(ui.Center).Opacity(1 - k*2).Children(func() {
					action := func(key string, icon *ui.SVG, label string, fn func()) {
						a := ui.Row(c).Key(key).Height(32).Padding(0, 6).Radius(16).AlignItems(ui.Center).Cursor(ui.CursorPointer).OnClick(fn)
						h := a.Spring("hover", b2f(a.Hovered()), ui.SpringOptions{Duration: 140 * time.Millisecond})
						a.Background(t.SurfaceSubtle.Alpha(h)).Children(func() {
							ui.Icon(c, icon).Size(16, 16).TextColor(t.Text).Margin(0, 8, 0, 0)
							ui.Text(c, label).FontSize(13.6).TextColor(t.Text).SingleLine()
						})
					}
					openNote := func() { st.note = true }
					action("explain", selectionActionsExplain, "Explain", openNote)
					ui.Box(c).Width(12)
					label := "Improve"
					if thinking {
						label = "Improving"
					}
					action("improve", selectionActionsImprove, label, func() {
						if st.asked.IsZero() {
							st.asked = c.Now()
						}
					})
					ui.Box(c).Grow(1)
					more := ui.Box(c).Key("more").Size(28, 28).Center().Cursor(ui.CursorPointer).OnClick(openNote)
					more.Children(func() { ui.Icon(c, icons.ChevronRight()).Size(16, 16).TextColor(t.TextSubtle) })
				})
				return
			}
			ui.Column(c).Key("note").Padding(14, 16, 16, 16).Gap(8).Opacity(clamp32((k-0.5)*2, 0, 1)).Children(func() {
				ui.Row(c).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icons.Sparkles()).Size(15, 15).TextColor(t.TextSubtle).Margin(0, 8, 0, 0)
					ui.Text(c, "Note").FontSize(13).FontWeight(500).TextColor(t.TextSubtle).Grow(1)
					x := ui.Box(c).Size(24, 24).Center().Radius(12).Cursor(ui.CursorPointer).OnClick(func() { st.note = false })
					x.Children(func() { ui.Icon(c, icons.X()).Size(14, 14).TextColor(t.TextSubtle) })
				})
				ui.Text(c, explanation).FontSize(14).LineHeight(1.4).TextColor(t.Text)
			})
		})
	})
	return col
}

func init() {
	registerDemo(
		"selection-actions",
		"A pill of AI actions under selected text: Improve brings a rewrite in out of a blur, Explain opens a note.",
		func(c *ui.Context) {
			SelectionActions(c,
				"Pistachio holds the top slot all weekend.",
				"Churn it first thing Saturday so the batch has time to firm up before the afternoon rush.",
				"Churn it early Saturday so the batch firms up well before the afternoon rush.",
				"Freshly churned ice cream needs a few hours to harden, so an early batch is at its best by the busy afternoon.",
			)
		},
	)
}
