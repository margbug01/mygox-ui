package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

func loadEaseInOut(t float32) float32 {
	t = clamp32(t, 0, 1)
	if t < 0.5 {
		return 4 * t * t * t
	}
	f := -2*t + 2
	return 1 - (f*f*f)/2
}

func loadEaseOutCubic(t float32) float32 {
	t = clamp32(t, 0, 1)
	f := 1 - t
	return 1 - f*f*f
}

func loadSmoothStep(t float32) float32 {
	t = clamp32(t, 0, 1)
	return t * t * (3 - 2*t)
}

// -----------------------------------------------------------------------------
// 2. ThinkingLoader ("thinking-loader")
// -----------------------------------------------------------------------------

// ThinkingLoaderOptions configures ThinkingLoader.
type ThinkingLoaderOptions struct {
	Phrases []string
	Accent  ui.Color
}

type thinkingLoaderState struct {
	phase     float64
	lastFrame time.Time
}

// ThinkingLoader renders a pulsing AI reasoning pill with a rotating orbital
// dot cluster and cycling status phrases.
func ThinkingLoader(c *ui.Context, opts ...ThinkingLoaderOptions) ui.Element {
	var o ThinkingLoaderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	phrases := o.Phrases
	if len(phrases) == 0 {
		phrases = []string{"Analyzing", "Synthesizing", "Verifying"}
	}
	accent := colorOr(o.Accent, t.Accent)
	accentAlt := ui.Hex("#38bdf8")

	pill := ui.Row(c).
		AlignItems(ui.Center).
		Justify(ui.SpaceBetween).
		Gap(10).
		Height(44).
		Padding(0, 16).
		Radius(999).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.BorderStrong).
		Cursor(ui.CursorPointer)

	hovered := pill.Hovered()
	hoverTarget := float32(0)
	if hovered {
		hoverTarget = 1
	}
	energy := pill.Spring("tl_energy", hoverTarget, ui.SpringOptions{
		Duration: 280 * time.Millisecond,
		Bounce:   0.2,
	})

	st := ui.Local(pill, "tl_clock", func() thinkingLoaderState {
		return thinkingLoaderState{lastFrame: c.Now()}
	})
	now := c.Now()
	dt := now.Sub(st.lastFrame).Seconds()
	if dt < 0 || dt > 0.25 {
		dt = 1.0 / 60.0
	}
	st.lastFrame = now
	st.phase += dt * (1.0 + 0.8*float64(energy))
	phase := st.phase

	c.AnimationFrame()

	// Cycle phrase every 2.2s.
	phrasePeriod := 2.2
	phraseIdx := int(math.Floor(phase/phrasePeriod)) % len(phrases)
	if phraseIdx < 0 {
		phraseIdx += len(phrases)
	}
	localT := float32(mod1(phase / phrasePeriod))

	// Fade smoothly at phrase entry and exit.
	textAlpha := float32(1)
	if localT < 0.14 {
		textAlpha = loadEaseOutCubic(localT / 0.14)
	} else if localT > 0.86 {
		textAlpha = loadEaseOutCubic((1 - localT) / 0.14)
	}

	pill.Children(func() {
		// Left: rotating orbital dot cluster.
		orb := ui.Box(c).Size(28, 28)
		orb.Draw(func(p *ui.Painter, r ui.Rect) {
			p.AnimationFrame()
			cx := r.X + r.W*0.5
			cy := r.Y + r.H*0.5

			// Ambient pulsing core halo.
			pulse := float32(0.5 + 0.5*math.Sin(phase*3.6))
			haloR := 11 + pulse*3 + energy*2
			p.FillRadialGradient(ui.Rect{X: cx - haloR, Y: cy - haloR, W: haloR * 2, H: haloR * 2}, ui.RadialGradient{
				Center: ui.Point{X: 0.5, Y: 0.5},
				Radius: haloR,
				Stops: []ui.GradientStop{
					{Color: accent.Alpha(0.28 + 0.14*pulse), At: 0},
					{Color: ui.Transparent, At: 1},
				},
			}, haloR)

			// Faint orbital ring.
			ringR := float32(10.5)
			p.Stroke(ui.Rect{X: cx - ringR, Y: cy - ringR, W: ringR * 2, H: ringR * 2}, accent.Alpha(0.20), ringR, 1)

			// 6 orbiting dots across two eccentric orbital streams.
			for i := 0; i < 6; i++ {
				a := phase*2.6 + float64(i)*(math.Pi/3)
				orbitR := float32(9.2)
				tilt := float64(0.55)
				if i%2 == 1 {
					a = -phase*2.1 + float64(i)*(math.Pi/3)
					orbitR = 6.6
					tilt = -0.55
				}
				ox := float32(math.Cos(a)) * orbitR
				oy := float32(math.Sin(a)) * orbitR * float32(math.Cos(tilt))
				depth := float32(0.5 + 0.5*math.Sin(a))
				dotR := 1.7 + 1.35*depth
				dotCol := lerpColor(accentAlt, accent, depth)
				p.Fill(ui.Rect{
					X: cx + ox - dotR,
					Y: cy + oy - dotR,
					W: dotR * 2,
					H: dotR * 2,
				}, dotCol.Alpha(0.45+0.55*depth), dotR)
			}

			// Pulsing center nucleus.
			coreR := float32(2.6 + 0.6*pulse)
			p.Fill(ui.Rect{X: cx - coreR, Y: cy - coreR, W: coreR * 2, H: coreR * 2}, accent, coreR)
		})

		// Middle: cycling reasoning phrase.
		labelBox := ui.Box(c).Size(102, 22)
		labelBox.Draw(func(p *ui.Painter, r ui.Rect) {
			p.AnimationFrame()
			yOff := float32(0)
			if localT < 0.14 {
				yOff = (1 - textAlpha) * 4
			} else if localT > 0.86 {
				yOff = -(1 - textAlpha) * 4
			}
			p.Text(r.X, r.Y+3+yOff, phrases[phraseIdx], 13, t.Text.Alpha(textAlpha))
		})

		// Right: staggered pulsing ellipsis dots.
		dots := ui.Box(c).Size(24, 14)
		dots.Draw(func(p *ui.Painter, r ui.Rect) {
			p.AnimationFrame()
			cy := r.Y + r.H*0.5
			for i := 0; i < 3; i++ {
				wave := float32(0.5 + 0.5*math.Sin(phase*4.8-float64(i)*0.9))
				dr := 1.6 + 1.1*wave
				dx := r.X + 4 + float32(i)*8
				p.Fill(ui.Rect{X: dx - dr, Y: cy - dr, W: dr * 2, H: dr * 2}, accent.Alpha(0.28+0.72*wave), dr)
			}
		})
	})

	return pill
}

// -----------------------------------------------------------------------------
// 3. LoadingState ("loading-state")
// -----------------------------------------------------------------------------

// LoadingPhase represents the finite state machine stages of LoadingState.
type LoadingPhase int

const (
	LoadingPhaseIdle LoadingPhase = iota
	LoadingPhaseFetching
	LoadingPhaseProcessing
	LoadingPhaseReady
)

// LoadingStateOptions configures LoadingState.
type LoadingStateOptions struct {
	AutoLoop bool
	Width    float32
}

type loadingStateFSM struct {
	phase   LoadingPhase
	phaseAt time.Time
}

// LoadingState renders a multi-phase progress card transitioning smoothly
// across Idle, Fetching, Processing, and Ready, with click-to-replay.
func LoadingState(c *ui.Context, opts ...LoadingStateOptions) ui.Element {
	o := LoadingStateOptions{AutoLoop: true, Width: 264}
	if len(opts) > 0 {
		if opts[0].Width > 0 {
			o.Width = opts[0].Width
		}
		o.AutoLoop = opts[0].AutoLoop
	}
	t := ThemeFrom(c)

	card := ui.Column(c).
		Width(o.Width).
		Padding(14, 16).
		Gap(12).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border).
		Cursor(ui.CursorPointer)

	st := ui.Local(card, "loading_state_fsm", func() loadingStateFSM {
		return loadingStateFSM{phase: LoadingPhaseFetching, phaseAt: c.Now()}
	})

	now := c.Now()
	if card.Clicked() {
		st.phase = LoadingPhaseFetching
		st.phaseAt = now
	}

	elapsed := now.Sub(st.phaseAt).Seconds()
	durations := [4]float64{1.2, 1.5, 1.5, 2.0}
	phaseDur := durations[int(st.phase)%4]
	if o.AutoLoop && elapsed >= phaseDur {
		st.phase = LoadingPhase((int(st.phase) + 1) % 4)
		st.phaseAt = now
		elapsed = 0
		phaseDur = durations[int(st.phase)%4]
	}

	c.AnimationFrame()

	localFrac := clamp32(float32(elapsed/phaseDur), 0, 1)
	var targetProg float32
	var phaseLabel, rightLabel string
	var phaseCol ui.Color

	switch st.phase {
	case LoadingPhaseIdle:
		targetProg = 0.04
		phaseLabel = "Idle"
		rightLabel = "Click to start"
		phaseCol = t.TextSubtle
	case LoadingPhaseFetching:
		targetProg = 0.08 + 0.42*loadEaseInOut(localFrac)
		phaseLabel = "Fetching"
		rightLabel = fmt.Sprintf("%d%%", int(targetProg*100))
		phaseCol = t.Accent
	case LoadingPhaseProcessing:
		targetProg = 0.50 + 0.46*loadEaseInOut(localFrac)
		phaseLabel = "Processing"
		rightLabel = fmt.Sprintf("%d%%", int(targetProg*100))
		phaseCol = ui.Hex("#8b5cf6")
	case LoadingPhaseReady:
		targetProg = 1.0
		phaseLabel = "Ready"
		rightLabel = "Replay"
		phaseCol = t.Success
	}

	smoothProg := card.Spring("ls_prog", targetProg, ui.SpringOptions{
		Duration: 260 * time.Millisecond,
		Bounce:   0.08,
	})
	readyPop := card.Spring("ls_ready", func() float32 {
		if st.phase == LoadingPhaseReady {
			return 1
		}
		return 0
	}(), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.35})

	sec := phaseSeconds(now)

	card.Children(func() {
		// Header row: animated status glyph + phase name + right readout.
		ui.Row(c).
			FillWidth().
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween).
			Children(func() {
				ui.Row(c).AlignItems(ui.Center).Gap(9).Children(func() {
					glyph := ui.Box(c).Size(22, 22)
					glyph.Draw(func(p *ui.Painter, r ui.Rect) {
						p.AnimationFrame()
						cx := r.X + r.W*0.5
						cy := r.Y + r.H*0.5
						rad := float32(9.5)

						switch st.phase {
						case LoadingPhaseIdle:
							p.Stroke(ui.Rect{X: cx - rad, Y: cy - rad, W: rad * 2, H: rad * 2}, t.BorderStrong, rad, 1.8)
							p.Fill(ui.Rect{X: cx - 2.5, Y: cy - 2.5, W: 5, H: 5}, t.TextSubtle, 2.5)

						case LoadingPhaseFetching:
							p.Stroke(ui.Rect{X: cx - rad, Y: cy - rad, W: rad * 2, H: rad * 2}, phaseCol.Alpha(0.20), rad, 2)
							arc := (&ui.Path{}).Cap(ui.RoundCap)
							startA := sec * 5.2
							steps := 14
							for i := 0; i <= steps; i++ {
								a := startA + float64(i)/float64(steps)*math.Pi*1.15
								px := cx + float32(math.Cos(a))*rad
								py := cy + float32(math.Sin(a))*rad
								if i == 0 {
									arc.MoveTo(px, py)
								} else {
									arc.LineTo(px, py)
								}
							}
							p.StrokePath(arc, 2.2, phaseCol)

						case LoadingPhaseProcessing:
							p.Fill(ui.Rect{X: cx - rad, Y: cy - rad, W: rad * 2, H: rad * 2}, phaseCol.Alpha(0.14), rad)
							for i := 0; i < 3; i++ {
								a := sec*4.4 + float64(i)*(2*math.Pi/3)
								dx := cx + float32(math.Cos(a))*5.4
								dy := cy + float32(math.Sin(a))*5.4
								p.Fill(ui.Rect{X: dx - 2.2, Y: dy - 2.2, W: 4.4, H: 4.4}, phaseCol, 2.2)
							}

						case LoadingPhaseReady:
							sc := 0.72 + 0.28*readyPop
							rr := rad * sc
							p.Fill(ui.Rect{X: cx - rr, Y: cy - rr, W: rr * 2, H: rr * 2}, phaseCol.Alpha(0.18), rr)
							p.Stroke(ui.Rect{X: cx - rr, Y: cy - rr, W: rr * 2, H: rr * 2}, phaseCol, rr, 1.8)
							chk := (&ui.Path{}).Cap(ui.RoundCap).Join(ui.RoundJoin)
							chk.MoveTo(cx-4.2*sc, cy+0.3*sc).
								LineTo(cx-1.1*sc, cy+3.4*sc).
								LineTo(cx+4.6*sc, cy-3.2*sc)
							p.StrokePath(chk, 2.1, phaseCol)
						}
					})

					ui.Text(c, phaseLabel).
						FontSize(13).
						FontWeight(600).
						TextColor(t.Text)
				})

				ui.Text(c, rightLabel).
					FontSize(11.5).
					FontWeight(500).
					TextColor(phaseCol)
			})

		// Progress track with sweeping luminous sheen.
		track := ui.Box(c).FillWidth().Height(8)
		track.Draw(func(p *ui.Painter, r ui.Rect) {
			p.AnimationFrame()
			rad := r.H * 0.5
			p.Fill(r, t.SurfaceSubtle, rad)

			fillW := r.W * clamp32(smoothProg, 0.03, 1)
			fillRect := ui.Rect{X: r.X, Y: r.Y, W: fillW, H: r.H}
			p.Clip(fillRect, rad, func() {
				p.FillGradient(fillRect, ui.LinearGradient{
					Angle: 90,
					Stops: []ui.GradientStop{
						{Color: phaseCol.Alpha(0.82), At: 0},
						{Color: phaseCol, At: 1},
					},
				}, rad)
				if st.phase == LoadingPhaseFetching || st.phase == LoadingPhaseProcessing {
					sheenU := float32(mod1(sec * 0.95))
					sx := r.X - 36 + sheenU*(fillW+72)
					p.FillGradient(ui.Rect{X: sx, Y: r.Y, W: 36, H: r.H}, ui.LinearGradient{
						Angle: 90,
						Stops: []ui.GradientStop{
							{Color: ui.Transparent, At: 0},
							{Color: ui.Hex("#ffffff").Alpha(0.45), At: 0.5},
							{Color: ui.Transparent, At: 1},
						},
					}, 0)
				}
			})
		})

		// Phase step indicator pills.
		stepsRow := ui.Row(c).FillWidth().AlignItems(ui.Center).Justify(ui.SpaceBetween)
		stepNames := [4]string{"Idle", "Fetch", "Process", "Ready"}
		stepsRow.Children(func() {
			for i, name := range stepNames {
				active := LoadingPhase(i) == st.phase
				done := LoadingPhase(i) < st.phase
				dotCol := t.BorderStrong
				textCol := t.TextSubtle
				if active {
					dotCol = phaseCol
					textCol = t.Text
				} else if done {
					dotCol = t.Success
					textCol = t.TextMuted
				}
				ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
					ui.Box(c).Size(6, 6).Radius(3).Background(dotCol)
					ui.Text(c, name).FontSize(10).FontWeight(500).TextColor(textCol)
				})
			}
		})
	})

	return card
}

// -----------------------------------------------------------------------------
// 4. BlueprintLoader ("blueprint-loader")
// -----------------------------------------------------------------------------

// BlueprintLoaderOptions configures BlueprintLoader.
type BlueprintLoaderOptions struct {
	Width  float32
	Height float32
}

// BlueprintLoader renders a technical drafting wireframe grid with an animated
// sweeping scanline, drafting crosshairs, and progressive stroke reveal.
func BlueprintLoader(c *ui.Context, opts ...BlueprintLoaderOptions) ui.Element {
	o := BlueprintLoaderOptions{Width: 264, Height: 154}
	if len(opts) > 0 {
		if opts[0].Width > 0 {
			o.Width = opts[0].Width
		}
		if opts[0].Height > 0 {
			o.Height = opts[0].Height
		}
	}
	t := ThemeFrom(c)

	var bg, gridMinor, gridMajor, wireCol, scanCol ui.Color
	if t.Dark {
		bg = ui.Hex("#0a1424")
		gridMinor = ui.Hex("#38bdf8").Alpha(0.08)
		gridMajor = ui.Hex("#38bdf8").Alpha(0.18)
		wireCol = ui.Hex("#7dd3fc")
		scanCol = ui.Hex("#38bdf8")
	} else {
		bg = ui.Hex("#edf4ff")
		gridMinor = ui.Hex("#2563eb").Alpha(0.08)
		gridMajor = ui.Hex("#2563eb").Alpha(0.18)
		wireCol = ui.Hex("#1d4ed8")
		scanCol = ui.Hex("#2563eb")
	}

	stage := ui.Box(c).
		Size(o.Width, o.Height).
		Radius(12).
		Background(bg).
		BorderWidth(1, 1, 1, 1).
		BorderColor(wireCol.Alpha(0.32))

	now := phaseSeconds(c.Now())
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()

		p.Clip(r, 12, func() {
			// Fine architectural drafting grid.
			step := float32(16)
			for i := 1; float32(i)*step < r.W; i++ {
				x := r.X + float32(i)*step
				col := gridMinor
				if i%3 == 0 {
					col = gridMajor
				}
				p.Line(x, r.Y, x, r.Y+r.H, 1, col)
			}
			for j := 1; float32(j)*step < r.H; j++ {
				y := r.Y + float32(j)*step
				col := gridMinor
				if j%3 == 0 {
					col = gridMajor
				}
				p.Line(r.X, y, r.X+r.W, y, 1, col)
			}

			// Progressive reveal cycle (0 -> 1 draw, brief hold, loop).
			cycle := float32(mod1(now * 0.28))
			reveal := loadEaseInOut(clamp32(cycle/0.72, 0, 1))

			cx := r.X + r.W*0.5
			cy := r.Y + r.H*0.5
			frameW := float32(148)
			frameH := float32(82)
			fx := cx - frameW*0.5
			fy := cy - frameH*0.5

			// Dashed blueprint outer guide.
			p.StrokeDashed(ui.Rect{X: fx - 8, Y: fy - 8, W: frameW + 16, H: frameH + 16}, wireCol.Alpha(0.28), 8, 1)

			// Progressive wireframe chassis path.
			chassis := (&ui.Path{}).Cap(ui.RoundCap).Join(ui.RoundJoin)
			chassis.MoveTo(fx, fy).
				LineTo(fx+frameW, fy).
				LineTo(fx+frameW, fy+frameH).
				LineTo(fx, fy+frameH).
				Close()
			chassis.Trim(0, reveal)
			p.StrokePath(chassis, 1.6, wireCol.Alpha(0.92))

			// Internal geometric blueprint primitives (circle + diagonal truss).
			innerProg := loadEaseOutCubic(clamp32((cycle-0.18)/0.56, 0, 1))
			if innerProg > 0.01 {
				ring := (&ui.Path{}).Cap(ui.RoundCap)
				ring.Circle(cx, cy, 24)
				ring.Trim(0, innerProg)
				p.StrokePath(ring, 1.4, wireCol.Alpha(0.85))

				truss := (&ui.Path{}).Cap(ui.RoundCap)
				truss.MoveTo(fx+14, fy+frameH-14).
					LineTo(cx-24, cy).
					MoveTo(cx+24, cy).
					LineTo(fx+frameW-14, fy+14).
					MoveTo(fx+14, fy+14).
					LineTo(fx+42, fy+14)
				truss.Trim(0, innerProg)
				p.StrokePath(truss, 1.15, wireCol.Alpha(0.65))
			}

			// Corner registration crosshairs.
			corners := [4][2]float32{
				{fx, fy},
				{fx + frameW, fy},
				{fx + frameW, fy + frameH},
				{fx, fy + frameH},
			}
			for _, pt := range corners {
				p.Line(pt[0]-5, pt[1], pt[0]+5, pt[1], 1, scanCol.Alpha(0.75))
				p.Line(pt[0], pt[1]-5, pt[0], pt[1]+5, 1, scanCol.Alpha(0.75))
			}

			// Animated vertical sweeping scanline + trailing luminescence.
			scanU := float32(mod1(now * 0.42))
			scanX := r.X + scanU*r.W
			trailW := float32(44)
			p.FillGradient(ui.Rect{X: scanX - trailW, Y: r.Y, W: trailW, H: r.H}, ui.LinearGradient{
				Angle: 90,
				Stops: []ui.GradientStop{
					{Color: ui.Transparent, At: 0},
					{Color: scanCol.Alpha(0.24), At: 1},
				},
			}, 0)
			p.Line(scanX, r.Y, scanX, r.Y+r.H, 1.5, scanCol.Alpha(0.92))

			// Tracking crosshair dot along horizontal axis.
			trackY := cy + float32(math.Sin(now*2.4))*22
			p.Line(r.X, trackY, r.X+r.W, trackY, 0.8, scanCol.Alpha(0.22))
			p.Fill(ui.Rect{X: scanX - 3, Y: trackY - 3, W: 6, H: 6}, scanCol, 3)

			// Corner telemetry labels.
			p.Text(r.X+10, r.Y+8, "CAD // WIREFRAME", 9.5, wireCol.Alpha(0.72))
			p.Text(r.X+r.W-52, r.Y+r.H-18, fmt.Sprintf("%03d%%", int(reveal*100)), 10, scanCol.Alpha(0.88))
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 5. DotLoader ("dot-loader")
// -----------------------------------------------------------------------------

// DotLoaderOptions configures DotLoader.
type DotLoaderOptions struct {
	Accent ui.Color
}

type dotLoaderState struct {
	manualMode int
}

// DotLoader renders a 5x3 precision dot matrix animating choreographed wave,
// spiral, and ripple luminance patterns.
func DotLoader(c *ui.Context, opts ...DotLoaderOptions) ui.Element {
	var o DotLoaderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	accent := colorOr(o.Accent, t.Accent)

	wrap := ui.Column(c).
		AlignItems(ui.Center).
		Gap(10).
		Cursor(ui.CursorPointer)

	st := ui.Local(wrap, "dot_loader_mode", func() dotLoaderState {
		return dotLoaderState{}
	})
	if wrap.Clicked() {
		st.manualMode = (st.manualMode + 1) % 3
	}

	now := phaseSeconds(c.Now())
	c.AnimationFrame()

	mode := (st.manualMode + int(math.Floor(now/4.5))) % 3
	modeNames := [3]string{"Wave matrix", "Spiral orbit", "Center ripple"}

	// Spiral traversal rank for a 5x3 matrix (0..14).
	spiralRank := [3][5]int{
		{0, 1, 2, 3, 4},
		{11, 12, 13, 14, 5},
		{10, 9, 8, 7, 6},
	}

	wrap.Children(func() {
		stage := ui.Box(c).
			Size(164, 96).
			Radius(14).
			Background(t.Surface).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.Border)

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			p.AnimationFrame()

			cols, rows := 5, 3
			pitch := float32(24)
			startX := r.X + (r.W-float32(cols-1)*pitch)*0.5
			startY := r.Y + (r.H-float32(rows-1)*pitch)*0.5

			for row := 0; row < rows; row++ {
				for col := 0; col < cols; col++ {
					var lum float32
					switch mode {
					case 0: // Diagonal wave
						d := float64(col)*0.55 + float64(row)*0.65
						lum = float32(0.5 + 0.5*math.Sin(now*4.6-d))
					case 1: // Spiral chase
						order := float64(spiralRank[row][col]) / 15.0
						head := mod1(now * 0.72)
						dist := mod1(head - order)
						lum = clamp32(float32(1-dist*2.2), 0, 1)
						lum = lum * lum
					case 2: // Radial ripple from center (2, 1)
						dx := float64(col - 2)
						dy := float64(row - 1)
						dist := math.Sqrt(dx*dx + dy*dy)
						lum = float32(0.5 + 0.5*math.Sin(now*5.0-dist*1.45))
					}

					cx := startX + float32(col)*pitch
					cy := startY + float32(row)*pitch
					dotR := lerp32(3.0, 6.2, lum)

					if lum > 0.35 {
						glowR := dotR * 2.1
						p.FillRadialGradient(ui.Rect{
							X: cx - glowR, Y: cy - glowR, W: glowR * 2, H: glowR * 2,
						}, ui.RadialGradient{
							Center: ui.Point{X: 0.5, Y: 0.5},
							Radius: glowR,
							Stops: []ui.GradientStop{
								{Color: accent.Alpha(0.28 * lum), At: 0},
								{Color: ui.Transparent, At: 1},
							},
						}, glowR)
					}

					dotCol := lerpColor(t.BorderStrong, accent, lum)
					if lum > 0.82 {
						dotCol = lerpColor(dotCol, ui.Hex("#ffffff"), (lum-0.82)*1.8)
					}
					p.Fill(ui.Rect{
						X: cx - dotR, Y: cy - dotR, W: dotR * 2, H: dotR * 2,
					}, dotCol, dotR)
				}
			}
		})

		ui.Text(c, modeNames[mode]+" · click to switch").
			FontSize(11).
			FontWeight(500).
			TextColor(t.TextSubtle)
	})

	return wrap
}

// -----------------------------------------------------------------------------
// 6. ImageShimmer ("image-shimmer")
// -----------------------------------------------------------------------------

// ImageShimmerOptions configures ImageShimmer.
type ImageShimmerOptions struct {
	Width  float32
	Height float32
	Aspect string
}

// ImageShimmer renders a rounded media skeleton frame with a crisp diagonal
// luminous shimmer sweep and subtle aspect ratio badge.
func ImageShimmer(c *ui.Context, opts ...ImageShimmerOptions) ui.Element {
	o := ImageShimmerOptions{Width: 248, Height: 148, Aspect: "16:9"}
	if len(opts) > 0 {
		if opts[0].Width > 0 {
			o.Width = opts[0].Width
		}
		if opts[0].Height > 0 {
			o.Height = opts[0].Height
		}
		if opts[0].Aspect != "" {
			o.Aspect = opts[0].Aspect
		}
	}
	t := ThemeFrom(c)

	card := ui.Box(c).
		Size(o.Width, o.Height).
		Radius(14).
		Background(t.SurfaceSubtle).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	now := phaseSeconds(c.Now())
	c.AnimationFrame()

	card.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()

		p.Clip(r, 14, func() {
			// Base gradient depth.
			baseA := t.SurfaceSubtle
			baseB := t.SurfaceHover
			p.FillGradient(r, ui.LinearGradient{
				Angle: 145,
				Stops: []ui.GradientStop{
					{Color: baseA, At: 0},
					{Color: baseB, At: 1},
				},
			}, 14)

			// Centered placeholder mountain + sun icon silhouette.
			cx := r.X + r.W*0.5
			cy := r.Y + r.H*0.44
			iconCol := t.BorderStrong.Alpha(0.75)

			// Sun dot.
			p.Fill(ui.Rect{X: cx + 12, Y: cy - 18, W: 9, H: 9}, iconCol, 4.5)

			// Mountain peaks.
			mtn := &ui.Path{}
			mtn.MoveTo(cx-34, cy+16).
				LineTo(cx-12, cy-10).
				LineTo(cx+4, cy+6).
				LineTo(cx+16, cy-4).
				LineTo(cx+34, cy+16).
				Close()
			p.FillPath(mtn, iconCol)

			// Subtle skeleton caption bars at the bottom.
			p.Fill(ui.Rect{X: r.X + 14, Y: r.Y + r.H - 32, W: r.W * 0.48, H: 7}, t.BorderStrong.Alpha(0.55), 3.5)
			p.Fill(ui.Rect{X: r.X + 14, Y: r.Y + r.H - 20, W: r.W * 0.30, H: 6}, t.BorderStrong.Alpha(0.38), 3)

			// Crisp diagonal luminous shimmer band.
			sweep := float32(mod1(now * 0.52))
			bandW := float32(84)
			travel := r.W + r.H + bandW*2
			bandX := r.X - bandW - r.H*0.4 + sweep*travel

			shimmerPeak := ui.Hex("#ffffff").Alpha(0.14)
			shimmerMid := ui.Hex("#ffffff").Alpha(0.05)
			if !t.Dark {
				shimmerPeak = ui.Hex("#ffffff").Alpha(0.68)
				shimmerMid = ui.Hex("#ffffff").Alpha(0.24)
			}

			p.Transform(ui.Rotation(22), bandX+bandW*0.5, r.Y+r.H*0.5, func() {
				p.FillGradient(ui.Rect{
					X: bandX,
					Y: r.Y - 60,
					W: bandW,
					H: r.H + 120,
				}, ui.LinearGradient{
					Angle: 90,
					Stops: []ui.GradientStop{
						{Color: ui.Transparent, At: 0},
						{Color: shimmerMid, At: 0.28},
						{Color: shimmerPeak, At: 0.5},
						{Color: shimmerMid, At: 0.72},
						{Color: ui.Transparent, At: 1},
					},
				}, 0)
			})

			// Top-right aspect ratio badge.
			badgeW := float32(38)
			badgeH := float32(20)
			badge := ui.Rect{X: r.X + r.W - badgeW - 10, Y: r.Y + 10, W: badgeW, H: badgeH}
			p.Fill(badge, t.Surface.Alpha(0.84), 6)
			p.Stroke(badge, t.BorderStrong.Alpha(0.65), 6, 1)
			p.Text(badge.X+7.5, badge.Y+3.5, o.Aspect, 10.5, t.TextMuted)
		})
	})

	return card
}

// -----------------------------------------------------------------------------
// 7. ShimmerText ("shimmer-text")
// -----------------------------------------------------------------------------

// ShimmerTextOptions configures ShimmerText.
type ShimmerTextOptions struct {
	Text    string
	Subtext string
	Size    float32
}

// ShimmerText renders a typography line with a continuous metallic highlight
// band gliding smoothly across individual characters.
func ShimmerText(c *ui.Context, opts ...ShimmerTextOptions) ui.Element {
	o := ShimmerTextOptions{
		Text:    "Generating response...",
		Subtext: "Streaming neural tokens",
		Size:    18,
	}
	if len(opts) > 0 {
		if opts[0].Text != "" {
			o.Text = opts[0].Text
		}
		if opts[0].Subtext != "" {
			o.Subtext = opts[0].Subtext
		}
		if opts[0].Size > 0 {
			o.Size = opts[0].Size
		}
	}
	t := ThemeFrom(c)

	runes := []rune(o.Text)
	stage := ui.Box(c).
		Size(276, 86).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	now := phaseSeconds(c.Now())
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 || len(runes) == 0 {
			return
		}
		p.AnimationFrame()

		// Measure total line width so we can center cleanly.
		fullW, fullH := p.MeasureText(1000, ui.Span{Text: o.Text, Size: o.Size, Weight: 600})
		startX := r.X + (r.W-fullW)*0.5
		startY := r.Y + (r.H-fullH)*0.5 - 9

		// Sweeping highlight position in [startX - pad, startX + fullW + pad].
		sweepU := float32(mod1(now * 0.46))
		pad := float32(48)
		bandX := startX - pad + sweepU*(fullW+pad*2)
		sigma := float32(26)

		baseCol := t.TextSubtle
		peakCol := ui.Hex("#ffffff")
		midCol := lerpColor(t.Accent, ui.Hex("#e0e7ff"), 0.55)
		if !t.Dark {
			baseCol = ui.Hex("#9a9a94")
			peakCol = ui.Hex("#111111")
			midCol = t.Accent
		}

		// Build per-character spans colored by distance to the gliding highlight band.
		spans := make([]ui.Span, len(runes))
		prevW := float32(0)
		for i := range runes {
			prefix := string(runes[:i+1])
			curW, _ := p.MeasureText(1000, ui.Span{Text: prefix, Size: o.Size, Weight: 600})
			charCenter := startX + (prevW+curW)*0.5
			prevW = curW

			dist := charCenter - bandX
			glow := float32(math.Exp(-float64(dist*dist) / float64(2*sigma*sigma)))
			var col ui.Color
			if glow < 0.5 {
				col = lerpColor(baseCol, midCol, glow*2)
			} else {
				col = lerpColor(midCol, peakCol, (glow-0.5)*2)
			}
			spans[i] = ui.Span{
				Text:   string(runes[i]),
				Size:   o.Size,
				Weight: 600,
				Color:  col,
			}
		}

		p.RichText(startX, startY, fullW+8, spans...)

		// Subtle metallic glint line beneath the text synced to the same band.
		lineY := startY + fullH + 9
		lineW := fullW * 0.85
		lineX := r.X + (r.W-lineW)*0.5
		p.Fill(ui.Rect{X: lineX, Y: lineY, W: lineW, H: 1.5}, t.Border, 1)
		p.Clip(ui.Rect{X: lineX, Y: lineY, W: lineW, H: 1.5}, 1, func() {
			p.FillGradient(ui.Rect{X: bandX - 36, Y: lineY, W: 72, H: 1.5}, ui.LinearGradient{
				Angle: 90,
				Stops: []ui.GradientStop{
					{Color: ui.Transparent, At: 0},
					{Color: t.Accent.Alpha(0.85), At: 0.5},
					{Color: ui.Transparent, At: 1},
				},
			}, 1)
		})

		if o.Subtext != "" {
			subW, _ := p.MeasureText(1000, ui.Span{Text: o.Subtext, Size: 11})
			p.Text(r.X+(r.W-subW)*0.5, lineY+7, o.Subtext, 11, t.TextSubtle)
		}
	})

	return stage
}

// -----------------------------------------------------------------------------
// 8. WaveLoader ("wave-loader")
// -----------------------------------------------------------------------------

// WaveLoaderOptions configures WaveLoader.
type WaveLoaderOptions struct {
	Bars   int
	Accent ui.Color
}

// WaveLoader renders smooth vertical equalizer bars undulating in a continuous
// sine phase wave.
func WaveLoader(c *ui.Context, opts ...WaveLoaderOptions) ui.Element {
	o := WaveLoaderOptions{Bars: 9}
	if len(opts) > 0 {
		if opts[0].Bars > 0 {
			o.Bars = opts[0].Bars
		}
		o.Accent = opts[0].Accent
	}
	t := ThemeFrom(c)
	accent := colorOr(o.Accent, t.Accent)
	accentHi := ui.Hex("#38bdf8")

	stage := ui.Box(c).
		Size(212, 118).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border).
		Cursor(ui.CursorPointer)

	hovered := stage.Hovered()
	boostTarget := float32(0)
	if hovered {
		boostTarget = 1
	}
	boost := stage.Spring("wl_boost", boostTarget, ui.SpringOptions{
		Duration: 300 * time.Millisecond,
		Bounce:   0.25,
	})

	now := phaseSeconds(c.Now())
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()

		n := o.Bars
		barW := float32(8)
		gap := float32(8)
		totalW := float32(n)*barW + float32(n-1)*gap
		startX := r.X + (r.W-totalW)*0.5
		cy := r.Y + r.H*0.5

		for i := 0; i < n; i++ {
			norm := float64(i) / float64(max(1, n-1))
			// Bell envelope gives an organic audio spectrum silhouette.
			env := float32(0.55 + 0.45*math.Sin(norm*math.Pi))
			waveA := math.Sin(now*4.2 - float64(i)*0.56)
			waveB := math.Sin(now*2.3 + float64(i)*0.38)
			amp := clamp32(float32(0.5+0.34*waveA+0.18*waveB), 0.05, 1)

			minH := float32(14)
			maxH := float32(66 + 14*boost)
			h := minH + (maxH-minH)*amp*env
			x := startX + float32(i)*(barW+gap)
			y := cy - h*0.5

			barRect := ui.Rect{X: x, Y: y, W: barW, H: h}
			topCol := lerpColor(accent, accentHi, clamp32(float32(norm)*0.7+amp*0.3, 0, 1))
			botCol := accent.Alpha(0.72)

			p.FillGradient(barRect, ui.LinearGradient{
				Angle: 180,
				Stops: []ui.GradientStop{
					{Color: topCol, At: 0},
					{Color: botCol, At: 1},
				},
			}, barW*0.5)

			// Subtle bright top cap highlight.
			p.Fill(ui.Rect{
				X: x + 1.5,
				Y: y + 1.5,
				W: barW - 3,
				H: 3,
			}, ui.Hex("#ffffff").Alpha(0.32+0.28*amp), 1.5)
		}
	})

	return stage
}

// -----------------------------------------------------------------------------
// 9. BounceLoader ("bounce-loader")
// -----------------------------------------------------------------------------

// BounceLoaderOptions configures BounceLoader.
type BounceLoaderOptions struct {
	Count  int
	Accent ui.Color
}

// BounceLoader renders staggered bouncing spheres with squash-and-stretch
// deformation and reactive ground shadows.
func BounceLoader(c *ui.Context, opts ...BounceLoaderOptions) ui.Element {
	o := BounceLoaderOptions{Count: 4}
	if len(opts) > 0 {
		if opts[0].Count > 0 {
			o.Count = opts[0].Count
		}
		o.Accent = opts[0].Accent
	}
	t := ThemeFrom(c)
	palette := [4]ui.Color{
		colorOr(o.Accent, t.Accent),
		ui.Hex("#38bdf8"),
		ui.Hex("#a855f7"),
		ui.Hex("#10b981"),
	}

	stage := ui.Box(c).
		Size(216, 122).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	now := phaseSeconds(c.Now())
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()

		n := o.Count
		if n > 4 {
			n = 4
		}
		baseR := float32(10.5)
		spacing := float32(38)
		startX := r.X + (r.W-float32(n-1)*spacing)*0.5
		groundY := r.Y + r.H - 26
		jumpH := float32(48)

		// Subtle floor reflection line.
		p.Line(r.X+24, groundY+4, r.X+r.W-24, groundY+4, 1, t.Border.Alpha(0.75))

		for i := 0; i < n; i++ {
			u := float32(mod1(now*1.12 - float64(i)*0.17))
			var lift, sx, sy float32
			const contact = float32(0.16)

			if u < contact {
				// Ground impact squash phase.
				k := float32(math.Sin(float64(u/contact) * math.Pi))
				lift = 0
				sx = 1 + 0.26*k
				sy = 1 - 0.24*k
			} else {
				// Parabolic flight with velocity-coupled stretch.
				f := (u - contact) / (1 - contact)
				lift = 4 * f * (1 - f)
				vel := float32(math.Abs(float64(1 - 2*f)))
				stretch := vel * 0.14
				sx = 1 - stretch*0.75
				sy = 1 + stretch
			}

			cx := startX + float32(i)*spacing
			rx := baseR * sx
			ry := baseR * sy
			cy := groundY - lift*jumpH - ry

			// Ground shadow shrinks and softens with altitude.
			shW := baseR * (1.15 - 0.48*lift) * sx
			shH := float32(4.2 - 1.6*lift)
			shAlpha := float32(0.28 - 0.18*lift)
			if t.Dark {
				shAlpha = 0.52 - 0.30*lift
			}
			p.Fill(ui.Rect{
				X: cx - shW,
				Y: groundY - shH*0.5 + 2,
				W: shW * 2,
				H: shH,
			}, ui.Hex("#000000").Alpha(shAlpha), shH)

			// Glossy sphere with top-left specular highlight.
			col := palette[i%len(palette)]
			ballRect := ui.Rect{X: cx - rx, Y: cy - ry, W: rx * 2, H: ry * 2}
			p.FillRadialGradient(ballRect, ui.RadialGradient{
				Center: ui.Point{X: 0.34, Y: 0.30},
				Radius: max(rx, ry) * 1.2,
				Stops: []ui.GradientStop{
					{Color: ui.Hex("#ffffff").Alpha(0.92), At: 0},
					{Color: col, At: 0.42},
					{Color: lerpColor(col, ui.Hex("#000000"), 0.28), At: 1},
				},
			}, min(rx, ry))
		}
	})

	return stage
}

// -----------------------------------------------------------------------------
// 10. TraceLoader ("trace-loader")
// -----------------------------------------------------------------------------

// TraceLoaderOptions configures TraceLoader.
type TraceLoaderOptions struct {
	Hexagon bool
	Accent  ui.Color
}

type traceLoaderState struct {
	hexagon bool
}

// loadRoundRectPoint evaluates a point along a rounded-rectangle perimeter for u in [0, 1).
func loadRoundRectPoint(cx, cy, w, h, rad, u float32) (float32, float32) {
	u = float32(mod1(float64(u)))
	rad = clamp32(rad, 0, min(w, h)*0.5)
	sw := w - 2*rad
	sh := h - 2*rad
	arcLen := float32(math.Pi * 0.5 * float64(rad))
	perim := 2*sw + 2*sh + 4*arcLen
	d := u * perim

	hw := w * 0.5
	hh := h * 0.5

	// Top edge.
	if d < sw {
		return cx - hw + rad + d, cy - hh
	}
	d -= sw
	// Top-right corner.
	if d < arcLen {
		a := -math.Pi*0.5 + float64(d/arcLen)*math.Pi*0.5
		return cx + hw - rad + float32(math.Cos(a))*rad, cy - hh + rad + float32(math.Sin(a))*rad
	}
	d -= arcLen
	// Right edge.
	if d < sh {
		return cx + hw, cy - hh + rad + d
	}
	d -= sh
	// Bottom-right corner.
	if d < arcLen {
		a := float64(d/arcLen) * math.Pi * 0.5
		return cx + hw - rad + float32(math.Cos(a))*rad, cy + hh - rad + float32(math.Sin(a))*rad
	}
	d -= arcLen
	// Bottom edge.
	if d < sw {
		return cx + hw - rad - d, cy + hh
	}
	d -= sw
	// Bottom-left corner.
	if d < arcLen {
		a := math.Pi*0.5 + float64(d/arcLen)*math.Pi*0.5
		return cx - hw + rad + float32(math.Cos(a))*rad, cy + hh - rad + float32(math.Sin(a))*rad
	}
	d -= arcLen
	// Left edge.
	if d < sh {
		return cx - hw, cy + hh - rad - d
	}
	d -= sh
	// Top-left corner.
	a := math.Pi + float64(d/arcLen)*math.Pi*0.5
	return cx - hw + rad + float32(math.Cos(a))*rad, cy - hh + rad + float32(math.Sin(a))*rad
}

// loadHexPoint evaluates a point along a regular hexagon perimeter for u in [0, 1).
func loadHexPoint(cx, cy, rad, u float32) (float32, float32) {
	u = float32(mod1(float64(u)))
	segF := u * 6
	seg := int(segF) % 6
	frac := segF - float32(seg)
	a0 := -math.Pi/2 + float64(seg)*math.Pi/3
	a1 := -math.Pi/2 + float64(seg+1)*math.Pi/3
	x0, y0 := cx+float32(math.Cos(a0))*rad, cy+float32(math.Sin(a0))*rad
	x1, y1 := cx+float32(math.Cos(a1))*rad, cy+float32(math.Sin(a1))*rad
	return lerp32(x0, x1, frac), lerp32(y0, y1, frac)
}

// TraceLoader renders a geometric rounded-square or hexagon perimeter with a
// glowing comet tail tracing seamlessly around its path.
func TraceLoader(c *ui.Context, opts ...TraceLoaderOptions) ui.Element {
	var o TraceLoaderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	accent := colorOr(o.Accent, t.Accent)
	tailCol := ui.Hex("#38bdf8")

	wrap := ui.Column(c).
		AlignItems(ui.Center).
		Gap(8).
		Cursor(ui.CursorPointer)

	st := ui.Local(wrap, "trace_loader_shape", func() traceLoaderState {
		return traceLoaderState{hexagon: o.Hexagon}
	})
	if wrap.Clicked() {
		st.hexagon = !st.hexagon
	}

	now := phaseSeconds(c.Now())
	c.AnimationFrame()

	wrap.Children(func() {
		stage := ui.Box(c).
			Size(148, 114).
			Radius(14).
			Background(t.Surface).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.Border)

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.AnimationFrame()

			cx := r.X + r.W*0.5
			cy := r.Y + r.H*0.5

			eval := func(u float32) (float32, float32) {
				if st.hexagon {
					return loadHexPoint(cx, cy, 40, u)
				}
				return loadRoundRectPoint(cx, cy, 78, 78, 18, u)
			}

			// Base closed perimeter track.
			track := (&ui.Path{}).Cap(ui.RoundCap).Join(ui.RoundJoin)
			const trackSteps = 72
			for i := 0; i <= trackSteps; i++ {
				px, py := eval(float32(i) / float32(trackSteps))
				if i == 0 {
					track.MoveTo(px, py)
				} else {
					track.LineTo(px, py)
				}
			}
			track.Close()
			p.StrokePath(track, 2.2, t.BorderStrong.Alpha(0.55))

			// Glowing comet tail tracing seamlessly across u=1->0 wraparound.
			headU := float32(mod1(now * 0.56))
			const tailLen = float32(0.34)
			const tailSegs = 28

			px0, py0 := eval(headU - tailLen)
			for i := 1; i <= tailSegs; i++ {
				frac := float32(i) / float32(tailSegs)
				u := headU - tailLen*(1-frac)
				px1, py1 := eval(u)
				segCol := lerpColor(tailCol, accent, frac)
				alpha := frac * frac
				width := 1.4 + 2.8*frac
				p.Line(px0, py0, px1, py1, width, segCol.Alpha(alpha))
				px0, py0 = px1, py1
			}

			// Comet head glow halo + bright core.
			hx, hy := eval(headU)
			glowR := float32(13)
			p.FillRadialGradient(ui.Rect{X: hx - glowR, Y: hy - glowR, W: glowR * 2, H: glowR * 2}, ui.RadialGradient{
				Center: ui.Point{X: 0.5, Y: 0.5},
				Radius: glowR,
				Stops: []ui.GradientStop{
					{Color: accent.Alpha(0.55), At: 0},
					{Color: tailCol.Alpha(0.18), At: 0.55},
					{Color: ui.Transparent, At: 1},
				},
			}, glowR)
			p.Fill(ui.Rect{X: hx - 3.2, Y: hy - 3.2, W: 6.4, H: 6.4}, ui.Hex("#ffffff"), 3.2)

			// Center telemetry percentage.
			pct := fmt.Sprintf("%02d%%", int(headU*100))
			p.Text(cx-12, cy-7, pct, 11.5, t.TextMuted)
		})

		shapeName := "Rounded square"
		if st.hexagon {
			shapeName = "Hexagon"
		}
		ui.Text(c, shapeName+" · click to toggle").
			FontSize(11).
			FontWeight(500).
			TextColor(t.TextSubtle)
	})

	return wrap
}

// -----------------------------------------------------------------------------
// 11. MorphLoader ("morph-loader")
// -----------------------------------------------------------------------------

// MorphLoaderOptions configures MorphLoader.
type MorphLoaderOptions struct {
	Accent ui.Color
}

// MorphLoader renders a geometric shape morphing continuously between circle,
// rounded square, diamond, and pill with a rotating conic gradient fill.
func MorphLoader(c *ui.Context, opts ...MorphLoaderOptions) ui.Element {
	var o MorphLoaderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	colA := colorOr(o.Accent, t.Accent)
	colB := ui.Hex("#38bdf8")
	colC := ui.Hex("#a855f7")

	type morphPose struct {
		name string
		w    float32
		h    float32
		rad  float32
		rot  float32
	}
	poses := [4]morphPose{
		{name: "Circle", w: 68, h: 68, rad: 34, rot: 0},
		{name: "Square", w: 66, h: 66, rad: 14, rot: 0},
		{name: "Diamond", w: 58, h: 58, rad: 12, rot: 45},
		{name: "Pill", w: 94, h: 42, rad: 21, rot: 0},
	}

	now := phaseSeconds(c.Now())
	c.AnimationFrame()

	stepDur := 1.25
	stepF := now / stepDur
	idx := int(math.Floor(stepF)) % len(poses)
	if idx < 0 {
		idx += len(poses)
	}
	nextIdx := (idx + 1) % len(poses)
	local := float32(mod1(stepF))
	// Hold briefly at each canonical shape, then morph smoothly.
	blend := loadEaseInOut(clamp32((local-0.28)/0.72, 0, 1))

	a := poses[idx]
	b := poses[nextIdx]
	w := lerp32(a.w, b.w, blend)
	h := lerp32(a.h, b.h, blend)
	rad := lerp32(a.rad, b.rad, blend)
	rot := lerp32(a.rot, b.rot, blend)
	label := a.name
	if blend > 0.55 {
		label = b.name
	}

	wrap := ui.Column(c).AlignItems(ui.Center).Gap(8)
	wrap.Children(func() {
		stage := ui.Box(c).Size(156, 114)
		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.AnimationFrame()

			cx := r.X + r.W*0.5
			cy := r.Y + r.H*0.5

			// Ambient glow halo behind the morphing body.
			haloR := float32(52)
			haloAlpha := float32(0.22)
			if !t.Dark {
				haloAlpha = 0.14
			}
			p.FillRadialGradient(ui.Rect{X: cx - haloR, Y: cy - haloR, W: haloR * 2, H: haloR * 2}, ui.RadialGradient{
				Center: ui.Point{X: 0.5, Y: 0.5},
				Radius: haloR,
				Stops: []ui.GradientStop{
					{Color: colA.Alpha(haloAlpha), At: 0},
					{Color: colC.Alpha(haloAlpha * 0.45), At: 0.55},
					{Color: ui.Transparent, At: 1},
				},
			}, haloR)

			shapeRect := ui.Rect{X: cx - w*0.5, Y: cy - h*0.5, W: w, H: h}
			gradAngle := float32(math.Mod(now*110, 360))

			p.Transform(ui.Rotation(rot), cx, cy, func() {
				p.Shadow(shapeRect, rad, 0, 8, 18, 0, colA.Alpha(0.26))
				p.FillConicGradient(shapeRect, ui.ConicGradient{
					Center: ui.Point{X: 0.5, Y: 0.5},
					Angle:  gradAngle,
					Oklab:  true,
					Stops: []ui.GradientStop{
						{Color: colA, At: 0},
						{Color: colB, At: 0.34},
						{Color: colC, At: 0.68},
						{Color: colA, At: 1},
					},
				}, rad)

				// Top-left gloss highlight and crisp specular rim.
				p.FillRadialGradient(shapeRect, ui.RadialGradient{
					Center: ui.Point{X: 0.32, Y: 0.28},
					Radius: max(w, h) * 0.72,
					Stops: []ui.GradientStop{
						{Color: ui.Hex("#ffffff").Alpha(0.34), At: 0},
						{Color: ui.Transparent, At: 1},
					},
				}, rad)
				p.Stroke(shapeRect, ui.Hex("#ffffff").Alpha(0.38), rad, 1.2)
			})
		})

		ui.Text(c, label).
			FontSize(11.5).
			FontWeight(500).
			TextColor(t.TextMuted)
	})

	return wrap
}

// -----------------------------------------------------------------------------
// 12. PixelLoader ("pixel-loader")
// -----------------------------------------------------------------------------

// PixelLoaderOptions configures PixelLoader.
type PixelLoaderOptions struct {
	Accent ui.Color
}

type pixelLoaderState struct {
	manualMode int
}

// Precomputed spiral path over a 6x6 grid (36 cells) for the retro snake pattern.
var pixelSnakePath = [36][2]int{
	{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0},
	{5, 1}, {5, 2}, {5, 3}, {5, 4}, {5, 5},
	{4, 5}, {3, 5}, {2, 5}, {1, 5}, {0, 5},
	{0, 4}, {0, 3}, {0, 2}, {0, 1},
	{1, 1}, {2, 1}, {3, 1}, {4, 1},
	{4, 2}, {4, 3}, {4, 4},
	{3, 4}, {2, 4}, {1, 4},
	{1, 3}, {1, 2},
	{2, 2}, {3, 2},
	{3, 3}, {2, 3},
}

// Four 6x6 retro bit-sprite frames (each row is a 6-bit mask, MSB on left).
var pixelSpriteFrames = [4][6]uint8{
	{0b001100, 0b011110, 0b110011, 0b110011, 0b011110, 0b001100}, // Diamond ring
	{0b010010, 0b111111, 0b101101, 0b111111, 0b010010, 0b100001}, // Space sprite
	{0b001100, 0b001100, 0b111111, 0b111111, 0b001100, 0b001100}, // Plus pulse
	{0b100001, 0b010010, 0b001100, 0b001100, 0b010010, 0b100001}, // X burst
}

// PixelLoader renders a retro 6x6 pixel grid cycling animated snake, radar, and
// sprite bit patterns.
func PixelLoader(c *ui.Context, opts ...PixelLoaderOptions) ui.Element {
	var o PixelLoaderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	accent := colorOr(o.Accent, ui.Hex("#22c55e"))

	wrap := ui.Column(c).
		AlignItems(ui.Center).
		Gap(8).
		Cursor(ui.CursorPointer)

	st := ui.Local(wrap, "pixel_loader_mode", func() pixelLoaderState {
		return pixelLoaderState{}
	})
	if wrap.Clicked() {
		st.manualMode = (st.manualMode + 1) % 3
	}

	now := phaseSeconds(c.Now())
	c.AnimationFrame()

	mode := (st.manualMode + int(math.Floor(now/4.2))) % 3
	modeLabels := [3]string{"SNAKE", "RADAR", "SPRITE"}

	var bezelBg, offCell ui.Color
	if t.Dark {
		bezelBg = ui.Hex("#0b0f0e")
		offCell = ui.Hex("#1a2421")
	} else {
		bezelBg = ui.Hex("#141a18")
		offCell = ui.Hex("#23302b")
	}

	wrap.Children(func() {
		bezel := ui.Box(c).
			Size(112, 112).
			Radius(14).
			Background(bezelBg).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.BorderStrong)

		bezel.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.AnimationFrame()

			const n = 6
			cell := float32(11.5)
			gap := float32(3.5)
			gridSize := float32(n)*cell + float32(n-1)*gap
			ox := r.X + (r.W-gridSize)*0.5
			oy := r.Y + (r.H-gridSize)*0.5

			var intensity [6][6]float32

			switch mode {
			case 0: // Spiral snake with smooth glowing tail.
				headF := mod1(now*0.36) * float64(len(pixelSnakePath))
				head := int(math.Floor(headF))
				const snakeLen = 11
				for k := 0; k < snakeLen; k++ {
					idx := (head - k + len(pixelSnakePath)) % len(pixelSnakePath)
					pt := pixelSnakePath[idx]
					val := 1 - float32(k)/float32(snakeLen)
					intensity[pt[1]][pt[0]] = max(intensity[pt[1]][pt[0]], val*val)
				}

			case 1: // Sweeping radar with phosphor decay + target blip.
				beam := mod1(now*0.62) * 2 * math.Pi
				for y := 0; y < n; y++ {
					for x := 0; x < n; x++ {
						dx := float64(x) - 2.5
						dy := float64(y) - 2.5
						ang := math.Atan2(dy, dx)
						if ang < 0 {
							ang += 2 * math.Pi
						}
						diff := beam - ang
						if diff < 0 {
							diff += 2 * math.Pi
						}
						decay := clamp32(float32(1-diff/(math.Pi*1.15)), 0, 1)
						intensity[y][x] = decay * decay
					}
				}
				// Blinking radar target blip.
				blip := float32(0.5 + 0.5*math.Sin(now*7))
				intensity[1][4] = max(intensity[1][4], 0.45+0.55*blip)

			case 2: // Cross-fading retro bit-sprite frames.
				frameF := now * 1.65
				fA := int(math.Floor(frameF)) % len(pixelSpriteFrames)
				fB := (fA + 1) % len(pixelSpriteFrames)
				blend := loadSmoothStep(float32(mod1(frameF)))
				for y := 0; y < n; y++ {
					rowA := pixelSpriteFrames[fA][y]
					rowB := pixelSpriteFrames[fB][y]
					for x := 0; x < n; x++ {
						bit := uint8(1 << (5 - x))
						vA := float32(0)
						if rowA&bit != 0 {
							vA = 1
						}
						vB := float32(0)
						if rowB&bit != 0 {
							vB = 1
						}
						intensity[y][x] = lerp32(vA, vB, blend)
					}
				}
			}

			for y := 0; y < n; y++ {
				for x := 0; x < n; x++ {
					lum := clamp32(intensity[y][x], 0, 1)
					cx := ox + float32(x)*(cell+gap)
					cy := oy + float32(y)*(cell+gap)
					rect := ui.Rect{X: cx, Y: cy, W: cell, H: cell}

					if lum > 0.25 {
						glow := ui.Rect{X: cx - 3, Y: cy - 3, W: cell + 6, H: cell + 6}
						p.Fill(glow, accent.Alpha(0.16*lum), 4)
					}
					col := lerpColor(offCell, accent, lum)
					if lum > 0.85 {
						col = lerpColor(col, ui.Hex("#ecfdf5"), (lum-0.85)*2.5)
					}
					p.Fill(rect, col, 2.5)
				}
			}
		})

		ui.Text(c, modeLabels[mode]+" · 6x6 · click to cycle").
			FontSize(10.5).
			FontWeight(500).
			TextColor(t.TextSubtle)
	})

	return wrap
}

// -----------------------------------------------------------------------------
// 13. BlocksLoader ("blocks-loader")
// -----------------------------------------------------------------------------

// BlocksLoaderOptions configures BlocksLoader.
type BlocksLoaderOptions struct {
	Accent ui.Color
}

// BlocksLoader renders 2x2 sliding puzzle tiles smoothly orbiting an empty slot
// in a continuous choreographed loop.
func BlocksLoader(c *ui.Context, opts ...BlocksLoaderOptions) ui.Element {
	var o BlocksLoaderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	tileColors := [3]ui.Color{
		colorOr(o.Accent, t.Accent),
		ui.Hex("#38bdf8"),
		ui.Hex("#a855f7"),
	}

	stage := ui.Box(c).
		Size(116, 116).
		Radius(18).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	now := phaseSeconds(c.Now())
	c.AnimationFrame()

	// 4 grid slots in clockwise order: 0=TL, 1=TR, 2=BR, 3=BL.
	slotPos := [4][2]float32{
		{0, 0},
		{1, 0},
		{1, 1},
		{0, 1},
	}

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()

		tileSize := float32(38)
		gap := float32(8)
		gridSpan := tileSize*2 + gap
		ox := r.X + (r.W-gridSpan)*0.5
		oy := r.Y + (r.H-gridSpan)*0.5

		// Step advances every 0.48s: on each step the empty slot orbits clockwise
		// by sliding the tile from the next clockwise slot back into the empty slot.
		stepDur := 0.48
		totalSteps := now / stepDur
		stepIdx := int(math.Floor(totalSteps))
		frac := float32(mod1(totalSteps))
		ease := loadEaseInOut(clamp32(frac/0.78, 0, 1))

		// Simulate exact tile slots over stepIdx mod 12 (full permutation cycle).
		cycleSteps := stepIdx % 12
		if cycleSteps < 0 {
			cycleSteps += 12
		}
		// Initial state: empty at slot 3 (BL); tiles 0,1,2 at slots 0,1,2.
		slots := [4]int{0, 1, 2, -1}
		empty := 3
		for s := 0; s < cycleSteps; s++ {
			moverSlot := (empty + 3) % 4 // counter-clockwise neighbor slides into empty
			slots[empty] = slots[moverSlot]
			slots[moverSlot] = -1
			empty = moverSlot
		}

		// Active transition: tile at fromSlot slides into empty.
		fromSlot := (empty + 3) % 4
		toSlot := empty
		movingTile := slots[fromSlot]

		// Recessed dashed empty-slot target outline.
		tx := ox + slotPos[toSlot][0]*(tileSize+gap)
		ty := oy + slotPos[toSlot][1]*(tileSize+gap)
		p.StrokeDashed(ui.Rect{X: tx + 2, Y: ty + 2, W: tileSize - 4, H: tileSize - 4}, t.BorderStrong.Alpha(0.55), 9, 1.2)

		drawTile := func(tileID int, gx, gy, lift float32) {
			x := ox + gx*(tileSize+gap)
			y := oy + gy*(tileSize+gap)
			grow := lift * 1.8
			rect := ui.Rect{
				X: x - grow,
				Y: y - grow,
				W: tileSize + grow*2,
				H: tileSize + grow*2,
			}
			col := tileColors[tileID%len(tileColors)]
			rad := float32(10)

			p.Shadow(rect, rad, 0, 3+lift*3, 8+lift*6, 0, col.Alpha(0.24))
			p.FillGradient(rect, ui.LinearGradient{
				Angle: 145,
				Stops: []ui.GradientStop{
					{Color: lerpColor(col, ui.Hex("#ffffff"), 0.16), At: 0},
					{Color: col, At: 0.65},
					{Color: lerpColor(col, ui.Hex("#000000"), 0.14), At: 1},
				},
			}, rad)

			// Crisp inner specular rim + center micro dot.
			p.Stroke(rect, ui.Hex("#ffffff").Alpha(0.32), rad, 1)
			dotR := float32(2.6)
			p.Fill(ui.Rect{
				X: rect.X + rect.W*0.5 - dotR,
				Y: rect.Y + rect.H*0.5 - dotR,
				W: dotR * 2,
				H: dotR * 2,
			}, ui.Hex("#ffffff").Alpha(0.62), dotR)
		}

		// Draw stationary tiles first, then the sliding tile on top.
		for slot := 0; slot < 4; slot++ {
			tid := slots[slot]
			if tid < 0 || tid == movingTile {
				continue
			}
			drawTile(tid, slotPos[slot][0], slotPos[slot][1], 0)
		}
		if movingTile >= 0 {
			gx := lerp32(slotPos[fromSlot][0], slotPos[toSlot][0], ease)
			gy := lerp32(slotPos[fromSlot][1], slotPos[toSlot][1], ease)
			lift := float32(math.Sin(float64(ease) * math.Pi))
			drawTile(movingTile, gx, gy, lift)
		}
	})

	return stage
}

func init() {
	registerDemo(
		"thinking-loader",
		"Pulsing AI reasoning pill with a rotating orbital dot cluster and cycling status phrases.",
		func(c *ui.Context) {
			ThinkingLoader(c)
		},
	)
	registerDemo(
		"loading-state",
		"Multi-phase progress card cycling smoothly through Idle, Fetching, Processing, and Ready with click-to-replay.",
		func(c *ui.Context) {
			LoadingState(c)
		},
	)
	registerDemo(
		"blueprint-loader",
		"Technical blueprint wireframe grid with a sweeping luminous scanline, crosshairs, and progressive stroke reveal.",
		func(c *ui.Context) {
			BlueprintLoader(c)
		},
	)
	registerDemo(
		"dot-loader",
		"Precision 5x3 dot matrix animating choreographed wave, spiral, and ripple luminance patterns.",
		func(c *ui.Context) {
			DotLoader(c)
		},
	)
	registerDemo(
		"image-shimmer",
		"Rounded media skeleton frame with a crisp diagonal luminous shimmer sweep and subtle aspect badge.",
		func(c *ui.Context) {
			ImageShimmer(c)
		},
	)
	registerDemo(
		"shimmer-text",
		"Typography line with a continuous metallic highlight band gliding across individual characters.",
		func(c *ui.Context) {
			ShimmerText(c)
		},
	)
	registerDemo(
		"wave-loader",
		"Smooth vertical equalizer bars undulating in a continuous sine phase wave.",
		func(c *ui.Context) {
			WaveLoader(c)
		},
	)
	registerDemo(
		"bounce-loader",
		"Staggered bouncing spheres with squash-and-stretch deformation and altitude-reactive ground shadows.",
		func(c *ui.Context) {
			BounceLoader(c)
		},
	)
	registerDemo(
		"trace-loader",
		"Geometric rounded-square and hexagon perimeter with a glowing comet tail tracing seamlessly around its path.",
		func(c *ui.Context) {
			TraceLoader(c)
		},
	)
	registerDemo(
		"morph-loader",
		"Geometric shape morphing continuously between circle, square, diamond, and pill with a rotating conic gradient fill.",
		func(c *ui.Context) {
			MorphLoader(c)
		},
	)
	registerDemo(
		"pixel-loader",
		"Retro 6x6 pixel grid cycling animated snake, radar, and bit-sprite patterns.",
		func(c *ui.Context) {
			PixelLoader(c)
		},
	)
	registerDemo(
		"blocks-loader",
		"2x2 sliding puzzle tiles smoothly orbiting an empty slot in a continuous choreographed loop.",
		func(c *ui.Context) {
			BlocksLoader(c)
		},
	)
}
