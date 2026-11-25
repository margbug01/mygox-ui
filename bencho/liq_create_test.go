package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The pill spreads into the menu, and picking an item pulls it back into a pill.
func TestCreateMenuSpreadsAndPicks(t *testing.T) {
	made := ""
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.CreateMenu(c, bencho.CreateMenuDefaultItems(), bencho.CreateMenuOptions{OnCreate: func(it bencho.CreateMenuItem) { made = it.Label }})
	}, 240, 200)
	step(tt, 100*time.Millisecond)
	if !tt.HasText("Create") || tt.HasText("Board") {
		t.Fatal("starts as the + Create pill")
	}
	tt.Click("Create")
	step(tt, 900*time.Millisecond)
	if tt.HasText("Create") || !tt.HasText("Board") {
		t.Fatal("pressing spreads the pill into the menu")
	}
	tt.Click("Board")
	step(tt, 900*time.Millisecond)
	if made != "Board" || !tt.HasText("Create") || tt.HasText("Board") {
		t.Fatalf("a pick makes the item and pulls the menu back in, made %q", made)
	}
}
