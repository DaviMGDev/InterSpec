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
```interspec
Button("Hover me") {
    on hover {
        Tooltip("This is a button")
    }
    on click {
        navigate AnotherPage()
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
*(This list will expand over time)*

**Components:**  
`Button`, `Text`, `Modal`, `Toast`, `Dialog`, `Tooltip`, `Checkbox`, `Slider`, `Toggle`

**Events:**  
`click`, `hover`

**Actions:**  
`navigate`

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
