package main

import (
	"fmt"
	"log"
	"strconv"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type task struct {
	Id int `json:ID`
	Name string `json:Name`
	Content string `json:Content`
}

type allTasks []task

var tasks = allTasks {
	{
		Id: 1,
		Name: "Task One",
		Content: "Some Content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	json.NewEncoder(w).Encode(tasks);
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r);

	taskId, err := strconv.Atoi(vars["id"]);

	if err != nil {
		fmt.Fprintf(w, "Invalid ID");
		return
	}

	for _, task := range tasks {
		if task.Id == taskId {
			w.Header().Set("Content-Type", "application/json");

			json.NewEncoder(w).Encode(task);
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r);

	taskId, err := strconv.Atoi(vars["id"]);

	if err != nil {
		fmt.Fprintf(w, "Invalid ID");
		return
	}

	for i, t := range tasks {
		if t.Id == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...);

			fmt.Fprintf(w, "The task ID %v has been removed successfully", taskId);
		}
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task;
	reqBody, err := ioutil.ReadAll(r.Body);

	if err != nil {
		fmt.Fprintf(w, "Insert a valid task");
	}

	json.Unmarshal(reqBody, &newTask);
	newTask.Id = len(tasks) + 1;
	tasks = append(tasks, newTask);

	w.Header().Set("Content-Type", "application/json");
	w.WriteHeader(http.StatusCreated);
	json.NewEncoder(w).Encode(newTask);
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r);

	taskId, err := strconv.Atoi(vars["id"]);
	var updatedTask task;

	if err != nil {
		fmt.Fprintf(w, "Invalid ID");
	}

	reqBody, err := ioutil.ReadAll(r.Body);
	if err != nil {
		fmt.Fprintf(w, "Please enter valid data");
	}

	json.Unmarshal(reqBody, &updatedTask);

	for i, t := range tasks {
		if t.Id == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...);
			updatedTask.Id = taskId;
			tasks = append(tasks, updatedTask);

			fmt.Fprintf(w, "The task With ID %v has been updated successfully", taskId);
		}
	}
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	fmt.Fprintf(w, "Welcome to my rest api");
}

func main() {
	router := mux.NewRouter().StrictSlash(true);

	router.HandleFunc("/", indexRoute);
	router.HandleFunc("/tasks", getTasks).Methods("GET");
	router.HandleFunc("/tasks", createTask).Methods("POST");
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET");
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE");
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT");

	log.Fatal(http.ListenAndServe(":3000", router));
}