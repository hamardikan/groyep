# Challenge 04: Loops

## Prerequisites

- [Control Flow](../../concepts/05-control-flow.md)
- Completed: [Challenge 03 — Conditionals](../03-conditionals/challenge.md)

## Overview

Loops let you repeat code. Combined with conditionals, they let you process sequences of data — which is most of what programs do. You'll implement FizzBuzz, the classic programming exercise that tests your ability to combine loops and conditionals.

---

## Concepts

### The `for` Loop

Go has only one looping construct: `for`. But it's flexible enough to cover all cases.

**Classic for loop (like C/Java):**

```go
for i := 1; i <= 10; i++ {
    fmt.Println(i)
}
```

Three parts separated by semicolons:
1. **Init**: `i := 1` — runs once before the loop starts
2. **Condition**: `i <= 10` — checked before each iteration; loop stops when false
3. **Post**: `i++` — runs after each iteration

**While-style loop (omit init and post):**

```go
i := 1
for i <= 10 {
    fmt.Println(i)
    i++
}
```

**Infinite loop (omit everything):**

```go
for {
    // runs forever until break or return
}
```

### Modulo Operator

The `%` operator gives you the remainder of division. It's essential for FizzBuzz:

```go
15 % 3  // = 0 (15 is divisible by 3)
15 % 5  // = 0 (15 is divisible by 5)
7 % 3   // = 1 (7 is not divisible by 3)
```

If `n % d == 0`, then `n` is evenly divisible by `d`.

### Converting Strings to Numbers (Reminder)

From the previous challenge, you know `strconv.Atoi` converts a string to an integer:

```go
n, err := strconv.Atoi(os.Args[1])
if err != nil {
    fmt.Println("error: not a number")
    return
}
```

### Printing Numbers with Println

`fmt.Println` can print integers directly — no format string needed:

```go
fmt.Println(42)     // prints: 42
fmt.Println("Fizz") // prints: Fizz
```

For more on loops and control flow, see [Control Flow](../../concepts/05-control-flow.md).

---

## Task

Create a FizzBuzz program that takes a number N as a command-line argument and prints from 1 to N.

### Setup

1. `go mod init fizzbuzz`
2. Create `main.go`

### Requirements

For each number from 1 to N (inclusive), print:

| Condition | Output |
|---|---|
| Divisible by both 3 and 5 | `FizzBuzz` |
| Divisible by 3 only | `Fizz` |
| Divisible by 5 only | `Buzz` |
| Neither | The number itself |

Each value on its own line.

**Error handling:**

| Input | Output |
|---|---|
| `./fizzbuzz` (no args) | `usage: fizzbuzz <n>` |
| `./fizzbuzz abc` | `error: not a number` |
| `./fizzbuzz 0` | (no output) |
| `./fizzbuzz -5` | (no output) |

### Example

```bash
$ ./fizzbuzz 15
1
2
Fizz
4
Buzz
Fizz
7
8
Fizz
Buzz
11
Fizz
13
14
FizzBuzz
```

### Verify Your Solution

```bash
bash test.sh
go test -v
```

---

## Hints

<details>
<summary>Hint 1: Order of checks matters</summary>

You must check "divisible by both 3 and 5" BEFORE checking "divisible by 3" or "divisible by 5" alone. Otherwise, 15 would print "Fizz" instead of "FizzBuzz".

```go
if i%3 == 0 && i%5 == 0 {
    // FizzBuzz
} else if i%3 == 0 {
    // Fizz
} else if i%5 == 0 {
    // Buzz
} else {
    // the number
}
```

</details>

<details>
<summary>Hint 2: Alternative check for "divisible by both"</summary>

Instead of `i%3 == 0 && i%5 == 0`, you can check `i%15 == 0` — since 15 is the least common multiple of 3 and 5.

</details>

<details>
<summary>Hint 3: Loop structure</summary>

```go
n, err := strconv.Atoi(os.Args[1])
// ... error handling ...

for i := 1; i <= n; i++ {
    // ... FizzBuzz logic ...
}
```

If N is 0 or negative, the loop simply doesn't execute — no special handling needed.

</details>

---

## Further Reading

- [Go by Example: For](https://gobyexample.com/for)
- [FizzBuzz on Wikipedia](https://en.wikipedia.org/wiki/Fizz_buzz)
- [Go Spec: For statements](https://go.dev/ref/spec#For_statements)
