package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

func hcSmoothStep(t float32) float32 {
	t = clamp32(t, 0, 1)
	return t * t * (3 - 2*t)
}

func hcRand(seed *uint32) float32 {
	x := *seed
	if x == 0 {
		x = 0x7f4a7c15
	}
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	*seed = x
	return float32(x&0x00ffffff) / float32(0x01000000)
}

func hcStarPath(cx, cy, outerR, innerR float32) *ui.Path {
	var p ui.Path
	for i := 0; i < 10; i++ {
		ang := -math.Pi/2 + float64(i)*math.Pi/5
		rad := outerR
		if i%2 == 1 {
			rad = innerR
		}
		x := cx + float32(math.Cos(ang))*rad
		y := cy + float32(math.Sin(ang))*rad
		if i == 0 {
			p.MoveTo(x, y)
		} else {
			p.LineTo(x, y)
		}
	}
	p.Close()
	return &p
}

// -----------------------------------------------------------------------------
// 7. AlbumShelf ("album-shelf")
// -----------------------------------------------------------------------------

// AlbumItem describes one vinyl sleeve in AlbumShelf.
type AlbumItem struct {
	Title  string
	Artist string
	ColorA ui.Color
	ColorB ui.Color
}

// AlbumShelfOptions configures AlbumShelf.
type AlbumShelfOptions struct {
	Width  float32
	Height float32
	Albums []AlbumItem
}

// AlbumShelf renders a row of vinyl sleeves standing side-by-side where hovering
// lifts and widens the hovered sleeve while sliding its grooved disc edge out.
func AlbumShelf(c *ui.Context, opts ...AlbumShelfOptions) ui.Element {
	var o AlbumShelfOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 298
	}
	if o.Height <= 0 {
		o.Height = 180
	}
	albums := o.Albums
	if len(albums) == 0 {
		albums = []AlbumItem{
			{Title: "Midnight Echo", Artist: "Kira Sato", ColorA: ui.Hex("#6366f1"), ColorB: ui.Hex("#312e81")},
			{Title: "Amber Waves", Artist: "Solstice", ColorA: ui.Hex("#f97316"), ColorB: ui.Hex("#9a3412")},
			{Title: "Glass Garden", Artist: "Nora Lin", ColorA: ui.Hex("#14b8a6"), ColorB: ui.Hex("#115e59")},
			{Title: "Velvet Static", Artist: "R. Vance", ColorA: ui.Hex("#ec4899"), ColorB: ui.Hex("#831843")},
		}
	}

	wrap := ui.Column(c).
		Size(o.Width, o.Height).
		Padding(10, 12).
		Justify(ui.SpaceBetween).
		AlignItems(ui.Center).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	hoveredAlbum := -1

	wrap.Children(func() {
		shelfRow := ui.Row(c).
			WidthPercent(100).
			Height(124).
			Justify(ui.Center).
			AlignItems(ui.End).
			Gap(8)

		sp := ui.SpringOptions{Duration: 340 * time.Millisecond, Bounce: 0.24}

		shelfRow.Children(func() {
			for i, alb := range albums {
				idx := i
				item := alb

				slot := ui.Box(c).
					Key(fmt.Sprintf("hc_alb_%d", idx)).
					Height(114).
					Cursor(ui.CursorPointer)

				hov := slot.Hovered()
				if hov {
					hoveredAlbum = idx
				}
				target := float32(0)
				if hov {
					target = 1
				}
				h := slot.Spring("h", target, sp)
				slotW := lerp32(52, 92, h)
				slot.Width(slotW)

				slot.Draw(func(p *ui.Painter, r ui.Rect) {
					sleeveW := lerp32(48, 66, h)
					sleeveH := float32(84)
					lift := lerp32(0, -16, h)
					sx := r.X + 2
					sy := r.Y + r.H - sleeveH + lift

					// Vinyl record disc peeking out behind the right edge of the sleeve.
					if h > 0.02 {
						discR := sleeveH * 0.43
						discCX := sx + sleeveW - 10 + lerp32(0, 22, h)
						discCY := sy + sleeveH*0.5
						discRect := ui.Rect{
							X: discCX - discR,
							Y: discCY - discR,
							W: discR * 2,
							H: discR * 2,
						}
						p.Fill(discRect, ui.Hex("#141418"), discR)
						// Grooves.
						for g := float32(0.84); g >= 0.44; g -= 0.16 {
							gr := discR * g
							p.Stroke(
								ui.Rect{X: discCX - gr, Y: discCY - gr, W: gr * 2, H: gr * 2},
								ui.Hex("#ffffff").Alpha(0.11),
								gr,
								0.9,
							)
						}
						// Center colored label + spindle hole.
						lblR := discR * 0.32
						p.Fill(ui.Rect{X: discCX - lblR, Y: discCY - lblR, W: lblR * 2, H: lblR * 2}, item.ColorA, lblR)
						p.Fill(ui.Rect{X: discCX - 2.5, Y: discCY - 2.5, W: 5, H: 5}, ui.Hex("#141418"), 2.5)
					}

					// Sleeve shadow & body.
					sleeveRect := ui.Rect{X: sx, Y: sy, W: sleeveW, H: sleeveH}
					p.Shadow(sleeveRect, 7, 0, lerp32(3, 9, h), lerp32(8, 16, h), 0, t.Shadow)
					p.FillGradient(sleeveRect, ui.LinearGradient{
						From:  item.ColorA,
						To:    item.ColorB,
						Angle: 150,
					}, 7)

					// Minimal geometric cover art inside sleeve.
					artR := sleeveW * 0.28
					artCX := sx + sleeveW*0.5
					artCY := sy + sleeveH*0.42
					p.Fill(
						ui.Rect{X: artCX - artR, Y: artCY - artR, W: artR * 2, H: artR * 2},
						ui.Hex("#ffffff").Alpha(0.24),
						artR,
					)
					p.Line(sx+8, sy+sleeveH-14, sx+sleeveW-8, sy+sleeveH-14, 2, ui.Hex("#ffffff").Alpha(0.45))
					// Spine highlight.
					p.Line(sx+2, sy+6, sx+2, sy+sleeveH-6, 1.2, ui.Hex("#ffffff").Alpha(0.32))
					p.Stroke(sleeveRect, ui.Hex("#ffffff").Alpha(0.18), 7, 1)
				})
			}
		})

		// Wooden/studio shelf bar.
		ui.Box(c).
			WidthPercent(94).
			Height(3).
			Radius(2).
			Background(t.BorderStrong)

		caption := "Hover a sleeve to pull out the vinyl"
		captionCol := t.TextMuted
		if hoveredAlbum >= 0 {
			a := albums[hoveredAlbum]
			caption = fmt.Sprintf("%s — %s", a.Title, a.Artist)
			captionCol = t.Text
		}
		ui.Text(c, caption).
			FontSize(11).
			FontWeight(600).
			TextColor(captionCol).
			SingleLine()
	})

	return wrap
}

// -----------------------------------------------------------------------------
// 8. ProgressiveBlur ("progressive-blur")
// -----------------------------------------------------------------------------

// ProgressiveBlurOptions configures ProgressiveBlur.
type ProgressiveBlurOptions struct {
	Width    float32
	Height   float32
	Badge    string
	Title    string
	Subtitle string
}

// ProgressiveBlur renders a layered visual card with a stepped bottom glass
// blur and scrim overlay that clears smoothly on hover.
func ProgressiveBlur(c *ui.Context, opts ...ProgressiveBlurOptions) ui.Element {
	var o ProgressiveBlurOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 290
	}
	if o.Height <= 0 {
		o.Height = 176
	}
	if o.Badge == "" {
		o.Badge = "OPTICAL GLASS"
	}
	if o.Title == "" {
		o.Title = "Alpine Dusk"
	}
	if o.Subtitle == "" {
		o.Subtitle = "Hover to clear progressive blur"
	}

	card := ui.Box(c).
		Size(o.Width, o.Height).
		Radius(16).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.BorderStrong).
		Cursor(ui.CursorPointer).
		Clip()

	hovered := card.Hovered()
	_, _, over := card.PointerPosition()
	clearTarget := float32(0)
	if hovered || over {
		clearTarget = 1
	}
	sp := ui.SpringOptions{Duration: 360 * time.Millisecond, Bounce: 0.12}
	cleared := card.Spring("hc_pb_clear", clearTarget, sp)
	blurAmt := 1 - cleared

	drawScene := func(p *ui.Painter, r ui.Rect) {
		// Vivid twilight sky gradient.
		p.FillGradient(r, ui.LinearGradient{
			Angle: 165,
			Stops: []ui.GradientStop{
				{Color: ui.Hex("#1e1b4b"), At: 0},
				{Color: ui.Hex("#7c3aed"), At: 0.42},
				{Color: ui.Hex("#f97316"), At: 0.80},
				{Color: ui.Hex("#fde047"), At: 1},
			},
		}, 0)

		// Glowing low sun.
		sunR := float32(34)
		sunX := r.X + r.W*0.68
		sunY := r.Y + r.H*0.42
		p.FillRadialGradient(
			ui.Rect{X: sunX - sunR*1.8, Y: sunY - sunR*1.8, W: sunR * 3.6, H: sunR * 3.6},
			ui.RadialGradient{
				Center: ui.Point{X: 0.5, Y: 0.5},
				Radius: sunR * 1.8,
				Stops: []ui.GradientStop{
					{Color: ui.Hex("#fef08a").Alpha(0.92), At: 0},
					{Color: ui.Hex("#fb923c").Alpha(0.45), At: 0.5},
					{Color: ui.Transparent, At: 1},
				},
			},
			sunR*1.8,
		)

		// Far mountain ridge.
		var far ui.Path
		far.MoveTo(r.X, r.Y+r.H)
		far.LineTo(r.X, r.Y+r.H*0.64)
		far.LineTo(r.X+r.W*0.24, r.Y+r.H*0.44)
		far.LineTo(r.X+r.W*0.48, r.Y+r.H*0.58)
		far.LineTo(r.X+r.W*0.76, r.Y+r.H*0.38)
		far.LineTo(r.X+r.W, r.Y+r.H*0.56)
		far.LineTo(r.X+r.W, r.Y+r.H)
		far.Close()
		p.FillPath(&far, ui.Hex("#312e81").Alpha(0.86))

		// Near foreground ridge + subtle horizon highlights.
		var near ui.Path
		near.MoveTo(r.X, r.Y+r.H)
		near.LineTo(r.X, r.Y+r.H*0.78)
		near.LineTo(r.X+r.W*0.34, r.Y+r.H*0.60)
		near.LineTo(r.X+r.W*0.62, r.Y+r.H*0.74)
		near.LineTo(r.X+r.W*0.88, r.Y+r.H*0.62)
		near.LineTo(r.X+r.W, r.Y+r.H*0.70)
		near.LineTo(r.X+r.W, r.Y+r.H)
		near.Close()
		p.FillPath(&near, ui.Hex("#0f172a").Alpha(0.92))
	}

	sc := lerp32(1.0, 1.04, cleared)

	card.Children(func() {
		// Base crisp artwork layer.
		ui.Box(c).
			Absolute().
			Left(0).
			Top(0).
			Size(o.Width, o.Height).
			Transform(ui.Scaling(sc, sc)).
			Draw(drawScene)

		// Stepped progressive blur bands across the lower half.
		type blurBand struct {
			topFrac float32
			sigma   float32
			scrim   float32
		}
		bands := [4]blurBand{
			{topFrac: 0.44, sigma: 2.2, scrim: 0.08},
			{topFrac: 0.56, sigma: 4.8, scrim: 0.16},
			{topFrac: 0.68, sigma: 8.0, scrim: 0.28},
			{topFrac: 0.80, sigma: 11.5, scrim: 0.44},
		}
		for i, b := range bands {
			sigma := b.sigma * blurAmt
			bandTop := o.Height * b.topFrac
			bandH := o.Height - bandTop
			if blurAmt > 0.02 {
				ui.Box(c).
					Key(fmt.Sprintf("hc_pb_b_%d", i)).
					Absolute().
					Left(0).
					Top(bandTop).
					Size(o.Width, bandH).
					Clip().
					Blur(sigma).
					Opacity(blurAmt).
					Children(func() {
						ui.Box(c).
							Absolute().
							Left(0).
							Top(-bandTop).
							Size(o.Width, o.Height).
							Transform(ui.Scaling(sc, sc)).
							Draw(drawScene)
					})
			}
		}

		// Smooth glass gradient scrim that fades on hover.
		ui.Box(c).
			Absolute().
			Left(0).
			Top(o.Height*0.36).
			Size(o.Width, o.Height*0.64).
			Draw(func(p *ui.Painter, r ui.Rect) {
				scrimAlpha := lerp32(0.76, 0.24, cleared)
				p.FillGradient(r, ui.LinearGradient{
					Angle: 180,
					Stops: []ui.GradientStop{
						{Color: ui.Transparent, At: 0},
						{Color: ui.Hex("#09090b").Alpha(scrimAlpha * 0.45), At: 0.48},
						{Color: ui.Hex("#09090b").Alpha(scrimAlpha), At: 1},
					},
				}, 0)
			})

		// Top glass badge + bottom caption text.
		ui.Column(c).
			Absolute().
			Left(0).
			Top(0).
			Size(o.Width, o.Height).
			Padding(12, 14).
			Justify(ui.SpaceBetween).
			Children(func() {
				ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
					ui.Box(c).
						Padding(3, 8).
						Radius(999).
						Background(ui.Hex("#09090b").Alpha(0.38)).
						BorderWidth(1, 1, 1, 1).
						BorderColor(ui.Hex("#ffffff").Alpha(0.20)).
						Children(func() {
							ui.Text(c, o.Badge).
								FontSize(9).
								FontWeight(700).
								LetterSpacing(0.6).
								TextColor(ui.Hex("#ffffff")).
								SingleLine()
						})
					stateTxt := "Blurred"
					if cleared > 0.5 {
						stateTxt = "Clear"
					}
					ui.Text(c, stateTxt).
						FontSize(10).
						FontWeight(600).
						TextColor(ui.Hex("#ffffff").Alpha(0.85)).
						SingleLine()
				})

				ui.Column(c).
					Gap(2).
					Transform(ui.Translation(0, lerp32(0, 4, cleared))).
					Children(func() {
						ui.Text(c, o.Title).
							FontSize(15).
							FontWeight(700).
							TextColor(ui.Hex("#ffffff")).
							SingleLine()
						sub := o.Subtitle
						if cleared > 0.55 {
							sub = "Full clarity revealed"
						}
						ui.Text(c, sub).
							FontSize(11).
							TextColor(ui.Hex("#ffffff").Alpha(0.78)).
							SingleLine()
					})
			})
	})

	return card
}

// -----------------------------------------------------------------------------
// 9. RadialPopOut ("radial-pop-out")
// -----------------------------------------------------------------------------

// RadialPopOutAction describes one satellite shortcut around RadialPopOut.
type RadialPopOutAction struct {
	Symbol string
	Label  string
	Color  ui.Color
}

// RadialPopOutOptions configures RadialPopOut.
type RadialPopOutOptions struct {
	Actions  []RadialPopOutAction
	OnSelect func(label string)
}

type hcRadialState struct {
	picked string
}

// RadialPopOut renders a compact center hub that springs 5 satellite shortcut
// badges radially outward while hovered.
func RadialPopOut(c *ui.Context, opts ...RadialPopOutOptions) ui.Element {
	var o RadialPopOutOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	actions := o.Actions
	if len(actions) == 0 {
		actions = []RadialPopOutAction{
			{Symbol: "★", Label: "Star", Color: ui.Hex("#f59e0b")},
			{Symbol: "✎", Label: "Edit", Color: ui.Hex("#6366f1")},
			{Symbol: "↗", Label: "Send", Color: ui.Hex("#0ea5e9")},
			{Symbol: "⊕", Label: "Pin", Color: ui.Hex("#10b981")},
			{Symbol: "⧉", Label: "Copy", Color: ui.Hex("#ec4899")},
		}
	}

	stageW := float32(276)
	stageH := float32(182)
	cx := stageW * 0.5
	cy := stageH*0.5 - 6

	stage := ui.Box(c).
		Size(stageW, stageH).
		Cursor(ui.CursorPointer)

	st := ui.Local(stage, "hc_radial_pop", func() hcRadialState {
		return hcRadialState{}
	})

	hovered := stage.Hovered()
	_, _, over := stage.PointerPosition()
	openTarget := float32(0)
	if hovered || over {
		openTarget = 1
	}

	sp := ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.32}
	open := stage.Spring("hc_rp_open", openTarget, sp)
	orbitR := float32(58)

	// Draw subtle orbit guide ring & connector spokes behind badges.
	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if open <= 0.02 {
			return
		}
		ringR := orbitR * open
		p.Stroke(
			ui.Rect{X: r.X + cx - ringR, Y: r.Y + cy - ringR, W: ringR * 2, H: ringR * 2},
			t.BorderStrong.Alpha(0.45*open),
			ringR,
			1,
		)
		for i, act := range actions {
			ang := -math.Pi/2 + float64(i)*2*math.Pi/float64(len(actions))
			sx := r.X + cx + float32(math.Cos(ang))*ringR
			sy := r.Y + cy + float32(math.Sin(ang))*ringR
			p.StrokePolyline([][2]float32{{r.X + cx, r.Y + cy}, {sx, sy}}, 1.1, act.Color.Alpha(0.28*open))
		}
	})

	stage.Children(func() {
		// Satellite badges.
		satSize := float32(34)
		for i, act := range actions {
			idx := i
			action := act
			ang := -math.Pi/2 + float64(idx)*2*math.Pi/float64(len(actions))

			// Slight per-item duration stagger for lively bloom.
			itemOpen := stage.Spring(
				fmt.Sprintf("hc_rp_%d", idx),
				openTarget,
				ui.SpringOptions{
					Duration: time.Duration(300+idx*38) * time.Millisecond,
					Bounce:   0.34,
				},
			)

			dist := orbitR * itemOpen
			bx := cx + float32(math.Cos(ang))*dist - satSize*0.5
			by := cy + float32(math.Sin(ang))*dist - satSize*0.5
			sc := lerp32(0.35, 1.0, itemOpen)

			sat := ui.Box(c).
				Key(fmt.Sprintf("hc_rp_sat_%d", idx)).
				Absolute().
				Left(bx).
				Top(by).
				Size(satSize, satSize).
				Radius(999).
				Center().
				Background(t.Surface).
				BorderWidth(1, 1, 1, 1).
				BorderColor(action.Color.Alpha(0.55)).
				Shadow(0, 4, 10, 0, t.Shadow).
				Opacity(clamp32(itemOpen*1.3, 0, 1)).
				Transform(ui.Scaling(sc, sc)).
				Cursor(ui.CursorPointer)

			if sat.Hovered() {
				sat.Background(action.Color.Alpha(0.18))
			}

			sat.OnClick(func() {
				st.picked = action.Label
				if o.OnSelect != nil {
					o.OnSelect(action.Label)
				}
			})

			sat.Children(func() {
				ui.Text(c, action.Symbol).
					FontSize(13).
					FontWeight(700).
					TextColor(action.Color).
					SingleLine()
			})
		}

		// Center hub button.
		hubSize := float32(46)
		hub := ui.Box(c).
			Absolute().
			Left(cx-hubSize*0.5).
			Top(cy-hubSize*0.5).
			Size(hubSize, hubSize).
			Radius(999).
			Center().
			Background(lerpColor(t.Surface, t.Accent, open)).
			BorderWidth(1, 1, 1, 1).
			BorderColor(lerpColor(t.BorderStrong, t.Accent, open)).
			Shadow(0, 6, 16, 0, t.Shadow).
			ZIndex(5).
			Cursor(ui.CursorPointer)

		hub.Children(func() {
			ui.Text(c, "+").
				FontSize(20).
				FontWeight(600).
				TextColor(lerpColor(t.Text, t.AccentText, open)).
				Transform(ui.Rotation(45 * open)).
				SingleLine()
		})

		// Subtle bottom status caption.
		status := "Hover hub to fan actions"
		if st.picked != "" {
			status = "Picked · " + st.picked
		}
		ui.Box(c).
			Absolute().
			Left(0).
			Bottom(2).
			Width(stageW).
			Center().
			Children(func() {
				ui.Text(c, status).
					FontSize(10.5).
					FontWeight(500).
					TextColor(t.TextMuted).
					SingleLine()
			})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 10. SharePill ("share-pill")
// -----------------------------------------------------------------------------

// SharePillOptions configures SharePill.
type SharePillOptions struct {
	Label   string
	OnShare func(channel string)
}

type hcSharePillState struct {
	lastChannel string
	at          time.Time
}

// SharePill renders a compact "Share" pill that smoothly expands horizontally
// on hover to reveal quick destination buttons with click feedback.
func SharePill(c *ui.Context, opts ...SharePillOptions) ui.Element {
	var o SharePillOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Label == "" {
		o.Label = "Share"
	}

	stage := ui.Column(c).
		Size(296, 128).
		Justify(ui.Center).
		AlignItems(ui.Center).
		Gap(10)

	st := ui.Local(stage, "hc_share_pill", func() hcSharePillState {
		return hcSharePillState{}
	})

	now := c.Now()
	feedbackActive := st.lastChannel != "" && now.Sub(st.at) < 1800*time.Millisecond
	if feedbackActive {
		c.AnimationFrame()
	}

	destinations := [4]struct {
		label string
		color ui.Color
	}{
		{"Copy", t.Accent},
		{"X", ui.Hex("#38bdf8")},
		{"Mail", ui.Hex("#f59e0b")},
		{"Link", ui.Hex("#10b981")},
	}

	stage.Children(func() {
		pill := ui.Box(c).
			Height(42).
			Radius(999).
			Background(t.Surface).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.BorderStrong).
			Shadow(0, 6, 16, 0, t.Shadow).
			Cursor(ui.CursorPointer).
			Clip()

		hov := pill.Hovered()
		_, _, over := pill.PointerPosition()
		expTarget := float32(0)
		if hov || over {
			expTarget = 1
		}
		sp := ui.SpringOptions{Duration: 340 * time.Millisecond, Bounce: 0.22}
		exp := pill.Spring("hc_sp_exp", expTarget, sp)
		pillW := lerp32(98, 252, exp)
		pill.Width(pillW)
		if exp > 0.1 {
			pill.BorderColor(lerpColor(t.BorderStrong, t.Accent.Alpha(0.55), exp))
		}

		pill.Children(func() {
			// Compact resting "Share" label.
			if exp < 0.72 {
				ui.Row(c).
					Absolute().
					Left(0).
					Top(0).
					Size(pillW, 42).
					Center().
					Gap(6).
					Opacity(clamp32(1-exp*1.9, 0, 1)).
					Children(func() {
						ui.Text(c, "↗").
							FontSize(13).
							FontWeight(700).
							TextColor(t.Accent).
							SingleLine()
						ui.Text(c, o.Label).
							FontSize(12).
							FontWeight(600).
							TextColor(t.Text).
							SingleLine()
					})
			}

			// Expanded quick destination buttons.
			if exp > 0.20 {
				ui.Row(c).
					Absolute().
					Left(0).
					Top(0).
					Size(252, 42).
					Padding(0, 8).
					Justify(ui.SpaceBetween).
					AlignItems(ui.Center).
					Opacity(clamp32((exp-0.20)/0.80, 0, 1)).
					Children(func() {
						for i, d := range destinations {
							dest := d
							btn := ui.Box(c).
								Key(fmt.Sprintf("hc_sp_d_%d", i)).
								Height(28).
								Padding(0, 11).
								Radius(999).
								Center().
								Background(t.SurfaceSubtle).
								Cursor(ui.CursorPointer)
							if btn.Hovered() {
								btn.Background(dest.color.Alpha(0.18))
							}
							btn.OnClick(func() {
								st.lastChannel = dest.label
								st.at = c.Now()
								if o.OnShare != nil {
									o.OnShare(dest.label)
								}
							})
							btn.Children(func() {
								txtCol := t.Text
								if btn.Hovered() {
									txtCol = dest.color
								}
								ui.Text(c, dest.label).
									FontSize(11).
									FontWeight(600).
									TextColor(txtCol).
									SingleLine()
							})
						}
					})
			}
		})

		msg := "Hover pill to reveal channels"
		msgCol := t.TextMuted
		if feedbackActive {
			msg = fmt.Sprintf("✓ Shared via %s", st.lastChannel)
			msgCol = t.Success
		}
		ui.Text(c, msg).
			FontSize(11).
			FontWeight(500).
			TextColor(msgCol).
			SingleLine()
	})

	return stage
}

// -----------------------------------------------------------------------------
// 11. SweepButton ("sweep-button")
// -----------------------------------------------------------------------------

// SweepButtonOptions configures SweepButton.
type SweepButtonOptions struct {
	Label   string
	OnClick func()
}

type hcSweepState struct {
	clicks int
}

// SweepButton renders a high-contrast CTA button with a diagonal specular sheen
// and inverted fill sweep gliding smoothly across the surface on hover.
func SweepButton(c *ui.Context, opts ...SweepButtonOptions) ui.Element {
	var o SweepButtonOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Label == "" {
		o.Label = "Launch workspace"
	}

	stage := ui.Column(c).
		Size(286, 132).
		Justify(ui.Center).
		AlignItems(ui.Center).
		Gap(10)

	st := ui.Local(stage, "hc_sweep_btn", func() hcSweepState {
		return hcSweepState{}
	})

	stage.Children(func() {
		btnW := float32(212)
		btnH := float32(46)
		btn := ui.Box(c).
			Size(btnW, btnH).
			Radius(14).
			Center().
			BorderWidth(1.25, 1.25, 1.25, 1.25).
			BorderColor(t.Text).
			Cursor(ui.CursorPointer).
			Clip()

		hov := btn.Hovered()
		_, _, over := btn.PointerPosition()
		target := float32(0)
		if hov || over {
			target = 1
		}
		sp := ui.SpringOptions{Duration: 360 * time.Millisecond, Bounce: 0.14}
		sweep := btn.Spring("hc_sw_p", target, sp)
		sheen := btn.Spring("hc_sw_s", target, ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: 0.04})

		btn.OnClick(func() {
			st.clicks++
			if o.OnClick != nil {
				o.OnClick()
			}
		})

		baseBg := t.Surface
		invBg := t.Text
		baseFg := t.Text
		invFg := t.Background

		btn.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, baseBg, 14)

			// Inverted high-contrast fill sweeping smoothly from left to right.
			if sweep > 0.005 {
				fillW := r.W * hcSmoothStep(sweep)
				p.Clip(r, 14, func() {
					var poly ui.Path
					skew := float32(18) * (1 - sweep)
					poly.MoveTo(r.X, r.Y)
					poly.LineTo(r.X+fillW+skew, r.Y)
					poly.LineTo(r.X+fillW-skew, r.Y+r.H)
					poly.LineTo(r.X, r.Y+r.H)
					poly.Close()
					p.FillPath(&poly, invBg)

					// Diagonal specular sheen ribbon gliding across the surface.
					sheenX := r.X - 48 + (r.W+96)*sheen
					sheenW := float32(34)
					var ribbon ui.Path
					ribbon.MoveTo(sheenX, r.Y)
					ribbon.LineTo(sheenX+sheenW, r.Y)
					ribbon.LineTo(sheenX+sheenW-22, r.Y+r.H)
					ribbon.LineTo(sheenX-22, r.Y+r.H)
					ribbon.Close()
					p.FillPath(&ribbon, ui.Hex("#ffffff").Alpha(0.28*float32(math.Sin(float64(clamp32(sheen, 0, 1))*math.Pi))))
				})
			}
		})

		textCol := lerpColor(baseFg, invFg, hcSmoothStep(sweep))
		btn.Children(func() {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, o.Label).
					FontSize(13).
					FontWeight(700).
					TextColor(textCol).
					SingleLine()
				ui.Text(c, "→").
					FontSize(14).
					FontWeight(700).
					TextColor(textCol).
					Transform(ui.Translation(lerp32(0, 4, sweep), 0)).
					SingleLine()
			})
		})

		sub := "Hover for inverted sweep & sheen"
		if st.clicks > 0 {
			sub = fmt.Sprintf("Triggered %d×", st.clicks)
		}
		ui.Text(c, sub).
			FontSize(11).
			TextColor(t.TextMuted).
			SingleLine()
	})

	return stage
}

// -----------------------------------------------------------------------------
// 12. ImageAccordion ("image-accordion")
// -----------------------------------------------------------------------------

// AccordionPanel describes one visual strip in ImageAccordion.
type AccordionPanel struct {
	Title  string
	Tag    string
	ColorA ui.Color
	ColorB ui.Color
}

// ImageAccordionOptions configures ImageAccordion.
type ImageAccordionOptions struct {
	Width  float32
	Height float32
	Panels []AccordionPanel
}

// ImageAccordion renders a horizontal row of 4 colorful vertical artwork panels
// where hovering expands the hovered panel wide with spring physics.
func ImageAccordion(c *ui.Context, opts ...ImageAccordionOptions) ui.Element {
	var o ImageAccordionOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 296
	}
	if o.Height <= 0 {
		o.Height = 172
	}
	panels := o.Panels
	if len(panels) == 0 {
		panels = []AccordionPanel{
			{Title: "Solstice", Tag: "01", ColorA: ui.Hex("#f97316"), ColorB: ui.Hex("#7c2d12")},
			{Title: "Lagoon", Tag: "02", ColorA: ui.Hex("#06b6d4"), ColorB: ui.Hex("#164e63")},
			{Title: "Orchid", Tag: "03", ColorA: ui.Hex("#a855f7"), ColorB: ui.Hex("#4c1d95")},
			{Title: "Emerald", Tag: "04", ColorA: ui.Hex("#10b981"), ColorB: ui.Hex("#064e3b")},
		}
	}

	row := ui.Row(c).
		Size(o.Width, o.Height).
		Gap(6).
		AlignItems(ui.Center)

	sp := ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.20}
	n := len(panels)

	// Measure pointer X across the row to determine which panel is hovered cleanly
	// without width feedback oscillation.
	px, _, over := row.PointerPosition()
	activeIdx := 0
	if over && n > 0 {
		band := o.Width / float32(n)
		idx := int(px / band)
		if idx < 0 {
			idx = 0
		}
		if idx >= n {
			idx = n - 1
		}
		activeIdx = idx
	}

	availW := o.Width - float32(n-1)*6

	// Compute spring weights first so total width always sums cleanly to availW.
	weights := make([]float32, n)
	sumW := float32(0)
	for i := 0; i < n; i++ {
		target := float32(1.0)
		if over {
			if i == activeIdx {
				target = 3.1
			} else {
				target = 0.85
			}
		} else if i == 0 {
			target = 1.6
		}
		w := row.Spring(fmt.Sprintf("hc_ia_w_%d", i), target, sp)
		if w < 0.2 {
			w = 0.2
		}
		weights[i] = w
		sumW += w
	}

	row.Children(func() {
		for i, pnl := range panels {
			idx := i
			panel := pnl
			pw := availW * (weights[idx] / sumW)
			expanded := clamp32((weights[idx]-1.0)/1.8, 0, 1)

			tile := ui.Box(c).
				Key(fmt.Sprintf("hc_ia_p_%d", idx)).
				Size(pw, o.Height).
				Radius(14).
				BorderWidth(1, 1, 1, 1).
				BorderColor(lerpColor(t.Border, ui.Hex("#ffffff").Alpha(0.30), expanded)).
				Cursor(ui.CursorPointer).
				Clip()

			tile.Draw(func(p *ui.Painter, r ui.Rect) {
				p.FillGradient(r, ui.LinearGradient{
					From:  panel.ColorA,
					To:    panel.ColorB,
					Angle: 155,
				}, 14)

				// Decorative glowing orb motif inside each panel.
				orbR := lerp32(22, 42, expanded)
				ox := r.X + r.W*0.65
				oy := r.Y + r.H*0.32
				p.FillRadialGradient(
					ui.Rect{X: ox - orbR, Y: oy - orbR, W: orbR * 2, H: orbR * 2},
					ui.RadialGradient{
						Center: ui.Point{X: 0.5, Y: 0.5},
						Radius: orbR,
						Stops: []ui.GradientStop{
							{Color: ui.Hex("#ffffff").Alpha(0.36), At: 0},
							{Color: ui.Transparent, At: 1},
						},
					},
					orbR,
				)

				// Soft bottom vignette for text legibility.
				p.FillGradient(
					ui.Rect{X: r.X, Y: r.Y + r.H*0.48, W: r.W, H: r.H * 0.52},
					ui.LinearGradient{
						Angle: 180,
						Stops: []ui.GradientStop{
							{Color: ui.Transparent, At: 0},
							{Color: ui.Hex("#09090b").Alpha(0.68), At: 1},
						},
					},
					0,
				)
			})

			tile.Children(func() {
				ui.Column(c).
					Size(pw, o.Height).
					Padding(10, 10).
					Justify(ui.SpaceBetween).
					Children(func() {
						ui.Text(c, panel.Tag).
							FontSize(10).
							FontWeight(700).
							TextColor(ui.Hex("#ffffff").Alpha(0.85)).
							SingleLine()

						if expanded > 0.25 {
							ui.Column(c).
								Gap(1).
								Opacity(clamp32((expanded-0.25)/0.75, 0, 1)).
								Children(func() {
									ui.Text(c, panel.Title).
										FontSize(13).
										FontWeight(700).
										TextColor(ui.Hex("#ffffff")).
										SingleLine()
									ui.Text(c, "Collection").
										FontSize(9.5).
										TextColor(ui.Hex("#ffffff").Alpha(0.75)).
										SingleLine()
								})
						}
					})
			})
		}
	})

	return row
}

// -----------------------------------------------------------------------------
// 13. TiltCard ("tilt")
// -----------------------------------------------------------------------------

// TiltCardOptions configures TiltCard.
type TiltCardOptions struct {
	Width    float32
	Height   float32
	Badge    string
	Title    string
	Subtitle string
	Accent   ui.Color
}

// TiltCard renders a glossy 3D-feel tilt card that shifts its perspective tilt,
// specular sheen, parallax inner badge, and dynamic shadow with the pointer.
func TiltCard(c *ui.Context, opts ...TiltCardOptions) ui.Element {
	var o TiltCardOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 246
	}
	if o.Height <= 0 {
		o.Height = 148
	}
	if o.Badge == "" {
		o.Badge = "STUDIO PASS"
	}
	if o.Title == "" {
		o.Title = "Holographic Card"
	}
	if o.Subtitle == "" {
		o.Subtitle = "3D perspective · parallax depth"
	}
	accent := colorOr(o.Accent, ui.Hex("#6366f1"))

	stage := ui.Box(c).
		Size(296, 182).
		Center().
		Cursor(ui.CursorPointer)

	hovered := stage.Hovered()
	px, py, over := stage.PointerPosition()
	bounds := stage.LayoutBounds()
	var nx, ny float32
	hoverTarget := float32(0)
	if (hovered || over) && bounds.W > 0 && bounds.H > 0 {
		nx = clamp32((px/bounds.W)*2-1, -1, 1)
		ny = clamp32((py/bounds.H)*2-1, -1, 1)
		hoverTarget = 1
	}

	sp := ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.18}
	tx := stage.Spring("hc_tc_x", nx, sp)
	ty := stage.Spring("hc_tc_y", ny, sp)
	hv := stage.Spring("hc_tc_h", hoverTarget, sp)

	transform := ui.Perspective(520).
		Then(ui.RotationX(-ty * 11)).
		Then(ui.RotationY(tx * 13)).
		Then(ui.Scaling(1+0.03*hv, 1+0.03*hv))

	stage.Children(func() {
		card := ui.Box(c).
			Size(o.Width, o.Height).
			Radius(18).
			Shadow(-tx*12, 10-ty*8, 22+8*hv, 0, t.Shadow).
			Transform(transform).
			Clip()

		card.Draw(func(p *ui.Painter, r ui.Rect) {
			// Deep iridescent studio gradient.
			p.FillGradient(r, ui.LinearGradient{
				Angle: 135 + tx*16,
				Stops: []ui.GradientStop{
					{Color: ui.Hex("#18181b"), At: 0},
					{Color: lerpColor(ui.Hex("#1e1b4b"), accent, 0.45), At: 0.55},
					{Color: ui.Hex("#0f172a"), At: 1},
				},
			}, 18)

			// Subtle geometric grid dots.
			for gx := float32(18); gx < r.W-12; gx += 22 {
				for gy := float32(18); gy < r.H-12; gy += 22 {
					p.Fill(ui.Rect{X: r.X + gx, Y: r.Y + gy, W: 1.5, H: 1.5}, ui.Hex("#ffffff").Alpha(0.08), 1)
				}
			}

			// Pointer-tracked glossy highlight sheen.
			sheenX := clamp32(0.5+tx*0.42, 0.05, 0.95)
			sheenY := clamp32(0.5+ty*0.42, 0.05, 0.95)
			p.FillRadialGradient(r, ui.RadialGradient{
				Center: ui.Point{X: sheenX, Y: sheenY},
				Radius: r.W * 0.65,
				Stops: []ui.GradientStop{
					{Color: ui.Hex("#ffffff").Alpha(0.28*hv + 0.06), At: 0},
					{Color: ui.Hex("#38bdf8").Alpha(0.14 * hv), At: 0.45},
					{Color: ui.Transparent, At: 1},
				},
			}, 18)

			p.Stroke(r, ui.Hex("#ffffff").Alpha(0.20+0.20*hv), 18, 1.2)
		})

		card.Children(func() {
			ui.Column(c).
				Size(o.Width, o.Height).
				Padding(16, 18).
				Justify(ui.SpaceBetween).
				Children(func() {
					// Floating parallax top badge + metallic chip.
					ui.Row(c).
						WidthPercent(100).
						Justify(ui.SpaceBetween).
						AlignItems(ui.Center).
						Transform(ui.Translation(tx*7, ty*5)).
						Children(func() {
							ui.Box(c).
								Padding(3, 9).
								Radius(999).
								Background(ui.Hex("#ffffff").Alpha(0.14)).
								BorderWidth(1, 1, 1, 1).
								BorderColor(ui.Hex("#ffffff").Alpha(0.24)).
								Children(func() {
									ui.Text(c, o.Badge).
										FontSize(9).
										FontWeight(700).
										LetterSpacing(0.8).
										TextColor(ui.Hex("#ffffff")).
										SingleLine()
								})

							ui.Box(c).
								Size(26, 18).
								Radius(5).
								Background(ui.Hex("#f59e0b").Alpha(0.75)).
								BorderWidth(1, 1, 1, 1).
								BorderColor(ui.Hex("#fde68a").Alpha(0.75))
						})

					// Parallax bottom title block.
					ui.Column(c).
						Gap(2).
						Transform(ui.Translation(tx*4.5, ty*3.5)).
						Children(func() {
							ui.Text(c, o.Title).
								FontSize(15).
								FontWeight(700).
								TextColor(ui.Hex("#ffffff")).
								SingleLine()
							ui.Text(c, o.Subtitle).
								FontSize(10.5).
								TextColor(ui.Hex("#ffffff").Alpha(0.72)).
								SingleLine()
						})
				})
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 14. EscapeButton ("escape")
// -----------------------------------------------------------------------------

// EscapeButtonOptions configures EscapeButton.
type EscapeButtonOptions struct {
	Width  float32
	Height float32
	Label  string
}

type hcEscapeState struct {
	tx       float32
	ty       float32
	evades   int
	caught   int
	seed     uint32
	lastJump time.Time
	inited   bool
}

// EscapeButton renders a playful "Catch me" button inside a bounded stage that
// springs briskly away from the approaching cursor and counts evades.
func EscapeButton(c *ui.Context, opts ...EscapeButtonOptions) ui.Element {
	var o EscapeButtonOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 296
	}
	if o.Height <= 0 {
		o.Height = 176
	}
	if o.Label == "" {
		o.Label = "Catch me"
	}

	btnW := float32(86)
	btnH := float32(32)
	padX := float32(12)
	minY := float32(34)
	maxX := o.Width - btnW - padX
	maxY := o.Height - btnH - 12

	stage := ui.Box(c).
		Size(o.Width, o.Height).
		Radius(15).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border).
		Clip()

	st := ui.Local(stage, "hc_escape_btn", func() hcEscapeState {
		return hcEscapeState{
			tx:     (o.Width - btnW) * 0.5,
			ty:     minY + (maxY-minY)*0.5,
			seed:   0x51a9c3d1,
			inited: true,
		}
	})

	px, py, over := stage.PointerPosition()
	now := c.Now()

	// Evade whenever pointer enters proximity around the button's current target.
	if over && now.Sub(st.lastJump) > 150*time.Millisecond {
		bcx := st.tx + btnW*0.5
		bcy := st.ty + btnH*0.5
		dx := px - bcx
		dy := py - bcy
		dist := float32(math.Hypot(float64(dx), float64(dy)))
		if dist < 56 {
			// Pick a bounded candidate farthest from pointer.
			bestX, bestY := st.tx, st.ty
			bestD := float32(-1)
			for tries := 0; tries < 6; tries++ {
				cx := padX + hcRand(&st.seed)*(maxX-padX)
				cy := minY + hcRand(&st.seed)*(maxY-minY)
				d := float32(math.Hypot(float64((cx+btnW*0.5)-px), float64((cy+btnH*0.5)-py)))
				if d > bestD {
					bestD = d
					bestX, bestY = cx, cy
				}
			}
			st.tx = clamp32(bestX, padX, maxX)
			st.ty = clamp32(bestY, minY, maxY)
			st.evades++
			st.lastJump = now
		}
	}

	sp := ui.SpringOptions{Duration: 240 * time.Millisecond, Bounce: 0.34}
	curX := stage.Spring("hc_esc_x", st.tx, sp)
	curY := stage.Spring("hc_esc_y", st.ty, sp)

	// Subtle playful dot-grid backdrop.
	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		dotCol := t.BorderStrong.Alpha(0.45)
		for gx := float32(16); gx < r.W-8; gx += 20 {
			for gy := float32(36); gy < r.H-8; gy += 20 {
				p.Fill(ui.Rect{X: r.X + gx, Y: r.Y + gy, W: 2, H: 2}, dotCol, 1)
			}
		}
	})

	stage.Children(func() {
		// Top counter header.
		ui.Row(c).
			Absolute().
			Left(12).
			Top(8).
			Width(o.Width - 24).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Children(func() {
				ui.Text(c, "Try to click it").
					FontSize(10.5).
					FontWeight(500).
					TextColor(t.TextMuted).
					SingleLine()

				badgeText := fmt.Sprintf("Evades: %d", st.evades)
				badgeCol := t.Accent
				if st.caught > 0 {
					badgeText = fmt.Sprintf("Caught %d · Evades %d", st.caught, st.evades)
					badgeCol = t.Success
				}
				ui.Box(c).
					Padding(2, 8).
					Radius(999).
					Background(badgeCol.Alpha(0.14)).
					Children(func() {
						ui.Text(c, badgeText).
							FontSize(10).
							FontWeight(700).
							TextColor(badgeCol).
							SingleLine()
					})
			})

		// Elusive spring-positioned button.
		btn := ui.Box(c).
			Absolute().
			Left(curX).
			Top(curY).
			Size(btnW, btnH).
			Radius(999).
			Center().
			Background(t.Accent).
			Shadow(0, 5, 12, 0, t.Accent.Alpha(0.32)).
			Cursor(ui.CursorPointer)

		btn.OnClick(func() {
			st.caught++
		})

		btn.Children(func() {
			ui.Text(c, o.Label).
				FontSize(11.5).
				FontWeight(700).
				TextColor(ui.Hex("#ffffff")).
				SingleLine()
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 15. MagnifyingDock ("dock")
// -----------------------------------------------------------------------------

// DockItem describes one app icon in MagnifyingDock.
type DockItem struct {
	Label  string
	Glyph  string
	ColorA ui.Color
	ColorB ui.Color
}

// MagnifyingDockOptions configures MagnifyingDock.
type MagnifyingDockOptions struct {
	Items    []DockItem
	OnSelect func(label string)
}

type hcDockState struct {
	active int
}

// MagnifyingDock renders a macOS-style horizontal icon dock with Gaussian
// proximity spring magnification and lift as the pointer glides across items.
func MagnifyingDock(c *ui.Context, opts ...MagnifyingDockOptions) ui.Element {
	var o MagnifyingDockOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	items := o.Items
	if len(items) == 0 {
		items = []DockItem{
			{Label: "Finder", Glyph: "☺", ColorA: ui.Hex("#38bdf8"), ColorB: ui.Hex("#2563eb")},
			{Label: "Notes", Glyph: "≡", ColorA: ui.Hex("#facc15"), ColorB: ui.Hex("#d97706")},
			{Label: "Music", Glyph: "♪", ColorA: ui.Hex("#fb7185"), ColorB: ui.Hex("#e11d48")},
			{Label: "Mail", Glyph: "✉", ColorA: ui.Hex("#60a5fa"), ColorB: ui.Hex("#4f46e5")},
			{Label: "Photos", Glyph: "✿", ColorA: ui.Hex("#34d399"), ColorB: ui.Hex("#059669")},
			{Label: "Term", Glyph: ">_", ColorA: ui.Hex("#52525b"), ColorB: ui.Hex("#18181b")},
		}
	}

	stageW := float32(300)
	stageH := float32(156)
	stage := ui.Box(c).
		Size(stageW, stageH).
		Cursor(ui.CursorPointer)

	st := ui.Local(stage, "hc_mag_dock", func() hcDockState {
		return hcDockState{active: 0}
	})

	px, py, over := stage.PointerPosition()
	inDockBand := over && py >= 46 && py <= stageH

	n := len(items)
	baseSize := float32(34)
	maxBoost := float32(18)
	gap := float32(8)
	sigma := float32(46)

	// Reference centers for Gaussian proximity calculation.
	restTotalW := float32(n)*baseSize + float32(n-1)*gap
	restStartX := (stageW - restTotalW) * 0.5

	sp := ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.24}

	sizes := make([]float32, n)
	lifts := make([]float32, n)
	closestIdx := -1
	closestDist := float32(1e9)
	totalCurW := float32(0)

	for i := 0; i < n; i++ {
		refCenter := restStartX + float32(i)*(baseSize+gap) + baseSize*0.5
		prox := float32(0)
		if inDockBand {
			dx := px - refCenter
			prox = float32(math.Exp(-float64(dx*dx) / (2 * float64(sigma*sigma))))
			absD := float32(math.Abs(float64(dx)))
			if absD < closestDist && absD < 34 {
				closestDist = absD
				closestIdx = i
			}
		}
		sz := stage.Spring(fmt.Sprintf("hc_dk_s_%d", i), baseSize+maxBoost*prox, sp)
		lf := stage.Spring(fmt.Sprintf("hc_dk_l_%d", i), -12*prox, sp)
		sizes[i] = sz
		lifts[i] = lf
		totalCurW += sz
	}
	totalCurW += float32(n-1) * gap

	dockPadX := float32(12)
	dockW := clamp32(totalCurW+dockPadX*2, restTotalW+dockPadX*2, stageW-8)
	dockH := float32(52)
	dockLeft := (stageW - dockW) * 0.5
	dockTop := stageH - dockH - 18

	// Compute hovered tooltip position.
	tipVisTarget := float32(0)
	tipXTarget := stageW * 0.5
	tipLabel := items[st.active].Label
	if closestIdx >= 0 {
		tipVisTarget = 1
		tipLabel = items[closestIdx].Label
		cursorX := dockLeft + dockPadX
		for i := 0; i < closestIdx; i++ {
			cursorX += sizes[i] + gap
		}
		tipXTarget = cursorX + sizes[closestIdx]*0.5
	}
	tipVis := stage.Spring("hc_dk_tv", tipVisTarget, ui.SpringOptions{Duration: 200 * time.Millisecond, Bounce: 0.14})
	tipX := stage.Spring("hc_dk_tx", tipXTarget, ui.SpringOptions{Duration: 190 * time.Millisecond, Bounce: 0.16})

	stage.Children(func() {
		// Floating tooltip pill above magnified item.
		tipW := float32(68)
		ui.Box(c).
			Absolute().
			Left(clamp32(tipX-tipW*0.5, 8, stageW-tipW-8)).
			Top(lerp32(dockTop-28, dockTop-40, tipVis)).
			Size(tipW, 22).
			Radius(999).
			Center().
			Background(t.Surface).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.BorderStrong).
			Shadow(0, 4, 10, 0, t.Shadow).
			Opacity(clamp32(tipVis, 0, 1)).
			Children(func() {
				ui.Text(c, tipLabel).
					FontSize(10.5).
					FontWeight(600).
					TextColor(t.Text).
					SingleLine()
			})

		// Frosted glass dock tray.
		trayBg := t.Surface.Alpha(0.88)
		if !t.Dark {
			trayBg = ui.Hex("#ffffff").Alpha(0.92)
		}
		ui.Box(c).
			Absolute().
			Left(dockLeft).
			Top(dockTop).
			Size(dockW, dockH).
			Radius(18).
			Background(trayBg).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.BorderStrong).
			Shadow(0, 10, 24, 0, t.Shadow)

		// App icon tiles positioned smoothly along the dock baseline.
		curX := (stageW - totalCurW) * 0.5
		baselineY := dockTop + dockH - 9
		for i, it := range items {
			idx := i
			item := it
			sz := sizes[idx]
			lf := lifts[idx]
			ix := curX
			iy := baselineY - sz + lf
			curX += sz + gap

			tile := ui.Box(c).
				Key(fmt.Sprintf("hc_dk_i_%d", idx)).
				Absolute().
				Left(ix).
				Top(iy).
				Size(sz, sz).
				Radius(sz*0.27).
				Center().
				Shadow(0, 4, 10, 0, t.Shadow).
				Cursor(ui.CursorPointer)

			tile.OnClick(func() {
				st.active = idx
				if o.OnSelect != nil {
					o.OnSelect(item.Label)
				}
			})

			tile.Draw(func(p *ui.Painter, r ui.Rect) {
				rad := r.W * 0.27
				p.FillGradient(r, ui.LinearGradient{
					From:  item.ColorA,
					To:    item.ColorB,
					Angle: 155,
				}, rad)
				// Subtle top inner gloss.
				p.Stroke(r, ui.Hex("#ffffff").Alpha(0.28), rad, 1)
			})

			tile.Children(func() {
				ui.Text(c, item.Glyph).
					FontSize(sz * 0.42).
					FontWeight(700).
					TextColor(ui.Hex("#ffffff")).
					SingleLine()
			})

			// Active app dot beneath selected item.
			if idx == st.active {
				ui.Box(c).
					Absolute().
					Left(ix+sz*0.5-2).
					Top(dockTop+dockH-5.5).
					Size(4, 4).
					Radius(99).
					Background(t.Text.Alpha(0.75))
			}
		}

		ui.Box(c).
			Absolute().
			Left(0).
			Bottom(0).
			Width(stageW).
			Center().
			Children(func() {
				ui.Text(c, "Glide pointer across the dock").
					FontSize(10).
					TextColor(t.TextSubtle).
					SingleLine()
			})
	})

	return stage
}

func init() {
	registerDemo(
		"album-shelf",
		"Side-by-side vinyl sleeves where hovering widens and lifts the sleeve to reveal its grooved disc edge.",
		func(c *ui.Context) { AlbumShelf(c) },
	)
	registerDemo(
		"progressive-blur",
		"Layered visual card with a stepped bottom glass blur scrim that clears smoothly on hover.",
		func(c *ui.Context) { ProgressiveBlur(c) },
	)
	registerDemo(
		"radial-pop-out",
		"Compact center hub that springs 5 satellite shortcut badges radially outward while hovered.",
		func(c *ui.Context) { RadialPopOut(c) },
	)
	registerDemo(
		"share-pill",
		"Compact Share pill that expands horizontally on hover to reveal quick destination buttons.",
		func(c *ui.Context) { SharePill(c) },
	)
	registerDemo(
		"sweep-button",
		"High-contrast CTA button with a diagonal specular sheen and inverted fill sweep on hover.",
		func(c *ui.Context) { SweepButton(c) },
	)
	registerDemo(
		"image-accordion",
		"Horizontal row of 4 colorful panels where hovering expands the active panel while compressing siblings.",
		func(c *ui.Context) { ImageAccordion(c) },
	)
	registerDemo(
		"tilt",
		"Glossy 3D perspective tilt card with pointer-tracked specular sheen, parallax inner badge, and dynamic shadow.",
		func(c *ui.Context) { TiltCard(c) },
	)
	registerDemo(
		"escape",
		"Playful Catch me button inside a bounded stage that springs briskly away from the approaching cursor.",
		func(c *ui.Context) { EscapeButton(c) },
	)
	registerDemo(
		"dock",
		"macOS-style horizontal icon dock with Gaussian proximity spring magnification and upward lift.",
		func(c *ui.Context) { MagnifyingDock(c) },
	)
}
