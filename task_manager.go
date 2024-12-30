package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Task struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
}

var (
	tasks []Task
	mu    sync.Mutex
)

const tasksFile = "tasks.txt"

func main() {
	err := loadTasksFromFile()
	if err != nil {
		fmt.Printf("Error loading tasks from file: %v\n", err)
	}

	http.HandleFunc("/tasks", tasksHandler)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		mu.Lock()
		defer mu.Unlock()
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		var newTask Task
		err := json.NewDecoder(r.Body).Decode(&newTask)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		mu.Lock()
		tasks = append(tasks, newTask)
		mu.Unlock()

		err = saveTaskToFile(newTask)
		if err != nil {
			http.Error(w, "Failed to save task to file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Task added successfully"))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func loadTasksFromFile() error {
	file, err := os.Open(tasksFile)
	if err != nil {

		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	var loadedTasks []Task
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "|", 3)
		if len(parts) != 3 {
			fmt.Printf("Skipping malformed line: %s\n", line)
			continue
		}
		loadedTasks = append(loadedTasks, Task{
			Title:       parts[0],
			Subtitle:    parts[1],
			Description: parts[2],
		})
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	mu.Lock()
	tasks = loadedTasks
	mu.Unlock()
	return nil
}

func saveTaskToFile(task Task) error {

	file, err := os.OpenFile(tasksFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	taskLine := fmt.Sprintf("%s|%s|%s\n", task.Title, task.Subtitle, task.Description)

	_, err = file.WriteString(taskLine)
	if err != nil {
		return err
	}

	return nil
}
