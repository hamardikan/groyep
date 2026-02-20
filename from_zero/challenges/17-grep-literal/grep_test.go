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

func TestLiteralMatching(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		file     string
		expected []string
		exitCode int
	}{
		{
			name:     "single char J",
			pattern:  "J",
			file:     "testdata/rockbands.txt",
			expected: []string{"Judas Priest", "Bon Jovi", "Junkyard"},
			exitCode: 0,
		},
		{
			name:     "multi-word Bad",
			pattern:  "Bad",
			file:     "testdata/rockbands.txt",
			expected: []string{"Bad English", "Bad Company"},
			exitCode: 0,
		},
		{
			name:     "exact band name",
			pattern:  "Nirvana",
			file:     "testdata/rockbands.txt",
			expected: []string{"Nirvana"},
			exitCode: 0,
		},
		{
			name:     "special chars AC/DC",
			pattern:  "AC/DC",
			file:     "testdata/rockbands.txt",
			expected: []string{"AC/DC"},
			exitCode: 0,
		},
		{
			name:     "no match",
			pattern:  "ZZZZNOTHERE",
			file:     "testdata/rockbands.txt",
			expected: []string{},
			exitCode: 1,
		},
		{
			name:     "match in BFS1985",
			pattern:  "Madonna",
			file:     "testdata/test-subdir/BFS1985.txt",
			expected: []string{"Springsteen, Madonna"},
			exitCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, _, exitCode := runGrep(tt.pattern, tt.file)
			lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")

			if exitCode != tt.exitCode {
				t.Errorf("expected exit code %d, got %d", tt.exitCode, exitCode)
			}

			for _, exp := range tt.expected {
				found := false
				for _, line := range lines {
					if strings.Contains(line, exp) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected output to contain %q\nGot:\n%s", exp, stdout)
				}
			}
		})
	}
}

func TestStdinPipe(t *testing.T) {
	tests := []struct {
		name     string
		stdin    string
		pattern  string
		expected []string
		exitCode int
	}{
		{
			name:     "pipe match",
			stdin:    "hello world\nfoo bar\nhello again\n",
			pattern:  "hello",
			expected: []string{"hello world", "hello again"},
			exitCode: 0,
		},
		{
			name:     "pipe no match",
			stdin:    "hello world\nfoo bar\n",
			pattern:  "zzz",
			expected: []string{},
			exitCode: 1,
		},
		{
			name:     "pipe empty pattern",
			stdin:    "line one\nline two\n",
			pattern:  "",
			expected: []string{"line one", "line two"},
			exitCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, exitCode := runGrepWithStdin(tt.stdin, tt.pattern)

			if exitCode != tt.exitCode {
				t.Errorf("expected exit code %d, got %d", tt.exitCode, exitCode)
			}

			for _, exp := range tt.expected {
				if !strings.Contains(stdout, exp) {
					t.Errorf("expected output to contain %q\nGot:\n%s", exp, stdout)
				}
			}
		})
	}
}

func TestExitCodes(t *testing.T) {
	// No args
	_, _, exitCode := runGrep()
	if exitCode != 2 {
		t.Errorf("no args: expected exit 2, got %d", exitCode)
	}

	// File not found
	_, _, exitCode = runGrep("pattern", "nonexistent.txt")
	if exitCode != 2 {
		t.Errorf("file not found: expected exit 2, got %d", exitCode)
	}

	// Match found
	_, _, exitCode = runGrep("Nirvana", "testdata/rockbands.txt")
	if exitCode != 0 {
		t.Errorf("match found: expected exit 0, got %d", exitCode)
	}

	// No match
	_, _, exitCode = runGrep("ZZZZZZ", "testdata/rockbands.txt")
	if exitCode != 1 {
		t.Errorf("no match: expected exit 1, got %d", exitCode)
	}
}
