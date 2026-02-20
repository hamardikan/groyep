# Challenge 14: Structs and Methods — Refactoring with Types

## Prerequisites

- [Structs and Interfaces](../../concepts/14-structs-and-interfaces.md)
- Completed: [Challenge 13 — Flag Parsing](../13-flag-parsing/challenge.md)

## Overview

Your search tool works. It handles flags, formats output, and manages errors. But all the logic lives in one function. As programs grow, this becomes hard to maintain.

This challenge introduces **refactoring** — changing the internal structure of code without changing its external behavior. You'll reorganize your search tool to use a **struct** with methods, which groups related data and behavior together.

You'll also add a new capability: searching **multiple files** at once.

---

## What is Refactoring?

Refactoring means restructuring existing code without changing what it does. The tests from challenge 13 still pass — the user sees identical behavior. But the code is better organized, easier to read, and easier to extend.

The key principle: **make the change easy, then make the easy change.** By refactoring first, adding multi-file support becomes straightforward.

---

## Task

Refactor your search tool to use a `Matcher` struct, then add multi-file support.

### Setup

1. `go mod init search`
2. Create `main.go` (you can start from your challenge 13 solution)

### Usage

```
./search [flags] <pattern> <filename> [filename2] ...
```

### Requirements

All behavior from challenge 13 must still work:
- Default output: `filename:matching line`
- `-n` flag: `filename:linenum:matching line`
- `-c` flag: `filename:count`
- `-n -c` combined: count wins
- Error handling: same exit codes and messages

**New: Multiple files.** When more than one file is given, search all files in order.

```bash
$ ./search "Bad" testdata/rockbands.txt testdata/fruits.txt
testdata/rockbands.txt:Bad English
testdata/rockbands.txt:Bad Company

$ ./search -n "Tora" testdata/rockbands.txt testdata/fruits.txt
testdata/rockbands.txt:47:Tora Tora
testdata/rockbands.txt:71:Tora Tora

$ ./search -c "Bad" testdata/rockbands.txt testdata/fruits.txt
testdata/rockbands.txt:2
```

**Multi-file error handling:**
- If one file cannot be opened, print the error for that file to stderr and continue with remaining files
- Exit code: `2` if any file had an error, `1` if no matches found in any file, `0` if at least one match

### Verify Your Solution

```bash
bash test.sh
go test -v
```

---

## Concepts

### Structs

A struct groups related data together:

```go
type Matcher struct {
    Pattern    string
    ShowNumbers bool
    CountOnly   bool
}
```

### Methods

Methods are functions attached to a struct:

```go
func (m *Matcher) Match(line string) bool {
    return strings.Contains(line, m.Pattern)
}

func (m *Matcher) SearchFile(filename string) ([]string, error) {
    // open file, scan lines, return matches
}
```

The `(m *Matcher)` part is the **receiver** — it's how Go attaches functions to types. Inside the method, `m` refers to the struct instance.

### Using Your Struct

```go
func main() {
    // parse flags...

    m := &Matcher{
        Pattern:     pattern,
        ShowNumbers: *showNumbers,
        CountOnly:   *countOnly,
    }

    // use m.SearchFile() for each file
}
```

### Why Structs?

Before structs, you pass many parameters to functions:

```go
func searchFile(pattern string, filename string, showNumbers bool, countOnly bool) { ... }
```

With structs, the configuration is bundled:

```go
func (m *Matcher) SearchFile(filename string) { ... }
```

This becomes critical as you add more features (case-insensitive matching, regex, invert match, etc.).

For more depth, see [Structs and Interfaces](../../concepts/14-structs-and-interfaces.md).

---

## Hints

<details>
<summary>Hint 1: Designing the Matcher struct</summary>

```go
type Matcher struct {
    Pattern     string
    ShowNumbers bool
    CountOnly   bool
}
```

Keep it simple. You can always add fields later.

</details>

<details>
<summary>Hint 2: The SearchFile method</summary>

```go
func (m *Matcher) SearchFile(filename string) (bool, error) {
    file, err := os.Open(filename)
    if err != nil {
        return false, err
    }
    defer file.Close()

    found := false
    lineNum := 0
    matchCount := 0
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        lineNum++
        line := scanner.Text()
        if m.Match(line) {
            found = true
            matchCount++
            if !m.CountOnly {
                // print the line (with or without number)
            }
        }
    }

    if m.CountOnly && matchCount > 0 {
        fmt.Printf("%s:%d\n", filename, matchCount)
    }

    return found, nil
}
```

</details>

<details>
<summary>Hint 3: Looping over multiple files</summary>

```go
files := flag.Args()[1:] // all args after pattern
anyMatch := false
anyError := false

for _, filename := range files {
    found, err := m.SearchFile(filename)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: cannot open '%s'\n", filename)
        anyError = true
        continue
    }
    if found {
        anyMatch = true
    }
}

if anyError {
    os.Exit(2)
} else if !anyMatch {
    os.Exit(1)
}
```

</details>

---

## Further Reading

- [Go by Example: Structs](https://gobyexample.com/structs)
- [Go by Example: Methods](https://gobyexample.com/methods)
- [Go by Example: Pointers](https://gobyexample.com/pointers)
- [Effective Go: Methods](https://go.dev/doc/effective_go#methods)
