package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The pill opens into the list; each decision folds a person away and the
// pill counts down until everyone is answered.
func TestAccessRequestApprovesAndRejects(t *testing.T) {
	decided := map[string]bool{}
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.AccessRequest(c, bencho.AccessRequestDefaultPeople(), bencho.AccessRequestOptions{
			OnDecide: func(h string, ok bool) { decided[h] = ok },
		})
	}, 300, 300)
	step(tt, 100*time.Millisecond)
	if !tt.HasText("2 others") || tt.HasText("Approve") {
		t.Fatal("starts as the pill")
	}
	tt.Click("requested access")
	step(tt, 900*time.Millisecond)
	if !tt.HasText("Approve") {
		t.Fatal("pressing opens the list")
	}
	tt.Click("Approve")
	step(tt, 600*time.Millisecond)
	if !decided["hanaito"] || !tt.HasText("1 other") {
		t.Fatalf("approving the first folds them away, decided %v", decided)
	}
	tt.Click("Reject")
	step(tt, 600*time.Millisecond)
	tt.Click("Reject")
	step(tt, 900*time.Millisecond)
	if len(decided) != 3 || decided["mira.k"] || !tt.HasText("All caught up") {
		t.Fatalf("answering everyone closes it on All caught up, decided %v", decided)
	}
}
