# Terminal and ANSI Escape Codes

## What Is It

A **terminal emulator** is a program that displays text and accepts keyboard input. It is the window where you run your shell, execute commands, and see output. Every time you type `go run main.go` and see colored output, a terminal emulator is rendering that text — interpreting special byte sequences to produce bold, colored, or underlined characters.

The **ANSI escape codes** are a standard set of byte sequences that control terminal formatting. They were formalized in 1979 (ANSI X3.64, later adopted as ECMA-48 and ISO 6429) and remain the universal mechanism for terminal colors and styling today.

When grep highlights matches in red, it is not using a graphics API. It is inserting invisible byte sequences into the text output that the terminal interprets as "switch to red" and "switch back to normal."

## What Problem It Solves

### The Display Problem

Programs produce text output. But plain text is flat — every character looks the same. When grep returns 50 lines, you need to instantly see **where** the match is within each line. Without color, you have to visually scan each line for the pattern. With color, the matches jump out.

This is not cosmetic. For a tool that may return thousands of lines, visual differentiation is a **usability requirement**.

### The Portability Problem

Different terminal hardware (and later, different terminal emulator software) supported different capabilities. A program written for a DEC VT100 would produce garbage on a Wyse 50. ANSI escape codes created a **common standard** — a shared language between programs and terminals.

### The Pipe Problem

Color codes are useful when a human is reading the output. But when output is piped to another program (`grep error log.txt | wc -l`), those invisible bytes corrupt the data. The downstream program sees `\033[31merror\033[0m` instead of `error`. Programs need to detect whether they are talking to a human (terminal) or a machine (pipe) and behave accordingly.

## First Principles

### TTY History

The term "TTY" comes from **teletypewriter** — a physical device from the early 1900s that printed characters on paper. When computers arrived:

1. **1960s**: Teletypewriters connected to mainframes as input/output devices
2. **1970s**: Video terminals (like the DEC VT100) replaced paper with screens
3. **1980s**: Terminal emulators — software that simulates a hardware terminal
4. **Today**: Programs like iTerm2, Windows Terminal, GNOME Terminal, and Alacritty are all terminal emulators

The abstraction persists: your shell talks to a "terminal" through a TTY device (`/dev/tty`), whether that terminal is hardware from 1978 or software from 2024.

### How Text Gets Displayed

When a program writes to stdout:

```
Program → stdout (fd 1) → terminal emulator → screen
```

The terminal emulator reads bytes from the program and:
1. Prints normal characters as text
2. Intercepts **escape sequences** and interprets them as formatting commands
3. Renders the formatted text on screen

An escape sequence starts with the **ESC** character (byte `0x1B`, octal `\033`, or `\x1b`). The terminal sees this byte and thinks: "the next bytes are a command, not text to display."

### Escape Sequence Format

The most common format is **CSI** (Control Sequence Introducer):

```
ESC [ <parameters> <command>
```

Where:
- `ESC` = byte `0x1B` (written as `\033` or `\x1b` in code)
- `[` = literal opening bracket
- `<parameters>` = numbers separated by semicolons
- `<command>` = a single letter that identifies the operation

For text formatting, the command letter is `m` (Select Graphic Rendition, or **SGR**):

```
\033[31m    → set foreground color to red
\033[0m     → reset all formatting
\033[1;31m  → set bold AND red
```

## ANSI Color Codes

### SGR (Select Graphic Rendition) Parameters

| Code | Effect |
|------|--------|
| `0` | Reset all attributes |
| `1` | Bold (or bright, depending on terminal) |
| `2` | Dim |
| `3` | Italic (not widely supported) |
| `4` | Underline |
| `7` | Reverse (swap foreground/background) |
| `9` | Strikethrough |

### 16-Color Mode (Original ANSI)

The original standard defines 8 colors, each with a normal and bright variant:

| Color | Foreground | Background | Bright FG | Bright BG |
|-------|-----------|------------|-----------|-----------|
| Black | 30 | 40 | 90 | 100 |
| Red | 31 | 41 | 91 | 101 |
| Green | 32 | 42 | 92 | 102 |
| Yellow | 33 | 43 | 93 | 103 |
| Blue | 34 | 44 | 94 | 104 |
| Magenta | 35 | 45 | 95 | 105 |
| Cyan | 36 | 46 | 96 | 106 |
| White | 37 | 47 | 97 | 107 |

Usage: `\033[<code>m`

```
\033[31m   → red text
\033[42m   → green background
\033[1;33m → bold yellow text
```

### 256-Color Mode

Extends the palette to 256 colors using a different parameter format:

```
\033[38;5;<n>m   → set foreground to color n (0–255)
\033[48;5;<n>m   → set background to color n (0–255)
```

The 256 colors are organized:
- 0–7: standard colors (same as 16-color mode)
- 8–15: bright colors
- 16–231: a 6×6×6 color cube (216 colors)
- 232–255: grayscale ramp (24 shades)

### True Color (24-bit)

Modern terminals support full 24-bit RGB color:

```
\033[38;2;<r>;<g>;<b>m   → set foreground to RGB
\033[48;2;<r>;<g>;<b>m   → set background to RGB
```

Example: `\033[38;2;255;100;0m` sets foreground to orange (R=255, G=100, B=0).

### Color Mode Comparison

| Mode | Colors | Format | Support |
|------|--------|--------|---------|
| 16-color | 16 | `\033[31m` | Universal |
| 256-color | 256 | `\033[38;5;196m` | Nearly universal |
| True color | 16.7M | `\033[38;2;255;0;0m` | Most modern terminals |

For grep, **16-color mode is sufficient**. GNU grep uses bold red (`\033[1;31m`) for matches and bold magenta (`\033[1;35m`) for filenames. These codes work everywhere.

## How Go Does It

### Writing ANSI Codes

ANSI escape codes are just bytes. Write them to stdout like any other string:

```go
package main

import "fmt"

const (
    Reset     = "\033[0m"
    Bold      = "\033[1m"
    Red       = "\033[31m"
    Green     = "\033[32m"
    Yellow    = "\033[33m"
    Blue      = "\033[34m"
    Magenta   = "\033[35m"
    Cyan      = "\033[36m"
    BoldRed   = "\033[1;31m"
    BoldGreen = "\033[1;32m"
)

func main() {
    fmt.Println(BoldRed + "error:" + Reset + " something went wrong")
    fmt.Println(Green + "success:" + Reset + " all tests passed")
    fmt.Println(Yellow + "warning:" + Reset + " deprecated function")
}
```

### Colorizing Grep Matches

```go
package main

import (
    "fmt"
    "regexp"
    "strings"
)

const (
    matchColor    = "\033[1;31m" // bold red — same as GNU grep
    resetColor    = "\033[0m"
    filenameColor = "\033[35m"   // magenta
    lineNumColor  = "\033[32m"   // green
    sepColor      = "\033[36m"   // cyan
)

func colorizeMatches(line string, re *regexp.Regexp) string {
    indices := re.FindAllStringIndex(line, -1)
    if indices == nil {
        return line
    }

    var buf strings.Builder
    prev := 0
    for _, loc := range indices {
        buf.WriteString(line[prev:loc[0]])  // text before match
        buf.WriteString(matchColor)          // start color
        buf.WriteString(line[loc[0]:loc[1]]) // matched text
        buf.WriteString(resetColor)          // end color
        prev = loc[1]
    }
    buf.WriteString(line[prev:]) // text after last match
    return buf.String()
}

func main() {
    re := regexp.MustCompile(`error|warning`)
    line := "2024-01-15 error: connection timeout, warning: retry limit"

    fmt.Println(colorizeMatches(line, re))
    // "error" and "warning" will appear in bold red
}
```

### Detecting Terminal vs Pipe

This is critical. You must not emit color codes when output goes to a pipe or file:

```go
package main

import (
    "fmt"
    "os"

    "golang.org/x/term"
)

func isTerminal(fd int) bool {
    return term.IsTerminal(fd)
}

func main() {
    if isTerminal(int(os.Stdout.Fd())) {
        fmt.Println("\033[32mThis is green (terminal detected)\033[0m")
    } else {
        fmt.Println("This is plain text (piped or redirected)")
    }
}
```

Without the `golang.org/x/term` package, you can use a syscall-based approach:

```go
import (
    "os"
    "golang.org/x/sys/unix"
)

func isTerminal(fd uintptr) bool {
    _, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
    return err == nil
}
```

The underlying mechanism: a TTY device responds to the `TCGETS` ioctl (or `TIOCGETA` on macOS). If the file descriptor is not a terminal (e.g., it is a pipe or regular file), the ioctl fails. This is what "isatty" means — "is this file descriptor attached to a TTY device?"

### Implementing `--color` Modes

GNU grep supports three color modes: `--color=auto`, `--color=always`, `--color=never`. Here is the pattern:

```go
type ColorMode int

const (
    ColorAuto   ColorMode = iota
    ColorAlways
    ColorNever
)

func shouldColorize(mode ColorMode) bool {
    switch mode {
    case ColorAlways:
        return true
    case ColorNever:
        return false
    case ColorAuto:
        // Color only if stdout is a terminal AND NO_COLOR is not set
        if os.Getenv("NO_COLOR") != "" {
            return false
        }
        return term.IsTerminal(int(os.Stdout.Fd()))
    }
    return false
}
```

### Formatted Output with Color

A complete example showing grep-style colored output:

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strings"

    "golang.org/x/term"
)

type Printer struct {
    colorize bool
}

func (p *Printer) printMatch(filename string, lineNum int, line string, re *regexp.Regexp) {
    var buf strings.Builder

    if p.colorize {
        // filename in magenta
        buf.WriteString("\033[35m")
        buf.WriteString(filename)
        buf.WriteString("\033[0m")
        buf.WriteString("\033[36m:\033[0m") // separator in cyan

        // line number in green
        buf.WriteString("\033[32m")
        buf.WriteString(fmt.Sprintf("%d", lineNum))
        buf.WriteString("\033[0m")
        buf.WriteString("\033[36m:\033[0m") // separator in cyan

        // line with matches highlighted
        buf.WriteString(colorizeMatches(line, re))
    } else {
        buf.WriteString(fmt.Sprintf("%s:%d:%s", filename, lineNum, line))
    }

    fmt.Println(buf.String())
}

func colorizeMatches(line string, re *regexp.Regexp) string {
    indices := re.FindAllStringIndex(line, -1)
    if indices == nil {
        return line
    }

    var buf strings.Builder
    prev := 0
    for _, loc := range indices {
        buf.WriteString(line[prev:loc[0]])
        buf.WriteString("\033[1;31m")
        buf.WriteString(line[loc[0]:loc[1]])
        buf.WriteString("\033[0m")
        prev = loc[1]
    }
    buf.WriteString(line[prev:])
    return buf.String()
}

func main() {
    useColor := term.IsTerminal(int(os.Stdout.Fd())) && os.Getenv("NO_COLOR") == ""

    printer := &Printer{colorize: useColor}
    re := regexp.MustCompile(`error`)

    scanner := bufio.NewScanner(os.Stdin)
    lineNum := 0
    for scanner.Scan() {
        lineNum++
        line := scanner.Text()
        if re.MatchString(line) {
            printer.printMatch("stdin", lineNum, line, re)
        }
    }
}
```

## Tradeoffs

| Decision | Benefit | Cost |
|----------|---------|------|
| Always emit color codes | Simplest code | Corrupts output in pipes |
| Never emit color codes | Safe everywhere | No visual benefit |
| Auto-detect terminal | Best UX (color when human reads, plain when piped) | Requires platform-specific isatty check |
| 16-color only | Universal support | Limited palette |
| 256-color or true color | Richer output | Not supported everywhere |
| Respect `NO_COLOR` | User control, accessibility | One more environment variable to check |
| Inline escape codes | Simple, no dependencies | Harder to test, codes scattered in logic |
| Color abstraction layer | Testable, clean | More code, possible dependency |

### Testing Colored Output

Color codes make testing harder. Strategy: separate the coloring decision from the output logic:

```go
// Testable: the function takes a flag, doesn't check the terminal itself
func formatMatch(line string, re *regexp.Regexp, colorize bool) string {
    if !colorize {
        return line
    }
    return colorizeMatches(line, re)
}

// In tests:
func TestFormatMatch(t *testing.T) {
    re := regexp.MustCompile(`error`)

    // Test without color — easy to assert
    got := formatMatch("an error occurred", re, false)
    want := "an error occurred"
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }

    // Test with color — assert the escape codes are present
    got = formatMatch("an error occurred", re, true)
    if !strings.Contains(got, "\033[1;31m") {
        t.Error("expected red color code in output")
    }
    if !strings.Contains(got, "\033[0m") {
        t.Error("expected reset code in output")
    }
}
```

## Cross-Platform Considerations

| Platform | ANSI Support | Notes |
|----------|-------------|-------|
| Linux | Native | All terminals support ANSI codes |
| macOS | Native | Terminal.app and iTerm2 fully support ANSI |
| Windows 10+ | Supported | Windows Terminal supports ANSI natively |
| Windows (legacy cmd.exe) | Limited | Requires `SetConsoleMode` to enable virtual terminal processing |
| CI/CD environments | Varies | Some CI systems set `TERM=dumb`; check before emitting color |

On Windows, you may need to enable ANSI support explicitly:

```go
// Windows-specific: enable virtual terminal processing
// This is only needed for legacy Windows console (cmd.exe)
if runtime.GOOS == "windows" {
    // golang.org/x/sys/windows provides SetConsoleMode
    // Modern Windows Terminal does not need this
}
```

For a cross-platform grep, the safest approach is: auto-detect the terminal, respect `NO_COLOR`, and stick to 16-color codes.

### The `NO_COLOR` Convention

[no-color.org](https://no-color.org) defines a simple convention: if the `NO_COLOR` environment variable is set (to any value), programs should not emit color. This is an accessibility feature — some users have visual impairments that make colored text harder to read, or they use terminal themes where certain colors are invisible against the background.

```go
func colorEnabled() bool {
    // NO_COLOR takes precedence — even over --color=always in strict implementations
    if _, exists := os.LookupEnv("NO_COLOR"); exists {
        return false
    }
    return term.IsTerminal(int(os.Stdout.Fd()))
}
```

## Why This Matters for Grep

Color is one of grep's most valuable features. Real-world grep invocations almost always use color because:

1. **Match visibility**: In a wall of log output, red highlights make matches instantly findable
2. **Context understanding**: When displaying context lines (`-B`, `-A`), color distinguishes matching lines from context lines
3. **Multi-file output**: Colored filenames and line numbers create visual structure in output spanning many files

The `--color=auto` behavior is the default in modern grep for a reason — it provides the best experience without breaking pipelines:

```bash
# Color ON — stdout is a terminal, human is reading
grep error /var/log/syslog

# Color OFF — stdout is a pipe, another program is reading
grep error /var/log/syslog | wc -l

# Color ON — forced, even through a pipe (for piping to `less -R`)
grep --color=always error /var/log/syslog | less -R
```

Implementing this correctly in your grep clone requires understanding everything in this document: escape sequences for coloring, terminal detection for auto mode, `NO_COLOR` for accessibility, and careful string building so escape codes don't break match position calculations.

## Further Reading

1. **[ANSI escape code — Wikipedia](https://en.wikipedia.org/wiki/ANSI_escape_code)** — Comprehensive reference covering CSI sequences, SGR parameters, and the full history from VT100 to modern terminals. The tables of color codes and formatting parameters are especially useful.

2. **[DEC VT100 User Guide](https://vt100.net/docs/vt100-ug/)** — The original terminal that popularized ANSI escape codes. Reading the VT100 documentation gives you a visceral understanding of why terminal interfaces work the way they do.

3. **[XTerm Control Sequences](https://invisible-island.net/xterm/ctlseqs/ctlseqs.html)** — The definitive reference for escape sequences supported by xterm and xterm-compatible terminals. This is the document terminal emulator authors implement against.

4. **[no-color.org](https://no-color.org/)** — The `NO_COLOR` convention. A simple, widely adopted standard for disabling color output. Lists hundreds of supporting programs.

5. **[Go `golang.org/x/term` package](https://pkg.go.dev/golang.org/x/term)** — Go's official terminal utility package. Provides `IsTerminal`, terminal size detection, and raw mode support. This is the standard way to do terminal detection in Go.

6. **["Build Your Own Shell" — CodeCrafters](https://codecrafters.io/challenges/shell)** — Hands-on exercises that build your understanding of terminals, TTY devices, and how shells interact with the terminal. Complements the concepts in this document.

7. **[Microsoft Windows Console documentation](https://learn.microsoft.com/en-us/windows/console/)** — If you need cross-platform terminal support, this documents Windows Console and Virtual Terminal Sequences. Covers `SetConsoleMode` and the differences between legacy console and Windows Terminal.

8. **[Chalk (Node.js) / termcolor (Rust) / fatih/color (Go)](https://github.com/fatih/color)** — Popular terminal color libraries across languages. Reading their source code shows how production software handles color detection, `NO_COLOR`, and cross-platform support. The Go `fatih/color` library is particularly instructive.

---

*Previous: [17 — Regular Expressions](./17-regular-expressions.md)*
