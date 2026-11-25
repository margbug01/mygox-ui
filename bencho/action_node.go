package bencho

import (
	"fmt"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// ActionNodeOptions configures ActionNode.
type ActionNodeOptions struct {
	Bounce  float32 // 0 to 40; 20 by default
	Stagger float32 // all at once (0) or one after another (100); 55 by default
	Count   int     // buttons, 2 to 4; 4 by default
	Corner  float32 // 0 to 32px; 20 by default
	// OnAction is called with the button pressed: "Link", "Add", "Run" or "More".
	OnAction func(name string)
}

const (
	actionNodeW     = 272
	actionNodeH     = 286
	actionNodeCardY = 52
	actionNodeCardH = 182
	actionNodeBtn   = 36
)

type actionNodeState struct {
	on bool      // the buttons are out
	at time.Time // when they were last sent out or folded
}

// actionNodeSpots are where the buttons come out to, around the card's top
// right corner: two along the top, two down the side.
var actionNodeSpots = [4][2]float32{{220, 22}, {268, 22}, {302, 58}, {302, 106}}

var actionNodeGlyphs = sync.OnceValue(func() [2]*ui.SVG {
	parse := func(body string) *ui.SVG {
		return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">` + body + `</svg>`))
	}
	return [2]*ui.SVG{
		// Lucide's git-commit-horizontal and link-2.
		parse(`<circle cx="12" cy="12" r="3"/><line x1="3" x2="9" y1="12" y2="12"/><line x1="15" x2="21" y1="12" y2="12"/>`),
		parse(`<path d="M9 17H7A5 5 0 0 1 7 7h2"/><path d="M15 7h2a5 5 0 1 1 0 10h-2"/><line x1="8" x2="16" y1="12" y2="12"/>`),
	}
})

// ActionNode is a workflow step's card: an icon, its name, what it does
// and who runs it. Its buttons sit folded into the top right corner and
// spring out round it while the card is pointed at, all at once or one
// after another.
func ActionNode(c *ui.Context, opts ...ActionNodeOptions) ui.Element {
	var o ActionNodeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	or := func(v, def, lo, hi float32) float32 {
		if v == 0 {
			return def
		}
		return clamp32(v, lo, hi)
	}
	bounce := or(o.Bounce, 20, 0, 40) / 100
	stagger := or(o.Stagger, 55, 0, 100) / 100
	corner := or(o.Corner, 20, 0, 32)
	if o.Corner < 0 {
		corner = 0
	}
	if o.Bounce < 0 {
		bounce = 0
	}
	if o.Stagger < 0 {
		stagger = 0
	}
	count := 4
	if o.Count != 0 {
		count = int(clamp32(float32(o.Count), 2, 4))
	}
	t := ThemeFrom(c)
	glyphs := actionNodeGlyphs()

	root := ui.Box(c).Size(actionNodeW, actionNodeH)
	st := ui.Local(root, "action_node", func() actionNodeState { return actionNodeState{} })
	now := c.Now()

	root.Children(func() {
		card := ui.Column(c).Key("card").Absolute().Left(0).Top(actionNodeCardY).Size(actionNodeW, actionNodeCardH).
			Radius(corner).Background(t.Surface).Padding(20, 14, 0, 20)
		// Open while the card or a button is pointed at.
		hover := card.Hovered()
		card.Children(func() {
			ui.Row(c).Height(30).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(30, 30).Radius(8).Background(t.SurfaceSubtle).Center().Children(func() {
					ui.Icon(c, glyphs[0]).Size(16, 16).TextColor(t.Text)
				})
				ui.Text(c, "Node 07").FontSize(16.5).TextColor(t.Text).Margin(0, 0, 0, 9)
			})
			ui.Text(c, "Reads the thread, writes a summary, and posts it back to the channel it came from.").
				FontSize(13.5).LineHeight(1.42).TextColor(t.TextMuted).Margin(12, 0, 0, 0)
			ui.Row(c).Height(28).AlignItems(ui.Center).Margin(14, 0, 0, 0).Children(func() {
				grads := [3][2]ui.Color{
					{ui.Hex("#38a8fc"), ui.Hex("#6a4cff")},
					{ui.Hex("#f56f4c"), ui.Hex("#e8475e")},
					{ui.Hex("#ffa7e8"), ui.Hex("#c445dc")},
				}
				for i, g := range grads {
					ui.Box(c).Key(fmt.Sprint("avatar", i)).Size(27, 27).Radius(13.5).Background(t.Surface).Margin(0, 0, 0, float32(min(i, 1))*-10).
						Draw(func(p *ui.Painter, r ui.Rect) {
							in := ui.Rect{X: r.X + 2, Y: r.Y + 2, W: r.W - 4, H: r.H - 4}
							p.FillGradient(in, ui.LinearGradient{From: g[0], To: g[1], Angle: 135}, in.W/2)
						})
				}
				ui.Text(c, "+3").FontSize(12.5).TextColor(t.TextMuted).Margin(0, 0, 0, 9)
			})
		})

		names := [4]string{"Link", "Add", "Run", "More"}
		svgs := [4]*ui.SVG{glyphs[1], icons.Plus(), icons.Play(), icons.Ellipsis()}
		for i := 0; i < count; i++ {
			// Each button waits its turn, then springs from the corner to its spot.
			delay := time.Duration(float32(i) * stagger * float32(90*time.Millisecond))
			if !st.on {
				delay = time.Duration(float32(count-1-i) * stagger * float32(40*time.Millisecond))
			}
			out := b2f(st.on)
			if now.Sub(st.at) < delay {
				out = b2f(!st.on)
			}
			if now.Sub(st.at) < delay+50*time.Millisecond {
				c.AnimationFrame()
			}
			btn := ui.Box(c).Key(names[i]).Absolute().Size(actionNodeBtn, actionNodeBtn).Radius(actionNodeBtn / 2).
				Background(t.Surface).Center().Cursor(ui.CursorPointer)
			k := btn.Spring("out", out, ui.SpringOptions{Duration: 450 * time.Millisecond, Bounce: float64(bounce)})
			if btn.Hovered() {
				hover = true
			}
			name := names[i]
			btn.OnClick(func() {
				if o.OnAction != nil {
					o.OnAction(name)
				}
			})
			sx, sy := actionNodeSpots[i][0], actionNodeSpots[i][1]
			x := lerp32(actionNodeW-actionNodeBtn/2-6, sx, k)
			y := lerp32(actionNodeCardY+actionNodeBtn/2+6, sy, k)
			s := clamp32(k, 0, 1.3)
			press := 1 - 0.08*b2f(btn.Pressed())
			btn.Left(x-actionNodeBtn/2).Top(y-actionNodeBtn/2).Scale(s*press, s*press).Opacity(clamp32(k*1.5, 0, 1)).Children(func() {
				ui.Icon(c, svgs[i]).Size(16, 16).TextColor(t.Text)
			})
		}
		if hover != st.on {
			st.on, st.at = hover, now
			c.AnimationFrame()
		}
	})
	return root
}

func init() {
	registerDemo(
		"action-node",
		"A workflow step card whose buttons sit folded into the corner and spring out on hover, all at once or one after another.",
		func(c *ui.Context) { ActionNode(c) },
	)
}
