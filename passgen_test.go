package passgen

import (
	"testing"
	"unicode"
)

func Benchmark128Complete(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Generate(128, Complete)
	}
}

func Benchmark128WithoutSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Generate(128, (Complete &^ Symbol))
	}
}

func Benchmark128Lower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Generate(128, Lower)
	}
}

func Benchmark32Complete(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Generate(32, Complete)
	}
}

func Benchmark32WithoutSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Generate(32, (Complete &^ Symbol))
	}
}

func TestRange(t *testing.T) {
	pass, err := Generate(128, Complete)
	if err != nil {
		t.Fatal(err)
	}

	var min, max byte = charOffset, charOffset + charRange

	for _, b := range pass {
		if b < min || b > max {
			t.Errorf("Encountered illegal character '%c' with value %d",
				b, b)
		}
	}
}

func TestWithoutSymbols(t *testing.T) {
	flags := Complete &^ Symbol
	pass, err := Generate(128, flags)
	if err != nil {
		t.Fatal(err)
	}

	for _, b := range pass {
		if unicode.IsSymbol(rune(b)) {
			t.Errorf("Encountered symbol '%c' with flag %b", b, flags)
		}
	}
}
