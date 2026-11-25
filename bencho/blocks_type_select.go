package bencho

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// ============================================================================
// 8. command — Command bar
// ============================================================================

// CommandItem describes a single runnable action in CommandBar.
type CommandItem struct {
	ID       string
	Title    string
	Shortcut string
}

// CommandBarOptions configures CommandBar.
type CommandBarOptions struct {
	Width       float32
	Placeholder string
	OnRun       func(CommandItem)
}

type commandBarState struct {
	query   string
	lastRan string
	ranAt   time.Time
}

// CommandBar renders a compact Cmd+K command palette with live search filtering
// and click-to-run action feedback.
func CommandBar(c *ui.Context, items []CommandItem, opts ...CommandBarOptions) ui.Element {
	var o CommandBarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 282
	}
	if o.Placeholder == "" {
		o.Placeholder = "Search actions..."
	}

	card := ui.Column(c).
		Width(o.Width).
		Radius(14).
		Background(t.Surface).
		Border(1, t.BorderStrong).
		Clip()

	st := ui.Local(card, "cmd_bar_st", func() commandBarState {
		return commandBarState{}
	})

	now := c.Now()
	q := strings.ToLower(strings.TrimSpace(st.query))
	filtered := make([]CommandItem, 0, len(items))
	for _, it := range items {
		if q == "" || strings.Contains(strings.ToLower(it.Title), q) || strings.Contains(strings.ToLower(it.Shortcut), q) {
			filtered = append(filtered, it)
		}
	}

	runItem := func(it CommandItem) {
		st.lastRan = it.Title
		st.ranAt = now
		if o.OnRun != nil {
			o.OnRun(it)
		}
	}

	card.Children(func() {
		// Search header
		ui.Row(c).
			WidthPercent(100).
			Padding(7, 10).
			Gap(7).
			AlignItems(ui.Center).
			BorderWidth(0, 0, 1, 0).
			BorderColor(t.Border).
			Children(func() {
				ui.Row(c).
					Center().
					Height(20).
					Padding(0, 5).
					Radius(5).
					Background(t.SurfaceSubtle).
					Border(1, t.Border).
					Children(func() {
						ui.Text(c, "⌘K").FontSize(9).FontWeight(700).TextColor(t.TextMuted)
					})
				inp := ui.TextInput(c, &st.query).
					Width(o.Width - 78).
					Placeholder(o.Placeholder).
					FontSize(12)
				if inp.Submitted() && len(filtered) > 0 {
					runItem(filtered[0])
				}
				if st.query != "" {
					clr := ui.Row(c).
						Key("cmd_clr").
						Center().
						Width(18).
						Height(18).
						Radius(999).
						Background(t.SurfaceSubtle).
						Cursor(ui.CursorPointer)
					clr.Children(func() {
						ui.Text(c, "×").FontSize(11).TextColor(t.TextMuted)
					})
					if clr.Clicked() {
						st.query = ""
					}
				}
			})

		// Filtered results (up to 3 visible)
		ui.Column(c).WidthPercent(100).Padding(5, 6).Gap(2).Children(func() {
			if len(filtered) == 0 {
				ui.Row(c).WidthPercent(100).Height(28).Center().Children(func() {
					ui.Text(c, "No matching commands").FontSize(11).TextColor(t.TextSubtle)
				})
			} else {
				limit := len(filtered)
				if limit > 3 {
					limit = 3
				}
				for i := 0; i < limit; i++ {
					it := filtered[i]
					row := ui.Row(c).
						Key("cmd_row_"+it.ID).
						WidthPercent(100).
						Height(28).
						Padding(0, 8).
						Radius(8).
						Justify(ui.SpaceBetween).
						AlignItems(ui.Center).
						Cursor(ui.CursorPointer)
					if row.Hovered() {
						row.Background(t.AccentSoft)
					} else if i == 0 {
						row.Background(t.SurfaceSubtle)
					}
					row.Children(func() {
						txtCol := t.Text
						if row.Hovered() {
							txtCol = t.Accent
						}
						ui.Text(c, it.Title).FontSize(11).FontWeight(500).TextColor(txtCol)
						if it.Shortcut != "" {
							ui.Row(c).
								Center().
								Height(18).
								Padding(0, 5).
								Radius(4).
								Background(t.SurfaceSubtle).
								Border(1, t.Border).
								Children(func() {
									ui.Text(c, it.Shortcut).FontSize(9).FontWeight(600).TextColor(t.TextSubtle)
								})
						}
					})
					if row.Clicked() {
						runItem(it)
					}
				}
			}
		})

		// Footer feedback bar
		ui.Row(c).
			WidthPercent(100).
			Height(24).
			Padding(0, 10).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Background(t.SurfaceSubtle).
			BorderWidth(1, 0, 0, 0).
			BorderColor(t.Border).
			Children(func() {
				if st.lastRan != "" {
					ui.Text(c, "✓ Ran "+st.lastRan).FontSize(10).FontWeight(600).TextColor(t.Success)
				} else {
					ui.Text(c, "Click or press Enter").FontSize(10).TextColor(t.TextSubtle)
				}
				ui.Text(c, fmt.Sprintf("%d/%d", len(filtered), len(items))).FontSize(10).TextColor(t.TextSubtle)
			})
	})

	return card
}

// ============================================================================
// 10. picker — Assignees
// ============================================================================

// Assignee represents one collaborator row in AssigneePicker.
type Assignee struct {
	ID       string
	Name     string
	Initials string
	Color    ui.Color
	Selected bool
}

// AssigneePickerOptions configures AssigneePicker.
type AssigneePickerOptions struct {
	Width    float32
	Title    string
	OnChange func([]Assignee)
}

// AssigneePicker renders a multi-select collaborator picker with colored
// initials avatars, spring checkmarks, and a live avatar stack preview.
func AssigneePicker(c *ui.Context, members []Assignee, opts ...AssigneePickerOptions) ui.Element {
	var o AssigneePickerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 278
	}
	if o.Title == "" {
		o.Title = "Assignees"
	}

	selectedMembers := make([]Assignee, 0, len(members))
	for _, m := range members {
		if m.Selected {
			selectedMembers = append(selectedMembers, m)
		}
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(10).
		Gap(6).
		Radius(14).
		Background(t.Surface).
		Border(1, t.Border)

	card.Children(func() {
		// Header + selected overlapping avatar stack preview
		ui.Row(c).WidthPercent(100).Height(24).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, o.Title).FontSize(12).FontWeight(600).TextColor(t.Text)

			ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
				if len(selectedMembers) == 0 {
					ui.Text(c, "Unassigned").FontSize(10).TextColor(t.TextSubtle)
				} else {
					stackW := float32(22 + (len(selectedMembers)-1)*14)
					stack := ui.Box(c).Width(stackW).Height(22)
					selCopy := append([]Assignee(nil), selectedMembers...)
					surfCol := t.Surface
					stack.Draw(func(p *ui.Painter, r ui.Rect) {
						for i, m := range selCopy {
							cx := r.X + float32(i)*14
							avRect := ui.Rect{X: cx, Y: r.Y, W: 22, H: 22}
							p.Fill(avRect, m.Color, 11)
							p.Stroke(avRect, surfCol, 11, 2)
							tw, th := p.MeasureText(22, ui.Span{Text: m.Initials, Size: 8.5, Weight: 700})
							p.Text(cx+(22-tw)*0.5, r.Y+(22-th)*0.5, m.Initials, 8.5, ui.Hex("#ffffff"))
						}
					})
					ui.Text(c, fmt.Sprintf("%d", len(selectedMembers))).FontSize(10).FontWeight(700).TextColor(t.Accent)
				}
			})
		})

		// Collaborator rows
		ui.Column(c).WidthPercent(100).Gap(3).Children(func() {
			for i := range members {
				idx := i
				m := members[idx]
				rowBg := ui.Color{}
				if m.Selected {
					rowBg = t.AccentSoft
				}
				row := ui.Row(c).
					Key("asg_"+m.ID).
					WidthPercent(100).
					Height(28).
					Padding(0, 8).
					Radius(8).
					Justify(ui.SpaceBetween).
					AlignItems(ui.Center).
					Cursor(ui.CursorPointer)
				if m.Selected {
					row.Background(rowBg)
				} else if row.Hovered() {
					row.Background(t.SurfaceSubtle)
				}

				chkTarget := float32(0)
				if m.Selected {
					chkTarget = 1
				}
				chk := row.Spring("asg_chk", chkTarget, ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.30})

				row.Children(func() {
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						ui.Row(c).
							Center().
							Width(20).
							Height(20).
							Radius(999).
							Background(m.Color).
							Children(func() {
								ui.Text(c, m.Initials).FontSize(8.5).FontWeight(700).TextColor(ui.Hex("#ffffff"))
							})
						ui.Text(c, m.Name).FontSize(11.5).FontWeight(500).TextColor(t.Text)
					})

					badgeBg := lerpColor(t.SurfaceSubtle, t.Accent, chk)
					badgeBorder := lerpColor(t.BorderStrong, t.Accent, chk)
					ui.Row(c).
						Center().
						Width(17).
						Height(17).
						Radius(5).
						Background(badgeBg).
						Border(1, badgeBorder).
						Children(func() {
							if chk > 0.35 {
								ui.Text(c, "✓").FontSize(10).FontWeight(700).TextColor(t.AccentText)
							}
						})
				})

				if row.Clicked() {
					members[idx].Selected = !members[idx].Selected
					if o.OnChange != nil {
						o.OnChange(members)
					}
				}
			}
		})
	})

	return card
}

type assigneePickerDemoState struct {
	members []Assignee
}

// ============================================================================
// 11. aspect — Aspect ratio
// ============================================================================

// AspectPreset defines one selectable aspect ratio.
type AspectPreset struct {
	Label string
	W     float32
	H     float32
}

// AspectRatioOptions configures AspectRatioSelector.
type AspectRatioOptions struct {
	Width    float32
	Presets  []AspectPreset
	OnSelect func(int, AspectPreset)
}

func defaultAspectPresets() []AspectPreset {
	return []AspectPreset{
		{Label: "1:1", W: 1, H: 1},
		{Label: "4:3", W: 4, H: 3},
		{Label: "16:9", W: 16, H: 9},
		{Label: "9:16", W: 9, H: 16},
		{Label: "21:9", W: 21, H: 9},
	}
}

func fitAspectBox(ratioW, ratioH, maxW, maxH float32) (w, h float32) {
	if ratioW <= 0 || ratioH <= 0 {
		return maxH, maxH
	}
	scaleW := maxW / ratioW
	scaleH := maxH / ratioH
	s := scaleW
	if scaleH < s {
		s = scaleH
	}
	return ratioW * s, ratioH * s
}

// AspectRatioSelector renders visual aspect ratio pills (1:1, 4:3, 16:9, 9:16,
// 21:9) paired with a live spring-morphing preview frame.
func AspectRatioSelector(c *ui.Context, selected *int, opts ...AspectRatioOptions) ui.Element {
	var o AspectRatioOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 282
	}
	if len(o.Presets) == 0 {
		o.Presets = defaultAspectPresets()
	}

	curSel := 2
	if selected != nil {
		curSel = *selected
	}
	if curSel < 0 || curSel >= len(o.Presets) {
		curSel = 0
	}
	preset := o.Presets[curSel]

	stageW := o.Width - 20
	stageH := float32(106)
	targetW, targetH := fitAspectBox(preset.W, preset.H, stageW-36, stageH-22)

	card := ui.Column(c).
		Width(o.Width).
		Padding(10).
		Gap(9).
		Radius(14).
		Background(t.Surface).
		Border(1, t.Border).
		AlignItems(ui.Center)

	card.Children(func() {
		stage := ui.Box(c).
			Key("aspect_stage").
			Width(stageW).
			Height(stageH).
			Radius(10).
			Background(t.SurfaceSubtle).
			Border(1, t.Border)

		curW := stage.Spring("aspect_w", targetW, ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.28})
		curH := stage.Spring("aspect_h", targetH, ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.28})

		label := preset.Label
		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			// Subtle center crosshair
			cx := r.X + r.W*0.5
			cy := r.Y + r.H*0.5
			guideCol := t.BorderStrong.Alpha(0.45)
			p.Line(r.X+10, cy, r.X+r.W-10, cy, 1, guideCol)
			p.Line(cx, r.Y+8, cx, r.Y+r.H-8, 1, guideCol)

			// Spring-morphing frame
			fw := clamp32(curW, 24, r.W-16)
			fh := clamp32(curH, 24, r.H-12)
			fr := ui.Rect{X: cx - fw*0.5, Y: cy - fh*0.5, W: fw, H: fh}
			p.Fill(fr, t.AccentSoft, 8)
			p.Stroke(fr, t.Accent, 8, 1.75)

			// Corner ticks
			tk := float32(7)
			p.Line(fr.X+4, fr.Y+4, fr.X+4+tk, fr.Y+4, 1.5, t.Accent)
			p.Line(fr.X+4, fr.Y+4, fr.X+4, fr.Y+4+tk, 1.5, t.Accent)
			p.Line(fr.X+fr.W-4-tk, fr.Y+fr.H-4, fr.X+fr.W-4, fr.Y+fr.H-4, 1.5, t.Accent)
			p.Line(fr.X+fr.W-4, fr.Y+fr.H-4-tk, fr.X+fr.W-4, fr.Y+fr.H-4, 1.5, t.Accent)

			tw, th := p.MeasureText(fw, ui.Span{Text: label, Size: 11, Weight: 700})
			p.Text(cx-tw*0.5, cy-th*0.5, label, 11, t.Accent)
		})

		// Ratio selector pills
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).Gap(4).Children(func() {
			for i, p := range o.Presets {
				idx := i
				pr := p
				isSel := idx == curSel
				bg := t.SurfaceSubtle
				bd := t.Border
				tc := t.TextMuted
				if isSel {
					bg = t.Accent
					bd = t.Accent
					tc = t.AccentText
				}
				pill := ui.Row(c).
					Key("ar_"+pr.Label).
					Center().
					Grow(1).
					Height(28).
					Radius(8).
					Background(bg).
					Border(1, bd).
					Cursor(ui.CursorPointer)
				pill.Children(func() {
					ui.Text(c, pr.Label).FontSize(10.5).FontWeight(600).TextColor(tc)
				})
				if pill.Clicked() {
					if selected != nil {
						*selected = idx
					}
					if o.OnSelect != nil {
						o.OnSelect(idx, pr)
					}
				}
			}
		})
	})

	return card
}

type aspectDemoState struct {
	selected int
}

// ============================================================================
// 12. toolbar — Canvas toolbar
// ============================================================================

// CanvasTool describes one studio tool in CanvasToolbar.
type CanvasTool struct {
	ID       string
	Name     string
	Shortcut string
	Hint     string
}

// CanvasToolbarOptions configures CanvasToolbar.
type CanvasToolbarOptions struct {
	Tools    []CanvasTool
	OnSelect func(int, CanvasTool)
}

func defaultCanvasTools() []CanvasTool {
	return []CanvasTool{
		{ID: "select", Name: "Select", Shortcut: "V", Hint: "Move & resize"},
		{ID: "frame", Name: "Frame", Shortcut: "F", Hint: "Draw artboard"},
		{ID: "pen", Name: "Pen", Shortcut: "P", Hint: "Bezier vector"},
		{ID: "text", Name: "Text", Shortcut: "T", Hint: "Type layer"},
		{ID: "hand", Name: "Hand", Shortcut: "H", Hint: "Pan canvas"},
	}
}

func drawCanvasToolIcon(p *ui.Painter, id string, r ui.Rect, col ui.Color) {
	cx := r.X + r.W*0.5
	cy := r.Y + r.H*0.5
	switch id {
	case "select":
		// Crisp pointer arrow
		p.Line(cx-4, cy-6, cx-4, cy+5, 1.6, col)
		p.Line(cx-4, cy-6, cx+5, cy+1, 1.6, col)
		p.Line(cx-4, cy+5, cx, cy+2, 1.6, col)
		p.Line(cx, cy+2, cx+5, cy+1, 1.6, col)
	case "frame":
		// Artboard # frame
		p.Line(cx-4, cy-7, cx-4, cy+7, 1.5, col)
		p.Line(cx+4, cy-7, cx+4, cy+7, 1.5, col)
		p.Line(cx-7, cy-4, cx+7, cy-4, 1.5, col)
		p.Line(cx-7, cy+4, cx+7, cy+4, 1.5, col)
	case "pen":
		// Vector nib diamond + tip line
		p.Line(cx, cy-7, cx+5, cy, 1.5, col)
		p.Line(cx+5, cy, cx, cy+6, 1.5, col)
		p.Line(cx, cy+6, cx-5, cy, 1.5, col)
		p.Line(cx-5, cy, cx, cy-7, 1.5, col)
		p.Fill(ui.Rect{X: cx - 1.5, Y: cy - 1.5, W: 3, H: 3}, col, 1.5)
	case "text":
		// Serif-free T
		p.Line(cx-6, cy-5.5, cx+6, cy-5.5, 1.7, col)
		p.Line(cx, cy-5.5, cx, cy+6.5, 1.7, col)
		p.Line(cx-3, cy+6.5, cx+3, cy+6.5, 1.4, col)
	case "hand":
		// Pan crosshair with cardinal ticks
		p.Line(cx-6.5, cy, cx+6.5, cy, 1.5, col)
		p.Line(cx, cy-6.5, cx, cy+6.5, 1.5, col)
		p.Stroke(ui.Rect{X: cx - 2.5, Y: cy - 2.5, W: 5, H: 5}, col, 2.5, 1.4)
	default:
		p.Stroke(ui.Rect{X: cx - 5, Y: cy - 5, W: 10, H: 10}, col, 3, 1.5)
	}
}

// CanvasToolbar renders a floating studio tool dock (Select, Frame, Pen, Text,
// Hand) with a spring-sliding active pill and live tool status badge.
func CanvasToolbar(c *ui.Context, selected *int, opts ...CanvasToolbarOptions) ui.Element {
	var o CanvasToolbarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if len(o.Tools) == 0 {
		o.Tools = defaultCanvasTools()
	}

	curSel := 0
	if selected != nil {
		curSel = *selected
	}
	if curSel < 0 || curSel >= len(o.Tools) {
		curSel = 0
	}
	activeTool := o.Tools[curSel]

	slotW := float32(44)
	slotH := float32(38)
	gap := float32(4)
	pad := float32(6)
	dockW := pad*2 + float32(len(o.Tools))*slotW + float32(len(o.Tools)-1)*gap

	wrap := ui.Column(c).
		AlignItems(ui.Center).
		Gap(12)

	wrap.Children(func() {
		dock := ui.Row(c).
			Key("tb_dock").
			Width(dockW).
			Padding(pad).
			Gap(gap).
			Radius(999).
			Background(t.Surface).
			Border(1, t.BorderStrong).
			AlignItems(ui.Center)

		targetX := pad + float32(curSel)*(slotW+gap)
		slideX := dock.Spring("tb_slide_x", targetX, ui.SpringOptions{
			Duration: 280 * time.Millisecond,
			Bounce:   0.30,
		})

		dock.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			pill := ui.Rect{X: r.X + slideX, Y: r.Y + pad, W: slotW, H: slotH}
			p.Fill(pill, t.Accent, slotH*0.5)
		})

		dock.Children(func() {
			for i, tool := range o.Tools {
				idx := i
				tl := tool
				isSel := idx == curSel
				btn := ui.Box(c).
					Key("tb_" + tl.ID).
					Width(slotW).
					Height(slotH).
					Radius(999).
					Cursor(ui.CursorPointer)
				if !isSel && btn.Hovered() {
					btn.Background(t.SurfaceSubtle)
				}
				iconCol := t.TextMuted
				if isSel {
					iconCol = ui.Hex("#ffffff")
				} else if btn.Hovered() {
					iconCol = t.Text
				}
				id := tl.ID
				btn.Draw(func(p *ui.Painter, r ui.Rect) {
					drawCanvasToolIcon(p, id, r, iconCol)
				})
				if btn.Clicked() {
					if selected != nil {
						*selected = idx
					}
					if o.OnSelect != nil {
						o.OnSelect(idx, tl)
					}
				}
			}
		})

		// Active tool status badge
		ui.Row(c).
			Height(28).
			Padding(0, 12).
			Gap(7).
			Radius(999).
			AlignItems(ui.Center).
			Background(t.Surface).
			Border(1, t.Border).
			Children(func() {
				ui.Box(c).Width(6).Height(6).Radius(999).Background(t.Accent)
				ui.Text(c, activeTool.Name).FontSize(11.5).FontWeight(600).TextColor(t.Text)
				ui.Row(c).
					Center().
					Height(18).
					Padding(0, 5).
					Radius(4).
					Background(t.SurfaceSubtle).
					Border(1, t.Border).
					Children(func() {
						ui.Text(c, activeTool.Shortcut).FontSize(9.5).FontWeight(700).TextColor(t.Accent)
					})
				ui.Text(c, activeTool.Hint).FontSize(10.5).TextColor(t.TextSubtle)
			})
	})

	return wrap
}

type canvasToolbarDemoState struct {
	selected int
}

// ============================================================================
// 13. icon-bar — Icon bar
// ============================================================================

// IconBarItem defines one icon button in IconBar.
type IconBarItem struct {
	ID    string
	Label string
	Icon  string
}

// IconBarOptions configures IconBar.
type IconBarOptions struct {
	OnSelect func(int)
}

type iconBarState struct {
	hovered int
}

func defaultIconBarItems() []IconBarItem {
	return []IconBarItem{
		{ID: "home", Label: "Home", Icon: "home"},
		{ID: "grid", Label: "Library", Icon: "grid"},
		{ID: "bolt", Label: "Actions", Icon: "bolt"},
		{ID: "bell", Label: "Alerts", Icon: "bell"},
		{ID: "gear", Label: "Settings", Icon: "gear"},
	}
}

func drawSegmentedIcon(p *ui.Painter, kind string, r ui.Rect, col ui.Color) {
	cx := r.X + r.W*0.5
	cy := r.Y + r.H*0.5
	switch kind {
	case "home":
		p.Line(cx-6.5, cy-0.5, cx, cy-6.5, 1.6, col)
		p.Line(cx, cy-6.5, cx+6.5, cy-0.5, 1.6, col)
		p.Stroke(ui.Rect{X: cx - 4.5, Y: cy - 0.5, W: 9, H: 6.5}, col, 1.5, 1.5)
	case "grid":
		p.Stroke(ui.Rect{X: cx - 6, Y: cy - 6, W: 5, H: 5}, col, 1.5, 1.4)
		p.Stroke(ui.Rect{X: cx + 1, Y: cy - 6, W: 5, H: 5}, col, 1.5, 1.4)
		p.Stroke(ui.Rect{X: cx - 6, Y: cy + 1, W: 5, H: 5}, col, 1.5, 1.4)
		p.Stroke(ui.Rect{X: cx + 1, Y: cy + 1, W: 5, H: 5}, col, 1.5, 1.4)
	case "bolt":
		p.Line(cx+1.5, cy-6.5, cx-4, cy+0.5, 1.6, col)
		p.Line(cx-4, cy+0.5, cx+1, cy+0.5, 1.6, col)
		p.Line(cx+1, cy+0.5, cx-1.5, cy+6.5, 1.6, col)
		p.Line(cx-1.5, cy+6.5, cx+4, cy-0.5, 1.6, col)
	case "bell":
		p.Stroke(ui.Rect{X: cx - 4.5, Y: cy - 5.5, W: 9, H: 8.5}, col, 4, 1.5)
		p.Line(cx-6.5, cy+3, cx+6.5, cy+3, 1.5, col)
		p.Fill(ui.Rect{X: cx - 1.5, Y: cy + 4.5, W: 3, H: 2.5}, col, 1.5)
	case "gear":
		p.Stroke(ui.Rect{X: cx - 4.5, Y: cy - 4.5, W: 9, H: 9}, col, 4.5, 1.5)
		p.Fill(ui.Rect{X: cx - 1.8, Y: cy - 1.8, W: 3.6, H: 3.6}, col, 1.8)
		p.Line(cx, cy-7, cx, cy-4.8, 1.6, col)
		p.Line(cx, cy+4.8, cx, cy+7, 1.6, col)
		p.Line(cx-7, cy, cx-4.8, cy, 1.6, col)
		p.Line(cx+4.8, cy, cx+7, cy, 1.6, col)
	default:
		p.Stroke(ui.Rect{X: cx - 5, Y: cy - 5, W: 10, H: 10}, col, 5, 1.5)
	}
}

// IconBar renders a compact glass/surface segmented icon bar with a spring
// active indicator and a tracking tooltip label.
func IconBar(c *ui.Context, items []IconBarItem, selected *int, opts ...IconBarOptions) ui.Element {
	var o IconBarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if len(items) == 0 {
		items = defaultIconBarItems()
	}

	curSel := 0
	if selected != nil {
		curSel = *selected
	}
	if curSel < 0 || curSel >= len(items) {
		curSel = 0
	}

	slotW := float32(42)
	slotH := float32(36)
	gap := float32(4)
	pad := float32(5)
	barW := pad*2 + float32(len(items))*slotW + float32(len(items)-1)*gap

	wrap := ui.Column(c).
		Width(barW).
		Gap(8).
		AlignItems(ui.Center)

	st := ui.Local(wrap, "icon_bar_st", func() iconBarState {
		return iconBarState{hovered: -1}
	})

	tipIdx := curSel
	if st.hovered >= 0 && st.hovered < len(items) {
		tipIdx = st.hovered
	}
	tipTargetX := pad + float32(tipIdx)*(slotW+gap) + slotW*0.5

	wrap.Children(func() {
		// Floating spring-tracking tooltip track
		tipTrack := ui.Box(c).
			Key("ib_tip_track").
			Width(barW).
			Height(24)

		tipX := tipTrack.Spring("ib_tip_x", tipTargetX, ui.SpringOptions{
			Duration: 230 * time.Millisecond,
			Bounce:   0.26,
		})
		tipLabel := items[tipIdx].Label
		tipTrack.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			tw, th := p.MeasureText(120, ui.Span{Text: tipLabel, Size: 10.5, Weight: 600})
			pw := tw + 16
			px := clamp32(r.X+tipX-pw*0.5, r.X, r.X+r.W-pw)
			pill := ui.Rect{X: px, Y: r.Y + 1, W: pw, H: 22}
			p.Fill(pill, t.Text, 8)
			p.Text(px+8, pill.Y+(pill.H-th)*0.5, tipLabel, 10.5, t.Background)
		})

		// Glass / surface segmented bar
		barBg := t.Surface.Alpha(0.88)
		bar := ui.Row(c).
			Key("ib_bar").
			Width(barW).
			Padding(pad).
			Gap(gap).
			Radius(15).
			Background(barBg).
			Border(1, t.BorderStrong).
			AlignItems(ui.Center)

		activeTargetX := pad + float32(curSel)*(slotW+gap)
		activeX := bar.Spring("ib_active_x", activeTargetX, ui.SpringOptions{
			Duration: 260 * time.Millisecond,
			Bounce:   0.28,
		})

		bar.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			ind := ui.Rect{X: r.X + activeX, Y: r.Y + pad, W: slotW, H: slotH}
			p.Fill(ind, t.AccentSoft, 10)
			p.Stroke(ind, t.Accent.Alpha(0.65), 10, 1.25)
		})

		nextHover := -1
		bar.Children(func() {
			for i, it := range items {
				idx := i
				item := it
				isSel := idx == curSel
				btn := ui.Box(c).
					Key("ib_btn_" + item.ID).
					Width(slotW).
					Height(slotH).
					Radius(10).
					Cursor(ui.CursorPointer)
				if btn.Hovered() {
					nextHover = idx
					if !isSel {
						btn.Background(t.SurfaceSubtle)
					}
				}
				col := t.TextMuted
				if isSel {
					col = t.Accent
				} else if btn.Hovered() {
					col = t.Text
				}
				iconKind := item.Icon
				btn.Draw(func(p *ui.Painter, r ui.Rect) {
					drawSegmentedIcon(p, iconKind, r, col)
				})
				if btn.Clicked() {
					if selected != nil {
						*selected = idx
					}
					if o.OnSelect != nil {
						o.OnSelect(idx)
					}
				}
			}
		})
		st.hovered = nextHover
	})

	return wrap
}

type iconBarDemoState struct {
	selected int
}

// ============================================================================
// 14. roster — Selection list
// ============================================================================

// RosterMember describes one selectable person in RosterList.
type RosterMember struct {
	ID       string
	Name     string
	Role     string
	Online   bool
	Selected bool
}

// RosterListOptions configures RosterList.
type RosterListOptions struct {
	Width    float32
	OnChange func([]RosterMember)
}

type rosterListState struct {
	filter string
}

// RosterList renders a filterable, multi-selectable team roster with role
// pills, live status dots, and a selection count badge.
func RosterList(c *ui.Context, items []RosterMember, opts ...RosterListOptions) ui.Element {
	var o RosterListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 282
	}

	selectedCount := 0
	for _, m := range items {
		if m.Selected {
			selectedCount++
		}
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(10).
		Gap(7).
		Radius(14).
		Background(t.Surface).
		Border(1, t.Border)

	st := ui.Local(card, "roster_st", func() rosterListState {
		return rosterListState{filter: "All"}
	})

	matchesFilter := func(m RosterMember) bool {
		switch st.filter {
		case "Online":
			return m.Online
		case "Eng":
			return m.Role == "Eng"
		default:
			return true
		}
	}

	card.Children(func() {
		// Filter pills + selection count badge
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Row(c).Gap(4).AlignItems(ui.Center).Children(func() {
				for _, f := range []string{"All", "Online", "Eng"} {
					name := f
					active := st.filter == name
					bg := t.SurfaceSubtle
					bd := t.Border
					tc := t.TextMuted
					if active {
						bg = t.Text
						bd = t.Text
						tc = t.Background
					}
					pill := ui.Row(c).
						Key("rf_"+name).
						Center().
						Height(22).
						Padding(0, 8).
						Radius(999).
						Background(bg).
						Border(1, bd).
						Cursor(ui.CursorPointer)
					pill.Children(func() {
						ui.Text(c, name).FontSize(10).FontWeight(600).TextColor(tc)
					})
					if pill.Clicked() {
						st.filter = name
					}
				}
			})

			cntBg := t.SurfaceSubtle
			cntTxt := t.TextSubtle
			if selectedCount > 0 {
				cntBg = t.AccentSoft
				cntTxt = t.Accent
			}
			ui.Row(c).
				Center().
				Height(22).
				Padding(0, 8).
				Radius(999).
				Background(cntBg).
				Children(func() {
					ui.Text(c, fmt.Sprintf("%d selected", selectedCount)).FontSize(10).FontWeight(700).TextColor(cntTxt)
				})
		})

		// Filtered rows
		ui.Column(c).WidthPercent(100).Gap(3).Children(func() {
			shown := 0
			for i := range items {
				idx := i
				m := items[idx]
				if !matchesFilter(m) {
					continue
				}
				shown++
				if shown > 4 {
					break
				}
				row := ui.Row(c).
					Key("rost_"+m.ID).
					WidthPercent(100).
					Height(29).
					Padding(0, 8).
					Radius(8).
					Justify(ui.SpaceBetween).
					AlignItems(ui.Center).
					Cursor(ui.CursorPointer)
				if m.Selected {
					row.Background(t.AccentSoft).Border(1, t.Accent.Alpha(0.40))
				} else if row.Hovered() {
					row.Background(t.SurfaceSubtle).Border(1, t.Border)
				} else {
					row.Border(1, ui.Color{})
				}

				row.Children(func() {
					ui.Row(c).Gap(7).AlignItems(ui.Center).Children(func() {
						dotCol := t.BorderStrong
						if m.Online {
							dotCol = t.Success
						}
						ui.Box(c).Width(7).Height(7).Radius(999).Background(dotCol)
						ui.Text(c, m.Name).FontSize(11.5).FontWeight(500).TextColor(t.Text)
					})

					ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
						ui.Row(c).
							Center().
							Height(18).
							Padding(0, 6).
							Radius(999).
							Background(t.SurfaceSubtle).
							Border(1, t.Border).
							Children(func() {
								ui.Text(c, m.Role).FontSize(9.5).FontWeight(600).TextColor(t.TextMuted)
							})

						chkCol := t.TextSubtle
						chkGlyph := "○"
						if m.Selected {
							chkCol = t.Accent
							chkGlyph = "●"
						}
						ui.Text(c, chkGlyph).FontSize(11).TextColor(chkCol)
					})
				})

				if row.Clicked() {
					items[idx].Selected = !items[idx].Selected
					if o.OnChange != nil {
						o.OnChange(items)
					}
				}
			}
			if shown == 0 {
				ui.Row(c).WidthPercent(100).Height(29).Center().Children(func() {
					ui.Text(c, "No matching members").FontSize(11).TextColor(t.TextSubtle)
				})
			}
		})
	})

	return card
}

type rosterDemoState struct {
	members []RosterMember
}

// ============================================================================
// Demo Registrations (14 blocks: 8 Type + 6 Select)
// ============================================================================

func init() {

	registerDemo(
		"command",
		"Compact Cmd+K command palette with live search filtering and click-to-run feedback.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			items := []CommandItem{
				{ID: "theme", Title: "Toggle theme", Shortcut: "⌘T"},
				{ID: "deploy", Title: "Deploy preview", Shortcut: "⌘D"},
				{ID: "export", Title: "Export tokens", Shortcut: "⌘E"},
				{ID: "branch", Title: "Switch branch", Shortcut: "⌘B"},
			}
			holder.Children(func() {
				CommandBar(c, items)
			})
		},
	)

	registerDemo(
		"picker",
		"Multi-select collaborator picker with initials avatars, checkmarks, and stack preview.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_picker", func() assigneePickerDemoState {
				return assigneePickerDemoState{
					members: []Assignee{
						{ID: "m1", Name: "Maya Lin", Initials: "ML", Color: ui.Hex("#6366f1"), Selected: true},
						{ID: "m2", Name: "Noah Sato", Initials: "NS", Color: ui.Hex("#0ea5e9"), Selected: true},
						{ID: "m3", Name: "Ava Brooks", Initials: "AB", Color: ui.Hex("#f59e0b"), Selected: false},
						{ID: "m4", Name: "Leo Park", Initials: "LP", Color: ui.Hex("#10b981"), Selected: false},
					},
				}
			})
			holder.Children(func() {
				AssigneePicker(c, st.members)
			})
		},
	)

	registerDemo(
		"aspect",
		"Aspect ratio selector with a live spring-morphing preview frame.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_aspect", func() aspectDemoState {
				return aspectDemoState{selected: 2}
			})
			holder.Children(func() {
				AspectRatioSelector(c, &st.selected)
			})
		},
	)

	registerDemo(
		"toolbar",
		"Floating studio canvas tool dock with sliding active pill and tool status badge.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_toolbar", func() canvasToolbarDemoState {
				return canvasToolbarDemoState{selected: 1}
			})
			holder.Children(func() {
				CanvasToolbar(c, &st.selected)
			})
		},
	)

	registerDemo(
		"icon-bar",
		"Segmented icon bar with smooth spring active indicator and tracking tooltip label.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_icon_bar", func() iconBarDemoState {
				return iconBarDemoState{selected: 2}
			})
			holder.Children(func() {
				IconBar(c, defaultIconBarItems(), &st.selected)
			})
		},
	)

	registerDemo(
		"roster",
		"Filterable team roster rows with role pills, status dots, and selection count badge.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_roster", func() rosterDemoState {
				return rosterDemoState{
					members: []RosterMember{
						{ID: "r1", Name: "Elena Cruz", Role: "Eng", Online: true, Selected: true},
						{ID: "r2", Name: "Samir Patel", Role: "Design", Online: true, Selected: false},
						{ID: "r3", Name: "Chloe Kim", Role: "Eng", Online: false, Selected: true},
						{ID: "r4", Name: "Owen Reed", Role: "PM", Online: true, Selected: false},
					},
				}
			})
			holder.Children(func() {
				RosterList(c, st.members)
			})
		},
	)
}
