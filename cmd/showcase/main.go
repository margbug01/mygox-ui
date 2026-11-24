package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/particles"
	"github.com/margbug01/mygox-ui/coss/theme"
)

const mobileBreakpoint float32 = 760

type activeSection string

const (
	sectionCoss   activeSection = "coss"
	sectionFlow   activeSection = "flow"
	sectionMotion activeSection = "motion"
	sectionBencho activeSection = "bencho"
)

type showcaseApp struct {
	dark          bool
	section       activeSection
	mobileNavOpen bool

	// Coss 530 Particles state
	selectedFam   string
	query         string
	page          int
	openCode      map[string]bool
	copiedUntil   map[string]time.Time
	sidebarScroll ui.ScrollState
	contentScroll ui.ScrollState

	// Flow Canvas, Motion Effects & Bencho Blocks states
	flow   *flowState
	motion *motionState
	bencho *benchoState
}

func newShowcaseApp() *showcaseApp {
	app := &showcaseApp{
		dark:        initialDark(),
		section:     sectionCoss,
		page:        1,
		openCode:    map[string]bool{},
		copiedUntil: map[string]time.Time{},
		flow:        newFlowState("schema"),
		motion:      newMotionState(""),
		bencho:      newBenchoState(""),
	}
	app.applyHash(initialHash())
	watchHashChange(func(h string) {
		app.applyHash(h)
	})
	return app
}

func (a *showcaseApp) applyHash(raw string) {
	raw = strings.Trim(strings.TrimSpace(raw), "/")
	if raw == "" {
		return
	}
	parts := strings.SplitN(raw, "/", 2)
	head := strings.ToLower(parts[0])
	sub := ""
	if len(parts) > 1 {
		sub = parts[1]
	}
	switch head {
	case "flow":
		a.section = sectionFlow
		if sub != "" {
			a.flow.loadPreset(sub)
		}
	case "motion":
		a.section = sectionMotion
		a.motion.selectCategory(sub)
	case "bencho":
		a.section = sectionBencho
		a.bencho = newBenchoState(sub)
	case "coss", "ui":
		a.section = sectionCoss
		if sub != "" && len(particles.ByFamily(sub)) > 0 {
			a.selectedFam = sub
		}
	default:
		// Backward-compatible direct particle family link like #button or #table.
		if len(particles.ByFamily(head)) > 0 {
			a.section = sectionCoss
			a.selectedFam = head
		}
	}
}

func (a *showcaseApp) updateHash() {
	switch a.section {
	case sectionFlow:
		syncHash("flow/" + a.flow.preset)
	case sectionMotion:
		if a.motion.selectedCat != "" {
			syncHash("motion/" + a.motion.selectedCat)
		} else {
			syncHash("motion")
		}
	case sectionBencho:
		if a.bencho.selectedKind != "" {
			syncHash("bencho/" + strings.ToLower(string(a.bencho.selectedKind)))
		} else {
			syncHash("bencho")
		}
	default:
		if a.selectedFam != "" {
			syncHash("coss/" + a.selectedFam)
		} else {
			syncHash("coss")
		}
	}
}

func (a *showcaseApp) setSection(sec activeSection) {
	if a.section == sec {
		return
	}
	a.section = sec
	a.mobileNavOpen = false
	if sec == sectionFlow {
		a.flow.fitted = false
	}
	a.updateHash()
}

func main() {
	app := newShowcaseApp()
	mygo.App.WhenReady(func() {
		w := mygo.NewWindow(mygo.WindowOptions{
			Title:   "mygox-ui — Components, Node Canvas & GPU Motion for MyGo",
			Width:   1280,
			Height:  840,
			Content: ui.View(app.view),
		})
		w.Center()
		w.Show()
		w.Focus()
		mygo.App.Focus()
		notifyWebReady()
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

func (a *showcaseApp) filteredParticles() []particles.Entry {
	all := particles.ByFamily(a.selectedFam)
	q := strings.ToLower(strings.TrimSpace(a.query))
	if q == "" {
		return all
	}
	out := make([]particles.Entry, 0, len(all))
	for _, e := range all {
		if strings.Contains(strings.ToLower(e.ID), q) ||
			strings.Contains(strings.ToLower(e.Description), q) ||
			strings.Contains(strings.ToLower(e.Family), q) {
			out = append(out, e)
		}
	}
	return out
}

func (a *showcaseApp) selectFamily(fam string) {
	a.selectedFam = fam
	a.page = 1
	a.contentScroll.Y = 0
	a.mobileNavOpen = false
	a.updateHash()
}

func (a *showcaseApp) view(c *ui.Context) {
	t := theme.Apply(c, a.dark)
	winW, winH := c.Size()
	if vw := webViewportWidth(); vw > 0 && (winW <= 0 || vw < winW) {
		winW = vw
	}
	isMobile := winW > 0 && winW < mobileBreakpoint

	ui.Column(c).
		WidthPercent(100).
		HeightPercent(100).
		Background(t.Background).
		TextColor(t.Foreground).
		Children(func() {
			// Unified Coss-style Suite Top Navigation Bar
			a.suiteTopBar(c, t, isMobile)

			// Active Module Workspace
			switch a.section {
			case sectionFlow:
				a.renderFlowTab(c, t, winW, winH, isMobile)
			case sectionMotion:
				a.renderMotionTab(c, t, isMobile)
			case sectionBencho:
				a.renderBenchoTab(c, t, isMobile)
			default:
				a.renderCossTab(c, t, isMobile)
			}
		})

	if isMobile {
		a.mobileDrawer(c, t)
	}
}

func (a *showcaseApp) suiteTopBar(c *ui.Context, t theme.Tokens, isMobile bool) {
	padH := float32(16)
	if isMobile {
		padH = 10
	}
	ui.Row(c).
		WidthPercent(100).
		Height(48).
		Shrink(0).
		Background(t.Sidebar).
		BorderWidth(0, 0, 1, 0).
		BorderColor(t.Border).
		Padding(0, padH).
		AlignItems(ui.Center).
		Justify(ui.SpaceBetween).
		Gap(8).
		Children(func() {
			// Left: Brand + Suite Segmented Tabs
			ui.Row(c).AlignItems(ui.Center).Gap(10).Children(func() {
				if !isMobile {
					ui.Text(c, "mygox-ui").FontSize(15).FontWeight(700).TextColor(t.Foreground)
					ui.Box(c).Width(1).Height(16).Background(t.Border)
				}
				// Coss segmented control for UI / Flow / Motion
				ui.Row(c).
					Background(t.Muted).
					Radius(t.RadiusMD).
					Padding(3).
					Gap(2).
					AlignItems(ui.Center).
					Children(func() {
						a.suiteTabButton(c, t, "UI", "530", sectionCoss, isMobile)
						a.suiteTabButton(c, t, "Flow", "Canvas", sectionFlow, isMobile)
						a.suiteTabButton(c, t, "Motion", "33", sectionMotion, isMobile)
						a.suiteTabButton(c, t, "Bencho", "184", sectionBencho, isMobile)
					})
			})

			// Right: Version badge + Global Theme Toggle
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if !isMobile {
					coss.Badge(c, "mygo v0.10.0").Variant(coss.BadgeOutline).Size(coss.BadgeSizeSM)
				}
				a.themeToggleButton(c)
			})
		})
}

func (a *showcaseApp) suiteTabButton(c *ui.Context, t theme.Tokens, label, count string, sec activeSection, isMobile bool) {
	active := a.section == sec
	padH := float32(10)
	if isMobile {
		padH = 8
	}
	b := ui.ButtonBase(c).
		Height(26).
		Padding(0, padH).
		Radius(t.RadiusSM).
		Center().
		Cursor(ui.CursorPointer)
	if active {
		b.Background(t.Card).TextColor(t.Foreground).Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.08))
	} else {
		b.TextColor(t.MutedForeground)
		if b.Hovered() {
			b.TextColor(t.Foreground)
		}
	}
	b.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(5).Children(func() {
			fw := 500
			if active {
				fw = 600
			}
			ui.Text(c, label).FontSize(12).FontWeight(fw).SingleLine()
			if !isMobile || count != "Canvas" {
				ui.Text(c, count).FontSize(10.5).TextColor(t.MutedForeground).SingleLine()
			}
		})
	})
	if b.Clicked() {
		a.setSection(sec)
	}
}

func (a *showcaseApp) renderCossTab(c *ui.Context, t theme.Tokens, isMobile bool) {
	allEntries := particles.All()
	families := particles.Families()
	visible := a.filteredParticles()

	ui.Row(c).WidthPercent(100).Grow(1).MinHeight(0).Children(func() {
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
						a.cossSearchInput(c)
					})
					a.cossFamilyNav(c, t, families, len(allEntries), false)
				})
		}

		ui.Column(c).Grow(1).MinWidth(0).HeightPercent(100).Children(func() {
			if isMobile {
				a.cossMobileHeader(c, t, len(visible))
			} else {
				a.cossDesktopHeader(c, t, len(visible))
			}

			ui.Scroll(c).TrackScroll(&a.contentScroll).WidthPercent(100).Grow(1).MinHeight(0).Children(func() {
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
								coss.EmptyTitle(c, "No matching particles")
								coss.EmptyDescription(c, "Try clearing your search query or choosing another family.")
							})
							coss.EmptyContent(c, func() {
								coss.Button(c, "Reset filters").Size(coss.SizeSM).Variant(coss.ButtonOutline).OnClick(func() {
									a.query = ""
									a.selectFamily("")
								})
							})
						})
						return
					}

					const pageSize = 24
					pageItems := visible
					totalPages := 1
					if a.selectedFam == "" && strings.TrimSpace(a.query) == "" && len(visible) > pageSize {
						totalPages = (len(visible) + pageSize - 1) / pageSize
						if a.page < 1 {
							a.page = 1
						}
						if a.page > totalPages {
							a.page = totalPages
						}
						start := (a.page - 1) * pageSize
						end := start + pageSize
						if end > len(visible) {
							end = len(visible)
						}
						pageItems = visible[start:end]
					}

					if isMobile {
						for i := range pageItems {
							a.particleCard(c, t, pageItems[i], true)
						}
					} else {
						for i := 0; i < len(pageItems); i += 2 {
							ui.Row(c).WidthPercent(100).AlignItems(ui.Start).Gap(16).Children(func() {
								a.particleCard(c, t, pageItems[i], false)
								if i+1 < len(pageItems) {
									a.particleCard(c, t, pageItems[i+1], false)
								} else {
									ui.Box(c).Grow(1).MinWidth(0)
								}
							})
						}
					}

					if totalPages > 1 {
						prevPage := a.page
						ui.Row(c).WidthPercent(100).Center().Padding(8, 0, 20, 0).Children(func() {
							if isMobile {
								a.mobilePagination(c, t, totalPages)
							} else {
								coss.Pagination(c, &a.page, totalPages)
							}
						})
						if a.page != prevPage {
							a.contentScroll.Y = 0
						}
					}
				})
			})
		})
	})
}

func (a *showcaseApp) themeToggleButton(c *ui.Context) {
	themeIcon := icons.Moon()
	themeLabel := "Switch to dark theme"
	if a.dark {
		themeIcon = icons.Sun()
		themeLabel = "Switch to light theme"
	}
	coss.IconButton(c, themeIcon, themeLabel).
		Size(coss.SizeIconSM).
		Variant(coss.ButtonGhost).
		OnClick(func() {
			a.dark = !a.dark
			syncWebTheme(a.dark)
		})
}

func (a *showcaseApp) cossSearchInput(c *ui.Context) {
	coss.InputGroup(c, func() {
		coss.InputGroupAddon(c, func() {
			ui.Icon(c, icons.Search()).Size(14, 14)
		})
		coss.InputGroupInput(c, &a.query).Placeholder("Search particles...")
		if a.query != "" {
			coss.IconButton(c, icons.X(), "Clear search").
				Size(coss.SizeIconXS).
				Variant(coss.ButtonGhost).
				OnClick(func() { a.query = "" })
		}
	})
}

func (a *showcaseApp) cossDesktopHeader(c *ui.Context, t theme.Tokens, count int) {
	ui.Row(c).
		WidthPercent(100).
		AlignItems(ui.Center).
		Justify(ui.SpaceBetween).
		Padding(16, 24).
		BorderWidth(0, 0, 1, 0).
		BorderColor(t.Border).
		Children(func() {
			title := "All Particles"
			if a.selectedFam != "" {
				title = a.selectedFam
			}
			ui.Column(c).Gap(2).Children(func() {
				ui.Text(c, title).FontSize(18).FontWeight(600).TextColor(t.Foreground)
				ui.Text(c, fmt.Sprintf("Showing %d native MyGo particles ported from coss.com/ui", count)).
					FontSize(13).TextColor(t.MutedForeground)
			})
			coss.Badge(c, "mygox-ui/coss").Variant(coss.BadgeOutline)
		})
}

func (a *showcaseApp) cossMobileHeader(c *ui.Context, t theme.Tokens, count int) {
	title := "All Particles"
	if a.selectedFam != "" {
		title = a.selectedFam
	}
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
					coss.IconButton(c, icons.PanelLeft(), "Open families").
						Size(coss.SizeIconSM).
						Variant(coss.ButtonOutline).
						OnClick(func() { a.mobileNavOpen = true })
					ui.Text(c, title).FontSize(14.5).FontWeight(600).TextColor(t.Foreground).SingleLine().Ellipsis("…")
					coss.Badge(c, fmt.Sprintf("%d", count)).Variant(coss.BadgeSecondary).Size(coss.BadgeSizeSM)
				})
			})
			a.cossSearchInput(c)
		})
}

func (a *showcaseApp) mobileDrawer(c *ui.Context, t theme.Tokens) {
	black := ui.RGB(0, 0, 0)
	ui.DialogBase(c, &a.mobileNavOpen, func(backdrop, panel ui.Element) {
		backdrop.Background(theme.Alpha(black, 0.42)).Justify(ui.Start).AlignItems(ui.Stretch)
		panel.
			Width(288).
			MaxWidthPercent(82).
			HeightPercent(100).
			AlignSelf(ui.Start).
			Background(t.Sidebar).
			TextColor(t.Foreground).
			BorderWidth(0, 1, 0, 0).
			BorderColor(t.Border).
			Shadow(4, 0, 24, 0, theme.Alpha(black, 0.28)).
			Children(func() {
				heading := "UI Families"
				if a.section == sectionFlow {
					heading = "Node Palette"
				} else if a.section == sectionMotion {
					heading = "Motion Categories"
				} else if a.section == sectionBencho {
					heading = "Bencho Interactions"
				}
				ui.Row(c).
					WidthPercent(100).
					AlignItems(ui.Center).
					Justify(ui.SpaceBetween).
					Padding(14, 14).
					BorderWidth(0, 0, 1, 0).
					BorderColor(t.SidebarBorder).
					Children(func() {
						ui.Text(c, heading).FontSize(15).FontWeight(600).TextColor(t.Foreground)
						coss.IconButton(c, icons.X(), "Close menu").
							Size(coss.SizeIconSM).
							Variant(coss.ButtonGhost).
							OnClick(func() { a.mobileNavOpen = false })
					})
				switch a.section {
				case sectionFlow:
					ui.Column(c).WidthPercent(100).Grow(1).MinHeight(0).Padding(12).Gap(10).Children(func() {
						a.flowPaletteContent(c, t)
					})
				case sectionMotion:
					a.motionCategoryNav(c, t, true)
				case sectionBencho:
					a.benchoCategoryNav(c, t, true)
				default:
					a.cossFamilyNav(c, t, particles.Families(), len(particles.All()), true)
				}
			})
	})
}

func (a *showcaseApp) cossFamilyNav(c *ui.Context, t theme.Tokens, families []string, total int, touch bool) {
	ui.Scroll(c).TrackScroll(&a.sidebarScroll).WidthPercent(100).Grow(1).MinHeight(0).Children(func() {
		ui.Column(c).WidthPercent(100).Padding(8, 14, 16, 12).Gap(2).Children(func() {
			a.genericNavItem(c, t, "All Particles", a.selectedFam == "", total, touch, func() {
				a.selectFamily("")
			})
			coss.Separator(c).Margin(6, 0)
			for _, fam := range families {
				f := fam
				a.genericNavItem(c, t, f, a.selectedFam == f, len(particles.ByFamily(f)), touch, func() {
					a.selectFamily(f)
				})
			}
		})
	})
}

func (a *showcaseApp) genericNavItem(c *ui.Context, t theme.Tokens, label string, active bool, count int, touch bool, onClick func()) {
	h := float32(30)
	fontSz := float32(13)
	if touch {
		h = 38
		fontSz = 14
	}
	row := ui.ButtonBase(c).
		WidthPercent(100).
		Height(h).
		Padding(0, 10).
		Radius(t.RadiusMD).
		Justify(ui.SpaceBetween).
		AlignItems(ui.Center).
		Cursor(ui.CursorPointer)
	if active {
		row.Background(t.SidebarAccent).TextColor(t.SidebarAccentForeground)
	} else if row.Hovered() {
		row.Background(theme.Alpha(t.SidebarAccent, 0.7)).TextColor(t.Foreground)
	} else {
		row.TextColor(t.SidebarForeground)
	}
	row.Children(func() {
		fw := 400
		if active {
			fw = 500
		}
		ui.Text(c, label).FontSize(fontSz).FontWeight(fw).SingleLine()
		ui.Text(c, fmt.Sprintf("%d", count)).FontSize(12).TextColor(t.MutedForeground)
	})
	if row.Clicked() && onClick != nil {
		onClick()
	}
}

func (a *showcaseApp) mobilePagination(c *ui.Context, t theme.Tokens, totalPages int) {
	ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(8).Children(func() {
		coss.Button(c, "Previous").
			Size(coss.SizeSM).
			Variant(coss.ButtonOutline).
			LeadingIcon(icons.ChevronLeft()).
			Disabled(a.page <= 1).
			OnClick(func() {
				if a.page > 1 {
					a.page--
				}
			})
		ui.Text(c, fmt.Sprintf("Page %d / %d", a.page, totalPages)).FontSize(12).TextColor(t.MutedForeground)
		coss.Button(c, "Next").
			Size(coss.SizeSM).
			Variant(coss.ButtonOutline).
			TrailingIcon(icons.ChevronRight()).
			Disabled(a.page >= totalPages).
			OnClick(func() {
				if a.page < totalPages {
					a.page++
				}
			})
	})
}

func (a *showcaseApp) particleCard(c *ui.Context, t theme.Tokens, e particles.Entry, isMobile bool) {
	showCode := a.openCode[e.ID]
	copied := c.Now().Before(a.copiedUntil[e.ID])
	if copied {
		c.After(250 * time.Millisecond)
	}

	card := coss.Card(c, func() {
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
					ui.Text(c, e.ID).FontSize(13).FontWeight(600).TextColor(t.Foreground).SingleLine()
					ui.Text(c, e.Description).FontSize(11).TextColor(t.MutedForeground).SingleLine().Ellipsis("…")
				})
				ui.Row(c).AlignItems(ui.Center).Gap(4).Shrink(0).Children(func() {
					codeLabel := "Code"
					if showCode {
						codeLabel = "Hide"
					}
					coss.Button(c, codeLabel).Size(coss.SizeXS).Variant(coss.ButtonGhost).OnClick(func() {
						a.openCode[e.ID] = !a.openCode[e.ID]
					})
					copyIcon := icons.Copy()
					if copied {
						copyIcon = icons.Check()
					}
					coss.IconButton(c, copyIcon, "Copy Go source").
						Size(coss.SizeIconXS).
						Variant(coss.ButtonGhost).
						OnClick(func() {
							c.WriteClipboard(e.Source)
							a.copiedUntil[e.ID] = c.Now().Add(1400 * time.Millisecond)
						})
				})
			})

		stagePadV, stagePadH := float32(22), float32(18)
		minH := float32(128)
		if isMobile {
			stagePadV, stagePadH = 16, 12
			minH = 104
		}
		ui.Column(c).
			WidthPercent(100).
			MinHeight(minH).
			Center().
			Padding(stagePadV, stagePadH).
			Children(func() {
				if e.View != nil {
					e.View(c)
				}
			})

		if showCode && e.Source != "" {
			ui.Box(c).
				WidthPercent(100).
				Padding(12, 14).
				Background(t.CodeHighlight).
				BorderWidth(1, 0, 0, 0).
				BorderColor(t.Border).
				Children(func() {
					ui.Text(c, e.Source).FontSize(11).TextColor(t.Foreground)
				})
		}
	})
	if !isMobile {
		card.Grow(1).MinWidth(0)
	}
}
