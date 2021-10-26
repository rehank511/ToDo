package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	Tid int64 `json:"ID"`
	Text string `json:"Text"`
	CreateTime time.Time `json:"Time"`
}

var counter int64

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db, err := sql.Open("sqlite3", "./tododata.db")
	checkErr(err)
	defer db.Close()
	var todos []*Todo
	rows, err := db.Query("SELECT * FROM Todos")
	checkErr(err)
	for rows.Next() {
		t := new(Todo)
		err := rows.Scan(&t.Tid, &t.Text, &t.CreateTime)
		checkErr(err)
		todos = append(todos, t)
	}
	json.NewEncoder(w).Encode(todos)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	db, err := sql.Open("sqlite3", "./tododata.db")
	checkErr(err)
	var todos []*Todo
	params := mux.Vars(r)

	counter += 1
	stmt, err := db.Prepare("INSERT INTO Todos(Text, CreateTime) VALUES(?,datetime('now'))")
	checkErr(err)
	_ , err = stmt.Exec(params["text"])
	checkErr(err)


	t := new(Todo)
	t.Tid = counter
	t.Text = params["text"]
	t.CreateTime = time.Now()
	todos = append(todos, t)
	db.Close()

	json.NewEncoder(w).Encode(todos)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	db, err := sql.Open("sqlite3", "./tododata.db")
	checkErr(err)
	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM Todos WHERE Tid=?")
	checkErr(err)
	_ , err = stmt.Exec(params["id"])
	checkErr(err)
	db.Close()
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	db, err := sql.Open("sqlite3", "./tododata.db")
	checkErr(err)
	params := mux.Vars(r)
	temp := params["text"]
	id := temp[0:1]
	text := temp[1:]

	stmt, err := db.Prepare("UPDATE Todos SET Text=? WHERE Tid=?")
    checkErr(err)
	_ , err = stmt.Exec(text, id)
    checkErr(err)
	db.Close()
}

func createTable(db *sql.DB) {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS Todos (Tid INTEGER PRIMARY KEY AUTOINCREMENT, Text TEXT, CreateTime DATETIME)")
	checkErr(err)
	statement.Exec()
	statement.Close()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() { 
	var _, err = os.Stat("./tododata.db")

	if os.IsNotExist(err) {
		log.Println("Creating Database File...")
		file, err := os.Create("tododata.db")
		checkErr(err)
		file.Close()
	}

	db, err := sql.Open("sqlite3", "./tododata.db")
	checkErr(err)
	createTable(db)
	db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/todo", homePage).Methods("GET")
	router.HandleFunc("/todo/{text}", addTodo).Methods("POST")
	router.HandleFunc("/todo/{id}", deleteTodo).Methods("DELETE")
	router.HandleFunc("/todo/{text}", updateTodo).Methods("PUT")
	
	fmt.Println("Starting New Server")
	log.Fatal(http.ListenAndServe(":8000", router))
}