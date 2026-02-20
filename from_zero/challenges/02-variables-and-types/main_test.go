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
	dir, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory: " + err.Error())
	}
	binaryPath = filepath.Join(dir, "variables_test_bin")

	build := exec.Command("go", "build", "-o", binaryPath, ".")
	build.Dir = dir
	if output, err := build.CombinedOutput(); err != nil {
		panic("failed to build binary: " + err.Error() + "\n" + string(output))
	}

	code := m.Run()

	os.Remove(binaryPath)
	os.Exit(code)
}

func runBinary(args ...string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestFullOutput(t *testing.T) {
	output, err := runBinary()
	if err != nil {
		t.Fatalf("program exited with error: %v\nOutput: %s", err, output)
	}

	expected := "Name: Gopher\nAge: 10\nHeight: 1.75\nIsAwesome: true\n"
	if output != expected {
		t.Errorf("output mismatch\n  expected: %q\n  actual:   %q", expected, output)
	}
}

func TestNameLine(t *testing.T) {
	output, err := runBinary()
	if err != nil {
		t.Fatalf("program exited with error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) < 1 {
		t.Fatal("expected at least 1 line of output")
	}
	if lines[0] != "Name: Gopher" {
		t.Errorf("line 1 mismatch\n  expected: %q\n  actual:   %q", "Name: Gopher", lines[0])
	}
}

func TestAgeLine(t *testing.T) {
	output, err := runBinary()
	if err != nil {
		t.Fatalf("program exited with error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) < 2 {
		t.Fatal("expected at least 2 lines of output")
	}
	if lines[1] != "Age: 10" {
		t.Errorf("line 2 mismatch\n  expected: %q\n  actual:   %q", "Age: 10", lines[1])
	}
}

func TestHeightLine(t *testing.T) {
	output, err := runBinary()
	if err != nil {
		t.Fatalf("program exited with error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) < 3 {
		t.Fatal("expected at least 3 lines of output")
	}
	if lines[2] != "Height: 1.75" {
		t.Errorf("line 3 mismatch\n  expected: %q\n  actual:   %q\n  hint: use fmt.Printf(\"Height: %%.2f\\n\", height) for 2 decimal places", "Height: 1.75", lines[2])
	}
}

func TestIsAwesomeLine(t *testing.T) {
	output, err := runBinary()
	if err != nil {
		t.Fatalf("program exited with error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) < 4 {
		t.Fatal("expected at least 4 lines of output")
	}
	if lines[3] != "IsAwesome: true" {
		t.Errorf("line 4 mismatch\n  expected: %q\n  actual:   %q", "IsAwesome: true", lines[3])
	}
}

func TestExactlyFourLines(t *testing.T) {
	output, err := runBinary()
	if err != nil {
		t.Fatalf("program exited with error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) != 4 {
		t.Errorf("expected exactly 4 lines, got %d\n  output: %q", len(lines), output)
	}
}
