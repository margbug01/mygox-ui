package bencho

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// ============================================================================
// Shared vector-icon helpers (prefixed with pf to avoid package collisions)
// ============================================================================

func pfDrawCheck(p *ui.Painter, cx, cy, size, trim float32, strokeWidth float32, col ui.Color) {
	if trim <= 0.01 {
		return
	}
	var path ui.Path
	path.Cap(ui.RoundCap).Join(ui.RoundJoin)
	path.MoveTo(cx-size*0.42, cy+size*0.02)
	path.LineTo(cx-size*0.10, cy+size*0.34)
	path.LineTo(cx+size*0.44, cy-size*0.32)
	if trim < 0.999 {
		path.Trim(0, clamp32(trim, 0, 1))
	}
	p.StrokePath(&path, strokeWidth, col)
}

func pfDrawCopyIcon(p *ui.Painter, cx, cy, size float32, col ui.Color) {
	s := size * 0.56
	back := ui.Rect{X: cx - s*0.58, Y: cy - s*0.68, W: s * 0.92, H: s * 1.06}
	front := ui.Rect{X: cx - s*0.26, Y: cy - s*0.34, W: s * 0.92, H: s * 1.06}
	p.Stroke(back, col.Alpha(0.65), 2, 1.3)
	p.Stroke(front, col, 2, 1.4)
}

func pfDrawHeart(p *ui.Painter, cx, cy, r float32, fillAmt float32, fillCol, strokeCol ui.Color) {
	var path ui.Path
	path.Cap(ui.RoundCap).Join(ui.RoundJoin)
	topY := cy - r*0.22
	botY := cy + r*0.88
	path.MoveTo(cx, botY)
	path.CubeTo(cx-r*1.18, cy+r*0.22, cx-r*1.08, cy-r*0.86, cx, topY)
	path.CubeTo(cx+r*1.08, cy-r*0.86, cx+r*1.18, cy+r*0.22, cx, botY)
	path.Close()

	if fillAmt > 0.02 {
		p.FillPath(&path, fillCol.Alpha(clamp32(fillAmt, 0, 1)))
	}
	strokeAlpha := clamp32(1-fillAmt*0.65, 0.35, 1)
	p.StrokePath(&path, 1.6, lerpColor(strokeCol, fillCol, fillAmt).Alpha(strokeAlpha))
}

func pfDrawSearchIcon(p *ui.Painter, cx, cy, r float32, col ui.Color) {
	var ring ui.Path
	ring.Circle(cx-r*0.16, cy-r*0.16, r*0.56)
	p.StrokePath(&ring, 1.55, col)
	p.Line(cx+r*0.24, cy+r*0.24, cx+r*0.72, cy+r*0.72, 1.65, col)
}

func pfDrawMoodIcon(p *ui.Painter, cx, cy, r float32, mood int, col ui.Color) {
	var head ui.Path
	head.Circle(cx, cy, r)
	p.StrokePath(&head, 1.4, col)

	var eyeL, eyeR ui.Path
	eyeL.Circle(cx-r*0.32, cy-r*0.18, r*0.11)
	eyeR.Circle(cx+r*0.32, cy-r*0.18, r*0.11)
	p.FillPath(&eyeL, col)
	p.FillPath(&eyeR, col)

	var mouth ui.Path
	mouth.Cap(ui.RoundCap)
	switch mood {
	case 0: // frown
		mouth.MoveTo(cx-r*0.36, cy+r*0.36)
		mouth.QuadTo(cx, cy+r*0.04, cx+r*0.36, cy+r*0.36)
	case 1: // neutral
		mouth.MoveTo(cx-r*0.34, cy+r*0.26)
		mouth.LineTo(cx+r*0.34, cy+r*0.26)
	case 2: // smile
		mouth.MoveTo(cx-r*0.36, cy+r*0.18)
		mouth.QuadTo(cx, cy+r*0.52, cx+r*0.36, cy+r*0.18)
	default: // big grin
		mouth.MoveTo(cx-r*0.42, cy+r*0.14)
		mouth.QuadTo(cx, cy+r*0.64, cx+r*0.42, cy+r*0.14)
	}
	p.StrokePath(&mouth, 1.35, col)
}

func pfDrawCloseIcon(p *ui.Painter, cx, cy, r float32, col ui.Color) {
	p.Line(cx-r, cy-r, cx+r, cy+r, 1.5, col)
	p.Line(cx+r, cy-r, cx-r, cy+r, 1.5, col)
}

// ============================================================================
// 1. copy-field — Copy field
// ============================================================================

// CopyFieldOptions configures CopyField.
type CopyFieldOptions struct {
	Label       string
	CopyLabel   string
	CopiedLabel string
	Width       float32
	ResetAfter  time.Duration
	OnCopy      func(string)
}

type pfCopyFieldState struct {
	copied   bool
	copiedAt time.Time
	count    int
}

// CopyField renders a monospace token or URL field with a tactile Copy button
// that writes to the clipboard and flips to an emerald Copied pill before resetting.
func CopyField(c *ui.Context, value string, opts ...CopyFieldOptions) ui.Element {
	var o CopyFieldOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 286
	}
	if o.CopyLabel == "" {
		o.CopyLabel = "Copy"
	}
	if o.CopiedLabel == "" {
		o.CopiedLabel = "Copied"
	}
	if o.ResetAfter <= 0 {
		o.ResetAfter = 1800 * time.Millisecond
	}

	root := ui.Column(c).
		Width(o.Width).
		Gap(6)

	st := ui.Local(root, "pf_copy_field_st", func() pfCopyFieldState {
		return pfCopyFieldState{}
	})

	now := c.Now()
	if st.copied {
		c.AnimationFrame()
		if now.Sub(st.copiedAt) >= o.ResetAfter {
			st.copied = false
		}
	}

	copiedTarget := float32(0)
	if st.copied {
		copiedTarget = 1
	}
	copiedAmt := root.Spring("pf_cf_copied", copiedTarget, ui.SpringOptions{
		Duration: 260 * time.Millisecond,
		Bounce:   0.26,
	})

	root.Children(func() {
		if o.Label != "" {
			ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
				ui.Text(c, o.Label).FontSize(11).FontWeight(600).TextColor(t.TextMuted)
				if st.count > 0 {
					status := "Ready"
					col := t.TextSubtle
					if st.copied {
						status = "In clipboard"
						col = t.Success
					}
					ui.Text(c, status).FontSize(10).FontWeight(500).TextColor(col)
				}
			})
		}

		fieldBorder := lerpColor(t.BorderStrong, t.Success.Alpha(0.65), copiedAmt)
		field := ui.Row(c).
			WidthPercent(100).
			Height(40).
			Padding(4, 4, 4, 10).
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween).
			Gap(8).
			Radius(11).
			Background(t.Surface).
			Border(1, fieldBorder)

		field.Children(func() {
			ui.Box(c).Grow(1).Shrink(1).MinWidth(0).Clip().Children(func() {
				ui.Text(c, value).
					Font("monospace").
					FontSize(11.5).
					FontWeight(500).
					TextColor(t.Text).
					NoWrap().
					TruncateMiddle()
			})

			btn := ui.Row(c).
				Key("copy_btn").
				Height(30).
				Padding(0, 10).
				Center().
				Gap(5).
				Radius(8).
				Cursor(ui.CursorPointer)

			pressTarget := float32(1)
			if btn.Pressed() {
				pressTarget = 0.94
			}
			btnScale := btn.Spring("pf_cf_scale", pressTarget, ui.SpringOptions{
				Duration: 160 * time.Millisecond,
				Bounce:   0.25,
			})

			idleBg := t.SurfaceSubtle
			if btn.Hovered() {
				idleBg = t.SurfaceHover
			}
			bg := lerpColor(idleBg, t.Success.Alpha(0.18), copiedAmt)
			bdr := lerpColor(t.Border, t.Success.Alpha(0.55), copiedAmt)
			ink := lerpColor(t.Text, t.Success, copiedAmt)

			btn.Background(bg).Border(1, bdr).Scale(btnScale, btnScale).Children(func() {
				ui.Box(c).Width(13).Height(13).Draw(func(p *ui.Painter, r ui.Rect) {
					cx, cy := r.X+r.W*0.5, r.Y+r.H*0.5
					if copiedAmt > 0.35 {
						pfDrawCheck(p, cx, cy, 11, (copiedAmt-0.35)/0.65, 1.7, t.Success)
					} else {
						pfDrawCopyIcon(p, cx, cy, 12, ink)
					}
				})
				lbl := o.CopyLabel
				if st.copied {
					lbl = o.CopiedLabel
				}
				ui.Text(c, lbl).FontSize(11.5).FontWeight(600).TextColor(ink)
			})

			if btn.Clicked() {
				c.WriteClipboard(value)
				st.copied = true
				st.copiedAt = now
				st.count++
				if o.OnCopy != nil {
					o.OnCopy(value)
				}
			}
		})
	})

	return root
}

// ============================================================================
// 3. feedback-bar — Feedback bar
// ============================================================================

// FeedbackBarOptions configures FeedbackBar.
type FeedbackBarOptions struct {
	Prompt   string
	Width    float32
	OnSubmit func(rating int, note string)
}

type pfFeedbackBarState struct {
	selected    int
	tagIdx      int
	submitted   bool
	submittedAt time.Time
}

var pfFeedbackMoods = []struct {
	label string
	tag   string
}{
	{label: "No", tag: "Confusing"},
	{label: "Meh", tag: "Needs detail"},
	{label: "Good", tag: "Clear"},
	{label: "Love", tag: "Super helpful"},
}

var pfFeedbackQuickTags = []string{"Clear", "Fast", "Needs docs"}

// FeedbackBar renders a compact "Was this helpful?" bar with mood rating pills
// that expands inline to collect a quick tag and confirm with a thank-you state.
func FeedbackBar(c *ui.Context, opts ...FeedbackBarOptions) ui.Element {
	var o FeedbackBarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}
	if o.Prompt == "" {
		o.Prompt = "Was this helpful?"
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(10, 12).
		Gap(8).
		Radius(14).
		Background(t.Surface).
		Border(1, t.Border)

	st := ui.Local(card, "pf_feedback_bar_st", func() pfFeedbackBarState {
		return pfFeedbackBarState{selected: -1, tagIdx: -1}
	})

	now := c.Now()
	if st.submitted && now.Sub(st.submittedAt) < 2600*time.Millisecond {
		c.AnimationFrame()
	}

	expandedTarget := float32(0)
	if st.selected >= 0 && !st.submitted {
		expandedTarget = 1
	}
	expanded := card.Spring("pf_fb_expand", expandedTarget, ui.SpringOptions{
		Duration: 280 * time.Millisecond,
		Bounce:   0.22,
	})

	card.Children(func() {
		if st.submitted {
			thanks := ui.Row(c).
				Key("fb_thanks").
				WidthPercent(100).
				Justify(ui.SpaceBetween).
				AlignItems(ui.Center).
				Padding(4, 2)
			thanks.Children(func() {
				ui.Row(c).AlignItems(ui.Center).Gap(7).Children(func() {
					ui.Box(c).Width(20).Height(20).Radius(999).Background(t.Success.Alpha(0.18)).Draw(func(p *ui.Painter, r ui.Rect) {
						pfDrawCheck(p, r.X+r.W*0.5, r.Y+r.H*0.5, 11, 1, 1.8, t.Success)
					})
					ui.Column(c).Gap(1).Children(func() {
						ui.Text(c, "Thanks for the feedback!").FontSize(12).FontWeight(600).TextColor(t.Text)
						picked := pfFeedbackMoods[st.selected].label
						if st.tagIdx >= 0 && st.tagIdx < len(pfFeedbackQuickTags) {
							picked += " · " + pfFeedbackQuickTags[st.tagIdx]
						}
						ui.Text(c, picked).FontSize(10.5).TextColor(t.Success)
					})
				})

				undo := ui.Row(c).
					Key("fb_reset").
					Height(24).
					Padding(0, 8).
					Center().
					Radius(999).
					Background(t.SurfaceSubtle).
					Cursor(ui.CursorPointer)
				undo.Children(func() {
					ui.Text(c, "Edit").FontSize(10.5).FontWeight(500).TextColor(t.TextMuted)
				})
				if undo.Clicked() {
					st.submitted = false
				}
			})
			return
		}

		// Top compact rating bar
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Gap(6).Children(func() {
			ui.Text(c, o.Prompt).FontSize(11.5).FontWeight(600).TextColor(t.Text)

			ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
				moodColors := []ui.Color{t.Danger, t.Warning, t.Accent, t.Success}
				for i, m := range pfFeedbackMoods {
					idx := i
					active := st.selected == idx
					col := moodColors[idx]

					pill := ui.Box(c).
						Key(idx).
						Width(28).
						Height(28).
						Center().
						Radius(999).
						Cursor(ui.CursorPointer)

					scTarget := float32(1)
					if active {
						scTarget = 1.08
					} else if pill.Pressed() {
						scTarget = 0.92
					}
					sc := pill.Spring("pf_fb_m", scTarget, ui.SpringOptions{
						Duration: 200 * time.Millisecond,
						Bounce:   0.35,
					})

					bg := t.SurfaceSubtle
					bdr := t.Border
					iconCol := t.TextMuted
					if active {
						bg = col.Alpha(0.18)
						bdr = col.Alpha(0.60)
						iconCol = col
					} else if pill.Hovered() {
						bg = t.SurfaceHover
						iconCol = t.Text
					}

					pill.Background(bg).Border(1, bdr).Scale(sc, sc).Draw(func(p *ui.Painter, r ui.Rect) {
						pfDrawMoodIcon(p, r.X+r.W*0.5, r.Y+r.H*0.5, 7.5, idx, iconCol)
					})
					if pill.Clicked() {
						if st.selected == idx {
							st.selected = -1
						} else {
							st.selected = idx
						}
					}
					_ = m
				}
			})
		})

		// Expandable quick-comment drawer
		if expanded > 0.04 {
			drawer := ui.Column(c).
				WidthPercent(100).
				Gap(8).
				Padding(8, 0, 0, 0).
				BorderWidth(1, 0, 0, 0).
				BorderColor(t.Border).
				Opacity(clamp32(expanded, 0, 1)).
				Translate(0, (1-expanded)*-6)

			drawer.Children(func() {
				ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Gap(6).Children(func() {
					ui.Row(c).Gap(4).AlignItems(ui.Center).Children(func() {
						for i, tag := range pfFeedbackQuickTags {
							tIdx := i
							picked := st.tagIdx == tIdx
							chip := ui.Row(c).
								Key(tIdx).
								Height(24).
								Padding(0, 8).
								Center().
								Radius(999).
								Cursor(ui.CursorPointer)
							if picked {
								chip.Background(t.AccentSoft).Border(1, t.Accent.Alpha(0.55))
							} else {
								chip.Background(t.SurfaceSubtle).Border(1, t.Border)
							}
							chip.Children(func() {
								txtCol := t.TextMuted
								if picked {
									txtCol = t.Accent
								}
								ui.Text(c, tag).FontSize(10.5).FontWeight(500).TextColor(txtCol)
							})
							if chip.Clicked() {
								if st.tagIdx == tIdx {
									st.tagIdx = -1
								} else {
									st.tagIdx = tIdx
								}
							}
						}
					})

					send := ui.Row(c).
						Key("fb_send").
						Height(25).
						Padding(0, 10).
						Center().
						Radius(999).
						Background(t.Text).
						Cursor(ui.CursorPointer)
					send.Children(func() {
						ui.Text(c, "Send").FontSize(10.5).FontWeight(600).TextColor(t.Background)
					})
					if send.Clicked() {
						st.submitted = true
						st.submittedAt = now
						if o.OnSubmit != nil {
							note := ""
							if st.tagIdx >= 0 && st.tagIdx < len(pfFeedbackQuickTags) {
								note = pfFeedbackQuickTags[st.tagIdx]
							}
							o.OnSubmit(st.selected, note)
						}
					}
				})
			})
		}
	})

	return card
}

// ============================================================================
// 4. price-flow — Price flow
// ============================================================================

// PriceTier defines one pricing plan in PriceFlow.
type PriceTier struct {
	Name    string
	Price   float32
	Period  string
	Tagline string
}

// PriceFlowOptions configures PriceFlow.
type PriceFlowOptions struct {
	Width    float32
	Currency string
	OnChange func(int)
}

func pfDefaultPriceTiers() []PriceTier {
	return []PriceTier{
		{Name: "Starter", Price: 19, Period: "/mo", Tagline: "For solo makers & side projects"},
		{Name: "Pro", Price: 49, Period: "/mo", Tagline: "Unlimited seats & priority GPU"},
		{Name: "Enterprise", Price: 149, Period: "/mo", Tagline: "Dedicated cluster & 99.99% SLA"},
	}
}

// PriceFlow renders an animated rolling-digit price tag cycling between
// Starter, Pro, and Enterprise tiers with smooth spring numeric interpolation.
func PriceFlow(c *ui.Context, tiers []PriceTier, active *int, opts ...PriceFlowOptions) ui.Element {
	var o PriceFlowOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if len(tiers) == 0 {
		tiers = pfDefaultPriceTiers()
	}
	if o.Width <= 0 {
		o.Width = 280
	}
	if o.Currency == "" {
		o.Currency = "$"
	}

	curIdx := 0
	if active != nil && *active >= 0 && *active < len(tiers) {
		curIdx = *active
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(14).
		Gap(12).
		Radius(16).
		Background(t.Surface).
		Border(1, t.Border)

	targetPrice := tiers[curIdx].Price
	animPrice := card.Spring("pf_price_val", targetPrice, ui.SpringOptions{
		Duration: 460 * time.Millisecond,
		Bounce:   0.18,
	})
	if animPrice < 0 {
		animPrice = 0
	}

	sliderPos := card.Spring("pf_price_tab", float32(curIdx), ui.SpringOptions{
		Duration: 300 * time.Millisecond,
		Bounce:   0.22,
	})

	card.Children(func() {
		// Segmented tier control with sliding spring pill
		segW := o.Width - 28
		tabW := (segW - 6) / float32(len(tiers))
		seg := ui.Row(c).
			Width(segW).
			Height(30).
			Padding(3).
			Radius(999).
			Background(t.SurfaceSubtle).
			Border(1, t.Border)

		seg.Draw(func(p *ui.Painter, r ui.Rect) {
			pillX := r.X + 3 + clamp32(sliderPos, 0, float32(len(tiers)-1))*tabW
			pillRect := ui.Rect{X: pillX, Y: r.Y + 3, W: tabW, H: r.H - 6}
			p.Fill(pillRect, t.Surface, pillRect.H*0.5)
			p.Stroke(pillRect, t.BorderStrong.Alpha(0.7), pillRect.H*0.5, 1)
		})

		seg.Children(func() {
			for i, tier := range tiers {
				idx := i
				isActive := idx == curIdx
				tab := ui.Row(c).
					Key(idx).
					Width(tabW).
					Height(24).
					Center().
					Cursor(ui.CursorPointer)
				tab.Children(func() {
					col := t.TextMuted
					w := 500
					if isActive {
						col = t.Text
						w = 600
					}
					ui.Text(c, tier.Name).FontSize(11).FontWeight(w).TextColor(col)
				})
				if tab.Clicked() && active != nil && *active != idx {
					*active = idx
					if o.OnChange != nil {
						o.OnChange(idx)
					}
				}
			}
		})

		// Rolling digit display
		val := float64(animPrice)
		whole := int(math.Floor(val))
		frac := float32(val - float64(whole))
		if whole < 0 {
			whole = 0
			frac = 0
		}

		// Determine digits needed (always show hundreds slot when >= 100 or animating through >= 99.5)
		showHundreds := whole >= 100 || targetPrice >= 100
		hundreds := (whole / 100) % 10
		tens := (whole / 10) % 10
		ones := whole % 10

		tensFrac := float32(0)
		if ones == 9 {
			tensFrac = frac
		}
		hundredsFrac := float32(0)
		if tens == 9 && ones == 9 {
			hundredsFrac = frac
		}

		digitH := float32(34)
		renderDigit := func(key string, digit int, rollFrac float32) {
			nextDigit := (digit + 1) % 10
			slot := ui.Box(c).
				Key(key).
				Width(18).
				Height(digitH).
				Clip()
			slot.Children(func() {
				ui.Box(c).
					Absolute().
					Left(0).
					Right(0).
					Top(0).
					Height(digitH).
					Center().
					Translate(0, -rollFrac*digitH).
					Opacity(1 - rollFrac*0.35).
					Children(func() {
						ui.Text(c, fmt.Sprintf("%d", digit)).
							FontSize(28).
							FontWeight(700).
							TextColor(t.Text)
					})
				if rollFrac > 0.005 {
					ui.Box(c).
						Absolute().
						Left(0).
						Right(0).
						Top(0).
						Height(digitH).
						Center().
						Translate(0, (1-rollFrac)*digitH).
						Opacity(0.65 + rollFrac*0.35).
						Children(func() {
							ui.Text(c, fmt.Sprintf("%d", nextDigit)).
								FontSize(28).
								FontWeight(700).
								TextColor(t.Text)
						})
				}
			})
		}

		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.End).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(1).Children(func() {
				ui.Text(c, o.Currency).FontSize(18).FontWeight(600).TextColor(t.TextMuted).Translate(0, -4)
				if showHundreds {
					hAlpha := clamp32((animPrice-88)/14, 0, 1)
					if targetPrice >= 100 && hAlpha < 0.2 {
						hAlpha = 0.2
					}
					ui.Box(c).Opacity(hAlpha).Children(func() {
						renderDigit("d_h", hundreds, hundredsFrac)
					})
				}
				renderDigit("d_t", tens, tensFrac)
				renderDigit("d_o", ones, frac)
				ui.Text(c, tiers[curIdx].Period).
					FontSize(12).
					FontWeight(500).
					TextColor(t.TextMuted).
					Translate(3, 5)
			})

			badge := ui.Row(c).
				Height(22).
				Padding(0, 8).
				Center().
				Radius(999).
				Background(t.AccentSoft)
			badge.Children(func() {
				ui.Text(c, tiers[curIdx].Name).FontSize(10).FontWeight(600).TextColor(t.Accent)
			})
		})

		ui.Text(c, tiers[curIdx].Tagline).FontSize(11).TextColor(t.TextSubtle)
	})

	return card
}

// ============================================================================
// 5. confirm — Inline confirm
// ============================================================================

// InlineConfirmOptions configures InlineConfirm.
type InlineConfirmOptions struct {
	Label      string
	Prompt     string
	YesLabel   string
	NoLabel    string
	DoneLabel  string
	ResetAfter time.Duration
	OnCancel   func()
}

type pfInlineConfirmState struct {
	phase   int // 0=idle, 1=confirming, 2=done
	phaseAt time.Time
}

// InlineConfirm renders a compact destructive inline confirmation button that
// expands smoothly into a "Sure? Yes / No" pill without abrupt layout jumps.
func InlineConfirm(c *ui.Context, onConfirm func(), opts ...InlineConfirmOptions) ui.Element {
	var o InlineConfirmOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Label == "" {
		o.Label = "Delete"
	}
	if o.Prompt == "" {
		o.Prompt = "Sure?"
	}
	if o.YesLabel == "" {
		o.YesLabel = "Yes"
	}
	if o.NoLabel == "" {
		o.NoLabel = "No"
	}
	if o.DoneLabel == "" {
		o.DoneLabel = "Deleted"
	}
	if o.ResetAfter <= 0 {
		o.ResetAfter = 1500 * time.Millisecond
	}

	holder := ui.Box(c).
		Width(236).
		Height(44).
		Center()

	st := ui.Local(holder, "pf_inline_confirm_st", func() pfInlineConfirmState {
		return pfInlineConfirmState{}
	})

	now := c.Now()
	if st.phase == 2 {
		c.AnimationFrame()
		if now.Sub(st.phaseAt) >= o.ResetAfter {
			st.phase = 0
		}
	}

	targetW := float32(106)
	switch st.phase {
	case 1:
		targetW = 212
	case 2:
		targetW = 118
	}
	pillW := holder.Spring("pf_ic_w", targetW, ui.SpringOptions{
		Duration: 280 * time.Millisecond,
		Bounce:   0.22,
	})

	confirmAmt := holder.Spring("pf_ic_conf", func() float32 {
		if st.phase == 1 {
			return 1
		}
		return 0
	}(), ui.SpringOptions{Duration: 240 * time.Millisecond, Bounce: 0.18})

	holder.Children(func() {
		bg := lerpColor(t.Danger.Alpha(0.12), t.Surface, confirmAmt)
		bdr := lerpColor(t.Danger.Alpha(0.45), t.BorderStrong, confirmAmt)
		if st.phase == 2 {
			bg = t.Success.Alpha(0.14)
			bdr = t.Success.Alpha(0.50)
		}

		pill := ui.Row(c).
			Width(pillW).
			Height(38).
			Padding(0, 6).
			Center().
			Gap(6).
			Radius(999).
			Background(bg).
			Border(1, bdr).
			Clip()

		pill.Children(func() {
			switch st.phase {
			case 0:
				delBtn := ui.Row(c).
					Key("ic_del").
					WidthPercent(100).
					HeightPercent(100).
					Center().
					Gap(6).
					Cursor(ui.CursorPointer)
				delBtn.Children(func() {
					ui.Box(c).Width(7).Height(7).Radius(999).Background(t.Danger)
					ui.Text(c, o.Label).FontSize(12.5).FontWeight(600).TextColor(t.Danger).NoWrap()
				})
				if delBtn.Clicked() {
					st.phase = 1
					st.phaseAt = now
				}

			case 1:
				ui.Text(c, o.Prompt).
					FontSize(12).
					FontWeight(600).
					TextColor(t.Text).
					NoWrap().
					Padding(0, 4, 0, 6)

				yesBtn := ui.Row(c).
					Key("ic_yes").
					Height(26).
					Padding(0, 12).
					Center().
					Radius(999).
					Background(t.Danger).
					Cursor(ui.CursorPointer)
				yesBtn.Children(func() {
					ui.Text(c, o.YesLabel).FontSize(11.5).FontWeight(600).TextColor(ui.Hex("#ffffff")).NoWrap()
				})
				if yesBtn.Clicked() {
					st.phase = 2
					st.phaseAt = now
					if onConfirm != nil {
						onConfirm()
					}
				}

				noBtn := ui.Row(c).
					Key("ic_no").
					Height(26).
					Padding(0, 11).
					Center().
					Radius(999).
					Background(t.SurfaceSubtle).
					Border(1, t.Border).
					Cursor(ui.CursorPointer)
				noBtn.Children(func() {
					ui.Text(c, o.NoLabel).FontSize(11.5).FontWeight(500).TextColor(t.TextMuted).NoWrap()
				})
				if noBtn.Clicked() {
					st.phase = 0
					st.phaseAt = now
					if o.OnCancel != nil {
						o.OnCancel()
					}
				}

			case 2:
				ui.Row(c).Key("ic_done").Center().Gap(6).Children(func() {
					ui.Box(c).Width(13).Height(13).Draw(func(p *ui.Painter, r ui.Rect) {
						pfDrawCheck(p, r.X+r.W*0.5, r.Y+r.H*0.5, 11, 1, 1.8, t.Success)
					})
					ui.Text(c, o.DoneLabel).FontSize(12).FontWeight(600).TextColor(t.Success).NoWrap()
				})
			}
		})
	})

	return holder
}

// ============================================================================
// 6. toasts — Notify
// ============================================================================

// ToastKind classifies a ToastNotifier notification.
type ToastKind int

const (
	ToastSuccess ToastKind = iota
	ToastWarning
	ToastInfo
)

// ToastNotifierOptions configures ToastNotifier.
type ToastNotifierOptions struct {
	Width       float32
	MaxVisible  int
	AutoDismiss time.Duration
}

type pfToastEntry struct {
	id        int
	kind      ToastKind
	title     string
	subtitle  string
	createdAt time.Time
	sticky    bool
}

type pfToastState struct {
	nextID int
	items  []pfToastEntry
}

// ToastNotifier renders trigger buttons (Success, Warning, Info) that push
// animated stacked toast notifications sliding in and auto-dismissing.
func ToastNotifier(c *ui.Context, opts ...ToastNotifierOptions) ui.Element {
	var o ToastNotifierOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}
	if o.MaxVisible <= 0 {
		o.MaxVisible = 3
	}
	if o.AutoDismiss <= 0 {
		o.AutoDismiss = 3200 * time.Millisecond
	}

	root := ui.Column(c).
		Width(o.Width).
		Gap(10)

	now := c.Now()
	st := ui.Local(root, "pf_toasts_st", func() pfToastState {
		return pfToastState{
			nextID: 2,
			items: []pfToastEntry{
				{
					id:        1,
					kind:      ToastInfo,
					title:     "Notifications ready",
					subtitle:  "Click a toast to dismiss",
					createdAt: now,
					sticky:    true,
				},
			},
		}
	})

	// Expire non-sticky toasts past AutoDismiss.
	hasLiveTimer := false
	kept := st.items[:0]
	for _, it := range st.items {
		if !it.sticky && now.Sub(it.createdAt) >= o.AutoDismiss {
			continue
		}
		if !it.sticky {
			hasLiveTimer = true
		}
		kept = append(kept, it)
	}
	st.items = kept
	if hasLiveTimer {
		c.AnimationFrame()
	}

	push := func(kind ToastKind, title, subtitle string) {
		entry := pfToastEntry{
			id:        st.nextID,
			kind:      kind,
			title:     title,
			subtitle:  subtitle,
			createdAt: now,
		}
		st.nextID++
		st.items = append([]pfToastEntry{entry}, st.items...)
		if len(st.items) > o.MaxVisible {
			st.items = st.items[:o.MaxVisible]
		}
	}

	kindColor := func(k ToastKind) ui.Color {
		switch k {
		case ToastSuccess:
			return t.Success
		case ToastWarning:
			return t.Warning
		default:
			return t.Accent
		}
	}

	root.Children(func() {
		// Trigger buttons row
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).Gap(6).Children(func() {
			triggers := []struct {
				label string
				kind  ToastKind
				title string
				sub   string
			}{
				{"Success", ToastSuccess, "Deploy complete", "v2.4 live in 1.2s"},
				{"Warning", ToastWarning, "Rate limit 85%", "Burst window active"},
				{"Info", ToastInfo, "Sync finished", "14 items indexed"},
			}
			for i, tr := range triggers {
				item := tr
				col := kindColor(item.kind)
				btn := ui.Row(c).
					Key(i).
					Grow(1).
					Basis(0).
					Height(28).
					Center().
					Gap(5).
					Radius(8).
					Background(t.Surface).
					Border(1, t.Border).
					Cursor(ui.CursorPointer)
				if btn.Hovered() {
					btn.Background(t.SurfaceHover).Border(1, col.Alpha(0.55))
				}
				btn.Children(func() {
					ui.Box(c).Width(6).Height(6).Radius(999).Background(col)
					ui.Text(c, item.label).FontSize(11).FontWeight(600).TextColor(t.Text)
				})
				if btn.Clicked() {
					push(item.kind, item.title, item.sub)
				}
			}
		})

		// Stacked toast stage
		stage := ui.Box(c).
			WidthPercent(100).
			Height(114).
			Clip()

		stage.Children(func() {
			if len(st.items) == 0 {
				ui.Box(c).
					WidthPercent(100).
					HeightPercent(100).
					Center().
					Radius(12).
					Background(t.SurfaceSubtle.Alpha(0.5)).
					Border(1, t.Border).
					Children(func() {
						ui.Text(c, "No active notifications").FontSize(11).TextColor(t.TextSubtle)
					})
				return
			}

			dismissID := -1
			// Draw back-to-front so newest card renders on top
			for i := len(st.items) - 1; i >= 0; i-- {
				it := st.items[i]
				depth := float32(i)
				col := kindColor(it.kind)

				age := float32(now.Sub(it.createdAt).Seconds())
				enter := clamp32(age/0.22, 0, 1)
				if age < 0.25 {
					c.AnimationFrame()
				}

				toastCard := ui.Row(c).
					Key(it.id).
					Absolute().
					Left(0).
					Right(0).
					Top(0).
					Height(44).
					Padding(0, 10).
					AlignItems(ui.Center).
					Justify(ui.SpaceBetween).
					Gap(8).
					Radius(11).
					Background(t.Surface).
					Border(1, t.BorderStrong).
					Shadow(0, 4, 12, 0, t.Shadow).
					Cursor(ui.CursorPointer)

				yTarget := depth * 32
				yPos := toastCard.Spring("pf_toast_y", yTarget, ui.SpringOptions{
					Duration: 280 * time.Millisecond,
					Bounce:   0.2,
				})
				scTarget := 1 - depth*0.045
				sc := toastCard.Spring("pf_toast_sc", scTarget, ui.SpringOptions{
					Duration: 280 * time.Millisecond,
					Bounce:   0.18,
				})

				slideIn := (1 - enter) * -14
				alpha := clamp32(enter*(1-depth*0.18), 0.2, 1)

				toastCard.Translate(0, yPos+slideIn).Scale(sc, sc).Opacity(alpha).Children(func() {
					ui.Row(c).AlignItems(ui.Center).Gap(8).Grow(1).Shrink(1).MinWidth(0).Children(func() {
						ui.Box(c).Width(3).Height(22).Radius(2).Background(col)
						ui.Column(c).Gap(1).Children(func() {
							ui.Text(c, it.title).FontSize(11.5).FontWeight(600).TextColor(t.Text).NoWrap()
							ui.Text(c, it.subtitle).FontSize(10).TextColor(t.TextMuted).NoWrap()
						})
					})
					ui.Box(c).Width(14).Height(14).Draw(func(p *ui.Painter, r ui.Rect) {
						pfDrawCloseIcon(p, r.X+r.W*0.5, r.Y+r.H*0.5, 3.6, t.TextSubtle)
					})
				})

				if toastCard.Clicked() {
					dismissID = it.id
				}
			}

			if dismissID >= 0 {
				next := st.items[:0]
				for _, it := range st.items {
					if it.id != dismissID {
						next = append(next, it)
					}
				}
				st.items = next
			}
		})
	})

	return root
}

// ============================================================================
// 7. stepper — Drag stepper
// ============================================================================

// DragStepperOptions configures DragStepper.
type DragStepperOptions struct {
	Min           int
	Max           int
	Step          int
	PixelsPerStep float32
	Label         string
	Width         float32
	OnChange      func(int)
}

type pfStepperState struct {
	dragging bool
	startVal int
	startPX  float32
	kickDir  float32
	kickAt   time.Time
}

// DragStepper renders a numeric counter pill with -/+ press buttons and
// horizontal drag-to-scrub on the center value with spring bounce feedback.
func DragStepper(c *ui.Context, value *int, opts ...DragStepperOptions) ui.Element {
	var o DragStepperOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Max <= o.Min {
		o.Min = 0
		o.Max = 100
	}
	if o.Step <= 0 {
		o.Step = 1
	}
	if o.PixelsPerStep <= 0 {
		o.PixelsPerStep = 12
	}
	if o.Width <= 0 {
		o.Width = 204
	}
	if o.Label == "" {
		o.Label = "Drag or tap"
	}

	cur := o.Min
	if value != nil {
		cur = *value
		if cur < o.Min {
			cur = o.Min
		} else if cur > o.Max {
			cur = o.Max
		}
	}

	wrap := ui.Column(c).
		AlignItems(ui.Center).
		Gap(8)

	st := ui.Local(wrap, "pf_stepper_st", func() pfStepperState {
		return pfStepperState{}
	})

	now := c.Now()
	setVal := func(next int, dir float32) {
		if next < o.Min {
			next = o.Min
		}
		if next > o.Max {
			next = o.Max
		}
		if next != cur {
			cur = next
			if value != nil {
				*value = cur
			}
			if dir != 0 {
				st.kickDir = dir
				st.kickAt = now
			}
			if o.OnChange != nil {
				o.OnChange(cur)
			}
		}
	}

	kickActive := !st.kickAt.IsZero() && now.Sub(st.kickAt) < 120*time.Millisecond
	if kickActive {
		c.AnimationFrame()
	}

	wrap.Children(func() {
		pill := ui.Row(c).
			Width(o.Width).
			Height(48).
			Padding(5).
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween).
			Radius(999).
			Background(t.Surface).
			Border(1, t.BorderStrong)

		pill.Children(func() {
			// Minus button
			minusBtn := ui.Row(c).
				Key("st_minus").
				Width(36).
				Height(36).
				Center().
				Radius(999).
				Background(t.SurfaceSubtle).
				Cursor(ui.CursorPointer)
			if minusBtn.Hovered() {
				minusBtn.Background(t.SurfaceHover)
			}
			mScale := minusBtn.Spring("pf_st_ms", func() float32 {
				if minusBtn.Pressed() {
					return 0.88
				}
				return 1
			}(), ui.SpringOptions{Duration: 160 * time.Millisecond, Bounce: 0.3})
			minusBtn.Scale(mScale, mScale).Children(func() {
				ui.Text(c, "−").FontSize(18).FontWeight(600).TextColor(t.Text)
			})
			if minusBtn.Clicked() {
				setVal(cur-o.Step, -1)
			}

			// Draggable center scrub zone
			center := ui.Row(c).
				Key("st_center").
				Grow(1).
				Height(36).
				Center().
				Radius(999).
				Cursor(ui.CursorGrab)

			dx, _, dragging := center.Dragged()
			px, _, _ := center.PointerPosition()
			_ = dx

			var dragShiftTarget float32
			if dragging {
				if !st.dragging {
					st.dragging = true
					st.startVal = cur
					st.startPX = px
				}
				deltaX := px - st.startPX
				steps := int(math.Round(float64(deltaX / o.PixelsPerStep)))
				setVal(st.startVal+steps*o.Step, 0)
				dragShiftTarget = rubberBand(deltaX*0.35, 16)
			} else if st.dragging {
				st.dragging = false
			}

			kickTarget := float32(0)
			if kickActive {
				kickTarget = st.kickDir * 5
			}
			shiftX := center.Spring("pf_st_shift", dragShiftTarget+kickTarget, ui.SpringOptions{
				Duration: 260 * time.Millisecond,
				Bounce:   0.42,
			})
			valScale := center.Spring("pf_st_vsc", func() float32 {
				if dragging {
					return 1.12
				}
				return 1.0
			}(), ui.SpringOptions{Duration: 200 * time.Millisecond, Bounce: 0.32})

			centerBg := ui.Transparent
			if dragging {
				centerBg = t.AccentSoft
			}
			center.Background(centerBg).Translate(shiftX, 0).Scale(valScale, valScale).Children(func() {
				col := t.Text
				if dragging {
					col = t.Accent
				}
				ui.Text(c, fmt.Sprintf("%d", cur)).
					FontSize(19).
					FontWeight(700).
					TextColor(col)
			})

			// Plus button
			plusBtn := ui.Row(c).
				Key("st_plus").
				Width(36).
				Height(36).
				Center().
				Radius(999).
				Background(t.SurfaceSubtle).
				Cursor(ui.CursorPointer)
			if plusBtn.Hovered() {
				plusBtn.Background(t.SurfaceHover)
			}
			pScale := plusBtn.Spring("pf_st_ps", func() float32 {
				if plusBtn.Pressed() {
					return 0.88
				}
				return 1
			}(), ui.SpringOptions{Duration: 160 * time.Millisecond, Bounce: 0.3})
			plusBtn.Scale(pScale, pScale).Children(func() {
				ui.Text(c, "+").FontSize(18).FontWeight(600).TextColor(t.Text)
			})
			if plusBtn.Clicked() {
				setVal(cur+o.Step, 1)
			}
		})

		hint := o.Label
		hintCol := t.TextSubtle
		if st.dragging {
			hint = "Scrubbing value..."
			hintCol = t.Accent
		}
		ui.Text(c, hint).FontSize(10.5).FontWeight(500).TextColor(hintCol)
	})

	return wrap
}

// ============================================================================
// 8. like — Like
// ============================================================================

// LikeButtonOptions configures LikeButton.
type LikeButtonOptions struct {
	Color    ui.Color
	OnToggle func(liked bool, count int)
}

type pfLikeState struct {
	burstAt  time.Time
	countAt  time.Time
	countDir float32
}

// LikeButton renders a tactile Heart/Like button with spring pop scale,
// bursting micro-sparkle ring on ui.Painter, and animated counter transition.
func LikeButton(c *ui.Context, liked *bool, count *int, opts ...LikeButtonOptions) ui.Element {
	var o LikeButtonOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	heartCol := colorOr(o.Color, ui.Hex("#f43f5e"))

	isLiked := false
	if liked != nil {
		isLiked = *liked
	}
	curCount := 128
	if count != nil {
		curCount = *count
	}

	stage := ui.Box(c).
		Width(148).
		Height(76).
		Center()

	st := ui.Local(stage, "pf_like_st", func() pfLikeState {
		return pfLikeState{}
	})

	now := c.Now()
	burstElapsed := float32(99)
	if !st.burstAt.IsZero() {
		burstElapsed = float32(now.Sub(st.burstAt).Seconds())
	}
	burstDur := float32(0.55)
	burstActive := burstElapsed < burstDur
	if burstActive {
		c.AnimationFrame()
	}

	countElapsed := float32(99)
	if !st.countAt.IsZero() {
		countElapsed = float32(now.Sub(st.countAt).Seconds())
	}
	if countElapsed < 0.28 {
		c.AnimationFrame()
	}

	fillAmt := stage.Spring("pf_like_fill", func() float32 {
		if isLiked {
			return 1
		}
		return 0
	}(), ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.3})

	popTarget := float32(1)
	if burstActive && burstElapsed < 0.14 {
		popTarget = 1.32
	}
	heartScale := stage.Spring("pf_like_pop", popTarget, ui.SpringOptions{
		Duration: 300 * time.Millisecond,
		Bounce:   0.52,
	})

	stage.Children(func() {
		btn := ui.Row(c).
			Key("like_pill").
			Height(42).
			Padding(0, 16, 0, 12).
			Center().
			Gap(9).
			Radius(999).
			Cursor(ui.CursorPointer)

		btnScale := btn.Spring("pf_like_btn_sc", func() float32 {
			if btn.Pressed() {
				return 0.94
			}
			return 1
		}(), ui.SpringOptions{Duration: 150 * time.Millisecond, Bounce: 0.25})

		bg := lerpColor(t.Surface, heartCol.Alpha(0.14), fillAmt)
		bdr := lerpColor(t.BorderStrong, heartCol.Alpha(0.55), fillAmt)
		if !isLiked && btn.Hovered() {
			bg = t.SurfaceHover
		}

		btn.Background(bg).Border(1, bdr).Scale(btnScale, btnScale).Children(func() {
			// Heart + bursting micro-sparkle ring
			ui.Box(c).Width(22).Height(22).DrawOver(func(p *ui.Painter, r ui.Rect) {
				cx, cy := r.X+r.W*0.5, r.Y+r.H*0.5

				if burstActive {
					prog := clamp32(burstElapsed/burstDur, 0, 1)
					ease := 1 - (1-prog)*(1-prog)
					fade := clamp32(1-prog, 0, 1)

					// Expanding sparkle ring
					ringR := lerp32(7, 23, ease)
					var ring ui.Path
					ring.Circle(cx, cy, ringR)
					p.StrokePath(&ring, lerp32(2.2, 0.4, prog), heartCol.Alpha(fade*0.65))

					// 8 radiating dual-tone micro-sparkles
					sparkCol2 := ui.Hex("#fbbf24")
					for i := 0; i < 8; i++ {
						ang := float64(i)*math.Pi/4.0 + 0.22
						dist := lerp32(9, 27, ease)
						sx := cx + float32(math.Cos(ang))*dist
						sy := cy + float32(math.Sin(ang))*dist
						dotR := lerp32(2.6, 0.4, prog)
						col := heartCol
						if i%2 == 1 {
							col = sparkCol2
						}
						var dot ui.Path
						dot.Circle(sx, sy, dotR)
						p.FillPath(&dot, col.Alpha(fade*0.95))
					}
				}

				pfDrawHeart(p, cx, cy, 8.2*heartScale, fillAmt, heartCol, t.TextMuted)
			})

			// Animated counter
			cProg := clamp32(countElapsed/0.22, 0, 1)
			cOffset := (1 - cProg) * st.countDir * 8
			cAlpha := 0.45 + cProg*0.55
			txtCol := lerpColor(t.Text, heartCol, fillAmt)

			ui.Text(c, fmt.Sprintf("%d", curCount)).
				FontSize(14).
				FontWeight(700).
				TextColor(txtCol).
				Translate(0, cOffset).
				Opacity(cAlpha)
		})

		if btn.Clicked() {
			nextLiked := !isLiked
			if nextLiked {
				curCount++
				st.burstAt = now
				st.countDir = -1
			} else {
				if curCount > 0 {
					curCount--
				}
				st.countDir = 1
			}
			st.countAt = now
			if liked != nil {
				*liked = nextLiked
			}
			if count != nil {
				*count = curCount
			}
			if o.OnToggle != nil {
				o.OnToggle(nextLiked, curCount)
			}
		}
	})

	return stage
}

// ============================================================================
// 9. checklist — Checklist
// ============================================================================

// ChecklistItem represents one task in Checklist.
type ChecklistItem struct {
	Label string
	Done  bool
}

// ChecklistOptions configures Checklist.
type ChecklistOptions struct {
	Title    string
	Width    float32
	OnToggle func(index int, done bool)
}

// Checklist renders an interactive task checklist with animated checkbox fills,
// strikethrough styling, and a live completion progress bar.
func Checklist(c *ui.Context, items []ChecklistItem, opts ...ChecklistOptions) ui.Element {
	var o ChecklistOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 282
	}
	if o.Title == "" {
		o.Title = "Release prep"
	}

	doneCount := 0
	for _, it := range items {
		if it.Done {
			doneCount++
		}
	}
	fracTarget := float32(0)
	if len(items) > 0 {
		fracTarget = float32(doneCount) / float32(len(items))
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(12).
		Gap(8).
		Radius(15).
		Background(t.Surface).
		Border(1, t.Border)

	prog := card.Spring("pf_cl_prog", fracTarget, ui.SpringOptions{
		Duration: 320 * time.Millisecond,
		Bounce:   0.18,
	})
	allDone := len(items) > 0 && doneCount == len(items)

	card.Children(func() {
		// Header + completion count
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, o.Title).FontSize(12).FontWeight(600).TextColor(t.Text)
			badgeBg := t.SurfaceSubtle
			badgeCol := t.TextMuted
			if allDone {
				badgeBg = t.Success.Alpha(0.16)
				badgeCol = t.Success
			}
			ui.Row(c).Height(20).Padding(0, 7).Center().Radius(999).Background(badgeBg).Children(func() {
				ui.Text(c, fmt.Sprintf("%d/%d", doneCount, len(items))).
					FontSize(10.5).
					FontWeight(600).
					TextColor(badgeCol)
			})
		})

		// Live progress bar
		ui.Box(c).WidthPercent(100).Height(5).Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, t.SurfaceSubtle, r.H*0.5)
			w := clamp32(prog, 0, 1) * r.W
			if w > 0.5 {
				fillCol := t.Accent
				if allDone {
					fillCol = t.Success
				}
				p.Fill(ui.Rect{X: r.X, Y: r.Y, W: w, H: r.H}, fillCol, r.H*0.5)
			}
		})

		// Task rows
		ui.Column(c).WidthPercent(100).Gap(3).Children(func() {
			for i := range items {
				idx := i
				done := items[idx].Done

				row := ui.Row(c).
					Key(idx).
					WidthPercent(100).
					Height(27).
					Padding(0, 6).
					AlignItems(ui.Center).
					Gap(8).
					Radius(8).
					Cursor(ui.CursorPointer)
				if row.Hovered() {
					row.Background(t.SurfaceSubtle)
				}

				chkAmt := row.Spring("pf_cl_item", func() float32 {
					if done {
						return 1
					}
					return 0
				}(), ui.SpringOptions{Duration: 240 * time.Millisecond, Bounce: 0.28})

				row.Children(func() {
					ui.Box(c).Width(16).Height(16).Draw(func(p *ui.Painter, r ui.Rect) {
						boxCol := lerpColor(t.SurfaceSubtle, t.Success, chkAmt)
						bdrCol := lerpColor(t.BorderStrong, t.Success, chkAmt)
						p.Fill(r, boxCol, 4.5)
						p.Stroke(r, bdrCol, 4.5, 1.2)
						if chkAmt > 0.1 {
							pfDrawCheck(p, r.X+r.W*0.5, r.Y+r.H*0.5, 10, chkAmt, 1.8, ui.Hex("#ffffff"))
						}
					})

					txt := ui.Text(c, items[idx].Label).
						FontSize(12).
						FontWeight(500).
						TextColor(lerpColor(t.Text, t.TextSubtle, chkAmt))
					if done {
						txt.Strikethrough()
					}
				})

				if row.Clicked() {
					items[idx].Done = !items[idx].Done
					if o.OnToggle != nil {
						o.OnToggle(idx, items[idx].Done)
					}
				}
			}
		})
	})

	return card
}

// ============================================================================
// 10. todo-tower — Todo tower
// ============================================================================

// TodoTowerOptions configures TodoTower.
type TodoTowerOptions struct {
	Width      float32
	OnComplete func(index int, task string)
	OnReset    func()
}

type pfTodoTowerState struct {
	completed  int
	peeling    bool
	peelAt     time.Time
	lastPeeled string
}

func pfDefaultTowerTasks() []string {
	return []string{
		"Review GPU shader cache",
		"Ship spring curve presets",
		"Trim dark mode borders",
		"Tag v0.10 release notes",
	}
}

// TodoTower renders a stacked card tower of pending tasks where clicking peels
// and drops the top task card with a live progress counter and Reset button.
func TodoTower(c *ui.Context, tasks []string, opts ...TodoTowerOptions) ui.Element {
	var o TodoTowerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if len(tasks) == 0 {
		tasks = pfDefaultTowerTasks()
	}
	if o.Width <= 0 {
		o.Width = 276
	}

	root := ui.Column(c).
		Width(o.Width).
		Gap(10)

	st := ui.Local(root, "pf_todo_tower_st", func() pfTodoTowerState {
		return pfTodoTowerState{}
	})

	now := c.Now()
	peelDur := float32(0.36)
	peelProg := float32(1)
	if st.peeling {
		peelProg = clamp32(float32(now.Sub(st.peelAt).Seconds())/peelDur, 0, 1)
		if peelProg < 1 {
			c.AnimationFrame()
		} else {
			st.peeling = false
		}
	}

	remaining := len(tasks) - st.completed
	if remaining < 0 {
		remaining = 0
	}

	root.Children(func() {
		// Header with progress counter + Reset button
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
				ui.Box(c).Width(7).Height(7).Radius(999).Background(func() ui.Color {
					if remaining == 0 {
						return t.Success
					}
					return t.Accent
				}())
				ui.Text(c, fmt.Sprintf("%d / %d done", st.completed, len(tasks))).
					FontSize(11.5).
					FontWeight(600).
					TextColor(t.Text)
			})

			resetBtn := ui.Row(c).
				Key("tt_reset").
				Height(24).
				Padding(0, 9).
				Center().
				Radius(999).
				Background(t.SurfaceSubtle).
				Border(1, t.Border).
				Cursor(ui.CursorPointer)
			if st.completed == 0 {
				resetBtn.Opacity(0.45)
			} else if resetBtn.Hovered() {
				resetBtn.Background(t.SurfaceHover)
			}
			resetBtn.Children(func() {
				ui.Text(c, "Reset").FontSize(10.5).FontWeight(600).TextColor(t.TextMuted)
			})
			if resetBtn.Clicked() && st.completed > 0 {
				st.completed = 0
				st.peeling = false
				if o.OnReset != nil {
					o.OnReset()
				}
			}
		})

		// Stacked tower stage
		stage := ui.Box(c).
			WidthPercent(100).
			Height(106)

		stage.Children(func() {
			if remaining == 0 {
				doneCard := ui.Column(c).
					Key("tt_clear").
					WidthPercent(100).
					Height(84).
					Center().
					Gap(4).
					Radius(14).
					Background(t.Success.Alpha(0.12)).
					Border(1, t.Success.Alpha(0.50)).
					Cursor(ui.CursorPointer)
				doneCard.Children(func() {
					ui.Text(c, "Tower cleared!").FontSize(13).FontWeight(700).TextColor(t.Success)
					ui.Text(c, "Click to rebuild stack").FontSize(10.5).TextColor(t.TextMuted)
				})
				if doneCard.Clicked() {
					st.completed = 0
					st.peeling = false
					if o.OnReset != nil {
						o.OnReset()
					}
				}
			} else {
				visible := remaining
				if visible > 3 {
					visible = 3
				}
				// Render bottom-to-top so depth=0 sits on top
				for d := visible - 1; d >= 0; d-- {
					taskIdx := st.completed + d
					isTop := d == 0
					depth := float32(d)

					card := ui.Row(c).
						Key(taskIdx).
						Absolute().
						Left(0).
						Right(0).
						Top(0).
						Height(62).
						Padding(0, 14).
						AlignItems(ui.Center).
						Justify(ui.SpaceBetween).
						Gap(10).
						Radius(14).
						Shadow(0, 6, 14, 0, t.Shadow)

					yPos := card.Spring("pf_tt_y", depth*14, ui.SpringOptions{
						Duration: 300 * time.Millisecond,
						Bounce:   0.26,
					})
					sc := card.Spring("pf_tt_sc", 1-depth*0.055, ui.SpringOptions{
						Duration: 300 * time.Millisecond,
						Bounce:   0.24,
					})

					bg := t.Surface
					if !isTop {
						bg = lerpColor(t.Surface, t.SurfaceSubtle, clamp32(depth*0.45, 0, 0.8))
					} else if card.Hovered() {
						bg = t.SurfaceHover
					}
					bdr := t.Border
					if isTop {
						bdr = t.BorderStrong
						card.Cursor(ui.CursorPointer)
					}

					card.Background(bg).Border(1, bdr).Translate(0, yPos).Scale(sc, sc).Children(func() {
						ui.Column(c).Gap(3).Grow(1).Shrink(1).MinWidth(0).Children(func() {
							ui.Text(c, fmt.Sprintf("TASK %02d", taskIdx+1)).
								FontSize(9.5).
								FontWeight(700).
								TextColor(t.Accent)
							ui.Text(c, tasks[taskIdx]).
								FontSize(12.5).
								FontWeight(600).
								TextColor(t.Text).
								NoWrap()
						})
						if isTop {
							ui.Box(c).
								Width(24).
								Height(24).
								Center().
								Radius(999).
								Background(t.AccentSoft).
								Draw(func(p *ui.Painter, r ui.Rect) {
									pfDrawCheck(p, r.X+r.W*0.5, r.Y+r.H*0.5, 11, 1, 1.7, t.Accent)
								})
						}
					})

					if isTop && card.Clicked() {
						completedTask := tasks[taskIdx]
						st.lastPeeled = completedTask
						st.peeling = true
						st.peelAt = now
						st.completed++
						if o.OnComplete != nil {
							o.OnComplete(taskIdx, completedTask)
						}
					}
				}
			}

			// Peeling/dropping overlay for the just-completed top card
			if st.peeling && peelProg < 1 {
				ease := peelProg * peelProg
				peelCard := ui.Row(c).
					Key("tt_peel").
					Absolute().
					Left(0).
					Right(0).
					Top(0).
					Height(62).
					Padding(0, 14).
					AlignItems(ui.Center).
					Justify(ui.SpaceBetween).
					Radius(14).
					Background(t.Success.Alpha(0.16)).
					Border(1, t.Success.Alpha(0.60)).
					Translate(ease*36, ease*54).
					Rotate(ease * 12).
					Opacity(clamp32(1-peelProg, 0, 1))
				peelCard.Children(func() {
					ui.Text(c, st.lastPeeled).FontSize(12.5).FontWeight(600).TextColor(t.Success).Strikethrough()
					ui.Box(c).Width(16).Height(16).Draw(func(p *ui.Painter, r ui.Rect) {
						pfDrawCheck(p, r.X+r.W*0.5, r.Y+r.H*0.5, 12, 1, 1.9, t.Success)
					})
				})
			}
		})
	})

	return root
}

// ============================================================================
// 11. palette — Palette
// ============================================================================

// PaletteSwatch defines one selectable color in Palette.
type PaletteSwatch struct {
	Name string
	Hex  string
}

// PaletteOptions configures Palette.
type PaletteOptions struct {
	Width    float32
	OnSelect func(index int, swatch PaletteSwatch)
}

type pfPaletteState struct {
	copiedAt time.Time
}

func pfDefaultSwatches() []PaletteSwatch {
	return []PaletteSwatch{
		{Name: "Indigo", Hex: "#6366F1"},
		{Name: "Rose", Hex: "#F43F5E"},
		{Name: "Amber", Hex: "#F59E0B"},
		{Name: "Emerald", Hex: "#10B981"},
		{Name: "Sky", Hex: "#0EA5E9"},
		{Name: "Slate", Hex: "#64748B"},
	}
}

// Palette renders an interactive color swatch strip where clicking selects and
// copies a swatch's hex code, springs the swatch larger, and updates a preview chip.
func Palette(c *ui.Context, swatches []PaletteSwatch, selected *int, opts ...PaletteOptions) ui.Element {
	var o PaletteOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if len(swatches) == 0 {
		swatches = pfDefaultSwatches()
	}
	if o.Width <= 0 {
		o.Width = 282
	}

	selIdx := 0
	if selected != nil && *selected >= 0 && *selected < len(swatches) {
		selIdx = *selected
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(14).
		Gap(12).
		Radius(16).
		Background(t.Surface).
		Border(1, t.Border)

	st := ui.Local(card, "pf_palette_st", func() pfPaletteState {
		return pfPaletteState{}
	})

	now := c.Now()
	justCopied := !st.copiedAt.IsZero() && now.Sub(st.copiedAt) < 1500*time.Millisecond
	if justCopied {
		c.AnimationFrame()
	}

	cur := swatches[selIdx]
	curCol := ui.Hex(cur.Hex)

	card.Children(func() {
		// Swatch strip
		ui.Row(c).
			WidthPercent(100).
			Height(48).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Children(func() {
				for i, sw := range swatches {
					idx := i
					isSel := idx == selIdx
					col := ui.Hex(sw.Hex)

					item := ui.Box(c).
						Key(idx).
						Width(34).
						Height(34).
						Center().
						Radius(11).
						Cursor(ui.CursorPointer)

					scTarget := float32(1)
					if isSel {
						scTarget = 1.24
					} else if item.Hovered() {
						scTarget = 1.10
					}
					sc := item.Spring("pf_pal_sc", scTarget, ui.SpringOptions{
						Duration: 260 * time.Millisecond,
						Bounce:   0.38,
					})

					item.Scale(sc, sc).Draw(func(p *ui.Painter, r ui.Rect) {
						if isSel {
							outer := ui.Rect{X: r.X - 2.5, Y: r.Y - 2.5, W: r.W + 5, H: r.H + 5}
							p.Stroke(outer, col.Alpha(0.70), 13, 1.8)
						}
						p.Fill(r, col, 10)
						p.Stroke(r, ui.Hex("#ffffff").Alpha(0.22), 10, 1)
						if isSel {
							pfDrawCheck(p, r.X+r.W*0.5, r.Y+r.H*0.5, 11, 1, 2.0, ui.Hex("#ffffff"))
						}
					})

					if item.Clicked() {
						if selected != nil {
							*selected = idx
						}
						c.WriteClipboard(sw.Hex)
						st.copiedAt = now
						if o.OnSelect != nil {
							o.OnSelect(idx, sw)
						}
					}
				}
			})

		// Live preview chip
		chip := ui.Row(c).
			WidthPercent(100).
			Height(36).
			Padding(0, 10).
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween).
			Radius(10).
			Background(t.SurfaceSubtle).
			Border(1, curCol.Alpha(0.45))

		chip.Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				ui.Box(c).Width(14).Height(14).Radius(999).Background(curCol)
				ui.Text(c, cur.Name).FontSize(12).FontWeight(600).TextColor(t.Text)
				ui.Text(c, cur.Hex).Font("monospace").FontSize(11).TextColor(t.TextMuted)
			})

			statusText := "Click swatch"
			statusCol := t.TextSubtle
			if justCopied {
				statusText = "Copied!"
				statusCol = t.Success
			}
			ui.Text(c, statusText).FontSize(10.5).FontWeight(600).TextColor(statusCol)
		})
	})

	return card
}

// ============================================================================
// 12. seek — Search
// ============================================================================

// SeekSearchOptions configures SeekSearch.
type SeekSearchOptions struct {
	Placeholder string
	Width       float32
	MaxChips    int
	OnSelect    func(string)
}

type pfSeekState struct {
	open   bool
	query  string
	picked string
}

func pfDefaultSeekItems() []string {
	return []string{"Springs", "Shaders", "Gestures", "Palette", "Toasts", "Layout"}
}

// SeekSearch renders a compact search pill button that smoothly springs open
// into an expandable search input bar with live match chips and a close/clear control.
func SeekSearch(c *ui.Context, items []string, opts ...SeekSearchOptions) ui.Element {
	var o SeekSearchOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if len(items) == 0 {
		items = pfDefaultSeekItems()
	}
	if o.Width <= 0 {
		o.Width = 286
	}
	if o.Placeholder == "" {
		o.Placeholder = "Seek blocks..."
	}
	if o.MaxChips <= 0 {
		o.MaxChips = 5
	}

	root := ui.Column(c).
		Width(o.Width).
		AlignItems(ui.Center).
		Gap(8)

	st := ui.Local(root, "pf_seek_st", func() pfSeekState {
		return pfSeekState{}
	})

	openTarget := float32(0)
	if st.open {
		openTarget = 1
	}
	openAmt := root.Spring("pf_seek_open", openTarget, ui.SpringOptions{
		Duration: 300 * time.Millisecond,
		Bounce:   0.22,
	})

	barW := lerp32(114, o.Width, openAmt)

	// Filter live matching items
	q := strings.ToLower(strings.TrimSpace(st.query))
	var matches []string
	for _, it := range items {
		if q == "" || strings.Contains(strings.ToLower(it), q) {
			matches = append(matches, it)
			if len(matches) >= o.MaxChips {
				break
			}
		}
	}

	root.Children(func() {
		bar := ui.Row(c).
			Width(barW).
			Height(38).
			Padding(0, 10).
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween).
			Gap(7).
			Radius(999).
			Background(t.Surface).
			Border(1, lerpColor(t.BorderStrong, t.Accent.Alpha(0.65), openAmt)).
			Clip()

		if !st.open {
			bar.Cursor(ui.CursorPointer)
			if bar.Hovered() {
				bar.Background(t.SurfaceHover)
			}
			bar.Children(func() {
				ui.Row(c).WidthPercent(100).Center().Gap(7).Children(func() {
					ui.Box(c).Width(14).Height(14).Draw(func(p *ui.Painter, r ui.Rect) {
						pfDrawSearchIcon(p, r.X+r.W*0.5, r.Y+r.H*0.5, 6.5, t.TextMuted)
					})
					lbl := "Search"
					if st.picked != "" {
						lbl = st.picked
					}
					ui.Text(c, lbl).FontSize(12.5).FontWeight(600).TextColor(t.Text).NoWrap()
				})
			})
			if bar.Clicked() {
				st.open = true
			}
		} else {
			bar.Children(func() {
				ui.Box(c).Width(14).Height(14).Draw(func(p *ui.Painter, r ui.Rect) {
					pfDrawSearchIcon(p, r.X+r.W*0.5, r.Y+r.H*0.5, 6.5, t.Accent)
				})

				inp := ui.TextInput(c, &st.query).
					Width(clamp32(barW-72, 60, o.Width-72)).
					Placeholder(o.Placeholder).
					FontSize(12)
				if inp.Submitted() && len(matches) > 0 {
					st.picked = matches[0]
					st.query = matches[0]
					if o.OnSelect != nil {
						o.OnSelect(matches[0])
					}
				}

				clearBtn := ui.Box(c).
					Key("seek_close").
					Width(22).
					Height(22).
					Center().
					Radius(999).
					Background(t.SurfaceSubtle).
					Cursor(ui.CursorPointer).
					Draw(func(p *ui.Painter, r ui.Rect) {
						pfDrawCloseIcon(p, r.X+r.W*0.5, r.Y+r.H*0.5, 3.8, t.TextMuted)
					})
				if clearBtn.Clicked() {
					if st.query != "" {
						st.query = ""
					} else {
						st.open = false
					}
				}
			})
		}

		// Live match chips when expanded
		if openAmt > 0.15 {
			chipsWrap := ui.Row(c).
				WidthPercent(100).
				Wrap().
				Gap(6).
				Opacity(clamp32((openAmt-0.15)/0.85, 0, 1)).
				Translate(0, (1-openAmt)*-6)

			chipsWrap.Children(func() {
				if len(matches) == 0 {
					ui.Text(c, "No matching blocks").FontSize(11).TextColor(t.TextSubtle)
					return
				}
				for i, m := range matches {
					val := m
					active := st.picked == val
					chip := ui.Row(c).
						Key(i).
						Height(25).
						Padding(0, 10).
						Center().
						Radius(999).
						Cursor(ui.CursorPointer)
					if active {
						chip.Background(t.AccentSoft).Border(1, t.Accent.Alpha(0.60))
					} else {
						chip.Background(t.Surface).Border(1, t.Border)
						if chip.Hovered() {
							chip.Background(t.SurfaceHover)
						}
					}
					chip.Children(func() {
						col := t.TextMuted
						if active {
							col = t.Accent
						}
						ui.Text(c, val).FontSize(11).FontWeight(500).TextColor(col)
					})
					if chip.Clicked() {
						st.picked = val
						st.query = val
						if o.OnSelect != nil {
							o.OnSelect(val)
						}
					}
				}
			})
		}
	})

	return root
}

// ============================================================================
// Demo state structs & catalog registrations
// ============================================================================

type pfPriceFlowDemoState struct {
	active int
}

type pfInlineConfirmDemoState struct {
	deletedCount int
}

type pfStepperDemoState struct {
	val int
}

type pfLikeDemoState struct {
	liked bool
	count int
}

type pfChecklistDemoState struct {
	items []ChecklistItem
}

type pfPaletteDemoState struct {
	selected int
}

func init() {
	registerDemo(
		"copy-field",
		"Monospace token field with a tactile Copy button that writes clipboard and flips to an emerald Copied pill.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				CopyField(c, "mg_live_9f4a8c2e7b1d3048", CopyFieldOptions{
					Label: "API TOKEN",
				})
			})
		},
	)

	registerDemo(
		"feedback-bar",
		"Compact Was this helpful? rating bar that expands inline for quick tags and confirmation.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				FeedbackBar(c)
			})
		},
	)

	registerDemo(
		"price-flow",
		"Animated rolling-digit price tag cycling between Starter, Pro, and Enterprise tiers with spring interpolation.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_pf_price_flow", func() pfPriceFlowDemoState {
				return pfPriceFlowDemoState{active: 1}
			})
			holder.Children(func() {
				PriceFlow(c, pfDefaultPriceTiers(), &st.active)
			})
		},
	)

	registerDemo(
		"confirm",
		"Compact destructive inline confirmation that smoothly expands into a Sure? Yes / No pill.",
		func(c *ui.Context) {
			holder := ui.Column(c).AlignItems(ui.Center).Gap(6)
			st := ui.Local(holder, "demo_pf_inline_confirm", func() pfInlineConfirmDemoState {
				return pfInlineConfirmDemoState{}
			})
			holder.Children(func() {
				InlineConfirm(c, func() {
					st.deletedCount++
				})
				if st.deletedCount > 0 {
					t := ThemeFrom(c)
					ui.Text(c, fmt.Sprintf("Confirmed %dx", st.deletedCount)).FontSize(10.5).TextColor(t.TextSubtle)
				}
			})
		},
	)

	registerDemo(
		"toasts",
		"Trigger buttons that push animated stacked toast notifications sliding in and auto-dismissing.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				ToastNotifier(c)
			})
		},
	)

	registerDemo(
		"stepper",
		"Numeric counter pill with -/+ press buttons and horizontal drag-to-scrub on the center value.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_pf_stepper", func() pfStepperDemoState {
				return pfStepperDemoState{val: 24}
			})
			holder.Children(func() {
				DragStepper(c, &st.val, DragStepperOptions{
					Min:   0,
					Max:   99,
					Step:  1,
					Label: "Tap −/+ or drag center",
				})
			})
		},
	)

	registerDemo(
		"like",
		"Heart/Like button with spring pop scale, bursting micro-sparkle ring on ui.Painter, and animated counter.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_pf_like", func() pfLikeDemoState {
				return pfLikeDemoState{liked: false, count: 142}
			})
			holder.Children(func() {
				LikeButton(c, &st.liked, &st.count)
			})
		},
	)

	registerDemo(
		"checklist",
		"Interactive task checklist with animated checkbox fills, strikethrough styling, and live progress bar.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_pf_checklist", func() pfChecklistDemoState {
				return pfChecklistDemoState{
					items: []ChecklistItem{
						{Label: "Audit spring timings", Done: true},
						{Label: "Verify dark contrast", Done: true},
						{Label: "Record demo clip", Done: false},
						{Label: "Publish catalog", Done: false},
					},
				}
			})
			holder.Children(func() {
				Checklist(c, st.items)
			})
		},
	)

	registerDemo(
		"todo-tower",
		"Stacked card tower of pending tasks where clicking peels and drops the top card with progress and Reset.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				TodoTower(c, pfDefaultTowerTasks())
			})
		},
	)

	registerDemo(
		"palette",
		"Interactive color swatch palette that copies hex, springs the active swatch larger, and updates a preview chip.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_pf_palette", func() pfPaletteDemoState {
				return pfPaletteDemoState{selected: 0}
			})
			holder.Children(func() {
				Palette(c, pfDefaultSwatches(), &st.selected)
			})
		},
	)

	registerDemo(
		"seek",
		"Compact search pill button that springs open into an expandable search bar with live match chips.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				SeekSearch(c, pfDefaultSeekItems())
			})
		},
	)
}
