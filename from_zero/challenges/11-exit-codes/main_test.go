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

	binary := filepath.Join(dir, "finder")
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

func getExitCode(err error) int {
	if err == nil {
		return 0
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return -1
	}
	return exitErr.ExitCode()
}

func TestFileSearchMatchFound(t *testing.T) {
	binary := buildBinary(t)

	tests := []struct {
		name         string
		pattern      string
		file         string
		expectedOut  string
		expectedExit int
	}{
		{
			name:         "match ap",
			pattern:      "ap",
			file:         "testdata/fruits.txt",
			expectedOut:  "apple\napricot\n",
			expectedExit: 0,
		},
		{
			name:         "match berry",
			pattern:      "berry",
			file:         "testdata/fruits.txt",
			expectedOut:  "blueberry\n",
			expectedExit: 0,
		},
		{
			name:         "match banana exact",
			pattern:      "banana",
			file:         "testdata/fruits.txt",
			expectedOut:  "banana\n",
			expectedExit: 0,
		},
		{
			name:         "match a (multiple lines)",
			pattern:      "a",
			file:         "testdata/fruits.txt",
			expectedOut:  "banana\napricot\navocado\n",
			expectedExit: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.pattern, tt.file)
			out, err := cmd.Output()
			exitCode := getExitCode(err)

			got := string(out)
			if got != tt.expectedOut {
				t.Errorf("finder %q %s\noutput got:  %q\noutput want: %q", tt.pattern, tt.file, got, tt.expectedOut)
			}
			if exitCode != tt.expectedExit {
				t.Errorf("finder %q %s\nexit code got: %d\nexit code want: %d", tt.pattern, tt.file, exitCode, tt.expectedExit)
			}
		})
	}
}

func TestFileSearchNoMatch(t *testing.T) {
	binary := buildBinary(t)

	tests := []struct {
		name    string
		pattern string
		file    string
	}{
		{"no match xyz", "xyz", "testdata/fruits.txt"},
		{"no match grape", "grape", "testdata/fruits.txt"},
		{"no match uppercase Apple", "Apple", "testdata/fruits.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.pattern, tt.file)
			out, err := cmd.Output()
			exitCode := getExitCode(err)

			got := string(out)
			if got != "" {
				t.Errorf("finder %q %s\nexpected no output, got: %q", tt.pattern, tt.file, got)
			}
			if exitCode != 1 {
				t.Errorf("finder %q %s\nexit code got: %d\nexit code want: 1", tt.pattern, tt.file, exitCode)
			}
		})
	}
}

func TestEmptyPatternMatchesAll(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "", "testdata/fruits.txt")
	out, err := cmd.Output()
	exitCode := getExitCode(err)

	got := string(out)
	expected := "apple\nbanana\ncherry\napricot\nblueberry\navocado\n"

	if got != expected {
		t.Errorf("finder '' fruits.txt\noutput got:  %q\noutput want: %q", got, expected)
	}
	if exitCode != 0 {
		t.Errorf("finder '' fruits.txt\nexit code got: %d\nexit code want: 0", exitCode)
	}
}

func TestStdinInput(t *testing.T) {
	binary := buildBinary(t)

	tests := []struct {
		name         string
		pattern      string
		stdin        string
		expectedOut  string
		expectedExit int
	}{
		{
			name:         "match from stdin",
			pattern:      "world",
			stdin:        "hello world\ngoodbye moon\n",
			expectedOut:  "hello world\n",
			expectedExit: 0,
		},
		{
			name:         "no match from stdin",
			pattern:      "xyz",
			stdin:        "hello world\ngoodbye moon\n",
			expectedOut:  "",
			expectedExit: 1,
		},
		{
			name:         "multiple matches from stdin",
			pattern:      "o",
			stdin:        "hello\nworld\nfoo\nbar\n",
			expectedOut:  "hello\nworld\nfoo\n",
			expectedExit: 0,
		},
		{
			name:         "empty pattern matches all stdin",
			pattern:      "",
			stdin:        "line1\nline2\n",
			expectedOut:  "line1\nline2\n",
			expectedExit: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.pattern)
			cmd.Stdin = strings.NewReader(tt.stdin)
			out, err := cmd.Output()
			exitCode := getExitCode(err)

			got := string(out)
			if got != tt.expectedOut {
				t.Errorf("echo ... | finder %q\noutput got:  %q\noutput want: %q", tt.pattern, got, tt.expectedOut)
			}
			if exitCode != tt.expectedExit {
				t.Errorf("echo ... | finder %q\nexit code got: %d\nexit code want: %d", tt.pattern, exitCode, tt.expectedExit)
			}
		})
	}
}

func TestCaseSensitive(t *testing.T) {
	binary := buildBinary(t)

	// "Apple" should NOT match "apple" (case sensitive)
	cmd := exec.Command(binary, "Apple", "testdata/fruits.txt")
	out, err := cmd.Output()
	exitCode := getExitCode(err)

	got := string(out)
	if got != "" {
		t.Errorf("case sensitivity: finder 'Apple' should not match 'apple', got: %q", got)
	}
	if exitCode != 1 {
		t.Errorf("case sensitivity: expected exit code 1, got %d", exitCode)
	}
}

func TestNoArgs(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	err := cmd.Run()

	exitCode := getExitCode(err)
	if exitCode != 2 {
		t.Errorf("no args: expected exit code 2, got %d", exitCode)
	}

	got := strings.TrimSpace(stderr.String())
	expected := "usage: finder <pattern> [file]"
	if got != expected {
		t.Errorf("no args stderr:\ngot:  %q\nwant: %q", got, expected)
	}
}

func TestFileNotFound(t *testing.T) {
	binary := buildBinary(t)

	cmd := exec.Command(binary, "test", "testdata/nonexistent.txt")
	var stderr strings.Builder
	cmd.Stderr = &stderr
	err := cmd.Run()

	exitCode := getExitCode(err)
	if exitCode != 2 {
		t.Errorf("file not found: expected exit code 2, got %d", exitCode)
	}

	got := strings.TrimSpace(stderr.String())
	if !strings.Contains(got, "error:") || !strings.Contains(got, "no such file or directory") {
		t.Errorf("file not found stderr:\ngot: %q\nexpected error about missing file", got)
	}
}

func TestExitCodes(t *testing.T) {
	binary := buildBinary(t)

	tests := []struct {
		name         string
		args         []string
		stdin        string
		expectedExit int
	}{
		{
			name:         "match found - exit 0",
			args:         []string{"apple", "testdata/fruits.txt"},
			expectedExit: 0,
		},
		{
			name:         "no match - exit 1",
			args:         []string{"xyz", "testdata/fruits.txt"},
			expectedExit: 1,
		},
		{
			name:         "no args - exit 2",
			args:         []string{},
			expectedExit: 2,
		},
		{
			name:         "file not found - exit 2",
			args:         []string{"test", "nonexistent.txt"},
			expectedExit: 2,
		},
		{
			name:         "stdin match - exit 0",
			args:         []string{"hello"},
			stdin:        "hello world\n",
			expectedExit: 0,
		},
		{
			name:         "stdin no match - exit 1",
			args:         []string{"xyz"},
			stdin:        "hello world\n",
			expectedExit: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			if tt.stdin != "" {
				cmd.Stdin = strings.NewReader(tt.stdin)
			}
			err := cmd.Run()
			exitCode := getExitCode(err)

			if exitCode != tt.expectedExit {
				t.Errorf("exit code got: %d, want: %d", exitCode, tt.expectedExit)
			}
		})
	}
}
