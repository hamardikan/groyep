# I/O and Streams

## What Is It

I/O (Input/Output) is how programs communicate with the outside world — reading data in, writing data out. In the Unix model, **everything is a file**. A file on disk, a network connection, your keyboard, your screen — they all look the same to a program. Data flows through these connections as **streams of bytes**, one after another, like water through a pipe.

Every Unix process starts with three streams already open:

| Stream   | File Descriptor | Go Variable  | Purpose              |
|----------|-----------------|--------------|----------------------|
| stdin    | 0               | `os.Stdin`   | Input (keyboard, pipe) |
| stdout   | 1               | `os.Stdout`  | Normal output        |
| stderr   | 2               | `os.Stderr`  | Error messages       |

These three streams are the foundation of the entire Unix tool ecosystem.

## What Problem It Solves

Without a standard I/O model, every program would need to know the specifics of every device it might talk to. Want to read from a file? One API. Read from the network? A different API. Read from the keyboard? Yet another API.

The Unix I/O model solves this by creating a **universal abstraction**: everything is a stream of bytes. A program that reads bytes doesn't need to know where those bytes come from. This enables:

- **Piping**: connecting one program's output to another's input
- **Redirection**: sending output to a file instead of the screen
- **Composability**: building complex behavior from simple tools
- **Testability**: feeding test input through stdin, capturing output from stdout

## First Principles

### File Descriptors

A file descriptor is just an integer — a handle the operating system gives your program to refer to an open I/O resource. When a process starts, the OS opens three for you:

```
fd 0 → stdin  (where input comes from)
fd 1 → stdout (where normal output goes)
fd 2 → stderr (where error messages go)
```

When you open a file, you get fd 3, then fd 4, and so on. The program never sees raw hardware — it just reads from and writes to file descriptors.

### Byte Streams

Data flows as a sequence of bytes. There's no inherent structure — no "lines," no "records." It's just bytes. If you want lines, you read bytes until you find a newline character (`\n`). If you want words, you read bytes until you find a space. The interpretation is up to the program.

### Buffered vs Unbuffered I/O

**Unbuffered I/O** sends every byte to the OS immediately. Writing "hello" means five separate system calls — five round-trips to the kernel. This is slow.

**Buffered I/O** collects bytes in a memory buffer and sends them in larger chunks. Writing "hello" puts five bytes into a buffer; the buffer gets flushed when it's full (or when you explicitly flush it). This is fast.

| Approach     | System Calls | Latency     | Use Case              |
|-------------|-------------|-------------|----------------------|
| Unbuffered  | Many        | Immediate   | Interactive prompts   |
| Buffered    | Few         | Delayed     | File processing       |

The tradeoff is latency vs throughput. Buffering increases throughput at the cost of delayed output.

## How Go Does It

### The Core Interfaces

Go's entire I/O system is built on two interfaces:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

That's it. Any type that has a `Read` method is a `Reader`. Any type that has a `Write` method is a `Writer`. This simplicity is what makes Go's I/O so powerful.

### Standard Streams

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Write to stdout (normal output)
    fmt.Fprintln(os.Stdout, "This is normal output")

    // Write to stderr (error messages)
    fmt.Fprintln(os.Stderr, "This is an error message")

    // fmt.Println writes to stdout by default
    fmt.Println("Also normal output")
}
```

### Reading with bufio.Scanner

The `bufio.Scanner` provides convenient line-by-line reading:

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        line := scanner.Text() // Get the current line (without \n)
        fmt.Println(line)      // Process it
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading input:", err)
    }
}
```

This reads from stdin line by line. If you pipe a file in (`cat file | ./program`), it reads the file. If you type on the keyboard, it reads your keystrokes. The code doesn't change.

### Reading from Files

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("example.txt")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer file.Close() // Ensure the file is closed when we're done

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }
}
```

Notice the pattern: `os.Open` returns an `*os.File`, which implements `io.Reader`. The `bufio.Scanner` accepts any `io.Reader`. So the same scanning code works for stdin and files — the abstraction is seamless.

### Connecting Streams with io.Copy

```go
package main

import (
    "io"
    "os"
)

func main() {
    // Copy everything from stdin to stdout
    // This is essentially the `cat` command
    io.Copy(os.Stdout, os.Stdin)
}
```

`io.Copy` reads from a `Reader` and writes to a `Writer` until there's nothing left. It handles buffering internally.

### Buffered Writing

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush() // Don't forget to flush!

    for i := 0; i < 1000; i++ {
        fmt.Fprintln(writer, "line", i)
    }
    // All those writes were buffered — flushed in batches
}
```

## Tradeoffs

| Decision                    | Benefit                           | Cost                              |
|----------------------------|-----------------------------------|-----------------------------------|
| Byte streams (no structure)| Universal — works with any data   | Program must parse structure      |
| Buffered I/O               | Much faster for bulk operations   | Output may be delayed             |
| Scanner (line-by-line)     | Simple, memory-efficient          | Can't handle very long lines (default 64KB limit) |
| Read entire file into memory | Simple code                     | Fails on large files              |
| Streaming (line-by-line)   | Handles any file size             | More complex code                 |

### Scanner Buffer Limits

`bufio.Scanner` has a default max line size of 64KB. For grep, this is usually fine — most text files have short lines. But if you need to handle very long lines:

```go
scanner := bufio.NewScanner(file)
scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB max line
```

## Why This Matters for Grep

Grep is fundamentally an I/O program. It:

1. **Reads input** — from stdin (pipe) or from files (arguments)
2. **Processes each line** — checks if it matches the pattern
3. **Writes matching lines** — to stdout
4. **Reports errors** — to stderr

The streaming model is essential: grep must handle files of any size. You cannot read a 10GB log file into memory. You read it line by line, check each line, and either print it or skip it. The line you just processed can be discarded — you never need it again.

```go
// The core of grep is this simple:
func grep(pattern *regexp.Regexp, reader io.Reader, writer io.Writer) {
    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        if pattern.MatchString(scanner.Text()) {
            fmt.Fprintln(writer, scanner.Text())
        }
    }
}
```

Because this function takes `io.Reader` and `io.Writer`, it works with:
- stdin → stdout (piped usage)
- file → stdout (normal usage)
- strings.Reader → bytes.Buffer (testing!)

This is the power of Go's I/O interfaces.

## Further Reading

1. **[Go `io` package documentation](https://pkg.go.dev/io)** — The official reference for `Reader`, `Writer`, and all the utility functions like `io.Copy`, `io.TeeReader`, etc.

2. **[Go `bufio` package documentation](https://pkg.go.dev/bufio)** — Covers `Scanner`, `Reader`, and `Writer` with buffering. Essential for line-by-line processing.

3. **"The Unix Programming Environment" by Brian Kernighan & Rob Pike (1984)** — The definitive book on the Unix I/O model. Chapter 1 alone will reshape how you think about programs. Still relevant 40 years later.

4. **[Rob Pike — "Lexical Scanning in Go" (YouTube)](https://www.youtube.com/watch?v=HxaD_trXwRE)** — A talk that demonstrates Go's I/O patterns in the context of building a scanner. Shows how `io.Reader` thinking leads to elegant designs.

5. **[Go by Example: Reading Files](https://gobyexample.com/reading-files)** — Concise, practical examples of different file reading patterns in Go.

6. **[Effective Go — I/O section](https://go.dev/doc/effective_go)** — The official guide on idiomatic Go I/O patterns.

7. **["Advanced Programming in the UNIX Environment" by W. Richard Stevens](https://www.apue.com/)** — Deep dive into Unix I/O at the system call level. Chapter 3 covers file descriptors and the I/O model in detail.

8. **[Go Blog — "io.Reader in depth"](https://medium.com/@matryer/golang-advent-calendar-day-seventeen-io-reader-in-depth-6f744bb4320b)** — Mat Ryer's exploration of the `io.Reader` interface and why it's the most important interface in Go.

9. **[Julia Evans — "How does cat work?"](https://jvns.ca/)** — Julia Evans writes excellent, approachable explanations of Unix fundamentals including I/O, file descriptors, and system calls.

10. **[The Linux Programming Interface — Chapter 4: File I/O](https://man7.org/tlpi/)** — Michael Kerrisk's comprehensive reference on Linux system-level I/O, covering file descriptors, open, read, write, and close in depth.
