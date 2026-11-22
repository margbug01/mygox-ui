package particles

import (
	_ "embed"
	"fmt"
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
}

func compactSwitch(c *ui.Context, on *bool, a11yLabel string) ui.Element {
	t := theme.Of(c)
	sw := ui.SwitchBase(c, on).
		AlignItems(ui.Center).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring).
		Label(a11yLabel)
	isOn := on != nil && *on
	var target float32
	if isOn {
		target = 1
	}
	prog := sw.Animate("compact-sw", target, 140*time.Millisecond)
	sw.Children(func() {
		trackColor := theme.Mix(t.Input, t.Primary, prog)
		if !t.Dark && !isOn {
			trackColor = theme.Alpha(ui.RGB(0, 0, 0), 0.16)
		}
		ui.Row(c).
			AlignItems(ui.Center).
			Size(26, 16).
			Padding(2).
			Radius(999).
			Background(trackColor).
			Shrink(0).
			Children(func() {
				thumbColor := ui.RGB(255, 255, 255)
				if t.Dark && isOn {
					thumbColor = t.PrimaryForeground
				}
				ui.Box(c).
					Size(12, 12).
					Radius(999).
					Background(thumbColor).
					Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.18)).
					Translate(10*prog, 0)
			})
	})
	return sw
}

func PSwitch1(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() {
		coss.Switch(c, on, "Marketing emails")
	})
}

func PSwitch2(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() {
		coss.Switch(c, on, "Marketing emails").Disabled(true).Opacity(0.64)
	})
}

func PSwitch3(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).AlignItems(ui.Start).Gap(10).MaxWidth(360)
	on := ui.Local(host, "on", func() bool { return true })
	host.Children(func() {
		coss.Switch(c, on).Margin(2, 0, 0, 0)
		ui.Column(c).Gap(4).Children(func() {
			coss.Label(c, "Marketing emails")
			ui.Text(c, "By enabling marketing emails, you agree to receive emails.").
				FontSize(12).TextColor(t.MutedForeground)
		})
	})
}

func PSwitch4(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).MaxWidth(380)
	on := ui.Local(host, "on", func() bool { return true })
	host.Children(func() {
		card := ui.Row(c).
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween).
			Gap(20).
			Padding(12).
			Radius(t.RadiusLG).
			Border(1, t.Border)
		if *on {
			card.Border(1, theme.Alpha(t.Primary, 0.48)).Background(theme.Alpha(t.Accent, 0.5))
		}
		card.Children(func() {
			ui.Column(c).Gap(4).Children(func() {
				ui.Text(c, "Enable notifications").FontSize(14).FontWeight(500).TextColor(t.Foreground)
				ui.Text(c, "You can enable or disable notifications at any time.").
					FontSize(12).TextColor(t.MutedForeground)
			})
			compactSwitch(c, on, "Enable notifications")
		})
	})
}

func PSwitch5(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return true })
	status := ui.Local(host, "status", func() string { return "" })
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() {
				coss.Switch(c, on, "Enable marketing emails")
			})
			coss.Button(c, "Submit").OnClick(func() {
				if *on {
					*status = "Marketing emails: on"
				} else {
					*status = "Marketing emails: off"
				}
			})
			if *status != "" {
				ui.Text(c, *status).FontSize(12).TextColor(t.MutedForeground)
			}
		})
	})
}

func PSwitch6(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() {
		compactSwitch(c, on, "Notifications")
	})
}

type availabilityRange struct {
	ID    int
	Start string
	End   string
}

type availabilityState struct {
	NextID int
	Days   map[string][]availabilityRange
}

var availabilityDays = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

var availabilityTimeOpts = func() []selectOption {
	opts := make([]selectOption, 0, 48)
	for i := 0; i < 48; i++ {
		h := i / 2
		m := (i % 2) * 30
		period := "AM"
		if h >= 12 {
			period = "PM"
		}
		dh := h % 12
		if dh == 0 {
			dh = 12
		}
		s := fmt.Sprintf("%d:%02d %s", dh, m, period)
		opts = append(opts, selectOption{Label: s, Value: s})
	}
	return opts
}()

func newAvailabilityState() *availabilityState {
	return &availabilityState{
		NextID: 10,
		Days: map[string][]availabilityRange{
			"Monday":    {{ID: 1, Start: "9:00 AM", End: "5:00 PM"}},
			"Tuesday":   {{ID: 2, Start: "9:00 AM", End: "1:00 PM"}, {ID: 3, Start: "3:00 PM", End: "5:00 PM"}},
			"Wednesday": {{ID: 4, Start: "9:00 AM", End: "5:00 PM"}},
			"Thursday":  {{ID: 5, Start: "9:00 AM", End: "5:00 PM"}},
			"Friday":    {{ID: 6, Start: "9:00 AM", End: "5:00 PM"}},
			"Saturday":  nil,
			"Sunday":    nil,
		},
	}
}

type availabilityStyle uint8

const (
	availStylePickers availabilityStyle = iota
	availStyleGroupedArrow
	availStyleGroupedFromTo
)

func renderAvailabilityEditor(c *ui.Context, mode availabilityStyle) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).MaxWidth(520)
	stPtr := ui.Local(host, "avail", newAvailabilityState)
	st := *stPtr

	host.Children(func() {
		for idx, day := range availabilityDays {
			dayName := day
			ranges := st.Days[dayName]
			row := ui.Row(c).WidthPercent(100).AlignItems(ui.Start).Gap(12).Padding(10, 0)
			if idx > 0 {
				row.BorderWidth(1, 0, 0, 0).BorderColor(t.Border)
			}
			row.Children(func() {
				enabled := len(ranges) > 0
				prevEnabled := enabled
				ui.Row(c).Width(116).Height(28).AlignItems(ui.Center).Shrink(0).Children(func() {
					coss.Switch(c, &enabled, dayName)
				})
				if enabled != prevEnabled {
					if enabled {
						st.NextID++
						st.Days[dayName] = []availabilityRange{{ID: st.NextID, Start: "9:00 AM", End: "5:00 PM"}}
					} else {
						st.Days[dayName] = nil
					}
					ranges = st.Days[dayName]
				}

				ui.Column(c).Grow(1).MinWidth(0).Gap(6).Children(func() {
					if len(ranges) == 0 {
						ui.Row(c).Height(28).AlignItems(ui.Center).Children(func() {
							ui.Text(c, "Unavailable").FontSize(13).TextColor(t.MutedForeground)
						})
					} else {
						for rIdx := range ranges {
							r := &st.Days[dayName][rIdx]
							ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
								switch mode {
								case availStylePickers:
									cossSelect(c, &r.Start, availabilityTimeOpts, 94).Height(28)
									ui.Text(c, "–").FontSize(13).TextColor(t.MutedForeground)
									cossSelect(c, &r.End, availabilityTimeOpts, 94).Height(28)
								case availStyleGroupedArrow:
									coss.Group(c, func() {
										cossSelect(c, &r.Start, availabilityTimeOpts, 88).Height(28).BorderWidth(0).Radius(0)
										coss.GroupSeparator(c)
										ui.Row(c).Center().Height(28).Padding(0, 7).Background(t.Muted).TextColor(t.MutedForeground).Children(func() {
											ui.Icon(c, icons.ArrowRight()).Size(13, 13)
										})
										coss.GroupSeparator(c)
										cossSelect(c, &r.End, availabilityTimeOpts, 88).Height(28).BorderWidth(0).Radius(0)
									})
								case availStyleGroupedFromTo:
									coss.Group(c, func() {
										coss.GroupText(c, "From").Height(28).Padding(0, 8)
										coss.GroupSeparator(c)
										cossSelect(c, &r.Start, availabilityTimeOpts, 84).Height(28).BorderWidth(0).Radius(0)
										coss.GroupSeparator(c)
										coss.GroupText(c, "To").Height(28).Padding(0, 8)
										coss.GroupSeparator(c)
										cossSelect(c, &r.End, availabilityTimeOpts, 84).Height(28).BorderWidth(0).Radius(0)
									})
								}
								rmBtn := coss.IconButton(c, icons.X(), "Delete range").
									Variant(coss.ButtonGhost).Size(coss.SizeIconSM).
									OnClick(func() {
										cur := st.Days[dayName]
										next := make([]availabilityRange, 0, len(cur))
										for _, item := range cur {
											if item.ID != r.ID {
												next = append(next, item)
											}
										}
										st.Days[dayName] = next
									})
								cossTooltip(c, rmBtn.Element, "Delete range")
							})
						}
					}
				})

				ui.Row(c).AlignItems(ui.Center).Gap(2).Shrink(0).Children(func() {
					addBtn := coss.IconButton(c, icons.Plus(), "Add range").
						Variant(coss.ButtonGhost).Size(coss.SizeIconSM).
						OnClick(func() {
							st.NextID++
							st.Days[dayName] = append(st.Days[dayName], availabilityRange{
								ID:    st.NextID,
								Start: "1:00 PM",
								End:   "5:00 PM",
							})
						})
					cossTooltip(c, addBtn.Element, "Add range")

					copyOpen := ui.Local(row, "copyOpen-"+dayName, func() bool { return false })
					copyTargets := ui.Local(row, "copyTargets-"+dayName, func() map[string]bool { return map[string]bool{} })
					cpBtn := coss.IconButton(c, icons.Copy(), "Copy to other days").
						Variant(coss.ButtonGhost).Size(coss.SizeIconSM).
						Disabled(len(ranges) == 0).
						OnClick(func() { *copyOpen = !*copyOpen })
					cossTooltip(c, cpBtn.Element, "Copy to other days")
					ui.PopoverBase(c, cpBtn.Element, copyOpen, func(panel ui.Element) {
						panel.Width(176).Padding(12).Gap(10).
							Radius(t.RadiusLG).Background(t.Popover).Border(1, t.Border).
							Shadow(0, 8, 20, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.14)).
							Children(func() {
								ui.Text(c, "Copy times to").FontSize(13).FontWeight(500).TextColor(t.Foreground)
								coss.CheckboxGroup(c, func() {
									for _, target := range availabilityDays {
										if target == dayName {
											continue
										}
										tgt := target
										v := (*copyTargets)[tgt]
										coss.Checkbox(c, &v, tgt)
										(*copyTargets)[tgt] = v
									}
								})
								coss.Button(c, "Apply").Size(coss.SizeSM).OnClick(func() {
									for _, target := range availabilityDays {
										if (*copyTargets)[target] {
											cp := make([]availabilityRange, len(ranges))
											for k, rr := range ranges {
												st.NextID++
												cp[k] = availabilityRange{ID: st.NextID, Start: rr.Start, End: rr.End}
											}
											st.Days[target] = cp
										}
									}
									*copyTargets = map[string]bool{}
									*copyOpen = false
								})
							})
					})
				})
			})
		}
	})
}

func PSwitch7(c *ui.Context) {
	renderAvailabilityEditor(c, availStylePickers)
}

func PSwitch8(c *ui.Context) {
	renderAvailabilityEditor(c, availStyleGroupedArrow)
}

func PSwitch9(c *ui.Context) {
	renderAvailabilityEditor(c, availStyleGroupedFromTo)
}
