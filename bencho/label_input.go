package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// LabelInputField is what LabelInput is set up to take.
type LabelInputField string

const (
	LabelInputEmail    LabelInputField = "Email"
	LabelInputPassword LabelInputField = "Password"
)

// LabelInputOptions configures LabelInput.
type LabelInputOptions struct {
	Field  LabelInputField // Email by default
	Corner float32         // 0 to 26; 20 by default
}

const (
	labelInputW    = 280
	labelInputH    = 52
	labelInputTop  = 12 // room above the field for the lifted label
	labelInputPadX = 24
)

type labelInputState struct {
	focused bool // last frame's focus
	want    bool // a press asked for the focus
}

// LabelInput is a floating label field: the label sits inside it until you
// need it, then lifts into its top edge, cutting a gap in the border, and
// leaves room to type. It stays up while the field holds text.
func LabelInput(c *ui.Context, value *string, opts ...LabelInputOptions) ui.Element {
	var o LabelInputOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Field == "" {
		o.Field = LabelInputEmail
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Corner = clamp32(o.Corner, 0, 26)
	t := ThemeFrom(c)

	root := ui.Box(c).Size(labelInputW, labelInputTop+labelInputH+1).Cursor(ui.CursorText)
	st := ui.Local(root, "label_input", func() labelInputState { return labelInputState{} })
	up := root.Spring("up", b2f(st.focused || *value != ""), ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.15})
	focus := root.Spring("focus", b2f(st.focused), ui.SpringOptions{Duration: 200 * time.Millisecond})
	hover := root.Spring("hover", b2f(root.Hovered()), ui.SpringOptions{Duration: 160 * time.Millisecond})
	lift := clamp32(up, 0, 1)

	label := string(o.Field)
	size := lerp32(16, 12, lift)
	labelW, _ := c.MeasureText(0, ui.Span{Text: label, Size: size})
	rest := lerpColor(t.BorderStrong, t.TextSubtle, 0.45)
	border := lerpColor(lerpColor(rest, t.TextSubtle, hover*0.6), t.Text, focus)

	var field ui.Element
	root.OnClick(func() { st.want = true })
	root.Children(func() {
		// The border, with a gap where the lifted label sits on it.
		ui.Box(c).Absolute().Left(0).Top(labelInputTop).Size(labelInputW, labelInputH).Draw(func(p *ui.Painter, r ui.Rect) {
			gap := (labelW + 8) * lift
			var cut ui.Path
			cut.MoveTo(r.X-2, r.Y-2).LineTo(r.X+r.W+2, r.Y-2).LineTo(r.X+r.W+2, r.Y+r.H+2).LineTo(r.X-2, r.Y+r.H+2).Close()
			if gap > 0.5 {
				gx := r.X + labelInputPadX - 4 + (labelW+8-gap)/2
				cut.MoveTo(gx, r.Y-2).LineTo(gx+gap, r.Y-2).LineTo(gx+gap, r.Y+3).LineTo(gx, r.Y+3).Close()
			}
			cut.FillRule(ui.EvenOdd)
			p.ClipPath(&cut, func() {
				p.Stroke(ui.Rect{X: r.X + 0.5, Y: r.Y + 0.5, W: r.W - 1, H: r.H - 1}, border, o.Corner, 1)
			})
		})

		field = ui.TextInputBase(c, value).FontSize(16).TextColor(t.Text).Background(ui.Color{}).Padding(0).
			Absolute().Left(labelInputPadX).Top(labelInputTop+(labelInputH-24)/2).Size(labelInputW-2*labelInputPadX, 24)
		if o.Field == LabelInputPassword {
			field.Password()
		}
		if st.want {
			field.Focus()
			st.want = false
		}

		// The label: inside the field at rest, on its top edge when up.
		y := lerp32(labelInputTop+labelInputH/2, labelInputTop, up)
		ui.Text(c, label).Absolute().Left(labelInputPadX).Top(y - size*0.62).FontSize(size).SingleLine().
			TextColor(lerpColor(lerpColor(t.TextMuted, t.TextSubtle, 0.2), t.Text, lift)).PassThrough()
	})
	if f := field.Focused(); f != st.focused {
		st.focused = f
		c.AnimationFrame()
	}
	return root
}

func init() {
	registerDemo(
		"label-input",
		"A floating label field: the label sits inside until you focus it, then lifts into the top edge and leaves room to type.",
		func(c *ui.Context) {
			holder := ui.Box(c)
			email := ui.Local(holder, "demo", func() string { return "" })
			holder.Children(func() { LabelInput(c, email) })
		},
	)
}
