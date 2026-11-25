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
// 3. pill-accordion ("Pill accordion")
// -----------------------------------------------------------------------------

// PillAccordionOptions configures PillAccordion.
type PillAccordionOptions struct {
	Width    float32
	Accent   ui.Color
	OnChange func(openIndex int)
}

// PillAccordion renders floating pill-shaped expandable items whose corner
// radius and padding smoothly morph when opened.
func PillAccordion(c *ui.Context, items []FAQItem, openIndex *int, opts ...PillAccordionOptions) ui.Element {
	var o PillAccordionOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 276
	}
	accent := colorOr(o.Accent, t.Accent)

	curOpen := -1
	if openIndex != nil {
		curOpen = *openIndex
	}

	root := ui.Column(c).
		Width(o.Width).
		Gap(7)

	root.Children(func() {
		for i := range items {
			idx := i
			item := items[idx]
			isOpen := curOpen == idx

			target := float32(0)
			if isOpen {
				target = 1
			}

			pill := ui.Column(c).
				Key(idx).
				WidthPercent(100).
				Clip().
				Cursor(ui.CursorPointer)

			p := clamp32(pill.Spring("pill_acc_p", target, ui.SpringOptions{
				Duration: 340 * time.Millisecond,
				Bounce:   0.22,
			}), 0, 1)

			rad := lerp32(18, 13, p)
			padY := lerp32(7, 10, p)
			padX := lerp32(14, 15, p)
			bg := lerpColor(t.Surface, t.SurfaceSubtle, p*0.45)
			borderCol := lerpColor(t.Border, accent.Alpha(0.60), p)

			pill.
				Padding(padY, padX).
				Radius(rad).
				Background(bg).
				Border(1, borderCol)

			pill.Children(func() {
				ui.Row(c).
					WidthPercent(100).
					Justify(ui.SpaceBetween).
					AlignItems(ui.Center).
					Gap(8).
					Children(func() {
						ui.Row(c).AlignItems(ui.Center).Gap(7).Children(func() {
							dotSize := lerp32(6, 7.5, p)
							ui.Box(c).
								Size(dotSize, dotSize).
								Radius(999).
								Background(lerpColor(t.BorderStrong, accent, p))
							ui.Text(c, item.Question).
								FontSize(11.5).
								FontWeight(600).
								TextColor(t.Text).
								SingleLine()
						})

						ui.Box(c).Size(14, 14).Draw(func(painter *ui.Painter, r ui.Rect) {
							drawChevron(painter, r, p*180, lerpColor(t.TextSubtle, accent, p))
						})
					})

				bodyH := p * 33
				if bodyH > 0.5 {
					bodyAlpha := clamp32((p-0.14)/0.86, 0, 1)
					ui.Box(c).
						WidthPercent(100).
						Height(bodyH).
						Clip().
						Opacity(bodyAlpha).
						Children(func() {
							ui.Box(c).Padding(5, 0, 0, 14).Children(func() {
								ui.Text(c, item.Answer).
									FontSize(10.5).
									LineHeight(1.35).
									TextColor(t.TextMuted)
							})
						})
				}
			})

			if pill.Clicked() {
				next := idx
				if isOpen {
					next = -1
				}
				if openIndex != nil {
					*openIndex = next
				}
				if o.OnChange != nil {
					o.OnChange(next)
				}
			}
		}
	})

	return root
}

// -----------------------------------------------------------------------------
// 4. tree-view ("Tree view")
// -----------------------------------------------------------------------------

// TreeNode describes a file or folder node in TreeView.
type TreeNode struct {
	ID       string
	Name     string
	IsFolder bool
	Children []TreeNode
}

// TreeViewOptions configures TreeView.
type TreeViewOptions struct {
	Width    float32
	Accent   ui.Color
	OnSelect func(id string)
}

func countVisibleTreeRows(nodes []TreeNode, expanded map[string]bool) int {
	total := 0
	for _, n := range nodes {
		total++
		if n.IsFolder && expanded[n.ID] {
			total += countVisibleTreeRows(n.Children, expanded)
		}
	}
	return total
}

func drawFolderIcon(p *ui.Painter, r ui.Rect, open bool, col ui.Color) {
	x := r.X + 1.5
	y := r.Y + 2.5
	w := r.W - 3
	h := r.H - 4.5
	// Top tab
	p.Fill(ui.Rect{X: x, Y: y - 1.2, W: w * 0.45, H: 2.5}, col.Alpha(0.85), 1)
	if open {
		p.Fill(ui.Rect{X: x, Y: y + 0.5, W: w, H: h}, col.Alpha(0.24), 2)
		p.Stroke(ui.Rect{X: x, Y: y + 0.5, W: w, H: h}, col, 2, 1.1)
	} else {
		p.Fill(ui.Rect{X: x, Y: y + 0.5, W: w, H: h}, col.Alpha(0.75), 2)
	}
}

func drawFileIcon(p *ui.Painter, r ui.Rect, col ui.Color) {
	x := r.X + 2.5
	y := r.Y + 1.5
	w := r.W - 5
	h := r.H - 3
	p.Stroke(ui.Rect{X: x, Y: y, W: w, H: h}, col, 1.5, 1.1)
	p.Line(x+2, y+4, x+w-2, y+4, 1, col.Alpha(0.7))
	p.Line(x+2, y+6.8, x+w-3, y+6.8, 1, col.Alpha(0.7))
}

// TreeView renders an interactive file/folder hierarchy tree with animated
// expand/collapse chevrons, indentation guide lines, and active file selection.
func TreeView(c *ui.Context, nodes []TreeNode, expanded map[string]bool, selectedID *string, opts ...TreeViewOptions) ui.Element {
	var o TreeViewOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 264
	}
	accent := colorOr(o.Accent, t.Accent)

	curSel := ""
	if selectedID != nil {
		curSel = *selectedID
	}

	const rowH float32 = 24
	const indentW float32 = 14

	card := ui.Column(c).
		Width(o.Width).
		Padding(8).
		Radius(12).
		Background(t.Surface).
		Border(1, t.Border)

	var renderNodes func(list []TreeNode, depth int)
	renderNodes = func(list []TreeNode, depth int) {
		for _, node := range list {
			n := node
			isOpen := n.IsFolder && expanded[n.ID]
			isSel := !n.IsFolder && curSel == n.ID

			branch := ui.Column(c).Key(n.ID).WidthPercent(100)
			branch.Children(func() {
				row := ui.Row(c).
					WidthPercent(100).
					Height(rowH).
					Padding(0, 6).
					AlignItems(ui.Center).
					Radius(6).
					Cursor(ui.CursorPointer)

				if isSel {
					row.Background(t.AccentSoft)
				} else if row.Hovered() {
					row.Background(t.SurfaceSubtle)
				}

				row.Draw(func(p *ui.Painter, r ui.Rect) {
					// Draw vertical indentation guide lines for ancestor depths.
					for d := 0; d < depth; d++ {
						gx := r.X + 11 + float32(d)*indentW
						p.Line(gx, r.Y, gx, r.Y+r.H, 1, t.Border)
					}
				})

				chevTarget := float32(0)
				if isOpen {
					chevTarget = 1
				}
				chevP := clamp32(row.Spring("tree_chev", chevTarget, ui.SpringOptions{
					Duration: 240 * time.Millisecond,
					Bounce:   0.18,
				}), 0, 1)

				row.Children(func() {
					if depth > 0 {
						ui.Box(c).Width(float32(depth) * indentW).Height(1)
					}
					if n.IsFolder {
						ui.Box(c).Size(12, 12).Margin(0, 4, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
							// Rotate from right-pointing (-90deg) to down-pointing (0deg).
							drawChevron(p, r, lerp32(-90, 0, chevP), t.TextSubtle)
						})
						ui.Box(c).Size(13, 13).Margin(0, 6, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
							drawFolderIcon(p, r, isOpen, lerpColor(t.TextMuted, accent, chevP*0.65))
						})
					} else {
						ui.Box(c).Size(12, 12).Margin(0, 4, 0, 0)
						ui.Box(c).Size(13, 13).Margin(0, 6, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
							ic := t.TextSubtle
							if isSel {
								ic = accent
							}
							drawFileIcon(p, r, ic)
						})
					}

					txtCol := t.Text
					weight := 400
					if n.IsFolder {
						weight = 500
					}
					if isSel {
						txtCol = accent
						weight = 600
					}
					ui.Text(c, n.Name).
						FontSize(11.5).
						FontWeight(weight).
						TextColor(txtCol).
						SingleLine()
				})

				if row.Clicked() {
					if n.IsFolder {
						if expanded != nil {
							expanded[n.ID] = !expanded[n.ID]
						}
					} else {
						if selectedID != nil {
							*selectedID = n.ID
						}
						if o.OnSelect != nil {
							o.OnSelect(n.ID)
						}
					}
				}

				if n.IsFolder && len(n.Children) > 0 {
					targetH := float32(0)
					if isOpen {
						targetH = float32(countVisibleTreeRows(n.Children, expanded)) * rowH
					}
					childWrap := ui.Column(c).
						Key("children").
						WidthPercent(100).
						Clip()
					animH := clamp32(childWrap.Spring("tree_h", targetH, ui.SpringOptions{
						Duration: 280 * time.Millisecond,
						Bounce:   0.12,
					}), 0, 400)
					if animH > 0.5 {
						alpha := float32(1)
						if targetH > 0 {
							alpha = clamp32(animH/targetH, 0, 1)
						} else {
							alpha = clamp32(animH/rowH, 0, 1)
						}
						childWrap.Height(animH).Opacity(alpha).Children(func() {
							renderNodes(n.Children, depth+1)
						})
					}
				}
			})
		}
	}

	card.Children(func() {
		renderNodes(nodes, 0)
	})

	return card
}

// -----------------------------------------------------------------------------
// 5. collapsing-sidebar ("Collapsing sidebar")
// -----------------------------------------------------------------------------

// SidebarNavItem describes one navigation entry in CollapsingSidebar.
type SidebarNavItem struct {
	ID    string
	Icon  string
	Label string
	Badge string
}

// CollapsingSidebarOptions configures CollapsingSidebar.
type CollapsingSidebarOptions struct {
	Width    float32
	Height   float32
	Accent   ui.Color
	OnSelect func(id string)
}

func drawSidebarGlyph(p *ui.Painter, r ui.Rect, kind string, col ui.Color) {
	cx := r.X + r.W*0.5
	cy := r.Y + r.H*0.5
	switch kind {
	case "grid":
		p.Stroke(ui.Rect{X: cx - 4.5, Y: cy - 4.5, W: 3.8, H: 3.8}, col, 1, 1.2)
		p.Stroke(ui.Rect{X: cx + 0.8, Y: cy - 4.5, W: 3.8, H: 3.8}, col, 1, 1.2)
		p.Stroke(ui.Rect{X: cx - 4.5, Y: cy + 0.8, W: 3.8, H: 3.8}, col, 1, 1.2)
		p.Stroke(ui.Rect{X: cx + 0.8, Y: cy + 0.8, W: 3.8, H: 3.8}, col, 1, 1.2)
	case "inbox":
		p.Stroke(ui.Rect{X: cx - 4.8, Y: cy - 3.8, W: 9.6, H: 7.6}, col, 1.6, 1.2)
		p.StrokePolyline([][2]float32{
			{cx - 4.2, cy - 1.5},
			{cx, cy + 1.2},
			{cx + 4.2, cy - 1.5},
		}, 1.2, col)
	case "tasks":
		p.Stroke(ui.Rect{X: cx - 4.5, Y: cy - 4.5, W: 9, H: 9}, col, 2, 1.2)
		p.StrokePolyline([][2]float32{
			{cx - 2.3, cy + 0.2},
			{cx - 0.5, cy + 2.0},
			{cx + 2.6, cy - 1.8},
		}, 1.3, col)
	case "toggle":
		p.Stroke(ui.Rect{X: cx - 5, Y: cy - 4, W: 10, H: 8}, col, 1.8, 1.2)
		p.Line(cx-1.6, cy-4, cx-1.6, cy+4, 1.2, col)
	default:
		// sliders / settings
		p.Line(cx-4.5, cy-2.4, cx+4.5, cy-2.4, 1.2, col)
		p.Fill(ui.Rect{X: cx - 2.2, Y: cy - 4.0, W: 3.2, H: 3.2}, col, 1.6)
		p.Line(cx-4.5, cy+2.4, cx+4.5, cy+2.4, 1.2, col)
		p.Fill(ui.Rect{X: cx + 0.4, Y: cy + 0.8, W: 3.2, H: 3.2}, col, 1.6)
	}
}

// CollapsingSidebar renders a mini workspace layout with a toggleable sidebar
// that springs smoothly between expanded (icons + labels) and rail (icons only).
func CollapsingSidebar(c *ui.Context, items []SidebarNavItem, collapsed *bool, activeID *string, opts ...CollapsingSidebarOptions) ui.Element {
	var o CollapsingSidebarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 292
	}
	if o.Height <= 0 {
		o.Height = 172
	}
	accent := colorOr(o.Accent, t.Accent)

	isCollapsed := false
	if collapsed != nil {
		isCollapsed = *collapsed
	}
	curActive := ""
	if activeID != nil {
		curActive = *activeID
	} else if len(items) > 0 {
		curActive = items[0].ID
	}

	const expandedW float32 = 114
	const railW float32 = 42

	targetW := expandedW
	if isCollapsed {
		targetW = railW
	}

	frame := ui.Row(c).
		Width(o.Width).
		Height(o.Height).
		Radius(12).
		Background(t.Surface).
		Border(1, t.Border).
		Clip()

	sbW := clamp32(frame.Spring("sidebar_w", targetW, ui.SpringOptions{
		Duration: 340 * time.Millisecond,
		Bounce:   0.16,
	}), railW, expandedW)

	expandFrac := clamp32((sbW-railW)/(expandedW-railW), 0, 1)
	labelAlpha := clamp32((expandFrac-0.25)/0.75, 0, 1)

	activeLabel := "Workspace"
	for _, it := range items {
		if it.ID == curActive {
			activeLabel = it.Label
			break
		}
	}

	frame.Children(func() {
		// Left sidebar
		sb := ui.Column(c).
			Width(sbW).
			HeightPercent(100).
			Padding(8, 7).
			Gap(4).
			Background(t.SurfaceSubtle).
			BorderWidth(0, 1, 0, 0).
			BorderColor(t.Border).
			Clip()

		sb.Children(func() {
			// Header row with brand dot + collapse toggle
			hdr := ui.Row(c).
				WidthPercent(100).
				Height(24).
				Padding(0, 3).
				Justify(ui.SpaceBetween).
				AlignItems(ui.Center)
			hdr.Children(func() {
				if expandFrac > 0.35 {
					ui.Row(c).AlignItems(ui.Center).Gap(5).Opacity(labelAlpha).Children(func() {
						ui.Box(c).Size(7, 7).Radius(999).Background(accent)
						ui.Text(c, "Studio").FontSize(11).FontWeight(700).TextColor(t.Text).SingleLine()
					})
				}
				toggleBtn := ui.Box(c).
					Key("sb_toggle").
					Size(22, 22).
					Radius(6).
					Center().
					Cursor(ui.CursorPointer)
				if toggleBtn.Hovered() {
					toggleBtn.Background(t.SurfaceHover)
				}
				toggleBtn.Draw(func(p *ui.Painter, r ui.Rect) {
					drawSidebarGlyph(p, r, "toggle", t.TextMuted)
				})
				if toggleBtn.Clicked() && collapsed != nil {
					*collapsed = !*collapsed
				}
			})

			// Nav items
			for i := range items {
				it := items[i]
				isSel := it.ID == curActive

				nav := ui.Row(c).
					Key(it.ID).
					WidthPercent(100).
					Height(28).
					Padding(0, 6).
					AlignItems(ui.Center).
					Gap(8).
					Radius(7).
					Clip().
					Cursor(ui.CursorPointer)

				if isSel {
					nav.Background(t.AccentSoft)
				} else if nav.Hovered() {
					nav.Background(t.SurfaceHover)
				}

				nav.Children(func() {
					icCol := t.TextMuted
					if isSel {
						icCol = accent
					}
					ui.Box(c).Size(15, 15).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
						drawSidebarGlyph(p, r, it.Icon, icCol)
					})
					if expandFrac > 0.12 {
						ui.Row(c).
							Grow(1).
							Justify(ui.SpaceBetween).
							AlignItems(ui.Center).
							Opacity(labelAlpha).
							Children(func() {
								txtCol := t.Text
								weight := 500
								if isSel {
									txtCol = accent
									weight = 600
								}
								ui.Text(c, it.Label).
									FontSize(11).
									FontWeight(weight).
									TextColor(txtCol).
									SingleLine()
								if it.Badge != "" {
									ui.Box(c).
										Padding(1, 5).
										Radius(999).
										Background(t.Surface).
										Border(1, t.Border).
										Children(func() {
											ui.Text(c, it.Badge).FontSize(9.5).FontWeight(600).TextColor(t.TextMuted)
										})
								}
							})
					}
				})

				if nav.Clicked() {
					if activeID != nil {
						*activeID = it.ID
					}
					if o.OnSelect != nil {
						o.OnSelect(it.ID)
					}
				}
			}
		})

		// Main workspace canvas
		main := ui.Column(c).
			Grow(1).
			HeightPercent(100).
			Padding(11).
			Gap(8).
			Clip()
		main.Children(func() {
			ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
				ui.Text(c, activeLabel).FontSize(12.5).FontWeight(700).TextColor(t.Text).SingleLine()
				modeTxt := "Full"
				if isCollapsed {
					modeTxt = "Rail"
				}
				ui.Box(c).Padding(2, 6).Radius(999).Background(t.SurfaceSubtle).Border(1, t.Border).Children(func() {
					ui.Text(c, modeTxt).FontSize(9.5).FontWeight(600).TextColor(t.TextSubtle)
				})
			})
			// Two compact content skeleton cards that respond to width changes
			ui.Row(c).WidthPercent(100).Gap(6).Children(func() {
				ui.Column(c).Grow(1).Height(46).Padding(7).Radius(8).Background(t.SurfaceSubtle).Border(1, t.Border).Justify(ui.SpaceBetween).Children(func() {
					ui.Box(c).WidthPercent(55).Height(5).Radius(3).Background(accent.Alpha(0.55))
					ui.Box(c).WidthPercent(85).Height(5).Radius(3).Background(t.BorderStrong)
				})
				ui.Column(c).Grow(1).Height(46).Padding(7).Radius(8).Background(t.SurfaceSubtle).Border(1, t.Border).Justify(ui.SpaceBetween).Children(func() {
					ui.Box(c).WidthPercent(45).Height(5).Radius(3).Background(t.Success.Alpha(0.60))
					ui.Box(c).WidthPercent(75).Height(5).Radius(3).Background(t.BorderStrong)
				})
			})
			ui.Column(c).WidthPercent(100).Grow(1).Padding(8).Radius(8).Background(t.SurfaceSubtle).Border(1, t.Border).Gap(6).Children(func() {
				ui.Box(c).WidthPercent(70).Height(5).Radius(3).Background(t.BorderStrong)
				ui.Box(c).WidthPercent(92).Height(5).Radius(3).Background(t.Border)
				ui.Box(c).WidthPercent(48).Height(5).Radius(3).Background(t.Border)
			})
		})
	})

	return frame
}

// -----------------------------------------------------------------------------
// 6. fold-dropdown ("Dropdown")
// -----------------------------------------------------------------------------

// FoldDropdownOptions configures FoldDropdown.
type FoldDropdownOptions struct {
	Label    string
	Width    float32
	Accent   ui.Color
	OnSelect func(index int)
}

// FoldDropdown renders a trigger button that unfolds an interactive option menu
// with a checkmarked active choice and spring reveal.
func FoldDropdown(c *ui.Context, options []string, selected *int, open *bool, opts ...FoldDropdownOptions) ui.Element {
	var o FoldDropdownOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 224
	}
	accent := colorOr(o.Accent, t.Accent)

	sel := 0
	if selected != nil && *selected >= 0 && *selected < len(options) {
		sel = *selected
	}
	isOpen := false
	if open != nil {
		isOpen = *open
	}

	curLabel := "Choose..."
	if len(options) > 0 {
		curLabel = options[sel]
	}

	const itemH float32 = 30
	fullMenuH := float32(len(options))*itemH + 10

	root := ui.Column(c).
		Width(o.Width).
		Height(36 + 6 + fullMenuH)

	openTarget := float32(0)
	if isOpen {
		openTarget = 1
	}
	openP := clamp32(root.Spring("fold_dd_p", openTarget, ui.SpringOptions{
		Duration: 320 * time.Millisecond,
		Bounce:   0.20,
	}), 0, 1)

	root.Children(func() {
		trigger := ui.Row(c).
			Key("trigger").
			WidthPercent(100).
			Height(36).
			Padding(0, 12).
			Radius(10).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Background(t.Surface).
			Border(1, lerpColor(t.BorderStrong, accent.Alpha(0.75), openP)).
			Cursor(ui.CursorPointer)

		trigger.Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
				if o.Label != "" {
					ui.Text(c, o.Label).FontSize(11).TextColor(t.TextSubtle)
				}
				ui.Text(c, curLabel).FontSize(12).FontWeight(600).TextColor(t.Text).SingleLine()
			})
			ui.Box(c).Size(14, 14).Draw(func(p *ui.Painter, r ui.Rect) {
				drawChevron(p, r, openP*180, lerpColor(t.TextMuted, accent, openP))
			})
		})

		if trigger.Clicked() && open != nil {
			*open = !*open
		}

		curMenuH := openP * fullMenuH
		if curMenuH > 1 {
			menu := ui.Column(c).
				Key("menu").
				WidthPercent(100).
				Height(curMenuH).
				Margin(6, 0, 0, 0).
				Padding(5).
				Radius(11).
				Background(t.Surface).
				Border(1, t.BorderStrong).
				Shadow(0, 6, 16, 0, t.Shadow).
				Opacity(clamp32(openP*1.35, 0, 1)).
				Translate(0, lerp32(-6, 0, openP)).
				Clip().
				IgnoreInput(!isOpen)

			menu.Children(func() {
				for i := range options {
					idx := i
					opt := options[idx]
					isSel := idx == sel

					row := ui.Row(c).
						Key(idx).
						WidthPercent(100).
						Height(itemH).
						Padding(0, 9).
						Radius(7).
						Justify(ui.SpaceBetween).
						AlignItems(ui.Center).
						Cursor(ui.CursorPointer)

					if isSel {
						row.Background(t.AccentSoft)
					} else if row.Hovered() {
						row.Background(t.SurfaceSubtle)
					}

					row.Children(func() {
						txtCol := t.Text
						weight := 500
						if isSel {
							txtCol = accent
							weight = 600
						}
						ui.Text(c, opt).FontSize(11.5).FontWeight(weight).TextColor(txtCol).SingleLine()
						if isSel {
							ui.Box(c).Size(14, 14).Draw(func(p *ui.Painter, r ui.Rect) {
								drawCheckmark(p, r, accent)
							})
						}
					})

					if row.Clicked() {
						if selected != nil {
							*selected = idx
						}
						if open != nil {
							*open = false
						}
						if o.OnSelect != nil {
							o.OnSelect(idx)
						}
					}
				}
			})
		}
	})

	return root
}

// -----------------------------------------------------------------------------
// 7. radial ("Radial menu")
// -----------------------------------------------------------------------------

// RadialAction defines one fan-out action node in RadialMenu.
type RadialAction struct {
	ID    string
	Label string
	Short string
	Color ui.Color
}

// RadialMenuOptions configures RadialMenu.
type RadialMenuOptions struct {
	Radius   float32
	Accent   ui.Color
	OnSelect func(index int)
}

// RadialMenu renders a circular center trigger button that fans out 5-6 action
// nodes around a ring using spring polar interpolation on click.
func RadialMenu(c *ui.Context, actions []RadialAction, open *bool, selected *int, opts ...RadialMenuOptions) ui.Element {
	var o RadialMenuOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Radius <= 0 {
		o.Radius = 56
	}
	accent := colorOr(o.Accent, t.Accent)

	isOpen := false
	if open != nil {
		isOpen = *open
	}
	curSel := -1
	if selected != nil {
		curSel = *selected
	}

	const stageW float32 = 196
	const stageH float32 = 166
	cx := stageW * 0.5
	cy := stageH * 0.5
	const centerSize float32 = 42
	const nodeSize float32 = 34

	holder := ui.Column(c).AlignItems(ui.Center).Gap(4)

	holder.Children(func() {
		stage := ui.Box(c).Size(stageW, stageH)

		openTarget := float32(0)
		if isOpen {
			openTarget = 1
		}
		ringP := clamp32(stage.Spring("radial_ring", openTarget, ui.SpringOptions{
			Duration: 360 * time.Millisecond,
			Bounce:   0.22,
		}), 0, 1)

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if ringP <= 0.02 {
				return
			}
			rad := o.Radius * ringP
			p.Stroke(
				ui.Rect{X: r.X + cx - rad, Y: r.Y + cy - rad, W: rad * 2, H: rad * 2},
				t.BorderStrong.Alpha(0.45*ringP),
				rad,
				1,
			)
		})

		stage.Children(func() {
			n := len(actions)
			for i := range actions {
				idx := i
				act := actions[idx]
				isSel := curSel == idx
				nodeCol := colorOr(act.Color, accent)

				node := ui.Box(c).
					Key(idx).
					Absolute().
					Size(nodeSize, nodeSize).
					Radius(nodeSize * 0.5).
					Center().
					Cursor(ui.CursorPointer).
					IgnoreInput(!isOpen)

				p := clamp32(node.Spring("radial_node", openTarget, ui.SpringOptions{
					Duration: time.Duration(300+idx*28) * time.Millisecond,
					Bounce:   0.26,
				}), 0, 1)

				angle := -math.Pi*0.5 + float64(idx)*(2*math.Pi/float64(max(n, 1))) - float64(1-p)*0.55
				dist := o.Radius * p
				nx := cx + float32(math.Cos(angle))*dist - nodeSize*0.5
				ny := cy + float32(math.Sin(angle))*dist - nodeSize*0.5

				bg := t.Surface
				borderCol := t.BorderStrong
				txtCol := t.Text
				if isSel {
					bg = nodeCol
					borderCol = nodeCol
					txtCol = ui.Hex("#ffffff")
				} else if node.Hovered() {
					bg = t.SurfaceHover
					borderCol = nodeCol.Alpha(0.75)
				}

				sc := lerp32(0.35, 1.0, p)
				node.
					Left(nx).
					Top(ny).
					Scale(sc, sc).
					Opacity(clamp32(p*1.4, 0, 1)).
					Background(bg).
					Border(1, borderCol).
					Shadow(0, 3, 8, 0, t.Shadow)

				node.Children(func() {
					ui.Text(c, act.Short).FontSize(10.5).FontWeight(700).TextColor(txtCol)
				})

				if node.Clicked() {
					if selected != nil {
						*selected = idx
					}
					if open != nil {
						*open = false
					}
					if o.OnSelect != nil {
						o.OnSelect(idx)
					}
				}
			}

			// Center trigger button
			centerBtn := ui.Box(c).
				Key("radial_center").
				Absolute().
				Left(cx-centerSize*0.5).
				Top(cy-centerSize*0.5).
				Size(centerSize, centerSize).
				Radius(centerSize*0.5).
				Center().
				Background(lerpColor(t.Text, accent, ringP)).
				Shadow(0, 4, 12, 0, t.Shadow).
				Cursor(ui.CursorPointer)

			pressTarget := float32(1)
			if centerBtn.Pressed() {
				pressTarget = 0.92
			}
			sc := centerBtn.Spring("radial_sc", pressTarget, ui.SpringOptions{
				Duration: 180 * time.Millisecond,
				Bounce:   0.25,
			})
			centerBtn.Scale(sc, sc).Draw(func(p *ui.Painter, r ui.Rect) {
				ink := lerpColor(t.Background, ui.Hex("#ffffff"), ringP)
				drawPlusRotated(p, r, ringP*45, 6.5, ink)
			})

			if centerBtn.Clicked() && open != nil {
				*open = !*open
			}
		})

		status := "Click + to fan out"
		statusCol := t.TextSubtle
		if curSel >= 0 && curSel < len(actions) {
			status = "Selected: " + actions[curSel].Label
			statusCol = colorOr(actions[curSel].Color, accent)
		} else if isOpen {
			status = "Pick an action"
			statusCol = t.TextMuted
		}
		ui.Text(c, status).FontSize(11).FontWeight(500).TextColor(statusCol)
	})

	return holder
}

// -----------------------------------------------------------------------------
// 8. liq-create ("Create menu")
// -----------------------------------------------------------------------------

// CreateMenuItem describes one action inside CreateMenu.
type CreateMenuItem struct {
	ID    string
	Label string
	Hint  string
	Color ui.Color
}

// CreateMenuOptions configures CreateMenu.
type CreateMenuOptions struct {
	Accent   ui.Color
	OnCreate func(item CreateMenuItem)
}

// CreateMenu renders a "+" pill button that smoothly morphs into a multi-action
// Create menu (Note, Task, Folder, Link) and collapses cleanly on selection.
func CreateMenu(c *ui.Context, items []CreateMenuItem, open *bool, opts ...CreateMenuOptions) ui.Element {
	var o CreateMenuOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	accent := colorOr(o.Accent, t.Accent)

	isOpen := false
	if open != nil {
		isOpen = *open
	}

	openTarget := float32(0)
	if isOpen {
		openTarget = 1
	}

	const closedW float32 = 104
	const closedH float32 = 36
	const openW float32 = 228
	openH := float32(34 + len(items)*30 + 10)

	morph := ui.Column(c).Clip()
	p := clamp32(morph.Spring("liq_create_p", openTarget, ui.SpringOptions{
		Duration: 360 * time.Millisecond,
		Bounce:   0.22,
	}), 0, 1)

	w := lerp32(closedW, openW, p)
	h := lerp32(closedH, openH, p)
	rad := lerp32(18, 14, p)
	bg := lerpColor(t.Text, t.Surface, p)
	borderCol := lerpColor(t.Text, t.BorderStrong, p)

	morph.
		Width(w).
		Height(h).
		Radius(rad).
		Background(bg).
		Border(1, borderCol).
		Shadow(0, 4*p+1, 14*p+2, 0, t.Shadow)

	morph.Children(func() {
		if p < 0.30 {
			pillBtn := ui.Row(c).
				Key("create_pill").
				WidthPercent(100).
				HeightPercent(100).
				Center().
				Gap(6).
				Opacity(clamp32(1-p/0.30, 0, 1)).
				Cursor(ui.CursorPointer)
			pillBtn.Children(func() {
				ui.Box(c).Size(12, 12).Draw(func(painter *ui.Painter, r ui.Rect) {
					drawPlusRotated(painter, r, 0, 4.5, t.Background)
				})
				ui.Text(c, "Create").FontSize(12.5).FontWeight(600).TextColor(t.Background)
			})
			if pillBtn.Clicked() && open != nil {
				*open = true
			}
		} else {
			contentAlpha := clamp32((p-0.25)/0.75, 0, 1)
			ui.Column(c).
				Key("create_body").
				WidthPercent(100).
				HeightPercent(100).
				Padding(6).
				Gap(2).
				Opacity(contentAlpha).
				Children(func() {
					// Menu header
					ui.Row(c).
						WidthPercent(100).
						Height(26).
						Padding(0, 6).
						Justify(ui.SpaceBetween).
						AlignItems(ui.Center).
						Children(func() {
							ui.Text(c, "Create new").FontSize(11).FontWeight(600).TextColor(t.TextSubtle)
							closeBtn := ui.Box(c).
								Key("close_create").
								Size(20, 20).
								Radius(10).
								Center().
								Cursor(ui.CursorPointer)
							if closeBtn.Hovered() {
								closeBtn.Background(t.SurfaceHover)
							}
							closeBtn.Draw(func(painter *ui.Painter, r ui.Rect) {
								drawPlusRotated(painter, r, 45, 4.2, t.TextMuted)
							})
							if closeBtn.Clicked() && open != nil {
								*open = false
							}
						})

					for i := range items {
						it := items[i]
						col := colorOr(it.Color, accent)
						row := ui.Row(c).
							Key(it.ID).
							WidthPercent(100).
							Height(28).
							Padding(0, 7).
							Radius(8).
							Justify(ui.SpaceBetween).
							AlignItems(ui.Center).
							Cursor(ui.CursorPointer)
						if row.Hovered() {
							row.Background(t.SurfaceSubtle)
						}
						row.Children(func() {
							ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
								ui.Box(c).
									Size(16, 16).
									Radius(5).
									Center().
									Background(col.Alpha(0.18)).
									Children(func() {
										ui.Box(c).Size(6, 6).Radius(999).Background(col)
									})
								ui.Text(c, it.Label).FontSize(12).FontWeight(500).TextColor(t.Text)
							})
							if it.Hint != "" {
								ui.Box(c).Padding(1, 5).Radius(4).Background(t.SurfaceSubtle).Border(1, t.Border).Children(func() {
									ui.Text(c, it.Hint).FontSize(9.5).FontWeight(600).TextColor(t.TextSubtle)
								})
							}
						})
						if row.Clicked() {
							if open != nil {
								*open = false
							}
							if o.OnCreate != nil {
								o.OnCreate(it)
							}
						}
					}
				})
		}
	})

	return morph
}

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
	sampleFAQs := []FAQItem{
		{
			Question: "How do springs stay 60fps?",
			Answer:   "Each transition integrates analytically on GPU frames without layout thrashing.",
		},
		{
			Question: "Does it support Dark mode?",
			Answer:   "Every token resolves dynamically from ThemeFrom(c) for light and dark surfaces.",
		},
		{
			Question: "Can I compose custom themes?",
			Answer:   "Pass custom Accent colors or swap context theme tokens at runtime.",
		},
	}

	// 1. faq-frames

	// 2. faq-accordion

	// 3. pill-accordion
	registerDemo(
		"pill-accordion",
		"Floating pill-shaped expandable items whose corner radius and padding morph when opened.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_pill_accordion", func() faqDemoState {
				return faqDemoState{open: 0}
			})
			holder.Children(func() {
				PillAccordion(c, sampleFAQs, &st.open)
			})
		},
	)

	// 4. tree-view
	treeNodes := []TreeNode{
		{
			ID:       "src",
			Name:     "src",
			IsFolder: true,
			Children: []TreeNode{
				{ID: "theme", Name: "theme.go"},
				{
					ID:       "ui",
					Name:     "components",
					IsFolder: true,
					Children: []TreeNode{
						{ID: "spring", Name: "spring.go"},
						{ID: "menu", Name: "menu.go"},
					},
				},
			},
		},
		{ID: "readme", Name: "README.md"},
	}
	registerDemo(
		"tree-view",
		"Interactive file/folder tree with animated chevrons, guide lines, and active file highlight.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_tree_view", func() treeDemoState {
				return treeDemoState{
					expanded: map[string]bool{"src": true, "ui": true},
					selected: "spring",
				}
			})
			holder.Children(func() {
				TreeView(c, treeNodes, st.expanded, &st.selected)
			})
		},
	)

	// 5. collapsing-sidebar
	sidebarItems := []SidebarNavItem{
		{ID: "overview", Icon: "grid", Label: "Overview"},
		{ID: "inbox", Icon: "inbox", Label: "Inbox", Badge: "4"},
		{ID: "tasks", Icon: "tasks", Label: "Tasks"},
		{ID: "settings", Icon: "sliders", Label: "Settings"},
	}
	registerDemo(
		"collapsing-sidebar",
		"Mini workspace with a sidebar that springs between expanded and icon-rail widths.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_collapsing_sidebar", func() sidebarDemoState {
				return sidebarDemoState{collapsed: false, active: "overview"}
			})
			holder.Children(func() {
				CollapsingSidebar(c, sidebarItems, &st.collapsed, &st.active)
			})
		},
	)

	// 6. fold-dropdown
	dropdownOptions := []string{"Spring physics", "Cubic ease", "Linear glide", "Instant snap"}
	registerDemo(
		"fold-dropdown",
		"Trigger button that unfolds an interactive option menu with spring reveal and checkmark.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_fold_dropdown", func() dropdownDemoState {
				return dropdownDemoState{selected: 0, open: true}
			})
			holder.Children(func() {
				FoldDropdown(c, dropdownOptions, &st.selected, &st.open, FoldDropdownOptions{
					Label: "Motion:",
				})
			})
		},
	)

	// 7. radial
	radialActions := []RadialAction{
		{ID: "share", Label: "Share", Short: "SH", Color: ui.Hex("#5865ff")},
		{ID: "copy", Label: "Copy", Short: "CP", Color: ui.Hex("#0ea5e9")},
		{ID: "star", Label: "Favorite", Short: "ST", Color: ui.Hex("#f59e0b")},
		{ID: "pin", Label: "Pin", Short: "PN", Color: ui.Hex("#16a34a")},
		{ID: "edit", Label: "Edit", Short: "ED", Color: ui.Hex("#a855f7")},
		{ID: "trash", Label: "Delete", Short: "DL", Color: ui.Hex("#ef4444")},
	}
	registerDemo(
		"radial",
		"Center circular trigger that fans out 6 action nodes with spring polar interpolation.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_radial", func() radialDemoState {
				return radialDemoState{open: true, selected: -1}
			})
			holder.Children(func() {
				RadialMenu(c, radialActions, &st.open, &st.selected)
			})
		},
	)

	// 8. liq-create
	createItems := []CreateMenuItem{
		{ID: "note", Label: "Note", Hint: "N", Color: ui.Hex("#5865ff")},
		{ID: "task", Label: "Task", Hint: "T", Color: ui.Hex("#16a34a")},
		{ID: "folder", Label: "Folder", Hint: "F", Color: ui.Hex("#f59e0b")},
		{ID: "link", Label: "Link", Hint: "L", Color: ui.Hex("#0ea5e9")},
	}
	registerDemo(
		"liq-create",
		"Morphing + Create pill that expands into a multi-action menu and collapses on selection.",
		func(c *ui.Context) {
			holder := ui.Column(c).AlignItems(ui.Center).Gap(6)
			st := ui.Local(holder, "demo_liq_create", func() createMenuDemoState {
				return createMenuDemoState{open: false}
			})
			holder.Children(func() {
				CreateMenu(c, createItems, &st.open, CreateMenuOptions{
					OnCreate: func(it CreateMenuItem) {
						st.lastCreated = "Created " + it.Label
					},
				})
				if st.lastCreated != "" {
					t := ThemeFrom(c)
					ui.Text(c, st.lastCreated).FontSize(10.5).FontWeight(500).TextColor(t.Success)
				}
			})
		},
	)

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
