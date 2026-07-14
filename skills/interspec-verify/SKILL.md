---
name: interspec-verify
description: >
  Verify that a frontend implementation conforms to its InterSpec (.is)
  specification. Checks structural conformance, viewport safety, hint
  compliance, and event wiring. Use after interspec-consume to validate
  that the build matches the spec.
license: MIT
metadata:
  author: interspec-community
  version: "0.1"
---

# InterSpec Verification

You verify that a frontend implementation (HTML page, React component, Vue app,
or any DOM-renderable UI) conforms to its InterSpec specification. You answer
the question: **"Does the build match the spec?"**

This is a skeleton skill — it defines the scope, contract, and verification
strategy. Detailed procedures will be filled in as the verification tooling
matures.

## What This Skill Covers

- **Structural conformance:** Every component declared in the `.is` file has
  a corresponding element in the implementation DOM.
- **Viewport safety:** The implementation respects scroll containers, max-width
  constraints, and overflow handling declared in the spec.
- **Hint conformance:** `@ role:` hints resolve to correct design tokens.
  `@ constraint:` and `@ responsive:` hints are structurally enforced.
- **Accessibility:** `@ a11y:` hints map to correct ARIA attributes, alt text,
  and focus management.
- **Component hierarchy:** Nesting order in the implementation matches the
  spec's parent-child structure.
- **Event wiring:** Declared events have corresponding handlers in the
  implementation source.

## What It Does NOT Cover

- Visual design review (colors, typography, spacing → DESIGN.md audit)
- Performance testing
- Cross-browser testing
- Security audit
- Business logic correctness

## When to Load This Skill

Load when:
- User asks to "verify", "check", "validate", "audit", or "test" a UI against
  a `.is` spec
- User asks "does this implementation match the spec?"
- As a CI pipeline step after `interspec-consume` generates code
- User wants a conformance report for a deployed page

## Input/Output Contract

**Input:**
- A `.is` file path (the specification)
- One or more of: a live URL, an HTML file, or component source code

**Output:**
- A conformance report organized by severity:
  - ✅ **Pass** — element matches spec
  - ⚠️ **Warning** — deviation from best practice, does not break functionality
  - 🚫 **Failure** — spec requirement not met, page may be broken
- Each finding includes: the spec component, the implementation element,
  file paths, and line numbers where available.

**Side effects:**
- May use Playwright or puppeteer for DOM inspection of live URLs.
- Reads `DESIGN.md` if present for hint-to-design-token mapping verification.
- Reads `interspec-write` for structural pattern checks.

## Relationship to Other Skills

| Skill | Relationship |
|-------|-------------|
| `interspec-reference` | The conformance baseline — you check that the implementation matches what `reference` declares as valid. |
| `interspec-write` | You check that `write` patterns (viewport root, modal placement, empty states) are reflected in the implementation. If the spec follows `write` patterns, verification is mechanical; if it doesn't, you flag the spec issues first. |
| `interspec-consume` | You are the quality gate after `consume`. The full pipeline: `write` → `consume` → `verify`. |
| `interspec-reverse` | A round-trip verification: `reverse(implementation)` should produce a spec structurally equivalent to the original. |

## Key Reference Files

- `LANGUAGE.md` — Sections 8 (Built-in Catalog), 10 (Hints), 11 (Viewport Safety)
- `skills/interspec-consume/SKILL.md` — implementation rules (check against these)
- `skills/interspec-write/SKILL.md` — patterns and anti-patterns (structural checks)
- `skills/interspec-reference/references/CATALOG.md` — component catalog for type mapping

## Verification Checklist

When performing verification, check each item below. Report severity and
suggested fix for every failure.

| # | Check | Method | Severity if Missing |
|---|-------|--------|---------------------|
| V1 | Every component type in `.is` appears in DOM | AST-to-DOM element mapping | 🚫 |
| V2 | Viewport root has bounded height | CSS inspection (`height`, `max-height`) | 🚫 |
| V3 | Scrollable containers have `overflow: auto` or `overflow-y: auto` | CSS inspection | 🚫 |
| V4 | Content area has `max-width` on wide viewports (≥1200px) | CSS inspection at desktop viewport | ⚠️ |
| V5 | Modals/drawers use `position: fixed` or portal rendering | CSS/DOM inspection | 🚫 |
| V6 | `@ role:` hints resolve to design tokens in output | DESIGN.md lookup + CSS class check | ⚠️ |
| V7 | `@ a11y:` hints have corresponding ARIA attributes | DOM attribute inspection | ⚠️ |
| V8 | `for` loops produce correct number of elements | DOM child count vs array length | 🚫 |
| V9 | `EmptyState` appears when data is empty | Conditional DOM inspection (set data to empty) | ⚠️ |
| V10 | `collapse: true` rows stack vertically on narrow viewports | Responsive viewport test (≤768px) | ⚠️ |
| V11 | Events map to handler functions in source | Source code inspection | ⚠️ |
| V12 | `loading: true` containers show loading indicators during async | State inspection during simulated delay | ⚠️ |

## Verification Workflow (Draft)

```
Input (.is spec + implementation URL/source)
    ↓
1. Parse .is spec → component tree, state declarations, event bindings
2. Inspect implementation → DOM tree, CSS computed styles, event handlers
3. Map spec components to DOM elements (by type, content, position)
4. Check each verification item (V1–V12)
5. Report findings with severity, file paths, and fix suggestions
```

## Notes

This skill is a **skeleton**. The verification workflow requires tooling that
does not yet exist — specifically, a `.is` parser (to produce an AST),
a DOM inspector (Playwright integration), and a mapping engine. As these
tools become available, this skill will be updated with concrete commands
and procedures.

For now, verification can be performed **manually** by an agent: read the
`.is` file, inspect the rendered page (via Playwright or screenshot), and
compare structure element by element against the checklist above.
