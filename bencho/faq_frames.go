package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// FAQItem is one question and its answer.
type FAQItem struct {
	Question string
	Answer   string
}

// FAQFramesOptions configures FAQFrames.
type FAQFramesOptions struct {
	Bounce float32 // how far a frame's height runs past its answer as it opens, 0 to 16; 3 by default
	Dim    bool    // while one is open, the closed frames step back
}

type faqFramesState struct {
	open    int
	answerH map[int]float32 // measured answer heights
}

// FAQFrames is an FAQ where every question sits in its own frame. Opening one
// slides its answer in with a little bounce; only one is open at a time.
func FAQFrames(c *ui.Context, items []FAQItem, opts ...FAQFramesOptions) ui.Element {
	var o FAQFramesOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Bounce == 0 {
		o.Bounce = 3
	}
	t := ThemeFrom(c)
	list := ui.Column(c).Width(320).Gap(4)
	st := ui.Local(list, "faq_frames", func() faqFramesState { return faqFramesState{answerH: map[int]float32{}} })
	list.Children(func() {
		for i, it := range items {
			open := st.open == i
			frame := ui.Column(c).Key(i).WidthPercent(100).Radius(16).Background(t.Surface).Clip()
			k := frame.Spring("open", b2f(open), ui.SpringOptions{Duration: 420 * time.Millisecond})
			// The height runs past the answer by Bounce on the way open, then settles.
			over := float32(0)
			if open && k < 1 {
				over = o.Bounce * float32(math.Sin(math.Pi*float64(clamp32(k, 0, 1))))
			}
			frame.Height(56 + st.answerH[i]*clamp32(k, 0, 1.2) + over)
			if o.Dim && st.open >= 0 && !open {
				frame.Opacity(lerp32(1, 0.5, frame.Spring("dim", 1, ui.SpringOptions{Duration: 300 * time.Millisecond})))
			}
			frame.Children(func() {
				head := ui.Row(c).WidthPercent(100).Height(56).Padding(0, 20).AlignItems(ui.Center).Cursor(ui.CursorPointer).OnClick(func() {
					if st.open == i {
						st.open = -1
					} else {
						st.open = i
					}
				})
				head.Children(func() {
					ui.Text(c, it.Question).FontSize(15).FontWeight(500).TextColor(t.Text).Grow(1)
					ui.Icon(c, icons.ChevronDown()).Size(16, 16).TextColor(t.TextSubtle).Rotate(180 * k)
				})
				if k > 0.01 {
					ans := ui.Text(c, it.Answer).FontSize(15).LineHeight(1.34).TextColor(t.TextSubtle).Padding(2, 20, 18, 20).Absolute().Left(0).Top(56).Width(320).Opacity(clamp32(k, 0, 1))
					if h := ans.LayoutBounds().H; h > 0 {
						st.answerH[i] = h
					} else if st.answerH[i] == 0 {
						st.answerH[i] = 60
					}
				}
			})
		}
	})
	return list
}

var faqFramesDemoItems = []FAQItem{
	{Question: "Can I use it in client work?", Answer: "Yes. Anything you build, personal or commercial, with no credit needed."},
	{Question: "Does it work with Tailwind?", Answer: "Yes. Every block is plain code you can restyle with your own classes."},
	{Question: "Is there a dark mode?", Answer: "Yes. The blocks follow your theme, light or dark."},
	{Question: "How often are blocks added?", Answer: "New blocks land every week, and they're yours as they come."},
}

func init() {
	registerDemo(
		"faq-frames",
		"An FAQ where every question sits in its own frame; open one and its answer slides in with a little bounce.",
		func(c *ui.Context) { FAQFrames(c, faqFramesDemoItems) },
	)
}
