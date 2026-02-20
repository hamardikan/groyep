# Challenge 09: Number Lines

## Prerequisites

- [I/O and Streams](../../concepts/10-io-and-streams.md)

## Objective

Build a `nl` (number lines) command called `mynl`. It reads a file (or stdin) and prints each line prefixed with its line number. This teaches you formatted output and line-by-line processing — exactly how grep processes files.

## Why This Matters for grep

grep reads files line by line and optionally prints line numbers (`grep -n`). This challenge teaches you the exact same pattern: read line by line, format output with line numbers, and write to stdout. The line-numbering format you build here is nearly identical to `grep -n` output.

## Requirements

### Output Format

Each line is prefixed with a right-aligned line number, 6 characters wide, followed by a tab character, then the line content.

```
     1	Hello, World!
     2	This is line two.
     3	
     4	This is line four.
```

Note: The space between the number and content is a **tab character** (`\t`), not spaces.

### Input Sources

- **File argument**: `./mynl file.txt` — number lines from the file
- **Stdin**: `cat file.txt | ./mynl` — number lines from piped input (when no file argument given)

### Error Handling

- **File not found**: print `error: open <filename>: no such file or directory` to stderr, exit with code 1
- **No args and no pipe**: read from stdin (will wait for input, like real `nl`)

### Details

- Empty lines still get line numbers
- Line numbers start at 1
- The number field is exactly 6 characters wide, right-aligned, padded with spaces

### Examples

```bash
$ ./mynl testdata/sample.txt
     1	Hello, World!
     2	This is line two.
     3	
     4	This is line four.
     5	The end.

$ echo -e "first\nsecond\nthird" | ./mynl
     1	first
     2	second
     3	third

$ ./mynl nonexistent.txt
error: open nonexistent.txt: no such file or directory
```

## Hints

- Use `bufio.NewScanner(reader)` to read line by line from any `io.Reader`
- Use `fmt.Printf("%6d\t%s\n", lineNum, line)` for formatted output
- `os.Stdin` is an `io.Reader` — you can pass it to `bufio.NewScanner` just like a file
- The `%6d` format specifier means: print an integer, right-aligned, in a field 6 characters wide

## Getting Started

```bash
cd 09-write-stdout
go mod init mynl
# Create main.go and implement the solution
go build -o mynl .
./mynl testdata/sample.txt
```

## Tests

```bash
# Run Go tests
go test -v

# Run shell tests
bash test.sh
```
