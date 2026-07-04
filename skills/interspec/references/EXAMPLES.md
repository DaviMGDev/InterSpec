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
page Main() {
    state modalOpen = false

    Button("Open Modal") {
        on click {
            toggle(modalOpen)
        }
    }

    if modalOpen {
        Modal("Confirmation") {
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
