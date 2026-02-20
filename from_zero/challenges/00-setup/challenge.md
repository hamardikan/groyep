# Challenge 00: Setup Your Go Environment

## Prerequisites

Before starting this challenge, make sure you've read:

- [How Computers Run Code](../../concepts/01-how-computers-run-code.md)
- [What is Go](../../concepts/02-what-is-go.md)
- [Modules and Packages](../../concepts/03-modules-and-packages.md)

## Overview

Before you write a single line of Go, you need a working environment. This challenge walks you through installing Go, understanding how Go projects are organized, and building your first program.

**Keep this challenge as your reference card.** Every future challenge assumes you know these commands.

---

## Concepts

### Checking if Go Is Installed

Open your terminal and run:

```bash
go version
```

You should see something like `go version go1.22.0 darwin/arm64` (your version and platform will differ). If you get "command not found", you need to install Go.

### Installing Go

Download and install from: **https://go.dev/dl/**

Choose the installer for your operating system. After installation, restart your terminal and run `go version` again to confirm it works.

### GOPATH vs Modules

Go has two systems for organizing code:

- **GOPATH** (legacy): All Go code lived in a single workspace directory (`~/go`). You don't need this anymore.
- **Go Modules** (modern, use this): Each project has its own `go.mod` file that declares the module name and dependencies. This is what you'll use.

### Initializing a Module

Every Go project starts with:

```bash
go mod init <module-name>
```

For these challenges, use a simple name like:

```bash
go mod init setup
```

This creates a `go.mod` file — the identity card of your project. It tells Go "this directory is a module."

### Creating Your First Go File

Create a file called `main.go`. Every executable Go program needs:

1. **`package main`** — declares this file belongs to the `main` package (the one that produces an executable)
2. **`func main()`** — the entry point; execution starts here

```go
package main

import "fmt"

func main() {
    fmt.Println("something")
}
```

### Running Your Program

You have two options:

**Run directly** (compile + run in one step, no binary saved):
```bash
go run .
```

**Build then run** (creates a binary you can distribute):
```bash
go build -o myprogram .
go build -o myprogram.exe .  # on Windows
./myprogram
```

The `-o` flag names the output binary. The `.` means "build the package in this directory."

### Testing

Go has a built-in testing framework. Test files:

- Must end in `_test.go` (e.g., `main_test.go`, `calc_test.go`)
- Are ignored by `go build` — they're only for testing
- Are run with:

```bash
go test ./...
```

The `./...` means "test all packages in this directory and subdirectories." You'll write tests starting in the next challenge.

### File Naming Conventions

| Convention | Example | Purpose |
|---|---|---|
| `main.go` | `main.go` | Entry point for executables |
| `*_test.go` | `main_test.go` | Test files (ignored by build) |
| `lowercase.go` | `parser.go` | Regular source files |
| No hyphens | `wordfreq.go` not `word-freq.go` | Go convention |

---

## Quick Reference Card

Save this — you'll use these commands in every challenge:

```bash
# Initialize a new module
go mod init <name>

# Run without building
go run .

# Build a binary
go build -o <binary-name> .

# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Check your Go version
go version
```

---

## Task

Create a Go program that prints exactly `ready` to stdout.

### Requirements

1. Initialize a Go module (pick any name)
2. Create `main.go` with `package main` and `func main()`
3. Print exactly `ready` (no extra spaces, no extra lines beyond the trailing newline from `Println`)
4. The program should build and run successfully

### Expected Output

```
ready
```

### Verify Your Solution

```bash
bash test.sh
```

---

## Hints

<details>
<summary>Hint 1: What function prints to stdout?</summary>

Use `fmt.Println("ready")`. You'll need to import the `"fmt"` package.

</details>

<details>
<summary>Hint 2: Complete structure</summary>

Your `main.go` should have exactly three things:
1. `package main`
2. An import for `"fmt"`
3. A `func main()` that calls `fmt.Println("ready")`

</details>

<details>
<summary>Hint 3: Module initialization</summary>

Run `go mod init setup` in the challenge directory before building.

</details>

---

## Further Reading

- [Go Getting Started](https://go.dev/doc/tutorial/getting-started)
- [How to Write Go Code](https://go.dev/doc/code)
- [Go Modules Reference](https://go.dev/ref/mod)
