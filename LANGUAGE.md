# InterSpec Language Specification

**File Extension:** `.is`  
**Purpose:** A declarative, text-based interface specification language focused on rapid UI wireframing, structural layout, and interaction prototyping.

## 1. Core Philosophy
InterSpec is designed to declare user interfaces without implementing the final visual version. 
* **No Styling:** There are no colors, fonts, shadows, or granular pixel controls. 
* **Focus on UX:** The language focuses purely on information architecture, structural layout, and interaction logic.
* **Sandboxed & Safe:** While Turing-complete, it has no external I/O (no network, file system, or database). It only executes InterSpec logic.

## 2. Syntax Basics
```interspec
// This is a single-line comment

/* 
  This is a 
  multi-line comment 
*/
```

## 3. Variables, State, and Data Types
InterSpec is dynamically and weakly typed. Supported data types include: `string`, `number`, `boolean`, `array`, `object`, and `any`.

### Declaration
* **Mutable variables:** Declared with `state`.
* **Immutable variables:** Declared with `const`.

```interspec
state myText = "Hello"
const maxItems = 10

myText = 123 // Valid: dynamic typing allows reassignment to a different type
```

### Variable Access & Interpolation
* **Direct Access:** Variables are accessed directly by their name. **No `$` prefix is used.**
* **String Interpolation:** The `${}` syntax is strictly reserved for injecting variables into strings.

```interspec
state userName = "Alice"

// Direct access (No $)
Text(userName) 

// String interpolation (Requires ${})
Text("Welcome, ${userName}!") 
```
*Note: To prevent parser ambiguity, use **PascalCase** for Components/Pages (e.g., `MyButton`) and **camelCase/snake_case** for variables (e.g., `userName`, `item_count`).*

## 4. Components and Pages
Components and Pages are the building blocks of the UI. Pages are simply components that can be navigated to. The entry point of any application is the `Main()` page.

### Declaration
```interspec
component MyButton(label) {
    Button() {
        text: label
        on click {
            Toast("Button clicked!") 
        }
    }
}

page Dashboard() {
    Text("Welcome to the dashboard")
}
```

### Instantiation and The "No Magic" Rule
When you instantiate a component, the code block `{ ... }` is used to **append children, override properties, or attach events**. 

**Crucial Rule:** Custom components are "neutral" (they do not render a wrapper node). Children passed during instantiation are **always appended to the very end** of the component's internal declaration. There are no implicit "slots" or magic placements.

```interspec
component Card(title) {
    Text(title)
    Text("Footer")
}

// Instantiation
Card("A Title") {
    Text("a child")
}

/* 
  RENDERED OUTPUT ORDER:
  1. Text("A Title")      <- From declaration
  2. Text("Footer")       <- From declaration
  3. Text("a child")      <- Appended from instantiation block
*/
```

### Navigation and Page Parameters
Pages can accept parameters, which are passed during navigation.
```interspec
// Triggering navigation
navigate UserProfile(userId: 123)

// Defining the page
page UserProfile(userId) {
    Text("Viewing user: ${userId}")
}
```

## 5. Layout and Alignment
InterSpec provides only two structural layouts: `row` and `column`. Granular spacing (like `gap` or `margin`) is intentionally omitted to maintain the wireframing philosophy.

### Layouts
```interspec
row {
    wrap: true      // Optional: allows items to wrap to the next line
    collapse: true   // Optional: auto-collapse to column on narrow viewports
    scrollable: true // Optional: content scrolls when it exceeds available space
    Text("Item 1")
    Text("Item 2")
}

column {
    scrollable: true // Optional: content scrolls when it exceeds available height
    Text("Item 1")
    Text("Item 2")
}
```

When `scrollable: true`, the implementer **must** constrain the container to a bounded height (viewport or parent) and apply `overflow: auto` (or equivalent). This is a structural signal — it tells the implementer that the container's content may exceed available space and must scroll rather than overflow.

### Component Properties
Components can be aligned and weighted within their parent layout.
```interspec
Text("This is a text") {
    align: (center, center) // (vertical, horizontal) -> top/bottom/center, left/right/center
    weight: both            // horizontal | vertical | both
}
```

## 6. Control Flow
InterSpec includes standard control flow structures. `for` loops are strictly enforced by the runtime to prevent infinite execution (e.g., via iteration limits or timeouts).

```interspec
state items = [1, 2, 3]
state condition = true 

for item in items {
    Text("Item: ${item}")
}

if condition {
    Text("Condition is true")
} else {
    Text("Condition is false")
}
```

## 7. Events and Reactivity
Components trigger built-in events, and state changes can be observed reactively.

### Standard Events
Components fire events in response to user interaction. The table in [Section 8](#8-built-in-catalog) lists all standard events and which components they apply to.

```interspec
Button("Hover me") {
    on hover {
        Tooltip("This is a button")
    }
    on click {
        navigate AnotherPage()
    }
}

Input("Enter your name") {
    on input {
        log("Current value: ${value}")
    }
    on focus {
        Toast("Input focused")
    }
    on blur {
        Toast("Input lost focus")
    }
}

Modal("Confirm") {
    on open {
        log("Modal opened")
    }
    on close {
        log("Modal closed")
    }
}
```

### Reactive State Watchers
Use `on change(variable)` to react to state updates. The `old(variable)` function provides the previous value.
```interspec
state val = true 

on change(val) {
    Text("Value changed from ${old(val)} to ${val}")
}

Button("Toggle") {
    on click {
        val = !val // Triggers the on change event
    }
}
```

## 8. Built-in Catalog

### Components

| Component | Parameters | Description |
|-----------|------------|-------------|
| `Button` | `label` | A clickable button that triggers actions. |
| `Text` | `content` | Renders a string of text. |
| `Input` | `placeholder` | A single-line text input field for user entry. |
| `DatePicker` | `placeholder` | A date selection input that opens a calendar picker for choosing a date. Use `min` and `max` properties to constrain the date range. |
| `Select` | `options` | A dropdown picker for selecting a value. Pass an array of strings or objects. Distinct from `DropdownMenu` which triggers actions. |
| `Checkbox` | `label` | A toggleable checkbox with an adjacent label. |
| `Toggle` | `label` | A boolean on/off toggle switch with a label. |
| `Slider` | — | A range slider for numeric selection. |
| `Image` | `src` | An image placeholder. Takes a URL string for `src`. |
| `Icon` | `name` | A semantic icon indicator (e.g., `"settings"`, `"user"`, `"search"`). The actual icon set is determined by the runtime. |
| `Alert` | `message` | An inline system message for success, error, warning, or info states. Use the `variant` property to set the semantic type. |
| `Card` | `title` | A content container that groups related information. Card is a neutral component (no wrapper node) — children are appended after the title. |
| `Modal` | `title` | A modal overlay dialog that blocks interaction with the rest of the page. |
| `Dialog` | `title` | A confirmation or information dialog that requires user action. |
| `Drawer` | `title` | A side panel that slides in from the edge of the screen for navigation, filters, or supplementary content. Use the `side` property to set the origin edge (`left` or `right`). |
| `Toast` | `message` | A brief, auto-dismissing pop-up notification. |
| `Tooltip` | `content` | A hover-revealed tooltip that displays contextual information. |
| `Table` | `columns`, `rows` | Tabular data display. `columns` is an array of header strings. `rows` is an array of row arrays. |
| `Tabs` | `tabs` | Tabbed content panel. `tabs` is an array of label strings. Children are panel bodies, matched by position to the labels. |
| `Accordion` | `items` | Expandable sections. `items` is an array of title strings. Children are section bodies, matched by position. Only one section is expanded at a time. |
| `TreeView` | `items` | A hierarchical tree display for nested data such as file systems, org charts, or category trees. Children define the content for each node level. |
| `Badge` | `label` | Inline indicator for status, count, or category. Use `variant` for semantic type. |
| `Link` | `label` | A navigational element distinct from action `Button`. Use `navigate` via `on click` to specify the target page. No `href` — InterSpec has no network access. |
| `Progress` | `value` | Read-only progress indicator. `value` is a number. Use `max` property to set the upper bound (default 100). |
| `EmptyState` | `message` | Placeholder for empty lists or search results. Children can provide a recovery action (e.g., a `Button` to create an item). |
| `Breadcrumb` | `items` | Navigation path trail. `items` is an array of label strings. The last item represents the current page. |
| `Stepper` | `steps` | Multi-step flow indicator (read-only). `steps` is an array of label strings. Use `current` property to mark the active step. |
| `Pagination` | — | Page navigation control for paginated lists. Use `current` and `total` properties. |
| `DropdownMenu` | `label` | A command menu that triggers actions (not a form input). Distinct from `Select` which returns a value. Children define the menu items as `Button` components. |
| `Divider` | — | A horizontal divider that visually separates content sections. Has no children and no interactive properties — purely structural. |
| `Section` | `title` | Structural content grouping under a heading. Groups related components without implying a card/border — purely architectural. |
| `Form` | — | Groups related inputs under a submission action. Fires a `submit` event when the user submits the form. |
| `FileUpload` | `label` | A file upload control that opens a file picker dialog or acts as a drop zone. Use `accept` to filter file types and `multiple` to allow multiple files. |

#### Usage Examples

```interspec
// Card — groups related content
Card("User Profile") {
    Text("Name: Alice")
    Text("Role: Admin")
}

// Input — text entry field
Input("Enter your name...")

// Select — dropdown for picking a value
Select(["Apple", "Banana", "Cherry"])

// Image — placeholder with source URL
Image("https://example.com/photo.png")

// Icon — semantic indicator
Icon("settings")

// Alert — inline system message
Alert("File saved successfully") {
    variant: success
}

// Table — data display
Table(["Name", "Role"], [
    ["Alice", "Admin"],
    ["Bob", "User"]
])

// Tabs — tabbed panels
Tabs(["First", "Second"]) {
    Text("Content for first tab")
    Text("Content for second tab")
}

// Accordion — expandable sections
Accordion(["FAQ 1", "FAQ 2"]) {
    Text("Answer to FAQ 1")
    Text("Answer to FAQ 2")
}

// Badge — status indicator
Badge("New") {
    variant: success
}

// Link — navigational (no href — use navigate)
Link("View Profile") {
    on click {
        navigate UserProfile(userId: 42)
    }
}

// Progress — read-only indicator
Progress(60) {
    max: 100
}

// EmptyState — nothing here yet
EmptyState("No results found") {
    Button("Create one") {
        on click { navigate NewItem() }
    }
}

// Breadcrumb — path trail
Breadcrumb(["Home", "Products", "Details"])

// Stepper — multi-step progress
Stepper(["Cart", "Shipping", "Payment"]) {
    current: 0
}

// Pagination — page control
Pagination {
    current: 1
    total: 10
}

// DropdownMenu — action menu
DropdownMenu("Actions") {
    Button("Edit") { on click { navigate Editor() } }
    Button("Delete") { on click { Dialog("Confirm delete") } }
}

// Drawer — side panel (slides in from the right by default, use `side: left` for left edge)
Drawer("Navigation") {
    side: left
    on open { log("Drawer opened") }
    on close { log("Drawer closed") }

    Text("Sidebar content goes here")
    Button("Close") {
        on click { toggle(drawerOpen) }
    }
}

// Section — structural grouping
Section("User Details") {
    Text("Name: Alice")
    Text("Role: Admin")
}

// Divider — horizontal separator
Divider()

// Form — grouped inputs with submit
Form {
    Input("Email") {
        required: true
        error: "Please enter a valid email"
    }
    Button("Submit") {
        on click { validate() }
    }
}

// FileUpload — file picker
FileUpload("Upload files") {
    accept: "image/*"
    multiple: true
    on input {
        log("Files selected")
    }
    on commit {
        log("Upload confirmed")
    }
}
```

### Events

| Event | Applies To | Description |
|-------|------------|-------------|
| `click` | `Button`, `Checkbox`, `Toggle`, `Card` | Fired when the user clicks or taps the component. |
| `hover` | `Button`, `Tooltip`, `Card`, `Icon` | Fired when the user hovers the cursor over the component. |
| `input` | `Input`, `Select`, `DatePicker`, `FileUpload` | Fired on every value change (e.g., each keystroke in an `Input`, date selection in a `DatePicker`, file selection in `FileUpload`). |
| `commit` | `Input`, `Select`, `DatePicker`, `FileUpload` | Fired when the user confirms a value (blur after change, Enter key, dropdown close, or date/file confirmation). Distinct from `input` which fires per interaction. |
| `focus` | `Input`, `Button`, `Select`, `DatePicker` | Fired when the component receives focus. |
| `blur` | `Input`, `Button`, `Select`, `DatePicker` | Fired when the component loses focus. |
| `open` | `Modal`, `Dialog`, `Drawer`, `Toast`, `TreeView` | Fired when the component becomes visible, or when a `TreeView` branch is expanded. |
| `close` | `Modal`, `Dialog`, `Drawer`, `Toast`, `TreeView` | Fired when the component is dismissed or hidden, or when a `TreeView` branch is collapsed. |
| `submit` | `Form` | Fired when the form is submitted (via Enter key or submit button). |
| `key` | `Input` | Fired on key press. Takes a key name argument: `on key("Enter") { ... }`, `on key("Escape") { ... }`. |
| `longpress` | `Button`, `Card`, any interactive component | Fired on touch-hold (mobile interaction). |
| `reachEnd` | scrollable `column`, `row` | Fired when the user scrolls to the end of the container. Use for infinite scroll or "load more" patterns. |

### Actions

| Action | Syntax | Description |
|--------|--------|-------------|
| `navigate` | `navigate PageName(param: value)` | Navigate to another page, optionally passing parameters. |
| `back` | `back()` | Navigate to the previous page. No-op if there is no history. |
| `toggle` | `toggle(variable)` | Shortcut for `variable = !variable`. Works with any boolean state variable. |
| `log` | `log(message)` | Print a debug message to the runtime console. No effect in production. |
| `validate` | `validate()` | Trigger form validation programmatically. Checks all `required` and `error` constraints within the enclosing `Form`. |
| `reset` | `reset(variable)` | Reset a state variable to its declared initial value. Consistent with `toggle()` as a convenience action. |
| `delay` | `delay(ms, action)` | Execute an action after a delay in milliseconds. Runtime-enforced to prevent infinite chains (same guard as `for` loops). |

### Component Properties

| Property | Applies To | Values | Description |
|----------|------------|--------|-------------|
| `align` | Any component in a `row` or `column` | `(vertical, horizontal)` where vertical: `top`, `center`, `bottom`; horizontal: `left`, `center`, `right` | Alignment within the parent layout. |
| `weight` | Any component in a `row` or `column` | `horizontal`, `vertical`, `both` | How the component fills available space. |
| `wrap` | `row` layout | `true`, `false` | Whether items wrap to the next line when they exceed the container width. |
| `collapse` | `row` layout | `true`, `false` | Whether the row automatically collapses to a column on narrow viewports. The runtime chooses the breakpoint. |
| `placeholder` | `Input` | `string` | Hint text displayed inside the field when it is empty. |
| `required` | `Input`, `Select`, `Checkbox` | `true`, `false` | Marks the field as required for form validation. |
| `side` | `Drawer` | `left`, `right` | Which edge the drawer slides in from. Default is `right`. |
| `accept` | `FileUpload` | `string` | Accepted file MIME types (e.g., `"image/*"`, `".pdf"`). |
| `multiple` | `FileUpload` | `true`, `false` | Whether multiple files can be selected at once. Default is `false`. |
| `disabled` | Any interactive component | `true`, `false` | Disables the component, preventing interaction. |
| `loading` | `Button`, `Card`, `Table`, `Image` | `true`, `false` | Indicates an async operation is in progress. Implies `disabled` when `true`. |
| `error` | `Input`, `Select` | `true`, `false`, or a message string | Marks the field as having a validation error. When a string is provided, it is shown as the error message. |
| `variant` | `Alert`, `Badge` | `info`, `success`, `warning`, `error` | Semantic variant indicating the type or severity. |
| `src` | `Image` | `string` | Source URL for the image content. |
| `min` | `DatePicker`, `Progress` | `number` | Lower bound. For `Progress`, defaults to 0. For `DatePicker`, a minimum date string. |
| `max` | `DatePicker`, `Progress` | `number` | Upper bound. For `Progress`, defaults to 100. For `DatePicker`, a maximum date string. |
| `current` | `Stepper`, `Pagination` | `number` (0-indexed for Stepper, 1-indexed for Pagination) | The active step or page. |
| `total` | `Pagination` | `number` | Total number of pages. |

## 9. Modularity and Imports
InterSpec supports modular file structures. You can import specific files, entire folders, or even remote files via URL. Remote imports are strictly sandboxed (they can only contain InterSpec code).

### File Structure Example
```text
myapp/
    index.is 
    components/
        mybutton.is
    pages/
        main.is
        anotherpage.is 
```

### Import Syntax
```interspec
// Import a file (Namespace becomes the filename without extension)
import "/components/mybutton.is" 

// Import a folder (Namespace becomes the folder name)
import "/components/"

// Import with an alias
import "/components/mybutton.is" as customBtn

// Import specific declarations
from "/pages/anotherpage.is" import AnotherPage, SomeState

// Import directly from a URL (CDN or Git raw)
import "https://example.com/ui-library/buttons.is"
```

***

### Complete Example

```interspec
import "/components/Card.is"

page Main() {
    state items = ["Apple", "Banana", "Cherry"]
    state selectedCount = 0

    on change(selectedCount) {
        Toast("You have selected ${selectedCount} items.")
    }

    column {
        align: (center, center)
        
        Text("My Fruit List") {
            weight: horizontal
        }

        for fruit in items {
            Card(fruit) {
                Button("Select ${fruit}") {
                    on click {
                        selectedCount = selectedCount + 1
                    }
                }
            }
        }
    }
}
```

## 10. Hints and Annotations

InterSpec provides a lightweight hint system for communicating **implementer guidance** directly inside `.is` files. Hints carry information about visual hierarchy, responsive intent, accessibility, spacing, animation, priority — anything that would help a developer or AI implement the interface faithfully.

### Philosophy

Hints occupy a tier between comments and properties:

| Construct | Syntax | Stripped by `isc strip`? | Transpiler sees? | Human/AI sees? |
|-----------|--------|------------------------|------------------|----------------|
| Comment | `//`, `/* */` | ✅ Removed | ❌ No | ❌ No |
| **Hint** | **`@ ...`**, **`@* ... *@`** | **❌ Kept** | **❌ Ignores** | **✅ Yes** |
| Property | `variant:`, `disabled:` | ❌ Kept | ✅ Enforces | ✅ Yes |

Hints are:
- **Freeform text** — no grammar, no parsing, no validation. Any text after `@` is valid.
- **Persistent** — they survive `isc strip` unchanged.
- **Non-enforced** — deterministic transpilers and runtimes ignore them entirely.
- **Human-first** — their audience is the person or AI that will translate the spec into a real UI.

### Syntax

#### Single-line hints

```interspec
@ This button is the primary call-to-action
Button("Save") {
    on click { submit() }
}
```

A single `@` begins a hint that runs to the end of the line. The space after `@` is conventional but not required.

#### Multi-line block hints

```interspec
@*
  On mobile, stack these fields vertically.
  On tablet and up, use a two-column grid.
  Keep labels above inputs for readability.
*@
Form {
    Input("First Name")
    Input("Last Name")
    Input("Email")
}
```

`@*` opens the block, `*@` closes it. Everything between is hint text. Blocks close at the first `*@` — nesting is not supported.

### Where hints can appear

Hints are valid anywhere in a `.is` file:

- **File-level** — top-of-file notes about the overall spec
- **Before a declaration** — single component or block
- **Inside a component block** — alongside properties and events
- **Between declarations** — separated by blank lines

```interspec
@* This page handles user onboarding.
   Keep the flow linear — one step at a time.
   Avoid modals on mobile. *@

page WelcomeWizard() {
    @ Compact layout — no extra spacing between sections
    column {
        Text("Step 1 of 3")

        Input("Full Name") {
            @ Consider adding autocomplete attributes for common name patterns
            required: true
            on commit { log("name: " + value) }
        }

        Button("Continue") {
            @ Primary action — make this visually prominent
            on click { navigate StepTwo() }
        }
    }
}
```

### Suggested uses

Hints have no fixed categories, but common applications include:

| Purpose | Example |
|---------|---------|
| **Visual hierarchy** | `@ Primary action — most prominent button on the page` |
| **Responsive intent** | `@ On mobile, replace this table with a card list` |
| **Accessibility** | `@ This image is decorative — use empty alt text` |
| **Spacing / density** | `@ Compact layout — minimize vertical gaps` |
| **Animation intent** | `@ Animate this card when it appears (fade + slide up)` |
| **Priority** | `@ High-priority: implement this first` |
| **Platform nuance** | `@ Desktop: show full table. Mobile: show first 3 columns` |
| **Destructive actions** | `@ Destructive — give this button prominent warning styling` |

### Notes on stripper interaction

The `isc strip` command strips `//` line comments and `/* */` block comments. Since hints use `@` and `@* *@` (not `/`), they are **not affected** by the stripper and pass through verbatim.

However, there is one interaction to be aware of:

- **Avoid `//` inside a single-line hint.** The character sequence `//` is always treated as a comment opener by `isc strip`, even if it appears after `@`. Everything from `//` to the end of the line will be removed.

  ```interspec
  @ Avoid this // everything after // is stripped
  @ Use a separate line instead — no slashes needed
  ```

- **`//` and `/*` inside `@* ... *@` blocks** are similarly treated as comment delimiters by the stripper. If your hint needs to mention URLs or code that contain slashes, rephrase to avoid literal `//` or `/*` inside the block.

  ```interspec
  @* Do this: reference the style guide at example.com
     Not this: https://example.com/docs (the // would cause issues) *@
  ```

### Relationship to comments

| Aspect | Comment (`//`, `/* */`) | Hint (`@`, `@* *@`) |
|--------|------------------------|---------------------|
| Audience | Developer of the `.is` file | Implementer of the final UI |
| Survival | Stripped by `isc strip` | Survives `isc strip` |
| Runtime | Removed before execution | Ignored (treated as code text) |
| Tone | Internal notes, TODOs, explanations | External guidance, design intent |
| Parsing | Recognized by the parser | Transparent to the parser |

Use comments for notes about the spec itself. Use hints for guidance about the final implementation.

---

## 11. Viewport Safety

InterSpec is style-free by design — it describes structure and interaction, not visual constraints. However, when a `.is` file is implemented into a real UI, the implementer must make decisions about viewport boundaries, overflow, and scrolling. This section provides patterns and conventions to prevent pages from overflowing the viewport or becoming unresponsive.

### The Problem

A `.is` file may declare a `column` with a `for` loop that renders 50 items, or a `Table` with many rows, or deeply nested layouts. Without explicit guidance, the implementer may render all of this content at full height — exceeding the viewport and creating an unscrollable, broken page.

### Recommended Page Pattern

Every page should establish a bounded root container. The recommended pattern:

```interspec
@* Viewport-safe page pattern:
   Root column fills the viewport height and scrolls internally.
   Max-width prevents edge-to-edge stretch on wide screens. *@
page Main() {
    column {
        scrollable: true
        align: (center, top)
        @ Constrain content width on desktop — center with max-width
        Text("My App")
        // ... content ...
    }
}
```

The implementer should translate this to:
- A root container with `height: 100vh` (or `dvh`) and `overflow: hidden`
- An inner scrollable area with `overflow-y: auto`
- A `max-width` constraint (e.g., `1200px`) with `margin: 0 auto` for desktop centering

### Table and List Overflow

Tables and long lists are the most common source of viewport overflow. Use `scrollable: true` on the parent layout:

```interspec
@ Data table — constrain to viewport height with internal scroll
column {
    scrollable: true
    Table(["Name", "Role"], rows)
}
```

For lists rendered with `for` loops, apply the same pattern:

```interspec
@ Long list — must scroll within bounded height
column {
    scrollable: true
    for item in items {
        Card(item)
    }
}
```

### Nesting Depth Warning

Deeply nested `row`/`column` layouts compound height. Each nesting level adds implicit vertical space. As a guideline:
- **3 or fewer nesting levels**: Safe for most implementations
- **4+ nesting levels**: Add a `@ constrained` hint to signal that the implementer should limit height at each level

```interspec
@ Constrained — limit height at each nesting level to prevent compounding
row {
    column {
        row {
            column {
                Text("Deeply nested content")
            }
        }
    }
}
```

### Viewport Hint Vocabulary

While hints are freeform, the following tokens are **conventional signals** that implementers (human or AI) should recognize and act on. Using these tokens makes viewport intent explicit and scannable.

| Token | Meaning | Implementer action |
|-------|---------|--------------------|
| `@ viewport-safe` | This section must not overflow the viewport | Constrain height, apply overflow handling |
| `@ scrollable` | This container should scroll when content exceeds space | Set bounded height, `overflow: auto` |
| `@ constrained` | Limit width/height to viewport or parent bounds | Apply max-width/max-height, prevent edge-to-edge stretch |
| `@ compact` | Minimize vertical space usage | Reduce padding, line-height, margins |
| `@ mobile-break` | Needs special handling on narrow screens | Implement responsive breakpoint behavior |

These tokens are **not** parsed or enforced — they are prose conventions. The consuming skill (see `interspec-consume`) teaches implementers to recognize and act on them.

```interspec
@ viewport-safe — constrain this section to viewport height
column {
    scrollable: true
    Table(columns, rows)
}

@ compact — minimize vertical space in this toolbar
row {
    wrap: true
    Button("Filter")
    Button("Sort")
    Button("Export")
}

@ mobile-break — stack into single column on narrow viewports
row {
    collapse: true
    Card("Left")
    Card("Right")
}
```
