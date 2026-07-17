# InterSpec Built-in Catalog

## Components

| Component    | Parameters  | Description |
|--------------|-------------|-------------|
| `Button`     | `label`     | Clickable button. |
| `Text`       | `content`   | Renders a string. |
| `Input`      | `placeholder` | Single-line text input. |
| `DatePicker` | `placeholder` | Date selection input that opens a calendar picker. Use `min` and `max` properties to constrain the date range. |
| `Select`     | `options`   | Dropdown for picking a value. Pass array of strings or objects. Distinct from `DropdownMenu`. |
| `Checkbox`   | `label`     | Toggleable checkbox with label. |
| `Toggle`     | `label`     | Boolean on/off switch. |
| `Slider`     | —           | Range slider for numeric selection. |
| `Image`      | `src`       | Image placeholder. Takes a URL string. |
| `Icon`       | `name`      | Semantic icon (e.g. "settings", "user"). |
| `Alert`      | `message`   | Inline system message. Use `variant` property. |
| `Card`       | `title`     | Content container. Neutral — children append after title. |
| `Modal`      | `title`     | Modal overlay that blocks page interaction. |
| `Dialog`     | `title`     | Confirmation/info dialog requiring user action. |
| `Drawer`     | `title`     | Side panel that slides in from the edge for navigation, filters, or supplementary content. Use the `side` property to set the origin edge (`left` or `right`). |
| `Toast`      | `message`   | Brief auto-dismissing popup. |
| `Tooltip`    | `content`   | Hover-revealed tooltip. |
| `Table`      | `columns`, `rows` | Tabular data. `columns` is array of header strings. `rows` is array of row arrays. |
| `Tabs`       | `tabs`      | Tabbed panels. `tabs` is array of labels. Children are panel bodies, matched by position. |
| `Accordion`  | `items`     | Expandable sections. `items` is array of titles. Children are bodies, matched by position. |
| `TreeView`   | `items`     | Hierarchical tree display for nested data. Children define node content. Events: `click`, `open`, `close`. |
| `Badge`      | `label`     | Status/count/category indicator. Use `variant` property. |
| `Link`       | `label`     | Navigational element (no href). Use `on click { navigate ... }`. |
| `Progress`   | `value`     | Read-only progress. Use `max` property (default 100). |
| `EmptyState` | `message`   | Placeholder for empty lists. Children can provide a recovery action. |
| `Breadcrumb` | `items`     | Navigation path trail. Array of label strings. Last item is current page. |
| `Stepper`    | `steps`     | Multi-step flow indicator (read-only). Use `current` property. |
| `Pagination` | —           | Page navigation. Use `current` and `total` properties. |
| `DropdownMenu` | `label`  | Action-triggering command menu. Children are menu items (Buttons). Distinct from `Select`. |
| `Divider`    | —           | Horizontal divider that visually separates content sections. Purely structural — no children or interactivity. |
| `Section`    | `title`     | Structural grouping under a heading. |
| `Form`       | —           | Groups inputs under a submission action. Fires `submit` event. |
| `FileUpload` | `label`     | File upload control with picker dialog or drop zone. Use `accept` to filter types, `multiple` for multi-file. |

## Events

| Event    | Applies To                          |
|----------|-------------------------------------|
| `click`  | Button, Checkbox, Toggle, Card      |
| `hover`  | Button, Tooltip, Card, Icon         |
| `input`  | Input, Select, DatePicker, FileUpload                       |
| `commit` | Input, Select, DatePicker, FileUpload                       |
| `focus`  | Input, Button, Select, DatePicker               |
| `blur`   | Input, Button, Select, DatePicker               |
| `open`   | Modal, Dialog, Drawer, Toast, TreeView                |
| `close`  | Modal, Dialog, Drawer, Toast, TreeView                |
| `submit` | Form                                |
| `key`    | Input (takes key name, e.g. `on key("Enter")`) |
| `longpress` | Any interactive component        |
| `reachEnd` | Scrollable column, row           |

## Actions

| Action     | Syntax                          |
|------------|---------------------------------|
| `navigate` | `navigate PageName(param: val)` |
| `back`     | `back()`                        |
| `toggle`   | `toggle(variable)`              |
| `log`      | `log(message)`                  |
| `validate` | `validate()`                    |
| `reset`    | `reset(variable)`               |
| `delay`    | `delay(ms, action)`             |

## Component Properties

| Property      | Applies To                       | Values |
|---------------|----------------------------------|--------|
| `align`       | Any component in row/column      | `(vertical, horizontal)`: top/center/bottom, left/center/right |
| `weight`      | Any component in row/column      | `horizontal`, `vertical`, `both` |
| `wrap`        | `row` layout                     | `true`, `false` |
| `collapse`    | `row` layout                     | `true`, `false` |
| `scrollable`  | `row`, `column`                  | `true`, `false` | Content scrolls when it exceeds available space. Implementer MUST constrain height and apply overflow. See [Viewport Safety](../../LANGUAGE.md#11-viewport-safety). |
| `placeholder` | `Input`                          | string |
| `required`    | Input, Select, Checkbox          | `true`, `false` |
| `side`        | Drawer                           | `left`, `right` |
| `accept`      | FileUpload                       | string | Accepted MIME types (e.g. `"image/*"`). |
| `multiple`    | FileUpload                       | `true`, `false` | Allow multiple file selection. |
| `disabled`    | Any interactive component        | `true`, `false` |
| `loading`     | Button, Card, Table, Image       | `true`, `false` |
| `error`       | Input, Select                    | `true`, `false`, or message string |
| `variant`     | Alert, Badge                     | `info`, `success`, `warning`, `error` |
| `src`         | `Image`                          | string (URL) |
| `min`         | DatePicker, Progress             | number | Lower bound (date string or 0 for progress). |
| `max`         | DatePicker, Progress             | number | Upper bound (date string or 100 for progress). |
| `current`     | Stepper, Pagination              | number (0-indexed for Stepper, 1-indexed for Pagination) |
| `total`       | `Pagination`                     | number |

## Hints

Any component, layout, or page can carry `@` hints — freeform annotations that
survive comment stripping and guide the implementer without affecting runtime behavior.

| Form | Syntax | Example |
|------|--------|---------|
| Single-line | `@ text` | `@ Primary action — make this button prominent` |
| Block | `@* ... *@` | `@* On mobile, stack vertically. *@` |

Hints have **no grammar, no validation, and no runtime effect.** They are
freeform prose for the developer or AI translating the spec into a real UI.
See [the language spec](../../LANGUAGE.md#10-hints-and-annotations) for full
details.
