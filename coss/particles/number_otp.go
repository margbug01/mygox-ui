package particles

import (
	_ "embed"
	"fmt"
	"time"
	"unicode"

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
	Register(Entry{ID: "p-number-field-6", Family: "number-field", Description: "Number field with scrub", View: PNumberField6, Source: extractFunc(numberOTPSource, "PNumberField6")})
	Register(Entry{ID: "p-number-field-7", Family: "number-field", Description: "Number field with range", View: PNumberField7, Source: extractFunc(numberOTPSource, "PNumberField7")})
	Register(Entry{ID: "p-number-field-8", Family: "number-field", Description: "Number field with formatted value", View: PNumberField8, Source: extractFunc(numberOTPSource, "PNumberField8")})
	Register(Entry{ID: "p-number-field-9", Family: "number-field", Description: "Number field with step", View: PNumberField9, Source: extractFunc(numberOTPSource, "PNumberField9")})
	Register(Entry{ID: "p-number-field-10", Family: "number-field", Description: "Number field in form", View: PNumberField10, Source: extractFunc(numberOTPSource, "PNumberField10")})
	Register(Entry{ID: "p-number-field-11", Family: "number-field", Description: "Pill-shaped number field", View: PNumberField11, Source: extractFunc(numberOTPSource, "PNumberField11")})

	Register(Entry{ID: "p-otp-field-1", Family: "otp-field", Description: "Basic OTP field", View: POTPField1, Source: extractFunc(numberOTPSource, "POTPField1")})
	Register(Entry{ID: "p-otp-field-2", Family: "otp-field", Description: "Large OTP field", View: POTPField2, Source: extractFunc(numberOTPSource, "POTPField2")})
	Register(Entry{ID: "p-otp-field-3", Family: "otp-field", Description: "OTP field with separator", View: POTPField3, Source: extractFunc(numberOTPSource, "POTPField3")})
	Register(Entry{ID: "p-otp-field-4", Family: "otp-field", Description: "OTP field with label", View: POTPField4, Source: extractFunc(numberOTPSource, "POTPField4")})
	Register(Entry{ID: "p-otp-field-6", Family: "otp-field", Description: "OTP field with custom sanitization", View: POTPField6, Source: extractFunc(numberOTPSource, "POTPField6")})
	Register(Entry{ID: "p-otp-field-7", Family: "otp-field", Description: "OTP field with auto validation", View: POTPField7, Source: extractFunc(numberOTPSource, "POTPField7")})
	Register(Entry{ID: "p-otp-field-8", Family: "otp-field", Description: "Alphanumeric OTP field", View: POTPField8, Source: extractFunc(numberOTPSource, "POTPField8")})
	Register(Entry{ID: "p-otp-field-9", Family: "otp-field", Description: "OTP field with placeholder hints", View: POTPField9, Source: extractFunc(numberOTPSource, "POTPField9")})
	Register(Entry{ID: "p-otp-field-10", Family: "otp-field", Description: "Masked OTP field", View: POTPField10, Source: extractFunc(numberOTPSource, "POTPField10")})
}

// numberField shows a NumberField starting at initial in the max-w-64
// column the number field particles are drawn in.
func numberField(c *ui.Context, initial float64, spec coss.NumberFieldSpec) coss.NumberFieldElement {
	var nf coss.NumberFieldElement
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() float64 { return initial })
	spec.Value = v
	host.Children(func() { nf = coss.NumberField(c, spec) })
	return nf
}

func PNumberField1(c *ui.Context) { numberField(c, 0, coss.NumberFieldSpec{}) }

func PNumberField2(c *ui.Context) { numberField(c, 0, coss.NumberFieldSpec{Size: coss.InputSM}) }

func PNumberField3(c *ui.Context) { numberField(c, 0, coss.NumberFieldSpec{Size: coss.InputLG}) }

func PNumberField4(c *ui.Context) { numberField(c, 42, coss.NumberFieldSpec{Disabled: true}) }

func PNumberField5(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).MaxWidth(256).Gap(8)
	v := ui.Local(host, "v", func() float64 { return 0 })
	host.Children(func() {
		coss.Label(c, "Quantity").LineHeight(16.0 / 14)
		coss.NumberField(c, coss.NumberFieldSpec{Value: v})
	})
}

func PNumberField6(c *ui.Context) { numberField(c, 0, coss.NumberFieldSpec{Label: "Quantity"}) }

func PNumberField7(c *ui.Context) { numberField(c, 5, coss.NumberFieldSpec{Min: 0, Max: 10}) }

// usd formats v as US dollars, as Intl's currency style.
func usd(v float64) string {
	if v < 0 {
		return "-" + usd(-v)
	}
	whole := int64(v)
	cents := fmt.Sprintf("%.2f", v-float64(whole))[1:]
	s := fmt.Sprint(whole)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return "$" + s + cents
}

func PNumberField8(c *ui.Context) { numberField(c, 0, coss.NumberFieldSpec{Format: usd}) }

func PNumberField9(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).MaxWidth(256).Gap(24)
	a := ui.Local(host, "a", func() float64 { return 0 })
	b := ui.Local(host, "b", func() float64 { return 0 })
	host.Children(func() {
		coss.NumberField(c, coss.NumberFieldSpec{Value: a, Step: 10, Label: "Step 10"})
		coss.NumberField(c, coss.NumberFieldSpec{Value: b, Step: 0.1, Label: "Step 0.1"})
	})
}

func PNumberField10(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() float64 { return 1 })
	until := ui.Local(host, "loadingUntil", func() time.Time { return time.Time{} })
	loading := c.Now().Before(*until)
	if loading {
		c.After(150 * time.Millisecond)
	}
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() {
				coss.NumberField(c, coss.NumberFieldSpec{Value: v, Min: 1, Max: 100, Label: "Quantity"})
			})
			coss.Button(c, "Submit").Loading(loading).OnClick(func() {
				*until = c.Now().Add(800 * time.Millisecond)
			})
		})
	})
}

func PNumberField11(c *ui.Context) { numberField(c, 0, coss.NumberFieldSpec{}).Group.Radius(9999) }

// otpField shows an OTPField in a centered Field under label, above the
// description fn draws, as the labelled OTP particles do.
func otpField(c *ui.Context, label string, spec func(code string) coss.OTPFieldSpec, desc func(code string)) {
	host := ui.Box(c)
	code := ui.Local(host, "code", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			if label != "" {
				coss.FieldLabel(c, label)
			}
			coss.OTPField(c, code, spec(*code))
			if desc != nil {
				desc(*code)
			}
		}).AlignItems(ui.Center)
	})
}

// otpCodeText is a description with an inline code span.
func otpCodeText(c *ui.Context, before, code, after string) {
	t := theme.Of(c)
	ui.RichText(c).FontSize(12).LineHeight(16.0 / 12).TextColor(t.MutedForeground).Children(func() {
		ui.Text(c, before)
		ui.Text(c, code).Font("monospace").TextColor(t.Foreground)
		ui.Text(c, after)
	})
}

func POTPField1(c *ui.Context) {
	otpField(c, "", func(string) coss.OTPFieldSpec { return coss.OTPFieldSpec{} }, nil)
}

func POTPField2(c *ui.Context) {
	otpField(c, "", func(string) coss.OTPFieldSpec { return coss.OTPFieldSpec{Length: 4, Large: true} }, nil)
}

func POTPField3(c *ui.Context) {
	otpField(c, "", func(string) coss.OTPFieldSpec { return coss.OTPFieldSpec{Group: 3} }, nil)
}

func POTPField4(c *ui.Context) {
	otpField(c, "Verification code", func(string) coss.OTPFieldSpec { return coss.OTPFieldSpec{Length: 4} },
		func(string) { coss.FieldDescription(c, "Enter the 4-digit code sent to your email.") })
}

func POTPField6(c *ui.Context) {
	tier := func(r rune) bool { return r >= '0' && r <= '3' }
	otpField(c, "Tier code", func(string) coss.OTPFieldSpec { return coss.OTPFieldSpec{Accept: tier} },
		func(string) { coss.FieldDescription(c, "Digits 0-3 only.") })
}

func POTPField7(c *ui.Context) {
	otpField(c, "Verification code",
		func(code string) coss.OTPFieldSpec {
			return coss.OTPFieldSpec{Invalid: len(code) == 6 && code != "123456"}
		},
		func(code string) {
			switch {
			case code == "123456":
				coss.FieldDescription(c, "Code verified.")
			case len(code) == 6:
				coss.FieldError(c, "Code must be 123456.")
			default:
				coss.FieldDescription(c, "Enter `123456` to pass validation.")
			}
		})
}

func POTPField8(c *ui.Context) {
	alnum := func(r rune) bool { return unicode.IsDigit(r) || unicode.IsLetter(r) }
	otpField(c, "Recovery code", func(string) coss.OTPFieldSpec { return coss.OTPFieldSpec{Accept: alnum} },
		func(string) { otpCodeText(c, "Accept letters and numbers for backup codes such as ", "A7C9XZ", ".") })
}

func POTPField9(c *ui.Context) {
	otpField(c, "Verification code", func(string) coss.OTPFieldSpec { return coss.OTPFieldSpec{Placeholder: "•"} },
		func(string) {
			coss.FieldDescription(c, "Placeholder hints stay visible until the focused slot is active.")
		})
}

func POTPField10(c *ui.Context) {
	otpField(c, "Access code", func(string) coss.OTPFieldSpec { return coss.OTPFieldSpec{Mask: true} },
		func(string) { otpCodeText(c, "Use ", "mask", " to obscure the code on shared screens.") })
}
