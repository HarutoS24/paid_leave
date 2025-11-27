package main

import (
	"backend/handler"
	"backend/middleware"
	"backend/service"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

var store *sessions.FilesystemStore

func init() {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		panic("secret keyが設定されていません")
	}
	store = sessions.NewFilesystemStore("/session_data", []byte(secretKey))
}

func main() {
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

	http.Handle("/auth/login",
		middleware.RequireAuthSession(store, middleware.UseDB(db,
			http.HandlerFunc(handler.LoginHandler),
		)),
	)
	http.Handle("/auth/info",
		middleware.RequireLogin(store, db, middleware.UseDB(db,
			http.HandlerFunc(handler.GetLoggedInEmployeeHandler),
		)),
	)
	http.Handle("/employee/add",
		middleware.RequireAdminLogin(store, db, middleware.UseDB(db,
			http.HandlerFunc(handler.AddEmployeeHandler),
		)),
	)

	http.Handle("/employee/add_demo", middleware.UseDB(db,
		http.HandlerFunc(handler.AddEmployeeHandler),
	))

	http.Handle("/vacation/add", middleware.UseDB(db, middleware.RequireLogin(store, db,
		http.HandlerFunc(handler.AddPaidLeaveHandler),
	)))
	fmt.Println("starting...")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Printf("ListenAndServe error: %v\n", err)
	}
}
