# Challenge 22: Grep - Case Insensitive Search

## Prerequisites

- Complete [Challenge 21](../21-grep-anchors/challenge.md)

## Goal

Add the `-i` flag for case-insensitive matching. This completes the core grep feature set from the original challenge spec.

## Your Task

Add the `-i` flag:

```bash
./mygrep -i <pattern> <file>
```

### The Definitive Test from the Spec

```bash
# Case-sensitive: only uppercase A
$ ./mygrep A testdata/rockbands.txt | wc -l
8

# Case-insensitive: both upper and lower 'a'
$ ./mygrep -i A testdata/rockbands.txt | wc -l
58
```

### Additional Test Cases

```bash
# Case-insensitive with anchors
$ ./mygrep -i "^a" testdata/rockbands.txt
AC/DC
Aerosmith
Accept
April Wine
Autograph

# Case-insensitive recursive
$ ./mygrep -ri "nirvana" testdata/
testdata/rockbands.txt:Nirvana
testdata/test-subdir/BFS1985.txt:Since Bruce Springsteen, Madonna, way before Nirvana
...

# Case-insensitive invert
$ ./mygrep -iv "the" testdata/test-subdir/BFS1985.txt
# (lines not containing "the" or "The" or "THE" etc.)
```

### All Flags Combined

Your grep should now support these flags in any combination:

| Flag | Meaning |
|------|---------|
| `-r` | Recursive directory search |
| `-v` | Invert match |
| `-i` | Case-insensitive |

Combined: `-riv`, `-ri`, `-rv`, `-iv`, etc.

## How to Test

```bash
go test -v
bash test.sh
```

## Hints

<details><summary>Hint 1: The regex approach</summary>

Go's regexp supports the `(?i)` flag prefix for case-insensitive matching:

```go
if caseInsensitive {
    pattern = "(?i)" + pattern
}
re, err := regexp.Compile(pattern)
```

This is the cleanest approach and works with all regex features.

</details>

<details><summary>Hint 2: The strings approach (literal only)</summary>

For literal patterns, you could also do:
```go
strings.Contains(strings.ToLower(line), strings.ToLower(pattern))
```

But this breaks with regex patterns. The `(?i)` approach is better.

</details>

<details><summary>Hint 3: Testing the exact count</summary>

To verify your count matches:
```bash
# Use real grep to verify
grep -c "A" testdata/rockbands.txt     # should be 8
grep -ic "A" testdata/rockbands.txt    # should be 58
```

</details>

## Congratulations!

After this challenge, you've implemented the **core grep** from the challenge spec:
- Empty match, literal match, regex
- Recursive search (`-r`)
- Invert match (`-v`)
- Case-insensitive (`-i`)
- Correct exit codes
- Stdin pipe support
- Anchors (`^`, `$`)
- Character classes, quantifiers, `\d`, `\w`

The next 3 challenges add polish: color, context lines, and the final integration.

---

Previous: [Challenge 21 - Anchors](../21-grep-anchors/challenge.md) | Next: [Challenge 23 - Color](../23-grep-color/challenge.md)
