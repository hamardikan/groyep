// Run: go test -v
//
// This test builds your FizzBuzz program and checks its output for various inputs.
// You do NOT need to edit this file — just create your main.go and run the tests.

package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	dir, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory: " + err.Error())
	}
	binaryPath = filepath.Join(dir, "fizzbuzz_test_bin")

	build := exec.Command("go", "build", "-o", binaryPath, ".")
	build.Dir = dir
	if output, err := build.CombinedOutput(); err != nil {
		panic("failed to build binary: " + err.Error() + "\n" + string(output))
	}

	code := m.Run()

	os.Remove(binaryPath)
	os.Exit(code)
}

func runBinary(args ...string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestFizzBuzz15(t *testing.T) {
	output, err := runBinary("15")
	if err != nil {
		t.Fatalf("program exited with error: %v\nOutput: %s", err, output)
	}

	expected := "1\n2\nFizz\n4\nBuzz\nFizz\n7\n8\nFizz\nBuzz\n11\nFizz\n13\n14\nFizzBuzz\n"
	if output != expected {
		t.Errorf("fizzbuzz 15\n  expected: %q\n  actual:   %q", expected, output)
	}
}

func TestFizzBuzzCases(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "N=1",
			args:     []string{"1"},
			expected: "1\n",
		},
		{
			name:     "N=3",
			args:     []string{"3"},
			expected: "1\n2\nFizz\n",
		},
		{
			name:     "N=5",
			args:     []string{"5"},
			expected: "1\n2\nFizz\n4\nBuzz\n",
		},
		{
			name:     "N=0 produces no output",
			args:     []string{"0"},
			expected: "",
		},
		{
			name:     "negative N produces no output",
			args:     []string{"-5"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinary(tt.args...)
			if err != nil && tt.expected != "" {
				t.Fatalf("program exited with error: %v\nOutput: %s", err, output)
			}
			if output != tt.expected {
				t.Errorf("fizzbuzz %v\n  expected: %q\n  actual:   %q", tt.args, tt.expected, output)
			}
		})
	}
}

func TestFizzBuzzSpecificValues(t *testing.T) {
	output, err := runBinary("30")
	if err != nil {
		t.Fatalf("program exited with error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	// Table of specific line checks (0-indexed)
	checks := []struct {
		lineNum  int // 1-indexed line number
		expected string
	}{
		{3, "Fizz"},      // 3 is divisible by 3
		{5, "Buzz"},      // 5 is divisible by 5
		{6, "Fizz"},      // 6 is divisible by 3
		{9, "Fizz"},      // 9 is divisible by 3
		{10, "Buzz"},     // 10 is divisible by 5
		{15, "FizzBuzz"}, // 15 is divisible by both
		{30, "FizzBuzz"}, // 30 is divisible by both
		{7, "7"},         // 7 is neither
		{22, "22"},       // 22 is neither
	}

	for _, c := range checks {
		idx := c.lineNum - 1
		if idx >= len(lines) {
			t.Errorf("line %d: expected %q but only got %d lines", c.lineNum, c.expected, len(lines))
			continue
		}
		if lines[idx] != c.expected {
			t.Errorf("line %d (number %d): expected %q, got %q", c.lineNum, c.lineNum, c.expected, lines[idx])
		}
	}
}

func TestFizzBuzzErrorCases(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "no arguments",
			args:     []string{},
			expected: "usage: fizzbuzz <n>",
		},
		{
			name:     "not a number",
			args:     []string{"abc"},
			expected: "error: not a number",
		},
		{
			name:     "float is not valid",
			args:     []string{"3.5"},
			expected: "error: not a number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, _ := runBinary(tt.args...)
			trimmed := strings.TrimRight(output, "\n")
			if trimmed != tt.expected {
				t.Errorf("fizzbuzz %v\n  expected: %q\n  actual:   %q", tt.args, tt.expected, trimmed)
			}
		})
	}
}
