package passgen

import (
	"bytes"
	crand "crypto/rand"
	"math/rand"
	"strings"
	"testing"
	"unicode"
)

type randSource struct {
	rand.Source
}

func (src *randSource) Read(p []byte) (int, error) {
	for i := range p {
		// Extract least significant byte from 63-bit integer
		p[i] = byte(src.Int63() & 0xff)
	}
	return len(p), nil
}

func Benchmark256Base32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Base32(256)
	}
}

func Benchmark256Hex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Hex(256)
	}
}

func Benchmark256Complete(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Ascii(256, SetComplete)
	}
}

func Benchmark256LowerDigit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Ascii(256, SetLower|SetDigit)
	}
}

func TestAsciiCharRange(t *testing.T) {
	pass, err := Ascii(1000, SetComplete)
	if err != nil {
		t.Fatal(err)
	}

	var min, max byte = 32, 126

	for _, b := range pass {
		if b < min || b > max {
			t.Errorf("Encountered illegal character '%c' with value %d",
				b, b)
		}
	}
}

func TestAsciiWithoutSymbols(t *testing.T) {
	mask := SetComplete &^ SetSymbol
	pass, err := Ascii(1000, mask)
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range pass {
		if unicode.IsSymbol(rune(b)) {
			t.Errorf("Encountered symbol '%c' with mask %b",
				b, mask)
		}
	}
}

func TestAsciiLen(t *testing.T) {
	n := 32
	pass, err := Ascii(n, SetComplete)
	if err != nil {
		t.Fatal(err)
	}

	if len(pass) != n {
		t.Errorf("len(pass) = %d, expected %d",
			len(pass), n)
	}

	pass, err = Ascii(0, SetComplete)
	if err != ErrLength {
		t.Fatalf("Expected ErrLength error")
	}
}

func TestHex(t *testing.T) {
	Reader = &randSource{rand.NewSource(0)}

	expected := []byte("01c073624aaf3978514ef8443bb2a859c75fc3cc6af26d5aaa20926f046baa66")
	pass, err := Hex(len(expected))
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(pass, expected) {
		t.Errorf(`Passwords don't match:
	Expected:  %s
	Generated: %s`, expected, pass)
	}

	Reader = crand.Reader
}

func TestBase32(t *testing.T) {
	Reader = &randSource{rand.NewSource(0)}

	expected := []byte("AHAHGYSKV44XQUKO7BCDXMVILHDV7Q6MNLZG2WVKECJG6BDLVJTOZENFUJ4UGI6C")
	pass, err := Base32(len(expected))
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(pass, expected) {
		t.Errorf(`Passwords don't match:
	Expected:  %s
	Generated: %s`, expected, pass)
	}

	Reader = crand.Reader
}

func TestAscii(t *testing.T) {
	Reader = &randSource{rand.NewSource(0)}

	expected := []byte("a5XL<C8ONID>}SV12T$q.3c1$Z-_]8HrGTpU.iDRw'?0`^0B]P9y>7TMA[FO\"jbe")

	pass, err := Ascii(len(expected), SetComplete)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(pass, expected) {
		t.Errorf(`Passwords don't match:
	Expected:  %s
	Generated: %s`, expected, pass)
	}

	Reader = crand.Reader
}

func TestDiceware(t *testing.T) {
	DicewareDict = ""
	n := 6
	if _, _, err := Diceware(n, ""); err == nil {
		t.Error("no error for empty DicewareDict")
	}

	DicewareDict = "diceware-en.txt"
	phrase, m, err := Diceware(n, "")
	if err != nil {
		t.Fatal(err)
	}

	// File has 7776 words, 26 words that end with "'s" are ignored
	if m != 7750 {
		t.Errorf("%s contains %d words, but m = %d",
			DicewareDict, 7750, m)
	}

	if len(phrase) < n {
		t.Errorf("len(phrase) = %d (n = %d)", len(phrase), n)
	}

	if strings.ContainsAny(string(phrase), " ") {
		t.Errorf("phrase '%s' contains spaces (sep = \"\")", phrase)
	}
}
