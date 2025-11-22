package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("starting...")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Printf("ListenAndServe error: %v\n", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	db_name := os.Getenv("")
	dsn := fmt.Sprintf("%s:%s@tcp(db:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, db_name)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var name string

	err = db.QueryRow("SELECT employee_name FROM kintai_db.employees_tbl").Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintln(w, "No result found")
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	fmt.Fprintf(w, "%s", name)
}
