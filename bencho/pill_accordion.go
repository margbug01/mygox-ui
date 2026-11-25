package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// PillAccordionItem is one item: a title, and the details it opens to.
type PillAccordionItem struct {
	Title, Details string
}

// PillAccordionOptions configures PillAccordion.
type PillAccordionOptions struct {
	Corner float32 // corner radius, 0 to 26; 26 by default
}

type pillAccordionState struct {
	open             int
	titleW, detailsH map[int]float32 // measured
}

// PillAccordion stacks items as pills. Opening one grows it into a card with
// its details while the one open before folds back into a pill.
func PillAccordion(c *ui.Context, items []PillAccordionItem, opts ...PillAccordionOptions) ui.Element {
	var o PillAccordionOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 26
	}
	t := ThemeFrom(c)
	col := ui.Column(c).Width(300).Padding(10).Gap(12)
	st := ui.Local(col, "pill_accordion", func() pillAccordionState {
		return pillAccordionState{titleW: map[int]float32{}, detailsH: map[int]float32{}}
	})
	col.Children(func() {
		for i, it := range items {
			open := st.open == i
			pill := ui.Box(c).Key(i).Radius(o.Corner).Background(t.Surface).Clip().Cursor(ui.CursorPointer)
			k := clamp32(pill.Spring("open", b2f(open), ui.SpringOptions{Duration: 460 * time.Millisecond, Bounce: 0.18}), 0, 1.06)
			closedW := 12 + 32 + 12 + st.titleW[i] + 18
			pill.Size(lerp32(closedW, 280, min(k, 1)), 52+st.detailsH[i]*k)
			pill.OnClick(func() {
				if st.open == i {
					st.open = -1
				} else {
					st.open = i
				}
			})
			pill.Children(func() {
				ui.Row(c).Absolute().Left(12).Top(10).Gap(12).AlignItems(ui.Center).Children(func() {
					knob := ui.Box(c).Size(32, 32).Radius(16).Center()
					knob.Background(lerpColor(t.SurfaceSubtle, t.Surface, min(k, 1)))
					if k > 0.5 {
						knob.Border(1, t.Border)
					}
					knob.Children(func() {
						ui.Icon(c, icons.Plus()).Size(16, 16).TextColor(t.Text).Rotate(45 * k)
					})
					title := ui.Text(c, it.Title).FontSize(16).FontWeight(500).TextColor(t.Text).SingleLine()
					if w := title.LayoutBounds().W; w > 0 {
						st.titleW[i] = w
					}
				})
				if k > 0.01 || st.detailsH[i] == 0 {
					d := ui.Text(c, it.Details).FontSize(14.5).LineHeight(1.4).TextColor(t.TextSubtle).
						Absolute().Left(20).Top(56).Width(240).Padding(0, 0, 20, 0).Opacity(clamp32(k*1.4-0.4, 0, 1))
					if h := d.LayoutBounds().H; h > 0 {
						st.detailsH[i] = h + 4
					}
				}
			})
		}
	})
	return col
}

var pillAccordionDemoItems = []PillAccordionItem{
	{Title: "Pick a block", Details: "Browse the wall and open the one that fits what you are making."},
	{Title: "Tune the knobs", Details: "Every block has a few settings; try them until it feels right."},
	{Title: "Copy the code", Details: "Take the source straight into your project and make it yours."},
	{Title: "Ship it", Details: "Drop it into your app and let people press, drag and play."},
}

func init() {
	registerDemo(
		"pill-accordion",
		"Each item starts as a pill in a column; open one and it grows into a card with its details while the open one folds back.",
		func(c *ui.Context) { PillAccordion(c, pillAccordionDemoItems) },
	)
}
