// Run: go test -v
//
// This test builds your program and checks its behavior with different inputs.
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
	binaryPath = filepath.Join(dir, "classify_test_bin")

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
	return strings.TrimRight(string(output), "\n"), err
}

func TestClassify(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "positive number",
			args:     []string{"5"},
			expected: "positive",
		},
		{
			name:     "large positive number",
			args:     []string{"100"},
			expected: "positive",
		},
		{
			name:     "one",
			args:     []string{"1"},
			expected: "positive",
		},
		{
			name:     "negative number",
			args:     []string{"-3"},
			expected: "negative",
		},
		{
			name:     "large negative number",
			args:     []string{"-999"},
			expected: "negative",
		},
		{
			name:     "negative one",
			args:     []string{"-1"},
			expected: "negative",
		},
		{
			name:     "zero",
			args:     []string{"0"},
			expected: "zero",
		},
		{
			name:     "not a number - letters",
			args:     []string{"abc"},
			expected: "error: not a number",
		},
		{
			name:     "not a number - mixed",
			args:     []string{"12abc"},
			expected: "error: not a number",
		},
		{
			name:     "not a number - float",
			args:     []string{"3.14"},
			expected: "error: not a number",
		},
		{
			name:     "not a number - empty string",
			args:     []string{""},
			expected: "error: not a number",
		},
		{
			name:     "no arguments",
			args:     []string{},
			expected: "usage: classify <number>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, _ := runBinary(tt.args...)
			if output != tt.expected {
				t.Errorf("classify %v\n  expected: %q\n  actual:   %q", tt.args, tt.expected, output)
			}
		})
	}
}
