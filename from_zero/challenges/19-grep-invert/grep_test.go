package main_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

var binaryPath string

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

func runGrepWithStdin(stdin string, args ...string) (stdout string, exitCode int) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Stdin = strings.NewReader(stdin)
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

func TestInvertBasic(t *testing.T) {
	// -v "o" on symbols.txt should exclude "pound" and "dollar" (both contain 'o')
	stdout, exitCode := runGrep("-v", "o", "testdata/symbols.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	// Should NOT contain "pound" or "dollar"
	if strings.Contains(stdout, "pound") || strings.Contains(stdout, "dollar") {
		t.Errorf("inverted match should exclude 'pound' and 'dollar', got:\n%s", stdout)
	}

	// Should contain symbol lines
	if !strings.Contains(stdout, "!") {
		t.Errorf("inverted match should include '!', got:\n%s", stdout)
	}
}

func TestInvertPipeComposition(t *testing.T) {
	// Simulate: mygrep -r Nirvana testdata/ | mygrep -v Madonna
	// First get recursive Nirvana results
	nirvanaOut, _ := runGrep("-r", "Nirvana", "testdata/")

	// Pipe those through -v Madonna
	filtered, exitCode := runGrepWithStdin(nirvanaOut, "-v", "Madonna")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	lines := strings.Split(strings.TrimRight(filtered, "\n"), "\n")

	// Should only have rockbands.txt:Nirvana left
	if len(lines) != 1 {
		t.Errorf("expected 1 line after filtering Madonna, got %d:\n%s", len(lines), filtered)
	}

	if !strings.Contains(filtered, "Nirvana") {
		t.Errorf("expected output to contain 'Nirvana', got: %s", filtered)
	}
}

func TestInvertWithRecursive(t *testing.T) {
	// -rv "e" in BFS1985 should give lines without 'e'
	stdout, exitCode := runGrep("-rv", "e", "testdata/test-subdir/BFS1985.txt")

	if exitCode != 0 && exitCode != 1 {
		t.Errorf("expected exit 0 or 1, got %d", exitCode)
	}

	// Every output line should NOT contain 'e'
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	for _, line := range lines {
		// Strip filepath prefix if present
		content := line
		if idx := strings.Index(line, ":"); idx >= 0 {
			content = line[idx+1:]
		}
		if content == "" {
			continue
		}
		if strings.Contains(content, "e") {
			t.Errorf("inverted match should not contain 'e': %s", line)
		}
	}
}

func TestInvertAllMatch(t *testing.T) {
	// If every line matches, -v should produce no output, exit 1
	// Empty pattern matches everything, so -v "" should match nothing
	_, exitCode := runGrep("-v", "", "testdata/symbols.txt")
	if exitCode != 1 {
		t.Errorf("invert with empty pattern (all match): expected exit 1, got %d", exitCode)
	}
}

func TestInvertExitCodes(t *testing.T) {
	// -v with a pattern that matches nothing: all lines pass through, exit 0
	stdout, exitCode := runGrep("-v", "ZZZNOTFOUND", "testdata/rockbands.txt")
	if exitCode != 0 {
		t.Errorf("invert with no-match pattern: expected exit 0, got %d", exitCode)
	}
	// Should output all lines
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	if len(lines) < 100 {
		t.Errorf("expected ~101 lines, got %d", len(lines))
	}
}
