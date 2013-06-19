package passgen

import (
	"fmt"
	"testing"
)

func TestDistribution(t *testing.T) {
	mask := SetSymbol | SetLower
	_ = mask
	/* pass, err := Generate(12800, mask) */
	pass, err := GenerateBase32(12800)
	if err != nil {
		t.Fatal(err)
	}

	chars := make(map[byte]int)

	/* for _, s := range charSets { */
	/* 	for _, c := range s { */
	/* 		chars[byte(c)] = 0 */
	/* 	} */
	/* } */

	for _, c := range pass {
		chars[c]++
	}

	min, max := 100000, 0

	for k, v := range chars {
		fmt.Printf("%c: %d\n", k, v)

		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	fmt.Printf("min: %d\n", min)
	fmt.Printf("max: %d\n", max)
}
