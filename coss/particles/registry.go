package particles

import (
	"sort"
	"strings"
	"sync"

	"github.com/egoist/mygo/ui"
)

// Entry describes one migrated coss Particle component.
type Entry struct {
	ID          string
	Family      string
	Description string
	View        func(c *ui.Context)
	Source      string
}

var (
	mu       sync.RWMutex
	entries  = map[string]Entry{}
	order    []string
)

// Register adds a particle to the global registry.
func Register(e Entry) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := entries[e.ID]; !exists {
		order = append(order, e.ID)
	}
	entries[e.ID] = e
}

// All returns all registered particles in registration order.
func All() []Entry {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]Entry, 0, len(order))
	for _, id := range order {
		out = append(out, entries[id])
	}
	return out
}

// Find looks up a particle by ID (e.g. "p-button-1").
func Find(id string) (Entry, bool) {
	mu.RLock()
	defer mu.RUnlock()
	e, ok := entries[id]
	return e, ok
}

// Families returns sorted unique family names of all registered particles.
func Families() []string {
	mu.RLock()
	defer mu.RUnlock()
	seen := map[string]struct{}{}
	for _, e := range entries {
		seen[e.Family] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for f := range seen {
		out = append(out, f)
	}
	sort.Strings(out)
	return out
}

// ByFamily returns all registered particles belonging to a family (empty returns all).
func ByFamily(family string) []Entry {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]Entry, 0, len(order))
	for _, id := range order {
		e := entries[id]
		if family == "" || e.Family == family {
			out = append(out, e)
		}
	}
	return out
}

// extractFunc extracts a top-level Go function block from an embedded source file.
func extractFunc(src, funcName string) string {
	marker := "func " + funcName + "("
	start := strings.Index(src, marker)
	if start < 0 {
		return ""
	}
	depth := 0
	inBody := false
	for i := start; i < len(src); i++ {
		switch src[i] {
		case '{':
			depth++
			inBody = true
		case '}':
			depth--
			if inBody && depth == 0 {
				return strings.TrimSpace(src[start : i+1])
			}
		}
	}
	return strings.TrimSpace(src[start:])
}
