package bencho

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
)

// OrdersGlobeTone is how OrdersGlobe colours its cities.
type OrdersGlobeTone string

const (
	OrdersGlobeColour OrdersGlobeTone = "Colour" // the top three cities in their own hues; the default
	OrdersGlobeMono   OrdersGlobeTone = "Mono"   // all in ink
)

// OrdersGlobeOptions configures OrdersGlobe.
type OrdersGlobeOptions struct {
	Tone  OrdersGlobeTone // "Colour" by default
	Rate  float32         // how often an order lands, 0 to 100, 0 for none; 40 by default (use NoOrders for none)
	Speed float32         // how fast the globe turns on its own, 0 to 100; 30 by default
	// NoOrders stops orders landing, as Rate 0 does on bencho.dev.
	NoOrders bool
}

const (
	ordersGlobeW    = 296
	ordersGlobeH    = 462
	ordersGlobeR    = 108 // the globe's radius
	ordersGlobeCX   = 148
	ordersGlobeCY   = 182
	ordersGlobeTilt = 18 // degrees the north pole leans toward the viewer
	ordersGlobeStep = 2.8
	ordersGlobeCols = 30 // dots across a city's bar
)

type ordersGlobeCity struct {
	name     string
	lat, lon float64
	today    int
}

var ordersGlobeCities = []ordersGlobeCity{
	{"London", 51.5, -0.1, 18}, {"Toronto", 43.7, -79.4, 15}, {"New York", 40.7, -74, 13},
	{"Berlin", 52.5, 13.4, 12}, {"Paris", 48.9, 2.35, 11}, {"Los Angeles", 34, -118.2, 10},
	{"São Paulo", -23.5, -46.6, 9}, {"Lagos", 6.5, 3.4, 8}, {"Tokyo", 35.7, 139.7, 8},
	{"Sydney", -33.9, 151.2, 6}, {"Mumbai", 19, 72.8, 5}, {"Singapore", 1.3, 103.8, 4},
}

// ordersGlobePeriods are the pill's choices: how many times today's
// orders each holds, and the change it reports.
var ordersGlobePeriods = []struct {
	label, versus string
	scale         float64
	change        string
}{
	{"Today", "vs yesterday", 1, "+3.1%"},
	{"This week", "vs last week", 6.6, "+5.4%"},
	{"This month", "vs last month", 27.3, "+12.8%"},
}

// ordersGlobeHues are the top three cities' colours, first to third.
var ordersGlobeHues = [3]ui.Color{ui.Hex("#4374ea"), ui.Hex("#895ff1"), ui.Hex("#ec5d9f")}

type ordersGlobeState struct {
	period  int
	counts  [3][]int // per period, per city
	lon     float64  // the meridian facing the viewer
	spin    float64  // degrees a second left from a flick
	last    time.Time
	next    time.Time // when the next order lands
	landed  []time.Time
	rng     *rand.Rand
	pressed bool
}

// ordersGlobeDot is a dot of the globe's lattice.
type ordersGlobeDot struct {
	lat, lon float64
	land     bool
}

var (
	ordersGlobeDotsOnce sync.Once
	ordersGlobeDots     []ordersGlobeDot
)

// ordersGlobeLand outlines the continents, coarsely, as lon/lat rings.
var ordersGlobeLand = [][][2]float64{
	{{-168, 65}, {-165, 70}, {-140, 70}, {-120, 72}, {-95, 73}, {-80, 70}, {-75, 62}, {-65, 60}, {-55, 52}, {-66, 45}, {-70, 42}, {-76, 35}, {-81, 31}, {-80, 25}, {-82, 28}, {-88, 30}, {-97, 27}, {-97, 22}, {-92, 18}, {-88, 16}, {-84, 10}, {-80, 8}, {-78, 8}, {-84, 14}, {-90, 14}, {-95, 16}, {-105, 20}, {-110, 24}, {-112, 30}, {-118, 34}, {-124, 40}, {-124, 48}, {-133, 56}, {-140, 60}, {-150, 61}, {-158, 58}, {-165, 60}},
	{{-55, 60}, {-43, 60}, {-20, 70}, {-18, 80}, {-35, 83}, {-60, 82}, {-70, 77}, {-55, 70}},
	{{-24, 64}, {-14, 64}, {-14, 66}, {-22, 66}},
	{{-80, 8}, {-75, 11}, {-62, 11}, {-52, 5}, {-35, -5}, {-35, -10}, {-40, -22}, {-48, -26}, {-58, -35}, {-65, -42}, {-68, -52}, {-72, -54}, {-75, -48}, {-73, -38}, {-71, -30}, {-70, -18}, {-76, -14}, {-81, -5}, {-80, 0}, {-78, 7}},
	{{-10, 36}, {-9, 43}, {-2, 44}, {-5, 48}, {2, 51}, {8, 54}, {8, 57}, {5, 62}, {12, 66}, {20, 70}, {30, 71}, {40, 68}, {45, 68}, {60, 70}, {60, 55}, {50, 45}, {40, 42}, {30, 41}, {26, 40}, {23, 36}, {20, 40}, {15, 38}, {12, 42}, {8, 44}, {3, 43}, {-1, 37}},
	{{-6, 50}, {2, 51}, {0, 53}, {-2, 57}, {-5, 59}, {-6, 56}, {-5, 54}, {-3, 52}},
	{{-10, 52}, {-6, 52}, {-6, 55}, {-10, 55}},
	{{-17, 21}, {-17, 14}, {-12, 8}, {-8, 4}, {0, 5}, {9, 4}, {10, -1}, {13, -10}, {12, -17}, {15, -27}, {18, -34}, {25, -34}, {33, -27}, {35, -20}, {40, -15}, {40, -5}, {44, 0}, {51, 11}, {43, 12}, {37, 20}, {33, 28}, {30, 31}, {20, 31}, {10, 33}, {10, 37}, {0, 36}, {-6, 36}, {-10, 30}},
	{{44, -25}, {47, -25}, {50, -15}, {49, -12}, {44, -17}},
	{{26, 40}, {36, 36}, {35, 32}, {43, 13}, {52, 16}, {57, 22}, {57, 26}, {62, 25}, {67, 25}, {73, 20}, {77, 8}, {80, 15}, {80, 20}, {87, 22}, {92, 22}, {98, 16}, {100, 13}, {104, 10}, {100, 2}, {104, 1}, {106, 10}, {109, 12}, {108, 21}, {117, 23}, {122, 30}, {120, 38}, {126, 38}, {129, 35}, {130, 42}, {140, 48}, {141, 53}, {135, 55}, {143, 59}, {155, 59}, {163, 60}, {170, 60}, {180, 65}, {180, 70}, {160, 70}, {140, 72}, {130, 72}, {115, 73}, {105, 78}, {90, 76}, {75, 73}, {70, 73}, {60, 70}, {60, 55}, {50, 45}, {40, 42}},
	{{130, 31}, {132, 34}, {137, 35}, {140, 36}, {141, 40}, {142, 45}, {144, 43}, {140, 41}, {139, 38}, {136, 36}, {131, 34}},
	{{95, 5}, {98, 4}, {106, -6}, {104, -5}},
	{{109, 1}, {117, 7}, {119, 1}, {116, -4}, {110, -3}},
	{{105, -6}, {114, -7}, {114, -8}, {106, -7}},
	{{131, -1}, {141, -3}, {150, -10}, {141, -9}, {137, -5}},
	{{114, -22}, {114, -34}, {118, -35}, {124, -33}, {131, -31}, {138, -35}, {141, -38}, {146, -39}, {150, -37}, {153, -28}, {153, -25}, {146, -19}, {142, -11}, {136, -12}, {136, -15}, {130, -12}, {126, -14}, {122, -17}},
	{{172, -34}, {178, -38}, {174, -41}, {171, -46}, {167, -46}, {172, -41}},
}

// ordersGlobeIsLand tells whether a point lies inside any outline.
func ordersGlobeIsLand(lat, lon float64) bool {
	if lat < -68 {
		return true
	}
	for _, ring := range ordersGlobeLand {
		in := false
		for i, j := 0, len(ring)-1; i < len(ring); j, i = i, i+1 {
			a, b := ring[i], ring[j]
			if (a[1] > lat) != (b[1] > lat) && lon < (b[0]-a[0])*(lat-a[1])/(b[1]-a[1])+a[0] {
				in = !in
			}
		}
		if in {
			return true
		}
	}
	return false
}

// ordersGlobeLattice lays dots in rings of latitude, evenly round each.
func ordersGlobeLattice() []ordersGlobeDot {
	ordersGlobeDotsOnce.Do(func() {
		for lat := -90 + ordersGlobeStep/2; lat < 90; lat += ordersGlobeStep {
			n := int(math.Round(360 * math.Cos(lat*math.Pi/180) / ordersGlobeStep))
			for i := 0; i < n; i++ {
				lon := -180 + 360*float64(i)/float64(n)
				ordersGlobeDots = append(ordersGlobeDots, ordersGlobeDot{lat, lon, ordersGlobeIsLand(lat, lon)})
			}
		}
	})
	return ordersGlobeDots
}

// ordersGlobeProject puts a point on the globe facing the viewer at lon0:
// across and up from the centre, in radii, and how far toward the viewer.
func ordersGlobeProject(lat, lon, lon0 float64) (x, y, z float64) {
	p, d := lat*math.Pi/180, (lon-lon0)*math.Pi/180
	x = math.Cos(p) * math.Sin(d)
	y0, z0 := math.Sin(p), math.Cos(p)*math.Cos(d)
	t := ordersGlobeTilt * math.Pi / 180
	return x, y0*math.Cos(t) - z0*math.Sin(t), y0*math.Sin(t) + z0*math.Cos(t)
}

// OrdersGlobe is a card of live orders: a dotted globe turning on its own,
// each city a dot sized by its share of orders, the total, and the top
// three cities as dot bars. Orders land now and then, the city's dot
// rippling and the counts going up; dragging spins the globe, and the
// pill switches between today, this week and this month.
func OrdersGlobe(c *ui.Context, opts ...OrdersGlobeOptions) ui.Element {
	var o OrdersGlobeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Rate == 0 {
		o.Rate = 40
	}
	if o.NoOrders {
		o.Rate = 0
	}
	o.Rate = clamp32(o.Rate, 0, 100)
	if o.Speed == 0 {
		o.Speed = 30
	}
	o.Speed = clamp32(o.Speed, 0, 100)
	t := ThemeFrom(c)
	card := ui.Box(c).Size(ordersGlobeW, ordersGlobeH).Background(t.Surface).Radius(22).Border(1, t.Border)
	st := ui.Local(card, "orders_globe", func() ordersGlobeState {
		s := ordersGlobeState{lon: -40, rng: rand.New(rand.NewSource(11)), landed: make([]time.Time, len(ordersGlobeCities))}
		for p, per := range ordersGlobePeriods {
			for _, city := range ordersGlobeCities {
				s.counts[p] = append(s.counts[p], int(math.Round(float64(city.today)*per.scale)))
			}
		}
		return s
	})
	now := c.Now()
	dt := 0.0
	if !st.last.IsZero() {
		dt = min(now.Sub(st.last).Seconds(), 0.05)
	}
	st.last = now

	// Orders land now and then, more often in the busiest cities.
	if o.Rate > 0 {
		gap := time.Duration(60 / float64(o.Rate) * float64(time.Second))
		if st.next.IsZero() {
			st.next = now.Add(gap)
		}
		if !now.Before(st.next) {
			counts := st.counts[0]
			sum := 0
			for _, n := range counts {
				sum += n
			}
			pick, k := st.rng.Intn(sum), 0
			for k = range counts {
				if pick -= counts[k]; pick < 0 {
					break
				}
			}
			for p := range st.counts {
				st.counts[p][k]++
			}
			st.landed[k] = now
			st.next = now.Add(time.Duration(float64(gap) * (0.6 + 0.8*st.rng.Float64())))
		}
	}

	c.AnimationFrame()

	counts := st.counts[st.period]
	order := make([]int, len(counts))
	total := 0
	for i := range order {
		order[i] = i
		total += counts[i]
	}
	sort.SliceStable(order, func(a, b int) bool { return counts[order[a]] > counts[order[b]] })
	hue := func(k int) ui.Color {
		for rank := 0; rank < 3; rank++ {
			if order[rank] == k && o.Tone != OrdersGlobeMono {
				return ordersGlobeHues[rank]
			}
		}
		return t.Text
	}
	per := ordersGlobePeriods[st.period]

	card.Children(func() {
		ui.Text(c, "Orders").Key("title").FontSize(16).FontWeight(500).TextColor(t.Text).
			Absolute().Left(17).Top(22)
		pill := ui.Box(c).Key("pill").Absolute().Right(16).Top(17).Height(30).Padding(0, 12).Radius(15).
			Background(t.SurfaceSubtle).Center().Cursor(ui.CursorPointer)
		if pill.Hovered() {
			pill.Background(t.SurfaceHover)
		}
		pill.OnClick(func() { st.period = (st.period + 1) % len(ordersGlobePeriods) })
		pill.Children(func() {
			ui.Text(c, per.label).FontSize(13).FontWeight(500).TextColor(t.TextMuted)
		})
		// The globe turns on its own; a drag spins it, and a flick keeps it
		// going a while.
		globe := ui.Box(c).Key("globe").Absolute().Left(ordersGlobeCX-ordersGlobeR).Top(ordersGlobeCY-ordersGlobeR).
			Size(2*ordersGlobeR, 2*ordersGlobeR).Cursor(ui.CursorGrab)
		dx, _, pressed := globe.Dragged()
		if pressed {
			globe.Cursor(ui.CursorGrabbing)
			st.lon -= float64(dx) * 0.55
			vx, _ := globe.DragVelocity()
			st.spin = -float64(vx) * 0.55
		} else {
			st.lon += (st.spin - float64(o.Speed)*0.2) * dt // the surface goes east, as the Earth turns
			st.spin *= math.Exp(-2.2 * dt)
		}
		st.lon = math.Mod(st.lon+540, 360) - 180
		globe.Draw(func(p *ui.Painter, r ui.Rect) {
			ordersGlobeDraw(p, t, r, st, counts, total, hue, now)
		})
		ui.Row(c).Key("total").Absolute().Left(15).Top(319).AlignItems(ui.End).Children(func() {
			ui.Text(c, strconv.Itoa(total)).FontSize(31).FontWeight(500).TextColor(t.Text).Margin(0, 8, 0, 0)
			ui.Text(c, per.change).FontSize(13.5).FontWeight(500).TextColor(ui.Hex("#16a34a")).Margin(0, 5, 5, 0)
			ui.Text(c, per.versus).FontSize(13.5).TextColor(t.TextSubtle).Margin(0, 0, 5, 0)
		})
		top := counts[order[0]]
		for rank := 0; rank < 3; rank++ {
			k := order[rank]
			n := counts[k]
			y := float32(368 + rank*27)
			ui.Text(c, ordersGlobeCities[k].name).Key(fmt.Sprint("name", rank)).FontSize(14).TextColor(t.TextMuted).
				Absolute().Left(17).Top(y)
			ui.Text(c, strconv.Itoa(n)).Key(fmt.Sprint("n", rank)).FontSize(14).TextColor(t.Text).
				Absolute().Left(230).Width(50).TextAlign(ui.End).Top(y)
			lit := int(math.Round(float64(n) / float64(top) * ordersGlobeCols))
			col := hue(k)
			ui.Box(c).Key(fmt.Sprint("bar", rank)).Absolute().Left(101).Top(y+1.5).Size(130, 14).Draw(func(p *ui.Painter, r ui.Rect) {
				var on, off ui.Path
				for i := 0; i < ordersGlobeCols; i++ {
					for j := 0; j < 3; j++ {
						path := &on
						if i >= lit {
							path = &off
						}
						x, y := r.X+float32(i)*130/ordersGlobeCols, r.Y+float32(j)*4.9
						path.MoveTo(x, y).LineTo(x+3, y).LineTo(x+3, y+3).LineTo(x, y+3).Close()
					}
				}
				p.FillPath(&on, col)
				p.FillPath(&off, t.Border)
			})
		}
	})
	return card
}

// ordersGlobeDraw draws the lattice, land dark and sea faint, fading at the
// rim, then each city's dot sized by its share, with a ripple where an
// order just landed.
func ordersGlobeDraw(p *ui.Painter, t Theme, r ui.Rect, st *ordersGlobeState, counts []int, total int, hue func(int) ui.Color, now time.Time) {
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	const shades = 3
	var land, sea [shades]ui.Path
	for _, d := range ordersGlobeLattice() {
		x, y, z := ordersGlobeProject(d.lat, d.lon, st.lon)
		if z <= 0.02 {
			continue
		}
		sx, sy := cx+float32(x)*ordersGlobeR, cy-float32(y)*ordersGlobeR
		k := min(int(z*shades*1.6), shades-1)
		if d.land {
			s := float32(0.55 + 0.45*z)
			land[k].MoveTo(sx-s, sy-s).LineTo(sx+s, sy-s).LineTo(sx+s, sy+s).LineTo(sx-s, sy+s).Close()
		} else {
			sea[k].Circle(sx, sy, 0.75)
		}
	}
	ink, faint := ui.Hex("#2a2a2a"), ui.Hex("#d9d9d9")
	if t.Dark {
		ink, faint = ui.Hex("#d4d4d8"), ui.Hex("#3a3a3e")
	}
	for k := 0; k < shades; k++ {
		a := 0.35 + 0.65*float32(k)/(shades-1)
		p.FillPath(&sea[k], faint.Alpha(a))
		p.FillPath(&land[k], ink.Alpha(a))
	}
	for k, city := range ordersGlobeCities {
		x, y, z := ordersGlobeProject(city.lat, city.lon, st.lon)
		if z <= 0.05 {
			continue
		}
		sx, sy := cx+float32(x)*ordersGlobeR, cy-float32(y)*ordersGlobeR
		col := hue(k)
		rad := 1.8 + 22*float32(counts[k])/float32(total)
		if age := float32(now.Sub(st.landed[k]).Seconds()); !st.landed[k].IsZero() && age < 0.9 {
			// A ring goes out from the city, and its dot swells a moment.
			rr := rad + 14*age/0.9
			p.Stroke(ui.Rect{X: sx - rr, Y: sy - rr, W: 2 * rr, H: 2 * rr}, col.Alpha(0.6*(1-age/0.9)), rr, 1.5)
			rad *= 1 + 0.35*float32(math.Sin(float64(min(age/0.35, 1))*math.Pi))
		}
		p.Fill(ui.Rect{X: sx - rad, Y: sy - rad, W: 2 * rad, H: 2 * rad}, col.Alpha(clamp32(float32(z)*3, 0, 1)), rad)
	}
}

func init() {
	registerDemo(
		"orders-globe",
		"A dotted globe of live orders: cities sized by share, drag to spin it, and a pill for today, this week or this month.",
		func(c *ui.Context) { OrdersGlobe(c) },
	)
}
