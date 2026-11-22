package flow

import (
	"encoding/json"
	"fmt"
	"math"
)

type interactMode int

const (
	modeIdle interactMode = iota
	modePan
	modeDragNodes
	modeBoxSelect
	modeConnect
	modeMiniMap
)

type alignmentGuide struct {
	Vertical bool
	Pos      float32
}

type connectPreview struct {
	FromNode   string
	FromPort   string
	FromSide   PortSide
	FromKind   PortKind
	FromWorld  Point
	CurWorld   Point
	SnapNode   string
	SnapPort   string
	SnapWorld  Point
	HasSnap    bool
}

type edgeButtonHit struct {
	EdgeID string
	Center Point
	Radius float32
}

type graphSnapshot struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// Graph holds the complete node-edge document, viewport, preferences, history,
// and live canvas interaction state.
type Graph struct {
	Nodes       []*Node
	Edges       []*Edge
	Viewport    Viewport
	SnapToGrid  bool
	GridSize    float32
	EdgeKind    EdgeKind
	Background  BackgroundKind
	ShowMiniMap bool
	ShowControls bool
	BoxSelectMode bool
	ReadOnly    bool
	MinZoom     float32
	MaxZoom     float32

	OnChange      func()
	OnConnect     func(e *Edge)
	OnSelectNode  func(n *Node)
	OnSelectEdge  func(e *Edge)
	OnDoubleClick func(n *Node)

	// Internal interaction state (unexported).
	mode          interactMode
	pressScreen   Point
	lastScreen    Point
	pressWorld    Point
	boxCurWorld   Point
	dragStartPos  map[string]Point
	dragMoved     bool
	hoverNodeID   string
	hoverPortNode string
	hoverPortID   string
	hoverEdgeID   string
	connect       connectPreview
	guides        []alignmentGuide
	lastCanvasW   float32
	lastCanvasH   float32
	miniMapRect   Rect
	edgeBtnHits   []edgeButtonHit
	nextID        int

	undoStack []graphSnapshot
	redoStack []graphSnapshot
}

// NewGraph creates an initialized Graph with sensible defaults.
func NewGraph() *Graph {
	return &Graph{
		Viewport:     DefaultViewport(),
		SnapToGrid:   true,
		GridSize:     20,
		EdgeKind:     EdgeBezier,
		Background:   BackgroundDots,
		ShowMiniMap:  true,
		ShowControls: true,
		MinZoom:      0.2,
		MaxZoom:      2.5,
		nextID:       100,
	}
}

// NodeByID returns the node with the given ID or nil.
func (g *Graph) NodeByID(id string) *Node {
	for _, n := range g.Nodes {
		if n.ID == id {
			return n
		}
	}
	return nil
}

// EdgeByID returns the edge with the given ID or nil.
func (g *Graph) EdgeByID(id string) *Edge {
	for _, e := range g.Edges {
		if e.ID == id {
			return e
		}
	}
	return nil
}

// AddNode inserts a node onto the canvas and records an undo snapshot.
func (g *Graph) AddNode(n *Node) *Node {
	if n == nil {
		return nil
	}
	if n.ID == "" {
		g.nextID++
		n.ID = fmt.Sprintf("node-%d", g.nextID)
	}
	if n.W <= 0 {
		n.W = 224
	}
	if n.H <= 0 {
		n.H = 112
	}
	g.PushUndo()
	g.Nodes = append(g.Nodes, n)
	g.notifyChange()
	return n
}

// RemoveEdge removes a single edge by ID.
func (g *Graph) RemoveEdge(id string) bool {
	if g.ReadOnly {
		return false
	}
	for i, e := range g.Edges {
		if e.ID == id {
			g.PushUndo()
			g.Edges = append(g.Edges[:i], g.Edges[i+1:]...)
			g.notifyChange()
			return true
		}
	}
	return false
}

// PromotePlaceholder turns a PlaceholderNode into a regular BaseNode on click
// (matching React Flow UI's PlaceholderNode behavior).
func (g *Graph) PromotePlaceholder(id string) bool {
	if g.ReadOnly {
		return false
	}
	n := g.NodeByID(id)
	if n == nil || n.Variant != VariantPlaceholder {
		return false
	}
	g.PushUndo()
	n.Variant = VariantBase
	n.Type = "node"
	n.Title = "New Step"
	n.Subtitle = "Configured from placeholder"
	n.Badge = "STEP"
	n.W = 216
	n.H = 104
	if len(n.Ports) == 0 {
		n.Ports = []Port{
			{ID: "in", Side: PortLeft, Kind: PortInput},
			{ID: "out", Side: PortRight, Kind: PortOutput},
		}
	}
	for _, e := range g.Edges {
		if e.Target == id {
			e.Animated = false
		}
	}
	g.SelectNode(n.ID, false)
	g.notifyChange()
	return true
}

// RemoveNode removes a node and all incident edges.
func (g *Graph) RemoveNode(id string) bool {
	idx := -1
	for i, n := range g.Nodes {
		if n.ID == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return false
	}
	g.PushUndo()
	g.Nodes = append(g.Nodes[:idx], g.Nodes[idx+1:]...)
	kept := g.Edges[:0]
	for _, e := range g.Edges {
		if e.Source != id && e.Target != id {
			kept = append(kept, e)
		}
	}
	g.Edges = kept
	g.notifyChange()
	return true
}

// Connect creates an edge if valid (no self-loop, no duplicate, matching kinds).
func (g *Graph) Connect(srcNode, srcPort, dstNode, dstPort string) (*Edge, bool) {
	if !g.CanConnect(srcNode, srcPort, dstNode, dstPort) {
		return nil, false
	}
	g.PushUndo()
	g.nextID++
	e := &Edge{
		ID:         fmt.Sprintf("edge-%d", g.nextID),
		Source:     srcNode,
		SourcePort: srcPort,
		Target:     dstNode,
		TargetPort: dstPort,
		Kind:       g.EdgeKind,
	}
	g.Edges = append(g.Edges, e)
	if g.OnConnect != nil {
		g.OnConnect(e)
	}
	g.notifyChange()
	return e, true
}

// CanConnect validates whether two ports can form a legal edge.
func (g *Graph) CanConnect(srcNode, srcPort, dstNode, dstPort string) bool {
	if srcNode == "" || dstNode == "" || srcNode == dstNode {
		return false
	}
	sn := g.NodeByID(srcNode)
	dn := g.NodeByID(dstNode)
	if sn == nil || dn == nil {
		return false
	}
	sp := findPort(sn, srcPort)
	dp := findPort(dn, dstPort)
	if sp == nil || dp == nil {
		return false
	}
	if sp.Kind == dp.Kind {
		return false
	}
	if sp.DataType != "" && dp.DataType != "" && sp.DataType != dp.DataType {
		return false
	}
	for _, e := range g.Edges {
		if e.Source == srcNode && e.SourcePort == srcPort && e.Target == dstNode && e.TargetPort == dstPort {
			return false
		}
	}
	return true
}

func findPort(n *Node, id string) *Port {
	for i := range n.Ports {
		if n.Ports[i].ID == id {
			return &n.Ports[i]
		}
	}
	return nil
}

// ClearSelection deselects all nodes and edges.
func (g *Graph) ClearSelection() {
	for _, n := range g.Nodes {
		n.Selected = false
	}
	for _, e := range g.Edges {
		e.Selected = false
	}
}

// SelectAll selects every node in the graph.
func (g *Graph) SelectAll() {
	for _, n := range g.Nodes {
		n.Selected = true
	}
	for _, e := range g.Edges {
		e.Selected = false
	}
}

// SelectNode selects a node by ID; additive preserves previous selection.
func (g *Graph) SelectNode(id string, additive bool) {
	if !additive {
		g.ClearSelection()
	}
	var picked *Node
	for _, n := range g.Nodes {
		if n.ID == id {
			if additive {
				n.Selected = !n.Selected
			} else {
				n.Selected = true
			}
			if n.Selected {
				picked = n
			}
		}
	}
	if g.OnSelectNode != nil {
		g.OnSelectNode(picked)
	}
}

// SelectEdge selects an edge by ID.
func (g *Graph) SelectEdge(id string) {
	g.ClearSelection()
	var picked *Edge
	for _, e := range g.Edges {
		if e.ID == id {
			e.Selected = true
			picked = e
		}
	}
	if g.OnSelectEdge != nil {
		g.OnSelectEdge(picked)
	}
}

// SelectedNodes returns all currently selected nodes.
func (g *Graph) SelectedNodes() []*Node {
	var out []*Node
	for _, n := range g.Nodes {
		if n.Selected {
			out = append(out, n)
		}
	}
	return out
}

// SelectedEdge returns the first selected edge, or nil.
func (g *Graph) SelectedEdge() *Edge {
	for _, e := range g.Edges {
		if e.Selected {
			return e
		}
	}
	return nil
}

// DeleteSelected removes all selected unlocked nodes and selected edges.
func (g *Graph) DeleteSelected() bool {
	if g.ReadOnly {
		return false
	}
	removedNodes := map[string]bool{}
	any := false
	for _, n := range g.Nodes {
		if n.Selected && !n.Locked {
			removedNodes[n.ID] = true
			any = true
		}
	}
	for _, e := range g.Edges {
		if e.Selected {
			any = true
		}
	}
	if !any {
		return false
	}
	g.PushUndo()

	keptNodes := g.Nodes[:0]
	for _, n := range g.Nodes {
		if !removedNodes[n.ID] {
			keptNodes = append(keptNodes, n)
		}
	}
	g.Nodes = keptNodes

	keptEdges := g.Edges[:0]
	for _, e := range g.Edges {
		if e.Selected || removedNodes[e.Source] || removedNodes[e.Target] {
			continue
		}
		keptEdges = append(keptEdges, e)
	}
	g.Edges = keptEdges
	g.notifyChange()
	return true
}

// DuplicateSelected clones selected nodes with a 36px offset and selects the clones.
func (g *Graph) DuplicateSelected() int {
	if g.ReadOnly {
		return 0
	}
	sel := g.SelectedNodes()
	if len(sel) == 0 {
		return 0
	}
	g.PushUndo()
	g.ClearSelection()
	idMap := make(map[string]string, len(sel))
	for _, n := range sel {
		g.nextID++
		cp := cloneNode(*n)
		cp.ID = fmt.Sprintf("node-%d", g.nextID)
		cp.X += 36
		cp.Y += 36
		cp.Selected = true
		idMap[n.ID] = cp.ID
		g.Nodes = append(g.Nodes, &cp)
	}
	// Duplicate edges internal to the selection.
	for _, e := range g.Edges {
		ns, okS := idMap[e.Source]
		nt, okT := idMap[e.Target]
		if okS && okT {
			g.nextID++
			ce := *e
			ce.ID = fmt.Sprintf("edge-%d", g.nextID)
			ce.Source = ns
			ce.Target = nt
			ce.Selected = false
			g.Edges = append(g.Edges, &ce)
		}
	}
	g.notifyChange()
	return len(sel)
}

// ContentBounds computes the world-space bounding box of all nodes.
func (g *Graph) ContentBounds() (Rect, bool) {
	if len(g.Nodes) == 0 {
		return Rect{}, false
	}
	minX := float32(math.Inf(1))
	minY := float32(math.Inf(1))
	maxX := float32(math.Inf(-1))
	maxY := float32(math.Inf(-1))
	for _, n := range g.Nodes {
		b := n.Bounds()
		if b.X < minX {
			minX = b.X
		}
		if b.Y < minY {
			minY = b.Y
		}
		if b.X+b.W > maxX {
			maxX = b.X + b.W
		}
		if b.Y+b.H > maxY {
			maxY = b.Y + b.H
		}
	}
	return Rect{X: minX, Y: minY, W: maxX - minX, H: maxY - minY}, true
}

// FitView centers and scales the viewport so every node fits inside canvasW x canvasH.
func (g *Graph) FitView(canvasW, canvasH, padding float32) {
	if canvasW <= 0 {
		canvasW = g.lastCanvasW
	}
	if canvasH <= 0 {
		canvasH = g.lastCanvasH
	}
	if canvasW <= 0 || canvasH <= 0 {
		return
	}
	b, ok := g.ContentBounds()
	if !ok {
		g.Viewport = DefaultViewport()
		return
	}
	availW := canvasW - padding*2
	availH := canvasH - padding*2
	if availW < 40 {
		availW = canvasW
	}
	if availH < 40 {
		availH = canvasH
	}
	zx := availW / max32(b.W, 1)
	zy := availH / max32(b.H, 1)
	zoom := clamp32(min32(zx, zy), g.MinZoom, 1.15)
	cx := b.X + b.W*0.5
	cy := b.Y + b.H*0.5
	g.Viewport = Viewport{
		X:    canvasW*0.5 - cx*zoom,
		Y:    canvasH*0.5 - cy*zoom,
		Zoom: zoom,
	}
}

// ZoomIn zooms in around the canvas center.
func (g *Graph) ZoomIn() {
	g.Viewport.ZoomAt(Point{X: g.lastCanvasW * 0.5, Y: g.lastCanvasH * 0.5}, 1.2, g.MinZoom, g.MaxZoom)
}

// ZoomOut zooms out around the canvas center.
func (g *Graph) ZoomOut() {
	g.Viewport.ZoomAt(Point{X: g.lastCanvasW * 0.5, Y: g.lastCanvasH * 0.5}, 1/1.2, g.MinZoom, g.MaxZoom)
}

// ZoomTo sets an exact zoom level around the canvas center (used by ZoomSlider).
func (g *Graph) ZoomTo(target float32) {
	cur := g.Viewport.safeZoom()
	if cur <= 0 {
		cur = 1
	}
	target = clamp32(target, g.MinZoom, g.MaxZoom)
	g.Viewport.ZoomAt(Point{X: g.lastCanvasW * 0.5, Y: g.lastCanvasH * 0.5}, target/cur, g.MinZoom, g.MaxZoom)
}

// PushUndo saves the current graph state onto the undo stack.
func (g *Graph) PushUndo() {
	const maxHistory = 60
	g.undoStack = append(g.undoStack, g.snapshot())
	if len(g.undoStack) > maxHistory {
		g.undoStack = g.undoStack[len(g.undoStack)-maxHistory:]
	}
	g.redoStack = nil
}

func (g *Graph) CanUndo() bool { return len(g.undoStack) > 0 }
func (g *Graph) CanRedo() bool { return len(g.redoStack) > 0 }

// Undo reverts the graph to the previous snapshot.
func (g *Graph) Undo() bool {
	if len(g.undoStack) == 0 {
		return false
	}
	g.redoStack = append(g.redoStack, g.snapshot())
	prev := g.undoStack[len(g.undoStack)-1]
	g.undoStack = g.undoStack[:len(g.undoStack)-1]
	g.restore(prev)
	g.notifyChange()
	return true
}

// Redo re-applies the last undone snapshot.
func (g *Graph) Redo() bool {
	if len(g.redoStack) == 0 {
		return false
	}
	g.undoStack = append(g.undoStack, g.snapshot())
	next := g.redoStack[len(g.redoStack)-1]
	g.redoStack = g.redoStack[:len(g.redoStack)-1]
	g.restore(next)
	g.notifyChange()
	return true
}

// SaveJSON serializes nodes, edges, and viewport to formatted JSON.
func (g *Graph) SaveJSON() ([]byte, error) {
	doc := struct {
		Viewport Viewport `json:"viewport"`
		Nodes    []Node   `json:"nodes"`
		Edges    []Edge   `json:"edges"`
	}{
		Viewport: g.Viewport,
		Nodes:    g.snapshot().Nodes,
		Edges:    g.snapshot().Edges,
	}
	return json.MarshalIndent(doc, "", "  ")
}

// LoadJSON replaces the graph contents from JSON produced by SaveJSON.
func (g *Graph) LoadJSON(data []byte) error {
	var doc struct {
		Viewport Viewport `json:"viewport"`
		Nodes    []Node   `json:"nodes"`
		Edges    []Edge   `json:"edges"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	g.PushUndo()
	g.restore(graphSnapshot{Nodes: doc.Nodes, Edges: doc.Edges})
	if doc.Viewport.Zoom > 0 {
		g.Viewport = doc.Viewport
	}
	g.notifyChange()
	return nil
}

func (g *Graph) snapshot() graphSnapshot {
	nodes := make([]Node, len(g.Nodes))
	for i, n := range g.Nodes {
		nodes[i] = cloneNode(*n)
	}
	edges := make([]Edge, len(g.Edges))
	for i, e := range g.Edges {
		edges[i] = *e
	}
	return graphSnapshot{Nodes: nodes, Edges: edges}
}

func (g *Graph) restore(s graphSnapshot) {
	g.Nodes = make([]*Node, len(s.Nodes))
	for i := range s.Nodes {
		cp := cloneNode(s.Nodes[i])
		g.Nodes[i] = &cp
	}
	g.Edges = make([]*Edge, len(s.Edges))
	for i := range s.Edges {
		cp := s.Edges[i]
		g.Edges[i] = &cp
	}
}

func cloneNode(n Node) Node {
	if len(n.Ports) > 0 {
		ports := make([]Port, len(n.Ports))
		copy(ports, n.Ports)
		n.Ports = ports
	}
	if len(n.Fields) > 0 {
		fields := make([]NodeField, len(n.Fields))
		copy(fields, n.Fields)
		n.Fields = fields
	}
	if len(n.SchemaRows) > 0 {
		rows := make([]SchemaRow, len(n.SchemaRows))
		copy(rows, n.SchemaRows)
		n.SchemaRows = rows
	}
	return n
}

func (g *Graph) notifyChange() {
	if g.OnChange != nil {
		g.OnChange()
	}
}

// HitTestPort returns the node and port closest to worldPt within radius.
func (g *Graph) HitTestPort(worldPt Point, radius float32) (*Node, *Port, Point) {
	bestDist := radius
	var bestN *Node
	var bestP *Port
	var bestPos Point
	for i := len(g.Nodes) - 1; i >= 0; i-- {
		n := g.Nodes[i]
		for j := range n.Ports {
			pos, _, _ := PortPosition(n, n.Ports[j].ID)
			d := dist2D(worldPt, pos)
			if d <= bestDist {
				bestDist = d
				bestN = n
				bestP = &n.Ports[j]
				bestPos = pos
			}
		}
	}
	return bestN, bestP, bestPos
}

// HitTestNode returns the topmost node covering worldPt (preferring non-group
// nodes over enclosing LabeledGroupNode backdrops).
func (g *Graph) HitTestNode(worldPt Point) *Node {
	var groupHit *Node
	for i := len(g.Nodes) - 1; i >= 0; i-- {
		n := g.Nodes[i]
		if n.Bounds().Contains(worldPt) {
			if n.Variant == VariantGroup {
				if groupHit == nil {
					groupHit = n
				}
				continue
			}
			return n
		}
	}
	return groupHit
}

// HitTestEdge returns the edge closest to screenPt within tolerance DIPs.
func (g *Graph) HitTestEdge(screenPt Point, tolerance float32) *Edge {
	bestDist := tolerance
	var best *Edge
	for i := len(g.Edges) - 1; i >= 0; i-- {
		e := g.Edges[i]
		sn := g.NodeByID(e.Source)
		dn := g.NodeByID(e.Target)
		if sn == nil || dn == nil {
			continue
		}
		spW, spSide, _ := PortPosition(sn, e.SourcePort)
		dpW, dpSide, _ := PortPosition(dn, e.TargetPort)
		spS := g.Viewport.WorldToScreen(spW)
		dpS := g.Viewport.WorldToScreen(dpW)
		kind := e.Kind
		if kind == "" {
			kind = g.EdgeKind
		}
		geom := BuildEdgeGeometry(spS, dpS, spSide, dpSide, kind)
		d := DistanceToPolyline(screenPt, geom.Samples)
		if d <= bestDist {
			bestDist = d
			best = e
		}
	}
	return best
}

func min32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
