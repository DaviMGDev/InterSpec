# InterSpec Built-in Catalog

## Components

| Component    | Parameters  | Description |
|--------------|-------------|-------------|
| `Button`     | `label`     | Clickable button. |
| `Text`       | `content`   | Renders a string. |
| `Input`      | `placeholder` | Single-line text input. |
| `Select`     | `options`   | Dropdown. Pass array of strings or objects. |
| `Checkbox`   | `label`     | Toggleable checkbox with label. |
| `Toggle`     | `label`     | Boolean on/off switch. |
| `Slider`     | â€”           | Range slider for numeric selection. |
| `Image`      | `src`       | Image placeholder. Takes a URL string. |
| `Icon`       | `name`      | Semantic icon (e.g. "settings", "user"). |
| `Alert`      | `message`   | Inline system message. Use `variant` property. |
| `Card`       | `title`     | Content container. Neutral â€” children append after title. |
| `Modal`      | `title`     | Modal overlay that blocks page interaction. |
| `Dialog`     | `title`     | Confirmation/info dialog requiring user action. |
| `Toast`      | `message`   | Brief auto-dismissing popup. |
| `Tooltip`    | `content`   | Hover-revealed tooltip. |

## Events

| Event    | Applies To                          |
|----------|-------------------------------------|
| `click`  | Button, Checkbox, Toggle, Card      |
| `hover`  | Button, Tooltip, Card, Icon         |
| `input`  | Input, Select                        |
| `focus`  | Input, Button, Select               |
| `blur`   | Input, Button, Select               |
| `open`   | Modal, Dialog, Toast                |
| `close`  | Modal, Dialog, Toast                |

## Actions

| Action     | Syntax                          |
|------------|---------------------------------|
| `navigate` | `navigate PageName(param: val)` |
| `back`     | `back()`                        |
| `toggle`   | `toggle(variable)`              |
| `log`      | `log(message)`                  |

## Component Properties

| Property      | Applies To                       | Values |
|---------------|----------------------------------|--------|
| `align`       | Any component in row/column      | `(vertical, horizontal)`: top/center/bottom, left/center/right |
| `weight`      | Any component in row/column      | `horizontal`, `vertical`, `both` |
| `wrap`        | `row` layout                     | `true`, `false` |
| `placeholder` | `Input`                          | string |
| `required`    | Input, Select, Checkbox          | `true`, `false` |
| `disabled`    | Any interactive component        | `true`, `false` |
| `variant`     | `Alert`                          | `info`, `success`, `warning`, `error` |
| `src`         | `Image`                          | string (URL) |
