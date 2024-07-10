package main

import (
	"math"
	"strings"

	"github.com/mattn/go-runewidth"
)


func char_wrap(s string, limit int) string {
	var result strings.Builder

	var rune_sl []rune
	var diff float64

	var aj_limit float64

	var fl_len float64

	var int_aj_limit int



	// if len(rune_sl) < limit {
	// 	return s
	// }

	string_split := strings.Split(s, "\n")


	// var inc int64 = 0

	for _, str := range string_split {
		rune_sl = []rune(str)

		for {
				// inc += 1
				// fmt.Println(inc)

				fl_len = float64(len(rune_sl))

				diff = float64(runewidth.StringWidth(string(rune_sl))) / fl_len

				if diff > 0 {
					aj_limit = float64(limit) / diff
					int_aj_limit = int(math.Floor(aj_limit))
				}


				if len(rune_sl) <= int_aj_limit {
					result.WriteString(string(rune_sl))
					result.WriteString("\n")
					break
				}

				// diff = runewidth.StringWidth(str) - len(rune_sl)
				// fmt.Printf("len: %v\n", len(rune_sl))
				// fmt.Printf("fl_len: %v\n", fl_len)
				// fmt.Printf("aj_limit: %v\n", aj_limit)
				// fmt.Printf("int_aj_limit: %v\n", int_aj_limit)
				// fmt.Printf("diff: %v\n", diff)



				result.WriteString(string(rune_sl[:int_aj_limit-1]))
				result.WriteString("⏎\n")

				rune_sl = rune_sl[int_aj_limit-1:]
		}



	}

	return result.String()
}