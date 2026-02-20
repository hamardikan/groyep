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

func countLines(s string) int {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return 0
	}
	return len(strings.Split(s, "\n"))
}

func TestCaseSensitiveCount(t *testing.T) {
	stdout, _ := runGrep("A", "testdata/rockbands.txt")
	count := countLines(stdout)
	if count != 8 {
		t.Errorf("case-sensitive 'A': expected 8 matches, got %d", count)
	}
}

func TestCaseInsensitiveCount(t *testing.T) {
	stdout, _ := runGrep("-i", "A", "testdata/rockbands.txt")
	count := countLines(stdout)
	if count != 58 {
		t.Errorf("case-insensitive 'A': expected 58 matches, got %d", count)
	}
}

func TestCaseInsensitiveBasic(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name:     "-i nirvana finds Nirvana",
			args:     []string{"-i", "nirvana", "testdata/rockbands.txt"},
			contains: []string{"Nirvana"},
		},
		{
			name:     "-i NIRVANA finds Nirvana",
			args:     []string{"-i", "NIRVANA", "testdata/rockbands.txt"},
			contains: []string{"Nirvana"},
		},
		{
			name:     "-i with anchor ^a",
			args:     []string{"-i", "^a", "testdata/rockbands.txt"},
			contains: []string{"AC/DC", "Aerosmith", "Accept", "April Wine", "Autograph"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, exitCode := runGrep(tt.args...)
			if exitCode != 0 {
				t.Errorf("expected exit 0, got %d", exitCode)
			}
			for _, exp := range tt.contains {
				if !strings.Contains(stdout, exp) {
					t.Errorf("expected output to contain %q\nGot:\n%s", exp, stdout)
				}
			}
		})
	}
}

func TestCaseInsensitiveRecursive(t *testing.T) {
	stdout, exitCode := runGrep("-ri", "nirvana", "testdata/")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	if !strings.Contains(stdout, "rockbands.txt") {
		t.Errorf("expected rockbands.txt in recursive output")
	}
	if !strings.Contains(stdout, "BFS1985.txt") {
		t.Errorf("expected BFS1985.txt in recursive output")
	}
}

func TestCaseInsensitiveInvert(t *testing.T) {
	// -iv "a" should exclude lines containing a or A
	stdout, exitCode := runGrep("-iv", "a", "testdata/symbols.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	// Should have the symbol lines (which don't have 'a' or 'A')
	if !strings.Contains(stdout, "!") {
		t.Errorf("expected '!' in output")
	}
}

func TestCaseInsensitiveRegex(t *testing.T) {
	// -i with regex character class
	stdout, exitCode := runGrep("-i", "^[a]c", "testdata/rockbands.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	if !strings.Contains(stdout, "AC/DC") {
		t.Errorf("expected AC/DC in output for -i ^[a]c, got: %s", stdout)
	}
	if !strings.Contains(stdout, "Accept") {
		t.Errorf("expected Accept in output for -i ^[a]c, got: %s", stdout)
	}
}

func TestWithoutCaseFlag(t *testing.T) {
	// Without -i, "nirvana" (lowercase) should NOT match "Nirvana"
	_, exitCode := runGrep("nirvana", "testdata/rockbands.txt")
	if exitCode != 1 {
		t.Errorf("'nirvana' without -i should not match 'Nirvana', expected exit 1, got %d", exitCode)
	}
}
