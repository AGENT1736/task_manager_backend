package main

import (
	"encoding/json"
	"net/http"
)

type Task struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
}

func main() {
	http.HandleFunc("/tasks", tasksHandler)
	http.ListenAndServe(":8080", nil)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := []Task{
		{"Task 1", "Complete the report", "we need to complete the reports"},
		{"Task 2", "Review PRs", "we need to complete these PRs"},
		{"Task 3", "Team meeting at 3 PM", "we have a team meeting at 3PM"},
		{"Task 4", "Client call at 5 PM", "we have a client call at 5PM"},
	}

	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
