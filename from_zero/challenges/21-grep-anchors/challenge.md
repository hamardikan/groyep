# Challenge 21: Grep - Anchors

## Prerequisites

- Complete [Challenge 20](../20-grep-regex/challenge.md)
- [Regular Expressions](../../concepts/17-regular-expressions.md) (anchors section)

## Goal

Ensure your regex implementation correctly handles **anchors** — `^` (start of line) and `$` (end of line).

If you used Go's `regexp` package in challenge 20, anchors likely already work. This challenge is about **understanding** them and **verifying** they work correctly.

## Your Task

Verify and test that these anchor patterns work:

### Test Cases from the Spec

```bash
# ^A — lines starting with A
$ ./mygrep "^A" testdata/rockbands.txt
AC/DC
Aerosmith
Accept
April Wine
Autograph

# na$ — lines ending with "na"
$ ./mygrep "na$" testdata/rockbands.txt
Nirvana
```

### Additional Patterns to Support

| Pattern | Meaning | Example Match |
|---------|---------|---------------|
| `^A` | Line starts with A | AC/DC, Aerosmith |
| `na$` | Line ends with "na" | Nirvana |
| `^$` | Empty line | (blank lines only) |
| `^.{3}$` | Line with exactly 3 chars | UFO, Kix, TNT |
| `^\w` | Line starts with a word char | (most lines) |

### Understanding Anchors

Anchors are **zero-width assertions**. They don't match a character — they match a *position*:

- `^` asserts "the beginning of the line is here"
- `$` asserts "the end of the line is here"

This is why `^A` means "A at the start" and `A` means "A anywhere".

### All Previous Features Still Work

All flags from previous challenges: `-r`, `-v`, combined.

```bash
# Lines NOT starting with a capital letter
$ ./mygrep -v "^[A-Z]" testdata/symbols.txt

# Recursively find empty lines
$ ./mygrep -r "^$" testdata/
```

## How to Test

```bash
go test -v
bash test.sh
```

## Hints

<details><summary>Hint 1: It might already work</summary>

If you used `regexp.Compile(pattern)` and `re.MatchString(line)` in challenge 20, anchors work automatically. Go's regex engine handles `^` and `$` correctly.

The purpose of this challenge is to verify and understand, not necessarily add new code.

</details>

<details><summary>Hint 2: Multi-line mode</summary>

By default, Go's `regexp` treats `^` and `$` as start/end of the entire string. Since you're matching one line at a time (with `bufio.Scanner`), this works correctly — each line IS the entire string.

If you were matching against a whole file at once, you'd need the `(?m)` flag.

</details>

---

Previous: [Challenge 20 - Regex](../20-grep-regex/challenge.md) | Next: [Challenge 22 - Case Insensitive](../22-grep-case/challenge.md)
