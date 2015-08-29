package passgen

import (
	"bytes"
	crand "crypto/rand"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"unicode"
)

type randSource struct {
	rand.Source
}

type dict []string

func (src *randSource) Read(p []byte) (int, error) {
	for i := range p {
		// Extract least significant byte from 63-bit integer
		p[i] = byte(src.Int63() & 0xff)
	}
	return len(p), nil
}

func (d dict) toBuffer() *bytes.Buffer {
	return bytes.NewBufferString(strings.Join(d, "\n"))
}

var words = dict{
	"then",
	"there",
	"these",
	"theta",
	"they",
	"thick",
	"thief",
	"thigh",
	"thin",
	"thine",
	"thing",
	"think",
	"third",
	"this",
}

func Benchmark64Base32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Base32(64)
	}
}

func Benchmark64Hex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Hex(64)
	}
}

func Benchmark64AsciiComplete(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Ascii(64, SetComplete)
	}
}

func Benchmark64AsciiLowerDigit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Ascii(64, SetLower|SetDigit)
	}
}

func Benchmark10Diceware(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Diceware(words.toBuffer(), 10, "")
	}
}

func ExampleAscii() {
	password, err := Ascii(20, SetComplete&^SetSymbol)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated password: %q\n", password)
}

func ExampleDiceware() {
	dict, err := os.Open("/usr/share/dict/words")
	if err != nil {
		log.Fatal(err)
	}

	password, _, err := Diceware(dict, 5, "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated password: %q\n", password)
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
	Rng = &randSource{rand.NewSource(0)}

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

	Rng = crand.Reader
}

func TestBase32(t *testing.T) {
	Rng = &randSource{rand.NewSource(0)}

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

	Rng = crand.Reader
}

func TestAscii(t *testing.T) {
	Rng = &randSource{rand.NewSource(0)}

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

	Rng = crand.Reader
}

func TestInvalidRng(t *testing.T) {
	var err error
	Rng = os.Stdin

	if _, err = Ascii(1, SetComplete); err == nil {
		t.Error("Ascii: err == nil with invalid Reader")
	}
	if _, err = Hex(1); err == nil {
		t.Error("Hex: err == nil with invalid Reader")
	}
	if _, err = Base32(1); err == nil {
		t.Error("Base32: err == nil with invalid Reader")
	}
	if _, _, err = Diceware(words.toBuffer(), 1, ""); err == nil {
		t.Error("Diceware: err == nil with invalid Reader")
	}

	Rng = crand.Reader
}

func TestDicewareInvalidDict(t *testing.T) {
	_, _, err := Diceware(os.Stdin, 1, "")
	if err != ErrDict {
		t.Fatal("Expected ErrDict for invalid dict")
	}
}

func TestDiceware(t *testing.T) {
	n := 6

	phrase, m, err := Diceware(words.toBuffer(), n, "")
	if err != nil {
		t.Fatal(err)
	}

	if m != len(words) {
		t.Errorf("dict contains %d words, but m = %d",
			len(words), m)
	}

	if len(phrase) < n {
		t.Errorf("len(phrase) = %d (n = %d)", len(phrase), n)
	}

	if strings.ContainsAny(string(phrase), " ") {
		t.Errorf(`phrase "%s" contains spaces (sep = "")`, phrase)
	}
}
