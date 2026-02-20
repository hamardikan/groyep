# 04 — Types and Memory

## What Is It

A **type** tells the compiler (and the programmer) what kind of data a value is and what operations are valid on it. When you say a variable has type `int`, you are saying: this is a whole number, you can add it, subtract it, compare it, but you cannot call methods on it like a string and you cannot use it as a boolean.

A **type system** is the set of rules a language uses to assign and check types. Go's type system is **static** — types are determined at compile time, before your program runs.

## What Problem It Solves

Without types, the computer sees everything as raw bytes. The number `65` and the character `A` are the same byte (`01000001`). Types give meaning to bytes. They tell the compiler:

- How much memory to allocate (an `int64` needs 8 bytes, a `bool` needs 1)
- What operations are valid (you can add two `int`s, but not an `int` and a `string`)
- How to interpret the bits (the same 8 bytes could be an `int64`, a `float64`, or two `int32`s)

A static type system catches a whole class of bugs **before your code runs**:

```go
count := 42
name := "hello"
result := count + name  // COMPILE ERROR: cannot add int and string
```

In Python, this would crash at runtime. In Go, it never compiles. You find the bug before you ship.

## First Principles

### Static vs Dynamic Typing

| Aspect | Static Typing (Go, Rust, Java) | Dynamic Typing (Python, Ruby, JS) |
|--------|-------------------------------|----------------------------------|
| **When types are checked** | Compile time | Runtime |
| **Type declarations** | Often required (or inferred) | Never required |
| **Bug detection** | Many bugs caught before running | Bugs found when code executes |
| **Refactoring** | Compiler tells you what broke | Must rely on tests |
| **Flexibility** | Less — must satisfy the type checker | More — anything goes until it crashes |
| **Tooling** | Better — IDE knows all types | Harder — types are not always known |
| **Performance** | Better — compiler can optimize knowing types | Worse — must check types at runtime |

Go's type system is static but relatively simple compared to languages like Rust or Haskell. It does not have sum types, pattern matching, or complex generics. This is intentional — Go values simplicity.

## How Go Does It

### Basic Types

Go has a small set of built-in types:

#### Boolean

```go
var active bool    // zero value: false
active = true
```

#### Numeric Types

| Type | Size | Range | Use case |
|------|------|-------|----------|
| `int` | Platform-dependent (32 or 64 bit) | Varies | Default integer type |
| `int8` | 1 byte | -128 to 127 | Rarely used directly |
| `int16` | 2 bytes | -32,768 to 32,767 | Rarely used directly |
| `int32` | 4 bytes | -2B to 2B | Also aliased as `rune` |
| `int64` | 8 bytes | ±9.2 quintillion | Large numbers, timestamps |
| `uint` | Platform-dependent | 0 to max | Unsigned default |
| `uint8` | 1 byte | 0 to 255 | Also aliased as `byte` |
| `uint16` | 2 bytes | 0 to 65,535 | Rarely used directly |
| `uint32` | 4 bytes | 0 to 4B | Rarely used directly |
| `uint64` | 8 bytes | 0 to 18.4 quintillion | Large unsigned numbers |
| `float32` | 4 bytes | ±3.4e38 | When precision is less critical |
| `float64` | 8 bytes | ±1.8e308 | Default float type |

Two type aliases are critical in Go:

```go
type byte = uint8  // a byte of data
type rune = int32  // a Unicode code point
```

You will see `byte` and `rune` everywhere. They are just `uint8` and `int32` with names that communicate intent.

#### Strings

```go
var name string     // zero value: "" (empty string)
name = "groyep"
```

Strings in Go are **immutable sequences of bytes**. Not characters — bytes. This distinction matters enormously and is covered in depth in [concept 09](./09-strings-and-encoding.md).

### Zero Values

In Go, every type has a **zero value** — the value a variable holds if you declare it without initializing it. There is no "undefined" or "null" for basic types.

| Type | Zero Value |
|------|-----------|
| `bool` | `false` |
| `int`, `int64`, etc. | `0` |
| `float64`, `float32` | `0.0` |
| `string` | `""` (empty string) |
| `byte` | `0` |
| `rune` | `0` |
| Pointer | `nil` |
| Slice | `nil` |
| Map | `nil` |
| Channel | `nil` |
| Interface | `nil` |

This is a deliberate design choice. It means you can always use a declared variable — it has a meaningful default, not garbage memory.

```go
var count int       // count is 0, not garbage
var name string     // name is "", not null
var found bool      // found is false, not undefined

fmt.Println(count)  // prints: 0
fmt.Println(name)   // prints: (empty)
fmt.Println(found)  // prints: false
```

### Variable Declaration

Go offers multiple ways to declare variables:

```go
// Full declaration with explicit type
var count int = 0

// Type inferred from value
var count = 0

// Short declaration (most common inside functions)
count := 0

// Multiple variables
var (
    pattern string
    count   int
    verbose bool
)
```

The **short declaration operator** `:=` is the most common inside functions. It declares and initializes in one step, with the type inferred:

```go
name := "groyep"       // type: string
count := 42            // type: int
ratio := 3.14          // type: float64
found := true          // type: bool
```

You cannot use `:=` outside a function. Package-level variables must use `var`.

### Type Inference

Go infers types from values. You do not need to write the type when the compiler can figure it out:

```go
x := 42              // int (not int64, not int32 — just int)
y := 3.14            // float64 (always float64 for untyped float constants)
z := "hello"         // string
ok := true           // bool
b := byte('A')       // byte (uint8) — explicit conversion needed
```

Note the subtlety: `42` becomes `int` (platform-dependent size), not `int64`. If you need a specific size, be explicit:

```go
var lineNumber int64 = 42
```

### Type Conversions

Go has **no implicit type conversions**. You must convert explicitly:

```go
var i int = 42
var f float64 = float64(i)    // explicit conversion
var u uint = uint(f)          // explicit conversion

// This will NOT compile:
// var f float64 = i   // ERROR: cannot use i (type int) as type float64
```

This verbosity prevents subtle bugs. In C, implicit conversions between signed and unsigned integers are a common source of security vulnerabilities. Go makes you be deliberate.

### Constants

Constants are immutable values known at compile time:

```go
const maxLineLength = 1024
const version = "0.1.0"
const pi = 3.14159265358979

// Grouped constants
const (
    exitSuccess = 0
    exitFailure = 1
    exitUsage   = 2
)
```

Go constants are **untyped** by default — they have more precision and flexibility than variables. An untyped constant adapts to whatever context it is used in:

```go
const x = 42

var i int = x        // works: x adapts to int
var f float64 = x    // works: x adapts to float64
var b byte = x       // works: x adapts to byte (42 fits in uint8)
```

## Memory: Stack vs Heap

When your program runs, data lives in two places:

### The Stack

- **Fast** — allocation is just moving a pointer
- **Automatic cleanup** — when a function returns, its stack frame is gone
- **Limited size** — typically a few MB per goroutine (Go starts at 8KB and grows)
- **Local variables** usually go here

```go
func search(line string) bool {
    count := 0    // lives on the stack — fast, automatic cleanup
    // ...
    return count > 0
}
```

### The Heap

- **Slower** — allocation requires finding free space
- **Garbage collected** — the runtime must track and free it
- **Unlimited** (up to available memory)
- **Shared data** and data that outlives the function goes here

```go
func createResult() *Result {
    r := Result{Line: "hello", Number: 42}
    return &r    // r must live on the heap — it outlives this function
}
```

### Escape Analysis

Go decides whether a variable lives on the stack or heap automatically through **escape analysis**. If the compiler can prove a variable does not escape the function, it goes on the stack. If it might be referenced after the function returns, it goes on the heap.

You do not control this directly. The Go compiler makes the decision. You can see what it decides with:

```bash
go build -gcflags="-m" main.go
```

For now, just know: Go handles this for you. You do not call `malloc` or `free`.

### Value Types vs Reference Types

In Go, some types are **value types** — assigning or passing them creates a copy:

```go
a := 42
b := a      // b is a copy — changing b does not change a
b = 100
fmt.Println(a)  // still 42
```

Some types behave like **reference types** — they contain an internal pointer to underlying data:

| Value Types | Reference-like Types |
|-------------|---------------------|
| `int`, `float64`, etc. | Slices |
| `bool` | Maps |
| `string` | Channels |
| Arrays | Pointers |
| Structs | Functions |

This distinction matters when you pass data to functions or assign between variables. Value types are copied; reference-like types share underlying data.

## Tradeoffs

| You gain | You give up |
|----------|-------------|
| Bugs caught at compile time | Must satisfy the type checker (more upfront work) |
| Clear documentation of intent | Verbosity when types are obvious |
| Better tooling (autocomplete, refactoring) | Less flexibility than dynamic typing |
| Performance (compiler knows sizes) | Explicit conversions required |
| Zero values prevent null-like bugs | Zero values can be surprising (empty string vs missing) |
| Automatic memory management | GC pauses, less control than manual management |

## Why This Matters for Grep

Your grep clone works with specific types:

- **`string`** — the pattern to search for, each line of input, file paths
- **`[]byte`** — the raw bytes read from files (more efficient than string for I/O)
- **`int`** — line numbers, match counts, exit codes
- **`bool`** — flags like case-insensitive, invert match, count-only
- **`byte` and `rune`** — individual characters when doing character-level matching

Understanding types helps you:
- Choose `[]byte` vs `string` for performance (reading files as bytes is faster)
- Use `rune` correctly when handling Unicode text
- Return `int` exit codes (0 for match found, 1 for no match, 2 for error — just like real grep)
- Use zero values wisely — a `string` flag defaults to `""`, a `bool` flag defaults to `false`

The type system also protects you. If you accidentally try to compare a line number with a file path, the compiler stops you. In a dynamically typed language, that bug would hide until a user hit it in production.

## Further Reading

1. **[The Go Programming Language Specification — Types](https://go.dev/ref/spec#Types)** — Official
   The authoritative reference for Go's type system. Every type, every rule, every edge case.

2. **[Effective Go — Data](https://go.dev/doc/effective_go#data)** — Official
   Practical guidance on using types, allocation with `new` and `make`, and data structures.

3. **[Go Blog — Constants](https://go.dev/blog/constants)** — Rob Pike
   A deep dive into how Go constants work, including untyped constants and their surprising flexibility.

4. **[Go Blog — The Laws of Reflection](https://go.dev/blog/laws-of-reflection)** — Rob Pike
   Understanding Go's type system at a deeper level through the lens of reflection. Advanced but illuminating.

5. **[A Tour of Go — Basic Types](https://go.dev/tour/basics/11)** — Official
   Interactive introduction to Go's basic types. Good for hands-on practice.

6. **[Understanding Allocations in Go](https://medium.com/eureka-engineering/understanding-allocations-in-go-stack-heap-memory-9a2631b5035d)** — Medium
   A practical explanation of stack vs heap allocation in Go with examples.

7. **[Go Memory Model](https://go.dev/ref/mem)** — Official
   The specification for how Go handles memory access across goroutines. Advanced reading for later.

8. **[William Kennedy — "Garbage Collection in Go"](https://www.ardanlabs.com/blog/2018/12/garbage-collection-in-go-part1-semantics.html)** — Ardan Labs
   A three-part series on how Go's garbage collector works. Not needed now, but valuable when you want to understand performance.

---

*Previous: [03 — Modules and Packages](./03-modules-and-packages.md) · Next: [05 — Control Flow](./05-control-flow.md)*
