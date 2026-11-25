package bencho

import (
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
		"toasts",
		"Trigger buttons that push animated stacked toast notifications sliding in and auto-dismissing.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				ToastNotifier(c)
			})
		},
	)

}
