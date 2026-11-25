package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

func hcSeriesColors(t Theme) [6]ui.Color {
	if t.Dark {
		return [6]ui.Color{
			t.Accent,
			ui.Hex("#38bdf8"),
			t.Success,
			t.Warning,
			ui.Hex("#c084fc"),
			ui.Hex("#fb7185"),
		}
	}
	return [6]ui.Color{
		t.Accent,
		ui.Hex("#0284c7"),
		t.Success,
		t.Warning,
		ui.Hex("#7c3aed"),
		ui.Hex("#e11d48"),
	}
}

func hcShade(c ui.Color, factor float32) ui.Color {
	scale := func(v uint8) uint8 {
		n := float32(v) * factor
		if n < 0 {
			return 0
		}
		if n > 255 {
			return 255
		}
		return uint8(n + 0.5)
	}
	return ui.Color{R: scale(c.R), G: scale(c.G), B: scale(c.B), A: c.A}
}

func hcFmtMoney(v float64) string {
	abs := math.Abs(v)
	sign := ""
	if v < 0 {
		sign = "-"
	}
	if abs >= 1000000 {
		return fmt.Sprintf("%s$%.1fM", sign, abs/1000000)
	}
	if abs >= 10000 {
		return fmt.Sprintf("%s$%.1fk", sign, abs/1000)
	}
	iv := int(abs + 0.5)
	if iv >= 1000 {
		return fmt.Sprintf("%s$%d,%03d", sign, iv/1000, iv%1000)
	}
	return fmt.Sprintf("%s$%d", sign, iv)
}

func hcFmtComma(v int) string {
	if v < 0 {
		return "-" + hcFmtComma(-v)
	}
	if v >= 1000000 {
		return fmt.Sprintf("%d,%03d,%03d", v/1000000, (v/1000)%1000, v%1000)
	}
	if v >= 1000 {
		return fmt.Sprintf("%d,%03d", v/1000, v%1000)
	}
	return fmt.Sprintf("%d", v)
}

func hcFmtTokensShort(v int) string {
	if v >= 1000 {
		return fmt.Sprintf("%.1fk", float64(v)/1000.0)
	}
	return fmt.Sprintf("%d", v)
}

func hcCubicY(y0, y1, u float32) float32 {
	// Horizontal cubic Bezier with symmetric tangents at x0+0.5*dx, x1-0.5*dx
	// yields smoothstep-like vertical interpolation: 3u^2 - 2u^3.
	u = clamp32(u, 0, 1)
	return y0 + (y1-y0)*(u*u*(3-2*u))
}

// ---------------------------------------------------------------------------
// 6. stacked-bars ("Stacked bars")
// ---------------------------------------------------------------------------

// StackedBarColumn defines one column of stacked values.
type StackedBarColumn struct {
	Label  string
	Values []float64
}

// StackedBarsOptions configures StackedBars.
type StackedBarsOptions struct {
	Width   float32
	Height  float32
	Title   string
	Series  []string
	Colors  []ui.Color
	Columns []StackedBarColumn
}

// StackedBars renders weekly stacked column bars where hovering a column dims
// siblings and updates the header breakdown legend readout.
func StackedBars(c *ui.Context, opts ...StackedBarsOptions) ui.Element {
	var o StackedBarsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 292
	}
	if o.Height <= 0 {
		o.Height = 114
	}
	if o.Title == "" {
		o.Title = "Weekly mix"
	}
	palette := hcSeriesColors(t)
	series := o.Series
	if len(series) == 0 {
		series = []string{"Web", "Mob", "API"}
	}
	colors := o.Colors
	if len(colors) < len(series) {
		colors = make([]ui.Color, len(series))
		for i := range series {
			colors[i] = palette[i%len(palette)]
		}
	}
	cols := o.Columns
	if len(cols) == 0 {
		cols = []StackedBarColumn{
			{Label: "Mon", Values: []float64{34, 24, 16}},
			{Label: "Tue", Values: []float64{42, 31, 22}},
			{Label: "Wed", Values: []float64{38, 28, 19}},
			{Label: "Thu", Values: []float64{56, 39, 27}},
			{Label: "Fri", Values: []float64{49, 34, 25}},
			{Label: "Sat", Values: []float64{29, 44, 15}},
			{Label: "Sun", Values: []float64{24, 38, 14}},
		}
	}

	maxSum := 1.0
	for _, col := range cols {
		s := 0.0
		for _, v := range col.Values {
			s += v
		}
		if s > maxSum {
			maxSum = s
		}
	}

	canvasW := o.Width - 24
	if canvasW < 180 {
		canvasW = 180
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(10, 12).
		Gap(6).
		Radius(12).
		Background(t.Surface).
		Border(1, t.Border)

	hoverPtr := ui.Local(card, "hc_stack_hover", func() int { return -1 })

	card.Children(func() {
		active := *hoverPtr

		// Header + crisp breakdown legend readout.
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			if active >= 0 && active < len(cols) {
				sum := 0.0
				for _, v := range cols[active].Values {
					sum += v
				}
				ui.Text(c, fmt.Sprintf("%s · %.0fk", cols[active].Label, sum)).
					FontSize(12).FontWeight(600).TextColor(t.Text)
			} else {
				ui.Text(c, o.Title).FontSize(12).FontWeight(600).TextColor(t.Text)
			}

			ui.Row(c).Gap(7).AlignItems(ui.Center).Children(func() {
				for si, name := range series {
					col := colors[si]
					lbl := name
					if active >= 0 && active < len(cols) && si < len(cols[active].Values) {
						lbl = fmt.Sprintf("%.0f", cols[active].Values[si])
					}
					ui.Row(c).Gap(3).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Size(6, 6).Radius(2).Background(col)
						txtCol := t.TextSubtle
						if active >= 0 {
							txtCol = t.Text
						}
						ui.Text(c, lbl).FontSize(10).FontWeight(600).TextColor(txtCol)
					})
				}
			})
		})

		stage := ui.Box(c).
			Size(canvasW, o.Height).
			Cursor(ui.CursorPointer)

		px, _, over := stage.PointerPosition()
		bounds := stage.LayoutBounds()
		nextHover := -1
		if over && bounds.W > 0 && len(cols) > 0 {
			padX := float32(6)
			slotW := (bounds.W - padX*2) / float32(len(cols))
			idx := int((px - padX) / slotW)
			if idx >= 0 && idx < len(cols) {
				nextHover = idx
			}
		}
		*hoverPtr = nextHover

		sp := ui.SpringOptions{Duration: 240 * time.Millisecond, Bounce: 0.2}
		hl := make([]float32, len(cols))
		for i := range cols {
			target := float32(0)
			if nextHover == i {
				target = 1
			}
			hl[i] = stage.Spring(fmt.Sprintf("stk_hl_%d", i), target, sp)
		}
		anyT := float32(0)
		if nextHover >= 0 {
			anyT = 1
		}
		anyH := stage.Spring("stk_any", anyT, ui.SpringOptions{Duration: 200 * time.Millisecond})

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 || len(cols) == 0 {
				return
			}
			padX := float32(6)
			topY := r.Y + 6
			baseY := r.Y + r.H - 18
			chartH := baseY - topY
			if chartH < 20 {
				chartH = 20
			}

			// Horizontal grid lines.
			for k := 0; k < 3; k++ {
				gy := topY + chartH*float32(k)/2.0
				p.Line(r.X+padX, gy, r.X+r.W-padX, gy, 1, t.Border.Alpha(0.55))
			}

			slotW := (r.W - padX*2) / float32(len(cols))
			baseBarW := clamp32(slotW*0.56, 12, 24)

			for i, col := range cols {
				h := hl[i]
				dim := 1 - anyH*0.68 + h*0.68
				cx := r.X + padX + slotW*(float32(i)+0.5)
				barW := baseBarW + h*2.5

				if h > 0.04 {
					bgRect := ui.Rect{X: cx - slotW*0.44, Y: topY - 2, W: slotW * 0.88, H: chartH + 4}
					p.Fill(bgRect, t.SurfaceSubtle.Alpha(0.85*h), 6)
				}

				segGap := float32(2)
				curY := baseY
				for si, val := range col.Values {
					if val <= 0 {
						continue
					}
					segH := float32(val/maxSum) * (chartH - segGap*float32(len(col.Values)-1))
					if segH < 3 {
						segH = 3
					}
					y0 := curY - segH
					rect := ui.Rect{X: cx - barW*0.5, Y: y0, W: barW, H: segH}
					sc := colors[si%len(colors)]
					p.Fill(rect, sc.Alpha(dim), 3)
					curY = y0 - segGap
				}

				lblCol := lerpColor(t.TextSubtle, t.Text, h)
				p.Text(cx-float32(len(col.Label))*2.6, baseY+5, col.Label, 9.5, lblCol)
			}
		})
	})

	return card
}

// ---------------------------------------------------------------------------
// 7. blast-radius ("Network graph")
// ---------------------------------------------------------------------------

// GraphNode represents a service node in BlastRadiusGraph.
type GraphNode struct {
	Label string
	X, Y  float32 // normalized 0..1
}

// GraphEdge represents a directed dependency edge From -> To.
type GraphEdge struct {
	From int
	To   int
}

// BlastRadiusOptions configures BlastRadiusGraph.
type BlastRadiusOptions struct {
	Width  float32
	Height float32
	Nodes  []GraphNode
	Edges  []GraphEdge
}

// BlastRadiusGraph renders an interactive service dependency graph on ui.Painter.
// Hovering any node illuminates its outgoing blast-radius edges and downstream nodes.
func BlastRadiusGraph(c *ui.Context, opts ...BlastRadiusOptions) ui.Element {
	var o BlastRadiusOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 292
	}
	if o.Height <= 0 {
		o.Height = 132
	}
	nodes := o.Nodes
	edges := o.Edges
	if len(nodes) == 0 {
		nodes = []GraphNode{
			{Label: "Auth", X: 0.14, Y: 0.28},
			{Label: "Edge", X: 0.14, Y: 0.74},
			{Label: "API", X: 0.44, Y: 0.48},
			{Label: "Cache", X: 0.48, Y: 0.16},
			{Label: "Orders", X: 0.74, Y: 0.24},
			{Label: "Billing", X: 0.74, Y: 0.56},
			{Label: "Notify", X: 0.86, Y: 0.84},
		}
		edges = []GraphEdge{
			{0, 2}, {0, 3}, {1, 2}, {3, 4},
			{2, 4}, {2, 5}, {4, 5}, {5, 6}, {1, 6},
		}
	}

	// Compute transitive downstream nodes & edges via BFS from start.
	blastFrom := func(start int) (nodeAff []bool, edgeAff []bool, count int) {
		nodeAff = make([]bool, len(nodes))
		edgeAff = make([]bool, len(edges))
		if start < 0 || start >= len(nodes) {
			return nodeAff, edgeAff, 0
		}
		queue := []int{start}
		visited := make([]bool, len(nodes))
		visited[start] = true
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for ei, e := range edges {
				if e.From == u && e.To >= 0 && e.To < len(nodes) {
					edgeAff[ei] = true
					if !visited[e.To] {
						visited[e.To] = true
						nodeAff[e.To] = true
						count++
						queue = append(queue, e.To)
					}
				}
			}
		}
		return nodeAff, edgeAff, count
	}

	canvasW := o.Width - 24
	if canvasW < 180 {
		canvasW = 180
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(10, 12).
		Gap(6).
		Radius(12).
		Background(t.Surface).
		Border(1, t.Border)

	hoverPtr := ui.Local(card, "hc_blast_hover", func() int { return -1 })

	card.Children(func() {
		active := *hoverPtr
		_, _, impacted := blastFrom(active)

		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Blast radius").FontSize(12).FontWeight(600).TextColor(t.Text)
			if active >= 0 && active < len(nodes) {
				status := fmt.Sprintf("%s -> %d impacted", nodes[active].Label, impacted)
				statusCol := t.Warning
				if impacted >= 3 {
					statusCol = t.Danger
				} else if impacted == 0 {
					status = fmt.Sprintf("%s · leaf node", nodes[active].Label)
					statusCol = t.Success
				}
				ui.Text(c, status).FontSize(11).FontWeight(600).TextColor(statusCol)
			} else {
				ui.Text(c, "Hover a service").FontSize(11).TextColor(t.TextSubtle)
			}
		})

		stage := ui.Box(c).
			Size(canvasW, o.Height).
			Cursor(ui.CursorPointer)

		px, py, over := stage.PointerPosition()
		bounds := stage.LayoutBounds()
		nextHover := -1
		if over && bounds.W > 0 && bounds.H > 0 {
			bestD := float32(28)
			for i, nd := range nodes {
				nx := nd.X * bounds.W
				ny := nd.Y * bounds.H
				d := float32(math.Hypot(float64(px-nx), float64(py-ny)))
				if d < bestD {
					bestD = d
					nextHover = i
				}
			}
		}
		*hoverPtr = nextHover

		nodeAff, edgeAff, _ := blastFrom(nextHover)
		sp := ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.22}
		nodeEnergy := make([]float32, len(nodes))
		for i := range nodes {
			target := float32(0)
			if i == nextHover {
				target = 1.0
			} else if nodeAff[i] {
				target = 0.7
			}
			nodeEnergy[i] = stage.Spring(fmt.Sprintf("bl_n_%d", i), target, sp)
		}
		anyT := float32(0)
		if nextHover >= 0 {
			anyT = 1
		}
		anyH := stage.Spring("bl_any", anyT, ui.SpringOptions{Duration: 200 * time.Millisecond})
		phase := phaseSeconds(c.Now())
		if nextHover >= 0 {
			c.AnimationFrame()
		}

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			if nextHover >= 0 {
				p.AnimationFrame()
			}
			nodePos := func(i int) (float32, float32) {
				return r.X + nodes[i].X*r.W, r.Y + nodes[i].Y*r.H
			}

			// Draw edges first.
			for ei, e := range edges {
				if e.From < 0 || e.From >= len(nodes) || e.To < 0 || e.To >= len(nodes) {
					continue
				}
				x0, y0 := nodePos(e.From)
				x1, y1 := nodePos(e.To)
				aff := edgeAff[ei]

				edgeCol := t.BorderStrong.Alpha(0.75 * (1 - anyH*0.55))
				w := float32(1.2)
				if aff {
					edgeCol = lerpColor(t.Warning, t.Danger, 0.45)
					w = 2.0
				}

				// Quadratic gentle arc for visual clarity.
				mx := (x0 + x1) * 0.5
				my := (y0+y1)*0.5 - 6
				var ep ui.Path
				ep.MoveTo(x0, y0)
				ep.QuadTo(mx, my, x1, y1)
				if aff {
					p.StrokePath(&ep, 5.0, edgeCol.Alpha(0.18))
				}
				p.StrokePath(&ep, w, edgeCol)

				// Small directional chevron near target.
				ang := math.Atan2(float64(y1-my), float64(x1-mx))
				tipX := x1 - float32(math.Cos(ang))*16
				tipY := y1 - float32(math.Sin(ang))*16
				wing := float32(4.2)
				p.Line(tipX, tipY, tipX-wing*float32(math.Cos(ang-0.55)), tipY-wing*float32(math.Sin(ang-0.55)), w, edgeCol)
				p.Line(tipX, tipY, tipX-wing*float32(math.Cos(ang+0.55)), tipY-wing*float32(math.Sin(ang+0.55)), w, edgeCol)

				// Animated blast-wave pulse dot traveling along affected edges.
				if aff {
					u := float32(mod1(phase*1.1 + float64(ei)*0.19))
					om := 1 - u
					qx := om*om*x0 + 2*om*u*mx + u*u*x1
					qy := om*om*y0 + 2*om*u*my + u*u*y1
					p.Fill(ui.Rect{X: qx - 2.6, Y: qy - 2.6, W: 5.2, H: 5.2}, t.Danger, 2.6)
				}
			}

			// Draw service nodes as sleek rounded pills.
			for i, nd := range nodes {
				cx, cy := nodePos(i)
				en := nodeEnergy[i]
				isSrc := i == nextHover
				isAff := nodeAff[i]
				dim := 1 - anyH*0.58 + clamp32(en*0.85, 0, 0.58)

				pillW := float32(len(nd.Label))*6.0 + 16 + en*2
				pillH := float32(22) + en*2
				rect := ui.Rect{X: cx - pillW*0.5, Y: cy - pillH*0.5, W: pillW, H: pillH}

				if en > 0.08 {
					glowCol := t.Warning
					if isSrc {
						glowCol = t.Danger
					}
					glowR := pillW * 0.65
					p.FillRadialGradient(
						ui.Rect{X: cx - glowR, Y: cy - glowR*0.6, W: glowR * 2, H: glowR * 1.2},
						ui.RadialGradient{
							Center: ui.Point{X: 0.5, Y: 0.5},
							Radius: glowR,
							Stops: []ui.GradientStop{
								{Color: glowCol.Alpha(0.32 * en), At: 0},
								{Color: ui.Transparent, At: 1},
							},
						},
						glowR,
					)
				}

				bg := t.SurfaceSubtle
				border := t.BorderStrong
				txt := t.Text
				if isSrc {
					bg = t.Danger.Alpha(0.18)
					border = t.Danger
				} else if isAff {
					bg = t.Warning.Alpha(0.16)
					border = t.Warning
				}

				p.Fill(rect, t.Surface, 11)
				p.Fill(rect, bg.Alpha(dim), 11)
				p.Stroke(rect, border.Alpha(dim), 11, 1.1+0.5*en)
				p.Text(cx-float32(len(nd.Label))*2.7, cy-5.2, nd.Label, 9.5, txt.Alpha(dim))
			}
		})
	})

	return card
}

// ---------------------------------------------------------------------------
// 8. hub-diagram ("Hub diagram")
// ---------------------------------------------------------------------------

// HubSpoke describes one surrounding spoke node in HubDiagram.
type HubSpoke struct {
	Label      string
	LatencyMS  int
	Throughput string
	Color      ui.Color
}

// HubDiagramOptions configures HubDiagram.
type HubDiagramOptions struct {
	Width    float32
	Height   float32
	HubLabel string
	Spokes   []HubSpoke
}

// HubDiagram renders a central hub with 6 surrounding spoke nodes connected by
// animated traveling signal pulses. Hovering a spoke intensifies its beam and telemetry.
func HubDiagram(c *ui.Context, opts ...HubDiagramOptions) ui.Element {
	var o HubDiagramOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 292
	}
	if o.Height <= 0 {
		o.Height = 136
	}
	if o.HubLabel == "" {
		o.HubLabel = "Core"
	}
	palette := hcSeriesColors(t)
	spokes := o.Spokes
	if len(spokes) == 0 {
		spokes = []HubSpoke{
			{Label: "US-East", LatencyMS: 12, Throughput: "4.2k/s", Color: palette[0]},
			{Label: "EU-West", LatencyMS: 38, Throughput: "2.8k/s", Color: palette[1]},
			{Label: "AP-Tok", LatencyMS: 84, Throughput: "1.9k/s", Color: palette[4]},
			{Label: "Cache", LatencyMS: 2, Throughput: "14k/s", Color: palette[2]},
			{Label: "Queue", LatencyMS: 7, Throughput: "6.1k/s", Color: palette[3]},
			{Label: "Search", LatencyMS: 19, Throughput: "3.4k/s", Color: palette[5]},
		}
	}

	canvasW := o.Width - 24
	if canvasW < 180 {
		canvasW = 180
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(10, 12).
		Gap(6).
		Radius(12).
		Background(t.Surface).
		Border(1, t.Border)

	hoverPtr := ui.Local(card, "hc_hub_hover", func() int { return -1 })

	spokePos := func(w, h float32, i int) (float32, float32) {
		cx := w * 0.5
		cy := h * 0.5
		rx := w * 0.37
		ry := h * 0.38
		ang := -math.Pi/2 + float64(i)*(math.Pi*2/float64(len(spokes)))
		return cx + rx*float32(math.Cos(ang)), cy + ry*float32(math.Sin(ang))
	}

	card.Children(func() {
		active := *hoverPtr

		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Hub mesh").FontSize(12).FontWeight(600).TextColor(t.Text)
			if active >= 0 && active < len(spokes) {
				sp := spokes[active]
				col := colorOr(sp.Color, palette[active%len(palette)])
				ui.Text(c, fmt.Sprintf("%s · %dms · %s", sp.Label, sp.LatencyMS, sp.Throughput)).
					FontSize(10.5).FontWeight(600).TextColor(col)
			} else {
				ui.Text(c, fmt.Sprintf("%d spokes online", len(spokes))).FontSize(10.5).TextColor(t.TextSubtle)
			}
		})

		stage := ui.Box(c).
			Size(canvasW, o.Height).
			Cursor(ui.CursorPointer)

		px, py, over := stage.PointerPosition()
		bounds := stage.LayoutBounds()
		nextHover := -1
		if over && bounds.W > 0 && bounds.H > 0 {
			bestD := float32(34)
			for i := range spokes {
				sx, sy := spokePos(bounds.W, bounds.H, i)
				d := float32(math.Hypot(float64(px-sx), float64(py-sy)))
				if d < bestD {
					bestD = d
					nextHover = i
				}
			}
		}
		*hoverPtr = nextHover

		spOpt := ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.25}
		beam := make([]float32, len(spokes))
		for i := range spokes {
			target := float32(0)
			if nextHover == i {
				target = 1
			}
			beam[i] = stage.Spring(fmt.Sprintf("hub_b_%d", i), target, spOpt)
		}
		anyT := float32(0)
		if nextHover >= 0 {
			anyT = 1
		}
		anyH := stage.Spring("hub_any", anyT, ui.SpringOptions{Duration: 220 * time.Millisecond})

		phase := phaseSeconds(c.Now())
		c.AnimationFrame()

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.AnimationFrame()

			cx := r.X + r.W*0.5
			cy := r.Y + r.H*0.5

			// Subtle orbital guide ring.
			rx := r.W * 0.37
			ry := r.H * 0.38
			var orbit ui.Path
			for k := 0; k <= 48; k++ {
				a := float64(k) * (math.Pi * 2 / 48.0)
				ox := cx + rx*float32(math.Cos(a))
				oy := cy + ry*float32(math.Sin(a))
				if k == 0 {
					orbit.MoveTo(ox, oy)
				} else {
					orbit.LineTo(ox, oy)
				}
			}
			orbit.Close()
			p.StrokePath(&orbit, 1, t.Border.Alpha(0.45))

			// Beams and traveling signal pulses.
			for i, sp := range spokes {
				lsx, lsy := spokePos(r.W, r.H, i)
				sx, sy := r.X+lsx, r.Y+lsy
				b := beam[i]
				dim := 1 - anyH*0.52 + b*0.52
				col := colorOr(sp.Color, palette[i%len(palette)])

				if b > 0.05 {
					p.Line(cx, cy, sx, sy, 5.5, col.Alpha(0.22*b))
				}
				p.Line(cx, cy, sx, sy, 1.2+1.4*b, lerpColor(t.BorderStrong, col, 0.35+0.65*b).Alpha(dim))

				// Traveling signal packets (outgoing & returning).
				speed := 0.55 + float64(b)*0.65
				for k := 0; k < 2; k++ {
					u := float32(mod1(phase*speed + float64(i)*0.17 + float64(k)*0.5))
					if k == 1 {
						u = 1 - u
					}
					pxPos := lerp32(cx, sx, u)
					pyPos := lerp32(cy, sy, u)
					pr := float32(2.1 + 1.5*b)
					if b > 0.1 {
						hr := pr * 3.0
						p.FillRadialGradient(
							ui.Rect{X: pxPos - hr, Y: pyPos - hr, W: hr * 2, H: hr * 2},
							ui.RadialGradient{
								Center: ui.Point{X: 0.5, Y: 0.5},
								Radius: hr,
								Stops: []ui.GradientStop{
									{Color: col.Alpha(0.55 * b), At: 0},
									{Color: ui.Transparent, At: 1},
								},
							},
							hr,
						)
					}
					p.Fill(ui.Rect{X: pxPos - pr, Y: pyPos - pr, W: pr * 2, H: pr * 2}, col.Alpha(dim), pr)
				}
			}

			// Central hub with breathing halo.
			pulse := float32(0.5 + 0.5*math.Sin(phase*2.8))
			hubGlowR := float32(28 + 5*pulse)
			p.FillRadialGradient(
				ui.Rect{X: cx - hubGlowR, Y: cy - hubGlowR, W: hubGlowR * 2, H: hubGlowR * 2},
				ui.RadialGradient{
					Center: ui.Point{X: 0.5, Y: 0.5},
					Radius: hubGlowR,
					Stops: []ui.GradientStop{
						{Color: t.Accent.Alpha(0.28 + 0.12*pulse), At: 0},
						{Color: ui.Transparent, At: 1},
					},
				},
				hubGlowR,
			)
			hubR := float32(17)
			hubRect := ui.Rect{X: cx - hubR, Y: cy - hubR, W: hubR * 2, H: hubR * 2}
			p.Fill(hubRect, t.Accent, hubR)
			p.Stroke(hubRect, ui.Hex("#ffffff").Alpha(0.35), hubR, 1.2)
			p.Text(cx-float32(len(o.HubLabel))*2.9, cy-5.2, o.HubLabel, 9.5, t.AccentText)

			// Surrounding spoke badges.
			for i, sp := range spokes {
				lsx, lsy := spokePos(r.W, r.H, i)
				sx, sy := r.X+lsx, r.Y+lsy
				b := beam[i]
				dim := 1 - anyH*0.48 + b*0.48
				col := colorOr(sp.Color, palette[i%len(palette)])

				pw := float32(len(sp.Label))*5.6 + 14 + b*4
				ph := float32(20) + b*2.5
				rect := ui.Rect{X: sx - pw*0.5, Y: sy - ph*0.5, W: pw, H: ph}
				p.Fill(rect, t.Surface, 10)
				p.Fill(rect, col.Alpha((0.08+0.16*b)*dim), 10)
				p.Stroke(rect, lerpColor(t.BorderStrong, col, 0.3+0.7*b).Alpha(dim), 10, 1.0+0.6*b)
				p.Text(sx-float32(len(sp.Label))*2.6, sy-5.0, sp.Label, 9, t.Text.Alpha(dim))
			}
		})
	})

	return card
}

// ---------------------------------------------------------------------------
// 9. heart-rate ("Heart rate")
// ---------------------------------------------------------------------------

// HeartRateOptions configures HeartRateMonitor.
type HeartRateOptions struct {
	Width    float32
	Height   float32
	BaseBPM  float64
	HoverBPM float64
	Accent   ui.Color
}

type hcHeartState struct {
	beats     float64
	lastFrame time.Time
}

func hcECGWave(cycle float64) float32 {
	// cycle in [0, 1): realistic P-Q-R-S-T Gaussian composite waveform.
	u := mod1(cycle)
	gauss := func(mu, sig, amp float64) float64 {
		d := (u - mu) / sig
		return amp * math.Exp(-0.5*d*d)
	}
	v := gauss(0.18, 0.032, 0.14) + // P wave
		gauss(0.31, 0.012, -0.16) + // Q dip
		gauss(0.35, 0.014, 1.00) + // R spike
		gauss(0.39, 0.014, -0.28) + // S dip
		gauss(0.58, 0.048, 0.24) // T wave
	return float32(v)
}

// HeartRateMonitor renders a live scrolling ECG heartbeat trace with a pulsing
// glowing lead dot and a BPM counter that accelerates smoothly on hover.
func HeartRateMonitor(c *ui.Context, opts ...HeartRateOptions) ui.Element {
	var o HeartRateOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 292
	}
	if o.Height <= 0 {
		o.Height = 118
	}
	if o.BaseBPM <= 0 {
		o.BaseBPM = 72
	}
	if o.HoverBPM <= 0 {
		o.HoverBPM = 108
	}
	accent := colorOr(o.Accent, ui.Hex("#f43f5e"))

	canvasW := o.Width - 24
	if canvasW < 180 {
		canvasW = 180
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(10, 12).
		Gap(6).
		Radius(12).
		Background(t.Surface).
		Border(1, t.Border).
		Cursor(ui.CursorPointer)

	hovered := card.Hovered()
	_, _, over := card.PointerPosition()
	activeTarget := float32(0)
	if hovered || over {
		activeTarget = 1
	}
	energy := card.Spring("hr_energy", activeTarget, ui.SpringOptions{
		Duration: 420 * time.Millisecond,
		Bounce:   0.15,
	})
	liveBPM := o.BaseBPM + (o.HoverBPM-o.BaseBPM)*float64(energy)

	st := ui.Local(card, "hc_hr_phase", func() hcHeartState {
		return hcHeartState{lastFrame: c.Now()}
	})
	now := c.Now()
	dt := now.Sub(st.lastFrame).Seconds()
	if dt < 0 || dt > 0.25 {
		dt = 1.0 / 60.0
	}
	st.lastFrame = now
	st.beats += dt * (liveBPM / 60.0)
	beats := st.beats
	c.AnimationFrame()

	leadVal := hcECGWave(beats)
	peakPulse := clamp32(leadVal, 0, 1)

	card.Children(func() {
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
				dotBox := ui.Box(c).Size(12, 12)
				dotBox.Draw(func(p *ui.Painter, r ui.Rect) {
					cx := r.X + r.W*0.5
					cy := r.Y + r.H*0.5
					hr := float32(3.2 + 2.2*peakPulse)
					p.Fill(ui.Rect{X: cx - hr - 2, Y: cy - hr - 2, W: (hr + 2) * 2, H: (hr + 2) * 2}, accent.Alpha(0.24+0.25*peakPulse), hr+2)
					p.Fill(ui.Rect{X: cx - hr, Y: cy - hr, W: hr * 2, H: hr * 2}, accent, hr)
				})
				ui.Text(c, "Heart rate").FontSize(12).FontWeight(600).TextColor(t.Text)
			})
			ui.Row(c).Gap(5).AlignItems(ui.Center).Children(func() {
				tag := "REST"
				if energy > 0.35 {
					tag = "ACTIVE"
				}
				ui.Text(c, tag).FontSize(9.5).FontWeight(700).TextColor(accent)
				ui.Text(c, fmt.Sprintf("%d BPM", int(liveBPM+0.5))).
					FontSize(13).FontWeight(700).TextColor(t.Text)
			})
		})

		stage := ui.Box(c).
			Size(canvasW, o.Height).
			Radius(8).
			Background(t.SurfaceSubtle).
			Border(1, t.Border).
			Clip()

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.AnimationFrame()

			// Subtle ECG millimeter grid.
			gridCol := t.Border.Alpha(0.55)
			step := float32(16)
			for x := r.X + step; x < r.X+r.W; x += step {
				p.Line(x, r.Y, x, r.Y+r.H, 0.75, gridCol)
			}
			for y := r.Y + step; y < r.Y+r.H; y += step {
				p.Line(r.X, y, r.X+r.W, y, 0.75, gridCol)
			}

			padL := float32(8)
			padR := float32(16)
			traceW := r.W - padL - padR
			baseY := r.Y + r.H*0.64
			amp := r.H * 0.44
			samples := 110
			windowBeats := 2.35

			// Draw segmented trailing phosphor trace so older tail fades cleanly.
			const chunks = 12
			ptsPerChunk := samples / chunks
			var leadX, leadY float32
			for ch := 0; ch < chunks; ch++ {
				startIdx := ch * ptsPerChunk
				endIdx := startIdx + ptsPerChunk
				if ch == chunks-1 {
					endIdx = samples
				}
				poly := make([][2]float32, 0, endIdx-startIdx+1)
				for i := startIdx; i <= endIdx; i++ {
					u := float32(i) / float32(samples)
					x := r.X + padL + u*traceW
					beatAt := beats - float64(1-u)*windowBeats
					y := baseY - hcECGWave(beatAt)*amp
					poly = append(poly, [2]float32{x, y})
					if i == samples {
						leadX, leadY = x, y
					}
				}
				fade := float32(ch+1) / float32(chunks)
				fade = fade * fade
				p.StrokePolyline(poly, 4.5, accent.Alpha(0.16*fade))
				p.StrokePolyline(poly, 1.9, accent.Alpha(0.14+0.86*fade))
			}

			// Pulsing glowing lead dot.
			haloR := float32(11 + 8*peakPulse + 3*energy)
			p.FillRadialGradient(
				ui.Rect{X: leadX - haloR, Y: leadY - haloR, W: haloR * 2, H: haloR * 2},
				ui.RadialGradient{
					Center: ui.Point{X: 0.5, Y: 0.5},
					Radius: haloR,
					Stops: []ui.GradientStop{
						{Color: accent.Alpha(0.55 + 0.35*peakPulse), At: 0},
						{Color: ui.Transparent, At: 1},
					},
				},
				haloR,
			)
			dotR := float32(3.4 + 1.0*peakPulse)
			p.Fill(ui.Rect{X: leadX - dotR, Y: leadY - dotR, W: dotR * 2, H: dotR * 2}, accent, dotR)
			p.Fill(ui.Rect{X: leadX - 1.5, Y: leadY - 1.5, W: 3, H: 3}, ui.Hex("#ffffff"), 1.5)
		})
	})

	return card
}

// ---------------------------------------------------------------------------
// 10. progress ("Progress ticks")
// ---------------------------------------------------------------------------

// ProgressTicksOptions configures ProgressTicks.
type ProgressTicksOptions struct {
	Width   float32
	Height  float32
	Count   int
	Initial float32
	Label   string
	Accent  ui.Color
}

// ProgressTicks renders a row of 28 fine vertical progress ticks that ripple in
// height near the pointer with Gaussian falloff and display a scrubbed percentage readout.
func ProgressTicks(c *ui.Context, opts ...ProgressTicksOptions) ui.Element {
	var o ProgressTicksOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 292
	}
	if o.Height <= 0 {
		o.Height = 90
	}
	if o.Count < 12 {
		o.Count = 28
	}
	if o.Count > 40 {
		o.Count = 40
	}
	if o.Initial <= 0 || o.Initial > 1 {
		o.Initial = 0.64
	}
	if o.Label == "" {
		o.Label = "Sync progress"
	}
	accent := colorOr(o.Accent, t.Accent)

	canvasW := o.Width - 24
	if canvasW < 180 {
		canvasW = 180
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(12, 12).
		Gap(8).
		Radius(12).
		Background(t.Surface).
		Border(1, t.Border)

	type tickState struct {
		saved   float32
		scrub   float32
		hovered bool
	}
	st := ui.Local(card, "hc_ticks_st", func() tickState {
		return tickState{saved: o.Initial, scrub: o.Initial}
	})

	card.Children(func() {
		displayVal := st.saved
		statusLbl := o.Label
		pctCol := t.Text
		if st.hovered {
			displayVal = st.scrub
			statusLbl = "Scrubbing · click to set"
			pctCol = accent
		}

		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, statusLbl).FontSize(11.5).FontWeight(600).TextColor(t.TextMuted)
			ui.Text(c, fmt.Sprintf("%d%%", int(displayVal*100+0.5))).
				FontSize(14).FontWeight(700).TextColor(pctCol)
		})

		stage := ui.Box(c).
			Size(canvasW, o.Height).
			Cursor(ui.CursorPointer)

		px, _, over := stage.PointerPosition()
		bounds := stage.LayoutBounds()
		padX := float32(6)
		if over && bounds.W > padX*2 {
			st.hovered = true
			st.scrub = clamp32((px-padX)/(bounds.W-padX*2), 0, 1)
			if stage.Clicked() {
				st.saved = st.scrub
			}
		} else {
			st.hovered = false
		}

		activeVal := st.saved
		if st.hovered {
			activeVal = st.scrub
		}
		sp := ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.25}
		heights := make([]float32, o.Count)
		for i := 0; i < o.Count; i++ {
			u := float32(i) / float32(o.Count-1)
			targetH := float32(18)
			if u <= activeVal+0.001 {
				targetH = 24
			}
			if st.hovered {
				dx := (u - st.scrub) * float32(o.Count)
				gauss := float32(math.Exp(float64(-dx * dx / 7.5)))
				targetH += gauss * 34
			}
			heights[i] = stage.Spring(fmt.Sprintf("pt_h_%d", i), targetH, sp)
		}

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			usableW := r.W - padX*2
			baseY := r.Y + r.H - 12
			tickW := float32(3.2)

			for i := 0; i < o.Count; i++ {
				u := float32(i) / float32(o.Count-1)
				cx := r.X + padX + u*usableW
				h := heights[i]
				filled := u <= activeVal+0.001

				var col ui.Color
				if filled {
					col = lerpColor(accent.Alpha(0.65), accent, clamp32((h-18)/36, 0, 1))
				} else {
					near := float32(0)
					if st.hovered {
						dx := (u - st.scrub) * float32(o.Count)
						near = float32(math.Exp(float64(-dx * dx / 9.0)))
					}
					col = lerpColor(t.BorderStrong.Alpha(0.75), accent.Alpha(0.45), near)
				}

				rect := ui.Rect{X: cx - tickW*0.5, Y: baseY - h, W: tickW, H: h}
				p.Fill(rect, col, tickW*0.5)
			}

			// Baseline indicator dot under active position.
			actX := r.X + padX + clamp32(activeVal, 0, 1)*usableW
			p.Fill(ui.Rect{X: actX - 2.5, Y: baseY + 4, W: 5, H: 5}, accent, 2.5)
		})
	})

	return card
}

// ---------------------------------------------------------------------------
// 11. heat-word ("Heat map")
// ---------------------------------------------------------------------------

// HeatMapOptions configures HeatMapGrid.
type HeatMapOptions struct {
	Width  float32
	Cols   int
	Rows   int
	Title  string
	Values []int
}

// HeatMapGrid renders a 12x5 contribution/activity matrix grid where hovering a
// cell illuminates its row/column neighborhood with a thermal glow and readout.
func HeatMapGrid(c *ui.Context, opts ...HeatMapOptions) ui.Element {
	var o HeatMapOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 292
	}
	if o.Cols <= 0 {
		o.Cols = 12
	}
	if o.Rows <= 0 {
		o.Rows = 5
	}
	if o.Title == "" {
		o.Title = "Activity heat"
	}

	cellCount := o.Cols * o.Rows
	vals := o.Values
	if len(vals) != cellCount {
		vals = make([]int, cellCount)
		// Deterministic rich activity matrix with realistic hotspots.
		for r := 0; r < o.Rows; r++ {
			for col := 0; col < o.Cols; col++ {
				h1 := math.Exp(-float64((col-3)*(col-3)+(r-2)*(r-2)) / 6.5)
				h2 := math.Exp(-float64((col-8)*(col-8)+(r-1)*(r-1)) / 5.0)
				h3 := math.Exp(-float64((col-10)*(col-10)+(r-3)*(r-3)) / 4.2)
				wave := 0.5 + 0.5*math.Sin(float64(col*5+r*7))
				v := int((h1*11 + h2*14 + h3*9 + wave*3.2) + 0.3)
				if (col*7+r*11)%9 == 0 {
					v = 0
				}
				vals[r*o.Cols+col] = v
			}
		}
	}
	maxV := 1
	totalEvents := 0
	for _, v := range vals {
		if v > maxV {
			maxV = v
		}
		totalEvents += v
	}

	thermalColor := func(norm float32) ui.Color {
		norm = clamp32(norm, 0, 1)
		if norm <= 0.01 {
			return t.SurfaceSubtle
		}
		if t.Dark {
			stops := [4]ui.Color{
				ui.Hex("#1e293b"),
				ui.Hex("#0ea5e9"),
				ui.Hex("#f59e0b"),
				ui.Hex("#f43f5e"),
			}
			scaled := norm * 3
			idx := int(scaled)
			if idx >= 3 {
				return stops[3]
			}
			return lerpColor(stops[idx], stops[idx+1], scaled-float32(idx))
		}
		stops := [4]ui.Color{
			ui.Hex("#e0f2fe"),
			ui.Hex("#38bdf8"),
			ui.Hex("#f59e0b"),
			ui.Hex("#e11d48"),
		}
		scaled := norm * 3
		idx := int(scaled)
		if idx >= 3 {
			return stops[3]
		}
		return lerpColor(stops[idx], stops[idx+1], scaled-float32(idx))
	}

	dayNames := [5]string{"Mon", "Tue", "Wed", "Thu", "Fri"}
	canvasW := o.Width - 24
	if canvasW < 180 {
		canvasW = 180
	}
	canvasH := float32(104)

	card := ui.Column(c).
		Width(o.Width).
		Padding(10, 12).
		Gap(6).
		Radius(12).
		Background(t.Surface).
		Border(1, t.Border)

	type heatHover struct {
		col, row int
	}
	hoverPtr := ui.Local(card, "hc_heat_hover", func() heatHover { return heatHover{-1, -1} })

	card.Children(func() {
		hh := *hoverPtr
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, o.Title).FontSize(12).FontWeight(600).TextColor(t.Text)
			if hh.col >= 0 && hh.row >= 0 {
				v := vals[hh.row*o.Cols+hh.col]
				day := dayNames[hh.row%len(dayNames)]
				col := thermalColor(float32(v) / float32(maxV))
				if v == 0 {
					col = t.TextMuted
				}
				ui.Text(c, fmt.Sprintf("W%02d · %s · %d ops", hh.col+1, day, v)).
					FontSize(10.5).FontWeight(600).TextColor(col)
			} else {
				ui.Text(c, fmt.Sprintf("%d ops · %dw", totalEvents, o.Cols)).
					FontSize(10.5).TextColor(t.TextSubtle)
			}
		})

		stage := ui.Box(c).
			Size(canvasW, canvasH).
			Cursor(ui.CursorPointer)

		px, py, over := stage.PointerPosition()
		bounds := stage.LayoutBounds()
		gap := float32(4)
		nextH := heatHover{-1, -1}
		if over && bounds.W > 0 && bounds.H > 0 {
			cellW := (bounds.W - gap*float32(o.Cols-1)) / float32(o.Cols)
			cellH := (bounds.H - gap*float32(o.Rows-1)) / float32(o.Rows)
			cc := int(px / (cellW + gap))
			rr := int(py / (cellH + gap))
			if cc >= 0 && cc < o.Cols && rr >= 0 && rr < o.Rows {
				nextH = heatHover{cc, rr}
			}
		}
		*hoverPtr = nextH

		anyT := float32(0)
		if nextH.col >= 0 {
			anyT = 1
		}
		anyH := stage.Spring("hm_any", anyT, ui.SpringOptions{Duration: 200 * time.Millisecond})
		focusX := stage.Spring("hm_fx", float32(nextH.col), ui.SpringOptions{Duration: 180 * time.Millisecond, Bounce: 0.15})
		focusY := stage.Spring("hm_fy", float32(nextH.row), ui.SpringOptions{Duration: 180 * time.Millisecond, Bounce: 0.15})

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			cellW := (r.W - gap*float32(o.Cols-1)) / float32(o.Cols)
			cellH := (r.H - gap*float32(o.Rows-1)) / float32(o.Rows)

			for row := 0; row < o.Rows; row++ {
				for col := 0; col < o.Cols; col++ {
					x := r.X + float32(col)*(cellW+gap)
					y := r.Y + float32(row)*(cellH+gap)
					v := vals[row*o.Cols+col]
					norm := float32(v) / float32(maxV)
					baseCol := thermalColor(norm)

					isExact := nextH.col == col && nextH.row == row
					sameAxis := nextH.col == col || nextH.row == row

					grow := float32(0)
					alpha := float32(1)
					drawCol := baseCol
					if anyH > 0.02 {
						if isExact {
							grow = 1.6 * anyH
						} else if sameAxis {
							dist := float32(math.Abs(float64(col-nextH.col)) + math.Abs(float64(row-nextH.row)))
							glowAmt := clamp32(1.0/(1.0+dist*0.32), 0.18, 0.75) * anyH
							drawCol = lerpColor(baseCol, ui.Hex("#fb923c"), glowAmt*0.42)
						} else {
							alpha = 1 - anyH*0.52
						}
					}

					rect := ui.Rect{
						X: x - grow,
						Y: y - grow,
						W: cellW + grow*2,
						H: cellH + grow*2,
					}
					p.Fill(rect, drawCol.Alpha(alpha), 4)
					if v == 0 {
						p.Stroke(rect, t.Border.Alpha(0.75*alpha), 4, 0.75)
					}
					if isExact {
						p.Stroke(rect, t.Text.Alpha(0.85), 4, 1.3)
					}
				}
			}

			// Ambient thermal halo centered on hovered cell.
			if anyH > 0.04 && nextH.col >= 0 {
				hx := r.X + (focusX+0.5)*cellW + focusX*gap
				hy := r.Y + (focusY+0.5)*cellH + focusY*gap
				rad := float32(36)
				p.FillRadialGradient(
					ui.Rect{X: hx - rad, Y: hy - rad, W: rad * 2, H: rad * 2},
					ui.RadialGradient{
						Center: ui.Point{X: 0.5, Y: 0.5},
						Radius: rad,
						Stops: []ui.GradientStop{
							{Color: ui.Hex("#fb7185").Alpha(0.24 * anyH), At: 0},
							{Color: ui.Transparent, At: 1},
						},
					},
					rad,
				)
			}
		})

		// Compact thermal scale footer.
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Mon–Fri · 12 weeks").FontSize(9.5).TextColor(t.TextSubtle)
			ui.Row(c).Gap(3).AlignItems(ui.Center).Children(func() {
				ui.Text(c, "Less").FontSize(9).TextColor(t.TextSubtle)
				for k := 0; k < 5; k++ {
					ui.Box(c).Size(8, 8).Radius(2).Background(thermalColor(float32(k) / 4.0))
				}
				ui.Text(c, "More").FontSize(9).TextColor(t.TextSubtle)
			})
		})
	})

	return card
}

// ---------------------------------------------------------------------------
// 12. action-node ("Action node")
// ---------------------------------------------------------------------------

// ActionNodeOptions configures ActionNodeCard.
type ActionNodeOptions struct {
	Width    float32
	Title    string
	Subtitle string
	Kind     string
}

type hcActionNodeState struct {
	runningUntil time.Time
	runCount     int
	inspecting   bool
	disabled     bool
}

// ActionNodeCard renders a workflow node card with input/output ports and an
// animated status ring that reveals inline action pills (Run, Inspect, Disable) on hover.
func ActionNodeCard(c *ui.Context, opts ...ActionNodeOptions) ui.Element {
	var o ActionNodeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 296
	}
	if o.Title == "" {
		o.Title = "Transform JSON"
	}
	if o.Subtitle == "" {
		o.Subtitle = "v2.4 · 14ms avg"
	}
	if o.Kind == "" {
		o.Kind = "NODE"
	}

	stage := ui.Box(c).
		Width(o.Width).
		Padding(14, 16).
		Center()

	st := ui.Local(stage, "hc_action_node", func() hcActionNodeState {
		return hcActionNodeState{}
	})
	now := c.Now()
	running := now.Before(st.runningUntil)
	if running {
		c.AnimationFrame()
	}

	hovered := stage.Hovered()
	_, _, over := stage.PointerPosition()
	showActions := hovered || over || st.inspecting || running
	revealTarget := float32(0)
	if showActions {
		revealTarget = 1
	}
	reveal := stage.Spring("an_reveal", revealTarget, ui.SpringOptions{
		Duration: 260 * time.Millisecond,
		Bounce:   0.22,
	})

	statusCol := t.Success
	statusTxt := "Ready"
	if st.disabled {
		statusCol = t.TextSubtle
		statusTxt = "Disabled"
	} else if running {
		statusCol = t.Accent
		statusTxt = "Running"
	} else if st.inspecting {
		statusCol = t.Warning
		statusTxt = "Inspect"
	}
	phase := phaseSeconds(now)

	// Draw incoming/outgoing workflow wire stubs + port rings around the card.
	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		if running {
			p.AnimationFrame()
		}
		midY := r.Y + r.H*0.5
		leftX := r.X + 16
		rightX := r.X + r.W - 16
		wireCol := lerpColor(t.BorderStrong, statusCol, 0.35+0.45*reveal)
		if st.disabled {
			wireCol = t.Border.Alpha(0.6)
		}

		// Input/output wire stubs.
		p.Line(r.X+2, midY, leftX, midY, 1.6, wireCol)
		p.Line(rightX, midY, r.X+r.W-2, midY, 1.6, wireCol)

		// Traveling packet pulse on wires while running.
		if running && !st.disabled {
			u := float32(mod1(phase * 2.2))
			ix := lerp32(r.X+2, leftX, u)
			ox := lerp32(rightX, r.X+r.W-2, u)
			p.Fill(ui.Rect{X: ix - 2.2, Y: midY - 2.2, W: 4.4, H: 4.4}, t.Accent, 2.2)
			p.Fill(ui.Rect{X: ox - 2.2, Y: midY - 2.2, W: 4.4, H: 4.4}, t.Accent, 2.2)
		}

		// Left input port & right output port.
		for _, px := range [2]float32{leftX, rightX} {
			pr := float32(4.5 + 1.0*reveal)
			if reveal > 0.1 && !st.disabled {
				p.Fill(ui.Rect{X: px - pr - 3, Y: midY - pr - 3, W: (pr + 3) * 2, H: (pr + 3) * 2}, statusCol.Alpha(0.22*reveal), pr+3)
			}
			p.Fill(ui.Rect{X: px - pr, Y: midY - pr, W: pr * 2, H: pr * 2}, t.Surface, pr)
			p.Stroke(ui.Rect{X: px - pr, Y: midY - pr, W: pr * 2, H: pr * 2}, wireCol, pr, 1.6)
		}
	})

	stage.Children(func() {
		borderCol := lerpColor(t.Border, statusCol.Alpha(0.65), reveal*0.85)
		cardOpacity := float32(1)
		if st.disabled {
			cardOpacity = 0.62
		}

		nodeCard := ui.Column(c).
			Width(o.Width-32).
			Padding(11, 13).
			Gap(8).
			Radius(12).
			Background(t.Surface).
			Border(1.2, borderCol).
			Opacity(cardOpacity)

		nodeCard.Children(func() {
			// Node header: status ring + title/subtitle + status pill.
			ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
				ui.Row(c).Gap(9).AlignItems(ui.Center).Children(func() {
					// Status ring drawn with ui.Painter.
					ring := ui.Box(c).Size(24, 24)
					ring.Draw(func(p *ui.Painter, r ui.Rect) {
						cx := r.X + r.W*0.5
						cy := r.Y + r.H*0.5
						rad := float32(9.5)
						rect := ui.Rect{X: cx - rad, Y: cy - rad, W: rad * 2, H: rad * 2}
						p.Stroke(rect, t.Border, rad, 2.0)
						if running {
							p.FillConicGradient(
								rect,
								ui.ConicGradient{
									Center: ui.Point{X: 0.5, Y: 0.5},
									Angle:  float32(mod1(phase*1.4) * 360),
									Stops: []ui.GradientStop{
										{Color: t.Accent, At: 0},
										{Color: t.Accent.Alpha(0.15), At: 0.55},
										{Color: ui.Transparent, At: 1},
									},
								},
								rad,
							)
							inner := float32(6.8)
							p.Fill(ui.Rect{X: cx - inner, Y: cy - inner, W: inner * 2, H: inner * 2}, t.Surface, inner)
						} else {
							p.Stroke(rect, statusCol.Alpha(0.85), rad, 2.0)
						}
						dot := float32(3.2)
						p.Fill(ui.Rect{X: cx - dot, Y: cy - dot, W: dot * 2, H: dot * 2}, statusCol, dot)
					})

					ui.Column(c).Gap(1).Children(func() {
						ui.Text(c, o.Title).FontSize(12.5).FontWeight(600).TextColor(t.Text)
						sub := o.Subtitle
						if st.runCount > 0 && !running && !st.disabled {
							sub = fmt.Sprintf("Ran %dx · 12ms", st.runCount)
						}
						ui.Text(c, sub).FontSize(10).TextColor(t.TextSubtle)
					})
				})

				ui.Box(c).
					Padding(2, 7).
					Radius(99).
					Background(statusCol.Alpha(0.14)).
					Border(1, statusCol.Alpha(0.35)).
					Children(func() {
						ui.Text(c, statusTxt).FontSize(9.5).FontWeight(600).TextColor(statusCol)
					})
			})

			// Optional inspect payload bar.
			if st.inspecting {
				ui.Row(c).
					WidthPercent(100).
					Padding(5, 8).
					Radius(7).
					Justify(ui.SpaceBetween).
					AlignItems(ui.Center).
					Background(t.SurfaceSubtle).
					Border(1, t.Border).
					Children(func() {
						ui.Text(c, "payload: {id, items[4]}").FontSize(9.5).TextColor(t.TextMuted)
						ui.Text(c, "1.2 KB").FontSize(9.5).FontWeight(600).TextColor(t.Warning)
					})
			}

			// Hover-revealed inline action pills (Run, Inspect, Disable).
			actionH := 28 * reveal
			if actionH > 1 {
				ui.Row(c).
					WidthPercent(100).
					Height(actionH).
					Gap(6).
					AlignItems(ui.Center).
					Opacity(clamp32(reveal, 0, 1)).
					Clip().
					BorderWidth(1, 0, 0, 0).
					BorderColor(t.Border.Alpha(0.8*reveal)).
					Padding(5, 0, 0, 0).
					Children(func() {
						makePill := func(key, label string, activeCol ui.Color, active bool) ui.Element {
							bg := t.SurfaceSubtle
							bd := t.Border
							tc := t.TextMuted
							if active {
								bg = activeCol.Alpha(0.16)
								bd = activeCol.Alpha(0.55)
								tc = activeCol
							}
							pill := ui.Box(c).
								Key(key).
								Padding(3, 9).
								Radius(99).
								Background(bg).
								Border(1, bd).
								Cursor(ui.CursorPointer)
							if pill.Hovered() && !active {
								pill.Background(t.SurfaceHover).Border(1, t.BorderStrong)
								tc = t.Text
							}
							pill.Children(func() {
								ui.Text(c, label).FontSize(10).FontWeight(600).TextColor(tc)
							})
							return pill
						}

						runLbl := "Run"
						if running {
							runLbl = "Running..."
						}
						if makePill("run", runLbl, t.Accent, running).Clicked() && !st.disabled {
							st.runningUntil = now.Add(1500 * time.Millisecond)
							st.runCount++
							c.AnimationFrame()
						}

						if makePill("inspect", "Inspect", t.Warning, st.inspecting).Clicked() {
							st.inspecting = !st.inspecting
						}

						disLbl := "Disable"
						if st.disabled {
							disLbl = "Enable"
						}
						if makePill("disable", disLbl, t.Danger, st.disabled).Clicked() {
							st.disabled = !st.disabled
							if st.disabled {
								st.runningUntil = time.Time{}
							}
						}
					})
			}
		})
	})

	return stage
}

func init() {
	registerDemo(
		"stacked-bars",
		"Weekly stacked column chart that dims siblings on hover and shows a crisp series breakdown.",
		func(c *ui.Context) {
			StackedBars(c)
		},
	)
	registerDemo(
		"blast-radius",
		"Interactive service dependency graph highlighting downstream blast-radius edges on hover.",
		func(c *ui.Context) {
			BlastRadiusGraph(c)
		},
	)
	registerDemo(
		"hub-diagram",
		"Central hub with 6 spokes and animated signal pulses that intensify with telemetry on hover.",
		func(c *ui.Context) {
			HubDiagram(c)
		},
	)
	registerDemo(
		"heart-rate",
		"Live scrolling ECG waveform with pulsing glowing lead dot and hover-accelerated BPM.",
		func(c *ui.Context) {
			HeartRateMonitor(c)
		},
	)
	registerDemo(
		"progress",
		"Row of 28 fine vertical progress ticks with Gaussian pointer ripple and scrubbed percentage.",
		func(c *ui.Context) {
			ProgressTicks(c)
		},
	)
	registerDemo(
		"heat-word",
		"12x5 activity heat map illuminating row/column thermal neighborhoods with live date readouts.",
		func(c *ui.Context) {
			HeatMapGrid(c)
		},
	)
	registerDemo(
		"action-node",
		"Workflow node card with ports and status ring revealing Run, Inspect, and Disable pills on hover.",
		func(c *ui.Context) {
			ActionNodeCard(c)
		},
	)
}
