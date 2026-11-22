package flow

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
)

// ZoomSelect renders React Flow UI's ZoomSelect dropdown anchored above the zoom button.
func ZoomSelect(c *ui.Context, g *Graph, t Theme) {
	levels := []struct {
		label string
		val   float32
	}{
		{"Best Fit", 0},
		{"50%", 0.50},
		{"75%", 0.75},
		{"100%", 1.00},
		{"150%", 1.50},
		{"200%", 2.00},
	}
	ui.Column(c).
		Width(104).
		Background(t.ControlsBG).
		Border(1, t.ControlsBorder).
		Radius(9).
		Padding(4).
		Gap(2).
		Shadow(0, 8, 24, 0, colorWithAlpha(ui.RGB(0, 0, 0), 0.22)).
		Children(func() {
			for _, item := range levels {
				lv := item
				row := ui.ButtonBase(c).
					WidthPercent(100).
					Height(26).
					Padding(0, 8).
					Radius(6).
					AlignItems(ui.Center).
					Justify(ui.SpaceBetween).
					Cursor(ui.CursorPointer)
				if row.Hovered() {
					row.Background(t.BadgeBG)
				}
				row.Children(func() {
					ui.Text(c, lv.label).FontSize(11.5).FontWeight(500).TextColor(t.ControlsText)
				})
				if row.Clicked() {
					if lv.val == 0 {
						g.FitView(0, 0, 56)
					} else {
						g.ZoomTo(lv.val)
					}
					g.zoomSelectOpen = false
				}
			}
		})
}

// NodeSearch renders React Flow UI's NodeSearch command palette overlay (Cmd/Ctrl+K).
// Typing filters nodes by title/id/subtitle; clicking any result focuses & centers it.
func NodeSearch(c *ui.Context, g *Graph, t Theme) {
	q := strings.ToLower(strings.TrimSpace(g.SearchQuery))
	var matches []*Node
	for _, n := range g.Nodes {
		if q == "" ||
			strings.Contains(strings.ToLower(n.Title), q) ||
			strings.Contains(strings.ToLower(n.ID), q) ||
			strings.Contains(strings.ToLower(n.Subtitle), q) {
			matches = append(matches, n)
		}
	}

	ui.Column(c).
		Width(296).
		Background(t.ControlsBG).
		Border(1, t.ControlsBorder).
		Radius(12).
		Shadow(0, 14, 34, 0, colorWithAlpha(ui.RGB(0, 0, 0), 0.28)).
		Clip().
		Children(func() {
			ui.Row(c).
				WidthPercent(100).
				Padding(8, 10).
				Gap(8).
				AlignItems(ui.Center).
				BorderWidth(0, 0, 1, 0).
				BorderColor(t.ControlsBorder).
				Children(func() {
					ui.TextInput(c, &g.SearchQuery).
						Placeholder("Search nodes (Esc to close)…").
						Grow(1).
						MinWidth(0).
						Height(28).
						Padding(0, 8).
						Radius(6).
						Background(t.NodeHeaderBG).
						Border(1, t.ControlsBorder).
						TextColor(t.TextPrimary)
					ctrlBtn(c, t, "×", "Close search", false, func() {
						g.ShowSearch = false
					})
				})

			ui.Scroll(c).WidthPercent(100).MaxHeight(220).Children(func() {
				ui.Column(c).WidthPercent(100).Padding(6).Gap(3).Children(func() {
					if len(matches) == 0 {
						ui.Box(c).WidthPercent(100).Padding(16, 8).Center().Children(func() {
							ui.Text(c, "No matching nodes").FontSize(12).TextColor(t.TextMuted)
						})
						return
					}
					for _, n := range matches {
						node := n
						row := ui.ButtonBase(c).
							WidthPercent(100).
							Padding(7, 9).
							Radius(7).
							AlignItems(ui.Center).
							Justify(ui.SpaceBetween).
							Cursor(ui.CursorPointer)
						if row.Hovered() {
							row.Background(t.BadgeBG)
						}
						row.Children(func() {
							ui.Column(c).Gap(1).Grow(1).MinWidth(0).Children(func() {
								ui.Text(c, node.Title).FontSize(12).FontWeight(500).TextColor(t.TextPrimary).SingleLine()
								if node.Subtitle != "" {
									ui.Text(c, node.Subtitle).FontSize(10.5).TextColor(t.TextMuted).SingleLine().Ellipsis("…")
								}
							})
							ui.Text(c, node.ID).FontSize(10).TextColor(t.TextMuted)
						})
						if row.Clicked() {
							g.FocusNode(node.ID)
							g.ShowSearch = false
						}
					}
				})
			})
		})
}

// DevTools renders React Flow UI's live DevTools overlay (ViewportLogger,
// NodeInspector, and ChangeLogger).
func DevTools(c *ui.Context, g *Graph, t Theme) {
	ui.Column(c).
		Width(248).
		Background(t.ControlsBG).
		Border(1, t.ControlsBorder).
		Radius(10).
		Padding(10, 12).
		Gap(8).
		Shadow(0, 8, 24, 0, colorWithAlpha(ui.RGB(0, 0, 0), 0.22)).
		Children(func() {
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
				ui.Text(c, "FLOW DEVTOOLS").FontSize(10).FontWeight(700).TextColor(t.TextMuted)
				ctrlBtn(c, t, "×", "Close DevTools", false, func() { g.ShowDevTools = false })
			})

			// ViewportLogger
			ui.Column(c).Gap(2).Children(func() {
				ui.Text(c, "Viewport").FontSize(10.5).FontWeight(600).TextColor(t.TextPrimary)
				ui.Text(c, fmt.Sprintf("x: %.1f  y: %.1f  zoom: %.2f", g.Viewport.X, g.Viewport.Y, g.Viewport.safeZoom())).
					FontSize(10.5).TextColor(t.TextMuted)
			})

			// NodeInspector
			ui.Box(c).WidthPercent(100).Height(1).Background(t.ControlsBorder)
			ui.Column(c).Gap(2).Children(func() {
				ui.Text(c, "Selected Node").FontSize(10.5).FontWeight(600).TextColor(t.TextPrimary)
				sel := g.SelectedNodes()
				if len(sel) == 0 {
					ui.Text(c, "none").FontSize(10.5).TextColor(t.TextMuted)
				} else {
					n := sel[0]
					v := n.Variant
					if v == "" {
						v = VariantBase
					}
					ui.Text(c, fmt.Sprintf("%s (%s) · %s", n.ID, v, n.Status)).FontSize(10.5).TextColor(t.TextMuted)
					ui.Text(c, fmt.Sprintf("pos: (%.0f, %.0f)  dim: %.0f×%.0f", n.X, n.Y, n.W, n.H)).FontSize(10.5).TextColor(t.TextMuted)
				}
			})

			// ChangeLogger
			ui.Box(c).WidthPercent(100).Height(1).Background(t.ControlsBorder)
			ui.Column(c).Gap(2).Children(func() {
				ui.Text(c, "Change Logger").FontSize(10.5).FontWeight(600).TextColor(t.TextPrimary)
				if len(g.ChangeLog) == 0 {
					ui.Text(c, "Interact with canvas…").FontSize(10.5).TextColor(t.TextMuted)
				} else {
					start := len(g.ChangeLog) - 4
					if start < 0 {
						start = 0
					}
					for i := len(g.ChangeLog) - 1; i >= start; i-- {
						ev := g.ChangeLog[i]
						ui.Text(c, fmt.Sprintf("• %s %s %s", ev.Kind, ev.Target, ev.Detail)).
							FontSize(10).TextColor(t.TextMuted).SingleLine().Ellipsis("…")
					}
				}
			})
		})
}
