package particles

import (
	_ "embed"
	"fmt"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed switch.go
var switchSource string

func init() {
	Register(Entry{ID: "p-switch-1", Family: "switch", Description: "Basic switch", View: PSwitch1, Source: extractFunc(switchSource, "PSwitch1")})
	Register(Entry{ID: "p-switch-2", Family: "switch", Description: "Disabled switch", View: PSwitch2, Source: extractFunc(switchSource, "PSwitch2")})
	Register(Entry{ID: "p-switch-3", Family: "switch", Description: "Switch with description", View: PSwitch3, Source: extractFunc(switchSource, "PSwitch3")})
	Register(Entry{ID: "p-switch-4", Family: "switch", Description: "Switch card", View: PSwitch4, Source: extractFunc(switchSource, "PSwitch4")})
	Register(Entry{ID: "p-switch-5", Family: "switch", Description: "Switch in form", View: PSwitch5, Source: extractFunc(switchSource, "PSwitch5")})
	Register(Entry{ID: "p-switch-6", Family: "switch", Description: "Custom size switch", View: PSwitch6, Source: extractFunc(switchSource, "PSwitch6")})
	Register(Entry{ID: "p-switch-7", Family: "switch", Description: "Weekly availability editor with time range combobox pickers", View: PSwitch7, Source: extractFunc(switchSource, "PSwitch7")})
	Register(Entry{ID: "p-switch-8", Family: "switch", Description: "Weekly availability editor with grouped time range controls", View: PSwitch8, Source: extractFunc(switchSource, "PSwitch8")})
	Register(Entry{ID: "p-switch-9", Family: "switch", Description: "Weekly availability editor with From/To labeled time groups", View: PSwitch9, Source: extractFunc(switchSource, "PSwitch9")})
	Register(Entry{ID: "p-switch-liquid", Family: "switch", Description: "Switch whose thumb flies across like a drop and can be dragged (mygox motion, not on coss.com)", View: PSwitchLiquid, Source: extractFunc(switchSource, "PSwitchLiquid")})
}

// PSwitchLiquid is mygox's own: the thumb springs across, stretching like
// a drop the faster it flies, and can be dragged across too.
func PSwitchLiquid(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() { coss.Switch(c, on, coss.SwitchOptions{Label: "Marketing emails", Liquid: true}) })
}

func PSwitch1(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() { coss.Switch(c, on, coss.SwitchOptions{Label: "Marketing emails"}) })
}

func PSwitch2(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() { coss.Switch(c, on, coss.SwitchOptions{Label: "Marketing emails"}).Disabled(true) })
}

func PSwitch3(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).AlignItems(ui.Start).Gap(8)
	on := ui.Local(host, "on", func() bool { return true })
	host.Children(func() {
		coss.Switch(c, on)
		ui.Column(c).Gap(4).Children(func() {
			coss.Label(c, "Marketing emails").LineHeight(16.0 / 14)
			ui.Text(c, "By enabling marketing emails, you agree to receive emails.").
				FontSize(12).LineHeight(16.0 / 12).TextColor(t.MutedForeground)
		})
	})
}

func PSwitch4(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return true })
	host.Children(func() {
		card := ui.ButtonBase(c).
			AlignItems(ui.Center).
			Gap(24).
			Padding(12).
			Radius(t.RadiusLG).
			Border(1, t.Border).
			Cursor(ui.CursorPointer)
		if *on {
			card.Border(1, theme.Alpha(t.Primary, 0.48)).Background(theme.Alpha(t.Accent, 0.5))
		} else if card.Hovered() {
			card.Background(theme.Alpha(t.Accent, 0.5))
		}
		card.Children(func() {
			ui.Column(c).Gap(4).Children(func() {
				ui.Text(c, "Enable notifications").FontSize(14).FontWeight(500).LineHeight(16.0 / 14).TextColor(t.Foreground)
				ui.Text(c, "You can enable or disable notifications at any time.").
					FontSize(12).LineHeight(16.0 / 12).TextColor(t.MutedForeground)
			})
			coss.Switch(c, on, coss.SwitchOptions{Thumb: 12})
		})
		if card.Clicked() {
			*on = !*on
		}
	})
}

func PSwitch5(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return true })
	until := ui.Local(host, "loadingUntil", func() time.Time { return time.Time{} })
	loading := c.Now().Before(*until)
	if loading {
		c.After(150 * time.Millisecond)
	}
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() { coss.Switch(c, on, coss.SwitchOptions{Label: "Enable marketing emails"}) })
			coss.Button(c, "Submit").Loading(loading).OnClick(func() {
				*until = c.Now().Add(800 * time.Millisecond)
			})
		})
	})
}

func PSwitch6(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() { coss.Switch(c, on, coss.SwitchOptions{Thumb: 12}).Label("Notifications") })
}

type availabilityRange struct {
	ID         int
	Start, End string
}

type availabilityState struct {
	NextID int
	Days   map[string][]availabilityRange
}

var availabilityDays = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

// availabilityTimes are the times of day in 15 minute steps, "9:00 AM".
var availabilityTimes = func() []string {
	out := make([]string, 0, 96)
	for i := range 96 {
		h, m := i/4, (i%4)*15
		period := "AM"
		if h >= 12 {
			period = "PM"
		}
		dh := h % 12
		if dh == 0 {
			dh = 12
		}
		out = append(out, fmt.Sprintf("%d:%02d %s", dh, m, period))
	}
	return out
}()

func availabilityIndex(time string) int { return slices.Index(availabilityTimes, time) }

func newAvailabilityState() *availabilityState {
	return &availabilityState{
		NextID: 10,
		Days: map[string][]availabilityRange{
			"Monday":    {{ID: 1, Start: "9:00 AM", End: "5:00 PM"}},
			"Tuesday":   {{ID: 2, Start: "9:00 AM", End: "1:00 PM"}, {ID: 3, Start: "3:00 PM", End: "5:00 PM"}},
			"Wednesday": {{ID: 4, Start: "9:00 AM", End: "5:00 PM"}},
			"Thursday":  {{ID: 5, Start: "9:00 AM", End: "5:00 PM"}},
			"Friday":    {{ID: 6, Start: "9:00 AM", End: "5:00 PM"}},
		},
	}
}

type availabilityStyle uint8

const (
	availStylePickers availabilityStyle = iota
	availStyleGroupedArrow
	availStyleGroupedFromTo
)

// availabilityTime is the time combobox of a range end: a button showing
// *v that opens a searchable list of items.
func availabilityTime(c *ui.Context, v *string, items []string, mode availabilityStyle) {
	t := theme.Of(c)
	btn := ui.ButtonBase(c)
	open := ui.Local(btn, "open", func() bool { return false })
	if mode == availStylePickers {
		coss.SelectTrigger(c, btn, coss.SelectSM, func() {
			coss.SelectValue(c, *v, "")
		}).Width(108).FontFeatures("tnum")
	} else {
		btn.Size(96, 28).Center().Cursor(ui.CursorPointer).FontFeatures("tnum").Children(func() {
			ui.Text(c, *v).FontSize(14).TextColor(t.Foreground).SingleLine()
		})
		if btn.Hovered() {
			btn.Background(t.Accent)
		}
	}
	if btn.Clicked() {
		*open = !*open
	}
	coss.ComboboxSearchPopup(c, btn, open, v, items, "Search time", "No times found.")
}

func renderAvailabilityEditor(c *ui.Context, mode availabilityStyle) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).MaxWidth(896)
	st := *ui.Local(host, "avail", newAvailabilityState)

	iconButton := func(ic *ui.SVG, tip string, disabled bool, click func()) ui.Element {
		btn := coss.IconButton(c, ic, tip).Variant(coss.ButtonGhost).Size(coss.SizeIconSM).
			Disabled(disabled).OnClick(click)
		cossTooltip(c, btn.Element, tip)
		return btn.Element
	}

	host.Children(func() {
		for idx, day := range availabilityDays {
			ranges := st.Days[day]
			row := ui.Row(c).WidthPercent(100).AlignItems(ui.Start).Gap(16).Wrap()
			var pt, pb float32 = 12, 12
			if idx == 0 {
				pt = 0
			} else {
				row.BorderWidth(1, 0, 0, 0).BorderColor(t.Border)
			}
			if idx == len(availabilityDays)-1 {
				pb = 0
			}
			row.Padding(pt, 0, pb, 0)
			row.Children(func() {
				enabled := len(ranges) > 0
				ui.Row(c).Width(120).Height(28).AlignItems(ui.Center).Gap(10).Shrink(0).Children(func() {
					coss.Switch(c, &enabled)
					coss.Label(c, day).LineHeight(16.0 / 14)
				})
				if enabled != (len(ranges) > 0) {
					st.Days[day] = nil
					if enabled {
						st.NextID++
						st.Days[day] = []availabilityRange{{ID: st.NextID, Start: "9:00 AM", End: "5:00 PM"}}
					}
					ranges = st.Days[day]
				}

				ui.Row(c).Grow(1).MinWidth(0).AlignItems(ui.Start).Gap(16).Children(func() {
					ui.Column(c).MinWidth(0).Gap(8).Children(func() {
						if len(ranges) == 0 {
							ui.Row(c).Height(28).AlignItems(ui.Center).Children(func() {
								ui.Text(c, "Unavailable").FontSize(14).TextColor(t.MutedForeground)
							})
						}
						for i := range ranges {
							r := &st.Days[day][i]
							ends := availabilityTimes[availabilityIndex(r.Start)+1:]
							start := r.Start
							ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
								switch mode {
								case availStylePickers:
									availabilityTime(c, &r.Start, availabilityTimes[:95], mode)
									ui.Text(c, "–").FontSize(14).TextColor(t.MutedForeground)
									availabilityTime(c, &r.End, ends, mode)
								case availStyleGroupedArrow:
									coss.Group(c, func() {
										availabilityTime(c, &r.Start, availabilityTimes[:95], mode)
										coss.GroupSeparator(c).Height(28).Background(t.Input)
										ui.Row(c).Center().Height(28).Padding(0, 8).Background(t.Muted).TextColor(t.MutedForeground).Children(func() {
											ui.Icon(c, icons.ArrowRight()).Size(14, 14)
										})
										coss.GroupSeparator(c).Height(28).Background(t.Input)
										availabilityTime(c, &r.End, ends, mode)
									})
								case availStyleGroupedFromTo:
									coss.Group(c, func() {
										coss.GroupText(c, "From").Height(28).Padding(0, 11).FontSize(14)
										coss.GroupSeparator(c).Height(28).Background(t.Input)
										availabilityTime(c, &r.Start, availabilityTimes[:95], mode)
										coss.GroupSeparator(c).Height(28).Background(t.Input)
										coss.GroupText(c, "To").Height(28).Padding(0, 11).FontSize(14)
										coss.GroupSeparator(c).Height(28).Background(t.Input)
										availabilityTime(c, &r.End, ends, mode)
									})
								}
								id := r.ID
								iconButton(icons.X(), "Delete range", false, func() {
									st.Days[day] = slices.DeleteFunc(st.Days[day], func(x availabilityRange) bool { return x.ID == id })
								})
							})
							// Moving the start past the end moves the end an hour after it.
							if r.Start != start && availabilityIndex(r.Start) >= availabilityIndex(r.End) {
								r.End = availabilityTimes[min(availabilityIndex(r.Start)+4, len(availabilityTimes)-1)]
							}
						}
					})

					ui.Row(c).Grow(1).Justify(ui.End).Gap(4).Shrink(0).Children(func() {
						last := len(availabilityTimes) - 1
						addDisabled := len(ranges) > 0 && availabilityIndex(ranges[len(ranges)-1].End) >= last-1
						iconButton(icons.Plus(), "Add range", addDisabled, func() {
							st.NextID++
							next := availabilityRange{ID: st.NextID, Start: "9:00 AM", End: "5:00 PM"}
							if n := len(st.Days[day]); n > 0 {
								s := min(availabilityIndex(st.Days[day][n-1].End)+4, last-1)
								next.Start, next.End = availabilityTimes[s], availabilityTimes[min(s+4, last)]
							}
							st.Days[day] = append(st.Days[day], next)
						})
						copyOpen := ui.Local(row, "copyOpen", func() bool { return false })
						copyTargets := ui.Local(row, "copyTargets", func() map[string]bool { return map[string]bool{} })
						cp := iconButton(icons.Copy(), "Copy to other days", len(ranges) == 0, func() {
							*copyOpen = !*copyOpen
							*copyTargets = map[string]bool{}
						})
						ui.PopoverBase(c, cp, copyOpen, func(panel ui.Element) {
							panel.Width(176).Padding(16).Gap(12).
								Radius(t.RadiusLG).Background(t.Popover).Border(1, t.Border).
								Shadow(0, 8, 20, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.14)).
								Children(func() {
									ui.Text(c, "Copy times to").FontSize(14).FontWeight(500).TextColor(t.Foreground)
									coss.CheckboxGroup(c, func() {
										for _, target := range availabilityDays {
											if target == day {
												continue
											}
											v := (*copyTargets)[target]
											coss.Checkbox(c, &v, target)
											(*copyTargets)[target] = v
										}
									})
									chosen := false
									for _, v := range *copyTargets {
										chosen = chosen || v
									}
									coss.Button(c, "Apply").Size(coss.SizeSM).Disabled(!chosen).OnClick(func() {
										for target, v := range *copyTargets {
											if !v {
												continue
											}
											cp := make([]availabilityRange, len(ranges))
											for k, rr := range ranges {
												st.NextID++
												cp[k] = availabilityRange{ID: st.NextID, Start: rr.Start, End: rr.End}
											}
											st.Days[target] = cp
										}
										*copyOpen = false
									})
								})
						})
					})
				})
			})
		}
	})
}

func PSwitch7(c *ui.Context) { renderAvailabilityEditor(c, availStylePickers) }

func PSwitch8(c *ui.Context) { renderAvailabilityEditor(c, availStyleGroupedArrow) }

func PSwitch9(c *ui.Context) { renderAvailabilityEditor(c, availStyleGroupedFromTo) }
