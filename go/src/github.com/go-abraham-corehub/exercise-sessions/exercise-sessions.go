package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var uuid string

func main() {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	uuid = strings.Trim(string(out), "\n")
	//fmt.Printf("%s", out)
	http.HandleFunc("/get", getCookie)
	http.HandleFunc("/set", setCookie)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getCookie(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value
	fmt.Println(sessionToken)
}

func setCookie(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	uuid = strings.Trim(string(out), "\n")
	sessionToken := uuid
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
	fmt.Println("Cookie Set")
}
