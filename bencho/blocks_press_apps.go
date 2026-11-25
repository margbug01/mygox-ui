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
// 4. Asset Swap ("asset-swap")
// =============================================================================

// SwapAsset describes one side of an AssetSwap currency pair.
type SwapAsset struct {
	Symbol   string
	Name     string
	PriceUSD float64
	Color    ui.Color
}

// AssetSwapOptions configures AssetSwap.
type AssetSwapOptions struct {
	Base            SwapAsset
	Quote           SwapAsset
	Amounts         []float64
	Width           float32
	OnSwapDirection func(flipped bool)
}

type assetSwapState struct {
	flipped   bool
	rotTarget float32
	amountIdx int
	swappedAt time.Time
}

// AssetSwap renders a compact crypto/currency swap widget with an animated
// swap-direction flip button that exchanges top/bottom rows and recalculates
// the quoted output.
func AssetSwap(c *ui.Context, opts ...AssetSwapOptions) ui.Element {
	var o AssetSwapOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Base.Symbol == "" {
		o.Base = SwapAsset{Symbol: "USDC", Name: "USD Coin", PriceUSD: 1.0, Color: ui.Hex("#2775ca")}
	}
	if o.Quote.Symbol == "" {
		o.Quote = SwapAsset{Symbol: "ETH", Name: "Ethereum", PriceUSD: 2600.0, Color: ui.Hex("#627eea")}
	}
	if len(o.Amounts) == 0 {
		o.Amounts = []float64{1000, 2500, 500}
	}
	if o.Width <= 0 {
		o.Width = 286
	}

	root := ui.Column(c).
		Width(o.Width).
		Padding(10).
		Gap(6).
		Radius(16).
		Background(t.Surface).
		Border(1, t.Border).
		Shadow(0, 6, 18, 0, t.Shadow)

	st := ui.Local(root, "pa_asset_swap", func() assetSwapState {
		return assetSwapState{}
	})

	now := c.Now()
	topAsset, botAsset := o.Base, o.Quote
	payAmt := o.Amounts[st.amountIdx%len(o.Amounts)]
	if st.flipped {
		topAsset, botAsset = o.Quote, o.Base
		// Scale ETH input cleanly so numbers stay human-sized.
		payAmt = payAmt / 1000.0
	}
	recvAmt := 0.0
	if botAsset.PriceUSD > 0 {
		recvAmt = payAmt * (topAsset.PriceUSD / botAsset.PriceUSD)
	}

	formatAmt := func(a SwapAsset, v float64) string {
		if a.PriceUSD >= 100 {
			return fmt.Sprintf("%.4f", v)
		}
		return fmt.Sprintf("%.2f", v)
	}

	// Brief spring nudge on swap so the rows visibly exchange.
	nudgeTarget := float32(0)
	if !st.swappedAt.IsZero() && now.Sub(st.swappedAt) < 140*time.Millisecond {
		c.AnimationFrame()
		nudgeTarget = 5
	}
	nudge := root.Spring("pa_swap_nudge", nudgeTarget, ui.SpringOptions{
		Duration: 260 * time.Millisecond,
		Bounce:   0.35,
	})

	renderRow := func(key, role string, asset SwapAsset, amtStr string, dy float32, clickable bool) ui.Element {
		row := ui.Row(c).
			Key(key).
			WidthPercent(100).
			Height(58).
			Padding(8, 12).
			Radius(12).
			Background(t.SurfaceSubtle).
			Border(1, t.Border).
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween).
			Translate(0, dy)
		if clickable {
			row.Cursor(ui.CursorPointer)
		}
		row.Children(func() {
			left := ui.Row(c).AlignItems(ui.Center).Gap(9)
			left.Children(func() {
				badge := ui.Row(c).
					Width(30).
					Height(30).
					Radius(15).
					Center().
					Background(asset.Color)
				badge.Children(func() {
					sym := asset.Symbol
					if len(sym) > 1 {
						sym = sym[:1]
					}
					ui.Text(c, sym).
						FontSize(13).
						FontWeight(700).
						TextColor(ui.Hex("#ffffff")).
						Unselectable()
				})
				meta := ui.Column(c).Gap(1)
				meta.Children(func() {
					ui.Text(c, role).FontSize(10).TextColor(t.TextSubtle).Unselectable()
					ui.Text(c, asset.Symbol).FontSize(13).FontWeight(700).TextColor(t.Text).Unselectable()
				})
			})

			right := ui.Column(c).AlignItems(ui.End).Gap(1)
			right.Children(func() {
				ui.Text(c, amtStr).FontSize(15).FontWeight(700).TextColor(t.Text).Unselectable()
				usdVal := payAmt * topAsset.PriceUSD
				ui.Text(c, fmt.Sprintf("$%.0f", usdVal)).FontSize(10).TextColor(t.TextSubtle).Unselectable()
			})
		})
		return row
	}

	root.Children(func() {
		stack := ui.Box(c).WidthPercent(100).Height(124)
		stack.Children(func() {
			rows := ui.Column(c).WidthPercent(100).Gap(8)
			rows.Children(func() {
				topRow := renderRow("swap_top", "You pay", topAsset, formatAmt(topAsset, payAmt), nudge, true)
				if topRow.Clicked() {
					st.amountIdx = (st.amountIdx + 1) % len(o.Amounts)
				}
				renderRow("swap_bot", "You receive", botAsset, formatAmt(botAsset, recvAmt), -nudge, false)
			})

			// Floating circular flip button centered between rows.
			flipBtn := ui.Row(c).
				Key("swap_flip_btn").
				Absolute().
				Top(46).
				Left((o.Width-20-32)*0.5).
				Width(32).
				Height(32).
				ZIndex(4).
				Radius(16).
				Center().
				Background(t.Surface).
				Border(1.5, t.BorderStrong).
				Shadow(0, 3, 8, 0, t.Shadow).
				Cursor(ui.CursorPointer)

			rot := flipBtn.Spring("pa_swap_rot", st.rotTarget, ui.SpringOptions{
				Duration: 320 * time.Millisecond,
				Bounce:   0.32,
			})
			btnScaleTarget := float32(1)
			if flipBtn.Pressed() {
				btnScaleTarget = 0.88
			} else if flipBtn.Hovered() {
				btnScaleTarget = 1.08
			}
			bSc := flipBtn.Spring("pa_swap_sc", btnScaleTarget, ui.SpringOptions{
				Duration: 180 * time.Millisecond,
				Bounce:   0.25,
			})
			flipBtn.Scale(bSc, bSc).Rotate(rot)
			flipBtn.Draw(func(p *ui.Painter, r ui.Rect) {
				cx, cy := r.X+r.W*0.5, r.Y+r.H*0.5
				col := t.Text
				if flipBtn.Hovered() {
					col = t.Accent
				}
				// Up arrow on left, down arrow on right.
				var up, down ui.Path
				up.Cap(ui.RoundCap).Join(ui.RoundJoin).
					MoveTo(cx-3.5, cy+4.5).LineTo(cx-3.5, cy-4.5).
					MoveTo(cx-6.2, cy-1.8).LineTo(cx-3.5, cy-4.5).LineTo(cx-0.8, cy-1.8)
				down.Cap(ui.RoundCap).Join(ui.RoundJoin).
					MoveTo(cx+3.5, cy-4.5).LineTo(cx+3.5, cy+4.5).
					MoveTo(cx+0.8, cy+1.8).LineTo(cx+3.5, cy+4.5).LineTo(cx+6.2, cy+1.8)
				p.StrokePath(&up, 1.7, col)
				p.StrokePath(&down, 1.7, col)
			})
			if flipBtn.Clicked() {
				st.flipped = !st.flipped
				st.rotTarget += 180
				st.swappedAt = now
				if o.OnSwapDirection != nil {
					o.OnSwapDirection(st.flipped)
				}
			}
		})

		// Rate footer line.
		rate := topAsset.PriceUSD / botAsset.PriceUSD
		rateStr := fmt.Sprintf("1 %s = %.4f %s", topAsset.Symbol, rate, botAsset.Symbol)
		if rate >= 10 {
			rateStr = fmt.Sprintf("1 %s = %.2f %s", topAsset.Symbol, rate, botAsset.Symbol)
		}
		footer := ui.Row(c).
			WidthPercent(100).
			Padding(2, 4, 0, 4).
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween)
		footer.Children(func() {
			ui.Text(c, rateStr).FontSize(10).TextColor(t.TextSubtle).Unselectable()
			ui.Text(c, "Tap amount to cycle").FontSize(10).TextColor(t.Accent).Unselectable()
		})
	})

	return root
}

// =============================================================================
// 5. Generate ("generate")
// =============================================================================

// GenerateOptions configures GenerateButton.
type GenerateOptions struct {
	Label           string
	GeneratingLabel string
	TagSets         [][]string
	Duration        time.Duration
	OnComplete      func(tags []string)
}

type generateState struct {
	generating bool
	startedAt  time.Time
	doneAt     time.Time
	setIdx     int
	completed  int
}

func paDrawSparkle(p *ui.Painter, cx, cy, rad float32, col ui.Color) {
	if rad <= 0 {
		return
	}
	in := rad * 0.28
	var path ui.Path
	path.MoveTo(cx, cy-rad).
		LineTo(cx+in, cy-in).
		LineTo(cx+rad, cy).
		LineTo(cx+in, cy+in).
		LineTo(cx, cy+rad).
		LineTo(cx-in, cy+in).
		LineTo(cx-rad, cy).
		LineTo(cx-in, cy-in).
		Close()
	p.FillPath(&path, col)
}

// GenerateButton renders an AI generation CTA button with an animated sparkle
// progress sweep while generating, cycling through crisp output tags on completion.
func GenerateButton(c *ui.Context, opts ...GenerateOptions) ui.Element {
	var o GenerateOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Label == "" {
		o.Label = "Generate"
	}
	if o.GeneratingLabel == "" {
		o.GeneratingLabel = "Generating..."
	}
	if o.Duration <= 0 {
		o.Duration = 1150 * time.Millisecond
	}
	if len(o.TagSets) == 0 {
		o.TagSets = [][]string{
			{"Spring physics", "Glass card", "Dark OLED"},
			{"Kinetic type", "Fluid tab", "Haptic pop"},
			{"Ambient glow", "Soft depth", "60fps GPU"},
			{"Micro badge", "Zero jank", "Crisp ink"},
		}
	}

	root := ui.Column(c).
		Width(286).
		AlignItems(ui.Center).
		Gap(12)

	st := ui.Local(root, "pa_generate", func() generateState {
		return generateState{}
	})

	now := c.Now()
	prog := float32(0)
	if st.generating {
		elapsed := now.Sub(st.startedAt)
		if elapsed >= o.Duration {
			st.generating = false
			st.doneAt = now
			st.completed++
			st.setIdx = st.completed % len(o.TagSets)
			if o.OnComplete != nil {
				o.OnComplete(o.TagSets[st.setIdx])
			}
		} else {
			c.AnimationFrame()
			prog = clamp32(float32(elapsed.Seconds()/o.Duration.Seconds()), 0, 1)
		}
	}

	justFinished := !st.doneAt.IsZero() && now.Sub(st.doneAt) < 420*time.Millisecond
	if justFinished {
		c.AnimationFrame()
	}

	root.Children(func() {
		btn := ui.Row(c).
			Key("pa_gen_btn").
			Width(196).
			Height(44).
			Padding(0, 18).
			Gap(9).
			Center().
			Radius(999).
			Clip().
			Cursor(ui.CursorPointer)

		scaleTarget := float32(1)
		if btn.Pressed() {
			scaleTarget = 0.95
		} else if btn.Hovered() {
			scaleTarget = 1.03
		}
		sc := btn.Spring("pa_gen_sc", scaleTarget, ui.SpringOptions{Duration: 180 * time.Millisecond, Bounce: 0.28})

		bg := t.Text
		textCol := t.Background
		borderCol := t.Text
		if st.generating {
			bg = t.SurfaceSubtle
			textCol = t.Text
			borderCol = t.Accent.Alpha(0.70)
		} else if justFinished {
			bg = t.Accent
			textCol = t.AccentText
			borderCol = t.Accent
		}

		btn.Background(bg).
			Border(1.2, borderCol).
			Shadow(0, 6, 16, 0, t.Shadow).
			Scale(sc, sc)

		btn.Draw(func(p *ui.Painter, r ui.Rect) {
			if !st.generating {
				return
			}
			// Progress fill bar.
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * prog, H: r.H}, t.Accent.Alpha(0.22), r.H*0.5)
			// Sweeping bright shimmer band.
			sweep := float32(mod1(phaseSeconds(now) * 1.45))
			bandW := r.W * 0.34
			bx := r.X - bandW + (r.W+bandW)*sweep
			p.FillGradient(ui.Rect{X: bx, Y: r.Y, W: bandW, H: r.H}, ui.LinearGradient{
				Angle: 90,
				Stops: []ui.GradientStop{
					{Color: ui.RGBA(255, 255, 255, 0), At: 0},
					{Color: t.Accent.Alpha(0.45), At: 0.5},
					{Color: ui.RGBA(255, 255, 255, 0), At: 1},
				},
			}, 0)
		})

		btn.Children(func() {
			// Sparkle icon box.
			spark := ui.Box(c).Width(16).Height(16)
			sparkRot := float32(0)
			if st.generating {
				sparkRot = float32(mod1(phaseSeconds(now)*1.2)) * 360
			}
			spark.Rotate(sparkRot).Draw(func(p *ui.Painter, r ui.Rect) {
				cx, cy := r.X+r.W*0.5, r.Y+r.H*0.5
				col := textCol
				if st.generating {
					col = t.Accent
				}
				paDrawSparkle(p, cx, cy, 6.5, col)
				paDrawSparkle(p, cx+5.2, cy-4.8, 2.2, col.Alpha(0.75))
			})

			lbl := o.Label
			if st.generating {
				lbl = fmt.Sprintf("%s %d%%", o.GeneratingLabel, int(prog*100))
			}
			ui.Text(c, lbl).
				FontSize(13).
				FontWeight(600).
				TextColor(textCol).
				Unselectable()
		})

		if btn.Clicked() && !st.generating {
			st.generating = true
			st.startedAt = now
		}

		// Generated output tag chips.
		tags := o.TagSets[st.setIdx]
		tagLiftTarget := float32(0)
		tagOpTarget := float32(1)
		if st.generating {
			tagLiftTarget = 6
			tagOpTarget = 0.42
		}
		tagY := root.Spring("pa_gen_tag_y", tagLiftTarget, ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.25})
		tagOp := root.Spring("pa_gen_tag_op", tagOpTarget, ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.1})

		tagRow := ui.Row(c).
			Gap(6).
			Center().
			Translate(0, tagY).
			Opacity(tagOp)
		tagRow.Children(func() {
			for i, tag := range tags {
				ui.Row(c).
					Key(fmt.Sprintf("gen_tag_%d", i)).
					Padding(4, 9).
					Radius(999).
					Background(t.SurfaceSubtle).
					Border(1, t.Border).
					Children(func() {
						ui.Text(c, tag).
							FontSize(10).
							FontWeight(500).
							TextColor(t.TextMuted).
							Unselectable()
					})
			}
		})

		status := "Click to synthesize"
		if st.generating {
			status = "Streaming tokens..."
		} else if st.completed > 0 {
			status = fmt.Sprintf("Set #%d ready · click for next", st.completed)
		}
		ui.Text(c, status).FontSize(10).TextColor(t.TextSubtle).Unselectable()
	})

	return root
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
		"asset-swap",
		"Currency swap widget with spring-rotating flip button that swaps pairs and recalculates quotes.",
		func(c *ui.Context) {
			AssetSwap(c)
		},
	)

	registerDemo(
		"generate",
		"AI generation CTA with sparkle progress sweep and cycling output tags.",
		func(c *ui.Context) {
			GenerateButton(c)
		},
	)

	registerDemo(
		"step-player",
		"Walkthrough step player with Play/Pause auto-advance, scrubber dots, and Prev/Next controls.",
		func(c *ui.Context) {
			StepPlayer(c)
		},
	)
}
