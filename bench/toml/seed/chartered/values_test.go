package toml

import (
	"testing"

	"example.com/toml/claims"
	"github.com/leroux/charter/chartertest"
)

func TestShortEscapesDecode(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_QZWX, "escapes-decoded")
	decodes(t, `s = "\b\t\n\f\r\"\\"`+"\n", map[string]any{"s": "\b\t\n\f\r\"\\"})
}

func TestUnicodeEscapesDecode(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_QZWX, "unicode-escapes")
	decodes(t, `s = "caf\u00E9 \U0001F600"`+"\n", map[string]any{"s": "café 😀"})
}

func TestAnUnknownEscapeIsRefused(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_QZWX, "unknown-escape-refused")
	for _, s := range []string{`\x41`, `\e`, `\q`, `\ `} {
		refused(t, `s = "`+s+`"`+"\n", "invalid escape")
	}
}

func TestASurrogateEscapeIsRefused(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_QZWX, "surrogate-refused")
	refused(t, `s = "\uD800"`+"\n", "scalar value")
	refused(t, `s = "\U00110000"`+"\n", "scalar value")
}

func TestABasicStringIsOneLine(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_QZWX, "one-line")
	refused(t, "s = \"no end\nt = 1\n", "unterminated")
	refused(t, "s = \"no end", "unterminated")
}

func TestAControlCharacterInABasicStringIsRefused(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_QZWX, "control-refused")
	refused(t, "s = \"a\x00b\"\n", "control character")
	decodes(t, "s = \"a\tb\"\n", map[string]any{"s": "a\tb"})
}

func TestALiteralStringKeepsItsBackslashes(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_WZ4K, "backslash-kept")
	decodes(t, `p = 'C:\Users\n'`+"\n", map[string]any{"p": `C:\Users\n`})
}

func TestAControlCharacterInALiteralStringIsRefused(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_WZ4K, "control-refused")
	refused(t, "p = 'a\x00b'\n", "control character")
}

func TestSignedDecimalIntegers(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_JD2N, "signed-decimal")
	decodes(t, "a = 42\nb = +17\nc = -5\nd = -0\ne = 0\n",
		map[string]any{"a": int64(42), "b": int64(17), "c": int64(-5), "d": int64(0), "e": int64(0)})
}

func TestALeadingZeroIsRefused(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_JD2N, "leading-zero-refused")
	refused(t, "a = 007\n")
	refused(t, "a = +01\n")
}

func TestAnIntegerBeyond64BitsIsRefused(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_JD2N, "out-of-range-refused")
	refused(t, "a = 9223372036854775808\n", "64 bits")
	decodes(t, "a = -9223372036854775808\n", map[string]any{"a": int64(-9223372036854775808)})
}

func TestBooleansAreLowerCase(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_WN2V, "lowercase-only")
	decodes(t, "a = true\nb = false\n", map[string]any{"a": true, "b": false})
	for _, v := range []string{"True", "FALSE", "truthy", "t"} {
		refused(t, "a = "+v+"\n")
	}
}

func TestFloatsAreRefusedByName(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_H72C, "float-refused")
	for _, v := range []string{"3.14", "1e6", "-2.5E-3", "inf", "nan", "+inf"} {
		refused(t, "a = "+v+"\n", "floats")
	}
}

func TestArraysAreRefusedByName(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_H72C, "array-refused")
	refused(t, "a = [1, 2]\n", "arrays")
	refused(t, "[[t]]\n", "arrays of tables")
}

func TestDatesAreRefusedByName(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_H72C, "date-refused")
	for _, v := range []string{"1979-05-27", "1979-05-27T07:32:00Z", "07:32:00"} {
		refused(t, "a = "+v+"\n", "dates and times")
	}
}

func TestOtherUnbuiltFormsAreRefusedByName(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_H72C, "other-forms-refused")
	refused(t, "a = 0xff\n", "hexadecimal")
	refused(t, "a = 1_000\n", "underscores")
	refused(t, "a = \"\"\"x\"\"\"\n", "multi-line strings")
	refused(t, "a = '''x'''\n", "multi-line strings")
	refused(t, "a = { b = 1 }\n", "inline tables")
}
