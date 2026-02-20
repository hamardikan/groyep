# Learning Pathway: Build Your Own Grep

A progressive, from-scratch learning pathway that takes you from zero Go knowledge to building a fully functional `grep` clone.

**Language:** Go
**Final Goal:** A working `grep` command-line tool with regex, recursive search, flags, and more.
**Philosophy:** You write every line of code yourself. Tests tell you *what* to build; you figure out *how*.

---

## How This Works

### Structure

```
pathways/grep/
├── concepts/       # Knowledge Library — the "why" and "what"
├── challenges/     # 26 progressive challenges — the "do"
│   ├── 00-setup/
│   ├── 01-hello-go/
│   ├── ...
│   └── 25-grep-final/
```

### Rules

1. **Read the concept docs** linked in each challenge before you start coding
2. **You create ALL source files** from scratch — no skeleton code is provided
3. **Each challenge is its own Go module** — you run `go mod init` yourself every time
4. **Tests are your spec** — `*_test.go` files tell you what functions/behavior is expected
5. **Two ways to verify your work:**
   - `go test ./...` — runs Go unit tests against your functions
   - `bash test.sh` — builds your binary and tests it as a CLI tool
6. **Don't skip ahead** — each challenge builds on the previous one
7. **Try before reading hints** — hints are hidden in collapsible sections

### What You'll Need

- Go 1.21+ installed ([go.dev/dl](https://go.dev/dl/))
- A terminal (bash/zsh)
- A text editor (VS Code, Vim, whatever you prefer)
- Curiosity and patience

---

## Skill Tree

```
                         [25: FINAL GREP]
                              |
                   +----------+----------+
                   |          |          |
              [23:Color] [24:Context] [22:Case-Insensitive]
                   |          |          |
                   +----------+----------+
                              |
                         [21:Anchors]
                              |
                         [20:Regex]
                              |
                         [19:Invert -v]
                              |
                         [18:Recursive -r]
                              |
                         [17:Literal Match]
                              |
                         [16:Empty Match]              ← GREP BEGINS
                              |
              +---------+-----+-----+---------+
              |         |           |         |
         [12:CLI   [13:Flag    [14:Structs [15:Packages]
          Args]    Parsing]    Methods]
              |         |           |         |
              +---------+-----+-----+---------+
                              |
                   +----------+----------+
                   |          |          |
              [08:Read   [09:Write  [10:Stdin  [11:Exit
               File]     Stdout]    Pipe]      Codes]
                   |          |          |         |
                   +----------+-----+----+---------+
                              |
                         [07:Strings]
                              |
                   +----------+----------+
                   |          |          |
              [04:Loops] [05:Functions] [06:Collections]
                   |          |          |
                   +----------+----------+
                              |
                         [03:Conditionals]
                              |
                         [02:Variables & Types]
                              |
                         [01:Hello Go]
                              |
                         [00:Setup]                    ← START HERE
```

---

## Roadmap

### Phase 0: Bootstrapping

| # | Challenge | You Learn | You Build |
|---|-----------|-----------|-----------|
| 00 | [Setup](challenges/00-setup/challenge.md) | Go toolchain, `go mod init`, project structure, `go test` | Verify your environment works |

### Phase 1: Go Foundations

| # | Challenge | You Learn | You Build |
|---|-----------|-----------|-----------|
| 01 | [Hello Go](challenges/01-hello-go/challenge.md) | `package main`, `func main`, `fmt.Println` | "Hello, World!" program |
| 02 | [Variables & Types](challenges/02-variables-and-types/challenge.md) | `var`, `:=`, `string`, `int`, `bool`, `fmt.Printf` | Formatted info printer |
| 03 | [Conditionals](challenges/03-conditionals/challenge.md) | `if/else`, `switch`, `os.Args`, comparisons | Number classifier |
| 04 | [Loops](challenges/04-loops/challenge.md) | `for` (all forms), `range`, `break`, `continue` | FizzBuzz |
| 05 | [Functions](challenges/05-functions/challenge.md) | Declarations, params, returns, multiple returns, errors | Mini calculator |
| 06 | [Collections](challenges/06-collections/challenge.md) | Slices, maps, `append`, `len`, iteration | Word frequency counter |

### Phase 2: Strings & I/O

| # | Challenge | You Learn | You Build |
|---|-----------|-----------|-----------|
| 07 | [Strings](challenges/07-strings/challenge.md) | `strings` package, `Contains`, runes vs bytes | String checker utility |
| 08 | [Read a File](challenges/08-read-file/challenge.md) | `os.Open`, `bufio.Scanner`, `defer` | Your own `cat` |
| 09 | [Write to Stdout](challenges/09-write-stdout/challenge.md) | `os.Stdout`, `fmt.Fprintln`, `bufio.Writer` | `cat` with line numbers |
| 10 | [Stdin & Pipes](challenges/10-stdin-pipe/challenge.md) | `os.Stdin`, piped input, Unix pipelines | Uppercase filter |
| 11 | [Exit Codes](challenges/11-exit-codes/challenge.md) | `os.Exit`, exit conventions | Substring finder with exit codes |

### Phase 3: CLI & Code Structure

| # | Challenge | You Learn | You Build |
|---|-----------|-----------|-----------|
| 12 | [CLI Arguments](challenges/12-cli-args/challenge.md) | `os.Args`, validation, usage messages | `./search pattern file` |
| 13 | [Flag Parsing](challenges/13-flag-parsing/challenge.md) | `flag` package, bool/string flags | Add `-n` and `-c` flags |
| 14 | [Structs & Methods](challenges/14-structs-methods/challenge.md) | `struct`, methods, receivers | `Matcher` struct |
| 15 | [Packages](challenges/15-packages/challenge.md) | Multiple packages, `internal/`, `cmd/` | Proper Go project layout |

### Phase 4: Building Grep

| # | Challenge | You Learn | You Build |
|---|-----------|-----------|-----------|
| 16 | [Grep: Empty Match](challenges/16-grep-empty/challenge.md) | Integration of all concepts | `mygrep "" file` |
| 17 | [Grep: Literal Match](challenges/17-grep-literal/challenge.md) | String matching + exit codes | `mygrep pattern file` |
| 18 | [Grep: Recursive](challenges/18-grep-recursive/challenge.md) | `filepath.WalkDir`, `-r` flag | `mygrep -r pattern dir` |
| 19 | [Grep: Invert](challenges/19-grep-invert/challenge.md) | Inverting logic, `-v` flag | `mygrep -v pattern file` |
| 20 | [Grep: Regex](challenges/20-grep-regex/challenge.md) | `regexp` package | `mygrep "\d" file` |
| 21 | [Grep: Anchors](challenges/21-grep-anchors/challenge.md) | `^` and `$` in regex | `mygrep "^A" file` |
| 22 | [Grep: Case Insensitive](challenges/22-grep-case/challenge.md) | `(?i)` regex, `-i` flag | `mygrep -i pattern file` |

### Phase 5: Polish & Beyond

| # | Challenge | You Learn | You Build |
|---|-----------|-----------|-----------|
| 23 | [Grep: Color](challenges/23-grep-color/challenge.md) | ANSI escape codes | `--color` flag |
| 24 | [Grep: Context](challenges/24-grep-context/challenge.md) | Line buffering, context | `-A`, `-B`, `-C` flags |
| 25 | [Grep: Final](challenges/25-grep-final/challenge.md) | Integration, edge cases | Complete grep binary |

---

## Knowledge Library

The [concepts/](concepts/) directory contains in-depth reading material for each topic.
Each doc covers: **what it is**, **what problem it solves**, **first principles**, **tradeoffs**, **how Go does it**, and **further reading** with comprehensive external links.

See the [Concepts README](concepts/README.md) for the full index.

---

## Test Data

The `root/` directory at the repo root contains test data files used by grep challenges:

```
root/
├── pg132.txt              # "The Art of War" by Sun Tzu (7,137 lines)
├── rockbands.txt          # 101 rock band names
├── symbols.txt            # Special characters + words
└── test-subdir/
    └── BFS1985.txt        # "1985" by Bowling for Soup (41 lines)
```

Grep challenges (16-25) include a local `testdata/` directory with copies of these files.

---

## Progress Tracker

Use this to track your progress. Copy this to a local file and check off as you go:

```
[ ] 00 - Setup
[ ] 01 - Hello Go
[ ] 02 - Variables & Types
[ ] 03 - Conditionals
[ ] 04 - Loops
[ ] 05 - Functions
[ ] 06 - Collections
[ ] 07 - Strings
[ ] 08 - Read a File
[ ] 09 - Write to Stdout
[ ] 10 - Stdin & Pipes
[ ] 11 - Exit Codes
[ ] 12 - CLI Arguments
[ ] 13 - Flag Parsing
[ ] 14 - Structs & Methods
[ ] 15 - Packages
[ ] 16 - Grep: Empty Match
[ ] 17 - Grep: Literal Match
[ ] 18 - Grep: Recursive
[ ] 19 - Grep: Invert
[ ] 20 - Grep: Regex
[ ] 21 - Grep: Anchors
[ ] 22 - Grep: Case Insensitive
[ ] 23 - Grep: Color
[ ] 24 - Grep: Context
[ ] 25 - Grep: Final
```
