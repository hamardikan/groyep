// Run: go test -v
//
// This test builds your calculator program and checks its output for various inputs.
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
	binaryPath = filepath.Join(dir, "calc_test_bin")

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

func TestCalcOperations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "add integers",
			args:     []string{"10", "add", "5"},
			expected: "Result: 15.00",
		},
		{
			name:     "add floats",
			args:     []string{"1.5", "add", "2.3"},
			expected: "Result: 3.80",
		},
		{
			name:     "add negative",
			args:     []string{"-5", "add", "3"},
			expected: "Result: -2.00",
		},
		{
			name:     "subtract",
			args:     []string{"10", "sub", "3"},
			expected: "Result: 7.00",
		},
		{
			name:     "subtract resulting in negative",
			args:     []string{"3", "sub", "10"},
			expected: "Result: -7.00",
		},
		{
			name:     "multiply integers",
			args:     []string{"4", "mul", "3"},
			expected: "Result: 12.00",
		},
		{
			name:     "multiply with float",
			args:     []string{"4", "mul", "2.5"},
			expected: "Result: 10.00",
		},
		{
			name:     "multiply by zero",
			args:     []string{"5", "mul", "0"},
			expected: "Result: 0.00",
		},
		{
			name:     "divide",
			args:     []string{"7", "div", "2"},
			expected: "Result: 3.50",
		},
		{
			name:     "divide even",
			args:     []string{"10", "div", "5"},
			expected: "Result: 2.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinary(tt.args...)
			if err != nil {
				t.Fatalf("program exited with error: %v\nOutput: %s", err, output)
			}
			if output != tt.expected {
				t.Errorf("calc %v\n  expected: %q\n  actual:   %q", tt.args, tt.expected, output)
			}
		})
	}
}

func TestCalcErrors(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "division by zero",
			args:     []string{"10", "div", "0"},
			expected: "error: division by zero",
		},
		{
			name:     "unknown operator",
			args:     []string{"10", "pow", "2"},
			expected: "error: unknown operator 'pow'",
		},
		{
			name:     "unknown operator mod",
			args:     []string{"10", "mod", "3"},
			expected: "error: unknown operator 'mod'",
		},
		{
			name:     "invalid first number",
			args:     []string{"abc", "add", "5"},
			expected: "error: invalid number 'abc'",
		},
		{
			name:     "invalid second number",
			args:     []string{"5", "add", "xyz"},
			expected: "error: invalid number 'xyz'",
		},
		{
			name:     "no arguments",
			args:     []string{},
			expected: "usage: calc <num1> <operator> <num2>",
		},
		{
			name:     "too few arguments",
			args:     []string{"1", "2"},
			expected: "usage: calc <num1> <operator> <num2>",
		},
		{
			name:     "too many arguments",
			args:     []string{"1", "add", "2", "3"},
			expected: "usage: calc <num1> <operator> <num2>",
		},
		{
			name:     "one argument",
			args:     []string{"5"},
			expected: "usage: calc <num1> <operator> <num2>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, _ := runBinary(tt.args...)
			if output != tt.expected {
				t.Errorf("calc %v\n  expected: %q\n  actual:   %q", tt.args, tt.expected, output)
			}
		})
	}
}
