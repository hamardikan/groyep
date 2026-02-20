# Challenge 17: Grep - Literal Match

## Prerequisites

- Complete [Challenge 16](../16-grep-empty/challenge.md)
- [Processes and Exit Codes](../../concepts/12-processes-and-exit-codes.md)
- [I/O and Streams](../../concepts/10-io-and-streams.md)

## Goal

Make literal pattern matching robust and add **stdin pipe support**. Your grep should now feel like a real Unix tool — it reads from a file OR from a pipe.

## Your Task

Extend `mygrep` so that:

1. **Literal substring matching works correctly** with various patterns
2. **Exit codes are correct**: 0 for match found, 1 for no match, 2 for error
3. **Stdin support**: when no file is given, read from stdin (enables piping)

### Expected Behavior

```bash
# Single character match
$ ./mygrep J testdata/rockbands.txt
Judas Priest
Bon Jovi
Junkyard

# Multi-word match
$ ./mygrep "Bad " testdata/rockbands.txt
Bad English
Bad Company

# No match — exit code 1
$ ./mygrep "ZZZZZ" testdata/rockbands.txt
$ echo $?
1

# Stdin pipe support
$ cat testdata/rockbands.txt | ./mygrep "J"
Judas Priest
Bon Jovi
Junkyard

# Pipe chaining
$ cat testdata/rockbands.txt | ./mygrep "an"
Van Halen
Damn Yankees
Vandenberg
Hanoi Rocks
Danger Danger
Dangerous Toys
Bang Tango
```

### Exit Code Behavior

| Scenario | Exit Code |
|----------|-----------|
| Pattern found in file | 0 |
| Pattern not found | 1 |
| No args (no pattern) | 2 |
| File not found | 2 |
| Pattern + no file = read stdin | (read from pipe) |

### Key: Stdin Support

When the user provides a pattern but no filename, your grep should read from stdin. This enables Unix pipe composition:

```bash
# Two-arg form: pattern + file
./mygrep "pattern" file.txt

# One-arg form: pattern only, reads stdin
echo "hello world" | ./mygrep "hello"
```

## How to Test

```bash
go test -v
bash test.sh
```

## Hints

<details><summary>Hint 1: Detecting stdin vs file</summary>

Check `len(os.Args)`:
- 1 arg (just program name): show usage
- 2 args (pattern only): read from `os.Stdin`
- 3 args (pattern + file): read from file

</details>

<details><summary>Hint 2: Unified reading</summary>

Both `os.Stdin` and `os.File` implement `io.Reader`. You can use `bufio.NewScanner()` on either one. Create a helper function that takes an `io.Reader`.

</details>

<details><summary>Hint 3: Special characters in test data</summary>

Note that rockbands.txt contains bands like `AC/DC`, `Guns N' Roses`, `W.A.S.P.` — your literal match should handle these since `strings.Contains` works with any substring.

</details>

---

Previous: [Challenge 16 - Empty Match](../16-grep-empty/challenge.md) | Next: [Challenge 18 - Recursive](../18-grep-recursive/challenge.md)
