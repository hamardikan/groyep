# Challenge 18: Grep - Recursive Directory Search

## Prerequisites

- Complete [Challenge 17](../17-grep-literal/challenge.md)
- [Filesystem](../../concepts/16-filesystem.md)

## Goal

Add the `-r` flag to search recursively through directories. This is where your grep becomes truly powerful — searching entire project trees.

## Your Task

Add the `-r` flag:

```bash
./mygrep -r <pattern> <path>
```

When `-r` is provided, if the path is a directory, walk through all files in it (and subdirectories) and search each one. Prefix each matching line with the relative file path.

### Expected Behavior

```bash
$ ./mygrep -r Nirvana testdata/
testdata/rockbands.txt:Nirvana
testdata/test-subdir/BFS1985.txt:Since Bruce Springsteen, Madonna, way before Nirvana
testdata/test-subdir/BFS1985.txt:On the radio was Springsteen, Madonna, way before Nirvana
testdata/test-subdir/BFS1985.txt:And bring back Springsteen, Madonna, way before Nirvana
testdata/test-subdir/BFS1985.txt:Bruce Springsteen, Madonna, way before Nirvana
```

### Rules

- With `-r` and a directory: search all files recursively, prefix output with `filepath:line`
- With `-r` and a file: just search that file (with prefix)
- Without `-r` and a single file: print matching lines (no prefix needed)
- Maintain line order within each file
- Exit 0 if any match found across all files, 1 if none

### Important: File Path Prefix

When searching multiple files (with `-r`), every matching line must be prefixed with the file path:

```
filepath:matching line content
```

## How to Test

```bash
go test -v
bash test.sh
```

## Hints

<details><summary>Hint 1: Walking directories</summary>

Use `filepath.WalkDir` (preferred over `filepath.Walk`):

```go
filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
    if d.IsDir() {
        return nil // skip directories themselves
    }
    // search this file
    return nil
})
```

</details>

<details><summary>Hint 2: Using the flag package</summary>

```go
recursive := flag.Bool("r", false, "recursive search")
flag.Parse()
args := flag.Args() // remaining args after flags
```

</details>

<details><summary>Hint 3: Checking if path is directory</summary>

```go
info, err := os.Stat(path)
if info.IsDir() {
    // walk it
} else {
    // search single file
}
```

</details>

---

Previous: [Challenge 17 - Literal Match](../17-grep-literal/challenge.md) | Next: [Challenge 19 - Invert](../19-grep-invert/challenge.md)
