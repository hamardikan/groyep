# Challenge 12: CLI Args — Line Search Tool

## Prerequisites

- [CLI Design](../../concepts/13-cli-design.md)
- Completed: [Challenge 11 — Exit Codes](../11-exit-codes/challenge.md)

## Overview

It's time to bring everything together. You've learned to read files, write to stdout, handle stdin, and use exit codes. Now you'll build a real line search tool — a program that takes a pattern and a filename, and prints every line that contains the pattern.

This is the beginning of your grep clone. The output format you'll implement here matches what real `grep` produces.

---

## Task

Build a command-line tool called `search` that finds lines matching a pattern in a file.

### Setup

1. `go mod init search`
2. Create `main.go`

### Usage

```
./search <pattern> <filename>
```

### Requirements

1. **Substring matching** — a line matches if it contains the pattern anywhere (case-sensitive)
2. **Output format** — prefix each matching line with the filename and a colon: `filename:matching line`
3. **Multiple matches** — print all matching lines, in the order they appear in the file
4. **No match** — if no lines match, print nothing and exit with code `1`
5. **Missing arguments** — if fewer than 2 arguments are given, print `usage: search <pattern> <filename>` to stderr and exit with code `2`
6. **File not found** — if the file cannot be opened, print `error: cannot open '<filename>'` to stderr and exit with code `2`

### Examples

```bash
$ ./search "Priest" testdata/rockbands.txt
testdata/rockbands.txt:Judas Priest

$ ./search "Bad" testdata/rockbands.txt
testdata/rockbands.txt:Bad English
testdata/rockbands.txt:Bad Company

$ ./search "White" testdata/rockbands.txt
testdata/rockbands.txt:Whitesnake
testdata/rockbands.txt:Great White
testdata/rockbands.txt:White Lion
testdata/rockbands.txt:Whitecross

$ ./search "zzz" testdata/rockbands.txt
# (no output, exit code 1)

$ ./search
usage: search <pattern> <filename>
# (exit code 2)

$ ./search "hello" nonexistent.txt
error: cannot open 'nonexistent.txt'
# (exit code 2)
```

### Verify Your Solution

```bash
bash test.sh
go test -v
```

---

## Concepts

### Reading Command-Line Arguments

You've used `os.Args` before. Recall:

```go
// os.Args[0] is the program name
// os.Args[1] is the first argument
// os.Args[2] is the second argument
args := os.Args[1:]  // everything after program name
```

### Checking for Substring

The `strings` package provides `strings.Contains`:

```go
import "strings"

if strings.Contains(line, pattern) {
    // line matches
}
```

### Opening a File

```go
file, err := os.Open(filename)
if err != nil {
    // handle error
}
defer file.Close()
```

### Scanning Lines

```go
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    // process line
}
```

### Writing to Stderr

```go
fmt.Fprintln(os.Stderr, "error message")
```

For more depth, see [CLI Design](../../concepts/13-cli-design.md).

---

## Hints

<details>
<summary>Hint 1: Structuring your program</summary>

```go
func main() {
    // 1. Check args length
    // 2. Extract pattern and filename
    // 3. Open file (handle error)
    // 4. Scan lines, check for match
    // 5. Track whether any match was found
    // 6. Exit with appropriate code
}
```

</details>

<details>
<summary>Hint 2: Tracking matches</summary>

Use a boolean to track whether you found any matches:

```go
found := false
for scanner.Scan() {
    line := scanner.Text()
    if strings.Contains(line, pattern) {
        fmt.Printf("%s:%s\n", filename, line)
        found = true
    }
}
if !found {
    os.Exit(1)
}
```

</details>

<details>
<summary>Hint 3: Error messages go to stderr</summary>

```go
fmt.Fprintln(os.Stderr, "usage: search <pattern> <filename>")
os.Exit(2)
```

Using `os.Stderr` instead of default stdout ensures error messages don't mix with normal output. This is standard Unix convention.

</details>

---

## Further Reading

- [Go by Example: Command-Line Arguments](https://gobyexample.com/command-line-arguments)
- [strings.Contains documentation](https://pkg.go.dev/strings#Contains)
- [bufio.Scanner documentation](https://pkg.go.dev/bufio#Scanner)
- [os.Open documentation](https://pkg.go.dev/os#Open)
