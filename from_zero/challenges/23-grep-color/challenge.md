# Challenge 23: Grep - Colorized Output

## Prerequisites

- Complete [Challenge 22](../22-grep-case/challenge.md)
- [Terminal and ANSI](../../concepts/18-terminal-and-ansi.md)

## Goal

Add `--color` flag to highlight the matched portion of each line in red. This makes grep output much easier to scan visually.

## Your Task

Add the `--color` flag with three modes:

| Mode | Behavior |
|------|----------|
| `--color=always` | Always colorize matched text |
| `--color=never` | Never colorize (default, same as no flag) |
| `--color=auto` | Colorize only when stdout is a terminal |
| `--color` (no value) | Same as `--color=always` |

### How Color Works

Wrap the **matched portion** (not the entire line) with ANSI escape codes:

- Start color: `\033[1;31m` (bold red)
- Reset color: `\033[0m`

Example: If pattern is "Nirvana" and line is "Nirvana", output is:
```
\033[1;31mNirvana\033[0m
```

If pattern is "an" and line is "Van Halen", output is:
```
V\033[1;31man\033[0m Halen
```

### Multiple Matches Per Line

If a pattern matches multiple times in one line, colorize each occurrence:

```bash
$ ./mygrep --color=always "an" testdata/rockbands.txt
# "Danger Danger" → "D\033[1;31man\033[0mger D\033[1;31man\033[0mger"
```

### With Regex

Regex matches should also be colorized:

```bash
$ ./mygrep --color=always "\d+" testdata/test-subdir/BFS1985.txt
# Each digit sequence gets colorized individually
```

### Implementation Note

Since the `flag` package doesn't natively support `--long-flags`, you have two options:

1. Parse `--color` manually from `os.Args` before `flag.Parse()`
2. Use `flag.String("color", "never", "...")` which accepts `-color=always`

Either approach is fine. The test suite accepts both `-color` and `--color`.

## How to Test

```bash
go test -v
bash test.sh
```

## Hints

<details><summary>Hint 1: Finding match positions</summary>

Use `re.FindAllStringIndex(line, -1)` to get all match positions, then rebuild the string with color codes inserted.

</details>

<details><summary>Hint 2: Building colorized output</summary>

```go
indices := re.FindAllStringIndex(line, -1)
var result strings.Builder
prev := 0
for _, loc := range indices {
    result.WriteString(line[prev:loc[0]])      // before match
    result.WriteString("\033[1;31m")            // start color
    result.WriteString(line[loc[0]:loc[1]])     // matched text
    result.WriteString("\033[0m")               // reset
    prev = loc[1]
}
result.WriteString(line[prev:])                 // after last match
```

</details>

<details><summary>Hint 3: Detecting terminal</summary>

```go
import "golang.org/x/term"
if term.IsTerminal(int(os.Stdout.Fd())) {
    // stdout is a terminal, use color
}
```

Or without the dependency, check if stdout is a char device (less portable).

</details>

---

Previous: [Challenge 22 - Case Insensitive](../22-grep-case/challenge.md) | Next: [Challenge 24 - Context](../24-grep-context/challenge.md)
