---
name: interspec-write
description: >
  Design judgment for InterSpec (.is) files — patterns to follow, anti-patterns
  to avoid, and decision rules for writing structurally sound, implementable
  specifications. Co-load with interspec-reference for syntax lookup. Use when
  authoring, reviewing, or refactoring .is files.
license: MIT
metadata:
  author: interspec-community
  version: "1.0"
---

# InterSpec Writing Guide

You are an expert InterSpec author with strong design judgment. Use this skill
whenever you create, review, or refactor `.is` files and need to make decisions
— not just avoid syntax errors.

This skill assumes you already have `interspec-reference` loaded for syntax,
component parameters, and the event catalog. It does not re-teach grammar.

## What This Skill Covers

- **Design patterns** — reusable solutions for common UI specification problems
- **Anti-patterns** — common mistakes and why they cause problems
- **Decision rules** — heuristics for choosing between alternatives

## What It Does NOT Cover

- Syntax rules, component parameters, or grammar (→ `interspec-reference`)
- Frontend implementation (→ `interspec-consume`)
- Verification or reverse-engineering

## When to Load This Skill

Load when:
- Creating a new `.is` file from scratch
- Reviewing or refactoring an existing `.is` file
- The user asks about InterSpec best practices, patterns, or "how should I..."
- The user asks "is this good InterSpec?" or "what's wrong with this?"
- A `.is` file exceeds ~20 lines (design decisions become necessary)

Always co-load `interspec-reference` for syntax lookup.

## Core Principles

Before applying specific patterns, internalise these principles:

1. **Every page must be viewport-safe.** The root layout must have `scrollable: true`
   and a hint describing the viewport strategy. Without this, the implementer
   cannot produce a working page.

2. **Structure communicates intent.** Use `Section` to group related content.
   Use `Divider` to separate zones. Component placement tells the implementer
   what goes where.

3. **State drives visibility.** Modals, drawers, and conditional content are
   controlled by boolean state variables, placed at page level (outside any
   scrollable container).

4. **Hints are semantic, not visual.** A hint must answer "what is this?" or
   "what constraint does it obey?" — never "what color is it?" If a CSS
   property can be derived from the hint alone, the hint is wrong.

5. **Plan for emptiness.** Every dynamic list must have an `EmptyState` branch.
   Every async operation must consider loading and error states.

---

## Part 1: Design Patterns

### P1: Viewport-Safe Page Root

**Problem:** Pages with dynamic content overflow the viewport, creating
unscrollable, broken layouts on real devices.

**Solution:** Every `page` declaration must root in a `column` with
`scrollable: true` and a top-of-file `@* ... *@` hint describing the
viewport strategy.

**When to use:** Always. Every page, no exceptions.

```interspec
@* Viewport-safe page.
   Root column fills viewport height and scrolls internally.
   Content centered with max-width on desktop. *@
page Main() {
    column {
        scrollable: true
        align: (center, top)
        // Content goes here
    }
}
```

**Why it works:** The scrollable column establishes a bounded container.
The implementer translates this to `height: 100%` (or `100dvh`) with
`overflow-y: auto`. Content that exceeds the viewport scrolls instead of
overflowing.

**Consequences of ignoring:** Pages render at full content height, exceed
the viewport, and the user cannot scroll. The page is broken.

---

### P2: Multi-Section Page

**Problem:** Pages with distinct content zones (header, body, actions) lack
visible structure, making them hard to navigate and implement.

**Solution:** Group related content with `Section(title)`. Separate sections
with `Divider()`. Place action buttons at the bottom of the section they
belong to.

**When to use:** Pages with 3+ visually distinct zones. Also useful for
any page where you want to communicate structure to the implementer.

```interspec
page Main() {
    column {
        scrollable: true

        @ Hero zone — top of page
        Section("Dashboard") {
            Text("Welcome back, ${userName}")
        }

        Divider()

        Section("Recent Activity") {
            for activity in recentItems {
                Card(activity) { /* ... */ }
            }
        }

        Divider()

        @ Action zone — secondary weight
        Section("Quick Actions") {
            row {
                wrap: true
                Button("New Item") { /* ... */ }
                Button("Settings") { /* ... */ }
            }
        }
    }
}
```

**Why it works:** `Section` gives each zone a heading. `Divider` visually
separates them. The implementer can implement one section at a time, and
the page's information architecture is clear at a glance.

---

### P3: State-Driven Visibility

**Problem:** Conditional UI (modals, drawers, dialogs) mixed into the content
flow causes layout shifts and scroll interference when toggled.

**Solution:** Declare a boolean `state` for visibility. Place the overlay
component at **page level** — outside the scrollable column. Toggle with
`toggle(variable)`. Use `on close` and `on open` events for lifecycle hooks.

**When to use:** Any overlay component (Modal, Dialog, Drawer, Toast) that
must float above page content and not participate in document flow.

```interspec
page Main() {
    state showModal = false

    column {
        scrollable: true

        Button("Open settings") {
            on click { toggle(showModal) }
        }
        // Main content here
    }

    // Page-level — outside scrollable column
    if showModal {
        Modal("Settings") {
            on close { showModal = false }

            Text("Configure your preferences")

            Button("Close") {
                on click { showModal = false }
            }
        }
    }
}
```

**Why it works:** The modal is a sibling of the scrollable column, not a
child. The implementer renders it with `position: fixed` or as a portal,
ensuring it stays centered regardless of scroll position.

---

### P4: Form with Validation

**Problem:** Forms need clear feedback for required fields, validation errors,
submission-in-progress state, and success/failure outcomes.

**Solution:** Wrap inputs in a `Form` component. Set `required: true` on
mandatory fields. Use the `error` property for field-level validation messages.
Use `loading` on the submit button. Call `validate()` on submit. Provide
feedback via `Toast` or `Alert`.

**When to use:** Any form with 2+ inputs and a submit action.

```interspec
page Main() {
    state email = ""
    state agreed = false
    state submitting = false

    column {
        scrollable: true

        Form {
            on submit { validate() }

            Input("Email") {
                required: true
                placeholder: "you@example.com"
                error: email == "" ? "Email is required" : false
                on input { email = value }
            }

            Checkbox("I agree to the terms") {
                required: true
            }

            Button("Submit") {
                loading: submitting
                on click {
                    submitting = true
                    delay(1500, Toast("Form submitted!"))
                    delay(1500, submitting = false)
                }
            }
        }
    }
}
```

**Why it works:** Every form state is accounted for: empty field → `required`
indicator; invalid data → `error` message; in-progress → `loading` spinner on
button; success → `Toast` notification. The implementer knows exactly what
feedback to show at each stage.

---

### P5: List with Empty, Loading, and Populated States

**Problem:** Dynamic lists that can be empty render nothing, leaving users
confused about whether the UI is broken or just has no data.

**Solution:** Always wrap `for` loops in a conditional that checks for empty
data. Provide an `EmptyState` with a recovery action. Use `loading: true` on
the container during data operations.

**When to use:** Any list whose data comes from state that can be empty
(dynamic data, search results, user-generated content, filtered views).

```interspec
page Main() {
    state items = fetchItems()
    state loading = false

    column {
        scrollable: true
        loading: loading

        if items.length == 0 {
            EmptyState("No items yet") {
                Button("Create your first item") {
                    on click { navigate NewItem() }
                }
            }
        } else {
            for item in items {
                Card(item) {
                    Button("View") {
                        on click { navigate Detail(itemId: item.id) }
                    }
                }
            }
        }
    }
}
```

**Why it works:** The empty branch explains what's missing and provides a
path forward. The `loading` property on the column signals to the implementer
that this container needs a loading skeleton or spinner.

---

### P6: Component Extraction

**Problem:** Knowing when to extract a custom `component` vs keeping content
inline in the page.

**Solution:** Extract when the element is **(a) used 3+ times** across the
spec, **OR (b) semantically distinct enough to deserve a name** (e.g., a
"UserRow" that appears once but represents a clear concept). Use PascalCase
for the component name. Pass data via parameters. Append children via the
instantiation block (they go to the end).

**When to use:** Same structure appears 3+ times, OR a block of UI has a
clear semantic identity that makes the page easier to read when named.

```interspec
// BEFORE — inline duplication (bad at 3+ repetitions)
page Main() {
    column {
        scrollable: true
        Card("Alice") {
            Text("Admin")
            Button("Edit") { /* ... */ }
        }
        Card("Bob") {
            Text("User")
            Button("Edit") { /* ... */ }
        }
        Card("Carol") {
            Text("User")
            Button("Edit") { /* ... */ }
        }
    }
}

// AFTER — extracted component (good)
component UserRow(name, role) {
    Text(role)
    Button("Edit ${name}") {
        on click { navigate EditUser(userName: name) }
    }
}

page Main() {
    state users = getAllUsers()

    column {
        scrollable: true
        for user in users {
            Card(user.name) {
                UserRow(user.name, user.role)
            }
        }
    }
}
```

**Why it works:** The component name documents intent. Changes to the
repeated structure happen in one place. The page becomes a high-level
outline rather than a wall of repeated code.

**Heuristic:** If you catch yourself copy-pasting a block of `.is` code,
extract it.

---

### P7: Navigation Structure

**Problem:** Multi-page apps need clear navigation — how pages relate,
how data passes between them, how the user goes back.

**Solution:** Keep pages flat (all at the same level, no nested page
hierarchies). Pass essential identifiers via parameters (IDs, not entire
objects). Use `back()` for detail→list returns. Use `Breadcrumb` for
paths deeper than 2 levels.

**When to use:** Any spec with 2+ pages.

```interspec
// List page
page ProductList() {
    state products = getProducts()

    column {
        scrollable: true
        Breadcrumb(["Home", "Products"])

        if products.length == 0 {
            EmptyState("No products") {
                Button("Add product") {
                    on click { navigate NewProduct() }
                }
            }
        } else {
            for product in products {
                Card(product.name) {
                    Button("View details") {
                        on click { navigate ProductDetail(productId: product.id) }
                    }
                }
            }
        }
    }
}

// Detail page
page ProductDetail(productId) {
    state product = getProduct(productId)

    column {
        scrollable: true
        Breadcrumb(["Home", "Products", product.name])

        Text("Details for ${product.name}")

        Button("← Back to list") {
            on click { back() }
        }
    }
}
```

**Why it works:** Flat navigation is predictable. Parameters carry just
enough data for the target page to fetch what it needs. `back()` avoids
hardcoding navigation targets in detail pages.

---

### P8: Hint Convention System

**Problem:** Hints are freeform text — without conventions they become
inconsistent, vague, or drift into describing visual appearance.

**Solution:** Use a structured prefix system. Every hint answers exactly
one question about the element: what it **is**, what **constraint** it
obeys, how it **responds** to breakpoints, or what **accessibility** it
needs. Never describe colors, fonts, spacing, or shadows.

**When to use:** On components where the default interpretation would be
wrong or ambiguous. NOT on every component — hint only where needed.

| Prefix | Domain | Example |
|--------|--------|---------|
| `@ role:` | Semantic purpose | `@ role: primary-action` |
| `@ constraint:` | Behavioral boundary | `@ constraint: viewport-safe` |
| `@ responsive:` | Cross-breakpoint behavior | `@ responsive: collapse-mobile` |
| `@ a11y:` | Accessibility requirement | `@ a11y: decorative-img` |
| `@ depends-on:` | Logical dependency | `@ depends-on: form-validation` |

**Hint quality test:** Before writing a hint, ask: "Can a CSS property-value
pair be generated from this text alone, without consulting a design system?"
If yes — rewrite the hint. If no — the hint is valid.

```interspec
// WRONG — describes visual appearance
@ Blue background, white text, 16px padding, rounded corners
Button("Save") { /* ... */ }

// RIGHT — describes semantic role (visuals come from DESIGN.md)
@ role: primary-action
Button("Save") { /* ... */ }
```

---

### P9: Responsive Row Collapse

**Problem:** Row layouts that work on desktop break on narrow screens
(typically phones at ~375px width).

**Solution:** Use `collapse: true` on rows that should stack vertically
on narrow viewports. Add a `@ responsive:` hint to communicate the
specific breakpoint behavior.

**When to use:** Any `row` with 2+ children that would overflow a phone
screen width.

```interspec
@ responsive: collapse to single column on narrow viewports
row {
    collapse: true
    Card("Left panel") { /* ... */ }
    Card("Right panel") { /* ... */ }
}
```

**Why it works:** `collapse: true` is a structural signal to the implementer.
At the chosen breakpoint (typically 768px), the implementer changes
`flex-direction` from `row` to `column`. The hint provides the rationale.

---

### P10: File and Module Organization

**Problem:** A single monolithic `.is` file becomes unmaintainable as the
spec grows beyond a single page.

**Solution:** Split by concern. One file per page under `pages/`. Shared
components under `components/`. An `index.is` entry point that imports
pages. Use aliased imports for clarity when paths are deep.

**When to use:** When the spec exceeds ~200 lines total, or has 3+ pages.

```text
myapp/
├── index.is              ← Entry point, imports pages
├── pages/
│   ├── main.is           ← Dashboard / landing
│   ├── settings.is       ← Settings page
│   └── profile.is        ← User profile
└── components/
    ├── UserCard.is       ← Reusable across pages
    └── FormField.is      ← Shared form element
```

```interspec
// index.is
import "/pages/main.is"
import "/pages/settings.is"

// Entry point is Main() — declared in pages/main.is
```

**Why it works:** Each file is small enough to understand at a glance.
Pages are independent — changing one doesn't risk breaking another.
Shared components are explicitly visible in the `components/` directory.

---

## Part 2: Anti-Patterns

### A1: Modal in Scrollable Container

**Problem:** Placing a `Modal`, `Dialog`, or `Drawer` inside a `scrollable`
column causes it to scroll away with content or get clipped by
`overflow: hidden`.

**Why it's wrong:** Overlay components need `position: fixed` to float
above the page. Nesting them inside scrollable containers creates competing
positioning contexts — the modal is trapped inside the scroll area.

**Fix:** P3 (State-Driven Visibility). Declare the overlay at page level,
outside the scrollable column.

```interspec
// WRONG — modal scrolls away with content
page Main() {
    column {
        scrollable: true
        Button("Open") { /* ... */ }
        if showModal { Modal("Title") { /* ... */ } }
    }
}

// RIGHT — modal is a page-level sibling
page Main() {
    column {
        scrollable: true
        Button("Open") { /* ... */ }
    }
    if showModal { Modal("Title") { /* ... */ } }
}
```

---

### A2: Deep Nesting Without Scroll Bounds

**Problem:** Nesting `column` > `row` > `column` > `row` (4+ levels)
without `scrollable` at each level compounds vertical height unpredictably.
Each layout level adds implicit space in the implementation.

**Why it's wrong:** Four levels of nested layout produce a total height
that depends on the implementation's spacing defaults — not the spec's
intent. The result may exceed the viewport with no scroll container to
catch it.

**Fix:** Either add `scrollable: true` at every level beyond 3, or flatten
the structure. Flatter is better.

```interspec
// WRONG — 4 levels, no scroll handling
column {
    row {
        column {
            row {
                Text("Deeply nested")
            }
        }
    }
}

// RIGHT — bounded or flattened
column {
    scrollable: true
    row {
        column {
            Text("Item 1")
            Text("Item 2")
        }
    }
}
```

---

### A3: God Page

**Problem:** A single `page` declaration exceeds ~200 lines, containing
multiple unrelated UI zones with no clear structure.

**Why it's wrong:** Large pages are hard to review, hard to implement
incrementally, and mask opportunities for reuse. The implementer cannot
tell which zones are independent and which are coupled.

**Fix:** Apply P2 (Multi-Section Page) for structure within the page.
Extract reusable elements with P6 (Component Extraction). If zones
represent distinct navigation targets, split into multiple pages.

---

### A4: Missing Empty State

**Problem:** A `for` loop over a dynamic array has no `if`/`else` branch
for the empty case. When data is empty, the UI renders a blank screen.

**Why it's wrong:** Blank screens look broken. Users cannot distinguish
"no data yet" from "the app crashed." They have no path forward.

**Fix:** P5 (List with Empty State). Always check `items.length == 0`
before the `for` loop, with an `EmptyState` in the true branch.

```interspec
// WRONG — blank screen when items is empty
for item in items {
    Card(item) { /* ... */ }
}

// RIGHT — empty state with recovery action
if items.length == 0 {
    EmptyState("No items found") {
        Button("Create one") { /* ... */ }
    }
} else {
    for item in items {
        Card(item) { /* ... */ }
    }
}
```

---

### A5: Premature Component Extraction

**Problem:** Creating a custom `component` for a UI element used only once,
adding a layer of indirection without benefit.

**Why it's wrong:** Each custom component is a cognitive hop — the reader
must find the declaration to understand what it renders. A one-off component
adds complexity without reuse.

**Fix:** P6 (Component Extraction). Keep it inline unless used 3+ times
or semantically distinct enough to justify the name.

```interspec
// WRONG — used once, adds indirection
component FancyHeader(title) {
    Text(title)
    Divider()
}
page Main() {
    FancyHeader("Dashboard")  // Only used here
}

// RIGHT — inline when used once
page Main() {
    Text("Dashboard")
    Divider()
}
```

---

### A6: Over-Parameterization

**Problem:** A custom component accepts 5+ parameters, making instantiation
verbose and the component's purpose unclear.

**Why it's wrong:** Too many parameters suggest the component does too much.
Every parameter adds cognitive load at the call site. The implementer must
trace each parameter to understand the component's behavior.

**Fix:** Split into smaller, single-purpose components. Use children and
properties instead of parameters for optional configuration.

```interspec
// WRONG — 6 params, unclear purpose
component DataCard(title, subtitle, value, unit, icon, variant) { /* ... */ }

// RIGHT — focused components
component MetricCard(title, value) { /* ... */ }
component ActionCard(title, action) { /* ... */ }
```

---

### A7: Styling via Hints

**Problem:** Using `@` hints to describe visual appearance — colors, fonts,
spacing values, pixel dimensions, shadows, or animation curves.

**Why it's wrong:** Hints are for semantic role and behavioral constraint,
not visual styling. Visual styling belongs in DESIGN.md. If an implementer
can write a CSS property from the hint alone, the hint is wrong.

```interspec
// WRONG — describes colors and spacing
@ Blue (#2563eb) background, white text, 16px padding all sides
Button("Save") { /* ... */ }

// RIGHT — describes semantic role
@ role: primary-action
Button("Save") { /* ... */ }
```

**Fix:** P8 (Hint Convention System). Use `@ role:`, `@ constraint:`,
`@ responsive:`, or `@ a11y:` prefixes. Let the design system provide
the visual tokens.

---

### A8: Vague Hints

**Problem:** Hints that say nothing actionable — `@ make it nice`,
`@ important`, `@ fancy`, `@ make it pop`.

**Why it's wrong:** The implementer gains zero information. The hint wastes
space and trains implementers to ignore all hints.

**Fix:** Either make the hint specific and actionable, or remove it.
Every hint must answer: "What should the implementer do differently?"

---

### A9: Hint Spam

**Problem:** Every single element in the spec carries a `@` hint, even
for obvious defaults.

**Why it's wrong:** Hint fatigue. When everything is annotated, nothing
stands out. The implementer stops reading hints — missing the ones that
actually matter.

**Fix:** Hint only where the default interpretation would be wrong or
ambiguous. `Text("Hello")` doesn't need `@ role: heading`. A
`Button("Save")` that is the primary CTA does need `@ role: primary-action`.

**Heuristic:** If removing the hint wouldn't change the implementer's
output, remove it.

---

### A10: Unbounded Iteration

**Problem:** A `for` loop iterates over a dataset with no upper bound —
e.g., "all users" or "all products" without pagination.

**Why it's wrong:** Even with the runtime's infinite-loop guard, a loop
over 10,000 items renders 10,000 components — hanging the runtime or
producing an un-scrollable, unusable page.

**Fix:** Paginate large lists. Limit displayed items to a reasonable count
(50–100 max). Use the `Pagination` component.

```interspec
// WRONG — no upper bound, could be 10000 items
for user in allUsers {
    Card(user.name) { /* ... */ }
}

// RIGHT — paginated subset
for user in visibleUsers {
    Card(user.name) { /* ... */ }
}
Pagination {
    current: currentPage
    total: totalPages
}
```

---

### A11: Card Grid Default

**Problem:** Every list of items is wrapped in a uniform grid of identical
`Card` components, producing layouts that look indistinguishable from every
other AI-generated UI. The spec defaults to `row { wrap: true }` full of
`Card(item)` with no variation in structure, weight, or presentation.

**Why it's wrong:** Uniform card grids are the most overused pattern in
AI-generated interfaces. They signal "no design decisions were made here."
Every item gets equal visual weight regardless of its importance. The
implementer has no signal about which items deserve prominence.

```interspec
// WRONG — uniform grid, every item is the same shape and weight
row {
    wrap: true
    for item in items {
        Card(item) {
            Text(item.description)
            Button("View") { /* ... */ }
        }
    }
}
```

**Fix:** Vary the structure. Not every item needs a `Card`. Mix layouts:
some items as rich rows with metadata, others as compact list entries, key
items as featured cards. Use `@ role:` hints to signal which items are
primary, secondary, or summary-level. Break the grid with a `Section`
separator or a featured item that spans the full width.

```interspec
// RIGHT — varied structure with clear hierarchy
column {
    scrollable: true

    @ Featured item — spans full width, more detail
    if featuredItem {
        Card(featuredItem.name) {
            @ role: hero-card
            Image(featuredItem.image)
            Text(featuredItem.description)
            row {
                Button("Primary action") { /* ... */ }
                Button("Secondary") { /* ... */ }
            }
        }
        Divider()
    }

    @ Remaining items — compact list, not cards
    Section("All items") {
        for item in remainingItems {
            row {
                Text(item.name) { weight: horizontal }
                @ Summary-level — minimal detail
                Badge(item.status) { variant: info }
                Button("View") { /* ... */ }
            }
        }
    }
}
```

**Heuristic:** If your spec has 3+ `Card` components in a single `row` or
`for` loop with identical structure, stop and ask: "Do all these items
deserve equal visual weight? Could some be richer and others more compact?"

---

### A12: Hero Section Cliché

**Problem:** Pages open with a full-width `Image` + centered `Text` title
— the startup landing page pattern that was distinctive in 2015 and is now
the single most generic way to open any page. In InterSpec, this manifests
as `Image(src)` followed by a centered `Text()` at the top of every page.

**Why it's wrong:** The hero-image-with-centered-title pattern is the
visual equivalent of "insert generic UI here." It wastes the most valuable
screen real estate (above the fold) on a decorative image that communicates
nothing structural. The implementer has no signal about what the page
does or how to navigate it.

```interspec
// WRONG — generic hero, says nothing about what the page does
page Main() {
    column {
        scrollable: true

        Image("hero-banner.png")
        Text("Welcome to MyApp") {
            align: (center, center)
        }
        Text("The best app for everything") {
            align: (center, center)
        }

        // ... actual content buried below ...
    }
}
```

**Fix:** Lead with structure that communicates purpose. Open with the
page's primary action, a data display, a navigation element, or a
search field — something the user can **use**, not just look at. Put
the title in a `Section` heading where it belongs. If an image is
essential, place it alongside content (in a `row`), not as a full-width
backdrop.

```interspec
// RIGHT — leads with actionable structure
page Main() {
    column {
        scrollable: true

        @ Page header — compact, actionable
        Section("Dashboard") {
            row {
                Text("Welcome back, ${userName}") {
                    weight: horizontal
                }
                Button("New project") {
                    @ role: primary-action
                    on click { navigate NewProject() }
                }
            }
        }

        Divider()

        @ Content starts here — immediate value
        Section("Recent projects") {
            if projects.length == 0 {
                EmptyState("No projects yet") {
                    Button("Create your first") {
                        on click { navigate NewProject() }
                    }
                }
            } else {
                for project in projects {
                    // Varied structure, not uniform cards
                    row {
                        Text(project.name) { weight: horizontal }
                        Text(project.updatedAt)
                    }
                }
            }
        }
    }
}
```

**Heuristic:** If your page's first visible elements are `Image` + centered
`Text`, delete them and start with something the user can act on.

---

## Part 3: Decision Rules

These are compact heuristics for common design decisions. When you reach a
fork, apply the rule.

| # | Rule | When | Action |
|---|------|------|--------|
| R1 | **Viewport root** | Any `page` declaration | Always wrap in `column { scrollable: true }` with viewport hint |
| R2 | **Extract component** | Same structure appears 3+ times OR has clear semantic identity | Create `component` with PascalCase name; pass data via parameters |
| R3 | **Use Section** | 3+ visually distinct zones on one page | Wrap each zone in `Section(title)`, separate with `Divider` |
| R4 | **Empty state** | Any `for` loop over dynamic/array data | Add `if empty` / `else` branch with `EmptyState` |
| R5 | **Modal placement** | Any Modal, Dialog, or Drawer | Place at page level, outside the scrollable column |
| R6 | **Hint prefix** | Any hint about semantic role or behavior | Use `@ role:`, `@ constraint:`, `@ responsive:`, or `@ a11y:` |
| R7 | **Hint quality** | Before writing any hint | Ask: "Can a CSS property be generated from this alone?" If yes → rewrite |
| R8 | **Nesting limit** | 4+ nested layout levels | Add `scrollable: true` at each level, or flatten the structure |
| R9 | **Page split** | Page exceeds ~200 lines OR 3+ distinct nav targets | Split into multiple `page` declarations with separate files |
| R10 | **Pagination** | List may contain >100 items | Add `Pagination` component; loop over current page subset only |
| R11 | **Loading state** | Any container with async data operations | Add `loading: true` property to the container |
| R12 | **accessibility hint** | Image with no informational content | Add `@ a11y: decorative-img` |
| R13 | **Card variety** | 3+ identical-structure Cards in one row or loop | Vary structure: some rich (image + text + actions), some compact (text + badge). Use `@ role:` to signal prominence. |
| R14 | **Page opening** | Tempted to start page with Image + centered Text | Delete both. Lead with action, data, navigation, or search — something the user can use immediately. |

---

## Part 4: Authoring Workflow

When creating a `.is` file from scratch, follow this order:

1. **Define pages.** List every `page` the spec needs. Start with `Main()`.
2. **Sketch state.** For each page, list the state variables. Keep them minimal.
3. **Structure each page.** Apply P1 (viewport root). Add sections and dividers.
4. **Add content.** Use built-in components. Apply P6 for repeated structures.
5. **Wire events.** Connect buttons to actions, inputs to state changes.
6. **Add overlays.** Place modals/drawers at page level (P3).
7. **Handle emptiness.** For every `for` loop, add an `EmptyState` branch (P5).
8. **Add hints.** Annotate only where the default interpretation is ambiguous (P8).
9. **Review against anti-patterns.** Check A1–A12. Pay special attention to A11 (card grid uniformity) and A12 (hero section cliché) — these are the most common AI defaultism patterns.
10. **Split files.** If the result exceeds ~200 lines, apply P10.

---

## Cross-References

| Skill | Relationship |
|-------|-------------|
| `interspec-reference` | Co-loaded. Provides syntax, component parameters, and the event catalog. |
| `interspec-consume` | `write` produces specs that `consume` implements. Specs following `write` patterns are easier to implement correctly. |
| `interspec-verify` | `write` patterns become `verify` conformance checks. A spec following this guide should pass verification. |
| `interspec-reverse` | `reverse` output should follow `write` patterns, not just `reference` syntax. |

**Key reference files:**
- `LANGUAGE.md` — Sections 8 (Built-in Catalog), 10 (Hints), 11 (Viewport Safety)
- `skills/interspec-reference/references/CATALOG.md` — full component/event/property catalog
- `skills/interspec-reference/references/EXAMPLES.md` — example snippets
- `example.is` — comprehensive feature demo (462 lines)
