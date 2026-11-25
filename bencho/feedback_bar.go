package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// FeedbackBarOptions configures FeedbackBar.
type FeedbackBarOptions struct {
	Corner   float32 // corner radius, 0 to 28px; 20 by default
	OnRate   func(helpful bool)
	OnReason func(reason string)
}

// FeedbackBarStep is where the bar is in asking.
type FeedbackBarStep int

const (
	FeedbackBarAsking    FeedbackBarStep = iota // Was this helpful?
	FeedbackBarWhy                              // opened after a thumbs down, asking why
	FeedbackBarThanks                           // said thanks
	FeedbackBarDismissed                        // closed with the X; it comes back after a moment
)

var feedbackBarReasons = []string{"Out of date", "Hard to follow", "Missing steps", "Something else"}

const (
	feedbackBarW, feedbackBarH = 300, 56
	feedbackBarOpenH           = 148
	feedbackBarReturn          = 1600 * time.Millisecond
)

type feedbackBarState struct {
	step FeedbackBarStep
	at   time.Time
}

// FeedbackBar asks whether a page was helpful. A thumbs up says thanks; a
// thumbs down opens the bar to ask why.
func FeedbackBar(c *ui.Context, opts ...FeedbackBarOptions) ui.Element {
	var o FeedbackBarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Corner = clamp32(o.Corner, 0, 28)
	t := ThemeFrom(c)

	holder := ui.Box(c).Size(feedbackBarW, feedbackBarOpenH)
	st := ui.Local(holder, "feedback_bar", func() feedbackBarState { return feedbackBarState{} })
	now := c.Now()
	if st.step == FeedbackBarDismissed {
		if now.Sub(st.at) >= feedbackBarReturn {
			*st = feedbackBarState{}
		} else {
			c.AnimationFrame()
		}
	}
	set := func(s FeedbackBarStep) { st.step, st.at = s, now }

	spring := ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.2}
	open := holder.Spring("fb_open", b2f(st.step == FeedbackBarWhy), spring)
	thanks := holder.Spring("fb_thanks", b2f(st.step == FeedbackBarThanks), ui.SpringOptions{Duration: 300 * time.Millisecond})
	gone := holder.Spring("fb_gone", b2f(st.step == FeedbackBarDismissed), ui.SpringOptions{Duration: 260 * time.Millisecond})
	h := feedbackBarH + (feedbackBarOpenH-feedbackBarH)*open

	holder.Children(func() {
		bar := ui.Box(c).Absolute().Left(0).Top((feedbackBarOpenH-feedbackBarH)/2-(h-feedbackBarH)/2).
			Size(feedbackBarW, h).Radius(o.Corner).Background(t.Surface).Clip().
			Opacity(clamp32(1-gone, 0, 1)).Scale(1-0.06*gone, 1-0.06*gone)
		bar.Children(func() {
			ask := clamp32(1-thanks, 0, 1)
			ui.Text(c, "Was this helpful?").FontSize(16).FontWeight(400).TextColor(t.Text).
				Absolute().Left(17).Top(feedbackBarH/2-11).Opacity(ask).Translate(0, -8*thanks)
			if thanks > 0.01 {
				ui.Box(c).Absolute().Left(17).Top(feedbackBarH/2-11).Row().Gap(8).Opacity(clamp32(thanks, 0, 1)).Translate(0, 8*(1-thanks)).
					Children(func() {
						ui.Icon(c, icons.CircleCheck()).Size(20, 20).TextColor(t.Text)
						ui.Text(c, "Thanks for your feedback").FontSize(16).FontWeight(400).TextColor(t.Text)
					})
			}
			if ask > 0.01 {
				feedbackBarThumb(c, t, 0, 171, ask, st.step == FeedbackBarAsking, func() {
					set(FeedbackBarThanks)
					if o.OnRate != nil {
						o.OnRate(true)
					}
				})
				feedbackBarThumb(c, t, 1, 214, ask, st.step == FeedbackBarAsking, func() {
					set(FeedbackBarWhy)
					if o.OnRate != nil {
						o.OnRate(false)
					}
				})
			}
			ui.Box(c).Absolute().Left(242).Top(feedbackBarH/2-14).Size(1, 28).Background(t.Border)
			x := ui.Box(c).Key("close").Absolute().Left(256).Top(feedbackBarH/2-14).Size(28, 28).Radius(14).Center().
				Cursor(ui.CursorPointer).Label("Dismiss").OnClick(func() { set(FeedbackBarDismissed) })
			xh := x.Spring("hover", b2f(x.Hovered()), ui.SpringOptions{Duration: 160 * time.Millisecond})
			x.Background(t.SurfaceSubtle.Alpha(xh)).Children(func() {
				ui.Icon(c, icons.X()).Size(20, 20).TextColor(lerpColor(t.TextSubtle, t.Text, xh))
			})
			if open > 0.01 {
				feedbackBarWhy(c, t, open, func(reason string) {
					set(FeedbackBarThanks)
					if o.OnReason != nil {
						o.OnReason(reason)
					}
				})
			}
		})
	})
	return holder
}

// feedbackBarThumb is the thumbs up (i 0) or down (i 1) button centred at cx.
func feedbackBarThumb(c *ui.Context, t Theme, i int, cx, alpha float32, active bool, onClick func()) {
	b := ui.Box(c).Key(i).Absolute().Left(cx-18).Top(feedbackBarH/2-18).Size(36, 36).Radius(18).Opacity(alpha)
	if active {
		b.Cursor(ui.CursorPointer).OnClick(onClick)
	}
	b.Label([]string{"Helpful", "Not helpful"}[i])
	hover := b.Spring("hover", b2f(active && b.Hovered()), ui.SpringOptions{Duration: 160 * time.Millisecond})
	press := b.Spring("press", b2f(active && b.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond, Bounce: 0.4})
	tilt := float32(-12)
	if i == 1 {
		tilt = 12
	}
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, t.SurfaceSubtle.Alpha(hover), r.W/2)
		col := lerpColor(t.TextSubtle, t.Text, hover)
		s := 1 - 0.12*press
		ic := ui.Rect{X: r.X + 8, Y: r.Y + 8, W: 20, H: 20}
		tr := ui.Rotation(tilt * hover).Then(ui.Scaling(s, s))
		if i == 1 {
			tr = ui.Scaling(1, -1).Then(tr)
		}
		p.Transform(tr, r.X+18, r.Y+18, func() { p.Icon(icons.ThumbsUp(), ic, col) })
	})
}

// feedbackBarWhy is the opened part: a question and chips with reasons.
func feedbackBarWhy(c *ui.Context, t Theme, open float32, choose func(string)) {
	ui.Box(c).Absolute().Left(17).Top(feedbackBarH).Width(feedbackBarW-34).Opacity(clamp32(open*1.4-0.4, 0, 1)).
		Translate(0, 6*(1-open)).Children(func() {
		ui.Text(c, "What went wrong?").FontSize(13).FontWeight(400).TextColor(t.TextSubtle)
		ui.Box(c).Row().Wrap().Gap(6).Margin(10, 0, 0, 0).Children(func() {
			for i, reason := range feedbackBarReasons {
				chip := ui.Box(c).Key(i).Padding(6, 11).Radius(14).Cursor(ui.CursorPointer).OnClick(func() { choose(reason) })
				hover := chip.Spring("hover", b2f(chip.Hovered()), ui.SpringOptions{Duration: 160 * time.Millisecond})
				chip.Background(lerpColor(t.SurfaceSubtle, t.SurfaceHover, hover)).Children(func() {
					ui.Text(c, reason).FontSize(13).FontWeight(400).TextColor(t.Text)
				})
			}
		})
	})
}

func init() {
	registerDemo(
		"feedback-bar",
		"Asks whether a page was helpful: a thumbs up says thanks, a thumbs down opens the bar to ask why.",
		func(c *ui.Context) { FeedbackBar(c) },
	)
}
