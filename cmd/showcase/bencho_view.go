package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

const benchoPageSize = 12

type benchoState struct {
	selectedKind  bencho.Interaction
	query         string
	lastQuery     string
	page          int
	sidebarScroll ui.ScrollState
	contentScroll ui.ScrollState
	catalog       []bencho.Block
}

func newBenchoState(kind string) *benchoState {
	sel := bencho.Interaction("")
	for _, k := range bencho.AllInteractions {
		if strings.EqualFold(string(k), kind) {
			sel = k
			break
		}
	}
	return &benchoState{
		selectedKind: sel,
		page:         1,
		catalog:      bencho.Catalog(),
	}
}

func (b *benchoState) selectKind(kind bencho.Interaction) {
	b.selectedKind = kind
	b.page = 1
	b.contentScroll.Y = 0
}

func (b *benchoState) countByKind(kind bencho.Interaction) int {
	if kind == "" {
		return len(b.catalog)
	}
	n := 0
	for _, blk := range b.catalog {
		if blk.Interaction == kind {
			n++
		}
	}
	return n
}

func (b *benchoState) filtered() []bencho.Block {
	if b.query != b.lastQuery {
		b.lastQuery = b.query
		b.page = 1
		b.contentScroll.Y = 0
	}
	q := strings.ToLower(strings.TrimSpace(b.query))
	out := make([]bencho.Block, 0, len(b.catalog))
	for _, blk := range b.catalog {
		if b.selectedKind != "" && blk.Interaction != b.selectedKind {
			continue
		}
		if q != "" {
			if !strings.Contains(strings.ToLower(blk.Title), q) &&
				!strings.Contains(strings.ToLower(blk.Slug), q) &&
				!strings.Contains(strings.ToLower(string(blk.Interaction)), q) {
				continue
			}
		}
		out = append(out, blk)
	}
	return out
}

func (a *showcaseApp) benchoSearchInput(c *ui.Context) {
	b := a.bencho
	coss.InputGroup(c, func() {
		coss.InputGroupAddon(c, func() {
			ui.Icon(c, icons.Search()).Size(14, 14)
		})
		coss.InputGroupInput(c, &b.query).Placeholder("Search 184 blocks...")
		if b.query != "" {
			coss.IconButton(c, icons.X(), "Clear search").
				Size(coss.SizeIconXS).
				Variant(coss.ButtonGhost).
				OnClick(func() { b.query = "" })
		}
	})
}

func (a *showcaseApp) benchoCategoryNav(c *ui.Context, t theme.Tokens, touch bool) {
	b := a.bencho
	ui.Scroll(c).TrackScroll(&b.sidebarScroll).WidthPercent(100).Grow(1).MinHeight(0).Children(func() {
		ui.Column(c).WidthPercent(100).Padding(8, 14, 16, 12).Gap(2).Children(func() {
			a.genericNavItem(c, t, "All Interactions", b.selectedKind == "", b.countByKind(""), touch, func() {
				b.selectKind("")
				a.mobileNavOpen = false
				a.updateHash()
			})
			coss.Separator(c).Margin(6, 0)
			for _, kind := range bencho.AllInteractions {
				k := kind
				a.genericNavItem(c, t, string(k), b.selectedKind == k, b.countByKind(k), touch, func() {
					b.selectKind(k)
					a.mobileNavOpen = false
					a.updateHash()
				})
			}
		})
	})
}

func (a *showcaseApp) renderBenchoTab(c *ui.Context, t theme.Tokens, isMobile bool) {
	b := a.bencho
	visible := b.filtered()
	liveCount := bencho.ImplementedCount()
	totalCount := len(b.catalog)

	totalPages := (len(visible) + benchoPageSize - 1) / benchoPageSize
	if totalPages < 1 {
		totalPages = 1
	}
	if b.page < 1 {
		b.page = 1
	}
	if b.page > totalPages {
		b.page = totalPages
	}
	start := (b.page - 1) * benchoPageSize
	end := start + benchoPageSize
	if end > len(visible) {
		end = len(visible)
	}
	pageItems := visible[start:end]

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
						a.benchoSearchInput(c)
					})
					a.benchoCategoryNav(c, t, false)
				})
		}

		ui.Column(c).Grow(1).MinWidth(0).HeightPercent(100).Children(func() {
			title := "Bencho Blocks"
			if b.selectedKind != "" {
				title = string(b.selectedKind) + " Blocks"
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
								coss.IconButton(c, icons.PanelLeft(), "Open interactions").
									Size(coss.SizeIconSM).
									Variant(coss.ButtonOutline).
									OnClick(func() { a.mobileNavOpen = true })
								ui.Text(c, title).FontSize(14.5).FontWeight(600).TextColor(t.Foreground).SingleLine().Ellipsis("…")
								coss.Badge(c, fmt.Sprintf("%d", len(visible))).Variant(coss.BadgeSecondary).Size(coss.BadgeSizeSM)
							})
							coss.Badge(c, fmt.Sprintf("%d / %d", liveCount, totalCount)).Variant(coss.BadgeOutline).Size(coss.BadgeSizeSM)
						})
						a.benchoSearchInput(c)
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
							ui.Text(c, fmt.Sprintf("Showing %d–%d of %d blocks (%d / %d live from bencho.dev catalog)", min(start+1, len(visible)), end, len(visible), liveCount, totalCount)).
								FontSize(13).TextColor(t.MutedForeground)
						})
						coss.Badge(c, "mygox-ui/bencho").Variant(coss.BadgeOutline)
					})
			}

			ui.Scroll(c).TrackScroll(&b.contentScroll).WidthPercent(100).Grow(1).MinHeight(0).Children(func() {
				pad := float32(24)
				gap := float32(16)
				if isMobile {
					pad = 12
					gap = 12
				}
				ui.Column(c).WidthPercent(100).Padding(pad).Gap(gap).Children(func() {
					if len(pageItems) == 0 {
						coss.Empty(c, func() {
							coss.EmptyHeader(c, func() {
								coss.EmptyMedia(c, icons.Search())
								coss.EmptyTitle(c, "No matching Bencho blocks")
								coss.EmptyDescription(c, "Clear the search query or pick another interaction category.")
							})
							coss.EmptyContent(c, func() {
								coss.Button(c, "Reset filters").Size(coss.SizeSM).Variant(coss.ButtonOutline).OnClick(func() {
									b.query = ""
									b.selectKind("")
								})
							})
						})
						return
					}

					if isMobile {
						for _, blk := range pageItems {
							a.benchoCard(c, t, blk, true)
						}
					} else {
						for i := 0; i < len(pageItems); i += 2 {
							left := pageItems[i]
							hasRight := i+1 < len(pageItems)
							ui.Row(c).WidthPercent(100).Gap(gap).AlignItems(ui.Stretch).Children(func() {
								a.benchoCard(c, t, left, false)
								if hasRight {
									a.benchoCard(c, t, pageItems[i+1], false)
								} else {
									ui.Box(c).Grow(1).MinWidth(0)
								}
							})
						}
					}

					if totalPages > 1 {
						a.benchoPagination(c, t, totalPages)
					}
				})
			})
		})
	})
}

func (a *showcaseApp) benchoPagination(c *ui.Context, t theme.Tokens, totalPages int) {
	b := a.bencho
	ui.Row(c).
		WidthPercent(100).
		Padding(8, 0, 4, 0).
		AlignItems(ui.Center).
		Justify(ui.SpaceBetween).
		Gap(8).
		Children(func() {
			coss.Button(c, "Previous").
				Size(coss.SizeSM).
				Variant(coss.ButtonOutline).
				LeadingIcon(icons.ChevronLeft()).
				Disabled(b.page <= 1).
				OnClick(func() {
					if b.page > 1 {
						b.page--
						b.contentScroll.Y = 0
					}
				})
			ui.Text(c, fmt.Sprintf("Page %d of %d", b.page, totalPages)).
				FontSize(12.5).
				TextColor(t.MutedForeground)
			coss.Button(c, "Next").
				Size(coss.SizeSM).
				Variant(coss.ButtonOutline).
				TrailingIcon(icons.ChevronRight()).
				Disabled(b.page >= totalPages).
				OnClick(func() {
					if b.page < totalPages {
						b.page++
						b.contentScroll.Y = 0
					}
				})
		})
}

func (a *showcaseApp) benchoCard(c *ui.Context, t theme.Tokens, blk bencho.Block, isMobile bool) {
	card := coss.Card(c, func() {
		ui.Row(c).
			WidthPercent(100).
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween).
			Padding(10, 14).
			Gap(8).
			BorderWidth(0, 0, 1, 0).
			BorderColor(t.Border).
			Children(func() {
				ui.Column(c).Grow(1).MinWidth(0).Gap(1).Children(func() {
					ui.Text(c, blk.Title).FontSize(13.5).FontWeight(600).TextColor(t.Foreground).SingleLine()
					sub := blk.Summary
					if sub == "" {
						sub = "/blocks/" + blk.Slug
					}
					ui.Text(c, sub).FontSize(11).TextColor(t.MutedForeground).SingleLine().Ellipsis("…")
				})
				badgeLabel := string(blk.Interaction)
				variant := coss.BadgeSecondary
				if !blk.Implemented() {
					badgeLabel = string(blk.Interaction) + " · Planned"
					variant = coss.BadgeOutline
				}
				coss.Badge(c, badgeLabel).Variant(variant).Size(coss.BadgeSizeSM)
			})

		if blk.Implemented() {
			minH := float32(190)
			if isMobile {
				minH = 160
			}
			ui.Column(c).
				WidthPercent(100).
				MinHeight(minH).
				Center().
				Padding(20, 16).
				Children(func() {
					blk.Demo(c)
				})
		}
	})
	if !isMobile {
		card.Grow(1).MinWidth(0)
	}
}
