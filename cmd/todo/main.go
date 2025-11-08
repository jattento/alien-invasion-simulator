package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/jattento/alien-invasion-simulator/internal/todo"
)

const defaultTodoFile = ".todos.json"

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	todoFilePath := filepath.Join(homeDir, defaultTodoFile)

	// Define commands
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addTitle := addCmd.String("title", "", "Title of the todo (required)")
	addDesc := addCmd.String("desc", "", "Description of the todo")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listStatus := listCmd.String("status", "", "Filter by status: pending, completed")

	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)
	completeID := completeCmd.Int("id", 0, "ID of the todo to complete (required)")

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateID := updateCmd.Int("id", 0, "ID of the todo to update (required)")
	updateTitle := updateCmd.String("title", "", "New title (required)")
	updateDesc := updateCmd.String("desc", "", "New description")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteID := deleteCmd.Int("id", 0, "ID of the todo to delete (required)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Create todo list
	todoList := todo.NewTodoList(todoFilePath)
	if err := todoList.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading todos: %v\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addTitle == "" {
			fmt.Fprintln(os.Stderr, "Error: title is required")
			addCmd.PrintDefaults()
			os.Exit(1)
		}
		handleAdd(todoList, *addTitle, *addDesc)

	case "list":
		listCmd.Parse(os.Args[2:])
		handleList(todoList, *listStatus)

	case "complete":
		completeCmd.Parse(os.Args[2:])
		if *completeID == 0 {
			fmt.Fprintln(os.Stderr, "Error: id is required")
			completeCmd.PrintDefaults()
			os.Exit(1)
		}
		handleComplete(todoList, *completeID)

	case "update":
		updateCmd.Parse(os.Args[2:])
		if *updateID == 0 || *updateTitle == "" {
			fmt.Fprintln(os.Stderr, "Error: id and title are required")
			updateCmd.PrintDefaults()
			os.Exit(1)
		}
		handleUpdate(todoList, *updateID, *updateTitle, *updateDesc)

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deleteID == 0 {
			fmt.Fprintln(os.Stderr, "Error: id is required")
			deleteCmd.PrintDefaults()
			os.Exit(1)
		}
		handleDelete(todoList, *deleteID)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Todo App - A simple command-line todo list manager")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  todo <command> [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add       Add a new todo")
	fmt.Println("  list      List all todos")
	fmt.Println("  complete  Mark a todo as completed")
	fmt.Println("  update    Update a todo")
	fmt.Println("  delete    Delete a todo")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  todo add -title=\"Buy groceries\" -desc=\"Milk, eggs, bread\"")
	fmt.Println("  todo list")
	fmt.Println("  todo list -status=pending")
	fmt.Println("  todo complete -id=1")
	fmt.Println("  todo update -id=1 -title=\"Buy groceries\" -desc=\"Milk, eggs, bread, cheese\"")
	fmt.Println("  todo delete -id=1")
}

func handleAdd(todoList *todo.TodoList, title, description string) {
	newTodo, err := todoList.Add(title, description)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding todo: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Todo added successfully (ID: %d)\n", newTodo.ID)
}

func handleList(todoList *todo.TodoList, status string) {
	var todos []todo.Todo

	if status != "" {
		var todoStatus todo.Status
		switch strings.ToLower(status) {
		case "pending":
			todoStatus = todo.StatusPending
		case "completed":
			todoStatus = todo.StatusCompleted
		default:
			fmt.Fprintf(os.Stderr, "Invalid status: %s (use 'pending' or 'completed')\n", status)
			os.Exit(1)
		}
		todos = todoList.ListByStatus(todoStatus)
	} else {
		todos = todoList.List()
	}

	if len(todos) == 0 {
		fmt.Println("No todos found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tStatus\tTitle\tDescription")
	fmt.Fprintln(w, "--\t------\t-----\t-----------")

	for _, t := range todos {
		statusIcon := "○"
		if t.Status == todo.StatusCompleted {
			statusIcon = "●"
		}
		desc := t.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", t.ID, statusIcon, t.Title, desc)
	}
	w.Flush()
}

func handleComplete(todoList *todo.TodoList, id int) {
	if err := todoList.Complete(id); err != nil {
		fmt.Fprintf(os.Stderr, "Error completing todo: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Todo %d marked as completed\n", id)
}

func handleUpdate(todoList *todo.TodoList, id int, title, description string) {
	if err := todoList.Update(id, title, description); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating todo: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Todo %d updated successfully\n", id)
}

func handleDelete(todoList *todo.TodoList, id int) {
	if err := todoList.Delete(id); err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting todo: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Todo %d deleted successfully\n", id)
}
