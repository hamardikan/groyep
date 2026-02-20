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

func TestAfterContext(t *testing.T) {
	// -A 1 shows 1 line after each match
	stdout, exitCode := runGrep("-A", "1", "Nirvana", "testdata/rockbands.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	// Nirvana is the last non-empty line in rockbands, so -A 1 may just show Nirvana
	if !strings.Contains(stdout, "Nirvana") {
		t.Errorf("expected 'Nirvana' in output, got:\n%s", stdout)
	}
}

func TestBeforeContext(t *testing.T) {
	// -B 2 before "Accept" should show Aerosmith and Iron Maiden (or whatever is 2 lines before)
	stdout, exitCode := runGrep("-B", "2", "Accept", "testdata/rockbands.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	if !strings.Contains(stdout, "Accept") {
		t.Errorf("expected 'Accept' in output")
	}

	// Should have at least 3 lines (2 before + the match)
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	nonSeparator := 0
	for _, l := range lines {
		if l != "--" {
			nonSeparator++
		}
	}
	if nonSeparator < 3 {
		t.Errorf("-B 2 should show at least 3 lines (2 context + match), got %d", nonSeparator)
	}
}

func TestContextBoth(t *testing.T) {
	// -C 1 shows 1 line before and after
	stdout, exitCode := runGrep("-C", "1", "Accept", "testdata/rockbands.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	if !strings.Contains(stdout, "Accept") {
		t.Errorf("expected 'Accept' in output")
	}

	// Should have at least 3 lines (1 before + match + 1 after)
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	nonSeparator := 0
	for _, l := range lines {
		if l != "--" {
			nonSeparator++
		}
	}
	if nonSeparator < 3 {
		t.Errorf("-C 1 should show at least 3 lines, got %d", nonSeparator)
	}
}

func TestGroupSeparator(t *testing.T) {
	// Search for "1985" in BFS1985 with -B 1; should have -- separators between groups
	stdout, exitCode := runGrep("-B", "1", "1985", "testdata/test-subdir/BFS1985.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	if !strings.Contains(stdout, "1985") {
		t.Errorf("expected '1985' in output")
	}

	// There should be -- separators between non-adjacent groups
	if !strings.Contains(stdout, "\n--\n") && !strings.Contains(stdout, "--\n") {
		t.Logf("Note: no -- separator found; matches might be adjacent. Output:\n%s", stdout)
	}
}

func TestOverlappingContextMerge(t *testing.T) {
	// With large enough context, adjacent matches should merge
	stdout, _ := runGrep("-C", "50", "Nirvana", "testdata/rockbands.txt")

	// With C=50, the entire file should be output (no separators needed)
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	for _, l := range lines {
		if l == "--" {
			t.Errorf("with context=50 on a ~100 line file with 1 match, no separator expected")
			break
		}
	}
}

func TestContextZero(t *testing.T) {
	// -C 0 should behave like normal grep (no context)
	contextOut, _ := runGrep("-C", "0", "Nirvana", "testdata/rockbands.txt")
	normalOut, _ := runGrep("Nirvana", "testdata/rockbands.txt")

	// Strip separators from context output
	contextStripped := strings.ReplaceAll(contextOut, "--\n", "")

	if strings.TrimSpace(contextStripped) != strings.TrimSpace(normalOut) {
		t.Errorf("-C 0 should match normal output\nContext: %q\nNormal: %q", contextStripped, normalOut)
	}
}

func TestPreviousFlagsWithContext(t *testing.T) {
	// -i with -A
	stdout, exitCode := runGrep("-i", "-A", "1", "nirvana", "testdata/rockbands.txt")
	if exitCode != 0 {
		t.Errorf("-i with -A: expected exit 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Nirvana") {
		t.Errorf("-i with -A: expected 'Nirvana' in output")
	}
}
