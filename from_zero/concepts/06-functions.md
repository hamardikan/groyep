# 06 — Functions

## What Is It

A **function** is a named block of code that performs a specific task. You define it once, and call it from anywhere. Functions are the fundamental unit of abstraction in programming — they let you give a name to a sequence of operations, hide the details, and reuse the logic.

In Go, functions are not just convenient — they are the primary way you structure your entire program. Go does not have classes or objects in the traditional sense. Functions (and methods on types) are how you organize behavior.

## What Problem It Solves

Without functions:
- You duplicate code every time you need the same logic
- Changing behavior means changing it everywhere
- Programs become impossible to read — a wall of sequential instructions
- You cannot test parts of your program in isolation
- You cannot give names to operations ("search" vs 50 lines of search logic)

Functions solve the **abstraction problem**: they let you think at a higher level. Instead of reasoning about bytes and loops, you can reason about "searchFile" and "printMatch."

## First Principles

A function has a **contract** defined by its signature:
- **Name** — what it is called
- **Parameters** — what it needs to do its job
- **Return values** — what it gives back

The function body is hidden from the caller. This is **encapsulation** — the caller does not need to know how the function works, only what it takes and what it returns.

```
     inputs           function           outputs
  ┌──────────┐    ┌──────────────┐    ┌──────────┐
  │ pattern   │───▶│              │───▶│ matched  │
  │ line      │    │   search()   │    │ (bool)   │
  └──────────┘    └──────────────┘    └──────────┘
```

## How Go Does It

### Basic Function

```go
func greet(name string) string {
    return "Hello, " + name
}
```

Anatomy:
- `func` — keyword that declares a function
- `greet` — the function name
- `(name string)` — parameters: name and type
- `string` — return type
- `{ ... }` — the function body

### Parameters

Parameters are the inputs to a function. In Go, the **type comes after the name**:

```go
func search(line string, pattern string, caseSensitive bool) bool {
    // ...
}
```

When consecutive parameters share a type, you can shorten:

```go
func search(line, pattern string, caseSensitive bool) bool {
    // ...
}
```

Parameters are **passed by value** in Go. The function gets a copy. Modifying a parameter inside the function does not affect the caller's variable:

```go
func double(x int) {
    x = x * 2  // modifies the copy, not the original
}

n := 5
double(n)
fmt.Println(n)  // still 5
```

To modify the caller's data, use a pointer (covered in later concepts) or return the new value.

### Multiple Return Values

This is one of Go's most distinctive features. Functions can return **more than one value**:

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

result, err := divide(10, 3)
if err != nil {
    fmt.Println("Error:", err)
} else {
    fmt.Println("Result:", result)
}
```

This is how Go handles errors — functions return the result **and** an error. No exceptions, no try/catch. The caller checks the error explicitly.

The convention is: **the error is always the last return value**.

```go
// Standard patterns you will see everywhere
file, err := os.Open("data.txt")
line, err := reader.ReadString('\n')
count, err := fmt.Fprintf(w, "matches: %d\n", n)
```

If you do not need one of the return values, use the blank identifier:

```go
_, err := fmt.Fprintf(w, "hello\n")  // ignore the byte count
```

### Named Return Values

You can name return values in the function signature:

```go
func searchFile(path, pattern string) (matches []string, err error) {
    file, err := os.Open(path)
    if err != nil {
        return  // returns current values: nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        if strings.Contains(scanner.Text(), pattern) {
            matches = append(matches, scanner.Text())
        }
    }
    err = scanner.Err()
    return  // returns matches, err
}
```

Named returns act as variables initialized to their zero values. A bare `return` returns whatever those variables currently hold.

Use named returns sparingly. They are useful for:
- Short functions where the names document the return values
- `defer` functions that need to modify the return value

Avoid them in long functions — a bare `return` far from the function signature is hard to follow.

### Variadic Functions

A function can accept a variable number of arguments using `...`:

```go
func searchAny(line string, patterns ...string) bool {
    for _, pattern := range patterns {
        if strings.Contains(line, pattern) {
            return true
        }
    }
    return false
}

// Call with any number of arguments
searchAny(line, "error")
searchAny(line, "error", "warning", "fatal")

// Spread a slice with ...
patterns := []string{"error", "warning", "fatal"}
searchAny(line, patterns...)
```

Inside the function, `patterns` is a `[]string` — a slice. The `...` is syntactic sugar for the caller.

`fmt.Println` is a variadic function — that is why it accepts any number of arguments.

### The `main` Function

Every Go executable has exactly one entry point:

```go
package main

func main() {
    // Execution starts here
    // No parameters, no return values
    // To exit with a code, use os.Exit(code)
}
```

`main` takes no arguments and returns nothing. Command-line arguments come from `os.Args`. Exit codes are set with `os.Exit()`.

### Function Scope

Variables declared inside a function are **local** to that function:

```go
func search(line, pattern string) bool {
    matched := strings.Contains(line, pattern)
    return matched
}

// matched does not exist out here
```

Variables declared at the package level are accessible to all functions in the package:

```go
package main

var verbose bool  // accessible to all functions in package main

func main() {
    verbose = true
    run()
}

func run() {
    if verbose {
        fmt.Println("verbose mode")
    }
}
```

Minimize package-level variables. They make it harder to reason about your code because any function might change them.

### Functions as First-Class Citizens

In Go, functions are values. You can assign them to variables, pass them as arguments, and return them from other functions:

```go
// Function assigned to a variable
matcher := func(line, pattern string) bool {
    return strings.Contains(line, pattern)
}

// Using the function variable
if matcher(line, pattern) {
    fmt.Println(line)
}
```

```go
// Function as a parameter
func searchLines(lines []string, pattern string, match func(string, string) bool) []string {
    var results []string
    for _, line := range lines {
        if match(line, pattern) {
            results = append(results, line)
        }
    }
    return results
}

// Call with a literal matcher
results := searchLines(lines, "error", strings.Contains)

// Call with a case-insensitive matcher
results := searchLines(lines, "error", func(line, pattern string) bool {
    return strings.Contains(strings.ToLower(line), strings.ToLower(pattern))
})
```

This is powerful for grep — you can swap out the matching strategy without changing the search loop.

### Closures

An anonymous function can capture variables from its enclosing scope:

```go
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

counter := makeCounter()
fmt.Println(counter())  // 1
fmt.Println(counter())  // 2
fmt.Println(counter())  // 3
```

The returned function "closes over" the `count` variable — it retains access even after `makeCounter` returns.

### Why No Function Overloading

In many languages, you can define multiple functions with the same name but different parameters:

```java
// Java — function overloading
void search(String line, String pattern) { ... }
void search(String line, String pattern, boolean caseSensitive) { ... }
```

Go does not allow this. One name, one function. The reasoning: if you need different behavior, give it a different name.

```go
func search(line, pattern string) bool { ... }
func searchCaseInsensitive(line, pattern string) bool { ... }
```

Or use options via a parameter:

```go
func search(line, pattern string, caseSensitive bool) bool { ... }
```

### The `init` Function

Each package can define one or more `init` functions that run automatically when the package is first imported:

```go
package matcher

var defaultPatterns []string

func init() {
    // Runs before main(), when matcher package is imported
    defaultPatterns = loadDefaultPatterns()
}
```

`init` takes no arguments and returns nothing. It runs before `main`. Multiple `init` functions in the same file run in the order they appear. Use it sparingly — implicit initialization is harder to debug than explicit.

### Defer

The `defer` keyword schedules a function call to run when the enclosing function returns:

```go
func searchFile(path, pattern string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()  // will run when searchFile returns, no matter how

    // ... search the file ...
    return nil
}
```

`defer` guarantees cleanup happens even if the function returns early due to an error. Multiple defers execute in **LIFO order** (last in, first out).

## Tradeoffs

| You gain | You give up |
|----------|-------------|
| Multiple return values — natural error handling | Must always check both values |
| Functions as values — flexible composition | Can be overused, making code hard to follow |
| No overloading — one name means one thing | Must create distinct names for variants |
| `defer` — reliable cleanup | Slight performance cost, runs in reverse order |
| Named returns — self-documenting | Bare returns can be confusing in long functions |
| Pass by value — safe, no surprise mutations | Must use pointers for large structs or when mutation is needed |

## Why This Matters for Grep

Your grep clone is structured as functions:

```go
func main()                                           // entry point — parse args, orchestrate
func searchFile(path, pattern string, opts Options) ([]Match, error)  // search one file
func searchStream(r io.Reader, pattern string, opts Options) ([]Match, error)  // search stdin
func matchLine(line, pattern string, opts Options) bool  // core matching logic
func printMatch(m Match, opts Options)                // format and print one result
```

Each function has a clear purpose, clear inputs, and clear outputs. You can test `matchLine` without a file. You can test `searchFile` without printing. You can swap `matchLine` for a regex version without touching the file reading logic.

Multiple return values mean every function that can fail returns `(result, error)`. No hidden exceptions. No surprise panics. You handle every error where it occurs.

Functions as values let you implement the `-e` flag (multiple patterns) by composing matchers, or implement `-i` (case insensitive) by wrapping the matcher function.

The function is your primary tool for managing complexity. Get comfortable with them.

## Further Reading

1. **[Effective Go — Functions](https://go.dev/doc/effective_go#functions)** — Official
   Covers multiple returns, named results, and defer. The essential reference.

2. **[A Tour of Go — Functions](https://go.dev/tour/basics/4)** — Official
   Interactive exercises on functions, multiple returns, and closures.

3. **[Go Blog — Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover)** — 2010
   The definitive explanation of `defer` semantics, including the three rules of defer.

4. **[Go by Example — Functions](https://gobyexample.com/functions)** — Mark McGranaghan
   Quick, runnable examples covering all function features.

5. **[Go FAQ — Why does Go not support overloading?](https://go.dev/doc/faq#overloading)** — Official
   The design rationale for disallowing function overloading.

6. **[The Go Programming Language — Chapter 5: Functions](https://www.gopl.io/)** — Donovan & Kernighan
   Comprehensive coverage of functions, including error handling strategies, anonymous functions, and variadic functions.

7. **[Dave Cheney — "Practical Go: Functional Options"](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)** — 2014
   A pattern for configuring functions with many optional parameters using closures. Advanced but extremely useful.

8. **[CodeReviewComments — Function Size](https://github.com/golang/go/wiki/CodeReviewComments#function-size)** — Go Wiki
   Community guidance on how big functions should be and when to extract new ones.

---

*Previous: [05 — Control Flow](./05-control-flow.md) · Next: [07 — Error Handling](./07-error-handling.md)*
