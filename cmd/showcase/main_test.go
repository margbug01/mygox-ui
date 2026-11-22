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

	// 2. Switch to Flow tab
	app.setSection(sectionFlow)
	tt.Frame()
	tt.Advance(50 * time.Millisecond)
	tt.Frame()
	if !tt.HasText("Node Palette") || !tt.HasText("Inspector") {
		t.Fatalf("expected Flow tab UI, got: %v", tt.Texts())
	}

	// 3. Switch to Motion tab
	app.setSection(sectionMotion)
	tt.Frame()
	tt.Advance(80 * time.Millisecond)
	tt.Frame()
	if !tt.HasText("All Motion Effects") || !tt.HasText("BorderBeam") {
		t.Fatalf("expected Motion tab UI, got: %v", tt.Texts())
	}

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

		// Flow Light & Dark (in coss style!)
		app.setSection(sectionFlow)
		app.flow.loadPreset("ai-agent")
		app.flow.graph.SelectNode("llm", false)
		app.dark = false
		tt.Frame()
		tt.Advance(40 * time.Millisecond)
		tt.Frame()
		savePNG(t, filepath.Join(dir, "flow-light.png"), tt)

		app.dark = true
		tt.Frame()
		tt.Advance(40 * time.Millisecond)
		tt.Frame()
		savePNG(t, filepath.Join(dir, "flow-dark.png"), tt)

		// Motion Light & Dark (in coss style!)
		app.setSection(sectionMotion)
		app.motion.selectCategory("")
		app.dark = false
		tt.Frame()
		tt.Advance(80 * time.Millisecond)
		tt.Frame()
		savePNG(t, filepath.Join(dir, "motion-light.png"), tt)

		app.dark = true
		app.motion.selectCategory("layout")
		app.motion.confettiTrigger++
		tt.Frame()
		tt.Advance(140 * time.Millisecond)
		tt.Frame()
		savePNG(t, filepath.Join(dir, "motion-dark.png"), tt)
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
	if !tt.HasText("BorderBeam") {
		t.Fatalf("expected mobile Motion tab, got: %v", tt.Texts())
	}

	if dir := os.Getenv("MYGOX_UI_SNAPSHOT_DIR"); dir != "" {
		app.dark = false
		app.setSection(sectionFlow)
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
