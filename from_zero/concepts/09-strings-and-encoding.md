# Strings and Encoding

## What Is It

A string is a sequence of characters — but what is a "character"? The answer is more complex than it appears, and understanding it is essential for writing correct text-processing tools.

At the hardware level, computers only understand numbers. Text is an **encoding** — a mapping from numbers to characters. The history of text encoding is a story of expanding scope:

| Era | Encoding | Characters | Bytes per Character |
|-----|----------|------------|---------------------|
| 1963 | ASCII | 128 (English letters, digits, symbols) | 1 |
| 1980s | Latin-1, Shift-JIS, etc. | 256 per encoding (regional) | 1 |
| 1991 | Unicode | 149,186+ (every writing system) | Varies |
| 1992 | UTF-8 | All Unicode code points | 1–4 |

Go strings are **byte slices** under the hood. The language has first-class support for Unicode through UTF-8, which is not an accident — Go's creators invented UTF-8.

## What Problem It Solves

### ASCII: The Starting Point

ASCII (American Standard Code for Information Interchange) maps 128 characters to the numbers 0–127. It covers:

- Uppercase and lowercase English letters (A–Z, a–z)
- Digits (0–9)
- Punctuation and symbols
- Control characters (newline, tab, etc.)

That is 128 characters. One byte (which can hold 0–255) is more than enough. For English text, ASCII works perfectly.

But the world has more than English. Japanese has thousands of kanji. Arabic flows right-to-left. Emoji exist. ASCII cannot represent any of this.

### The Chaos of Regional Encodings

Before Unicode, every region created its own encoding. Latin-1 for Western Europe. Shift-JIS for Japanese. Big5 for Traditional Chinese. GB2312 for Simplified Chinese. KOI8-R for Russian.

This created a fundamental problem: **the same byte could mean different characters in different encodings**. Byte `0xC0` is `À` in Latin-1 but `タ` (katakana ta) in a Japanese encoding. Without knowing which encoding a file uses, you cannot interpret its contents. Files were routinely garbled when transferred between systems.

### Unicode: One Set to Rule Them All

Unicode solves this by assigning a **unique number** (called a **code point**) to every character in every writing system. There is exactly one number for `A` (U+0041), one for `漢` (U+6F22), one for `🔍` (U+1F50D).

But Unicode is a character set, not an encoding. The question remains: how do you store code point U+1F50D (which requires at least 3 bytes) in a file?

### UTF-8: The Elegant Solution

UTF-8 is a **variable-width encoding** that represents each Unicode code point in 1 to 4 bytes:

| Code Point Range | Bytes | Bit Pattern | Example |
|-----------------|-------|-------------|---------|
| U+0000 – U+007F | 1 | `0xxxxxxx` | `A` = `0x41` |
| U+0080 – U+07FF | 2 | `110xxxxx 10xxxxxx` | `ñ` = `0xC3 0xB1` |
| U+0800 – U+FFFF | 3 | `1110xxxx 10xxxxxx 10xxxxxx` | `漢` = `0xE6 0xBC 0xA2` |
| U+10000 – U+10FFFF | 4 | `11110xxx 10xxxxxx 10xxxxxx 10xxxxxx` | `🔍` = `0xF0 0x9F 0x94 0x8D` |

The genius of UTF-8:

1. **Backward compatible with ASCII** — any valid ASCII file is also valid UTF-8
2. **Self-synchronizing** — you can jump to any byte and find the start of a character
3. **No byte-order issues** — unlike UTF-16, there is no endianness problem
4. **Compact for English text** — ASCII characters still take just 1 byte

Rob Pike and Ken Thompson invented UTF-8 in 1992 on a placemat at a New Jersey diner. These are the same people who created Go. UTF-8 is in Go's DNA.

## First Principles

### Characters, Code Points, and Bytes

These three concepts are distinct:

- **Character**: What a human reads — the letter "A", the emoji "🔍"
- **Code point**: The Unicode number for that character — U+0041, U+1F50D
- **Bytes**: How the code point is stored in memory using an encoding like UTF-8

A single character can be multiple code points (e.g., `é` can be `e` + combining accent), and a single code point can be multiple bytes. Never assume a 1:1 relationship between any of these.

### Why Variable Width Matters

Fixed-width encodings (like UTF-32, which uses 4 bytes per character) waste space. English text in UTF-32 uses 4x the memory compared to ASCII. UTF-8's variable width means English text stays compact while still supporting every character in existence.

The tradeoff: you cannot jump to "the nth character" in O(1) time. Finding the 100th character means scanning from the beginning, counting characters. This is a fundamental design decision that affects every string operation.

## How Go Does It

### Strings Are Byte Slices

In Go, a `string` is a read-only slice of bytes. It is **not** an array of characters:

```go
s := "Hello"
fmt.Println(len(s))    // 5 — five bytes (ASCII, so one byte per character)

s2 := "Hello, 世界"
fmt.Println(len(s2))   // 13 — NOT 9 characters! "世" and "界" are 3 bytes each
```

You can index into a string, but you get **bytes**, not characters:

```go
s := "Hello, 世界"
fmt.Println(s[0])      // 72 — the byte value of 'H'
fmt.Println(s[7])      // 228 — the FIRST byte of '世', not the character itself
```

### The `rune` Type

Go introduces the `rune` type to represent a single Unicode code point:

```go
// rune is an alias for int32
var r rune = '世'
fmt.Println(r)          // 19990 — the Unicode code point for '世'
fmt.Printf("%c\n", r)   // 世 — printed as a character
fmt.Printf("U+%04X\n", r) // U+4E16
```

A `rune` can represent any Unicode code point. It always takes 4 bytes in memory (int32), but this is for in-memory processing — not for storage on disk (where UTF-8's variable width is used).

### Indexing vs Ranging

This is one of Go's most important string behaviors:

```go
s := "café"

// Indexing gives BYTES
for i := 0; i < len(s); i++ {
    fmt.Printf("byte[%d] = %d (%c)\n", i, s[i], s[i])
}
// byte[0] = 99 (c)
// byte[1] = 97 (a)
// byte[2] = 102 (f)
// byte[3] = 195 (Ã)  ← first byte of 'é'
// byte[4] = 169 (©)  ← second byte of 'é'

// Ranging gives RUNES
for i, r := range s {
    fmt.Printf("rune at byte %d = %c (U+%04X)\n", i, r, r)
}
// rune at byte 0 = c (U+0063)
// rune at byte 1 = a (U+0061)
// rune at byte 2 = f (U+0066)
// rune at byte 3 = é (U+00E9)  ← correct character, byte index 3
```

Notice: `range` on a string yields `(byte_index, rune)` pairs. The byte index jumps from 3 to the next character — there is no index 4 in the range output because `é` is 2 bytes.

### `len()` Returns Bytes, Not Characters

```go
s := "🔍 search"
fmt.Println(len(s))                    // 11 — four bytes for 🔍, one for space, six for "search"
fmt.Println(utf8.RuneCountInString(s)) // 8 — eight characters
```

This is a common source of bugs. Always ask yourself: do I need the byte count or the character count?

### Converting Between Strings, Bytes, and Runes

```go
s := "Hello, 世界"

// String → byte slice
b := []byte(s)
fmt.Println(b) // [72 101 108 108 111 44 32 228 184 150 231 ... ]

// Byte slice → string
s2 := string(b)

// String → rune slice
r := []rune(s)
fmt.Println(len(r))  // 9 — nine characters
fmt.Println(r[7])    // 30028 — code point for '界'

// Rune slice → string
s3 := string(r)
```

Converting to `[]rune` allocates a new slice. For large strings, this is expensive. Prefer `range` when you just need to iterate.

### The `strings` Package

The `strings` package is your primary tool for string manipulation:

```go
import "strings"

s := "Hello, World"

// Searching
strings.Contains(s, "World")      // true
strings.HasPrefix(s, "Hello")     // true
strings.HasSuffix(s, "World")     // true
strings.Index(s, "World")         // 7 (byte position)
strings.Count(s, "l")             // 3

// Case conversion
strings.ToLower(s)                // "hello, world"
strings.ToUpper(s)                // "HELLO, WORLD"

// Splitting and joining
parts := strings.Split("a:b:c", ":") // ["a", "b", "c"]
joined := strings.Join(parts, "-")    // "a-b-c"

// Trimming
strings.TrimSpace("  hello  ")        // "hello"
strings.Trim("***hello***", "*")      // "hello"

// Replacing
strings.Replace(s, "World", "Go", 1)  // "Hello, Go" (replace first)
strings.ReplaceAll(s, "l", "L")       // "HeLLo, WorLd" (replace all)
```

These functions are UTF-8 aware. `strings.ToUpper("café")` correctly returns `"CAFÉ"`, not garbled bytes.

### The `unicode/utf8` Package

For low-level UTF-8 operations:

```go
import "unicode/utf8"

s := "Hello, 世界"

// Count characters (runes), not bytes
utf8.RuneCountInString(s)              // 9

// Decode the first rune
r, size := utf8.DecodeRuneInString(s)  // r='H', size=1
r, size = utf8.DecodeRuneInString(s[7:]) // r='世', size=3

// Check if bytes are valid UTF-8
utf8.ValidString(s)                     // true
utf8.ValidString("\xff\xfe")            // false

// Check if a byte is a rune start
utf8.RuneStart(s[0])                    // true (start of 'H')
utf8.RuneStart(s[8])                    // false (middle of '世')
```

### The `strings.Builder` for Efficient Construction

When building strings incrementally, use `strings.Builder` to avoid allocating many intermediate strings:

```go
var b strings.Builder
for i := 0; i < 1000; i++ {
    b.WriteString("line ")
    b.WriteString(strconv.Itoa(i))
    b.WriteByte('\n')
}
result := b.String() // one allocation at the end
```

## Tradeoffs

| Approach | Benefit | Cost |
|----------|---------|------|
| Byte-based processing (`s[i]`, `len(s)`) | Fast, zero allocation | Breaks on multi-byte characters |
| Rune-based processing (`range s`, `[]rune`) | Correct for all Unicode | Slower, `[]rune` conversion allocates |
| `strings.Contains` (byte-level search) | Very fast | Cannot do rune-aware partial matching |
| `utf8.RuneCountInString` | Accurate character count | Must scan entire string — O(n) |
| Fixed `[]rune` indexing | O(1) character access | O(n) conversion cost, 4x memory for ASCII |
| `strings.Builder` | Efficient string construction | Slightly more verbose than `+=` |

### When to Use Which

```go
// USE BYTES when you know the content is ASCII or you need raw speed
if line[0] == '#' { // checking for comment lines — '#' is ASCII
    continue
}

// USE RUNES when correctness across languages matters
for _, r := range line {
    if unicode.IsLetter(r) {
        // Works for Latin, Cyrillic, CJK, etc.
    }
}

// USE strings PACKAGE for most string operations
if strings.Contains(line, pattern) {
    // Already handles UTF-8 correctly
}
```

## Common Gotchas

### Gotcha 1: String Length

```go
s := "café"
fmt.Println(len(s))                    // 5 (bytes), not 4 (characters)
fmt.Println(utf8.RuneCountInString(s)) // 4 (characters)
```

### Gotcha 2: Slicing Splits a Character

```go
s := "café"
fmt.Println(s[:4])  // "caf" + invalid byte — this cuts 'é' in half!
// Output: "caf\xc3" — garbled
```

To safely slice, you need to know character boundaries. Use `range` to find them:

```go
s := "café"
i := 0
for j, _ := range s {
    if i == 3 { // want first 3 characters
        fmt.Println(s[:j]) // "caf" — up to but not including 'é'
        break
    }
    i++
}
```

### Gotcha 3: String Comparison with Different Unicode Representations

The character `é` can be represented two ways in Unicode:

- **Precomposed**: U+00E9 (a single code point, `é`)
- **Decomposed**: U+0065 U+0301 (`e` + combining acute accent)

These look identical on screen but are different byte sequences. Go compares strings byte-by-byte, so they are **not equal**:

```go
precomposed := "caf\u00e9"       // "café" — one code point for é
decomposed := "cafe\u0301"       // "café" — e + combining accent

fmt.Println(precomposed == decomposed) // false!
fmt.Println(len(precomposed))          // 5 bytes
fmt.Println(len(decomposed))          // 6 bytes
```

For Unicode normalization, use `golang.org/x/text/unicode/norm`.

### Gotcha 4: Iterating Bytes When You Mean Runes

```go
// BUG: this counts bytes, not characters
func charCount(s string) int {
    count := 0
    for i := 0; i < len(s); i++ {
        count++
    }
    return count
}

// CORRECT: count runes
func charCount(s string) int {
    return utf8.RuneCountInString(s)
}

// ALSO CORRECT: range counts runes
func charCount(s string) int {
    count := 0
    for range s {
        count++
    }
    return count
}
```

## Why This Matters for Grep

Grep is a text search tool. Every line it processes is a string. Getting string handling wrong means:

1. **Incorrect match positions**: If you calculate match positions using byte offsets but display them as character positions, colorized output will break on non-ASCII text.

2. **Case-insensitive matching across languages**: `strings.ToLower` handles Unicode case folding, but `é` vs `e` + combining accent requires normalization.

3. **Pattern matching on multi-byte characters**: The regex engine (`regexp`) works correctly on UTF-8 strings — `\w` matches Unicode word characters, `.` matches any rune (not any byte). But if you do manual byte-level processing alongside regex, the two can disagree.

4. **Line length counting**: If your grep reports line length (for column-based features), you need to decide: byte length or character length?

```go
// Grep uses byte-level operations where safe and rune-level where needed
func highlightMatches(line string, re *regexp.Regexp) string {
    // FindAllStringIndex returns BYTE positions — this is correct
    // because Go strings are bytes and slicing by bytes works with
    // properly aligned UTF-8 (which regexp guarantees)
    matches := re.FindAllStringIndex(line, -1)
    if matches == nil {
        return line
    }

    var buf strings.Builder
    prev := 0
    for _, m := range matches {
        buf.WriteString(line[prev:m[0]])     // bytes before match
        buf.WriteString("\033[1;31m")         // red
        buf.WriteString(line[m[0]:m[1]])      // the match
        buf.WriteString("\033[0m")            // reset
        prev = m[1]
    }
    buf.WriteString(line[prev:])
    return buf.String()
}
```

The key insight: Go's `regexp` package returns byte offsets that align to UTF-8 boundaries, so slicing strings at those offsets is safe. This is by design — Go's regex operates on UTF-8 encoded strings and always matches at valid rune boundaries.

## Further Reading

1. **[Joel Spolsky — "The Absolute Minimum Every Software Developer Absolutely, Positively Must Know About Unicode and Character Sets (No Excuses!)"](https://www.joelonsoftware.com/2003/10/08/the-absolute-minimum-every-software-developer-absolutely-positively-must-know-about-unicode-and-character-sets-no-excuses/)** — The classic essay that explains why encoding matters. If you read one thing on this list, make it this one.

2. **[Go Blog — "Strings, bytes, runes and characters in Go"](https://go.dev/blog/strings)** — Rob Pike's authoritative explanation of how Go handles strings. Covers the distinction between bytes and runes, the `range` keyword behavior, and the design decisions behind Go's string model.

3. **[Rob Pike — "UTF-8 History"](https://www.cl.cam.ac.uk/~mgk25/ucs/utf-8-history.txt)** — The story of how UTF-8 was designed on a placemat at a New Jersey diner. A short, fascinating read about one of computing's most important encoding decisions.

4. **[The Unicode Consortium](https://home.unicode.org/)** — The official source for the Unicode standard. The code charts and FAQ are useful references when you need to understand specific character properties.

5. **[Go `strings` package documentation](https://pkg.go.dev/strings)** — The complete API reference for string manipulation functions. You will use this package constantly.

6. **[Go `unicode/utf8` package documentation](https://pkg.go.dev/unicode/utf8)** — Low-level UTF-8 encoding and decoding functions. Essential when you need to work at the byte level while respecting character boundaries.

7. **["Programming in Go" by Mark Summerfield — Chapter 3: Strings](https://www.qtrac.eu/gobook.html)** — A thorough treatment of Go's string handling with practical examples, covering the `strings`, `strconv`, `unicode`, and `unicode/utf8` packages.

8. **[Go `unicode` package documentation](https://pkg.go.dev/unicode)** — Functions for testing Unicode properties: `IsLetter`, `IsDigit`, `IsSpace`, `ToUpper`, `ToLower`. Useful for character classification beyond ASCII.

---

*Previous: [08 — Collections](./08-collections.md) · Next: [10 — I/O and Streams](./10-io-and-streams.md)*
