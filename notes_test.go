package main

import (
	"testing"
)

var noteTests = []struct {
	in  string
	out string
}{
	{"E0", "E"},
	{"e0", "E"},
	{"Ab0", "Ab"},
	{"ab0", "Ab"},
	{"a12", "A"},
	{"E3", "G"},
	{"e4", "Ab"},
}

func TestParse(t *testing.T) {
	for _, tt := range noteTests {
		t.Run(tt.in, func(t *testing.T) {
			note, fret, err := parse(tt.in)
			if err != nil {
				t.Errorf("%v: parse error", tt.in)
			}
			output := half_step_plus(note, fret)
			if output != tt.out {
				t.Errorf("Want %s got %s", tt.out, output)
			}
		})
	}
}
