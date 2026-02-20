# Challenge 15: Packages — Project Structure

## Prerequisites

- [Modules and Packages](../../concepts/03-modules-and-packages.md)
- [Testing Philosophy](../../concepts/15-testing-philosophy.md)
- Completed: [Challenge 14 — Structs and Methods](../14-structs-methods/challenge.md)

## Overview

Your search tool works and is well-structured with a `Matcher` struct. But everything still lives in a single file in a single package. Real Go projects split code into multiple packages — separating the CLI layer from the core logic.

This challenge teaches **project structure**. The behavior is identical to challenge 14. The learning is all about organization.

---

## Task

Restructure your search tool into a proper multi-package Go project.

### Target Layout

```
search/
  go.mod
  cmd/
    search/
      main.go           # Entry point — thin: parses args, calls matcher
  internal/
    matcher/
      matcher.go        # Core matching logic (Matcher struct + methods)
      matcher_test.go   # Unit tests for the matcher package
  testdata/
    rockbands.txt
    fruits.txt
```

### Setup

1. `go mod init search`
2. Create the directory structure above
3. Move your `Matcher` struct and methods into `internal/matcher/`
4. Create a thin `main.go` in `cmd/search/` that imports and uses the matcher
5. Copy `matcher_test.go` from this challenge directory into `internal/matcher/`

### The `matcher_test.go` File

This challenge includes a `matcher_test.go` file alongside this `challenge.md`. This file tests the matcher package directly — not through the binary. **You must copy it to `internal/matcher/matcher_test.go`** once you've created the package.

The test file defines the API your matcher package must implement:

```go
// matcher.Match(pattern, line) returns true if line contains pattern
matcher.Match("Bad", "Bad English")  // true
matcher.Match("Bad", "Aerosmith")    // false

// matcher.SearchLines(pattern, lines) returns all matching lines
matcher.SearchLines("ap", []string{"apple", "banana", "apricot"})
// returns: ["apple", "apricot"]
```

These are **exported functions** — your package must provide them with exactly these signatures.

### Binary Behavior

The compiled binary must behave identically to challenge 14:

```bash
# Build from the cmd/search directory
go build -o search ./cmd/search

# All previous behavior works
./search "Bad" testdata/rockbands.txt
./search -n "Bad" testdata/rockbands.txt
./search -c "Bad" testdata/rockbands.txt
./search "Bad" testdata/rockbands.txt testdata/fruits.txt
```

### Verify Your Solution

```bash
# Run unit tests for the matcher package
go test -v ./internal/matcher/

# Run the binary tests
bash test.sh
```

---

## Concepts

### Why Multiple Packages?

A single-file program is fine for small tools. But as code grows, splitting into packages gives you:

1. **Testability** — you can unit test the matching logic without building and running a binary
2. **Reusability** — other programs could import your matcher package
3. **Clarity** — the CLI concerns (flag parsing, exit codes) are separate from the core logic (string matching, file scanning)

### The `cmd/` Convention

Go projects typically put executables in `cmd/<name>/main.go`:

```
cmd/
  search/
    main.go     # package main — the entry point
```

This is a convention, not a requirement. It keeps the root clean and allows multiple commands in one repo.

### The `internal/` Directory

The `internal/` directory is special in Go. Packages under `internal/` can only be imported by code in the parent of `internal/`. This means:

- `search/cmd/search/main.go` can import `search/internal/matcher` (same module)
- External projects **cannot** import `search/internal/matcher`

This is Go's way of marking packages as "for internal use only" — you get a compilation error if you try to import them from outside.

### Importing Your Package

In `cmd/search/main.go`:

```go
package main

import (
    "search/internal/matcher"
)

func main() {
    // Use matcher.Match(), matcher.SearchLines(), etc.
}
```

The import path is `<module-name>/internal/matcher` where `<module-name>` is whatever you put in `go.mod`.

### Thin `main.go`

The entry point should be thin — its job is:
1. Parse command-line flags
2. Validate arguments
3. Call into the matcher package
4. Handle exit codes

All the "interesting" logic lives in the matcher package.

### Exported vs. Unexported

In Go, names that start with a capital letter are **exported** (visible outside the package):

```go
package matcher

// Exported — usable by other packages
func Match(pattern, line string) bool { ... }

// Unexported — only usable within this package
func formatOutput(filename, line string) string { ... }
```

For more depth, see [Modules and Packages](../../concepts/03-modules-and-packages.md) and [Testing Philosophy](../../concepts/15-testing-philosophy.md).

---

## Hints

<details>
<summary>Hint 1: The matcher package API</summary>

Your `internal/matcher/matcher.go` should export at least:

```go
package matcher

// Match returns true if line contains pattern (case-sensitive substring match).
func Match(pattern, line string) bool {
    // ...
}

// SearchLines returns all lines that match the pattern.
func SearchLines(pattern string, lines []string) []string {
    // ...
}
```

You'll likely also want your `Matcher` struct here, with methods like `SearchFile`.

</details>

<details>
<summary>Hint 2: Importing in main.go</summary>

```go
package main

import (
    "flag"
    "fmt"
    "os"

    "search/internal/matcher"
)

func main() {
    showNumbers := flag.Bool("n", false, "prefix with line numbers")
    countOnly := flag.Bool("c", false, "count matches only")
    flag.Parse()

    args := flag.Args()
    if len(args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: search <pattern> <filename>")
        os.Exit(2)
    }

    m := &matcher.Matcher{
        Pattern:     args[0],
        ShowNumbers: *showNumbers,
        CountOnly:   *countOnly,
    }

    // loop over args[1:] files...
}
```

</details>

<details>
<summary>Hint 3: Building the binary</summary>

Since `main.go` is in `cmd/search/`, you build with:

```bash
go build -o search ./cmd/search
```

Or from the `cmd/search` directory:

```bash
cd cmd/search && go build -o ../../search .
```

</details>

<details>
<summary>Hint 4: Running package tests</summary>

```bash
# Run matcher unit tests
go test -v ./internal/matcher/

# Run all tests in the module
go test -v ./...
```

</details>

---

## Further Reading

- [Go Blog: Organizing Go Code](https://go.dev/blog/organizing-go-code)
- [Standard Go Project Layout (community)](https://github.com/golang-standards/project-layout)
- [Go Modules Reference](https://go.dev/ref/mod)
- [Go by Example: Testing](https://gobyexample.com/testing)
