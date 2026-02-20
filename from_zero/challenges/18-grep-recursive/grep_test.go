package main_test

import (
	"os"
	"os/exec"
	"sort"
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

func runGrep(args ...string) (stdout, stderr string, exitCode int) {
	cmd := exec.Command(binaryPath, args...)
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func sortedLines(s string) []string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	sort.Strings(lines)
	return lines
}

func TestRecursiveNirvana(t *testing.T) {
	stdout, _, exitCode := runGrep("-r", "Nirvana", "testdata/")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	expectedLines := []string{
		"testdata/rockbands.txt:Nirvana",
		"testdata/test-subdir/BFS1985.txt:Since Bruce Springsteen, Madonna, way before Nirvana",
		"testdata/test-subdir/BFS1985.txt:On the radio was Springsteen, Madonna, way before Nirvana",
		"testdata/test-subdir/BFS1985.txt:And bring back Springsteen, Madonna, way before Nirvana",
		"testdata/test-subdir/BFS1985.txt:Bruce Springsteen, Madonna, way before Nirvana",
	}

	for _, exp := range expectedLines {
		if !strings.Contains(stdout, exp) {
			t.Errorf("expected output to contain:\n  %s\nGot:\n%s", exp, stdout)
		}
	}
}

func TestRecursiveNoMatch(t *testing.T) {
	_, _, exitCode := runGrep("-r", "ZZZNOTFOUND", "testdata/")
	if exitCode != 1 {
		t.Errorf("no match recursive: expected exit 1, got %d", exitCode)
	}
}

func TestRecursiveSingleFile(t *testing.T) {
	stdout, _, exitCode := runGrep("-r", "Nirvana", "testdata/rockbands.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	if !strings.Contains(stdout, "Nirvana") {
		t.Errorf("expected output to contain 'Nirvana', got: %s", stdout)
	}
}

func TestRecursiveFilePrefix(t *testing.T) {
	stdout, _, _ := runGrep("-r", "Nirvana", "testdata/")
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")

	for _, line := range lines {
		if !strings.Contains(line, ":") {
			t.Errorf("recursive output should have filepath prefix, got: %s", line)
		}
	}
}

func TestNonRecursiveSingleFile(t *testing.T) {
	// Without -r, single file should still work
	stdout, _, exitCode := runGrep("Nirvana", "testdata/rockbands.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	if !strings.Contains(stdout, "Nirvana") {
		t.Errorf("expected output to contain 'Nirvana', got: %s", stdout)
	}
}

func TestRecursiveMTV(t *testing.T) {
	stdout, _, exitCode := runGrep("-r", "MTV", "testdata/")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	// BFS1985 mentions MTV 4 times
	mtvCount := 0
	for _, line := range lines {
		if strings.Contains(line, "MTV") {
			mtvCount++
		}
	}
	if mtvCount < 4 {
		t.Errorf("expected at least 4 MTV matches, got %d", mtvCount)
	}
}
