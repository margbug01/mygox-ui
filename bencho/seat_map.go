package bencho

import (
	"fmt"
	"sort"
	"time"

	"github.com/egoist/mygo/ui"
)

// SeatMapSetting is where SeatMap's seats are.
type SeatMapSetting string

const (
	SeatMapPlane  SeatMapSetting = "Plane"  // a plane cabin with seat letters, row numbers and a price
	SeatMapCinema SeatMapSetting = "Cinema" // a cinema hall with its screen and a ticket count
)

// SeatMapTone is the colour of your seats.
type SeatMapTone string

const (
	SeatMapColour SeatMapTone = "Colour" // your seats in blue
	SeatMapMono   SeatMapTone = "Mono"   // your seats in the card's ink
)

// SeatMapOptions configures SeatMap; the zero value is bencho.dev's default.
type SeatMapOptions struct {
	Setting  SeatMapSetting // Plane by default
	Tone     SeatMapTone    // Colour by default
	OnChange func(seats []string)
}

const (
	seatMapW, seatMapH = 268, 338
	seatMapSeat        = 28
	seatMapPitch       = 34 // seat plus gap
	seatMapLeft        = 24
	seatMapPricePlane  = 38
	seatMapPricePlus   = 64 // the first two rows of the cabin
	seatMapPriceCinema = 14
)

// seatMapTaken lists the taken seats per row, columns 0 to 5.
var seatMapTaken = [7][6]bool{
	{true, true, false, false, false, true},
	{false, false, true, true, false, false},
	{true, false, false, false, true, true},
	{false, true, false, false, false, false},
	{true, false, true, true, true, false},
	{false, false, false, false, false, true},
	{true, true, false, true, false, false},
}

// seatMapPassengers are the travellers whose initials mark the seats you hold.
var seatMapPassengers = [2]string{"A", "L"}

type seatMapPick struct {
	row, col  int
	passenger int
}

type seatMapState struct {
	picks []seatMapPick // oldest first, at most two
}

func seatMapName(s SeatMapSetting, row, col int) string {
	if s == SeatMapCinema {
		return fmt.Sprintf("%c%d", 'A'+row, col+1)
	}
	return fmt.Sprintf("%d%c", 11+row, 'A'+col)
}

func seatMapRowY(row int) float32 {
	y := float32(38 + row*seatMapPitch)
	if row >= 2 {
		y += 10 // the exit row gap
	}
	return y
}

func seatMapColX(col int) float32 {
	x := float32(seatMapLeft + col*seatMapPitch)
	if col >= 3 {
		x += seatMapSeat - 6 // the aisle
	}
	return x
}

// SeatMap is a seat picker laid out like a plane cabin (or a cinema hall).
// Tap a free seat to choose it and tap again to let it go; you hold up to two
// seats, a third moves your earliest pick, and the price shows underneath.
func SeatMap(c *ui.Context, opts ...SeatMapOptions) ui.Element {
	var o SeatMapOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Setting == "" {
		o.Setting = SeatMapPlane
	}
	if o.Tone == "" {
		o.Tone = SeatMapColour
	}
	t := ThemeFrom(c)
	cinema := o.Setting == SeatMapCinema

	pick, pickInk, chipBg, chipInk := ui.Hex("#397cff"), ui.Hex("#ffffff"), ui.Hex("#397cff").Alpha(0.12), ui.Hex("#336cd3")
	if t.Dark {
		chipInk = ui.Hex("#7aa7ff")
	}
	if o.Tone == SeatMapMono {
		pick, pickInk, chipBg, chipInk = t.Text, t.Surface, t.SurfaceSubtle, t.Text
	}
	plus, plusBg := ui.Hex("#a78bfc"), ui.Hex("#a78bfc").Alpha(0.08)
	label := t.TextSubtle.Alpha(0.55)
	if t.Dark {
		label = t.TextSubtle
	}

	card := ui.Box(c).Size(seatMapW, seatMapH).Radius(24).Background(t.Surface)
	st := ui.Local(card, "seat_map", func() seatMapState {
		return seatMapState{picks: []seatMapPick{{2, 2, 0}, {2, 3, 1}}}
	})
	held := func(row, col int) (int, bool) {
		for _, p := range st.picks {
			if p.row == row && p.col == col {
				return p.passenger, true
			}
		}
		return 0, false
	}
	names := func() []string {
		ps := append([]seatMapPick(nil), st.picks...)
		sort.Slice(ps, func(i, j int) bool { return ps[i].row*6+ps[i].col < ps[j].row*6+ps[j].col })
		out := make([]string, len(ps))
		for i, p := range ps {
			out[i] = seatMapName(o.Setting, p.row, p.col)
		}
		return out
	}
	toggle := func(row, col int) {
		kept := st.picks[:0:0]
		for _, p := range st.picks {
			if p.row != row || p.col != col {
				kept = append(kept, p)
			}
		}
		if len(kept) == len(st.picks) { // a new pick
			used := [2]bool{}
			for _, p := range kept {
				used[p.passenger] = true
			}
			passenger := 0
			if used[0] {
				passenger = 1
			}
			if len(kept) == 2 { // the earliest pick moves here
				passenger = kept[0].passenger
				kept = kept[1:]
			}
			kept = append(kept, seatMapPick{row, col, passenger})
		}
		st.picks = kept
		if o.OnChange != nil {
			o.OnChange(names())
		}
	}

	total := 0
	for _, p := range st.picks {
		switch {
		case cinema:
			total += seatMapPriceCinema
		case p.row < 2:
			total += seatMapPricePlus
		default:
			total += seatMapPricePlane
		}
	}
	shownTotal := card.Spring("total", float32(total), ui.SpringOptions{Duration: 450 * time.Millisecond})

	card.Children(func() {
		if cinema {
			ui.Box(c).Absolute().Left(seatMapLeft).Top(10).Size(seatMapW-2*seatMapLeft, 20).Draw(func(p *ui.Painter, r ui.Rect) {
				var path ui.Path
				path.Cap(ui.RoundCap).MoveTo(r.X+4, r.Y+12).QuadTo(r.X+r.W/2, r.Y+2, r.X+r.W-4, r.Y+12)
				p.StrokePath(&path, 3, label)
			})
		} else {
			for col := range 6 {
				ui.Text(c, string(rune('A'+col))).Key(fmt.Sprint("h", col)).Absolute().Left(seatMapColX(col)).Top(16).Width(seatMapSeat).TextAlign(ui.Center).FontSize(11.5).TextColor(label)
			}
		}
		for row := range 7 {
			rowLabel := fmt.Sprint(11 + row)
			if cinema {
				rowLabel = string(rune('A' + row))
			}
			ui.Text(c, rowLabel).Key(fmt.Sprint("r", row)).Absolute().Left(seatMapColX(2) + seatMapSeat - 3).Top(seatMapRowY(row) + 6).Width(seatMapSeat).TextAlign(ui.Center).FontSize(11.5).TextColor(label)
			for col := range 6 {
				passenger, mine := held(row, col)
				seatMapSeatBox(c, t, seatMapSeatLook{
					x: seatMapColX(col), y: seatMapRowY(row), key: row*6 + col,
					name: seatMapName(o.Setting, row, col), taken: seatMapTaken[row][col], mine: mine,
					initial: seatMapPassengers[passenger], plus: !cinema && row < 2,
					pick: pick, pickInk: pickInk, plusInk: plus, plusBg: plusBg,
				}, func() { toggle(row, col) })
			}
		}

		ui.Row(c).Absolute().Left(16).Top(296).Gap(4).Children(func() {
			for _, n := range names() {
				ui.Box(c).Key(n).Height(26).Padding(0, 11).Radius(13).Background(chipBg).Center().Children(func() {
					ui.Text(c, n).FontSize(13).FontWeight(600).TextColor(chipInk)
				})
			}
		})
		right := fmt.Sprintf("$%d", int(shownTotal+0.5))
		if cinema {
			right = fmt.Sprintf("%d tickets", len(st.picks))
			if len(st.picks) == 1 {
				right = "1 ticket"
			}
		}
		ui.Text(c, right).Absolute().Right(16).Top(294).FontSize(21).TextColor(t.Text)
	})
	return card
}

type seatMapSeatLook struct {
	x, y          float32
	key           int
	name, initial string
	taken, mine   bool
	plus          bool
	pick, pickInk ui.Color
	plusInk       ui.Color
	plusBg        ui.Color
}

func seatMapSeatBox(c *ui.Context, t Theme, l seatMapSeatLook, toggle func()) {
	seat := ui.Box(c).Key(l.key).Absolute().Left(l.x).Top(l.y).Size(seatMapSeat, seatMapSeat).Label("Seat " + l.name)
	if !l.taken {
		seat.Cursor(ui.CursorPointer).OnClick(toggle)
	}
	on := seat.Spring("on", b2f(l.mine), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.35})
	hover := seat.Spring("hover", b2f(seat.Hovered() && !l.taken && !l.mine), ui.SpringOptions{Duration: 160 * time.Millisecond})
	press := seat.Spring("press", b2f(seat.Pressed() && !l.taken), ui.SpringOptions{Duration: 140 * time.Millisecond})
	s := 1 - 0.1*press
	seat.Scale(s, s)
	seat.Draw(func(p *ui.Painter, r ui.Rect) {
		switch {
		case l.taken:
			p.Fill(r, t.SurfaceHover, 8)
			return
		case l.plus:
			p.Fill(r, l.plusBg, 8)
			p.Stroke(r, lerpColor(l.plusInk, l.pick, hover), 8, 2)
		default:
			p.Stroke(r, lerpColor(t.BorderStrong, l.pick.Alpha(0.55), hover), 8, 2)
		}
		if hover > 0.01 {
			p.Fill(seatMapInset(r, 2), l.pick.Alpha(0.06*hover), 6)
		}
		if on > 0.01 {
			k := clamp32(on, 0, 1.2)
			box := seatMapInset(r, r.W/2*(1-k))
			p.Fill(box, l.pick, 8*k)
		}
	})
	if l.mine {
		seat.Center().Children(func() {
			ui.Text(c, l.initial).FontSize(15).FontWeight(600).TextColor(l.pickInk).Opacity(clamp32(on, 0, 1))
		})
	}
}

func init() {
	registerDemo(
		"seat-map",
		"A seat picker laid out like a plane cabin: tap a free seat to choose it, tap again to let it go, up to two seats with the price underneath.",
		func(c *ui.Context) {
			SeatMap(c)
		},
	)
}

func seatMapInset(r ui.Rect, d float32) ui.Rect {
	return ui.Rect{X: r.X + d, Y: r.Y + d, W: r.W - 2*d, H: r.H - 2*d}
}
