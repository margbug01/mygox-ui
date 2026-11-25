package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// AIChatInputOptions configures AIChatInput.
type AIChatInputOptions struct {
	Bounce     float32 // how far the box overshoots as it opens, 0 to 100; 40 by default
	HideAttach bool    // drops the + for adding files
	HideModel  bool    // drops the model menu
	HideEffort bool    // drops the Low / Medium / High step
	Models     []string
	OnSend     func(prompt string)
}

// AIChatInputEffort is how hard the model thinks.
type AIChatInputEffort int

const (
	AIChatInputLow AIChatInputEffort = iota
	AIChatInputMedium
	AIChatInputHigh
)

func (e AIChatInputEffort) String() string { return [...]string{"Low", "Medium", "High"}[e] }

const (
	aiChatInputPillW = 240
	aiChatInputPillH = 48
	aiChatInputBoxW  = 300
	aiChatInputBoxH  = 112
	aiChatInputBtn   = 32
)

var aiChatInputDefaultModels = []string{"Atlas 2", "Atlas 2 Mini", "Atlas 2 Pro"}

type aiChatInputState struct {
	open   bool
	prompt string
	model  int
	effort AIChatInputEffort
	menu   bool
	focus  bool
}

// AIChatInput is a prompt bar that stretches from a pill into a box when
// pressed, with a model menu and an effort step below; the mic turns into a
// send arrow as you type.
func AIChatInput(c *ui.Context, opts ...AIChatInputOptions) ui.Element {
	var o AIChatInputOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Bounce == 0 {
		o.Bounce = 40
	}
	if len(o.Models) == 0 {
		o.Models = aiChatInputDefaultModels
	}
	t := ThemeFrom(c)

	// The block fills its stage, so a press beside the box closes it.
	root := ui.Box(c).Size(340, 200).Center()
	st := ui.Local(root, "ai_chat_input", func() aiChatInputState { return aiChatInputState{effort: AIChatInputMedium} })
	st.model = min(st.model, len(o.Models)-1)
	root.OnClick(func() {
		st.menu = false
		if st.prompt == "" {
			st.open = false
		}
	})
	send := func() {
		if st.prompt == "" {
			return
		}
		if o.OnSend != nil {
			o.OnSend(st.prompt)
		}
		st.prompt = ""
	}

	spring := ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: float64(clamp32(o.Bounce, 0, 100)) / 100 * 0.5}
	p := root.Spring("open", b2f(st.open), spring)
	in := clamp32(p, 0, 1)
	w := lerp32(aiChatInputPillW, aiChatInputBoxW, p)
	h := lerp32(aiChatInputPillH, aiChatInputBoxH, p)
	typing := root.Spring("typing", b2f(st.prompt != ""), ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.3})

	root.Children(func() {
		bar := ui.Box(c).Size(w, h).Radius(lerp32(aiChatInputPillH/2, 22, in)).Background(t.Surface).Clip().
			Cursor(ui.CursorText).OnClick(func() { st.open, st.focus, st.menu = true, true, false })
		bar.Children(func() {
			field := ui.TextInputBase(c, &st.prompt).Placeholder("Ask anything...").PlaceholderColor(lerpColor(t.TextSubtle, t.Surface, 0.15)).
				FontSize(15).TextColor(t.Text).Background(ui.Color{}).Padding(0).
				Absolute().Left(lerp32(20, 16, in)).Top(lerp32((aiChatInputPillH-22)/2, 14, in)).Size(w-lerp32(20+aiChatInputBtn+16, 32, in), 22).
				OnSubmit(send).
				OnTextKey(func(k ui.TextKey) bool {
					if k.Key == ui.KeyEscape && k.Text() == "" {
						st.open, st.menu = false, false
						return true
					}
					return false
				})
			if st.focus {
				field.Focus()
				st.focus = false
			}
			if field.Focused() && !st.open {
				st.open = true
			}

			// The tools along the bottom, once open.
			row := clamp32((p-0.45)/0.55, 0, 1)
			if row > 0 {
				ui.Row(c).Absolute().Left(10).Top(h-10-aiChatInputBtn).Height(aiChatInputBtn).Gap(6).AlignItems(ui.Center).
					Opacity(row).Translate(0, (1-row)*6).Children(func() {
					if !o.HideAttach {
						aiChatInputRound(c, "attach", icons.Plus(), t.TextMuted, ui.Color{}, t.Border, func() {})
					}
					if !o.HideModel {
						aiChatInputChip(c, "model", func() {
							ui.Text(c, o.Models[st.model]).FontSize(13).FontWeight(500).TextColor(t.Text)
							ui.Icon(c, icons.ChevronDown()).Size(14, 14).TextColor(t.TextSubtle).Margin(0, 0, 0, 4).Rotate(180 * b2f(st.menu))
						}, func() { st.menu = !st.menu })
					}
					if !o.HideEffort {
						aiChatInputChip(c, "effort", func() {
							ui.Box(c).Size(12, 12).Margin(0, 6, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
								for i := range 3 {
									bh := float32(4 + 4*i)
									col := t.BorderStrong
									if AIChatInputEffort(i) <= st.effort {
										col = t.Text
									}
									p.Fill(ui.Rect{X: r.X + float32(i)*4.5, Y: r.Y + r.H - bh, W: 3, H: bh}, col, 1)
								}
							})
							ui.Text(c, st.effort.String()).FontSize(13).FontWeight(500).TextColor(t.Text)
						}, func() { st.effort = (st.effort + 1) % 3; st.menu = false })
					}
				})
			}

			// Mic while empty, a send arrow once there is something to send.
			bx := w - 8 - aiChatInputBtn
			by := lerp32((aiChatInputPillH-aiChatInputBtn)/2, h-10-aiChatInputBtn, p)
			btn := ui.Box(c).Key("send").Absolute().Left(bx).Top(by).Size(aiChatInputBtn, aiChatInputBtn).Radius(aiChatInputBtn / 2).
				Background(lerpColor(t.SurfaceSubtle, t.Text, clamp32(typing, 0, 1))).Cursor(ui.CursorPointer).KeepFocus().
				OnClick(func() {
					if st.prompt == "" {
						st.open, st.focus = true, true
						return
					}
					send()
				})
			btn.Children(func() {
				ui.Icon(c, icons.Mic()).Absolute().Left(7).Top(7).Size(18, 18).TextColor(t.Text).
					Opacity(clamp32(1-typing, 0, 1)).Scale(lerp32(1, 0.5, clamp32(typing, 0, 1)), lerp32(1, 0.5, clamp32(typing, 0, 1)))
				ui.Icon(c, icons.ArrowUp()).Absolute().Left(7).Top(7).Size(18, 18).TextColor(t.Surface).
					Opacity(clamp32(typing, 0, 1)).Scale(lerp32(0.5, 1, typing), lerp32(0.5, 1, typing))
			})
		})

		// The model menu opens upward from its chip.
		menu := root.Spring("menu", b2f(st.menu && st.open), ui.SpringOptions{Duration: 220 * time.Millisecond})
		if menu > 0.01 && !o.HideModel {
			left := (340-w)/2 + 10 + b2f(!o.HideAttach)*(aiChatInputBtn+6)
			rowsH := float32(len(o.Models))*30 + 8
			bottom := (200-h)/2 + h - 10 - aiChatInputBtn - 6
			ui.Column(c).Key("menu").Absolute().Left(left).Top(bottom-rowsH).Size(160, rowsH).Padding(4, 4).
				Radius(12).Background(t.Surface).Border(1, t.Border).Shadow(0, 8, 24, 0, ui.RGBA(0, 0, 0, 0.1)).
				Opacity(clamp32(menu, 0, 1)).Scale(lerp32(0.94, 1, menu), lerp32(0.94, 1, menu)).Translate(0, (1-menu)*6).
				KeepFocus().Children(func() {
				for i, m := range o.Models {
					item := ui.Row(c).Key(i).Height(30).Padding(0, 8).Radius(8).AlignItems(ui.Center).Cursor(ui.CursorPointer).
						OnClick(func() { st.model, st.menu = i, false })
					item.Background(t.SurfaceSubtle.Alpha(b2f(item.Hovered()))).Children(func() {
						ui.Text(c, m).FontSize(13).TextColor(t.Text).Grow(1)
						if i == st.model {
							ui.Icon(c, icons.Check()).Size(14, 14).TextColor(t.Text)
						}
					})
				}
			})
		}
	})
	return root
}

// aiChatInputRound is a round tool button.
func aiChatInputRound(c *ui.Context, key string, icon *ui.SVG, ink, bg, border ui.Color, fn func()) {
	b := ui.Box(c).Key(key).Size(aiChatInputBtn, aiChatInputBtn).Radius(aiChatInputBtn/2).Border(1, border).Center().
		Cursor(ui.CursorPointer).KeepFocus().OnClick(fn)
	t := ThemeFrom(c)
	b.Background(lerpColor(bg, t.SurfaceSubtle, b2f(b.Hovered()))).Children(func() {
		ui.Icon(c, icon).Size(16, 16).TextColor(ink)
	})
}

// aiChatInputChip is a pill button along the bottom of the box.
func aiChatInputChip(c *ui.Context, key string, body func(), fn func()) {
	t := ThemeFrom(c)
	b := ui.Row(c).Key(key).Height(aiChatInputBtn).Padding(0, 12).Radius(aiChatInputBtn / 2).AlignItems(ui.Center).
		Cursor(ui.CursorPointer).KeepFocus().OnClick(fn)
	b.Background(lerpColor(t.SurfaceSubtle, t.SurfaceHover, b2f(b.Hovered())*0.6)).Children(body)
}

func init() {
	registerDemo(
		"ai-chat-input",
		"A prompt bar that stretches from a pill into a box when pressed, with a model menu and an effort step below; the mic turns into a send arrow as you type.",
		func(c *ui.Context) { AIChatInput(c) },
	)
}
