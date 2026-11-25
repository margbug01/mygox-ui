package particles_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/particles"
	"github.com/margbug01/mygox-ui/coss/theme"
)

func TestAllRegisteredParticlesRenderHeadless(t *testing.T) {
	all := particles.All()
	if len(all) != 534 {
		t.Fatalf("expected the 530 particles of coss.com and mygox's own motion ones, 534, got %d", len(all))
	}

	for _, e := range all {
		entry := e
		t.Run(entry.ID, func(t *testing.T) {
			if entry.View == nil {
				t.Fatalf("particle %s has nil View", entry.ID)
			}
			if entry.Source == "" {
				t.Fatalf("particle %s has empty Source", entry.ID)
			}
			for _, dark := range []bool{false, true} {
				tt := ui.NewTester(func(c *ui.Context) {
					tk := theme.Apply(c, dark)
					ui.Column(c).WidthPercent(100).HeightPercent(100).Center().Padding(20).Background(tk.Background).Children(func() {
						entry.View(c)
					})
				}, 520, 320)
				tt.Frame()
				tt.Advance(100 * time.Millisecond)
				tt.Frame()
			}
		})
	}
}

func TestParticleInteractions(t *testing.T) {
	t.Run("p-accordion-4 open first two", func(t *testing.T) {
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, false)
			ui.Column(c).WidthPercent(100).HeightPercent(100).Padding(20).Children(func() {
				particles.PAccordion4(c)
			})
		}, 520, 420)
		tt.Frame()
		if !tt.HasText("Open items: None") {
			t.Fatalf("expected initial None state, got %v", tt.Texts())
		}
		tt.Click("Open First Two")
		tt.Frame()
		if !tt.HasText("Open items: item-1, item-2") {
			t.Fatalf("expected item-1, item-2 open after click, got %v", tt.Texts())
		}
	})

	t.Run("p-button-26 star toggle", func(t *testing.T) {
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, false)
			ui.Column(c).WidthPercent(100).HeightPercent(100).Center().Children(func() {
				particles.PButton26(c)
			})
		}, 320, 120)
		tt.Frame()
		if !tt.HasText("Star") || !tt.HasText("729") {
			t.Fatalf("expected Star 729, got %v", tt.Texts())
		}
		tt.Click("Star")
		tt.Frame()
		if !tt.HasText("Starred") || !tt.HasText("730") {
			t.Fatalf("expected Starred 730 after click, got %v", tt.Texts())
		}
	})

	t.Run("p-alert-dialog-1 open and close", func(t *testing.T) {
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, false)
			ui.Column(c).WidthPercent(100).HeightPercent(100).Center().Children(func() {
				particles.PAlertDialog1(c)
			})
		}, 600, 400)
		tt.Frame()
		tt.Click("Delete Account")
		tt.Frame()
		if !tt.HasText("Are you absolutely sure?") {
			t.Fatalf("expected dialog title after click, got %v", tt.Texts())
		}
		tt.Click("Cancel")
		tt.Frame()
		if tt.HasText("Are you absolutely sure?") {
			t.Fatalf("expected dialog closed after Cancel, got %v", tt.Texts())
		}
	})

	t.Run("no duplicate children on chained builders", func(t *testing.T) {
		checkOnce := func(name string, view func(*ui.Context), needle string) {
			t.Helper()
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, false)
				ui.Column(c).WidthPercent(100).HeightPercent(100).Padding(20).Children(func() {
					view(c)
				})
			}, 520, 240)
			tt.Frame()
			count := 0
			for _, txt := range tt.Texts() {
				if txt == needle {
					count++
				}
			}
			if count != 1 {
				t.Fatalf("%s: expected %q exactly once, got %d in %v", name, needle, count, tt.Texts())
			}
		}
		checkOnce("p-alert-2", particles.PAlert2, "Heads up!")
		checkOnce("p-alert-4", particles.PAlert4, "Heads up!")
		checkOnce("p-button-16", particles.PButton16, "Download")
		checkOnce("p-badge-11", particles.PBadge11, "Verified")
	})
}
