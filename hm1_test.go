package main

import (
	"strings"
	"testing"
)

func TestInteractiveApplicationWithFakeInput(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		"1", "E001", "100", "Ava", "123", "IT", "not-a-number", "2500",
		"1", "E001", "E002", "Ben", "HR", "1500",
		"4", "E001",
		"4", "E001",
		"5",
		"0",
	}, "\n"))
	var output strings.Builder

	runInteractiveApplication(input, &output)

	for _, expected := range []string{
		"Invalid input: employee name must contain at least one letter",
		"Invalid input: department must contain at least one letter",
		"Invalid input: salary must be a valid number greater than zero.",
		"Invalid input: employee already exists",
		"SUCCESS: E001 added.",
		"SUCCESS: E002 added.",
		"SUCCESS: payment of 2500.00 processed for E001.",
		"ERROR: payment already processed",
		"Total employees: 2",
		"Payments processed: 1",
	} {
		if !strings.Contains(output.String(), expected) {
			t.Errorf("interactive output does not contain %q", expected)
		}
	}
}
