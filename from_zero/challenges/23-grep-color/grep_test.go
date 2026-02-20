package main_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

var binaryPath string

const (
	colorStart = "\033[1;31m"
	colorReset = "\033[0m"
)

func TestMain(m *testing.M) {
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

func runGrep(args ...string) (stdout string, exitCode int) {
	cmd := exec.Command(binaryPath, args...)
	var outBuf strings.Builder
	cmd.Stdout = &outBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}
	return outBuf.String(), exitCode
}

func TestColorAlways(t *testing.T) {
	// --color=always or -color=always should wrap matched text in ANSI codes
	stdout, exitCode := runGrep("-color=always", "Nirvana", "testdata/rockbands.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	if !strings.Contains(stdout, colorStart) {
		t.Errorf("--color=always should contain ANSI start code, got:\n%q", stdout)
	}
	if !strings.Contains(stdout, colorReset) {
		t.Errorf("--color=always should contain ANSI reset code")
	}
	if !strings.Contains(stdout, colorStart+"Nirvana"+colorReset) {
		t.Errorf("expected 'Nirvana' wrapped in color codes, got:\n%q", stdout)
	}
}

func TestColorNever(t *testing.T) {
	stdout, exitCode := runGrep("-color=never", "Nirvana", "testdata/rockbands.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	if strings.Contains(stdout, "\033[") {
		t.Errorf("--color=never should NOT contain ANSI codes, got:\n%q", stdout)
	}
}

func TestNoColorFlag(t *testing.T) {
	// Without --color flag, no ANSI codes (default is never)
	stdout, _ := runGrep("Nirvana", "testdata/rockbands.txt")

	if strings.Contains(stdout, "\033[") {
		t.Errorf("no color flag should NOT produce ANSI codes, got:\n%q", stdout)
	}
}

func TestColorMultipleMatches(t *testing.T) {
	// "an" in "Danger Danger" should colorize both occurrences
	stdout, exitCode := runGrep("-color=always", "an", "testdata/rockbands.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	// Find the "Danger Danger" line
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Contains(line, "Danger") && strings.Contains(line, "Danger") {
			// Should have at least 2 color start codes
			count := strings.Count(line, colorStart)
			if count < 2 {
				t.Errorf("'Danger Danger' should have >=2 color starts for 'an', got %d in: %q", count, line)
			}
			break
		}
	}
}

func TestColorWithRegex(t *testing.T) {
	stdout, exitCode := runGrep("-color=always", `\d+`, "testdata/test-subdir/BFS1985.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	if !strings.Contains(stdout, colorStart) {
		t.Errorf("color with regex should contain ANSI codes")
	}
}

func TestColorPreservesContent(t *testing.T) {
	// With color, the underlying text content should be the same
	colorOut, _ := runGrep("-color=always", "Nirvana", "testdata/rockbands.txt")
	plainOut, _ := runGrep("Nirvana", "testdata/rockbands.txt")

	// Strip ANSI codes from color output
	stripped := strings.ReplaceAll(colorOut, colorStart, "")
	stripped = strings.ReplaceAll(stripped, colorReset, "")

	if stripped != plainOut {
		t.Errorf("color output (stripped) should match plain output\nStripped: %q\nPlain: %q", stripped, plainOut)
	}
}

func TestPreviousFeaturesStillWork(t *testing.T) {
	// -r still works
	stdout, exitCode := runGrep("-r", "Nirvana", "testdata/")
	if exitCode != 0 {
		t.Errorf("-r: expected exit 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Nirvana") {
		t.Errorf("-r: expected Nirvana in output")
	}

	// -v still works
	_, exitCode = runGrep("-v", "ZZZZ", "testdata/rockbands.txt")
	if exitCode != 0 {
		t.Errorf("-v: expected exit 0, got %d", exitCode)
	}

	// -i still works
	stdout, exitCode = runGrep("-i", "nirvana", "testdata/rockbands.txt")
	if exitCode != 0 {
		t.Errorf("-i: expected exit 0, got %d", exitCode)
	}
}
