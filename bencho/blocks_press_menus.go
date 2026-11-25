package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// -----------------------------------------------------------------------------
// Shared small vector painters for press menu, accordion & card blocks
// -----------------------------------------------------------------------------

func drawChevron(p *ui.Painter, r ui.Rect, angleDeg float32, col ui.Color) {
	cx := r.X + r.W*0.5
	cy := r.Y + r.H*0.5
	rad := float64(angleDeg) * math.Pi / 180.0
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))
	rot := func(dx, dy float32) [2]float32 {
		return [2]float32{cx + dx*cos - dy*sin, cy + dx*sin + dy*cos}
	}
	// Down-pointing chevron by default; angleDeg rotates smoothly.
	pts := [][2]float32{
		rot(-3.5, -1.5),
		rot(0, 2.0),
		rot(3.5, -1.5),
	}
	p.StrokePolyline(pts, 1.5, col)
}

func drawCheckmark(p *ui.Painter, r ui.Rect, col ui.Color) {
	cx := r.X + r.W*0.5
	cy := r.Y + r.H*0.5
	pts := [][2]float32{
		{cx - 3.5, cy + 0.2},
		{cx - 1.0, cy + 2.8},
		{cx + 3.8, cy - 2.6},
	}
	p.StrokePolyline(pts, 1.6, col)
}

func drawPlusMinus(p *ui.Painter, r ui.Rect, openP float32, col ui.Color) {
	cx := r.X + r.W*0.5
	cy := r.Y + r.H*0.5
	half := float32(4.2)
	// Horizontal bar stays constant.
	p.Line(cx-half, cy, cx+half, cy, 1.5, col)
	// Vertical bar rotates 90 degrees into the horizontal bar as openP -> 1.
	rad := float64(clamp32(openP, 0, 1)) * (math.Pi * 0.5)
	dx := float32(math.Sin(rad)) * half
	dy := float32(math.Cos(rad)) * half
	alpha := clamp32(1-openP*0.85, 0, 1)
	if alpha > 0.02 {
		p.Line(cx-dx, cy-dy, cx+dx, cy+dy, 1.5, col.Alpha(alpha))
	}
}

func drawPlusRotated(p *ui.Painter, r ui.Rect, angleDeg float32, half float32, col ui.Color) {
	cx := r.X + r.W*0.5
	cy := r.Y + r.H*0.5
	rad := float64(angleDeg) * math.Pi / 180.0
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))
	p.Line(cx-half*cos, cy-half*sin, cx+half*cos, cy+half*sin, 1.8, col)
	p.Line(cx+half*sin, cy-half*cos, cx-half*sin, cy+half*cos, 1.8, col)
}

// -----------------------------------------------------------------------------
// 1. faq-frames ("FAQ frames")
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// 9. selection-actions ("Selection actions")
// -----------------------------------------------------------------------------

// SelectableItem represents one selectable row in SelectionActions.
type SelectableItem struct {
	ID   string
	Name string
	Meta string
}

// SelectionActionsOptions configures SelectionActions.
type SelectionActionsOptions struct {
	Width    float32
	Accent   ui.Color
	OnAction func(action string, selectedIDs []string)
}

// SelectionActions renders a multi-select item list paired with a floating bulk
// action bar (Archive, Tag, Delete, Clear) that springs in when >=1 item is selected.
func SelectionActions(c *ui.Context, items []SelectableItem, selected map[string]bool, opts ...SelectionActionsOptions) ui.Element {
	var o SelectionActionsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}
	accent := colorOr(o.Accent, t.Accent)

	selIDs := make([]string, 0, len(items))
	for _, it := range items {
		if selected[it.ID] {
			selIDs = append(selIDs, it.ID)
		}
	}
	selCount := len(selIDs)

	root := ui.Box(c).
		Width(o.Width).
		Height(180).
		Clip()

	barTarget := float32(0)
	if selCount > 0 {
		barTarget = 1
	}
	barP := clamp32(root.Spring("sel_bar_p", barTarget, ui.SpringOptions{
		Duration: 320 * time.Millisecond,
		Bounce:   0.22,
	}), 0, 1)

	root.Children(func() {
		// Item list
		list := ui.Column(c).
			WidthPercent(100).
			Radius(12).
			Background(t.Surface).
			Border(1, t.Border).
			Clip()

		list.Children(func() {
			for i := range items {
				it := items[i]
				isSel := selected[it.ID]
				isLast := i == len(items)-1

				row := ui.Row(c).
					Key(it.ID).
					WidthPercent(100).
					Height(32).
					Padding(0, 10).
					Justify(ui.SpaceBetween).
					AlignItems(ui.Center).
					Cursor(ui.CursorPointer)

				if !isLast {
					row.BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
				}
				if isSel {
					row.Background(t.AccentSoft)
				} else if row.Hovered() {
					row.Background(t.SurfaceSubtle)
				}

				row.Children(func() {
					ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
						boxBg := t.Surface
						boxBorder := t.BorderStrong
						if isSel {
							boxBg = accent
							boxBorder = accent
						}
						cb := ui.Box(c).
							Size(15, 15).
							Radius(4).
							Background(boxBg).
							Border(1, boxBorder)
						if isSel {
							cb.Draw(func(p *ui.Painter, r ui.Rect) {
								drawCheckmark(p, r, t.AccentText)
							})
						}
						txtCol := t.Text
						weight := 500
						if isSel {
							txtCol = accent
							weight = 600
						}
						ui.Text(c, it.Name).FontSize(11.5).FontWeight(weight).TextColor(txtCol).SingleLine()
					})
					ui.Text(c, it.Meta).FontSize(10).TextColor(t.TextSubtle).SingleLine()
				})

				if row.Clicked() && selected != nil {
					if isSel {
						delete(selected, it.ID)
					} else {
						selected[it.ID] = true
					}
				}
			}
		})

		// Floating bulk action bar
		if barP > 0.02 {
			bar := ui.Row(c).
				Key("bulk_bar").
				Absolute().
				Left(2).
				Right(2).
				Bottom(4).
				Height(36).
				Padding(0, 8).
				Radius(11).
				Justify(ui.SpaceBetween).
				AlignItems(ui.Center).
				Background(t.Text).
				Shadow(0, 6, 16, 0, t.Shadow).
				Translate(0, lerp32(28, 0, barP)).
				Opacity(barP).
				IgnoreInput(selCount == 0)

			bar.Children(func() {
				ui.Box(c).Padding(2, 6).Radius(999).Background(accent).Children(func() {
					ui.Text(c, fmt.Sprintf("%d", selCount)).FontSize(10.5).FontWeight(700).TextColor(t.AccentText)
				})

				ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
					actions := []struct {
						name   string
						danger bool
					}{
						{name: "Archive"},
						{name: "Tag"},
						{name: "Delete", danger: true},
						{name: "Clear"},
					}
					for _, act := range actions {
						a := act
						btn := ui.Box(c).
							Key(a.name).
							Padding(4, 7).
							Radius(7).
							Cursor(ui.CursorPointer)
						ink := t.Background
						if a.danger {
							ink = t.Danger
						}
						if btn.Hovered() {
							btn.Background(t.Background.Alpha(0.14))
						}
						btn.Children(func() {
							ui.Text(c, a.name).FontSize(10.5).FontWeight(600).TextColor(ink)
						})
						if btn.Clicked() {
							if a.name == "Clear" {
								for k := range selected {
									delete(selected, k)
								}
							} else if o.OnAction != nil {
								o.OnAction(a.name, selIDs)
							}
						}
					}
				})
			})
		}
	})

	return root
}

// -----------------------------------------------------------------------------
// 10. connect-sources ("Integrations")
// -----------------------------------------------------------------------------

// IntegrationSource describes one integration tile in ConnectSources.
type IntegrationSource struct {
	ID        string
	Name      string
	Badge     string
	Color     ui.Color
	Connected bool
}

// ConnectSourcesOptions configures ConnectSources.
type ConnectSourcesOptions struct {
	Width    float32
	OnToggle func(id string, connected bool)
}

// ConnectSources renders a 2x2 grid of integration source tiles (GitHub, Linear,
// Slack, Figma) with interactive Connect/Connected toggles, live status dots,
// and a connected count badge.
func ConnectSources(c *ui.Context, sources []IntegrationSource, opts ...ConnectSourcesOptions) ui.Element {
	var o ConnectSourcesOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}

	connCount := 0
	for _, s := range sources {
		if s.Connected {
			connCount++
		}
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(10).
		Gap(8).
		Radius(14).
		Background(t.Surface).
		Border(1, t.Border)

	card.Children(func() {
		// Header row with connected count
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Integrations").FontSize(12).FontWeight(700).TextColor(t.Text)
			ui.Row(c).AlignItems(ui.Center).Gap(5).Padding(2, 7).Radius(999).Background(t.SurfaceSubtle).Border(1, t.Border).Children(func() {
				dotCol := t.TextSubtle
				if connCount > 0 {
					dotCol = t.Success
				}
				ui.Box(c).Size(6, 6).Radius(999).Background(dotCol)
				ui.Text(c, fmt.Sprintf("%d/%d connected", connCount, len(sources))).
					FontSize(10).
					FontWeight(600).
					TextColor(t.TextMuted)
			})
		})

		// 2-column grid of tiles
		for rowStart := 0; rowStart < len(sources); rowStart += 2 {
			rs := rowStart
			ui.Row(c).Key(rs).WidthPercent(100).Gap(7).Children(func() {
				for col := 0; col < 2; col++ {
					idx := rs + col
					if idx >= len(sources) {
						ui.Box(c).Grow(1)
						continue
					}
					src := &sources[idx]
					brand := colorOr(src.Color, t.Accent)

					tile := ui.Column(c).
						Key(src.ID).
						Grow(1).
						Padding(8, 9).
						Gap(7).
						Radius(10).
						Cursor(ui.CursorPointer)

					target := float32(0)
					if src.Connected {
						target = 1
					}
					p := clamp32(tile.Spring("conn_p", target, ui.SpringOptions{
						Duration: 280 * time.Millisecond,
						Bounce:   0.18,
					}), 0, 1)

					tileBg := lerpColor(t.SurfaceSubtle, lerpColor(t.SurfaceSubtle, brand, 0.08), p)
					tileBorder := lerpColor(t.Border, brand.Alpha(0.55), p)
					tile.Background(tileBg).Border(1, tileBorder)

					tile.Children(func() {
						ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
							ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
								ui.Box(c).
									Size(20, 20).
									Radius(6).
									Center().
									Background(brand.Alpha(0.18)).
									Children(func() {
										ui.Text(c, src.Badge).FontSize(9.5).FontWeight(700).TextColor(brand)
									})
								ui.Text(c, src.Name).FontSize(11.5).FontWeight(600).TextColor(t.Text).SingleLine()
							})
							ui.Box(c).
								Size(6.5, 6.5).
								Radius(999).
								Background(lerpColor(t.BorderStrong, t.Success, p))
						})

						btnLabel := "Connect"
						btnBg := t.Surface
						btnBorder := t.BorderStrong
						btnTxt := t.Text
						if src.Connected {
							btnLabel = "Connected"
							btnBg = t.Success.Alpha(0.16)
							btnBorder = t.Success.Alpha(0.50)
							btnTxt = t.Success
						}
						ui.Row(c).
							WidthPercent(100).
							Height(22).
							Center().
							Radius(6).
							Background(btnBg).
							Border(1, btnBorder).
							Children(func() {
								ui.Text(c, btnLabel).FontSize(10).FontWeight(600).TextColor(btnTxt)
							})
					})

					if tile.Clicked() {
						src.Connected = !src.Connected
						if o.OnToggle != nil {
							o.OnToggle(src.ID, src.Connected)
						}
					}
				}
			})
		}
	})

	return card
}

// -----------------------------------------------------------------------------
// 11. access-request ("Access request")
// -----------------------------------------------------------------------------

// AccessRequestPhase enumerates the finite state machine states of AccessRequest.
type AccessRequestPhase int

const (
	AccessRequestPending AccessRequestPhase = iota
	AccessRequestApproved
	AccessRequestDenied
)

// AccessRequestOptions configures AccessRequest.
type AccessRequestOptions struct {
	Requester     string
	Initials      string
	Resource      string
	Role          string
	Width         float32
	Accent        ui.Color
	OnPhaseChange func(AccessRequestPhase)
}

// AccessRequest renders a compact permission request card with requester avatar,
// role badge, Approve/Deny state machine transitions, and inline Undo.
func AccessRequest(c *ui.Context, phase *AccessRequestPhase, opts ...AccessRequestOptions) ui.Element {
	var o AccessRequestOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Requester == "" {
		o.Requester = "Ava Kim"
	}
	if o.Initials == "" {
		o.Initials = "AK"
	}
	if o.Resource == "" {
		o.Resource = "Production DB"
	}
	if o.Role == "" {
		o.Role = "Editor"
	}
	if o.Width <= 0 {
		o.Width = 282
	}
	accent := colorOr(o.Accent, t.Accent)

	curPhase := AccessRequestPending
	if phase != nil {
		curPhase = *phase
	}

	setPhase := func(next AccessRequestPhase) {
		if phase != nil {
			*phase = next
		}
		if o.OnPhaseChange != nil {
			o.OnPhaseChange(next)
		}
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(12).
		Gap(10).
		Radius(14).
		Background(t.Surface)

	resolvedTarget := float32(0)
	if curPhase != AccessRequestPending {
		resolvedTarget = 1
	}
	resP := clamp32(card.Spring("ar_res_p", resolvedTarget, ui.SpringOptions{
		Duration: 300 * time.Millisecond,
		Bounce:   0.18,
	}), 0, 1)

	statusCol := accent
	if curPhase == AccessRequestApproved {
		statusCol = t.Success
	} else if curPhase == AccessRequestDenied {
		statusCol = t.Danger
	}
	card.Border(1, lerpColor(t.Border, statusCol.Alpha(0.55), resP))

	card.Children(func() {
		// Requester row
		ui.Row(c).
			WidthPercent(100).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Gap(8).
			Children(func() {
				ui.Row(c).AlignItems(ui.Center).Gap(9).Children(func() {
					ui.Box(c).
						Size(34, 34).
						Radius(17).
						Center().
						Background(lerpColor(accent.Alpha(0.18), statusCol.Alpha(0.20), resP)).
						Border(1, lerpColor(accent.Alpha(0.45), statusCol.Alpha(0.55), resP)).
						Children(func() {
							ui.Text(c, o.Initials).FontSize(12).FontWeight(700).TextColor(lerpColor(accent, statusCol, resP))
						})
					ui.Column(c).Gap(2).Children(func() {
						ui.Text(c, o.Requester).FontSize(12.5).FontWeight(600).TextColor(t.Text).SingleLine()
						ui.Text(c, "Requests "+o.Resource).FontSize(10.5).TextColor(t.TextMuted).SingleLine()
					})
				})

				ui.Box(c).
					Padding(3, 8).
					Radius(999).
					Background(t.AccentSoft).
					Children(func() {
						ui.Text(c, o.Role).FontSize(10).FontWeight(600).TextColor(accent)
					})
			})

		// Hairline divider
		ui.Box(c).WidthPercent(100).Height(1).BorderWidth(0, 0, 1, 0).BorderColor(t.Border)

		// Action footer driven by state machine
		switch curPhase {
		case AccessRequestPending:
			ui.Row(c).Key("ar_pending").WidthPercent(100).Gap(8).Children(func() {
				denyBtn := ui.Row(c).
					Key("deny").
					Grow(1).
					Height(32).
					Center().
					Radius(9).
					Background(t.SurfaceSubtle).
					Border(1, t.BorderStrong).
					Cursor(ui.CursorPointer)
				if denyBtn.Hovered() {
					denyBtn.Background(t.Danger.Alpha(0.12)).Border(1, t.Danger.Alpha(0.45))
				}
				denyBtn.Children(func() {
					ui.Text(c, "Deny").FontSize(12).FontWeight(600).TextColor(t.TextMuted)
				})
				if denyBtn.Clicked() {
					setPhase(AccessRequestDenied)
				}

				approveBtn := ui.Row(c).
					Key("approve").
					Grow(1).
					Height(32).
					Center().
					Radius(9).
					Background(accent).
					Cursor(ui.CursorPointer)
				approveBtn.Children(func() {
					ui.Text(c, "Approve").FontSize(12).FontWeight(600).TextColor(t.AccentText)
				})
				if approveBtn.Clicked() {
					setPhase(AccessRequestApproved)
				}
			})

		case AccessRequestApproved, AccessRequestDenied:
			bannerText := "Approved · " + o.Role + " granted"
			if curPhase == AccessRequestDenied {
				bannerText = "Denied · Request declined"
			}
			ui.Row(c).
				Key("ar_resolved").
				WidthPercent(100).
				Height(32).
				Padding(0, 10).
				Radius(9).
				Justify(ui.SpaceBetween).
				AlignItems(ui.Center).
				Background(statusCol.Alpha(0.14)).
				Border(1, statusCol.Alpha(0.45)).
				Children(func() {
					ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
						ui.Box(c).Size(7, 7).Radius(999).Background(statusCol)
						ui.Text(c, bannerText).FontSize(11).FontWeight(600).TextColor(statusCol).SingleLine()
					})

					undoBtn := ui.Box(c).
						Key("undo").
						Padding(3, 7).
						Radius(6).
						Background(t.Surface).
						Border(1, t.BorderStrong).
						Cursor(ui.CursorPointer)
					undoBtn.Children(func() {
						ui.Text(c, "Undo").FontSize(10.5).FontWeight(600).TextColor(t.Text)
					})
					if undoBtn.Clicked() {
						setPhase(AccessRequestPending)
					}
				})
		}
	})

	return card
}

// -----------------------------------------------------------------------------
// 12. billing-card ("Billing card")
// -----------------------------------------------------------------------------

// BillingCycle enumerates monthly vs yearly billing intervals.
type BillingCycle int

const (
	BillingMonthly BillingCycle = iota
	BillingYearly
)

// BillingCardOptions configures BillingCard.
type BillingCardOptions struct {
	PlanName     string
	TierBadge    string
	MonthlyPrice float32
	YearlyPrice  float32
	Width        float32
	Accent       ui.Color
	OnUpgrade    func(cycle BillingCycle)
}

// BillingCard renders an interactive subscription plan card with Monthly/Yearly
// toggle, spring-animated price update, plan tier badge, and upgrade state.
func BillingCard(c *ui.Context, cycle *BillingCycle, upgraded *bool, opts ...BillingCardOptions) ui.Element {
	var o BillingCardOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.PlanName == "" {
		o.PlanName = "Pro Studio"
	}
	if o.TierBadge == "" {
		o.TierBadge = "PRO"
	}
	if o.MonthlyPrice <= 0 {
		o.MonthlyPrice = 24
	}
	if o.YearlyPrice <= 0 {
		o.YearlyPrice = 19
	}
	if o.Width <= 0 {
		o.Width = 278
	}
	accent := colorOr(o.Accent, t.Accent)

	curCycle := BillingMonthly
	if cycle != nil {
		curCycle = *cycle
	}
	isUpgraded := false
	if upgraded != nil {
		isUpgraded = *upgraded
	}

	targetPrice := o.MonthlyPrice
	cycleTarget := float32(0)
	if curCycle == BillingYearly {
		targetPrice = o.YearlyPrice
		cycleTarget = 1
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(12).
		Gap(9).
		Radius(14).
		Background(t.Surface)

	upgTarget := float32(0)
	if isUpgraded {
		upgTarget = 1
	}
	upgP := clamp32(card.Spring("bill_upg_p", upgTarget, ui.SpringOptions{
		Duration: 280 * time.Millisecond,
		Bounce:   0.18,
	}), 0, 1)

	animPrice := card.Spring("bill_price", targetPrice, ui.SpringOptions{
		Duration: 340 * time.Millisecond,
		Bounce:   0.14,
	})
	cycleP := clamp32(card.Spring("bill_cycle_p", cycleTarget, ui.SpringOptions{
		Duration: 300 * time.Millisecond,
		Bounce:   0.20,
	}), 0, 1)

	card.Border(1, lerpColor(t.Border, t.Success.Alpha(0.55), upgP))

	card.Children(func() {
		// Top row: Plan name + Tier badge + Monthly/Yearly toggle
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
				ui.Text(c, o.PlanName).FontSize(12.5).FontWeight(700).TextColor(t.Text)
				badgeLabel := o.TierBadge
				badgeBg := t.AccentSoft
				badgeTxt := accent
				if isUpgraded {
					badgeLabel = "ACTIVE"
					badgeBg = t.Success.Alpha(0.16)
					badgeTxt = t.Success
				}
				ui.Box(c).Padding(2, 6).Radius(999).Background(badgeBg).Children(func() {
					ui.Text(c, badgeLabel).FontSize(9.5).FontWeight(700).TextColor(badgeTxt)
				})
			})

			// Segmented Monthly / Yearly switch
			const segW float32 = 108
			const thumbW float32 = 51
			seg := ui.Box(c).
				Size(segW, 24).
				Radius(999).
				Background(t.SurfaceSubtle).
				Border(1, t.Border)
			seg.Children(func() {
				// Sliding thumb
				ui.Box(c).
					Absolute().
					Left(lerp32(2, segW-thumbW-2, cycleP)).
					Top(2).
					Size(thumbW, 18).
					Radius(999).
					Background(t.Surface).
					Border(1, t.BorderStrong)

				ui.Row(c).WidthPercent(100).HeightPercent(100).Children(func() {
					mBtn := ui.Box(c).Key("bill_m").Grow(1).HeightPercent(100).Center().Cursor(ui.CursorPointer)
					mBtn.Children(func() {
						ui.Text(c, "Mo").FontSize(10).FontWeight(600).TextColor(lerpColor(t.Text, t.TextSubtle, cycleP))
					})
					if mBtn.Clicked() && cycle != nil {
						*cycle = BillingMonthly
					}

					yBtn := ui.Box(c).Key("bill_y").Grow(1).HeightPercent(100).Center().Cursor(ui.CursorPointer)
					yBtn.Children(func() {
						ui.Text(c, "Yr").FontSize(10).FontWeight(600).TextColor(lerpColor(t.TextSubtle, t.Text, cycleP))
					})
					if yBtn.Clicked() && cycle != nil {
						*cycle = BillingYearly
					}
				})
			})
		})

		// Animated price display + savings badge
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Row(c).AlignItems(ui.End).Gap(3).Children(func() {
				shown := int(animPrice + 0.5)
				ui.Text(c, fmt.Sprintf("$%d", shown)).FontSize(26).FontWeight(700).TextColor(t.Text)
				ui.Box(c).Padding(0, 0, 4, 0).Children(func() {
					ui.Text(c, "/mo").FontSize(11).TextColor(t.TextSubtle)
				})
			})
			if cycleP > 0.05 {
				ui.Box(c).
					Padding(2, 7).
					Radius(999).
					Background(t.Success.Alpha(0.15)).
					Border(1, t.Success.Alpha(0.45)).
					Opacity(cycleP).
					Scale(lerp32(0.85, 1, cycleP), lerp32(0.85, 1, cycleP)).
					Children(func() {
						ui.Text(c, "Save 20%").FontSize(10).FontWeight(600).TextColor(t.Success)
					})
			}
		})

		// Compact plan summary
		ui.Text(c, "Unlimited blocks · GPU shaders · Team sync").
			FontSize(10.5).
			TextColor(t.TextMuted).
			SingleLine()

		// Upgrade CTA button
		cta := ui.Row(c).
			Key("bill_cta").
			WidthPercent(100).
			Height(32).
			Center().
			Radius(9).
			Cursor(ui.CursorPointer)

		pressTarget := float32(1)
		if cta.Pressed() {
			pressTarget = 0.97
		}
		sc := cta.Spring("bill_cta_sc", pressTarget, ui.SpringOptions{
			Duration: 160 * time.Millisecond,
			Bounce:   0.20,
		})

		if isUpgraded {
			cta.Background(t.Success.Alpha(0.16)).Border(1, t.Success.Alpha(0.55))
		} else {
			cta.Background(t.Text)
		}
		cta.Scale(sc, sc).Children(func() {
			if isUpgraded {
				ui.Text(c, "Active plan · Manage").FontSize(11.5).FontWeight(600).TextColor(t.Success)
			} else {
				ui.Text(c, "Upgrade to "+o.PlanName).FontSize(11.5).FontWeight(600).TextColor(t.Background)
			}
		})
		if cta.Clicked() {
			if upgraded != nil {
				*upgraded = !*upgraded
			}
			if !isUpgraded && o.OnUpgrade != nil {
				o.OnUpgrade(curCycle)
			}
		}
	})

	return card
}

// -----------------------------------------------------------------------------
// Catalog Demo Registrations
// -----------------------------------------------------------------------------

type faqDemoState struct {
	open int
}

type treeDemoState struct {
	expanded map[string]bool
	selected string
}

type sidebarDemoState struct {
	collapsed bool
	active    string
}

type dropdownDemoState struct {
	selected int
	open     bool
}

type radialDemoState struct {
	open     bool
	selected int
}

type createMenuDemoState struct {
	open        bool
	lastCreated string
}

type selectionActionsDemoState struct {
	selected   map[string]bool
	lastAction string
}

type connectSourcesDemoState struct {
	sources []IntegrationSource
}

type accessRequestDemoState struct {
	phase AccessRequestPhase
}

type billingCardDemoState struct {
	cycle    BillingCycle
	upgraded bool
}

func init() {

	// 1. faq-frames

	// 2. faq-accordion

	// 3. pill-accordion

	// 9. selection-actions
	selectableItems := []SelectableItem{
		{ID: "1", Name: "Design tokens.json", Meta: "12 KB"},
		{ID: "2", Name: "Spring curves.md", Meta: "4 KB"},
		{ID: "3", Name: "Release notes.txt", Meta: "8 KB"},
		{ID: "4", Name: "Shader pipeline.wgsl", Meta: "19 KB"},
	}
	registerDemo(
		"selection-actions",
		"Multi-select list paired with a floating bulk action bar that springs in on selection.",
		func(c *ui.Context) {
			holder := ui.Column(c).AlignItems(ui.Center).Gap(4)
			st := ui.Local(holder, "demo_selection_actions", func() selectionActionsDemoState {
				return selectionActionsDemoState{
					selected: map[string]bool{"1": true, "2": true},
				}
			})
			holder.Children(func() {
				SelectionActions(c, selectableItems, st.selected, SelectionActionsOptions{
					OnAction: func(action string, ids []string) {
						st.lastAction = fmt.Sprintf("%s applied to %d items", action, len(ids))
						for _, id := range ids {
							delete(st.selected, id)
						}
					},
				})
				if st.lastAction != "" && len(st.selected) == 0 {
					t := ThemeFrom(c)
					ui.Text(c, st.lastAction).FontSize(10.5).TextColor(t.TextMuted)
				}
			})
		},
	)

	// 10. connect-sources
	registerDemo(
		"connect-sources",
		"Integration source tiles with live Connect/Connected state, status dots, and active count.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_connect_sources", func() connectSourcesDemoState {
				return connectSourcesDemoState{
					sources: []IntegrationSource{
						{ID: "github", Name: "GitHub", Badge: "GH", Color: ui.Hex("#6366f1"), Connected: true},
						{ID: "linear", Name: "Linear", Badge: "LN", Color: ui.Hex("#5e6ad2"), Connected: true},
						{ID: "slack", Name: "Slack", Badge: "SL", Color: ui.Hex("#e01e5a"), Connected: false},
						{ID: "figma", Name: "Figma", Badge: "FG", Color: ui.Hex("#f24e1e"), Connected: false},
					},
				}
			})
			holder.Children(func() {
				ConnectSources(c, st.sources)
			})
		},
	)

	// 11. access-request
	registerDemo(
		"access-request",
		"Compact permission request card with Approve/Deny state machine transitions and undo.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_access_request", func() accessRequestDemoState {
				return accessRequestDemoState{phase: AccessRequestPending}
			})
			holder.Children(func() {
				AccessRequest(c, &st.phase)
			})
		},
	)

	// 12. billing-card
	registerDemo(
		"billing-card",
		"Subscription plan card with Monthly/Yearly spring toggle, animated price, and upgrade state.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_billing_card", func() billingCardDemoState {
				return billingCardDemoState{cycle: BillingMonthly, upgraded: false}
			})
			holder.Children(func() {
				BillingCard(c, &st.cycle, &st.upgraded)
			})
		},
	)
}
