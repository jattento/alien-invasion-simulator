# Todo App

A simple and efficient command-line todo list manager built in Go.

## Features

- ✓ Add new todos with title and description
- ✓ List all todos or filter by status (pending/completed)
- ✓ Mark todos as completed
- ✓ Update existing todos
- ✓ Delete todos
- ✓ Persistent storage (JSON file in home directory)
- ✓ Clean and intuitive CLI interface

## Installation

### Build from source

```bash
go build -o bin/todo ./cmd/todo/
```

### Add to PATH (optional)

To use the `todo` command from anywhere:

```bash
# Copy to a directory in your PATH
sudo cp bin/todo /usr/local/bin/

# Or add the bin directory to your PATH
export PATH=$PATH:$(pwd)/bin
```

## Usage

### Add a new todo

```bash
todo add -title="Buy groceries" -desc="Milk, eggs, bread"
```

### List all todos

```bash
todo list
```

### List todos by status

```bash
# List only pending todos
todo list -status=pending

# List only completed todos
todo list -status=completed
```

### Complete a todo

```bash
todo complete -id=1
```

### Update a todo

```bash
todo update -id=1 -title="Buy groceries" -desc="Milk, eggs, bread, cheese"
```

### Delete a todo

```bash
todo delete -id=1
```

## Data Storage

Todos are stored in a JSON file located at `~/.todos.json`. This file is automatically created when you add your first todo.

## Status Indicators

- `○` - Pending todo
- `●` - Completed todo

## Examples

```bash
# Add some todos
todo add -title="Write documentation" -desc="Document the new feature"
todo add -title="Fix bug #42" -desc="Memory leak in the parser"
todo add -title="Review PR" -desc="Check the new authentication code"

# List all todos
todo list

# Complete a todo
todo complete -id=1

# List only pending todos
todo list -status=pending

# Update a todo
todo update -id=2 -title="Fix critical bug #42" -desc="Memory leak in parser - high priority"

# Delete a todo
todo delete -id=3
```

## Testing

Run the tests:

```bash
go test ./internal/todo/... -v
```

## Architecture

The todo app follows Go best practices with a clean separation of concerns:

- `internal/todo/` - Core business logic and data structures
  - `todo.go` - Todo data structures and CRUD operations
  - `todo_test.go` - Comprehensive unit tests
- `cmd/todo/` - CLI application
  - `main.go` - Command-line interface and user interaction

## Data Model

```go
type Todo struct {
    ID          int        // Unique identifier
    Title       string     // Todo title (required)
    Description string     // Todo description (optional)
    Status      Status     // pending or completed
    CreatedAt   time.Time  // Creation timestamp
    CompletedAt *time.Time // Completion timestamp (nil if not completed)
}
```

## Contributing

Feel free to submit issues or pull requests to improve the todo app!
