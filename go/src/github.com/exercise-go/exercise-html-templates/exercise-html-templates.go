package main

import (
	"html/template"
	"net/http"
)

//Todo is a custom type
type Todo struct {
	Title string
	Done  bool
}

//TodoPageData is a custom type
type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

var tmpl *template.Template

const tmplDir = "tmpl/01"

func main() {
	tmpl = template.Must(template.ParseFiles(tmplDir + "/" + "layout.html"))
	http.Handle("/static/", //final url can be anything
		http.StripPrefix("/static/",
			http.FileServer(http.Dir(tmplDir+"/"+"static"))))
	http.HandleFunc("/a", handlerPageA)
	http.HandleFunc("/b", handlerPageB)
	http.ListenAndServe(":8080", nil)
}

func handlerPageA(w http.ResponseWriter, r *http.Request) {
	data := TodoPageData{
		PageTitle: "My TODO list A",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
	tmpl.Execute(w, data)
}

func handlerPageB(w http.ResponseWriter, r *http.Request) {
	data := TodoPageData{
		PageTitle: "My TODO list B",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
			{Title: "Task 4", Done: true},
		},
	}
	tmpl.Execute(w, data)
}
