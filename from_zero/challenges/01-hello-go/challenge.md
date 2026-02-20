# Challenge 01: Hello, Go!

## Prerequisites

- [What is Go](../../concepts/02-what-is-go.md)
- [Modules and Packages](../../concepts/03-modules-and-packages.md)
- Completed: [Challenge 00 — Setup](../00-setup/challenge.md)

## Overview

Every programming journey starts with Hello World. In Go, this teaches you three fundamental concepts: packages, imports, and printing to the screen.

---

## Concepts

### `package main`

Every Go file starts with a package declaration. The `main` package is special — it tells Go "this package produces an executable program" (as opposed to a library that other code imports).

```go
package main
```

### `func main()`

The `main()` function inside `package main` is where execution starts. When you run your program, Go calls this function first. No arguments, no return value.

```go
func main() {
    // your code here
}
```

### `import "fmt"`

Go's standard library is organized into packages. To use a package, you import it. The `fmt` package provides formatted I/O — you'll use it constantly.

```go
import "fmt"
```

For multiple imports, use parentheses:

```go
import (
    "fmt"
    "os"
)
```

### `fmt.Println()`

`Println` prints its arguments to stdout, followed by a newline. You can print multiple values:

```go
fmt.Println("Hello, World!")    // prints: Hello, World!
fmt.Println("Age:", 25)         // prints: Age: 25
```

For more on these concepts, see the [What is Go](../../concepts/02-what-is-go.md) concept doc.

---

## Task

Create a Go program that prints two lines to stdout, exactly as shown:

```
Hello, World!
Hello, Go!
```

### Setup Reminder

Remember from Challenge 00 — you need to:

1. Initialize the module yourself: `go mod init hello`
2. Create `main.go` yourself

### Requirements

1. The program must print exactly two lines
2. First line: `Hello, World!`
3. Second line: `Hello, Go!`
4. No extra blank lines, no trailing spaces

### Verify Your Solution

```bash
# Run the shell tests
bash test.sh

# Run the Go tests
go test -v
```

---

## Hints

<details>
<summary>Hint 1: How to print two lines</summary>

Call `fmt.Println()` twice — once for each line.

</details>

<details>
<summary>Hint 2: Exact strings matter</summary>

Make sure you match the punctuation exactly: `Hello, World!` has a comma after "Hello" and an exclamation mark at the end. Same for `Hello, Go!`.

</details>

<details>
<summary>Hint 3: Full skeleton</summary>

```go
package main

import "fmt"

func main() {
    fmt.Println(/* first line */)
    fmt.Println(/* second line */)
}
```

</details>

---

## Further Reading

- [fmt package documentation](https://pkg.go.dev/fmt)
- [A Tour of Go](https://go.dev/tour/welcome/1)
- [Go by Example: Hello World](https://gobyexample.com/hello-world)
