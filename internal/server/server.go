package server

import (
	"dalennod/internal/archive"
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
)

const PORT string = "41415"

var (
	data *sql.DB
	tmpl *template.Template

	deleteID = regexp.MustCompile("^/delete/([0-9]+)$")
)

type InsertEntry struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	Reason   string `json:"reason"`
	Keywords string `json:"keywords"`
	Group    string `json:"group"`
	Archive  int    `json:"archive"`
}

func Start(database *sql.DB) {
	data = database
	var mux *http.ServeMux = http.NewServeMux()

	var staticFileServer http.Handler = http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

	mux.HandleFunc("/", root)
	mux.HandleFunc("/delete/", deleteHandler)
	mux.HandleFunc("/add/", addHandler)

	err := http.ListenAndServe(":"+PORT, mux)
	if err != nil {
		logger.Error.Println(err)
	}
	logger.Info.Println("Started on http://localhost:" + PORT)
}

func root(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl = template.Must(template.ParseFiles("index.html"))
		var bookmarks []db.Bookmark = db.ViewAll(data, "s")
		tmpl.Execute(w, bookmarks)
	} else {
		internalServerErrorHandler(w, r)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var match []string = deleteID.FindStringSubmatch(r.URL.Path)
		if len(match) < 2 {
			internalServerErrorHandler(w, r)
			return
		}
		matchInt, err := strconv.Atoi(match[1])
		if err != nil {
			logger.Error.Println(err)
		}

		db.Remove(data, matchInt)
		w.WriteHeader(http.StatusOK)
	} else {
		internalServerErrorHandler(w, r)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var insData InsertEntry
		err := json.NewDecoder(r.Body).Decode(&insData)
		if err != nil {
			logger.Error.Println(err)
		}

		if insData.Archive == 0 {
			db.Add(data, insData.URL, insData.Title, insData.Reason, insData.Keywords, insData.Group, insData.Archive)
		} else {
			db.Add(data, insData.URL, insData.Title, insData.Reason, insData.Keywords, insData.Group, insData.Archive)
			go archive.SendSnapshot(insData.URL)
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		internalServerErrorHandler(w, r)
	}
}

func internalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

// func notFoundHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusNotFound)
// 	w.Write([]byte("404 Not Found"))
// }
