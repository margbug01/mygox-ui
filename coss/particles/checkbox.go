package particles

import (
	_ "embed"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed checkbox.go
var checkboxSource string

func init() {
	Register(Entry{ID: "p-checkbox-1", Family: "checkbox", Description: "Basic checkbox", View: PCheckbox1, Source: extractFunc(checkboxSource, "PCheckbox1")})
	Register(Entry{ID: "p-checkbox-2", Family: "checkbox", Description: "Disabled checkbox", View: PCheckbox2, Source: extractFunc(checkboxSource, "PCheckbox2")})
	Register(Entry{ID: "p-checkbox-3", Family: "checkbox", Description: "Checkbox with description", View: PCheckbox3, Source: extractFunc(checkboxSource, "PCheckbox3")})
	Register(Entry{ID: "p-checkbox-4", Family: "checkbox", Description: "Card-style checkbox", View: PCheckbox4, Source: extractFunc(checkboxSource, "PCheckbox4")})
	Register(Entry{ID: "p-checkbox-5", Family: "checkbox", Description: "Checkbox form", View: PCheckbox5, Source: extractFunc(checkboxSource, "PCheckbox5")})

	Register(Entry{ID: "p-checkbox-group-1", Family: "checkbox-group", Description: "Basic checkbox group", View: PCheckboxGroup1, Source: extractFunc(checkboxSource, "PCheckboxGroup1")})
	Register(Entry{ID: "p-checkbox-group-2", Family: "checkbox-group", Description: "Checkbox group with disabled items", View: PCheckboxGroup2, Source: extractFunc(checkboxSource, "PCheckboxGroup2")})
	Register(Entry{ID: "p-checkbox-group-3", Family: "checkbox-group", Description: "Checkbox group with parent checkbox", View: PCheckboxGroup3, Source: extractFunc(checkboxSource, "PCheckboxGroup3")})
	Register(Entry{ID: "p-checkbox-group-4", Family: "checkbox-group", Description: "Nested checkbox group with parent", View: PCheckboxGroup4, Source: extractFunc(checkboxSource, "PCheckboxGroup4")})
	Register(Entry{ID: "p-checkbox-group-5", Family: "checkbox-group", Description: "Checkbox group form", View: PCheckboxGroup5, Source: extractFunc(checkboxSource, "PCheckboxGroup5")})
}

func triStateCheckbox(c *ui.Context, state ui.CheckState, label string, onToggle func()) ui.Element {
	t := theme.Of(c)
	btn := ui.ButtonBase(c).
		AlignItems(ui.Center).
		Gap(8).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring)
	btn.Children(func() {
		box := ui.Row(c).Center().Size(16, 16).Radius(4).Shrink(0)
		switch state {
		case ui.CheckOn:
			box.Background(t.Primary).Border(1, t.Primary).TextColor(t.PrimaryForeground).Children(func() {
				ui.Icon(c, icons.Check()).Size(12, 12)
			})
		case ui.CheckMixed:
			box.Background(t.Primary).Border(1, t.Primary).TextColor(t.PrimaryForeground).Children(func() {
				ui.Icon(c, icons.Minus()).Size(12, 12)
			})
		default:
			bg := t.Background
			if t.Dark {
				bg = theme.Alpha(t.Input, 0.32)
			}
			box.Background(bg).Border(1, t.Input)
		}
		ui.Text(c, label).FontSize(14).FontWeight(500).TextColor(t.Foreground)
	})
	if btn.Clicked() && onToggle != nil {
		onToggle()
	}
	return btn
}

func PCheckbox1(c *ui.Context) {
	host := ui.Box(c)
	checked := ui.Local(host, "checked", func() bool { return false })
	host.Children(func() {
		coss.Checkbox(c, checked, "Accept terms and conditions")
	})
}

func PCheckbox2(c *ui.Context) {
	host := ui.Box(c)
	checked := ui.Local(host, "checked", func() bool { return true })
	host.Children(func() {
		coss.Checkbox(c, checked, "Accept terms and conditions").Disabled(true).Opacity(0.64)
	})
}

func PCheckbox3(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).AlignItems(ui.Start).Gap(10).MaxWidth(360)
	checked := ui.Local(host, "checked", func() bool { return true })
	host.Children(func() {
		coss.Checkbox(c, checked).Margin(2, 0, 0, 0)
		ui.Column(c).Gap(4).Children(func() {
			coss.Label(c, "Accept terms and conditions")
			ui.Text(c, "By clicking this checkbox, you agree to the terms and conditions.").
				FontSize(12).TextColor(t.MutedForeground)
		})
	})
}

func PCheckbox4(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).MaxWidth(360)
	checked := ui.Local(host, "checked", func() bool { return true })
	host.Children(func() {
		card := ui.ButtonBase(c).
			AlignItems(ui.Start).
			Gap(10).
			Padding(12).
			Radius(t.RadiusLG).
			Border(1, t.Border).
			Cursor(ui.CursorPointer)
		if *checked {
			card.Border(1, theme.Alpha(t.Primary, 0.48)).Background(theme.Alpha(t.Accent, 0.5))
		} else if card.Hovered() {
			card.Background(theme.Alpha(t.Accent, 0.5))
		}
		card.Children(func() {
			coss.Checkbox(c, checked).Margin(2, 0, 0, 0)
			ui.Column(c).Gap(4).Children(func() {
				ui.Text(c, "Enable notifications").FontSize(14).FontWeight(500).TextColor(t.Foreground)
				ui.Text(c, "You can enable or disable notifications at any time.").
					FontSize(12).TextColor(t.MutedForeground)
			})
		})
		if card.Clicked() {
			*checked = !*checked
		}
	})
}

func PCheckbox5(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	checked := ui.Local(host, "checked", func() bool { return true })
	submitted := ui.Local(host, "submitted", func() string { return "" })
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() {
				coss.Checkbox(c, checked, "Accept terms and conditions")
			})
			coss.Button(c, "Submit").OnClick(func() {
				if *checked {
					*submitted = "Terms: yes"
				} else {
					*submitted = "Terms: no"
				}
			})
			if *submitted != "" {
				ui.Text(c, *submitted).FontSize(12).TextColor(t.MutedForeground)
			}
		})
	})
}

func PCheckboxGroup1(c *ui.Context) {
	host := ui.Box(c)
	next := ui.Local(host, "next", func() bool { return true })
	vite := ui.Local(host, "vite", func() bool { return false })
	astro := ui.Local(host, "astro", func() bool { return false })
	host.Children(func() {
		coss.CheckboxGroup(c, func() {
			coss.Checkbox(c, next, "Next.js")
			coss.Checkbox(c, vite, "Vite")
			coss.Checkbox(c, astro, "Astro")
		})
	})
}

func PCheckboxGroup2(c *ui.Context) {
	host := ui.Box(c)
	next := ui.Local(host, "next", func() bool { return true })
	vite := ui.Local(host, "vite", func() bool { return false })
	astro := ui.Local(host, "astro", func() bool { return false })
	host.Children(func() {
		coss.CheckboxGroup(c, func() {
			coss.Checkbox(c, next, "Next.js")
			coss.Checkbox(c, vite, "Vite").Disabled(true).Opacity(0.64)
			coss.Checkbox(c, astro, "Astro")
		})
	})
}

func PCheckboxGroup3(c *ui.Context) {
	host := ui.Box(c)
	next := ui.Local(host, "next", func() bool { return false })
	vite := ui.Local(host, "vite", func() bool { return false })
	astro := ui.Local(host, "astro", func() bool { return false })
	host.Children(func() {
		count := 0
		if *next {
			count++
		}
		if *vite {
			count++
		}
		if *astro {
			count++
		}
		parentState := ui.CheckOff
		if count == 3 {
			parentState = ui.CheckOn
		} else if count > 0 {
			parentState = ui.CheckMixed
		}
		coss.CheckboxGroup(c, func() {
			triStateCheckbox(c, parentState, "Frameworks", func() {
				v := count != 3
				*next, *vite, *astro = v, v, v
			})
			ui.Column(c).Gap(10).Margin(0, 0, 0, 16).Children(func() {
				coss.Checkbox(c, next, "Next.js")
				coss.Checkbox(c, vite, "Vite")
				coss.Checkbox(c, astro, "Astro")
			})
		})
	})
}

func PCheckboxGroup4(c *ui.Context) {
	host := ui.Box(c)
	viewDash := ui.Local(host, "viewDash", func() bool { return false })
	accessRep := ui.Local(host, "accessRep", func() bool { return false })
	createUsr := ui.Local(host, "createUsr", func() bool { return false })
	editUsr := ui.Local(host, "editUsr", func() bool { return false })
	deleteUsr := ui.Local(host, "deleteUsr", func() bool { return false })
	assignRol := ui.Local(host, "assignRol", func() bool { return false })
	host.Children(func() {
		mgmtCount := 0
		for _, b := range []*bool{createUsr, editUsr, deleteUsr, assignRol} {
			if *b {
				mgmtCount++
			}
		}
		mgmtState := ui.CheckOff
		if mgmtCount == 4 {
			mgmtState = ui.CheckOn
		} else if mgmtCount > 0 {
			mgmtState = ui.CheckMixed
		}
		totalChecked := mgmtCount
		if *viewDash {
			totalChecked++
		}
		if *accessRep {
			totalChecked++
		}
		rootState := ui.CheckOff
		if totalChecked == 6 {
			rootState = ui.CheckOn
		} else if totalChecked > 0 {
			rootState = ui.CheckMixed
		}
		coss.CheckboxGroup(c, func() {
			triStateCheckbox(c, rootState, "User Permissions", func() {
				v := totalChecked != 6
				*viewDash, *accessRep = v, v
				*createUsr, *editUsr, *deleteUsr, *assignRol = v, v, v, v
			})
			ui.Column(c).Gap(10).Margin(0, 0, 0, 16).Children(func() {
				coss.Checkbox(c, viewDash, "View Dashboard")
				coss.Checkbox(c, accessRep, "Access Reports")
				triStateCheckbox(c, mgmtState, "Manage Users", func() {
					v := mgmtCount != 4
					*createUsr, *editUsr, *deleteUsr, *assignRol = v, v, v, v
				})
				ui.Column(c).Gap(10).Margin(0, 0, 0, 16).Children(func() {
					coss.Checkbox(c, createUsr, "Create User")
					coss.Checkbox(c, editUsr, "Edit User")
					coss.Checkbox(c, deleteUsr, "Delete User")
					coss.Checkbox(c, assignRol, "Assign Roles")
				})
			})
		})
	})
}

func PCheckboxGroup5(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(180)
	next := ui.Local(host, "next", func() bool { return true })
	vite := ui.Local(host, "vite", func() bool { return false })
	astro := ui.Local(host, "astro", func() bool { return false })
	result := ui.Local(host, "result", func() string { return "" })
	host.Children(func() {
		coss.Form(c, func() {
			coss.Fieldset(c, func() {
				coss.FieldsetLegend(c, "Frameworks")
				coss.CheckboxGroup(c, func() {
					coss.Checkbox(c, next, "Next.js")
					coss.Checkbox(c, vite, "Vite")
					coss.Checkbox(c, astro, "Astro")
				})
			})
			coss.Button(c, "Submit").OnClick(func() {
				var sel []string
				if *next {
					sel = append(sel, "next")
				}
				if *vite {
					sel = append(sel, "vite")
				}
				if *astro {
					sel = append(sel, "astro")
				}
				if len(sel) == 0 {
					*result = "Selected: none"
				} else {
					*result = "Selected: " + strings.Join(sel, ", ")
				}
			})
			if *result != "" {
				ui.Text(c, *result).FontSize(12).TextColor(t.MutedForeground)
			}
		})
	})
}
