# Challenge 07: String Checker

## Prerequisites

- [Strings and Encoding](../../concepts/09-strings-and-encoding.md)

## Objective

Build a string checker utility called `strtool` that performs various string operations. This will teach you Go's `strings` package and how to work with runes (Unicode code points) rather than raw bytes.

## Why This Matters for grep

grep operates on text — matching, comparing, and transforming strings. Understanding how Go handles strings, bytes, and runes is essential. When you later implement case-insensitive matching (`grep -i`), you'll use the same `strings` package functions you'll learn here.

## Requirements

Your program takes a command as the first argument, followed by a string (and optionally more arguments depending on the command).

### Commands

| Command | Args | Output Format | Example |
|---------|------|--------------|---------|
| `length` | `<string>` | `Length: N` | `Length: 5` |
| `upper` | `<string>` | `Upper: STRING` | `Upper: HELLO` |
| `lower` | `<string>` | `Lower: string` | `Lower: hello` |
| `reverse` | `<string>` | `Reverse: gnirts` | `Reverse: olleh` |
| `contains` | `<string> <substring>` | `Contains: true` or `Contains: false` | `Contains: true` |

### Edge Cases

- **No arguments**: print `usage: strtool <command> <string> [args...]` and exit with code 1
- **Unknown command**: print `error: unknown command 'xxx'` and exit with code 1
- **Unicode**: `length` must count **runes** (characters), not bytes. For example, `café` has 4 characters but 5 bytes.
- **Reverse**: must correctly reverse Unicode strings (reverse runes, not bytes)

### Examples

```bash
$ ./strtool length "Hello"
Length: 5

$ ./strtool upper "hello world"
Upper: HELLO WORLD

$ ./strtool lower "Go Is FUN"
Lower: go is fun

$ ./strtool reverse "hello"
Reverse: olleh

$ ./strtool contains "hello world" "world"
Contains: true

$ ./strtool contains "hello world" "xyz"
Contains: false

$ ./strtool length "café"
Length: 4

$ ./strtool reverse "café"
Reverse: éfac

$ ./strtool
usage: strtool <command> <string> [args...]

$ ./strtool explode "test"
error: unknown command 'explode'
```

## Hints

- Use `strings.ToUpper()` and `strings.ToLower()` for case conversion
- Use `[]rune(s)` to convert a string to a slice of runes for correct Unicode handling
- `len(string)` returns byte count; `len([]rune(string))` returns character count
- Use `strings.Contains()` for substring checking
- To reverse a string, convert to `[]rune`, swap from both ends, convert back to `string`

## Getting Started

```bash
cd 07-strings
go mod init strtool
# Create main.go and implement the solution
go build -o strtool .
./strtool length "Hello"
```

## Tests

```bash
# Run Go tests
go test -v

# Run shell tests
bash test.sh
```
