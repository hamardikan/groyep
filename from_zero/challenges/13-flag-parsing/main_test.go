// Run: go test -v
//
// This test builds your search program and checks flag behavior.
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

// ---------- Default behavior (no flags) — same as challenge 12 ----------

func TestDefaultBehavior(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		file     string
		expected string
	}{
		{
			name:     "single match",
			pattern:  "Priest",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:Judas Priest",
		},
		{
			name:     "multiple matches",
			pattern:  "Bad",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:Bad English\ntestdata/rockbands.txt:Bad Company",
		},
		{
			name:     "match in fruits",
			pattern:  "cherry",
			file:     "testdata/fruits.txt",
			expected: "testdata/fruits.txt:cherry",
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

// ---------- -n flag (line numbers) ----------

func TestFlagN(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		file     string
		expected string
	}{
		{
			name:     "line number — Priest is line 1",
			pattern:  "Priest",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:1:Judas Priest",
		},
		{
			name:     "line numbers — Bad at lines 19 and 25",
			pattern:  "Bad",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:19:Bad English\ntestdata/rockbands.txt:25:Bad Company",
		},
		{
			name:     "line number — Nirvana is line 101",
			pattern:  "Nirvana",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:101:Nirvana",
		},
		{
			name:     "line numbers in fruits — cherry is line 3",
			pattern:  "cherry",
			file:     "testdata/fruits.txt",
			expected: "testdata/fruits.txt:3:cherry",
		},
		{
			name:     "line numbers — multiple in fruits",
			pattern:  "a",
			file:     "testdata/fruits.txt",
			expected: "testdata/fruits.txt:1:apple\ntestdata/fruits.txt:2:banana\ntestdata/fruits.txt:4:apricot\ntestdata/fruits.txt:6:avocado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, exitCode := runBinary("-n", tt.pattern, tt.file)
			if exitCode != 0 {
				t.Fatalf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
			}
			if stdout != tt.expected {
				t.Errorf("search -n %q %s\n  expected:\n%s\n  actual:\n%s", tt.pattern, tt.file, tt.expected, stdout)
			}
		})
	}
}

// ---------- -c flag (count only) ----------

func TestFlagC(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		file     string
		expected string
	}{
		{
			name:     "count — Priest (1 match)",
			pattern:  "Priest",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:1",
		},
		{
			name:     "count — Bad (2 matches)",
			pattern:  "Bad",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:2",
		},
		{
			name:     "count — White (4 matches)",
			pattern:  "White",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:4",
		},
		{
			name:     "count — a in fruits (4 matches)",
			pattern:  "a",
			file:     "testdata/fruits.txt",
			expected: "testdata/fruits.txt:4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, exitCode := runBinary("-c", tt.pattern, tt.file)
			if exitCode != 0 {
				t.Fatalf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
			}
			if stdout != tt.expected {
				t.Errorf("search -c %q %s\n  expected:\n%s\n  actual:\n%s", tt.pattern, tt.file, tt.expected, stdout)
			}
		})
	}
}

// ---------- -n and -c combined (count wins) ----------

func TestFlagNCCombined(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "-n -c together",
			args:     []string{"-n", "-c", "Bad", "testdata/rockbands.txt"},
			expected: "testdata/rockbands.txt:2",
		},
		{
			name:     "-c -n together (reversed order)",
			args:     []string{"-c", "-n", "Bad", "testdata/rockbands.txt"},
			expected: "testdata/rockbands.txt:2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, exitCode := runBinary(tt.args...)
			if exitCode != 0 {
				t.Fatalf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
			}
			if stdout != tt.expected {
				t.Errorf("search %v\n  expected:\n%s\n  actual:\n%s", tt.args, tt.expected, stdout)
			}
		})
	}
}

// ---------- No match with flags ----------

func TestFlagNoMatch(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "-n no match",
			args: []string{"-n", "zzz", "testdata/rockbands.txt"},
		},
		{
			name: "-c no match",
			args: []string{"-c", "zzz", "testdata/rockbands.txt"},
		},
		{
			name: "-n -c no match",
			args: []string{"-n", "-c", "zzz", "testdata/rockbands.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, _, exitCode := runBinary(tt.args...)
			if exitCode != 1 {
				t.Errorf("expected exit code 1 for no match, got %d", exitCode)
			}
			if stdout != "" {
				t.Errorf("expected no stdout for no match, got: %q", stdout)
			}
		})
	}
}

// ---------- Error handling (same as challenge 12) ----------

func TestErrorMissingArgs(t *testing.T) {
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
		{
			name: "flag but no positional args",
			args: []string{"-n"},
		},
		{
			name: "flag and pattern but no filename",
			args: []string{"-n", "hello"},
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

func TestErrorFileNotFound(t *testing.T) {
	_, stderr, exitCode := runBinary("hello", "nonexistent.txt")
	if exitCode != 2 {
		t.Errorf("expected exit code 2, got %d", exitCode)
	}
	expected := "error: cannot open 'nonexistent.txt'"
	if stderr != expected {
		t.Errorf("expected stderr: %q, got: %q", expected, stderr)
	}
}

func TestErrorFileNotFoundWithFlags(t *testing.T) {
	_, stderr, exitCode := runBinary("-n", "hello", "nonexistent.txt")
	if exitCode != 2 {
		t.Errorf("expected exit code 2, got %d", exitCode)
	}
	expected := "error: cannot open 'nonexistent.txt'"
	if stderr != expected {
		t.Errorf("expected stderr: %q, got: %q", expected, stderr)
	}
}
