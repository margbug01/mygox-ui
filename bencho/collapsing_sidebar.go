package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// CollapsingSidebarItem is one entry in CollapsingSidebar's navigation.
type CollapsingSidebarItem struct {
	Label string
	Icon  *ui.SVG
	Badge string // a red count after the label; a dot on the icon when folded
}

// CollapsingSidebarOptions configures CollapsingSidebar.
type CollapsingSidebarOptions struct {
	Corner float32 // 0 to 28; 20 by default
	Bounce float32 // how far the panel overshoots as it opens, 0 to 100; 70 by default
	Speed  float32 // 0.5 to 2; 1 by default
	Title  string  // the workspace name beside the avatar; "Bencho" by default
}

type collapsingSidebarState struct {
	folded bool
	active int
	words  float32 // last frame's label opacity
	width  float32 // last frame's 0 (folded) to 1 (open)
}

const (
	collapsingSidebarOpenW   = 220
	collapsingSidebarFoldedW = 56
	collapsingSidebarRow     = 36
	collapsingSidebarGap     = 40
	collapsingSidebarIconX   = 28 // icon centers stay on this column as the panel folds
	collapsingSidebarLabelX  = 50
)

// Outline glyphs as bencho.dev draws them (a plain calendar, a toggle for
// settings, a panel with a bar for collapse), in Lucide's stroke style.
var (
	collapsingSidebarCalendar = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4.5" width="18" height="16" rx="3.5"/><path d="M3 9.5h18"/></svg>`))
	collapsingSidebarToggle   = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="2.5" y="6.5" width="19" height="11" rx="5.5"/><circle cx="8" cy="12" r="2.3"/></svg>`))
	collapsingSidebarPanel    = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="3.5" y="3.5" width="17" height="17" rx="3.5"/><path d="M9 9v6"/></svg>`))
)

// CollapsingSidebarDefaultItems is the navigation bencho.dev shows.
func CollapsingSidebarDefaultItems() []CollapsingSidebarItem {
	return []CollapsingSidebarItem{
		{Label: "Home", Icon: icons.Home()},
		{Label: "Inbox", Icon: icons.Inbox(), Badge: "3"},
		{Label: "Projects", Icon: icons.Folder()},
		{Label: "Calendar", Icon: collapsingSidebarCalendar},
		{Label: "Settings", Icon: collapsingSidebarToggle},
	}
}

// CollapsingSidebar is a navigation panel that folds down to its icons.
// Pressing Collapse fades the words out first, then the panel closes over an
// icon column that stays still while the highlight settles as a square;
// opening springs the panel back past its width before the words return.
func CollapsingSidebar(c *ui.Context, items []CollapsingSidebarItem, opts ...CollapsingSidebarOptions) ui.Element {
	var o CollapsingSidebarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	if o.Bounce == 0 {
		o.Bounce = 70
	}
	if o.Speed == 0 {
		o.Speed = 1
	}
	o.Title = orDefault(o.Title, "Bencho")
	t := ThemeFrom(c)
	dur := func(ms float32) time.Duration { return time.Duration(ms/o.Speed) * time.Millisecond }

	listH := float32(len(items)) * collapsingSidebarGap
	h := 60 + listH + 6 + 1 + 52
	panel := ui.Box(c).Height(h).Radius(o.Corner).Background(t.Surface).Clip()
	st := ui.Local(panel, "collapsing_sidebar", func() collapsingSidebarState {
		return collapsingSidebarState{words: 1, width: 1}
	})

	// Folding: words go first, then the panel. Opening: the panel, then words.
	wordsTarget := b2f(!st.folded && st.width > 0.75)
	widthTarget := b2f(!st.folded || st.words > 0.15)
	words := clamp32(panel.Spring("words", wordsTarget, ui.SpringOptions{Duration: dur(180)}), 0, 1)
	width := panel.Spring("width", widthTarget, ui.SpringOptions{Duration: dur(520), Bounce: float64(clamp32(o.Bounce, 0, 100)) / 100 * 0.4 * float64(widthTarget)})
	st.words, st.width = words, width
	panel.Width(lerp32(collapsingSidebarFoldedW, collapsingSidebarOpenW, max(width, 0)))

	open := clamp32(width, 0, 1)
	hiW := lerp32(collapsingSidebarRow, collapsingSidebarOpenW-20, open)
	hiY := panel.Spring("active", float32(st.active), ui.SpringOptions{Duration: dur(380), Bounce: 0.15})
	label := func(s string, col ui.Color, weight int) {
		ui.Text(c, s).FontSize(15).FontWeight(weight).TextColor(col).SingleLine().Opacity(words)
	}

	panel.Children(func() {
		// Highlight under the active row; a square once folded.
		ui.Box(c).Absolute().Left(10).Top(60+hiY*collapsingSidebarGap).Size(hiW, collapsingSidebarRow).
			Radius(min(10, o.Corner)).Background(lerpColor(t.SurfaceSubtle, t.SurfaceHover, 0.25))

		// Workspace header.
		ui.Box(c).Absolute().Left(collapsingSidebarIconX-12).Top(17).Size(24, 24).Draw(func(p *ui.Painter, r ui.Rect) {
			p.FillGradient(r, ui.LinearGradient{From: ui.Hex("#f7866f"), To: ui.Hex("#e8c544"), Angle: 135}, 12)
		})
		ui.Row(c).Absolute().Left(collapsingSidebarLabelX).Top(17).Height(24).AlignItems(ui.Center).Children(func() {
			label(o.Title, t.Text, 500)
		})

		for i, it := range items {
			on := i == st.active
			row := ui.Box(c).Key(i).Absolute().Left(10).Top(60+float32(i)*collapsingSidebarGap).Size(hiW, collapsingSidebarRow).
				Radius(min(10, o.Corner)).Cursor(ui.CursorPointer).OnClick(func() { st.active = i })
			hover := row.Spring("hover", b2f(row.Hovered() && !on), ui.SpringOptions{Duration: dur(160)})
			col, weight := t.TextSubtle, 400
			if on {
				col, weight = t.Text, 500
			}
			row.Background(t.SurfaceSubtle.Alpha(hover * 0.7))
			row.Children(func() {
				ui.Icon(c, it.Icon).Absolute().Left(collapsingSidebarIconX-10-9).Top(9).Size(18, 18).TextColor(col)
				ui.Row(c).Absolute().Left(collapsingSidebarLabelX - 10).Top(0).Height(collapsingSidebarRow).AlignItems(ui.Center).Children(func() {
					label(it.Label, col, weight)
					if it.Badge != "" {
						ui.Box(c).Size(18, 18).Margin(0, 0, 0, 8).Radius(9).Background(t.Danger).Center().Opacity(words).Children(func() {
							ui.Text(c, it.Badge).FontSize(11).FontWeight(600).TextColor(ui.Hex("#ffffff"))
						})
					}
				})
				if it.Badge != "" && words < 1 {
					// Folded, the count becomes a dot on the icon's corner.
					ui.Box(c).Absolute().Left(collapsingSidebarIconX-10+6).Top(7).Size(7, 7).Radius(3.5).
						Background(t.Danger).Opacity(1-words).Border(1.5, t.Surface)
				}
			})
		}

		top := 60 + listH + 6
		ui.Box(c).Absolute().Left(0).Top(top).Size(collapsingSidebarOpenW+40, 1).Background(t.Border)
		toggle := ui.Box(c).Key("collapse").Absolute().Left(10).Top(top+9).Size(hiW, collapsingSidebarRow).
			Radius(min(10, o.Corner)).Cursor(ui.CursorPointer).OnClick(func() { st.folded = !st.folded })
		hover := toggle.Spring("hover", b2f(toggle.Hovered()), ui.SpringOptions{Duration: dur(160)})
		toggle.Background(t.SurfaceSubtle.Alpha(hover * 0.7))
		toggle.Children(func() {
			ui.Icon(c, collapsingSidebarPanel).Absolute().Left(collapsingSidebarIconX-10-9).Top(9).Size(18, 18).TextColor(t.TextSubtle).
				Rotate(180 * (1 - open))
			ui.Row(c).Absolute().Left(collapsingSidebarLabelX - 10).Top(0).Height(collapsingSidebarRow).AlignItems(ui.Center).Children(func() {
				label("Collapse", t.TextSubtle, 400)
			})
		})
	})
	return panel
}

func init() {
	registerDemo(
		"collapsing-sidebar",
		"A navigation panel that folds down to its icons: the words leave first, then the panel closes over a still icon column.",
		func(c *ui.Context) { CollapsingSidebar(c, CollapsingSidebarDefaultItems()) },
	)
}
