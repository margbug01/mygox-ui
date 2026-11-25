package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing sends the next message across the board, and it settles on it.
func TestFlapBoardFlipsToTheNextMessage(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.FlapBoard(c, bencho.FlapBoardOptions{Messages: []string{"AB", "CD"}})
	}, 286, 140)
	step(tt, 50*time.Millisecond)
	before := tt.Image()
	tt.ClickAt(140, 52)
	step(tt, 260*time.Millisecond) // the sweep reaches the middle tiles
	if changed(before, tt.Image(), 140) == 0 {
		t.Fatal("the flaps start falling")
	}
	step(tt, 2*time.Second)
	settled := tt.Image()
	step(tt, 200*time.Millisecond)
	if changed(settled, tt.Image(), 140) != 0 || changed(before, settled, 140) == 0 {
		t.Fatal("then the board rests on the new message")
	}
}
