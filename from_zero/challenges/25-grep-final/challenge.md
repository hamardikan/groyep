# Challenge 25: Grep - The Final Build

**Congratulations.** You've made it to the last challenge. Everything you've built across 25 challenges comes together here into one polished, complete grep clone.

## Prerequisites

- Complete all previous challenges (00-24)
- All concept documents in [concepts/](../../concepts/)

## Goal

Build the **final version** of `mygrep` — a complete, well-structured, production-quality grep tool.

## Complete Feature Set

| Flag | Meaning | Since |
|------|---------|-------|
| (none) | Search pattern in file | Ch 16 |
| (stdin) | Read from pipe when no file given | Ch 17 |
| `-r` | Recursive directory search | Ch 18 |
| `-v` | Invert match (exclude matching lines) | Ch 19 |
| `-i` | Case-insensitive matching | Ch 22 |
| `-n` | Prefix with line numbers | Ch 13 |
| `-c` | Print count of matches only | Ch 13 |
| `-A N` | Show N lines after match | Ch 24 |
| `-B N` | Show N lines before match | Ch 24 |
| `-C N` | Show N lines before and after | Ch 24 |
| `-color=MODE` | Colorize matches (always/never/auto) | Ch 23 |

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | At least one match found |
| 1 | No matches found |
| 2 | Error (bad args, bad regex, file not found) |

## The Complete Test Suite

These are the definitive tests from the original challenge spec. Your grep must pass **all** of them:

### Step 1: Empty Match
```bash
./mygrep "" testdata/test.txt | diff testdata/test.txt -
# (no diff = pass)
```

### Step 2: Literal Match + Exit Codes
```bash
./mygrep J testdata/rockbands.txt
# Judas Priest
# Bon Jovi
# Junkyard
echo $?  # 0
```

### Step 3: Recursive
```bash
./mygrep -r Nirvana testdata/
# testdata/rockbands.txt:Nirvana
# testdata/test-subdir/BFS1985.txt:Since Bruce Springsteen, Madonna, way before Nirvana
# testdata/test-subdir/BFS1985.txt:On the radio was Springsteen, Madonna, way before Nirvana
# testdata/test-subdir/BFS1985.txt:And bring back Springsteen, Madonna, way before Nirvana
# testdata/test-subdir/BFS1985.txt:Bruce Springsteen, Madonna, way before Nirvana
```

### Step 4: Invert
```bash
./mygrep -r Nirvana testdata/ | ./mygrep -v Madonna
# testdata/rockbands.txt:Nirvana
```

### Step 5: Regex
```bash
./mygrep "\d" testdata/test-subdir/BFS1985.txt
# (9 lines with digits)

./mygrep "\w" testdata/symbols.txt
# pound
# dollar
```

### Step 6: Anchors
```bash
./mygrep "^A" testdata/rockbands.txt
# AC/DC
# Aerosmith
# Accept
# April Wine
# Autograph

./mygrep "na$" testdata/rockbands.txt
# Nirvana
```

### Step 7: Case Insensitive
```bash
./mygrep A testdata/rockbands.txt | wc -l       # 8
./mygrep -i A testdata/rockbands.txt | wc -l    # 58
```

## Recommended Project Structure

For the final build, structure your code properly:

```
25-grep-final/
├── go.mod
├── cmd/
│   └── mygrep/
│       └── main.go          # Entry point: parse flags, call library
├── internal/
│   ├── matcher/
│   │   ├── matcher.go       # Core matching logic (regex, case-insensitive)
│   │   └── matcher_test.go  # Unit tests for matcher
│   ├── searcher/
│   │   ├── searcher.go      # File/stdin searching, recursive walking
│   │   └── searcher_test.go # Unit tests for searcher
│   └── output/
│       ├── formatter.go     # Output formatting (color, line numbers, context)
│       └── formatter_test.go
├── testdata/
└── test.sh
```

Build with:
```bash
go build -o mygrep ./cmd/mygrep
```

## How to Test

```bash
go test ./...
bash test.sh
```

## Edge Cases to Handle

- Empty files
- Files with no trailing newline
- Very long lines (> 10K chars)
- Binary files (skip or warn)
- Patterns that match empty strings
- `-c` with `-v` (count non-matching lines)
- `-n` with `-r` (line numbers per file)
- `-c` with `-r` (count per file)

## Going Beyond

You've built a working grep. Here's what you could do next:

### More Features
- `-l` — print only filenames with matches
- `-L` — print only filenames without matches
- `-w` — whole word match only
- `-x` — whole line match only
- `-m N` — stop after N matches
- `--include=PATTERN` — only search files matching glob
- `--exclude=PATTERN` — skip files matching glob

### Performance
- Benchmark your grep against real GNU grep on large files
- Profile with `go tool pprof`
- Try memory-mapped file reading
- Parallel file searching with goroutines

### The Hard Path: Build Your Own Regex Engine
The ultimate challenge. Instead of using Go's `regexp` package, implement your own:
- [Let's Build a Regex Engine](https://kean.blog/post/lets-build-regex)
- [Russ Cox: Regular Expression Matching Can Be Simple And Fast](https://swtch.com/~rsc/regexp/regexp1.html)
- [Ken Thompson's 1968 CACM paper](https://dl.acm.org/doi/10.1145/363347.363387)

### Share Your Work
Put it on GitHub. You built this from scratch — that's worth sharing.

---

**You started with `fmt.Println("Hello, World!")` and ended with a full Unix tool. Well done.**

---

Previous: [Challenge 24 - Context](../24-grep-context/challenge.md) | Back to [Roadmap](../../README.md)
