package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	hashTest("admin")
	hashTest("abey")
	getTimeStamp(time.Now())
}

func hashTest(str string) []byte {
	pW := str
	pWH := sha1.New()
	pWH.Write([]byte(pW))

	pWHS := hex.EncodeToString(pWH.Sum(nil))

	fmt.Println(pW, pWHS)
	return []byte(pWHS)
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

func getTimeStamp(t time.Time) string {
	tD := t.Format("20060102")
	tT := strings.Replace(t.Format("15.04.05.000"), ".", "", 3)
	dTS := tD + tT
	fmt.Println(dTS)
	return dTS
}
