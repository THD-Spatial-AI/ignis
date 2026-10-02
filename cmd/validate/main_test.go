package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestReport_countsPerTableAndListsFailures(t *testing.T) {
	results := map[string][]TestResult{
		"germany": {
			{BuildingID: "DE.A", CalculatedQHND: 100, ExpectedQHND: 100, Passed: true},
		},
		"spain": {
			{BuildingID: "ES.A", CalculatedQHND: 50, ExpectedQHND: 50, Passed: true},
			{BuildingID: "ES.B", CalculatedQHND: 97, ExpectedQHND: 100, PercentError: 3, Passed: false},
			{BuildingID: "ES.C", ErrorMessage: "pipeline error: boom"},
		},
	}

	var out bytes.Buffer
	failed := report(&out, []string{"germany", "spain"}, results)

	if failed != 2 {
		t.Errorf("failed = %d, want 2", failed)
	}
	got := out.String()
	for _, want := range []string{"germany", "1/1", "spain", "1/3", "ES.B", "3.00%", "ES.C", "pipeline error: boom", "2/4"} {
		if !strings.Contains(got, want) {
			t.Errorf("report output missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "DE.A") {
		t.Errorf("report lists a passing row:\n%s", got)
	}
}

func TestExitCode_nonZeroOnlyWhenStrictAndFailing(t *testing.T) {
	cases := []struct {
		failed int
		strict bool
		want   int
	}{
		{0, false, 0},
		{0, true, 0},
		{3, false, 0},
		{3, true, 1},
	}
	for _, c := range cases {
		if got := exitCode(c.failed, c.strict); got != c.want {
			t.Errorf("exitCode(%d, %v) = %d, want %d", c.failed, c.strict, got, c.want)
		}
	}
}
