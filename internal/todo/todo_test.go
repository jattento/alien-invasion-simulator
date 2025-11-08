package todo

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewTodoList(t *testing.T) {
	tl := NewTodoList("test.json")
	if tl == nil {
		t.Fatal("NewTodoList returned nil")
	}
	if len(tl.Todos) != 0 {
		t.Errorf("Expected empty todos, got %d", len(tl.Todos))
	}
	if tl.NextID != 1 {
		t.Errorf("Expected NextID to be 1, got %d", tl.NextID)
	}
}

func TestAdd(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "todos.json")
	tl := NewTodoList(filePath)

	todo, err := tl.Add("Test Todo", "This is a test")
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if todo.ID != 1 {
		t.Errorf("Expected ID 1, got %d", todo.ID)
	}
	if todo.Title != "Test Todo" {
		t.Errorf("Expected title 'Test Todo', got '%s'", todo.Title)
	}
	if todo.Status != StatusPending {
		t.Errorf("Expected status pending, got %s", todo.Status)
	}
	if len(tl.Todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(tl.Todos))
	}
}

func TestAddEmptyTitle(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "todos.json")
	tl := NewTodoList(filePath)

	_, err := tl.Add("", "Description")
	if err == nil {
		t.Error("Expected error for empty title, got nil")
	}
}

func TestGet(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "todos.json")
	tl := NewTodoList(filePath)

	added, _ := tl.Add("Test Todo", "Description")

	retrieved, err := tl.Get(added.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.Title != added.Title {
		t.Errorf("Expected title '%s', got '%s'", added.Title, retrieved.Title)
	}
}

func TestGetNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "todos.json")
	tl := NewTodoList(filePath)

	_, err := tl.Get(999)
	if err == nil {
		t.Error("Expected error for non-existent ID, got nil")
	}
}

func TestList(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "todos.json")
	tl := NewTodoList(filePath)

	tl.Add("Todo 1", "Description 1")
	tl.Add("Todo 2", "Description 2")
	tl.Add("Todo 3", "Description 3")

	todos := tl.List()
	if len(todos) != 3 {
		t.Errorf("Expected 3 todos, got %d", len(todos))
	}
}

func TestListByStatus(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "todos.json")
	tl := NewTodoList(filePath)

	todo1, _ := tl.Add("Todo 1", "Description 1")
	tl.Add("Todo 2", "Description 2")
	tl.Add("Todo 3", "Description 3")

	tl.Complete(todo1.ID)

	pending := tl.ListByStatus(StatusPending)
	if len(pending) != 2 {
		t.Errorf("Expected 2 pending todos, got %d", len(pending))
	}

	completed := tl.ListByStatus(StatusCompleted)
	if len(completed) != 1 {
		t.Errorf("Expected 1 completed todo, got %d", len(completed))
	}
}

func TestComplete(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "todos.json")
	tl := NewTodoList(filePath)

	added, _ := tl.Add("Test Todo", "Description")

	err := tl.Complete(added.ID)
	if err != nil {
		t.Fatalf("Complete failed: %v", err)
	}

	retrieved, _ := tl.Get(added.ID)
	if retrieved.Status != StatusCompleted {
		t.Errorf("Expected status completed, got %s", retrieved.Status)
	}
	if retrieved.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set")
	}
}

func TestCompleteAlreadyCompleted(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "todos.json")
	tl := NewTodoList(filePath)

	added, _ := tl.Add("Test Todo", "Description")
	tl.Complete(added.ID)

	err := tl.Complete(added.ID)
	if err == nil {
		t.Error("Expected error when completing already completed todo")
	}
}

func TestUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "todos.json")
	tl := NewTodoList(filePath)

	added, _ := tl.Add("Original Title", "Original Description")

	err := tl.Update(added.ID, "Updated Title", "Updated Description")
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, _ := tl.Get(added.ID)
	if updated.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", updated.Title)
	}
	if updated.Description != "Updated Description" {
		t.Errorf("Expected description 'Updated Description', got '%s'", updated.Description)
	}
}

func TestUpdateEmptyTitle(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "todos.json")
	tl := NewTodoList(filePath)

	added, _ := tl.Add("Original Title", "Original Description")

	err := tl.Update(added.ID, "", "Updated Description")
	if err == nil {
		t.Error("Expected error for empty title, got nil")
	}
}

func TestDelete(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "todos.json")
	tl := NewTodoList(filePath)

	added, _ := tl.Add("Test Todo", "Description")

	err := tl.Delete(added.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if len(tl.Todos) != 0 {
		t.Errorf("Expected 0 todos after delete, got %d", len(tl.Todos))
	}

	_, err = tl.Get(added.ID)
	if err == nil {
		t.Error("Expected error when getting deleted todo")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "todos.json")

	// Create and save todos
	tl1 := NewTodoList(filePath)
	tl1.Add("Todo 1", "Description 1")
	tl1.Add("Todo 2", "Description 2")

	// Load in a new instance
	tl2 := NewTodoList(filePath)
	err := tl2.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(tl2.Todos) != 2 {
		t.Errorf("Expected 2 todos after load, got %d", len(tl2.Todos))
	}
	if tl2.NextID != 3 {
		t.Errorf("Expected NextID to be 3, got %d", tl2.NextID)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "nonexistent.json")

	tl := NewTodoList(filePath)
	err := tl.Load()
	if err != nil {
		t.Errorf("Load should not fail for non-existent file: %v", err)
	}
}

func TestLoadEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "empty.json")

	// Create empty file
	os.WriteFile(filePath, []byte{}, 0644)

	tl := NewTodoList(filePath)
	err := tl.Load()
	if err != nil {
		t.Errorf("Load should not fail for empty file: %v", err)
	}
}

func TestPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "todos.json")

	// Create and add todos
	tl1 := NewTodoList(filePath)
	todo1, _ := tl1.Add("Todo 1", "Description 1")
	tl1.Complete(todo1.ID)

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	// Load in new instance and verify
	tl2 := NewTodoList(filePath)
	tl2.Load()

	retrieved, err := tl2.Get(todo1.ID)
	if err != nil {
		t.Fatalf("Failed to get todo: %v", err)
	}

	if retrieved.Status != StatusCompleted {
		t.Errorf("Expected status completed, got %s", retrieved.Status)
	}
	if retrieved.CompletedAt == nil {
		t.Error("Expected CompletedAt to be persisted")
	}
}
