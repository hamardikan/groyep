# Regular Expressions

## What Is It

A regular expression (regex) is a pattern language for matching text. It's a concise way to describe a set of strings. Instead of asking "does this line contain 'error'?", you can ask "does this line contain 'error' or 'Error' or 'ERROR' or any capitalization thereof?" — and you express that as a single pattern: `(?i)error`.

Regular expressions are the heart of grep. The name `grep` literally comes from the `ed` editor command `g/re/p` — "**g**lobally search for a **r**egular **e**xpression and **p**rint matching lines."

## What Problem It Solves

Text searching has a spectrum of complexity:

| Need                              | Solution          | Example                          |
|----------------------------------|-------------------|----------------------------------|
| Find exact text                  | String search     | `strings.Contains(line, "error")`|
| Find text with variations        | Regex             | `err(or)?` matches "err" or "error" |
| Find structured patterns         | Regex             | `\d{3}-\d{4}` matches phone numbers |
| Find complex nested structures   | Parser            | Regex can't parse HTML           |

Regex fills the middle ground: more powerful than literal string matching, simpler than writing a full parser. For the class of problems it handles — finding patterns in text — it's unmatched in expressiveness per character.

## First Principles

### History

1. **1956**: Stephen Kleene formalized regular expressions mathematically
2. **1968**: Ken Thompson implemented regex in the `ed` editor on Unix — this was the first practical use in software
3. **1973**: Thompson created `grep` as a standalone tool extracted from `ed`
4. **1979**: `egrep` added extended regex (alternation, quantifiers)
5. **1986**: Henry Spencer wrote a portable regex library
6. **1997**: Philip Hazel created PCRE (Perl-Compatible Regular Expressions)
7. **2007**: Russ Cox published his RE2 papers, leading to Go's regex engine

### The Chomsky Hierarchy (Brief)

Regular expressions describe **regular languages** — the simplest class in the Chomsky hierarchy of formal languages:

| Language Class    | Recognized By          | Example                    |
|------------------|------------------------|----------------------------|
| Regular          | Finite automaton       | `a*b+` (any a's followed by b's) |
| Context-Free     | Pushdown automaton     | Balanced parentheses `((()))` |
| Context-Sensitive| Linear bounded automaton| Some natural language grammar |
| Recursively Enumerable | Turing machine   | Anything computable        |

This is why regex **cannot parse HTML** or match balanced brackets — those are context-free languages, one level above what regex can handle. Knowing this boundary prevents you from misapplying regex.

### Basic Patterns

| Pattern  | Matches                          | Example        | Matches               |
|---------|----------------------------------|----------------|----------------------|
| `abc`   | Literal characters               | `hello`        | "hello"              |
| `.`     | Any single character             | `h.llo`        | "hello", "hallo"     |
| `*`     | Zero or more of previous         | `ab*c`         | "ac", "abc", "abbc"  |
| `+`     | One or more of previous          | `ab+c`         | "abc", "abbc" (not "ac") |
| `?`     | Zero or one of previous          | `colou?r`      | "color", "colour"    |
| `^`     | Start of line                    | `^Error`       | "Error: ..." (at start) |
| `$`     | End of line                      | `\.go$`        | "main.go" (at end)   |
| `\`     | Escape special character         | `\\.`          | Literal "."          |

### Character Classes

| Pattern     | Matches                       | Example                        |
|------------ |-------------------------------|-------------------------------|
| `[abc]`     | Any of a, b, or c            | `[aeiou]` matches any vowel  |
| `[a-z]`     | Any lowercase letter          | `[a-zA-Z]` any letter        |
| `[^abc]`    | Any character except a, b, c  | `[^0-9]` any non-digit       |
| `\d`        | Any digit (same as `[0-9]`)   | `\d{3}` three digits         |
| `\w`        | Word character `[a-zA-Z0-9_]` | `\w+` a word                 |
| `\s`        | Whitespace (space, tab, etc)  | `\s+` one or more spaces     |
| `\D`        | Non-digit                     | `\D+` non-digit sequence     |
| `\W`        | Non-word character            | `\W` punctuation, etc.       |
| `\S`        | Non-whitespace                | `\S+` non-whitespace run     |

### Quantifiers

| Quantifier | Meaning                  | Example          |
|-----------|--------------------------|------------------|
| `{n}`     | Exactly n times          | `\d{4}` four digits |
| `{n,}`    | At least n times         | `\d{2,}` two or more digits |
| `{n,m}`   | Between n and m times    | `\d{2,4}` two to four digits |
| `*`       | Same as `{0,}`           | Zero or more     |
| `+`       | Same as `{1,}`           | One or more      |
| `?`       | Same as `{0,1}`          | Zero or one      |

### Alternation and Groups

```
cat|dog          → matches "cat" or "dog"
(cat|dog) food   → matches "cat food" or "dog food"
(ab)+            → matches "ab", "abab", "ababab", ...
```

Groups `()` serve two purposes:
1. **Grouping** for alternation and quantifiers
2. **Capturing** for extracting matched substrings

### Greedy vs Lazy

By default, quantifiers are **greedy** — they match as much as possible:

```
Pattern: ".*"
Input:   He said "hello" and "goodbye"
Match:   "hello" and "goodbye"    (greedy — matches from first " to last ")
```

Adding `?` makes them **lazy** — match as little as possible:

```
Pattern: ".*?"
Input:   He said "hello" and "goodbye"
Match:   "hello"                  (lazy — stops at first closing ")
```

## How Go Does It

### The `regexp` Package

Go uses the **RE2** engine, created by Russ Cox (who also works on Go). RE2 guarantees **linear time** matching — no matter what the pattern or input, it will never take exponentially long.

```go
package main

import (
    "fmt"
    "regexp"
)

func main() {
    // Compile a pattern (returns error if invalid)
    re, err := regexp.Compile(`\d{3}-\d{4}`)
    if err != nil {
        fmt.Println("Bad pattern:", err)
        return
    }

    // Test if a string matches
    fmt.Println(re.MatchString("Call 555-1234"))    // true
    fmt.Println(re.MatchString("No number here"))   // false

    // Find the first match
    fmt.Println(re.FindString("Call 555-1234 or 555-5678"))  // "555-1234"

    // Find all matches
    fmt.Println(re.FindAllString("Call 555-1234 or 555-5678", -1))
    // ["555-1234", "555-5678"]
}
```

### `Compile` vs `MustCompile`

```go
// Compile: returns an error (use when pattern comes from user input)
re, err := regexp.Compile(userPattern)
if err != nil {
    fmt.Fprintf(os.Stderr, "groyep: invalid regex: %v\n", err)
    os.Exit(2)
}

// MustCompile: panics on error (use for hard-coded patterns)
var emailPattern = regexp.MustCompile(`[\w.+-]+@[\w-]+\.[\w.]+`)
```

For grep, the pattern comes from the user, so **always use `Compile`** and handle the error gracefully.

### Key Methods

```go
re := regexp.MustCompile(`error|warning`)

// Does it match?
re.MatchString("error: file not found")     // true

// Find the match
re.FindString("error: file not found")      // "error"

// Find match position [start, end]
re.FindStringIndex("error: file not found") // [0, 5]

// Find all matches
re.FindAllString("error and warning", -1)   // ["error", "warning"]

// Find submatch groups
re2 := regexp.MustCompile(`(\w+):(\d+)`)
re2.FindStringSubmatch("file.go:42")        // ["file.go:42", "file.go", "42"]

// Replace
re.ReplaceAllString("error and warning", "ISSUE") // "ISSUE and ISSUE"
```

### Method Naming Convention

Go's regexp methods follow a naming pattern:

| Component    | Meaning                          |
|-------------|----------------------------------|
| `Find`      | Find a match                     |
| `All`       | Find all matches (not just first)|
| `String`    | Input and output are strings     |
| `Submatch`  | Include capture group results    |
| `Index`     | Return position, not text        |

So `FindAllStringSubmatchIndex` means: find all matches, in a string, with capture groups, returning positions.

### Case-Insensitive Matching

```go
// Method 1: Flag in the pattern
re := regexp.MustCompile(`(?i)error`)
re.MatchString("Error")  // true
re.MatchString("ERROR")  // true
re.MatchString("error")  // true

// Method 2: Convert input to lowercase (for literal matching)
pattern := strings.ToLower("Error")
line := strings.ToLower("ERROR occurred")
strings.Contains(line, pattern) // true
```

For grep's `-i` flag, the `(?i)` prefix approach is cleaner because it works with regex patterns too:

```go
func compilePattern(pattern string, ignoreCase bool) (*regexp.Regexp, error) {
    if ignoreCase {
        pattern = "(?i)" + pattern
    }
    return regexp.Compile(pattern)
}
```

### RE2: Why Go's Regex Is Different

Most regex engines (Perl, Python, Java, JavaScript) use **backtracking**. Go uses **RE2**, which uses finite automata. The difference matters:

| Property              | Backtracking (Perl, Python) | RE2 (Go)                   |
|----------------------|----------------------------|-----------------------------|
| Time complexity      | Exponential (worst case)   | Linear (guaranteed)         |
| Backreferences       | Supported (`\1`)           | Not supported               |
| Lookahead/lookbehind | Supported                  | Not supported               |
| Performance          | Fast for simple cases      | Consistently fast           |
| Catastrophic backtracking | Possible             | Impossible                  |

**Catastrophic backtracking** is when a malicious or careless regex causes exponential time:

```
Pattern: (a+)+$
Input:   aaaaaaaaaaaaaaaaaaaaaaaaaab
```

In Perl/Python, this takes seconds or minutes. In Go, it takes microseconds. For a tool like grep that takes user-supplied patterns, this safety guarantee is critical — you don't want a user's regex to hang your program.

The tradeoff: Go doesn't support backreferences (`\1`) or lookahead/lookbehind (`(?=...)`, `(?!...)`). For grep, this is almost never a problem — these features are rare in grep usage.

### When NOT to Use Regex

| Task                              | Use Instead                      |
|----------------------------------|----------------------------------|
| Find an exact substring          | `strings.Contains`               |
| Check if string starts/ends with | `strings.HasPrefix`/`HasSuffix`  |
| Split on a simple delimiter      | `strings.Split`                  |
| Parse structured formats (JSON)  | A proper parser                  |
| Validate complex structures (HTML)| A proper parser                 |

Regex has overhead. For literal string matching (grep's most common use case), `strings.Contains` is faster. A production grep would optimize for the literal case:

```go
func matchLine(line string, pattern *regexp.Regexp, literal string, isLiteral bool) bool {
    if isLiteral {
        return strings.Contains(line, literal) // Fast path
    }
    return pattern.MatchString(line) // Regex path
}
```

## Tradeoffs

| Decision                        | Benefit                           | Cost                               |
|--------------------------------|-----------------------------------|------------------------------------|
| RE2 engine (Go's choice)      | Guaranteed linear time            | No backreferences or lookahead     |
| Backtracking engine (Perl)    | Full regex features               | Risk of exponential time           |
| Compile regex once             | Amortize compilation cost         | Need to pass compiled regex around |
| Compile regex per line         | Simpler code flow                 | Extremely slow (don't do this)     |
| Literal string optimization   | Much faster for simple patterns   | Two code paths to maintain         |
| Always use regex               | One code path, simple             | Slower for literal patterns        |

### Compile Once, Match Many

This is crucial for performance:

```go
// WRONG: compiles on every line — O(n * compilation_time)
for scanner.Scan() {
    matched, _ := regexp.MatchString(pattern, scanner.Text())
    // ...
}

// RIGHT: compile once, match many — O(compilation_time + n * match_time)
re, err := regexp.Compile(pattern)
if err != nil {
    // handle error
}
for scanner.Scan() {
    if re.MatchString(scanner.Text()) {
        // ...
    }
}
```

## Why This Matters for Grep

Regex is grep's core competency. The entire program exists to apply a regular expression to lines of text. Your grep clone needs to:

1. **Accept a pattern from the user**: The first non-flag argument
2. **Compile the pattern**: Using `regexp.Compile` (not `MustCompile` — handle errors)
3. **Apply the pattern to each line**: Using `MatchString` for basic grep
4. **Find match positions**: Using `FindStringIndex` for colorized output
5. **Handle invalid patterns gracefully**: Report the error, exit with code 2

```go
func main() {
    cfg, err := parseArgs()
    if err != nil {
        fmt.Fprintf(os.Stderr, "groyep: %v\n", err)
        os.Exit(2)
    }

    // Compile the pattern
    pattern := cfg.Pattern
    if cfg.IgnoreCase {
        pattern = "(?i)" + pattern
    }

    re, err := regexp.Compile(pattern)
    if err != nil {
        fmt.Fprintf(os.Stderr, "groyep: invalid regex %q: %v\n", cfg.Pattern, err)
        os.Exit(2)
    }

    // Search
    found := false
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        line := scanner.Text()
        if re.MatchString(line) != cfg.InvertMatch {
            fmt.Println(line)
            found = true
        }
    }

    if !found {
        os.Exit(1)
    }
}
```

For colorized output, you need match positions:

```go
func colorize(line string, re *regexp.Regexp) string {
    // Find all match positions
    matches := re.FindAllStringIndex(line, -1)
    if len(matches) == 0 {
        return line
    }

    var result strings.Builder
    lastEnd := 0
    for _, match := range matches {
        start, end := match[0], match[1]
        result.WriteString(line[lastEnd:start])        // Text before match
        result.WriteString("\033[1;31m")                // Red bold
        result.WriteString(line[start:end])             // The match
        result.WriteString("\033[0m")                   // Reset
        lastEnd = end
    }
    result.WriteString(line[lastEnd:])                  // Text after last match
    return result.String()
}
```

## Further Reading

1. **[Russ Cox — "Regular Expression Matching Can Be Simple And Fast"](https://swtch.com/~rsc/regexp/regexp1.html)** — The paper that explains why Go uses RE2. Demonstrates how backtracking engines can be exponentially slow and how Thompson's NFA algorithm guarantees linear time. Essential reading for understanding Go's regex design.

2. **[Go `regexp` package documentation](https://pkg.go.dev/regexp)** — The official reference. Pay attention to the RE2 syntax link and the method naming conventions.

3. **[Go `regexp/syntax` package documentation](https://pkg.go.dev/regexp/syntax)** — Documents the exact regex syntax Go supports. This is your definitive reference for "does Go support X?"

4. **["Mastering Regular Expressions" by Jeffrey Friedl (3rd edition)](https://www.oreilly.com/library/view/mastering-regular-expressions/0596528124/)** — The most comprehensive book on regex ever written. Covers theory, engines, and practical use across languages. Even the Go-relevant chapters alone are worth the book.

5. **[Ken Thompson — "Regular Expression Search Algorithm" (CACM, 1968)](https://dl.acm.org/doi/10.1145/363347.363387)** — The original paper where Thompson describes his regex implementation. This is the algorithm that grep was built on. A piece of computing history.

6. **[regex101.com](https://regex101.com/)** — Interactive regex tester with real-time explanation of what each part of your pattern does. Select the "Golang" flavor for Go-compatible syntax.

7. **[RE2 syntax reference](https://github.com/google/re2/wiki/Syntax)** — The exact syntax that Go's regexp package supports. Bookmark this.

8. **[Russ Cox — "Regular Expression Matching in the Wild"](https://swtch.com/~rsc/regexp/regexp3.html)** — Part 3 of Cox's series, discussing real-world regex engine implementation. Goes deeper into the theory and practice behind RE2.

9. **[Julia Evans — "How does grep work?" (zine)](https://wizardzines.com/)** — Julia Evans makes complex topics approachable. Her zines on command-line tools and text processing are excellent visual learning resources.

10. **[Go Playground regex examples](https://go.dev/play/)** — Use the Go Playground to experiment with regex patterns interactively. No installation needed.
