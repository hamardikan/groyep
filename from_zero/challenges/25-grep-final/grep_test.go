package main_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Try to build from cmd/mygrep if it exists, otherwise from current dir
	var build *exec.Cmd
	if _, err := os.Stat("cmd/mygrep/main.go"); err == nil {
		build = exec.Command("go", "build", "-o", "mygrep", "./cmd/mygrep")
	} else {
		build = exec.Command("go", "build", "-o", "mygrep", ".")
	}
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

func countLines(s string) int {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return 0
	}
	return len(strings.Split(s, "\n"))
}

// === STEP 1: Empty match ===

func TestEmptyPattern(t *testing.T) {
	stdout, _, exitCode := runGrep("", "testdata/test.txt")
	data, _ := os.ReadFile("testdata/test.txt")

	if stdout != string(data) {
		t.Errorf("empty pattern should output entire file")
	}
	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}
}

// === STEP 2: Literal match + exit codes ===

func TestLiteralJ(t *testing.T) {
	stdout, _, exitCode := runGrep("J", "testdata/rockbands.txt")

	expected := []string{"Judas Priest", "Bon Jovi", "Junkyard"}
	for _, exp := range expected {
		if !strings.Contains(stdout, exp) {
			t.Errorf("expected %q in output", exp)
		}
	}
	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}
}

func TestExitCodes(t *testing.T) {
	tests := []struct {
		name string
		args []string
		exit int
	}{
		{"match found", []string{"Nirvana", "testdata/rockbands.txt"}, 0},
		{"no match", []string{"ZZZZZ", "testdata/rockbands.txt"}, 1},
		{"no args", []string{}, 2},
		{"file not found", []string{"pattern", "nonexistent.txt"}, 2},
		{"invalid regex", []string{"[bad", "testdata/rockbands.txt"}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, exitCode := runGrep(tt.args...)
			if exitCode != tt.exit {
				t.Errorf("expected exit %d, got %d", tt.exit, exitCode)
			}
		})
	}
}

// === STEP 3: Recursive ===

func TestRecursiveNirvana(t *testing.T) {
	stdout, _, exitCode := runGrep("-r", "Nirvana", "testdata/")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	expected := []string{
		"testdata/rockbands.txt:Nirvana",
		"testdata/test-subdir/BFS1985.txt:",
	}
	for _, exp := range expected {
		if !strings.Contains(stdout, exp) {
			t.Errorf("expected output to contain %q", exp)
		}
	}

	// 5 total Nirvana lines
	nirvanaCount := 0
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Contains(line, "Nirvana") {
			nirvanaCount++
		}
	}
	if nirvanaCount != 5 {
		t.Errorf("expected 5 Nirvana matches, got %d", nirvanaCount)
	}
}

// === STEP 4: Invert ===

func TestInvertPipe(t *testing.T) {
	nirvanaOut, _, _ := runGrep("-r", "Nirvana", "testdata/")
	filtered, exitCode := runGrepWithStdin(nirvanaOut, "-v", "Madonna")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	lines := strings.Split(strings.TrimRight(filtered, "\n"), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line after -v Madonna, got %d:\n%s", len(lines), filtered)
	}
	if !strings.Contains(filtered, "Nirvana") {
		t.Errorf("expected 'Nirvana' in filtered output")
	}
}

// === STEP 5: Regex ===

func TestRegexDigits(t *testing.T) {
	stdout, _, exitCode := runGrep(`\d`, "testdata/test-subdir/BFS1985.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}
	if countLines(stdout) != 9 {
		t.Errorf("expected 9 digit-matching lines, got %d", countLines(stdout))
	}
}

func TestRegexWordChar(t *testing.T) {
	stdout, _, exitCode := runGrep(`\w`, "testdata/symbols.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 \\w matches in symbols.txt, got %d:\n%s", len(lines), stdout)
	}
}

// === STEP 6: Anchors ===

func TestAnchorCaret(t *testing.T) {
	stdout, _, _ := runGrep("^A", "testdata/rockbands.txt")

	expected := []string{"AC/DC", "Aerosmith", "Accept", "April Wine", "Autograph"}
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	if len(lines) != len(expected) {
		t.Errorf("^A: expected %d matches, got %d", len(expected), len(lines))
	}
}

func TestAnchorDollar(t *testing.T) {
	stdout, _, _ := runGrep("na$", "testdata/rockbands.txt")

	if strings.TrimSpace(stdout) != "Nirvana" {
		t.Errorf("na$: expected 'Nirvana', got %q", strings.TrimSpace(stdout))
	}
}

// === STEP 7: Case insensitive ===

func TestCaseSensitiveA(t *testing.T) {
	stdout, _, _ := runGrep("A", "testdata/rockbands.txt")
	if countLines(stdout) != 8 {
		t.Errorf("case-sensitive A: expected 8, got %d", countLines(stdout))
	}
}

func TestCaseInsensitiveA(t *testing.T) {
	stdout, _, _ := runGrep("-i", "A", "testdata/rockbands.txt")
	if countLines(stdout) != 58 {
		t.Errorf("case-insensitive A: expected 58, got %d", countLines(stdout))
	}
}

// === FLAG COMBINATIONS ===

func TestCombinedRI(t *testing.T) {
	stdout, _, exitCode := runGrep("-ri", "nirvana", "testdata/")
	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Nirvana") {
		t.Errorf("-ri: expected Nirvana in output")
	}
}

func TestCombinedRV(t *testing.T) {
	_, _, exitCode := runGrep("-rv", "ZZZZZ", "testdata/")
	if exitCode != 0 {
		t.Errorf("-rv with no-match pattern: expected exit 0, got %d", exitCode)
	}
}

func TestCombinedRIV(t *testing.T) {
	stdout, _, exitCode := runGrep("-riv", "nirvana", "testdata/rockbands.txt")
	if exitCode != 0 {
		t.Errorf("-riv: expected exit 0, got %d", exitCode)
	}
	// Should NOT contain Nirvana (inverted case-insensitive)
	for _, line := range strings.Split(stdout, "\n") {
		content := line
		if idx := strings.Index(line, ":"); idx >= 0 {
			content = line[idx+1:]
		}
		if strings.Contains(strings.ToLower(content), "nirvana") {
			t.Errorf("-riv: should NOT contain nirvana, found: %s", line)
		}
	}
}

// === STDIN ===

func TestStdinPipe(t *testing.T) {
	stdout, exitCode := runGrepWithStdin("hello world\nfoo bar\nhello again\n", "hello")
	if exitCode != 0 {
		t.Errorf("stdin: expected exit 0, got %d", exitCode)
	}
	if countLines(stdout) != 2 {
		t.Errorf("stdin: expected 2 matches, got %d", countLines(stdout))
	}
}

// === LINE NUMBERS (-n) ===

func TestLineNumbers(t *testing.T) {
	stdout, _, exitCode := runGrep("-n", "Nirvana", "testdata/rockbands.txt")
	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}
	// Should contain line number prefix like "101:Nirvana"
	if !strings.Contains(stdout, ":Nirvana") {
		t.Errorf("-n should have line number prefix, got: %s", stdout)
	}
	// Check it starts with a number
	line := strings.TrimSpace(stdout)
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		t.Errorf("-n output should be NUMBER:LINE, got: %s", line)
	}
}

// === COUNT (-c) ===

func TestCount(t *testing.T) {
	stdout, _, exitCode := runGrep("-c", "Bad", "testdata/rockbands.txt")
	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}
	if strings.TrimSpace(stdout) != "2" {
		t.Errorf("-c Bad: expected '2', got %q", strings.TrimSpace(stdout))
	}
}
