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

	binary := filepath.Join(dir, "mycat")
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

func TestSingleFile(t *testing.T) {
	binary := buildBinary(t)

	tests := []struct {
		name     string
		file     string
		expected string
	}{
		{"hello file", "testdata/hello.txt", "Hello, World!\n"},
		{"numbers file", "testdata/numbers.txt", "1\n2\n3\n"},
		{"empty file", "testdata/empty.txt", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.file)
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("command failed: %v", err)
			}

			got := string(out)
			if got != tt.expected {
				t.Errorf("mycat %s\ngot:  %q\nwant: %q", tt.file, got, tt.expected)
			}
		})
	}
}

func TestMultipleFiles(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "testdata/hello.txt", "testdata/numbers.txt")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	expected := "Hello, World!\n1\n2\n3\n"
	got := string(out)
	if got != expected {
		t.Errorf("mycat hello.txt numbers.txt\ngot:  %q\nwant: %q", got, expected)
	}
}

func TestMultipleFilesWithEmpty(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "testdata/hello.txt", "testdata/empty.txt", "testdata/numbers.txt")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	expected := "Hello, World!\n1\n2\n3\n"
	got := string(out)
	if got != expected {
		t.Errorf("mycat hello.txt empty.txt numbers.txt\ngot:  %q\nwant: %q", got, expected)
	}
}

func TestFileNotFound(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "testdata/nonexistent.txt")
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("expected non-zero exit code for missing file")
	}

	got := strings.TrimSpace(string(out))
	if !strings.Contains(got, "error:") || !strings.Contains(got, "no such file or directory") {
		t.Errorf("expected error message about missing file, got: %q", got)
	}
}

func TestMixedValidAndInvalidFiles(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "testdata/hello.txt", "testdata/nonexistent.txt", "testdata/numbers.txt")
	// Capture stdout and stderr separately
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	// Should exit with code 1 because one file failed
	if err == nil {
		t.Fatal("expected non-zero exit code when a file is missing")
	}

	// stdout should contain valid file contents
	gotStdout := stdout.String()
	expectedStdout := "Hello, World!\n1\n2\n3\n"
	if gotStdout != expectedStdout {
		t.Errorf("stdout:\ngot:  %q\nwant: %q", gotStdout, expectedStdout)
	}

	// stderr should contain error message
	gotStderr := stderr.String()
	if !strings.Contains(gotStderr, "error:") || !strings.Contains(gotStderr, "no such file or directory") {
		t.Errorf("stderr should mention missing file, got: %q", gotStderr)
	}
}

func TestExitCodeSuccess(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "testdata/hello.txt")
	err := cmd.Run()
	if err != nil {
		t.Errorf("expected exit code 0 for valid file, got error: %v", err)
	}
}

func TestExitCodeFailure(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "testdata/nonexistent.txt")
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
}

func TestNoArgs(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err == nil {
		t.Fatal("expected non-zero exit code with no args")
	}

	got := strings.TrimSpace(stderr.String())
	expected := "usage: mycat <file1> [file2] ..."
	if got != expected {
		t.Errorf("no args stderr:\ngot:  %q\nwant: %q", got, expected)
	}
}
