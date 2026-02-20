# 07 — Error Handling

## What Is It

Error handling is how a program detects, reports, and recovers from things going wrong. Files that do not exist. Network connections that drop. Invalid input. Permission denied. Every real program must deal with errors.

Go takes a distinctive approach: **errors are values**, not exceptions. There is no `try/catch`. There is no special control flow. A function that can fail returns an `error` as its last return value. The caller checks it. Explicitly. Every time.

This is a deliberate, philosophical choice — not a limitation.

## What Problem It Solves

Programs fail. The question is: **how does the programmer find out?**

There are two major schools of thought:

### Exceptions (Java, Python, C#, JavaScript)

```python
try:
    file = open("data.txt")
    data = file.read()
except FileNotFoundError:
    print("file not found")
except PermissionError:
    print("permission denied")
```

Problems with exceptions:
- Any function in the call stack might throw — you cannot tell from the signature
- Exceptions create invisible control flow — code jumps to a catch block somewhere
- Easy to forget to handle an exception — the program crashes
- Hard to know which exceptions a function might throw (checked exceptions in Java tried to solve this but added painful verbosity)

### Return Values (Go, C, Rust)

```go
file, err := os.Open("data.txt")
if err != nil {
    fmt.Fprintf(os.Stderr, "error: %v\n", err)
    os.Exit(1)
}
```

The error is right there in the return value. You cannot call `os.Open` without being confronted with the possibility of failure. The control flow is visible and explicit.

## First Principles

Go's error philosophy rests on a few principles:

1. **Errors are values** — an error is just a value of type `error`, like an `int` or a `string`. You can store it, pass it, inspect it, compose it.
2. **Explicit is better than implicit** — if a function can fail, its signature tells you. No hidden exceptions.
3. **Handle errors where they occur** — do not let them propagate silently up the call stack.
4. **Add context as errors propagate** — each layer should add information about what it was trying to do.

## How Go Does It

### The `error` Interface

In Go, `error` is a built-in interface:

```go
type error interface {
    Error() string
}
```

Any type that has an `Error() string` method satisfies this interface. The simplest implementation is a string:

```go
// errors.New creates a simple error from a string
err := errors.New("file not found")
fmt.Println(err)        // "file not found"
fmt.Println(err.Error()) // "file not found"
```

### The Fundamental Pattern

The pattern you will write hundreds of times:

```go
result, err := someFunction()
if err != nil {
    // handle the error
    return fmt.Errorf("doing X: %w", err)
}
// use result
```

This is `if err != nil` — the most common three words in Go. Some people find it verbose. Others find it comforting: you can see exactly where every error is handled.

### Creating Errors

#### `errors.New` — Simple Errors

```go
import "errors"

func validate(pattern string) error {
    if pattern == "" {
        return errors.New("empty pattern")
    }
    return nil  // nil means no error — success
}
```

#### `fmt.Errorf` — Formatted Errors

```go
import "fmt"

func openFile(path string) (*os.File, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("opening %s: %v", path, err)
    }
    return file, nil
}
```

The `%v` verb formats the error using its `Error()` method.

### Wrapping Errors with `%w`

Go 1.13 introduced error wrapping. The `%w` verb (instead of `%v`) **wraps** the original error, preserving it for later inspection:

```go
func searchFile(path, pattern string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("searching %s: %w", path, err)
    }
    defer file.Close()
    // ...
}
```

The difference between `%v` and `%w`:
- `%v` formats the error as a string — the original error is lost
- `%w` wraps the error — the original error is preserved and can be unwrapped

### `errors.Is` — Checking for Specific Errors

When you wrap errors, you need a way to check if a specific error is somewhere in the chain:

```go
_, err := searchFile("/etc/shadow", "password")
if errors.Is(err, os.ErrPermission) {
    fmt.Println("permission denied — try running as root")
} else if errors.Is(err, os.ErrNotExist) {
    fmt.Println("file does not exist")
} else if err != nil {
    fmt.Println("unexpected error:", err)
}
```

`errors.Is` walks the error chain (unwrapping each layer) looking for a match. Without wrapping, you would lose the ability to identify the root cause.

### `errors.As` — Extracting Error Types

Sometimes you need the actual error value, not just a match:

```go
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Printf("operation %s failed on path %s: %v\n",
        pathErr.Op, pathErr.Path, pathErr.Err)
}
```

`errors.As` walks the chain looking for an error that can be assigned to the target type.

### Sentinel Errors

A **sentinel error** is a predefined error value used for comparison:

```go
var (
    ErrNotFound    = errors.New("not found")
    ErrNoMatch     = errors.New("no match")
    ErrInvalidFlag = errors.New("invalid flag")
)

func search(pattern string) error {
    // ...
    if matchCount == 0 {
        return ErrNoMatch
    }
    return nil
}

// Caller
err := search("foo")
if errors.Is(err, ErrNoMatch) {
    os.Exit(1)  // grep convention: exit 1 = no match
}
```

Sentinel errors are exported package-level variables. By convention, they start with `Err`:
- `io.EOF` — end of file
- `os.ErrNotExist` — file does not exist
- `os.ErrPermission` — permission denied

### Custom Error Types

For richer errors, define a struct that implements the `error` interface:

```go
type SearchError struct {
    File    string
    Line    int
    Pattern string
    Err     error
}

func (e *SearchError) Error() string {
    return fmt.Sprintf("%s:%d: search for %q failed: %v",
        e.File, e.Line, e.Pattern, e.Err)
}

func (e *SearchError) Unwrap() error {
    return e.Err  // enables errors.Is and errors.As on the wrapped error
}
```

The `Unwrap()` method tells the `errors` package how to traverse the chain.

### When to Panic vs Return Error

Go has `panic` for truly unrecoverable situations:

```go
panic("unreachable code reached")
```

**Use `panic` for:**
- Programming bugs (index out of range, nil pointer dereference)
- Situations that should literally never happen in correct code
- During `init` when setup fails and the program cannot continue

**Use `return error` for:**
- Everything else
- File not found, network timeout, invalid input, permission denied
- Anything the **user** or **environment** can cause

Rule of thumb: if the caller can reasonably be expected to handle it, return an error. If it is a bug in the programmer's code, panic.

### The Error Handling Debate

Go's error handling is the most debated aspect of the language. The arguments:

| For Explicit Errors | Against Explicit Errors |
|--------------------|-----------------------|
| You can see every error path | `if err != nil` is repetitive |
| No hidden control flow | 50%+ of lines may be error handling |
| Errors are self-documenting | Easy to accidentally ignore (assign to `_`) |
| Forces you to think about failures | Wrapping adds boilerplate |
| Simple — no exception hierarchy | No stack traces by default |
| Composable — errors are values | Generic error handling is hard to factor out |

Regardless of where you fall in this debate, idiomatic Go uses explicit error returns. Learn the pattern, internalize it, and you will find it becomes second nature.

### Common Patterns

#### Handle Then Continue

The most common pattern — handle the error immediately:

```go
file, err := os.Open(path)
if err != nil {
    return nil, fmt.Errorf("opening file: %w", err)
}
// continue with file
```

#### Collect and Continue

Sometimes you want to continue despite errors:

```go
var errs []error
for _, file := range files {
    if err := searchFile(file, pattern); err != nil {
        errs = append(errs, err)
        continue  // skip this file, try the next
    }
}
```

#### Log and Exit

At the top level, when there is nothing more to do:

```go
func main() {
    if err := run(); err != nil {
        fmt.Fprintf(os.Stderr, "groyep: %v\n", err)
        os.Exit(2)
    }
}
```

## Tradeoffs

| You gain | You give up |
|----------|-------------|
| Every error path is visible | Verbosity — many lines of error checking |
| No hidden exceptions | Easy to forget to check (assign to `_`) |
| Errors compose like any value | No stack traces without extra work |
| Simple mental model | Can feel tedious for rapid prototyping |
| Forces thinking about failure | Repetitive boilerplate |
| No exception hierarchy to learn | Must build custom error types for complex cases |

## Why This Matters for Grep

Grep encounters errors constantly:

- **File not found** — `grep pattern nonexistent.txt` should print an error and exit with code 2
- **Permission denied** — `grep pattern /etc/shadow` should report the error, not crash
- **Binary file** — real grep detects binary files and handles them differently
- **Broken pipe** — `grep pattern huge-file.txt | head -1` closes the pipe early; grep must handle SIGPIPE
- **Invalid regex** — `grep '[invalid' file.txt` should report a clear error, not panic
- **Directory instead of file** — `grep pattern /tmp/` should recurse or error depending on flags

Every one of these is an error you handle with `if err != nil`. Your grep must:

```go
file, err := os.Open(path)
if errors.Is(err, os.ErrNotExist) {
    fmt.Fprintf(os.Stderr, "groyep: %s: No such file or directory\n", path)
    return 2  // exit code 2 = error
}
if errors.Is(err, os.ErrPermission) {
    fmt.Fprintf(os.Stderr, "groyep: %s: Permission denied\n", path)
    return 2
}
if err != nil {
    fmt.Fprintf(os.Stderr, "groyep: %s: %v\n", path, err)
    return 2
}
```

The real `grep` handles errors gracefully — it does not crash, it does not silently skip. It reports the problem, uses the correct exit code, and keeps going if there are more files to search. Your grep should too.

Error handling is not the glamorous part of building grep. But it is what separates a toy from a tool.

## Further Reading

1. **[Error Handling and Go](https://go.dev/blog/error-handling-and-go)** — Go Blog, 2011
   The foundational blog post on Go's error handling philosophy. Start here.

2. **[Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)** — Go Blog, 2019
   Introduces `errors.Is`, `errors.As`, and the `%w` wrapping verb. Essential reading.

3. **[Errors Are Values](https://go.dev/blog/errors-are-values)** — Rob Pike, 2015
   Rob Pike shows how treating errors as values enables elegant error handling patterns that reduce the `if err != nil` repetition.

4. **[Don't Just Check Errors, Handle Them Gracefully](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)** — Dave Cheney, 2016
   A GopherCon talk on error handling strategies including sentinel errors, error types, and opaque errors.

5. **[Effective Go — Errors](https://go.dev/doc/effective_go#errors)** — Official
   The official guide to idiomatic error handling in Go.

6. **[Go FAQ — Why does Go not have exceptions?](https://go.dev/doc/faq#exceptions)** — Official
   The design rationale straight from the Go team.

7. **[Go by Example — Errors](https://gobyexample.com/errors)** — Mark McGranaghan
   Concise, runnable examples of creating and handling errors.

8. **[Error Handling in Upspin](https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html)** — Rob Pike, 2017
   A real-world case study of error handling patterns in a large Go project.

9. **[The Go Programming Language — Chapter 5.4: Errors](https://www.gopl.io/)** — Donovan & Kernighan
   Comprehensive treatment of error handling strategies with practical examples.

---

*Previous: [06 — Functions](./06-functions.md) · Next: [08 — Collections](./08-collections.md)*
