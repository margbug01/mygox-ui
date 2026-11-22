package particles

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed number_otp.go
var numberOTPSource string

func init() {
	Register(Entry{ID: "p-number-field-1", Family: "number-field", Description: "Basic number field", View: PNumberField1, Source: extractFunc(numberOTPSource, "PNumberField1")})
	Register(Entry{ID: "p-number-field-2", Family: "number-field", Description: "Small number field", View: PNumberField2, Source: extractFunc(numberOTPSource, "PNumberField2")})
	Register(Entry{ID: "p-number-field-3", Family: "number-field", Description: "Large number field", View: PNumberField3, Source: extractFunc(numberOTPSource, "PNumberField3")})
	Register(Entry{ID: "p-number-field-4", Family: "number-field", Description: "Disabled number field", View: PNumberField4, Source: extractFunc(numberOTPSource, "PNumberField4")})
	Register(Entry{ID: "p-number-field-5", Family: "number-field", Description: "Number field with label", View: PNumberField5, Source: extractFunc(numberOTPSource, "PNumberField5")})
	Register(Entry{ID: "p-number-field-6", Family: "number-field", Description: "Number field with scrub label", View: PNumberField6, Source: extractFunc(numberOTPSource, "PNumberField6")})
	Register(Entry{ID: "p-number-field-7", Family: "number-field", Description: "Number field with min/max range", View: PNumberField7, Source: extractFunc(numberOTPSource, "PNumberField7")})
	Register(Entry{ID: "p-number-field-8", Family: "number-field", Description: "Number field with decimal formatting", View: PNumberField8, Source: extractFunc(numberOTPSource, "PNumberField8")})
	Register(Entry{ID: "p-number-field-9", Family: "number-field", Description: "Number field with step 5", View: PNumberField9, Source: extractFunc(numberOTPSource, "PNumberField9")})
	Register(Entry{ID: "p-number-field-10", Family: "number-field", Description: "Number field in a form", View: PNumberField10, Source: extractFunc(numberOTPSource, "PNumberField10")})
	Register(Entry{ID: "p-number-field-11", Family: "number-field", Description: "Pill-shaped number field", View: PNumberField11, Source: extractFunc(numberOTPSource, "PNumberField11")})

	Register(Entry{ID: "p-otp-field-1", Family: "otp-field", Description: "Basic 6-digit OTP field", View: POTPField1, Source: extractFunc(numberOTPSource, "POTPField1")})
	Register(Entry{ID: "p-otp-field-2", Family: "otp-field", Description: "Large OTP field", View: POTPField2, Source: extractFunc(numberOTPSource, "POTPField2")})
	Register(Entry{ID: "p-otp-field-3", Family: "otp-field", Description: "OTP field with dash separator", View: POTPField3, Source: extractFunc(numberOTPSource, "POTPField3")})
	Register(Entry{ID: "p-otp-field-4", Family: "otp-field", Description: "OTP field with label", View: POTPField4, Source: extractFunc(numberOTPSource, "POTPField4")})
	Register(Entry{ID: "p-otp-field-6", Family: "otp-field", Description: "OTP field with sanitization", View: POTPField6, Source: extractFunc(numberOTPSource, "POTPField6")})
	Register(Entry{ID: "p-otp-field-7", Family: "otp-field", Description: "OTP field with auto validation", View: POTPField7, Source: extractFunc(numberOTPSource, "POTPField7")})
	Register(Entry{ID: "p-otp-field-8", Family: "otp-field", Description: "Alphanumeric OTP field", View: POTPField8, Source: extractFunc(numberOTPSource, "POTPField8")})
	Register(Entry{ID: "p-otp-field-9", Family: "otp-field", Description: "OTP field with placeholder hints", View: POTPField9, Source: extractFunc(numberOTPSource, "POTPField9")})
	Register(Entry{ID: "p-otp-field-10", Family: "otp-field", Description: "Masked PIN OTP field", View: POTPField10, Source: extractFunc(numberOTPSource, "POTPField10")})
}

func PNumberField1(c *ui.Context) {
	host := ui.Box(c)
	v := ui.Local(host, "v", func() float64 { return 0 })
	host.Children(func() { coss.NumberField(c, coss.NumberFieldSpec{Value: v}) })
}

func PNumberField2(c *ui.Context) {
	host := ui.Box(c)
	v := ui.Local(host, "v", func() float64 { return 1 })
	host.Children(func() { coss.NumberField(c, coss.NumberFieldSpec{Value: v}).Height(28) })
}

func PNumberField3(c *ui.Context) {
	host := ui.Box(c)
	v := ui.Local(host, "v", func() float64 { return 1 })
	host.Children(func() { coss.NumberField(c, coss.NumberFieldSpec{Value: v}).Height(38) })
}

func PNumberField4(c *ui.Context) {
	host := ui.Box(c)
	v := ui.Local(host, "v", func() float64 { return 5 })
	host.Children(func() { coss.NumberField(c, coss.NumberFieldSpec{Value: v, Disabled: true}) })
}

func PNumberField5(c *ui.Context) {
	host := ui.Box(c)
	v := ui.Local(host, "v", func() float64 { return 2 })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Quantity")
			coss.NumberField(c, coss.NumberFieldSpec{Value: v, Min: 0, Max: 99})
		})
	})
}

func PNumberField6(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).Gap(6)
	v := ui.Local(host, "v", func() float64 { return 24 })
	host.Children(func() {
		ui.Text(c, "Font size (px)").FontSize(13).FontWeight(500).TextColor(t.Foreground)
		coss.NumberField(c, coss.NumberFieldSpec{Value: v, Min: 8, Max: 72})
	})
}

func PNumberField7(c *ui.Context) {
	host := ui.Box(c)
	v := ui.Local(host, "v", func() float64 { return 5 })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Between 0 and 10")
			coss.NumberField(c, coss.NumberFieldSpec{Value: v, Min: 0, Max: 10})
		})
	})
}

func PNumberField8(c *ui.Context) {
	host := ui.Box(c)
	v := ui.Local(host, "v", func() float64 { return 12.5 })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Price ($)")
			coss.NumberField(c, coss.NumberFieldSpec{Value: v, Min: 0, Max: 999, Step: 0.5, Decimals: 2})
		})
	})
}

func PNumberField9(c *ui.Context) {
	host := ui.Box(c)
	v := ui.Local(host, "v", func() float64 { return 10 })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Step by 5")
			coss.NumberField(c, coss.NumberFieldSpec{Value: v, Min: 0, Max: 100, Step: 5})
		})
	})
}

func PNumberField10(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	v := ui.Local(host, "v", func() float64 { return 3 })
	msg := ui.Local(host, "msg", func() string { return "" })
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() {
				coss.FieldLabel(c, "Seats")
				coss.NumberField(c, coss.NumberFieldSpec{Value: v, Min: 1, Max: 50})
			})
			coss.Button(c, "Confirm").Size(coss.SizeSM).OnClick(func() { *msg = fmt.Sprintf("Reserved %.0f seats", *v) })
			if *msg != "" {
				ui.Text(c, *msg).FontSize(12).TextColor(t.SuccessForeground)
			}
		})
	})
}

func PNumberField11(c *ui.Context) {
	host := ui.Box(c)
	v := ui.Local(host, "v", func() float64 { return 1 })
	host.Children(func() { coss.NumberField(c, coss.NumberFieldSpec{Value: v}).Radius(999) })
}

func POTPField1(c *ui.Context) {
	host := ui.Box(c)
	code := ui.Local(host, "code", func() string { return "" })
	host.Children(func() { coss.OTPField(c, code, 6) })
}

func POTPField2(c *ui.Context) {
	host := ui.Box(c)
	code := ui.Local(host, "code", func() string { return "12" })
	host.Children(func() { coss.OTPField(c, code, 6).Height(46) })
}

func POTPField3(c *ui.Context) {
	host := ui.Box(c)
	code := ui.Local(host, "code", func() string { return "482" })
	host.Children(func() { coss.OTPField(c, code, 6, 3) })
}

func POTPField4(c *ui.Context) {
	host := ui.Box(c)
	code := ui.Local(host, "code", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Verification code")
			coss.OTPField(c, code, 6, 3)
			coss.FieldDescription(c, "Enter the 6-digit code sent to your email.")
		})
	})
}

func POTPField6(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).Gap(6)
	code := ui.Local(host, "code", func() string { return "A1B2" })
	host.Children(func() {
		*code = strings.ToUpper(*code)
		coss.OTPField(c, code, 6, 3)
		ui.Text(c, "Sanitized: "+*code).FontSize(12).TextColor(t.MutedForeground)
	})
}

func POTPField7(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).Gap(8)
	code := ui.Local(host, "code", func() string { return "123456" })
	host.Children(func() {
		coss.OTPField(c, code, 6, 3)
		if len(*code) == 6 {
			if *code == "123456" {
				ui.Text(c, "Code verified").FontSize(12).TextColor(t.SuccessForeground)
			} else {
				coss.FieldError(c, "Invalid code (try 123456)")
			}
		}
	})
}

func POTPField8(c *ui.Context) {
	host := ui.Box(c)
	code := ui.Local(host, "code", func() string { return "G0X9" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Alphanumeric invite code")
			coss.OTPField(c, code, 4)
		})
	})
}

func POTPField9(c *ui.Context) {
	host := ui.Box(c)
	code := ui.Local(host, "code", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "One-time passcode")
			coss.OTPField(c, code, 6)
			coss.FieldDescription(c, "Format: 000000")
		})
	})
}

func POTPField10(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).Gap(6)
	code := ui.Local(host, "code", func() string { return "8492" })
	host.Children(func() {
		coss.FieldLabel(c, "Security PIN")
		coss.OTPField(c, code, 4)
		ui.Text(c, "Masked: "+strings.Repeat("•", len(*code))).FontSize(12).TextColor(t.MutedForeground)
	})
}
