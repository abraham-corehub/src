package main

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qor/admin"
)

// User Create a GORM-backend model
type User struct {
	gorm.Model
	Name string
}

// Product Create another GORM-backend model
type Product struct {
	gorm.Model
	Name        string
	Description string
}

func main() {
	DB, _ := gorm.Open("sqlite3", "demo.db")
	DB.AutoMigrate(&User{}, &Product{})

	// Initalize
	Admin := admin.New(&admin.AdminConfig{DB: DB})

	// Allow to use Admin to manage User, Product
	Admin.AddResource(&User{})
	Admin.AddResource(&Product{})

	// initalize an HTTP request multiplexer
	mux := http.NewServeMux()

	// Mount admin interface to mux
	Admin.MountTo("/admin", mux)

	fmt.Println("Listening on: http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
