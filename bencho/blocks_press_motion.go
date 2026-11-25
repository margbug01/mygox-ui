package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

func pmSmoothStep(t float32) float32 {
	t = clamp32(t, 0, 1)
	return t * t * (3 - 2*t)
}

func pmEaseOutCubic(t float32) float32 {
	t = clamp32(t, 0, 1)
	f := 1 - t
	return 1 - f*f*f
}

func pmEaseOutQuart(t float32) float32 {
	t = clamp32(t, 0, 1)
	f := 1 - t
	return 1 - f*f*f*f
}

func pmRand(seed *uint32) float32 {
	x := *seed
	if x == 0 {
		x = 0x9e3779b9
	}
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	*seed = x
	return float32(x&0x00ffffff) / float32(0x01000000)
}

func pmPillButton(c *ui.Context, t Theme, label string, primary bool) ui.Element {
	bg := t.SurfaceSubtle
	fg := t.Text
	border := t.BorderStrong
	if primary {
		bg = t.Accent
		fg = t.AccentText
		border = t.Accent
	}
	btn := ui.Box(c).
		Height(24).
		Padding(0, 10).
		Radius(999).
		Center().
		Background(bg).
		BorderWidth(1, 1, 1, 1).
		BorderColor(border).
		Cursor(ui.CursorPointer)
	if btn.Hovered() {
		if primary {
			btn.Background(t.Accent.Alpha(0.88))
		} else {
			btn.Background(t.SurfaceHover)
		}
	}
	btn.Children(func() {
		ui.Text(c, label).
			FontSize(11).
			FontWeight(600).
			TextColor(fg).
			SingleLine()
	})
	return btn
}

// -----------------------------------------------------------------------------
// 12. ThinkingBlock ("thinking")
// -----------------------------------------------------------------------------

// ThinkingBlockOptions configures ThinkingBlock.
type ThinkingBlockOptions struct {
	Width float32
	Steps []string
}

type pmThinkingState struct {
	thinking  bool
	expanded  bool
	startedAt time.Time
	thoughtS  float64
}

// ThinkingBlock renders a collapsible AI "Thought for 1.8s" reasoning block
// with a pulsing shimmer indicator while thinking and click-to-expand trace steps.
func ThinkingBlock(c *ui.Context, opts ...ThinkingBlockOptions) ui.Element {
	var o ThinkingBlockOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}
	steps := o.Steps
	if len(steps) == 0 {
		steps = []string{
			"Parse card bounds & theme tokens",
			"Evaluate spring curves & stagger",
			"Compose accessible contrast pairs",
		}
	}

	root := ui.Column(c).
		Width(o.Width).
		Padding(12).
		Gap(8).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	st := ui.Local(root, "pm_thinking_state", func() pmThinkingState {
		return pmThinkingState{
			thinking: false,
			expanded: false,
			thoughtS: 1.8,
		}
	})

	const thinkDuration = 1.8
	liveElapsed := 0.0
	if st.thinking {
		liveElapsed = c.Now().Sub(st.startedAt).Seconds()
		if liveElapsed < 0 {
			liveElapsed = 0
		}
		if liveElapsed >= thinkDuration {
			st.thinking = false
			st.thoughtS = thinkDuration
		} else {
			c.AnimationFrame()
		}
	}

	openTarget := float32(0)
	if st.expanded {
		openTarget = 1
	}
	openP := clamp32(root.Spring("pm_think_open", openTarget, ui.SpringOptions{
		Duration: 340 * time.Millisecond,
		Bounce:   0.14,
	}), 0, 1)

	root.Children(func() {
		// Collapsible reasoning header bar.
		header := ui.Row(c).
			WidthPercent(100).
			Height(36).
			Padding(0, 10).
			Radius(10).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Background(t.SurfaceSubtle).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.Border).
			Cursor(ui.CursorPointer)

		if header.Hovered() {
			header.Background(t.SurfaceHover)
		}
		if header.Clicked() {
			st.expanded = !st.expanded
		}

		// Subtle sweeping shimmer bar across the header while thinking.
		header.Draw(func(p *ui.Painter, r ui.Rect) {
			if !st.thinking || r.W <= 0 || r.H <= 0 {
				return
			}
			p.AnimationFrame()
			sweep := float32(mod1(phaseSeconds(p.Now()) * 0.95))
			sx := r.X - 60 + sweep*(r.W+120)
			p.Clip(r, 10, func() {
				p.FillGradient(ui.Rect{X: sx - 45, Y: r.Y, W: 90, H: r.H}, ui.LinearGradient{
					Angle: 90,
					Stops: []ui.GradientStop{
						{Color: ui.Transparent, At: 0},
						{Color: t.Accent.Alpha(0.18), At: 0.5},
						{Color: ui.Transparent, At: 1},
					},
				}, 0)
			})
		})

		header.Children(func() {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				// Pulsing shimmer orb indicator.
				dot := ui.Box(c).Size(14, 14)
				dot.Draw(func(p *ui.Painter, r ui.Rect) {
					cx := r.X + r.W*0.5
					cy := r.Y + r.H*0.5
					if st.thinking {
						p.AnimationFrame()
						pulse := float32(0.5 + 0.5*math.Sin(phaseSeconds(p.Now())*6.2))
						ringR := 4.5 + pulse*2.8
						p.Fill(ui.Rect{X: cx - ringR, Y: cy - ringR, W: ringR * 2, H: ringR * 2}, t.Accent.Alpha(0.28*(1-pulse*0.4)), ringR)
						coreR := float32(3.4)
						p.Fill(ui.Rect{X: cx - coreR, Y: cy - coreR, W: coreR * 2, H: coreR * 2}, t.Accent, coreR)
					} else {
						coreR := float32(4.0)
						p.Fill(ui.Rect{X: cx - coreR, Y: cy - coreR, W: coreR * 2, H: coreR * 2}, t.Accent.Alpha(0.22), coreR)
						p.Fill(ui.Rect{X: cx - 2.3, Y: cy - 2.3, W: 4.6, H: 4.6}, t.Accent, 2.3)
					}
				})

				title := fmt.Sprintf("Thought for %.1fs", st.thoughtS)
				titleCol := t.Text
				if st.thinking {
					title = fmt.Sprintf("Thinking... %.1fs", liveElapsed)
					titleCol = t.Accent
				}
				ui.Text(c, title).
					FontSize(12).
					FontWeight(600).
					TextColor(titleCol).
					SingleLine()
			})

			// Rotating chevron.
			chev := ui.Box(c).Size(14, 14)
			chev.Draw(func(p *ui.Painter, r ui.Rect) {
				cx := r.X + r.W*0.5
				cy := r.Y + r.H*0.5
				ang := openP * 90
				p.Transform(ui.Rotation(ang), cx, cy, func() {
					var path ui.Path
					path.MoveTo(cx-2, cy-4)
					path.LineTo(cx+2.5, cy)
					path.LineTo(cx-2, cy+4)
					p.StrokePath(&path, 1.6, t.TextMuted)
				})
			})
		})

		// Expandable reasoning trace steps.
		if openP > 0.02 {
			traceH := openP * 78
			trace := ui.Box(c).
				WidthPercent(100).
				Height(traceH).
				Clip().
				Opacity(openP)

			trace.Children(func() {
				ui.Column(c).
					WidthPercent(100).
					Padding(4, 8, 4, 12).
					Gap(6).
					BorderWidth(0, 0, 0, 2).
					BorderColor(t.Accent.Alpha(0.45)).
					Children(func() {
						for i, step := range steps {
							ui.Row(c).Key(i).Gap(6).AlignItems(ui.Center).Children(func() {
								ui.Text(c, fmt.Sprintf("%d.", i+1)).
									FontSize(10).
									FontWeight(700).
									TextColor(t.Accent).
									SingleLine()
								ui.Text(c, step).
									FontSize(11).
									TextColor(t.TextMuted).
									SingleLine()
							})
						}
					})
			})
		}

		// Footer bar with concise summary & Re-think trigger.
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			hint := "Click header to inspect trace"
			if st.expanded {
				hint = "3 reasoning steps verified"
			}
			ui.Text(c, hint).
				FontSize(10).
				TextColor(t.TextSubtle).
				SingleLine()

			rethink := pmPillButton(c, t, "Think", st.thinking)
			if rethink.Clicked() {
				st.thinking = true
				st.startedAt = c.Now()
				c.AnimationFrame()
			}
		})
	})

	return root
}

func init() {
	registerDemo(
		"thinking",
		"Collapsible AI 'Thought for 1.8s' reasoning block with pulsing shimmer indicator and expandable trace steps.",
		func(c *ui.Context) {
			ThinkingBlock(c)
		},
	)
}
