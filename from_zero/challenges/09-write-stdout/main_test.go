package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func buildBinary(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	binary := filepath.Join(dir, "mynl")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, out)
	}

	t.Cleanup(func() {
		os.Remove(binary)
	})

	return binary
}

func TestFileInput(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "testdata/sample.txt")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	got := string(out)
	expected := "     1\tHello, World!\n     2\tThis is line two.\n     3\t\n     4\tThis is line four.\n     5\tThe end.\n"

	if got != expected {
		t.Errorf("mynl testdata/sample.txt\ngot:\n%s\nwant:\n%s", got, expected)
	}
}

func TestLineNumberFormat(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "testdata/sample.txt")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	lines := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")

	tests := []struct {
		lineIdx int
		prefix  string
		content string
	}{
		{0, "     1\t", "Hello, World!"},
		{1, "     2\t", "This is line two."},
		{2, "     3\t", ""},
		{3, "     4\t", "This is line four."},
		{4, "     5\t", "The end."},
	}

	if len(lines) != len(tests) {
		t.Fatalf("expected %d lines, got %d", len(tests), len(lines))
	}

	for _, tt := range tests {
		line := lines[tt.lineIdx]
		expectedLine := tt.prefix + tt.content
		if line != expectedLine {
			t.Errorf("line %d:\ngot:  %q\nwant: %q", tt.lineIdx+1, line, expectedLine)
		}
	}
}

func TestEmptyLineGetsNumber(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "testdata/sample.txt")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	lines := strings.Split(string(out), "\n")
	// Line 3 (index 2) should be the empty line with number
	if len(lines) < 3 {
		t.Fatal("not enough output lines")
	}

	emptyLine := lines[2]
	// Should be exactly "     3\t" (number + tab, no trailing content)
	expected := "     3\t"
	if emptyLine != expected {
		t.Errorf("empty line:\ngot:  %q\nwant: %q", emptyLine, expected)
	}
}

func TestStdinInput(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader("first\nsecond\nthird\n")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	got := string(out)
	expected := "     1\tfirst\n     2\tsecond\n     3\tthird\n"

	if got != expected {
		t.Errorf("echo ... | mynl\ngot:\n%s\nwant:\n%s", got, expected)
	}
}

func TestSingleLineInput(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader("only one line\n")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	got := string(out)
	expected := "     1\tonly one line\n"

	if got != expected {
		t.Errorf("single line stdin\ngot:  %q\nwant: %q", got, expected)
	}
}

func TestFileNotFound(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "testdata/nonexistent.txt")
	var stderr strings.Builder
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err == nil {
		t.Fatal("expected non-zero exit code for missing file")
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected *exec.ExitError, got %T", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", exitErr.ExitCode())
	}

	got := strings.TrimSpace(stderr.String())
	if !strings.Contains(got, "error:") || !strings.Contains(got, "no such file or directory") {
		t.Errorf("expected error message about missing file, got: %q", got)
	}
}

func TestTabSeparator(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader("test\n")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	got := string(out)
	// Verify the tab character is present
	if !strings.Contains(got, "\t") {
		t.Errorf("output should contain tab character, got: %q", got)
	}

	// Verify 6-char wide number field + tab
	if !strings.HasPrefix(got, "     1\t") {
		t.Errorf("expected 6-char right-aligned number + tab, got prefix: %q", got[:10])
	}
}
