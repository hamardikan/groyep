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

func TestCaretAnchor(t *testing.T) {
	stdout, exitCode := runGrep("^A", "testdata/rockbands.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	expected := []string{"AC/DC", "Aerosmith", "Accept", "April Wine", "Autograph"}

	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")

	if len(lines) != len(expected) {
		t.Errorf("expected %d matches for ^A, got %d:\n%s", len(expected), len(lines), stdout)
	}

	for _, exp := range expected {
		if !strings.Contains(stdout, exp) {
			t.Errorf("expected %q in output", exp)
		}
	}
}

func TestDollarAnchor(t *testing.T) {
	stdout, exitCode := runGrep("na$", "testdata/rockbands.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")

	if len(lines) != 1 {
		t.Errorf("expected 1 match for na$, got %d:\n%s", len(lines), stdout)
	}

	if !strings.Contains(stdout, "Nirvana") {
		t.Errorf("expected 'Nirvana' for na$ pattern, got: %s", stdout)
	}
}

func TestEmptyLineAnchor(t *testing.T) {
	// ^$ should match empty lines in BFS1985
	stdout, exitCode := runGrep("^$", "testdata/test-subdir/BFS1985.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	// BFS1985 has 5 empty lines (between verses)
	if len(lines) < 4 {
		t.Errorf("expected at least 4 empty line matches, got %d", len(lines))
	}
}

func TestExactLengthAnchor(t *testing.T) {
	// ^.{3}$ matches lines with exactly 3 characters
	stdout, exitCode := runGrep("^.{3}$", "testdata/rockbands.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	// UFO, Kix, TNT are 3-char bands
	expected := []string{"UFO", "Kix", "TNT"}
	for _, exp := range expected {
		if !strings.Contains(stdout, exp) {
			t.Errorf("expected %q for ^.{3}$, got:\n%s", exp, stdout)
		}
	}
}

func TestAnchorWithInvert(t *testing.T) {
	// -v "^[A-Z]" on symbols should return lines NOT starting with uppercase letter
	stdout, exitCode := runGrep("-v", "^[A-Z]", "testdata/symbols.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	// All symbol characters don't start with uppercase, but "pound" and "dollar" are lowercase
	if !strings.Contains(stdout, "!") {
		t.Errorf("expected '!' in output")
	}
}

func TestCaretNotFirstChar(t *testing.T) {
	// "A" without anchor matches anywhere
	stdoutAnchored, _ := runGrep("^A", "testdata/rockbands.txt")
	stdoutUnanchored, _ := runGrep("A", "testdata/rockbands.txt")

	anchoredCount := len(strings.Split(strings.TrimRight(stdoutAnchored, "\n"), "\n"))
	unanchoredCount := len(strings.Split(strings.TrimRight(stdoutUnanchored, "\n"), "\n"))

	if anchoredCount >= unanchoredCount {
		t.Errorf("^A should match fewer lines than A: ^A=%d, A=%d", anchoredCount, unanchoredCount)
	}
}
