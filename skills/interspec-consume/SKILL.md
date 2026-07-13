---
name: interspec-consume
description: >
  Implement InterSpec (.is) files into production frontend code with viewport
  safety, responsive behavior, and overflow prevention. Use when the user asks
  to build, implement, or render a .is file into HTML/CSS/JS, React, Vue, or
  any frontend framework.
license: MIT
metadata:
  author: interspec-community
  version: "1.0"
---

# InterSpec Consuming Skill

You are an expert at translating InterSpec (.is) files into working frontend
code. This skill teaches you how to implement .is files **safely** — ensuring
pages never overflow the viewport, never become unresponsive, and always handle
content growth gracefully.

## The Viewport Contract

When implementing a `.is` file, you MUST follow this contract:

1. **Every page MUST have a bounded root container.** The root layout must have
   a fixed height (viewport or parent) and `overflow: hidden` (or equivalent).
2. **Long content MUST scroll.** Any container with `scrollable: true`, `for`
   loops, many children (>5), tables, or lists MUST be constrained and scrollable.
3. **Desktop content MUST have max-width.** On wide screens, content should be
   centered with a max-width constraint (e.g., 1200px) — not stretch edge-to-edge.
4. **Nesting MUST be height-aware.** Deeply nested `row`/`column` layouts must
   not compound height unboundedly. Each level should respect its parent's bounds.

**Violating this contract produces broken pages.** Content that exceeds the
viewport without scrolling is unusable.

## Property Translation Guide

### `scrollable: true`

The most important property for viewport safety. When you see `scrollable: true`
on a `column` or `row`:

**HTML/CSS:**
```css
.scrollable-container {
  height: 100%;           /* or a fixed value */
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;  /* smooth scroll on iOS */
}
```

**React:**
```jsx
<div style={{ height: '100%', overflowY: 'auto' }}>
  {/* children */}
</div>
```

**Vue:**
```vue
<div class="scrollable-container">
  <!-- children -->
</div>
```

### `wrap: true` (on `row`)

Items wrap to the next line when they exceed the container width.

**CSS:** `flex-wrap: wrap` or `flex-wrap: wrap-reverse`

### `collapse: true` (on `row`)

The row collapses to a column on narrow viewports. The runtime chooses the
breakpoint — typically 768px.

**CSS:** Use a media query or container query:
```css
@media (max-width: 768px) {
  .collapsible-row {
    flex-direction: column;
  }
}
```

### `align: (vertical, horizontal)`

Maps to flexbox alignment:
- `vertical`: `align-items` (top → flex-start, center → center, bottom → flex-end)
- `horizontal`: `justify-content` (left → flex-start, center → center, right → flex-end)

### `weight: both`

The component fills available space in both directions. Maps to `flex: 1` or
`flex-grow: 1`.

## Default Safety Rules

Apply these rules to EVERY page you implement from a `.is` file:

### Rule 1: Root viewport container

Every `page Main()` (or equivalent entry point) MUST render inside a container
that fills the viewport:

```css
/* Applied to the root of every page */
html, body, #root {
  margin: 0;
  padding: 0;
  height: 100vh;     /* or 100dvh for dynamic viewport */
  overflow: hidden;
}
```

The innermost layout of the page gets `overflow-y: auto` to enable scrolling.

### Rule 2: Max-width on content

Desktop layouts should constrain content width:

```css
.page-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}
```

If the `.is` file has `@ constrained` hints, honor them. If not, apply this
as a sensible default.

### Rule 3: Table overflow

Tables with visible rows MUST be in a scrollable container:

```css
.table-wrapper {
  overflow-x: auto;    /* horizontal scroll for wide tables */
  max-height: 60vh;    /* vertical constraint */
  overflow-y: auto;
}
```

### Rule 4: List virtualization threshold

If a `for` loop iterates over an array with >20 items, consider:
- Using a virtualized list (react-window, vue-virtual-scroller, etc.)
- Or at minimum, ensuring the list is in a scrollable container with bounded height

If the `.is` file has a `@ viewport-safe` hint on the list, prioritize
virtualization.

### Rule 5: Modal and drawer containment

Modals and drawers MUST NOT push page content down. They must:
- Use `position: fixed` or `position: absolute` with a backdrop
- Have their own scroll containment
- Never exceed viewport height

## Viewport Hint Recognition

When you encounter these hint tokens in `.is` files, take the specified action:

| Token | Action |
|-------|--------|
| `@ viewport-safe` | Constrain this section to viewport height, apply overflow handling |
| `@ scrollable` | Set bounded height on container, apply `overflow: auto` |
| `@ constrained` | Apply `max-width` and/or `max-height`, prevent edge-to-edge stretch |
| `@ compact` | Minimize padding, line-height, and margins in this section |
| `@ mobile-break` | Implement responsive breakpoint behavior (stack on narrow screens) |

These are **not** optional suggestions — they are signals that the spec author
identified a viewport risk. Always act on them.

## Common Pitfalls

### Pitfall 1: Nested scroll containers

**Problem:** A scrollable column inside another scrollable column creates
competing scroll contexts. The inner one captures scroll events, making the
outer one unreachable.

**Fix:** Only one scroll container per scroll axis in a given nesting path.
If you need multiple scrollable sections, give each a fixed height and place
them side by side (row) rather than nested (column).

### Pitfall 2: Dynamic content height

**Problem:** A container is initially short, but after a `for` loop renders
or a state change reveals content, it exceeds the viewport.

**Fix:** Always set `height: 100%` (or equivalent) on scroll containers
**before** content loads. Don't rely on content to set the height.

### Pitfall 3: Flex children ignoring parent bounds

**Problem:** A flex child with `flex: 1` or `height: 100%` ignores the
parent's overflow constraint and grows unbounded.

**Fix:** Ensure every flex parent in the chain has `overflow: hidden` or
`overflow: auto` with a defined height.

### Pitfall 4: 100vh on mobile

**Problem:** `100vh` on mobile includes the browser chrome (address bar),
causing content to be slightly taller than the visible area.

**Fix:** Use `100dvh` (dynamic viewport height) when available, with
`100vh` as a fallback:
```css
height: 100dvh;
height: 100vh; /* fallback */
```

### Pitfall 5: Missing min-height on empty states

**Problem:** An empty page or loading state has no content, so the container
collapses to zero height.

**Fix:** Set `min-height: 100vh` (or `100dvh`) on the root container so
the page always fills the viewport.

## Framework-Specific Notes

### React / Next.js

- Use `100dvh` on `<html>` or the root layout div
- For App Router, set viewport height on the `<body>` or root layout
- Use `sticky` positioning for headers/footers that should stay in view

### Vue / Nuxt

- Set viewport height on `<div id="app">` or the root layout component
- Use `Teleport` for modals/drawers to escape scroll contexts
- Consider `<Teleport to="body">` for overlays

### Vanilla HTML/CSS

- Apply viewport rules to `html, body, #root`
- Use CSS custom properties for consistent spacing:
  ```css
  :root {
    --page-max-width: 1200px;
    --viewport-height: 100dvh;
  }
  ```

### Tailwind CSS

- Use `h-dvh` or `h-screen` for viewport height
- Use `overflow-auto` for scrollable containers
- Use `max-w-screen-xl mx-auto` for centered content with max-width

## Checklist Before Delivering

- [ ] Root container fills viewport with `overflow: hidden`
- [ ] Inner content area has `overflow-y: auto` for scrolling
- [ ] Content has `max-width` on desktop (centered)
- [ ] Tables are in scrollable wrappers
- [ ] Long lists (>20 items) are virtualized or scroll-constrained
- [ ] Modals/drawers use fixed positioning, not flow positioning
- [ ] No nested competing scroll containers
- [ ] Mobile uses `dvh` with `vh` fallback
- [ ] Empty/loading states have `min-height` to prevent collapse
