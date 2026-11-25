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

func pfDrawSearchIcon(p *ui.Painter, cx, cy, r float32, col ui.Color) {
	var ring ui.Path
	ring.Circle(cx-r*0.16, cy-r*0.16, r*0.56)
	p.StrokePath(&ring, 1.55, col)
	p.Line(cx+r*0.24, cy+r*0.24, cx+r*0.72, cy+r*0.72, 1.65, col)
}

func pfDrawCloseIcon(p *ui.Painter, cx, cy, r float32, col ui.Color) {
	p.Line(cx-r, cy-r, cx+r, cy+r, 1.5, col)
	p.Line(cx+r, cy-r, cx-r, cy+r, 1.5, col)
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

type pfPaletteDemoState struct {
	selected int
}

func init() {
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
