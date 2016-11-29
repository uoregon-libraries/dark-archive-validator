package main

import (
	"fmt"
	"unicode"
)

var equivmap = make(map[rune][]rune)

func main() {
	var r rune
	var equiv = make([]rune, 1, 8)
	var lowest rune
	for r = 0; r < unicode.MaxRune; r++ {
		var folded = unicode.SimpleFold(r)
		if folded == r {
			continue
		}

		if r < folded {
			lowest = r
		} else {
			lowest = folded
		}

		equiv[0] = r
		for folded != r {
			equiv = append(equiv, folded)
			folded = unicode.SimpleFold(folded)

			if folded < lowest {
				lowest = folded
			}
		}

		if _, exists := equivmap[lowest]; !exists {
			equivmap[lowest] = equiv
		}
		equiv = make([]rune, 1, 8)
	}

	for r, rList := range equivmap {
		fmt.Printf("%c: %#v (%#v)\n", r, string(rList), rList)
	}
}
