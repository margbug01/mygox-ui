package flow

import "fmt"

// NodeTemplate describes a reusable palette item users can spawn onto a Graph.
type NodeTemplate struct {
	Type      string
	Title     string
	Subtitle  string
	Badge     string
	AccentHex uint32
	W, H      float32
	Ports     []Port
	Fields    []NodeField
}

// DefaultPalette returns ready-to-use node templates for AI & data workflow builders.
func DefaultPalette() []NodeTemplate {
	return []NodeTemplate{
		{
			Type:      "trigger",
			Title:     "Webhook Trigger",
			Subtitle:  "POST /api/v1/run",
			Badge:     "IN",
			AccentHex: 0x22c55e,
			W:         216,
			H:         106,
			Ports: []Port{
				{ID: "out", Label: "event", Side: PortRight, Kind: PortOutput, ColorHex: 0x22c55e},
			},
			Fields: []NodeField{
				{Key: "Auth", Value: "Bearer JWT"},
			},
		},
		{
			Type:      "llm",
			Title:     "LLM Agent",
			Subtitle:  "Reasoning & tool calls",
			Badge:     "AI",
			AccentHex: 0x6366f1,
			W:         236,
			H:         122,
			Ports: []Port{
				{ID: "prompt", Label: "prompt", Side: PortLeft, Kind: PortInput, ColorHex: 0x6366f1},
				{ID: "context", Label: "ctx", Side: PortLeft, Kind: PortInput, ColorHex: 0x38bdf8},
				{ID: "out", Label: "reply", Side: PortRight, Kind: PortOutput, ColorHex: 0x6366f1},
			},
			Fields: []NodeField{
				{Key: "Model", Value: "claude-opus-5.5"},
				{Key: "Temp", Value: "0.2"},
			},
		},
		{
			Type:      "router",
			Title:     "Intent Router",
			Subtitle:  "Conditional branch",
			Badge:     "IF",
			AccentHex: 0xf59e0b,
			W:         224,
			H:         122,
			Ports: []Port{
				{ID: "in", Label: "in", Side: PortLeft, Kind: PortInput, ColorHex: 0xf59e0b},
				{ID: "true", Label: "fast", Side: PortRight, Kind: PortOutput, ColorHex: 0x22c55e},
				{ID: "false", Label: "deep", Side: PortRight, Kind: PortOutput, ColorHex: 0xf59e0b},
			},
			Fields: []NodeField{
				{Key: "Rule", Value: "complexity > 0.7"},
			},
		},
		{
			Type:      "tool",
			Title:     "Vector Search",
			Subtitle:  "Top-K semantic recall",
			Badge:     "DB",
			AccentHex: 0x06b6d4,
			W:         224,
			H:         114,
			Ports: []Port{
				{ID: "query", Label: "query", Side: PortLeft, Kind: PortInput, ColorHex: 0x06b6d4},
				{ID: "docs", Label: "docs", Side: PortRight, Kind: PortOutput, ColorHex: 0x06b6d4},
			},
			Fields: []NodeField{
				{Key: "Index", Value: "knowledge_v2"},
			},
		},
		{
			Type:      "code",
			Title:     "Go Transform",
			Subtitle:  "Pure Go sandbox",
			Badge:     "FN",
			AccentHex: 0xec4899,
			W:         220,
			H:         114,
			Ports: []Port{
				{ID: "in", Label: "data", Side: PortLeft, Kind: PortInput, ColorHex: 0xec4899},
				{ID: "out", Label: "json", Side: PortRight, Kind: PortOutput, ColorHex: 0xec4899},
			},
			Fields: []NodeField{
				{Key: "Runtime", Value: "go 1.24"},
			},
		},
		{
			Type:      "output",
			Title:     "Response Sink",
			Subtitle:  "Stream SSE / Slack",
			Badge:     "OUT",
			AccentHex: 0x8b5cf6,
			W:         216,
			H:         106,
			Ports: []Port{
				{ID: "in", Label: "msg", Side: PortLeft, Kind: PortInput, ColorHex: 0x8b5cf6},
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
		Title:     tpl.Title,
		Subtitle:  tpl.Subtitle,
		Badge:     tpl.Badge,
		AccentHex: tpl.AccentHex,
		W:         tpl.W,
		H:         tpl.H,
		X:         world.X - tpl.W*0.5 + offset,
		Y:         world.Y - tpl.H*0.5 + offset,
		Status:    StatusIdle,
	}
	n.Ports = append([]Port(nil), tpl.Ports...)
	n.Fields = append([]NodeField(nil), tpl.Fields...)
	added := g.AddNode(n)
	g.SelectNode(added.ID, false)
	return added
}

// PresetAIAgent builds a multi-branch AI Agent orchestration graph.
func PresetAIAgent() *Graph {
	g := NewGraph()
	g.Nodes = []*Node{
		{
			ID: "trigger", Type: "trigger", Title: "User Prompt", Subtitle: "Chat / API ingress",
			Badge: "IN", AccentHex: 0x22c55e, X: 40, Y: 180, W: 210, H: 108, Status: StatusSuccess,
			Ports:  []Port{{ID: "out", Label: "msg", Side: PortRight, Kind: PortOutput, ColorHex: 0x22c55e}},
			Fields: []NodeField{{Key: "Channel", Value: "Desktop + Web"}},
		},
		{
			ID: "router", Type: "router", Title: "Intent Router", Subtitle: "Classify & dispatch",
			Badge: "IF", AccentHex: 0xf59e0b, X: 330, Y: 174, W: 224, H: 122, Status: StatusSuccess,
			Ports: []Port{
				{ID: "in", Label: "in", Side: PortLeft, Kind: PortInput, ColorHex: 0xf59e0b},
				{ID: "rag", Label: "recall", Side: PortRight, Kind: PortOutput, ColorHex: 0x06b6d4},
				{ID: "web", Label: "search", Side: PortRight, Kind: PortOutput, ColorHex: 0xec4899},
			},
			Fields: []NodeField{{Key: "Strategy", Value: "Hybrid Parallel"}},
		},
		{
			ID: "vector", Type: "tool", Title: "Vector Memory", Subtitle: "pgvector HNSW top-8",
			Badge: "DB", AccentHex: 0x06b6d4, X: 640, Y: 72, W: 228, H: 112, Status: StatusRunning,
			Ports: []Port{
				{ID: "q", Label: "query", Side: PortLeft, Kind: PortInput, ColorHex: 0x06b6d4},
				{ID: "docs", Label: "chunks", Side: PortRight, Kind: PortOutput, ColorHex: 0x06b6d4},
			},
			Fields: []NodeField{{Key: "Latency", Value: "14 ms"}},
		},
		{
			ID: "search", Type: "tool", Title: "Exa Web Search", Subtitle: "Live citation fetch",
			Badge: "WEB", AccentHex: 0xec4899, X: 640, Y: 266, W: 228, H: 112, Status: StatusRunning,
			Ports: []Port{
				{ID: "q", Label: "query", Side: PortLeft, Kind: PortInput, ColorHex: 0xec4899},
				{ID: "res", Label: "pages", Side: PortRight, Kind: PortOutput, ColorHex: 0xec4899},
			},
			Fields: []NodeField{{Key: "Results", Value: "Top 6 URLs"}},
		},
		{
			ID: "llm", Type: "llm", Title: "Synthesizer LLM", Subtitle: "Grounded final answer",
			Badge: "AI", AccentHex: 0x6366f1, X: 960, Y: 162, W: 242, H: 132, Status: StatusIdle,
			Ports: []Port{
				{ID: "ctx", Label: "memory", Side: PortLeft, Kind: PortInput, ColorHex: 0x06b6d4},
				{ID: "web", Label: "web", Side: PortLeft, Kind: PortInput, ColorHex: 0xec4899},
				{ID: "out", Label: "stream", Side: PortRight, Kind: PortOutput, ColorHex: 0x6366f1},
			},
			Fields: []NodeField{
				{Key: "Model", Value: "opus-5.5"},
				{Key: "Tokens", Value: "4k max"},
			},
		},
		{
			ID: "sink", Type: "output", Title: "UI Stream Sink", Subtitle: "Render tokens @ 60fps",
			Badge: "OUT", AccentHex: 0x8b5cf6, X: 1290, Y: 174, W: 214, H: 108, Status: StatusIdle,
			Ports:  []Port{{ID: "in", Label: "stream", Side: PortLeft, Kind: PortInput, ColorHex: 0x8b5cf6}},
			Fields: []NodeField{{Key: "Target", Value: "MyGo Native"}},
		},
	}
	g.Edges = []*Edge{
		{ID: "e1", Source: "trigger", SourcePort: "out", Target: "router", TargetPort: "in", Label: "JSON", Animated: true},
		{ID: "e2", Source: "router", SourcePort: "rag", Target: "vector", TargetPort: "q", Label: "embed", Animated: true, ColorHex: 0x06b6d4},
		{ID: "e3", Source: "router", SourcePort: "web", Target: "search", TargetPort: "q", Label: "live", Animated: true, ColorHex: 0xec4899},
		{ID: "e4", Source: "vector", SourcePort: "docs", Target: "llm", TargetPort: "ctx", ColorHex: 0x06b6d4},
		{ID: "e5", Source: "search", SourcePort: "res", Target: "llm", TargetPort: "web", ColorHex: 0xec4899},
		{ID: "e6", Source: "llm", SourcePort: "out", Target: "sink", TargetPort: "in", Label: "SSE", Animated: true, ColorHex: 0x6366f1},
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
			ID: "kafka", Type: "trigger", Title: "Kafka Stream", Subtitle: "orders.events.v3",
			Badge: "SRC", AccentHex: 0x22c55e, X: 60, Y: 80, W: 218, H: 108, Status: StatusSuccess,
			Ports:  []Port{{ID: "out", Label: "batch", Side: PortRight, Kind: PortOutput, ColorHex: 0x22c55e}},
			Fields: []NodeField{{Key: "Rate", Value: "18.4k msg/s"}},
		},
		{
			ID: "pg", Type: "tool", Title: "Postgres CDC", Subtitle: "users & pricing dim",
			Badge: "SQL", AccentHex: 0x06b6d4, X: 60, Y: 250, W: 218, H: 108, Status: StatusSuccess,
			Ports:  []Port{{ID: "out", Label: "rows", Side: PortRight, Kind: PortOutput, ColorHex: 0x06b6d4}},
			Fields: []NodeField{{Key: "Lag", Value: "42 ms"}},
		},
		{
			ID: "join", Type: "code", Title: "Enrich & Filter", Subtitle: "Deduplicate + currency",
			Badge: "ETL", AccentHex: 0xf59e0b, X: 380, Y: 162, W: 234, H: 122, Status: StatusRunning,
			Ports: []Port{
				{ID: "a", Label: "events", Side: PortLeft, Kind: PortInput, ColorHex: 0x22c55e},
				{ID: "b", Label: "dim", Side: PortLeft, Kind: PortInput, ColorHex: 0x06b6d4},
				{ID: "out", Label: "clean", Side: PortRight, Kind: PortOutput, ColorHex: 0xf59e0b},
			},
			Fields: []NodeField{{Key: "Window", Value: "5s tumbling"}},
		},
		{
			ID: "duck", Type: "tool", Title: "DuckDB Rollup", Subtitle: "Parquet columnar agg",
			Badge: "OLAP", AccentHex: 0x6366f1, X: 710, Y: 162, W: 226, H: 118, Status: StatusIdle,
			Ports: []Port{
				{ID: "in", Label: "in", Side: PortLeft, Kind: PortInput, ColorHex: 0xf59e0b},
				{ID: "bi", Label: "kpi", Side: PortRight, Kind: PortOutput, ColorHex: 0x6366f1},
				{ID: "alert", Label: "anom", Side: PortRight, Kind: PortOutput, ColorHex: 0xef4444},
			},
			Fields: []NodeField{{Key: "Memory", Value: "512 MB"}},
		},
		{
			ID: "dash", Type: "output", Title: "Metrics Board", Subtitle: "Live GMV & P99",
			Badge: "BI", AccentHex: 0x8b5cf6, X: 1030, Y: 80, W: 212, H: 104, Status: StatusIdle,
			Ports: []Port{{ID: "in", Label: "kpi", Side: PortLeft, Kind: PortInput, ColorHex: 0x6366f1}},
		},
		{
			ID: "pager", Type: "output", Title: "Anomaly Alert", Subtitle: "PagerDuty / Telegram",
			Badge: "OPS", AccentHex: 0xef4444, X: 1030, Y: 250, W: 212, H: 104, Status: StatusIdle,
			Ports: []Port{{ID: "in", Label: "alert", Side: PortLeft, Kind: PortInput, ColorHex: 0xef4444}},
		},
	}
	g.Edges = []*Edge{
		{ID: "d1", Source: "kafka", SourcePort: "out", Target: "join", TargetPort: "a", Animated: true},
		{ID: "d2", Source: "pg", SourcePort: "out", Target: "join", TargetPort: "b"},
		{ID: "d3", Source: "join", SourcePort: "out", Target: "duck", TargetPort: "in", Label: "Arrow", Animated: true},
		{ID: "d4", Source: "duck", SourcePort: "bi", Target: "dash", TargetPort: "in"},
		{ID: "d5", Source: "duck", SourcePort: "alert", Target: "pager", TargetPort: "in", ColorHex: 0xef4444},
	}
	return g
}

// PresetCICD builds a multi-stage build and progressive rollout pipeline.
func PresetCICD() *Graph {
	g := NewGraph()
	g.Background = BackgroundCross
	g.Nodes = []*Node{
		{
			ID: "git", Type: "trigger", Title: "Git Tag Push", Subtitle: "v1.4.0 on main",
			Badge: "GIT", AccentHex: 0x22c55e, X: 40, Y: 170, W: 204, H: 106, Status: StatusSuccess,
			Ports:  []Port{{ID: "out", Side: PortRight, Kind: PortOutput, ColorHex: 0x22c55e}},
			Fields: []NodeField{{Key: "Commit", Value: "a628e92"}},
		},
		{
			ID: "test", Type: "code", Title: "Go Test + Race", Subtitle: "go test -race ./...",
			Badge: "CI", AccentHex: 0x06b6d4, X: 330, Y: 70, W: 214, H: 108, Status: StatusSuccess,
			Ports: []Port{
				{ID: "in", Side: PortLeft, Kind: PortInput},
				{ID: "out", Side: PortRight, Kind: PortOutput},
			},
			Fields: []NodeField{{Key: "Passed", Value: "148 / 148"}},
		},
		{
			ID: "wasm", Type: "code", Title: "Build WASM + Mac", Subtitle: "darwin/arm64 + js/wasm",
			Badge: "BIN", AccentHex: 0x6366f1, X: 330, Y: 256, W: 214, H: 108, Status: StatusSuccess,
			Ports: []Port{
				{ID: "in", Side: PortLeft, Kind: PortInput},
				{ID: "out", Side: PortRight, Kind: PortOutput},
			},
			Fields: []NodeField{{Key: "Size", Value: "4.1 MB gz"}},
		},
		{
			ID: "canary", Type: "router", Title: "Canary 10%", Subtitle: "Error rate < 0.1%",
			Badge: "GATE", AccentHex: 0xf59e0b, X: 640, Y: 164, W: 214, H: 112, Status: StatusRunning,
			Ports: []Port{
				{ID: "a", Side: PortLeft, Kind: PortInput},
				{ID: "b", Side: PortLeft, Kind: PortInput},
				{ID: "out", Side: PortRight, Kind: PortOutput, ColorHex: 0x22c55e},
			},
			Fields: []NodeField{{Key: "Bake", Value: "15m"}},
		},
		{
			ID: "prod", Type: "output", Title: "Release Global", Subtitle: "GitHub + CDN edge",
			Badge: "LIVE", AccentHex: 0x8b5cf6, X: 940, Y: 166, W: 214, H: 108, Status: StatusIdle,
			Ports:  []Port{{ID: "in", Side: PortLeft, Kind: PortInput, ColorHex: 0x8b5cf6}},
			Fields: []NodeField{{Key: "Regions", Value: "US / EU / AP"}},
		},
	}
	g.Edges = []*Edge{
		{ID: "c1", Source: "git", SourcePort: "out", Target: "test", TargetPort: "in"},
		{ID: "c2", Source: "git", SourcePort: "out", Target: "wasm", TargetPort: "in"},
		{ID: "c3", Source: "test", SourcePort: "out", Target: "canary", TargetPort: "a", Animated: true},
		{ID: "c4", Source: "wasm", SourcePort: "out", Target: "canary", TargetPort: "b", Animated: true},
		{ID: "c5", Source: "canary", SourcePort: "out", Target: "prod", TargetPort: "in", Label: "promote", Animated: true},
	}
	return g
}

// TopologicalOrder returns node IDs ordered from sources to sinks (best-effort on cycles).
func (g *Graph) TopologicalOrder() []string {
	inDeg := make(map[string]int, len(g.Nodes))
	adj := make(map[string][]string, len(g.Nodes))
	for _, n := range g.Nodes {
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
	q := make([]string, 0, len(g.Nodes))
	for _, n := range g.Nodes {
		if inDeg[n.ID] == 0 {
			q = append(q, n.ID)
		}
	}
	seen := map[string]bool{}
	out := make([]string, 0, len(g.Nodes))
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
		if !seen[n.ID] {
			out = append(out, n.ID)
		}
	}
	return out
}

// SummaryString formats a concise status line for toolbars and tests.
func (g *Graph) SummaryString() string {
	return fmt.Sprintf("%d nodes · %d edges · %d%%", len(g.Nodes), len(g.Edges), int(g.Viewport.safeZoom()*100))
}
