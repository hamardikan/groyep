# CLI Design

## What Is It

A Command-Line Interface (CLI) is a text-based interface where users type commands to interact with a program. Every CLI command follows a common anatomy:

```
program [flags] [arguments]
```

For example:

```bash
grep -i -n "error" server.log access.log
│     │  │  │       │          │
│     │  │  │       └──────────┴── arguments (files to search)
│     │  │  └── argument (the pattern)
│     │  └── flag: show line numbers
│     └── flag: case-insensitive
└── program name
```

CLI design is the art of making this interface intuitive, consistent, and discoverable. Good CLI design means users can guess how your tool works because it follows conventions they already know.

## What Problem It Solves

A program is useless if people can't figure out how to use it. CLI design solves the communication problem between human intent and program behavior:

- How does the user tell the program **what** to do? (arguments)
- How does the user tell the program **how** to do it? (flags)
- What happens when the user makes a mistake? (error messages, help text)
- How does the user learn what's possible? (`--help`, man pages)

Bad CLI design creates tools that are powerful but impenetrable. Good CLI design creates tools that feel obvious — you can often guess the flag name on the first try.

## First Principles

### The Three Parts of a Command

| Part       | Purpose                    | Example in `grep -i "error" log.txt` |
|-----------|----------------------------|---------------------------------------|
| Program    | What tool to run           | `grep`                                |
| Flags      | Modify behavior            | `-i` (case-insensitive)               |
| Arguments  | What to operate on         | `"error"` (pattern), `log.txt` (file) |

### Flag Conventions

There are two major conventions:

**POSIX (short flags):**
- Single dash, single letter: `-v`, `-i`, `-n`
- Can be combined: `-vin` is the same as `-v -i -n`
- Flags with values: `-e pattern` or `-f file`

**GNU (long flags):**
- Double dash, full word: `--verbose`, `--ignore-case`, `--line-number`
- Values with `=`: `--color=always` or `--color always`
- `--` alone means "end of flags, everything after is an argument"

Most modern tools support both:

```bash
grep -i "hello" file.txt       # POSIX short flag
grep --ignore-case "hello" file.txt  # GNU long flag
```

### The `--` Separator

What if your search pattern starts with a dash? Without `--`, the tool thinks it's a flag:

```bash
grep -v file.txt          # -v is a flag (invert match)
grep -- -v file.txt       # -v is the PATTERN (search for literal "-v")
```

The `--` convention says: "everything after this is an argument, not a flag." It's essential for robust CLI tools.

### Help and Usage Messages

Every CLI tool should respond to `--help` (and conventionally `-h`) with a usage message:

```
Usage: groyep [OPTIONS] PATTERN [FILE...]

Search for PATTERN in each FILE or standard input.

Options:
  -i, --ignore-case    Ignore case distinctions in patterns
  -v, --invert-match   Select non-matching lines
  -n, --line-number    Prefix each line with its line number
  -c, --count          Only print a count of matching lines
  -l, --files-with-matches  Only print names of files with matches
  -r, --recursive      Recursively search directories
      --color[=WHEN]   Highlight matches; WHEN is 'always', 'never', or 'auto'
  -h, --help           Display this help and exit

Exit status:
  0  if any match was found
  1  if no match was found
  2  if an error occurred
```

A good usage message tells you:
1. **Synopsis** — the command structure
2. **Description** — what the tool does (one sentence)
3. **Options** — every flag with a brief explanation
4. **Exit codes** — what the return values mean

### Designing Flag Names

| Principle           | Good                         | Bad                      |
|--------------------|------------------------------|--------------------------|
| Predictable        | `-n` for line **n**umber     | `-l` for line number     |
| Consistent         | `-v` for in**v**ert (like grep) | `-x` for invert       |
| Mnemonic           | `-i` for **i**gnore case     | `-z` for ignore case     |
| Follow precedent   | `-r` for **r**ecursive       | `-d` for recursive       |

When in doubt, check what the established tool (real grep) uses and match it. Users have muscle memory.

## How Go Does It

### The `flag` Package

Go's standard library includes the `flag` package for parsing command-line flags:

```go
package main

import (
    "flag"
    "fmt"
)

func main() {
    // Declare flags
    ignoreCase := flag.Bool("i", false, "ignore case distinctions")
    lineNumber := flag.Bool("n", false, "prefix each line with line number")
    count := flag.Bool("c", false, "only print a count of matching lines")
    invertMatch := flag.Bool("v", false, "select non-matching lines")

    // Parse the command line
    flag.Parse()

    // Remaining arguments (after flags)
    args := flag.Args()

    fmt.Printf("ignoreCase: %v\n", *ignoreCase)
    fmt.Printf("lineNumber: %v\n", *lineNumber)
    fmt.Printf("count: %v\n", *count)
    fmt.Printf("invertMatch: %v\n", *invertMatch)
    fmt.Printf("arguments: %v\n", args)
}
```

```bash
$ ./groyep -i -n "hello" file.txt
ignoreCase: true
lineNumber: true
count: false
invertMatch: false
arguments: [hello file.txt]
```

### Flag Types

```go
// Boolean flags
verbose := flag.Bool("v", false, "verbose output")

// String flags
pattern := flag.String("e", "", "use PATTERN as the pattern")

// Integer flags
maxCount := flag.Int("m", 0, "stop after NUM matches")

// You can also bind to existing variables
var color string
flag.StringVar(&color, "color", "auto", "highlight matches: always, never, auto")
```

### Custom Usage Message

```go
func main() {
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage: groyep [OPTIONS] PATTERN [FILE...]\n\n")
        fmt.Fprintf(os.Stderr, "Search for PATTERN in each FILE or standard input.\n\n")
        fmt.Fprintf(os.Stderr, "Options:\n")
        flag.PrintDefaults()
        fmt.Fprintf(os.Stderr, "\nExit status:\n")
        fmt.Fprintf(os.Stderr, "  0  if any match was found\n")
        fmt.Fprintf(os.Stderr, "  1  if no match was found\n")
        fmt.Fprintf(os.Stderr, "  2  if an error occurred\n")
    }

    // ... declare and parse flags ...

    if flag.NArg() < 1 {
        flag.Usage()
        os.Exit(2)
    }
}
```

### Raw Argument Access with `os.Args`

Sometimes you need access to the raw arguments before flag parsing:

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // os.Args[0] is the program name
    // os.Args[1:] are the arguments
    fmt.Println("Program:", os.Args[0])
    fmt.Println("Args:", os.Args[1:])
}
```

```bash
$ ./groyep -i "hello" file.txt
Program: ./groyep
Args: [-i hello file.txt]
```

### A Complete CLI Skeleton

```go
package main

import (
    "flag"
    "fmt"
    "os"
)

// Config holds all the parsed command-line options
type Config struct {
    IgnoreCase  bool
    InvertMatch bool
    LineNumber  bool
    Count       bool
    Color       string
    Pattern     string
    Files       []string
}

func parseArgs() (*Config, error) {
    cfg := &Config{}

    flag.BoolVar(&cfg.IgnoreCase, "i", false, "ignore case distinctions")
    flag.BoolVar(&cfg.InvertMatch, "v", false, "select non-matching lines")
    flag.BoolVar(&cfg.LineNumber, "n", false, "prefix each line with line number")
    flag.BoolVar(&cfg.Count, "c", false, "only print count of matching lines")
    flag.StringVar(&cfg.Color, "color", "auto", "highlight matches: always, never, auto")

    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage: groyep [OPTIONS] PATTERN [FILE...]\n\n")
        fmt.Fprintf(os.Stderr, "Search for PATTERN in each FILE or standard input.\n\n")
        fmt.Fprintf(os.Stderr, "Options:\n")
        flag.PrintDefaults()
    }

    flag.Parse()

    args := flag.Args()
    if len(args) < 1 {
        return nil, fmt.Errorf("missing pattern argument")
    }

    cfg.Pattern = args[0]
    cfg.Files = args[1:]

    return cfg, nil
}

func main() {
    cfg, err := parseArgs()
    if err != nil {
        fmt.Fprintf(os.Stderr, "groyep: %v\n", err)
        fmt.Fprintf(os.Stderr, "Try 'groyep --help' for more information.\n")
        os.Exit(2)
    }

    // Use cfg to run the search...
    _ = cfg
}
```

### Limitations of the `flag` Package

Go's `flag` package has some quirks compared to POSIX/GNU conventions:

| Feature               | POSIX/GNU grep         | Go `flag` package        |
|-----------------------|------------------------|--------------------------|
| Short flags           | `-i`                   | `-i` (supported)         |
| Long flags            | `--ignore-case`        | `-ignore-case` (single dash!) |
| Combined short flags  | `-vin`                 | Not supported            |
| `--` separator        | Supported              | Supported                |
| Flag=value            | `--color=always`       | `-color=always`          |

Go's `flag` package uses single-dash for all flags. It doesn't support GNU-style double-dash long flags or POSIX-style flag combining. For a learning project, this is fine. For a production tool, you might consider third-party packages like `pflag` or `cobra`.

## Tradeoffs

| Decision                              | Benefit                           | Cost                                 |
|--------------------------------------|-----------------------------------|--------------------------------------|
| Use `flag` package (stdlib)          | Simple, no dependencies           | No GNU-style long flags              |
| Use `pflag`/`cobra` (third-party)   | Full POSIX/GNU compliance         | External dependency                  |
| Many flags                           | Flexible, powerful                | Overwhelming for new users           |
| Few flags                            | Simple, approachable              | May not cover all use cases          |
| Required flags                       | Explicit behavior                 | More typing for common operations    |
| Sensible defaults                    | Less typing, "just works"         | Implicit behavior can surprise       |
| Flags-only (no positional args)      | Explicit, self-documenting        | Verbose for simple cases             |
| Positional arguments                 | Concise, natural                  | Order matters, less discoverable     |

### Error Messages: A Subtle Art

Good error messages tell the user what went wrong AND how to fix it:

```
# Bad
groyep: error

# Better
groyep: missing pattern argument

# Best
groyep: missing pattern argument
Try 'groyep --help' for more information.
```

The convention for error messages: `program: description`. Lowercase. No period. Point the user toward help.

## Why This Matters for Grep

The CLI is your user's first impression of your tool. Grep's CLI is one of the most well-known in computing — millions of developers use it daily. Your grep clone should feel familiar:

```bash
# These should all work exactly as a user expects:
groyep "pattern" file.txt          # Basic search
groyep -i "pattern" file.txt       # Case-insensitive
groyep -n "pattern" file.txt       # With line numbers
groyep -v "pattern" file.txt       # Invert match
groyep -c "pattern" file.txt       # Count matches
groyep "pattern" file1.txt file2.txt  # Multiple files
cat file.txt | groyep "pattern"    # Read from stdin
groyep "pattern"                   # Read from stdin (no file)
```

The pattern is always the first non-flag argument. Files are the remaining arguments. If no files are given, read from stdin. This convention is so deeply ingrained that violating it would confuse every user.

Your `Config` struct is the bridge between CLI parsing and program logic:

```go
// The CLI parsing produces a Config
// The search logic consumes a Config
// They never know about each other's details
cfg, err := parseArgs()
if err != nil {
    fmt.Fprintf(os.Stderr, "groyep: %v\n", err)
    os.Exit(2)
}
exitCode := search(cfg)
os.Exit(exitCode)
```

This separation makes both halves testable independently.

## Further Reading

1. **[POSIX Utility Conventions](https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html)** — The formal specification for how Unix command-line utilities should behave. Dry but authoritative — this is the standard that grep follows.

2. **[Go `flag` package documentation](https://pkg.go.dev/flag)** — The official reference for Go's built-in flag parsing. Read the package-level documentation, not just the function signatures.

3. **[GNU Coding Standards — Command-Line Interfaces](https://www.gnu.org/prep/standards/html_node/Command_002dLine-Interfaces.html)** — GNU's extensions to the POSIX conventions. Explains `--long-flags`, `--help`, and `--version` conventions.

4. **[Command Line Interface Guidelines (clig.dev)](https://clig.dev/)** — A modern, comprehensive guide to CLI design. Covers everything from help text to output formatting. Excellent resource.

5. **[12 Factor CLI Apps](https://medium.com/@jdxcode/12-factor-cli-apps-dd3c227a0e46)** — Jeff Dickey's adaptation of the 12-factor app methodology for CLI tools. Practical advice for building professional-quality CLIs.

6. **[Cobra library documentation](https://cobra.dev/)** — While you won't use Cobra for this project (learning the `flag` package is more valuable), understanding what it provides gives context for what production CLIs need.

7. **[grep man page](https://man7.org/linux/man-pages/man1/grep.1.html)** — Study the original. Every flag, every option, every convention. Your clone's CLI should be a subset of this.

8. **[Rob Pike — "Go Command-Line Programming" (dotGo 2015)](https://www.youtube.com/watch?v=1B71SL6Y0kA)** — Rob Pike discusses building command-line tools in Go, covering the design decisions that shaped Go's standard library.

9. **[Effective Go — Package flag example](https://go.dev/doc/effective_go)** — The official guide shows idiomatic flag usage in context.

10. **[Writing Friendly Command Line Applications (Carolyn Van Slyck)](https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/)** — Practical advice on making CLIs that users actually enjoy using.
