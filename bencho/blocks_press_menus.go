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
