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
// 10. SkeletonMorph ("skeleton-morph")
// -----------------------------------------------------------------------------

// SkeletonMorphOptions configures SkeletonMorph.
type SkeletonMorphOptions struct {
	Width float32
}

type pmSkeletonMorphState struct {
	loaded bool
}

// SkeletonMorph renders a shimmering loading skeleton card that smoothly
// cross-morphs into real populated profile and metric content on click.
func SkeletonMorph(c *ui.Context, opts ...SkeletonMorphOptions) ui.Element {
	var o SkeletonMorphOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}

	card := ui.Column(c).
		Size(o.Width, 178).
		Padding(12, 14).
		Justify(ui.SpaceBetween).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	st := ui.Local(card, "pm_skel_morph_state", func() pmSkeletonMorphState {
		return pmSkeletonMorphState{}
	})

	target := float32(0)
	if st.loaded {
		target = 1
	}
	m := clamp32(card.Spring("pm_skel_m", target, ui.SpringOptions{
		Duration: 460 * time.Millisecond,
		Bounce:   0.14,
	}), 0, 1)

	skelAlpha := 1 - m
	if skelAlpha > 0.02 {
		c.AnimationFrame()
	}

	card.Children(func() {
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			stateText := "Loading skeleton..."
			if st.loaded {
				stateText = "Profile ready"
			}
			ui.Text(c, stateText).
				FontSize(11).
				FontWeight(600).
				TextColor(t.TextMuted).
				SingleLine()

			btnLabel := "Load"
			if st.loaded {
				btnLabel = "Reset"
			}
			btn := pmPillButton(c, t, btnLabel, !st.loaded)
			if btn.Clicked() {
				st.loaded = !st.loaded
			}
		})

		body := ui.Box(c).
			WidthPercent(100).
			Height(118).
			Radius(10).
			Background(t.SurfaceSubtle).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.Border).
			Clip()

		// Shimmering skeleton bones layer drawn on Painter with fading opacity.
		body.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 || skelAlpha <= 0.01 {
				return
			}
			p.AnimationFrame()
			shimmer := float32(mod1(phaseSeconds(p.Now()) * 0.72))
			baseBone := t.BorderStrong.Alpha(0.55 * skelAlpha)
			hiBone := t.Surface.Alpha(0.85 * skelAlpha)
			if t.Dark {
				hiBone = ui.Hex("#ffffff").Alpha(0.16 * skelAlpha)
			}

			drawBone := func(rect ui.Rect, rad float32) {
				p.Fill(rect, baseBone, rad)
				sweepX := r.X - 60 + shimmer*(r.W+120)
				p.Clip(rect, rad, func() {
					p.FillGradient(ui.Rect{X: sweepX - 42, Y: rect.Y, W: 84, H: rect.H}, ui.LinearGradient{
						Angle: 90,
						Stops: []ui.GradientStop{
							{Color: ui.Transparent, At: 0},
							{Color: hiBone, At: 0.5},
							{Color: ui.Transparent, At: 1},
						},
					}, 0)
				})
			}

			// Avatar circle bone.
			drawBone(ui.Rect{X: r.X + 12, Y: r.Y + 14, W: 36, H: 36}, 18)
			// Name + subtitle bars.
			drawBone(ui.Rect{X: r.X + 58, Y: r.Y + 17, W: 112, H: 12}, 5)
			drawBone(ui.Rect{X: r.X + 58, Y: r.Y + 35, W: 146, H: 9}, 4)
			// Two metric cards.
			halfW := (r.W - 32) * 0.5
			drawBone(ui.Rect{X: r.X + 12, Y: r.Y + 64, W: halfW, H: 40}, 8)
			drawBone(ui.Rect{X: r.X + 20 + halfW, Y: r.Y + 64, W: halfW, H: 40}, 8)
		})

		// Real populated profile/metric content cross-morphing in.
		if m > 0.02 {
			body.Children(func() {
				ui.Column(c).
					WidthPercent(100).
					HeightPercent(100).
					Padding(12).
					Justify(ui.SpaceBetween).
					Opacity(m).
					Translate(0, lerp32(6, 0, m)).
					Scale(lerp32(0.97, 1, m), lerp32(0.97, 1, m)).
					Children(func() {
						ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
							ui.Box(c).
								Size(36, 36).
								Radius(18).
								Center().
								Background(t.Accent).
								Children(func() {
									ui.Text(c, "MK").
										FontSize(12).
										FontWeight(700).
										TextColor(t.AccentText)
								})
							ui.Column(c).Gap(2).Children(func() {
								ui.Text(c, "Mina Kovacs").
									FontSize(13).
									FontWeight(600).
									TextColor(t.Text).
									SingleLine()
								ui.Text(c, "Design Systems · Active").
									FontSize(10).
									TextColor(t.TextMuted).
									SingleLine()
							})
						})

						ui.Row(c).WidthPercent(100).Gap(8).Children(func() {
							metrics := [2]struct {
								val string
								lbl string
								col ui.Color
							}{
								{"99.4%", "Uptime", t.Success},
								{"18 ms", "Latency", t.Accent},
							}
							for i, mt := range metrics {
								ui.Column(c).
									Key(i).
									Grow(1).
									Padding(6, 10).
									Radius(8).
									Background(t.Surface).
									BorderWidth(1, 1, 1, 1).
									BorderColor(t.Border).
									Children(func() {
										ui.Text(c, mt.val).
											FontSize(13).
											FontWeight(700).
											TextColor(mt.col).
											SingleLine()
										ui.Text(c, mt.lbl).
											FontSize(10).
											TextColor(t.TextSubtle).
											SingleLine()
									})
							}
						})
					})
			})
		}
	})

	return card
}

// -----------------------------------------------------------------------------
// 11. LoadCards ("load-cards")
// -----------------------------------------------------------------------------

// LoadCardsOptions configures LoadCards.
type LoadCardsOptions struct {
	Width float32
}

type pmLoadCardsState struct {
	run     int
	startAt time.Time
}

// LoadCards renders a staggered cascade entrance of 4 mini metric cards that
// spring in one after another when clicking Reload.
func LoadCards(c *ui.Context, opts ...LoadCardsOptions) ui.Element {
	var o LoadCardsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}

	type metricSpec struct {
		label string
		value string
		delta string
		fill  float32
		col   ui.Color
	}
	metrics := [4]metricSpec{
		{"Latency", "14 ms", "-18%", 0.78, t.Accent},
		{"Frame rate", "120 fps", "Locked", 0.94, t.Success},
		{"Builds", "48", "+6 today", 0.64, ui.Hex("#0ea5e9")},
		{"Quality", "98%", "+2.4%", 0.92, t.Warning},
	}

	root := ui.Column(c).
		Size(o.Width, 178).
		Padding(12).
		Justify(ui.SpaceBetween).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	st := ui.Local(root, "pm_load_cards_state", func() pmLoadCardsState {
		return pmLoadCardsState{startAt: c.Now()}
	})

	elapsed := c.Now().Sub(st.startAt)
	if elapsed < 0 {
		elapsed = 0
	}
	cardSpring := ui.SpringOptions{Duration: 450 * time.Millisecond, Bounce: 0.28}
	const stagger = 85 * time.Millisecond
	if elapsed < 3*stagger+cardSpring.SettleTime() {
		c.AnimationFrame()
	}

	root.Children(func() {
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Telemetry cascade").
				FontSize(12).
				FontWeight(600).
				TextColor(t.Text).
				SingleLine()

			reload := pmPillButton(c, t, "Reload", true)
			if reload.Clicked() {
				st.run++
				st.startAt = c.Now()
				c.AnimationFrame()
			}
		})

		renderCard := func(i int) {
			mt := metrics[i]
			local := elapsed - time.Duration(i)*stagger
			p := float32(0)
			if local > 0 {
				p = float32(cardSpring.Progress(local))
			}
			alpha := clamp32(p*1.25, 0, 1)
			yOff := (1 - p) * 18
			sc := 0.88 + 0.12*clamp32(p, 0, 1.04)

			box := ui.Column(c).
				Key(i).
				Grow(1).
				Height(56).
				Padding(7, 9).
				Justify(ui.SpaceBetween).
				Radius(10).
				Background(t.SurfaceSubtle).
				BorderWidth(1, 1, 1, 1).
				BorderColor(t.Border).
				Opacity(alpha).
				Translate(0, yOff).
				Scale(sc, sc)

			box.Children(func() {
				ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
					ui.Text(c, mt.label).
						FontSize(10).
						TextColor(t.TextMuted).
						SingleLine()
					ui.Box(c).Size(6, 6).Radius(3).Background(mt.col)
				})

				ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
					ui.Text(c, mt.value).
						FontSize(13).
						FontWeight(700).
						TextColor(t.Text).
						SingleLine()
					ui.Text(c, mt.delta).
						FontSize(9).
						FontWeight(600).
						TextColor(mt.col).
						SingleLine()
				})

				// Mini progress bar fills as the card springs in.
				bar := ui.Box(c).
					WidthPercent(100).
					Height(3).
					Radius(2).
					Background(t.Border)
				fillFrac := mt.fill * clamp32(p, 0, 1)
				bar.Draw(func(painter *ui.Painter, r ui.Rect) {
					if r.W <= 0 || r.H <= 0 || fillFrac <= 0.01 {
						return
					}
					painter.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * fillFrac, H: r.H}, mt.col, 2)
				})
			})
		}

		ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
			ui.Row(c).WidthPercent(100).Gap(8).Children(func() {
				renderCard(0)
				renderCard(1)
			})
			ui.Row(c).WidthPercent(100).Gap(8).Children(func() {
				renderCard(2)
				renderCard(3)
			})
		})
	})

	return root
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
		"skeleton-morph",
		"Shimmering loading skeleton card that cross-morphs smoothly into populated profile and metric content.",
		func(c *ui.Context) {
			SkeletonMorph(c)
		},
	)
	registerDemo(
		"load-cards",
		"Staggered cascade entrance of 4 mini metric cards springing into place one after another on Reload.",
		func(c *ui.Context) {
			LoadCards(c)
		},
	)
	registerDemo(
		"thinking",
		"Collapsible AI 'Thought for 1.8s' reasoning block with pulsing shimmer indicator and expandable trace steps.",
		func(c *ui.Context) {
			ThinkingBlock(c)
		},
	)
}
