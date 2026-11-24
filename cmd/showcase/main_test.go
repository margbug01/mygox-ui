package main

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestUnifiedShowcaseAllThreeSuites(t *testing.T) {
	app := newShowcaseApp()
	tt := ui.NewTester(app.view, 1280, 840)

	// 1. UI (530 particles)
	tt.Frame()
	if !tt.HasText("mygox-ui") || !tt.HasText("All Particles") {
		t.Fatalf("expected UI tab header, got: %v", tt.Texts())
	}

	// 2. Switch to Flow tab (React Flow UI DB Schema & Sub-Flows)
	app.setSection(sectionFlow)
	tt.Frame()
	tt.Advance(50 * time.Millisecond)
	tt.Frame()
	if !tt.HasText("Node Palette") || !tt.HasText("DB Schema") {
		t.Fatalf("expected Flow tab UI, got: %v", tt.Texts())
	}

	// 3. Switch to Motion tab (33 Motion Primitives)
	app.setSection(sectionMotion)
	tt.Frame()
	tt.Advance(80 * time.Millisecond)
	tt.Frame()
	if !tt.HasText("All Motion Effects") || !tt.HasText("BorderTrail") {
		t.Fatalf("expected Motion tab UI with BorderTrail, got: %v", tt.Texts())
	}

	// 4. Switch to Bencho tab (184 / 184 catalog blocks, paginated 12 per page)
	app.setSection(sectionBencho)
	tt.Frame()
	tt.Advance(80 * time.Millisecond)
	tt.Frame()
	if !tt.HasText("Bencho Blocks") || !tt.HasText("Thinking orb") || !tt.HasText("Elastic slider") || !tt.HasText("Page 1 of 16") {
		t.Fatalf("expected Bencho tab UI with live blocks and pagination, got: %v", tt.Texts())
	}
	app.bencho.page = 2
	tt.Frame()
	if !tt.HasText("Page 2 of 16") {
		t.Fatalf("expected Bencho page 2 of 16, got: %v", tt.Texts())
	}
	app.bencho.selectKind("Loading")
	tt.Frame()
	if !tt.HasText("Loading Blocks") || !tt.HasText("Race bars") {
		t.Fatalf("expected filtered Loading blocks, got: %v", tt.Texts())
	}
	app.bencho.selectKind("")

	if dir := os.Getenv("MYGOX_UI_SNAPSHOT_DIR"); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir snapshot dir: %v", err)
		}
		// UI Light & Dark
		app.dark = false
		app.setSection(sectionCoss)
		app.selectedFam = ""
		tt.Frame()
		savePNG(t, filepath.Join(dir, "ui-light.png"), tt)

		app.dark = true
		tt.Frame()
		savePNG(t, filepath.Join(dir, "ui-dark.png"), tt)

		// Flow Light (React Flow UI DatabaseSchemaNode) & Dark (Sub-Flows + Spinning Conic Status + Placeholder)
		app.setSection(sectionFlow)
		app.flow.loadPreset("schema")
		app.flow.graph.SelectNode("products", false)
		app.dark = false
		tt.Frame()
		tt.Advance(50 * time.Millisecond)
		tt.Frame()
		savePNG(t, filepath.Join(dir, "flow-light.png"), tt)

		app.flow.loadPreset("subflows")
		app.flow.graph.SelectNode("synth", false)
		app.dark = true
		tt.Frame()
		tt.Advance(60 * time.Millisecond)
		tt.Frame()
		savePNG(t, filepath.Join(dir, "flow-dark.png"), tt)

		// Motion Primitives 33 — Light (Cards & Surfaces) & Dark (Text & Numbers)
		app.setSection(sectionMotion)
		app.motion.selectCategory("surfaces")
		app.dark = false
		tt.Frame()
		tt.Advance(100 * time.Millisecond)
		tt.Frame()
		savePNG(t, filepath.Join(dir, "motion-light.png"), tt)

		app.dark = true
		app.motion.selectCategory("text")
		tt.Frame()
		tt.Advance(140 * time.Millisecond)
		tt.Frame()
		savePNG(t, filepath.Join(dir, "motion-dark.png"), tt)

		// Bencho 184 — Light (Page 1 live blocks) & Dark (Drag sliders & physics)
		app.setSection(sectionBencho)
		app.bencho.selectKind("")
		app.dark = false
		tt.Frame()
		tt.Advance(100 * time.Millisecond)
		tt.Frame()
		savePNG(t, filepath.Join(dir, "bencho-light.png"), tt)

		app.dark = true
		app.bencho.selectKind("Drag")
		tt.Frame()
		tt.Advance(100 * time.Millisecond)
		tt.Frame()
		savePNG(t, filepath.Join(dir, "bencho-dark.png"), tt)
	}
}

func TestUnifiedShowcaseMobileResponsive(t *testing.T) {
	app := newShowcaseApp()
	tt := ui.NewTester(app.view, 390, 844)
	tt.Frame()
	if !tt.HasText("All Particles") {
		t.Fatalf("expected mobile UI tab, got: %v", tt.Texts())
	}

	app.setSection(sectionFlow)
	tt.Frame()
	tt.Advance(40 * time.Millisecond)
	tt.Frame()

	app.setSection(sectionMotion)
	tt.Frame()
	tt.Advance(60 * time.Millisecond)
	tt.Frame()
	if !tt.HasText("BorderTrail") {
		t.Fatalf("expected mobile Motion tab with BorderTrail, got: %v", tt.Texts())
	}

	app.setSection(sectionBencho)
	tt.Frame()
	tt.Advance(60 * time.Millisecond)
	tt.Frame()
	if !tt.HasText("Bencho Blocks") || !tt.HasText("Thinking orb") {
		t.Fatalf("expected mobile Bencho tab with Thinking orb, got: %v", tt.Texts())
	}

	if dir := os.Getenv("MYGOX_UI_SNAPSHOT_DIR"); dir != "" {
		app.dark = false
		app.setSection(sectionFlow)
		app.flow.loadPreset("schema")
		app.flow.fitted = false
		tt.Frame()
		tt.Advance(40 * time.Millisecond)
		tt.Frame()
		savePNG(t, filepath.Join(dir, "mobile-flow-light.png"), tt)
	}
}

func savePNG(t *testing.T, path string, tt *ui.Tester) {
	t.Helper()
	img := tt.Image()
	if img == nil {
		t.Fatalf("Tester.Image() returned nil for %s", path)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create %s: %v", path, err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode %s: %v", path, err)
	}
}
