---
name: interspec
description: >
  Write and edit InterSpec (.is) files â€” a declarative, text-based interface
  specification language for rapid UI wireframing, structural layout, and
  interaction prototyping. Use when the user asks to create UI wireframes,
  prototypes, interface specs, or .is files, or mentions InterSpec.
license: MIT
metadata:
  author: interspec-community
  version: "1.0"
---

# InterSpec Agent Skill

You are an expert InterSpec author. Use this skill whenever the user asks you to
create, modify, or review `.is` files.

## What InterSpec Is

InterSpec is a **declarative UI specification language**. It describes
*information architecture*, *structural layout*, and *interaction logic* â€”
**not** visual styling (no colors, fonts, shadows, or pixel values). Think of
it as executable wireframes.

**Golden rule:** If the user asks for visual styling, explain that InterSpec
is intentionally style-free and focus on structure and interaction. If they
need to communicate visual intent without breaking the no-styling rule, use
`@` hints instead.

## Critical Rules (Read First)

These are the rules agents get wrong most often. Internalise them before
writing any code.

### 1. No `$` prefix on variable access
Variables are referenced directly by name. The `${}` syntax is **only** used
inside string literals for interpolation.

```interspec
state userName = "Alice"
Text(userName)                     // direct access â€” NO $
Text("Hi ${userName}!")           // interpolation â€” REQUIRES ${}
```

### 2. Custom components are neutral (no wrapper node)
When you instantiate a custom component, children passed in the `{}` block
are **always appended to the end** of the component's internal declaration.
There are no implicit slots.

```interspec
component Card(title) {
    Text(title)
    Text("Footer")
}

Card("My Title") {
    Text("child")
}
// Renders: Text("My Title"), Text("Footer"), Text("child")
```

### 3. PascalCase for components/pages, camelCase for variables
This prevents parser ambiguity between a component call and a variable
reference.

### 4. Only two layout primitives: `row` and `column`
There is no `flex`, `grid`, `stack`, etc. Use `wrap: true` inside `row` for
wrapping.

### 5. No external I/O
InterSpec is sandboxed. No network, file system, or database access. Only
InterSpec logic executes.

### 6. Hints use `@` / `@* *@` — they survive stripping and guide implementers
Single-line hints start with `@` and run to end of line. Multi-line hints use
`@* ... *@` blocks. Hints are freeform text — no grammar, no validation. They
pass through `isc strip` unchanged and are ignored by deterministic transpilers.
Use hints to communicate visual hierarchy, responsive intent, accessibility,
spacing, animation, or any other implementer guidance.

```interspec
@ This is the primary action on the page — make it prominent
Button("Save") { ... }

@*
  On mobile: single column, full width.
  On tablet and up: two-column grid.
*@
row { ... }
```

## Syntax Quick Reference

### Hints
```interspec
@ Single-line hint — guides the implementer

@*
  Multi-line hint block.
  Everything between @* and *@ is guidance.
*@
```
Hints have no runtime effect. They survive `isc strip` and are ignored by
transpilers. Use them to communicate visual hierarchy, responsive intent,
accessibility, or spacing decisions.

### Variables
```interspec
state myText = "Hello"       // mutable
const maxItems = 10          // immutable
```
Dynamically typed: reassigning to a different type is valid.

### Components & Pages
```interspec
component MyButton(label) {
    Button() {
        text: label
        on click { Toast("Clicked!") }
    }
}

page Dashboard() {
    Text("Welcome")
}
```
The entry point is always `Main()`.

### Layout
```interspec
row {
    wrap: true
    collapse: true    // auto-collapse to column on narrow viewports
    Text("Item 1")
    Text("Item 2")
}

column {
    Text("Item 1")
}
```

### Alignment & Weight (properties inside the component block)
```interspec
Text("Centered") {
    align: (center, center)   // (vertical, horizontal)
    weight: both             // horizontal | vertical | both
}
```

### Control Flow
```interspec
state items = [1, 2, 3]
state show = true

for item in items {
    Text("Item: ${item}")
}

if show {
    Text("Visible")
} else {
    Text("Hidden")
}
```

### Events
```interspec
Button("Click me") {
    on click { navigate OtherPage() }
    on hover { Tooltip("Info") }
    on longpress { Dialog("Options") }
}

Input("Name") {
    on input  { log(value) }
    on commit { log("Confirmed: " + value) }
    on focus  { Toast("Focused") }
    on blur   { Toast("Blurred") }
    on key("Enter") { validate() }
    on key("Escape") { searchTerm = "" }
}

Form {
    on submit { validate() }
    // ... inputs ...
}
```

### Reactive Watchers
```interspec
state val = true
on change(val) {
    Text("Changed from ${old(val)} to ${val}")
}
```

### Navigation
```interspec
navigate UserProfile(userId: 123)
back()
```

## Built-in Components (quick list)

Button, Text, Input, Select, Checkbox, DatePicker, Toggle, Slider, Image, Icon, Alert,
Card, Modal, Dialog, Drawer, Toast, Tooltip, Table, Tabs, Accordion, TreeView, Badge, Link,
Progress, EmptyState, Breadcrumb, Stepper, Pagination, DropdownMenu, Divider, Section,
Form, FileUpload.

For the full catalog with parameters, events, and properties, see
[the component reference](references/CATALOG.md).

## File Structure & Imports
```interspec
import "/components/mybutton.is"
import "/components/" as components
from "/pages/main.is" import MainPage
import "https://cdn.example.com/ui/buttons.is"
```

## Authoring Checklist
- Entry page is `page Main()`. Add a top-of-file `@* ... *@` hint describing the overall purpose.
- All custom components use PascalCase; all variables use camelCase or snake_case.
- No `$` prefix on variable access; `${}` only inside strings.
- Only `row` and `column` for layout; use `wrap: true` for wrapping and `collapse: true` for responsive collapse.
- Children passed at instantiation append to the **end** of the component body.
- No styling properties (colors, fonts, spacing, pixel values). Use `@` hints if the implementer needs visual guidance.
- `for` loops iterate over arrays â€” never write unbounded loops.
- Use `@` for brief hints (one line) and `@* ... *@` for detailed guidance (multiple sentences).
- Prefer hints over comments for anything the implementer needs to see.

## Common Patterns

### List with interactive items
```interspec
@* Show items as cards with a select button.
   Each selection increments the count. *@
page Main() {
    state items = ["Apple", "Banana", "Cherry"]
    state selected = 0

    on change(selected) {
        Toast("Selected ${selected} items")
    }

    column {
        @ Compact cards — minimal padding between items
        for fruit in items {
            Card(fruit) {
                Button("Pick ${fruit}") {
                    on click { selected = selected + 1 }
                }
            }
        }
    }
}
```

### Form with validation
```interspec
@* This form collects user email and consent.
   Keep the layout simple — labels above inputs,
   with clear error states for validation. *@
page Main() {
    state email = ""
    state agreed = false

    column {
        @ Primary input — wide, full width
        Input("Email") {
            required: true
            on input { log("email: ${value}") }
        }
        Checkbox("I agree") {
            required: true
        }
        @ Primary CTA — most visually prominent button on the page
        Button("Submit") {
            on click { Dialog("Confirm submission") }
        }
    }
}
```

For more examples see [the examples file](references/EXAMPLES.md).
