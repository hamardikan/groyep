// Run: go test -v
//
// This test builds your search program and checks its output against expected behavior.
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
	binaryPath = filepath.Join(dir, "search_test_bin")

	build := exec.Command("go", "build", "-o", binaryPath, ".")
	build.Dir = dir
	if output, err := build.CombinedOutput(); err != nil {
		panic("failed to build binary: " + err.Error() + "\n" + string(output))
	}

	code := m.Run()

	os.Remove(binaryPath)
	os.Exit(code)
}

// runBinary runs the search binary and returns stdout, stderr, and the exit code.
func runBinary(args ...string) (stdout string, stderr string, exitCode int) {
	cmd := exec.Command(binaryPath, args...)
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	exitCode = 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		exitCode = -1
	}

	stdout = strings.TrimRight(outBuf.String(), "\n")
	stderr = strings.TrimRight(errBuf.String(), "\n")
	return
}

func TestSearchMatches(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		file     string
		expected string
	}{
		{
			name:     "single match — Priest",
			pattern:  "Priest",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:Judas Priest",
		},
		{
			name:     "multiple matches — Bad",
			pattern:  "Bad",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:Bad English\ntestdata/rockbands.txt:Bad Company",
		},
		{
			name:     "multiple matches — White",
			pattern:  "White",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:Whitesnake\ntestdata/rockbands.txt:Great White\ntestdata/rockbands.txt:White Lion\ntestdata/rockbands.txt:Whitecross",
		},
		{
			name:     "match at start of file",
			pattern:  "Judas",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:Judas Priest",
		},
		{
			name:     "match at end of file",
			pattern:  "Nirvana",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:Nirvana",
		},
		{
			name:     "match in fruits file",
			pattern:  "ap",
			file:     "testdata/fruits.txt",
			expected: "testdata/fruits.txt:apple\ntestdata/fruits.txt:apricot",
		},
		{
			name:     "multiple matches in fruits — a",
			pattern:  "a",
			file:     "testdata/fruits.txt",
			expected: "testdata/fruits.txt:apple\ntestdata/fruits.txt:banana\ntestdata/fruits.txt:apricot\ntestdata/fruits.txt:avocado",
		},
		{
			name:     "duplicate line in file — Tora Tora",
			pattern:  "Tora Tora",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:Tora Tora\ntestdata/rockbands.txt:Tora Tora",
		},
		{
			name:     "partial match within word — ois",
			pattern:  "ois",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:Poison",
		},
		{
			name:     "special chars in pattern — /",
			pattern:  "/",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:AC/DC\ntestdata/rockbands.txt:Love/Hate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, exitCode := runBinary(tt.pattern, tt.file)
			if exitCode != 0 {
				t.Fatalf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
			}
			if stdout != tt.expected {
				t.Errorf("search %q %s\n  expected:\n%s\n  actual:\n%s", tt.pattern, tt.file, tt.expected, stdout)
			}
		})
	}
}

func TestSearchNoMatch(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		file    string
	}{
		{
			name:    "no match — zzz",
			pattern: "zzz",
			file:    "testdata/rockbands.txt",
		},
		{
			name:    "case sensitive — priest (lowercase)",
			pattern: "priest",
			file:    "testdata/rockbands.txt",
		},
		{
			name:    "no match in fruits — grape",
			pattern: "grape",
			file:    "testdata/fruits.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, _, exitCode := runBinary(tt.pattern, tt.file)
			if exitCode != 1 {
				t.Errorf("expected exit code 1 for no match, got %d", exitCode)
			}
			if stdout != "" {
				t.Errorf("expected no output for no match, got: %q", stdout)
			}
		})
	}
}

func TestSearchMissingArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no arguments",
			args: []string{},
		},
		{
			name: "only pattern",
			args: []string{"hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, stderr, exitCode := runBinary(tt.args...)
			if exitCode != 2 {
				t.Errorf("expected exit code 2, got %d", exitCode)
			}
			expected := "usage: search <pattern> <filename>"
			if stderr != expected {
				t.Errorf("expected stderr: %q, got: %q", expected, stderr)
			}
		})
	}
}

func TestSearchFileNotFound(t *testing.T) {
	_, stderr, exitCode := runBinary("hello", "nonexistent.txt")
	if exitCode != 2 {
		t.Errorf("expected exit code 2, got %d", exitCode)
	}
	expected := "error: cannot open 'nonexistent.txt'"
	if stderr != expected {
		t.Errorf("expected stderr: %q, got: %q", expected, stderr)
	}
}

func TestSearchOutputFormat(t *testing.T) {
	stdout, _, exitCode := runBinary("Priest", "testdata/rockbands.txt")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	// Verify format is "filename:line"
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "testdata/rockbands.txt:") {
			t.Errorf("output line should be prefixed with filename, got: %q", line)
		}
	}
}

func TestSearchPreservesOrder(t *testing.T) {
	stdout, _, exitCode := runBinary("White", "testdata/rockbands.txt")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	lines := strings.Split(stdout, "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 matching lines, got %d: %q", len(lines), stdout)
	}

	// Whitesnake (line 27) should come before Great White (line 49)
	if !strings.Contains(lines[0], "Whitesnake") {
		t.Errorf("first match should be Whitesnake, got: %s", lines[0])
	}
	if !strings.Contains(lines[1], "Great White") {
		t.Errorf("second match should be Great White, got: %s", lines[1])
	}
}
