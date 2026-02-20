package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildBinary compiles the student's solution and returns the path to the binary.
func buildBinary(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	binary := filepath.Join(dir, "strtool")
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

func TestLength(t *testing.T) {
	binary := buildBinary(t)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple word", "Hello", "Length: 5"},
		{"empty string", "", "Length: 0"},
		{"single char", "a", "Length: 1"},
		{"with spaces", "hello world", "Length: 11"},
		{"unicode cafe", "café", "Length: 4"},
		{"unicode emoji", "Go!", "Length: 3"},
		{"unicode japanese", "日本語", "Length: 3"},
		{"numbers", "12345", "Length: 5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, "length", tt.input)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("command failed: %v\n%s", err, out)
			}

			got := strings.TrimSpace(string(out))
			if got != tt.expected {
				t.Errorf("strtool length %q\ngot:  %q\nwant: %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestUpper(t *testing.T) {
	binary := buildBinary(t)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "hello", "Upper: HELLO"},
		{"already upper", "HELLO", "Upper: HELLO"},
		{"mixed case", "Hello World", "Upper: HELLO WORLD"},
		{"with numbers", "abc123", "Upper: ABC123"},
		{"empty string", "", "Upper: "},
		{"unicode", "café", "Upper: CAFÉ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, "upper", tt.input)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("command failed: %v\n%s", err, out)
			}

			got := strings.TrimSpace(string(out))
			if got != tt.expected {
				t.Errorf("strtool upper %q\ngot:  %q\nwant: %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestLower(t *testing.T) {
	binary := buildBinary(t)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "HELLO", "Lower: hello"},
		{"already lower", "hello", "Lower: hello"},
		{"mixed case", "Go Is FUN", "Lower: go is fun"},
		{"with numbers", "ABC123", "Lower: abc123"},
		{"empty string", "", "Lower: "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, "lower", tt.input)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("command failed: %v\n%s", err, out)
			}

			got := strings.TrimSpace(string(out))
			if got != tt.expected {
				t.Errorf("strtool lower %q\ngot:  %q\nwant: %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	binary := buildBinary(t)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "hello", "Reverse: olleh"},
		{"palindrome", "racecar", "Reverse: racecar"},
		{"single char", "a", "Reverse: a"},
		{"empty string", "", "Reverse: "},
		{"with spaces", "hello world", "Reverse: dlrow olleh"},
		{"unicode cafe", "café", "Reverse: éfac"},
		{"unicode japanese", "日本語", "Reverse: 語本日"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, "reverse", tt.input)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("command failed: %v\n%s", err, out)
			}

			got := strings.TrimSpace(string(out))
			if got != tt.expected {
				t.Errorf("strtool reverse %q\ngot:  %q\nwant: %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestContains(t *testing.T) {
	binary := buildBinary(t)

	tests := []struct {
		name      string
		input     string
		substring string
		expected  string
	}{
		{"found", "hello world", "world", "Contains: true"},
		{"not found", "hello world", "xyz", "Contains: false"},
		{"empty substring", "hello", "", "Contains: true"},
		{"empty string", "", "hello", "Contains: false"},
		{"both empty", "", "", "Contains: true"},
		{"exact match", "hello", "hello", "Contains: true"},
		{"case sensitive", "Hello", "hello", "Contains: false"},
		{"unicode", "café latte", "café", "Contains: true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, "contains", tt.input, tt.substring)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("command failed: %v\n%s", err, out)
			}

			got := strings.TrimSpace(string(out))
			if got != tt.expected {
				t.Errorf("strtool contains %q %q\ngot:  %q\nwant: %q", tt.input, tt.substring, got, tt.expected)
			}
		})
	}
}

func TestUnknownCommand(t *testing.T) {
	binary := buildBinary(t)

	tests := []struct {
		name    string
		command string
	}{
		{"explode", "explode"},
		{"split", "split"},
		{"trim", "trim"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.command, "test")
			out, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatal("expected non-zero exit code for unknown command")
			}

			got := strings.TrimSpace(string(out))
			expected := "error: unknown command '" + tt.command + "'"
			if got != expected {
				t.Errorf("strtool %s\ngot:  %q\nwant: %q", tt.command, got, expected)
			}
		})
	}
}

func TestNoArgs(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit code when called with no args")
	}

	got := strings.TrimSpace(string(out))
	expected := "usage: strtool <command> <string> [args...]"
	if got != expected {
		t.Errorf("strtool (no args)\ngot:  %q\nwant: %q", got, expected)
	}
}
