# 08 — Collections

## What Is It

A **collection** is a data structure that holds multiple values. Instead of declaring fifty separate variables for fifty lines of text, you put them in a collection and work with them as a group.

Go has three primary collection types:
- **Arrays** — fixed-size, rarely used directly
- **Slices** — dynamic-size, the workhorse of Go collections
- **Maps** — key-value pairs for lookups

These are the data structures you will use in nearly every Go program.

## What Problem It Solves

Programs work with groups of things. Grep reads **lines** from a file — many lines. It may search across **multiple files**. It may need to track which **files** had matches. It might count **matches per file**.

Without collections, you are stuck with:

```go
line1 := "first line"
line2 := "second line"
line3 := "third line"
// How many lines are there? What if there are 10,000?
```

Collections give you a single variable that holds many values, with operations to add, remove, access, and iterate.

## First Principles

All collections face the same fundamental tradeoffs:

| Operation | Array/Slice | Map |
|-----------|-------------|-----|
| Access by position | O(1) — instant | N/A — no positions |
| Access by key | O(n) — must search | O(1) — instant (amortized) |
| Insert at end | O(1) amortized | O(1) amortized |
| Insert at position | O(n) — must shift elements | N/A |
| Search for value | O(n) — must scan | O(n) for values, O(1) for keys |
| Ordered? | Yes — insertion order | No — random iteration order |
| Memory | Contiguous — cache friendly | Scattered — hash table overhead |

Choose the collection based on how you need to access the data.

## How Go Does It

### Arrays

An array has a **fixed size** determined at compile time:

```go
var lines [3]string
lines[0] = "first"
lines[1] = "second"
lines[2] = "third"

// Or with a literal
lines := [3]string{"first", "second", "third"}

// Let the compiler count
lines := [...]string{"first", "second", "third"}
```

The size is part of the type: `[3]string` and `[5]string` are **different types**. You cannot pass a `[3]string` to a function expecting `[5]string`.

Arrays are rarely used directly in Go. They exist mainly as the backing store for slices. When you think "I need a collection," you almost always want a slice.

### Slices

A slice is a **dynamically-sized, flexible view** into an array. Slices are Go's primary collection type.

#### Creating Slices

```go
// From a literal
lines := []string{"first", "second", "third"}

// Empty slice (length 0, will grow as needed)
var matches []string

// With make — preallocate capacity
buffer := make([]byte, 0, 4096)  // length 0, capacity 4096

// From an existing array or slice
all := []int{1, 2, 3, 4, 5}
subset := all[1:3]  // [2, 3] — elements at index 1 and 2
```

#### Length and Capacity

Every slice has two properties:

- **Length** (`len`) — how many elements it currently holds
- **Capacity** (`cap`) — how many elements it can hold before needing to reallocate

```go
s := make([]int, 3, 10)
fmt.Println(len(s))  // 3 — three elements (all zero-valued)
fmt.Println(cap(s))  // 10 — room for 10 before reallocation
```

#### Append

`append` adds elements to a slice. If the slice has enough capacity, it uses the existing backing array. If not, it allocates a new, larger array and copies everything over:

```go
var matches []string
matches = append(matches, "first match")
matches = append(matches, "second match")
matches = append(matches, "third", "fourth")  // append multiple

fmt.Println(matches)
// [first match second match third fourth]
```

**Critical:** `append` may return a new slice pointing to a different backing array. Always assign the result back:

```go
// CORRECT
matches = append(matches, "new")

// WRONG — discards the possibly new slice
append(matches, "new")  // compiler warning: result not used
```

#### Slice Internals

A slice is a small struct — a **header** — containing three fields:

```
┌────────────────────────────────────────┐
│              Slice Header              │
├──────────┬──────────┬─────────────────┤
│ Pointer  │  Length   │    Capacity     │
│ to array │   (3)    │      (5)        │
└────┬─────┴──────────┴─────────────────┘
     │
     ▼
┌────┬────┬────┬────┬────┐
│ 10 │ 20 │ 30 │    │    │  ← backing array in memory
└────┴────┴────┴────┴────┘
  0    1    2    3    4
```

When you pass a slice to a function, the header is copied (it is a value type), but the pointer still points to the same backing array. This means:

```go
func addMatch(matches []string, line string) []string {
    // matches is a copy of the header, but the backing array is shared
    return append(matches, line)  // may allocate a new array
}
```

#### Slicing Syntax

```go
s := []int{0, 1, 2, 3, 4, 5}

s[1:4]   // [1, 2, 3]     — from index 1 up to (not including) 4
s[:3]    // [0, 1, 2]     — from start up to 3
s[2:]    // [2, 3, 4, 5]  — from index 2 to end
s[:]     // [0, 1, 2, 3, 4, 5] — the entire slice
```

Slices created this way **share the backing array** with the original. Modifying one affects the other:

```go
original := []int{1, 2, 3, 4, 5}
sub := original[1:3]  // [2, 3]
sub[0] = 99
fmt.Println(original) // [1, 99, 3, 4, 5] — original changed!
```

To get an independent copy, use `copy` or the `slices.Clone` function (Go 1.21+):

```go
original := []int{1, 2, 3, 4, 5}
independent := make([]int, len(original))
copy(independent, original)
```

#### Nil Slices vs Empty Slices

```go
var s []string      // nil slice: s == nil, len(s) == 0, cap(s) == 0
s2 := []string{}    // empty slice: s2 != nil, len(s2) == 0, cap(s2) == 0
s3 := make([]string, 0)  // empty slice: s3 != nil, len(s3) == 0
```

Both nil and empty slices have length 0. Both work with `append`, `len`, `range`, and `for`. The only difference is the nil check:

```go
// Both work fine
for _, v := range s { }     // no iterations — length is 0
s = append(s, "value")      // works on nil slices

// The difference
fmt.Println(s == nil)       // true for nil slice
fmt.Println(s2 == nil)      // false for empty slice
```

Convention: use `var s []string` (nil) when declaring. Use the literal `[]string{}` only when you need a non-nil empty slice (rare, mostly for JSON serialization where `null` vs `[]` matters).

#### How Append Grows

When a slice runs out of capacity, `append` allocates a new backing array. The growth strategy (as of Go 1.18+):

- For small slices: roughly double the capacity
- For larger slices: grow by ~25%

This means `append` is **amortized O(1)** — most calls are instant, occasional calls require copying, but averaged over many appends, the cost per operation is constant.

If you know the final size, preallocate to avoid unnecessary copies:

```go
// BAD — many small allocations as the slice grows
var lines []string
for scanner.Scan() {
    lines = append(lines, scanner.Text())
}

// GOOD — if you know the approximate size
lines := make([]string, 0, 1000)  // preallocate for ~1000 lines
for scanner.Scan() {
    lines = append(lines, scanner.Text())
}
```

### Maps

A map is an unordered collection of key-value pairs. Access by key is O(1) on average.

#### Creating Maps

```go
// With a literal
exitCodes := map[string]int{
    "match":    0,
    "no_match": 1,
    "error":    2,
}

// With make
matchCount := make(map[string]int)  // empty map, ready to use
```

#### Basic Operations

```go
// Set a value
matchCount["file1.txt"] = 5

// Get a value
count := matchCount["file1.txt"]  // 5
count := matchCount["nonexistent"] // 0 (zero value for int)

// Delete a key
delete(matchCount, "file1.txt")

// Length
fmt.Println(len(matchCount))
```

#### The Comma-Ok Idiom

The zero value problem: `matchCount["nonexistent"]` returns `0`. But does that mean "zero matches" or "key does not exist"? The comma-ok idiom distinguishes:

```go
count, ok := matchCount["file1.txt"]
if ok {
    fmt.Printf("%s: %d matches\n", "file1.txt", count)
} else {
    fmt.Println("file1.txt not searched yet")
}

// Often in a single if statement
if count, ok := matchCount["file1.txt"]; ok {
    fmt.Printf("found %d matches\n", count)
}
```

#### Iterating Maps

```go
for filename, count := range matchCount {
    fmt.Printf("%s: %d\n", filename, count)
}

// Keys only
for filename := range matchCount {
    fmt.Println(filename)
}
```

**Map iteration order is randomized** in Go. Each time you iterate, the order may differ. This is deliberate — it prevents you from depending on a specific order. If you need sorted output, collect the keys into a slice and sort:

```go
keys := make([]string, 0, len(matchCount))
for k := range matchCount {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Printf("%s: %d\n", k, matchCount[k])
}
```

#### Map Restrictions

- Keys must be **comparable** (`==` must work): strings, numbers, bools, structs of comparable types. Slices, maps, and functions **cannot** be keys.
- Maps are **not safe for concurrent access**. If multiple goroutines read and write the same map, you need synchronization (covered in concept 18).
- A nil map can be read (returns zero values) but **cannot be written to** — it will panic:

```go
var m map[string]int  // nil map
_ = m["key"]          // ok — returns 0
m["key"] = 1          // PANIC: assignment to entry in nil map
```

Always initialize maps before writing to them.

### When to Use Slice vs Map

| Use Case | Use |
|----------|-----|
| Ordered sequence of items | Slice |
| Fast lookup by key | Map |
| Iterating in order | Slice |
| Checking membership | Map (existence check is O(1)) |
| Counting occurrences | Map |
| Stack or queue behavior | Slice |
| Small collections (< ~10 items) | Slice (even for lookups — linear scan is fast for small n) |

## Tradeoffs

| You gain | You give up |
|----------|-------------|
| Slices: dynamic growth, cache-friendly | Slice: append may reallocate (copy cost) |
| Maps: O(1) lookup | Maps: more memory overhead, unordered |
| Simple API: append, len, range | No built-in generics-based methods (filter, map, reduce) until Go 1.21+ |
| Zero values work naturally | Nil map writes panic |
| Slice headers are cheap to pass | Shared backing array can cause surprises |
| Range works uniformly | Range copies values (use index for large structs) |

## Why This Matters for Grep

Grep uses collections everywhere:

**Slices:**
- `os.Args[1:]` — command-line arguments (a slice of strings)
- Reading all lines from a file — each line is an element
- Collecting all matches before printing
- Storing multiple patterns (`-e pattern1 -e pattern2`)

**Maps:**
- Counting matches per file (`map[string]int`)
- Tracking which files had matches (`map[string]bool`)
- Parsing flags and options

A typical grep flow:

```go
// os.Args is already a slice
files := os.Args[2:]  // all arguments after the pattern

// Collect results per file
results := make(map[string][]string)

for _, file := range files {
    matches, err := searchFile(file, pattern)
    if err != nil {
        fmt.Fprintf(os.Stderr, "groyep: %v\n", err)
        continue
    }
    if len(matches) > 0 {
        results[file] = matches
    }
}

// Print results
for file, matches := range results {
    for _, line := range matches {
        fmt.Printf("%s:%s\n", file, line)
    }
}
```

Understanding slice internals helps you make performance-conscious decisions. When searching large files, preallocating a matches slice avoids many small allocations. When counting matches across hundreds of files, a map gives you O(1) lookups.

## Further Reading

1. **[Go Slices: Usage and Internals](https://go.dev/blog/slices-intro)** — Go Blog, 2011
   The definitive guide to how slices work under the hood. Covers the slice header, append behavior, and common patterns.

2. **[Go Maps in Action](https://go.dev/blog/maps)** — Go Blog, 2013
   Everything about Go maps: creation, access, iteration, and common patterns.

3. **[Effective Go — Data Structures](https://go.dev/doc/effective_go#data)** — Official
   Covers arrays, slices, maps, and the `make`/`new` distinction.

4. **[A Tour of Go — Slices](https://go.dev/tour/moretypes/7)** — Official
   Interactive exercises on slices, slicing, and append.

5. **[Go Blog — Go Slices: Grow and Append](https://go.dev/blog/slices)** — 2013
   Deeper dive into slice growth behavior and the append function.

6. **[Arrays, Slices, and Strings](https://go.dev/blog/slices)** — Rob Pike
   Rob Pike's explanation of the relationship between arrays, slices, and strings in Go.

7. **[The Go Programming Language — Chapter 4: Composite Types](https://www.gopl.io/)** — Donovan & Kernighan
   Comprehensive coverage of arrays, slices, maps, and structs with practical examples.

8. **[Go by Example — Slices](https://gobyexample.com/slices)** — Mark McGranaghan
   Quick, runnable examples covering all slice operations.

9. **[Go by Example — Maps](https://gobyexample.com/maps)** — Mark McGranaghan
   Quick, runnable examples covering all map operations.

---

*Previous: [07 — Error Handling](./07-error-handling.md) · Next: [09 — Strings and Encoding](./09-strings-and-encoding.md)*
