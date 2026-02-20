# Testing Philosophy

## What Is It

Testing is the practice of writing code that verifies your other code works correctly. A test says: "given this input, I expect this output." If the output doesn't match, something is broken.

Go has testing built into the language and toolchain. No external framework needed. No configuration files. No test runner to install. You write a function that starts with `Test`, put it in a file that ends with `_test.go`, and run `go test`. That's it.

```go
// math.go
func Add(a, b int) int {
    return a + b
}

// math_test.go
func TestAdd(t *testing.T) {
    got := Add(2, 3)
    want := 5
    if got != want {
        t.Errorf("Add(2, 3) = %d, want %d", got, want)
    }
}
```

## What Problem It Solves

Without tests, you discover bugs in production — from users, at the worst possible time. Testing inverts this: you discover bugs before they ship, on your terms, automatically.

Tests solve several problems simultaneously:

- **Correctness**: Does the code do what it should?
- **Regression prevention**: Does new code break old behavior?
- **Documentation**: Tests show how code is intended to be used
- **Design feedback**: Hard-to-test code is usually poorly designed
- **Confidence**: You can refactor knowing tests will catch mistakes

## First Principles

### The Testing Pyramid

```
        /\
       /  \        E2E Tests (few)
      /    \       - Test the whole system
     /------\      - Slow, expensive, brittle
    /        \
   /  Integ.  \   Integration Tests (some)
  /    Tests   \  - Test component interactions
 /--------------\ - Medium speed
/                \
/   Unit Tests    \ Unit Tests (many)
/------------------\- Test individual functions
                     - Fast, cheap, reliable
```

| Test Type    | Scope                    | Speed    | Reliability | Quantity |
|-------------|--------------------------|----------|-------------|----------|
| Unit        | Single function/method   | Fast     | High        | Many     |
| Integration | Multiple components      | Medium   | Medium      | Some     |
| End-to-End  | Entire program           | Slow     | Lower       | Few      |

For a grep clone, this translates to:

- **Unit tests**: Does the regex matching work? Does line numbering format correctly?
- **Integration tests**: Does reading a file and matching produce correct output?
- **E2E tests**: Does running the compiled binary with flags produce expected output?

### Test-Driven Development (TDD)

TDD is a discipline where you write the test **before** the code:

1. **Red**: Write a test that fails (because the code doesn't exist yet)
2. **Green**: Write the minimum code to make the test pass
3. **Refactor**: Clean up the code while keeping tests green

```go
// Step 1: Red — write the test first
func TestMatchLine(t *testing.T) {
    got := MatchLine("hello world", "hello")
    if !got {
        t.Error("expected 'hello world' to match 'hello'")
    }
}

// Step 2: Green — write the minimum code
func MatchLine(line, pattern string) bool {
    return strings.Contains(line, pattern)
}

// Step 3: Refactor — improve without changing behavior
// (In this case, the code is already simple enough)
```

TDD isn't mandatory, but it's powerful: it forces you to think about the interface before the implementation, and it ensures every line of code has a corresponding test.

## How Go Does It

### File Convention

Test files live next to the code they test:

```
groyep/
  ├── matcher.go       # Production code
  ├── matcher_test.go  # Tests for matcher.go
  ├── output.go        # Production code
  └── output_test.go   # Tests for output.go
```

Files ending in `_test.go` are **only compiled during `go test`**. They're excluded from the final binary.

### Test Function Signature

Every test function must:
- Be in a `_test.go` file
- Start with `Test` followed by an uppercase letter
- Take exactly one parameter: `*testing.T`

```go
func TestSomething(t *testing.T) {
    // Test code here
}
```

### Reporting Failures

```go
func TestMatchLine(t *testing.T) {
    // t.Error reports a failure but continues the test
    if !MatchLine("hello", "hello") {
        t.Error("expected match")
    }

    // t.Errorf reports with formatting
    got := MatchLine("hello", "world")
    if got {
        t.Errorf("MatchLine(%q, %q) = true, want false", "hello", "world")
    }

    // t.Fatal reports and STOPS the test immediately
    cfg, err := ParseConfig([]string{"-i", "pattern"})
    if err != nil {
        t.Fatal(err) // No point continuing if config parsing failed
    }

    // t.Fatalf with formatting
    if cfg.Pattern != "pattern" {
        t.Fatalf("Pattern = %q, want %q", cfg.Pattern, "pattern")
    }
}
```

| Method       | Behavior                             | Use When                           |
|-------------|--------------------------------------|------------------------------------|
| `t.Error`   | Report failure, continue             | Failure is informative but not blocking |
| `t.Errorf`  | Report formatted failure, continue   | Same, with format string           |
| `t.Fatal`   | Report failure, stop test            | Continuing would cause crashes/meaningless errors |
| `t.Fatalf`  | Report formatted failure, stop test  | Same, with format string           |

### Table-Driven Tests

Table-driven tests are **the** idiomatic Go testing pattern. Instead of writing one test function per case, you define a table of inputs and expected outputs:

```go
func TestMatchLine(t *testing.T) {
    tests := []struct {
        name    string
        line    string
        pattern string
        want    bool
    }{
        {
            name:    "exact match",
            line:    "hello",
            pattern: "hello",
            want:    true,
        },
        {
            name:    "partial match",
            line:    "hello world",
            pattern: "hello",
            want:    true,
        },
        {
            name:    "no match",
            line:    "hello",
            pattern: "world",
            want:    false,
        },
        {
            name:    "empty pattern matches everything",
            line:    "hello",
            pattern: "",
            want:    true,
        },
        {
            name:    "empty line",
            line:    "",
            pattern: "hello",
            want:    false,
        },
        {
            name:    "case sensitive by default",
            line:    "Hello",
            pattern: "hello",
            want:    false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := MatchLine(tt.line, tt.pattern)
            if got != tt.want {
                t.Errorf("MatchLine(%q, %q) = %v, want %v",
                    tt.line, tt.pattern, got, tt.want)
            }
        })
    }
}
```

Why table-driven tests are preferred:

| Benefit              | Explanation                                           |
|---------------------|-------------------------------------------------------|
| Easy to add cases   | Just add another entry to the table                   |
| Consistent format   | Every case has the same structure                     |
| Named cases         | `t.Run(name, ...)` makes failures easy to identify    |
| Reduces boilerplate | The test logic is written once                        |
| Documents behavior  | The table is a specification                          |

### Subtests with `t.Run`

`t.Run` creates named subtests that can be run individually:

```go
// Run all tests:
// go test ./...

// Run just the "no match" case:
// go test -run TestMatchLine/no_match
```

Subtests also enable parallel execution:

```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        t.Parallel() // Run this subtest in parallel with others
        got := MatchLine(tt.line, tt.pattern)
        if got != tt.want {
            t.Errorf("got %v, want %v", got, tt.want)
        }
    })
}
```

### Test Helpers

When multiple tests share setup logic, extract a helper:

```go
// t.Helper() marks this as a helper — failure line numbers point to the caller
func assertMatch(t *testing.T, line, pattern string) {
    t.Helper()
    if !MatchLine(line, pattern) {
        t.Errorf("expected %q to match %q", line, pattern)
    }
}

func assertNoMatch(t *testing.T, line, pattern string) {
    t.Helper()
    if MatchLine(line, pattern) {
        t.Errorf("expected %q to NOT match %q", line, pattern)
    }
}
```

Without `t.Helper()`, a failure in `assertMatch` would report the line inside the helper. With it, the failure points to the test that called the helper — much more useful.

### The `testdata` Directory

Go has a special convention: a directory named `testdata` is ignored by `go build` but accessible to tests:

```
groyep/
  ├── search.go
  ├── search_test.go
  └── testdata/
      ├── simple.txt
      ├── multiline.txt
      └── unicode.txt
```

```go
func TestSearchFile(t *testing.T) {
    // testdata files are accessed relative to the test file
    results, err := SearchFile("hello", "testdata/simple.txt")
    if err != nil {
        t.Fatal(err)
    }
    // ... assert results ...
}
```

### Benchmarks

Go's testing package includes benchmarking:

```go
func BenchmarkMatchLine(b *testing.B) {
    line := "the quick brown fox jumps over the lazy dog"
    pattern := "fox"

    for i := 0; i < b.N; i++ {
        MatchLine(line, pattern)
    }
}
```

Run with `go test -bench=.`:

```
BenchmarkMatchLine-8    50000000    25.3 ns/op
```

### Testing CLI Programs

For grep, you need to test the full command-line interface:

```go
func TestGrepCLI(t *testing.T) {
    tests := []struct {
        name     string
        args     []string
        stdin    string
        wantOut  string
        wantErr  string
        wantCode int
    }{
        {
            name:     "match from stdin",
            args:     []string{"hello"},
            stdin:    "hello world\ngoodbye\n",
            wantOut:  "hello world\n",
            wantCode: 0,
        },
        {
            name:     "no match",
            args:     []string{"missing"},
            stdin:    "hello world\n",
            wantOut:  "",
            wantCode: 1,
        },
        {
            name:     "missing pattern",
            args:     []string{},
            wantCode: 2,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var stdout, stderr bytes.Buffer
            stdin := strings.NewReader(tt.stdin)

            code := run(tt.args, stdin, &stdout, &stderr)

            if code != tt.wantCode {
                t.Errorf("exit code = %d, want %d", code, tt.wantCode)
            }
            if stdout.String() != tt.wantOut {
                t.Errorf("stdout = %q, want %q", stdout.String(), tt.wantOut)
            }
        })
    }
}
```

The key insight: your `run` function should accept `io.Reader` and `io.Writer` for stdin/stdout/stderr, not use `os.Stdin`/`os.Stdout` directly. This makes it testable.

## Tradeoffs

| Decision                          | Benefit                           | Cost                                |
|----------------------------------|-----------------------------------|-------------------------------------|
| Write tests first (TDD)         | Better design, full coverage      | Slower initial development          |
| Write tests after               | Faster initial development        | May miss edge cases, harder to retrofit |
| Table-driven tests              | Easy to add cases, clear          | Requires upfront structure          |
| Individual test functions        | Simple, no pattern to learn       | Lots of repeated boilerplate        |
| Test with real files (`testdata`)| Tests realistic scenarios         | Slower, more to maintain            |
| Test with in-memory data         | Fast, self-contained              | Doesn't catch file-system issues    |
| Test private functions           | Thorough coverage                 | Tests break when refactoring internals |
| Test only public API             | Tests survive refactoring          | May miss internal bugs              |
| External test frameworks         | Richer assertions, BDD style      | Extra dependency, non-standard      |
| Standard `testing` package       | No dependencies, Go-standard      | Verbose assertions                  |

### The No-Framework Philosophy

Go intentionally doesn't have assertion functions like `assertEquals`. The Go team believes:

```go
// Java-style (not Go)
assertEquals(5, Add(2, 3))

// Go-style (explicit, clear error messages)
got := Add(2, 3)
want := 5
if got != want {
    t.Errorf("Add(2, 3) = %d, want %d", got, want)
}
```

The Go way is more verbose but produces better error messages. When a test fails, you see exactly what happened, not just "expected 5 but got 4."

## Why This Matters for Grep

Grep is a program with clearly defined inputs and outputs, which makes it highly testable:

| Input                | Output                    | Testable Aspect                  |
|---------------------|---------------------------|----------------------------------|
| Pattern + text      | Matching lines            | Core search logic                |
| `-i` flag + text    | Case-insensitive matches  | Flag behavior                    |
| `-n` flag           | Lines with numbers        | Output formatting                |
| `-v` flag           | Non-matching lines        | Invert logic                     |
| No matches          | Empty output, exit code 1 | Exit code behavior               |
| Invalid regex       | Error message, exit code 2| Error handling                   |
| Multiple files      | Prefixed filenames        | Multi-file output                |
| stdin input         | Matching lines            | Pipe compatibility               |

Table-driven tests are especially natural for grep because you can encode entire test scenarios:

```go
func TestGrep(t *testing.T) {
    tests := []struct {
        name       string
        pattern    string
        input      string
        ignoreCase bool
        invert     bool
        want       string
    }{
        {"literal match", "hello", "hello\nworld\n", false, false, "hello\n"},
        {"no match", "xyz", "hello\nworld\n", false, false, ""},
        {"case insensitive", "HELLO", "hello\nworld\n", true, false, "hello\n"},
        {"invert", "hello", "hello\nworld\n", false, true, "world\n"},
        {"regex dot", "h.llo", "hello\nhallo\nworld\n", false, false, "hello\nhallo\n"},
        {"empty input", "hello", "", false, false, ""},
        {"multiple matches", "o", "one\ntwo\nthree\n", false, false, "one\ntwo\n"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ... run grep with tt inputs, compare to tt.want ...
        })
    }
}
```

Every feature you add gets a row in this table. Every bug you fix gets a test case. Over time, this table becomes a comprehensive specification of your grep's behavior.

## Further Reading

1. **[Go `testing` package documentation](https://pkg.go.dev/testing)** — The official reference. Read the package overview, not just the function docs — it explains conventions and best practices.

2. **[Go Blog — "Using Subtests and Sub-benchmarks"](https://go.dev/blog/subtests)** — Explains `t.Run`, table-driven tests, and how to run specific test cases. Essential for productive testing.

3. **[Mitchell Hashimoto — "Advanced Testing with Go" (YouTube)](https://www.youtube.com/watch?v=8hQG7QlcLBk)** — The creator of Terraform, Vault, and Vagrant shares hard-won testing patterns in Go. Covers table-driven tests, test helpers, golden files, and more.

4. **[Dave Cheney — "Writing table driven tests in Go"](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)** — A focused post on why table-driven tests are the preferred pattern and how to write them well.

5. **["Learn Go with Tests" by Chris James](https://quii.gitbook.io/learn-go-with-tests/)** — A free online book that teaches Go through TDD. Each chapter introduces a concept by writing tests first. Excellent for beginners.

6. **[Go Wiki — TableDrivenTests](https://go.dev/wiki/TableDrivenTests)** — The community wiki page on the table-driven test pattern, with examples and discussion.

7. **[Go Blog — "Testable Examples in Go"](https://go.dev/blog/examples)** — How to write `Example` functions that serve as both documentation and tests. A unique Go feature.

8. **[Testify library (for reference)](https://github.com/stretchr/testify)** — The most popular testing helper library for Go. Worth knowing about even if you stick with the standard library. Understand what it offers and why some teams prefer it.

9. **[Mat Ryer — "Writing Beautiful Packages in Go" (YouTube)](https://www.youtube.com/watch?v=cmkKxNN7cs4)** — Includes excellent advice on testing patterns, package design, and making code testable.

10. **[The Go Programming Language — Chapter 11: Testing](https://www.gopl.io/)** — Donovan and Kernighan's treatment of Go's testing tools and philosophy. Thorough and well-explained.
