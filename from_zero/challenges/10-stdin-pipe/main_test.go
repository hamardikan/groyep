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

	binary := filepath.Join(dir, "upper")
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

func TestPipedInput(t *testing.T) {
	binary := buildBinary(t)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple lowercase", "hello world\n", "HELLO WORLD\n"},
		{"already uppercase", "HELLO\n", "HELLO\n"},
		{"mixed case", "Hello World\n", "HELLO WORLD\n"},
		{"multiple lines", "hello\nworld\n", "HELLO\nWORLD\n"},
		{"with numbers", "abc 123 def\n", "ABC 123 DEF\n"},
		{"empty input", "", ""},
		{"unicode cafe", "café au lait\n", "CAFÉ AU LAIT\n"},
		{"special characters", "hello! @#$ world?\n", "HELLO! @#$ WORLD?\n"},
		{"tabs and spaces", "hello\tworld  foo\n", "HELLO\tWORLD  FOO\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary)
			cmd.Stdin = strings.NewReader(tt.input)
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("command failed: %v", err)
			}

			got := string(out)
			if got != tt.expected {
				t.Errorf("echo %q | upper\ngot:  %q\nwant: %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFileInput(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "testdata/mixed.txt")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	got := string(out)
	expected := "HELLO WORLD\nTHIS IS LOWERCASE\nTHIS IS UPPERCASE\nMIXED CASE TEXT\nGO IS A GREAT LANGUAGE\n123 NUMBERS STAY THE SAME\nCAFÉ AU LAIT\n"

	if got != expected {
		t.Errorf("upper testdata/mixed.txt\ngot:\n%s\nwant:\n%s", got, expected)
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

func TestFileTakesPriorityOverStdin(t *testing.T) {
	binary := buildBinary(t)

	// When both a file arg and stdin are provided, file should take priority
	cmd := exec.Command(binary, "testdata/mixed.txt")
	cmd.Stdin = strings.NewReader("this should be ignored\n")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	got := string(out)
	// Should contain file content, not stdin content
	if strings.Contains(got, "THIS SHOULD BE IGNORED") {
		t.Error("file argument should take priority over stdin")
	}
	if !strings.Contains(got, "HELLO WORLD") {
		t.Errorf("expected file content in output, got: %q", got)
	}
}

func TestExitCodeSuccess(t *testing.T) {
	binary := buildBinary(t)

	// Piped input should exit 0
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader("hello\n")
	err := cmd.Run()
	if err != nil {
		t.Errorf("expected exit code 0 for piped input, got error: %v", err)
	}

	// File input should exit 0
	cmd = exec.Command(binary, "testdata/mixed.txt")
	err = cmd.Run()
	if err != nil {
		t.Errorf("expected exit code 0 for file input, got error: %v", err)
	}
}

func TestMultilinePipe(t *testing.T) {
	binary := buildBinary(t)

	input := "line one\nline two\nline three\n"
	expected := "LINE ONE\nLINE TWO\nLINE THREE\n"

	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	got := string(out)
	if got != expected {
		t.Errorf("multiline pipe\ngot:  %q\nwant: %q", got, expected)
	}
}
