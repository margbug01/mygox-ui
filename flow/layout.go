package flow

import "sort"

// AutoLayout arranges all nodes into clean topological ranks with reduced edge crossings.
func (g *Graph) AutoLayout(dir LayoutDirection) {
	if len(g.Nodes) == 0 {
		return
	}
	g.PushUndo()

	if dir == "" {
		dir = LayoutLR
	}

	byID := make(map[string]*Node, len(g.Nodes))
	inDeg := make(map[string]int, len(g.Nodes))
	adj := make(map[string][]string, len(g.Nodes))
	parents := make(map[string][]string, len(g.Nodes))
	orderIdx := make(map[string]int, len(g.Nodes))

	for i, n := range g.Nodes {
		byID[n.ID] = n
		inDeg[n.ID] = 0
		orderIdx[n.ID] = i
	}
	for _, e := range g.Edges {
		if byID[e.Source] == nil || byID[e.Target] == nil || e.Source == e.Target {
			continue
		}
		adj[e.Source] = append(adj[e.Source], e.Target)
		parents[e.Target] = append(parents[e.Target], e.Source)
		inDeg[e.Target]++
	}

	// Longest-path rank assignment via Kahn topological sweep (cycle-safe fallback).
	rank := make(map[string]int, len(g.Nodes))
	queue := make([]string, 0, len(g.Nodes))
	for _, n := range g.Nodes {
		if inDeg[n.ID] == 0 {
			queue = append(queue, n.ID)
		}
	}
	visited := make(map[string]bool, len(g.Nodes))
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		visited[cur] = true
		for _, nxt := range adj[cur] {
			if rank[cur]+1 > rank[nxt] {
				rank[nxt] = rank[cur] + 1
			}
			inDeg[nxt]--
			if inDeg[nxt] == 0 {
				queue = append(queue, nxt)
			}
		}
	}
	// Handle any cycle members deterministically.
	for _, n := range g.Nodes {
		if !visited[n.ID] {
			best := 0
			for _, p := range parents[n.ID] {
				if visited[p] && rank[p]+1 > best {
					best = rank[p] + 1
				}
			}
			rank[n.ID] = best
		}
	}

	maxRank := 0
	layers := map[int][]*Node{}
	for _, n := range g.Nodes {
		r := rank[n.ID]
		if r > maxRank {
			maxRank = r
		}
		layers[r] = append(layers[r], n)
	}

	// Barycenter ordering pass to minimize edge crossings between adjacent ranks.
	posInLayer := map[string]float32{}
	for r := 0; r <= maxRank; r++ {
		layer := layers[r]
		sort.SliceStable(layer, func(i, j int) bool {
			bi := barycenter(layer[i].ID, parents[layer[i].ID], posInLayer, float32(orderIdx[layer[i].ID]))
			bj := barycenter(layer[j].ID, parents[layer[j].ID], posInLayer, float32(orderIdx[layer[j].ID]))
			if bi != bj {
				return bi < bj
			}
			return orderIdx[layer[i].ID] < orderIdx[layer[j].ID]
		})
		for idx, n := range layer {
			posInLayer[n.ID] = float32(idx)
		}
	}

	const (
		rankGapLR = float32(96)
		nodeGapLR = float32(42)
		rankGapTB = float32(84)
		nodeGapTB = float32(56)
	)

	if dir == LayoutLR {
		for _, n := range g.Nodes {
			for i := range n.Ports {
				if n.Ports[i].Kind == PortInput {
					n.Ports[i].Side = PortLeft
				} else {
					n.Ports[i].Side = PortRight
				}
			}
		}
		curX := float32(60)
		for r := 0; r <= maxRank; r++ {
			layer := layers[r]
			if len(layer) == 0 {
				continue
			}
			var maxW, totalH float32
			for i, n := range layer {
				b := n.Bounds()
				if b.W > maxW {
					maxW = b.W
				}
				totalH += b.H
				if i > 0 {
					totalH += nodeGapLR
				}
			}
			curY := -totalH * 0.5
			for _, n := range layer {
				b := n.Bounds()
				n.X = curX
				n.Y = curY
				curY += b.H + nodeGapLR
			}
			curX += maxW + rankGapLR
		}
	} else {
		for _, n := range g.Nodes {
			for i := range n.Ports {
				if n.Ports[i].Kind == PortInput {
					n.Ports[i].Side = PortTop
				} else {
					n.Ports[i].Side = PortBottom
				}
			}
		}
		curY := float32(60)
		for r := 0; r <= maxRank; r++ {
			layer := layers[r]
			if len(layer) == 0 {
				continue
			}
			var maxH, totalW float32
			for i, n := range layer {
				b := n.Bounds()
				if b.H > maxH {
					maxH = b.H
				}
				totalW += b.W
				if i > 0 {
					totalW += nodeGapTB
				}
			}
			curX := -totalW * 0.5
			for _, n := range layer {
				b := n.Bounds()
				n.X = curX
				n.Y = curY
				curX += b.W + nodeGapTB
			}
			curY += maxH + rankGapTB
		}
	}
}

func barycenter(id string, parents []string, pos map[string]float32, fallback float32) float32 {
	var sum float32
	count := 0
	for _, p := range parents {
		if v, ok := pos[p]; ok {
			sum += v
			count++
		}
	}
	if count == 0 {
		return fallback
	}
	return sum / float32(count)
}
