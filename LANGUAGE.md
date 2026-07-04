Here is the newly organized, official specification for **InterSpec**. It is structured as a clean, professional language reference, incorporating all your design decisions, removing the `$` sigil for variables, and strictly documenting the "no magic" component behavior.

***

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
    wrap: true // Optional: allows items to wrap to the next line
    Text("Item 1")
    Text("Item 2")
}

column {
    Text("Item 1")
    Text("Item 2")
}
```

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
| `Select` | `options` | A dropdown picker. Pass an array of strings or objects as options. |
| `Checkbox` | `label` | A toggleable checkbox with an adjacent label. |
| `Toggle` | `label` | A boolean on/off toggle switch with a label. |
| `Slider` | — | A range slider for numeric selection. |
| `Image` | `src` | An image placeholder. Takes a URL string for `src`. |
| `Icon` | `name` | A semantic icon indicator (e.g., `"settings"`, `"user"`, `"search"`). The actual icon set is determined by the runtime. |
| `Alert` | `message` | An inline system message for success, error, warning, or info states. Use the `variant` property to set the semantic type. |
| `Card` | `title` | A content container that groups related information. Card is a neutral component (no wrapper node) — children are appended after the title. |
| `Modal` | `title` | A modal overlay dialog that blocks interaction with the rest of the page. |
| `Dialog` | `title` | A confirmation or information dialog that requires user action. |
| `Toast` | `message` | A brief, auto-dismissing pop-up notification. |
| `Tooltip` | `content` | A hover-revealed tooltip that displays contextual information. |

#### Usage Examples

```interspec
// Card — groups related content
Card("User Profile") {
    Text("Name: Alice")
    Text("Role: Admin")
}

// Input — text entry field
Input("Enter your name...")

// Select — dropdown with options
Select(["Apple", "Banana", "Cherry"])

// Image — placeholder with source URL
Image("https://example.com/photo.png")

// Icon — semantic indicator
Icon("settings")

// Alert — inline system message
Alert("File saved successfully") {
    variant: success
}
```

### Events

| Event | Applies To | Description |
|-------|------------|-------------|
| `click` | `Button`, `Checkbox`, `Toggle`, `Card` | Fired when the user clicks or taps the component. |
| `hover` | `Button`, `Tooltip`, `Card`, `Icon` | Fired when the user hovers the cursor over the component. |
| `input` | `Input`, `Select` | Fired on every value change (e.g., each keystroke in an `Input`). |
| `focus` | `Input`, `Button`, `Select` | Fired when the component receives focus. |
| `blur` | `Input`, `Button`, `Select` | Fired when the component loses focus. |
| `open` | `Modal`, `Dialog`, `Toast` | Fired when the component becomes visible. |
| `close` | `Modal`, `Dialog`, `Toast` | Fired when the component is dismissed or hidden. |

### Actions

| Action | Syntax | Description |
|--------|--------|-------------|
| `navigate` | `navigate PageName(param: value)` | Navigate to another page, optionally passing parameters. |
| `back` | `back()` | Navigate to the previous page. No-op if there is no history. |
| `toggle` | `toggle(variable)` | Shortcut for `variable = !variable`. Works with any boolean state variable. |
| `log` | `log(message)` | Print a debug message to the runtime console. No effect in production. |

### Component Properties

| Property | Applies To | Values | Description |
|----------|------------|--------|-------------|
| `align` | Any component in a `row` or `column` | `(vertical, horizontal)` where vertical: `top`, `center`, `bottom`; horizontal: `left`, `center`, `right` | Alignment within the parent layout. |
| `weight` | Any component in a `row` or `column` | `horizontal`, `vertical`, `both` | How the component fills available space. |
| `wrap` | `row` layout | `true`, `false` | Whether items wrap to the next line when they exceed the container width. |
| `placeholder` | `Input` | `string` | Hint text displayed inside the field when it is empty. |
| `required` | `Input`, `Select`, `Checkbox` | `true`, `false` | Marks the field as required for form validation. |
| `disabled` | Any interactive component | `true`, `false` | Disables the component, preventing interaction. |
| `variant` | `Alert` | `info`, `success`, `warning`, `error` | Semantic variant indicating the type or severity of the message. |
| `src` | `Image` | `string` | Source URL for the image content. |

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
