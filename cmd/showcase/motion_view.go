package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
	"github.com/margbug01/mygox-ui/motion"
)

type motionState struct {
	selectedCat   string
	query         string
	openCode      map[string]bool
	copiedUntil   map[string]time.Time
	sidebarScroll ui.ScrollState
	contentScroll ui.ScrollState
	catalog       []motionDemoEntry

	// Kept for showcase test compatibility
	confettiTrigger int

	// Interactive states for all 33 motion-primitives demos
	animNumVal       float64
	slidingCount     int
	textEffectTrig   int
	textEffectPreset motion.Preset
	textMorphConfirm bool
	textRollTrig     int
	textScrambleTrig int

	glowMode   motion.GlowMode
	comparePos float32

	animBgTab       int
	dockActive      int
	magneticClicks  int
	toolbarDynOpen  bool
	toolbarDynQuery string
	toolbarExpIdx   int

	accordionOpen   string
	disclosureOpen  bool
	dialogOpen      bool
	dialogEmail     string
	morphDialogOpen bool
	morphPopOpen    bool
	morphPopNote    string
	transPanelTab   int

	groupTrigger   int
	groupPreset    motion.Preset
	carouselIdx    int
	inViewTrigger  int
	scrollDemoStat ui.ScrollState
}

func newMotionState(cat string) *motionState {
	valid := false
	for _, info := range motionCategories {
		if info.ID == cat {
			valid = true
			break
		}
	}
	if !valid {
		cat = ""
	}
	return &motionState{
		selectedCat:      cat,
		openCode:         map[string]bool{},
		copiedUntil:      map[string]time.Time{},
		catalog:          buildMotionCatalog(),
		animNumVal:       1000,
		slidingCount:     42,
		textEffectPreset: motion.PresetBlurSlide,
		glowMode:         motion.GlowRotate,
		comparePos:       50,
		animBgTab:        0,
		dockActive:       0,
		toolbarExpIdx:    -1,
		accordionOpen:    "getting-started",
		groupPreset:      motion.PresetScale,
	}
}

func (m *motionState) countByCategory(cat string) int {
	if cat == "" {
		return len(m.catalog)
	}
	n := 0
	for _, e := range m.catalog {
		if e.Category == cat {
			n++
		}
	}
	return n
}

func (m *motionState) filtered() []motionDemoEntry {
	q := strings.ToLower(strings.TrimSpace(m.query))
	out := make([]motionDemoEntry, 0, len(m.catalog))
	for _, e := range m.catalog {
		if m.selectedCat != "" && e.Category != m.selectedCat {
			continue
		}
		if q != "" {
			if !strings.Contains(strings.ToLower(e.Name), q) &&
				!strings.Contains(strings.ToLower(e.Description), q) &&
				!strings.Contains(strings.ToLower(e.Category), q) {
				continue
			}
		}
		out = append(out, e)
	}
	return out
}

func (m *motionState) selectCategory(cat string) {
	m.selectedCat = cat
	m.contentScroll.Y = 0
}

func (a *showcaseApp) renderMotionTab(c *ui.Context, t theme.Tokens, isMobile bool) {
	m := a.motion
	mt := motion.ApplyTheme(c, a.dark)
	// Keep motion surface & border tokens aligned with coss neutral zinc theme.
	mt.Background = t.Background
	mt.Surface = t.Card
	mt.SurfaceSubtle = t.Sidebar
	mt.SurfaceHover = t.Accent
	mt.Border = t.Border
	mt.Text = t.Foreground
	mt.TextMuted = t.MutedForeground

	visible := m.filtered()

	ui.Row(c).WidthPercent(100).Grow(1).MinHeight(0).Children(func() {
		// Left Coss-style Sidebar
		if !isMobile {
			ui.Column(c).
				Width(256).
				HeightPercent(100).
				Shrink(0).
				Background(t.Sidebar).
				BorderWidth(0, 1, 0, 0).
				BorderColor(t.SidebarBorder).
				Padding(14, 0, 0, 0).
				Gap(12).
				Children(func() {
					ui.Box(c).WidthPercent(100).Padding(0, 12).Children(func() {
						a.motionSearchInput(c)
					})
					a.motionCategoryNav(c, t, false)
				})
		}

		// Right Main Content (coss card-grid visual grammar)
		ui.Column(c).Grow(1).MinWidth(0).HeightPercent(100).Children(func() {
			title := "All Motion Effects"
			for _, cat := range motionCategories {
				if cat.ID == m.selectedCat {
					title = cat.Label
					break
				}
			}

			if isMobile {
				ui.Column(c).
					WidthPercent(100).
					Background(t.Sidebar).
					BorderWidth(0, 0, 1, 0).
					BorderColor(t.Border).
					Padding(10, 12).
					Gap(8).
					Children(func() {
						ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(8).Children(func() {
							ui.Row(c).AlignItems(ui.Center).Gap(8).Grow(1).MinWidth(0).Children(func() {
								coss.IconButton(c, icons.PanelLeft(), "Open categories").
									Size(coss.SizeIconSM).
									Variant(coss.ButtonOutline).
									OnClick(func() { a.mobileNavOpen = true })
								ui.Text(c, title).FontSize(14.5).FontWeight(600).TextColor(t.Foreground).SingleLine().Ellipsis("…")
								coss.Badge(c, fmt.Sprintf("%d", len(visible))).Variant(coss.BadgeSecondary).Size(coss.BadgeSizeSM)
							})
							coss.Badge(c, "33 primitives").Variant(coss.BadgeOutline).Size(coss.BadgeSizeSM)
						})
						a.motionSearchInput(c)
					})
			} else {
				ui.Row(c).
					WidthPercent(100).
					AlignItems(ui.Center).
					Justify(ui.SpaceBetween).
					Padding(16, 24).
					BorderWidth(0, 0, 1, 0).
					BorderColor(t.Border).
					Children(func() {
						ui.Column(c).Gap(2).Children(func() {
							ui.Text(c, title).FontSize(18).FontWeight(600).TextColor(t.Foreground)
							ui.Text(c, fmt.Sprintf("Showing %d of %d motion-primitives components in monochrome coss style", len(visible), len(m.catalog))).
								FontSize(13).TextColor(t.MutedForeground)
						})
						coss.Badge(c, "mygox-ui/motion").Variant(coss.BadgeOutline)
					})
			}

			ui.Scroll(c).TrackScroll(&m.contentScroll).WidthPercent(100).Grow(1).MinHeight(0).Children(func() {
				pad := float32(24)
				gap := float32(16)
				if isMobile {
					pad = 12
					gap = 12
				}
				ui.Column(c).WidthPercent(100).Padding(pad).Gap(gap).Children(func() {
					if len(visible) == 0 {
						coss.Empty(c, func() {
							coss.EmptyHeader(c, func() {
								coss.EmptyMedia(c, icons.Search())
								coss.EmptyTitle(c, "No matching motion primitives")
								coss.EmptyDescription(c, "Try clearing your search query or choosing another category.")
							})
							coss.EmptyContent(c, func() {
								coss.Button(c, "Reset filters").Size(coss.SizeSM).Variant(coss.ButtonOutline).OnClick(func() {
									m.query = ""
									m.selectCategory("")
									a.updateHash()
								})
							})
						})
						return
					}

					if isMobile {
						for i := range visible {
							a.motionCard(c, t, mt, visible[i], true)
						}
					} else {
						for i := 0; i < len(visible); i += 2 {
							ui.Row(c).WidthPercent(100).AlignItems(ui.Start).Gap(16).Children(func() {
								a.motionCard(c, t, mt, visible[i], false)
								if i+1 < len(visible) {
									a.motionCard(c, t, mt, visible[i+1], false)
								} else {
									ui.Box(c).Grow(1).MinWidth(0)
								}
							})
						}
					}
				})
			})
		})
	})
}

func (a *showcaseApp) motionSearchInput(c *ui.Context) {
	m := a.motion
	coss.InputGroup(c, func() {
		coss.InputGroupAddon(c, func() {
			ui.Icon(c, icons.Search()).Size(14, 14)
		})
		coss.InputGroupInput(c, &m.query).Placeholder("Search 33 primitives...")
		if m.query != "" {
			coss.IconButton(c, icons.X(), "Clear search").
				Size(coss.SizeIconXS).
				Variant(coss.ButtonGhost).
				OnClick(func() { m.query = "" })
		}
	})
}

func (a *showcaseApp) motionCategoryNav(c *ui.Context, t theme.Tokens, touch bool) {
	m := a.motion
	ui.Scroll(c).TrackScroll(&m.sidebarScroll).WidthPercent(100).Grow(1).MinHeight(0).Children(func() {
		ui.Column(c).WidthPercent(100).Padding(8, 14, 16, 12).Gap(2).Children(func() {
			a.genericNavItem(c, t, "All Primitives", m.selectedCat == "", len(m.catalog), touch, func() {
				m.selectCategory("")
				a.mobileNavOpen = false
				a.updateHash()
			})
			coss.Separator(c).Margin(6, 0)
			for _, cat := range motionCategories {
				id := cat.ID
				a.genericNavItem(c, t, cat.Label, m.selectedCat == id, m.countByCategory(id), touch, func() {
					m.selectCategory(id)
					a.mobileNavOpen = false
					a.updateHash()
				})
			}
		})
	})
}

func (a *showcaseApp) motionCard(c *ui.Context, t theme.Tokens, mt motion.Theme, e motionDemoEntry, isMobile bool) {
	m := a.motion
	showCode := m.openCode[e.ID]
	copied := c.Now().Before(m.copiedUntil[e.ID])
	if copied {
		c.After(250 * time.Millisecond)
	}

	card := coss.Card(c, func() {
		// Card Header
		ui.Row(c).
			WidthPercent(100).
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween).
			Padding(10, 12).
			Gap(8).
			BorderWidth(0, 0, 1, 0).
			BorderColor(t.Border).
			Children(func() {
				ui.Column(c).Grow(1).MinWidth(0).Gap(1).Children(func() {
					ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
						ui.Text(c, e.Name).FontSize(13).FontWeight(600).TextColor(t.Foreground).SingleLine()
						coss.Badge(c, e.Category).Variant(coss.BadgeSecondary).Size(coss.BadgeSizeSM)
					})
					ui.Text(c, e.Description).FontSize(11).TextColor(t.MutedForeground).SingleLine().Ellipsis("…")
				})
				ui.Row(c).AlignItems(ui.Center).Gap(4).Shrink(0).Children(func() {
					codeLabel := "Code"
					if showCode {
						codeLabel = "Hide"
					}
					coss.Button(c, codeLabel).Size(coss.SizeXS).Variant(coss.ButtonGhost).OnClick(func() {
						m.openCode[e.ID] = !m.openCode[e.ID]
					})
					copyIcon := icons.Copy()
					if copied {
						copyIcon = icons.Check()
					}
					coss.IconButton(c, copyIcon, "Copy Go snippet").
						Size(coss.SizeIconXS).
						Variant(coss.ButtonGhost).
						OnClick(func() {
							c.WriteClipboard(e.Code)
							m.copiedUntil[e.ID] = c.Now().Add(1400 * time.Millisecond)
						})
				})
			})

		// Live Primitive Stage
		padV, padH := float32(20), float32(18)
		minH := float32(188)
		if isMobile {
			padV, padH = 14, 12
			minH = 160
		}
		ui.Column(c).
			WidthPercent(100).
			MinHeight(minH).
			Center().
			Padding(padV, padH).
			Children(func() {
				if e.Render != nil {
					e.Render(c, mt, m)
				}
			})

		// Expandable Go Source Block
		if showCode && e.Code != "" {
			ui.Box(c).
				WidthPercent(100).
				Padding(12, 14).
				Background(t.CodeHighlight).
				BorderWidth(1, 0, 0, 0).
				BorderColor(t.Border).
				Children(func() {
					ui.Text(c, e.Code).FontSize(11).TextColor(t.Foreground)
				})
		}
	})
	if !isMobile {
		card.Grow(1).MinWidth(0)
	}
}
