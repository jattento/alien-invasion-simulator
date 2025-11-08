package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Status represents the status of a todo item
type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
)

// Todo represents a single todo item
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// TodoList manages a collection of todos
type TodoList struct {
	Todos    []Todo `json:"todos"`
	NextID   int    `json:"next_id"`
	filePath string
}

// NewTodoList creates a new TodoList instance
func NewTodoList(filePath string) *TodoList {
	return &TodoList{
		Todos:    []Todo{},
		NextID:   1,
		filePath: filePath,
	}
}

// Load loads todos from the file
func (tl *TodoList) Load() error {
	data, err := os.ReadFile(tl.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, that's okay
		}
		return fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return nil // Empty file, that's okay
	}

	if err := json.Unmarshal(data, tl); err != nil {
		return fmt.Errorf("failed to unmarshal todos: %w", err)
	}

	return nil
}

// Save saves todos to the file
func (tl *TodoList) Save() error {
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal todos: %w", err)
	}

	if err := os.WriteFile(tl.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Add adds a new todo item
func (tl *TodoList) Add(title, description string) (*Todo, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	todo := Todo{
		ID:          tl.NextID,
		Title:       title,
		Description: description,
		Status:      StatusPending,
		CreatedAt:   time.Now(),
	}

	tl.Todos = append(tl.Todos, todo)
	tl.NextID++

	if err := tl.Save(); err != nil {
		return nil, err
	}

	return &todo, nil
}

// Get retrieves a todo by ID
func (tl *TodoList) Get(id int) (*Todo, error) {
	for i := range tl.Todos {
		if tl.Todos[i].ID == id {
			return &tl.Todos[i], nil
		}
	}
	return nil, fmt.Errorf("todo with ID %d not found", id)
}

// List returns all todos
func (tl *TodoList) List() []Todo {
	return tl.Todos
}

// ListByStatus returns todos filtered by status
func (tl *TodoList) ListByStatus(status Status) []Todo {
	var filtered []Todo
	for _, todo := range tl.Todos {
		if todo.Status == status {
			filtered = append(filtered, todo)
		}
	}
	return filtered
}

// Complete marks a todo as completed
func (tl *TodoList) Complete(id int) error {
	for i := range tl.Todos {
		if tl.Todos[i].ID == id {
			if tl.Todos[i].Status == StatusCompleted {
				return fmt.Errorf("todo %d is already completed", id)
			}
			now := time.Now()
			tl.Todos[i].Status = StatusCompleted
			tl.Todos[i].CompletedAt = &now
			return tl.Save()
		}
	}
	return fmt.Errorf("todo with ID %d not found", id)
}

// Update updates a todo's title and description
func (tl *TodoList) Update(id int, title, description string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}

	for i := range tl.Todos {
		if tl.Todos[i].ID == id {
			tl.Todos[i].Title = title
			tl.Todos[i].Description = description
			return tl.Save()
		}
	}
	return fmt.Errorf("todo with ID %d not found", id)
}

// Delete removes a todo
func (tl *TodoList) Delete(id int) error {
	for i := range tl.Todos {
		if tl.Todos[i].ID == id {
			tl.Todos = append(tl.Todos[:i], tl.Todos[i+1:]...)
			return tl.Save()
		}
	}
	return fmt.Errorf("todo with ID %d not found", id)
}
