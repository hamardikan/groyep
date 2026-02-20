# Challenge 24: Grep - Context Lines

## Prerequisites

- Complete [Challenge 23](../23-grep-color/challenge.md)

## Goal

Add context line flags: `-A` (after), `-B` (before), `-C` (both). This is the most algorithmically complex challenge in the series.

## Your Task

Add three new flags:

| Flag | Meaning | Example |
|------|---------|---------|
| `-A N` | Print N lines **after** each match | `-A 2` shows 2 lines after |
| `-B N` | Print N lines **before** each match | `-B 2` shows 2 lines before |
| `-C N` | Print N lines **before and after** (same as `-A N -B N`) | `-C 1` shows 1 before + 1 after |

### Group Separator

When showing context, separate **non-adjacent groups** of output with `--` on its own line:

```bash
$ ./mygrep -B 1 "1985" testdata/test-subdir/BFS1985.txt
'Cause she's still preoccupied with 19, 19, 1985, 1985
--
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 19, 19, 1985
--
...
```

### Overlapping Context

If context regions overlap, **merge them** — don't print duplicate lines:

```bash
# If matches are close together, their context merges into one block
$ ./mygrep -C 2 "Madonna" testdata/test-subdir/BFS1985.txt
# Lines around all Madonna occurrences, merged when overlapping
```

### Key Rules

1. Context lines are printed even if they don't match the pattern
2. `--` separator only between non-adjacent groups (not at start or end)
3. Overlapping context is merged (no duplicate lines)
4. Works with all other flags: `-r`, `-v`, `-i`, `--color`
5. With `-r`, the file prefix is on every line (both matching and context)

## How to Test

```bash
go test -v
bash test.sh
```

## Hints

<details><summary>Hint 1: Algorithm approach</summary>

One approach:
1. Read ALL lines into a slice first
2. Find all matching line indices
3. For each match, calculate the range of lines to show (match - B to match + A)
4. Merge overlapping ranges
5. Print ranges with `--` separators between non-adjacent ranges

</details>

<details><summary>Hint 2: Merging ranges</summary>

```go
type Range struct { start, end int }

// Sort ranges by start, then merge overlapping ones
func mergeRanges(ranges []Range) []Range {
    sort.Slice(ranges, func(i, j int) bool {
        return ranges[i].start < ranges[j].start
    })
    var merged []Range
    for _, r := range ranges {
        if len(merged) > 0 && r.start <= merged[len(merged)-1].end+1 {
            // Overlapping or adjacent — extend
            if r.end > merged[len(merged)-1].end {
                merged[len(merged)-1].end = r.end
            }
        } else {
            merged = append(merged, r)
        }
    }
    return merged
}
```

</details>

<details><summary>Hint 3: The separator</summary>

Print `--` between groups, not before the first or after the last:

```go
for i, rng := range ranges {
    if i > 0 {
        fmt.Println("--")
    }
    // print lines in rng
}
```

</details>

---

Previous: [Challenge 23 - Color](../23-grep-color/challenge.md) | Next: [Challenge 25 - Final](../25-grep-final/challenge.md)
