package server

import (
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"database/sql"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
)

const PORT string = "41415"

var (
	data *sql.DB
	tmpl *template.Template

	deleteID = regexp.MustCompile("/delete/([0-9]+)$")
)

func Start(database *sql.DB) {
	data = database
	var mux *http.ServeMux = http.NewServeMux()

	var assetsFileServer http.Handler = http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assetsFileServer))
	var jsFileServer http.Handler = http.FileServer(http.Dir("js"))
	mux.Handle("/js/", http.StripPrefix("/js/", jsFileServer))

	mux.HandleFunc("/", root)
	mux.HandleFunc("/delete/", deleteHandler)

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
		// http.ServeFile(w, r, "delete.html")
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
