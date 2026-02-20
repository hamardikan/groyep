# 03 — Modules and Packages

## What Is It

Every non-trivial program needs a way to organize code into logical units. In Go, there are two levels of organization:

- A **package** is a directory of `.go` files that belong together. It is Go's unit of compilation and code reuse.
- A **module** is a collection of packages that are released, versioned, and distributed together. It is Go's unit of dependency management.

If a package is a chapter, a module is the book.

## What Problem It Solves

Without code organization:
- All code lives in one massive file — impossible to navigate
- You cannot reuse code across projects without copying it
- You cannot manage external dependencies (libraries other people wrote)
- Two files might define the same name and collide
- There is no way to control what is public and what is private

Go's package and module system solves all of these problems with minimal ceremony.

### A Brief History: GOPATH

Before modules (pre-Go 1.11), Go used **GOPATH** — a single workspace directory where all Go code lived. Every project, every dependency, everything went into `$GOPATH/src/`. This had serious problems:

- You could only have one version of a dependency across all projects
- No lockfiles — builds were not reproducible
- You had to set `GOPATH` correctly or nothing worked
- All code had to live in a specific directory structure

**Go modules** (introduced in Go 1.11, default since Go 1.13) replaced GOPATH. You no longer need to put your code in any special directory. A `go.mod` file at the root of your project declares what module it is and what it depends on.

You will only ever use modules. GOPATH is history. But you will see references to it in older tutorials and Stack Overflow answers — now you know to ignore them.

## First Principles

### What Is a Package?

A package is a **directory** containing one or more `.go` files that all declare the same `package` name at the top.

```
groyep/
├── main.go          ← package main
├── matcher/
│   ├── literal.go   ← package matcher
│   └── regex.go     ← package matcher
└── output/
    └── printer.go   ← package output
```

Rules:
1. All `.go` files in one directory must declare the same package name
2. The package name is the last element of the import path (by convention)
3. Package names are lowercase, single words — no underscores, no mixedCaps

```go
// In matcher/literal.go
package matcher

func ContainsLiteral(line, pattern string) bool {
    return strings.Contains(line, pattern)
}
```

### What Is a Module?

A module is defined by a `go.mod` file at the root of your project. It declares:
- The **module path** — a unique identifier, usually a URL-like string
- The **Go version** the module requires
- Any **dependencies** on other modules

```
// go.mod
module github.com/yourusername/groyep

go 1.22
```

### The `go.mod` File Explained

```
module github.com/yourusername/groyep    ← module path (unique identifier)

go 1.22                                   ← minimum Go version

require (                                 ← dependencies (none yet)
    golang.org/x/text v0.14.0            ← example: module path + version
)
```

Key points:
- The module path does not have to be a real URL, but using your repository URL is the convention
- The `go` directive specifies the minimum Go version, not the exact version
- The `require` block lists dependencies with their exact versions (semantic versioning)
- A `go.sum` file is auto-generated alongside `go.mod` — it contains cryptographic hashes of dependencies for verification. Do not edit it manually, but do commit it to version control.

### Import Paths

You import packages using their full module path + package path:

```go
import (
    "fmt"                                    // standard library
    "os"                                     // standard library
    "github.com/yourusername/groyep/matcher" // your own package
)
```

Standard library packages are imported by their short name. Your own packages and third-party packages use the full module path.

## How Go Does It

### Creating a Module

Every standalone Go project starts with `go mod init`:

```bash
mkdir groyep
cd groyep
go mod init github.com/yourusername/groyep
```

This creates `go.mod`. You are now ready to write Go code.

### The `main` Package

Every executable Go program must have exactly one `package main` with a `main` function:

```go
package main

import "fmt"

func main() {
    fmt.Println("groyep starting")
}
```

- `package main` tells the compiler "this package should produce an executable"
- `func main()` is the entry point — where execution begins
- A project can have multiple `package main` files in different directories (for multiple commands), but within one directory, all files must agree on the package name

### Package Visibility

Go uses a beautifully simple rule for visibility:

- **Exported** (public): name starts with an uppercase letter
- **Unexported** (private): name starts with a lowercase letter

```go
package matcher

// Search is exported — other packages can call matcher.Search()
func Search(line, pattern string) bool {
    return containsLiteral(line, pattern)
}

// containsLiteral is unexported — only usable within package matcher
func containsLiteral(line, pattern string) bool {
    return strings.Contains(line, pattern)
}
```

No keywords like `public`, `private`, or `protected`. No access modifiers. Just capitalization.

### Package Naming Conventions

| Convention | Example | Rationale |
|-----------|---------|-----------|
| Short, lowercase, one word | `matcher`, `output`, `config` | Easy to type and remember |
| No underscores or mixedCaps | `matcher` not `string_matcher` | Go convention |
| Noun, not verb | `matcher` not `matching` | Packages are things, not actions |
| Do not stutter | `matcher.Search()` not `matcher.MatcherSearch()` | The package name is already context |
| Avoid `util`, `common`, `misc` | Be specific | These names carry no meaning |

### Importing Packages

Go imports are explicit. You list exactly what you use:

```go
import (
    "bufio"
    "fmt"
    "os"
    "strings"
)
```

If you import a package and do not use it, the compiler **rejects your code**. This is not a warning — it is an error. This keeps imports clean and makes dependencies obvious.

If you need to import a package only for its side effects (its `init` function), use the blank identifier:

```go
import _ "image/png" // register PNG decoder
```

### `go mod tidy`

Over time, your `go.mod` might list dependencies you no longer use, or you might use packages without adding them to `go.mod`. The `go mod tidy` command fixes both:

```bash
go mod tidy
```

It:
1. Adds any missing module requirements
2. Removes any unused module requirements
3. Updates `go.sum` accordingly

Run it whenever you add or remove imports.

### Project Structure

For a small CLI tool like groyep, a simple structure works:

```
groyep/
├── go.mod
├── go.sum
├── main.go              ← package main, entry point
├── matcher/
│   ├── literal.go       ← package matcher
│   └── regex.go         ← package matcher
├── reader/
│   └── reader.go        ← package reader
└── output/
    └── formatter.go     ← package output
```

For more complex projects, the community has conventions, but there is no single mandated layout. Start simple. Split into packages when you have a clear reason to.

### The `init` Function

Each package can define an `init` function that runs automatically when the package is imported:

```go
package matcher

func init() {
    // runs before main(), when this package is first imported
    // useful for setup that must happen before any other code
}
```

Use `init` sparingly. It makes code harder to reason about because it runs implicitly.

## Tradeoffs

| You gain | You give up |
|----------|-------------|
| Clear dependency graph — imports are explicit | Must list every import (tooling helps) |
| Fast compilation — packages compile independently | Cannot have circular imports (package A imports B and B imports A) |
| Simple visibility — capitalization | No fine-grained access control (e.g., "visible to subpackages only") |
| Reproducible builds — `go.sum` locks dependency hashes | Must manage `go.mod` and `go.sum` (but `go mod tidy` handles most of it) |
| No GOPATH — code lives anywhere | Module path naming requires thought for published modules |
| Single convention for code organization | Less flexibility than languages with more organizational primitives |

### The Circular Import Restriction

Go does not allow circular dependencies between packages. If `matcher` imports `output` and `output` imports `matcher`, the compiler refuses. This forces you to think about dependency direction and often leads to cleaner architecture.

## Why This Matters for Grep

Your grep clone will be organized as a Go module with multiple packages:

- `package main` — the entry point, parses arguments, orchestrates everything
- A matching package — contains the search logic (literal, regex)
- Possibly an output package — handles formatting results with colors, line numbers
- Possibly a reader package — handles reading from files and stdin

The `go.mod` file makes your project self-contained. Anyone can clone your repository, run `go build`, and get a working binary. No dependency setup, no environment configuration.

Understanding packages helps you make good decisions about where code lives. Should the regex matcher and the literal matcher be in the same package? Probably — they solve the same problem. Should the output formatter be in the same package as the matcher? Probably not — they are different concerns.

Good package boundaries make your code easier to test, easier to change, and easier to understand.

## Further Reading

1. **[How to Write Go Code](https://go.dev/doc/code)** — Official
   The official guide to organizing Go code. Start here for the basics of modules, packages, and project structure.

2. **[Using Go Modules](https://go.dev/blog/using-go-modules)** — Go Blog, 2019
   Part 1 of a series on Go modules. Walks through creating a module, adding dependencies, and upgrading them.

3. **[Go Modules Reference](https://go.dev/ref/mod)** — Official
   The complete, authoritative reference for the module system. Dense but comprehensive — use it when you need to understand a specific behavior.

4. **[Organizing a Go Module](https://go.dev/doc/modules/layout)** — Official
   Official guidance on project layout and structuring packages within a module.

5. **[Go Package Naming](https://go.dev/blog/package-names)** — Go Blog
   Best practices for naming packages. Short, clear, and no stuttering.

6. **[Effective Go — Package Names](https://go.dev/doc/effective_go#package-names)** — Official
   The section from Effective Go on how to name packages idiomatically.

7. **[Go Modules: v2 and Beyond](https://go.dev/blog/v2-go-modules)** — Go Blog
   How Go handles major version changes in modules. Not needed immediately, but important to understand.

8. **[Russ Cox — "The Principles of Versioning in Go"](https://research.swtch.com/vgo-principles)** — 2018
   The design principles behind Go modules, by the person who designed them.

---

*Previous: [02 — What Is Go](./02-what-is-go.md) · Next: [04 — Types and Memory](./04-types-and-memory.md)*
