package bencho_test

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

func TestCatalogCompletenessAndUniqueness(t *testing.T) {
	cat := bencho.Catalog()
	if len(cat) != 184 {
		t.Fatalf("expected 184 bencho.dev blocks, got %d", len(cat))
	}

	seen := make(map[string]bool, len(cat))
	validKind := make(map[bencho.Interaction]bool, len(bencho.AllInteractions))
	for _, k := range bencho.AllInteractions {
		validKind[k] = true
	}

	live := 0
	for _, b := range cat {
		if b.Slug == "" || b.Title == "" {
			t.Fatalf("block has empty slug or title: %+v", b)
		}
		if seen[b.Slug] {
			t.Fatalf("duplicate block slug: %s", b.Slug)
		}
		seen[b.Slug] = true
		if !validKind[b.Interaction] {
			t.Fatalf("block %s has unknown interaction %q", b.Slug, b.Interaction)
		}
		if b.Implemented() {
			live++
		}
	}

	if live != benchogoTotal || live != bencho.ImplementedCount() {
		t.Fatalf("expected all %d blocks implemented: ImplementedCount() = %d, counted %d", benchogoTotal, bencho.ImplementedCount(), live)
	}
}

const benchogoTotal = 184

func TestEverySingleBlockRendersInCardAndHandlesPointer(t *testing.T) {
	cat := bencho.Catalog()
	for _, b := range cat {
		block := b
		if !block.Implemented() {
			t.Fatalf("block %s (%s) is not implemented", block.Slug, block.Title)
		}
		for _, dark := range []bool{false, true} {
			view := func(c *ui.Context) {
				th := bencho.ApplyTheme(c, dark)
				ui.Box(c).Size(320, 210).Center().Background(th.Surface).Children(func() {
					block.Demo(c)
				})
			}
			tt := ui.NewTester(view, 320, 210)
			tt.Frame()
			tt.Move(160, 105)
			tt.Advance(32 * time.Millisecond)
			tt.Frame()
			tt.Press(160, 105)
			tt.Move(190, 105)
			tt.Advance(32 * time.Millisecond)
			tt.Frame()
			tt.Release(190, 105)
			tt.Advance(32 * time.Millisecond)
			tt.Frame()
		}
	}
}

func TestAllImplementedDemosRenderDarkAndLight(t *testing.T) {
	cat := bencho.Catalog()
	for _, dark := range []bool{false, true} {
		mode := "light"
		if dark {
			mode = "dark"
		}
		t.Run(mode, func(t *testing.T) {
			view := func(c *ui.Context) {
				th := bencho.ApplyTheme(c, dark)
				ui.Scroll(c).WidthPercent(100).HeightPercent(100).Background(th.Background).Children(func() {
					ui.Column(c).WidthPercent(100).Padding(16).Gap(16).Children(func() {
						for _, b := range cat {
							if !b.Implemented() {
								continue
							}
							block := b
							ui.Column(c).Key(block.Slug).Gap(6).Children(func() {
								ui.Text(c, block.Title)
								block.Demo(c)
							})
						}
					})
				})
			}

			tt := ui.NewTester(view, 600, 1200)
			tt.Frame()
			for i := 0; i < 6; i++ {
				tt.Advance(50 * time.Millisecond)
				tt.Frame()
			}
			for _, b := range cat {
				if b.Implemented() && !tt.HasText(b.Title) {
					t.Fatalf("expected rendered block title %q, got texts: %v", b.Title, tt.Texts())
				}
			}
		})
	}
}

// Pointing at the orb wakes it up: its dots travel further in the same time.
func TestThinkingOrbWakesOnHover(t *testing.T) {
	for _, style := range []bencho.ThinkingOrbStyle{bencho.ThinkingOrbSpinning, bencho.ThinkingOrbWeaving} {
		motion := func(hover bool) int {
			view := func(c *ui.Context) {
				bencho.ApplyTheme(c, false)
				ui.Box(c).Size(140, 172).Children(func() { bencho.ThinkingOrb(c, bencho.ThinkingOrbOptions{Style: style}) })
			}
			tt := ui.NewTester(view, 140, 172)
			tt.Frame()
			if hover {
				tt.Move(70, 70)
			}
			step(tt, 900*time.Millisecond)
			before := tt.Image()
			before = &image.RGBA{Pix: append([]uint8(nil), before.Pix...), Stride: before.Stride, Rect: before.Rect}
			step(tt, 32*time.Millisecond)
			return changed(before, tt.Image(), 140)
		}
		idle, awake := motion(false), motion(true)
		if idle == 0 || awake < idle*3/2 {
			t.Fatalf("%s: hovered orb should move more than idle: idle moved %d, hovered %d", style, idle, awake)
		}
	}
}

func step(tt *ui.Tester, d time.Duration) {
	for ; d > 0; d -= 16 * time.Millisecond {
		tt.Advance(16 * time.Millisecond)
		tt.Frame()
	}
}

// changed sums how much the green channel moved between a and b above maxY,
// which for small motions grows with how far things moved.
func changed(a, b *image.RGBA, maxY int) int {
	n := 0
	for i := 1; i < min(len(a.Pix), maxY*a.Stride); i += 4 {
		n += max(int(a.Pix[i])-int(b.Pix[i]), int(b.Pix[i])-int(a.Pix[i]))
	}
	return n
}

// Dragging sets the value; pulling past an end stretches the pill, which
// springs back to its width once let go.
func TestElasticSliderStretchesPastEndsAndSpringsBack(t *testing.T) {
	val := float32(64)
	view := func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(400, 100).Padding(26, 70).Background(th.SurfaceSubtle).Children(func() { bencho.ElasticSlider(c, &val) })
	}
	tt := ui.NewTester(view, 400, 100)
	tt.Frame()
	width := func() int { // the pill's extent along its middle row
		img := tt.Image()
		bg := img.RGBAAt(0, 50)
		x0, x1 := -1, -1
		for x := 0; x < 400; x++ {
			if img.RGBAAt(x, 50) != bg {
				if x0 < 0 {
					x0 = x
				}
				x1 = x
			}
		}
		return x1 - x0 + 1
	}
	if w := width(); w < 258 || w > 262 {
		t.Fatalf("pill at rest is %dpx wide, want 260", w)
	}

	tt.Press(200, 50) // 50% across
	tt.Frame()
	if val != 50 {
		t.Fatalf("pressing the middle set %.0f, want 50", val)
	}
	tt.Move(395, 50) // 65px past the right end
	step(tt, 500*time.Millisecond)
	if val != 100 {
		t.Fatalf("dragging past the end set %.0f, want 100", val)
	}
	if w := width(); w < 270 {
		t.Fatalf("pulled pill is %dpx wide, want it stretched past 260", w)
	}
	tt.Release(395, 50)
	step(tt, time.Second)
	if w := width(); w < 258 || w > 262 {
		t.Fatalf("released pill is %dpx wide, want it back at 260", w)
	}
}

// Voting turns rows into results while a few more voters come in; pressing
// another row moves the vote, pressing yours again takes it back.
func TestLiveVoteMovesAndTakesBackVote(t *testing.T) {
	opts := []bencho.VoteOption{{Label: "Alpha", Votes: 6}, {Label: "Beta", Votes: 3}}
	vote := -1
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.LiveVote(c, opts, &vote, bencho.LiveVoteOptions{Question: "Pick one"})
	}
	tt := ui.NewTester(view, 360, 260)
	tt.Frame()
	if !tt.HasText("9 votes · Closes Friday") || tt.HasText("67%") {
		t.Fatalf("before voting want 9 votes and no results, got %v", tt.Texts())
	}

	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	step(tt, 300*time.Millisecond)
	if vote != 1 || !tt.HasText("10 votes · Closes Friday") {
		t.Fatalf("after voting Beta want vote 1 and 10 votes; vote %d, got %v", vote, tt.Texts())
	}
	step(tt, 3*time.Second)
	if !tt.HasText("13 votes · Closes Friday") || !tt.HasText("38%") {
		t.Fatalf("three voters should have come in, Beta at 5 of 13, got %v", tt.Texts())
	}

	tt.Click("Alpha")
	tt.Frame()
	if vote != 0 {
		t.Fatalf("pressing Alpha should move the vote there, vote %d", vote)
	}
	tt.Click("Alpha")
	step(tt, 600*time.Millisecond)
	if vote != -1 || !tt.HasText("12 votes · Closes Friday") || tt.HasText("%") {
		t.Fatalf("pressing your vote again should take it back, vote %d, got %v", vote, tt.Texts())
	}
}

// Pressing send works for a moment, then flips to a receipt with the arrival
// time; Done flips back.
func TestConfirmSendFlipsToReceipt(t *testing.T) {
	sent := 0
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ConfirmSend(c, bencho.ConfirmSendOptions{OnSend: func() { sent++ }})
	}
	tt := ui.NewTester(view, 320, 366)
	tt.Frame()
	if err := tt.Click("Send €150"); err != nil {
		t.Fatal(err)
	}
	step(tt, 500*time.Millisecond)
	if !tt.HasText("Sending…") || sent != 0 {
		t.Fatalf("half a second in want Sending… and nothing sent yet (%d), got %v", sent, tt.Texts())
	}
	arrived := tt.Now()
	step(tt, 2*time.Second)
	if sent != 1 || !tt.HasText("€150 sent") || !tt.HasText(arrived.Add(1300*time.Millisecond).Format("15:04")) {
		t.Fatalf("after the work want one send and the receipt with its time, sent %d, got %v", sent, tt.Texts())
	}
	if err := tt.Click("Done"); err != nil {
		t.Fatal(err)
	}
	step(tt, time.Second)
	if !tt.HasText("Send €150") {
		t.Fatalf("Done should flip back to the transfer, got %v", tt.Texts())
	}
}

func TestHoldDeleteEarlyReleaseCancelsAndFullHoldCompletes(t *testing.T) {
	deleted := 0
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, true)
		ui.Box(c).Padding(16).Children(func() {
			bencho.HoldDelete(c, func() {
				deleted++
			}, bencho.HoldDeleteOptions{
				HoldDuration: 300 * time.Millisecond,
				ResetAfter:   400 * time.Millisecond,
			})
		})
	}

	tt := ui.NewTester(view, 280, 120)
	tt.Frame()

	r, ok := tt.Find("Hold to delete")
	if !ok {
		t.Fatalf("expected Hold to delete button, got %v", tt.Texts())
	}
	cx, cy := r.X+r.W*0.5, r.Y+r.H*0.5

	// Short hold should cancel without firing onDelete.
	tt.Press(cx, cy)
	tt.Frame()
	tt.Advance(100 * time.Millisecond)
	tt.Frame()
	tt.Release(cx, cy)
	tt.Frame()
	if deleted != 0 {
		t.Fatalf("early release should not trigger delete, got %d", deleted)
	}

	// Full hold past 300ms should trigger onDelete once.
	tt.Press(cx, cy)
	tt.Frame()
	for i := 0; i < 4; i++ {
		tt.Advance(100 * time.Millisecond)
		tt.Frame()
	}
	tt.Release(cx, cy)
	tt.Frame()
	if deleted != 1 || !tt.HasText("Deleted") {
		t.Fatalf("expected deleted=1 and Deleted text, got deleted=%d texts=%v", deleted, tt.Texts())
	}
}

// Joining pops you into the room under Leave and speakers take turns; Leave
// takes you out again.
func TestHuddleJoinSpeakAndLeave(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Huddle(c)
	}
	tt := ui.NewTester(view, 300, 202)
	tt.Frame()
	if tt.HasText("You") || !tt.HasText("Join huddle") {
		t.Fatalf("before joining want only Ines and Tom under Join huddle, got %v", tt.Texts())
	}
	if err := tt.Click("Join huddle"); err != nil {
		t.Fatal(err)
	}
	step(tt, 600*time.Millisecond)
	if !tt.HasText("You") || !tt.HasText("Leave") {
		t.Fatalf("after joining want You in the room and Leave, got %v", tt.Texts())
	}
	// You speak first: your avatar darkens under the wave; then Tom's turn.
	you := tt.Image().RGBAAt(227, 72)
	step(tt, 1300*time.Millisecond)
	if now := tt.Image().RGBAAt(227, 72); now == you {
		t.Fatalf("your avatar should brighten when Tom takes the turn, still %v", now)
	}
	if err := tt.Click("Leave"); err != nil {
		t.Fatal(err)
	}
	step(tt, 700*time.Millisecond)
	if tt.HasText("You") || !tt.HasText("Join huddle") {
		t.Fatalf("after leaving want Join huddle back without you, got %v", tt.Texts())
	}
}
