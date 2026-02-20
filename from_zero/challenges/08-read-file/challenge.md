# Challenge 08: Read File

## Prerequisites

- [I/O and Streams](../../concepts/10-io-and-streams.md)

## Objective

Build your own `cat` command called `mycat`. It reads one or more files and prints their contents to stdout. This teaches you how to open files, read their contents, and handle errors — all fundamental operations for grep.

## Why This Matters for grep

grep's core job is reading files and scanning their contents line by line. Before you can search for patterns, you need to know how to open a file, read its contents, and handle the case where a file doesn't exist. This challenge builds exactly those skills.

## Requirements

### Basic Usage

```bash
# Print a single file
$ ./mycat file.txt
(contents of file.txt)

# Print multiple files concatenated
$ ./mycat file1.txt file2.txt
(contents of file1.txt immediately followed by contents of file2.txt)
```

### Error Handling

- **No arguments**: print `usage: mycat <file1> [file2] ...` to stderr and exit with code 1
- **File not found**: print `error: open <filename>: no such file or directory` to stderr, then continue processing remaining files
- **Exit code**: 0 if all files were read successfully, 1 if any file failed

### Behavior Details

- Files are printed in the order given on the command line
- No separator between files — contents are concatenated directly
- If a file ends without a newline, the next file's content starts on the same line (just like real `cat`)
- Empty files produce no output

### Examples

```bash
$ ./mycat testdata/hello.txt
Hello, World!

$ ./mycat testdata/hello.txt testdata/numbers.txt
Hello, World!
1
2
3

$ ./mycat testdata/nonexistent.txt
error: open nonexistent.txt: no such file or directory

$ ./mycat testdata/hello.txt testdata/nonexistent.txt testdata/numbers.txt
Hello, World!
error: open nonexistent.txt: no such file or directory
1
2
3

$ ./mycat
usage: mycat <file1> [file2] ...
```

## Hints

- Use `os.Open()` to open a file — it returns a `*os.File` and an `error`
- Use `io.Copy(os.Stdout, file)` to efficiently copy file contents to stdout
- Or use `os.ReadFile()` for a simpler approach (reads entire file into memory)
- Always close files after reading: `defer file.Close()`
- Use `fmt.Fprintf(os.Stderr, ...)` to write error messages to stderr
- Use `os.Exit()` to set the exit code

## Getting Started

```bash
cd 08-read-file
go mod init mycat
# Create main.go and implement the solution
go build -o mycat .
./mycat testdata/hello.txt
```

## Tests

```bash
# Run Go tests
go test -v

# Run shell tests
bash test.sh
```
