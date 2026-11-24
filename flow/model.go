package flow

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// Point represents a 2D coordinate in either screen or world space.
type Point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

// Rect represents an axis-aligned rectangle in world or screen coordinates.
type Rect struct {
	X, Y, W, H float32
}

func (r Rect) Contains(p Point) bool {
	return p.X >= r.X && p.X <= r.X+r.W && p.Y >= r.Y && p.Y <= r.Y+r.H
}

func (r Rect) Intersects(o Rect) bool {
	return r.X <= o.X+o.W && r.X+r.W >= o.X && r.Y <= o.Y+o.H && r.Y+r.H >= o.Y
}

// Viewport maps between world space (infinite canvas) and screen space (DIPs).
type Viewport struct {
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
	Zoom float32 `json:"zoom"`
}

func DefaultViewport() Viewport {
	return Viewport{X: 0, Y: 0, Zoom: 1}
}

func (v Viewport) WorldToScreen(p Point) Point {
	z := v.safeZoom()
	return Point{X: p.X*z + v.X, Y: p.Y*z + v.Y}
}

func (v Viewport) ScreenToWorld(p Point) Point {
	z := v.safeZoom()
	return Point{X: (p.X - v.X) / z, Y: (p.Y - v.Y) / z}
}

func (v *Viewport) ZoomAt(screen Point, factor, minZoom, maxZoom float32) {
	if minZoom <= 0 {
		minZoom = 0.2
	}
	if maxZoom <= minZoom {
		maxZoom = 3.0
	}
	oldZoom := v.safeZoom()
	newZoom := clamp32(oldZoom*factor, minZoom, maxZoom)
	if math.Abs(float64(newZoom-oldZoom)) < 1e-4 {
		return
	}
	world := Point{
		X: (screen.X - v.X) / oldZoom,
		Y: (screen.Y - v.Y) / oldZoom,
	}
	v.Zoom = newZoom
	v.X = screen.X - world.X*newZoom
	v.Y = screen.Y - world.Y*newZoom
}

func (v Viewport) safeZoom() float32 {
	if v.Zoom <= 0.01 {
		return 1
	}
	return v.Zoom
}

// PortSide defines which edge of a node a port sits on.
type PortSide string

const (
	PortLeft   PortSide = "left"
	PortRight  PortSide = "right"
	PortTop    PortSide = "top"
	PortBottom PortSide = "bottom"
)

// PortKind distinguishes input sockets from output sockets.
type PortKind string

const (
	PortInput  PortKind = "input"
	PortOutput PortKind = "output"
)

// Port represents a connection handle on a node.
type Port struct {
	ID         string   `json:"id"`
	Label      string   `json:"label,omitempty"`
	Side       PortSide `json:"side"`
	Kind       PortKind `json:"kind"`
	DataType   string   `json:"dataType,omitempty"`
	ColorHex   uint32   `json:"color,omitempty"`
	ShowButton bool     `json:"showButton,omitempty"`
}

// NodeStatus displays live execution state on a node card.
type NodeStatus string

const (
	StatusIdle    NodeStatus = "idle"
	StatusRunning NodeStatus = "running"
	StatusSuccess NodeStatus = "success"
	StatusError   NodeStatus = "error"
)

// StatusVariant mirrors React Flow UI NodeStatusIndicator loadingVariant ("border" | "overlay").
type StatusVariant string

const (
	StatusBorder  StatusVariant = "border"
	StatusOverlay StatusVariant = "overlay"
)

// NodeField is a lightweight key/value line rendered inside a node card.
type NodeField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NodeVariant selects the React Flow UI visual card style of a Node.
type NodeVariant string

const (
	VariantBase        NodeVariant = "base"
	VariantSchema      NodeVariant = "schema"
	VariantGroup       NodeVariant = "group"
	VariantPlaceholder NodeVariant = "placeholder"
	VariantAnnotation  NodeVariant = "annotation"
)

// SchemaRow represents a typed column row inside a DatabaseSchemaNode.
type SchemaRow struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	InID   string `json:"inId,omitempty"`
	OutID  string `json:"outId,omitempty"`
	IsKey  bool   `json:"isKey,omitempty"`
}

// Node represents a draggable, connectable unit on the canvas.
type Node struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Variant       NodeVariant   `json:"variant,omitempty"`
	Title         string        `json:"title"`
	Subtitle      string        `json:"subtitle,omitempty"`
	Badge         string        `json:"badge,omitempty"`
	Footer        string        `json:"footer,omitempty"`
	Appendix      string        `json:"appendix,omitempty"`
	AppendixPos   PortSide      `json:"appendixPos,omitempty"`
	Tooltip       string        `json:"tooltip,omitempty"`
	X             float32       `json:"x"`
	Y             float32       `json:"y"`
	W             float32       `json:"w"`
	H             float32       `json:"h"`
	AccentHex     uint32        `json:"accent,omitempty"`
	Ports         []Port        `json:"ports,omitempty"`
	Fields        []NodeField   `json:"fields,omitempty"`
	SchemaRows    []SchemaRow   `json:"schemaRows,omitempty"`
	Selected      bool          `json:"selected,omitempty"`
	Locked        bool          `json:"locked,omitempty"`
	Status        NodeStatus    `json:"status,omitempty"`
	StatusVariant StatusVariant `json:"statusVariant,omitempty"`
}

func (n *Node) Bounds() Rect {
	w := n.W
	if w <= 0 {
		w = 224
	}
	h := n.H
	if n.Variant == VariantSchema && len(n.SchemaRows) > 0 {
		minH := float32(34 + len(n.SchemaRows)*26 + 8)
		if h < minH {
			h = minH
		}
	}
	if h <= 0 {
		h = 112
	}
	return Rect{X: n.X, Y: n.Y, W: w, H: h}
}

// syncSchemaHandles gives every row of a schema node a target handle on the
// left and a source handle on the right, as the LabeledHandles of React Flow
// UI's DatabaseSchemaNode do. Rows keep the ids they already have.
func (n *Node) syncSchemaHandles() {
	if n.Variant != VariantSchema {
		return
	}
	has := map[string]bool{}
	for _, p := range n.Ports {
		has[p.ID] = true
	}
	add := func(id string, side PortSide, kind PortKind) {
		if !has[id] {
			n.Ports = append(n.Ports, Port{ID: id, Side: side, Kind: kind})
			has[id] = true
		}
	}
	for i := range n.SchemaRows {
		r := &n.SchemaRows[i]
		if r.InID == "" {
			r.InID = r.Name + "-in"
		}
		if r.OutID == "" {
			r.OutID = r.Name + "-out"
		}
		add(r.InID, PortLeft, PortInput)
		add(r.OutID, PortRight, PortOutput)
	}
}

// EdgeKind selects the curve routing algorithm for edges.
type EdgeKind string

const (
	EdgeBezier     EdgeKind = "bezier"
	EdgeSmoothStep EdgeKind = "smoothstep"
	EdgeStraight   EdgeKind = "straight"
)

// Edge connects an output port on Source to an input port on Target.
type Edge struct {
	ID              string   `json:"id"`
	Source          string   `json:"source"`
	SourcePort      string   `json:"sourcePort"`
	Target          string   `json:"target"`
	TargetPort      string   `json:"targetPort"`
	Label           string   `json:"label,omitempty"`
	Kind            EdgeKind `json:"kind,omitempty"`
	Animated        bool     `json:"animated,omitempty"`
	PacketMotion    bool     `json:"packetMotion,omitempty"`
	PacketDuration  float32  `json:"packetDuration,omitempty"`
	PacketDirection string   `json:"packetDirection,omitempty"` // "forward", "reverse", "alternate"
	ColorHex        uint32   `json:"color,omitempty"`
	Selected        bool     `json:"selected,omitempty"`
}

// ChangeEvent records graph mutations for React Flow UI DevTools ChangeLogger.
type ChangeEvent struct {
	Kind   string
	Target string
	Detail string
}

// BackgroundKind selects the canvas backdrop pattern.
type BackgroundKind string

const (
	BackgroundDots  BackgroundKind = "dots"
	BackgroundLines BackgroundKind = "lines"
	BackgroundCross BackgroundKind = "cross"
	BackgroundNone  BackgroundKind = "none"
)

// LayoutDirection configures automatic topological layout orientation.
type LayoutDirection string

const (
	LayoutLR LayoutDirection = "LR"
	LayoutTB LayoutDirection = "TB"
)

// Theme defines all surface, edge, port, selection, and minimap colors.
type Theme struct {
	Dark             bool
	CanvasBG         ui.Color
	GridColor        ui.Color
	GridMajorColor   ui.Color
	NodeBG           ui.Color
	NodeHeaderBG     ui.Color
	NodeBorder       ui.Color
	NodeSelected     ui.Color
	NodeShadow       ui.Color
	TextPrimary      ui.Color
	TextMuted        ui.Color
	BadgeBG          ui.Color
	BadgeText        ui.Color
	PortFill         ui.Color
	PortBorder       ui.Color
	PortSnap         ui.Color
	EdgeDefault      ui.Color
	EdgeHover        ui.Color
	EdgeSelected     ui.Color
	EdgeLabelBG      ui.Color
	EdgeLabelBorder  ui.Color
	EdgeLabelText    ui.Color
	SelectionFill    ui.Color
	SelectionBorder  ui.Color
	GuideLine        ui.Color
	MiniMapBG        ui.Color
	MiniMapBorder    ui.Color
	MiniMapNode      ui.Color
	MiniMapNodeSel   ui.Color
	MiniMapViewport  ui.Color
	MiniMapViewEdge  ui.Color
	ControlsBG       ui.Color
	ControlsBorder   ui.Color
	ControlsText     ui.Color
	ControlsActiveBG ui.Color
	ControlsActiveFG ui.Color
	StatusIdle       ui.Color
	StatusRunning    ui.Color
	StatusSuccess    ui.Color
	StatusError      ui.Color
}

// RGBHex converts a 0xRRGGBB integer literal into an opaque ui.Color.
func RGBHex(v uint32) ui.Color {
	return ui.RGB(uint8((v>>16)&0xff), uint8((v>>8)&0xff), uint8(v&0xff))
}

// DarkTheme returns a refined obsidian canvas palette inspired by Linear & XYFlow.
func DarkTheme() Theme {
	return Theme{
		Dark:             true,
		CanvasBG:         RGBHex(0x0b0c10),
		GridColor:        ui.RGBA(255, 255, 255, 0.13),
		GridMajorColor:   ui.RGBA(255, 255, 255, 0.24),
		NodeBG:           RGBHex(0x14161d),
		NodeHeaderBG:     RGBHex(0x1a1d26),
		NodeBorder:       ui.RGBA(255, 255, 255, 0.11),
		NodeSelected:     RGBHex(0x6366f1),
		NodeShadow:       ui.RGBA(0, 0, 0, 0.52),
		TextPrimary:      RGBHex(0xf4f4f5),
		TextMuted:        RGBHex(0x9ca3af),
		BadgeBG:          ui.RGBA(255, 255, 255, 0.08),
		BadgeText:        RGBHex(0xd4d4d8),
		PortFill:         RGBHex(0x27272a),
		PortBorder:       RGBHex(0x27272a),
		PortSnap:         RGBHex(0x22c55e),
		EdgeDefault:      RGBHex(0x525866),
		EdgeHover:        RGBHex(0x9ca3af),
		EdgeSelected:     RGBHex(0x818cf8),
		EdgeLabelBG:      RGBHex(0x181b24),
		EdgeLabelBorder:  ui.RGBA(255, 255, 255, 0.14),
		EdgeLabelText:    RGBHex(0xe4e4e7),
		SelectionFill:    ui.RGBA(99, 102, 241, 0.15),
		SelectionBorder:  ui.RGBA(129, 140, 248, 0.82),
		GuideLine:        ui.RGBA(244, 114, 182, 0.78),
		MiniMapBG:        ui.RGBA(18, 20, 28, 0.92),
		MiniMapBorder:    ui.RGBA(255, 255, 255, 0.12),
		MiniMapNode:      ui.RGBA(148, 163, 184, 0.55),
		MiniMapNodeSel:   RGBHex(0x818cf8),
		MiniMapViewport:  ui.RGBA(99, 102, 241, 0.18),
		MiniMapViewEdge:  ui.RGBA(129, 140, 248, 0.85),
		ControlsBG:       ui.RGBA(20, 22, 29, 0.94),
		ControlsBorder:   ui.RGBA(255, 255, 255, 0.11),
		ControlsText:     RGBHex(0xe4e4e7),
		ControlsActiveBG: RGBHex(0x6366f1),
		ControlsActiveFG: RGBHex(0xffffff),
		StatusIdle:       RGBHex(0x64748b),
		StatusRunning:    RGBHex(0xf59e0b),
		StatusSuccess:    RGBHex(0x22c55e),
		StatusError:      RGBHex(0xef4444),
	}
}

// LightTheme returns a crisp neutral canvas palette matching coss/XYFlow light mode.
func LightTheme() Theme {
	return Theme{
		Dark:             false,
		CanvasBG:         RGBHex(0xf7f8fa),
		GridColor:        ui.RGBA(15, 23, 42, 0.13),
		GridMajorColor:   ui.RGBA(15, 23, 42, 0.24),
		NodeBG:           RGBHex(0xffffff),
		NodeHeaderBG:     RGBHex(0xf8fafc),
		NodeBorder:       ui.RGBA(15, 23, 42, 0.11),
		NodeSelected:     RGBHex(0x4f46e5),
		NodeShadow:       ui.RGBA(15, 23, 42, 0.11),
		TextPrimary:      RGBHex(0x0f172a),
		TextMuted:        RGBHex(0x64748b),
		BadgeBG:          ui.RGBA(15, 23, 42, 0.06),
		BadgeText:        RGBHex(0x334155),
		PortFill:         RGBHex(0xf1f5f9),
		PortBorder:       RGBHex(0xcbd5e1),
		PortSnap:         RGBHex(0x16a34a),
		EdgeDefault:      RGBHex(0x94a3b8),
		EdgeHover:        RGBHex(0x475569),
		EdgeSelected:     RGBHex(0x4f46e5),
		EdgeLabelBG:      RGBHex(0xffffff),
		EdgeLabelBorder:  ui.RGBA(15, 23, 42, 0.12),
		EdgeLabelText:    RGBHex(0x334155),
		SelectionFill:    ui.RGBA(79, 70, 229, 0.10),
		SelectionBorder:  ui.RGBA(79, 70, 229, 0.75),
		GuideLine:        ui.RGBA(236, 72, 153, 0.72),
		MiniMapBG:        ui.RGBA(255, 255, 255, 0.93),
		MiniMapBorder:    ui.RGBA(15, 23, 42, 0.11),
		MiniMapNode:      ui.RGBA(100, 116, 139, 0.52),
		MiniMapNodeSel:   RGBHex(0x4f46e5),
		MiniMapViewport:  ui.RGBA(79, 70, 229, 0.12),
		MiniMapViewEdge:  ui.RGBA(79, 70, 229, 0.78),
		ControlsBG:       ui.RGBA(255, 255, 255, 0.96),
		ControlsBorder:   ui.RGBA(15, 23, 42, 0.10),
		ControlsText:     RGBHex(0x1e293b),
		ControlsActiveBG: RGBHex(0x4f46e5),
		ControlsActiveFG: RGBHex(0xffffff),
		StatusIdle:       RGBHex(0x94a3b8),
		StatusRunning:    RGBHex(0xd97706),
		StatusSuccess:    RGBHex(0x16a34a),
		StatusError:      RGBHex(0xdc2626),
	}
}

func colorWithAlpha(c ui.Color, factor float32) ui.Color {
	a := (float32(c.A) / 255.0) * factor
	return ui.RGBA(c.R, c.G, c.B, clamp32(a, 0, 1))
}

func clamp32(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func dist2D(a, b Point) float32 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return float32(math.Hypot(float64(dx), float64(dy)))
}
