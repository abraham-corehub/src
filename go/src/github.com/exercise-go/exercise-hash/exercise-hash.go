package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	getHash("abey")
}
func getHash(str string) []byte {
	pW := str
	pWH := sha1.New()
	pWH.Write([]byte(pW))

	pWHS := hex.EncodeToString(pWH.Sum(nil))

	fmt.Println(pW, pWHS)
	return []byte(pWHS)
}
