# Challenge 16: Grep - Empty Match

**You've arrived.** Everything you've learned so far — Go basics, file I/O, stdin, exit codes, CLI arguments, structs, packages — it all leads here. From this point forward, you're building **your own grep**.

## Prerequisites

Read these before starting:
- [I/O and Streams](../../concepts/10-io-and-streams.md)
- [Unix Philosophy](../../concepts/11-unix-philosophy.md)
- [Filesystem](../../concepts/16-filesystem.md)
- [Testing Philosophy](../../concepts/15-testing-philosophy.md)

## Goal

Build the first version of `mygrep` — the simplest possible grep.

From the grep man page:
> The grep utility searches any given input files, selecting lines that match one or more patterns. By default, a pattern matches an input line if the regular expression in the pattern matches the input line without its trailing newline. **An empty expression matches every line.**

## Your Task

Create a program called `mygrep` that:

1. Takes a pattern and a filename as arguments: `./mygrep <pattern> <file>`
2. Prints every line in the file that contains the pattern (literal substring match for now)
3. An **empty pattern** (`""`) matches every line — prints the entire file

### Expected Behavior

```bash
# Empty pattern — prints every line of the file
./mygrep "" testdata/test.txt

# Literal pattern — prints matching lines
./mygrep "Gutenberg" testdata/test.txt

# Verify empty match is identical to the file itself
./mygrep "" testdata/test.txt | diff testdata/test.txt -
# (no output means they're identical)
```

### Exit Codes

| Situation | Exit Code |
|-----------|-----------|
| Match found | 0 |
| No match found | 1 |
| Error (no args, file not found) | 2 |

### Error Handling

- No arguments: print `usage: mygrep <pattern> <file>` to stderr, exit 2
- File not found: print error to stderr, exit 2

## Project Structure

From this challenge forward, each challenge is its own module:

```
16-grep-empty/
├── testdata/          # provided (don't modify)
├── challenge.md       # this file
├── grep_test.go       # provided (don't modify)
├── test.sh            # provided (don't modify)
├── go.mod             # YOU create this
└── main.go            # YOU create this
```

Initialize your module:
```bash
go mod init mygrep
```

## How to Test

```bash
# Go tests
go test -v

# Shell tests
bash test.sh
```

## Hints

<details><summary>Hint 1: Reading the file</summary>

You've done this before in challenge 08. Use `os.Open` and `bufio.Scanner` to read line by line.

</details>

<details><summary>Hint 2: Matching</summary>

For now, `strings.Contains(line, pattern)` is all you need. An empty pattern is contained in every string.

</details>

<details><summary>Hint 3: Program structure</summary>

```
func main():
  1. Parse os.Args (need at least 2 args after program name)
  2. Open the file
  3. Read line by line
  4. If line contains pattern, print it and set a "found" flag
  5. Exit with appropriate code
```

</details>

## What You're Building Toward

This is step 1 of 10. By challenge 25, your grep will support:
- Regex patterns, recursive search, inverted matching
- Case-insensitive search, colored output, context lines
- All the flags real grep has

But right now: just make the empty match work. One step at a time.

---

Next: [Challenge 17 - Grep: Literal Match](../17-grep-literal/challenge.md)
