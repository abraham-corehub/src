package main

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qor/admin"
)

// User is GORM backend model
type User struct {
	gorm.Model
	Role     int
	Name     string
	Username string
	Password string
}

// Image is GORM backend model
type Image struct {
	gorm.Model
	Name  string
	User  User
	Album Album
}

// Album is GORM backend model
type Album struct {
	gorm.Model
	Name string
	User string
}

const dirDB = `db/`
const fNDB = `dbpbqor.db`

func main() {
	DB, _ := gorm.Open(`sqlite3`, fNDB)
	DB.AutoMigrate(&User{}, &Image{}, &Album{})

	Admin := admin.New(&admin.AdminConfig{DB: DB})

	Admin.AddResource(&User{})
	Admin.AddResource(&Image{})
	Admin.AddResource(&Album{})

	mux := http.NewServeMux()

	Admin.MountTo("/admin", mux)

	fmt.Println("Listening on: http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
