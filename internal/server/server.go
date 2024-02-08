package server

import (
	"dalennod/internal/archive"
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"dalennod/internal/setup"
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
		var bookmarks []setup.Bookmark = db.ViewAll(data, true)
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
		var (
			insData       setup.Bookmark
			archiveResult bool   = false
			snapshotURL   string = ""
		)
		err := json.NewDecoder(r.Body).Decode(&insData)
		if err != nil {
			logger.Error.Println(err)
		}

		if !insData.Archived {
			db.Add(data, insData.URL, insData.Title, insData.Note, insData.Keywords, insData.BGroup, insData.Archived, snapshotURL)
		} else {
			archiveResult, snapshotURL = archive.SendSnapshot(insData.URL)
			if archiveResult {
				db.Add(data, insData.URL, insData.Title, insData.Note, insData.Keywords, insData.BGroup, insData.Archived, snapshotURL)
			} else {
				logger.Warn.Println("Snapshot failed.")
				db.Add(data, insData.URL, insData.Title, insData.Note, insData.Keywords, insData.BGroup, false, snapshotURL)
			}
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
