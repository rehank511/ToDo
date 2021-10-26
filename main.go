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
	ID int64 `json:"ID"`
	Text string `json:"Text"`
	CreateTime time.Time `json:"Time"`
}

var counter int64

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("sqlite3", "./tododata.db")
	checkErr(err)
	defer db.Close()
	var todos []*Todo
	rows, err := db.Query("SELECT * FROM Todos")
	checkErr(err)
	for rows.Next() {
		t := new(Todo)
		fmt.Print(rows.Scan())
		err := rows.Scan(&t.ID, &t.Text, &t.CreateTime)
		checkErr(err)
		todos = append(todos, t)
	}
	json.NewEncoder(w).Encode(todos)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("sqlite3", "./tododata.db")
	checkErr(err)
	var todos []*Todo
	params := mux.Vars(r)
	counter += 1
	stmt, err := db.Prepare("INSERT INTO Todos(ID, Text, CreateTime) VALUES(?,?,?)")
	checkErr(err)
	_ , err = stmt.Exec(counter, params["text"], time.Now())
	checkErr(err)


	t := new(Todo)
	t.ID = counter
	t.Text = params["text"]
	t.CreateTime = time.Now()
	todos = append(todos, t)
	db.Close()

	json.NewEncoder(w).Encode(todos)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("sqlite3", "./tododata.db")
	checkErr(err)
	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM Todos WHERE ID=?")
	checkErr(err)
	_ , err = stmt.Exec(params["id"])
	checkErr(err)
	db.Close()
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	// get id from user
	// delete that todo from list
	// get the updaetd text
	// make a new todo with new text
	db, err := sql.Open("sqlite3", "./tododata.db")
	checkErr(err)
	params := mux.Vars(r)

	stmt, err := db.Prepare("update Todos set Text=? where ID=?")
    checkErr(err)
	_ , err = stmt.Exec(params["text"], params["id"])
    checkErr(err)
	db.Close()

}

func createTable(db *sql.DB) {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS Todos (ID INT AUTO_INCREMENT(1,1), Text varchar(255), CreateTime DATETIME, PRIMARY KEY (ID))")
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
		file, err := os.Create("tododata.db") // Create SQLite file
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
	router.HandleFunc("todo/{id}, {text}", updateTodo).Methods("PUT")
	
	fmt.Println("Starting New Server")
	log.Fatal(http.ListenAndServe(":3000", router))
}