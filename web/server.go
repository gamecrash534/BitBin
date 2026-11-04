package web

import (
	"bitbin/storage"
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

//go:embed page.gohtml
var page embed.FS

const keyCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var link string
var linkTitle string

var tmpl *template.Template
var rnd *rand.Rand

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

	src := rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixMicro()))
	rnd = rand.New(src)

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/{path}", pathHandler)

	log.Fatal(http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), nil))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := tmpl.Execute(w, Page{Link: link, LinkTitle: linkTitle})
		if err != nil {
			printErr(w, err, "template execution")
			return
		}
	} else if r.Method == http.MethodPost {
		log.Println("Received post")

		data, err := io.ReadAll(r.Body)
		if err != nil {
			printErr(w, err, "POST body retrieval")
			return
		}

		var key string
		for {
			key = randomStr()
			exists, err := storage.Exists(key)
			if err != nil {
				printErr(w, err, "key gen")
				return
			}

			if !exists {
				break
			}
		}

		err = storage.StoreFile(key, data, r.Header.Get("Content-Encoding") == "gzip")
		if err != nil {
			printErr(w, err, "file storing")
			return
		}

		w.WriteHeader(201)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(Response{Key: key})
		if err != nil {
			printErr(w, err, "POST response")
			return
		}
	}
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("path")

	data, err := storage.GetContents(key)
	if err != nil {
		printErr(w, err, "file content retrieval")
		return
	}

	if data == nil {
		w.WriteHeader(404)
		w.Write([]byte("Given path could not be found"))
		return
	}

	w.WriteHeader(200)
	w.Write(data)
}

func printErr(w http.ResponseWriter, err error, at string) {
	w.WriteHeader(500)
	w.Write([]byte("There was an internal error :("))
	log.Println("Error occurred at "+at+": ", err)
}

func randomStr() string {
	res := make([]byte, 8)
	for i := range res {
		res[i] = keyCharset[rnd.IntN(len(keyCharset))]
	}

	return string(res)
}
