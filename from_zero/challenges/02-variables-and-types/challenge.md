# Challenge 02: Variables and Types

## Prerequisites

- [Types and Memory](../../concepts/04-types-and-memory.md)
- Completed: [Challenge 01 — Hello, Go!](../01-hello-go/challenge.md)

## Overview

Variables are how programs remember things. Go is **statically typed** — every variable has a fixed type that's checked at compile time. This catches bugs before your program ever runs.

---

## Concepts

### Declaring Variables

Go gives you two ways to declare variables:

**Using `var` (explicit):**

```go
var name string = "Gopher"
var age int = 10
var height float64 = 1.75
var isAwesome bool = true
```

The `var` keyword declares a variable with a specific type. You can also let Go infer the type from the value:

```go
var name = "Gopher"    // Go infers: string
var age = 10           // Go infers: int
```

**Using `:=` (short declaration, inside functions only):**

```go
name := "Gopher"       // same as: var name string = "Gopher"
age := 10              // same as: var age int = 10
height := 1.75         // same as: var height float64 = 1.75
isAwesome := true      // same as: var isAwesome bool = true
```

The `:=` operator declares AND assigns in one step. It only works inside functions.

### Basic Types

| Type | Description | Example |
|---|---|---|
| `string` | Text | `"Gopher"` |
| `int` | Whole numbers | `10`, `-3`, `0` |
| `float64` | Decimal numbers | `1.75`, `3.14` |
| `bool` | True or false | `true`, `false` |

### Formatting Output

`fmt.Printf` gives you control over formatting (note: no automatic newline — you add `\n` yourself):

```go
fmt.Printf("Name: %s\n", name)        // %s for strings
fmt.Printf("Age: %d\n", age)          // %d for integers
fmt.Printf("Height: %.2f\n", height)  // %.2f for floats (2 decimal places)
fmt.Printf("Active: %t\n", isActive)  // %t for booleans
```

Common format verbs:

| Verb | Type | Example Output |
|---|---|---|
| `%s` | string | `Gopher` |
| `%d` | int | `10` |
| `%f` | float64 | `1.750000` |
| `%.2f` | float64 (2 decimals) | `1.75` |
| `%t` | bool | `true` |
| `%v` | any (default format) | varies |

`fmt.Println` is simpler — it uses default formatting and adds a newline:

```go
fmt.Println("Name:", name)  // prints: Name: Gopher
```

For more depth on types and how they're stored in memory, see [Types and Memory](../../concepts/04-types-and-memory.md).

---

## Task

Create a Go program that declares variables of different types and prints them in a specific format.

### Setup

1. `go mod init variables`
2. Create `main.go`

### Requirements

Your program must print exactly this output:

```
Name: Gopher
Age: 10
Height: 1.75
IsAwesome: true
```

**Rules:**

- You must declare at least one variable using `var`
- You must declare at least one variable using `:=`
- Height must be printed with exactly 2 decimal places
- The output must match exactly (spacing, capitalization, values)

### Verify Your Solution

```bash
bash test.sh
go test -v
```

---

## Hints

<details>
<summary>Hint 1: Which format verb for the height?</summary>

Use `fmt.Printf("Height: %.2f\n", height)` to get exactly 2 decimal places. `%f` alone would give you `1.750000`.

</details>

<details>
<summary>Hint 2: Mixing var and :=</summary>

You could do something like:
```go
var name string = "Gopher"
age := 10
```
As long as you use each style at least once.

</details>

<details>
<summary>Hint 3: Println vs Printf for booleans</summary>

`fmt.Println("IsAwesome:", true)` prints `IsAwesome: true` — that works. Or use `fmt.Printf("IsAwesome: %t\n", isAwesome)`.

</details>

---

## Further Reading

- [Go by Example: Variables](https://gobyexample.com/variables)
- [Go by Example: Constants](https://gobyexample.com/constants)
- [fmt package — format verbs](https://pkg.go.dev/fmt#hdr-Printing)
