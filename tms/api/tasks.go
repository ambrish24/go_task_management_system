package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"tms.zinkworks.com/model"
)

func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching All Tasks")

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	filePath := filepath.Join(wd, "tasks.json")

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	var tasks []model.Task
	err = json.Unmarshal(jsonData, &tasks)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, tasks, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

func (app *application) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Println("Fetching Task with ID:", id)

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	filePath := filepath.Join(wd, "tasks.json")

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	var tasks []model.Task
	err = json.Unmarshal(jsonData, &tasks)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	desiredID := id

	// Find the task with the desired ID
	var foundTask *model.Task
	for _, task := range tasks {
		if int64(task.ID) == desiredID {
			foundTask = &task
			break
		}
	}

	// Check if the task was found
	if foundTask == nil {
		fmt.Printf("Task with ID %d not found\n", desiredID)
		return
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, foundTask, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
