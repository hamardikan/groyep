# Concepts Library

This is the **knowledge library** for the grep clone pathway. Each document covers a fundamental concept you need to understand deeply before (or while) you tackle the challenges.

## How to Use This Library

These concept docs are **not** tutorials you follow step-by-step. They are **reference material** — read them alongside the challenges. When a challenge says "Prerequisite: read concept 07", that means the concept doc gives you the mental model you need to solve the problem yourself.

Each concept doc follows the same structure:

1. **What is it** — the concept explained from first principles
2. **What problem it solves** — why this exists
3. **How Go does it** — Go's specific approach and idioms
4. **Tradeoffs** — what you gain and what you give up
5. **Why it matters for grep** — tying everything back to the end goal
6. **Further Reading** — curated links to go deeper

You do not need to read these in order, but they are numbered roughly from foundational to advanced.

---

## Table of Contents

### Foundations (Concepts 01–03)

| #  | Concept | Description |
|----|---------|-------------|
| [01](./01-how-computers-run-code.md) | How Computers Run Code | Source code to binary to execution. What compilation means and why Go is a compiled language. |
| [02](./02-what-is-go.md) | What Is Go | Go's history, design philosophy, and where it fits among programming languages. |
| [03](./03-modules-and-packages.md) | Modules and Packages | How Go organizes code — modules, packages, imports, and the `go.mod` file. |

### Core Language (Concepts 04–09)

| #  | Concept | Description |
|----|---------|-------------|
| [04](./04-types-and-memory.md) | Types and Memory | Go's type system, basic types, zero values, and how data lives in memory. |
| [05](./05-control-flow.md) | Control Flow | Branching, looping, and decision-making in Go — `if`, `switch`, `for`, and why there is nothing else. |
| [06](./06-functions.md) | Functions | Defining behavior, multiple returns, scope, and Go's approach to keeping functions simple. |
| [07](./07-error-handling.md) | Error Handling | Go's explicit error philosophy — no exceptions, just values. The `error` interface and idiomatic patterns. |
| [08](./08-collections.md) | Collections | Arrays, slices, and maps — Go's core data structures for grouping and organizing data. |
| [09](./09-strings-and-encoding.md) | Strings and Encoding | How text actually works — ASCII, Unicode, UTF-8, and why Go strings are byte slices. |

### I/O and System Interaction (Concepts 10–13)

| #  | Concept | Description |
|----|---------|-------------|
| 10 | Reading Files and Streams | `os.Open`, `bufio.Scanner`, `io.Reader` — how grep consumes input line by line. |
| 11 | Command-Line Arguments | `os.Args`, flag parsing, and how CLI tools receive instructions from the user. |
| 12 | Standard I/O | stdin, stdout, stderr — the three streams that Unix tools live and breathe. |
| 13 | The `fmt` Package | Formatted output, string formatting verbs, and printing results. |

### Pattern Matching (Concepts 14–16)

| #  | Concept | Description |
|----|---------|-------------|
| 14 | String Searching Algorithms | Brute force, Boyer-Moore, and how `strings.Contains` works under the hood. |
| 15 | Regular Expressions | Formal language theory, finite automata, and Go's `regexp` package. |
| 16 | The `regexp` Package Deep Dive | RE2 syntax, compilation, performance characteristics, and practical usage. |

### Architecture (Concepts 17–18)

| #  | Concept | Description |
|----|---------|-------------|
| 17 | Interfaces and Polymorphism | Go's implicit interfaces, the `io.Reader`/`io.Writer` pattern, and composable design. |
| 18 | Concurrency | Goroutines, channels, and `sync` — searching multiple files in parallel. |

---

## A Note on Depth

Some of these documents go deeper than you strictly need for the challenges. That is intentional. Understanding **why** something works the way it does makes you a better engineer than just knowing **how** to use it. Read as deep as you want — the core content will always be marked clearly, and the "further reading" sections are there when you are hungry for more.

---

*These documents are part of the [groyep](../../../README.md) learning pathway.*
