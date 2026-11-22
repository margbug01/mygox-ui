package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
	"github.com/margbug01/mygox-ui/flow"
)

type flowState struct {
	preset        string
	graph         *flow.Graph
	palette       []flow.NodeTemplate
	fitted        bool
	showInspector bool
	showJSON      bool
	jsonPreview   string
	running       bool
	runOrder      []string
	runStep       int
	runNextAt     time.Time
	toastMsg      string
	toastUntil    time.Time
}

func newFlowState(preset string) *flowState {
	fs := &flowState{
		palette:       flow.DefaultPalette(),
		showInspector: true,
	}
	fs.loadPreset(preset)
	return fs
}

func (f *flowState) loadPreset(name string) {
	f.running = false
	f.fitted = false
	switch name {
	case "etl":
		f.preset = "etl"
		f.graph = flow.PresetDataPipeline()
	case "cicd":
		f.preset = "cicd"
		f.graph = flow.PresetCICD()
	case "blank":
		f.preset = "blank"
		f.graph = flow.NewGraph()
	default:
		f.preset = "ai-agent"
		f.graph = flow.PresetAIAgent()
	}
	f.graph.OnConnect = func(e *flow.Edge) {
		f.flash("Connected " + e.Source + " -> " + e.Target)
	}
}

func (f *flowState) flash(msg string) {
	f.toastMsg = msg
	f.toastUntil = time.Now().Add(2200 * time.Millisecond)
}

func (f *flowState) startSimulation(now time.Time) {
	order := f.graph.TopologicalOrder()
	if len(order) == 0 {
		return
	}
	for _, n := range f.graph.Nodes {
		n.Status = flow.StatusIdle
	}
	for _, e := range f.graph.Edges {
		e.Animated = false
	}
	f.running = true
	f.runOrder = order
	f.runStep = 0
	f.runNextAt = now
	f.stepSimulation(now)
}

func (f *flowState) stepSimulation(now time.Time) {
	if !f.running {
		return
	}
	if f.runStep >= len(f.runOrder) {
		for _, n := range f.graph.Nodes {
			if n.Status == flow.StatusRunning {
				n.Status = flow.StatusSuccess
			}
		}
		f.running = false
		f.flash("Workflow run completed")
		return
	}
	if f.runStep > 0 {
		if prev := f.graph.NodeByID(f.runOrder[f.runStep-1]); prev != nil {
			prev.Status = flow.StatusSuccess
		}
	}
	curID := f.runOrder[f.runStep]
	if cur := f.graph.NodeByID(curID); cur != nil {
		cur.Status = flow.StatusRunning
	}
	for _, e := range f.graph.Edges {
		if e.Target == curID || e.Source == curID {
			e.Animated = true
		}
	}
	f.runStep++
	f.runNextAt = now.Add(520 * time.Millisecond)
}

// cossFlowTheme maps coss/theme.Tokens onto flow.Theme so the node canvas shares
// the exact minimalist neutral palette, borders, and surfaces of cossgo.
func cossFlowTheme(t theme.Tokens) flow.Theme {
	ft := flow.LightTheme()
	if t.Dark {
		ft = flow.DarkTheme()
	}
	ft.Dark = t.Dark
	ft.CanvasBG = t.Background
	ft.NodeBG = t.Card
	ft.NodeHeaderBG = t.Sidebar
	ft.NodeBorder = t.Border
	ft.NodeSelected = t.Primary
	ft.TextPrimary = t.Foreground
	ft.TextMuted = t.MutedForeground
	ft.BadgeBG = t.Secondary
	ft.BadgeText = t.SecondaryForeground
	ft.PortFill = t.Card
	ft.PortBorder = t.Primary
	ft.EdgeLabelBG = t.Card
	ft.EdgeLabelBorder = t.Border
	ft.EdgeLabelText = t.Foreground
	ft.MiniMapBG = t.Card
	ft.MiniMapBorder = t.Border
	ft.MiniMapNodeSel = t.Primary
	ft.ControlsBG = t.Card
	ft.ControlsBorder = t.Border
	ft.ControlsText = t.Foreground
	ft.ControlsActiveBG = t.Primary
	ft.ControlsActiveFG = t.PrimaryForeground
	if t.Dark {
		ft.GridColor = ui.RGBA(255, 255, 255, 0.10)
		ft.EdgeDefault = ui.RGBA(255, 255, 255, 0.28)
		ft.EdgeHover = ui.RGBA(255, 255, 255, 0.55)
		ft.EdgeSelected = t.Primary
		ft.SelectionFill = ui.RGBA(255, 255, 255, 0.08)
		ft.SelectionBorder = ui.RGBA(255, 255, 255, 0.45)
		ft.MiniMapViewport = ui.RGBA(255, 255, 255, 0.10)
		ft.MiniMapViewEdge = ui.RGBA(255, 255, 255, 0.45)
	} else {
		ft.GridColor = ui.RGBA(0, 0, 0, 0.11)
		ft.EdgeDefault = ui.RGBA(0, 0, 0, 0.26)
		ft.EdgeHover = ui.RGBA(0, 0, 0, 0.55)
		ft.EdgeSelected = t.Primary
		ft.SelectionFill = ui.RGBA(0, 0, 0, 0.06)
		ft.SelectionBorder = ui.RGBA(0, 0, 0, 0.45)
		ft.MiniMapViewport = ui.RGBA(0, 0, 0, 0.06)
		ft.MiniMapViewEdge = ui.RGBA(0, 0, 0, 0.40)
	}
	return ft
}

func (a *showcaseApp) renderFlowTab(c *ui.Context, t theme.Tokens, winW, winH float32, isMobile bool) {
	f := a.flow
	ft := cossFlowTheme(t)

	if f.running {
		if !c.Now().Before(f.runNextAt) {
			f.stepSimulation(c.Now())
		}
		if f.running {
			c.After(120 * time.Millisecond)
		}
	}

	if !f.fitted && winW > 0 && winH > 0 {
		if isMobile {
			f.graph.AutoLayout(flow.LayoutTB)
			f.graph.FitView(winW, winH-108, 24)
		} else {
			cw := winW - 256
			if f.showInspector {
				cw -= 264
			}
			f.graph.FitView(cw, winH-104, 48)
		}
		f.fitted = true
	}

	ui.Row(c).WidthPercent(100).Grow(1).MinHeight(0).Children(func() {
		// Left Coss-style Sidebar (Node Palette + Presets)
		if !isMobile {
			ui.Column(c).
				Width(256).
				HeightPercent(100).
				Shrink(0).
				Background(t.Sidebar).
				BorderWidth(0, 1, 0, 0).
				BorderColor(t.SidebarBorder).
				Padding(14, 12, 12, 12).
				Gap(10).
				Children(func() {
					a.flowPaletteContent(c, t)
				})
		}

		// Center Flow Workspace
		ui.Column(c).Grow(1).MinWidth(0).HeightPercent(100).Children(func() {
			// Coss-style Secondary Action Bar
			ui.Row(c).
				WidthPercent(100).
				AlignItems(ui.Center).
				Justify(ui.SpaceBetween).
				Padding(10, 16).
				BorderWidth(0, 0, 1, 0).
				BorderColor(t.Border).
				Background(t.Background).
				Gap(8).
				Children(func() {
					ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
						if isMobile {
							coss.IconButton(c, icons.PanelLeft(), "Open node palette").
								Size(coss.SizeIconSM).
								Variant(coss.ButtonOutline).
								OnClick(func() { a.mobileNavOpen = true })
						}
						for _, p := range []struct{ id, label string }{
							{"ai-agent", "AI Agent"},
							{"etl", "Data ETL"},
							{"cicd", "CI/CD"},
							{"blank", "Blank"},
						} {
							if isMobile && p.id == "blank" {
								continue
							}
							v := coss.ButtonGhost
							if f.preset == p.id {
								v = coss.ButtonSecondary
							}
							id := p.id
							coss.Button(c, p.label).Size(coss.SizeXS).Variant(v).OnClick(func() {
								f.loadPreset(id)
								a.updateHash()
							})
						}
					})

					ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
						runLabel := "Run"
						if !isMobile {
							runLabel = "Run Flow"
						}
						if f.running {
							runLabel = "Running…"
						}
						coss.Button(c, runLabel).Size(coss.SizeXS).Variant(coss.ButtonDefault).OnClick(func() {
							f.startSimulation(c.Now())
						})

						if !isMobile {
							coss.Button(c, "Layout LR").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() {
								f.graph.AutoLayout(flow.LayoutLR)
								f.graph.FitView(0, 0, 56)
							})
							coss.Button(c, "Layout TB").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() {
								f.graph.AutoLayout(flow.LayoutTB)
								f.graph.FitView(0, 0, 56)
							})

							edgeLabel := "Bezier"
							if f.graph.EdgeKind == flow.EdgeSmoothStep {
								edgeLabel = "Step"
							} else if f.graph.EdgeKind == flow.EdgeStraight {
								edgeLabel = "Straight"
							}
							coss.Button(c, "Edge: "+edgeLabel).Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() {
								switch f.graph.EdgeKind {
								case flow.EdgeBezier:
									f.graph.EdgeKind = flow.EdgeSmoothStep
								case flow.EdgeSmoothStep:
									f.graph.EdgeKind = flow.EdgeStraight
								default:
									f.graph.EdgeKind = flow.EdgeBezier
								}
								for _, e := range f.graph.Edges {
									e.Kind = f.graph.EdgeKind
								}
							})

							bgLabel := "Dots"
							if f.graph.Background == flow.BackgroundLines {
								bgLabel = "Lines"
							} else if f.graph.Background == flow.BackgroundCross {
								bgLabel = "Cross"
							}
							coss.Button(c, "Grid: "+bgLabel).Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() {
								switch f.graph.Background {
								case flow.BackgroundDots:
									f.graph.Background = flow.BackgroundLines
								case flow.BackgroundLines:
									f.graph.Background = flow.BackgroundCross
								default:
									f.graph.Background = flow.BackgroundDots
								}
							})

							inspVar := coss.ButtonOutline
							if f.showInspector {
								inspVar = coss.ButtonSecondary
							}
							coss.Button(c, "Inspector").Size(coss.SizeXS).Variant(inspVar).OnClick(func() {
								f.showInspector = !f.showInspector
							})
						}
					})
				})

			// Live GPU Node Canvas
			ui.Box(c).WidthPercent(100).Grow(1).MinHeight(0).Children(func() {
				flow.Canvas(c, f.graph, flow.CanvasOptions{Theme: ft})

				ui.Box(c).Absolute().Top(12).Right(14).Children(func() {
					coss.Badge(c, f.graph.SummaryString()).Variant(coss.BadgeOutline)
				})

				if f.toastMsg != "" && time.Now().Before(f.toastUntil) {
					c.After(250 * time.Millisecond)
					ui.Box(c).Absolute().Top(12).Left(14).Children(func() {
						coss.Badge(c, f.toastMsg).Variant(coss.BadgeDefault)
					})
				}
			})
		})

		// Right Coss-style Property Inspector
		if !isMobile && f.showInspector {
			ui.Column(c).
				Width(264).
				HeightPercent(100).
				Shrink(0).
				Background(t.Sidebar).
				BorderWidth(0, 0, 0, 1).
				BorderColor(t.SidebarBorder).
				Padding(14).
				Gap(12).
				Children(func() {
					a.flowInspectorContent(c, t)
				})
		}
	})
}

func (a *showcaseApp) flowPaletteContent(c *ui.Context, t theme.Tokens) {
	f := a.flow
	ui.Column(c).Gap(2).Padding(0, 4).Children(func() {
		ui.Text(c, "Node Palette").FontSize(13).FontWeight(600).TextColor(t.Foreground)
		ui.Text(c, "Click to add · drag ports to connect").FontSize(11).TextColor(t.MutedForeground)
	})

	ui.Scroll(c).WidthPercent(100).Grow(1).MinHeight(0).Children(func() {
		ui.Column(c).WidthPercent(100).Gap(6).Padding(2, 4, 12, 0).Children(func() {
			for i := range f.palette {
				tpl := f.palette[i]
				accent := flow.RGBHex(tpl.AccentHex)
				row := ui.ButtonBase(c).
					WidthPercent(100).
					Padding(9, 10).
					Radius(t.RadiusMD).
					Background(t.Card).
					Border(1, t.Border).
					Cursor(ui.CursorPointer)
				if row.Hovered() {
					row.Border(1, t.Ring)
				}
				row.Children(func() {
					ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(10).Children(func() {
						ui.Box(c).Width(6).Height(26).Radius(3).Background(accent).Shrink(0)
						ui.Column(c).Grow(1).MinWidth(0).Gap(1).Children(func() {
							ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
								ui.Text(c, tpl.Title).FontSize(12.5).FontWeight(500).TextColor(t.Foreground).SingleLine()
								coss.Badge(c, tpl.Badge).Variant(coss.BadgeSecondary).Size(coss.BadgeSizeSM)
							})
							ui.Text(c, tpl.Subtitle).FontSize(11).TextColor(t.MutedForeground).SingleLine().Ellipsis("…")
						})
					})
				})
				if row.Clicked() {
					n := f.graph.SpawnFromTemplate(tpl)
					a.mobileNavOpen = false
					f.flash("Added " + n.Title)
				}
			}
		})
	})

	coss.Separator(c)
	ui.Column(c).WidthPercent(100).Gap(3).Padding(0, 4).Children(func() {
		ui.Text(c, "Shortcuts").FontSize(11).FontWeight(600).TextColor(t.Foreground)
		ui.Text(c, "Wheel/Pinch: Zoom · Drag BG: Pan").FontSize(10.5).TextColor(t.MutedForeground)
		ui.Text(c, "Shift+Drag: Box Select · Del: Delete").FontSize(10.5).TextColor(t.MutedForeground)
		ui.Text(c, "Cmd+Z/D/A: Undo / Dup / Select All").FontSize(10.5).TextColor(t.MutedForeground)
	})
}

func (a *showcaseApp) flowInspectorContent(c *ui.Context, t theme.Tokens) {
	f := a.flow
	ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
		ui.Text(c, "Inspector").FontSize(13).FontWeight(600).TextColor(t.Foreground)
		v := coss.ButtonOutline
		if f.showJSON {
			v = coss.ButtonSecondary
		}
		coss.Button(c, "JSON").Size(coss.SizeXS).Variant(v).OnClick(func() {
			f.showJSON = !f.showJSON
			if f.showJSON {
				if data, err := f.graph.SaveJSON(); err == nil {
					f.jsonPreview = string(data)
				}
			}
		})
	})

	if f.showJSON {
		ui.Scroll(c).WidthPercent(100).Grow(1).MinHeight(0).Children(func() {
			ui.Box(c).WidthPercent(100).Padding(10).Radius(t.RadiusMD).Background(t.CodeHighlight).Border(1, t.Border).Children(func() {
				ui.Text(c, f.jsonPreview).FontSize(10.5).TextColor(t.Foreground)
			})
		})
		return
	}

	selNodes := f.graph.SelectedNodes()
	selEdge := f.graph.SelectedEdge()

	if len(selNodes) == 1 {
		n := selNodes[0]
		ui.Column(c).WidthPercent(100).Gap(10).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
				ui.Text(c, "Node").FontSize(12).FontWeight(600).TextColor(t.Foreground)
				coss.Badge(c, n.ID).Variant(coss.BadgeOutline).Size(coss.BadgeSizeSM)
			})
			ui.Column(c).WidthPercent(100).Gap(4).Children(func() {
				coss.Label(c, "Title")
				coss.Input(c, &n.Title).Size(coss.InputSM)
			})
			ui.Column(c).WidthPercent(100).Gap(4).Children(func() {
				coss.Label(c, "Subtitle")
				coss.Input(c, &n.Subtitle).Size(coss.InputSM)
			})
			ui.Column(c).WidthPercent(100).Gap(4).Children(func() {
				coss.Label(c, "Status")
				ui.Row(c).Gap(6).Wrap().Children(func() {
					for _, st := range []flow.NodeStatus{flow.StatusIdle, flow.StatusRunning, flow.StatusSuccess, flow.StatusError} {
						v := coss.ButtonOutline
						if n.Status == st {
							v = coss.ButtonDefault
						}
						stVal := st
						coss.Button(c, string(st)).Size(coss.SizeXS).Variant(v).OnClick(func() { n.Status = stVal })
					}
				})
			})
			ui.Text(c, fmt.Sprintf("Pos (%.0f, %.0f) · Size %.0f×%.0f", n.X, n.Y, n.W, n.H)).
				FontSize(11).TextColor(t.MutedForeground)
			ui.Row(c).Gap(6).Children(func() {
				coss.Button(c, "Duplicate").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { f.graph.DuplicateSelected() })
				coss.Button(c, "Delete").Size(coss.SizeXS).Variant(coss.ButtonDestructiveOutline).OnClick(func() { f.graph.DeleteSelected() })
			})
		})
		return
	}

	if len(selNodes) > 1 {
		ui.Column(c).WidthPercent(100).Gap(10).Children(func() {
			ui.Text(c, fmt.Sprintf("%d nodes selected", len(selNodes))).FontSize(13).FontWeight(600).TextColor(t.Foreground)
			ui.Row(c).Gap(6).Children(func() {
				coss.Button(c, "Duplicate").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { f.graph.DuplicateSelected() })
				coss.Button(c, "Delete").Size(coss.SizeXS).Variant(coss.ButtonDestructiveOutline).OnClick(func() { f.graph.DeleteSelected() })
			})
		})
		return
	}

	if selEdge != nil {
		ui.Column(c).WidthPercent(100).Gap(10).Children(func() {
			ui.Text(c, fmt.Sprintf("Edge · %s -> %s", selEdge.Source, selEdge.Target)).FontSize(12).FontWeight(600).TextColor(t.Foreground)
			ui.Column(c).WidthPercent(100).Gap(4).Children(func() {
				coss.Label(c, "Label")
				coss.Input(c, &selEdge.Label).Size(coss.InputSM)
			})
			ui.Row(c).Gap(6).Children(func() {
				v := coss.ButtonOutline
				if selEdge.Animated {
					v = coss.ButtonDefault
				}
				coss.Button(c, "Animated").Size(coss.SizeXS).Variant(v).OnClick(func() { selEdge.Animated = !selEdge.Animated })
				coss.Button(c, "Delete").Size(coss.SizeXS).Variant(coss.ButtonDestructiveOutline).OnClick(func() { f.graph.DeleteSelected() })
			})
		})
		return
	}

	ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
		ui.Text(c, "No Selection").FontSize(13).FontWeight(600).TextColor(t.Foreground)
		ui.Text(c, "Click any node or edge to inspect properties, edit labels, toggle animated flow, or inspect JSON.").
			FontSize(11.5).TextColor(t.MutedForeground)
		ui.Row(c).Gap(6).Wrap().Children(func() {
			mv := coss.ButtonOutline
			if f.graph.ShowMiniMap {
				mv = coss.ButtonSecondary
			}
			coss.Button(c, "MiniMap").Size(coss.SizeXS).Variant(mv).OnClick(func() { f.graph.ShowMiniMap = !f.graph.ShowMiniMap })
			coss.Button(c, "Fit View").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { f.graph.FitView(0, 0, 56) })
		})
	})
}
