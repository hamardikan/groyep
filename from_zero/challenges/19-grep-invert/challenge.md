# Challenge 19: Grep - Invert Match

## Prerequisites

- Complete [Challenge 18](../18-grep-recursive/challenge.md)

## Goal

Add the `-v` flag to **invert** the match — print lines that do NOT match the pattern.

This is one of grep's most powerful features. Combined with pipes, it lets you filter out noise from results.

## Your Task

Add the `-v` flag:

```bash
./mygrep -v <pattern> <file>
```

When `-v` is provided, print lines that do **not** contain the pattern.

### Expected Behavior

```bash
# The classic pipe composition from the challenge spec:
$ ./mygrep -r Nirvana testdata/ | ./mygrep -v Madonna
testdata/rockbands.txt:Nirvana

# Invert on symbols — exclude lines containing 'o'
$ ./mygrep -v "o" testdata/symbols.txt
!
@
£
$
%
^
&
*
(
)
```

### Flag Combinations

`-v` should work with `-r`:

```bash
# Recursive search, then invert
$ ./mygrep -rv "the" testdata/test-subdir/BFS1985.txt
```

This prints all lines in BFS1985 that do NOT contain "the".

### Exit Codes with -v

| Scenario | Exit Code |
|----------|-----------|
| At least one line passed through (didn't match) | 0 |
| All lines matched (nothing to output) | 1 |
| Error | 2 |

## How to Test

```bash
go test -v
bash test.sh
```

## Hints

<details><summary>Hint 1: Simple logic flip</summary>

If your matching logic looks like:
```go
if strings.Contains(line, pattern) {
    // print line
}
```

With `-v`, just flip the condition:
```go
if invert != strings.Contains(line, pattern) {
    // print line
}
```

Or use XOR logic: `matched != invert`

</details>

<details><summary>Hint 2: Combined flags</summary>

With the `flag` package, each flag is independent:
```go
recursive := flag.Bool("r", false, "recursive")
invert := flag.Bool("v", false, "invert match")
flag.Parse()
```

Users can write `-rv`, `-r -v`, or `-v -r` — the flag package handles all forms.

</details>

---

Previous: [Challenge 18 - Recursive](../18-grep-recursive/challenge.md) | Next: [Challenge 20 - Regex](../20-grep-regex/challenge.md)
