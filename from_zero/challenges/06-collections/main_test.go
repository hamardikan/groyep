// Run: go test -v
//
// This test builds your word frequency program and checks its output.
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
	binaryPath = filepath.Join(dir, "wordfreq_test_bin")

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

func TestWordFreq(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "three distinct words with repeats",
			args:     []string{"hello", "world", "hello", "go", "world", "hello"},
			expected: "go: 1\nhello: 3\nworld: 2",
		},
		{
			name:     "single word",
			args:     []string{"apple"},
			expected: "apple: 1",
		},
		{
			name:     "two words alternating",
			args:     []string{"a", "b", "a", "b", "a"},
			expected: "a: 3\nb: 2",
		},
		{
			name:     "all same word",
			args:     []string{"go", "go", "go"},
			expected: "go: 3",
		},
		{
			name:     "alphabetical sorting",
			args:     []string{"cherry", "apple", "banana"},
			expected: "apple: 1\nbanana: 1\ncherry: 1",
		},
		{
			name:     "case sensitive",
			args:     []string{"Hello", "hello", "HELLO"},
			expected: "HELLO: 1\nHello: 1\nhello: 1",
		},
		{
			name:     "many duplicates",
			args:     []string{"x", "y", "x", "y", "x", "y", "z"},
			expected: "x: 3\ny: 3\nz: 1",
		},
		{
			name:     "single character words",
			args:     []string{"c", "b", "a", "b", "c", "c"},
			expected: "a: 1\nb: 2\nc: 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinary(tt.args...)
			if err != nil {
				t.Fatalf("program exited with error: %v\nOutput: %s", err, output)
			}
			if output != tt.expected {
				t.Errorf("wordfreq %v\n  expected:\n%s\n  actual:\n%s", tt.args, tt.expected, output)
			}
		})
	}
}

func TestWordFreqSorted(t *testing.T) {
	// Specifically test that output is sorted alphabetically
	output, err := runBinary("zebra", "apple", "mango", "banana")
	if err != nil {
		t.Fatalf("program exited with error: %v", err)
	}

	lines := strings.Split(output, "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d: %q", len(lines), output)
	}

	expected := []string{"apple: 1", "banana: 1", "mango: 1", "zebra: 1"}
	for i, exp := range expected {
		if lines[i] != exp {
			t.Errorf("line %d: expected %q, got %q", i+1, exp, lines[i])
		}
	}
}

func TestWordFreqNoArgs(t *testing.T) {
	output, _ := runBinary()
	expected := "usage: wordfreq <word1> [word2] ..."
	if output != expected {
		t.Errorf("no args\n  expected: %q\n  actual:   %q", expected, output)
	}
}

func TestWordFreqOutputFormat(t *testing.T) {
	output, err := runBinary("test")
	if err != nil {
		t.Fatalf("program exited with error: %v", err)
	}

	// Check format is "word: count" (with colon and space)
	if !strings.Contains(output, ": ") {
		t.Errorf("output format should be 'word: count', got: %q", output)
	}

	if output != "test: 1" {
		t.Errorf("expected %q, got %q", "test: 1", output)
	}
}
