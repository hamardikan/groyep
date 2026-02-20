# Challenge 03: Conditionals

## Prerequisites

- [Control Flow](../../concepts/05-control-flow.md)
- Completed: [Challenge 02 — Variables and Types](../02-variables-and-types/challenge.md)

## Overview

Programs need to make decisions. Conditionals let your code take different paths depending on the data it encounters. In Go, you'll also learn how to read command-line arguments and convert strings to numbers.

---

## Concepts

### Command-Line Arguments

When you run `./classify 42`, the `42` is a command-line argument. In Go, you access these through `os.Args`:

```go
import "os"

// os.Args[0] is the program name ("./classify")
// os.Args[1] is the first argument ("42")
// os.Args[2] would be the second argument, etc.

// len(os.Args) tells you how many args there are (including the program name)
```

### Converting Strings to Integers

Arguments arrive as strings. To use them as numbers, convert with `strconv.Atoi`:

```go
import "strconv"

num, err := strconv.Atoi("42")  // num = 42, err = nil
num, err := strconv.Atoi("abc") // num = 0, err = non-nil error
```

`Atoi` returns **two values** — the number and an error. This is Go's way of handling operations that can fail. If `err` is not `nil`, the conversion failed.

### If/Else

Go's `if` statement doesn't use parentheses around the condition, but the braces are required:

```go
if num > 0 {
    fmt.Println("positive")
} else if num < 0 {
    fmt.Println("negative")
} else {
    fmt.Println("zero")
}
```

### Error Checking Pattern

Go uses a pervasive pattern for error handling:

```go
value, err := someFunction()
if err != nil {
    // handle the error
    fmt.Println("something went wrong")
    return
}
// use value safely here
```

You'll see this pattern hundreds of times in Go code. Get comfortable with it.

For more on control flow, see the [Control Flow](../../concepts/05-control-flow.md) concept doc.

---

## Task

Create a number classifier program that takes a single number as a command-line argument and categorizes it.

### Setup

1. `go mod init classify`
2. Create `main.go`

### Requirements

The program should handle these cases:

| Input | Output |
|---|---|
| `./classify 5` | `positive` |
| `./classify -3` | `negative` |
| `./classify 0` | `zero` |
| `./classify abc` | `error: not a number` |
| `./classify` (no args) | `usage: classify <number>` |

**Rules:**

- Exactly one argument expected (besides the program name)
- If the argument can't be parsed as an integer, print the error message
- Output must match exactly (lowercase, exact text)
- The program should exit with code 0 in all cases (don't call `os.Exit(1)`)

### Verify Your Solution

```bash
bash test.sh
go test -v
```

---

## Hints

<details>
<summary>Hint 1: Checking argument count</summary>

```go
if len(os.Args) < 2 {
    fmt.Println("usage: classify <number>")
    return
}
```

Remember: `os.Args[0]` is the program name, so you need at least 2 elements.

</details>

<details>
<summary>Hint 2: Converting the argument</summary>

```go
num, err := strconv.Atoi(os.Args[1])
if err != nil {
    fmt.Println("error: not a number")
    return
}
```

</details>

<details>
<summary>Hint 3: The complete logic flow</summary>

1. Check `len(os.Args)` — if < 2, print usage
2. Try `strconv.Atoi(os.Args[1])` — if error, print error
3. Check if num > 0, < 0, or == 0

</details>

---

## Further Reading

- [os.Args documentation](https://pkg.go.dev/os#pkg-variables)
- [strconv.Atoi documentation](https://pkg.go.dev/strconv#Atoi)
- [Go by Example: If/Else](https://gobyexample.com/if-else)
- [Go by Example: Command-Line Arguments](https://gobyexample.com/command-line-arguments)
