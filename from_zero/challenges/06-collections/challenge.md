# Challenge 06: Collections

## Prerequisites

- [Collections](../../concepts/08-collections.md)
- Completed: [Challenge 05 — Functions](../05-functions/challenge.md)

## Overview

So far, you've worked with individual values. But real programs work with groups of data — lists of words, sets of numbers, mappings from keys to values. Go provides two essential collection types: **slices** (dynamic arrays) and **maps** (key-value stores). This challenge uses both.

---

## Concepts

### Slices

A slice is a dynamically-sized sequence of elements. Think of it as a flexible array.

```go
// Create a slice
numbers := []int{1, 2, 3, 4, 5}

// Access elements (0-indexed)
first := numbers[0]    // 1
third := numbers[2]    // 3

// Length
fmt.Println(len(numbers))  // 5

// Append
numbers = append(numbers, 6)  // [1, 2, 3, 4, 5, 6]
```

### Maps

A map associates keys with values. Like a dictionary or hash table.

```go
// Create a map
counts := map[string]int{}

// Set values
counts["hello"] = 1
counts["world"] = 2

// Read values
fmt.Println(counts["hello"])  // 1

// Increment (useful for counting!)
counts["hello"]++             // now 2

// Check if key exists
val, exists := counts["missing"]
// val = 0 (zero value), exists = false
```

### Iterating with `range`

The `range` keyword iterates over collections:

```go
// Over a slice
words := []string{"hello", "world"}
for i, word := range words {
    fmt.Printf("%d: %s\n", i, word)
}

// Over a map
counts := map[string]int{"a": 1, "b": 2}
for key, value := range counts {
    fmt.Printf("%s: %d\n", key, value)
}
```

**Important:** Map iteration order is random in Go. If you need sorted output, you must sort the keys yourself.

### Sorting

The `sort` package provides sorting:

```go
import "sort"

// Sort a slice of strings alphabetically
words := []string{"banana", "apple", "cherry"}
sort.Strings(words)  // ["apple", "banana", "cherry"]

// Sort a slice of ints
nums := []int{3, 1, 2}
sort.Ints(nums)  // [1, 2, 3]
```

### Extracting Keys from a Map

To get sorted keys from a map:

```go
// Collect keys into a slice
keys := []string{}
for k := range counts {
    keys = append(keys, k)
}
// Sort the keys
sort.Strings(keys)
// Iterate in order
for _, k := range keys {
    fmt.Printf("%s: %d\n", k, counts[k])
}
```

### Command-Line Arguments as a Slice

Remember `os.Args`? It's a `[]string` — a slice of strings. You can slice it:

```go
// os.Args[0] is the program name
// os.Args[1:] is everything after the program name
words := os.Args[1:]
```

For more depth, see [Collections](../../concepts/08-collections.md).

---

## Task

Create a word frequency counter. The program takes words as command-line arguments, counts how many times each word appears, and prints the results sorted alphabetically.

### Setup

1. `go mod init wordfreq`
2. Create `main.go`

### Requirements

Usage: `./wordfreq <word1> [word2] ...`

The program should:

1. Count how many times each word appears in the arguments
2. Print each unique word with its count, one per line
3. Format: `word: count`
4. Sort output alphabetically by word

| Input | Output |
|---|---|
| `./wordfreq hello world hello go world hello` | `go: 1`<br>`hello: 3`<br>`world: 2` |
| `./wordfreq apple` | `apple: 1` |
| `./wordfreq a b a b a` | `a: 3`<br>`b: 2` |
| `./wordfreq` (no args) | `usage: wordfreq <word1> [word2] ...` |

**Rules:**

- Words are case-sensitive (`Hello` and `hello` are different words)
- Output must be sorted alphabetically
- Format is `word: count` (word, colon, space, count)
- One word per line

### Verify Your Solution

```bash
bash test.sh
go test -v
```

---

## Hints

<details>
<summary>Hint 1: Building the frequency map</summary>

```go
counts := map[string]int{}
for _, word := range os.Args[1:] {
    counts[word]++
}
```

The `++` works even if the key doesn't exist yet — Go initializes missing map values to the zero value (0 for int).

</details>

<details>
<summary>Hint 2: Sorting the keys</summary>

```go
keys := []string{}
for k := range counts {
    keys = append(keys, k)
}
sort.Strings(keys)
```

Don't forget to `import "sort"`.

</details>

<details>
<summary>Hint 3: Printing in order</summary>

```go
for _, k := range keys {
    fmt.Printf("%s: %d\n", k, counts[k])
}
```

</details>

---

## Further Reading

- [Go by Example: Maps](https://gobyexample.com/maps)
- [Go by Example: Slices](https://gobyexample.com/slices)
- [Go by Example: Range](https://gobyexample.com/range)
- [Go by Example: Sorting](https://gobyexample.com/sorting)
- [sort package documentation](https://pkg.go.dev/sort)
