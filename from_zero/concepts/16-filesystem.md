# Filesystem

## What Is It

A filesystem is a hierarchical structure for organizing data on a storage device. Files hold data. Directories (folders) hold files and other directories. Together they form a tree:

```
/
├── home/
│   └── user/
│       ├── Documents/
│       │   └── notes.txt
│       └── code/
│           └── groyep/
│               ├── main.go
│               ├── search.go
│               └── testdata/
│                   └── sample.txt
└── tmp/
    └── output.log
```

Every file has a **path** — its address in the tree. Filesystem operations (open, read, write, delete, list) are how programs interact with persistent data.

## What Problem It Solves

Programs need to work with data that outlasts a single execution. The filesystem provides:

- **Persistence**: Data survives after the program exits
- **Organization**: Hierarchical structure for managing files
- **Naming**: Paths provide human-readable addresses
- **Access control**: Permissions determine who can read, write, execute
- **Sharing**: Multiple programs can access the same files

For grep specifically, the filesystem is one of two input sources (the other being stdin). Grep needs to open files, read their contents line by line, and optionally walk directory trees to find files to search.

## First Principles

### Absolute vs Relative Paths

| Type     | Starts From    | Example                          | Portable? |
|----------|---------------|----------------------------------|-----------|
| Absolute | Root of filesystem | `/home/user/code/main.go`    | No (OS-specific root) |
| Relative | Current directory  | `code/main.go` or `./main.go`| Yes (if structure is consistent) |

```bash
# Absolute: unambiguous, works from anywhere
grep "hello" /home/user/code/main.go

# Relative: depends on your current directory
grep "hello" main.go
grep "hello" ./main.go
grep "hello" ../other/file.txt
```

### File Permissions (Brief)

Unix files have three permission types for three categories:

| Permission | File Meaning       | Directory Meaning       |
|-----------|-------------------|------------------------|
| Read (r)  | View contents     | List contents          |
| Write (w) | Modify contents   | Create/delete files    |
| Execute (x)| Run as program   | Enter directory        |

| Category | Description          |
|----------|---------------------|
| Owner    | The file's creator   |
| Group    | Members of a group   |
| Other    | Everyone else        |

Displayed as: `-rwxr-xr--` (owner: rwx, group: r-x, other: r--)

For grep, permissions matter when: a file can't be read (permission denied), or a directory can't be entered during recursive search.

### Symlinks (Symbolic Links)

A symlink is a file that points to another file or directory:

```bash
ln -s /var/log/syslog ~/syslog_link
```

When grep encounters a symlink, it needs to decide: follow it or skip it? Real grep follows symlinks to files by default but not to directories during recursive search (unless `-R` is used instead of `-r`). This is a subtle but important design decision.

## How Go Does It

### Opening and Reading Files

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    // Open a file for reading
    file, err := os.Open("example.txt")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer file.Close() // Always close files when done

    // Read line by line with bufio.Scanner
    scanner := bufio.NewScanner(file)
    lineNum := 0
    for scanner.Scan() {
        lineNum++
        fmt.Printf("%d: %s\n", lineNum, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading file:", err)
    }
}
```

### Reading an Entire File

For small files, `os.ReadFile` reads everything into memory:

```go
data, err := os.ReadFile("config.txt")
if err != nil {
    log.Fatal(err)
}
// data is []byte — the entire file contents
fmt.Println(string(data))
```

For grep, this is usually the **wrong** approach. Files can be gigabytes. Use `bufio.Scanner` for line-by-line processing.

### Creating and Writing Files

```go
// Create a new file (truncates if exists)
file, err := os.Create("output.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

// Write to it
fmt.Fprintln(file, "first line")
fmt.Fprintln(file, "second line")

// Or write the whole thing at once
err = os.WriteFile("output.txt", []byte("contents\n"), 0644)
```

### The `defer file.Close()` Pattern

Every opened file must be closed. `defer` ensures this happens even if the function returns early due to an error:

```go
func processFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err // File never opened, nothing to close
    }
    defer file.Close() // Will run when processFile returns

    // ... process file ...
    // Even if we return an error here, file.Close() runs
    return nil
}
```

Without `defer`, you'd need to call `file.Close()` before every return statement — easy to forget, leading to resource leaks.

### Checking File Properties

```go
info, err := os.Stat("example.txt")
if err != nil {
    if os.IsNotExist(err) {
        fmt.Println("File does not exist")
    } else {
        fmt.Println("Error:", err)
    }
    return
}

fmt.Println("Name:", info.Name())         // "example.txt"
fmt.Println("Size:", info.Size())         // Size in bytes
fmt.Println("Mode:", info.Mode())         // Permissions
fmt.Println("IsDir:", info.IsDir())       // true if directory
fmt.Println("ModTime:", info.ModTime())   // Last modification time
```

### Path Manipulation with `filepath`

Never build paths by string concatenation. Use `filepath`:

```go
import "path/filepath"

// Join paths (handles OS separators correctly)
path := filepath.Join("home", "user", "code", "main.go")
// Linux/Mac: "home/user/code/main.go"
// Windows:   "home\user\code\main.go"

// Extract components
dir := filepath.Dir("/home/user/code/main.go")    // "/home/user/code"
base := filepath.Base("/home/user/code/main.go")   // "main.go"
ext := filepath.Ext("/home/user/code/main.go")     // ".go"

// Clean up messy paths
clean := filepath.Clean("/home/user/../user/./code") // "/home/user/code"

// Make absolute
abs, err := filepath.Abs("./main.go") // "/current/working/dir/main.go"

// Match against a pattern
matched, err := filepath.Match("*.go", "main.go") // true
```

### Glob Patterns

Find files matching a pattern:

```go
// Find all Go files in the current directory
files, err := filepath.Glob("*.go")
// ["main.go", "search.go", "output.go"]

// Find all test files
testFiles, err := filepath.Glob("*_test.go")
// ["search_test.go", "output_test.go"]

// Find Go files in subdirectories (one level)
files, err = filepath.Glob("*/*.go")
```

Note: `filepath.Glob` does not support `**` for recursive matching. For recursive search, use `filepath.WalkDir`.

### Walking Directory Trees

`filepath.WalkDir` (Go 1.16+) recursively visits every file and directory:

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    root := "."

    err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            // Handle permission errors gracefully
            fmt.Fprintf(os.Stderr, "groyep: %s: %v\n", path, err)
            return nil // Continue walking despite errors
        }

        // Skip hidden directories
        if d.IsDir() && d.Name()[0] == '.' && d.Name() != "." {
            return filepath.SkipDir
        }

        // Only process regular files
        if !d.IsDir() {
            fmt.Println(path)
        }

        return nil
    })

    if err != nil {
        fmt.Fprintf(os.Stderr, "error walking: %v\n", err)
    }
}
```

### `WalkDir` vs `Walk`

| Feature        | `filepath.Walk` (old)      | `filepath.WalkDir` (Go 1.16+)  |
|---------------|---------------------------|---------------------------------|
| Callback arg  | `os.FileInfo` (calls `Stat`) | `os.DirEntry` (lazy `Stat`)  |
| Performance   | Slower (stats every file)  | Faster (stats only when needed) |
| Skip directory| Return `filepath.SkipDir`  | Return `filepath.SkipDir`       |
| Recommended   | No (legacy)                | Yes                             |

`WalkDir` is faster because it doesn't call `os.Stat()` on every file. For recursive grep, where you may be visiting thousands of files, this difference matters.

### Special Return Values in WalkDir

```go
filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
    // Skip this directory and all its children
    if d.IsDir() && d.Name() == "vendor" {
        return filepath.SkipDir
    }

    // Stop walking entirely (Go 1.20+)
    if someCondition {
        return filepath.SkipAll
    }

    // Continue normally
    return nil

    // Any other error stops the walk
    // return fmt.Errorf("something went wrong")
})
```

### The `testdata` Convention

Go's tooling ignores directories named `testdata`:

```
groyep/
├── search.go
├── search_test.go
└── testdata/           # Ignored by go build, available to tests
    ├── empty.txt       # Empty file
    ├── simple.txt      # A few lines of text
    ├── binary.dat      # Binary file (should not match)
    ├── unicode.txt     # Unicode characters
    └── large/          # Subdirectory for recursive tests
        ├── file1.txt
        └── file2.txt
```

```go
func TestSearchFile(t *testing.T) {
    matches, err := SearchFile("hello", "testdata/simple.txt")
    if err != nil {
        t.Fatal(err)
    }
    // ... verify matches ...
}
```

This is the standard place for test fixtures. Use it.

## Tradeoffs

| Decision                           | Benefit                           | Cost                               |
|-----------------------------------|-----------------------------------|------------------------------------|
| Read file line-by-line (stream)   | Handles any file size             | More complex code                  |
| Read entire file (`ReadFile`)     | Simple code                       | Fails on large files (OOM)         |
| `filepath.WalkDir` (recursive)    | Finds all files automatically     | Slow for large trees, may follow symlinks |
| Explicit file list (arguments)    | User controls scope               | User must list every file          |
| Follow symlinks                   | Complete search                   | Risk of infinite loops             |
| Skip symlinks                     | Safe, predictable                 | May miss files the user expects    |
| Skip binary files                 | Cleaner output                    | Need heuristic to detect binary    |
| Search binary files               | Thorough                          | Output may be garbled              |

### Error Handling Philosophy

Real grep doesn't stop when one file fails. It reports the error and continues:

```bash
$ grep "hello" file1.txt nonexistent.txt file2.txt
file1.txt:hello world
grep: nonexistent.txt: No such file or directory
file2.txt:hello there
```

In Go:

```go
func searchFiles(pattern string, files []string) int {
    exitCode := 1 // No match (default)
    hasError := false

    for _, filename := range files {
        file, err := os.Open(filename)
        if err != nil {
            fmt.Fprintf(os.Stderr, "groyep: %v\n", err)
            hasError = true
            continue // Skip this file, keep going
        }

        found := searchReader(pattern, file, filename)
        file.Close()

        if found {
            exitCode = 0 // At least one match
        }
    }

    if hasError {
        return 2
    }
    return exitCode
}
```

## Why This Matters for Grep

Grep interacts with the filesystem in several ways:

1. **Opening files from arguments**: `groyep "pattern" file1.txt file2.txt`
2. **Reading file contents**: Line-by-line streaming
3. **Recursive directory search**: `groyep -r "pattern" ./src`
4. **Handling errors gracefully**: Permission denied, file not found
5. **Multi-file output format**: Prefixing lines with filenames

```go
// When searching a single file, no filename prefix:
// hello world

// When searching multiple files, prefix with filename:
// file1.txt:hello world
// file2.txt:hello there
```

The filesystem operations form the "outer layer" of grep — they determine what gets fed into the search logic. Getting this layer right means:

- Handling arbitrarily large files (streaming, not `ReadFile`)
- Continuing past errors (report and move on)
- Respecting file types (skip directories when not in recursive mode)
- Cross-platform paths (use `filepath.Join`, not string concatenation)

```go
func grepFile(cfg *Config, filename string) (bool, error) {
    file, err := os.Open(filename)
    if err != nil {
        return false, err
    }
    defer file.Close()

    return grepReader(cfg, file, filename)
}

func grepRecursive(cfg *Config, root string) (bool, error) {
    found := false

    err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            fmt.Fprintf(os.Stderr, "groyep: %s: %v\n", path, err)
            return nil
        }

        if d.IsDir() {
            return nil // Continue into subdirectories
        }

        matched, fileErr := grepFile(cfg, path)
        if fileErr != nil {
            fmt.Fprintf(os.Stderr, "groyep: %v\n", fileErr)
            return nil
        }

        if matched {
            found = true
        }
        return nil
    })

    return found, err
}
```

## Further Reading

1. **[Go `os` package documentation](https://pkg.go.dev/os)** — Covers `os.Open`, `os.Create`, `os.ReadFile`, `os.Stat`, and all file operations. The primary reference.

2. **[Go `path/filepath` package documentation](https://pkg.go.dev/path/filepath)** — Path manipulation, globbing, and directory walking. Essential for cross-platform file operations.

3. **[Go Blog — "Working with Errors in Go 1.13"](https://go.dev/blog/go1.13-errors)** — Error wrapping and `errors.Is`/`errors.As` are important when dealing with file system errors (e.g., checking for "not found" vs "permission denied").

4. **["Learn Go with Tests" — Reading Files chapter](https://quii.gitbook.io/learn-go-with-tests/)** — Practical, test-driven approach to file I/O in Go. Great for seeing how to test filesystem code.

5. **[Go by Example: Reading Files](https://gobyexample.com/reading-files)** — Concise examples of different file reading approaches in Go.

6. **[Go by Example: Directories](https://gobyexample.com/directories)** — Creating, listing, and walking directories.

7. **[Go by Example: File Paths](https://gobyexample.com/file-paths)** — Using the `filepath` package for portable path manipulation.

8. **["Advanced Programming in the UNIX Environment" by W. Richard Stevens — Chapters 3-4](https://www.apue.com/)** — The definitive reference on file I/O at the system call level. Explains what Go's `os` package is abstracting over.

9. **[Go `io/fs` package documentation (Go 1.16+)](https://pkg.go.dev/io/fs)** — The filesystem abstraction layer that `WalkDir` is built on. Understanding `fs.FS` is useful for testing (you can create in-memory filesystems).

10. **[The Linux Programming Interface — Chapter 14-18: File Systems](https://man7.org/tlpi/)** — Deep technical reference on how filesystems work, including inodes, directories, links, and permissions.
