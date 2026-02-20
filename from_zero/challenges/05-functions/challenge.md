# Challenge 05: Functions

## Prerequisites

- [Functions](../../concepts/06-functions.md)
- [Error Handling](../../concepts/07-error-handling.md)
- Completed: [Challenge 04 — Loops](../04-loops/challenge.md)

## Overview

Functions break your code into reusable, testable pieces. Combined with Go's error handling patterns, they let you build programs that handle real-world messiness gracefully. You'll build a calculator that takes command-line arguments and handles every edge case.

---

## Concepts

### Defining Functions

A Go function has a name, parameters, and return types:

```go
func add(a float64, b float64) float64 {
    return a + b
}
```

When consecutive parameters share a type, you can shorten:

```go
func add(a, b float64) float64 {
    return a + b
}
```

### Multiple Return Values

Go functions can return multiple values. This is most commonly used for returning a result AND an error:

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```

Call it like this:

```go
result, err := divide(10, 3)
if err != nil {
    fmt.Println("error:", err)
    return
}
fmt.Printf("Result: %.2f\n", result)
```

### Creating Errors

Use `fmt.Errorf` to create error values with formatted messages:

```go
fmt.Errorf("unknown operator '%s'", op)
fmt.Errorf("invalid number '%s'", s)
```

Or use `errors.New` for simple error messages:

```go
import "errors"
errors.New("division by zero")
```

### Converting Strings to Floats

For this challenge, you'll need `strconv.ParseFloat` instead of `strconv.Atoi`:

```go
num, err := strconv.ParseFloat("3.14", 64)
// 64 means 64-bit precision (float64)
```

### Formatting Floats

To print a float with exactly 2 decimal places:

```go
fmt.Printf("Result: %.2f\n", 10.0)  // prints: Result: 10.00
```

For more on functions, see [Functions](../../concepts/06-functions.md). For error handling patterns, see [Error Handling](../../concepts/07-error-handling.md).

---

## Task

Create a calculator program that performs basic arithmetic on two numbers.

### Setup

1. `go mod init calc`
2. Create `main.go`

### Requirements

Usage: `./calc <num1> <operator> <num2>`

Supported operators: `add`, `sub`, `mul`, `div`

| Input | Output |
|---|---|
| `./calc 10 add 5` | `Result: 15.00` |
| `./calc 10 sub 3` | `Result: 7.00` |
| `./calc 4 mul 2.5` | `Result: 10.00` |
| `./calc 7 div 2` | `Result: 3.50` |
| `./calc 10 div 0` | `error: division by zero` |
| `./calc 10 pow 2` | `error: unknown operator 'pow'` |
| `./calc abc add 5` | `error: invalid number 'abc'` |
| `./calc 5 add xyz` | `error: invalid number 'xyz'` |
| `./calc` (wrong arg count) | `usage: calc <num1> <operator> <num2>` |
| `./calc 1 2` (wrong arg count) | `usage: calc <num1> <operator> <num2>` |

**Rules:**

- Results always formatted with 2 decimal places: `Result: X.XX`
- Error messages must match exactly
- The unknown operator error must include the operator name in single quotes
- The invalid number error must include the invalid value in single quotes
- Print the usage message for any argument count other than exactly 3

### Verify Your Solution

```bash
bash test.sh
go test -v
```

---

## Hints

<details>
<summary>Hint 1: Structuring your code</summary>

Consider writing a `calculate` function that takes two floats and an operator string, returning `(float64, error)`:

```go
func calculate(a float64, op string, b float64) (float64, error) {
    switch op {
    case "add":
        return a + b, nil
    case "sub":
        // ...
    default:
        return 0, fmt.Errorf("unknown operator '%s'", op)
    }
}
```

</details>

<details>
<summary>Hint 2: Parsing arguments</summary>

```go
if len(os.Args) != 4 {
    fmt.Println("usage: calc <num1> <operator> <num2>")
    return
}
num1, err := strconv.ParseFloat(os.Args[1], 64)
if err != nil {
    fmt.Printf("error: invalid number '%s'\n", os.Args[1])
    return
}
```

</details>

<details>
<summary>Hint 3: Switch vs if/else</summary>

Go's `switch` statement is cleaner than multiple if/else for matching operator names:

```go
switch op {
case "add":
    return a + b, nil
case "sub":
    return a - b, nil
case "mul":
    return a * b, nil
case "div":
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
default:
    return 0, fmt.Errorf("unknown operator '%s'", op)
}
```

</details>

---

## Further Reading

- [Go by Example: Functions](https://gobyexample.com/functions)
- [Go by Example: Multiple Return Values](https://gobyexample.com/multiple-return-values)
- [Go by Example: Errors](https://gobyexample.com/errors)
- [strconv.ParseFloat documentation](https://pkg.go.dev/strconv#ParseFloat)
