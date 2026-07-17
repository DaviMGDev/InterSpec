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

## Agent Skills

InterSpec comes with a set of agent skills that guide AI coding agents in working with `.is` files:

| Skill | Purpose |
|-------|---------|
| **interspec-reference** | Language syntax reference, component catalog, event catalog |
| **interspec-write** | Patterns, anti-patterns, and design judgment for authoring `.is` files |
| **interspec-consume** | Implementing `.is` files into production frontend code (HTML/CSS/JS, React, Vue) |
| **interspec-verify** | Verifying that a frontend implementation conforms to its `.is` specification |
| **interspec-reverse** | Extracting `.is` specifications from existing UIs |

Together they form a complete workflow: **write → consume → verify → reverse**.

## Project Structure

```
InterSpec/
├── LANGUAGE.md          # Full language specification
├── README.md            # This file
├── CONTRIBUTING.md      # Contribution guidelines
├── LICENSE              # MIT license
├── assets/              # Demo assets and design tokens
│   ├── demo.css
│   ├── demo.js
│   └── tokens.json
└── skills/              # Agent skill definitions
    ├── interspec-consume/
    ├── interspec-reference/
    ├── interspec-reverse/
    ├── interspec-verify/
    └── interspec-write/
```

## Getting Started

1. Read the **[full language specification](LANGUAGE.md)** — covers syntax, components, events, hints, and viewport safety.
2. Browse the **agent skills** under `skills/` — each skill teaches an AI agent how to work with `.is` files.
3. Explore the **built-in catalog** ([skill reference](skills/interspec-reference/references/CATALOG.md)) for the complete list of components, events, and properties.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on contributing to the specification, skills, and documentation.

## License

MIT — see [LICENSE](LICENSE) for details.
