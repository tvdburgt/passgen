package passgen

import (
	"crypto/rand"
	"math/big"
	"unicode"
)

const (
	charOffset = 32 // Printable ASCII range: [32, 126]
	charRange  = 95
)

const (
	// Bitmasks for controlling password output
	Lower    = 1 << iota       // Lower case letters [a-z]
	Upper                      // Upper case letters [A-Z]
	Digit                      // Decimal digits [0-9]
	Punct                      // Punctuation characters
	Space                      // White space characters
	Symbol                     // Symbolic characters (incomplete set: [<>|^`~])
	Complete = (1 << iota) - 1 // Set of all printable characters
)

// Generates a cryptographically secure password, by reading from crypto's
// rand.Reader. Charset constraints can be applied using the flags parameter.
func Generate(size, flags int) ([]byte, error) {

	pass := make([]byte, size)

	for i := 0; i < size; {
		n, err := rand.Int(rand.Reader, big.NewInt(charRange))
		if err != nil {
			return nil, err
		}

		b := byte(charOffset + n.Int64())

		switch {
		case flags&Digit == 0 && unicode.IsDigit(rune(b)):
			continue

		case flags&Lower == 0 && unicode.IsLower(rune(b)):
			continue

		case flags&Upper == 0 && unicode.IsUpper(rune(b)):
			continue

		case flags&Punct == 0 && unicode.IsPunct(rune(b)):
			continue

		case flags&Symbol == 0 && unicode.IsSymbol(rune(b)):
			continue

		case flags&Space == 0 && unicode.IsSpace(rune(b)):
			continue
		}

		pass[i] = b
		i++
	}
	return pass, nil
}
