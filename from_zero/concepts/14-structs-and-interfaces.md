# Structs and Interfaces

## What Is It

Go has two primary mechanisms for building abstractions:

**Structs** group related data together. A struct is a composite type — it bundles multiple values (fields) into a single unit.

**Interfaces** define behavior. An interface is a set of method signatures. Any type that implements those methods satisfies the interface — automatically, without declaring it.

Together, structs and interfaces are how you model the world in Go. Structs are nouns (things with properties). Interfaces are verbs (things that can do something).

```go
// Struct: a thing with properties
type Match struct {
    Line     string
    Number   int
    Filename string
}

// Interface: a capability
type Searcher interface {
    Search(line string) bool
}
```

## What Problem It Solves

Programs need to organize data and define contracts between components.

**Structs solve the data organization problem.** Without structs, you'd pass around loose variables — a line number here, a filename there, a match string somewhere else. Structs bundle related values so they travel together and can't get out of sync.

**Interfaces solve the coupling problem.** Without interfaces, your search function would need to know the concrete type of everything it works with. Want to search a file? Hard-code `*os.File`. Want to search a string? Different function. Want to test with fake data? Another function. Interfaces let you write code against behavior, not concrete types.

## First Principles

### Structs: Grouping Data

A struct is defined with the `type` keyword:

```go
type GrepResult struct {
    Filename string
    LineNum  int
    Line     string
    Matches  []string
}
```

Creating and using struct values:

```go
// Named fields (preferred — explicit and readable)
result := GrepResult{
    Filename: "server.log",
    LineNum:  42,
    Line:     "ERROR: connection timeout",
    Matches:  []string{"ERROR"},
}

// Access fields with dot notation
fmt.Printf("%s:%d:%s\n", result.Filename, result.LineNum, result.Line)
```

### Zero Values

Struct fields get zero values if not explicitly set:

```go
var result GrepResult
// result.Filename == ""   (zero value for string)
// result.LineNum  == 0    (zero value for int)
// result.Matches  == nil  (zero value for slice)
```

Go's zero values are useful defaults. A `GrepResult` with all zero values is a valid (empty) result, not a broken one.

### Methods on Structs

Methods give structs behavior. A method is a function with a **receiver**:

```go
type Config struct {
    IgnoreCase  bool
    InvertMatch bool
    LineNumber  bool
    Count       bool
    Pattern     string
}

// Method with a value receiver
func (c Config) String() string {
    return fmt.Sprintf("pattern=%q ignoreCase=%v", c.Pattern, c.IgnoreCase)
}

// Usage
cfg := Config{Pattern: "hello", IgnoreCase: true}
fmt.Println(cfg) // Calls cfg.String() automatically
```

### Value Receivers vs Pointer Receivers

This is one of the most important distinctions in Go:

```go
// Value receiver: gets a COPY of the struct
func (c Config) IsValid() bool {
    return c.Pattern != ""
}

// Pointer receiver: gets a POINTER to the struct (can modify it)
func (c *Config) SetPattern(p string) {
    c.Pattern = p  // Modifies the original
}
```

| Receiver Type | Gets a...         | Can Modify? | When to Use                          |
|--------------|-------------------|-------------|--------------------------------------|
| Value `(c Config)` | Copy of struct | No          | Read-only methods, small structs     |
| Pointer `(c *Config)` | Pointer to struct | Yes    | Methods that modify, large structs   |

**Rule of thumb**: if any method needs a pointer receiver, make all methods use pointer receivers. Consistency prevents subtle bugs.

```go
// Be consistent: all pointer receivers
type Matcher struct {
    pattern    *regexp.Regexp
    ignoreCase bool
    matches    int
}

func (m *Matcher) Match(line string) bool {
    if m.pattern.MatchString(line) {
        m.matches++  // Modifies the struct
        return true
    }
    return false
}

func (m *Matcher) Count() int {
    return m.matches  // Read-only, but consistent with Match
}
```

### Composition Over Inheritance

Go has no inheritance. No classes. No `extends`. Instead, Go uses **composition** — you embed one struct inside another:

```go
type FileResult struct {
    Filename string
    Results  []LineResult
}

type LineResult struct {
    Number  int
    Text    string
    Matched bool
}
```

You can also **embed** structs to promote their fields and methods:

```go
type TimestampedResult struct {
    LineResult           // Embedded — fields are promoted
    Timestamp  time.Time
}

// You can access LineResult's fields directly:
tr := TimestampedResult{
    LineResult: LineResult{Number: 42, Text: "hello", Matched: true},
    Timestamp:  time.Now(),
}
fmt.Println(tr.Number) // 42 — promoted from LineResult
fmt.Println(tr.Text)   // "hello" — promoted from LineResult
```

Embedding is not inheritance. There's no "is-a" relationship. A `TimestampedResult` is not a `LineResult` — it **has** a `LineResult`.

## Interfaces: Defining Behavior

### Implicit Satisfaction

This is Go's most distinctive feature. In Java or C#, you write `class Foo implements Bar`. In Go, you just... implement the methods:

```go
// Define an interface
type Searcher interface {
    Search(line string) bool
}

// This struct satisfies Searcher — no declaration needed
type LiteralSearcher struct {
    pattern string
}

func (s *LiteralSearcher) Search(line string) bool {
    return strings.Contains(line, s.pattern)
}

// This also satisfies Searcher — completely different implementation
type RegexSearcher struct {
    re *regexp.Regexp
}

func (s *RegexSearcher) Search(line string) bool {
    return s.re.MatchString(line)
}
```

Neither `LiteralSearcher` nor `RegexSearcher` mentions `Searcher`. They don't even need to know `Searcher` exists. They satisfy it just by having a `Search(string) bool` method.

### Why Implicit Interfaces Matter

Implicit satisfaction enables **decoupling**. The interface and the implementation don't need to know about each other at compile time. This means:

- You can define interfaces in the **consumer** package, not the provider package
- Third-party types can satisfy your interfaces without modification
- Testing is trivial — any type with the right methods works

```go
// In your grep package, you define what you need:
type Searcher interface {
    Search(line string) bool
}

// Your grep function depends on behavior, not a concrete type:
func grep(searcher Searcher, reader io.Reader, writer io.Writer) {
    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        if searcher.Search(scanner.Text()) {
            fmt.Fprintln(writer, scanner.Text())
        }
    }
}
```

### Common Standard Library Interfaces

| Interface       | Methods                           | Used For                    |
|----------------|-----------------------------------|-----------------------------|
| `io.Reader`    | `Read([]byte) (int, error)`      | Anything you can read from  |
| `io.Writer`    | `Write([]byte) (int, error)`     | Anything you can write to   |
| `io.Closer`    | `Close() error`                  | Resources that need cleanup |
| `error`        | `Error() string`                 | Error values                |
| `fmt.Stringer` | `String() string`                | Human-readable representation |
| `sort.Interface` | `Len()`, `Less()`, `Swap()`   | Sortable collections        |

### The Empty Interface: `any`

The empty interface has zero methods, so every type satisfies it:

```go
// These are equivalent (any is an alias for interface{})
var x interface{}
var y any

x = 42
x = "hello"
x = Config{Pattern: "test"}
```

Use `any` sparingly. It throws away type safety. You almost always want a specific interface.

### Interface Segregation

Go culture favors **small interfaces**:

```go
// Good: small, focused interfaces
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Composed when needed
type ReadWriter interface {
    Reader
    Writer
}
```

```go
// Bad: fat interface
type FileSystem interface {
    Open(path string) (*File, error)
    Create(path string) (*File, error)
    Remove(path string) error
    Rename(old, new string) error
    Stat(path string) (FileInfo, error)
    ReadDir(path string) ([]DirEntry, error)
    Walk(root string, fn WalkFunc) error
    // ... 15 more methods
}
```

The Go proverb: **"The bigger the interface, the weaker the abstraction."** A one-method interface like `io.Reader` is satisfied by hundreds of types. A 20-method interface is satisfied by almost nothing.

### Accepting Interfaces, Returning Structs

A key Go idiom:

```go
// Accept an interface (flexible — works with any Reader)
func processInput(r io.Reader) error { ... }

// Return a concrete type (useful — caller gets full access)
func NewMatcher(pattern string) (*Matcher, error) { ... }
```

Functions should accept the narrowest interface they need and return the most specific type they can. This maximizes both flexibility (callers can pass many types) and usefulness (callers get a fully-featured return value).

## Tradeoffs

| Decision                        | Benefit                              | Cost                              |
|--------------------------------|--------------------------------------|-----------------------------------|
| Structs (composition)          | Explicit, flexible, no hierarchy     | No polymorphism without interfaces|
| Classes (inheritance)          | Code reuse via hierarchy             | Fragile base class, tight coupling|
| Implicit interfaces            | Decoupled, testable                  | Harder to find all implementors   |
| Explicit interfaces (`implements`) | Clear documentation of intent    | Coupling between packages         |
| Small interfaces (1-2 methods) | Many types satisfy them              | More interfaces to name/manage    |
| Large interfaces               | One contract covers everything       | Few types can satisfy them        |
| Value receivers                | Safe (no mutation), copyable         | Can't modify, may be slow for large structs |
| Pointer receivers              | Can modify, efficient for large structs | Nil pointer risk, shared state |

## Why This Matters for Grep

Structs and interfaces let you design a grep that is **modular, testable, and extensible**:

### Structs for Organizing Grep's Data

```go
// Configuration from CLI flags
type Config struct {
    Pattern     string
    IgnoreCase  bool
    InvertMatch bool
    LineNumber  bool
    Count       bool
    Color       string
    Files       []string
}

// A single match result
type Match struct {
    Filename string
    LineNum  int
    Line     string
}
```

Without structs, your search function would need 8+ parameters. With structs, it needs one:

```go
// Without structs: nightmare to call and maintain
func search(pattern string, ignoreCase bool, invert bool, lineNum bool,
    count bool, color string, filename string) { ... }

// With structs: clean and extensible
func search(cfg *Config, reader io.Reader) ([]Match, error) { ... }
```

### Interfaces for Testable Design

```go
// Define what a "searcher" does
type Searcher interface {
    Search(line string) bool
}

// Production: uses regexp
type RegexSearcher struct {
    re *regexp.Regexp
}

func (s *RegexSearcher) Search(line string) bool {
    return s.re.MatchString(line)
}

// Test: always matches (for testing output formatting)
type AlwaysMatch struct{}

func (a AlwaysMatch) Search(line string) bool { return true }

// Test: never matches (for testing "no match" exit code)
type NeverMatch struct{}

func (n NeverMatch) Search(line string) bool { return false }
```

Because `grep()` accepts a `Searcher` interface, you can test every aspect of your grep without involving real regex or real files:

```go
func TestGrepOutputsMatchingLines(t *testing.T) {
    input := strings.NewReader("hello\nworld\nhello world\n")
    var output bytes.Buffer

    grep(&AlwaysMatch{}, input, &output)

    if output.String() != "hello\nworld\nhello world\n" {
        t.Errorf("expected all lines, got %q", output.String())
    }
}
```

This is the power of interfaces: your test doesn't create files, doesn't compile regexes, doesn't touch the filesystem. It tests the grep logic in isolation.

## Further Reading

1. **[Effective Go — Interfaces and Types](https://go.dev/doc/effective_go#interfaces)** — The official guide on how to think about interfaces in Go. Required reading for any Go developer.

2. **[Go Blog — "The Laws of Reflection"](https://go.dev/blog/laws-of-reflection)** — While focused on reflection, this post explains Go's type system and interfaces at a deep level.

3. **[Rob Pike — "Go Proverbs" (YouTube)](https://www.youtube.com/watch?v=PAAkCSZUG1c)** — "The bigger the interface, the weaker the abstraction" and other design wisdom directly applicable to struct/interface design.

4. **[Dave Cheney — "SOLID Go Design" (YouTube)](https://www.youtube.com/watch?v=zzAdEt3xZ1M)** — Dave Cheney applies SOLID principles to Go, with particular attention to interface segregation and dependency inversion.

5. **[Go Wiki — CodeReviewComments](https://go.dev/wiki/CodeReviewComments)** — The community's collected wisdom on Go style, including sections on interfaces, receiver types, and struct design.

6. **[Go Blog — "Go Data Structures: Interfaces"](https://research.swtch.com/interfaces)** — Russ Cox explains how interfaces work under the hood. Understanding the implementation helps you use them well.

7. **[Ardan Labs — "Interface Values Are Valueless"](https://www.ardanlabs.com/blog/2018/03/interface-values-are-valueless.html)** — Bill Kennedy's deep dive into the mechanics of interface values, nil interfaces, and common pitfalls.

8. **[Jordan Orelli — "How to use interfaces in Go"](https://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go)** — A practical guide that starts from confusion and builds to clarity. Good for newcomers.

9. **[Go by Example: Structs](https://gobyexample.com/structs)** and **[Go by Example: Interfaces](https://gobyexample.com/interfaces)** — Concise, executable examples of both concepts.

10. **[Francesc Campoy — "Understanding Go Interfaces" (YouTube)](https://www.youtube.com/watch?v=F4wUrj6pmSI)** — A talk that walks through interface design from beginner concepts to advanced patterns.
