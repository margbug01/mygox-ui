package flow

import "fmt"

// NodeTemplate describes a reusable palette item users can spawn onto a Graph.
type NodeTemplate struct {
	Type       string
	Variant    NodeVariant
	Title      string
	Subtitle   string
	Badge      string
	Footer     string
	AccentHex  uint32
	W, H       float32
	Ports      []Port
	Fields     []NodeField
	SchemaRows []SchemaRow
}

// DefaultPalette returns React Flow UI-styled node templates for workflow & schema builders.
func DefaultPalette() []NodeTemplate {
	return []NodeTemplate{
		{
			Type:      "base",
			Variant:   VariantBase,
			Title:     "LLM Agent",
			Subtitle:  "Reasoning & tool calls",
			Badge:     "AI",
			Footer:    "model: opus-5.5",
			AccentHex: 0x6366f1,
			W:         232,
			H:         124,
			Ports: []Port{
				{ID: "prompt", Label: "prompt", Side: PortLeft, Kind: PortInput},
				{ID: "context", Label: "ctx", Side: PortLeft, Kind: PortInput},
				{ID: "out", Label: "reply", Side: PortRight, Kind: PortOutput},
			},
			Fields: []NodeField{
				{Key: "Temp", Value: "0.2"},
			},
		},
		{
			Type:      "schema",
			Variant:   VariantSchema,
			Title:     "Orders",
			Subtitle:  "Database table",
			Badge:     "SQL",
			AccentHex: 0x06b6d4,
			W:         224,
			H:         146,
			SchemaRows: []SchemaRow{
				{Name: "id", Type: "uuid", InID: "in-id", OutID: "out-id", IsKey: true},
				{Name: "product_id", Type: "uuid", OutID: "out-pid"},
				{Name: "total", Type: "money"},
				{Name: "created_at", Type: "timestamptz"},
			},
			Ports: []Port{
				{ID: "in-id", Side: PortLeft, Kind: PortInput},
				{ID: "out-id", Side: PortRight, Kind: PortOutput},
				{ID: "out-pid", Side: PortRight, Kind: PortOutput},
			},
		},
		{
			Type:      "group",
			Variant:   VariantGroup,
			Title:     "Sub-flow Group",
			Subtitle:  "Drag moves enclosed nodes",
			Badge:     "GRP",
			AccentHex: 0x64748b,
			W:         320,
			H:         200,
		},
		{
			Type:      "placeholder",
			Variant:   VariantPlaceholder,
			Title:     "+ Add step",
			Subtitle:  "Click to convert to node",
			Badge:     "NEW",
			AccentHex: 0x94a3b8,
			W:         168,
			H:         72,
			Ports: []Port{
				{ID: "in", Side: PortLeft, Kind: PortInput},
				{ID: "out", Side: PortRight, Kind: PortOutput},
			},
		},
		{
			Type:      "annotation",
			Variant:   VariantAnnotation,
			Title:     "Canvas note",
			Subtitle:  "Drag to document this stage",
			Badge:     "NOTE",
			AccentHex: 0xf59e0b,
			W:         210,
			H:         64,
		},
		{
			Type:      "trigger",
			Variant:   VariantBase,
			Title:     "Webhook Trigger",
			Subtitle:  "POST /api/v1/run",
			Badge:     "IN",
			AccentHex: 0x22c55e,
			W:         216,
			H:         106,
			Ports: []Port{
				{ID: "out", Label: "event", Side: PortRight, Kind: PortOutput},
			},
			Fields: []NodeField{
				{Key: "Auth", Value: "Bearer JWT"},
			},
		},
		{
			Type:      "output",
			Variant:   VariantBase,
			Title:     "Response Sink",
			Subtitle:  "Stream SSE / Slack",
			Badge:     "OUT",
			AccentHex: 0x8b5cf6,
			W:         216,
			H:         106,
			Ports: []Port{
				{ID: "in", Label: "msg", Side: PortLeft, Kind: PortInput},
			},
			Fields: []NodeField{
				{Key: "Format", Value: "Markdown"},
			},
		},
	}
}

// SpawnFromTemplate instantiates a template near the center of the current viewport.
func (g *Graph) SpawnFromTemplate(tpl NodeTemplate) *Node {
	cx := g.lastCanvasW * 0.5
	cy := g.lastCanvasH * 0.5
	if cx == 0 {
		cx = 400
	}
	if cy == 0 {
		cy = 260
	}
	world := g.Viewport.ScreenToWorld(Point{X: cx, Y: cy})
	offset := float32((len(g.Nodes) % 5) * 24)
	n := &Node{
		Type:      tpl.Type,
		Variant:   tpl.Variant,
		Title:     tpl.Title,
		Subtitle:  tpl.Subtitle,
		Badge:     tpl.Badge,
		Footer:    tpl.Footer,
		AccentHex: tpl.AccentHex,
		W:         tpl.W,
		H:         tpl.H,
		X:         world.X - tpl.W*0.5 + offset,
		Y:         world.Y - tpl.H*0.5 + offset,
		Status:    StatusIdle,
	}
	n.Ports = append([]Port(nil), tpl.Ports...)
	n.Fields = append([]NodeField(nil), tpl.Fields...)
	n.SchemaRows = append([]SchemaRow(nil), tpl.SchemaRows...)
	added := g.AddNode(n)
	g.SelectNode(added.ID, false)
	return added
}

// PresetDatabaseSchema replicates React Flow UI's canonical DatabaseSchemaNode
// showcase (Products linked by foreign-key rows to Warehouses and Suppliers).
func PresetDatabaseSchema() *Graph {
	g := NewGraph()
	g.Background = BackgroundDots
	g.EdgeKind = EdgeSmoothStep
	g.Nodes = []*Node{
		{
			ID: "note-schema", Type: "annotation", Variant: VariantAnnotation,
			Badge: "01. React Flow UI — DatabaseSchemaNode",
			Title: "Row-level foreign key handles",
			Subtitle: "Hover any edge to delete, or drag a handle to connect",
			X: 40, Y: -54, W: 320, H: 60,
		},
		{
			ID: "products", Type: "schema", Variant: VariantSchema,
			Title: "Products", X: 40, Y: 70, W: 248, H: 226,
			SchemaRows: []SchemaRow{
				{Name: "id", Type: "uuid", OutID: "id", IsKey: true},
				{Name: "name", Type: "varchar"},
				{Name: "description", Type: "varchar"},
				{Name: "warehouse_id", Type: "uuid", OutID: "warehouse_id"},
				{Name: "supplier_id", Type: "uuid", OutID: "supplier_id"},
				{Name: "price", Type: "money"},
				{Name: "quantity", Type: "int4"},
			},
			Ports: []Port{
				{ID: "id", Side: PortRight, Kind: PortOutput},
				{ID: "warehouse_id", Side: PortRight, Kind: PortOutput},
				{ID: "supplier_id", Side: PortRight, Kind: PortOutput},
			},
		},
		{
			ID: "warehouses", Type: "schema", Variant: VariantSchema,
			Title: "Warehouses", X: 420, Y: 20, W: 236, H: 148,
			SchemaRows: []SchemaRow{
				{Name: "id", Type: "uuid", InID: "id", IsKey: true},
				{Name: "name", Type: "varchar"},
				{Name: "address", Type: "varchar"},
				{Name: "capacity", Type: "int4"},
			},
			Ports: []Port{
				{ID: "id", Side: PortLeft, Kind: PortInput},
			},
		},
		{
			ID: "suppliers", Type: "schema", Variant: VariantSchema,
			Title: "Suppliers", X: 420, Y: 210, W: 236, H: 148,
			SchemaRows: []SchemaRow{
				{Name: "id", Type: "uuid", InID: "id", IsKey: true},
				{Name: "name", Type: "varchar"},
				{Name: "description", Type: "varchar"},
				{Name: "country", Type: "varchar"},
			},
			Ports: []Port{
				{ID: "id", Side: PortLeft, Kind: PortInput},
			},
		},
		{
			ID: "inventory", Type: "schema", Variant: VariantSchema,
			Title: "StockMovements", X: 760, Y: 100, W: 242, H: 148, Status: StatusSuccess,
			SchemaRows: []SchemaRow{
				{Name: "id", Type: "uuid", InID: "id", IsKey: true},
				{Name: "warehouse_id", Type: "uuid", InID: "wh"},
				{Name: "delta", Type: "int4"},
				{Name: "recorded_at", Type: "timestamp"},
			},
			Ports: []Port{
				{ID: "id", Side: PortLeft, Kind: PortInput},
				{ID: "wh", Side: PortLeft, Kind: PortInput},
			},
		},
	}
	// Add output handle on Warehouses.id to feed StockMovements.
	g.Nodes[2].SchemaRows[0].OutID = "out-id"
	g.Nodes[2].Ports = append(g.Nodes[2].Ports, Port{ID: "out-id", Side: PortRight, Kind: PortOutput})

	g.Edges = []*Edge{
		{ID: "fk-wh", Source: "products", SourcePort: "warehouse_id", Target: "warehouses", TargetPort: "id", Label: "1:N"},
		{ID: "fk-sup", Source: "products", SourcePort: "supplier_id", Target: "suppliers", TargetPort: "id", Label: "1:N"},
		{ID: "fk-mov", Source: "warehouses", SourcePort: "out-id", Target: "inventory", TargetPort: "wh", Animated: true},
	}
	return g
}

// PresetSubFlows showcases LabeledGroupNode, NodeStatusIndicator (spinning conic border),
// PlaceholderNode (click to promote), AnnotationNode, and ButtonEdge together.
func PresetSubFlows() *Graph {
	g := NewGraph()
	g.Background = BackgroundDots
	g.EdgeKind = EdgeBezier
	g.Nodes = []*Node{
		{
			ID: "note-1", Type: "annotation", Variant: VariantAnnotation,
			Badge: "01. React Flow UI — Sub-flows & Status",
			Title: "Drag the group container to move enclosed nodes",
			Subtitle: "Click the dashed Placeholder node to convert it into a step",
			X: 40, Y: -24, W: 360, H: 60,
		},
		{
			ID: "grp-retrieval", Type: "group", Variant: VariantGroup,
			Title: "Parallel Retrieval Sub-flow",
			X: 330, Y: 54, W: 310, H: 310,
		},
		{
			ID: "ingress", Type: "trigger", Variant: VariantBase,
			Title: "API Request", Subtitle: "POST /v1/agent/run",
			Badge: "IN", AccentHex: 0x22c55e, X: 40, Y: 156, W: 214, H: 104, Status: StatusSuccess,
			Ports: []Port{{ID: "out", Label: "req", Side: PortRight, Kind: PortOutput}},
		},
		{
			ID: "vec", Type: "tool", Variant: VariantBase,
			Title: "Vector Recall", Subtitle: "HNSW cosine top-8",
			Badge: "DB", AccentHex: 0x06b6d4, X: 368, Y: 102, W: 234, H: 104, Status: StatusRunning,
			Ports: []Port{
				{ID: "in", Side: PortLeft, Kind: PortInput},
				{ID: "out", Side: PortRight, Kind: PortOutput},
			},
			Fields: []NodeField{{Key: "Latency", Value: "12 ms"}},
		},
		{
			ID: "web", Type: "tool", Variant: VariantBase,
			Title: "Live Search", Subtitle: "Exa neural fetch",
			Badge: "WEB", AccentHex: 0xec4899, X: 368, Y: 236, W: 234, H: 104, Status: StatusRunning,
			Ports: []Port{
				{ID: "in", Side: PortLeft, Kind: PortInput},
				{ID: "out", Side: PortRight, Kind: PortOutput},
			},
			Fields: []NodeField{{Key: "Sources", Value: "5 docs"}},
		},
		{
			ID: "synth", Type: "llm", Variant: VariantBase,
			Title: "Synthesizer", Subtitle: "Grounded streaming reply",
			Badge: "AI", Footer: "Status: idle · Click Run", AccentHex: 0x6366f1,
			X: 720, Y: 146, W: 232, H: 126, Status: StatusIdle,
			Ports: []Port{
				{ID: "a", Label: "vec", Side: PortLeft, Kind: PortInput},
				{ID: "b", Label: "web", Side: PortLeft, Kind: PortInput},
				{ID: "out", Label: "next", Side: PortRight, Kind: PortOutput},
			},
		},
		{
			ID: "placeholder-next", Type: "placeholder", Variant: VariantPlaceholder,
			Title: "+ Add output step",
			X: 1030, Y: 172, W: 176, H: 74,
			Ports: []Port{
				{ID: "in", Side: PortLeft, Kind: PortInput},
				{ID: "out", Side: PortRight, Kind: PortOutput},
			},
		},
	}
	g.Edges = []*Edge{
		{ID: "s1", Source: "ingress", SourcePort: "out", Target: "vec", TargetPort: "in", Animated: true},
		{ID: "s2", Source: "ingress", SourcePort: "out", Target: "web", TargetPort: "in", Animated: true},
		{ID: "s3", Source: "vec", SourcePort: "out", Target: "synth", TargetPort: "a"},
		{ID: "s4", Source: "web", SourcePort: "out", Target: "synth", TargetPort: "b"},
		{ID: "s5", Source: "synth", SourcePort: "out", Target: "placeholder-next", TargetPort: "in", Animated: true},
	}
	return g
}

// PresetAIAgent builds a multi-branch AI Agent orchestration graph.
func PresetAIAgent() *Graph {
	g := NewGraph()
	g.Nodes = []*Node{
		{
			ID: "trigger", Type: "trigger", Variant: VariantBase, Title: "User Prompt", Subtitle: "Chat / API ingress",
			Badge: "IN", AccentHex: 0x22c55e, X: 40, Y: 180, W: 210, H: 108, Status: StatusSuccess,
			Ports:  []Port{{ID: "out", Label: "msg", Side: PortRight, Kind: PortOutput}},
			Fields: []NodeField{{Key: "Channel", Value: "Desktop + Web"}},
		},
		{
			ID: "router", Type: "router", Variant: VariantBase, Title: "Intent Router", Subtitle: "Classify & dispatch",
			Badge: "IF", AccentHex: 0xf59e0b, X: 330, Y: 174, W: 224, H: 122, Status: StatusSuccess,
			Ports: []Port{
				{ID: "in", Label: "in", Side: PortLeft, Kind: PortInput},
				{ID: "rag", Label: "recall", Side: PortRight, Kind: PortOutput},
				{ID: "web", Label: "search", Side: PortRight, Kind: PortOutput},
			},
			Fields: []NodeField{{Key: "Strategy", Value: "Hybrid Parallel"}},
		},
		{
			ID: "vector", Type: "tool", Variant: VariantBase, Title: "Vector Memory", Subtitle: "pgvector HNSW top-8",
			Badge: "DB", AccentHex: 0x06b6d4, X: 640, Y: 72, W: 228, H: 112, Status: StatusRunning,
			Ports: []Port{
				{ID: "q", Label: "query", Side: PortLeft, Kind: PortInput},
				{ID: "docs", Label: "chunks", Side: PortRight, Kind: PortOutput},
			},
			Fields: []NodeField{{Key: "Latency", Value: "14 ms"}},
		},
		{
			ID: "search", Type: "tool", Variant: VariantBase, Title: "Exa Web Search", Subtitle: "Live citation fetch",
			Badge: "WEB", AccentHex: 0xec4899, X: 640, Y: 266, W: 228, H: 112, Status: StatusRunning,
			Ports: []Port{
				{ID: "q", Label: "query", Side: PortLeft, Kind: PortInput},
				{ID: "res", Label: "pages", Side: PortRight, Kind: PortOutput},
			},
			Fields: []NodeField{{Key: "Results", Value: "Top 6 URLs"}},
		},
		{
			ID: "llm", Type: "llm", Variant: VariantBase, Title: "Synthesizer LLM", Subtitle: "Grounded final answer",
			Badge: "AI", AccentHex: 0x6366f1, X: 960, Y: 162, W: 242, H: 132, Status: StatusIdle,
			Ports: []Port{
				{ID: "ctx", Label: "memory", Side: PortLeft, Kind: PortInput},
				{ID: "web", Label: "web", Side: PortLeft, Kind: PortInput},
				{ID: "out", Label: "stream", Side: PortRight, Kind: PortOutput},
			},
			Fields: []NodeField{
				{Key: "Model", Value: "opus-5.5"},
				{Key: "Tokens", Value: "4k max"},
			},
		},
		{
			ID: "sink", Type: "output", Variant: VariantBase, Title: "UI Stream Sink", Subtitle: "Render tokens @ 60fps",
			Badge: "OUT", AccentHex: 0x8b5cf6, X: 1290, Y: 174, W: 214, H: 108, Status: StatusIdle,
			Ports:  []Port{{ID: "in", Label: "stream", Side: PortLeft, Kind: PortInput}},
			Fields: []NodeField{{Key: "Target", Value: "MyGo Native"}},
		},
	}
	g.Edges = []*Edge{
		{ID: "e1", Source: "trigger", SourcePort: "out", Target: "router", TargetPort: "in", Label: "JSON", Animated: true},
		{ID: "e2", Source: "router", SourcePort: "rag", Target: "vector", TargetPort: "q", Label: "embed", Animated: true},
		{ID: "e3", Source: "router", SourcePort: "web", Target: "search", TargetPort: "q", Label: "live", Animated: true},
		{ID: "e4", Source: "vector", SourcePort: "docs", Target: "llm", TargetPort: "ctx"},
		{ID: "e5", Source: "search", SourcePort: "res", Target: "llm", TargetPort: "web"},
		{ID: "e6", Source: "llm", SourcePort: "out", Target: "sink", TargetPort: "in", Label: "SSE", Animated: true},
	}
	return g
}

// PresetDataPipeline builds a real-time streaming ETL workflow.
func PresetDataPipeline() *Graph {
	g := NewGraph()
	g.EdgeKind = EdgeSmoothStep
	g.Background = BackgroundLines
	g.Nodes = []*Node{
		{
			ID: "kafka", Type: "trigger", Variant: VariantBase, Title: "Kafka Stream", Subtitle: "orders.events.v3",
			Badge: "SRC", AccentHex: 0x22c55e, X: 60, Y: 80, W: 218, H: 108, Status: StatusSuccess,
			Ports:  []Port{{ID: "out", Label: "batch", Side: PortRight, Kind: PortOutput}},
			Fields: []NodeField{{Key: "Rate", Value: "18.4k msg/s"}},
		},
		{
			ID: "pg", Type: "tool", Variant: VariantBase, Title: "Postgres CDC", Subtitle: "users & pricing dim",
			Badge: "SQL", AccentHex: 0x06b6d4, X: 60, Y: 250, W: 218, H: 108, Status: StatusSuccess,
			Ports:  []Port{{ID: "out", Label: "rows", Side: PortRight, Kind: PortOutput}},
			Fields: []NodeField{{Key: "Lag", Value: "42 ms"}},
		},
		{
			ID: "join", Type: "code", Variant: VariantBase, Title: "Enrich & Filter", Subtitle: "Deduplicate + currency",
			Badge: "ETL", AccentHex: 0xf59e0b, X: 380, Y: 162, W: 234, H: 122, Status: StatusRunning,
			Ports: []Port{
				{ID: "a", Label: "events", Side: PortLeft, Kind: PortInput},
				{ID: "b", Label: "dim", Side: PortLeft, Kind: PortInput},
				{ID: "out", Label: "clean", Side: PortRight, Kind: PortOutput},
			},
			Fields: []NodeField{{Key: "Window", Value: "5s tumbling"}},
		},
		{
			ID: "duck", Type: "tool", Variant: VariantBase, Title: "DuckDB Rollup", Subtitle: "Parquet columnar agg",
			Badge: "OLAP", AccentHex: 0x6366f1, X: 710, Y: 162, W: 226, H: 118, Status: StatusIdle,
			Ports: []Port{
				{ID: "in", Label: "in", Side: PortLeft, Kind: PortInput},
				{ID: "bi", Label: "kpi", Side: PortRight, Kind: PortOutput},
				{ID: "alert", Label: "anom", Side: PortRight, Kind: PortOutput},
			},
			Fields: []NodeField{{Key: "Memory", Value: "512 MB"}},
		},
		{
			ID: "dash", Type: "output", Variant: VariantBase, Title: "Metrics Board", Subtitle: "Live GMV & P99",
			Badge: "BI", AccentHex: 0x8b5cf6, X: 1030, Y: 80, W: 212, H: 104, Status: StatusIdle,
			Ports: []Port{{ID: "in", Label: "kpi", Side: PortLeft, Kind: PortInput}},
		},
		{
			ID: "pager", Type: "output", Variant: VariantBase, Title: "Anomaly Alert", Subtitle: "PagerDuty / Telegram",
			Badge: "OPS", AccentHex: 0xef4444, X: 1030, Y: 250, W: 212, H: 104, Status: StatusError,
			Ports: []Port{{ID: "in", Label: "alert", Side: PortLeft, Kind: PortInput}},
		},
	}
	g.Edges = []*Edge{
		{ID: "d1", Source: "kafka", SourcePort: "out", Target: "join", TargetPort: "a", Animated: true},
		{ID: "d2", Source: "pg", SourcePort: "out", Target: "join", TargetPort: "b"},
		{ID: "d3", Source: "join", SourcePort: "out", Target: "duck", TargetPort: "in", Label: "Arrow", Animated: true},
		{ID: "d4", Source: "duck", SourcePort: "bi", Target: "dash", TargetPort: "in"},
		{ID: "d5", Source: "duck", SourcePort: "alert", Target: "pager", TargetPort: "in"},
	}
	return g
}

// PresetCICD builds a multi-stage build and progressive rollout pipeline.
func PresetCICD() *Graph {
	return PresetSubFlows()
}

// TopologicalOrder returns node IDs ordered from sources to sinks (best-effort on cycles).
func (g *Graph) TopologicalOrder() []string {
	inDeg := make(map[string]int, len(g.Nodes))
	adj := make(map[string][]string, len(g.Nodes))
	for _, n := range g.Nodes {
		if n.Variant == VariantGroup || n.Variant == VariantAnnotation {
			continue
		}
		inDeg[n.ID] = 0
	}
	for _, e := range g.Edges {
		if _, ok1 := inDeg[e.Source]; !ok1 {
			continue
		}
		if _, ok2 := inDeg[e.Target]; !ok2 {
			continue
		}
		if e.Source == e.Target {
			continue
		}
		adj[e.Source] = append(adj[e.Source], e.Target)
		inDeg[e.Target]++
	}
	q := make([]string, 0, len(inDeg))
	for _, n := range g.Nodes {
		if deg, ok := inDeg[n.ID]; ok && deg == 0 {
			q = append(q, n.ID)
		}
	}
	seen := map[string]bool{}
	out := make([]string, 0, len(inDeg))
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		seen[cur] = true
		out = append(out, cur)
		for _, nxt := range adj[cur] {
			inDeg[nxt]--
			if inDeg[nxt] == 0 {
				q = append(q, nxt)
			}
		}
	}
	for _, n := range g.Nodes {
		if _, ok := inDeg[n.ID]; ok && !seen[n.ID] {
			out = append(out, n.ID)
		}
	}
	return out
}

// SummaryString formats a concise status line for toolbars and tests.
func (g *Graph) SummaryString() string {
	return fmt.Sprintf("%d nodes · %d edges · %d%%", len(g.Nodes), len(g.Edges), int(g.Viewport.safeZoom()*100))
}
