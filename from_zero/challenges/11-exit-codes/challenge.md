# Challenge 11: Substring Finder (Mini Grep)

## Prerequisites

- [Processes and Exit Codes](../../concepts/12-processes-and-exit-codes.md)

## Objective

Build a substring finder called `finder`. It searches for a pattern in lines from a file or stdin and prints matching lines. This is a **mini grep** — the precursor to everything you'll build in Phase 4.

## Why This Matters for grep

This IS grep (in miniature). The real `grep` searches for patterns in text and uses exit codes to communicate results. Exit code 0 means "matches found," exit code 1 means "no matches," and exit code 2 means "error." You're building the exact same contract here.

## Requirements

### Basic Usage

```bash
# Search in a file
$ ./finder "pattern" file.txt

# Search from stdin
$ echo "some text" | ./finder "pattern"
```

### Matching Rules

- Print each line that **contains** the pattern (case-sensitive)
- An **empty pattern** (`""`) matches every line (just like `grep ""`)
- Matching is done as a simple substring search (not regex)

### Exit Codes

| Condition | Exit Code |
|-----------|-----------|
| At least one match found | 0 |
| No matches found | 1 |
| Error (no args, bad usage) | 2 |

This matches the real `grep` exit code convention.

### Error Handling

- **No arguments**: print `usage: finder <pattern> [file]` to stderr, exit 2
- **File not found**: print `error: open <filename>: no such file or directory` to stderr, exit 2

### Examples

```bash
$ ./finder "ap" testdata/fruits.txt
apple
apricot

$ echo $?
0

$ ./finder "berry" testdata/fruits.txt
blueberry

$ ./finder "xyz" testdata/fruits.txt
$ echo $?
1

$ ./finder "" testdata/fruits.txt
apple
banana
cherry
apricot
blueberry
avocado

$ echo "hello world" | ./finder "world"
hello world

$ echo "hello world" | ./finder "xyz"
$ echo $?
1

$ ./finder
usage: finder <pattern> [file]
$ echo $?
2

$ ./finder "test" nonexistent.txt
error: open nonexistent.txt: no such file or directory
$ echo $?
2
```

## Hints

- Use `strings.Contains(line, pattern)` for substring matching
- An empty string is contained in every string: `strings.Contains("anything", "")` returns `true`
- Use `os.Exit(0)`, `os.Exit(1)`, or `os.Exit(2)` to set exit codes
- Remember: exit code indicates whether matches were found, not whether the program ran successfully
- Use `bufio.NewScanner()` to read line by line from a file or stdin

## Getting Started

```bash
cd 11-exit-codes
go mod init finder
# Create main.go and implement the solution
go build -o finder .
./finder "ap" testdata/fruits.txt
```

## Tests

```bash
# Run Go tests
go test -v

# Run shell tests
bash test.sh
```

## Reflection

After completing this challenge, you've built a tool that:
- Reads from files or stdin
- Searches for patterns in text
- Reports results via exit codes

This is the foundation of grep. In the upcoming Phase 4 challenges, you'll extend this pattern with regex support, line numbers, recursive search, and more.
