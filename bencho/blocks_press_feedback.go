package bencho

import (
	"fmt"
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

}
