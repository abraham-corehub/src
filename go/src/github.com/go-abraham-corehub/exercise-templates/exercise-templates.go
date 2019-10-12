package main

import (
	"log"
	"net/http"
	"strings"
	"text/template"
)

const dirTmpl = "templates/"
const dirStatic = "static/"

var uITemplates *template.Template

// TemplateData type
type TemplateData struct {
	Title string
}

func main() {
	startWebApp()
}

func parseTemplates() {

	nUITs := []string{
		"layout",
		"head",
		"body",
	}

	for i := 0; i < len(nUITs); i++ {
		nUITs[i] = dirTmpl + nUITs[i] + ".html"
	}

	uITemplates = template.Must(template.ParseFiles(nUITs...))
}

func startWebApp() {
	parseTemplates()
	mux := http.NewServeMux()
	fileServer := http.FileServer(neuteredFileSystem{http.Dir(dirStatic)})
	mux.Handle("/"+dirStatic, http.StripPrefix("/"+dirStatic, fileServer))

	mux.HandleFunc("/", handlerHome)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handlerHome(w http.ResponseWriter, r *http.Request) {
	state := "layout"
	renderTemplate(w, state, &TemplateData{})
}

//To disable Directory Listing
//https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
type neuteredFileSystem struct {
	fs http.FileSystem
}

//To disable Directory Listing
//https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := nfs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, tD *TemplateData) {
	err := uITemplates.ExecuteTemplate(w, tmpl, tD)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
