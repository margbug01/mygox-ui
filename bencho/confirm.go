package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// InlineConfirmOptions configures InlineConfirm.
type InlineConfirmOptions struct {
	Corner   float32 // 0 to 22px; 22 by default
	OnDelete func()
	OnUndo   func()
}

// InlineConfirmStep is where the button is in deleting.
type InlineConfirmStep int

const (
	InlineConfirmIdle    InlineConfirmStep = iota // Delete
	InlineConfirmAsking                           // the button is its own confirm
	InlineConfirmDeleted                          // gone, with an undo
)

const (
	inlineConfirmH       = 44
	inlineConfirmTimeout = 3 * time.Second // an unanswered confirm backs off
	inlineConfirmRest    = 5 * time.Second // the demo's item comes back
	inlineConfirmRed     = "#e5484d"
)

var inlineConfirmW = map[InlineConfirmStep]float32{InlineConfirmIdle: 112, InlineConfirmAsking: 160, InlineConfirmDeleted: 144}

type inlineConfirmState struct {
	step InlineConfirmStep
	at   time.Time
}

// InlineConfirm is a delete button that needs no dialog. Press Delete and the
// button turns into its own confirm; once the item is gone it offers an undo.
func InlineConfirm(c *ui.Context, opts ...InlineConfirmOptions) ui.Element {
	var o InlineConfirmOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 22
	}
	o.Corner = clamp32(o.Corner, 0, 22)
	t := ThemeFrom(c)

	holder := ui.Box(c).Size(inlineConfirmW[InlineConfirmAsking], inlineConfirmH)
	st := ui.Local(holder, "inline_confirm", func() inlineConfirmState { return inlineConfirmState{} })
	now := c.Now()
	switch {
	case st.step == InlineConfirmAsking && now.Sub(st.at) >= inlineConfirmTimeout,
		st.step == InlineConfirmDeleted && now.Sub(st.at) >= inlineConfirmRest:
		st.step, st.at = InlineConfirmIdle, now
	}
	if st.step != InlineConfirmIdle {
		c.AnimationFrame()
	}
	moveTo := func(s InlineConfirmStep) { st.step, st.at = s, now }

	w := holder.Spring("ic_w", inlineConfirmW[st.step], ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.25})
	red := holder.Spring("ic_red", b2f(st.step == InlineConfirmAsking), ui.SpringOptions{Duration: 260 * time.Millisecond})
	// Each face fades in as the button reaches it.
	face := func(s InlineConfirmStep) float32 {
		return holder.Spring(s, b2f(st.step == s), ui.SpringOptions{Duration: 240 * time.Millisecond})
	}
	idle, asking, deleted := face(InlineConfirmIdle), face(InlineConfirmAsking), face(InlineConfirmDeleted)

	holder.Children(func() {
		btn := ui.Box(c).Absolute().Left((inlineConfirmW[InlineConfirmAsking]-w)/2).Top(0).Size(w, inlineConfirmH).
			Radius(o.Corner).Background(lerpColor(t.Surface, ui.Hex(inlineConfirmRed), red)).Clip().Cursor(ui.CursorPointer)
		press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond, Bounce: 0.4})
		btn.Scale(1-0.04*press, 1-0.04*press)
		switch st.step {
		case InlineConfirmIdle:
			btn.Label("Delete").OnClick(func() { moveTo(InlineConfirmAsking) })
		case InlineConfirmAsking:
			btn.Label("Confirm delete").OnClick(func() {
				moveTo(InlineConfirmDeleted)
				if o.OnDelete != nil {
					o.OnDelete()
				}
			})
		case InlineConfirmDeleted:
			btn.Label("Undo").OnClick(func() {
				moveTo(InlineConfirmIdle)
				if o.OnUndo != nil {
					o.OnUndo()
				}
			})
		}
		white := ui.Hex("#ffffff")
		btn.Children(func() {
			inlineConfirmFace(c, "idle", w, idle, -6*(1-idle), func() {
				ui.Icon(c, icons.Trash()).Size(18, 18).TextColor(t.Text)
				ui.Text(c, "Delete").FontSize(14).FontWeight(400).TextColor(t.Text)
			})
			inlineConfirmFace(c, "asking", w, asking, 6*(1-asking), func() {
				ui.Icon(c, icons.Trash()).Size(18, 18).TextColor(white)
				ui.Text(c, "Confirm delete").FontSize(14).FontWeight(500).TextColor(white)
			})
			inlineConfirmFace(c, "deleted", w, deleted, 6*(1-deleted), func() {
				ui.Text(c, "Deleted").FontSize(14).FontWeight(400).TextColor(t.TextSubtle)
				ui.Box(c).Size(1, 16).Background(t.Border)
				ui.Icon(c, icons.RotateCcw()).Size(16, 16).TextColor(t.Text)
				ui.Text(c, "Undo").FontSize(14).FontWeight(500).TextColor(t.Text)
			})
		})
	})
	return holder
}

// inlineConfirmFace is one face of the button, centred, sliding up as it fades in.
func inlineConfirmFace(c *ui.Context, key string, w, alpha, dy float32, children func()) {
	if alpha < 0.01 {
		return
	}
	ui.Box(c).Key(key).Absolute().Left(0).Top(0).Size(w, inlineConfirmH).Row().Center().Gap(8).
		Opacity(clamp32(alpha, 0, 1)).Translate(0, dy).Children(children)
}

func init() {
	registerDemo(
		"confirm",
		"An inline delete confirmation: press Delete and the button turns into its own confirm, then offers an undo once the item is gone.",
		func(c *ui.Context) { InlineConfirm(c) },
	)
}
