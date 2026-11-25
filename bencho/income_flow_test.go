package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pointing at a month shows its total against the month before; pressing
// moves on to the next year.
func TestIncomeFlowMarksMonthAndMorphsYear(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.IncomeFlow(c)
	}, 360, 326)
	step(tt, 50*time.Millisecond)
	if !tt.HasText("$82,229") || !tt.HasText("$76,330") || !tt.HasText("7.7%") {
		t.Fatal("rests on March against February")
	}
	tt.Move(200, 172) // the ribbon over December
	step(tt, time.Second)
	if !tt.HasText("$53,870") || !tt.HasText("$28,310") || !tt.HasText("Nov") {
		t.Fatalf("pointing at December shows its total against November, got %v", tt.Texts())
	}
	tt.ClickAt(200, 172)
	tt.Move(5, 5)
	step(tt, time.Second)
	if !tt.HasText("$131,604") || !tt.HasText("$118,930") {
		t.Fatalf("pressing morphs into the next year, got %v", tt.Texts())
	}
}
