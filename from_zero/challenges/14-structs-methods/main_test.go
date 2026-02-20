// Run: go test -v
//
// This test builds your search program and checks all behavior from challenge 13
// plus new multi-file support. You do NOT need to edit this file.

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

// ======================================================================
// SINGLE FILE — backward compatible with challenge 13
// ======================================================================

func TestSingleFileDefault(t *testing.T) {
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

func TestSingleFileFlagN(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		file     string
		expected string
	}{
		{
			name:     "-n Priest line 1",
			pattern:  "Priest",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:1:Judas Priest",
		},
		{
			name:     "-n Bad lines 19, 25",
			pattern:  "Bad",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:19:Bad English\ntestdata/rockbands.txt:25:Bad Company",
		},
		{
			name:     "-n cherry line 3 in fruits",
			pattern:  "cherry",
			file:     "testdata/fruits.txt",
			expected: "testdata/fruits.txt:3:cherry",
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

func TestSingleFileFlagC(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		file     string
		expected string
	}{
		{
			name:     "-c Bad count 2",
			pattern:  "Bad",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:2",
		},
		{
			name:     "-c White count 4",
			pattern:  "White",
			file:     "testdata/rockbands.txt",
			expected: "testdata/rockbands.txt:4",
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

func TestSingleFileFlagNCCombined(t *testing.T) {
	stdout, _, exitCode := runBinary("-n", "-c", "Bad", "testdata/rockbands.txt")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	expected := "testdata/rockbands.txt:2"
	if stdout != expected {
		t.Errorf("-n -c combined\n  expected: %s\n  actual: %s", expected, stdout)
	}
}

func TestSingleFileNoMatch(t *testing.T) {
	stdout, _, exitCode := runBinary("zzz", "testdata/rockbands.txt")
	if exitCode != 1 {
		t.Errorf("expected exit code 1 for no match, got %d", exitCode)
	}
	if stdout != "" {
		t.Errorf("expected no output for no match, got: %q", stdout)
	}
}

// ======================================================================
// MULTI-FILE SUPPORT — new in challenge 14
// ======================================================================

func TestMultiFileDefault(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		files    []string
		expected string
	}{
		{
			name:     "match in first file only",
			pattern:  "Bad",
			files:    []string{"testdata/rockbands.txt", "testdata/fruits.txt"},
			expected: "testdata/rockbands.txt:Bad English\ntestdata/rockbands.txt:Bad Company",
		},
		{
			name:     "match in second file only",
			pattern:  "cherry",
			files:    []string{"testdata/rockbands.txt", "testdata/fruits.txt"},
			expected: "testdata/fruits.txt:cherry",
		},
		{
			name:    "match in both files",
			pattern: "an",
			files:   []string{"testdata/rockbands.txt", "testdata/fruits.txt"},
			expected: strings.Join([]string{
				"testdata/rockbands.txt:Van Halen",
				"testdata/rockbands.txt:Damn Yankees",
				"testdata/rockbands.txt:Bad Company",
				"testdata/rockbands.txt:Vandenberg",
				"testdata/rockbands.txt:Hanoi Rocks",
				"testdata/rockbands.txt:Danger Danger",
				"testdata/rockbands.txt:Dangerous Toys",
				"testdata/rockbands.txt:Night Ranger",
				"testdata/rockbands.txt:Giant",
				"testdata/rockbands.txt:Danzig",
				"testdata/rockbands.txt:Hurricane",
				"testdata/rockbands.txt:SouthGang",
				"testdata/rockbands.txt:Bang Tango",
				"testdata/rockbands.txt:Nirvana",
				"testdata/fruits.txt:banana",
			}, "\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := append([]string{tt.pattern}, tt.files...)
			stdout, stderr, exitCode := runBinary(args...)
			if exitCode != 0 {
				t.Fatalf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
			}
			if stdout != tt.expected {
				t.Errorf("search %q %v\n  expected:\n%s\n  actual:\n%s", tt.pattern, tt.files, tt.expected, stdout)
			}
		})
	}
}

func TestMultiFileFlagN(t *testing.T) {
	args := []string{"-n", "cherry", "testdata/rockbands.txt", "testdata/fruits.txt"}
	stdout, _, exitCode := runBinary(args...)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	expected := "testdata/fruits.txt:3:cherry"
	if stdout != expected {
		t.Errorf("search -n cherry (multi-file)\n  expected: %s\n  actual: %s", expected, stdout)
	}
}

func TestMultiFileFlagC(t *testing.T) {
	args := []string{"-c", "a", "testdata/rockbands.txt", "testdata/fruits.txt"}
	stdout, _, exitCode := runBinary(args...)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	expected := "testdata/rockbands.txt:52\ntestdata/fruits.txt:4"
	if stdout != expected {
		t.Errorf("search -c a (multi-file)\n  expected:\n%s\n  actual:\n%s", expected, stdout)
	}
}

func TestMultiFileNoMatchAnywhere(t *testing.T) {
	args := []string{"zzz", "testdata/rockbands.txt", "testdata/fruits.txt"}
	stdout, _, exitCode := runBinary(args...)
	if exitCode != 1 {
		t.Errorf("expected exit code 1 when no match in any file, got %d", exitCode)
	}
	if stdout != "" {
		t.Errorf("expected no output, got: %q", stdout)
	}
}

func TestMultiFileOneFileMissing(t *testing.T) {
	args := []string{"Bad", "testdata/rockbands.txt", "nonexistent.txt"}
	stdout, stderr, exitCode := runBinary(args...)

	// Should still print matches from the file that exists
	if !strings.Contains(stdout, "testdata/rockbands.txt:Bad English") {
		t.Errorf("should still output matches from valid file, got stdout: %q", stdout)
	}

	// Should print error for missing file
	if !strings.Contains(stderr, "error: cannot open 'nonexistent.txt'") {
		t.Errorf("should print error for missing file, got stderr: %q", stderr)
	}

	// Exit code should be 2 when there's a file error
	if exitCode != 2 {
		t.Errorf("expected exit code 2 when a file is missing, got %d", exitCode)
	}
}

func TestMultiFileOrder(t *testing.T) {
	// Results from first file should come before second file
	args := []string{"apple", "testdata/fruits.txt", "testdata/rockbands.txt"}
	stdout, _, exitCode := runBinary(args...)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	expected := "testdata/fruits.txt:apple"
	if stdout != expected {
		t.Errorf("expected: %q, got: %q", expected, stdout)
	}
}

// ======================================================================
// ERROR HANDLING
// ======================================================================

func TestErrorMissingArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "no arguments", args: []string{}},
		{name: "only pattern", args: []string{"hello"}},
		{name: "flag but no args", args: []string{"-n"}},
		{name: "flag and pattern but no file", args: []string{"-n", "hello"}},
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
