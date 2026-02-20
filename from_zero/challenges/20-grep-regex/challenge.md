# Challenge 20: Grep - Regular Expressions

**This is a milestone.** You're moving from literal substring matching to regex — the heart of what makes grep, *grep* (Global Regular Expression Print).

## Prerequisites

- Complete [Challenge 19](../19-grep-invert/challenge.md)
- **Read this carefully**: [Regular Expressions](../../concepts/17-regular-expressions.md)

## Goal

Replace your literal `strings.Contains` matching with Go's `regexp` package. The pattern argument is now treated as a regular expression.

## Your Task

Upgrade `mygrep` so patterns are interpreted as regular expressions. Support:

| Pattern | Meaning | Example |
|---------|---------|---------|
| `\d` | Any digit (0-9) | `\d` matches "abc**1**def" |
| `\w` | Any word char (letter, digit, _) | `\w` matches "**a**!@#" |
| `.` | Any single character | `a.c` matches "a**b**c" |
| `*` | Zero or more of previous | `ab*c` matches "ac", "abc", "abbc" |
| `+` | One or more of previous | `ab+c` matches "abc", "abbc" but not "ac" |
| `?` | Zero or one of previous | `ab?c` matches "ac", "abc" |
| `[abc]` | Character class | `[BJ]` matches "**B**ob" or "**J**im" |
| `[a-z]` | Character range | `[0-9]` matches any digit |

### Test Cases from the Challenge Spec

```bash
# \d matches lines with digits
$ ./mygrep "\d" testdata/test-subdir/BFS1985.txt
Her dreams went out the door when she turned 24
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 19, 19, 1985, 1985
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 19, 19, 1985
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 1985
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 19, 19, 1985

# \w matches lines with word characters
$ ./mygrep "\w" testdata/symbols.txt
pound
dollar
```

### All Previous Features Still Work

- `-r` recursive
- `-v` invert
- Combined flags: `-rv`
- Stdin pipe
- Exit codes: 0/1/2

### Important: regexp.Compile vs regexp.MustCompile

- `regexp.Compile(pattern)` returns `(*Regexp, error)` — use this since patterns come from user input
- If the pattern is invalid regex, print an error and exit 2

```bash
$ ./mygrep "[invalid" testdata/rockbands.txt
error: invalid regex pattern: ...
$ echo $?
2
```

## How to Test

```bash
go test -v
bash test.sh
```

## Hints

<details><summary>Hint 1: Switching to regexp</summary>

Replace:
```go
if strings.Contains(line, pattern) {
```
With:
```go
re, err := regexp.Compile(pattern)
// handle err
if re.MatchString(line) {
```

</details>

<details><summary>Hint 2: Go's regexp and \d</summary>

Go's `regexp` package uses RE2 syntax. `\d` and `\w` work as expected. Empty pattern `""` compiles fine and matches everything (just like before).

</details>

<details><summary>Hint 3: Compile once</summary>

Compile the regex once before the loop, not on every line:
```go
re, err := regexp.Compile(pattern)
if err != nil { /* handle error */ }
// then in loop:
if re.MatchString(line) { ... }
```

</details>

## Why This Matters

After this challenge, your grep can do things like:
- `./mygrep "[0-9]{4}" file` — find 4-digit numbers
- `./mygrep "^[A-Z]" file` — lines starting with uppercase
- `./mygrep "\.txt$" file` — lines ending in .txt

You're building a real tool now.

---

Previous: [Challenge 19 - Invert](../19-grep-invert/challenge.md) | Next: [Challenge 21 - Anchors](../21-grep-anchors/challenge.md)
