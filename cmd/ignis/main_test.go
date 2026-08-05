package main

import (
	"strings"
	"testing"

	"github.com/thd-spatial-ai/ignis/internal/version"
)

func TestVersionString(t *testing.T) {
	got := versionString()
	for _, want := range []string{version.Version, version.Commit, version.Date} {
		if !strings.Contains(got, want) {
			t.Errorf("versionString() = %q, want it to contain %q", got, want)
		}
	}
}

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
