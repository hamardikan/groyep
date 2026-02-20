# Processes and Exit Codes

## What Is It

A **process** is a running instance of a program. When you type `./groyep "hello" file.txt` and press Enter, the operating system creates a new process: it loads your compiled binary into memory, sets up the standard streams (stdin, stdout, stderr), and begins executing your `main()` function.

Every process has:
- A **Process ID (PID)** — a unique integer
- A **parent process** — the process that created it (usually your shell)
- **Standard streams** — stdin, stdout, stderr
- **Environment variables** — inherited key-value pairs
- An **exit code** — a number between 0 and 255, reported when the process terminates

The exit code is the process's final word. It's how a program tells the world whether it succeeded or failed. This small integer is the foundation of shell scripting and tool composition.

## What Problem It Solves

Programs need to communicate their result status to the system that ran them. Did the command succeed? Did it find what it was looking for? Did something go wrong?

Without a standardized way to report success or failure:
- Shell scripts couldn't make decisions (`if grep pattern file; then ...`)
- Pipelines couldn't short-circuit on failure (`set -e`)
- Build systems couldn't detect compilation errors
- CI/CD systems couldn't know if tests passed

Exit codes solve this with radical simplicity: **0 means success, anything else means failure.** That's enough to build an entire automation ecosystem.

## First Principles

### Process Lifecycle

```
Created  →  Running  →  Terminated
(fork)     (executing)  (exit code returned to parent)
```

1. **Created**: The parent process (usually a shell) creates a new process
2. **Running**: The process executes its code
3. **Terminated**: The process calls `exit()` with a code, or returns from `main()`

### Exit Code Conventions

| Exit Code | Meaning                              | Example                         |
|-----------|--------------------------------------|---------------------------------|
| 0         | Success                              | `grep` found a match            |
| 1         | General error / failure              | `grep` found no match           |
| 2         | Misuse of command (bad arguments)    | `grep` given invalid regex      |
| 126       | Command found but not executable     | Permission denied               |
| 127       | Command not found                    | Typo in command name            |
| 128+N     | Killed by signal N                   | 130 = killed by Ctrl+C (SIGINT) |

The exit code is a single byte (0-255). By convention, 0 is the only success code. Everything else is a flavor of failure.

### Grep's Specific Exit Codes

Grep defines its exit codes precisely:

| Exit Code | Meaning                                    |
|-----------|--------------------------------------------|
| 0         | One or more matches were found             |
| 1         | No matches were found (not an error)       |
| 2         | An error occurred (bad regex, file not found, etc.) |

This three-way distinction is important. "No match" is not the same as "error." A grep that finds no matches in a log file isn't broken — there's just nothing to report.

### How the Shell Reads Exit Codes

After a command finishes, the shell stores its exit code in the special variable `$?`:

```bash
grep "hello" file.txt
echo $?   # 0 if found, 1 if not found, 2 if error
```

This enables conditional execution:

```bash
# && runs the next command only if the previous one succeeded (exit 0)
grep "error" log.txt && echo "Errors found!"

# || runs the next command only if the previous one failed (exit non-zero)
grep "error" log.txt || echo "No errors found"

# if/then uses the exit code directly
if grep -q "error" log.txt; then
    echo "Found errors"
else
    echo "Clean log"
fi
```

### Signals

Signals are asynchronous notifications sent to a process. The most common ones:

| Signal   | Number | Trigger         | Default Action | Description                 |
|----------|--------|-----------------|----------------|-----------------------------|
| SIGINT   | 2      | Ctrl+C          | Terminate      | Interrupt from keyboard     |
| SIGTERM  | 15     | `kill <pid>`    | Terminate      | Polite termination request  |
| SIGKILL  | 9      | `kill -9 <pid>` | Terminate      | Forced kill (can't catch)   |
| SIGPIPE  | 13     | Broken pipe     | Terminate      | Writing to a closed pipe    |

When a process is killed by signal N, its exit code is 128 + N. So Ctrl+C (SIGINT = signal 2) gives exit code 130.

**SIGPIPE** is particularly relevant for grep. If you pipe grep's output to `head -1`, `head` closes the pipe after reading one line. When grep tries to write the next match, it gets SIGPIPE and terminates. This is normal and expected — it's how `grep pattern file | head -5` works efficiently.

### Environment Variables

Processes inherit environment variables from their parent. These are key-value pairs available to the program:

```bash
GREP_OPTIONS="--color=auto" grep "hello" file.txt
```

In Go, you access them with `os.Getenv`:

```go
colorMode := os.Getenv("GREP_COLOR")
```

Environment variables are how users configure tools without flags. But they're secondary to command-line arguments — flags should always take precedence.

## How Go Does It

### Exiting with a Code

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Success
    os.Exit(0)

    // This line never executes — os.Exit terminates immediately
    fmt.Println("unreachable")
}
```

**Important:** `os.Exit` does NOT run deferred functions. If you have `defer file.Close()`, it won't execute. This is why you should structure your code so that `os.Exit` is only called in `main()`, after all cleanup has happened.

### A Better Pattern: Return from Main

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    os.Exit(run())
}

func run() int {
    file, err := os.Open("data.txt")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return 2 // Error exit code
    }
    defer file.Close() // This WILL execute because we return, not os.Exit

    // ... do work ...

    if matchFound {
        return 0 // Match found
    }
    return 1 // No match
}
```

This pattern separates "exit code determination" from "program logic" and ensures deferred functions run.

### Handling Signals

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // Create a channel to receive signals
    sigChan := make(chan os.Signal, 1)

    // Register for SIGINT and SIGTERM
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Wait for a signal (in a real program, this would be in a goroutine)
    sig := <-sigChan
    fmt.Fprintf(os.Stderr, "\nReceived signal: %v\n", sig)

    // Clean up and exit
    os.Exit(130) // Convention for SIGINT
}
```

For a simple grep, you usually don't need explicit signal handling — Go's default behavior (terminate on SIGINT) is fine. But knowing it's there is valuable.

### Reading Environment Variables

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Get a single variable
    home := os.Getenv("HOME")
    fmt.Println("Home:", home)

    // Check if a variable is set (empty string could be intentional)
    value, exists := os.LookupEnv("GREP_COLOR")
    if !exists {
        fmt.Println("GREP_COLOR is not set")
    } else {
        fmt.Println("GREP_COLOR:", value)
    }
}
```

### Getting the Process ID

```go
pid := os.Getpid()
ppid := os.Getppid() // Parent process ID
```

## Tradeoffs

| Decision                         | Benefit                              | Cost                                    |
|---------------------------------|--------------------------------------|-----------------------------------------|
| Exit immediately with `os.Exit` | Clear, simple                        | Deferred functions don't run            |
| Return codes from `run()`       | Defers run, testable                 | Slightly more code                      |
| Three exit codes (0, 1, 2)      | Distinguishes "no match" from "error"| Callers must know the convention        |
| Two exit codes (0, non-zero)    | Simpler for callers                  | Can't distinguish expected vs unexpected|
| Catch signals explicitly        | Graceful cleanup                     | More complex code                       |
| Default signal handling         | Zero code, works fine                | No cleanup on interrupt                 |

### The `os.Exit` vs `return` Tension

This is worth understanding deeply:

```go
func main() {
    f, _ := os.Open("file.txt")
    defer f.Close()  // WARNING: won't run if os.Exit is called!

    // ... work ...

    os.Exit(1) // f.Close() is NEVER called
}
```

The Go team made `os.Exit` skip defers intentionally — it mirrors the C `exit()` behavior. The solution is the `run()` pattern shown above.

## Why This Matters for Grep

Exit codes are how grep communicates with the rest of the Unix ecosystem. They're not a minor detail — they're part of grep's **interface contract**.

```bash
# These common patterns depend on correct exit codes:

# "Do something if pattern is found"
if groyep -q "error" server.log; then
    send_alert
fi

# "Count files containing a pattern"
for f in *.log; do
    groyep -q "FATAL" "$f" && echo "$f"
done | wc -l

# "Fail the build if forbidden patterns exist"
groyep "TODO(hack)" src/*.go && exit 1

# "Process only if matches exist"
groyep "pattern" data.txt | process_matches
# (if no matches, process_matches gets empty input — that's fine)
```

Your grep clone needs these exit codes:

```go
const (
    ExitMatch   = 0 // At least one match was found
    ExitNoMatch = 1 // No matches found (not an error)
    ExitError   = 2 // Something went wrong
)
```

Getting this wrong breaks scripts. If your grep exits with 0 when it finds no matches, `if groyep ... ; then` always succeeds. If it exits with 2 instead of 1 for "no match," error-handling scripts will think something crashed.

The process model also explains why grep should write errors to **stderr**, not stdout:

```bash
# If grep writes an error message to stdout, it pollutes the data:
groyep "pattern" nonexistent.txt | wc -l
# Should print error to stderr and count should be 0, not 1
```

## Further Reading

1. **["Advanced Programming in the UNIX Environment" by W. Richard Stevens (3rd ed.)](https://www.apue.com/)** — Chapters 7-10 cover process control, signals, and process relationships in definitive detail. The bible of Unix systems programming.

2. **[Go `os` package documentation](https://pkg.go.dev/os)** — Covers `os.Exit`, `os.Getenv`, `os.Getpid`, and process-related functions. The official reference.

3. **[Bash Reference Manual — Exit Status](https://www.gnu.org/software/bash/manual/html_node/Exit-Status.html)** — How the shell interprets exit codes, the `$?` variable, and how `&&`/`||` work. Essential for understanding how your tool will be used.

4. **[grep man page — EXIT STATUS section](https://man7.org/linux/man-pages/man1/grep.1.html)** — The specification for grep's exit codes. Your clone should match this behavior exactly.

5. **[Go `os/signal` package documentation](https://pkg.go.dev/os/signal)** — How to handle signals in Go. The `Notify` function and signal channel patterns.

6. **["The Linux Programming Interface" by Michael Kerrisk — Chapter 25: Process Termination](https://man7.org/tlpi/)** — Deep dive into how processes terminate, exit codes, and wait status.

7. **[Go Blog — "Defer, Panic, and Recover"](https://go.dev/blog/defer-panic-and-recover)** — Understanding defer is critical for the `os.Exit` vs return discussion. This blog post covers the mechanics.

8. **[Shell Scripting: Expert Recipes for Linux, Bash, and More (Steve Parker)](https://www.amazon.com/Shell-Scripting-Expert-Recipes-Linux/dp/1118024486)** — Practical examples of how scripts use exit codes, which helps you understand why getting them right matters.

9. **[POSIX Process Environment specification](https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap08.html)** — The formal specification for environment variables, exit codes, and process behavior. Dense but authoritative.
