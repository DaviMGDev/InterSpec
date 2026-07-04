# InterSpec

A declarative, text-based interface specification language for rapid UI wireframing, structural layout, and interaction prototyping.

## What is InterSpec?

InterSpec lets you describe user interfaces without implementing the final visual version. It focuses purely on **information architecture**, **structural layout**, and **interaction logic** — no colors, fonts, shadows, or pixel-level styling.

Files use the `.is` extension.

```interspec
page Main() {
    state count = 0

    column {
        Text("Counter: ${count}")

        Button("Increment") {
            on click {
                count = count + 1
            }
        }
    }
}
```

### Core Principles

- **No Styling** — No colors, fonts, shadows, or granular pixel controls. Pure structure.
- **Focus on UX** — Information architecture and interaction logic first.
- **Sandboxed & Safe** — No external I/O (no network, file system, or database access). Only InterSpec logic runs.
- **Component-based** — Build reusable components with parameters, events, and composition.

## Language Highlights

| Feature | Syntax |
|---------|--------|
| Mutable state | `state name = "value"` |
| Constants | `const limit = 10` |
| Components | `component Card(title) { ... }` |
| Pages (navigable) | `page Dashboard() { ... }` |
| Layouts | `row { ... }` / `column { ... }` |
| Control flow | `for`, `if`/`else` |
| Events | `on click { ... }`, `on change(var) { ... }` |
| Navigation | `navigate PageName(param: value)` |
| String interpolation | `"Hello, ${name}!"` |

## `isc` — The CLI Tool

The `isc` command-line tool processes `.is` files. Currently it provides:

### `isc strip`

Removes all comments (single-line `//` and block `/* */`) from InterSpec source files, preserving string literals and code.

```bash
# Strip comments from a file (output to stdout)
isc strip input.is

# Strip and write to a file
isc strip -o clean.is input.is

# Pipe from stdin
cat input.is | isc strip
```

**Features:**
- Preserves comment delimiters inside string literals (`"http://example.com"` stays intact)
- Handles escaped quotes inside strings
- Collapses comment-only lines (no blank lines left behind)
- Reports unterminated block comments with file, line, and column
- Supports Unicode content

## Project Structure

```
InterSpec/
├── LANGUAGE.md          # Full language specification
├── example.is           # Comprehensive example demonstrating all features
└── isc/                 # CLI tool (Go)
    ├── main.go          # Entry point
    ├── cmd/
    │   └── strip.go     # `strip` command implementation
    ├── stripper/
    │   ├── stripper.go  # Comment stripping state machine
    │   └── stripper_test.go
    ├── testdata/        # Test fixtures
    └── go.mod
```

## Building `isc`

```bash
cd isc
go build -o isc .
```

## Running Tests

```bash
cd isc
go test ./...
```

## License

TBD
