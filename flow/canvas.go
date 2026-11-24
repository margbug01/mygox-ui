package flow

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"
)

// CanvasOptions configures rendering and overlays of a flow Canvas.
type CanvasOptions struct {
	Theme        Theme
	NodeRenderer func(p *ui.Painter, n *Node, r ui.Rect, zoom float32, t Theme)
}

// Canvas renders an interactive, GPU-accelerated infinite node graph canvas
// modeled after React Flow UI: BaseNode, DatabaseSchemaNode, LabeledGroupNode,
// PlaceholderNode, AnnotationNode, NodeStatusIndicator, ButtonEdge, NodeToolbar,
// ZoomSlider, MiniMap, alignment guides, and keyboard shortcuts.
func Canvas(c *ui.Context, g *Graph, opts ...CanvasOptions) ui.Element {
	if g == nil {
		return ui.Box(c)
	}
	cfg := CanvasOptions{Theme: DarkTheme()}
	if len(opts) > 0 {
		cfg = opts[0]
	}
	t := cfg.Theme

	root := ui.Box(c).
		WidthPercent(100).
		HeightPercent(100).
		Background(t.CanvasBG).
		Clip()

	root.Children(func() {
		// Interactive GPU drawing surface.
		ui.Box(c).
			WidthPercent(100).
			HeightPercent(100).
			Focusable().
			Cursor(canvasCursor(g)).
			HandleInput(func(ev ui.InputEvent) bool {
				return g.handleInput(ev)
			}).
			Draw(func(p *ui.Painter, r ui.Rect) {
				g.lastCanvasW = r.W
				g.lastCanvasH = r.H
				g.paint(p, r, t, cfg.NodeRenderer)
			})

		// React Flow UI NodeToolbar floating above a single selected node.
		if !g.ReadOnly {
			sel := g.SelectedNodes()
			if len(sel) == 1 && sel[0].Variant != VariantGroup && sel[0].Variant != VariantPlaceholder && g.lastCanvasW > 0 {
				n := sel[0]
				b := n.Bounds()
				topCenter := g.Viewport.WorldToScreen(Point{X: b.X + b.W*0.5, Y: b.Y})
				if topCenter.Y > 46 && topCenter.Y < g.lastCanvasH-60 && topCenter.X > 90 && topCenter.X < g.lastCanvasW-90 {
					left := clamp32(topCenter.X-84, 10, max32(10, g.lastCanvasW-180))
					top := max32(8, topCenter.Y-38)
					ui.Box(c).Absolute().Left(left).Top(top).Children(func() {
						NodeToolbar(c, g, n, t)
					})
				}
			}
		}

		// React Flow UI NodeSearch overlay (top-center, Cmd/Ctrl+K).
		if g.ShowSearch {
			left := max32(14, (g.lastCanvasW-296)*0.5)
			ui.Box(c).Absolute().Left(left).Top(16).Children(func() {
				NodeSearch(c, g, t)
			})
		}

		// React Flow UI DevTools overlay (top-left).
		if g.ShowDevTools {
			ui.Box(c).Absolute().Left(14).Top(14).Children(func() {
				DevTools(c, g, t)
			})
		}

		// Floating bottom-left React Flow UI ZoomSlider + ZoomSelect + controls overlay.
		if g.ShowControls {
			if g.zoomSelectOpen {
				ui.Box(c).Absolute().Left(120).Bottom(54).Children(func() {
					ZoomSelect(c, g, t)
				})
			}
			ui.Box(c).
				Absolute().
				Left(14).
				Bottom(14).
				Children(func() {
					ControlsBar(c, g, t)
				})
		}
	})

	return root
}

// NodeToolbar renders a compact floating toolbar anchored above the selected node.
func NodeToolbar(c *ui.Context, g *Graph, n *Node, t Theme) {
	ui.Row(c).
		Background(t.ControlsBG).
		Border(1, t.ControlsBorder).
		Radius(8).
		Padding(3, 5).
		Gap(3).
		AlignItems(ui.Center).
		Shadow(0, 4, 14, 0, colorWithAlpha(ui.RGB(0, 0, 0), 0.16)).
		Children(func() {
			nextSt := StatusRunning
			stLabel := "Run"
			if n.Status == StatusRunning {
				nextSt = StatusSuccess
				stLabel = "Done"
			} else if n.Status == StatusSuccess {
				nextSt = StatusIdle
				stLabel = "Reset"
			}
			ctrlBtn(c, t, stLabel, "Toggle node status", n.Status == StatusRunning, func() {
				n.Status = nextSt
			})
			ctrlBtn(c, t, "Dup", "Duplicate node", false, func() {
				g.DuplicateSelected()
			})
			ctrlBtn(c, t, "Del", "Delete node", false, func() {
				g.DeleteSelected()
			})
		})
}

func canvasCursor(g *Graph) ui.Cursor {
	switch g.mode {
	case modePan:
		return ui.CursorGrabbing
	case modeDragNodes:
		return ui.CursorMove
	case modeConnect, modeBoxSelect:
		return ui.CursorCrosshair
	default:
		if g.hoverPortID != "" || g.BoxSelectMode {
			return ui.CursorCrosshair
		}
		if g.hoverNodeID != "" || g.hoverEdgeID != "" {
			return ui.CursorPointer
		}
		return ui.CursorGrab
	}
}

// ControlsBar renders React Flow UI's ZoomSlider combined with Fit, Box, Snap, and Undo/Redo.
func ControlsBar(c *ui.Context, g *Graph, t Theme) {
	ui.Row(c).
		Background(t.ControlsBG).
		Border(1, t.ControlsBorder).
		Radius(10).
		Padding(4, 6).
		Gap(4).
		AlignItems(ui.Center).
		Shadow(0, 6, 20, 0, colorWithAlpha(ui.RGB(0, 0, 0), 0.18)).
		Children(func() {
			ctrlBtn(c, t, "-", "Zoom out", false, func() { g.ZoomOut() })

			// Interactive ZoomSlider track (hidden on very narrow canvases to keep bar compact).
			if g.lastCanvasW == 0 || g.lastCanvasW >= 540 {
				zoomSliderTrack(c, g, t)
			}

			ctrlBtn(c, t, "+", "Zoom in", false, func() { g.ZoomIn() })
			// React Flow UI ZoomSelect trigger: click percentage to pick preset zoom levels.
			ctrlBtn(c, t, fmt.Sprintf("%d%% ▾", int(math.Round(float64(g.Viewport.safeZoom()*100)))), "Zoom levels", g.zoomSelectOpen, func() {
				g.zoomSelectOpen = !g.zoomSelectOpen
			})
			ctrlSep(c, t)
			ctrlBtn(c, t, "Fit", "Fit view (0)", false, func() { g.FitView(0, 0, 56) })
			ctrlBtn(c, t, "Box", "Box select mode", g.BoxSelectMode, func() { g.BoxSelectMode = !g.BoxSelectMode })
			ctrlBtn(c, t, "Snap", "Snap to grid", g.SnapToGrid, func() { g.SnapToGrid = !g.SnapToGrid })
			ctrlSep(c, t)
			ctrlBtn(c, t, "Undo", "Undo", false, func() { g.Undo() })
			ctrlBtn(c, t, "Redo", "Redo", false, func() { g.Redo() })
		})
}

func zoomSliderTrack(c *ui.Context, g *Graph, t Theme) {
	const trackW = float32(84)
	minZ := g.MinZoom
	if minZ <= 0 {
		minZ = 0.2
	}
	maxZ := g.MaxZoom
	if maxZ <= minZ {
		maxZ = 2.5
	}
	frac := clamp32((g.Viewport.safeZoom()-minZ)/(maxZ-minZ), 0, 1)

	ui.Box(c).
		Width(trackW).
		Height(24).
		Cursor(ui.CursorPointer).
		HandleInput(func(ev ui.InputEvent) bool {
			if ev.Kind == ui.InputPointerDown || (ev.Kind == ui.InputPointerMove && ev.Button == 0) {
				f := clamp32((ev.X-6)/(trackW-12), 0, 1)
				g.ZoomTo(minZ + f*(maxZ-minZ))
				return true
			}
			return false
		}).
		Draw(func(p *ui.Painter, r ui.Rect) {
			barY := r.Y + r.H*0.5 - 1.5
			innerW := r.W - 12
			p.Fill(ui.Rect{X: r.X + 6, Y: barY, W: innerW, H: 3}, t.ControlsBorder, 1.5)
			p.Fill(ui.Rect{X: r.X + 6, Y: barY, W: innerW * frac, H: 3}, t.ControlsActiveBG, 1.5)
			kx := r.X + 6 + innerW*frac
			ky := r.Y + r.H*0.5
			kr := ui.Rect{X: kx - 5, Y: ky - 5, W: 10, H: 10}
			p.Fill(kr, t.ControlsBG, 5)
			p.Stroke(kr, t.ControlsActiveBG, 5, 1.8)
		})
}

func ctrlSep(c *ui.Context, t Theme) {
	ui.Box(c).Width(1).Height(16).Background(t.ControlsBorder).Margin(0, 2)
}

func ctrlBtn(c *ui.Context, t Theme, label, tip string, active bool, onClick func()) {
	b := ui.ButtonBase(c).
		Height(26).
		Padding(0, 7).
		Radius(6).
		Center().
		Cursor(ui.CursorPointer)
	if active {
		b.Background(t.ControlsActiveBG).TextColor(t.ControlsActiveFG)
	} else if b.Hovered() {
		b.Background(t.BadgeBG).TextColor(t.TextPrimary)
	} else {
		b.TextColor(t.ControlsText)
	}
	b.Children(func() {
		ui.Text(c, label).FontSize(11).FontWeight(500).SingleLine()
	})
	if b.Clicked() && onClick != nil {
		onClick()
	}
}

func (g *Graph) handleInput(ev ui.InputEvent) bool {
	switch ev.Kind {
	case ui.InputKeyDown:
		return g.handleKeyDown(ev)

	case ui.InputZoom:
		g.Viewport.ZoomAt(Point{X: ev.X, Y: ev.Y}, ev.Scale, g.MinZoom, g.MaxZoom)
		return true

	case ui.InputScroll:
		if ev.Mods&(ui.Cmd|ui.Ctrl) != 0 {
			factor := float32(math.Exp(float64(-ev.DY) * 0.004))
			g.Viewport.ZoomAt(Point{X: ev.X, Y: ev.Y}, factor, g.MinZoom, g.MaxZoom)
			return true
		}
		if !ev.Precise && ev.DX == 0 && ev.DY != 0 {
			factor := float32(math.Exp(float64(-ev.DY) * 0.0022))
			g.Viewport.ZoomAt(Point{X: ev.X, Y: ev.Y}, factor, g.MinZoom, g.MaxZoom)
			return true
		}
		g.Viewport.X -= ev.DX
		g.Viewport.Y -= ev.DY
		return true

	case ui.InputPointerDown:
		return g.handlePointerDown(ev)

	case ui.InputPointerMove:
		return g.handlePointerMove(ev)

	case ui.InputPointerUp:
		return g.handlePointerUp(ev)
	}
	return false
}

func (g *Graph) handleKeyDown(ev ui.InputEvent) bool {
	cmd := ev.Mods&(ui.Cmd|ui.Ctrl) != 0
	shift := ev.Mods&ui.Shift != 0
	switch ev.Key {
	case ui.KeyDelete, ui.KeyBackspace:
		return g.DeleteSelected()
	case ui.KeyEscape:
		if g.ShowSearch {
			g.ShowSearch = false
			return true
		}
		if g.zoomSelectOpen {
			g.zoomSelectOpen = false
			return true
		}
		if g.mode == modeConnect {
			g.mode = modeIdle
			return true
		}
		g.ClearSelection()
		return true
	case ui.KeyK:
		if cmd {
			g.ShowSearch = !g.ShowSearch
			return true
		}
	case ui.KeyA:
		if cmd {
			g.SelectAll()
			return true
		}
	case ui.KeyD:
		if cmd {
			return g.DuplicateSelected() > 0
		}
	case ui.KeyZ:
		if cmd && shift {
			return g.Redo()
		}
		if cmd {
			return g.Undo()
		}
	case ui.KeyY:
		if cmd {
			return g.Redo()
		}
	case ui.Key0:
		g.FitView(0, 0, 56)
		return true
	case ui.KeyEqual:
		g.ZoomIn()
		return true
	case ui.KeyMinus:
		g.ZoomOut()
		return true
	}
	return false
}

func (g *Graph) handlePointerDown(ev ui.InputEvent) bool {
	ptScreen := Point{X: ev.X, Y: ev.Y}
	ptWorld := g.Viewport.ScreenToWorld(ptScreen)
	g.pressScreen = ptScreen
	g.lastScreen = ptScreen
	g.pressWorld = ptWorld
	g.guides = nil

	// 0. Check MiniMap hit first.
	if g.ShowMiniMap && g.miniMapRect.W > 0 && g.miniMapRect.Contains(ptScreen) {
		g.mode = modeMiniMap
		g.panToMiniMapPoint(ptScreen)
		return true
	}

	// 0b. Check ButtonEdge midpoint delete button.
	if !g.ReadOnly && ev.Button == 0 {
		for _, btn := range g.edgeBtnHits {
			if dist2D(ptScreen, btn.Center) <= btn.Radius+2 {
				g.RemoveEdge(btn.EdgeID)
				return true
			}
		}
		// 0c. Check ButtonHandle '+' spawn button.
		for _, hb := range g.handleBtnHits {
			if dist2D(ptScreen, hb.Center) <= hb.Radius+3 {
				g.SpawnConnectedNode(hb.NodeID, hb.PortID)
				return true
			}
		}
	}

	if ev.Button != 0 {
		g.mode = modePan
		return true
	}

	// 1. Check Port hit to start wiring an edge.
	if !g.ReadOnly {
		portRadius := 14 / g.Viewport.safeZoom()
		if n, p, pos := g.HitTestPort(ptWorld, portRadius); n != nil && p != nil {
			side := p.Side
			if side == "" {
				if p.Kind == PortInput {
					side = PortLeft
				} else {
					side = PortRight
				}
			}
			g.mode = modeConnect
			g.connect = connectPreview{
				FromNode:  n.ID,
				FromPort:  p.ID,
				FromSide:  side,
				FromKind:  p.Kind,
				FromWorld: pos,
				CurWorld:  ptWorld,
			}
			return true
		}
	}

	// 2. Check Node hit (PlaceholderNode promotes on click; GroupNode drags contained children).
	if n := g.HitTestNode(ptWorld); n != nil {
		if n.Variant == VariantPlaceholder && !g.ReadOnly {
			g.PromotePlaceholder(n.ID)
			return true
		}
		shift := ev.Mods&ui.Shift != 0
		if shift {
			g.SelectNode(n.ID, true)
		} else if !n.Selected {
			g.SelectNode(n.ID, false)
		} else if g.OnSelectNode != nil {
			g.OnSelectNode(n)
		}
		if ev.Clicks >= 2 && g.OnDoubleClick != nil {
			g.OnDoubleClick(n)
			return true
		}
		if n.Variant != VariantGroup {
			g.bringToFront(n.ID)
		}
		if !g.ReadOnly && !n.Locked {
			g.mode = modeDragNodes
			g.dragMoved = false
			g.dragStartPos = map[string]Point{}
			for _, sel := range g.SelectedNodes() {
				if !sel.Locked {
					g.dragStartPos[sel.ID] = Point{X: sel.X, Y: sel.Y}
					// Dragging a LabeledGroupNode moves all non-group nodes enclosed within it.
					if sel.Variant == VariantGroup {
						gb := sel.Bounds()
						for _, child := range g.Nodes {
							if child.ID == sel.ID || child.Variant == VariantGroup || child.Locked {
								continue
							}
							cb := child.Bounds()
							center := Point{X: cb.X + cb.W*0.5, Y: cb.Y + cb.H*0.5}
							if gb.Contains(center) {
								g.dragStartPos[child.ID] = Point{X: child.X, Y: child.Y}
							}
						}
					}
				}
			}
		}
		return true
	}

	// 3. Check Edge hit to select an edge.
	if e := g.HitTestEdge(ptScreen, 9); e != nil {
		g.SelectEdge(e.ID)
		return true
	}

	// 4. Empty background: box-select or pan.
	if ev.Mods&ui.Shift != 0 || g.BoxSelectMode {
		g.ClearSelection()
		g.mode = modeBoxSelect
		g.boxCurWorld = ptWorld
		return true
	}

	g.ClearSelection()
	g.mode = modePan
	return true
}

func (g *Graph) handlePointerMove(ev ui.InputEvent) bool {
	ptScreen := Point{X: ev.X, Y: ev.Y}
	ptWorld := g.Viewport.ScreenToWorld(ptScreen)

	switch g.mode {
	case modePan:
		g.Viewport.X += ptScreen.X - g.lastScreen.X
		g.Viewport.Y += ptScreen.Y - g.lastScreen.Y
		g.lastScreen = ptScreen
		return true

	case modeMiniMap:
		g.panToMiniMapPoint(ptScreen)
		g.lastScreen = ptScreen
		return true

	case modeBoxSelect:
		g.boxCurWorld = ptWorld
		selRect := normalizedRect(g.pressWorld, g.boxCurWorld)
		for _, n := range g.Nodes {
			if n.Variant == VariantGroup {
				continue
			}
			n.Selected = selRect.Intersects(n.Bounds())
		}
		g.lastScreen = ptScreen
		return true

	case modeDragNodes:
		dx := (ptScreen.X - g.pressScreen.X) / g.Viewport.safeZoom()
		dy := (ptScreen.Y - g.pressScreen.Y) / g.Viewport.safeZoom()
		if !g.dragMoved && (math.Abs(float64(dx)) > 1.5 || math.Abs(float64(dy)) > 1.5) {
			g.PushUndo()
			g.dragMoved = true
		}
		if g.dragMoved {
			g.guides = nil
			grid := g.GridSize
			if grid <= 0 {
				grid = 20
			}
			var snapDX, snapDY float32
			if len(g.dragStartPos) == 1 {
				for id, start := range g.dragStartPos {
					n := g.NodeByID(id)
					if n == nil {
						continue
					}
					candX := start.X + dx
					candY := start.Y + dy
					if g.SnapToGrid {
						candX = float32(math.Round(float64(candX/grid))) * grid
						candY = float32(math.Round(float64(candY/grid))) * grid
					}
					b := n.Bounds()
					const alignTol = float32(7)
					for _, other := range g.Nodes {
						if other.ID == id || other.Variant == VariantGroup {
							continue
						}
						ob := other.Bounds()
						if math.Abs(float64(candX-ob.X)) <= float64(alignTol) {
							snapDX = ob.X - candX
							g.guides = append(g.guides, alignmentGuide{Vertical: true, Pos: ob.X})
						} else if math.Abs(float64((candX+b.W*0.5)-(ob.X+ob.W*0.5))) <= float64(alignTol) {
							snapDX = (ob.X + ob.W*0.5) - (candX + b.W*0.5)
							g.guides = append(g.guides, alignmentGuide{Vertical: true, Pos: ob.X + ob.W*0.5})
						}
						if math.Abs(float64(candY-ob.Y)) <= float64(alignTol) {
							snapDY = ob.Y - candY
							g.guides = append(g.guides, alignmentGuide{Vertical: false, Pos: ob.Y})
						} else if math.Abs(float64((candY+b.H*0.5)-(ob.Y+ob.H*0.5))) <= float64(alignTol) {
							snapDY = (ob.Y + ob.H*0.5) - (candY + b.H*0.5)
							g.guides = append(g.guides, alignmentGuide{Vertical: false, Pos: ob.Y + ob.H*0.5})
						}
					}
				}
			}

			for id, start := range g.dragStartPos {
				n := g.NodeByID(id)
				if n == nil {
					continue
				}
				nx := start.X + dx
				ny := start.Y + dy
				if g.SnapToGrid {
					nx = float32(math.Round(float64(nx/grid))) * grid
					ny = float32(math.Round(float64(ny/grid))) * grid
				}
				n.X = nx + snapDX
				n.Y = ny + snapDY
			}
		}
		g.lastScreen = ptScreen
		return true

	case modeConnect:
		g.connect.CurWorld = ptWorld
		g.connect.HasSnap = false
		snapR := 24 / g.Viewport.safeZoom()
		if n, p, pos := g.HitTestPort(ptWorld, snapR); n != nil && p != nil {
			srcN, srcP, dstN, dstP := g.connect.FromNode, g.connect.FromPort, n.ID, p.ID
			if g.connect.FromKind == PortInput {
				srcN, srcP, dstN, dstP = n.ID, p.ID, g.connect.FromNode, g.connect.FromPort
			}
			if g.CanConnect(srcN, srcP, dstN, dstP) {
				g.connect.HasSnap = true
				g.connect.SnapNode = n.ID
				g.connect.SnapPort = p.ID
				g.connect.SnapWorld = pos
			}
		}
		g.lastScreen = ptScreen
		return true

	default:
		g.hoverPortNode = ""
		g.hoverPortID = ""
		g.hoverNodeID = ""
		g.hoverEdgeID = ""
		portR := 14 / g.Viewport.safeZoom()
		if n, p, _ := g.HitTestPort(ptWorld, portR); n != nil && p != nil {
			g.hoverPortNode = n.ID
			g.hoverPortID = p.ID
			return false
		}
		if n := g.HitTestNode(ptWorld); n != nil {
			g.hoverNodeID = n.ID
			return false
		}
		if e := g.HitTestEdge(ptScreen, 9); e != nil {
			g.hoverEdgeID = e.ID
		}
		return false
	}
}

func (g *Graph) handlePointerUp(ev ui.InputEvent) bool {
	if g.mode == modeIdle {
		return false
	}
	prevMode := g.mode
	g.mode = modeIdle
	g.guides = nil

	if ev.Canceled {
		return true
	}

	switch prevMode {
	case modeConnect:
		if g.connect.HasSnap {
			srcN, srcP, dstN, dstP := g.connect.FromNode, g.connect.FromPort, g.connect.SnapNode, g.connect.SnapPort
			if g.connect.FromKind == PortInput {
				srcN, srcP, dstN, dstP = g.connect.SnapNode, g.connect.SnapPort, g.connect.FromNode, g.connect.FromPort
			}
			g.Connect(srcN, srcP, dstN, dstP)
		}
	case modeDragNodes:
		if g.dragMoved {
			g.notifyChange()
		}
	}
	return true
}

func (g *Graph) bringToFront(id string) {
	idx := -1
	for i, n := range g.Nodes {
		if n.ID == id {
			idx = i
			break
		}
	}
	if idx < 0 || idx == len(g.Nodes)-1 {
		return
	}
	n := g.Nodes[idx]
	copy(g.Nodes[idx:], g.Nodes[idx+1:])
	g.Nodes[len(g.Nodes)-1] = n
}

func (g *Graph) paint(p *ui.Painter, r ui.Rect, t Theme, customNode func(*ui.Painter, *Node, ui.Rect, float32, Theme)) {
	p.Fill(r, t.CanvasBG, 0)

	// 1. Background pattern.
	g.paintBackground(p, r, t)

	// 2. Group backdrops (LabeledGroupNode sits underneath edges & child nodes).
	needsAnim := false
	nowSec := float64(p.Now().UnixNano()) / 1e9
	for _, n := range g.Nodes {
		if n.Variant == VariantGroup {
			g.paintGroupNode(p, r, n, t)
		}
	}

	// 3. Alignment guides during drag.
	for _, guide := range g.guides {
		if guide.Vertical {
			sp := g.Viewport.WorldToScreen(Point{X: guide.Pos, Y: 0})
			p.Line(r.X+sp.X, r.Y, r.X+sp.X, r.Y+r.H, 1, t.GuideLine)
		} else {
			sp := g.Viewport.WorldToScreen(Point{X: 0, Y: guide.Pos})
			p.Line(r.X, r.Y+sp.Y, r.X+r.W, r.Y+sp.Y, 1, t.GuideLine)
		}
	}

	// 4. Edges + ButtonEdge midpoint controls + AnimatedSvgEdge packets.
	g.edgeBtnHits = g.edgeBtnHits[:0]
	g.handleBtnHits = g.handleBtnHits[:0]
	for _, e := range g.Edges {
		if e.Animated || e.PacketMotion {
			needsAnim = true
		}
		g.paintEdge(p, r, e, t, nowSec)
	}

	// 5. Live connection preview wire.
	if g.mode == modeConnect {
		needsAnim = true
		g.paintConnectPreview(p, r, t, nowSec)
	}

	// 6. Regular, Schema, Placeholder & Annotation nodes.
	for _, n := range g.Nodes {
		if n.Variant == VariantGroup {
			continue
		}
		if n.Status == StatusRunning {
			needsAnim = true
		}
		g.paintNode(p, r, n, t, nowSec, customNode)
	}

	// 6b. React Flow UI NodeTooltip rendered on top of hovered node.
	if g.hoverNodeID != "" && g.mode == modeIdle {
		if hn := g.NodeByID(g.hoverNodeID); hn != nil && hn.Tooltip != "" {
			g.paintNodeTooltip(p, r, hn, t)
		}
	}

	// 7. Box selection marquee.
	if g.mode == modeBoxSelect {
		wRect := normalizedRect(g.pressWorld, g.boxCurWorld)
		tl := g.Viewport.WorldToScreen(Point{X: wRect.X, Y: wRect.Y})
		br := g.Viewport.WorldToScreen(Point{X: wRect.X + wRect.W, Y: wRect.Y + wRect.H})
		sr := ui.Rect{X: r.X + tl.X, Y: r.Y + tl.Y, W: br.X - tl.X, H: br.Y - tl.Y}
		p.Fill(sr, t.SelectionFill, 4)
		p.Stroke(sr, t.SelectionBorder, 4, 1.25)
	}

	// 8. Interactive MiniMap (hidden on narrow viewports to avoid overlapping ZoomSlider).
	if g.ShowMiniMap && r.W >= 620 && r.H >= 280 {
		g.paintMiniMap(p, r, t)
	} else {
		g.miniMapRect = Rect{}
	}

	if needsAnim {
		p.AnimationFrame()
	}
}

func (g *Graph) paintBackground(p *ui.Painter, r ui.Rect, t Theme) {
	if g.Background == BackgroundNone {
		return
	}
	step := g.GridSize
	if step <= 0 {
		step = 20
	}
	z := g.Viewport.safeZoom()
	screenStep := step * z
	for screenStep < 14 {
		screenStep *= 2
	}
	for screenStep > 56 {
		screenStep *= 0.5
	}

	offX := float32(math.Mod(float64(g.Viewport.X), float64(screenStep)))
	if offX < 0 {
		offX += screenStep
	}
	offY := float32(math.Mod(float64(g.Viewport.Y), float64(screenStep)))
	if offY < 0 {
		offY += screenStep
	}

	switch g.Background {
	case BackgroundLines:
		for x := offX; x < r.W; x += screenStep {
			p.Line(r.X+x, r.Y, r.X+x, r.Y+r.H, 1, t.GridColor)
		}
		for y := offY; y < r.H; y += screenStep {
			p.Line(r.X, r.Y+y, r.X+r.W, r.Y+y, 1, t.GridColor)
		}
	case BackgroundCross:
		arm := clamp32(2.5*z, 2, 4.5)
		for x := offX; x < r.W; x += screenStep {
			for y := offY; y < r.H; y += screenStep {
				p.Line(r.X+x-arm, r.Y+y, r.X+x+arm, r.Y+y, 1, t.GridColor)
				p.Line(r.X+x, r.Y+y-arm, r.X+x, r.Y+y+arm, 1, t.GridColor)
			}
		}
	default:
		dotSize := clamp32(2.0*z, 1.5, 3.0)
		half := dotSize * 0.5
		for x := offX; x < r.W; x += screenStep {
			for y := offY; y < r.H; y += screenStep {
				p.Fill(ui.Rect{X: r.X + x - half, Y: r.Y + y - half, W: dotSize, H: dotSize}, t.GridColor, half)
			}
		}
	}
}

func (g *Graph) paintEdge(p *ui.Painter, r ui.Rect, e *Edge, t Theme, nowSec float64) {
	sn := g.NodeByID(e.Source)
	dn := g.NodeByID(e.Target)
	if sn == nil || dn == nil {
		return
	}
	spW, spSide, _ := PortPosition(sn, e.SourcePort)
	dpW, dpSide, _ := PortPosition(dn, e.TargetPort)
	sp := g.Viewport.WorldToScreen(spW)
	dp := g.Viewport.WorldToScreen(dpW)
	sp.X += r.X
	sp.Y += r.Y
	dp.X += r.X
	dp.Y += r.Y

	kind := e.Kind
	if kind == "" {
		kind = g.EdgeKind
	}
	geom := BuildEdgeGeometry(sp, dp, spSide, dpSide, kind)

	col := t.EdgeDefault
	if e.ColorHex != 0 {
		col = RGBHex(e.ColorHex)
	}
	width := clamp32(1.8*g.Viewport.safeZoom(), 1.25, 2.8)
	if e.ID == g.hoverEdgeID {
		col = t.EdgeHover
		width += 0.5
	}
	if e.Selected {
		col = t.EdgeSelected
		width += 0.75
		p.StrokePath(geom.Path, width+4.5, colorWithAlpha(t.EdgeSelected, 0.18))
	}

	if e.Animated {
		p.StrokePath(geom.Path, width, colorWithAlpha(col, 0.30))
		animPath := BuildEdgeGeometry(sp, dp, spSide, dpSide, kind).Path
		dashOffset := float32(-math.Mod(nowSec*36, 16))
		animPath.Dash(dashOffset, 7, 7)
		animCol := col
		if !e.Selected && e.ColorHex == 0 {
			animCol = t.NodeSelected
		}
		p.StrokePath(animPath, width+0.25, animCol)
	} else {
		p.StrokePath(geom.Path, width, col)
	}

	// React Flow UI AnimatedSvgEdge: glowing data packet traveling smoothly along curve.
	if e.PacketMotion && len(geom.Samples) >= 2 {
		dur := float64(e.PacketDuration)
		if dur <= 0.2 {
			dur = 2.2
		}
		rawT := math.Mod(nowSec/dur, 1.0)
		frac := float32(rawT)
		switch e.PacketDirection {
		case "reverse":
			frac = 1 - frac
		case "alternate":
			cycle := math.Mod(nowSec/dur, 2.0)
			if cycle > 1.0 {
				frac = float32(2.0 - cycle)
			} else {
				frac = float32(cycle)
			}
		}
		pkt := SamplePolylineAt(geom.Samples, frac)
		pktCol := t.NodeSelected
		if e.ColorHex != 0 {
			pktCol = RGBHex(e.ColorHex)
		}
		pr := clamp32(4.2*g.Viewport.safeZoom(), 2.8, 5.6)
		glowR := pr * 2.1
		p.Fill(ui.Rect{X: pkt.X - glowR, Y: pkt.Y - glowR, W: glowR * 2, H: glowR * 2}, colorWithAlpha(pktCol, 0.28), glowR)
		p.Fill(ui.Rect{X: pkt.X - pr, Y: pkt.Y - pr, W: pr * 2, H: pr * 2}, pktCol, pr)
	}

	arrowSize := clamp32(7.5*g.Viewport.safeZoom(), 4.5, 10)
	tipBack := clamp32(6.0*g.Viewport.safeZoom(), 4, 8.5)
	tip := Point{
		X: geom.Tip.X - float32(math.Cos(float64(geom.Angle)))*tipBack,
		Y: geom.Tip.Y - float32(math.Sin(float64(geom.Angle)))*tipBack,
	}
	p.FillPath(ArrowPath(tip, geom.Angle, arrowSize), col)

	// ButtonEdge: show clickable delete button on selected or hovered edges.
	if (e.Selected || e.ID == g.hoverEdgeID) && !g.ReadOnly {
		rad := float32(10)
		bx, by := geom.LabelPos.X, geom.LabelPos.Y
		br := ui.Rect{X: bx - rad, Y: by - rad, W: rad * 2, H: rad * 2}
		p.Shadow(br, rad, 0, 2, 8, 0, colorWithAlpha(ui.RGB(0, 0, 0), 0.18))
		p.Fill(br, t.EdgeLabelBG, rad)
		p.Stroke(br, t.EdgeSelected, rad, 1.25)
		p.Line(bx-3.2, by-3.2, bx+3.2, by+3.2, 1.4, t.EdgeLabelText)
		p.Line(bx+3.2, by-3.2, bx-3.2, by+3.2, 1.4, t.EdgeLabelText)
		g.edgeBtnHits = append(g.edgeBtnHits, edgeButtonHit{
			EdgeID: e.ID,
			Center: Point{X: bx - r.X, Y: by - r.Y},
			Radius: rad,
		})
		return
	}

	// Centered DataEdge label pill.
	if e.Label != "" && g.Viewport.safeZoom() >= 0.45 {
		fontSize := clamp32(10*g.Viewport.safeZoom(), 8.5, 12.5)
		pillW := float32(len(e.Label))*fontSize*0.60 + 14
		pillH := fontSize + 9
		lr := ui.Rect{
			X: geom.LabelPos.X - pillW*0.5,
			Y: geom.LabelPos.Y - pillH*0.5,
			W: pillW,
			H: pillH,
		}
		p.Fill(lr, t.EdgeLabelBG, pillH*0.5)
		p.Stroke(lr, t.EdgeLabelBorder, pillH*0.5, 1)
		p.Text(lr.X+7, lr.Y+(pillH-fontSize)*0.5-1, e.Label, fontSize, t.EdgeLabelText)
	}
}

func (g *Graph) paintConnectPreview(p *ui.Painter, r ui.Rect, t Theme, nowSec float64) {
	dstWorld := g.connect.CurWorld
	dstSide := PortLeft
	if g.connect.FromKind == PortInput {
		dstSide = PortRight
	}
	if g.connect.HasSnap {
		dstWorld = g.connect.SnapWorld
	}
	srcS := g.Viewport.WorldToScreen(g.connect.FromWorld)
	dstS := g.Viewport.WorldToScreen(dstWorld)
	srcS.X += r.X
	srcS.Y += r.Y
	dstS.X += r.X
	dstS.Y += r.Y

	geom := BuildEdgeGeometry(srcS, dstS, g.connect.FromSide, dstSide, g.EdgeKind)
	dashOffset := float32(-math.Mod(nowSec*40, 14))
	geom.Path.Dash(dashOffset, 7, 7)

	wireCol := t.NodeSelected
	if g.connect.HasSnap {
		wireCol = t.PortSnap
		ring := ui.Rect{X: dstS.X - 11, Y: dstS.Y - 11, W: 22, H: 22}
		p.Fill(ring, colorWithAlpha(t.PortSnap, 0.22), 11)
		p.Stroke(ring, t.PortSnap, 11, 1.8)
	}
	p.StrokePath(geom.Path, 2.0, wireCol)
}

func (g *Graph) paintGroupNode(p *ui.Painter, r ui.Rect, n *Node, t Theme) {
	b := n.Bounds()
	tl := g.Viewport.WorldToScreen(Point{X: b.X, Y: b.Y})
	z := g.Viewport.safeZoom()
	nr := ui.Rect{X: r.X + tl.X, Y: r.Y + tl.Y, W: b.W * z, H: b.H * z}
	if nr.X+nr.W < r.X-40 || nr.Y+nr.H < r.Y-40 || nr.X > r.X+r.W+40 || nr.Y > r.Y+r.H+40 {
		return
	}

	rad := clamp32(10*z, 4, 14)
	bg := colorWithAlpha(t.NodeHeaderBG, 0.52)
	if !t.Dark {
		bg = colorWithAlpha(t.NodeHeaderBG, 0.72)
	}
	p.Fill(nr, bg, rad)

	borderCol := t.NodeBorder
	borderW := float32(1)
	if n.Selected {
		borderCol = t.NodeSelected
		borderW = 1.5
	}
	p.Stroke(nr, borderCol, rad, borderW)

	// React Flow LabeledGroupNode top-left corner pill tag.
	if n.Title != "" && z >= 0.4 {
		fs := clamp32(11*z, 8, 13)
		lw := float32(len(n.Title))*fs*0.60 + 18*z
		lh := fs + 10*z
		p.Clip(nr, rad, func() {
			tagR := ui.Rect{X: nr.X, Y: nr.Y, W: lw, H: lh}
			p.Fill(tagR, t.BadgeBG, 0)
			p.Line(nr.X, nr.Y+lh, nr.X+lw, nr.Y+lh, 1, t.NodeBorder)
			p.Line(nr.X+lw, nr.Y, nr.X+lw, nr.Y+lh, 1, t.NodeBorder)
			p.Text(nr.X+9*z, nr.Y+(lh-fs)*0.5-0.5, n.Title, fs, t.TextPrimary)
		})
	}
}

func (g *Graph) paintNode(p *ui.Painter, r ui.Rect, n *Node, t Theme, nowSec float64, custom func(*ui.Painter, *Node, ui.Rect, float32, Theme)) {
	b := n.Bounds()
	tl := g.Viewport.WorldToScreen(Point{X: b.X, Y: b.Y})
	z := g.Viewport.safeZoom()
	nr := ui.Rect{
		X: r.X + tl.X,
		Y: r.Y + tl.Y,
		W: b.W * z,
		H: b.H * z,
	}
	if nr.X+nr.W < r.X-40 || nr.Y+nr.H < r.Y-40 || nr.X > r.X+r.W+40 || nr.Y > r.Y+r.H+40 {
		return
	}

	if custom != nil {
		custom(p, n, nr, z, t)
		g.paintNodePorts(p, r, n, t, z)
		return
	}

	switch n.Variant {
	case VariantPlaceholder:
		g.paintPlaceholderNode(p, nr, n, t, z)
		return
	case VariantAnnotation:
		g.paintAnnotationNode(p, nr, n, t, z)
		return
	case VariantSchema:
		g.paintSchemaNode(p, r, nr, n, t, z, nowSec)
		return
	}

	// React Flow UI BaseNode + NodeStatusIndicator.
	rad := clamp32(10*z, 5, 14)
	g.paintNodeStatusFrame(p, nr, rad, n, t, z, nowSec)

	// Card body.
	p.Fill(nr, t.NodeBG, rad)

	headerH := clamp32(34*z, 18, 42)
	p.Clip(nr, rad, func() {
		// Subtle coss-style header strip (no loud thick bar; clean hairline border).
		p.Fill(ui.Rect{X: nr.X, Y: nr.Y, W: nr.W, H: headerH}, t.NodeHeaderBG, 0)
		p.Line(nr.X, nr.Y+headerH, nr.X+nr.W, nr.Y+headerH, 1, t.NodeBorder)

		if z >= 0.42 {
			// Subtle accent/status indicator dot.
			dotCol := statusColor(n.Status, t)
			if n.Status == "" || n.Status == StatusIdle {
				if n.AccentHex != 0 {
					dotCol = RGBHex(n.AccentHex)
				}
			}
			dotD := clamp32(7.0*z, 4, 8.5)
			p.Fill(ui.Rect{
				X: nr.X + 11*z,
				Y: nr.Y + headerH*0.5 - dotD*0.5,
				W: dotD,
				H: dotD,
			}, dotCol, dotD*0.5)

			// Title.
			titleSize := clamp32(12.5*z, 8, 15.5)
			p.Text(nr.X+23*z, nr.Y+(headerH-titleSize)*0.5-0.5, truncateStr(n.Title, 22), titleSize, t.TextPrimary)

			// Right-aligned restrained badge pill.
			if n.Badge != "" && nr.W >= 92 {
				bs := clamp32(9.2*z, 7.2, 11)
				bw := float32(len(n.Badge))*bs*0.62 + 11*z
				bh := bs + 6*z
				br := ui.Rect{
					X: nr.X + nr.W - bw - 8*z,
					Y: nr.Y + (headerH-bh)*0.5,
					W: bw,
					H: bh,
				}
				p.Fill(br, t.BadgeBG, bh*0.5)
				p.Text(br.X+5.5*z, br.Y+(bh-bs)*0.5-0.5, n.Badge, bs, t.BadgeText)
			}

			// BaseNodeContent: subtitle + structured fields.
			curY := nr.Y + headerH + 9*z
			subSize := clamp32(10.5*z, 7.5, 13)
			if n.Subtitle != "" {
				p.Text(nr.X+12*z, curY, truncateStr(n.Subtitle, 28), subSize, t.TextMuted)
				curY += subSize + 7*z
			}
			for _, f := range n.Fields {
				if curY+subSize > nr.Y+nr.H-6*z {
					break
				}
				p.Text(nr.X+12*z, curY, f.Key+":", subSize, t.TextMuted)
				valX := nr.X + 12*z + float32(len(f.Key)+2)*subSize*0.58
				p.Text(valX, curY, truncateStr(f.Value, 18), subSize, t.TextPrimary)
				curY += subSize + 5*z
			}

			// Optional BaseNodeFooter.
			if n.Footer != "" && nr.H >= 96*z {
				footH := clamp32(24*z, 16, 28)
				fy := nr.Y + nr.H - footH
				p.Line(nr.X, fy, nr.X+nr.W, fy, 1, t.NodeBorder)
				p.Text(nr.X+12*z, fy+(footH-subSize)*0.5-0.5, truncateStr(n.Footer, 26), clamp32(9.5*z, 7.5, 11.5), t.TextMuted)
			}
		}
	})

	// Outer crisp border.
	g.paintNodeOuterBorder(p, nr, rad, n, t)

	// React Flow UI NodeStatusIndicator "overlay" variant: frosted scrim + spinning ring.
	if n.Status == StatusRunning && n.StatusVariant == StatusOverlay {
		p.Fill(nr, colorWithAlpha(t.CanvasBG, 0.58), rad)
		cx, cy := nr.X+nr.W*0.5, nr.Y+nr.H*0.5
		ping := float32(math.Mod(nowSec*1.2, 1.0))
		pingR := clamp32((10+14*ping)*z, 6, 26)
		p.Fill(ui.Rect{X: cx - pingR, Y: cy - pingR, W: pingR * 2, H: pingR * 2}, ui.RGBA(59, 130, 246, (1-ping)*0.28), pingR)
		sr := clamp32(11*z, 6, 15)
		ringRect := ui.Rect{X: cx - sr, Y: cy - sr, W: sr * 2, H: sr * 2}
		p.Stroke(ringRect, ui.RGBA(59, 130, 246, 0.25), sr, 2.2)
		spinPath := new(ui.Path)
		startA := nowSec * 5.2
		const segs = 12
		for i := 0; i <= segs; i++ {
			a := startA + float64(i)/segs*math.Pi*1.15
			px := cx + float32(math.Cos(a))*sr
			py := cy + float32(math.Sin(a))*sr
			if i == 0 {
				spinPath.MoveTo(px, py)
			} else {
				spinPath.LineTo(px, py)
			}
		}
		p.StrokePath(spinPath, 2.2, ui.RGB(59, 130, 246))
	}

	// React Flow UI NodeAppendix attached outside top/bottom/left/right edge of node.
	if n.Appendix != "" && z >= 0.45 {
		g.paintNodeAppendix(p, nr, n, t, z)
	}

	// Handles.
	g.paintNodePorts(p, r, n, t, z)
}

func (g *Graph) paintNodeAppendix(p *ui.Painter, nr ui.Rect, n *Node, t Theme, z float32) {
	fs := clamp32(9.8*z, 7.5, 11.5)
	ah := fs + 11*z
	gap := 4 * z
	var ar ui.Rect
	switch n.AppendixPos {
	case PortTop:
		ar = ui.Rect{X: nr.X, Y: nr.Y - ah - gap, W: nr.W, H: ah}
	case PortLeft:
		aw := float32(len(n.Appendix))*fs*0.58 + 16*z
		ar = ui.Rect{X: nr.X - aw - gap, Y: nr.Y + (nr.H-ah)*0.5, W: aw, H: ah}
	case PortRight:
		aw := float32(len(n.Appendix))*fs*0.58 + 16*z
		ar = ui.Rect{X: nr.X + nr.W + gap, Y: nr.Y + (nr.H-ah)*0.5, W: aw, H: ah}
	default: // PortBottom
		ar = ui.Rect{X: nr.X, Y: nr.Y + nr.H + gap, W: nr.W, H: ah}
	}
	rad := clamp32(7*z, 4, 9)
	p.Fill(ar, t.NodeBG, rad)
	p.Stroke(ar, t.NodeBorder, rad, 1)
	lw := float32(len(n.Appendix)) * fs * 0.56
	tx := ar.X + max32(8*z, (ar.W-lw)*0.5)
	p.Text(tx, ar.Y+(ar.H-fs)*0.5-0.5, truncateStr(n.Appendix, 32), fs, t.TextMuted)
}

func (g *Graph) paintNodeTooltip(p *ui.Painter, r ui.Rect, n *Node, t Theme) {
	b := n.Bounds()
	z := g.Viewport.safeZoom()
	topMid := g.Viewport.WorldToScreen(Point{X: b.X + b.W*0.5, Y: b.Y})
	fs := float32(11)
	tw := float32(len(n.Tooltip))*fs*0.58 + 20
	th := fs + 12
	tr := ui.Rect{
		X: r.X + topMid.X - tw*0.5,
		Y: r.Y + topMid.Y - th - 10,
		W: tw,
		H: th,
	}
	if tr.Y < r.Y+6 {
		tr.Y = r.Y + topMid.Y + b.H*z + 8
	}
	p.Shadow(tr, 7, 0, 4, 12, 0, colorWithAlpha(ui.RGB(0, 0, 0), 0.20))
	p.Fill(tr, t.ControlsActiveBG, 7)
	p.Text(tr.X+10, tr.Y+(th-fs)*0.5-0.5, n.Tooltip, fs, t.ControlsActiveFG)
}

func (g *Graph) paintSchemaNode(p *ui.Painter, r, nr ui.Rect, n *Node, t Theme, z float32, nowSec float64) {
	rad := clamp32(9*z, 4, 12)
	g.paintNodeStatusFrame(p, nr, rad, n, t, z, nowSec)
	p.Fill(nr, t.NodeBG, rad)

	headerH := 34 * z
	rowH := 26 * z

	p.Clip(nr, rad, func() {
		// DatabaseSchemaNodeHeader: bg-secondary text-muted-foreground p-2
		// text-sm, its title at the start of the strip.
		p.Fill(ui.Rect{X: nr.X, Y: nr.Y, W: nr.W, H: headerH}, t.BadgeBG, 0)

		if z >= 0.42 {
			ts := clamp32(13.5*z, 8.5, 16)
			p.Text(nr.X+10*z, nr.Y+(headerH-ts)*0.5-0.5, n.Title, ts, t.TextMuted)

			// DatabaseSchemaTableRow entries: left column name, right column muted type.
			fs := clamp32(10.8*z, 7.5, 13)
			for idx, row := range n.SchemaRows {
				ry := nr.Y + headerH + float32(idx)*rowH + 2*z
				if idx > 0 {
					p.Line(nr.X, ry, nr.X+nr.W, ry, 1, colorWithAlpha(t.NodeBorder, 0.7))
				}
				p.Text(nr.X+13*z, ry+(rowH-fs)*0.5-0.5, row.Name, fs, t.TextPrimary)
				typeW := float32(len(row.Type)) * fs * 0.58
				p.Text(nr.X+nr.W-typeW-13*z, ry+(rowH-fs)*0.5-0.5, row.Type, clamp32(fs*0.92, 7, 12), t.TextMuted)
			}
		}
	})

	g.paintNodeOuterBorder(p, nr, rad, n, t)
	g.paintNodePorts(p, r, n, t, z)
}

func (g *Graph) paintPlaceholderNode(p *ui.Painter, nr ui.Rect, n *Node, t Theme, z float32) {
	rad := clamp32(9*z, 4, 12)
	fill := colorWithAlpha(t.NodeBG, 0.72)
	if n.ID == g.hoverNodeID {
		fill = t.NodeHeaderBG
	}
	p.Fill(nr, fill, rad)
	borderCol := t.TextMuted
	if n.ID == g.hoverNodeID || n.Selected {
		borderCol = t.NodeSelected
	}
	p.StrokeDashed(nr, borderCol, rad, 1.3)

	if z >= 0.42 {
		label := n.Title
		if label == "" {
			label = "+ Add node"
		}
		fs := clamp32(11.5*z, 8, 14)
		lw := float32(len(label)) * fs * 0.56
		p.Text(nr.X+(nr.W-lw)*0.5, nr.Y+(nr.H-fs)*0.5-1, label, fs, t.TextMuted)
	}
}

func (g *Graph) paintAnnotationNode(p *ui.Painter, nr ui.Rect, n *Node, t Theme, z float32) {
	if z < 0.35 {
		return
	}
	if n.Selected {
		p.StrokeDashed(nr, t.NodeSelected, 8, 1)
	}
	numSize := clamp32(10.5*z, 8, 12.5)
	textSize := clamp32(11.5*z, 8.5, 13.5)
	if n.Badge != "" {
		p.Text(nr.X+4*z, nr.Y+4*z, n.Badge, numSize, t.TextMuted)
	}
	p.Text(nr.X+4*z, nr.Y+20*z, n.Title, textSize, t.TextPrimary)
	if n.Subtitle != "" {
		p.Text(nr.X+4*z, nr.Y+38*z, n.Subtitle, clamp32(10.5*z, 7.5, 12), t.TextMuted)
	}
}

func (g *Graph) paintNodeStatusFrame(p *ui.Painter, nr ui.Rect, rad float32, n *Node, t Theme, z float32, nowSec float64) {
	// Elevation shadow (stronger when selected, matching React Flow BaseNode).
	shadowAlpha := float32(1.0)
	blur := clamp32(14*z, 5, 22)
	if n.Selected {
		blur = clamp32(22*z, 8, 30)
		shadowAlpha = 1.35
	}
	p.Shadow(nr, rad, 0, clamp32(5*z, 2, 10), blur, 0, colorWithAlpha(t.NodeShadow, shadowAlpha))

	// React Flow UI NodeStatusIndicator: 2s rotating conic-gradient spinning border when Running.
	if n.Status == StatusRunning {
		outer := ui.Rect{X: nr.X - 2.2, Y: nr.Y - 2.2, W: nr.W + 4.4, H: nr.H + 4.4}
		angle := float32(math.Mod(nowSec*180.0, 360.0)) // 2s per revolution
		p.FillConicGradient(outer, ui.ConicGradient{
			Angle: angle,
			Stops: []ui.GradientStop{
				{Color: ui.RGBA(59, 130, 246, 0.95), At: 0.0},
				{Color: ui.RGBA(99, 102, 241, 0.55), At: 0.25},
				{Color: ui.RGBA(59, 130, 246, 0.0), At: 0.65},
				{Color: ui.RGBA(59, 130, 246, 0.0), At: 1.0},
			},
		}, rad+2.2)
	} else if n.Selected {
		halo := ui.Rect{X: nr.X - 2, Y: nr.Y - 2, W: nr.W + 4, H: nr.H + 4}
		p.Stroke(halo, colorWithAlpha(t.NodeSelected, 0.22), rad+2, 2.5)
	}
}

func (g *Graph) paintNodeOuterBorder(p *ui.Painter, nr ui.Rect, rad float32, n *Node, t Theme) {
	borderCol := t.NodeBorder
	borderW := float32(1)
	switch n.Status {
	case StatusSuccess:
		borderCol = t.StatusSuccess
		borderW = 1.5
	case StatusError:
		borderCol = t.StatusError
		borderW = 1.5
	}
	if n.Selected {
		borderCol = t.NodeSelected
		borderW = 1.6
	} else if n.ID == g.hoverNodeID && n.Status != StatusSuccess && n.Status != StatusError {
		borderCol = colorWithAlpha(t.NodeSelected, 0.5)
		borderW = 1.25
	}
	p.Stroke(nr, borderCol, rad, borderW)
}

func (g *Graph) paintNodePorts(p *ui.Painter, r ui.Rect, n *Node, t Theme, z float32) {
	for i := range n.Ports {
		port := &n.Ports[i]
		wPos, side, _ := PortPosition(n, port.ID)
		sPos := g.Viewport.WorldToScreen(wPos)
		sx, sy := r.X+sPos.X, r.Y+sPos.Y

		hovered := g.hoverPortNode == n.ID && g.hoverPortID == port.ID
		snapped := g.mode == modeConnect && g.connect.HasSnap && g.connect.SnapNode == n.ID && g.connect.SnapPort == port.ID
		// BaseHandle: an 11px circle, bg-slate-100 with a border-slate-300 hairline.
		pr := clamp32(5.5*z, 3.0, 7)
		if hovered || snapped {
			pr += 1.6
		}
		borderCol := t.PortBorder
		if port.ColorHex != 0 {
			borderCol = RGBHex(port.ColorHex)
		}
		if snapped {
			borderCol = t.PortSnap
		}
		prRect := ui.Rect{X: sx - pr, Y: sy - pr, W: pr * 2, H: pr * 2}
		p.Fill(prRect, t.PortFill, pr)
		p.Stroke(prRect, borderCol, pr, 1)

		// React Flow UI ButtonHandle: stem line + circular '+' button to spawn & wire next node.
		if port.ShowButton && !g.ReadOnly && z >= 0.45 {
			sn := sideNormal(side)
			stem := 26 * z
			bx := sx + sn.X*stem
			by := sy + sn.Y*stem
			p.Line(sx+sn.X*pr, sy+sn.Y*pr, bx, by, 1.2, t.NodeBorder)
			br := clamp32(8.5*z, 5.5, 10)
			brect := ui.Rect{X: bx - br, Y: by - br, W: br * 2, H: br * 2}
			p.Shadow(brect, br, 0, 2, 6, 0, colorWithAlpha(ui.RGB(0, 0, 0), 0.14))
			p.Fill(brect, t.NodeBG, br)
			p.Stroke(brect, t.NodeBorder, br, 1.1)
			arm := br * 0.42
			p.Line(bx-arm, by, bx+arm, by, 1.3, t.TextPrimary)
			p.Line(bx, by-arm, bx, by+arm, 1.3, t.TextPrimary)
			g.handleBtnHits = append(g.handleBtnHits, handleButtonHit{
				NodeID: n.ID,
				PortID: port.ID,
				Center: Point{X: bx - r.X, Y: by - r.Y},
				Radius: br,
			})
		}

		// For BaseNode (not SchemaNode which already shows column names inline), show crisp handle labels outside.
		if n.Variant != VariantSchema && port.Label != "" && z >= 0.72 && (side == PortLeft || side == PortRight) {
			ls := clamp32(8.5*z, 7.5, 10.5)
			if side == PortLeft {
				lw := float32(len(port.Label)) * ls * 0.56
				p.Text(sx-pr-4-lw, sy-pr-ls-1, port.Label, ls, t.TextMuted)
			} else {
				p.Text(sx+pr+4, sy-pr-ls-1, port.Label, ls, t.TextMuted)
			}
		}
	}
}

func (g *Graph) paintMiniMap(p *ui.Painter, r ui.Rect, t Theme) {
	const (
		mw     = float32(174)
		mh     = float32(110)
		margin = float32(14)
	)
	mr := ui.Rect{
		X: r.X + r.W - mw - margin,
		Y: r.Y + r.H - mh - margin,
		W: mw,
		H: mh,
	}
	g.miniMapRect = Rect{X: mr.X - r.X, Y: mr.Y - r.Y, W: mw, H: mh}

	p.Shadow(mr, 10, 0, 6, 18, 0, colorWithAlpha(ui.RGB(0, 0, 0), 0.18))
	p.Fill(mr, t.MiniMapBG, 10)
	p.Stroke(mr, t.MiniMapBorder, 10, 1)

	_, scale, offX, offY := g.miniMapTransform(r, mr)

	p.Clip(mr, 10, func() {
		for _, e := range g.Edges {
			sn := g.NodeByID(e.Source)
			dn := g.NodeByID(e.Target)
			if sn == nil || dn == nil {
				continue
			}
			sb, db := sn.Bounds(), dn.Bounds()
			x0 := mr.X + offX + (sb.X+sb.W*0.5)*scale
			y0 := mr.Y + offY + (sb.Y+sb.H*0.5)*scale
			x1 := mr.X + offX + (db.X+db.W*0.5)*scale
			y1 := mr.Y + offY + (db.Y+db.H*0.5)*scale
			p.Line(x0, y0, x1, y1, 1, colorWithAlpha(t.EdgeDefault, 0.55))
		}
		for _, n := range g.Nodes {
			b := n.Bounds()
			col := t.MiniMapNode
			if n.Variant == VariantGroup {
				col = colorWithAlpha(t.MiniMapNode, 0.25)
			} else if n.Selected {
				col = t.MiniMapNodeSel
			} else if n.AccentHex != 0 {
				col = colorWithAlpha(RGBHex(n.AccentHex), 0.72)
			}
			nr := ui.Rect{
				X: mr.X + offX + b.X*scale,
				Y: mr.Y + offY + b.Y*scale,
				W: max32(b.W*scale, 4),
				H: max32(b.H*scale, 3),
			}
			p.Fill(nr, col, 2)
		}
		vTL := g.Viewport.ScreenToWorld(Point{X: 0, Y: 0})
		vBR := g.Viewport.ScreenToWorld(Point{X: r.W, Y: r.H})
		vr := ui.Rect{
			X: mr.X + offX + vTL.X*scale,
			Y: mr.Y + offY + vTL.Y*scale,
			W: (vBR.X - vTL.X) * scale,
			H: (vBR.Y - vTL.Y) * scale,
		}
		p.Fill(vr, t.MiniMapViewport, 3)
		p.Stroke(vr, t.MiniMapViewEdge, 3, 1.2)
	})
}

func (g *Graph) miniMapTransform(canvasR, miniR ui.Rect) (Rect, float32, float32, float32) {
	vTL := g.Viewport.ScreenToWorld(Point{X: 0, Y: 0})
	vBR := g.Viewport.ScreenToWorld(Point{X: canvasR.W, Y: canvasR.H})
	minX, minY, maxX, maxY := vTL.X, vTL.Y, vBR.X, vBR.Y
	if cb, ok := g.ContentBounds(); ok {
		minX = min32(minX, cb.X)
		minY = min32(minY, cb.Y)
		maxX = max32(maxX, cb.X+cb.W)
		maxY = max32(maxY, cb.Y+cb.H)
	}
	pad := float32(80)
	minX -= pad
	minY -= pad
	maxX += pad
	maxY += pad
	w := max32(maxX-minX, 1)
	h := max32(maxY-minY, 1)

	innerPad := float32(10)
	sx := (miniR.W - innerPad*2) / w
	sy := (miniR.H - innerPad*2) / h
	scale := min32(sx, sy)
	offX := (miniR.W-w*scale)*0.5 - minX*scale
	offY := (miniR.H-h*scale)*0.5 - minY*scale
	return Rect{X: minX, Y: minY, W: w, H: h}, scale, offX, offY
}

func (g *Graph) panToMiniMapPoint(screenPt Point) {
	if g.miniMapRect.W <= 0 || g.lastCanvasW <= 0 || g.lastCanvasH <= 0 {
		return
	}
	canvasR := ui.Rect{W: g.lastCanvasW, H: g.lastCanvasH}
	miniR := ui.Rect{X: g.miniMapRect.X, Y: g.miniMapRect.Y, W: g.miniMapRect.W, H: g.miniMapRect.H}
	_, scale, offX, offY := g.miniMapTransform(canvasR, miniR)
	if scale <= 0 {
		return
	}
	wx := (screenPt.X - miniR.X - offX) / scale
	wy := (screenPt.Y - miniR.Y - offY) / scale
	z := g.Viewport.safeZoom()
	g.Viewport.X = g.lastCanvasW*0.5 - wx*z
	g.Viewport.Y = g.lastCanvasH*0.5 - wy*z
}

func statusColor(s NodeStatus, t Theme) ui.Color {
	switch s {
	case StatusRunning:
		return t.StatusRunning
	case StatusSuccess:
		return t.StatusSuccess
	case StatusError:
		return t.StatusError
	default:
		return t.StatusIdle
	}
}

func normalizedRect(a, b Point) Rect {
	x0, x1 := a.X, b.X
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	y0, y1 := a.Y, b.Y
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rect{X: x0, Y: y0, W: x1 - x0, H: y1 - y0}
}

func truncateStr(s string, maxChars int) string {
	r := []rune(s)
	if len(r) <= maxChars {
		return s
	}
	if maxChars <= 1 {
		return string(r[:1])
	}
	return string(r[:maxChars-1]) + "…"
}
