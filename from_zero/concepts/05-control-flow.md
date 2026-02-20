# 05 — Control Flow

## What Is It

**Control flow** is the order in which individual statements, instructions, or function calls are executed in a program. Without control flow, a program runs from top to bottom, executing every line exactly once. That is almost never what you want.

You need your program to make decisions ("if the pattern matches, print the line"), repeat actions ("for every line in the file, search for the pattern"), and skip ahead or bail out when appropriate.

Control flow is how you teach a program to make choices and repeat work.

## What Problem It Solves

Consider what grep does:

1. Read a line from input
2. Check if the pattern exists in the line
3. If yes, print the line
4. Repeat until there are no more lines
5. Exit with an appropriate code

Steps 2-3 require **branching** (making a decision). Step 4 requires **looping** (repeating). Step 5 requires knowing when to **stop**. Without control flow constructs, you cannot express any of this.

## First Principles

All control flow comes down to two capabilities:

1. **Conditional execution** — run this code only if a condition is true
2. **Repeated execution** — run this code multiple times

Everything else — `switch`, `for range`, `break`, `continue`, `goto` — is syntactic convenience built on these two ideas.

At the CPU level, control flow is implemented with **jumps** — instructions that change which instruction the CPU executes next instead of just moving forward sequentially. An `if` compiles down to a conditional jump. A `for` compiles down to a conditional jump back to an earlier instruction.

## How Go Does It

### If/Else

The most basic branching construct:

```go
if matched {
    fmt.Println(line)
}
```

```go
if count > 0 {
    fmt.Printf("Found %d matches\n", count)
} else {
    fmt.Println("No matches found")
}
```

```go
if count == 0 {
    os.Exit(1)
} else if count == 1 {
    fmt.Println("1 match")
} else {
    fmt.Printf("%d matches\n", count)
}
```

Go-specific rules:
- Braces `{}` are **always required**, even for single statements
- The condition **does not** use parentheses (unlike C/Java)
- The opening brace must be on the same line as `if` (enforced by `gofmt`)

### If with Init Statement

Go's most distinctive `if` feature — you can run a short statement before the condition:

```go
if err := scanner.Err(); err != nil {
    fmt.Fprintf(os.Stderr, "error: %v\n", err)
    os.Exit(2)
}
```

The variable `err` is scoped to the `if/else` block. It does not leak into the surrounding function. This pattern is everywhere in Go:

```go
if file, err := os.Open(path); err != nil {
    return fmt.Errorf("cannot open %s: %w", path, err)
} else {
    defer file.Close()
    // use file...
}
```

More commonly, you will see it without the `else`:

```go
file, err := os.Open(path)
if err != nil {
    return fmt.Errorf("cannot open %s: %w", path, err)
}
defer file.Close()
// use file...
```

This "handle the error, then continue with the happy path" style is idiomatic Go. The success path flows straight down.

### Switch

Go's `switch` is more powerful and safer than in most languages:

```go
switch exitCode {
case 0:
    fmt.Println("matches found")
case 1:
    fmt.Println("no matches")
case 2:
    fmt.Println("error occurred")
default:
    fmt.Println("unknown exit code")
}
```

Key differences from C/Java:

| Behavior | C / Java | Go |
|----------|----------|-----|
| **Fallthrough** | Falls through by default (need `break`) | Does NOT fall through (need `fallthrough`) |
| **Break** | Required to prevent fallthrough | Implicit at end of each case |
| **Case values** | Must be constants | Can be any expression |
| **Multiple values per case** | No | Yes |

```go
// Multiple values in one case
switch char {
case 'a', 'e', 'i', 'o', 'u':
    fmt.Println("vowel")
default:
    fmt.Println("consonant")
}
```

#### Expression-less Switch

Go allows `switch` without a value — it acts like a chain of `if/else`:

```go
switch {
case line == "":
    // skip empty lines
case strings.HasPrefix(line, "#"):
    // skip comments
case strings.Contains(line, pattern):
    fmt.Println(line)
}
```

This is cleaner than nested `if/else` when you have multiple unrelated conditions.

#### The `fallthrough` Keyword

If you explicitly want to fall through to the next case (rare):

```go
switch level {
case "error":
    errorCount++
    fallthrough
case "warning":
    warningCount++
}
// If level is "error", BOTH counters increment
```

Use `fallthrough` sparingly. It exists but is almost never the right tool.

### Why No Ternary Operator

Many languages have `condition ? valueIfTrue : valueIfFalse`. Go deliberately omits this. From the FAQ:

> The reason `?:` is absent from Go is that the language's designers had seen the operation used too often to create impenetrably complex expressions.

In Go, you write:

```go
result := ""
if matched {
    result = "MATCH"
} else {
    result = "NO MATCH"
}
```

More verbose? Yes. Easier to read at 2 AM when debugging production? Also yes.

### For Loops

Go has **one** looping construct: `for`. But it covers every case.

#### C-style For Loop

```go
for i := 0; i < len(lines); i++ {
    fmt.Println(lines[i])
}
```

Three components: `init; condition; post`. All three are optional.

#### While-style Loop

Omit the init and post — you get a while loop:

```go
for scanner.Scan() {
    line := scanner.Text()
    if strings.Contains(line, pattern) {
        fmt.Println(line)
    }
}
```

This is how grep reads input — keep scanning while there are more lines.

#### Infinite Loop

Omit everything:

```go
for {
    // runs forever until break, return, or os.Exit
    conn, err := listener.Accept()
    if err != nil {
        break
    }
    go handle(conn)
}
```

#### Range Loop

The `range` keyword iterates over data structures:

```go
// Over a slice
for index, value := range lines {
    fmt.Printf("Line %d: %s\n", index, value)
}

// Over a map
for key, value := range config {
    fmt.Printf("%s = %s\n", key, value)
}

// Over a string (gives runes, not bytes!)
for index, char := range "hello 世界" {
    fmt.Printf("byte %d: %c\n", index, char)
}

// When you only need the index
for i := range lines {
    process(lines[i])
}

// When you only need the value
for _, line := range lines {
    fmt.Println(line)
}
```

The blank identifier `_` discards a value you do not need. This is required — Go does not let you declare a variable and not use it.

### Break and Continue

**`break`** exits the innermost loop:

```go
for scanner.Scan() {
    line := scanner.Text()
    if strings.Contains(line, pattern) {
        fmt.Println(line)
        break  // stop after first match (like grep -m 1)
    }
}
```

**`continue`** skips to the next iteration:

```go
for scanner.Scan() {
    line := scanner.Text()
    if line == "" {
        continue  // skip blank lines
    }
    if strings.Contains(line, pattern) {
        fmt.Println(line)
    }
}
```

#### Labeled Break

When you have nested loops, `break` only exits the innermost one. Labels let you break out of an outer loop:

```go
outer:
for _, file := range files {
    scanner := createScanner(file)
    for scanner.Scan() {
        if strings.Contains(scanner.Text(), pattern) {
            fmt.Println(file)
            break outer  // found a match, move to next file
        }
    }
}
```

### Early Returns

Go encourages **early returns** to handle edge cases and errors at the top of a function, keeping the main logic at low indentation:

```go
func searchFile(path, pattern string) ([]string, error) {
    // Handle error cases early
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Main logic flows straight down
    var matches []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        if strings.Contains(scanner.Text(), pattern) {
            matches = append(matches, scanner.Text())
        }
    }
    return matches, scanner.Err()
}
```

This is called the **"happy path"** pattern. Errors are handled immediately and returned. The successful code path is the one that falls through all the checks.

### Go Has No While or Do-While

This is intentional. `for` does everything:

```go
// "while" equivalent
for condition {
    // ...
}

// "do-while" equivalent (loop body always runs at least once)
for {
    // ...
    if !condition {
        break
    }
}
```

One keyword. Multiple uses. Nothing to memorize.

## Tradeoffs

| You gain | You give up |
|----------|-------------|
| One loop keyword (`for`) — simple mental model | Slightly more verbose for some patterns |
| No accidental fallthrough in switch | Must explicitly opt in when fallthrough is needed |
| If-init scopes variables tightly | Unfamiliar syntax for newcomers |
| Mandatory braces — no dangling-else bugs | More keystrokes for one-liners |
| No ternary — always readable | More lines for simple conditional assignment |
| Range handles iteration uniformly | Must remember range behavior varies by type |

## Why This Matters for Grep

Grep's core algorithm is almost entirely control flow:

```go
// This is the skeleton of grep, using only concepts from this document
lineNumber := 0
for scanner.Scan() {
    lineNumber++
    line := scanner.Text()

    if matchesPattern(line, pattern) {
        if invertMatch {
            continue  // -v flag: skip matches
        }
        fmt.Printf("%s:%d:%s\n", filename, lineNumber, line)
        matchCount++
    } else if invertMatch {
        fmt.Printf("%s:%d:%s\n", filename, lineNumber, line)
        matchCount++
    }
}
```

Every flag you add to your grep changes the control flow:
- `-c` (count only) — do not print lines, just count
- `-l` (files only) — print the filename and `break` to the next file on first match
- `-m N` (max count) — `break` after N matches
- `-v` (invert) — flip the condition
- `-q` (quiet) — do not print anything, just set exit code

Understanding control flow is understanding how to teach your program to make the right decision for every line of every file.

## Further Reading

1. **[A Tour of Go — Flow Control](https://go.dev/tour/flowcontrol/1)** — Official
   Interactive exercises covering `for`, `if`, `switch`, and `defer`. Hands-on practice.

2. **[Effective Go — Control Structures](https://go.dev/doc/effective_go#control-structures)** — Official
   Idiomatic patterns for using Go's control flow constructs, including the if-init pattern.

3. **[The Go Programming Language Specification — Statements](https://go.dev/ref/spec#Statements)** — Official
   The authoritative reference for exactly how each statement behaves, including all edge cases.

4. **[Go FAQ — Why does Go not have the ?: operator?](https://go.dev/doc/faq#Does_Go_have_a_ternary_form)** — Official
   The design rationale for omitting the ternary operator.

5. **[Go by Example — For, If/Else, Switch](https://gobyexample.com/for)** — Mark McGranaghan
   Concise, runnable examples of every control flow construct.

6. **[CodeReviewComments — Indent Error Flow](https://github.com/golang/go/wiki/CodeReviewComments#indent-error-flow)** — Go Wiki
   The official style guidance for keeping the happy path left-aligned.

7. **[Rob Pike — "Go Proverbs"](https://go-proverbs.github.io/)** — 2015
   "The bigger the interface, the weaker the abstraction" and other principles that inform Go's control flow design.

---

*Previous: [04 — Types and Memory](./04-types-and-memory.md) · Next: [06 — Functions](./06-functions.md)*
