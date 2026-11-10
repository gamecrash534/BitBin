package web

import (
	"bitbin/storage"
	"embed"
	"encoding/json"
	"html/template"
	"io"
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

	mux := http.NewServeMux()

	mux.HandleFunc("/", viewHandler)
	mux.HandleFunc("/favicon.svg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/favicon.svg")
	})
	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed on this path", http.StatusMethodNotAllowed)
			return
		}
		postHandler(w, r)
	})
	mux.HandleFunc("/{path}", pathHandler)

	log.Fatal(http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), corsHandler(mux)))
}

func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := tmpl.Execute(w, Page{Link: link, LinkTitle: linkTitle})
		if err != nil {
			printErr(w, err, "template execution")
			return
		}
	} else if r.Method == http.MethodPost {
		postHandler(w, r)
	}
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("path")

	c, data, err := storage.Get(key)
	if err != nil {
		printErr(w, err, "file content retrieval")
		return
	}

	if data == nil || c == nil {
		w.WriteHeader(404)
		w.Write([]byte("Given path could not be found"))
		return
	}

	w.Header().Set("Content-Type", data.MimeType)
	w.WriteHeader(200)
	w.Write(c)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		printErr(w, err, "POST body retrieval")
		return
	}

	ctype := r.Header.Get("Content-Type")
	if ctype == "" {
		ctype = "applicatiom/octet-stream"
	}

	key, _, err := storage.Create(data, ctype, r.Header.Get("Content-Encoding") == "gzip")
	if err != nil {
		printErr(w, err, "file storing")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	err = json.NewEncoder(w).Encode(Response{Key: *key})
	if err != nil {
		printErr(w, err, "POST response")
		return
	}
}

func printErr(w http.ResponseWriter, err error, at string) {
	w.WriteHeader(500)
	w.Write([]byte("There was an internal error :("))
	log.Println("Error occurred at "+at+": ", err)
}
