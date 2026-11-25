package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

func TestProfileCardOpensAfterARestAndFollows(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ProfileCard(c)
	}
	tt := ui.NewTester(view, 300, 324)
	tt.Frame()
	tt.Move(26, 162) // Mara's face
	step(tt, 200*time.Millisecond)
	if tt.HasText("Mara Quinn") {
		t.Fatal("the card should wait for the 400ms rest before opening")
	}
	step(tt, 400*time.Millisecond)
	if !tt.HasText("Mara Quinn") {
		t.Fatalf("after resting on Mara want her card, got %v", tt.Texts())
	}
	tt.Move(98, 162) // run along to Juno's face: the card changes at once
	step(tt, 100*time.Millisecond)
	if !tt.HasText("Juno Park") || tt.HasText("Mara Quinn") {
		t.Fatalf("running along the faces want Juno's card, got %v", tt.Texts())
	}
	if err := tt.Click("Follow"); err != nil {
		t.Fatal(err)
	}
	step(tt, 100*time.Millisecond)
	if !tt.HasText("Following") {
		t.Fatalf("after Follow want Following, got %v", tt.Texts())
	}
	tt.Move(150, 300)
	step(tt, 700*time.Millisecond)
	if tt.HasText("Juno Park") {
		t.Fatalf("leaving the faces and the card want it closed, got %v", tt.Texts())
	}
}
