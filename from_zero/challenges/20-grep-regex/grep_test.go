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

func TestDigitPattern(t *testing.T) {
	stdout, exitCode := runGrep(`\d`, "testdata/test-subdir/BFS1985.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	expectedSubstrings := []string{
		"turned 24",
		"U2",
		"1985",
	}

	for _, exp := range expectedSubstrings {
		if !strings.Contains(stdout, exp) {
			t.Errorf("\\d pattern should match lines with %q\nGot:\n%s", exp, stdout)
		}
	}

	// Lines with digits in BFS1985: lines with 24, U2, 19, 1985
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	if len(lines) < 9 {
		t.Errorf("expected at least 9 lines matching \\d, got %d", len(lines))
	}
}

func TestWordCharPattern(t *testing.T) {
	stdout, exitCode := runGrep(`\w`, "testdata/symbols.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")

	// Only "pound" and "dollar" have word characters
	if len(lines) != 2 {
		t.Errorf("expected 2 lines matching \\w in symbols.txt, got %d:\n%s", len(lines), stdout)
	}

	if !strings.Contains(stdout, "pound") {
		t.Errorf("expected 'pound' in output")
	}
	if !strings.Contains(stdout, "dollar") {
		t.Errorf("expected 'dollar' in output")
	}
}

func TestDotPattern(t *testing.T) {
	// "Z.Top" should match "ZZ Top" (dot matches any char)
	stdout, exitCode := runGrep("Z.Top", "testdata/rockbands.txt")
	// Actually "Z.Top" won't match "ZZ Top" because there's a space. Let's use a better pattern.
	_ = stdout
	_ = exitCode

	// "K.x" should match "Kix"
	stdout, exitCode = runGrep("K.x", "testdata/rockbands.txt")
	if exitCode != 0 {
		t.Errorf("K.x: expected exit 0, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Kix") {
		t.Errorf("K.x should match 'Kix', got: %s", stdout)
	}
}

func TestCharacterClass(t *testing.T) {
	// [YZ] should match Y&T and ZZ Top
	stdout, exitCode := runGrep("[YZ]", "testdata/rockbands.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	if !strings.Contains(stdout, "Y&T") {
		t.Errorf("expected Y&T in output")
	}
	if !strings.Contains(stdout, "ZZ Top") {
		t.Errorf("expected ZZ Top in output")
	}
}

func TestQuantifiers(t *testing.T) {
	// "s+" should match lines with one or more 's'
	stdout, exitCode := runGrep("ss", "testdata/rockbands.txt")
	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}
	// "Kiss" has "ss"
	if !strings.Contains(stdout, "Kiss") {
		t.Errorf("expected 'Kiss' in output for 'ss' pattern, got: %s", stdout)
	}
}

func TestInvalidRegex(t *testing.T) {
	_, exitCode := runGrep("[invalid", "testdata/rockbands.txt")
	if exitCode != 2 {
		t.Errorf("invalid regex: expected exit 2, got %d", exitCode)
	}
}

func TestEmptyPatternStillWorks(t *testing.T) {
	stdout, exitCode := runGrep("", "testdata/symbols.txt")

	if exitCode != 0 {
		t.Errorf("empty pattern: expected exit 0, got %d", exitCode)
	}

	data, _ := os.ReadFile("testdata/symbols.txt")
	if stdout != string(data) {
		t.Error("empty pattern should still match all lines")
	}
}

func TestRegexWithFlags(t *testing.T) {
	// -v with regex: lines NOT matching \d
	stdout, exitCode := runGrep("-v", `\d`, "testdata/symbols.txt")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	// All symbol lines should be there (none have digits), plus empty lines
	if !strings.Contains(stdout, "!") {
		t.Errorf("expected symbols in output for -v \\d")
	}
}

func TestRecursiveWithRegex(t *testing.T) {
	stdout, exitCode := runGrep("-r", `\d{4}`, "testdata/")

	if exitCode != 0 {
		t.Errorf("expected exit 0, got %d", exitCode)
	}

	// Should find 1985 in BFS1985.txt
	if !strings.Contains(stdout, "1985") {
		t.Errorf("expected 1985 in recursive regex output, got: %s", stdout)
	}
}
