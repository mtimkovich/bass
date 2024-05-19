package main

import (
	"strings"
)

var note_to_int = map[string]int{
	"A":  0,
	"A#": 1,
	"Bb": 1,
	"B":  2,
	"C":  3,
	"C#": 4,
	"Db": 4,
	"D":  5,
	"D#": 6,
	"Eb": 6,
	"E":  7,
	"F":  8,
	"F#": 9,
	"Gb": 9,
	"G":  10,
	"G#": 11,
	"Ab": 11,
}

var int_to_note = map[int]string{
	0:  "A",
	1:  "Bb",
	2:  "B",
	3:  "C",
	4:  "Db",
	5:  "D",
	6:  "Eb",
	7:  "E",
	8:  "F",
	9:  "Gb",
	10: "G",
	11: "Ab",
}

const OCTAVE = 12

// Add half steps to a note and return the resulting note.
func half_step_plus(start string, n int) string {
	start = strings.ToUpper(start)
	tone := note_to_int[start] + n
	return int_to_note[tone%OCTAVE]
}

// func main() {
// 	strings := []string{"E", "A", "D", "G"}
// 	// dots := []int{0, 3, 5, 7, 9, 12, 15, 17, 19, 21}
// 	dots := []int{0, 3, 5, 7, 9}

// 	for _, dot := range dots {
// 		fmt.Printf("%v: ", dot)
// 		for _, s := range strings {
// 			fmt.Printf("% 2v ", half_step_plus(s, dot))
// 		}
// 		fmt.Println()
// 	}
// }
