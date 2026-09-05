package main

import (
	"testing"
)

func TestParseVersionFlag(t *testing.T) {
	cases := map[string][]string{
		"no flags":     {},
		"-version":     {"-version"},
		"-v shorthand": {"-v"},
	}
	for name, args := range cases {
		want := len(args) > 0
		if got := parseVersionFlag(args); got != want {
			t.Errorf("%s: parseVersionFlag(%v) = %v, want %v", name, args, got, want)
		}
	}
}
