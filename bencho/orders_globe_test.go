package bencho_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Orders land and the total goes up; dragging spins the globe past where
// it turns on its own; the pill switches to this week's numbers.
func TestOrdersGlobeLandsOrdersSpinsAndSwitchesPeriod(t *testing.T) {
	view := func(o bencho.OrdersGlobeOptions) *ui.Tester {
		return ui.NewTester(func(c *ui.Context) {
			th := bencho.ApplyTheme(c, false)
			ui.Box(c).Size(296, 462).Background(th.SurfaceSubtle).Children(func() { bencho.OrdersGlobe(c, o) })
		}, 296, 462)
	}
	live := view(bencho.OrdersGlobeOptions{Rate: 100})
	step(live, 100*time.Millisecond)
	if !live.HasText("119") || !live.HasText("vs yesterday") {
		t.Fatalf("today starts at 119 orders, have %v", live.Texts())
	}
	step(live, 1200*time.Millisecond)
	if live.HasText("119") {
		t.Fatal("orders land and the total goes up")
	}

	still, spun := view(bencho.OrdersGlobeOptions{NoOrders: true}), view(bencho.OrdersGlobeOptions{NoOrders: true})
	step(still, 100*time.Millisecond)
	step(spun, 100*time.Millisecond)
	spun.Press(148, 182)
	for x := float32(148); x < 240; x += 8 {
		spun.Move(x, 182)
		step(spun, 16*time.Millisecond)
		step(still, 16*time.Millisecond)
	}
	spun.Release(240, 182)
	step(spun, 100*time.Millisecond)
	step(still, 100*time.Millisecond)
	if changed(still.Image(), spun.Image(), 300) < 20000 {
		t.Fatal("dragging spins the globe")
	}

	if err := still.Click("Today"); err != nil {
		t.Fatal(err)
	}
	step(still, 100*time.Millisecond)
	if !still.HasText("This week") || !still.HasText("vs last week") || !still.HasText(strconv.Itoa(786)) {
		t.Fatalf("the pill switches to this week, have %v", still.Texts())
	}
}
