package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// ConfirmSendPhase enumerates the finite state machine states of ConfirmSend.
type ConfirmSendPhase int

const (
	ConfirmSendIdle ConfirmSendPhase = iota
	ConfirmSendConfirming
	ConfirmSendSending
	ConfirmSendSent
)

// ConfirmSendOptions configures ConfirmSend.
type ConfirmSendOptions struct {
	IdleLabel      string
	ConfirmLabel   string
	CancelLabel    string
	SendingLabel   string
	SentLabel      string
	ConfirmTimeout time.Duration
	SendDuration   time.Duration
	ResetAfter     time.Duration
	OnPhaseChange  func(ConfirmSendPhase)
}

type confirmSendState struct {
	phase   ConfirmSendPhase
	phaseAt time.Time
}

// ConfirmSend renders a two-step inline confirmation send button driven by an
// explicit state machine (Idle -> Confirming -> Sending -> Sent).
func ConfirmSend(c *ui.Context, onSend func(), opts ...ConfirmSendOptions) ui.Element {
	var o ConfirmSendOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.IdleLabel == "" {
		o.IdleLabel = "Send invoice"
	}
	if o.ConfirmLabel == "" {
		o.ConfirmLabel = "Confirm"
	}
	if o.CancelLabel == "" {
		o.CancelLabel = "Cancel"
	}
	if o.SendingLabel == "" {
		o.SendingLabel = "Sending..."
	}
	if o.SentLabel == "" {
		o.SentLabel = "Sent"
	}
	if o.ConfirmTimeout <= 0 {
		o.ConfirmTimeout = 4 * time.Second
	}
	if o.SendDuration <= 0 {
		o.SendDuration = 650 * time.Millisecond
	}
	if o.ResetAfter <= 0 {
		o.ResetAfter = 1800 * time.Millisecond
	}

	root := ui.Row(c).
		AlignItems(ui.Center).
		Gap(8)

	st := ui.Local(root, "confirm_send_fsm", func() confirmSendState {
		return confirmSendState{phase: ConfirmSendIdle, phaseAt: c.Now()}
	})

	now := c.Now()
	transition := func(next ConfirmSendPhase) {
		if st.phase == next {
			return
		}
		st.phase = next
		st.phaseAt = now
		if o.OnPhaseChange != nil {
			o.OnPhaseChange(next)
		}
	}

	// Time-based transitions.
	elapsed := now.Sub(st.phaseAt)
	switch st.phase {
	case ConfirmSendConfirming:
		c.AnimationFrame()
		if elapsed >= o.ConfirmTimeout {
			transition(ConfirmSendIdle)
		}
	case ConfirmSendSending:
		c.AnimationFrame()
		if elapsed >= o.SendDuration {
			transition(ConfirmSendSent)
			if onSend != nil {
				onSend()
			}
		}
	case ConfirmSendSent:
		c.AnimationFrame()
		if elapsed >= o.ResetAfter {
			transition(ConfirmSendIdle)
		}
	}

	root.Children(func() {
		switch st.phase {
		case ConfirmSendIdle:
			btn := ui.Row(c).
				Key("idle").
				Center().
				Height(38).
				Padding(0, 18).
				Radius(999).
				Background(t.Text).
				Cursor(ui.CursorPointer)
			scaleTarget := float32(1)
			if btn.Pressed() {
				scaleTarget = 0.96
			}
			sc := btn.Spring("cs_scale", scaleTarget, ui.SpringOptions{Duration: 180 * time.Millisecond, Bounce: 0.2})
			btn.Scale(sc, sc).Children(func() {
				ui.Text(c, o.IdleLabel).
					FontSize(13).
					FontWeight(600).
					TextColor(t.Background)
			})
			if btn.Clicked() {
				transition(ConfirmSendConfirming)
			}

		case ConfirmSendConfirming:
			cancelBtn := ui.Row(c).
				Key("cancel").
				Center().
				Height(38).
				Padding(0, 14).
				Radius(999).
				Background(t.SurfaceSubtle).
				Border(1, t.BorderStrong).
				Cursor(ui.CursorPointer)
			cancelBtn.Children(func() {
				ui.Text(c, o.CancelLabel).
					FontSize(13).
					FontWeight(500).
					TextColor(t.TextMuted)
			})
			if cancelBtn.Clicked() {
				transition(ConfirmSendIdle)
			}

			confirmBtn := ui.Row(c).
				Key("confirm").
				Center().
				Height(38).
				Padding(0, 18).
				Radius(999).
				Background(t.Accent).
				Cursor(ui.CursorPointer)
			confirmBtn.Children(func() {
				ui.Text(c, o.ConfirmLabel).
					FontSize(13).
					FontWeight(600).
					TextColor(t.AccentText)
			})
			if confirmBtn.Clicked() {
				transition(ConfirmSendSending)
			}

		case ConfirmSendSending:
			frac := clamp32(float32(elapsed.Seconds()/o.SendDuration.Seconds()), 0, 1)
			pill := ui.Row(c).
				Key("sending").
				Center().
				Height(38).
				Padding(0, 18).
				Radius(999).
				Background(t.SurfaceSubtle).
				Border(1, t.Accent.Alpha(0.55)).
				Clip()
			pill.Draw(func(p *ui.Painter, r ui.Rect) {
				if r.W <= 0 || r.H <= 0 {
					return
				}
				p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * frac, H: r.H}, t.Accent.Alpha(0.22), r.H*0.5)
			})
			pill.Children(func() {
				ui.Text(c, o.SendingLabel).
					FontSize(13).
					FontWeight(600).
					TextColor(t.Text)
			})

		case ConfirmSendSent:
			doneBtn := ui.Row(c).
				Key("sent").
				Center().
				Height(38).
				Padding(0, 18).
				Radius(999).
				Background(t.Success.Alpha(0.16)).
				Border(1, t.Success.Alpha(0.60)).
				Cursor(ui.CursorPointer)
			doneBtn.Children(func() {
				ui.Text(c, o.SentLabel).
					FontSize(13).
					FontWeight(600).
					TextColor(t.Success)
			})
			if doneBtn.Clicked() {
				transition(ConfirmSendIdle)
			}
		}
	})

	return root
}

type confirmSendDemoState struct {
	sentCount int
}

func init() {
	registerDemo(
		"confirm-send",
		"Two-step inline send confirmation powered by a four-state machine.",
		func(c *ui.Context) {
			holder := ui.Column(c).AlignItems(ui.Center).Gap(8)
			st := ui.Local(holder, "demo_confirm_send", func() confirmSendDemoState {
				return confirmSendDemoState{}
			})
			holder.Children(func() {
				ConfirmSend(c, func() {
					st.sentCount++
				})
				if st.sentCount > 0 {
					t := ThemeFrom(c)
					ui.Text(c, "Dispatched").FontSize(11).TextColor(t.Success)
				}
			})
		},
	)
}
