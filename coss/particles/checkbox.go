package particles

import (
	_ "embed"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
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

// disabledCheckbox is coss's <Label><Checkbox disabled />…</Label>: the
// box fades to 64%, the label stays solid.
func disabledCheckbox(c *ui.Context, checked *bool, label string) {
	state := ui.CheckOff
	if *checked {
		state = ui.CheckOn
	}
	ui.CheckboxBase(c, checked).Disabled(true).DimDisabled(false).AlignItems(ui.Center).Gap(8).Children(func() {
		coss.CheckboxBox(c, state, true)
		coss.CheckboxLabel(c, label)
	})
}

// groupState is a parent checkbox's state over its group's boxes.
func groupState(boxes ...*bool) ui.CheckState {
	n := 0
	for _, b := range boxes {
		if *b {
			n++
		}
	}
	switch n {
	case 0:
		return ui.CheckOff
	case len(boxes):
		return ui.CheckOn
	}
	return ui.CheckMixed
}

// setAll checks or clears every box.
func setAll(on bool, boxes ...*bool) {
	for _, b := range boxes {
		*b = on
	}
}

// formLoading starts an 800ms submit, as coss's form particles fake one,
// and reports whether one is running.
func formLoading(c *ui.Context, host ui.Element) (loading bool, start func()) {
	until := ui.Local(host, "loading", func() time.Time { return time.Time{} })
	loading = c.Now().Before(*until)
	if loading {
		c.After(150 * time.Millisecond)
	}
	return loading, func() { *until = c.Now().Add(800 * time.Millisecond) }
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
	host.Children(func() { disabledCheckbox(c, checked, "Accept terms and conditions") })
}

func PCheckbox3(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).AlignItems(ui.Start).Gap(8)
	checked := ui.Local(host, "checked", func() bool { return true })
	host.Children(func() {
		coss.Checkbox(c, checked).Label("Accept terms and conditions")
		ui.Column(c).Gap(4).Children(func() {
			coss.CheckboxLabel(c, "Accept terms and conditions")
			ui.Text(c, "By clicking this checkbox, you agree to the terms and conditions.").
				FontSize(12).FixedLineHeight(16).TextColor(t.MutedForeground)
		})
	})
}

func PCheckbox4(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	checked := ui.Local(host, "checked", func() bool { return true })
	host.Children(func() {
		card := ui.CheckboxBase(c, checked).AlignItems(ui.Start).Gap(8).Padding(12).
			Radius(t.RadiusLG).Border(1, t.Border).Cursor(ui.CursorPointer).FocusRingColor(t.Ring)
		state := ui.CheckOff
		if *checked {
			state = ui.CheckOn
			card.Border(1, theme.Alpha(t.Primary, 0.48)).Background(theme.Alpha(t.Accent, 0.5))
		} else if card.Hovered() {
			card.Background(theme.Alpha(t.Accent, 0.5))
		}
		card.Children(func() {
			coss.CheckboxBox(c, state, false)
			ui.Column(c).Gap(4).Children(func() {
				coss.CheckboxLabel(c, "Enable notifications")
				ui.Text(c, "You can enable or disable notifications at any time.").
					FontSize(12).FixedLineHeight(16).TextColor(t.MutedForeground)
			})
		})
	})
}

func PCheckbox5(c *ui.Context) {
	host := ui.Box(c)
	checked := ui.Local(host, "checked", func() bool { return true })
	loading, submit := formLoading(c, host)
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() {
				coss.Checkbox(c, checked, "Accept terms and conditions")
			})
			coss.Button(c, "Submit").Loading(loading).OnClick(submit)
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
			disabledCheckbox(c, vite, "Vite")
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
		coss.CheckboxGroup(c, func() {
			all := groupState(next, vite, astro)
			if coss.CheckboxParent(c, &all, "Frameworks").Changed() {
				setAll(all == ui.CheckOn, next, vite, astro)
			}
			coss.Checkbox(c, next, "Next.js").Margin(0, 0, 0, 16)
			coss.Checkbox(c, vite, "Vite").Margin(0, 0, 0, 16)
			coss.Checkbox(c, astro, "Astro").Margin(0, 0, 0, 16)
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
		management := []*bool{createUsr, editUsr, deleteUsr, assignRol}
		manage := groupState(management...)
		// "Manage Users" counts toward the main group once every user
		// management permission is on; the parent is mixed while some are.
		manageUsers := manage == ui.CheckOn
		root := groupState(viewDash, accessRep, &manageUsers)
		if manage == ui.CheckMixed {
			root = ui.CheckMixed
		}
		coss.CheckboxGroup(c, func() {
			if coss.CheckboxParent(c, &root, "User Permissions").Changed() {
				on := root == ui.CheckOn
				setAll(on, viewDash, accessRep)
				setAll(on, management...)
			}
			coss.Checkbox(c, viewDash, "View Dashboard").Margin(0, 0, 0, 16)
			coss.Checkbox(c, accessRep, "Access Reports").Margin(0, 0, 0, 16)
			coss.CheckboxGroup(c, func() {
				if coss.CheckboxParent(c, &manage, "Manage Users").Changed() {
					setAll(manage == ui.CheckOn, management...)
				}
				coss.Checkbox(c, createUsr, "Create User").Margin(0, 0, 0, 16)
				coss.Checkbox(c, editUsr, "Edit User").Margin(0, 0, 0, 16)
				coss.Checkbox(c, deleteUsr, "Delete User").Margin(0, 0, 0, 16)
				coss.Checkbox(c, assignRol, "Assign Roles").Margin(0, 0, 0, 16)
			}).Margin(0, 0, 0, 16)
		})
	})
}

func PCheckboxGroup5(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(160)
	next := ui.Local(host, "next", func() bool { return true })
	vite := ui.Local(host, "vite", func() bool { return false })
	astro := ui.Local(host, "astro", func() bool { return false })
	loading, submit := formLoading(c, host)
	host.Children(func() {
		coss.Form(c, func() {
			coss.Fieldset(c, func() {
				ui.Text(c, "Frameworks").FontSize(14).FixedLineHeight(20).FontWeight(500).TextColor(t.Foreground)
				coss.CheckboxGroup(c, func() {
					coss.Checkbox(c, next, "Next.js")
					coss.Checkbox(c, vite, "Vite")
					coss.Checkbox(c, astro, "Astro")
				})
			}).Gap(8)
			coss.Button(c, "Submit").Loading(loading).AlignSelf(ui.Start).OnClick(submit)
		})
	})
}
