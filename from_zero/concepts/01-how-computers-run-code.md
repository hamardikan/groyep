# 01 — How Computers Run Code

## What Is It

Every program you have ever used — your web browser, your terminal, grep — started as text a human wrote. That text, called **source code**, is meaningless to a computer. A CPU does not understand English, Go, Python, or any programming language. It understands one thing: **machine code** — a sequence of binary instructions specific to the hardware architecture it was built on.

The entire history of programming is the story of bridging the gap between what humans can read and what machines can execute.

## What Problem It Solves

Humans cannot write machine code productively. It is a stream of numbers — opcodes and operands — that directly manipulate registers, memory addresses, and hardware flags. Writing a program in machine code is like building a house by individually placing atoms.

We need a way to write instructions in a language humans can reason about, and then **translate** those instructions into something the CPU can execute. That translation process is the foundation of all software engineering.

## First Principles

### The CPU Execution Model

A CPU does three things in a loop, billions of times per second:

1. **Fetch** — read the next instruction from memory
2. **Decode** — figure out what the instruction means
3. **Execute** — do the thing (add numbers, move data, compare values, jump to a different instruction)

This is the **fetch-decode-execute cycle**. Every program, from a "Hello, World" to an operating system, is just a sequence of these tiny operations.

### Memory

Programs need a place to store data while they run. **Memory** (RAM) is that workspace — a giant array of bytes, each with an address. When your program creates a variable, it occupies some bytes in memory. When you read a file, its contents get loaded into memory. When you search for a pattern in text, both the pattern and the text sit in memory while the CPU compares them.

### Machine Code and Assembly

**Machine code** is the raw binary the CPU executes. Each CPU architecture (x86-64, ARM, RISC-V) has its own instruction set — its own "language."

**Assembly language** is a thin human-readable layer over machine code. Each assembly instruction maps (roughly) 1:1 to a machine instruction:

```asm
MOV  RAX, 1      ; move the value 1 into register RAX
MOV  RDI, 1      ; file descriptor 1 (stdout)
SYSCALL           ; ask the operating system to do something
```

You will not write assembly in this pathway. But knowing it exists helps you understand what your Go code becomes.

### The Translation Problem

We need to go from this:

```go
fmt.Println("hello")
```

To a sequence of machine instructions the CPU can execute. There are two fundamental approaches.

## Compiled vs Interpreted Languages

| Aspect | Compiled | Interpreted |
|--------|----------|-------------|
| **Translation** | Entire program translated before execution | Translated line-by-line during execution |
| **Output** | A standalone binary (machine code) | No binary — needs the interpreter to run |
| **Startup speed** | Fast — binary is ready to go | Slower — must parse and translate at runtime |
| **Runtime speed** | Generally faster — optimized at compile time | Generally slower — translation overhead |
| **Distribution** | Ship the binary | Ship the source + require the interpreter |
| **Error detection** | Many errors caught at compile time | Errors found at runtime |
| **Examples** | Go, C, C++, Rust | Python, Ruby, JavaScript* |

*JavaScript is a special case — modern engines (V8) use Just-In-Time (JIT) compilation, blurring the line.

### Compilation

A **compiler** reads your entire source code, analyzes it, optimizes it, and outputs a binary file containing machine code for a specific platform. Once compiled, you no longer need the source code or the compiler to run the program.

### Interpretation

An **interpreter** reads your source code at runtime, translating and executing one statement at a time. You always need the interpreter present to run the program.

### The Tradeoff

Compilation trades **build time** for **runtime performance**. You wait once during compilation, then every execution is fast. Interpretation gives you **instant feedback** (no build step) but pays the translation cost every single time you run.

## How Go Does It

Go is a **compiled language**. The Go toolchain compiles your source code into a **statically linked binary** — a single executable file with no external dependencies.

### The Go Build Process

```
┌─────────────┐     ┌──────────┐     ┌──────────────┐     ┌────────────┐
│ Source Code  │ ──▶ │ Compiler │ ──▶ │    Binary     │ ──▶ │  OS runs   │
│  (.go files) │     │ (go build)│     │ (machine code)│     │  the binary │
└─────────────┘     └──────────┘     └──────────────┘     └────────────┘
```

Step by step:

1. **You write** `.go` source files
2. **`go build`** invokes the Go compiler
3. **The compiler** parses your code, checks types, optimizes, and emits machine code
4. **The output** is a single binary file for your target OS and architecture
5. **The OS** loads the binary into memory and the CPU starts executing

```bash
# Write source code
vim main.go

# Compile it
go build -o groyep main.go

# Run the binary — no Go installation needed on the target machine
./groyep "pattern" file.txt
```

### What Makes Go's Compilation Special

- **Fast compilation** — Go was designed from the ground up for fast compile times. The dependency model ensures the compiler only looks at what it needs.
- **Static linking** — the binary includes everything it needs. No "install the runtime" step. No dependency hell on the target machine.
- **Cross-compilation** — you can compile for a different OS or architecture from your own machine:

```bash
# Compile for Linux on a Mac
GOOS=linux GOARCH=amd64 go build -o groyep-linux
```

- **Garbage collection** — unlike C, you do not manually manage memory. The Go runtime (compiled into your binary) handles allocation and cleanup.

### What the Binary Contains

A Go binary includes:
- Your compiled code (machine instructions)
- The Go **runtime** (garbage collector, goroutine scheduler, etc.)
- Any standard library packages you imported
- Debugging information (unless stripped)

This is why even a simple "Hello, World" Go binary is a few megabytes — it includes the runtime. This tradeoff (larger binary, zero dependencies) is deliberate.

## Tradeoffs

| You gain | You give up |
|----------|-------------|
| Fast execution — no translation at runtime | Must wait for compilation before running |
| Single binary deployment — ship one file | Larger file size (runtime included) |
| Type errors caught at compile time | Slower iteration loop vs interpreted languages |
| Cross-compilation for any platform | Must recompile for each target |
| No runtime dependency on the target machine | Cannot modify behavior at runtime as easily |

## Why This Matters for Grep

Grep is a command-line tool. It needs to:

- **Start fast** — when you type `grep pattern file.txt`, you do not want to wait for an interpreter to spin up. A compiled binary starts in milliseconds.
- **Run fast** — grep may scan millions of lines. Compiled code runs at near-hardware speed.
- **Be portable** — you want to copy the binary to any machine and have it work. Static compilation gives you that.
- **Have no dependencies** — the real `grep` ships as a single binary. Yours will too.

When you run `go build -o groyep`, you are turning your Go source code into a tool that works exactly like any other command on the system. It is not a script. It is not an app that needs a runtime. It is a native binary, indistinguishable from tools written in C.

## Further Reading

1. **[Code: The Hidden Language of Computer Hardware and Software](https://www.charlespetzold.com/code/)** — Charles Petzold
   The best book for understanding how computers work from first principles. Starts with flashlights and Morse code, ends with a working CPU. Read this if you want to truly understand the machine.

2. **[Ben Eater — Building an 8-bit Computer](https://www.youtube.com/playlist?list=PLowKtXNTBypGqImE405J2565dvjafglHU)** — YouTube playlist
   Watch someone build a CPU from scratch on breadboards. Nothing makes the fetch-decode-execute cycle click like seeing it happen in hardware.

3. **[How Go Compiles Down to Machine Code](https://go.dev/blog/go1.17)** — Go Blog
   Details on Go's compilation model and the SSA-based compiler backend.

4. **[Rob Pike — "Go at Google: Language Design in the Service of Software Engineering"](https://go.dev/talks/2012/splash.article)** — 2012
   Why Go was designed the way it was, including the emphasis on fast compilation.

5. **[Computer Science from the Bottom Up](https://bottomupcs.com/)** — Ian Wienand
   Free online book covering how programs interact with the operating system and hardware. Good companion material.

6. **[The Go Compiler](https://github.com/golang/go/tree/master/src/cmd/compile)** — Go source code
   If you are curious what the compiler itself looks like, it is written in Go. The README in this directory explains the compilation phases.

7. **[Putting the "Go" in "Go Binary"](https://www.youtube.com/watch?v=GOmVEm0XLWE)** — GopherCon talk
   A deep dive into what is inside a Go binary and how the linker works.

---

*Next: [02 — What Is Go](./02-what-is-go.md)*
