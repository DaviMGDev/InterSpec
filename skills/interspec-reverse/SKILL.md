---
name: interspec-reverse
description: >
  Extract an InterSpec (.is) specification from an existing UI — a live page,
  screenshot, or HTML source. Produces structurally clean .is files that follow
  interspec-write patterns. The inverse of interspec-consume. Use when migrating
  existing UIs to InterSpec-driven workflows or documenting legacy interfaces.
license: MIT
metadata:
  author: interspec-community
  version: "0.1"
---

# InterSpec Reverse Engineering

You extract InterSpec specifications from existing user interfaces. Given a
live URL, a screenshot, an HTML file, or component source code, you produce
a `.is` file that structurally describes the UI. You answer the question:
**"What would the InterSpec look like for this UI?"**

This is a skeleton skill — it defines the scope, contract, and extraction
pipeline. Detailed procedures will be filled in as reverse-engineering
tooling matures.

## What This Skill Covers

- **Structural extraction:** Identify components, layouts, and hierarchy
  from DOM inspection or visual analysis of screenshots.
- **Component classification:** Map HTML elements and ARIA roles to InterSpec
  built-in components.
- **Layout inference:** Detect `row` vs `column` orientation from flexbox/grid
  CSS or visual arrangement in screenshots.
- **State inference:** Identify interactive elements that suggest state
  variables (toggles, inputs with values, conditional visibility, list data).
- **Hint generation:** Produce `@ role:`, `@ responsive:`, and `@ a11y:`
  hints from semantic HTML, ARIA roles, and responsive breakpoint behavior.
- **Pattern application:** Restructure extracted content to follow
  `interspec-write` patterns — add viewport roots, lift modals to page level,
  add empty states to lists, split large pages.

## What It Does NOT Cover

- Extracting visual design tokens (colors, fonts, spacing → DESIGN.md extraction)
- Extracting business logic or backend API contracts
- Producing pixel-perfect visual clones (InterSpec is style-free)
- Extracting animations or transition behaviors

## When to Load This Skill

Load when:
- User asks to "extract", "reverse-engineer", "generate a spec for", or
  "document" an existing UI
- User provides a URL, screenshot, or HTML and asks "what's the InterSpec
  for this?"
- User wants to migrate an existing app to an InterSpec-driven workflow
- User wants to document a legacy interface for reference or redesign

## Input/Output Contract

**Input:** One or more of:
- A live URL (inspected via Playwright or web_fetch)
- An HTML file or snippet
- A screenshot or image of the UI (for visual analysis)
- Component source code (React, Vue, Svelte, etc.)

**Output:**
- A `.is` file (or set of files organized by page/component) that
  structurally describes the UI.
- The output **must** follow `interspec-write` patterns — it should not
  blindly mirror the input's structure, which may contain anti-patterns.
- Each generated `.is` file includes appropriate `@` hints for semantic
  roles, responsive behavior, and accessibility.

**Side effects:**
- Uses Playwright for DOM inspection of live URLs.
- Uses screenshot/image analysis for visual structure detection when DOM
  access is not available.
- Reads `interspec-reference` for the component catalog.
- Reads `interspec-write` for pattern guidance during cleanup.

## Relationship to Other Skills

| Skill | Relationship |
|-------|-------------|
| `interspec-reference` | The target language. Your output must use only valid InterSpec constructs from the catalog. |
| `interspec-write` | You are an author, so you apply `write` patterns. Extracted structure is cleaned up before output. |
| `interspec-consume` | You are the inverse of `consume`. A round-trip `consume(reverse(ui))` should produce a UI structurally equivalent to the original. |
| `interspec-verify` | After extraction, `verify` can validate that the generated spec accurately describes the original UI. |

## Key Reference Files

- `LANGUAGE.md` — Section 8 (Built-in Catalog) for component mapping
- `skills/interspec-reference/references/CATALOG.md` — full component/event/property catalog
- `skills/interspec-write/SKILL.md` — patterns to apply during cleanup
- `skills/interspec-write/SKILL.md` — anti-patterns to detect and fix in extracted structure

## Extraction Pipeline

```
Input (URL / HTML / screenshot / component source)
    ↓
Stage 1: DOM Inspection or Visual Analysis
    - Identify layout regions (rows, columns, sections)
    - Classify elements into InterSpec built-in components
    - Detect hierarchy and nesting depth
    - Note responsive breakpoints (if inspecting at multiple widths)
    ↓
Stage 2: State Inference
    - Identify interactive elements → candidate state variables
    - Detect conditional visibility → boolean state flags
    - Find repeated elements → array state (list data)
    - Note form inputs → form state variables
    ↓
Stage 3: Structure Generation
    - Build page tree using row/column layouts
    - Map HTML elements to InterSpec component calls
    - Generate preliminary .is code (syntactically valid, but may contain
      structural anti-patterns from the original UI)
    ↓
Stage 4: Pattern Application (interspec-write cleanup)
    - Add viewport-safe root wrapper (P1)
    - Lift modals/drawers to page level (P3)
    - Add empty states for lists (P5)
    - Add Section/Divider grouping (P2)
    - Add @ role: hints from semantic HTML/ARIA roles (P8)
    - Add @ responsive: hints from observed breakpoint behavior (P9)
    - Split pages exceeding ~200 lines (R9)
    - Remove styling-oriented content (if any leaked into structure)
    ↓
Stage 5: Output
    - Organized .is file(s) following P10 (file and module organization)
    - Entry point index.is that imports all pages
    - Shared components extracted into components/
```

## Component Mapping Table

When classifying HTML elements into InterSpec components, use this mapping
as your primary reference. When multiple mappings are possible, prefer the
one that best matches the element's semantic role.

| HTML Element / ARIA Role | InterSpec Component |
|--------------------------|---------------------|
| `<button>` (generic action) | `Button(label)` |
| `<button type="submit">` | `Button(label)` inside `Form` |
| `<input type="text">`, `<input type="email">` | `Input(placeholder)` |
| `<input type="checkbox">` | `Checkbox(label)` |
| `<input type="date">` | `DatePicker(placeholder)` |
| `<select>` (form input context) | `Select(options)` |
| `<select>` / `<menu>` (action/command context) | `DropdownMenu(label)` |
| `<table>` | `Table(columns, rows)` |
| `<nav>` / `<ol class="breadcrumb">` | `Breadcrumb(items)` |
| `<dialog>`, `role="dialog"` | `Dialog(title)` or `Modal(title)` |
| `<aside>`, drawer/sidebar pattern | `Drawer(title)` |
| `<img>` (content image) | `Image(src)` |
| `<img>` (decorative, `alt=""`) | `Image(src)` with `@ a11y: decorative-img` |
| `<progress>` | `Progress(value)` |
| `<hr>` | `Divider()` |
| `<section>`, `<fieldset>` | `Section(title)` |
| `<form>` | `Form` |
| `<input type="file">` | `FileUpload(label)` |
| Tab panel pattern (`role="tablist"`) | `Tabs(tabs)` |
| Accordion pattern (`<details>`) | `Accordion(items)` |
| Toast/snackbar (`role="status"`) | `Toast(message)` |
| Tooltip (`title` attribute, `role="tooltip"`) | `Tooltip(content)` |
| Alert/banner (`role="alert"`) | `Alert(message)` |
| Badge/chip (`<span class="badge">`) | `Badge(label)` |
| Pagination nav (`<nav aria-label="pagination">`) | `Pagination` |
| Stepper/wizard indicator | `Stepper(steps)` |
| Tree view (nested `<ul>`) | `TreeView(items)` |
| Toggle switch (`role="switch"`) | `Toggle(label)` |
| Range slider (`<input type="range">`) | `Slider` |
| Icon (`<svg>`, icon font, `<i class="icon-*">`) | `Icon(name)` |
| Link (`<a>` with navigation behavior) | `Link(label)` |
| Card/panel wrapper (`<div class="card">`) | `Card(title)` |
| Empty state / zero state | `EmptyState(message)` |

## Layout Detection Heuristics

| Visual Arrangement / CSS | InterSpec Layout |
|--------------------------|------------------|
| `display: flex; flex-direction: row` | `row { ... }` |
| `display: flex; flex-direction: column` | `column { ... }` |
| `flex-wrap: wrap` | `row { wrap: true }` |
| `@media (max-width: …) { flex-direction: column }` | `row { collapse: true }` with `@ responsive:` hint |
| `overflow-y: auto` on container | `column { scrollable: true }` |
| Centered content with `max-width` | `align: (center, top)` with `@ constrained` hint |

## Notes

This skill is a **skeleton**. The extraction pipeline requires tooling that
does not yet exist — specifically, DOM inspection (Playwright integration),
visual screenshot analysis (for non-DOM sources), and a `.is` code generator.
As these tools become available, this skill will be updated with concrete
commands and procedures.

For now, reverse-engineering can be performed **manually** by an agent:
inspect the source (HTML or screenshot), classify elements against the
mapping table, and write the `.is` file following the pipeline stages above.
