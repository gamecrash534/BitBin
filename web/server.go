package web

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

//go:embed page.gohtml
var page embed.FS

var link string
var linkTitle string

var tmpl *template.Template

type Page struct {
	Link      string
	LinkTitle string
}

type Response struct {
	Key string `json:"key"`
}

func RunServer(port uint, l string, lt string) {
	link = l
	linkTitle = lt
	tmpl = template.Must(template.ParseFS(page, "*.gohtml"))

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/{path}", pathHandler)

	log.Fatal(http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), nil))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := tmpl.Execute(w, Page{Link: link, LinkTitle: linkTitle})
		if err != nil {
			printErr(w, err, "template execution")
		}
	} else if r.Method == http.MethodPost {
		w.WriteHeader(201)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Response{Key: "returnedasdf"})
		if err != nil {
			printErr(w, err, "POST request handling")
		}
	}
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	w.WriteHeader(200)
	w.Write([]byte(path))
}

func printErr(w http.ResponseWriter, err error, at string) {
	w.WriteHeader(500)
	w.Write([]byte("There was an internal error :("))
	log.Println("Error occurred at "+at+": ", err)
}
