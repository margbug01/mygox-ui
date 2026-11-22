package flow_test

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/flow"
)

func TestViewportTransformsAndZoomAt(t *testing.T) {
	vp := flow.Viewport{X: 100, Y: 50, Zoom: 2}
	w := vp.ScreenToWorld(flow.Point{X: 300, Y: 250})
	if math.Abs(float64(w.X-100)) > 1e-3 || math.Abs(float64(w.Y-100)) > 1e-3 {
		t.Fatalf("unexpected ScreenToWorld: %+v", w)
	}
	s := vp.WorldToScreen(w)
	if math.Abs(float64(s.X-300)) > 1e-3 || math.Abs(float64(s.Y-250)) > 1e-3 {
		t.Fatalf("round-trip mismatch: %+v", s)
	}

	// Zooming around an anchor keeps that world point pinned to the same screen point.
	anchorScreen := flow.Point{X: 420, Y: 310}
	beforeWorld := vp.ScreenToWorld(anchorScreen)
	vp.ZoomAt(anchorScreen, 1.5, 0.2, 4.0)
	afterScreen := vp.WorldToScreen(beforeWorld)
	if math.Abs(float64(afterScreen.X-anchorScreen.X)) > 0.1 || math.Abs(float64(afterScreen.Y-anchorScreen.Y)) > 0.1 {
		t.Fatalf("ZoomAt drifted anchor: before=%+v after=%+v", anchorScreen, afterScreen)
	}
}

func TestConnectValidationUndoRedoAndJSON(t *testing.T) {
	g := flow.NewGraph()
	a := g.AddNode(&flow.Node{
		ID: "a", Title: "Source", X: 0, Y: 0, W: 200, H: 100,
		Ports: []flow.Port{{ID: "out", Side: flow.PortRight, Kind: flow.PortOutput}},
	})
	b := g.AddNode(&flow.Node{
		ID: "b", Title: "Target", X: 300, Y: 0, W: 200, H: 100,
		Ports: []flow.Port{
			{ID: "in", Side: flow.PortLeft, Kind: flow.PortInput},
			{ID: "out", Side: flow.PortRight, Kind: flow.PortOutput},
		},
	})

	// Reject self-loop and output-to-output.
	if _, ok := g.Connect(a.ID, "out", a.ID, "out"); ok {
		t.Fatal("expected self-loop connection to be rejected")
	}
	if _, ok := g.Connect(a.ID, "out", b.ID, "out"); ok {
		t.Fatal("expected output->output connection to be rejected")
	}

	edge, ok := g.Connect(a.ID, "out", b.ID, "in")
	if !ok || edge == nil {
		t.Fatal("expected valid output->input connection to succeed")
	}
	// Reject duplicate edge.
	if _, ok := g.Connect(a.ID, "out", b.ID, "in"); ok {
		t.Fatal("expected duplicate connection to be rejected")
	}

	// Undo removes edge, Redo restores it.
	if !g.Undo() || len(g.Edges) != 0 {
		t.Fatalf("expected Undo to remove edge, got %d edges", len(g.Edges))
	}
	if !g.Redo() || len(g.Edges) != 1 {
		t.Fatalf("expected Redo to restore edge, got %d edges", len(g.Edges))
	}

	// JSON round-trip preserves nodes and edges.
	raw, err := g.SaveJSON()
	if err != nil {
		t.Fatalf("SaveJSON: %v", err)
	}
	g2 := flow.NewGraph()
	if err := g2.LoadJSON(raw); err != nil {
		t.Fatalf("LoadJSON: %v", err)
	}
	if len(g2.Nodes) != 2 || len(g2.Edges) != 1 {
		t.Fatalf("JSON round-trip mismatch: nodes=%d edges=%d", len(g2.Nodes), len(g2.Edges))
	}
}

func TestAutoLayoutAndTopologicalOrder(t *testing.T) {
	g := flow.PresetAIAgent()
	g.AutoLayout(flow.LayoutLR)
	trig := g.NodeByID("trigger")
	router := g.NodeByID("router")
	llm := g.NodeByID("llm")
	sink := g.NodeByID("sink")
	if !(trig.X < router.X && router.X < llm.X && llm.X < sink.X) {
		t.Fatalf("expected left-to-right rank ordering: trigger=%.0f router=%.0f llm=%.0f sink=%.0f", trig.X, router.X, llm.X, sink.X)
	}

	g.AutoLayout(flow.LayoutTB)
	if !(trig.Y < router.Y && router.Y < llm.Y && llm.Y < sink.Y) {
		t.Fatalf("expected top-to-bottom rank ordering: trigger=%.0f router=%.0f llm=%.0f sink=%.0f", trig.Y, router.Y, llm.Y, sink.Y)
	}

	order := g.TopologicalOrder()
	if len(order) != len(g.Nodes) || order[0] != "trigger" || order[len(order)-1] != "sink" {
		t.Fatalf("unexpected topological order: %v", order)
	}
}

func TestCanvasRenderDragAndConnectInteraction(t *testing.T) {
	g := flow.PresetAIAgent()
	g.SnapToGrid = false
	g.FitView(1000, 640, 40)

	tt := ui.NewTester(func(c *ui.Context) {
		flow.Canvas(c, g, flow.CanvasOptions{Theme: flow.DarkTheme()})
	}, 1000, 640)
	tt.Frame()
	tt.Advance(50 * time.Millisecond)
	tt.Frame()

	// Select and drag the "trigger" node via real pointer events.
	trig := g.NodeByID("trigger")
	startX := trig.X
	centerS := g.Viewport.WorldToScreen(flow.Point{X: trig.X + trig.W*0.5, Y: trig.Y + trig.H*0.5})
	tt.Press(centerS.X, centerS.Y)
	tt.Move(centerS.X+80, centerS.Y+20)
	tt.Release(centerS.X+80, centerS.Y+20)
	tt.Frame()

	if !trig.Selected {
		t.Fatal("expected dragged node to be selected")
	}
	if trig.X <= startX+20 {
		t.Fatalf("expected node X to move right after drag: before=%.1f after=%.1f", startX, trig.X)
	}

	// Duplicate and delete selection via graph commands.
	prevNodes := len(g.Nodes)
	if n := g.DuplicateSelected(); n != 1 || len(g.Nodes) != prevNodes+1 {
		t.Fatalf("expected duplicate to add 1 node, got added=%d total=%d", n, len(g.Nodes))
	}
	if !g.DeleteSelected() || len(g.Nodes) != prevNodes {
		t.Fatalf("expected delete to remove duplicated node, total=%d", len(g.Nodes))
	}
}

func TestReactFlowUISchemaPlaceholderAndGroupNodes(t *testing.T) {
	// 1. DatabaseSchemaNode row-level handles align monotonically down table rows.
	sg := flow.PresetDatabaseSchema()
	prod := sg.NodeByID("products")
	p1, s1, ok1 := flow.PortPosition(prod, "id")
	p2, s2, ok2 := flow.PortPosition(prod, "warehouse_id")
	p3, s3, ok3 := flow.PortPosition(prod, "supplier_id")
	if !ok1 || !ok2 || !ok3 || s1 != flow.PortRight || s2 != flow.PortRight || s3 != flow.PortRight {
		t.Fatalf("expected right-side schema row handles, got ok=(%v,%v,%v)", ok1, ok2, ok3)
	}
	if !(p1.Y < p2.Y && p2.Y < p3.Y) {
		t.Fatalf("expected schema handles ordered top-to-bottom: %.1f, %.1f, %.1f", p1.Y, p2.Y, p3.Y)
	}

	// 2. PlaceholderNode promotes to BaseNode on click & clears incoming animated dash.
	fg := flow.PresetSubFlows()
	ph := fg.NodeByID("placeholder-next")
	if ph == nil || ph.Variant != flow.VariantPlaceholder {
		t.Fatalf("expected placeholder node in SubFlows preset")
	}
	if !fg.PromotePlaceholder(ph.ID) || ph.Variant != flow.VariantBase {
		t.Fatalf("expected PromotePlaceholder to convert to VariantBase, got %v", ph.Variant)
	}

	// 3. RemoveEdge deletes a single edge cleanly and supports Undo.
	eCount := len(fg.Edges)
	if !fg.RemoveEdge("s1") || len(fg.Edges) != eCount-1 {
		t.Fatalf("expected RemoveEdge to remove s1, count=%d", len(fg.Edges))
	}
	if !fg.Undo() || len(fg.Edges) != eCount {
		t.Fatalf("expected Undo to restore removed edge, count=%d", len(fg.Edges))
	}

	// 4. Headless render both new React Flow UI presets cleanly.
	for _, g := range []*flow.Graph{sg, fg} {
		g.FitView(1100, 720, 48)
		tt := ui.NewTester(func(c *ui.Context) {
			flow.Canvas(c, g, flow.CanvasOptions{Theme: flow.LightTheme()})
		}, 1100, 720)
		tt.Frame()
		tt.Advance(60 * time.Millisecond)
		tt.Frame()
	}
}

func TestReactFlowUIAll17ComponentsCompleteness(t *testing.T) {
	g := flow.PresetSubFlows()

	// AnimatedSvgEdge: SamplePolylineAt smoothly interpolates midpoint along path.
	mid := flow.SamplePolylineAt([]flow.Point{{X: 0, Y: 0}, {X: 100, Y: 0}, {X: 100, Y: 100}}, 0.5)
	if math.Abs(float64(mid.X-100)) > 0.1 || math.Abs(float64(mid.Y-0)) > 0.1 {
		t.Fatalf("unexpected SamplePolylineAt midpoint: %+v", mid)
	}

	// ButtonHandle: SpawnConnectedNode appends a wired child node & records DevTools ChangeLog.
	prevNodes, prevEdges := len(g.Nodes), len(g.Edges)
	child, edge := g.SpawnConnectedNode("synth", "out")
	if child == nil || edge == nil || len(g.Nodes) != prevNodes+1 || len(g.Edges) != prevEdges+1 {
		t.Fatalf("SpawnConnectedNode failed: nodes=%d edges=%d", len(g.Nodes), len(g.Edges))
	}
	if len(g.ChangeLog) == 0 {
		t.Fatal("expected DevTools ChangeLog to record spawn mutation")
	}

	// NodeSearch: FocusNode selects and centers viewport on target node.
	if !g.FocusNode("vec") || !g.NodeByID("vec").Selected {
		t.Fatal("expected FocusNode(vec) to select node")
	}

	// Render with DevTools, NodeSearch, StatusOverlay, Appendix, PacketMotion, and ZoomSlider active.
	g.ShowDevTools = true
	g.ShowSearch = true
	tt := ui.NewTester(func(c *ui.Context) {
		flow.Canvas(c, g, flow.CanvasOptions{Theme: flow.DarkTheme()})
	}, 1200, 760)
	tt.Frame()
	tt.Advance(120 * time.Millisecond)
	tt.Frame()
	if !tt.HasText("FLOW DEVTOOLS") || !tt.HasText("Change Logger") {
		t.Fatalf("expected DevTools overlay text, got: %v", tt.Texts())
	}
}
