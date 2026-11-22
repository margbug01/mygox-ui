package motion

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// ============================================================================
// 17. AnimatedBackground
// ============================================================================

// AnimatedBackgroundOptions configures AnimatedBackground.
type AnimatedBackgroundOptions struct {
	EnableHover bool
	Radius      float32
	PillColor   ui.Color
	TrackColor  ui.Color
	Padding     float32
	OnChange    func(index int)
	Spring      ui.SpringOptions
}

type animBgState struct {
	rects   []ui.Rect
	hovered int
}

// AnimatedBackground renders a horizontal group of items with a single shared
// spring-sliding highlight pill that glides smoothly behind the active (or hovered) item.
func AnimatedBackground(c *ui.Context, items []string, active *int, opts ...AnimatedBackgroundOptions) ui.Element {
	var o AnimatedBackgroundOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Radius <= 0 {
		o.Radius = 8
	}
	if o.Padding == 0 {
		o.Padding = 3
	}
	trackCol := colorOr(o.TrackColor, t.SurfaceSubtle)
	defaultPill := t.Surface
	if t.Dark {
		defaultPill = t.SurfaceHover
	}
	pillCol := colorOr(o.PillColor, defaultPill)
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 380, Damping: 32, Mass: 0.7})

	activeIdx := 0
	if active != nil {
		activeIdx = *active
	}
	if len(items) > 0 && (activeIdx < 0 || activeIdx >= len(items)) {
		activeIdx = 0
	}

	track := ui.Row(c).
		AlignItems(ui.Center).
		Padding(o.Padding).
		Radius(o.Radius + 2).
		Background(trackCol).
		Border(1, t.Border)

	st := ui.Local(track, "ab-state", func() animBgState {
		return animBgState{hovered: -1}
	})
	if len(st.rects) != len(items) {
		st.rects = make([]ui.Rect, len(items))
	}

	targetIdx := activeIdx
	if o.EnableHover && st.hovered >= 0 && st.hovered < len(items) {
		targetIdx = st.hovered
	}

	tb := track.Bounds()
	var targetX, targetW, targetH float32
	hasLayout := false
	if targetIdx >= 0 && targetIdx < len(st.rects) && st.rects[targetIdx].W > 0 && tb.W > 0 {
		r := st.rects[targetIdx]
		targetX = r.X - tb.X
		targetW = r.W
		targetH = r.H
		hasLayout = true
	}

	px := track.Spring("ab-x", targetX, spring)
	pw := track.Spring("ab-w", targetW, spring)
	if !hasLayout {
		px = track.Jump("ab-x", targetX)
		pw = track.Jump("ab-w", targetW)
	}

	st.hovered = -1
	track.Children(func() {
		if hasLayout && pw > 1 {
			pill := ui.Box(c).
				Key("pill").
				Absolute().
				Left(px).
				Top(o.Padding).
				Width(pw).
				Height(max(targetH, 26)).
				Radius(o.Radius).
				Background(pillCol).
				PassThrough()
			if !t.Dark {
				pill.Shadow(0, 1, 2, 0, ui.RGBA(9, 9, 11, 0.08))
			} else {
				pill.Border(1, t.BorderStrong)
			}
		}

		for i, label := range items {
			idx := i
			isSel := idx == targetIdx
			fg := t.TextMuted
			if isSel {
				fg = t.Text
			}
			btn := ui.Box(c).
				Key(idx).
				ZIndex(1).
				Padding(6, 13).
				Radius(o.Radius).
				Cursor(ui.CursorPointer)
			btn.Children(func() {
				ui.Text(c, label).
					FontSize(12.5).
					FontWeight(500).
					TextColor(fg).
					SingleLine()
			})
			if b := btn.Bounds(); b.W > 0 {
				st.rects[idx] = b
			}
			if btn.Hovered() {
				st.hovered = idx
			}
			if btn.Clicked() {
				if active != nil {
					*active = idx
				}
				if o.OnChange != nil {
					o.OnChange(idx)
				}
			}
		}
	})
	return track
}

// ============================================================================
// 18. Dock
// ============================================================================

// DockItem describes a single icon entry in Dock.
type DockItem struct {
	Label   string
	Icon    *ui.SVG
	Glyph   string
	OnClick func()
}

// DockOptions configures Dock.
type DockOptions struct {
	BaseSize      float32
	Magnification float32
	Distance      float32
	ActiveIndex   int
	Spring        ui.SpringOptions
}

// Dock renders a macOS-style Gaussian spring magnification dock in a restrained
// monochrome zinc aesthetic with floating tooltips above hovered icons.
func Dock(c *ui.Context, items []DockItem, opts ...DockOptions) ui.Element {
	var o DockOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.BaseSize <= 0 {
		o.BaseSize = 40
	}
	if o.Magnification <= 0 {
		o.Magnification = 62
	}
	if o.Distance <= 0 {
		o.Distance = 120
	}
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 220, Damping: 18, Mass: 0.25})

	gap := float32(10)
	padH := float32(14)

	outer := ui.Box(c).
		Height(o.Magnification + 36).
		AlignItems(ui.End).
		Justify(ui.Center)

	outer.Children(func() {
		bar := ui.Row(c).
			AlignItems(ui.End).
			Gap(gap).
			Padding(8, padH).
			Radius(16).
			Background(t.SurfaceSubtle).
			Border(1, t.Border)
		if !t.Dark {
			bar.Shadow(0, 6, 20, -4, ui.RGBA(9, 9, 11, 0.08))
		}

		ptrX, _, hasPtr := bar.PointerPosition()
		barHovered := bar.Hovered()

		bar.Children(func() {
			for i, item := range items {
				idx := i
				centerX := padH + float32(idx)*(o.BaseSize+gap) + o.BaseSize*0.5
				targetSz := o.BaseSize
				if hasPtr && barHovered {
					dist := float64(ptrX - centerX)
					sigma := float64(o.Distance) * 0.48
					mag := math.Exp(-(dist * dist) / (2 * sigma * sigma))
					targetSz = o.BaseSize + float32(mag)*(o.Magnification-o.BaseSize)
				}

				slot := ui.Column(c).
					Key(idx).
					AlignItems(ui.Center).
					Justify(ui.End).
					Gap(5).
					Cursor(ui.CursorPointer)

				sz := slot.Spring("dock-sz", targetSz, spring)
				hov := slot.Hovered()
				tipAlpha := slot.Spring("dock-tip", boolFloat(hov), SpringSnappy)

				wellBg := t.SurfaceHover
				iconFg := t.TextMuted
				if hov || idx == o.ActiveIndex {
					iconFg = t.Text
				}
				if !t.Dark {
					wellBg = ui.Hex("#e4e4e7")
					if hov {
						wellBg = ui.Hex("#d4d4d8")
					}
				} else if hov {
					wellBg = ui.Hex("#3f3f46")
				}

				slot.Children(func() {
					if tipAlpha > 0.05 && item.Label != "" {
						ui.Box(c).
							Absolute().
							Top(-26).
							Padding(2, 7).
							Radius(6).
							Background(t.Surface).
							Border(1, t.Border).
							Opacity(tipAlpha).
							Translate(0, (1-tipAlpha)*4).
							PassThrough().
							Children(func() {
								ui.Text(c, item.Label).FontSize(10.5).FontWeight(500).TextColor(t.Text).SingleLine()
							})
					}

					iconSize := sz * 0.44
					ui.Box(c).
						Width(sz).
						Height(sz).
						Radius(999).
						Background(wellBg).
						Center().
						Children(func() {
							if item.Icon != nil {
								ui.Icon(c, item.Icon).Size(iconSize, iconSize).TextColor(iconFg)
							} else {
								ui.Text(c, item.Glyph).FontSize(iconSize).FontWeight(600).TextColor(iconFg).SingleLine()
							}
						})

					// Minimal active indicator dot
					dotOp := float32(0)
					if idx == o.ActiveIndex {
						dotOp = 1
					}
					ui.Box(c).Width(4).Height(4).Radius(999).Background(t.TextMuted).Opacity(dotOp)
				})

				if slot.Clicked() && item.OnClick != nil {
					item.OnClick()
				}
			}
		})
	})
	return outer
}

// ============================================================================
// 19. Magnetic
// ============================================================================

// MagneticOptions configures Magnetic.
type MagneticOptions struct {
	Intensity float32
	Range     float32
	Spring    ui.SpringOptions
}

// Magnetic wraps content in a spring-physics container that smoothly attracts
// toward the pointer within Range and springs back to rest on leave.
func Magnetic(c *ui.Context, content func(), opts ...MagneticOptions) ui.Element {
	var o MagneticOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Intensity <= 0 {
		o.Intensity = 0.45
	}
	if o.Range <= 0 {
		o.Range = 90
	}
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 260, Damping: 14, Mass: 0.4})

	// Outer capture zone provides a comfortable attraction radius around the button.
	zone := ui.Box(c).Padding(16).Center()
	b := zone.Bounds()
	var tx, ty float32
	if px, py, ok := zone.PointerPosition(); ok && zone.Hovered() && b.W > 0 && b.H > 0 {
		dx := px - b.W*0.5
		dy := py - b.H*0.5
		dist := float32(math.Hypot(float64(dx), float64(dy)))
		if dist <= o.Range {
			scale := (1 - dist/o.Range) * o.Intensity
			tx = dx * scale
			ty = dy * scale
		}
	}

	sx := zone.Spring("mag-x", tx, spring)
	sy := zone.Spring("mag-y", ty, spring)

	zone.Children(func() {
		ui.Box(c).Translate(sx, sy).Children(func() {
			if content != nil {
				content()
			}
		})
	})
	return zone
}

// ============================================================================
// 20. Cursor
// ============================================================================

// CursorOptions configures Cursor.
type CursorOptions struct {
	OffsetX float32
	OffsetY float32
	Radius  float32
	Height  float32
	Spring  ui.SpringOptions
}

// Cursor renders an interactive stage where a custom floating badge/pill smoothly
// follows the pointer via spring physics while hovering over the stage.
func Cursor(c *ui.Context, stageContent func(), cursorContent func(), opts ...CursorOptions) ui.Element {
	var o CursorOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Radius <= 0 {
		o.Radius = 12
	}
	if o.Height <= 0 {
		o.Height = 148
	}
	if o.OffsetX == 0 && o.OffsetY == 0 {
		o.OffsetX = 12
		o.OffsetY = 12
	}
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 300, Damping: 26, Mass: 0.45})

	stage := ui.Box(c).
		WidthPercent(100).
		Height(o.Height).
		Radius(o.Radius).
		Background(t.SurfaceSubtle).
		Border(1, t.Border).
		Clip().
		Center()

	b := stage.Bounds()
	tx := b.W * 0.5
	ty := b.H * 0.5
	hov := stage.Hovered()
	if px, py, ok := stage.PointerPosition(); ok && hov {
		tx = px + o.OffsetX
		ty = py + o.OffsetY
	}
	sx := stage.Spring("cur-x", tx, spring)
	sy := stage.Spring("cur-y", ty, spring)
	vis := stage.Spring("cur-vis", boolFloat(hov), SpringSnappy)

	stage.Children(func() {
		if stageContent != nil {
			stageContent()
		}
		if vis > 0.02 && cursorContent != nil {
			sc := lerp32(0.7, 1.0, vis)
			ui.Box(c).
				Key("floating-cursor").
				Absolute().
				Left(sx).
				Top(sy).
				Opacity(vis).
				Scale(sc, sc).
				ZIndex(10).
				PassThrough().
				Children(cursorContent)
		}
	})
	return stage
}

// ============================================================================
// 21. ToolbarDynamic
// ============================================================================

// ToolbarDynamicOptions configures ToolbarDynamic.
type ToolbarDynamicOptions struct {
	Open          *bool
	Placeholder   string
	CompactWidth  float32
	ExpandedWidth float32
	UserIcon      *ui.SVG
	SearchIcon    *ui.SVG
	BackIcon      *ui.SVG
	Spring        ui.SpringOptions
}

type dynToolbarLocal struct {
	open bool
}

// ToolbarDynamic renders a compact monochrome pill toolbar that smoothly springs
// horizontally into a live search input bar on click, matching upstream toolbar-dynamic.
func ToolbarDynamic(c *ui.Context, query *string, opts ...ToolbarDynamicOptions) ui.Element {
	var o ToolbarDynamicOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Placeholder == "" {
		o.Placeholder = "Search notes"
	}
	if o.CompactWidth <= 0 {
		o.CompactWidth = 94
	}
	if o.ExpandedWidth <= 0 {
		o.ExpandedWidth = 272
	}
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 340, Damping: 28, Mass: 0.7})

	pill := ui.Box(c).
		Radius(12).
		Background(t.Surface).
		Border(1, t.Border).
		Clip()
	if !t.Dark {
		pill.Shadow(0, 2, 10, -2, ui.RGBA(9, 9, 11, 0.07))
	}

	local := ui.Local(pill, "td-open", func() dynToolbarLocal { return dynToolbarLocal{} })
	isOpen := local.open
	if o.Open != nil {
		isOpen = *o.Open
	}
	setOpen := func(v bool) {
		local.open = v
		if o.Open != nil {
			*o.Open = v
		}
	}

	targetW := o.CompactWidth
	if isOpen {
		targetW = o.ExpandedWidth
	}
	curW := pill.Spring("td-w", targetW, spring)
	pill.Width(curW).Padding(6)

	iconBtn := func(key string, icon *ui.SVG, fallback string, onClick func()) {
		b := ui.Box(c).
			Key(key).
			Width(34).
			Height(34).
			Shrink(0).
			Radius(8).
			Center().
			Cursor(ui.CursorPointer)
		fg := t.TextMuted
		if b.Hovered() {
			b.Background(t.SurfaceHover)
			fg = t.Text
		}
		b.Children(func() {
			if icon != nil {
				ui.Icon(c, icon).Size(16, 16).TextColor(fg)
			} else {
				ui.Text(c, fallback).FontSize(14).FontWeight(600).TextColor(fg).SingleLine()
			}
		})
		if b.Clicked() && onClick != nil {
			onClick()
		}
	}

	pill.Children(func() {
		if !isOpen {
			ui.Row(c).Key("closed").AlignItems(ui.Center).Gap(6).Children(func() {
				iconBtn("user", o.UserIcon, "U", nil)
				iconBtn("search", o.SearchIcon, "S", func() { setOpen(true) })
			})
		} else {
			ui.Row(c).Key("open").WidthPercent(100).AlignItems(ui.Center).Gap(6).Children(func() {
				iconBtn("back", o.BackIcon, "←", func() { setOpen(false) })
				q := query
				if q == nil {
					dummy := ""
					q = &dummy
				}
				inp := ui.TextInput(c, q).
					Grow(1).
					MinWidth(0).
					Height(34).
					Padding(0, 10).
					Radius(8).
					Background(t.SurfaceSubtle).
					Border(1, t.Border).
					FontSize(12.5).
					TextColor(t.Text).
					Placeholder(o.Placeholder).
					PlaceholderColor(t.TextSubtle)
				if inp.Shortcut(0, ui.KeyEscape) {
					setOpen(false)
				}
			})
		}
	})
	return pill
}

// ============================================================================
// 22. ToolbarExpandable
// ============================================================================

// ExpandableToolbarItem defines one tab/drawer in ToolbarExpandable.
type ExpandableToolbarItem struct {
	ID      string
	Label   string
	Icon    *ui.SVG
	Glyph   string
	Content func()
}

// ToolbarExpandableOptions configures ToolbarExpandable.
type ToolbarExpandableOptions struct {
	Width       float32
	PanelHeight float32
	Spring      ui.SpringOptions
}

// ToolbarExpandable renders a monochrome icon toolbar with a smoothly expanding
// drawer panel directly above it whenever an item is selected.
func ToolbarExpandable(c *ui.Context, items []ExpandableToolbarItem, activeIndex *int, opts ...ToolbarExpandableOptions) ui.Element {
	var o ToolbarExpandableOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 280
	}
	if o.PanelHeight <= 0 {
		o.PanelHeight = 112
	}
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 320, Damping: 28, Mass: 0.7})

	active := -1
	if activeIndex != nil {
		active = *activeIndex
	}
	isOpen := active >= 0 && active < len(items)

	shell := ui.Column(c).
		Width(o.Width).
		Radius(12).
		Background(t.Surface).
		Border(1, t.Border).
		Clip()
	if !t.Dark {
		shell.Shadow(0, 2, 10, -2, ui.RGBA(9, 9, 11, 0.07))
	}

	targetH := float32(0)
	if isOpen {
		targetH = o.PanelHeight
	}
	drawerH := clamp32(shell.Spring("te-h", targetH, spring), 0, o.PanelHeight*1.2)
	drawerOp := clamp32(drawerH/max(o.PanelHeight, 1), 0, 1)

	shell.Children(func() {
		if drawerH > 1 {
			ui.Box(c).
				Key("drawer").
				WidthPercent(100).
				Height(drawerH).
				BorderWidth(0, 0, 1, 0).
				BorderColor(t.Border).
				Clip().
				Opacity(drawerOp).
				Children(func() {
					if isOpen && items[active].Content != nil {
						ui.Column(c).WidthPercent(100).Padding(12, 14).Gap(6).Children(items[active].Content)
					}
				})
		}

		ui.Row(c).
			Key("bar").
			WidthPercent(100).
			Padding(6).
			Justify(ui.SpaceAround).
			AlignItems(ui.Center).
			Children(func() {
				for i, item := range items {
					idx := i
					sel := idx == active
					btn := ui.Box(c).
						Key(idx).
						Width(38).
						Height(34).
						Radius(8).
						Center().
						Cursor(ui.CursorPointer)
					fg := t.TextMuted
					if sel {
						btn.Background(t.SurfaceHover)
						fg = t.Text
					} else if btn.Hovered() {
						btn.Background(t.SurfaceSubtle)
						fg = t.Text
					}
					btn.Children(func() {
						if item.Icon != nil {
							ui.Icon(c, item.Icon).Size(16, 16).TextColor(fg)
						} else {
							ui.Text(c, item.Glyph).FontSize(13).FontWeight(600).TextColor(fg).SingleLine()
						}
					})
					if btn.Clicked() && activeIndex != nil {
						if *activeIndex == idx {
							*activeIndex = -1
						} else {
							*activeIndex = idx
						}
					}
				}
			})
	})
	return shell
}

func boolFloat(v bool) float32 {
	if v {
		return 1
	}
	return 0
}
