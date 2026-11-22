package flow

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"
)

// CanvasOptions configures rendering and overlays of a flowgo Canvas.
type CanvasOptions struct {
	Theme        Theme
	NodeRenderer func(p *ui.Painter, n *Node, r ui.Rect, zoom float32, t Theme)
}

// Canvas renders an interactive, GPU-accelerated infinite node graph canvas
// complete with pan/zoom, grid background, animated edges, port connecting,
// multi-node drag with alignment guides, box selection, minimap, and shortcuts.
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

		// Floating bottom-left zoom / fit / snap / undo controls overlay.
		if g.ShowControls {
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

// ControlsBar renders a compact floating toolbar for zoom, fit, mode, snap, and undo/redo.
func ControlsBar(c *ui.Context, g *Graph, t Theme) {
	ui.Row(c).
		Background(t.ControlsBG).
		Border(1, t.ControlsBorder).
		Radius(10).
		Padding(4, 6).
		Gap(4).
		AlignItems(ui.Center).
		Shadow(0, 6, 20, 0, colorWithAlpha(ui.RGB(0, 0, 0), 0.25)).
		Children(func() {
			ctrlBtn(c, t, "-", "Zoom out", false, func() { g.ZoomOut() })
			ui.Text(c, fmt.Sprintf("%d%%", int(math.Round(float64(g.Viewport.safeZoom()*100))))).
				FontSize(11).
				FontWeight(600).
				TextColor(t.TextMuted).
				Width(38).
				TextAlign(ui.Center)
			ctrlBtn(c, t, "+", "Zoom in", false, func() { g.ZoomIn() })
			ctrlSep(c, t)
			ctrlBtn(c, t, "Fit", "Fit view (0)", false, func() { g.FitView(0, 0, 56) })
			ctrlBtn(c, t, "Box", "Box select mode", g.BoxSelectMode, func() { g.BoxSelectMode = !g.BoxSelectMode })
			ctrlBtn(c, t, "Snap", "Snap to grid", g.SnapToGrid, func() { g.SnapToGrid = !g.SnapToGrid })
			ctrlSep(c, t)
			ctrlBtn(c, t, "Undo", "Undo", false, func() { g.Undo() })
			ctrlBtn(c, t, "Redo", "Redo", false, func() { g.Redo() })
		})
}

func ctrlSep(c *ui.Context, t Theme) {
	ui.Box(c).Width(1).Height(16).Background(t.ControlsBorder).Margin(0, 2)
}

func ctrlBtn(c *ui.Context, t Theme, label, tip string, active bool, onClick func()) {
	b := ui.ButtonBase(c).
		Height(26).
		Padding(0, 8).
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
		// Mouse wheel (non-precise) zooms; smooth trackpad scrolls pan cleanly.
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
		if g.mode == modeConnect {
			g.mode = modeIdle
			return true
		}
		g.ClearSelection()
		return true
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

	// Check MiniMap hit first.
	if g.ShowMiniMap && g.miniMapRect.W > 0 && g.miniMapRect.Contains(ptScreen) {
		g.mode = modeMiniMap
		g.panToMiniMapPoint(ptScreen)
		return true
	}

	// Middle or right mouse button always pans.
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

	// 2. Check Node hit to select and drag.
	if n := g.HitTestNode(ptWorld); n != nil {
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
		// Raise clicked node to top of z-order.
		g.bringToFront(n.ID)
		if !g.ReadOnly && !n.Locked {
			g.mode = modeDragNodes
			g.dragMoved = false
			g.dragStartPos = map[string]Point{}
			for _, sel := range g.SelectedNodes() {
				if !sel.Locked {
					g.dragStartPos[sel.ID] = Point{X: sel.X, Y: sel.Y}
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

	// 4. Empty background: box-select (Shift or BoxSelectMode) or pan.
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
			// Compute optional alignment snapping against unselected nodes when single node dragged.
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
						if other.ID == id {
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
		// Idle hover hit-testing for crisp cursor & port feedback.
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

	// 2. Alignment guides during drag.
	for _, guide := range g.guides {
		if guide.Vertical {
			sp := g.Viewport.WorldToScreen(Point{X: guide.Pos, Y: 0})
			p.Line(r.X+sp.X, r.Y, r.X+sp.X, r.Y+r.H, 1, t.GuideLine)
		} else {
			sp := g.Viewport.WorldToScreen(Point{X: 0, Y: guide.Pos})
			p.Line(r.X, r.Y+sp.Y, r.X+r.W, r.Y+sp.Y, 1, t.GuideLine)
		}
	}

	// 3. Edges (with animated flowing dashes when Animated == true).
	needsAnim := false
	nowSec := float64(p.Now().UnixNano()) / 1e9
	for _, e := range g.Edges {
		if e.Animated {
			needsAnim = true
		}
		g.paintEdge(p, r, e, t, nowSec)
	}

	// 4. Live connection preview wire.
	if g.mode == modeConnect {
		needsAnim = true
		g.paintConnectPreview(p, r, t, nowSec)
	}

	// 5. Nodes.
	for _, n := range g.Nodes {
		if n.Status == StatusRunning {
			needsAnim = true
		}
		g.paintNode(p, r, n, t, nowSec, customNode)
	}

	// 6. Box selection marquee.
	if g.mode == modeBoxSelect {
		wRect := normalizedRect(g.pressWorld, g.boxCurWorld)
		tl := g.Viewport.WorldToScreen(Point{X: wRect.X, Y: wRect.Y})
		br := g.Viewport.WorldToScreen(Point{X: wRect.X + wRect.W, Y: wRect.Y + wRect.H})
		sr := ui.Rect{X: r.X + tl.X, Y: r.Y + tl.Y, W: br.X - tl.X, H: br.Y - tl.Y}
		p.Fill(sr, t.SelectionFill, 4)
		p.Stroke(sr, t.SelectionBorder, 4, 1.25)
	}

	// 7. Interactive MiniMap (bottom-right corner; hidden on narrow viewports so it never overlaps ControlsBar).
	if g.ShowMiniMap && r.W >= 560 && r.H >= 280 {
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
	default: // BackgroundDots
		dotSize := clamp32(2.0*z, 1.5, 3.2)
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
	width := clamp32(2.0*g.Viewport.safeZoom(), 1.3, 3.2)
	if e.ID == g.hoverEdgeID {
		col = t.EdgeHover
		width += 0.6
	}
	if e.Selected {
		col = t.EdgeSelected
		width += 0.9
		// Subtle selection glow halo.
		p.StrokePath(geom.Path, width+5, colorWithAlpha(t.EdgeSelected, 0.25))
	}

	if e.Animated {
		// Base track + animated flowing dashed highlight.
		p.StrokePath(geom.Path, width, colorWithAlpha(col, 0.35))
		animPath := BuildEdgeGeometry(sp, dp, spSide, dpSide, kind).Path
		dashOffset := float32(-math.Mod(nowSec*36, 16))
		animPath.Dash(dashOffset, 8, 8)
		animCol := col
		if !e.Selected && e.ColorHex == 0 {
			animCol = t.NodeSelected
		}
		p.StrokePath(animPath, width+0.3, animCol)
	} else {
		p.StrokePath(geom.Path, width, col)
	}

	// Arrowhead right outside destination port.
	arrowSize := clamp32(8.5*g.Viewport.safeZoom(), 5, 11)
	tipBack := clamp32(6.5*g.Viewport.safeZoom(), 4, 9)
	tip := Point{
		X: geom.Tip.X - float32(math.Cos(float64(geom.Angle)))*tipBack,
		Y: geom.Tip.Y - float32(math.Sin(float64(geom.Angle)))*tipBack,
	}
	p.FillPath(ArrowPath(tip, geom.Angle, arrowSize), col)

	// Centered edge label pill.
	if e.Label != "" && g.Viewport.safeZoom() >= 0.45 {
		fontSize := clamp32(10.5*g.Viewport.safeZoom(), 8.5, 13)
		pillW := float32(len(e.Label))*fontSize*0.62 + 16
		pillH := fontSize + 10
		lr := ui.Rect{
			X: geom.LabelPos.X - pillW*0.5,
			Y: geom.LabelPos.Y - pillH*0.5,
			W: pillW,
			H: pillH,
		}
		p.Fill(lr, t.EdgeLabelBG, pillH*0.5)
		p.Stroke(lr, t.EdgeLabelBorder, pillH*0.5, 1)
		p.Text(lr.X+8, lr.Y+(pillH-fontSize)*0.5-1, e.Label, fontSize, t.EdgeLabelText)
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
	p.StrokePath(geom.Path, 2.2, wireCol)
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
	// Frustum cull off-screen nodes.
	if nr.X+nr.W < r.X-40 || nr.Y+nr.H < r.Y-40 || nr.X > r.X+r.W+40 || nr.Y > r.Y+r.H+40 {
		return
	}

	if custom != nil {
		custom(p, n, nr, z, t)
		g.paintNodePorts(p, r, n, t, z)
		return
	}

	rad := clamp32(12*z, 5, 16)
	accent := t.NodeSelected
	if n.AccentHex != 0 {
		accent = RGBHex(n.AccentHex)
	}

	// Elevation shadow.
	p.Shadow(nr, rad, 0, clamp32(6*z, 2, 12), clamp32(18*z, 6, 28), 0, t.NodeShadow)

	// Selected / Running halo ring.
	if n.Status == StatusRunning {
		pulse := float32(0.5 + 0.5*math.Sin(nowSec*6.0))
		halo := ui.Rect{X: nr.X - 3, Y: nr.Y - 3, W: nr.W + 6, H: nr.H + 6}
		p.Stroke(halo, colorWithAlpha(t.StatusRunning, 0.35+0.45*pulse), rad+3, 2.2)
	} else if n.Selected {
		halo := ui.Rect{X: nr.X - 2.5, Y: nr.Y - 2.5, W: nr.W + 5, H: nr.H + 5}
		p.Stroke(halo, colorWithAlpha(t.NodeSelected, 0.32), rad+2.5, 3.0)
	}

	// Card body.
	p.Fill(nr, t.NodeBG, rad)

	// Header strip + top accent bar inside clipped card.
	headerH := clamp32(34*z, 18, 44)
	p.Clip(nr, rad, func() {
		p.Fill(ui.Rect{X: nr.X, Y: nr.Y, W: nr.W, H: headerH}, t.NodeHeaderBG, 0)
		p.Fill(ui.Rect{X: nr.X, Y: nr.Y, W: nr.W, H: clamp32(3*z, 2, 4)}, accent, 0)
		p.Line(nr.X, nr.Y+headerH, nr.X+nr.W, nr.Y+headerH, 1, t.NodeBorder)

		if z >= 0.42 {
			// Status dot.
			dotCol := statusColor(n.Status, t)
			dotD := clamp32(7.5*z, 4, 9)
			p.Fill(ui.Rect{
				X: nr.X + 11*z,
				Y: nr.Y + headerH*0.5 - dotD*0.5 + 1,
				W: dotD,
				H: dotD,
			}, dotCol, dotD*0.5)

			// Title.
			titleSize := clamp32(12.5*z, 8, 16)
			p.Text(nr.X+24*z, nr.Y+(headerH-titleSize)*0.5+0.5, truncateStr(n.Title, 22), titleSize, t.TextPrimary)

			// Badge pill on top-right of header.
			if n.Badge != "" && nr.W >= 92 {
				bs := clamp32(9.5*z, 7.5, 11.5)
				bw := float32(len(n.Badge))*bs*0.62 + 12*z
				bh := bs + 6*z
				br := ui.Rect{
					X: nr.X + nr.W - bw - 8*z,
					Y: nr.Y + (headerH-bh)*0.5 + 1,
					W: bw,
					H: bh,
				}
				p.Fill(br, colorWithAlpha(accent, 0.18), bh*0.5)
				p.Text(br.X+6*z, br.Y+(bh-bs)*0.5-0.5, n.Badge, bs, accent)
			}

			// Body subtitle & fields.
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
		}
	})

	// Outer border.
	borderCol := t.NodeBorder
	borderW := float32(1)
	if n.Selected {
		borderCol = t.NodeSelected
		borderW = 1.6
	} else if n.ID == g.hoverNodeID {
		borderCol = colorWithAlpha(accent, 0.65)
		borderW = 1.25
	}
	p.Stroke(nr, borderCol, rad, borderW)

	// Ports.
	g.paintNodePorts(p, r, n, t, z)
}

func (g *Graph) paintNodePorts(p *ui.Painter, r ui.Rect, n *Node, t Theme, z float32) {
	for i := range n.Ports {
		port := &n.Ports[i]
		wPos, side, _ := PortPosition(n, port.ID)
		sPos := g.Viewport.WorldToScreen(wPos)
		sx, sy := r.X+sPos.X, r.Y+sPos.Y

		hovered := g.hoverPortNode == n.ID && g.hoverPortID == port.ID
		snapped := g.mode == modeConnect && g.connect.HasSnap && g.connect.SnapNode == n.ID && g.connect.SnapPort == port.ID
		pr := clamp32(5.2*z, 3.2, 7)
		if hovered || snapped {
			pr += 1.8
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
		p.Stroke(prRect, borderCol, pr, 2.0)

		if port.Label != "" && z >= 0.72 && (side == PortLeft || side == PortRight) {
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
		mw     = float32(178)
		mh     = float32(114)
		margin = float32(14)
	)
	mr := ui.Rect{
		X: r.X + r.W - mw - margin,
		Y: r.Y + r.H - mh - margin,
		W: mw,
		H: mh,
	}
	g.miniMapRect = Rect{X: mr.X - r.X, Y: mr.Y - r.Y, W: mw, H: mh}

	p.Shadow(mr, 10, 0, 6, 18, 0, colorWithAlpha(ui.RGB(0, 0, 0), 0.25))
	p.Fill(mr, t.MiniMapBG, 10)
	p.Stroke(mr, t.MiniMapBorder, 10, 1)

	worldB, scale, offX, offY := g.miniMapTransform(r, mr)
	_ = worldB

	p.Clip(mr, 10, func() {
		// Mini edges.
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
		// Mini nodes.
		for _, n := range g.Nodes {
			b := n.Bounds()
			col := t.MiniMapNode
			if n.Selected {
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
		// Current viewport rectangle.
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
