package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var Connector *gorm.DB                                                                                                                                                                                                                                                                                                                 

type Todo struct {
	ID string `json:"ID"`
	Text string `json:"Description"`
	CreateTime time.Time `json:"CreationTime"`
}

var todos []Todo

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	todo.ID = strconv.Itoa(len(todos) + 1)
	todos = append(todos, todo)
	json.NewEncoder(w).Encode(todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range todos {
		if item.ID == params["id"] {
			todos = append(todos[:index], todos[index+1:]...)
			break
		}
	}
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	// get id from user
	// delete that todo from list
	// get the updaetd text
	// make a new todo with new text
	
}

func main() { 
	router := mux.NewRouter()

	layout := "01-02-2006 15:04:05"
	timenow := time.Now().Format(layout)
	now, error := time.Parse(layout, timenow)
	if error != nil {
		fmt.Println(error)
	}

	todos = append(todos, Todo{ID: "1", Text: "Get Groceries", CreateTime: now})

	router.HandleFunc("/todo", homePage).Methods("GET")
	router.HandleFunc("/todo", addTodo).Methods("POST")
	router.HandleFunc("/todo/{id}", deleteTodo).Methods("DELETE")
	router.HandleFunc("todo/{id}", updateTodo).Methods("PUT")
	
	fmt.Println("Starting New Server")
	log.Fatal(http.ListenAndServe(":3000", router))
}