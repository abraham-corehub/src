package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	testDb()
}

func testDb() {
	db, err := sql.Open("sqlite3", "../photobook/db/pb.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select username, password from user")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var un string
		var pw string
		err = rows.Scan(&un, &pw)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(un, pw)
	}
}
