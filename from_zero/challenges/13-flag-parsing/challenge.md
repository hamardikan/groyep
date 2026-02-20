# Challenge 13: Flag Parsing

## Prerequisites

- [CLI Design](../../concepts/13-cli-design.md)
- Completed: [Challenge 12 — CLI Args](../12-cli-args/challenge.md)

## Overview

Your search tool works, but real command-line tools have **flags** — optional switches that modify behavior. Think `ls -l`, `grep -n`, or `git commit -m`. Go's standard library includes the `flag` package, purpose-built for parsing command-line flags.

In this challenge, you'll add `-n` (line numbers) and `-c` (count) flags to your search tool.

---

## Task

Extend your search tool to support flags.

### Setup

1. `go mod init search`
2. Create `main.go` (you can start from your challenge 12 solution)

### Usage

```
./search [flags] <pattern> <filename>
```

### Flags

| Flag | Description |
|------|-------------|
| `-n` | Prefix each matching line with its 1-based line number |
| `-c` | Print only the count of matching lines (not the lines themselves) |
| `-h` | Print usage/help message |

### Output Format

**Default (no flags):** Same as challenge 12.
```
filename:matching line
```

**With `-n`:** Add the line number between filename and content.
```
filename:linenum:matching line
```

**With `-c`:** Print only the count.
```
filename:count
```

**With `-n` and `-c` together:** Count wins. Print only the count (same as `-c` alone).

### Examples

```bash
$ ./search "Bad" testdata/rockbands.txt
testdata/rockbands.txt:Bad English
testdata/rockbands.txt:Bad Company

$ ./search -n "Bad" testdata/rockbands.txt
testdata/rockbands.txt:19:Bad English
testdata/rockbands.txt:25:Bad Company

$ ./search -c "Bad" testdata/rockbands.txt
testdata/rockbands.txt:2

$ ./search -n -c "Bad" testdata/rockbands.txt
testdata/rockbands.txt:2

$ ./search -h
Usage of search:
  ./search [flags] <pattern> <filename>
Flags:
  -c    print only count of matching lines
  -n    prefix output with line numbers
```

### Error Behavior

Same as challenge 12:
- Missing args after flags: `usage: search <pattern> <filename>` to stderr, exit code 2
- File not found: `error: cannot open '<filename>'` to stderr, exit code 2
- No match: exit code 1, no output

### Verify Your Solution

```bash
bash test.sh
go test -v
```

---

## Concepts

### The `flag` Package

Go's `flag` package handles command-line flag parsing:

```go
import "flag"

func main() {
    // Define flags — these return *bool or *string pointers
    showNumbers := flag.Bool("n", false, "prefix output with line numbers")
    countOnly := flag.Bool("c", false, "print only count of matching lines")

    // Parse flags — this MUST be called before accessing flag values
    flag.Parse()

    // After Parse(), flag.Args() returns the remaining non-flag arguments
    args := flag.Args()
    // args[0] is pattern, args[1] is filename
}
```

**Key insight:** `flag.Parse()` consumes all flags (arguments starting with `-`) and leaves the positional arguments accessible via `flag.Args()`. This is why you call `flag.Args()` instead of `os.Args[1:]` — the flag package has already consumed the flags.

### Checking Flag Values

Since `flag.Bool` returns a `*bool` (pointer to bool), you dereference it with `*`:

```go
if *showNumbers {
    // -n was passed
}

if *countOnly {
    // -c was passed
}
```

### Line Numbers

When scanning a file, track the line number with a counter:

```go
lineNum := 0
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    lineNum++
    line := scanner.Text()
    // lineNum is now the 1-based line number
}
```

For more depth, see [CLI Design](../../concepts/13-cli-design.md).

---

## Hints

<details>
<summary>Hint 1: Checking positional args after flag parsing</summary>

After `flag.Parse()`, use `flag.Args()` instead of `os.Args`:

```go
flag.Parse()
args := flag.Args()
if len(args) < 2 {
    fmt.Fprintln(os.Stderr, "usage: search <pattern> <filename>")
    os.Exit(2)
}
pattern := args[0]
filename := args[1]
```

</details>

<details>
<summary>Hint 2: Formatting output based on flags</summary>

```go
if *countOnly {
    fmt.Printf("%s:%d\n", filename, matchCount)
} else if *showNumbers {
    fmt.Printf("%s:%d:%s\n", filename, lineNum, line)
} else {
    fmt.Printf("%s:%s\n", filename, line)
}
```

When both `-n` and `-c` are set, count wins — so check `*countOnly` first.

</details>

<details>
<summary>Hint 3: Implementing -c (count mode)</summary>

In count mode, don't print individual lines. Instead, count matches in the loop, then print the total at the end:

```go
matchCount := 0
for scanner.Scan() {
    if strings.Contains(scanner.Text(), pattern) {
        matchCount++
    }
}
if matchCount > 0 {
    fmt.Printf("%s:%d\n", filename, matchCount)
} else {
    os.Exit(1)
}
```

</details>

---

## Further Reading

- [Go by Example: Command-Line Flags](https://gobyexample.com/command-line-flags)
- [flag package documentation](https://pkg.go.dev/flag)
