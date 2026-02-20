// Run: go test -v
//
// This test builds your program and checks that it prints the correct output.
// You do NOT need to edit this file — just create your main.go and run the tests.

package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Build the binary before running tests
	dir, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory: " + err.Error())
	}
	binaryPath = filepath.Join(dir, "hello_test_bin")

	build := exec.Command("go", "build", "-o", binaryPath, ".")
	build.Dir = dir
	if output, err := build.CombinedOutput(); err != nil {
		panic("failed to build binary: " + err.Error() + "\n" + string(output))
	}

	code := m.Run()

	// Cleanup
	os.Remove(binaryPath)
	os.Exit(code)
}

func runBinary(args ...string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestHelloOutput(t *testing.T) {
	output, err := runBinary()
	if err != nil {
		t.Fatalf("program exited with error: %v\nOutput: %s", err, output)
	}

	expected := "Hello, World!\nHello, Go!\n"
	if output != expected {
		t.Errorf("output mismatch\n  expected: %q\n  actual:   %q", expected, output)
	}
}

func TestHelloFirstLine(t *testing.T) {
	output, err := runBinary()
	if err != nil {
		t.Fatalf("program exited with error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) < 1 {
		t.Fatal("expected at least 1 line of output, got none")
	}
	if lines[0] != "Hello, World!" {
		t.Errorf("first line mismatch\n  expected: %q\n  actual:   %q", "Hello, World!", lines[0])
	}
}

func TestHelloSecondLine(t *testing.T) {
	output, err := runBinary()
	if err != nil {
		t.Fatalf("program exited with error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected 2 lines of output, got %d", len(lines))
	}
	if lines[1] != "Hello, Go!" {
		t.Errorf("second line mismatch\n  expected: %q\n  actual:   %q", "Hello, Go!", lines[1])
	}
}

func TestHelloLineCount(t *testing.T) {
	output, err := runBinary()
	if err != nil {
		t.Fatalf("program exited with error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("expected exactly 2 lines, got %d\n  output: %q", len(lines), output)
	}
}
