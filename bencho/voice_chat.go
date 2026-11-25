package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// VoiceChatOptions configures VoiceChat.
type VoiceChatOptions struct {
	Speed  float32 // how quickly the pill opens into the room, 0 to 100; 50 by default
	Corner float32 // the room's corner radius, 0 to 32; 20 by default
}

// voiceChatPeople are the seven in the room; the first three show on the
// pill. Each orb blends two colours, top left to bottom right.
var voiceChatPeople = []struct{ name, from, to string }{
	{"Mara", "#f2416a", "#f2804c"},
	{"Kai", "#fe6ec7", "#9b59e8"},
	{"Rosa", "#27b4fd", "#5a4ff3"},
	{"Juno", "#8cf05a", "#2be6c6"},
	{"Theo", "#fbbf24", "#f97316"},
	{"Lena", "#22d3ee", "#0ea5e9"},
	{"Ivo", "#c084fc", "#7c3aed"},
}

const (
	voiceChatW, voiceChatH = 300, 330
	voiceChatTurn          = 1400 * time.Millisecond
)

var (
	voiceChatPill = ui.Rect{X: 59, Y: 139, W: 182, H: 52}
	voiceChatRoom = ui.Rect{X: 0, Y: 0, W: voiceChatW, H: voiceChatH}
)

type voiceChatState struct {
	open, muted bool
}

// VoiceChat is a small pill of faces that opens into the full voice room:
// press it and the faces travel into a grid of everyone there, whoever is
// talking ringed and their bars moving. Leave folds the room back into the
// pill.
func VoiceChat(c *ui.Context, opts ...VoiceChatOptions) ui.Element {
	var o VoiceChatOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Speed == 0 {
		o.Speed = 50
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Speed, o.Corner = clamp32(o.Speed, 0, 100), clamp32(o.Corner, 0, 32)
	t := ThemeFrom(c)
	now := c.Now()
	c.AnimationFrame() // someone is always talking

	root := ui.Box(c).Size(voiceChatW, voiceChatH)
	st := ui.Local(root, "voice_chat", func() voiceChatState { return voiceChatState{} })
	dur := time.Duration(lerp32(1000, 260, o.Speed/100)) * time.Millisecond
	k := root.Spring("open", b2f(st.open), ui.SpringOptions{Duration: dur, Bounce: 0.18})
	kc := clamp32(k, 0, 1)
	speaker := int(phaseSeconds(now)/voiceChatTurn.Seconds()) % len(voiceChatPeople)

	frame := ui.Rect{
		X: lerp32(voiceChatPill.X, voiceChatRoom.X, k), Y: lerp32(voiceChatPill.Y, voiceChatRoom.Y, k),
		W: lerp32(voiceChatPill.W, voiceChatRoom.W, k), H: lerp32(voiceChatPill.H, voiceChatRoom.H, k),
	}
	root.Children(func() {
		shell := ui.Box(c).Absolute().Left(frame.X).Top(frame.Y).Size(frame.W, frame.H).
			Radius(lerp32(voiceChatPill.H/2, o.Corner, kc)).Background(t.Surface).Clip()
		if !st.open {
			shell.Cursor(ui.CursorPointer).OnClick(func() { st.open = true })
		}
		press := shell.Spring("press", b2f(shell.Pressed() && !st.open), ui.SpringOptions{Duration: 140 * time.Millisecond})
		shell.Scale(1-0.04*press, 1-0.04*press)
		shell.Children(func() {
			// The pill's own parts fade as it opens: the +4 and the bars.
			if kc < 0.5 {
				fade := 1 - kc*2
				ui.Box(c).Absolute().Left(104).Top(6).Size(40, 40).Radius(20).Background(t.SurfaceSubtle).Center().Opacity(fade).
					Children(func() { ui.Text(c, "+4").FontSize(15).FontWeight(500).TextColor(t.Text) })
				ui.Box(c).Absolute().Left(151).Top(16).Size(14, 20).Opacity(fade).Draw(func(p *ui.Painter, r ui.Rect) {
					voiceChatBars(p, r, t.Text, now, 1)
				})
			}
			if kc > 0.3 {
				voiceChatRoomChrome(c, t, st, clamp32((kc-0.3)/0.6, 0, 1))
			}
			// Faces travel from the pill into the room's grid; the rest
			// of the room comes in around them.
			for i := range voiceChatPeople { // later faces sit on top, as on the pill
				pp := voiceChatPeople[i]
				gx, gy := voiceChatSlot(i)
				px, py, show := float32(6+i*32), float32(6), float32(1)
				if i >= 3 {
					px, py, show = gx, gy, kc
				}
				x, y := lerp32(px, gx, k), lerp32(py, gy, k)
				d := lerp32(40, 52, kc)
				if show < 0.02 {
					continue
				}
				talking := i == speaker && kc > 0.9
				ui.Box(c).Key(i).Absolute().Left(x-2).Top(y-2).Size(d+4, d+4).Opacity(show).Draw(func(p *ui.Painter, r ui.Rect) {
					if talking {
						s := float32(mod1(phaseSeconds(now) / 1.1))
						grow := 2 + s*6
						p.Stroke(ui.Rect{X: r.X - grow, Y: r.Y - grow, W: r.W + 2*grow, H: r.H + 2*grow}, ui.Hex(pp.from).Alpha(0.5*(1-s)), r.W/2+grow, 2)
					}
					p.Fill(r, t.Surface, r.W/2)
					in := ui.Rect{X: r.X + 2, Y: r.Y + 2, W: d, H: d}
					p.FillGradient(in, ui.LinearGradient{Angle: 135, From: ui.Hex(pp.from), To: ui.Hex(pp.to)}, d/2)
				})
				if kc > 0.6 {
					col := t.TextSubtle
					if talking {
						col = t.Text
					}
					ui.Text(c, pp.name).Key(100 + i).FontSize(12).FontWeight(500).TextColor(col).Opacity((kc - 0.6) / 0.4).
						Absolute().Left(gx - 10).Top(gy + 58).Width(72).TextAlign(ui.Center)
				}
			}
		})
	})
	return root
}

// voiceChatSlot is where person i sits in the room: four to a row.
func voiceChatSlot(i int) (x, y float32) {
	return 26 + float32(i%4)*65, 84 + float32(i/4)*92
}

// voiceChatRoomChrome is the room's header and controls, shown by a.
func voiceChatRoomChrome(c *ui.Context, t Theme, st *voiceChatState, a float32) {
	ui.Text(c, "Design crit").FontSize(16).FontWeight(600).TextColor(t.Text).Opacity(a).Absolute().Left(20).Top(20)
	ui.Text(c, "7 in the room").FontSize(13).TextColor(t.TextSubtle).Opacity(a).Absolute().Left(20).Top(42)
	ui.Box(c).Absolute().Left(248).Top(22).Size(8, 8).Radius(5).Background(ui.Hex("#30a46c")).Opacity(a)
	ui.Text(c, "Live").FontSize(12).FontWeight(500).TextColor(t.TextSubtle).Opacity(a).Absolute().Left(260).Top(18)
	// The eighth seat invites someone in.
	gx, gy := voiceChatSlot(7)
	ui.Box(c).Absolute().Left(gx).Top(gy).Size(52, 52).Radius(26).Opacity(a).Center().
		Draw(func(p *ui.Painter, r ui.Rect) { p.StrokeDashed(r, t.BorderStrong, 26, 1.5) }).
		Children(func() { ui.Icon(c, icons.Plus()).Size(18, 18).TextColor(t.TextSubtle) })
	ui.Text(c, "Invite").FontSize(12).FontWeight(500).TextColor(t.TextSubtle).Opacity(a).
		Absolute().Left(gx - 10).Top(gy + 58).Width(72).TextAlign(ui.Center)

	mic := ui.Box(c).Absolute().Left(20).Top(266).Size(44, 44).Radius(22).Background(t.SurfaceSubtle).Center().Opacity(a).
		Cursor(ui.CursorPointer).OnClick(func() { st.muted = !st.muted })
	muted := st.muted
	mic.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if muted {
			p.Line(r.X+14, r.Y+14, r.X+30, r.Y+30, 1.8, t.Text)
		}
	})
	mic.Children(func() { ui.Icon(c, icons.Mic()).Size(18, 18).TextColor(t.Text) })
	leave := ui.Box(c).Absolute().Left(76).Top(266).Size(204, 44).Radius(22).Background(ui.Hex("#e5484d")).Center().Opacity(a).
		Cursor(ui.CursorPointer).OnClick(func() { st.open = false })
	lp := leave.Spring("press", b2f(leave.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
	leave.Scale(1-0.03*lp, 1-0.03*lp)
	leave.Children(func() { ui.Text(c, "Leave").FontSize(14).FontWeight(500).TextColor(ui.Hex("#ffffff")) })
}

// voiceChatBars draws three sound bars bouncing in r.
func voiceChatBars(p *ui.Painter, r ui.Rect, col ui.Color, now time.Time, a float32) {
	p.AnimationFrame()
	s := phaseSeconds(now)
	for b := range 3 {
		h := r.H * float32(0.55+0.45*math.Sin(s*7+float64(b)*1.9))
		h = max(h, 6)
		p.Fill(ui.Rect{X: r.X + float32(b)*5, Y: r.Y + r.H/2 - h/2, W: 4, H: h}, col.Alpha(a), 2)
	}
}

func init() {
	registerDemo(
		"voice-chat",
		"A voice chat pill of faces that opens into the full voice room, showing who is talking; Leave folds it back.",
		func(c *ui.Context) { VoiceChat(c) },
	)
}
