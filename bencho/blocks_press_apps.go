package bencho

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

func paAbs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func paDrawCircle(p *ui.Painter, cx, cy, r float32, col ui.Color) {
	if r <= 0 {
		return
	}
	p.Fill(ui.Rect{X: cx - r, Y: cy - r, W: r * 2, H: r * 2}, col, r)
}

func paStrokeCircle(p *ui.Painter, cx, cy, r, w float32, col ui.Color) {
	if r <= 0 || w <= 0 {
		return
	}
	p.Stroke(ui.Rect{X: cx - r, Y: cy - r, W: r * 2, H: r * 2}, col, r, w)
}

// =============================================================================
// 6. Step Player ("step-player")
// =============================================================================

// WalkthroughStep defines one slide inside StepPlayer.
type WalkthroughStep struct {
	Title   string
	Caption string
	Accent  ui.Color
}

// StepPlayerOptions configures StepPlayer.
type StepPlayerOptions struct {
	Steps        []WalkthroughStep
	StepDuration time.Duration
	Width        float32
	OnStepChange func(index int)
}

type stepPlayerState struct {
	index     int
	playing   bool
	startedAt time.Time
	changedAt time.Time
}

// StepPlayer renders an interactive walkthrough player with Play/Pause
// auto-advance, scrubber dots, Prev/Next buttons, and step captions.
func StepPlayer(c *ui.Context, opts ...StepPlayerOptions) ui.Element {
	var o StepPlayerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 294
	}
	if o.StepDuration <= 0 {
		o.StepDuration = 2400 * time.Millisecond
	}
	if len(o.Steps) == 0 {
		o.Steps = []WalkthroughStep{
			{Title: "Connect workspace", Caption: "Link your repository to index design tokens automatically.", Accent: ui.Hex("#6366f1")},
			{Title: "Tune spring curves", Caption: "Adjust stiffness and bounce with live 60fps GPU previews.", Accent: ui.Hex("#0ea5e9")},
			{Title: "Inspect dark mode", Caption: "Verify surface contrast and ink legibility in both themes.", Accent: ui.Hex("#10b981")},
			{Title: "Ship native UI", Caption: "Export zero-dependency Go components ready for production.", Accent: ui.Hex("#f59e0b")},
		}
	}

	root := ui.Column(c).
		Width(o.Width).
		Padding(12).
		Gap(10).
		Radius(16).
		Background(t.Surface).
		Border(1, t.Border).
		Shadow(0, 6, 18, 0, t.Shadow)

	now := c.Now()
	st := ui.Local(root, "pa_step_player", func() stepPlayerState {
		return stepPlayerState{index: 0, playing: true, startedAt: now, changedAt: now}
	})

	n := len(o.Steps)
	if st.index >= n {
		st.index = 0
	}

	setStep := func(next int) {
		next = ((next % n) + n) % n
		st.index = next
		st.startedAt = now
		st.changedAt = now
		if o.OnStepChange != nil {
			o.OnStepChange(next)
		}
	}

	stepFrac := float32(0)
	if st.playing {
		elapsed := now.Sub(st.startedAt)
		if elapsed >= o.StepDuration {
			setStep(st.index + 1)
			stepFrac = 0
		} else {
			c.AnimationFrame()
			stepFrac = clamp32(float32(elapsed.Seconds()/o.StepDuration.Seconds()), 0, 1)
		}
	}

	cur := o.Steps[st.index]
	accent := colorOr(cur.Accent, t.Accent)

	// Brief caption entrance nudge on step change.
	slideTarget := float32(0)
	if now.Sub(st.changedAt) < 120*time.Millisecond {
		c.AnimationFrame()
		slideTarget = 6
	}
	slideY := root.Spring("pa_sp_slide", slideTarget, ui.SpringOptions{
		Duration: 220 * time.Millisecond,
		Bounce:   0.25,
	})

	root.Children(func() {
		// Top row: step counter pill + title.
		header := ui.Row(c).
			WidthPercent(100).
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween)
		header.Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(7).Children(func() {
				ui.Box(c).Width(8).Height(8).Radius(4).Background(accent)
				ui.Text(c, cur.Title).
					FontSize(13).
					FontWeight(600).
					TextColor(t.Text).
					Unselectable()
			})
			ui.Row(c).
				Padding(2, 7).
				Radius(999).
				Background(accent.Alpha(0.14)).
				Children(func() {
					ui.Text(c, fmt.Sprintf("%02d / %02d", st.index+1, n)).
						FontSize(10).
						FontWeight(600).
						TextColor(accent).
						Unselectable()
				})
		})

		// Caption card with live progress underline.
		capBox := ui.Column(c).
			WidthPercent(100).
			Height(56).
			Padding(9, 10).
			Radius(10).
			Background(t.SurfaceSubtle).
			Border(1, t.Border).
			Clip().
			Justify(ui.Center)
		capBox.Draw(func(p *ui.Painter, r ui.Rect) {
			barH := float32(2.5)
			p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - barH, W: r.W, H: barH}, t.Border.Alpha(0.6), 0)
			if st.playing {
				p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - barH, W: r.W * stepFrac, H: barH}, accent, 0)
			}
		})
		capBox.Children(func() {
			ui.Text(c, cur.Caption).
				FontSize(11).
				TextColor(t.TextMuted).
				Translate(0, slideY).
				Unselectable()
		})

		// Transport bar: Prev, Play/Pause, Scrubber dots, Next.
		bar := ui.Row(c).
			WidthPercent(100).
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween)
		bar.Children(func() {
			navBtns := ui.Row(c).AlignItems(ui.Center).Gap(5)
			navBtns.Children(func() {
				prevBtn := ui.Row(c).
					Key("sp_prev").
					Width(26).
					Height(26).
					Center().
					Radius(8).
					Background(t.SurfaceSubtle).
					Border(1, t.Border).
					Cursor(ui.CursorPointer)
				prevBtn.Children(func() {
					ui.Text(c, "<").FontSize(12).FontWeight(700).TextColor(t.Text).Unselectable()
				})
				if prevBtn.Clicked() {
					setStep(st.index - 1)
				}

				playLbl := "Play"
				playBg := accent
				playCol := t.AccentText
				if st.playing {
					playLbl = "Pause"
					playBg = t.SurfaceSubtle
					playCol = t.Text
				}
				playBtn := ui.Row(c).
					Key("sp_play").
					Height(26).
					Padding(0, 10).
					Center().
					Radius(8).
					Background(playBg).
					Border(1, lerpColor(t.Border, accent, func() float32 {
						if st.playing {
							return 0
						}
						return 1
					}())).
					Cursor(ui.CursorPointer)
				playBtn.Children(func() {
					ui.Text(c, playLbl).FontSize(11).FontWeight(600).TextColor(playCol).Unselectable()
				})
				if playBtn.Clicked() {
					st.playing = !st.playing
					st.startedAt = now
				}

				nextBtn := ui.Row(c).
					Key("sp_next").
					Width(26).
					Height(26).
					Center().
					Radius(8).
					Background(t.SurfaceSubtle).
					Border(1, t.Border).
					Cursor(ui.CursorPointer)
				nextBtn.Children(func() {
					ui.Text(c, ">").FontSize(12).FontWeight(700).TextColor(t.Text).Unselectable()
				})
				if nextBtn.Clicked() {
					setStep(st.index + 1)
				}
			})

			// Scrubber dots with spring width expansion on active dot.
			dots := ui.Row(c).AlignItems(ui.Center).Gap(5)
			dots.Children(func() {
				for i := 0; i < n; i++ {
					idx := i
					active := idx == st.index
					wTarget := float32(7)
					if active {
						wTarget = 20
					}
					dot := ui.Box(c).
						Key(fmt.Sprintf("sp_dot_%d", idx)).
						Height(7).
						Radius(4).
						Cursor(ui.CursorPointer)
					dw := dot.Spring(fmt.Sprintf("pa_sp_dw_%d", idx), wTarget, ui.SpringOptions{
						Duration: 240 * time.Millisecond,
						Bounce:   0.3,
					})
					col := t.BorderStrong
					if active {
						col = accent
					}
					dot.Width(dw).Background(col)
					if dot.Clicked() {
						setStep(idx)
					}
				}
			})
		})
	})

	return root
}

// =============================================================================
// Demo registrations
// =============================================================================

func init() {
	registerDemo(
		"step-player",
		"Walkthrough step player with Play/Pause auto-advance, scrubber dots, and Prev/Next controls.",
		func(c *ui.Context) {
			StepPlayer(c)
		},
	)
}
