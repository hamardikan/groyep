# Challenge 10: Uppercase Filter

## Prerequisites

- [I/O and Streams](../../concepts/10-io-and-streams.md)
- [Unix Philosophy](../../concepts/11-unix-philosophy.md)

## Objective

Build an uppercase filter called `upper`. It reads text from stdin (or a file) and writes it to stdout with all characters converted to uppercase. This teaches you the Unix pipe concept — building small tools that compose together.

## Why This Matters for grep

grep is designed to work in Unix pipelines: you pipe data into it, it filters, and pipes data out. Understanding how to read from stdin and write to stdout is essential. This is also how `grep -i` (case-insensitive) works internally — by comparing uppercased versions of strings.

## Requirements

### Basic Usage

```bash
# Pipe input
$ echo "hello world" | ./upper
HELLO WORLD

# File argument
$ ./upper file.txt
(file contents in uppercase)

# Multi-line pipe
$ printf "hello\nworld\n" | ./upper
HELLO
WORLD
```

### Input Priority

- If a **filename argument** is provided, read from the file
- If **no arguments**, read from stdin
- If the file doesn't exist: print `error: open <filename>: no such file or directory` to stderr, exit 1

### Details

- Convert all text to uppercase using Unicode-aware conversion
- Preserve line breaks and whitespace structure
- Numbers and special characters pass through unchanged
- Process the entire input (not just line by line — preserve exact whitespace)

### Examples

```bash
$ echo "hello world" | ./upper
HELLO WORLD

$ echo "café au lait" | ./upper
CAFÉ AU LAIT

$ echo "123 abc" | ./upper
123 ABC

$ cat testdata/mixed.txt | ./upper
HELLO WORLD
THIS IS LOWERCASE
THIS IS UPPERCASE
MIXED CASE TEXT
GO IS A GREAT LANGUAGE
123 NUMBERS STAY THE SAME
CAFÉ AU LAIT

$ ./upper testdata/mixed.txt
HELLO WORLD
THIS IS LOWERCASE
THIS IS UPPERCASE
MIXED CASE TEXT
GO IS A GREAT LANGUAGE
123 NUMBERS STAY THE SAME
CAFÉ AU LAIT

$ ./upper nonexistent.txt
error: open nonexistent.txt: no such file or directory
```

## Hints

- Use `strings.ToUpper()` for the conversion
- Use `io.ReadAll(reader)` to read all input at once from any `io.Reader`
- `os.Stdin` implements `io.Reader`, so you can read from it just like a file
- Use `os.Open()` to open a file if a filename argument is given
- The `strings.ToUpper()` function is Unicode-aware — it handles accented characters correctly

## Getting Started

```bash
cd 10-stdin-pipe
go mod init upper
# Create main.go and implement the solution
go build -o upper .
echo "hello world" | ./upper
```

## Tests

```bash
# Run Go tests
go test -v

# Run shell tests
bash test.sh
```
