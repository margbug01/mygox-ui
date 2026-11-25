package bencho

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// AccessRequester is one person waiting to join.
type AccessRequester struct {
	Handle   string // shown with an @
	Role     string // what they asked for, under the handle in the list
	From, To ui.Color
}

// AccessRequestOptions configures AccessRequest.
type AccessRequestOptions struct {
	Speed    float32 // how quickly the pill opens into the list, 0 to 100; 90 by default
	Bounce   float32 // how much the card swells as it opens and gathers as it folds, 0 to 100; 40 by default
	Corner   float32 // the open card's corner, 8 to 28; 20 by default
	OnDecide func(handle string, approved bool)
}

type accessRequestState struct {
	open    bool
	decided map[string]bool // handle → approved
	pressed time.Time
}

const (
	accessRequestW    = 300
	accessRequestPill = 56
	accessRequestRow  = 52
)

// AccessRequestDefaultPeople are the three people bencho.dev shows waiting.
func AccessRequestDefaultPeople() []AccessRequester {
	return []AccessRequester{
		{Handle: "hanaito", Role: "Editor", From: ui.Hex("#36b0ff"), To: ui.Hex("#5b45dd")},
		{Handle: "leoroux", Role: "Viewer", From: ui.Hex("#74f25a"), To: ui.Hex("#22d3bf")},
		{Handle: "mira.k", Role: "Commenter", From: ui.Hex("#ee3c88"), To: ui.Hex("#f8603f")},
	}
}

// accessRequestAvatar is a gradient disc with a ring in the card's colour.
func accessRequestAvatar(c *ui.Context, t Theme, key any, p AccessRequester, size float32) ui.Element {
	return ui.Box(c).Key(key).Size(size, size).Draw(func(pt *ui.Painter, r ui.Rect) {
		pt.Fill(ui.Rect{X: r.X - 2, Y: r.Y - 2, W: r.W + 4, H: r.H + 4}, t.Surface, r.W/2+2)
		pt.FillGradient(r, ui.LinearGradient{From: p.From, To: p.To, Angle: 150}, r.W/2)
	})
}

// AccessRequest starts as a pill of people waiting to join. Pressing it
// opens it into a list with Approve and Reject for each; a decision folds
// that person's row away, and the pill says so once everyone is answered.
func AccessRequest(c *ui.Context, people []AccessRequester, opts ...AccessRequestOptions) ui.Element {
	var o AccessRequestOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Speed == 0 {
		o.Speed = 90
	}
	if o.Bounce == 0 {
		o.Bounce = 40
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	t := ThemeFrom(c)
	holder := ui.Box(c).Width(accessRequestW)
	st := ui.Local(holder, "access_request", func() accessRequestState { return accessRequestState{decided: map[string]bool{}} })
	var waiting []AccessRequester
	for _, p := range people {
		if _, ok := st.decided[p.Handle]; !ok {
			waiting = append(waiting, p)
		}
	}
	if len(waiting) == 0 {
		st.open = false
	}

	dur := time.Duration(lerp32(900, 300, clamp32(o.Speed, 0, 100)/100)) * time.Millisecond
	bounce := float64(clamp32(o.Bounce, 0, 100)) / 100
	k := holder.Spring("open", b2f(st.open), ui.SpringOptions{Duration: dur, Bounce: 0.45 * bounce})
	ko := clamp32(k, 0, 1)
	// The swell: a press squeezes the pill; it bulges a touch as it opens.
	squeeze := holder.Spring("press", b2f(c.Now().Sub(st.pressed) < 120*time.Millisecond), ui.SpringOptions{Duration: 160 * time.Millisecond})
	if c.Now().Sub(st.pressed) < 200*time.Millisecond {
		c.AnimationFrame()
	}
	scale := 1 - 0.03*squeeze*float32(bounce)*2

	listH := float32(0)
	for _, p := range people {
		_, gone := st.decided[p.Handle]
		listH += accessRequestRow * clamp32(holder.Spring("row/"+p.Handle, b2f(!gone), ui.SpringOptions{Duration: 320 * time.Millisecond}), 0, 1)
	}
	h := accessRequestPill + max(k, 0)*(listH+8)
	corner := lerp32(accessRequestPill/2, clamp32(o.Corner, 8, 28), ko)
	holder.Height(accessRequestPill)

	holder.Children(func() {
		card := ui.Column(c).Absolute().Left(0).Top(0).Size(accessRequestW, h).Radius(corner).Background(t.Surface).Clip().
			Scale(scale, scale)
		if k > 0.02 {
			card.Shadow(0, 10, 30, 0, ui.RGBA(0, 0, 0, 0.08*ko))
		}
		card.Children(func() {
			head := ui.Row(c).Key("head").Height(accessRequestPill).Padding(0, 20, 0, 12).AlignItems(ui.Center).Cursor(ui.CursorPointer).
				OnClick(func() {
					if len(waiting) > 0 {
						st.open = !st.open
						st.pressed = c.Now()
					}
				})
			head.Children(func() {
				if len(waiting) == 0 {
					ui.Box(c).Size(32, 32).Radius(16).Background(t.SurfaceSubtle).Center().Margin(0, 12, 0, 0).Children(func() {
						ui.Icon(c, icons.Check()).Size(16, 16).TextColor(t.Text)
					})
					ui.Column(c).Grow(1).Children(func() {
						ui.Text(c, "All caught up").FontSize(14).TextColor(t.Text)
						ui.Text(c, "No one is waiting").FontSize(14).TextColor(t.TextSubtle)
					})
					return
				}
				ui.Box(c).Size(32+18*float32(len(waiting)-1), 32).Margin(0, 12, 0, 0).Children(func() {
					for i := len(waiting) - 1; i >= 0; i-- {
						accessRequestAvatar(c, t, waiting[i].Handle, waiting[i], 32).Absolute().Left(18 * float32(i)).Top(0)
					}
				})
				ui.Column(c).Grow(1).Children(func() {
					ui.Row(c).Children(func() {
						ui.Text(c, "@"+waiting[0].Handle).FontSize(14).TextColor(t.Text)
						if n := len(waiting) - 1; n > 0 {
							others := fmt.Sprintf("%d others", n)
							if n == 1 {
								others = "1 other"
							}
							ui.Text(c, "and").FontSize(14).TextColor(t.TextSubtle).Margin(0, 4)
							ui.Text(c, others).FontSize(14).TextColor(t.Text)
						}
					})
					ui.Text(c, "requested access").FontSize(14).TextColor(t.TextSubtle)
				})
				ui.Icon(c, icons.ChevronDown()).Size(16, 16).TextColor(t.TextSubtle).Rotate(180 * ko)
			})

			if k <= 0.02 {
				return
			}
			ui.Box(c).Key("rule").Height(1).Margin(0, 16).Background(t.Border).Opacity(ko)
			for i, p := range people {
				_, gone := st.decided[p.Handle]
				rk := clamp32(holder.Spring("row/"+p.Handle, b2f(!gone), ui.SpringOptions{Duration: 320 * time.Millisecond}), 0, 1)
				if rk < 0.01 {
					continue
				}
				// Rows come in one after another as the card opens.
				in := clamp32(k*1.6-float32(i)*0.2, 0, 1)
				ui.Row(c).Key(p.Handle).Height(accessRequestRow*rk).Padding(0, 12, 0, 12).AlignItems(ui.Center).
					Opacity(in*rk).Translate(0, (1-in)*-8).Children(func() {
					accessRequestAvatar(c, t, "av", p, 32).Margin(0, 10, 0, 0)
					ui.Column(c).Grow(1).MinWidth(0).Children(func() {
						ui.Text(c, "@"+p.Handle).FontSize(13.5).TextColor(t.Text).SingleLine()
						ui.Text(c, p.Role).FontSize(12.5).TextColor(t.TextSubtle).SingleLine()
					})
					decide := func(approved bool) func() {
						return func() {
							st.decided[p.Handle] = approved
							if o.OnDecide != nil {
								o.OnDecide(p.Handle, approved)
							}
						}
					}
					button := func(label string, bg, fg ui.Color, fn func()) {
						b := ui.Box(c).Key(label).Height(30).Padding(0, 12).Radius(15).Center().Cursor(ui.CursorPointer).OnClick(fn)
						hv := b.Spring("hover", b2f(b.Hovered()), ui.SpringOptions{Duration: 140 * time.Millisecond})
						b.Background(bg).Opacity(1 - 0.12*hv).Children(func() {
							ui.Text(c, label).FontSize(13).FontWeight(500).TextColor(fg)
						})
					}
					button("Reject", t.SurfaceSubtle, t.TextMuted, decide(false))
					ui.Box(c).Width(6)
					button("Approve", t.Text, t.Surface, decide(true))
				})
			}
		})
	})
	return holder
}

func init() {
	registerDemo(
		"access-request",
		"A pill of people waiting to join that opens into a list with Approve and Reject for each.",
		func(c *ui.Context) { AccessRequest(c, AccessRequestDefaultPeople()) },
	)
}
