package icons_test

import (
	"testing"

	"github.com/margbug01/mygox-ui/coss/icons"
)

func TestAllIconsValidSVG(t *testing.T) {
	names := icons.Names()
	if len(names) == 0 {
		t.Fatal("icons.Names() returned empty slice")
	}

	for _, name := range names {
		svg, ok := icons.ByName(name)
		if !ok || svg == nil {
			t.Fatalf("ByName(%q) = (%v, %v), want non-nil SVG", name, svg, ok)
		}
		w, h := svg.Size()
		if w <= 0 || h <= 0 {
			t.Errorf("icon %q has non-positive size (%v, %v)", name, w, h)
		}
		if got := icons.Get(name); got != svg {
			t.Errorf("Get(%q) did not return cached *ui.SVG instance", name)
		}
	}
}

func TestLookupFormsAndTypedAccessors(t *testing.T) {
	want := icons.ChevronDown()
	if want == nil {
		t.Fatal("ChevronDown() returned nil")
	}
	for _, query := range []string{"chevron-down", "ChevronDown", "ChevronDownIcon"} {
		got, ok := icons.ByName(query)
		if !ok || got != want {
			t.Errorf("ByName(%q) = (%v, %v), want same *ui.SVG as ChevronDown()", query, got, ok)
		}
	}

	// Alias names resolve to the same cached SVG as their canonical form.
	if icons.AlertCircle() != icons.CircleAlert() {
		t.Error("AlertCircle() and CircleAlert() should share the same cached SVG")
	}
	if icons.Loader2() != icons.LoaderCircle() {
		t.Error("Loader2() and LoaderCircle() should share the same cached SVG")
	}
	if icons.MoreHorizontal() != icons.Ellipsis() {
		t.Error("MoreHorizontal() and Ellipsis() should share the same cached SVG")
	}

	if svg, ok := icons.ByName("does-not-exist-icon"); ok || svg != nil {
		t.Errorf("unknown icon returned (%v, %v), want (nil, false)", svg, ok)
	}
	if got := icons.Get("does-not-exist-icon"); got != nil {
		t.Errorf("Get(unknown) = %v, want nil", got)
	}
}

func TestNamesReturnsDefensiveCopy(t *testing.T) {
	n1 := icons.Names()
	orig := n1[0]
	n1[0] = "mutated"
	n2 := icons.Names()
	if n2[0] != orig {
		t.Fatalf("Names() leaked internal slice: got %q, want %q", n2[0], orig)
	}
}
