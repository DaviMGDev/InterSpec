# InterSpec Examples

## Interactive List with Selection Counter
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

## Navigation with Page Parameters
```interspec
page Main() {
    Button("View Profile") {
        on click {
            navigate UserProfile(userId: 42)
        }
    }
}

page UserProfile(userId) {
    Text("Viewing user: ${userId}")
}
```

## Modal with Toggle State
```interspec
@* Modal appears as an overlay on top of the page content.
   On mobile, it should fill the screen. On desktop,
   show it as a centered dialog with a backdrop. *@
page Main() {
    state modalOpen = false

    @ Primary trigger — clearly indicate this opens a dialog
    Button("Open Modal") {
        on click {
            toggle(modalOpen)
        }
    }

    if modalOpen {
        Modal("Confirmation") {
            @ Destructive action — give "Confirm" prominent warning styling
            Text("Are you sure?")
            Button("Confirm") {
                on click {
                    toggle(modalOpen)
                    Toast("Confirmed!")
                }
            }
            Button("Cancel") {
                on click {
                    toggle(modalOpen)
                }
            }
        }
    }
}
```

## Reusable Component with Children
```interspec
component Section(title) {
    Text(title)
    Text("---")
}

page Main() {
    column {
        Section("First Section") {
            Text("Child appended AFTER '---'")
        }
    }
}
```

## Table with Empty State
```interspec
@*
  Desktop: show the full table with all columns visible.
  Mobile: collapse to a stacked card layout (label-value pairs).
  The EmptyState should be centered with breathing room. *@
page Main() {
    state users = [["Alice", "Admin"], ["Bob", "User"]]

    column {
        Text("User List") {
            weight: horizontal
        }

        if users.length > 0 {
            Table(["Name", "Role"], users)
        } else {
            EmptyState("No users found") {
                Button("Add User") {
                    on click { navigate NewUser() }
                }
            }
        }
    }
}
```

## Tabs and Accordion
```interspec
page Main() {
    column {
        // Tabs — children match tab labels by position
        Tabs(["Overview", "Settings", "Logs"]) {
            Text("Dashboard overview content")
            Text("Application settings")
            Text("Recent activity logs")
        }

        // Accordion — only one section open at a time
        Accordion(["What is InterSpec?", "How do I start?"]) {
            Text("A declarative UI specification language.")
            Text("Create a Main() page and run isc.")
        }
    }
}
```

## DropdownMenu vs Select
```interspec
page Main() {
    state selected = ""

    column {
        // Select — returns a value (form input)
        Select(["Apple", "Banana", "Cherry"]) {
            on commit {
                selected = value
            }
        }

        // DropdownMenu — triggers actions (command menu)
        DropdownMenu("Actions") {
            Button("Edit")   { on click { navigate Editor() } }
            Button("Share")  { on click { Dialog("Share this item") } }
            Button("Delete") { on click { Dialog("Confirm delete") } }
        }
    }
}
```

## Form with Validation
```interspec
page Main() {
    state email = ""
    state agreed = false
    state submitting = false

    Form {
        on submit {
            validate()
        }

        Input("Email") {
            required: true
            placeholder: "you@example.com"
            error: email == "" ? "Email is required" : false
            on input {
                email = value
            }
            on commit {
                log("Email committed: " + value)
            }
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
```

## Stepper and Pagination
```interspec
page Main() {
    state currentStep = 0

    column {
        // Stepper — read-only flow indicator
        Stepper(["Cart", "Shipping", "Payment"]) {
            current: currentStep
        }

        // Pagination — page navigation
        Pagination {
            current: 1
            total: 10
        }
    }
}
```

## Breadcrumb and Link
```interspec
page Main() {
    column {
        @ Breadcrumb: small muted text, current page ("Widget") not clickable
        Breadcrumb(["Home", "Products", "Widget"])

        @ Link: underlined text, no button styling — distinct from action buttons
        Link("Back to Products") {
            on click {
                navigate Products()
            }
        }

        @ Primary CTA — filled button, full-width on mobile
        Button("Buy Now") {
            on click {
                navigate Checkout()
            }
        }
    }
}
```

## Progress and Badge
```interspec
page Main() {
    column {
        // Progress — read-only indicator
        Progress(60) {
            max: 100
        }

        // Badge — status indicator with variants
        Badge("New") {
            variant: success
        }
        Badge("Pending") {
            variant: warning
        }
        Badge("3") {
            variant: error
        }
    }
}
```

## Delay and Reset
```interspec
page Main() {
    state message = "Waiting..."
    state count = 0

    column {
        Text(message)
        Text("Count: ${count}")

        Button("Increment") {
            on click {
                count = count + 1
            }
        }

        Button("Reset counter") {
            on click {
                reset(count)   // resets to declared initial value (0)
            }
        }

        Button("Show delayed message") {
            on click {
                delay(2000, message = "Delayed!")
            }
        }
    }
}
```

## Keyboard Shortcuts
```interspec
page Main() {
    state searchTerm = ""

    Input("Search") {
        placeholder: "Type and press Enter"
        on key("Enter") {
            log("Searching: " + searchTerm)
            Toast("Search submitted")
        }
        on key("Escape") {
            searchTerm = ""
            Toast("Search cleared")
        }
    }
}
```
