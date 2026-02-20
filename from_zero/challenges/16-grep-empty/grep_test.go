package main_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Build the binary
	build := exec.Command("go", "build", "-o", "mygrep", ".")
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		os.Exit(1)
	}

	binaryPath = "./mygrep"
	code := m.Run()

	os.Remove("mygrep")
	os.Exit(code)
}

func runGrep(args ...string) (string, int) {
	cmd := exec.Command(binaryPath, args...)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}
	return string(out), exitCode
}

func runGrepStdout(args ...string) (string, string, int) {
	cmd := exec.Command(binaryPath, args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}
	return stdout.String(), stderr.String(), exitCode
}

// Test: empty pattern matches every line
func TestEmptyPattern(t *testing.T) {
	stdout, _, exitCode := runGrepStdout("", "testdata/test.txt")

	// Read the original file
	data, err := os.ReadFile("testdata/test.txt")
	if err != nil {
		t.Fatalf("could not read test.txt: %v", err)
	}

	if stdout != string(data) {
		t.Errorf("empty pattern should output entire file.\nExpected %d bytes, got %d bytes", len(data), len(stdout))
	}

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
}

// Test: empty pattern on rockbands.txt
func TestEmptyPatternRockbands(t *testing.T) {
	stdout, _, exitCode := runGrepStdout("", "testdata/rockbands.txt")

	data, err := os.ReadFile("testdata/rockbands.txt")
	if err != nil {
		t.Fatalf("could not read rockbands.txt: %v", err)
	}

	if stdout != string(data) {
		t.Errorf("empty pattern should output entire file")
	}

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
}

// Test: literal pattern match
func TestLiteralMatch(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		file     string
		contains []string
		exitCode int
	}{
		{
			name:     "match Gutenberg in test.txt",
			pattern:  "Gutenberg",
			file:     "testdata/test.txt",
			contains: []string{"Gutenberg"},
			exitCode: 0,
		},
		{
			name:     "match Nirvana in rockbands",
			pattern:  "Nirvana",
			file:     "testdata/rockbands.txt",
			contains: []string{"Nirvana"},
			exitCode: 0,
		},
		{
			name:     "no match returns exit 1",
			pattern:  "ZZZZNOTFOUND",
			file:     "testdata/rockbands.txt",
			contains: []string{},
			exitCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, _, exitCode := runGrepStdout(tt.pattern, tt.file)

			for _, expected := range tt.contains {
				if !strings.Contains(stdout, expected) {
					t.Errorf("expected output to contain %q, got:\n%s", expected, stdout)
				}
			}

			if exitCode != tt.exitCode {
				t.Errorf("expected exit code %d, got %d", tt.exitCode, exitCode)
			}
		})
	}
}

// Test: no arguments shows usage
func TestNoArgs(t *testing.T) {
	_, stderr, exitCode := runGrepStdout()

	if !strings.Contains(strings.ToLower(stderr), "usage") {
		t.Errorf("expected usage message on stderr, got: %s", stderr)
	}

	if exitCode != 2 {
		t.Errorf("expected exit code 2, got %d", exitCode)
	}
}

// Test: file not found
func TestFileNotFound(t *testing.T) {
	_, stderr, exitCode := runGrepStdout("pattern", "nonexistent.txt")

	if stderr == "" {
		t.Error("expected error message on stderr for missing file")
	}

	if exitCode != 2 {
		t.Errorf("expected exit code 2, got %d", exitCode)
	}
}
